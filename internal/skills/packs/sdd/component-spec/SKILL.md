---
name: component-spec
description: "Guides component team engineers through creating an Sdd-OpenSpec Component Implementation Spec. Trigger when receiving a fan-out task, writing local 'how' for a service, or creating local ADRs."
---

# skill:component-spec

## Does exactly this

Helps component teams write a correct, MCP-grounded Component Implementation Spec — the local "how" derived from the Platform Spec.

---

## When to use

- You received a fan-out task from the Platform Repo with a platform_spec_id
- You need to write your component's implementation spec
- You need to create a local ADR scoped to your component only
- You're unsure what goes in a component spec or how to structure it

---

## Steps — in order, no skipping

1. **Verify required inputs** — You must have: fan-out task, Platform Spec access, Context Pack file, domain MCP file. If any are missing or `blocked_by` is non-empty, stop.

2. **Create the Metadata block** — Open spec with ID, Implements, Context Pack version, Contracts Referenced, Blocked By (empty to start), Status: Draft.

3. **Write each section with a Source: line** — Problem Statement, Goals, Domain Understanding, Cross-Domain Interactions, Contracts, Technical Approach, NFRs, Observability. See resources for per-section guidance.

4. **Run the Gate Check** — All 5 gates must PASS before approval. See resources for gate details and remediation if any fail.

5. **After approval, update spec-graph.json** — Mark status as Approved. Handle local ADRs if needed (see resources).

---

## Output

Component Spec file in `.specify/memory/specs/` or component repo. spec-graph.json updated with status: Approved.

---

## Done when

- All 5 gates PASS
- `Blocked By` is empty
- spec-graph.json shows status: Approved
- Component team can pick up implementation spec without ambiguity

---

## If you need more detail

→ `resources/component-spec-guide.md` — Full section walkthroughs, all 5 gates with checklists, spec-graph.json format, local ADR guidance
