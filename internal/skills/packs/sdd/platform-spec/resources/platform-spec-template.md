# Platform Spec Template

Full template for creating a Platform Spec with all sections, worked examples, and acceptance criteria patterns. Referenced by the main platform-spec skill.

---

## Platform Spec File Structure

File location: Platform Repo (not component repo)
File name: `SPEC-<DOMAIN>-<NUMBER>-<short-title>.md` (e.g., `SPEC-CHECKOUT-001-guest-checkout.md`)

```markdown
# Platform Spec: SPEC-<DOMAIN>-<NUMBER>

## Metadata
- ID: SPEC-<DOMAIN>-<NUMBER>
- Initiative: <ECO-XXX or equivalent>
- Context Pack: cp-vX
- Constitution: v3.0
- Status: Draft | In Review | Approved | Blocked
- Blocked By: [ADR-XXX, if applicable]
- Owners: <PM Name>, <Architect Name>

---

## Problem & Success Criteria

### Problem Statement
<What is the user pain? Why does this matter now?>

Example:
"Guest users currently must create an account before completing purchase. This creates friction:
- 23% cart abandonment before login step
- Support load: 40 support tickets/day about 'why must I create account?'
- Competitive disadvantage: competitors offer 1-click guest checkout"

### Success Metric
<How we'll know this worked — measurable KPI>

Example:
- Reduce pre-login cart abandonment from 23% → 8%
- Increase guest checkout volume by 40%
- Reduce checkout time by 60 seconds (10 min → 4 min)
- Support tickets about account creation drop by 90%

### Acceptance Criteria (Given/When/Then format)

```
AC1: Guest can checkout without account
  Given a guest user with items in cart
  When they proceed to checkout
  Then they can complete purchase without creating an account
  And they receive order confirmation via email

AC2: Guest can track order without account
  Given a guest user who completed purchase
  When they visit the order tracking page with order ID + email
  Then they can view order status and tracking info
  And no account is required

AC3: Guest data not persisted after 30 days
  Given a guest checkout order created 31 days ago
  When the privacy job runs
  Then guest email is deleted from the system
  And no trace of guest data remains (except order audit trail)
```

---

## Domain Model

### Domains Involved
- Checkout domain: Owns guest session, cart, payment
- Order domain: Owns order state
- User domain: Owns user account lifecycle
- Notification domain: Owns email delivery
- PII domain: Owns data retention and privacy

### Domain Invariants
Each domain has invariants that MUST be respected:

**Checkout domain invariants:**
- Guest session TTL: 30 minutes of inactivity
- Cart belongs to exactly one session
- Guest session cannot convert to account (by design)

**Order domain invariants:**
- Order has exactly one owner (user_id or guest_id)
- Order total = sum of items (always in sync)
- Once confirmed, order cannot be modified

**User domain invariants:**
- Email is unique per active user
- User account ≠ guest checkout (they are separate)

**PII domain invariants:**
- Guest email deleted after 30 days
- No guest PII logged in plain text
- Retention policy: 30 days for guest, per-user policy for account users

### Cross-Domain Interactions

```
Guest → Checkout (manages session/cart) →
  → Order (creates order) →
    → Notification (sends confirmation email) →
      → PII (schedules email deletion)
```

Contracts (events):
- `GuestSessionCreated` v1 — when guest starts checkout
- `CartConfirmed` v1 — when guest confirms cart
- `GuestOrderCreated` v1 — when guest completes payment
- `GuestOrderConfirmed` v1 — after payment confirmed
- `OrderTrackedRequested` v1 — when guest checks order status

---

## User Experience & Flows

### User Flow: Guest Checkout

```
1. Guest visits site (no account)
2. Adds items to cart
3. Proceeds to checkout
   ├─ Enter email (not password)
   ├─ Confirm shipping address
   ├─ Confirm billing address
   ├─ Choose payment method
   └─ Review order
4. Confirm & pay (Stripe payment modal)
5. See order confirmation page
6. Receive confirmation email with order ID
7. [Later] Visit order tracking page with order ID + email
```

### Edge Cases to Handle

**Guest session expires:**
- Cart is lost after 30 min inactivity
- User can restart checkout (new session, old items not recovered)

**Guest wants to create account:**
- After checkout, guest receives email: "Create account to track orders faster"
- If guest creates account with same email, they see past guest orders in their account

**Payment fails:**
- Guest can retry payment with same email/cart
- No new account created, same session continues

**Duplicate email addresses:**
- If guest email matches existing account, payment succeeds
- But guest order is separate from account (linked by email only)

---

## Technical Approach (High-Level)

### Architecture Decision: Guest Session Isolation

**Decision:** Guest orders are completely separate from user accounts. A guest can never "convert" to a user account linking that order.

**Rationale:**
- Simplifies PII deletion (no accounts to migrate)
- Prevents data leakage (guest data ≠ account data)
- Reduces support burden (no "where's my guest order?" from account users)

**Trade-off:**
- If guest creates account with same email, they won't see guest order history in their account
- Accepted trade-off for privacy/simplicity

### Components

**Checkout Service:**
- Manages guest sessions (in-memory cache + Redis backup)
- Generates guest IDs (UUID format)
- Session TTL: 30 minutes inactivity
- No account lookup (purely anonymous)

**Order Service:**
- Accepts orders with `guest_id` or `user_id`
- Order.owner = guest_id | user_id (mutually exclusive)
- Tracks guest orders in separate table for PII deletion

**Notification Service:**
- Sends order confirmation to guest email
- No email validation (can be fake)

**PII Service:**
- Scheduled job: delete guest records after 30 days
- Runs daily at 2am UTC
- Criteria: `created_at < now() - 30 days AND order_status IN (DELIVERED, CANCELLED)`

---

## Non-Functional Requirements (NFRs)

### Observability

**Logging:**
- Event: `guest.session.created` — guest starts checkout
- Event: `guest.order.created` — guest completes payment
- Event: `guest.data.deleted` — PII deletion job deletes record
- Masking: Guest email masked in all logs: `g****@example.com`

**Metrics:**
- `guest_checkout_attempts_total` — how many guest checkouts started
- `guest_checkout_completion_rate` — % of started → completed
- `guest_session_duration_seconds` — how long from start to payment
- `guest_data_deletion_latency_days` — average days from order to deletion

**Tracing:**
- Span: `guest.checkout.start` → `guest.checkout.confirm_payment` → `guest.checkout.complete`
- Attributes: `guest_id`, `cart_total`, `payment_status`

### Performance

- Checkout flow (add item → payment): p95 < 5 seconds
- Payment processing: p95 < 3 seconds
- Order confirmation email send: p95 < 10 seconds
- Order tracking lookup: p95 < 500ms

### Security

- Guest email not logged in plain text
- Guest orders not accessible without email + order ID
- No account enumeration (can't guess valid emails from guest orders)
- Payment token never stored (Stripe handles PCI compliance)
- Guest email deleted after 30 days (GDPR compliance)

### PII Handling

- Guest email masked in logs: `g****@example.com`
- Guest IP never stored
- Guest payment method details never stored (Stripe owns this)
- Guest orders kept for 30 days, then deleted
- Exception: Order audit trail kept forever (legal requirement)

---

## Contracts & Data Models

### Events (to be emitted by checkout service)

**GuestSessionCreated-v1**
```json
{
  "guest_id": "guest-uuid-123",
  "session_id": "sess-uuid-456",
  "created_at": "2025-02-28T14:00:00Z",
  "ttl_seconds": 1800
}
```

**GuestOrderCreated-v1**
```json
{
  "guest_id": "guest-uuid-123",
  "order_id": "ord-789",
  "total_cents": 9999,
  "currency": "USD",
  "guest_email": "guest@temp.com",
  "created_at": "2025-02-28T14:15:00Z"
}
```

**GuestDataDeleted-v1**
```json
{
  "guest_id": "guest-uuid-123",
  "deletion_reason": "retention_policy_30_days",
  "deleted_at": "2025-03-30T02:00:00Z",
  "audit_trail_retained": true
}
```

### API Contracts (to be implemented by components)

**GET /api/guest/order/{order_id}**
- Query params: `email=guest@example.com`
- Response: `{ order_id, status, items, tracking_url, created_at }`
- Auth: None (email + order ID is the auth token)
- Rate limit: 10 requests per minute per IP

---

## Gates Check (Self-Review Before Approval)

- [ ] **Gate 1 (Context):** Every section has `Source:` line referencing Constitution
- [ ] **Gate 2 (Domain):** No invariant violations, cross-domain communication via contracts only
- [ ] **Gate 3 (Integration):** All contracts identified, consumers listed for each contract
- [ ] **Gate 4 (NFR):** Observability declared (logging, metrics, tracing), security specified, PII handling explicit
- [ ] **Gate 5 (Ready):** No TBD sections, all ACs testable, no blocking ADRs

---

## Next Steps (After Approval)

1. **Create Component Specs** — One per service implementing this feature:
   - Checkout Service Spec (implements guest session + payment)
   - Order Service Spec (implements guest order creation)
   - Notification Service Spec (implements email sending)
   - PII Service Spec (implements 30-day deletion job)

2. **Fan-Out Tasks** — Send to each component team with:
   - Platform Spec ID + version
   - Component Spec ID (which service to implement)
   - Contracts this component produces/consumes
   - Timeline and blockers

3. **Parallel Implementation** — Each component team:
   - Writes Component Spec (must pass all 5 gates)
   - Implements with tests + observability
   - Produces verify.md with evidence

4. **Integration Testing** — After all components done:
   - End-to-end test: guest checkout flow
   - Contract verification: all events flowing correctly
   - PII deletion: verify 30-day cleanup works

5. **Go-Live** — Feature flag OFF initially:
   - Canary: 5% of traffic
   - Monitor: metrics, errors, performance
   - Ramp: 10% → 25% → 50% → 100%
   - Validation: success metrics met

This template is your spec. Customize per your platform's needs.
```

All sections must be completed before this spec can pass Gate 1 (Context Completeness).
