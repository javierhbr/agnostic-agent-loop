# Platform MCP Boilerplate

Use this boilerplate as a reference for how platform truth is organized so that
`platform-spec` can index and search it effectively.

See also:

- `../../local-platform-mcp-model.md`
- `../../../platform-spec/SKILL.md`
- `~/.agentic/skills/` — agentic skills directory

## Files

- `platform-mcp-config.yaml`
- `component-client-config.sample.json`
- `tool-catalog.md`
- `component-consumption-checklist.md`

## Intended use

- keep the platform repository as the source of truth
- run the MCP server locally against a local platform clone
- let component repositories query and validate platform truth without copying it
