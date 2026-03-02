# Component Spec: Billing Service

**Component ID:** SBL-BILLING-001
**Implements:** Feature SBL-2025-Q1 (Subscription Billing)
**Context Pack:** cp-v2.1
**Status:** Draft (blocked by ADR-001)
**Blocked By:** ADR-001 (card vault decision)

---

## Overview

The Billing Service is a new microservice responsible for the subscription lifecycle: creating subscriptions, maintaining billing schedules, triggering recurring charges, and managing cancellations. It is the source of truth for subscription state.

**Source:** Feature-spec Problem Statement + Architecture review

---

## Domain Boundaries

**Owns:** Subscription records, billing schedules, subscription lifecycle state

**Does NOT own:** Card storage (Payments Service), customer data (User Service), email sending (Email Service)

**Contracts:**
- **Consumes:** `subscription.renew_requested` (from Billing scheduler)
- **Emits:** `subscription.created`, `subscription.renewed`, `subscription.cancelled`, `subscription.suspended`
- **Uses:** Payments Service API for charges (see component-spec-payments)

---

## Data Model

### Subscription Record

```typescript
subscription {
  id: string                    // unique, e.g. "sub_2xj9k2..."
  merchant_id: string           // which merchant owns this
  customer_id: string           // which customer is subscribed
  plan_id: string               // which plan (monthly $9.99, etc.)
  status: ACTIVE | CANCELLED | SUSPENDED_PAYMENT_FAILED | PAUSED
  stripe_token_id: string       // token from Stripe vault (ADR-001)
  next_charge_date: date        // when to charge next
  created_at: timestamp
  updated_at: timestamp
  cancelled_at: timestamp (nullable)
  cancellation_reason: string (nullable)
}
```

### Plan Record (Merchant-facing)

```typescript
plan {
  id: string
  merchant_id: string
  name: string                  // "Monthly Pro"
  amount_cents: number          // 999 = $9.99
  currency: "USD"
  interval: "MONTHLY" | "QUARTERLY" | "ANNUAL"
  description: string
  is_active: boolean
  created_at: timestamp
}
```

---

## API Contracts

**Source:** Feature-spec Acceptance Criteria + Platform MCP REST standards

### Merchant APIs

#### Create Plan
```
POST /api/v1/subscriptions/plans
{
  name: "Monthly Pro",
  amount_cents: 999,
  interval: "MONTHLY",
  description: "Pro plan with advanced features"
}
→ 201 {
  id: "plan_abc123",
  merchant_id: "merchant_456",
  created_at: "2025-02-25T10:00:00Z"
}
```

#### List Plans
```
GET /api/v1/subscriptions/plans
→ 200 {
  plans: [
    { id, name, amount_cents, interval, is_active, created_at },
    ...
  ]
}
```

#### Create Subscription
```
POST /api/v1/subscriptions
{
  plan_id: "plan_abc123",
  customer_id: "cust_789",
  stripe_token_id: "tok_visa_from_stripe"  // from Payments Service
}
→ 201 {
  id: "sub_xyz789",
  status: "ACTIVE",
  next_charge_date: "2025-03-25",
  created_at: "2025-02-25T10:00:00Z"
}
// Emits event: subscription.created
```

#### Cancel Subscription
```
DELETE /api/v1/subscriptions/:id
→ 200 {
  id: "sub_xyz789",
  status: "CANCELLED",
  cancelled_at: "2025-02-25T10:30:00Z"
}
// Emits event: subscription.cancelled
```

---

## Implementation Responsibilities

**Source:** Feature-spec feature flag + gate requirements

1. **Subscription Lifecycle Management**
   - Create, read, update, cancel operations
   - State machine: ACTIVE → CANCELLED | SUSPENDED (on payment failure)
   - Idempotent creation (same request twice = same subscription, not duplicate)

2. **Billing Schedule**
   - Daily scheduler job (runs every 6 hours) to find subscriptions due today
   - For each due subscription, emit `subscription.renew_requested` event
   - Retry logic: if charge fails, retry up to 3 times over next 2 days

3. **Event Emission**
   - Publish to event bus after every state change
   - Event schema: `{ subscription_id, customer_id, merchant_id, event_type, timestamp }`
   - At-least-once delivery (exactly-once preferred)

4. **PII Handling**
   - Never store card numbers or customer emails in subscription records
   - Log only: subscription_id, customer_id, card_last_4 (from Payments Service), amount
   - Mask email in logs

---

## Acceptance Criteria

### AC1: Merchants Can Create Plans
**Given** a merchant POSTs to `/subscriptions/plans` with valid data,
**When** the request is processed,
**Then** a plan is created, assigned a unique ID, and returned with 201.

### AC2: Subscriptions Can Be Created
**Given** a valid plan and Stripe token exist,
**When** a POST to `/subscriptions` includes plan_id and stripe_token_id,
**Then** a subscription record is created with status ACTIVE, next_charge_date set to interval, and event emitted.

### AC3: Scheduler Triggers Charges
**Given** a subscription exists with next_charge_date <= today,
**When** the daily scheduler runs,
**Then** a `subscription.renew_requested` event is emitted to Payments Service.

### AC4: Cancellation Works
**Given** a subscription exists,
**When** a DELETE request is made,
**Then** status becomes CANCELLED and event is emitted.

---

## Observability

**Source:** Feature-spec NFR section + Discovery.md observability needs

### Logging
```
event: "subscription.created"
subscription_id: "sub_xyz789"
merchant_id: "merchant_456"
customer_id: "cust_789"
amount_cents: 999
interval: "MONTHLY"
next_charge_date: "2025-03-25"
timestamp: "2025-02-25T10:00:00Z"
```

### Metrics
- Counter: `subscriptions.created` (per merchant, per plan)
- Counter: `subscriptions.cancelled` (reason: user_request, payment_failed, etc.)
- Gauge: `active_subscriptions_count` (per merchant)

### Tracing
- Span: `subscription.create` (parent: API request)
- Span: `scheduler.renew` (triggered daily)

---

## Deployment

**Source:** Feature-spec Feature Flag section

- Feature flag: `subscriptions_enabled` (default: false during dev)
- Rollout: Phased (5% → 50% → 100% over 1 week)
- Kill switch: `subscriptions_kill_switch` (disables new subscriptions in seconds)

---

## Gates Check

- ✅ Gate 1: Context Completeness (all sections have Source:)
- ✅ Gate 2: Domain Validity (owns only subscriptions, events for cross-domain)
- ✅ Gate 3: Integration Safety (events identified, consumers listed)
- ✅ Gate 4: NFR Compliance (logging, metrics, tracing specified)
- ⏳ Gate 5: Ready to Implement (blocked by ADR-001 — card vault not yet decided)

---

## Handoff to Developer

Once ADR-001 is approved, this spec is unblocked and can be handed to the Developer agent for `impl-spec-billing.md` and `tasks.yaml` creation.

**Inputs for Developer:**
- This component-spec-billing.md
- ADR-001 resolution (Stripe chosen)
- Feature-spec constraints (acceptance criteria, NFRs)

**Outputs Expected:**
- `impl-spec-billing.md` (data model detail, code structure, edge cases)
- `tasks.yaml` (decomposed into development tasks)
- `verify.md` (test plan and verification)
