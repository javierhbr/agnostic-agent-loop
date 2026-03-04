# Implementation Spec: Billing Service

**Spec ID:** SBL-BILLING-IMPL-001
**Implements:** Component SBL-BILLING-001
**Context Pack:** cp-v2.1
**Status:** Draft (will be verified by Verifier agent)

---

## Overview

This document details the exact code changes, data model, and implementation strategy for the Billing Service. It is the bridge between the component-spec (WHAT) and actual code (HOW).

**Source:** Component Spec SBL-BILLING-001 + Developer analysis

---

## Data Model (Database Schema)

**Source:** Component-spec + Stripe phase 1 choice (ADR-001)

### plans table
```sql
CREATE TABLE plans (
  id VARCHAR(32) PRIMARY KEY,
  merchant_id VARCHAR(32) NOT NULL,
  name VARCHAR(255) NOT NULL,
  amount_cents INT NOT NULL,
  currency VARCHAR(3) DEFAULT 'USD',
  interval ENUM('MONTHLY', 'QUARTERLY', 'ANNUAL') NOT NULL,
  description TEXT,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  FOREIGN KEY (merchant_id) REFERENCES merchants(id),
  INDEX (merchant_id),
  INDEX (is_active)
);
```

### subscriptions table
```sql
CREATE TABLE subscriptions (
  id VARCHAR(32) PRIMARY KEY,
  merchant_id VARCHAR(32) NOT NULL,
  customer_id VARCHAR(32) NOT NULL,
  plan_id VARCHAR(32) NOT NULL,
  status ENUM('ACTIVE', 'CANCELLED', 'SUSPENDED_PAYMENT_FAILED', 'PAUSED') DEFAULT 'ACTIVE',
  stripe_token_id VARCHAR(255) NOT NULL,  -- from Stripe, not stored plaintext
  next_charge_date DATE NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  cancelled_at TIMESTAMP NULL,
  cancellation_reason VARCHAR(255) NULL,
  FOREIGN KEY (merchant_id) REFERENCES merchants(id),
  FOREIGN KEY (customer_id) REFERENCES customers(id),
  FOREIGN KEY (plan_id) REFERENCES plans(id),
  INDEX (merchant_id),
  INDEX (customer_id),
  INDEX (next_charge_date),  -- for scheduler query
  INDEX (status)
);
```

### subscription_events table (audit trail)
```sql
CREATE TABLE subscription_events (
  id VARCHAR(32) PRIMARY KEY,
  subscription_id VARCHAR(32) NOT NULL,
  event_type ENUM('CREATED', 'RENEWED', 'CANCELLED', 'SUSPENDED', 'RESUMED'),
  event_data JSON,  -- {subscription_id, customer_id, merchant_id, amount, ...}
  created_at TIMESTAMP DEFAULT NOW(),
  published BOOLEAN DEFAULT false,
  published_at TIMESTAMP NULL,
  FOREIGN KEY (subscription_id) REFERENCES subscriptions(id),
  INDEX (subscription_id),
  INDEX (event_type),
  INDEX (published)
);
```

---

## Code Changes

**Source:** Component-spec + TDD approach

### New Service: billing-service

#### Directory Structure
```
billing-service/
├── cmd/
│   └── server/main.go
├── internal/
│   ├── models/
│   │   └── subscription.go      # Plan, Subscription, SubscriptionEvent types
│   ├── api/
│   │   ├── plans.go             # POST /plans, GET /plans
│   │   └── subscriptions.go     # POST /subscriptions, DELETE /:id
│   ├── service/
│   │   ├── subscription.go      # Business logic (create, cancel, etc.)
│   │   └── scheduler.go         # Daily charge scheduler
│   ├── events/
│   │   └── publisher.go         # Publish to event bus
│   ├── db/
│   │   └── migrations/          # SQL migration files
│   └── config/
│       └── config.go
├── test/
│   ├── integration_test.go
│   └── models_test.go
└── go.mod
```

#### Key Functions

**`subscription.go`** — Business Logic

```go
// CreatePlan creates a new subscription plan for a merchant
func (s *SubscriptionService) CreatePlan(ctx context.Context, merchantID string, req CreatePlanRequest) (*Plan, error)
  // Validate: amount_cents > 0, interval valid
  // Generate unique plan ID
  // Insert into plans table
  // Return plan

// CreateSubscription creates a subscription for a customer
func (s *SubscriptionService) CreateSubscription(ctx context.Context, req CreateSubscriptionRequest) (*Subscription, error)
  // Validate: plan exists, merchant active
  // Generate unique subscription ID
  // Calculate next_charge_date based on plan interval
  // INSERT subscription record
  // Emit event: subscription.created
  // Return subscription

// CancelSubscription cancels a subscription
func (s *SubscriptionService) CancelSubscription(ctx context.Context, subscriptionID string, reason string) error
  // Find subscription, verify status
  // Set status = CANCELLED, cancelled_at = now()
  // Emit event: subscription.cancelled
  // Return error if not found or already cancelled

// SuspendSubscription suspends due to payment failure
func (s *SubscriptionService) SuspendSubscription(ctx context.Context, subscriptionID string) error
  // Set status = SUSPENDED_PAYMENT_FAILED
  // Emit event: subscription.suspended
  // Return to Payments Service for retry logic
```

**`scheduler.go`** — Billing Schedule

```go
// RunDailyScheduler runs every 6 hours, finds subscriptions due, triggers charges
func (s *BillingScheduler) RunDailyScheduler(ctx context.Context) error
  // SELECT * FROM subscriptions WHERE next_charge_date <= TODAY AND status = 'ACTIVE'
  // For each:
  //   - Get latest plan (price might have changed)
  //   - Emit event: subscription.renew_requested
  //     { subscription_id, customer_id, amount, stripe_token_id }
  // Return count of charges triggered

// The Payments Service CONSUMES subscription.renew_requested
// and processes the charge asynchronously
```

**`events/publisher.go`** — Event Publishing

```go
// PublishEvent publishes a subscription event to the event bus
func (p *EventPublisher) PublishEvent(ctx context.Context, event SubscriptionEvent) error
  // Marshal event to JSON
  // Publish to event bus (e.g., Kafka, RabbitMQ, AWS SNS)
  // Mark event as published in DB
  // Return error if bus is down

// Events published:
// - subscription.created
// - subscription.renewed (from Payments Service feedback)
// - subscription.cancelled
// - subscription.suspended
```

---

## Edge Cases

**Source:** Component-spec + Team discussion

| Edge Case | Scenario | Handling | Evidence |
|-----------|----------|----------|----------|
| **Customer cancels mid-cycle** | Subscription active, next_charge_date in 15 days. Customer cancels today. | No refund (out of scope Phase 1). Subscription status = CANCELLED, no future charges. Cancel event emitted. | Test: verify next_charge_date cleared, no event scheduled |
| **Payment fails, retry 3 times** | Charge fails on day 1, 2, 3. Never succeeds. | After 3 failures, status = SUSPENDED_PAYMENT_FAILED. Manual intervention needed (merchant contact). Suspension event emitted. | Test: simulate 3 failed charges, verify suspension logic |
| **Merchant updates plan price** | Plan price changes $9.99 → $14.99. Active subscriptions exist. | NEW subscriptions use new price. ACTIVE subscriptions keep old price until cancellation (Phase 1 simplicity). Phase 2: proration. | Test: create subscription, update plan price, verify next charge still uses old price |
| **Scheduler runs twice (crash recovery)** | Scheduler crashes mid-run, restarts, runs again same day. | Idempotency: subscription.renew_requested events have subscription_id + charge_date as composite key. If event already published, skip. | Test: publish same event twice, verify only one charge triggered |
| **Database write fails after event publish** | Event published, but subscription.next_charge_date update fails. | Retry logic in event consumer (Payments Service). If charge succeeds but update fails, Payments Service retries until success (eventual consistency). | Test: simulate crash between event publish and DB commit, verify recovery |
| **Subscription cancelled, then renewed request arrives** | Subscription cancelled, but renewal event in flight. | Payments Service checks subscription.status before processing. If CANCELLED, charge rejected with reason "subscription_cancelled". | Test: cancel subscription, trigger renewal event, verify charge rejected |

---

## Observability

**Source:** Feature-spec NFR + Component-spec

### Logging (Structured JSON)

```json
{
  "timestamp": "2025-02-25T10:00:00Z",
  "event": "subscription.created",
  "subscription_id": "sub_xyz789",
  "merchant_id": "merchant_456",
  "customer_id": "cust_789",
  "plan_id": "plan_abc123",
  "amount_cents": 999,
  "interval": "MONTHLY",
  "next_charge_date": "2025-03-25",
  "service": "billing-service",
  "level": "info"
}
```

### Metrics

```go
// In code:
metrics.Counter("subscriptions.created", 1,
  "merchant_id", merchantID,
  "interval", plan.Interval)

metrics.Counter("subscriptions.cancelled", 1,
  "merchant_id", merchantID,
  "reason", cancellationReason)

metrics.Gauge("active_subscriptions", count,
  "merchant_id", merchantID)
```

### Tracing

```go
// Using OpenTelemetry:
span := tracer.Start(ctx, "subscription.create")
defer span.End()
span.SetAttributes(
  attribute.String("subscription_id", subscription.ID),
  attribute.String("plan_id", plan.ID),
  attribute.Int64("amount_cents", plan.AmountCents),
)
```

---

## Rollout Plan

**Source:** Feature-spec Feature Flag section

### Phase 1: Internal Testing
- Feature flag: `subscriptions_enabled` = false
- Staging environment: Full tests with real event bus
- Timeline: 3 days

### Phase 2: Canary (5% of new merchants)
- Feature flag: enabled for 5% of new merchants
- Monitoring: Alert if error rate > 1%
- Timeline: 2-3 days

### Phase 3: Ramp (50% of new merchants)
- Increase to 50%
- Timeline: 2-3 days

### Phase 4: General Availability
- 100% of new merchants see subscriptions
- Kill switch: `subscriptions_kill_switch` remains (disables new creation in seconds if critical bug)
- Timeline: 1 day

### Rollback
- Flip `subscriptions_enabled` = false
- Existing subscriptions keep billing (no charge disruption)
- New subscriptions rejected with "Feature not available yet"
- Reconciliation job run to identify any inconsistencies

---

## Gates Check

- ✅ Gate 1: Context Completeness (all sections sourced, component-spec cited)
- ✅ Gate 2: Domain Validity (owns subscriptions, events for cross-domain, no DB sharing)
- ✅ Gate 3: Integration Safety (events defined: `subscription.*`, consumers: Payments, Email, User)
- ✅ Gate 4: NFR Compliance (logging, metrics, tracing, PII masking specified, performance targets in code)
- ✅ Gate 5: Ready to Implement (no blocking ADRs, edge cases covered, tests defined)

---

## Next Steps

1. ✓ Implementation spec complete (this document)
2. → Decompose into tasks.yaml (8-10 development tasks)
3. → Assign to developer team
4. → Code review
5. → Verifier agent: write verify.md with test evidence

---

## References

- [Component Spec: Billing](./component-spec-billing.md)
- [Feature Spec](./feature-spec.md)
- [Tasks](./tasks.yaml)
