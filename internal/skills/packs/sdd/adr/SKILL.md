---
name: adr
description: "Create, track, and resolve Architecture Decision Records that gate implementation. ADRs are first-class artifacts that block specs until resolved."
---

# skill:adr

## Does exactly this

Help create, track, and resolve ADRs that gate implementation in the SDD Operating Model. An ADR is not optional documentation — it is a first-class artifact that blocks implementation until resolved.

---

## When to use

- You need to create an ADR to resolve a cross-cutting or component-local decision
- `/speckit.clarify` surfaces an unresolved question
- A spec is `BlockedBy` an ADR
- You're unsure whether a decision should be Global (Platform Repo) or Local (Component Repo)
- Reviewing specs before fan-out to identify blocking decisions early

---

## Step 1 — Determine ADR Scope

**Is this decision cross-cutting or component-local?**

| Question | If Yes → |
|---|---|
| Does it affect more than one component/domain? | **Global ADR** → Platform Repo `/adr/` |
| Does it affect a platform policy (security, NFR, versioning)? | **Global ADR** → Platform Repo `/adr/` |
| Does it affect a contract or shared event schema? | **Global ADR** → Platform Repo `/adr/` |
| Does it only affect one component's internal implementation? | **Local ADR** → Component Repo `/adr/` |

**Rule:** Global blocks any dependent component. Local blocks only this component.

---

## Steps 2-5 — Draft → Link → Review → Resolve

1. **Draft the ADR** — File: `ADR-<number>-<short-title>.md`. Use template. Status: Proposed. See resources for full template.

2. **Link to blocked specs** — Add `Blocked By: ADR-<number>` to each spec Metadata. Update spec-graph.json with `blocked_by` list.

3. **Move through review** — States: Proposed → In Review → Approved (or Rejected). In Review is required before fan-out.

4. **Resolve the ADR** — When Approved: fill Decision + Consequences, remove `BlockedBy` from all specs, update spec-graph.json, notify teams.

---

## Common ADR Triggers

These are the most frequent ambiguities that need ADRs:

| Ambiguity | ADR Topic |
|---|---|
| How are retried operations made safe? | Idempotency strategy |
| When does a guest session expire? | Session lifecycle / TTL policy |
| Which service generates a shared ID? | ID ownership |
| How does Service B handle duplicate events? | Event deduplication strategy |
| Which version of the contract is canonical? | Contract versioning decision |
| Can Service A read Service B's data directly? | Domain boundary decision |
| What happens when external provider is down? | Failure mode / circuit breaker strategy |

---

## Anti-Patterns to Avoid

- Don't fill in the Decision before review — bypasses the review process
- Don't create a Global ADR for a local concern — it blocks unrelated components
- Don't leave ADRs in Proposed indefinitely — assign an owner now
- Don't skip the ADR because "we'll figure it out in code" — that is tribal knowledge
- Don't close an ADR as Approved without updating all blocked specs — orphaned BlockedBy entries stall silently

---

## If you need more detail

→ `resources/adr-reference.md` — Full ADR template with all fields, state machine details, rejection flow, spec-graph.json linking, detailed anti-patterns
