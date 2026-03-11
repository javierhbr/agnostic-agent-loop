# Verifier Templates

Full verify.md template, evidence examples, and AC-to-evidence mapping. Referenced by the main verifier skill.

---

## Verify.md Full Template

File: `.agentic/specs/[component-spec-id]/verify.md`

```markdown
# Verification Report: [Component Spec ID]

## Metadata
- Spec ID: SPEC-SERVICE-001
- Verified by: [Your name]
- Date: [Date]
- Status: [PASSED | BLOCKED]

## Acceptance Criteria Verification

### AC1: [Exact criterion from component spec]
Status: PASS ✓
Evidence: [test name or command that produces evidence]
Output: [actual result showing AC is satisfied]

### AC2: [Next exact criterion]
Status: PASS ✓
Evidence: [test name]
Output: [actual result]

[... repeat for each AC ...]

## Code Quality Checks
- [ ] All tests pass: ✓
- [ ] Lint passes: ✓
- [ ] Build passes: ✓
- [ ] Coverage ≥ 80%: ✓
- [ ] No REQUIRES HUMAN APPROVAL items: ✓

## Observability Verification
- [ ] Logging implemented: ✓ (see log sample below)
- [ ] Metrics exported: ✓ (see metric sample below)
- [ ] Tracing enabled: ✓ (see trace sample below)
- [ ] Alerts configured: ✓

## Sign-Off
All acceptance criteria verified with evidence. Ready for production deployment.
```

---

## What Counts as Evidence

Evidence is observable, reproducible proof that an AC is satisfied. Different types:

### Unit Test Evidence
```
AC1: Cart deletion removes all line items

Evidence: test_cart_deletion() in tests/cart_test.go
Output: PASS (0.23s)
```

### Integration Test Evidence
```
AC2: Order endpoint returns 201 Created with order_id in response

Evidence: POST /orders (tests/integration/order_creation_test.go)
Output:
  curl -X POST http://localhost:8080/orders -d '{"items": [...]}' -H "Content-Type: application/json"
  {
    "order_id": "ORD-12345",
    "status": "PENDING",
    "created_at": "2025-02-28T14:23:45Z"
  }
  HTTP 201 Created
```

### Metrics Evidence
```
AC3: Metrics record order processing latency with p95 < 300ms

Evidence: Prometheus query: histogram_quantile(0.95, order_processing_duration_ms)
Output:
  order_processing_duration_ms{quantile="0.95"} = 245.8
```

### Log Evidence
```
AC4: Each order emit a structured log event with order_id, user_id, amount

Evidence: tail -f logs/orders.log | grep "order.created"
Output:
  {
    "timestamp": "2025-02-28T14:23:45Z",
    "event": "order.created",
    "order_id": "ORD-12345",
    "user_id": "USR-789",
    "amount": 99.99
  }
```

### Lint/Build Evidence
```
AC5: Code passes linting and builds without errors

Evidence: make lint && make build
Output:
  All checks passed. Build output in ./bin/service
  No errors or warnings.
```

---

## Acceptance Criteria Mapping

For each AC in the component spec, trace which code changes satisfy it:

### Example Mapping

**AC1: Cart deletion removes all line items**
→ impl-spec code change: `impl/cart.go` lines 45-67 (DeleteCart function)
→ evidence: test_cart_deletion() PASS

**AC2: Order endpoint returns 201 Created**
→ impl-spec code change: `impl/order_handler.go` lines 12-35 (CreateOrder handler)
→ evidence: integration test POST /orders PASS + HTTP response

**AC3: Metrics recorded with p95 < 300ms**
→ impl-spec code change: `impl/order_handler.go` line 28 (metrics.RecordLatency call)
→ evidence: Prometheus histogram query shows p95 = 245.8ms

---

## Worked Example: Single AC End-to-End

**Component Spec, AC3:**
```
Given a guest checkout session, when the user submits the order,
then the order is created with status PENDING and returned within 300ms p95 latency.
```

**Implementation Spec, Code Changes:**
- File: `pkg/checkout/order_service.go` lines 34-78
  - CreateOrder function: validates session, creates order row, publishes OrderCreated event
  - Latency instrumentation: defer metrics.RecordLatency("checkout.order.create")(ctx)

**Verification Steps:**

1. **Run the integration test:**
   ```bash
   go test -run TestCheckoutOrderCreation ./tests/integration
   Output: PASS (0.18s)
   ```

2. **Check the metrics (live system test):**
   ```bash
   curl -s http://localhost:8081/metrics | grep 'checkout_order_create_duration_seconds_bucket'
   Output:
     checkout_order_create_duration_seconds_bucket{le="0.1"} 145
     checkout_order_create_duration_seconds_bucket{le="0.3"} 2847
     checkout_order_create_duration_seconds_bucket{le="+"} 2850

   p95 calculation: 2847/2850 = 99.9% of requests < 300ms ✓
   ```

3. **Check the logs:**
   ```bash
   tail -f logs/checkout.log | grep "order.created"
   Output:
     {"ts":"2025-02-28T14:32:15Z","event":"order.created","order_id":"ORD-98765","session_id":"SESS-abc123","status":"PENDING","latency_ms":87}
   ```

**Verify.md Entry:**
```
### AC3: Order created within 300ms p95

Status: PASS ✓
Evidence: Integration test + Prometheus histogram + logs
Output:
  - TestCheckoutOrderCreation PASS (0.18s)
  - Prometheus p95 = 245ms (99.9% of requests < 300ms)
  - Sample log: {"event":"order.created","order_id":"ORD-98765","latency_ms":87}
```

---

## REQUIRES HUMAN APPROVAL Cases

Mark and escalate if the change touches:

### Payment Processing
```
**REQUIRES HUMAN APPROVAL** — Payment processing logic change

What: Stripe charge logic modified (pkg/payment/stripe.go lines 45-67)
Why: Financial transaction — must be manually verified by Payment Owner
How to escalate: Set Status = REQUIRES HUMAN APPROVAL, ping #payment-squad
```

### Authentication/Authorization
```
**REQUIRES HUMAN APPROVAL** — Session token validation changed

What: JWT validation logic (pkg/auth/jwt.go)
Why: Security boundary change — must be reviewed by Security team
How to escalate: Set Status = REQUIRES HUMAN APPROVAL, ping #security-team
```

### PII Handling
```
**REQUIRES HUMAN APPROVAL** — User data masking logic changed

What: Email masking at API boundary (pkg/api/response.go)
Why: PII exposure risk — must be reviewed by Privacy/Compliance
How to escalate: Set Status = REQUIRES HUMAN APPROVAL, ping #privacy-squad
```

---

## Blocking Conditions Checklist

**Do not merge if ANY are true:**

- [ ] Any AC is UNTESTABLE (reword the AC to be testable, or open an ADR)
- [ ] Any test FAILS (debug and fix root cause)
- [ ] Lint or build FAILS (fix all warnings/errors)
- [ ] Coverage < 80% (add tests for uncovered branches)
- [ ] Change touches payment/auth/PII without human approval (escalate and get sign-off)
- [ ] Spec Graph not updated to Done (run `agentic-agent specifyify sync-graph`)

**If ANY are true:** BLOCK the merge. Return the spec to the Architect with specific remediation steps.
