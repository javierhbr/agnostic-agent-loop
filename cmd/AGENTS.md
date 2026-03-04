# AGENTS.md - cmd

## Purpose

- Define CLI entrypoints and command wiring for the `agentic-agent` executable.

## Key Entrypoints

- `cmd/agentic-agent/main.go`: process entrypoint calling `Execute()`.
- `cmd/agentic-agent/root.go`: root Cobra command, global flags, command registration, startup bootstrapping.
- Command files in `cmd/agentic-agent/*.go`:
  `init`, `task`, `context`, `validate`, `skills`, `token`, `run`, `learnings`, `spec`, `autopilot`, `status`, `track`, `plan`, `simplify`, `openspec`, `prompts`, `platform`, `work`, `start`, `version`.

## File Map

- `cmd/agentic-agent/main.go`
- `cmd/agentic-agent/root.go`
- `cmd/agentic-agent/*.go` command implementations

## Important Flows

- CLI startup:
  `main` -> `Execute()` -> `rootCmd`.
- In `PersistentPreRunE` (`root.go`):
  config is loaded (or defaulted), active agent is detected, and mandatory skill packs may be auto-ensured.
- Commands delegate business logic to `internal/*` packages.

## Local Conventions / Invariants

- Keep flag parsing and user-facing command behavior in `cmd/`; keep core logic in `internal/`.
- Register all commands in `root.go` to expose them via CLI help.
- Respect global flags: `--config`, `--agent`, `--no-interactive`, `--interactive`.

## How to Run / Test

- `go run ./cmd/agentic-agent --help`
- `go build -o build/agentic-agent ./cmd/agentic-agent`
- `go test ./cmd/...`
- `make run ARGS='task list'`

## Do / Don’t

- Do keep command UX and help text clear and stable.
- Do route reusable logic to `internal/*` packages.
- Don’t embed heavy domain logic directly in command handlers.
- Don’t bypass config/agent setup paths defined in `root.go`.

## Common Tasks

- Add a new subcommand file and wire it in `cmd/agentic-agent/root.go`.
- Extend existing command flags and pass options into `internal/*` services.
- Add command-level tests when introducing parsing/UX behavior that is hard to cover through integration tests alone.
