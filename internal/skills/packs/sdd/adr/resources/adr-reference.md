# ADR Reference

Full ADR template, state machine details, rejection flow, and anti-patterns. Referenced by the main adr skill.

---

## Full ADR Template

File name: `ADR-<number>-<short-kebab-title>.md` (e.g., `ADR-219-idempotency-strategy.md`)

```markdown
# ADR-<number>: <Title>

## Status
Proposed

## Date
<YYYY-MM-DD>

## Context

What is the situation? Why is a decision needed now? What constraints exist?

Example:
"When a guest checks out and submits an order, the CheckoutService sends an OrderCreated event to downstream subscribers. If the event publish fails midway (network timeout), CheckoutService retries the publish. But downstream services have already started processing based on the first publish attempt. We need to decide: what is our idempotency strategy?"

## Decision Drivers

- <Driver 1: e.g., consistency with Platform MCP security policy>
- <Driver 2: e.g., performance target of p95 < 300ms>
- <Driver 3: e.g., existing contract consumers who cannot be broken>

Example:
- Existing consumers of OrderCreated v1 are already deployed and cannot change behavior immediately
- Platform NFR requires p95 event latency < 500ms
- Event deduplication must not block the happy path

## Options Considered

### Option A: <Name>
<Description>

Pros:
- <pro 1>
- <pro 2>

Cons:
- <con 1>
- <con 2>

### Option B: <Name>
<Description>

Pros:
- <pro 1>

Cons:
- <con 1>

## Decision

[Leave as PENDING until approved — do not fill in prematurely]

Example (when approved):
"We will implement idempotency using event.idempotency_key. CheckoutService generates a stable UUID (hash of order_id + timestamp) as the idempotency_key and includes it in every OrderCreated event. Downstream subscribers use this key to deduplicate. CheckoutService may emit the same event multiple times on retry, but each subscriber will process it only once."

## Consequences

[Fill in after decision is made]

Example:
"Downstream services must add idempotency_key to their event handlers. Existing consumers of OrderCreated v0 (without idempotency_key) will receive events with an extra field and must treat it as optional. A compatibility plan is needed for v0 → v1 migration."

## Owner

<Name / team responsible for resolving this ADR>

## Blocks

- <SPEC-ID or PLAT-ID that is waiting on this ADR>

Example:
```
- SPEC-CHECKOUT-01
- SPEC-FULFILLMENT-02
```

## References

- <Link to relevant Platform MCP section, if applicable>
- <Link to relevant contract, if applicable>

Example:
```
- Platform MCP v2.0 — Integration Contracts section
- Integration MCP — OrderCreated event schema v1
```
```

---

## ADR State Machine

ADRs move through these states:

```
Proposed → In Review → Approved
                    → Rejected
```

### Proposed

An ADR is drafted. Needs review. Fill in Context, Decision Drivers, Options A & B. Leave Decision and Consequences blank. Assign an Owner.

### In Review

Owner has been assigned. Stakeholders (PM, Platform Architect, Domain Owner) are reviewing. **This state is required before any fan-out task proceeds** — it signals the decision is being actively worked.

### Approved

Decision is made. Fill in the Decision and Consequences sections. Unblock all dependent specs by removing `BlockedBy` entries.

### Rejected

Decision went another way. Affected specs must be updated to reflect the rejected path or sent back to Design.

---

## Resolving an ADR

### When the ADR Reaches Approved

1. **Fill in Decision and Consequences** — Document the chosen option and its implications
2. **Set status to Approved** — Update the ADR file header
3. **Unblock all specs** — In every spec's Metadata block, remove `Blocked By: ADR-<number>` or set it to `[]`
4. **Update spec-graph.json** — For each spec, set `blocked_by: []` and `status` to the prior state (usually `Approved` or `Draft`)
5. **Notify component teams** — "ADR-219 is resolved. Your spec is unblocked. Recheck Gate 5."

### When the ADR is Rejected

1. **Document the rejection reason** — Update the ADR file with notes on why Option A/B was rejected
2. **Update each blocked spec** — Reflect the rejected path or send spec back to Design
3. **If rejection requires fundamental redesign** — Spec may need to go back to Draft state

---

## Linking ADR to Specs in spec-graph.json

Example: Three specs blocked by ADR-219

```json
{
  "SPEC-CHECKOUT-01": {
    "implements": "PLAT-124 v1",
    "blocked_by": ["ADR-219"],
    "status": "Blocked"
  },
  "SPEC-FULFILLMENT-02": {
    "implements": "PLAT-125 v1",
    "blocked_by": ["ADR-219"],
    "status": "Blocked"
  },
  "SPEC-INVENTORY-01": {
    "implements": "PLAT-126 v1",
    "blocked_by": ["ADR-219"],
    "status": "Blocked"
  },
  "ADR-219": {
    "title": "idempotency-strategy",
    "status": "In Review",
    "blocks": ["SPEC-CHECKOUT-01", "SPEC-FULFILLMENT-02", "SPEC-INVENTORY-01"],
    "owner": "platform-architect@example.com"
  }
}
```

When ADR-219 is approved, update:

```json
{
  "SPEC-CHECKOUT-01": {
    "implements": "PLAT-124 v1",
    "blocked_by": [],
    "status": "Approved"
  },
  "SPEC-FULFILLMENT-02": {
    "implements": "PLAT-125 v1",
    "blocked_by": [],
    "status": "Approved"
  },
  "SPEC-INVENTORY-01": {
    "implements": "PLAT-126 v1",
    "blocked_by": [],
    "status": "Approved"
  },
  "ADR-219": {
    "title": "idempotency-strategy",
    "status": "Approved",
    "blocks": [],
    "owner": "platform-architect@example.com"
  }
}
```

---

## ADR Anti-Patterns

**❌ Don't fill in the Decision before review.** Premature decisions bypass the review process. Stakeholders must weigh in on Options A & B before you commit.

**❌ Don't create a Global ADR for a local concern.** A Global ADR blocks unrelated components. If the decision is scoped to Component Repo A only, make it Local (Component Repo A `/adr/`).

**❌ Don't leave ADRs in Proposed indefinitely.** If no owner is assigned, assign one now. If it's been Proposed for > 1 week, escalate.

**❌ Don't skip the ADR because "we'll figure it out in code."** That is exactly the tribal knowledge this system prevents. The ADR forces consensus _before_ coding.

**❌ Don't close an ADR as Approved without updating all blocked specs.** Orphaned `BlockedBy` entries stall implementations silently. Always update all 3:
  1. ADR status → Approved
  2. Spec Metadata → remove `BlockedBy`
  3. spec-graph.json → both entries → unblocked

---

## Common ADR Topics from /speckit.clarify

These are the most frequent ambiguities that become ADRs:

| Ambiguity | ADR Topic |
|---|---|
| How are retried operations made safe? | Idempotency strategy |
| When does a guest session expire? | Session lifecycle / TTL policy |
| Which service generates a shared ID? | ID ownership |
| How does Service B handle duplicate events? | Event deduplication strategy |
| Which version of the contract is the canonical one? | Contract versioning decision |
| Can Service A read Service B's data directly? | Domain boundary decision |
| What happens when the external payment provider is down? | Failure mode / circuit breaker strategy |
