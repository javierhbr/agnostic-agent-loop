# Risk Assessment: Subscription Billing Platform

**Initiative ID:** SBL-2025-Q1
**Risk Level Determined:** **CRITICAL**
**Workflow Selected:** **FULL (Analyst → Architect → Developer(s) → Verifier)**
**Estimated Timeline:** 3-4 weeks

---

## Risk Classification Interview

### Question 1: Payment or PII Involved?

**Q:** Does this initiative touch payment processing, authentication, or personally identifiable information (PII)?

**A:** Yes. This initiative directly involves:
- **PII:** Customer email addresses, names, billing addresses
- **Payment Data:** Storing tokenized credit card information for recurring charges
- **Sensitive:** Processing charges automatically without per-charge authorization

**Risk Indicator:** HIGH

---

### Question 2: Breaking Existing Contracts?

**Q:** Does this change any existing API contracts or data models that other services depend on?

**A:** Mostly no, with one caveat:
- Adding new `/subscriptions/` API endpoints — this is additive, not breaking
- BUT: The Payments Service will add a new `is_recurring_charge` flag to charge records — **minor contract evolution**
- Email Service will add new subscription event types — additive webhooks

**Risk Indicator:** MEDIUM (additive changes, no breaking contracts)

---

### Question 3: Multiple Services in Parallel?

**Q:** Will this feature require changes across 3+ services?

**A:** Yes:
- **Billing Service** (new service)
- **Payments Service** (card vault, tokenization)
- **Email Service** (subscription lifecycle emails)
- **Platform API** (endpoint layer)
- **User Service** (webhooks for subscription events)

**Risk Indicator:** HIGH (5 affected components = high coordination risk)

---

### Question 4: Rollback Complexity?

**Q:** If this feature breaks in production, how hard is rollback?

**A:** DIFFICULT:
- If a bug exists in the Billing Service and we disable it, customers with active subscriptions see no charges (revenue loss)
- If the Payments Service has a card tokenization bug, we must re-tokenize all stored cards
- No simple feature flag off switch — rollback requires careful data migration

**Risk Indicator:** CRITICAL (financial impact, data consistency impact)

---

### Question 5: Observability & Monitoring?

**Q:** Do we have clear metrics to detect issues in production?

**A:** Needs explicit design:
- Subscription state transitions must be logged with traceability
- Charge processing must emit metrics: success rate, latency, error breakdown
- Card storage must audit every tokenization and retrieval
- Missing until explicitly specified in the spec

**Risk Indicator:** MEDIUM-HIGH (requires intentional observability spec)

---

## Risk Summary

| Factor | Assessment | Impact |
|--------|-----------|--------|
| Payment/PII | Yes — tokenized cards + addresses | **CRITICAL** |
| Contract Impact | Additive (safe) | Medium |
| Service Coordination | 5 services + new service | **HIGH** |
| Rollback Risk | Difficult (data consistency) | **CRITICAL** |
| Observability | Not yet defined | High |
| **Overall Risk** | **CRITICAL** | **Requires Full Workflow** |

---

## Workflow Justification

**Full Workflow REQUIRED because:**

1. **Payment + PII Touch** — Automatic disqualifies Quick workflow; requires Analyst + Architect rigor
2. **Data Consistency Risk** — Multiple services writing to shared subscription state requires careful contract design
3. **Coordination Complexity** — 5-component fan-out requires Architect's careful component-spec writing
4. **Rollback Uncertainty** — No clear feature-flag boundary; needs Architect's rollback strategy
5. **Regulatory** — Card storage falls under PCI-DSS; requires explicit security and audit spec

**Why not Standard?**
- Standard is for "cross-service, medium risk"
- This is "critical risk, multiple services, high financial impact"
- Analyst phase is essential to confirm with stakeholders that the problem statement and risk understanding are aligned before architecture

---

## Workflow Selection: Full

### Phases

#### Phase 1: Analyst (3-4 days)
- Conduct team interview: confirm problem statement, gather evidence of merchant demand
- Confirm affected components are correct
- Classify risk (completed above: CRITICAL)
- Identify blocking decisions: "Where will we store tokenized cards?" (Stripe or self-hosted vault?)
- Produce `discovery.md` with exit gate checklist

#### Phase 2: Architect (5 days)
- Read `discovery.md`
- Call Platform MCP for payment standards, PCI requirements, card storage policies
- Call Payment Service Component MCP for tokenization patterns
- Design `feature-spec.md`: subscription lifecycle, acceptance criteria, NFRs
- Design component-spec for each affected service:
  - `component-spec-billing.md`
  - `component-spec-payments.md`
  - `component-spec-email.md`
- Create `ADR-001-card-vault.md`: self-hosted vs. Stripe Vault decision (blocks implementation until resolved)

#### Phase 3: Developer (8-10 days, parallel)
- Billing Developer reads `component-spec-billing.md`, produces `impl-spec-billing.md`
- Payments Developer reads `component-spec-payments.md`, produces `impl-spec-payments.md`
- Email Developer (or Payments Developer) reads `component-spec-email.md`, produces `impl-spec-email.md`
- All three wait if ADR-001 blocks them; all three blocked until ADR resolved

#### Phase 4: Verifier (3-4 days)
- Verify all acceptance criteria with tests, logs, metrics
- Confirm PCI audit readiness: all card data handling logged, masked appropriately
- Confirm rollback plan works: can disable subscriptions safely?
- Produce `verify.md` with REQUIRES HUMAN APPROVAL section

---

## Next Steps

1. ✓ Risk classification complete (this document)
2. → Invoke Analyst skill (see `03-analyst/prompts.md`)
3. → Schedule team interview for discovery phase
4. → Set up spec tracking in `.agentic/spec-graph.json`
