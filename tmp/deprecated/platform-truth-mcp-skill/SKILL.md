---
name: platform-truth-mcp
summary: Use for defining, setting up, or applying a local read-only MCP gateway that gives component teams access to canonical platform truth and JIRA-linked metadata.
triggers:
  - platform mcp
  - local mcp
  - platform truth mcp
  - component alignment validation
  - local platform repo query
  - platform requirements gateway
---

# Platform Truth MCP

Use this skill when the team wants a local MCP server that reads canonical
platform truth from a local platform clone and helps component teams query and
validate alignment without hosted infrastructure.

## What this skill is for

Use this skill to:

- define the local MCP operating model
- define which platform artifacts the MCP server reads
- define how JIRA mapping is exposed without becoming the spec source of truth
- define the first tool surface for query and validation
- keep the MCP server read-only in v1
- define how component teams use the MCP server during Specify, Plan, and Deliver

## Core operating rule

- platform repo stays the canonical source of truth
- component repos stay the local OpenSpec source of truth
- the local MCP server is a read-only gateway between them
- default validation uses the component's pinned platform version, not latest

## Recommended workflow

```text
[platform repo clone]
        |
        v
[local MCP server]
        |
        v
[component repo uses OpenSpec locally]
```

## Output structure

When using this skill, structure the output as:

1. Platform repo inputs
2. Component repo inputs
3. MCP tool surface
4. Versioning mode
5. JIRA integration mode
6. Safety rules
7. Local workflow
8. Risks
9. Next step

## Rules

- `rules/local-read-only-rules.md`

## References

- `references/workflow.md`
- `references/tool-surface.md`
- `references/configuration.md`
- `../unified-sdd-methodology/local-platform-mcp-model.md`
- `../unified-sdd-methodology/templates/platform-mcp-boilerplate/README.md`
- `../platform-truth-mcp-server/README.md`
