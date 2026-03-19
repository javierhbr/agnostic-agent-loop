# Example: Plan Phase

Goal:

- convert the approved spec into an implementation-ready component design and
  task plan

Rule:

- once the work is inside the component repo, use `OpenSpec` only

## Flow

```text
[Approved spec package]
        |
        v
[Architect leads Plan]
        |
        +--> OpenSpec: design.md + tasks.md + story mapping
        |
        v
[Implementation-ready plan]
  design + tasks + story mapping
```

## `Architect`

When to act:

- first in Plan

Skills in order:

1. `openspec-skill`
2. `explain-code-skill`

Example prompts:

- "Using the OpenSpec skill, turn the platform plan handoff for validated customer email updates into `design.md` for `profile-service`, including shared contract implications, service boundaries, and rollout constraints."
- "Using the OpenSpec skill, draft `tasks.md` for `PROF-456` and make the platform refs and JIRA story mapping explicit."
- "Using the explain-code skill, explain the planned architecture path with an analogy, an ASCII diagram, a walkthrough, and one implementation gotcha."

Expected outputs:

- `design.md`
- `tasks.md`
- ADRs if needed
- task-to-story mapping

## `Team Lead`

When to act:

- during Plan to validate sequencing and roadmap

Skills in order:

1. `openspec-skill`
2. `explain-code-skill`

Example prompts:

- "Using the OpenSpec skill, break the validated-email work into reviewable delivery slices and map each slice to a story key under PROF-456 and AUTH-234."
- "Using the explain-code skill, explain how the planned slices map to the current code flow and identify the biggest coordination risk."

Expected outputs:

- team roadmap
- delivery slices
- story mapping

## `Product`

When to act:

- during Plan to confirm the design still serves the approved intent

Skills in order:

1. `openspec-skill`
2. `explain-code-skill`

Example prompts:

- "Using the OpenSpec skill, review the design and task plan against the approved user stories and acceptance criteria, and flag any missing behavior."
- "Using the explain-code skill, explain how the planned implementation will satisfy the approved behavior and point out one tradeoff product should understand."

Expected outputs:

- product approval notes
- behavior-gap review

## `Developer`

When to act:

- during Plan to validate that tasks are executable

Skills in order:

1. `openspec-skill`
2. `explain-code-skill`

Example prompts:

- "Using the OpenSpec skill, review the tasks for validated customer email updates and confirm that each task is executable, testable, and small enough for a reviewable PR."
- "Using the explain-code skill, explain the affected code path before implementation begins and highlight one technical gotcha the first slice should avoid."

Expected outputs:

- feasibility review
- task quality feedback

## Output targets

Use these concrete artifacts as the Plan target:

- `component-repo/design.md`
- `component-repo/tasks.md`
- updated `component-repo/platform-ref.yaml`
- updated `component-repo/jira-traceability.yaml`
