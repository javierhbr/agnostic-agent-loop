# AGENTS.md - configs

## Purpose

- Hold default configuration and initialization templates used by `agentic-agent init`.

## Key Entrypoints

- `configs/agnostic-agent.yaml`: base configuration template.
- `configs/templates/init/`: canonical project-init template set.
- `configs/README.md`: template architecture and sync expectations.

## File Map

- `configs/agnostic-agent.yaml`
- `configs/templates/init/agnostic-agent.yaml`
- `configs/templates/init/tasks/*.yaml`
- `configs/templates/init/context/*.md`
- `configs/templates/init/agent-rules/base.md`
- `configs/templates/init/tool-rules/*.md`
- `configs/templates/init/AGENT_RULES.md`

## Important Flows

- These templates are the source content for project initialization.
- Runtime initialization uses embedded templates in `internal/project/templates/`; template changes here must be reflected there.

## Local Conventions / Invariants

- Keep template paths and filenames stable unless `internal/project` mapping is updated.
- Preserve placeholder/template semantics expected by init logic.
- `scripts/sync-templates.sh` is referenced in docs but not present in this repo (`Unknown - verify` preferred sync command for contributors).

## How to Run / Test

- `go test ./internal/project/...`
- `go test ./cmd/agentic-agent/...`
- Manual validation: `go run ./cmd/agentic-agent init <project-name>`

## Do / Don’t

- Do treat `configs/templates/` as canonical source templates.
- Do copy updates into `internal/project/templates/` before shipping.
- Don’t change template structure without testing `init` end-to-end.
- Don’t add secrets or environment-specific values to default templates.

## Common Tasks

- Update starter task/context files for new project defaults.
- Adjust default agent/workflow settings in `configs/agnostic-agent.yaml`.
- Keep template parity between `configs/templates/` and `internal/project/templates/`.
