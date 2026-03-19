# How to Use the platform-spec Tool

## What is platform-spec?

platform-spec is a CLI tool and agentic skill for querying platform truth locally. It provides read-only access to platform specifications without requiring hosted infrastructure.

Key characteristics:

- Operates on a local clone of the platform repo
- Builds and queries a SQLite FTS5 (full-text search) index
- v1 implementation uses Go stdio with newline-delimited JSON-RPC
- Designed for offline, low-latency lookups against platform truth

## Setup

### 1. Clone the platform repo

Clone the platform repo alongside your component repo. The two repos should sit side-by-side, not nested:

```
workspace/
  platform/          # platform truth repo
  my-component/      # your component repo
```

### 2. Build the index

Run the indexer to create the SQLite FTS5 index:

```bash
platform-spec index
```

This scans all platform specs and builds a searchable index. Re-run after pulling new platform changes.

### 3. Configure

Point platform-spec to your platform clone path. This is typically auto-detected from the side-by-side layout, but can be set explicitly if your workspace differs.

## Core Commands

| Command | Description |
|---------|-------------|
| `platform-spec search "<query>"` | Full-text search across all platform specs |
| `platform-spec read <file-path>` | Read a specific platform file by path |
| `platform-spec list <directory>` | List files in a platform directory |
| `platform-spec related <term>` | Find specs related to a given term |
| `platform-spec json <query>` | JSON output for programmatic consumption |

### Examples

```bash
# Search for authentication-related specs
platform-spec search "authentication"

# Read the glossary
platform-spec read glossary.md

# List all component ownership files
platform-spec list component-ownership/

# Find specs related to "billing"
platform-spec related billing

# Get JSON output for scripting
platform-spec json "api gateway" | jq '.results[]'
```

## Usage by Phase

### Platform phase

Register new components in the platform repo and rebuild the index so they become searchable.

### Assess phase

Read component-ownership files to determine the primary owner of a component. Read dependency-map.md to understand impact tiers before proposing changes.

### Specify phase

Check glossary.md terms before writing a proposal. Use consistent terminology that matches platform truth.

### Plan phase

Read impact tiers from platform-ref.yaml to constrain your design. Higher-impact changes require broader review and more conservative rollout plans.

### Deliver phase

Detect drift between your component and platform truth. If your component has diverged from what the platform specifies, platform-spec surfaces the mismatch.

## Workspace Model

The platform repo and component repos sit side-by-side in the same workspace directory. Component repos pin to a specific platform version using a `platform-ref.yaml` file at the component root.

platform-spec validates queries against the pinned version by default. This ensures that the truth you query matches the truth your component was built against.

## Versioning Rule

Pin the platform version in your component's `platform-ref.yaml`:

```yaml
platform:
  version: "2026-03-15"
  repo: "../platform"
```

platform-spec resolves all queries against the pinned version. If the platform repo has evolved beyond your pinned version, platform-spec warns about the version mismatch so you can decide whether to upgrade.

## Adoption Path

1. **Clone** the platform repo locally, side-by-side with your component repo.
2. **Build the index** by running `platform-spec index`.
3. **Start with search** in the Assess phase -- use `platform-spec search` to find relevant specs before making design decisions.
4. **Integrate with agentic skills** for automated lookups. The platform-spec skill can be invoked by agents during task execution to ground decisions in platform truth.

## Troubleshooting

- **Index out of date:** Re-run `platform-spec index` after pulling platform repo changes.
- **Version mismatch warning:** Update `platform-ref.yaml` to the latest platform version, or acknowledge the drift if intentional.
- **Path not found:** Verify the platform repo is cloned side-by-side and the path in your configuration is correct.
