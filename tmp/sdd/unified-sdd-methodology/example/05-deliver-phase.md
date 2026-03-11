# Example: Deliver Phase

Goal:

- execute the plan through build, PR, review, verification, deploy, and archive

Rule:

- once the work is inside the component repo, use `OpenSpec` only

## Flow

```text
[Tasks + stories]
        |
        v
[Team Lead coordinates]
        |
        +--> Developer: build and create PR
        +--> Architect: review design integrity
        +--> Product: verify delivered behavior
        |
        v
[Reviewed and verified delivery]
  stories -> PRs -> validation -> deploy -> archive
```

## `Team Lead`

When to act:

- throughout Deliver

Skills in order:

1. `openspec-codex-skill`
2. `explain-code-codex-skill`

Example prompts:

- "Using the OpenSpec skill, monitor the current delivery slices for PROF-456 and AUTH-234, keep task status current, and summarize which stories are ready for PR."
- "Using the explain-code skill, explain the current PR scope and review focus so reviewers understand the exact change."

Expected outputs:

- current slice status
- PR coordination
- release-readiness summary

## `Developer`

When to act:

- during Build, PR creation, and review resolution

Skills in order:

1. `openspec-codex-skill`
2. `explain-code-codex-skill`

Example prompts:

- "Using the OpenSpec skill, implement the current PROF-789 task slice according to the approved local spec and keep the change narrow enough for one reviewable PR."
- "Using the OpenSpec skill, update the task state, link the current story and PR, and record the validation performed."
- "Using the explain-code skill, explain this PR change with an analogy, an ASCII diagram, a step-by-step walkthrough, and one reviewer gotcha."

Expected outputs:

- code
- tests
- PR description
- updated task and story links

## `Architect`

When to act:

- during review and when design integrity is at risk

Skills in order:

1. `openspec-codex-skill`
2. `explain-code-codex-skill`

Example prompts:

- "Using the OpenSpec skill, review this PR for design drift, contract safety, and alignment with the approved platform refs."
- "Using the explain-code skill, explain how this PR changes the architecture path and highlight one risk reviewers should watch."

Expected outputs:

- architecture review notes
- approval or correction requests

## `Product`

When to act:

- during final behavior review and release decisions

Skills in order:

1. `openspec-codex-skill`
2. `explain-code-codex-skill`

Example prompts:

- "Using the OpenSpec skill, review the delivered behavior against the approved stories and acceptance criteria for validated customer email updates."
- "Using the explain-code skill, explain how the delivered code behaves now and highlight one business-facing behavior to verify before deploy."

Expected outputs:

- acceptance review
- go / no-go input

## Output targets

Use these concrete artifacts as the Deliver target:

- `component-repo/pr-description.md`
- updated `component-repo/jira-traceability.yaml`
- final archive-ready OpenSpec artifacts
