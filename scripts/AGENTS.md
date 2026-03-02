# AGENTS.md - scripts

## Purpose

- Hold project-level helper scripts for local verification and coverage reporting.

## Key Entrypoints

- `scripts/verify_all.sh`: end-to-end smoke script that builds the CLI and runs a full task lifecycle in a temp workspace.
- `scripts/coverage-report.sh`: coverage report generator with threshold enforcement and optional browser open.

## File Map

- `scripts/verify_all.sh`
- `scripts/coverage-report.sh`

## Important Flows

- `verify_all.sh` builds `./cmd/agentic-agent`, initializes a project, creates/claims/completes workflow artifacts, and cleans up.
- `coverage-report.sh` runs `go test ./...` with coverage profile output, prints summaries, and can open `coverage/coverage.html`.

## Local Conventions / Invariants

- Scripts are Bash and use `set -e`; failing commands should terminate early.
- Keep scripts runnable from repo root (paths assume root working directory).
- `coverage-report.sh` enforces a minimum total coverage threshold (currently `50%`).

## How to Run / Test

- `bash scripts/verify_all.sh`
- `bash scripts/coverage-report.sh`
- `bash scripts/coverage-report.sh --open`

## Do / Don’t

- Do keep scripts idempotent and cleanup-aware for temporary files.
- Do update script messages when behavior changes.
- Don’t hardcode machine-specific absolute paths.
- Don’t assume GUI browser availability when using `--open`.

## Common Tasks

- Run a quick CLI smoke test: `bash scripts/verify_all.sh`
- Generate local coverage reports: `bash scripts/coverage-report.sh`
