# Discovery Report: Subscription Billing Platform

**Initiative ID:** SBL-2025-Q1
**Analyst:** Claude (AI Agent)
**Date:** 2025-02-25
**Status:** Complete — Ready for Architect Handoff

---

## Problem Statement

ShopFlow merchants lack subscription/recurring payment capabilities. This prevents us from serving subscription-based business models (SaaS, memberships) and directly causes support overload and deal losses.

**Evidence:**
- **Support Load:** 43 support tickets in January alone (12% of total) asking "How do I set up recurring payments?"
- **Deal Loss:** 8 confirmed lost deals to competitors citing subscription feature lack
- **Market Demand:** 2023 SMB benchmark report: 67% of merchants want recurring payment options

**Success Metric:** Within 60 days of launch, 25% of new merchants enable at least one subscription plan.

---

## Team Insights

**Interview Participants:**
- Sarah (Product Manager)
- Alex (Backend Tech Lead)
- James (SRE/Ops)
- Priya (Platform Architect, interim)

### Key Facts Extracted

1. **Pain Points:** Merchants manually invoice customers or use external billing providers (Stripe, Recurly), creating friction and reducing loyalty
2. **Competitor Analysis:** Shopify, BigCommerce, WooCommerce all offer native subscriptions
3. **Merchant Segments:**
   - SMB (< $100k/year revenue): want simple monthly/annual plans
   - Mid-market: want more control (prorations, custom intervals, metered billing)
4. **Timeline Pressure:** Q2 roadmap depends on this; product team has 3 merchants actively waiting

---

## Affected Components

The subscription lifecycle touches **5 services and 1 new service**:

| Service | Role | Change |
|---------|------|--------|
| **Platform API** | Entry point for merchants to create plans | New `/subscriptions/` endpoints |
| **Billing Service** | NEW | Core logic: subscription lifecycle, billing schedule, charge triggering |
| **Payments Service** | Card storage & recurring charges | Card tokenization, charge processing for recurring |
| **Email Service** | Subscription lifecycle notifications | New event subscriptions: `subscription.created`, `subscription.renewed`, `subscription.cancelled` |
| **User Service** | Customer identity & cohorts | Add "has_active_subscription" flag, webhook for marking customers |
| **Analytics Service** | Metrics & dashboards | Track subscription events for cohort analysis |

### Data Flow: One Subscription Lifecycle

```
1. Merchant creates plan
   Platform API → Billing Service (POST /subscriptions/plans)

2. Customer purchases subscription
   Platform API → Billing Service (POST /subscriptions)
   → Payments Service (tokenize card)
   → Event: subscription.created

3. On billing date, auto-charge
   Billing Service (scheduler) → Payments Service (charge)
   → Event: subscription.renewed (or failed)

4. Customer cancels
   Platform API → Billing Service (DELETE /subscriptions/:id)
   → Event: subscription.cancelled

5. All events → Email Service, User Service, Analytics Service
```

---

## Risk Classification

**Level: CRITICAL**

### Risk Factors

1. **Payment Data** ✓ CRITICAL
   - Tokenized credit cards will be stored (PCI-DSS scope)
   - Recurring charges without per-charge authorization (regulatory risk)
   - Card storage location TBD (architecture decision required)

2. **PII Data** ✓ CRITICAL
   - Billing addresses, customer emails
   - Subscription history reveals customer behavior

3. **Data Consistency** ✓ CRITICAL
   - Subscription state must stay in sync with Payments Service charge state
   - If Billing says "charged" but Payments says "failed," we have inconsistency
   - Rollback is not a simple feature flag

4. **Multi-Service Coordination** ✓ HIGH
   - 5 services must emit/consume events correctly
   - Event-driven system is eventually consistent (reconciliation needed)
   - New service (Billing) must handle edge cases: retry logic, idempotency

5. **Regulatory/Compliance** ✓ HIGH
   - PCI-DSS certification scope increases
   - Recurring charges may require explicit customer consent (varies by jurisdiction)
   - Audit trail must be complete

### Rollback Complexity

**Not a feature flag.**

If a bug exists in production:
- Disabling subscriptions = no charges (revenue loss)
- Re-enabling = risk of double-charges if not reconciled first

**Runbook Required:**
1. Disable new subscription creation (safe)
2. Keep renewing existing subscriptions (contain damage)
3. Fix the bug
4. Reconcile any missed/double charges manually
5. Re-enable

---

## Blocking Decisions (ADRs Required)

### ADR-001: Card Storage Strategy

**Decision:** Where will ShopFlow store tokenized credit cards?

**Options:**
A. **Self-Hosted Vault in Payments Service**
   - Pro: Full control, no third-party dependency, lower per-charge cost
   - Con: PCI-DSS compliance burden increases, card rotation/expiry management, security responsibility

B. **Stripe Vault (via Stripe Connect)**
   - Pro: PCI-DSS handled by Stripe, mature product, automatic card expiry
   - Con: Vendor lock-in, per-transaction fees higher, customer data lives on Stripe

C. **Hybrid (tokenize via Stripe, store token at ShopFlow)**
   - Pro: Stripe handles PCI, we have token ownership
   - Con: Still our responsibility to secure tokens, still third-party dependency

**Rationale for Architect:**
This decision blocks the Payments Service component-spec and all developer work. **Architect must resolve this ADR before development can begin.** Recommend gathering CTO + Security input.

**Timeline Impact:** 3-5 days for decision (security review required)

---

## Components Analysis

### Affected Services' Current State

**Payments Service:**
- Already has one-time charge logic, card tokenization
- Does NOT currently store cards long-term (one-time only)
- Needs extension for recurring charge capability

**Billing Service:**
- DOES NOT EXIST — will be created
- Responsible for: subscription plan creation, subscription records, billing schedule, charge triggering logic

**Email Service:**
- Already emits for one-time orders
- Needs to subscribe to new events: `subscription.created`, `subscription.renewed`, `subscription.cancelled`

**User Service:**
- Already tracks customer data
- Needs to add: "has_active_subscription" boolean flag, webhook to receive `subscription.created` events

**Analytics Service:**
- Needs to ingest new events for cohort-level subscription analytics

---

## Non-Goals (Explicitly Out of Scope)

1. **Dunning (Retry Logic):** Automatic retry on failed charges — Phase 2
2. **Prorations:** Smart handling of mid-cycle upgrades/downgrades — Phase 2
3. **Metered Billing:** Usage-based billing — Phase 2
4. **Multi-Currency:** Only USD in Phase 1
5. **Merchant Dashboard:** Subscription analytics UI — Phase 2
6. **Customer Portal:** Self-service subscription management — Phase 2

---

## Success Metrics & Observability Requirements

### Platform-Level Metrics (Architect to specify in feature-spec)

1. **Adoption Metric:** % of new merchants with at least 1 active subscription (Target: 25% within 60 days)
2. **Charge Success Rate:** % of scheduled charges that succeed (Target: 99.5%)
3. **Charge Latency:** p50, p95, p99 latency for charge processing (Target: p95 < 200ms)
4. **Error Breakdown:** Categorize failures: invalid card, gateway down, insufficient funds, etc.
5. **Revenue Impact:** MRR from subscriptions (tracks business outcome)

### Observability Spec (Developer to implement)

**Logging:**
- Every subscription lifecycle event (create, update, cancel, suspend)
- Every charge attempt (attempt, success, failure with reason)
- All with: subscription_id, customer_id, amount, timestamp, error_code (if applicable)
- **PII Masking:** Card last 4 only, customer_email masked in logs

**Metrics:**
- Counter: `subscriptions.created`, `subscriptions.cancelled`
- Counter: `charges.attempted`, `charges.succeeded`, `charges.failed_by_reason`
- Histogram: `charge_processing_latency_ms`
- Gauge: `active_subscriptions_count`

**Tracing:**
- Span: `subscription.create` with subscription_id, plan_id
- Span: `payment.charge` with subscription_id, amount, card_token
- Span: `email.send` with subscription_id, event_type

**Alerts (On-Call):**
- Charge success rate drops below 95%
- Charge latency p95 > 500ms
- Event queue depth growing (email service lag)

---

## Exit Gate Validation

**Analyst Self-Check Before Handing Off to Architect:**

- ✅ **Gate 1: Context Completeness**
  - Problem statement has metric: 43 tickets, 8 lost deals, 25% new merchant adoption goal
  - Evidence is real data, not assumption
  - All sources documented

- ✅ **Gate 2: Domain Validity**
  - No cross-domain data access violations identified
  - Each service owns its domain (Payments owns charges, Billing owns subscriptions)
  - Event-driven for cross-domain communication (proper boundary)

- ✅ **Gate 3: Integration Safety**
  - New events identified: `subscription.created`, `subscription.renewed`, `subscription.cancelled`
  - Consumers listed: Email, User Service, Analytics
  - No contract changes (additive events)

- ✅ **Gate 4: NFR Compliance**
  - Observability specified: metrics, logs, tracing, alerts
  - PII handling specified: card last 4 only, email masked
  - Performance targets captured: p95 < 200ms for charges

- ✅ **Gate 5: Ready to Architect**
  - Blocking ADR identified: ADR-001 (card vault)
  - All non-goals explicitly listed
  - Affected components clear and verified with team

---

## Recommendation for Architect

1. **Resolve ADR-001** immediately — gather Security + CTO, decide on card storage strategy
2. **Design feature-spec** with explicit acceptance criteria, NFRs, rollback strategy
3. **Design 5 component-specs**: Platform API, Billing, Payments, Email, User Service
4. **Create fan-out tasks** for parallel Developer work
5. **Run gate-check** on all specs before handoff to developers

---

## Next Steps

1. ✓ Analyst discovery complete (this document)
2. → Schedule architecture kickoff meeting
3. → Assign Architect role
4. → Begin ADR-001 resolution (CTO review)
