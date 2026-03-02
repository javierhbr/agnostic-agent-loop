---
name: sdd-skills-router
description: Sub-router for SDD v3.0 methodology — maps 15 role-specific skills
---

# sdd/SKILLS.md — SDD v3.0 Skills Router

## Prime Directive

One trigger. One skill per phase. Load only what the current role needs.

---

## SDD SKILL ROUTER

| Trigger | Role | Maps To | Purpose |
|---|---|---|---|
| `sdd:initiative` | Product Manager | `./initiative-definition/SKILL.md` | Define epic problem statements |
| `sdd:risk` | Product Manager | `./risk-assessment/SKILL.md` | Classify initiatives (Low/Medium/High/Critical) |
| `sdd:analyst` | Analyst | `./analyst/SKILL.md` | Conduct discovery interviews, produce discovery.md |
| `sdd:architect` | Architect | `./architect/SKILL.md` | Design feature-spec.md and component-specs |
| `sdd:component-spec` | Component Team | `./component-spec/SKILL.md` | Guide team through creating component impl-spec |
| `sdd:developer` | Developer | `./developer/SKILL.md` | Read component-spec, produce impl-spec.md and tasks |
| `sdd:verifier` | Verifier | `./verifier/SKILL.md` | Verify all ACs, write verify.md, hard stop before merge |
| `sdd:gate-check` | Any role | `./gate-check/SKILL.md` | Diagnose why a gate failed and fix it |
| `sdd:adr` | Architect/Lead | `./adr/SKILL.md` | Create Architecture Decision Records |
| `sdd:hotfix` | Developer | `./hotfix/SKILL.md` | Production hotfix path with minimal spec |
| `sdd:workflow` | Any role | `./workflow-router/SKILL.md` | Route initiative to correct SDD workflow (Quick/Standard/Full) |
| `sdd:constitution` | Platform Architect | `./platform-constitution/SKILL.md` | Author platform governance policy |
| `sdd:platform-spec` | Platform Architect | `./platform-spec/SKILL.md` | Create platform specs with SpecKit MCP |
| `sdd:process-guide` | All roles | `./process-guide/SKILL.md` | Complete step-by-step guide for entire SDD workflow |
| `sdd:stakeholder` | Product Manager | `./stakeholder-communication/SKILL.md` | Communicate SDD status to non-technical stakeholders |

---

## HARD RULES

- **One role per task** — a task belongs to one agent (Analyst, Architect, Developer, Verifier, etc.)
- **Never skip gates** — all 5 gates (Context, Domain, Integration, NFR, Ready) must pass before handoff
- **ADRs block implementation** — no code until the ADR is approved and blocking is cleared
- **Sources required** — every section of a spec MUST cite where it came from (MCP call, PR, existing contract, etc.)

---

## SDD Workflow Routing

Before choosing a skill, route your initiative through the risk classification:

- **LOW risk** → `Quick` workflow (Developer + Verifier only)
- **MEDIUM risk** → `Standard` workflow (Architect + Developers + Verifier)
- **HIGH+ risk** → `Full` workflow (Analyst + Architect + Developers + Verifier + human approval)

Start with `sdd:workflow` or `sdd:risk` to determine which path you're on.

---

## For More Detail

→ `./process-guide/SKILL.md` for the complete SDD methodology walkthrough
