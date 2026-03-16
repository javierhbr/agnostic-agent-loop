---
name: skills-router
description: Master router for all CLI skills — maps triggers to files
---

# SKILLS.md — Master Skills Router

## Prime Directive

One trigger. One skill. One output. Load nothing else.

---

## SKILL ROUTER

| Trigger | Maps To | Purpose |
|---|---|---|
| `sdd-openspec` | `./sdd-openspec/SKILL.md` | Spec-driven development from requirements (Unified SDD) |
| `api-docs` | `./api-docs/SKILL.md` | Generate comprehensive API documentation |
| `tdd` | `./tdd/SKILL.md` | Test-driven development orchestrator |
| `atdd` | `./atdd/SKILL.md` | Acceptance test-driven development |
| `product-wizard` | `./product-wizard/SKILL.md` | Generate PRDs for products and features |
| `dev-plans` | `./dev-plans/SKILL.md` | Create structured development plans |
| `code-simplification` | `./code-simplification/SKILL.md` | Review code for simplicity and maintainability |
| `extract-wisdom` | `./extract-wisdom/SKILL.md` | Extract insights from text sources |
| `context-manager` | `./context-manager/SKILL.md` | Enforce mandatory context workflow |
| `run-with-ralph` | `./run-with-ralph/SKILL.md` | Execute tasks using Ralph Wiggum loops |
| `agentic-helper` | `./agentic-helper/SKILL.md` | CLI guide, executor, and automator |
| `diataxis` | `./diataxis/SKILL.md` | Apply Diataxis framework to documentation |
| `superpowers-bridge` | `./superpowers-bridge/SKILL.md` | Bridge CLI and Superpowers plugin |
| `tier-check` | `./tier-enforcer/SKILL.md` | Audit, create, or fix skill files for 3-tier compliance |
| `tier-enforcer` | `./tier-enforcer/SKILL.md` | Audit, create, or fix skill files for 3-tier compliance |
| `sdd:*` | `./sdd/SKILLS.md` | SDD v3.0 sub-router (15 role-specific skills) |
| `sdd-bmad` | `./sdd-bmad/SKILL.md` | Progressive planning and role-based execution (Unified SDD) |
| `explain-code` | `./explain-code-skill/SKILL.md` | Visual diagrams and analogies for codebase architecture |
| `platform-contextualizer` | `./platform-contextualizer-skill/SKILL.md` | Documenting existing platforms |
| `sdd-speckit` | `./sdd-speckit/SKILL.md` | Turning ideas into executable specs and task lists (Unified SDD) |
| `unified-sdd` | `./unified-sdd-skill/SKILL.md` | Platform-scale approach combining Sdd-Bmad, Sdd-OpenSpec, and Sdd-Speckit |

---

## HARD RULES

- **One trigger per invocation** — never chain skills inside a skill file
- **Never load more than ONE skill per task** — if a task touches two skills, split into two invocations
- **No inline examples** — examples live in `resources/` files, linked and loaded only when needed
- **No tech stack details in skill files** — those live in playbooks or references, linked from skills
- **Skills are routers, not manuals** — triggers map to files that are ~1 page, with detail extracted to references

---

## SDD Sub-Router

For all `sdd:*` triggers, see `./sdd/SKILLS.md` which maps 15 SDD v3.0 skills (Analyst, Architect, Developer, Verifier, ADR, Gate-Check, Component-Spec, Workflow-Router, Risk-Assessment, Initiative-Definition, Platform-Constitution, Platform-Spec, Process-Guide, Hotfix, Stakeholder-Communication).

---

## For More Detail

→ `./sdd/SKILLS.md` for SDD methodology
