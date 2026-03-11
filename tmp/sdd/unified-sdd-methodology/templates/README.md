# Templates

This directory contains reusable starting points for the unified SDD methodology.

Use these when you want a clean structure instead of a filled example.

## Available templates

- `platform-template/`
- `component-boilerplate/`
- `platform-mcp-boilerplate/`
- `platform-ref.yaml`
- `jira-traceability.yaml`

## When to use each one

- use `platform-template/` to start or refresh the master platform repository
- use `component-boilerplate/` to start a component repository structure that aligns to the platform truth
- use `platform-mcp-boilerplate/` to define a small local read-only MCP server that exposes platform truth to component teams
- use `../example/` when you want a filled worked example instead of placeholders

## Quick usage docs

- `platform-template/how-to-use.md`
- `component-boilerplate/how-to-use.md`
- `platform-mcp-boilerplate/README.md`

## Small filled samples

- `platform-template/sample-platform-baseline.md`
- `component-boilerplate/sample-component-package.md`

## Relationship to the example package

```text
[templates/]
  reusable starting points
        |
        v
[example/]
  filled example with real sample IDs and artifacts
```
