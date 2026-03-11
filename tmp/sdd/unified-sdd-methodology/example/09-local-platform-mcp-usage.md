# Example: Local Platform MCP Usage

Goal:

- show how a developer uses a local platform MCP server while working in a
  component repository

Implementation scaffold:

- `../../platform-truth-mcp-server/README.md`

## Flow

```text
[developer updates platform clone]
        |
        v
[local platform MCP]
        |
        v
[component repo reads pinned platform version]
        |
        v
[OpenSpec local work]
  proposal -> spec -> design -> tasks -> PR
```

## Example prompts

- `Team Lead`
  - "Using the platform truth MCP, validate that `profile-service` is aligned to platform version `2026.03` and list the platform refs that constrain the current change."
- `Architect`
  - "Using the platform truth MCP, retrieve the shared contract and planning constraints for `contracts.customer-profile.v2` and summarize what must appear in `design.md`."
- `Developer`
  - "Using the platform truth MCP, check whether the current component package is drifting from the pinned platform version and list the exact refs to verify before PR creation."

## Example local run

```bash
cd ../../platform-truth-mcp-server
go run ./cmd/platform-truth-mcp serve --config ./examples/demo-platform-mcp-config.yaml
```

## What good looks like

- the local MCP reads from the platform clone
- the component repo still keeps local truth in OpenSpec artifacts
- validation uses the pinned version by default
- JIRA mapping is visible without becoming the spec source of truth
