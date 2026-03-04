# Analyst Phase: Prompts

## How to Invoke the Analyst Agent

**Context:** You have completed risk assessment and determined workflow = Full. Now hand off to Analyst to conduct deep discovery before architecture begins.

### Prompt for Analyst Agent

```
You are running the Analyst phase of the SDD Full workflow. Your job is to conduct
a one-question-at-a-time team interview about the Subscription Billing initiative and
produce a discovery.md document that the Architect can read and trust.

INITIATIVE: SBL-2025-Q1 - Subscription Billing Platform

CONTEXT: ShopFlow wants to enable merchants to offer subscription plans. This initiative
touches PII and payment data (risk = CRITICAL), so we are running the Full workflow.

YOUR JOB:
1. Call Platform MCP: get_context_pack("Subscription Billing") to get payment policies
2. Interview the team ONE QUESTION AT A TIME about:
   - Why subscriptions matter (evidence: support tickets, merchant requests, market data)
   - Which merchants will benefit most (SMB vs. Enterprise)
   - What contract changes are needed (new endpoints, email events)
   - What rollback looks like if this breaks
   - What observability is critical (which metrics MUST we track?)
3. Classify risk as Low/Medium/High/Critical with evidence
4. Identify blocking architectural decisions (e.g., where to store cards?)
5. Produce discovery.md with all exit gate conditions met

REQUIRED EXIT GATES:
- [ ] Problem statement has a concrete metric (success = 25% of merchants adopt subscriptions)
- [ ] Evidence is a real data point (not assumption): e.g., "43 support tickets in Jan alone"
- [ ] Affected components named explicitly
- [ ] Risk classified (CRITICAL) with rationale
- [ ] Blocking ADRs identified (card vault storage location)
- [ ] Recommended workflow = Full

RULES:
- Ask ONE question per exchange
- Don't overwhelm the team
- Demand evidence, not opinions ("I think..." is not evidence)
- Document the team's answers as you go

START THE INTERVIEW NOW. Ask the first question about the problem statement.
```

---

## Example Analyst Interview (Annotated)

The following is a simulated transcript of the Analyst interviewing the ShopFlow team.

### Exchange 1: Problem Validation

**Analyst (Q1):**
> What user problem does subscription billing solve, and what metric will prove we solved it?

**Team (Answer from Product):**
> We see three problems: (1) Merchants ask for subscriptions in every support ticket, (2) we lose deals to competitors who offer it, (3) churn is high because we can't offer sticky subscription models.

**Analyst (Follow-up):**
> Those are strong signals. Let me be more specific: in January, how many support tickets mentioned "subscription" or "recurring"?

**Team (Evidence):**
> 43 tickets. That's 12% of all support load. And we tracked 8 lost deals explicitly to competitors who had subscriptions.

**Analyst (Note for discovery.md):**
> ✓ Concrete evidence: 43 support tickets (12% of load), 8 lost deals. Metrics: reduce subscription-related support < 2%, capture subscriptions from new merchants.

---

### Exchange 2: Scope & Affected Services

**Analyst (Q2):**
> Which services will this change touch? Walk me through one subscription: from creation to first charge.

**Team (Architect, interim):**
> A merchant creates a plan in the Platform API → stores in Billing Service → on the billing date, Billing Service calls Payments Service to charge the stored card → if successful, Billing updates subscription state → Email Service sends a receipt.

**Analyst (Clarification):**
> When you "charge the stored card," where is the card stored?

**Team:**
> That's TBD. We could use Stripe's vault or build our own vault in Payments Service.

**Analyst (Note for discovery.md):**
> Components: Platform API, Billing Service (new), Payments Service (card storage), Email Service.
> **BLOCKING DECISION:** Where to store tokenized cards? Self-hosted or Stripe vault? Needs ADR.

---

### Exchange 3: Contract Changes

**Analyst (Q3):**
> When Billing Service tells Email Service "send a subscription confirmation," how does that message flow?

**Team (Backend Lead):**
> We'd use the same event system we use for one-time orders: Billing Service emits a `subscription.created` event to the event bus, and Email Service subscribes to it.

**Analyst (Clarification):**
> Will any OTHER service need to know about subscriptions? User Service? Analytics?

**Team:**
> Good point. User Service will want to know so it can mark a customer as "has active subscription" for marketing. Analytics wants to track subscription events for cohort analysis.

**Analyst (Note for discovery.md):**
> New contracts: `subscription.created`, `subscription.renewed`, `subscription.cancelled` events. Consumers: Email Service, User Service, Analytics. Architect must design these contracts explicitly.

---

### Exchange 4: Rollback & Production Risk

**Analyst (Q4):**
> Imagine we launch subscriptions on Monday and discover a critical bug Tuesday morning (e.g., the charge was processed but subscription state wasn't updated). What does rollback look like?

**Team (CTO, interim response):**
> Oof. We can't just flip a feature flag. Customers have active subscriptions. If we disable billing, they see no charges (revenue loss). If we try to re-process, we might double-charge.

**Analyst (Severity check):**
> Right. So rollback is "data-aware" — we need careful steps. What would your runbook be?

**Team:**
> We'd need to: (1) Disable new subscription creation, (2) Keep renewing existing subscriptions (live with the bug), (3) Fix the bug, (4) Reconcile missed charges manually.

**Analyst (Note for discovery.md):**
> CRITICAL RISK: Rollback is data-dependent, not a simple feature flag. Architect must design an explicit "safe disable" mode and reconciliation process. **REQUIRES human approval before merge.**

---

### Exchange 5: Observability Requirements

**Analyst (Q5):**
> In production, which events must we see in logs or metrics to know "subscriptions are working"?

**Team (Ops/SRE):**
> We need to track: (1) subscription creation count (success/failure), (2) charge success rate, (3) charge latency (p50, p95, p99), (4) cards successfully tokenized, (5) errors by type (invalid card, gateway down, etc).

**Analyst (Rigor):**
> Beyond metrics — what about logs? If a charge fails, what details must the log contain?

**Team:**
> Subscription ID, customer ID, amount, card last 4, reason for failure, timestamp. Everything masked appropriately.

**Analyst (Note for discovery.md):**
> OBSERVABILITY SPEC REQUIRED: Explicit logging for each lifecycle event (create, renew, cancel, charge), metrics for success/failure/latency, PII masking rules. Architect must include in feature-spec.

---

## Exit Gate Checklist (from Analyst Interview)

- ✓ Problem statement has concrete metric: 43 support tickets, 8 lost deals, goal = reduce subscription tickets to < 2%
- ✓ Evidence is real data: January support ticket count, lost deals to competitors
- ✓ Affected components listed: Platform API, Billing (new), Payments, Email, User Service, Analytics
- ✓ Risk classified: CRITICAL (payment/PII, data consistency risk, rollback complexity)
- ✓ Blocking ADRs identified: ADR-001 (card vault storage location)
- ✓ Recommended workflow: Full (justified by risk level)

**STATUS:** Ready to hand off to Architect.
