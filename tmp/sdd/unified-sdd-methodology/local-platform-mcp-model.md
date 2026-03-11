# Local Platform MCP Model

## Purpose

This document defines how teams can use a local MCP server to query and
validate canonical platform truth without needing hosted infrastructure.

Use this model when:

- the platform repository is the source of truth
- component repositories need fast local access to platform rules and refs
- developers can run a small local process on their machine
- the organization is not ready yet for a centrally hosted MCP service

Related assets:

- [templates/platform-mcp-boilerplate/README.md](templates/platform-mcp-boilerplate/README.md)
- [../platform-truth-mcp-codex-skill/SKILL.md](../platform-truth-mcp-codex-skill/SKILL.md)
- [example/09-local-platform-mcp-usage.md](example/09-local-platform-mcp-usage.md)
- [../platform-truth-mcp-server/README.md](../platform-truth-mcp-server/README.md)

## Core idea

Do not copy platform truth into component repositories. Instead, keep a local
clone of the platform repository and run a read-only MCP server against it.

```text
[platform repo clone]
  platform specs
  contracts
  ADRs
  refs
  JIRA metadata
        |
        v
[local platform MCP]
  read-only query tools
  alignment validation tools
        |
        v
[component repo]
  platform-ref.yaml
  jira-traceability.yaml
  proposal.md
  spec.md
  design.md
  tasks.md
```

## Why this model fits the methodology

This model preserves the source-of-truth separation already defined in the
methodology:

- platform repo = canonical shared truth
- component repo = local OpenSpec truth
- JIRA = workflow and coordination truth

The local MCP server does not change those roles. It only makes the platform
truth easier to consume and validate from a developer workstation.

That makes it a good fit for the current stage of adoption:

- local enough to start now
- structured enough to reduce drift
- small enough to replace later with a hosted MCP service if needed

## Workspace model

Use a side-by-side workspace, not nested clones.

```text
/workspaces
  /platform-repo
  /profile-service
  /auth-service
  /notification-service
```

Recommended rule:

- developers keep `platform-repo` cloned locally
- component repositories do not vendor or mirror the platform repository
- the local MCP server points to the local `platform-repo` path

## Versioning rule

This is the most important rule.

Developers may keep the local platform clone updated. That is good. But local
component work must validate against the pinned platform version in
`platform-ref.yaml` by default.

That means the MCP server should support two modes:

- `pinned`
  - default for component work
  - answers questions from the version or ref declared by the component
- `latest`
  - optional for platform exploration or impact analysis
  - should not be the default validation mode for component delivery

Without this rule, two developers can get different answers from the same local
platform clone.

## What the MCP server should read

The first MCP version should read both platform truth and JIRA-linked metadata.

### Platform truth inputs

- platform baseline
- capability refs
- contract refs
- ADRs
- version markers or tags
- ownership maps
- shared planning notes that are already part of platform truth

### JIRA-linked inputs

Read only the JIRA information that is already represented in the platform or
component artifact chain, or a locally synced export.

Examples:

- platform issue keys
- component epic keys
- story mappings
- issue-link relationships

Do not make JIRA the detailed spec source. The MCP server should help teams
locate the right issue chain, not replace the platform or component specs.

## Recommended first tool surface

Keep the first version small and useful.

### Full target tool surface

- `get_platform_version`
- `list_platform_refs`
- `get_platform_ref`
- `get_contract`
- `get_capability`
- `get_platform_adr`
- `list_platform_changes_since_version`
- `get_jira_mapping`

### Validation tools

- `validate_platform_ref_exists`
- `validate_component_alignment`
- `validate_component_jira_chain`
- `explain_constraints_for_component`
- `detect_platform_drift_from_pinned_version`

### Implemented v1 tools

The current Go server scaffold implements:

- `get_platform_version`
- `list_platform_refs`
- `get_platform_ref`
- `get_jira_mapping`
- `validate_component_alignment`
- `validate_component_jira_chain`
- `detect_platform_drift_from_pinned_version`

## How component teams should use it

Use the MCP server as a local read-only gateway during component work.

### In Specify

- confirm the pinned platform version
- confirm the platform refs used by the component spec
- confirm whether the change is local-only or shared

### In Plan

- read the platform planning decisions for the pinned version
- confirm contract expectations
- validate that local `design.md` and `tasks.md` still align to platform refs

### In Deliver

- validate the active component package against the pinned platform refs
- confirm the JIRA chain is still aligned
- capture PR and verification traceability without re-reading platform docs by hand

## Operational rules

Keep the local MCP server simple.

- read-only only
- no writes to the platform repo
- no writes to the component repo
- no database required for v1
- local filesystem + git metadata + optional local cache only
- small startup cost
- clear failure messages when the local platform clone is stale or missing

## Recommended implementation approach

Both Rust and Go fit this. For v1, prefer the language the team can ship
fastest.

Recommended default:

- use `Go` if the goal is the fastest path to a small cross-platform binary
- use `Rust` if the team already has strong Rust experience and wants stricter
  control of binary size and safety

For a first internal tool, Go is usually the faster path.

## Suggested local workflow

```text
developer updates local platform clone
        |
        v
component repo loads pinned platform version
        |
        v
local MCP answers against pinned version
        |
        v
component team writes local OpenSpec artifacts
        |
        v
validation runs before PR and before archive
```

## Adoption path

Start small.

### Step 1

Use the local MCP server for platform query only.

### Step 2

Add alignment validation for `platform-ref.yaml` and `jira-traceability.yaml`.

### Step 3

Add plan and delivery support for component teams.

### Step 4

Reuse the same tool contract later for a hosted MCP server if the platform team
adds shared infrastructure.
