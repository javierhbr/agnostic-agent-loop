---
name: op
summary: Use this skill for spec-driven development with OpenSpec, including proposing changes, generating and updating artifacts, applying tasks, archiving completed changes, configuring context and rules, and customizing schemas.
triggers:
  - openspec
  - opsx
  - spec-driven development
  - sdd
  - openspec/project.md
  - openspec/specs
  - openspec/changes
  - proposal.md
  - design.md
  - tasks.md
  - delta specs
  - openspec/config.yaml
  - schema.yaml
  - /openspec-proposal
  - /openspec-apply
  - /openspec-archive
  - /opsx:propose
  - /opsx:apply
  - /opsx:archive
---

# OpenSpec Spec-Driven Development

Use this skill when the user wants to:

- use OpenSpec or OPSX commands correctly
- plan work before implementation
- create or refine `proposal.md`, `design.md`, `tasks.md`, or delta specs
- set up `openspec/project.md` or `openspec/config.yaml`
- add project context or per-artifact rules
- customize a schema or fork `spec-driven`
- explain OpenSpec workflows to another coding agent

## Layer 1 — Core operating rules

1. Separate **core OpenSpec** from the **experimental OPSX workflow**.
   - Core OpenSpec uses `openspec/project.md`, `openspec/specs/`, `openspec/changes/`, and `openspec/changes/archive/`.
   - In Codex, prefer `/openspec-proposal`, `/openspec-apply`, and `/openspec-archive` for the default core flow.
   - Treat `/opsx:*`, `openspec/config.yaml`, and schema customization as extended or experimental workflow features, not the default assumption.

2. Treat OpenSpec as a **spec-first workflow**.
   - Align on the change before touching implementation.
   - Prefer artifacts that clarify scope, behavior, approach, and task order.

3. Think in **actions, not locked phases**.
   - You can move forward and backward between proposal, specs, design, and tasks.
   - Update earlier artifacts whenever implementation reveals better information.

4. Keep the **source of truth** explicit.
   - Main specs describe the current system behavior.
   - Change folders describe proposed deltas.
   - Archive merges approved deltas back into main specs.

5. Use the right command for the job.
   - Use `/openspec-proposal` for the normal Codex quick path.
   - Use `/openspec-apply` to implement approved tasks while keeping artifacts aligned.
   - Use `/openspec-archive` after the change is complete.
   - Use `/opsx:explore` only for explicit discovery work when the team is already using OPSX.
   - Use `/opsx:new`, `/opsx:continue`, `/opsx:ff`, `/opsx:apply`, and `/opsx:archive` only when the repo has opted into the expanded OPSX profile.

6. Use the canonical OpenSpec tree when writing artifacts.
   - Durable repo context lives in `openspec/project.md`.
   - Current truth lives in `openspec/specs/<capability>/spec.md`.
   - Proposed changes live in `openspec/changes/<change-slug>/`.
   - A core change package normally contains `proposal.md`, `tasks.md`, optional `design.md`, and `specs/<capability>/spec.md`.

7. Prefer stable project context over ad hoc instruction drift.
   - Put durable project guidance in `openspec/project.md`.
   - Use `openspec/config.yaml` only when the workflow explicitly supports extra config and rules.
   - Keep context concise enough to stay useful and maintainable.

8. Write artifacts so another agent can execute them.
   - Requirements should be concrete and testable.
   - Design should record decisions, tradeoffs, and affected systems.
   - Tasks should be actionable, ordered, and easy to verify.

## Layer 2 — Execution workflow

### A. Default quick path

Use this when the user wants the normal OpenSpec flow.

1. Clarify the change goal.
2. Create or update the change with `/openspec-proposal`.
3. Review generated artifacts:
   - `proposal.md`
   - `specs/` delta specs
   - `tasks.md`
   - optional `design.md`
4. Refine artifacts if scope, requirements, or design need changes.
5. Implement with `/openspec-apply`.
6. Validate the change.
7. Archive with `/openspec-archive`.

### B. Expanded OPSX workflow

Use this only when the project has enabled the experimental or expanded OPSX profile.

1. `/opsx:explore` for investigation or framing.
2. `/opsx:new` to scaffold a change only.
3. `/opsx:continue` to generate the next missing artifact.
4. `/opsx:ff` to fast-forward planning artifacts.
5. `/opsx:apply` to implement tasks.
6. `/opsx:verify` or project validation steps when relevant.
7. `/opsx:sync` if the workflow uses sync.
8. `/opsx:archive` when done.

### C. Artifact decision rules

#### When to write a proposal

Write or refine `proposal.md` when the user needs:

- problem framing
- goals and non-goals
- rollout or rollback thinking
- affected scope and stakeholders

#### When to write delta specs

Write `specs/` deltas when the user is changing behavior.
Use explicit sections such as:

- `ADDED Requirements`
- `MODIFIED Requirements`
- `REMOVED Requirements`

Use the exact heading pattern:

- `### Requirement: ...`
- `#### Scenario: ...`

Each requirement should use concrete SHALL or MUST language, and each scenario should be written in Given/When/Then form.

#### When to write design

Write `design.md` when the change needs:

- architecture decisions
- data flow or sequence explanation
- tradeoffs
- interfaces, dependencies, risks, or rollout design

#### When to write tasks

Write `tasks.md` when implementation must be broken into verifiable steps.
Tasks should:

- be small enough to complete predictably
- map back to requirements or design decisions
- include validation where useful

### D. Project context and optional config

Use `openspec/project.md` for durable project guidance in the core workflow.
Use `openspec/config.yaml` only when the repo explicitly uses extended config or OPSX-specific customization.

Recommended pattern:

- `openspec/project.md` stores shared project conventions and durable context
- `schema:` selects the default workflow schema when config is enabled
- `context:` stores shared project conventions when config is enabled
- `rules:` adds per-artifact guidance when config is enabled

When asked to generate config:

1. capture stable tech stack and conventions in `context:`
2. keep artifact rules scoped under `rules:`
3. avoid stuffing ephemeral ticket details into config
4. prefer schema defaults unless the team truly needs a custom workflow

### E. Custom schema workflow

Use a custom schema only when the default schema is not enough.

1. Start from the built-in `spec-driven` schema.
2. Fork it into `openspec/schemas/<name>/`.
3. Modify `schema.yaml` only for real workflow differences.
4. Modify templates to shape output quality.
5. Keep custom workflows understandable to both humans and agents.

## Layer 3 — References

- Architecture and lifecycle: `references/openspec-overview.md`
- Commands and workflow mapping: `references/commands-and-workflows.md`
- Config, rules, and schemas: `references/config-and-customization.md`
- Artifact writing guide: `references/artifact-authoring.md`
- Source notes: `references/sources.md`
- Reusable project rules: `rules/project-config-template.yaml`
- Reusable artifact rules: `rules/artifact-rules.md`

## Response behavior when applying this skill

When you use this skill in an answer:

- identify whether the user needs explore, propose, apply, archive, or customization help
- keep terminology aligned with OpenSpec artifacts and commands
- produce concrete artifact content, not vague theory
- distinguish current specs from delta specs
- distinguish core OpenSpec from experimental OPSX behavior
- recommend config and rules when conventions should persist across future changes
- mention whether the advice assumes the default `core` profile or an expanded profile
