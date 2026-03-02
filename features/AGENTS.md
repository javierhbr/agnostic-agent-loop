# AGENTS.md - features

## Purpose

- Store BDD feature specifications (Gherkin) for CLI behavior and tutorial workflows.

## Key Entrypoints

- `features/init/project_initialization.feature`
- `features/tasks/task_lifecycle.feature`
- `features/tasks/error_handling.feature`
- `features/context/context_generation.feature`
- `features/skills/agent_detection.feature`
- `features/skills/gemini_skills.feature`
- `features/workflows/*.feature`

## File Map

- `features/context/`: context generation scenarios.
- `features/init/`: project bootstrap scenarios.
- `features/skills/`: agent detection and skills generation scenarios.
- `features/tasks/`: task lifecycle and error-handling scenarios.
- `features/workflows/`: beginner/intermediate/advanced tutorial flows.
- `features/validation/`: currently empty.

## Important Flows

- BDD runner in `test/bdd/features_test.go` loads these files via `Paths: []string{"../../features"}`.
- Step implementations live in `test/bdd/steps/*`; keep feature phrasing aligned with existing steps or add new step defs.

## Local Conventions / Invariants

- Keep scenarios deterministic and self-contained (most use clean temporary environments).
- Prefer behavior-focused language over implementation details.
- Use tags (for example tutorial tags) only when needed for suite filtering.

## How to Run / Test

- `go test ./test/bdd -v`
- `make test-bdd`
- `make test-bdd-verbose`

## Do / Don’t

- Do add/adjust step definitions when introducing new Given/When/Then phrases.
- Do keep feature folder placement consistent with domain.
- Don’t reference local machine paths or environment assumptions in scenarios.
- Don’t duplicate near-identical scenarios without additional behavior coverage.

## Common Tasks

- Add a new CLI behavior scenario in the relevant `features/<domain>/` folder.
- Expand workflow tutorials under `features/workflows/`.
- Add missing validation scenarios under `features/validation/` (currently empty).
