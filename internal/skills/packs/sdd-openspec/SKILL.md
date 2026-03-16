---
name: sdd-openspec
summary: Use this skill for spec-driven development with OpenSpec. It merges the artifact-first design philosophy with strict CLI-first task execution at the component level. Use for proposing changes, generating delta specs, and executing tasks.
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
---

# Unified OpenSpec Skill (The Golden Path)

This skill merges the **Artifact-First** design philosophy with the **CLI-First** execution rigor. It operates seamlessly at the component level using `.sdd-spec/`.

Use this skill when the user wants to:
- implement a change from a spec or requirements file
- plan work using `proposal.md`, `design.md`, `tasks.md`, and delta specs
- safely combine artifact generation with `agentic-agent` task tracking

## Layer 1 — Core Operating Rules (The Artifact Philosophy)

1. **Spec-first workflow:** Align on the change before touching implementation. Prefer artifacts that clarify scope, behavior, approach, and task order.
2. **Hierarchical Source of Truth:** 
   - The *current truth* often comes from the platform (e.g., `vendor/platform-source/specs`), treated as read-only.
   - The *proposed deltas* happen locally in your component (e.g., `.sdd-spec/changes/`).
   - The local component overrides or extends the platform via `agnostic-agent.yaml` configuration.
3. **Think in actions, not locked phases:** You can move forward and backward between proposal, specs, design, and tasks. Update earlier artifacts whenever implementation reveals better information.

## Layer 2 — The Golden Path Execution Workflow (The CLI Bridge)

This is the concrete sequence that merges markdown design with CLI execution.

### Phase 1: Context & Proposal
1. **Context Sync:** Read the platform or component specs to understand current state.
2. **Scaffold Change:** Use `/openspec-proposal` to create a new change package in the local component (e.g., `.sdd-spec/changes/<slug>/proposal.md`).
3. **Delta Specs:** Write the "MODIFIED" or "ADDED" requirements into `.sdd-spec/changes/<slug>/specs/`. Use explicit `MUST/SHALL` behavioral language.

### Phase 2: Design & Tasks
4. **Design:** Write `design.md` if the change requires architecture decisions or tradeoffs.
5. **Task Definition:** Write `tasks.md` in the change folder. Break the work down into small, verifiable steps.

### Phase 3: CLI Ingestion (The Merge Point)
6. **Import:** Run `agentic-agent task list`. The CLI will automatically parse your `.sdd-spec/changes/<slug>/tasks.md` and ingest them into its internal YAML tracking state. *Your fluid design is now locked into strict CLI execution.*

### Phase 4: Execution Loop
7. **Claim:** Run `agentic-agent task claim <ID>`. **(Mandatory: Locks branch, starts trace)**.
8. **Implement:** Write code against the delta specs.
9. **Validate:** Run `agentic-agent validate`. **(Mandatory: Runs gate checks against platform rules)**.
10. **Complete:** Run `agentic-agent task complete <ID>`. Repeat for all tasks.

### Phase 5: Archive
11. **Merge Truth:** Once all tasks are complete, run `/openspec-archive`. This merges the local delta specs into the component's main `.sdd-spec/specs/` directory, becoming the new current truth.

## Layer 3 — Hard Stops (Non-Negotiable)

- **NEVER** edit `.agentic/` YAML files directly. Let the CLI manage state.
- **NEVER** skip `task claim <ID>`. It breaks traceability.
- **NEVER** merge or complete a change without running `agentic-agent validate`.

## References

- Architecture and lifecycle: `references/openspec-overview.md`
- Config, rules, and schemas: `references/config-and-customization.md`
- Artifact writing guide: `references/artifact-authoring.md`
