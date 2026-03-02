# AGENTS.md - pkg

## Purpose

- Provide shared, import-safe domain models used by CLI and internal packages.

## Key Entrypoints

- `pkg/models/config.go`: configuration schema (`models.Config`) used by config loading/defaulting.
- `pkg/models/task.go`: task lifecycle models (`Task`, `TaskStatus`, `SubTask`).
- `pkg/models/agent.go`, `pkg/models/context.go`, `pkg/models/track.go`: additional domain types used across the app.

## File Map

- `pkg/models/config.go`
- `pkg/models/task.go`
- `pkg/models/agent.go`
- `pkg/models/context.go`
- `pkg/models/track.go`
- `pkg/models/task_test.go`

## Important Flows

- CLI config loading (`internal/config`) unmarshals YAML into `pkg/models` types, then applies defaults.
- Task workflows (`internal/tasks`, `cmd/agentic-agent/task.go`) read/write YAML using `pkg/models.Task` and status constants.

## Local Conventions / Invariants

- Keep model structs serialization-safe (`yaml` tags aligned with persisted files).
- Additive model changes should preserve backward compatibility with existing `.agentic/*.yaml` files when possible.
- Runtime-only fields should be explicitly marked (example: `ActiveAgent` has `yaml:"-"`).

## How to Run / Test

- `go test ./pkg/...`
- `go test ./pkg/models/...`

## Do / Don’t

- Do keep this layer free of CLI/UI concerns.
- Do update tests when changing model fields or status enums.
- Don’t add business logic that belongs in `internal/*`.
- Don’t break YAML field names without migration planning.

## Common Tasks

- Add a new persisted field to task/config models.
- Add/update tests in `pkg/models/task_test.go` for schema-sensitive behavior.
