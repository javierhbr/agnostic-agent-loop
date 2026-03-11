# SDD v3.0 Example Workflow: Guest Checkout

This document shows a **real walkthrough** of how to use the new SDD v3.0 methodology from start to finish.

---

## Phase 0: Initiative Definition (Product Manager)

A product manager wants to reduce cart abandonment by enabling guest checkout.

### Step 1: Create the Initiative

```bash
# Classify the risk level using the risk-assessment skill
# User answers 5 questions one at a time...
#
# Q1: Does this change any cross-service contract? YES (CartUpdated event schema)
# Q2: Does this touch payment/auth/PII? NO (just cart)
# Q3: Will multiple teams work on this? YES (3-4 services)
# Q4: Do we have a rollback plan? YES (feature flag)
# Q5: ADR dependencies? NO (clear)
#
# Result: MEDIUM risk → STANDARD workflow

agentic-agent specifyify start "Enable Guest Checkout" --risk medium
```

**Output:**
```
✓ Initiative created: .agentic/sdd/initiatives/enable-guest-checkout.yaml

Workflow: STANDARD (Medium Risk)
├─ Architect → Design feature spec
├─ Developers (parallel, per component) → Implement
└─ Verifier → Test and verify

Expected timeline: 3-5 days
Current agent: Architect (ready to receive specs)
```

### Step 2: PM Hands Off to Architect

The initiative file contains:
- Problem statement: "30% cart abandonment from guests"
- Success metrics: "< 20% abandonment within 30 days"
- Affected components: [checkout-service, payment-service, user-service, email-service]
- Risk context: medium, standard workflow

---

## Phase 1: Architecture Design (Solution Architect)

The architect reads the initiative and uses the **architect/SKILL.md** to understand what to produce.

### Step 1: Architect Reads the Initiative

```bash
# Architect reads the initiative file
cat .agentic/sdd/initiatives/enable-guest-checkout.yaml

# Then reads the Platform Constitution to understand constraints
cat openclaw-specs/constitution/policies.md

# Key policies this spec must follow:
# - PII handling: guest email must be encrypted at rest
# - Observability: all auth flows must log (no PII), emit metrics
# - Performance: payment processing < 2 seconds (p95)
# - Security: session TTL max 24 hours
```

### Step 2: Architect Designs the Solution

Creates `feature-spec.md`:
- What: "Enable guest checkout flow"
- Why: "Reduce cart abandonment 30% → 20%"
- How: "No account required, guest email collected, payment processed, confirmation email sent"
- UX flows: Step 1 (guest flag) → Step 2 (email collection) → Step 3 (payment) → Step 4 (confirmation)
- Key invariants: "Guest checkout must not create user account; payment idempotent"
- Security: "Guest email encrypted at rest per Constitution/Security"
- Observability: "Log auth attempt, emit `guest_checkout_started` metric, trace payment flow"

Creates `component-spec.md` for each service:
- **checkout-service**: Add guest checkout flow, skip account creation step
- **payment-service**: No changes (existing API works)
- **user-service**: Add guest customer tracking (no account = no login)
- **email-service**: Send confirmation to guest email
- **analytics-service**: Track checkout_guest metric

Each component spec includes:
- Acceptance criteria (GWT format)
- NFRs (performance, security, observability)
- Contract changes (if any)
- Dependencies on other specs

### Step 3: Gate Check Architecture

```bash
# Architect validates their design passes all 5 gates before handing to developers
agentic-agent specifyify gate-check SPEC-CHECKOUT-GUEST

✓ Gate 1: Context Completeness PASS
  - feature-spec.md has implements, context_pack, status fields
  - component-spec.md files present for all 4 services

✓ Gate 2: Domain Validity PASS
  - No cross-domain DB access
  - User domain invariant respected: "Email unique, Active XOR Deleted"

✓ Gate 3: Integration Safety PASS
  - Contract change declared: CartUpdated event adds guest_flag
  - Consumers identified: checkout-service, analytics-service
  - Dual-publish plan documented in spec

✓ Gate 4: NFR Compliance PASS
  - Security: Guest email encryption per Constitution §2
  - Observability: logging, metrics, tracing all declared
  - Performance: payment < 2s baseline respected

✓ Gate 5: Ready to Implement PASS
  - blocked_by is empty (no pending ADRs)
  - All 8 acceptance criteria in GWT format
  - No ambiguity in component responsibilities
```

### Step 4: Create Task Fan-Out

Architect produces a task fan-out (automatically from spec):

```bash
# Each component team gets a task with impl-spec.md template
agentic-agent task create --from feature-spec.md --fanout-by-component

# Creates 4 parallel tasks:
# - TASK-001: Implement checkout-service guest flow
# - TASK-002: Implement user-service guest tracking
# - TASK-003: Implement email-service guest confirmation
# - TASK-004: Implement analytics-service guest metrics
```

---

## Phase 2: Implementation (Developers, Parallel)

Each component team receives their `component-spec.md` and uses the **developer/SKILL.md** to understand what to produce.

### Developer: Checkout Service

**Input:** `component-spec-checkout-service.md`

**What they produce:**
- `impl-spec-checkout-service.md` with:
  - Exact code changes (add GuestCheckoutFlow class)
  - All 6 acceptance criteria verified
  - Edge cases (invalid email, network timeout, duplicate submission)
  - Observability: logging points, metrics emitted, tracing spans
  - Rollback: feature flag check before enabling guest flow

```go
// Example code snippet in impl-spec
if !featureFlags.IsEnabled("GuestCheckoutEnabled") {
  return ErrGuestCheckoutDisabled
}

guestEmail := req.Email
// Log attempt without PII
log.Info("guest_checkout_attempt", "request_id", requestID)

// Emit metric
metrics.Increment("guest_checkout_total")
```

- `tasks.yaml` with implementation steps:
  1. Add GuestCheckout flag to Feature Flags system
  2. Implement GuestCheckoutFlow class
  3. Add guest email validation
  4. Integrate with payment service
  5. Add feature flag guard rails
  6. Write unit tests (min 4 edge cases)
  7. Write integration tests
  8. Performance test: payment < 2s latency

**Gate Check:**
```bash
agentic-agent specifyify gate-check SPEC-CHECKOUT-GUEST-IMPL

✓ Gate 4: NFR Compliance PASS
  - Logging: "guest_checkout_attempt" emitted with request_id
  - Metrics: "guest_checkout_total" counter, "payment_latency_ms" histogram
  - Tracing: Span "GuestCheckout" with tags guest_email_hash, result
  - Security: Email not logged raw, hashed in observability

✓ Gate 5: Ready to Implement PASS
  - All 6 ACs in GWT format with test evidence
  - No blocking ADRs
  - Feature flag default: false (safe)
```

### Developer: User Service

**Input:** `component-spec-user-service.md`

**What they produce:**
- `impl-spec-user-service.md`:
  - Add GuestCustomer table (minimal: email, created_at, last_order_at)
  - Do NOT create User account (guest remains anonymous)
  - Maintain invariant: "Email unique, Active XOR Deleted" (guests don't violate this)

### Developer: Email Service

**Input:** `component-spec-email-service.md`

**What they produce:**
- `impl-spec-email-service.md`:
  - Add guest confirmation email template
  - Send to email address from guest_checkout event
  - Track delivery in metrics

### Developer: Analytics Service

**Input:** `component-spec-analytics-service.md`

**What they produce:**
- `impl-spec-analytics-service.md`:
  - Consume CartUpdated v2 event (has guest_flag)
  - Emit "checkout_guest_flow" dimension for reports
  - Calculate abandonment rate by flow type

---

## Phase 3: Verification (Verifier)

After all developers commit their code, the verifier uses the **verifier/SKILL.md** to prove it works.

### Step 1: Verify Every Acceptance Criterion

```bash
# Run all 6 acceptance criteria with observable evidence

✓ AC1: Given guest user, When they reach checkout, Then they can proceed without account
   Evidence: Unit test GuestCheckoutFlow_NoAccountRequired passes
   Trace: curl -X POST /api/v1/checkout/guest -d '{"email":"test@example.com"}' → 200 OK

✓ AC2: Given guest provides email, When they enter payment, Then payment is processed
   Evidence: Integration test GuestCheckout_PaymentFlow passes
   Trace: Payment service logs "payment_processed" for guest email

✓ AC3: Given payment completes, When guest receives confirmation, Then email sent < 30s
   Evidence: E2E test confirms delivery in 20ms average
   Trace: Email service metrics "email_delivery_latency_ms" p95 < 500

✓ AC4-6: Similar evidence for other ACs...
```

### Step 2: Update Spec Graph

```bash
# Mark initiative as Done in spec-graph.json
agentic-agent specifyify sync-graph

# This creates/updates .agentic/spec-graph.json:
{
  "SPEC-CHECKOUT-GUEST": {
    "implements": "INIT-GUEST-CHECKOUT",
    "status": "Done",
    "depends_on": [],
    "affects": ["SPEC-ANALYTICS-GUEST"],
    "contracts_referenced": ["CartUpdated-v2"],
    "blocked_by": [],
    "updated_at": "2026-02-28T15:00:00Z"
  }
}
```

### Step 3: Create Verify.md

```markdown
# Verification Report: Guest Checkout

## All Acceptance Criteria Verified ✓

### AC1: Guest Checkout Without Account
- **Test:** GuestCheckoutFlow_NoAccountRequired
- **Result:** PASS
- **Evidence:** curl output shows HTTP 200, no account created
- **Latency:** 150ms p95

### AC2: Guest Email → Payment Processing
- **Test:** GuestCheckout_PaymentFlow_Integration
- **Result:** PASS
- **Evidence:** Payment service logs guest transaction, metrics confirm
- **Latency:** 850ms p95 (within budget)

### AC3: Confirmation Email < 30s
- **Test:** GuestCheckout_EmailConfirmation_E2E
- **Result:** PASS
- **Evidence:** Email received in 20ms average
- **Latency:** 500ms p99

### AC4-6: Additional criteria...

## Gate Verification

- Gate 1 (Context): ✓ All metadata complete
- Gate 2 (Domain): ✓ Invariants maintained
- Gate 3 (Integration): ✓ Consumers notified and updated
- Gate 4 (NFR): ✓ Observability deployed and working
- Gate 5 (Ready): ✓ All blocking ADRs resolved

## Deployment

- Feature Flag: GuestCheckoutEnabled = OFF (safe default)
- Rollout Plan: Monitor metrics for 24h before enabling in production
- Rollback: Disable feature flag (instant rollback, no deploy required)
```

---

## Phase 4: Deployment & Success Metrics

Once verifier confirms all ACs and gates:

```bash
# 1. Merge to main
git merge feature/guest-checkout

# 2. Deploy with feature flag OFF (safe)
agentic-agent deploy

# 3. Monitor metrics
agentic-agent metrics watch guest_checkout_total

# 4. Enable feature flag at 10% traffic
agentic-agent flags set GuestCheckoutEnabled=10pct

# 5. Monitor for 24 hours
# - guest_checkout_total steadily increasing
# - cart_abandonment_rate trending down
# - payment_latency_ms stays < 2s
# - no error spikes

# 6. If metrics look good, enable to 100%
agentic-agent flags set GuestCheckoutEnabled=100pct

# 7. Measure final success at day 30
# Expected: cart_abandonment 30% → 18% (below 20% target)
```

---

## Key Artifacts Created

| Phase | Owner | Artifact | Location |
|-------|-------|----------|----------|
| **0** | PM | Initiative | `.agentic/sdd/initiatives/enable-guest-checkout.yaml` |
| **1** | Architect | Feature Spec | `openclaw-specs/features/guest-checkout/feature-spec.md` |
| **1** | Architect | Component Specs (x4) | `openclaw-specs/features/guest-checkout/component-spec-*.md` |
| **2** | Developers | Impl Specs (x4) | `openclaw-specs/features/guest-checkout/impl-spec-*.md` |
| **2** | Developers | Tasks | `.agentic/tasks/in-progress.yaml` |
| **3** | Verifier | Verify Report | `openclaw-specs/features/guest-checkout/verify.md` |
| **3** | Verifier | Spec Graph | `.agentic/spec-graph.json` |

---

## CLI Cheat Sheet for This Feature

```bash
# PM: Define and classify risk
agentic-agent specifyify start "Enable Guest Checkout" --risk medium

# PM: View progress
agentic-agent specifyify workflow show enable-guest-checkout

# Architect: Gate check design
agentic-agent specifyify gate-check SPEC-CHECKOUT-GUEST

# Developers: Claim component task
agentic-agent task claim TASK-001
agentic-agent task complete TASK-001

# Verifier: Run all gates before merge
agentic-agent specifyify gate-check SPEC-CHECKOUT-GUEST
agentic-agent validate

# All: Sync traceability to platform repo
agentic-agent specifyify sync-graph

# Deploy when verified
git merge feature/guest-checkout
agentic-agent deploy --with-feature-flag GuestCheckoutEnabled=off
```

---

## Success: 30 Days Later

```
METRIC                   BEFORE    TARGET    ACTUAL    STATUS
Cart Abandonment Rate    30%       < 20%     18%       ✓ EXCEEDS
Payment Latency (p95)    3.2s      < 2s      1.8s      ✓ EXCEEDS
Guest Checkout Adoption  0%        > 10%     12.5%     ✓ EXCEEDS
Logged-in Conversion     5.2%      >= 5.2%   5.3%      ✓ MAINTAINED
```

**Result:** Initiative marked as Done. Architecture, implementation, and verification all documented in spec-graph. Team ships guest checkout with confidence.
