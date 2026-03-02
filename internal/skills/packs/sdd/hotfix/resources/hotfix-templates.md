# Hotfix Templates

Full Bug Spec and Hotfix Spec templates, Follow-up Spec template, and incident response examples. Referenced by the main hotfix skill.

---

## Bug Spec Template (Non-Urgent)

File name: `BUG-<number>-<short-description>.md` (e.g., `BUG-451-cart-checkout-race-condition.md`)

```markdown
# Bug Spec: BUG-<number> — <Short Description>

## Metadata
- ID: BUG-<number>
- Component: <service name>
- Implements: [link to original spec that introduced this behaviour, if known]
- Status: Draft
- Severity: [Low | Medium | High | Critical]

## Reproduction

**Steps to reproduce:**
1. [Step 1]
2. [Step 2]
3. [Step 3]

**Expected:** <what should happen>
**Actual:** <what happens instead>
**Impact:** <how many users affected, revenue impact, error rate>

## Root Cause Hypothesis

<Based on logs, metrics, or code review — what do you think caused this?>

Include relevant error messages, stack traces, or metric spikes if available.

## Fix Plan

<Minimal change to correct the behaviour without introducing new risk>

Be specific: file paths, function names, configuration keys, database migrations.

## Tests

**Unit tests:**
- [ ] Test case 1: <description>
- [ ] Test case 2: <description>

**Integration tests:**
- [ ] Test case 1: <description>

**Regression test:**
<How to verify this specific bug doesn't recur>

## Gate Check (Quick)

- [ ] Fix is scoped to this component only: yes / no
- [ ] No contract changes required: yes / no
- [ ] No breaking API changes: yes / no
- [ ] Logging/observability preserved: yes / no
- [ ] No PII exposure introduced: yes / no

## Sign-off

All tests pass. Fix ready for code review.
```

---

## Hotfix Spec Template (Production, Urgent)

File name: `HOTFIX-<number>-<short-description>.md` (e.g., `HOTFIX-892-payment-processor-down.md`)

**Design principle:** Get this created FAST. Do not let perfection block speed. This is short by design.

```markdown
# Hotfix Spec: HOTFIX-<number> — <Short Description>

## Metadata
- ID: HOTFIX-<number>
- Component: <service name>
- Started: <approximate time incident began>
- Status: In Progress
- Follow-up Spec: [to be created post-incident]

## Issue

**What is failing:** <specific behaviour or error>

Example: "POST /api/checkout returns 503, payments not processing, error rate 98%"

**Impact:** <users affected | revenue | SLA breach | error rate>

Example: "All customer checkouts blocked. Estimated $50k/hour revenue loss."

**Started:** <approximate time>

## Root Cause Hypothesis

<Based on logs/metrics — what do you think caused this?>

Include:
- Error message from logs
- Metric spikes (latency, error rate)
- Recent deployments or changes
- External service dependencies (payment gateway status, etc.)

## Fix

<Minimal change to restore service. Be specific.>

Examples:
- "Revert commit abc123d (added caching layer, causing cache invalidation bug)"
- "Roll back feature flag FeatureFlag_PaymentV2 to 0%"
- "Update config: PAYMENT_TIMEOUT_MS from 1000 to 5000"
- "Restart pod payment-processor-pod-123 (stuck connection pool)"

**Do NOT:** Try to fix the root cause during hotfix. Just restore service.

## Rollback

<Exact steps to revert if this fix makes things worse>

Examples:
```bash
git revert <commit>
kubectl rollout undo deployment/payment-processor
feature flag disable FeatureFlag_PaymentV2
```

**Test this rollback plan BEFORE deploying the fix.**

## Validation

<Which specific metric or log confirms the fix is working>

Examples:
- "Error rate on /api/checkout drops below 0.1%"
- "Log event 'payment.authorized' appears in Datadog within 30 seconds"
- "p95 latency on POST /orders returns to < 500ms"
- "Payment processing queue depth decreases from 50k to < 100"

**You MUST verify this metric/log in production before declaring the hotfix successful.**

## Follow-up Spec

[REQUIRED after incident is resolved]

Link: [to be filled in]

Create this within 1 business day of incident resolution.
```

---

## Follow-up Spec Template

File name: `FOLLOWUP-<hotfix-number>-hardening.md` (e.g., `FOLLOWUP-892-payment-hardening.md`)

Create this within 1 business day of incident resolution. Do NOT skip this.

```markdown
# Follow-up Spec: FOLLOWUP-<n> — Hardening

## Metadata
- ID: FOLLOWUP-<n>
- Related Hotfix: HOTFIX-<n>
- Status: Draft

## What the Hotfix Did

<Summary of the emergency fix in 2-3 sentences>

Example: "Reverted the caching layer added in commit abc123d that caused cache invalidation bugs. Restored service availability and payment processing within 15 minutes."

## What Still Needs to Be Done

### [ ] Full Test Coverage

<Add comprehensive tests for the fixed behaviour>

- Unit tests for [specific function]
- Integration tests for [specific workflow]
- Stress tests for [concurrent operations]
- Edge case tests for [boundary conditions]

### [ ] Refactor (If Applicable)

<If the fix introduced tech debt, document it>

Example: "The emergency revert left the caching layer disabled. We need to re-implement it with proper invalidation logic and test coverage (target: 1 sprint)."

### [ ] ADR (If Decision Was Made Under Pressure)

<Document any architectural decision made during incident>

Example: "ADR-234: Cache invalidation strategy for payment processor. Why we chose X over Y under time pressure. Why we need to revisit this."

### [ ] Contract Spec (If a Contract Was Touched)

<Document if any event/API contract changed or was violated>

Example: "Contract Spec for PaymentProcessor/CheckoutAPI v2. The hotfix temporarily disabled async payment processing; need to specify new contract."

### [ ] Domain MCP Update (If Invariant Was Missing)

<Document if the bug revealed a missing invariant>

Example: "Payment domain invariant: Idempotency key must be stable across retries. This was missing from Domain MCP; payment processor assumed it."

### [ ] Platform MCP Update (If Policy Was Insufficient)

<Document if the bug revealed a policy gap>

Example: "Platform Constitution now requires: All services with external dependencies must have circuit breaker + fallback documented in Constitution/Resilience section."

### [ ] Post-Mortem Notes

<What went wrong, why the gate didn't catch this, how to prevent recurrence>

Template:
- **What was the root cause?** [answer]
- **Why did the gate not catch this?** [answer]
- **What invariant was missing?** [answer]
- **What test class would have caught this?** [answer]
- **How can we prevent this class of bug?** [answer]

Example:
```
Root cause: Cache invalidation bug introduced in caching layer. No tests verified cache consistency under concurrent updates.

Why Gate 4 didn't catch it: NFR section declared "Caching enabled" but didn't specify invalidation strategy or test requirements.

Missing invariant: "Cache TTL ≤ 5 minutes OR explicit invalidation on write" (now in Domain MCP).

Test class: Concurrent cache update tests with race condition injection.

Prevention: Constitution now requires "All caching must be tested with concurrent update scenarios" in Gate 4 (NFR).
```

## Gate Check

- [ ] All acceptance criteria from original spec still met
- [ ] Tests added for the fixed behaviour
- [ ] No regressions introduced
- [ ] All 5 gates pass on this Follow-up Spec

## Sign-off

Ready for implementation within the next sprint.
```

---

## Incident Response Examples

### Example 1: Payment Processor Down (Critical)

**Incident:** POST /api/checkout returns 503, payment processing blocked

**HOTFIX-892 created immediately:**
- Issue: "Payment processor service returned 503 errors, checkout unavailable"
- Root Cause Hypothesis: "Recent deploy of caching layer (commit abc123d) caused cache invalidation bug, all cache lookups fail"
- Fix: "Revert commit abc123d (remove caching layer temporarily)"
- Rollback: "git revert abc123d"
- Validation: "Error rate on /api/checkout drops below 0.1% within 2 minutes"

**Deployed:** 3 minutes after incident started
**Service restored:** 5 minutes (verification)
**FOLLOWUP-892 created:** 1 day later

Follow-up covers:
- Full test coverage for cache invalidation logic
- ADR documenting cache strategy decision made under pressure
- Constitution update requiring concurrent cache tests
- Post-mortem: "Why no tests verified cache consistency?"

---

### Example 2: Data Corruption Bug (High)

**Incident:** User orders showing wrong totals in dashboard

**BUG-451 created:**
- Reproduction: "1. Create order with 3 items. 2. View order total. 3. Order total ≠ sum of items"
- Root Cause Hypothesis: "Integer overflow in calculation or stale read from replica database"
- Fix Plan: "Add database constraint and fix calculation with proper decimal types"
- Tests: "Race condition test with concurrent order updates, replica lag test"

**Status:** In progress, no immediate impact to new orders, old orders show garbage data

---

### Example 3: Security Issue (Critical)

**Incident:** PII (user emails) visible in audit logs unmasked

**HOTFIX-893 created immediately:**
- Issue: "User email addresses appearing in plain text in audit logs"
- Fix: "Mask email field at API boundary before writing to audit log"
- Validation: "Grep audit logs for email patterns — none found"

**Deployed:** 2 minutes
**FOLLOWUP-893 created:** 1 day later

Follow-up covers:
- Full audit log masking test suite
- Constitution update: "PII must be masked at API boundary, enforced in Gate 4"
- ADR: Why we log emails at all (compliance requirements)
- Post-mortem: "How did unmasked PII reach logs?"

---

## Hotfix Governance

### When to Use Hotfix Path

✅ **Production incident** — Service down, data loss, security breach, PII exposure
✅ **Payment/auth/critical** — Any issue in payment, auth, or critical user path
✅ **Need speed** — Restore service ASAP, document thoroughly afterward
✅ **No time for TDD** — Write Hotfix Spec, fix, validate, then Follow-up Spec with tests

### When to Use Bug Path (Non-Urgent)

✅ **Low severity** — Non-critical feature broken, only affects small number of users
✅ **No immediate impact** — Bug exists but workaround available
✅ **Can wait for tests** — Schedule for next sprint with full test coverage
✅ **Non-critical path** — Admin feature, internal tool, non-customer-facing

---

## Key Reminders

1. **Hotfix Spec takes 10 minutes to write** — Do NOT skip it. It's your insurance policy.
2. **Fix is minimal** — Restore service, don't solve root cause during hotfix.
3. **Rollback plan is required** — Test it before deploying the fix.
4. **Validation metric is non-negotiable** — You MUST verify recovery in production.
5. **Follow-up Spec is mandatory** — Create it within 1 business day. This is where the system learns.
6. **Post-mortem is required** — Why did the gate not catch this? What's the fix to prevent it next time?
