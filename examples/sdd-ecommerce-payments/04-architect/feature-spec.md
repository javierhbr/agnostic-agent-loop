# Feature Spec: Subscription Billing Platform

**Feature ID:** SBL-2025-Q1
**Implements:** Initiative SBL-2025-Q1
**Context Pack:** cp-v2.1 (Payment + PII, Critical Risk)
**Status:** Draft → Approved (after gate validation)
**Blocked By:** ADR-001 (card vault decision)

---

## Metadata

**Source:** Platform MCP v2.1 — Payment Initiation & PCI Scope

- **Initiative:** SBL-2025-Q1 (see `01-initiative/initiative.md`)
- **Risk Level:** Critical (payment + PII + data consistency)
- **Workflow:** Full (Analyst → Architect → Developer × 5 → Verifier)
- **Affected Services:** Billing (new), Payments, Email, User Service, Platform API
- **Timeline:** 3-4 weeks
- **Deployment:** Feature-flagged rollout (see Feature Flag Strategy below)

---

## Problem Statement

**Source:** Discovery.md + Platform MCP Market Context

ShopFlow merchants cannot offer recurring payments, limiting ShopFlow's addressable market to one-time purchase merchants. This causes:
- **Support Overload:** 43 support tickets/month (12% of volume) asking "How do I do subscriptions?"
- **Deal Loss:** 8 confirmed lost deals to competitors (Shopify, BigCommerce, WooCommerce) that have native subscriptions
- **Market Gap:** 67% of SMB merchants want subscription options (2023 benchmark)

**Success Metric:** Within 60 days of launch, 25% of new merchants enable at least one subscription plan.

---

## Goals & Non-Goals

**Source:** Initiative SBL-2025-Q1

### Goals

1. Enable merchants to create, view, and manage subscription plans (monthly, quarterly, annual)
2. Process recurring automatic charges on schedule with tokenized cards
3. Handle common subscription lifecycle events: create, upgrade, downgrade, cancel, suspend
4. Reduce subscription-related support tickets from 43/month to < 8/month (< 2% of load)
5. Enable 25% of new merchants to launch subscriptions within 60 days

### Non-Goals (Explicitly Out of Scope — Phase 2)

- Dunning (automatic retry on failed charges) — too complex for Phase 1
- Prorations (smart handling of mid-cycle downgrades) — implementation complexity
- Metered billing (usage-based charges) — requires rate ingestion system
- Multi-currency support — Phase 2, after core model is stable
- Merchant dashboard for subscription analytics — merchant portal feature, Phase 2
- Customer self-service portal — customer UX, Phase 2

---

## User Experience

**Source:** Discovery.md team interviews + Platform MCP UX guidelines

### Merchant Perspective: Creating a Subscription Plan

```
Merchant Flow:
1. Products → Create New Product → Subscription Checkbox (not a one-time sale)
2. Configure Plan:
   - Plan name (e.g., "Monthly Pro")
   - Price ($9.99)
   - Billing interval (Monthly, Quarterly, Annual)
   - First charge date (immediate, or custom date)
3. Review and Publish
4. See plan in products list with subscription badge
5. When customers purchase, subscription record created, billing scheduled automatically
```

### Customer Perspective: Starting a Subscription

```
Customer Flow:
1. Browse products
2. See "Pro Plan - $9.99/month" (subscription badge visible)
3. Add to cart
4. Checkout (same as one-time, but see "auto-renews monthly" notice)
5. Enter card (tokenized in Payments Service)
6. Confirmation: "Your subscription is active. Next charge: 2025-03-25"
7. Auto-charge on 2025-03-25 (silent, unless it fails)
```

### Observability: Seeing It Work

```
Operations Team:
- Dashboard metric: "Active subscriptions: 1,247"
- Alert: "Charge success rate dropped to 94% (threshold: 95%)"
- Log search: "subscription_id:sub-12345 event:charge_failed reason:insufficient_funds"
- Tracing: Drill into a charge: customer ID, amount, card last 4, result
```

---

## Domain Responsibilities

**Source:** Discovery.md component analysis + Platform MCP Domain Model

Each service owns one domain; cross-domain communication is event-driven.

| Service | Domain | Responsibility | Boundary |
|---------|--------|---|---|
| **Billing Service** | Subscription Lifecycle | • Create subscriptions • Schedule charges • Record renewals • Handle cancellations • Maintain subscription state | Event: emits `subscription.*` events |
| **Payments Service** | Card Processing & Storage | • Tokenize cards (PCI scope) • Store tokens (if self-hosted) or issue tokens (if Stripe) • Process charges (recurring) • Emit charge results | Event: consumes `subscription.renewed` requests, emits `charge.*` events |
| **Email Service** | Customer Notifications | • Subscribe to `subscription.created`, `subscription.renewed`, `subscription.cancelled` • Compose and send emails | Event: consumes `subscription.*` and `charge.*` events |
| **User Service** | Customer Identity | • Mark customers with `has_active_subscription` • Track subscription history (cohort analysis) • Send webhooks on subscription events | Event: consumes `subscription.*` events |
| **Platform API** | Entry Point | • Expose `/subscriptions/plans` (create, list, update) • Expose `/subscriptions` (list, create, view) • Expose `/subscriptions/:id/cancel` | REST API |

---

## Cross-Domain Interactions

**Source:** Platform MCP Integration Model + Discovery.md data flow

Event-driven architecture. No direct database access across domains.

```
Subscription Lifecycle (Data Flow):

1. SUBSCRIPTION CREATION
   Platform API (Merchant)
     → Billing Service: POST /subscriptions
       (merchant_id, plan_id, customer_id, card_token)
     → Payments Service: Store token (if self-hosted) or validate Stripe token
     → Billing Service: Create subscription record (status: ACTIVE, next_charge: T+30d)
     → Emit Event: subscription.created
       (subscription_id, customer_id, merchant_id, plan_id, next_charge_date)
     ↓ (consumed by)
     → Email Service: Sends "Subscription Confirmed" email
     → User Service: Marks customer has_active_subscription=true
     → Analytics: Ingests for cohort

2. SCHEDULED CHARGE (Billing Scheduler)
   Billing Service (cron every 6 hours, checks subscriptions due today)
     → Identify subscriptions with next_charge_date <= today
     → Emit Event: subscription.renew_requested
       (subscription_id, amount, customer_id)
     ↓ (consumed by)
     → Payments Service: Processes charge against stored token
     → Payments Service: Emit Event: charge.completed OR charge.failed
     ↓ (consumed by)
     → Billing Service: Updates subscription.next_charge_date or suspension_reason
     → Email Service: Sends receipt (success) or alert (failure)
     → User Service: Updates subscription_status in customer record

3. SUBSCRIPTION CANCELLATION
   Platform API (Merchant or Customer)
     → Billing Service: DELETE /subscriptions/:id
     → Billing Service: Set status: CANCELLED, next_charge_date: null
     → Emit Event: subscription.cancelled
       (subscription_id, customer_id, cancellation_reason)
     ↓ (consumed by)
     → Email Service: Sends "Subscription Cancelled" email
     → User Service: Sets has_active_subscription=false
     → Analytics: Cohort update
```

---

## Non-Functional Requirements

**Source:** Platform MCP v2.1 — Payment Processing Standards + Discovery.md observability needs

### Performance

- **Charge Processing Latency:** p50 < 50ms, p95 < 200ms, p99 < 500ms (credit card API is slow)
  - Source: Platform MCP — payment integrations
- **Subscription Creation:** p95 < 100ms (mostly database write)
  - Source: Platform MCP — API response time SLA
- **Batch Job (Daily Charge Schedule):** Complete all charges due within 1 hour
  - Source: Discovery.md — team expectation

### Reliability

- **Charge Success Rate:** Target 99.5% (failures due to invalid card, insufficient funds, gateway down)
  - Source: Platform MCP — payment SLA for production
- **Event Delivery:** At-least-once (exactly-once preferred, but at-least-once acceptable with idempotency keys)
  - Source: Platform MCP — event bus guarantees
- **Rollback Strategy:** See Feature Flag & Rollback section below

### Security & Compliance

- **PCI-DSS Scope:** All card data handling must comply with PCI-DSS Level 1
  - Source: Platform MCP v2.1 — payment security
- **Card Storage:** Tokenized only; card numbers NEVER stored plaintext
  - Decision: ADR-001 (self-hosted vault or Stripe vault)
- **PII Masking:** Logs contain card last 4 only, email masked with **** characters
  - Source: Platform MCP — PII handling policy
- **Audit Trail:** Every subscription event logged with timestamp, actor, result
  - Source: Platform MCP — compliance audit requirements
- **Data Retention:** Subscription records kept for 7 years (financial law)
  - Source: Platform MCP — data retention policy

### Observability

- **Logging:** Structured JSON logs for every lifecycle event, with traceable IDs
- **Metrics:** Counters for subscriptions created/cancelled, histograms for charge latency, gauges for active count
- **Tracing:** Distributed traces for charge processing across Billing → Payments → Email
- **Alerting:** Page on-call if charge success rate drops below 95%

---

## Acceptance Criteria

**Source:** Initiative SBL-2025-Q1 + Discovery.md + Platform MCP examples

All criteria are testable and observable.

### AC1: Merchants Can Create Subscription Plans

**Given** a merchant is logged in and on the Products page,
**When** they click "Create New Product" and enable "Subscription Plan",
**Then** they can specify plan name, price, and billing interval (monthly/quarterly/annual) and publish the plan.

**Evidence:** Unit test for `POST /subscriptions/plans` with valid payload returns 201. Integration test shows plan in products list with subscription badge.

---

### AC2: Customers Can Start Subscriptions

**Given** a subscription plan exists,
**When** a customer adds it to cart and completes checkout with a valid card,
**Then** a subscription record is created with status ACTIVE, next_charge_date set to subscription interval, and confirmation email sent within 30 seconds.

**Evidence:** Integration test: simulate customer purchase flow, verify subscription record in database, verify email sent via log.

---

### AC3: Automatic Charges Execute on Schedule

**Given** a subscription has a next_charge_date <= today,
**When** the Billing Service charge scheduler runs,
**Then** a charge is processed against the stored card, the transaction succeeds or fails, and Billing Service updates the subscription's next_charge_date for the next cycle (or sets suspension_reason if failed).

**Evidence:** Integration test: schedule a charge for today, run scheduler, verify charge processed and subscription state updated.

---

### AC4: Subscription Lifecycle Events Are Emitted

**Given** any subscription event (created, renewed, cancelled, suspended),
**When** the event occurs,
**Then** an event is emitted to the event bus with subscription_id, customer_id, timestamp, and event type.

**Evidence:** Unit test: verify events are published to event bus with correct schema. Integration test: verify Email Service receives and processes event.

---

### AC5: PII Is Masked in Logs

**Given** a subscription lifecycle event is logged,
**When** the event includes PII (customer email, card, address),
**Then** logs contain only non-sensitive identifiers (subscription_id, customer_id, card_last_4) and email is masked (**** characters).

**Evidence:** Log output inspection: search logs for a subscription event, verify no plaintext email or card found.

---

### AC6: Charge Failures Are Retryable and Observable

**Given** a charge attempt fails (invalid card, gateway down, insufficient funds),
**When** the failure is recorded,
**Then** the subscription status reflects the failure, an alert is emitted, and the Billing Service scheduler will retry on the next cycle (up to 3 times).

**Evidence:** Integration test: simulate failed charge, verify subscription status is "SUSPENDED_PAYMENT_FAILED", verify alert emitted, verify retry on next cycle.

---

### AC7: Rollback Strategy Is Tested

**Given** the system is running subscriptions in production,
**When** a critical bug is discovered,
**Then** the rollback procedure can disable new subscriptions without affecting existing ones, and a reconciliation report identifies any missed/double charges.

**Evidence:** Runbook test: execute disable procedure, verify new subscriptions rejected, existing subscriptions continue, verify reconciliation job produces accurate report.

---

## Feature Flag & Rollback Strategy

**Source:** Platform MCP Deployment Practices + Discovery.md rollback concerns

### Feature Flag

```yaml
flag_name: subscriptions_enabled
default: false (disabled during development)
rollout_strategy: phased
  phase_1: 5% of new merchants (canary, 2 days)
  phase_2: 50% of new merchants (ramp, 3 days)
  phase_3: 100% of new merchants (general availability)
kill_switch: subscriptions_kill_switch (disables new subscription creation in seconds)
```

### Rollback Procedure

**If a critical bug is discovered:**

1. **Immediate (0-30 min):** Flip `subscriptions_kill_switch` → new subscription creation disabled
   - Existing subscriptions continue (live with bug, contain damage)
   - No new customers affected

2. **Investigation (30 min - 2 hours):** Understand the bug
   - E.g., "Charges are succeeding but subscription state not updating"

3. **Fix & Deploy (1-2 hours):** Deploy hotfix to Billing or Payments Service

4. **Reconciliation (2-4 hours):** Run reconciliation job
   - Compare Billing Service subscription records vs. Payments Service charge records
   - Identify missed charges (should have happened but didn't) or double-charges
   - Generate report for Finance team

5. **Re-enable (post-reconciliation):** Flip `subscriptions_kill_switch` back on

**Acceptance:** Rollback tested in staging before production release.

---

## Blocking Decisions

**Source:** Discovery.md + CTO review

### ADR-001: Card Storage Strategy (BLOCKS IMPLEMENTATION)

**Decision Title:** Where will ShopFlow store tokenized credit cards?

**Status:** PENDING (blocks Payments component-spec and all developer work)

**Options:**
1. **Self-Hosted Vault in Payments Service**
   - Pro: Full control, no vendor lock-in, lower per-charge costs
   - Con: PCI-DSS burden (security + audit), card rotation logic, our liability

2. **Stripe Vault (via Stripe Connect)**
   - Pro: PCI-DSS handled by Stripe, proven product, automatic card expiry
   - Con: Vendor lock-in, higher transaction fees, customer data on Stripe

3. **Hybrid: Tokenize via Stripe, Store Token at ShopFlow**
   - Pro: Stripe handles PCI, we have token ownership
   - Con: Still our responsibility to secure tokens, still third-party dependency

**Recommendation:** Option 2 (Stripe Vault) for Phase 1
- Justification: Stripe reduces time-to-market, PCI compliance risk is lower, feature parity achieved
- Phase 2 can evaluate self-hosted if cost becomes concern

**Decision Timeline:** CTO + Security review by EOD (3 days)

**Impact:** Entire Payments component-spec depends on this. Cannot proceed to developer without decision.

---

## Gates Validation

**Self-Check Performed by Architect:**

- **Gate 1 — Context Completeness:** ✅ PASS
  - Every section has Source: line
  - Context Pack: cp-v2.1 declared in Metadata
  - All MCP sources cited

- **Gate 2 — Domain Validity:** ✅ PASS
  - No cross-domain DB access (all event-driven)
  - Each service owns its domain (Billing owns subscriptions, Payments owns charges)
  - Event boundaries clear

- **Gate 3 — Integration Safety:** ✅ PASS
  - Events identified: `subscription.created`, `subscription.renewed`, `subscription.cancelled`, `charge.completed`, `charge.failed`
  - Consumers listed: Email, User Service, Analytics
  - No contract breaking (additive events)

- **Gate 4 — NFR Compliance:** ✅ PASS
  - Logging specified: structured JSON, PII masked, audit trail
  - Metrics specified: success rate, latency, counts
  - Tracing specified: distributed traces across services
  - Alerting specified: on-call threshold

- **Gate 5 — Ready to Implement:** ⏳ BLOCKED
  - **Blocker:** ADR-001 (card vault decision) not yet resolved
  - Cannot hand off to Developer until ADR decision made

---

## Next Steps

1. ✓ Feature-spec complete (this document)
2. ✓ Component-specs complete (see separate component-spec files)
3. → **Resolve ADR-001** (CTO + Security meeting, 3 days)
4. → Re-run Gate 5 check once ADR resolved
5. → Create developer fan-out tasks (one per component: Billing, Payments, Email, User)
6. → Hand off to Developer agents for parallel implementation

---

## References

- [Initiative](../../01-initiative/initiative.md)
- [Discovery](../../03-analyst/discovery.md)
- [ADR-001](./ADR-001-card-vault.md)
- [Component Spec: Billing](./component-spec-billing.md)
- [Component Spec: Payments](./component-spec-payments.md)
- [Component Spec: Email](./component-spec-email.md)
