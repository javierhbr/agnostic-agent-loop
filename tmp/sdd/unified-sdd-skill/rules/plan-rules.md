# Plan Rules

Use these rules in the Plan phase.

## Goal

Turn the approved spec into an execution-ready technical plan.

## Apply these skills in this order

1. `bmad-codex-skill`
2. `openspec-codex-skill`
3. `speckit-codex-skill`

## Must do

- map major technical choices back to the approved spec
- define architecture, data flow, interfaces, and testing strategy
- keep platform refs visible in the design when they constrain the solution
- capture dependencies, failure modes, and rollout concerns
- create ADRs when a significant technical decision is introduced
- break the work into reviewable delivery slices
- map tasks to stories or story groups
- produce executable tasks with validation notes

## Avoid

- planning that drifts from the approved spec
- giant vague tasks
- architecture detail that adds no delivery value
- slice boundaries that are too large to review safely

## Required outputs

- `design.md`
- `tasks.md`
- ADRs when needed
- dependency and rollout notes
- delivery slices
- PR strategy notes when useful
- updated `platform-ref.yaml`
- updated `jira-traceability.yaml`

## Exit gate

Move to Deliver only when:

- the delivery team understands the plan
- tasks are executable without large hidden gaps
- validation needs are visible
- slice boundaries are clear and reviewable
