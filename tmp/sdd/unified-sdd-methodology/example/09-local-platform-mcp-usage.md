# Example: Local Platform Access with platform-spec

Goal:

- show how a developer uses `platform-spec` to query platform truth while
  working in a component repository

Related skill:

- `../../platform-spec/SKILL.md`

Agentic skills:

- `~/.agentic/skills/` — locally installed skills available across sessions

## Flow

```text
[developer registers local platform clone]
  platform-spec repo add /path/to/platform-repo platform
        |
        v
[platform-spec indexes .md/.mdx/.txt files]
  SQLite FTS5 index built locally
        |
        v
[component repo reads pinned platform version]
  platform-ref.yaml pins the version
        |
        v
[OpenSpec local work]
  proposal -> spec -> design -> tasks -> PR
```

## Example prompts

- `Team Lead`
  - "Using platform-spec, search for the component ownership file for `profile-service` and confirm which team owns this request before opening a JIRA epic."
  - "Using platform-spec, read the dependency map and populate the impact tier fields in platform-ref.yaml for this change."
- `Architect`
  - "Using platform-spec, search for all tier 1 entries that affect `profile-service` and list them as named coordination requirements in design.md."
  - "Using platform-spec, read the shared contract for `contracts.customer-profile.v2` and summarize what must appear in design.md."
- `Developer`
  - "Using platform-spec, check whether the component package has any drift from the pinned platform version by searching for the ownership and contract refs used in proposal.md."

## Example local setup

```bash
# Register the local platform clone once
platform-spec repo add /path/to/platform-repo platform

# Verify it indexed
platform-spec list platform

# Search during Assess
platform-spec read ownership/component-ownership-profile-service
platform-spec read ownership/dependency-map

# Search during Specify
platform-spec read ownership/glossary
platform-spec search "eligibility" platform

# Search during Plan
platform-spec json read ownership/dependency-map

# Search during Deliver
platform-spec search "contract" platform
platform-spec related component-ownership-profile-service
```

## What good looks like

- `platform-spec` reads from the local platform clone
- the component repo still keeps local truth in OpenSpec artifacts
- searches are scoped to the pinned platform context by default
- ownership and glossary lookups happen before JIRA epics and before writing proposals
- JIRA mapping stays as coordination tracking, not the spec source of truth
