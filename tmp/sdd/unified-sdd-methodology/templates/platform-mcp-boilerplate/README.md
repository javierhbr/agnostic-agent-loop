# Platform MCP Boilerplate

Use this boilerplate to define a small local MCP server that gives component
teams read-only access to canonical platform truth.

This is not a full server implementation. It is a starter shape for the first
local version.

See also:

- `../../local-platform-mcp-model.md`
- `../../../platform-truth-mcp-codex-skill/SKILL.md`
- `../../../platform-truth-mcp-server/README.md`

## Files

- `platform-mcp-config.yaml`
- `component-client-config.sample.json`
- `tool-catalog.md`
- `component-consumption-checklist.md`

## Intended use

- keep the platform repository as the source of truth
- run the MCP server locally against a local platform clone
- let component repositories query and validate platform truth without copying it
