# How to Use the SDD Process (Unified 5-Phase Model)

The Specification-Driven Development (SDD) process is a five-phase methodology
that takes a change from platform setup through delivery. Each phase has a
designated role, specific artifacts, and CLI commands.

This guide walks through all five phases step by step.

---

## Core Concepts Before You Start

**Change package** -- the canonical execution unit. Every change flows through
a change package directory under `.agentic/sdd/changes/<name>/`. It holds
proposals, delta specs, designs, tasks, traceability files, and evidence.

**Boundary rule** -- Platform repos use BMAD + OpenSpec + Speckit. Component
repos use OpenSpec ONLY. If you are working in a component repo, skip Speckit
and BMAD classification; go straight to OpenSpec for your change package.

**Two-iteration adoption** -- Teams new to SDD adopt in two iterations:

- Iteration 1: Platform + Assess + Specify (phases 1-3)
- Iteration 2: Plan + Deliver (phases 4-5)

Do not attempt all five phases on your first project.

---

## Phase 1: Platform (Architect)

Goal: Establish the durable platform baseline that all future changes reference.

### Step 1.1 -- Draft the platform constitution

Use Speckit to produce a constitution document. This captures the non-negotiable
decisions for your platform: tech stack, deployment topology, observability
standards, security posture, and API conventions.

```bash
# Ensure SDD skills are installed
agentic-agent skills ensure
```

The constitution lives at `.agentic/sdd/platform/constitution.md`. Write it once,
update it rarely.

### Step 1.2 -- Encode durable configuration

Create `config.yaml` using OpenSpec. This file captures versions, environment
references, and integration points that remain stable across changes.

```text
.agentic/sdd/platform/config.yaml
```

### Step 1.3 -- Write ownership artifacts

Create three ownership files in `.agentic/sdd/platform/`:

| File                              | Purpose                                                             |
| --------------------------------- | ------------------------------------------------------------------- |
| `component-ownership-<name>.md`   | One per component. Lists owner, repo, API surface, SLOs.            |
| `dependency-map.md`               | Directed graph of component dependencies with impact tiers.         |
| `glossary.md`                     | Canonical terms used across specs. Prevents ambiguity in proposals. |

Impact tiers in `dependency-map.md` classify each dependency as:

- **Tier 1** -- Hard constraint. Breaking this dependency blocks the change.
- **Tier 2** -- Rollout risk. Can be mitigated with sequencing or feature flags.
- **Tier 3** -- Informational. Downstream consumers should be notified.

### When is Phase 1 done?

You have a constitution, a config.yaml, ownership files for every component,
a dependency map with impact tiers, and a glossary. These artifacts are
referenced by every subsequent phase.

---

## Phase 2: Assess (Team Lead)

Goal: Classify the incoming work, identify the owning component, and open a
change package with traceability.

### Step 2.1 -- Classify using BMAD

Answer three questions:

1. **Greenfield or brownfield?** New component vs. modifying existing code.
2. **Small, medium, or large?** Based on number of components affected.
3. **Quick, standard, or full?** Based on risk level and review depth needed.

The classification determines which artifacts are mandatory vs. optional in
later phases.

### Step 2.2 -- Identify the primary owner

Read the relevant `component-ownership-<name>.md` files from Phase 1. The
primary owner is the team or person listed in the ownership file for the most
affected component.

### Step 2.3 -- Open a change package

Use OpenSpec to create the change package directory:

```text
.agentic/sdd/changes/<change-name>/
```

### Step 2.4 -- Populate platform-ref.yaml

Create `platform-ref.yaml` inside the change package. It references the
platform baseline:

```yaml
platform_version: "1.2.0"
constitution_ref: "../../platform/constitution.md"
config_ref: "../../platform/config.yaml"
impact_tiers:
  tier_1:
    - component: auth-service
      constraint: "Must maintain backward-compatible token format"
  tier_2:
    - component: notification-service
      risk: "Email templates change; rollout behind feature flag"
```

Pull impact tiers directly from `dependency-map.md`.

### Step 2.5 -- Populate jira-traceability.yaml

Link the change package to your project management system:

```yaml
epic: "PROJ-100"
stories:
  - key: "PROJ-101"
    summary: "Add guest checkout flow"
  - key: "PROJ-102"
    summary: "Update payment adapter for guest users"
```

### Step 2.6 -- Review the backlog

```bash
agentic-agent task list
```

Confirm no conflicting changes are in progress for the same components.

### When is Phase 2 done?

You have a BMAD classification, identified owner, open change package with
`platform-ref.yaml` and `jira-traceability.yaml`, and no conflicts in the
backlog.

---

## Phase 3: Specify (Product)

Goal: Define what changes, using delta specs anchored to the glossary.

### Step 3.1 -- Check the glossary

Before writing any spec, read `glossary.md`. Use its terms exactly. If you
need a new term, add it to the glossary first and get agreement.

### Step 3.2 -- Write proposal.md

Create `proposal.md` inside the change package:

```markdown
# Proposal: <Change Name>

## Problem
What is broken or missing, and for whom.

## Goals
What success looks like, stated as measurable outcomes.

## Non-Goals
What is explicitly out of scope.

## Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2
```

Keep it concise. The proposal is the contract between Product and the rest
of the team.

### Step 3.3 -- Write delta specs

For each affected component, write a delta spec file inside the change package.
Delta specs use three sections:

```markdown
# Delta Spec: <component-name>

## ADDED
- New endpoint POST /checkout/guest
- New event GuestCheckoutCompleted

## MODIFIED
- PATCH /checkout now accepts optional `guest_token` field

## REMOVED
- Nothing removed in this change
```

Delta specs describe differences from the current baseline, not the full
component specification.

### Step 3.4 -- Scope and alignment check

The Team Lead reviews the proposal and delta specs for:

- Scope creep (does the delta match the proposal goals?)
- Glossary compliance (are terms used consistently?)
- Impact alignment (do delta specs cover all tier-1 dependencies?)

### Boundary rule reminder

Component repos use OpenSpec ONLY. Do not introduce Speckit or BMAD artifacts
into component-level specifications.

### When is Phase 3 done?

You have a reviewed proposal.md, delta specs for every affected component,
glossary updated if needed, and Team Lead sign-off on scope.

---

## Phase 4: Plan (Architect)

Goal: Turn the specification into a concrete design and task breakdown.

### Step 4.1 -- Read impact tiers before designing

Open `platform-ref.yaml` from the change package. Understand tier-1 hard
constraints before making any design decisions.

### Step 4.2 -- Write design.md

Create `design.md` inside the change package:

```markdown
# Design: <Change Name>

## Approach
High-level design describing how the delta specs will be implemented.

## Tier 1 Constraints
List each hard constraint and how the design satisfies it.

## Tier 2 Rollout Risks
List each rollout risk and the mitigation strategy (feature flags, phased
rollout, backward compatibility, etc.).

## Sequence
Order of implementation if there are dependencies between tasks.
```

### Step 4.3 -- Write tasks.md

Break the design into discrete tasks, mapped to story keys from
`jira-traceability.yaml`:

```markdown
# Tasks

| ID | Description | Story Key | Owner | Estimate |
|----|------------|-----------|-------|----------|
| TASK-001 | Implement guest checkout endpoint | PROJ-101 | auth-team | 3d |
| TASK-002 | Update payment adapter | PROJ-102 | payments-team | 2d |
| TASK-003 | Add integration tests | PROJ-101 | auth-team | 1d |
```

### Step 4.4 -- Developer feasibility review

Developers review the design and tasks for:

- Technical feasibility of the approach
- Accuracy of estimates
- Missing tasks or dependencies

Adjust design.md and tasks.md based on feedback before moving to Phase 5.

### When is Phase 4 done?

You have a reviewed design.md with constraint/risk analysis, tasks.md mapped
to story keys, and developer sign-off on feasibility.

---

## Phase 5: Deliver (Team Lead)

Goal: Build, review, verify, deploy, and archive the change.

### Step 5.1 -- Build

For each task in tasks.md:

```bash
# Claim the task
agentic-agent task claim TASK-001

# Implement the change
# ... write code, tests, documentation ...

# Validate before completing
agentic-agent validate

# Mark complete
agentic-agent task complete TASK-001
```

### Step 5.2 -- Create the pull request

Link the PR to the change package, tasks, and stories:

```bash
agentic-agent pr create --task TASK-001
```

The PR description should reference:

- Change package path
- Task ID from tasks.md
- Story key from jira-traceability.yaml

### Step 5.3 -- Review the pull request

Reviewers check for:

- **Design integrity** -- Does the code match design.md?
- **Scope check** -- Does the code stay within the delta spec boundaries?
- **Tier-1 compliance** -- Are hard constraints respected?

### Step 5.4 -- Verify

Record evidence that acceptance criteria from proposal.md are met. Store
verification artifacts (test results, screenshots, metrics) inside the
change package under a `evidence/` directory.

### Step 5.5 -- Deploy

Coordinate deployment timing, respecting:

- Dependencies between tasks (sequence from design.md)
- Tier-2 rollout risks (feature flags, phased rollout)
- Rollback plan for each component

### Step 5.6 -- Archive the change package

After successful deployment:

1. Merge delta specs into the component baseline specifications.
2. Close the change package (mark status as `completed`).
3. Update `dependency-map.md` if the change altered component relationships.
4. Update `component-ownership-<name>.md` if ownership changed.

### When is Phase 5 done?

All tasks are complete, PRs merged, acceptance criteria verified with evidence,
deployment successful, delta specs merged into baselines, and change package
closed.

---

## Quick Reference

| Phase        | Role      | Key Artifacts                                                                    | Key CLI Commands                                                                                                      |
| ------------ | --------- | -------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------- |
| 1. Platform  | Architect | constitution.md, config.yaml, ownership files, dependency-map.md, glossary.md    | `agentic-agent skills ensure`                                                                                         |
| 2. Assess    | Team Lead | BMAD classification, platform-ref.yaml, jira-traceability.yaml                   | `agentic-agent task list`                                                                                             |
| 3. Specify   | Product   | proposal.md, delta specs                                                         | --                                                                                                                    |
| 4. Plan      | Architect | design.md, tasks.md                                                              | --                                                                                                                    |
| 5. Deliver   | Team Lead | Code, PRs, evidence, archived package                                            | `agentic-agent task claim ID`, `agentic-agent validate`, `agentic-agent task complete ID`, `agentic-agent pr create`  |

---

## Troubleshooting

**"Which phase am I in?"** -- Check the change package directory. If it has
only platform-ref.yaml, you are in Assess. If it has proposal.md, you are in
Specify. If it has design.md and tasks.md, you are in Plan. If tasks are
claimed, you are in Deliver.

**"Do I need all five phases?"** -- For component repos, skip Phase 1 (use the
existing platform baseline) and simplify Phase 2 (no BMAD, just open the change
package). For small/quick changes, the Assess classification will indicate
which artifacts are optional.

**"Where are the old commands?"** -- Commands like `specifyify start`,
`specifyify gate-check`, and `specifyify sync-graph` are replaced by the
unified workflow above. Use `agentic-agent task` and `agentic-agent validate`
instead.
