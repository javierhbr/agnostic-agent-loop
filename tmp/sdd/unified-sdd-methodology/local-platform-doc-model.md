# Local Platform Access Model

## Purpose

This document defines how teams query and validate canonical platform truth
locally without needing hosted infrastructure.

Use this model when:

- the platform repository is the source of truth
- component repositories need fast local access to platform rules and refs
- developers keep a local clone of the platform repository
- the organization is not ready yet for a centrally hosted service

Related assets:

- [../platform-spec/SKILL.md](../platform-spec/SKILL.md)
- [canonical-platform-truth-and-component-alignment.md](canonical-platform-truth-and-component-alignment.md)
- [example/09-local-platform-mcp-usage.md](example/09-local-platform-mcp-usage.md)

## Core idea

Do not copy platform truth into component repositories. Instead, keep a local
clone of the platform repository and use `platform-spec` to index and search it.

```text
[platform repo clone]
  platform specs (.md, .mdx, .txt)
  contracts
  ADRs
  refs
  ownership artifacts
        |
        v
[platform-spec index]
  SQLite FTS5 full-text index
  built from local clone
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

## The platform-spec tool

`platform-spec` is a CLI tool and a Claude skill. It indexes `.md`, `.mdx`,
and `.txt` files across one or more local repositories and provides instant
full-text search via SQLite FTS5. Single bash script, zero dependencies beyond
`sqlite3`.

Skill location: `../platform-spec/SKILL.md`

### Setup

Register your local platform clone as a repo:

```bash
platform-spec repo add /path/to/platform-repo platform
platform-spec search "component ownership"   # ready immediately
```

The index auto-rebuilds when repos are added and auto-detects file changes on
next search.

### Core commands

```bash
# Search
platform-spec search "payment eligibility"
platform-spec search "tier 1 OR tier 2" platform
platform-spec search "contract NOT deprecated"

# Read a spec
platform-spec read component-ownership
platform-spec read dependency-map

# Browse
platform-spec list platform
platform-spec related component-ownership

# JSON output for agent pipelines
platform-spec json search "impact tier"
platform-spec json read glossary
```

## Why this model fits the methodology

This model preserves the source-of-truth separation defined in the methodology:

- platform repo = canonical shared truth
- component repo = local OpenSpec truth
- JIRA = workflow and coordination truth

`platform-spec` does not change those roles. It makes the platform truth
searchable and queryable from a developer workstation with no server process
required.

## Workspace model

Use a side-by-side workspace, not nested clones.

```text
/workspaces
  /platform-repo       ← register with platform-spec
  /profile-service
  /auth-service
  /notification-service
```

Recommended rule:

- developers keep `platform-repo` cloned locally
- component repositories do not vendor or mirror the platform repository
- `platform-spec` points to the local `platform-repo` path

## Versioning rule

Local component work must validate against the pinned platform version in
`platform-ref.yaml` by default, even when the local clone is more recent.

Use `platform-spec` with a scoped search to stay aligned with the pinned
version:

```bash
# Search only within the pinned platform context
platform-spec search "contract" platform
platform-spec read ownership/dependency-map
```

Without this discipline, two developers can get different answers from the same
local platform clone.

## How component teams use platform-spec by phase

### Platform

- register the platform repo: `platform-spec repo add /path/to/platform platform`
- verify ownership artifacts are indexed: `platform-spec list platform`

### Assess

- confirm component ownership before opening a JIRA epic (rule O-1):

```bash
platform-spec search "component-ownership auth-service"
platform-spec read ownership/component-ownership-auth-service
```

- read the dependency map and classify impact tiers:

```bash
platform-spec read ownership/dependency-map
```

### Specify

- look up glossary terms before writing `proposal.md` (rule O-2):

```bash
platform-spec search "eligibility" platform
platform-spec read ownership/glossary
```

### Plan

- read platform impact tiers before designing:

```bash
platform-spec json read ownership/dependency-map
```

- validate that `design.md` and `tasks.md` align to platform refs

### Deliver

- detect drift against the pinned platform version before merge:

```bash
platform-spec search "contract" platform
platform-spec related component-ownership-auth-service
```

## Agentic directory integration

When working inside a `.agentic/` project directory, `platform-spec` integrates
with the agentic workflow:

```text
/.agentic/
  skills/
    openspec/          ← local OpenSpec task execution
    ...
  tasks/               ← current task state
```

Register the platform repo once per workspace:

```bash
platform-spec repo add /path/to/platform-repo platform
```

From that point, any agentic skill or Claude session in the workspace can query
platform truth during Assess, Specify, Plan, and Deliver without manual file
lookups.

## Operational rules

Keep the local setup simple.

- read-only queries only — `platform-spec` never writes to the platform repo
- no server process required — single SQLite file, single bash script
- no external dependencies beyond `sqlite3` (pre-installed on macOS and Linux)
- small startup cost — index rebuilds automatically on file change
- clear failure messages when the local platform clone is stale or missing

## Adoption path

Start small.

### Step 1

Register the platform repo and run your first search.

```bash
platform-spec repo add /path/to/platform-repo platform
platform-spec search "ownership"
```

### Step 2

Use `platform-spec` during Assess to confirm ownership and classify impact
tiers before opening JIRA epics.

### Step 3

Use `platform-spec` during Specify to verify all proposal terms exist in the
shared glossary.

### Step 4

Use `platform-spec` during Plan and Deliver to validate alignment and detect
drift before PRs and before archive.
