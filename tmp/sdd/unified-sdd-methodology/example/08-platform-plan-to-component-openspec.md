# Example: Platform Plan to Component OpenSpec

Goal:

- show how a component team turns a platform `Plan` handoff into local
  OpenSpec artifacts and local delivery work

## Core rule

- platform planning can use the combined platform methodology
- once the work enters the component repository, use `OpenSpec` only

## Flow

```text
[Platform Plan output]
  platform version
  platform refs
  shared decisions
  rollout constraints
  JIRA chain
        |
        v
[Component repo aligns]
  platform-ref.yaml
  jira-traceability.yaml
        |
        v
[OpenSpec local package]
  proposal.md
  spec.md
  design.md
  tasks.md
        |
        v
[OpenSpec delivery updates]
  story status -> PR -> verification -> archive
```

## Platform handoff example

Assume the platform `Plan` phase hands this to `profile-service`:

- platform version: `2026.03`
- platform issue: `PLAT-123`
- component epic: `PROF-456`
- platform refs:
  - `capabilities.customer-identity`
  - `contracts.customer-profile.v2`
- shared rules:
  - validate email before persistence
  - keep contract `v2` backward compatible
  - emit explicit validation errors and metrics

## Step 1: Pin alignment in the component repo

Use OpenSpec to record the platform baseline before writing local detail.

Example prompts:

- `Team Lead`
  - "Using the OpenSpec skill, create `platform-ref.yaml` for `profile-service` from platform version `2026.03`, and record the required refs and shared rules for this change."
- `Team Lead`
  - "Using the OpenSpec skill, create `jira-traceability.yaml` for `PLAT-123` and `PROF-456`, and map the expected local stories for this component."

Expected outputs:

- `component-repo/platform-ref.yaml`
- `component-repo/jira-traceability.yaml`

## Step 2: Turn platform intent into a local proposal

Use OpenSpec to convert shared planning intent into local component intent.

Example prompts:

- `Product`
  - "Using the OpenSpec skill, create `proposal.md` for `profile-service` from the platform plan for validated email updates. State what changes in this component, what does not change, and how the behavior stays aligned with `capabilities.customer-identity`."

Expected output:

- `component-repo/proposal.md`

## Step 3: Turn required behavior into a local component spec

Use OpenSpec to define the local behavior in testable terms.

Example prompts:

- `Product`
  - "Using the OpenSpec skill, write the component spec delta for validated email updates in `profile-service`, including success behavior, failure behavior, and compatibility expectations with `contracts.customer-profile.v2`."
- `Developer`
  - "Using the OpenSpec skill, refine the local component spec with explicit scenarios for invalid email, duplicate email, and downstream sync failure."

Expected output:

- `component-repo/specs/validated-email-updates/spec.md`

## Step 4: Turn the shared decisions into a local design

Use OpenSpec to explain how the component will implement the shared rules.

Example prompts:

- `Architect`
  - "Using the OpenSpec skill, draft `design.md` for `profile-service` from the platform plan handoff. Show where validation happens, how contract `v2` remains compatible, and which metrics or logs must be emitted."
- `Team Lead`
  - "Using the OpenSpec skill, review `design.md` and confirm that the local design matches the platform refs, team boundaries, and planned delivery slices."

Expected output:

- `component-repo/design.md`

## Step 5: Break the design into local tasks and stories

Use OpenSpec to create reviewable, local execution slices.

Example prompts:

- `Team Lead`
  - "Using the OpenSpec skill, create `tasks.md` for `profile-service` from the approved local design, and map each task to a reviewable story under `PROF-456`."
- `Developer`
  - "Using the OpenSpec skill, review the proposed tasks and make each one small enough for a single reviewable PR."

Expected output:

- `component-repo/tasks.md`

## Step 6: Deliver through local OpenSpec updates

Use OpenSpec to keep local delivery traceability current.

Example prompts:

- `Developer`
  - "Using the OpenSpec skill, update task status for `PROF-789`, link the PR, and record the validation performed for the current slice."
- `Team Lead`
  - "Using the OpenSpec skill, review the component change package and confirm it is ready for PR review, deploy, and archive."
- `Product`
  - "Using the OpenSpec skill, review the delivered component behavior against the approved proposal and spec, and record any acceptance gaps before deploy."

Expected outputs:

- updated `component-repo/tasks.md`
- updated `component-repo/jira-traceability.yaml`
- `component-repo/pr-description.md`
- archive-ready local OpenSpec artifacts

## What good looks like

The handoff is working when:

- the component repo pins platform version and refs before local design starts
- the local proposal, spec, design, and tasks all reference the same platform baseline
- component roles use `OpenSpec` only for local artifacts
- JIRA, PRs, and archive records stay aligned to the same component change package
