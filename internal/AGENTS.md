# AGENTS.md - internal

## Purpose

- Contain core application logic that powers the `agentic-agent` CLI.
- Implement orchestration, tasks, config handling, skills, validation, project initialization, and TUI support.

## Key Entrypoints

- `internal/config/config.go`: load YAML config and apply defaults.
- `internal/project/init.go`: initialize `.agentic` project structure from embedded templates.
- `internal/tasks/manager.go`: task lifecycle operations (claim, complete, list, decompose).
- `internal/orchestrator/loop.go` and `internal/orchestrator/autopilot.go`: work-loop and autopilot orchestration.
- `internal/skills/ensure.go`, `internal/skills/generator.go`, `internal/skills/detect.go`: agent detection and skill generation/installation.
- `internal/validator/validator.go`: validation pipeline and rule execution.

## File Map

- `internal/agents/`: agent execution abstractions and multi-agent execution.
- `internal/checkpoint/`: checkpoint manager logic.
- `internal/config/`: config loading/default merging.
- `internal/context/`: context generation/rolling/global updates.
- `internal/encoding/`: bundle/toon encoding utilities.
- `internal/gitops/`: git tracking utilities.
- `internal/openspec/`: OpenSpec change management and templates.
- `internal/orchestrator/`: iterative run/autopilot flow and state transitions.
- `internal/plans/`: plan parsing, generation, and updating.
- `internal/project/`: project scaffolding and embedded templates.
- `internal/prompts/`: builtin prompt registry.
- `internal/simplify/`: code/context simplification bundle helpers.
- `internal/skills/`: skill registry, packs, ensure/install/symlink logic.
- `internal/specs/`: spec reference resolution.
- `internal/status/`: status dashboard builders.
- `internal/tasks/`: task lifecycle, locking, decomposition, progress writing.
- `internal/token/`: token counting, budgets, summarization.
- `internal/tracks/`: track management and validation.
- `internal/ui/`: Bubble Tea models/components/styles/helpers.
- `internal/validator/`: validation core and individual rule implementations.

## Important Flows

- Config + agent bootstrapping:
  `cmd` calls `internal/config.LoadConfig` and `internal/skills.DetectAgent`; mandatory skill packs are ensured before most commands.
- Project initialization:
  `internal/project` creates `.agentic` directories/files and writes templates from `internal/project/templates/`.
- Task lifecycle:
  `internal/tasks` moves task YAML across backlog/in-progress/done and maintains progress files.
- Run loop:
  `internal/orchestrator` coordinates iterative execution, transitions, and archival/checkpoint behavior.

## Local Conventions / Invariants

- Keep packages cohesive by domain; avoid cross-package cycles.
- Prefer table-driven tests colocated with implementation (`*_test.go` is pervasive here).
- Template-backed initialization depends on `internal/project/templates/*` being kept in sync with source templates in `configs/templates/*`.
- `internal/*` is not importable by external modules; reusable public data types should stay in `pkg/*`.

## How to Run / Test

- `go test ./internal/...`
- `go test ./internal/tasks/...`
- `go test ./internal/skills/...`
- `go test ./internal/orchestrator/...`

## Do / Don’t

- Do update/add tests with behavior changes.
- Do keep interfaces small across package boundaries.
- Don’t add CLI argument parsing here (keep that in `cmd/`).
- Don’t duplicate model types that already exist in `pkg/models`.

## Common Tasks

- Extend task lifecycle behavior in `internal/tasks/`.
- Add new validation rules under `internal/validator/rules/`.
- Add skill packs under `internal/skills/packs/` and registry wiring in `internal/skills/`.
- Update init templates in `configs/templates/` and embedded copies in `internal/project/templates/`.
