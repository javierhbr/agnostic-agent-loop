# Component Spec: Payments Service (Recurring Charges)

**Component ID:** SBL-PAYMENTS-001
**Implements:** Feature SBL-2025-Q1 (Subscription Billing → Payment Processing)
**Context Pack:** cp-v2.1
**Status:** Draft (blocked by ADR-001)
**Blocked By:** ADR-001 (card vault decision — determines entire implementation)

---

## Overview

The Payments Service handles card tokenization and recurring charge processing. In Phase 1, using Stripe Vault (per recommended ADR-001), we integrate with Stripe's API for secure token storage and charge processing.

**Source:** Feature-spec + ADR-001 (Stripe Vault recommended)

---

## Domain Boundaries

**Owns:** Card tokenization, charge processing, Stripe integration, charge state transitions

**Does NOT own:** Subscription records (Billing Service), customer identity (User Service)

**Contracts:**
- **Consumes:** `subscription.renew_requested` from Billing Service
- **Emits:** `charge.completed`, `charge.failed` to event bus
- **Uses:** Stripe API (external system)

---

## Data Model

### Charge Record

```typescript
charge {
  id: string                    // unique, e.g. "charge_abc123"
  subscription_id: string       // which subscription triggered this
  amount_cents: number
  currency: "USD"
  stripe_charge_id: string      // Stripe's charge ID
  stripe_token_id: string       // Token from Stripe vault
  status: SUCCESS | FAILED
  failure_reason: string (nullable)  // "insufficient_funds", "card_expired", etc.
  attempt_number: number        // 1-3 (retries)
  created_at: timestamp
}
```

### Stripe Integration (if self-hosted in Phase 2)

```typescript
// NOT used in Phase 1 (Stripe vault)
// Placeholder for Phase 2 self-hosted consideration

stripe_customer {
  id: string                    // Stripe Customer ID (maps to our customer)
  shopflow_customer_id: string
  stripe_id: string
}

stripe_token {
  // NOT stored in Phase 1 (Stripe keeps it)
}
```

---

## API Contracts

**Source:** Feature-spec + ADR-001 decision

### Stripe Integration APIs (Internal to Payments Service)

#### Tokenize Card
```
POST /api/v1/tokens/create
{
  card_number: "4242424242424242",
  expiry_month: "12",
  expiry_year: "2026",
  cvc: "123"
}
→ 200 {
  stripe_token_id: "tok_visa_from_stripe",
  card_last_4: "4242",
  created_at: "2025-02-25T10:00:00Z"
}
// Note: Card data in request is sent to Stripe, not stored at ShopFlow
```

#### Process Recurring Charge
```
POST /api/v1/charges/process-recurring
{
  subscription_id: "sub_xyz789",
  stripe_token_id: "tok_visa_from_stripe",
  amount_cents: 999,
  idempotency_key: "sub_xyz789_2025-03-25"  // for retry idempotency
}
→ 200 {
  id: "charge_abc123",
  stripe_charge_id: "ch_stripe_id_from_stripe",
  status: "SUCCESS",
  amount_cents: 999,
  created_at: "2025-02-25T10:00:00Z"
}
// Emits event: charge.completed

OR

→ 400 {
  id: "charge_abc123",
  status: "FAILED",
  failure_reason: "card_declined",
  created_at: "2025-02-25T10:00:00Z"
}
// Emits event: charge.failed
```

---

## Implementation Responsibilities

**Source:** Feature-spec + ADR-001 (Stripe Vault)

1. **Card Tokenization**
   - Accept card details from customer at checkout
   - Send to Stripe API to tokenize (card never touches our servers)
   - Return Stripe token to client, then to Billing Service
   - Log tokenization with metadata (last 4, expiry, merchant_id) but NOT card data

2. **Recurring Charge Processing**
   - Consume `subscription.renew_requested` event from Billing Service
   - Look up Stripe token from Billing Service (subscription.stripe_token_id)
   - Call Stripe API to charge the token
   - Record charge success/failure with reason
   - Emit `charge.completed` or `charge.failed` event

3. **Idempotency**
   - Use idempotency keys (subscription_id + charge_date) to ensure retry safety
   - If same charge retried, Stripe returns same result (no double-charge)

4. **Error Handling**
   - Categorize Stripe errors: card_declined, insufficient_funds, expired_card, gateway_error
   - Non-retryable (card_declined, expired_card): set subscription status to SUSPENDED
   - Retryable (gateway_error, timeout): schedule retry for next cycle

5. **PII/PCI Handling**
   - Never log or store card numbers
   - Never store plain Stripe tokens (use references only)
   - Log charge metadata: subscription_id, amount, card_last_4, result
   - Compliance: card data handled per Stripe's PCI standards

---

## Acceptance Criteria

### AC1: Cards Can Be Tokenized
**Given** a customer provides card details,
**When** POST to `/tokens/create`,
**Then** Stripe returns a token and we store only the token (not the card).

### AC2: Recurring Charges Process Successfully
**Given** a Stripe token exists,
**When** Payments Service receives `subscription.renew_requested`,
**Then** charge is processed, Stripe API is called, and `charge.completed` or `charge.failed` is emitted.

### AC3: Idempotency Works
**Given** the same charge is requested twice (network retry),
**When** both requests use same idempotency key,
**Then** Stripe returns the same result (no double-charge).

### AC4: Errors Are Categorized
**Given** a charge fails,
**When** Stripe returns error code,
**Then** error is categorized (retryable vs. non-retryable) and Billing Service is notified.

---

## Observability

**Source:** Feature-spec NFR + Discovery.md observability

### Logging
```
event: "charge.processed"
subscription_id: "sub_xyz789"
charge_id: "charge_abc123"
stripe_charge_id: "ch_stripe_id"
amount_cents: 999
card_last_4: "4242"
status: "SUCCESS" | "FAILED"
failure_reason: "card_declined" (if failed)
attempt_number: 1
timestamp: "2025-02-25T10:00:00Z"
latency_ms: 250
```

### Metrics
- Counter: `charges.attempted` (per merchant, per result)
- Counter: `charges.succeeded` / `charges.failed` (by failure_reason)
- Histogram: `charge_latency_ms` (p50, p95, p99)

### Tracing
- Span: `charge.process` (call Stripe API)
- Attributes: subscription_id, amount, card_last_4

### Alerts
- Page on-call if charge success rate < 95%
- Page on-call if charge latency p95 > 500ms (Stripe is slow)

---

## Dependencies

**Source:** Feature-spec integration

- **Stripe API:** For tokenization and charging
  - Endpoint: https://api.stripe.com/v1/
  - Auth: API key (stored in secrets)
  - Timeout: 30 seconds (Stripe can be slow)

- **Billing Service:** Consumes `subscription.renew_requested` event
- **Event Bus:** Publishes `charge.completed`, `charge.failed`

---

## Deployment

**Source:** Feature-spec feature flag

- Feature flag: `subscriptions_enabled` (controls whether charges are triggered)
- Stripe account: Separate test + production credentials
- Rate limiting: 100 charges/sec per merchant (configurable)

---

## Phase 2 Optionality (Notes for Future)

If Phase 2 evaluates self-hosted vault:
- This spec would change significantly
- Stripe integration becomes optional (fallback only)
- Card rotation, expiry, encryption logic added
- PCI-DSS audit scope increases

---

## Gates Check

- ✅ Gate 1: Context Completeness (all sections sourced)
- ✅ Gate 2: Domain Validity (owns charges, events for cross-domain)
- ✅ Gate 3: Integration Safety (charge events clear, consumers listed)
- ✅ Gate 4: NFR Compliance (logging, metrics, tracing, alerts specified)
- ⏳ Gate 5: Ready to Implement (blocked by ADR-001 resolution)

Once ADR-001 is approved (Stripe chosen), this spec is unblocked.

---

## Handoff to Developer

**Inputs:**
- This component-spec-payments.md
- ADR-001 resolution
- Feature-spec (acceptance criteria, NFRs)

**Outputs Expected:**
- `impl-spec-payments.md` (Stripe integration code structure, edge cases)
- `tasks.yaml` (decomposed tasks)
- `verify.md` (integration tests with Stripe sandbox)
