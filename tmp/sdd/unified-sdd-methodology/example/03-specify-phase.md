# Example: Specify Phase

Goal:

- define the required behavior before planning starts

Rule:

- once the work is inside the component repo, use `OpenSpec` only

## Flow

```text
[Routed change]
  PLAT-123 + component epics
        |
        v
[Product leads Specify]
        |
        +--> OpenSpec: proposal.md + delta specs + local traceability
        |
        v
[Approved spec package]
  behavior + platform refs + shared/local decision
```

## `Product`

When to act:

- first in Specify

Skills in order:

1. `openspec-codex-skill`
2. `explain-code-codex-skill`

Example prompts:

- "Using the OpenSpec skill, create proposal.md for validated customer email updates, define the user stories, and write acceptance criteria aligned to the business goal and the customer-identity capability."
- "Using the OpenSpec skill, create delta specs for profile-service and auth-service, and make the shared contract references explicit."
- "Using the explain-code skill, explain the current email-update behavior and highlight the behavior mismatch this new specification is correcting."

Expected outputs:

- `proposal.md`
- delta specs
- acceptance criteria
- clarified assumptions
- confirmed platform refs

## `Team Lead`

When to act:

- during Specify to protect scope and readiness

Skills in order:

1. `openspec-codex-skill`
2. `explain-code-codex-skill`

Example prompts:

- "Using the OpenSpec skill, review this spec package for scope control, platform alignment, and readiness to move into planning."
- "Using the OpenSpec skill, confirm that the component specs reference the correct platform version, contracts, and JIRA issue chain."

Expected outputs:

- ready-for-plan signal
- scope protection notes
- alignment check

## `Architect`

When to act:

- during Specify when hard constraints or shared contracts must be called out

Skills in order:

1. `openspec-codex-skill`
2. `explain-code-codex-skill`

Example prompts:

- "Using the OpenSpec skill, review the validated-email spec package and identify the architectural constraints, dependency boundaries, and contract implications that must be reflected before planning."
- "Using the explain-code skill, explain the current architecture path behind this spec and highlight one design constraint planning must respect."

Expected outputs:

- architectural constraints
- contract implications
- non-functional requirements

## `Developer`

When to act:

- during Specify when edge cases and testability need review

Skills in order:

1. `openspec-codex-skill`
2. `explain-code-codex-skill`

Example prompts:

- "Using the OpenSpec skill, review the approved user stories and turn them into executable acceptance expectations, including edge cases for duplicate emails, invalid formats, and downstream sync failures."
- "Using the explain-code skill, explain the current implementation behavior and call out one testing gotcha that the executable specs must cover."

Expected outputs:

- executable acceptance expectations
- edge-case notes

## Output targets

Use these concrete artifacts as the Specify target:

- `component-repo/proposal.md`
- `component-repo/specs/validated-email-updates/spec.md`
- updated `component-repo/platform-ref.yaml`
- updated `component-repo/jira-traceability.yaml`
