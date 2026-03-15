---
name: speckit-spec-driven-development
description: Use for turning product ideas into executable specifications, plans, task lists, and implementation workflows using Spec Kit / Speckit. Triggers on spec-driven development, requirements, constitution, clarify, plan, tasks, implement, brownfield/greenfield workflows, and agent-guided delivery.
---

# Speckit / Spec Kit Skill

## Purpose
Use this skill when the user wants to:
- define a feature before coding
- create a spec, plan, tasks, and implementation flow
- establish project rules and engineering principles
- run greenfield or brownfield work with a predictable workflow
- reduce vague prompting and turn intent into structured artifacts

This skill follows a **3-layer context model**:
- **Layer 1:** Operating rules and non-negotiables
- **Layer 2:** Workflow/playbook for how to execute the work
- **Layer 3:** Deep reference material and templates

---

# Layer 1 — Operating Rules

## 1) Specs define the "what" before the "how"
Start from user intent, outcomes, constraints, and acceptance behavior.
Do not jump straight into implementation details unless the workflow is already in the planning phase.

## 2) Constitution first for meaningful work
Before generating major specs or plans, define the project constitution:
- engineering principles
- quality bar
- testing expectations
- UX or accessibility constraints
- security / performance / compliance guardrails

## 3) Use the workflow in order unless there is a strong reason not to
Preferred sequence:
1. constitution
2. specify
3. clarify
4. checklist (optional but recommended)
5. plan
6. tasks
7. analyze (recommended)
8. implement

## 4) Clarify ambiguity before planning
If the requirements are underspecified, conflicting, or missing scope boundaries, resolve that through a clarification pass before producing technical plans.

## 5) Keep artifacts separated by responsibility
- constitution = governing principles
- spec = requirements, user stories, scope, acceptance
- plan = architecture and technical decisions
- tasks = executable implementation steps
- analyze/checklist = quality and consistency validation

## 6) Prefer phased delivery for large features
For complex work, split into small phases so the agent does not saturate context and so each phase can be validated before continuing.

## 7) Protect the specs directory
Do not treat `specs/` as disposable scaffolding. It is the durable source of intent, plans, and task decomposition.

---

# Layer 2 — Workflow Playbook

## Standard playbook

### A. Constitution
Create or refresh project principles.
Use when:
- starting a new project
- entering a repo with weak engineering rules
- the user wants explicit standards for all future work

Expected output:
- concise list of governing principles
- clear quality bar
- explicit tradeoffs and priorities

### B. Specify
Capture the feature as a product/problem statement.
Focus on:
- user goals
- required behaviors
- constraints
- non-goals
- acceptance conditions

Avoid:
- framework decisions
- storage choices
- premature API or schema design

### C. Clarify
Run this when the spec contains ambiguity.
Clarify:
- actors and permissions
- edge cases
- data lifecycle
- failure behavior
- performance/security expectations
- out-of-scope boundaries

### D. Checklist
Use a checklist when you need an "English unit test" for the spec.
Verify completeness, consistency, and testability before planning.

### E. Plan
Convert the approved spec into implementation strategy.
Include:
- architecture
- major components
- technical choices
- integration points
- data flow
- testing approach
- rollout notes if relevant

### F. Tasks
Break the plan into sequenced, executable tasks.
Tasks should be:
- concrete
- reviewable
- small enough for an agent or engineer to complete safely
- traceable back to the plan and spec

### G. Analyze
Run cross-artifact analysis before implementation.
Look for:
- missing requirements coverage
- contradictions between spec/plan/tasks
- hidden dependencies
- missing validation or testing work

### H. Implement
Only implement after the spec and plan are strong enough.
Prefer incremental execution, validating at each phase.

---

## Greenfield mode
Use when building from scratch.
Pattern:
- define constitution
- write a high-signal spec
- clarify missing behavior
- create plan
- break into tasks
- implement in phases

## Brownfield mode
Use when extending an existing codebase.
Add to the workflow:
- map current system constraints first
- identify compatibility and migration risks
- keep tasks small and reversible
- be explicit about what changes versus what stays untouched

## Output format when applying this skill
When you answer with this skill, structure the result as:
1. **Phase** — where the user currently is
2. **Recommended next artifact**
3. **Questions or clarifications needed**
4. **Proposed artifact content**
5. **Risks / assumptions**
6. **Next command or next step**

---

# Layer 3 — Deep References

Use these references as needed:
- `references/overview-and-philosophy.md`
- `references/workflow-and-commands.md`
- `references/installation-and-init.md`
- `references/upgrade-and-safety.md`
- `references/templates.md`
- `references/sources.md`

Also use the reusable rule files:
- `rules/constitution-rules.md`
- `rules/spec-rules.md`
- `rules/plan-rules.md`
- `rules/task-rules.md`
