# Refactoring Plan: Unified SDD Methodology & Delegate Architecture

## 1. Vision & Strategy
Transition the `agentic-agent` CLI from a "monolithic generator" to a **Phase-based Orchestrator**. The CLI will no longer duplicate logic; it will delegate to specialized authority tools via adapters.

### The Authority Tools
- **BMAD**: Routing, sizing, track selection.
- **Spec Kit**: Platform constitution, clarity/checklist, durable rules.
- **OpenSpec**: Component proposal, specs, design, tasks, archive.
- **Delivery Planner (such as JIRA)**: Initiative/Epic/Story/PR traceability.
- **Agentic-Agent**: Context, task claiming, validation, worktree isolation, and sync/bridge.

---

## 2. Source-of-Truth Contract
Exactly which tool owns each phase and its primary artifacts.

| Phase | Authority Tool | CLI Responsibility | Key Artifacts |
| :--- | :--- | :--- | :--- |
| **Platform** | **Spec Kit** | Context & Rule hydration | `constitution.md`, `config.yaml` |
| **Route** | **BMAD** | Track selection & Scoping | `change-package.yaml`, `platform-ref.yaml` |
| **Specify** | **OpenSpec** | Artifact bridging & Clarity | `proposal.md`, `delta-specs/` |
| **Plan** | **Planner (such as JIRA)** | Local Task -> Planner Sync | `design.md`, `planner-traceability.yaml` |
| **Deliver** | **Agentic-Agent** | Implementation & Validation | `tasks.md`, PRs, `archive/` |

---

## 3. Core Architecture Changes

### A. Adapter Boundaries (`internal/adapters/`)
Create thin wrappers for external CLIs to avoid chasing their internal models.
- `OpenSpecAdapter`: Wraps `opsx` CLI.
- `SpecKitAdapter`: Wraps `speckit` CLI.
- `BMADAdapter`: Wraps `bmad` logic/CLI.
- `PlannerAdapter`: Interface for project management (JIRA, Linear, GitHub).

### B. Artifact-Aware Validation
Refactor `internal/validator` to verify:
1. **Platform Alignment**: Code matches `platform-ref.yaml` constraints.
2. **Behavioral Integrity**: Code matches OpenSpec `delta-specs`.
3. **Workflow State**: Local task state matches the **Delivery Planner (such as JIRA)**.

---

## 4. Command Rationalization

### Commands to DELETE / Deprecate
- `sdd`: **DELETE**. Replaced by phase-specific commands.
- `track`: **Deprecate**. Replaced by `route`.
- `spec`: **Deprecate**. Replaced by `specify`.
- `work`: **DELETE**. Functionality moved to `deliver`.
- `start`: **DELETE**. Redundant.

### NEW / Repurposed Commands
- `platform`: Initializes Spec Kit rules and Platform links.
- `route`: The BMAD intake gateway (Size vs. Impact).
- `specify`: OpenSpec proposal + Spec Kit clarity pass.
- `plan`: Design creation + **Sync to Delivery Planner (such as JIRA)**.
- `deliver`: The execution loop (Build -> PR -> Review -> Verify -> Deploy -> Archive).
- `sync`: Bidirectional state sync between local OpenSpec and the Planner.

---

## 5. Implementation Phases

### Phase 1: The Contract (Adapters)
1. Define interfaces in `internal/adapters/interfaces.go`.
2. Implement `Mock` adapters for testing the orchestrator loop.
3. Update `pkg/models/config.go` for Planner credentials and Tool paths.

### Phase 2: Iteration 1 (Front-Half)
1. Implement `platform` command (Spec Kit bridge).
2. Implement `route` command (BMAD bridge).
3. Implement `specify` command (OpenSpec + Spec Kit bridge).

### Phase 3: Iteration 2 (Back-Half)
1. Implement `PlannerAdapter` for JIRA/Linear.
2. Implement `plan` command (Sync to Planner).
3. Implement `deliver` dashboard (Mapping tasks to implementation slices).

### Phase 4: Cleanup
1. Remove `internal/sdd` (legacy logic).
2. Update `internal/tasks` to be a sub-package of `deliver`.
3. Final verification using `tmp/sdd/unified-sdd-methodology/example`.

---

## 6. Pilot Verification Case
- **Input**: A raw requirement in `tmp/pilot-request.md`.
- **Expected Outcome**:
  1. `route` selects "Standard Track".
  2. `specify` creates an OpenSpec package.
  3. `plan` syncs tasks to a mock JIRA JSON file.
  4. `deliver` completes one task and archives the result.
