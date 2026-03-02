---
name: creating-development-plans
description: Create structured development plans with phased task breakdowns, requirements, and QA checklists
---

# skill:dev-plans

## Does exactly this

Create a detailed, actionable development plan broken into independent phases with tasks, technical decisions, testing strategy, and QA checklists. Output a single `DEVELOPMENT_PLAN.md` file ready for agent execution.

---

## When to use

- User explicitly asks to "create a development plan", "plan the work", or "create a dev plan"
- User provides requirements and wants structured breakdown before implementation
- Project has complex scope needing phased delivery with human review checkpoints

---

## Steps — in order, no skipping

1. **Gather Context** — Read existing code, docs, tests, architecture, conventions, tech stack.

2. **Analyse Requirements** — Identify goals, hard requirements, unknowns, edge cases, integration points, trade-offs.

3. **Break Down Tasks** — Organize into 3-5 phases with 3-8 atomic tasks per phase. Identify dependencies. Add self-review and STOP checkpoints per phase.

4. **Plan QA** — Build checklist with standard items + project-specific checks + technology-specific rules.

5. **Deep Review** — Ultrathink soundness, architecture, constraints, requirements alignment, missing items. Adjust plan. Ensure British English.

---

## Output

Single file: `DEVELOPMENT_PLAN.md`

Contains:
- Project Purpose and Goals
- Context and Background
- Development Tasks (phases + tasks with checkboxes)
- Important Considerations & Requirements
- Technical Decisions
- Testing Strategy
- Debugging Protocol
- QA Checklist

---

## Done when

- `DEVELOPMENT_PLAN.md` file created and saved
- All phases have tasks, self-review, and STOP checkpoints
- Task count 10-20 for medium feature (aggressive splitting applied)
- QA checklist is complete and project-specific
- Technical decisions documented with trade-offs
- Plan ready for agent execution (clear enough for independent execution)

---

## Core Principles

- **Planning before code** — Fully understand context and requirements first
- **Phased approach** — Discrete, independently testable, reviewable phases with human checkpoints
- **Atomic tasks** — One concern per task. Split aggressively (frontend/backend = 2 tasks, infrastructure/logic = 2 tasks).
- **Simplicity over complexity** — No unnecessary abstractions
- **Actionable output** — Clear enough for another AI to execute independently

---

## If you need more detail

→ `resources/plan-template.md` — Full DEVELOPMENT_PLAN.md template, task granularity rules, phase structure, writing guidelines, QA gates by risk, testing philosophy

