# Configuration Files Reference

This document lists every configuration file used across platform and component
repositories in the unified SDD methodology. Use it to answer: "what files do I
need, where do they go, and when do I create or update them?"

For the alignment model behind these files, see
[canonical-platform-truth-and-component-alignment.md](canonical-platform-truth-and-component-alignment.md).

---

## Quick Reference

| File | Repo | Purpose | Created | Updated |
|------|------|---------|---------|---------|
| `platform-baseline.md` | Platform | Durable platform truth | Platform | When platform evolves |
| `ownership/component-ownership-<name>.md` | Platform | What each component owns | Platform | When ownership changes |
| `ownership/dependency-map.md` | Platform | 3-tier impact classification | Platform | When dependencies change |
| `ownership/glossary.md` | Platform | Shared terms | Platform | Specify (incremental) |
| `jira-conventions.md` | Platform | JIRA hierarchy rules | Platform | Rarely |
| `refs-index.md` | Platform | Index of all platform refs | Platform | When refs added |
| `adrs/ADR-NNN.md` | Platform | Architecture decisions | Any phase | Never (immutable) |
| `contracts/*.md` | Platform | Shared API/event contracts | Platform/Specify | When contracts change |
| `capabilities/*.md` | Platform | Shared capability definitions | Platform | When capabilities change |
| `platform-ref.yaml` | Component | Platform version alignment | Assess | Plan, Deliver |
| `jira-traceability.yaml` | Component | JIRA chain to stories | Assess | Plan, Deliver |
| `openspec/config.yaml` | Component | OpenSpec project config | Platform (once) | Rarely |
| `.sdd-spec/changes/<slug>/` | Component | Change package artifacts | Specify | Plan, Deliver |
| `platform-mcp-config.yaml` | Local infra | MCP gateway config | Platform (once) | Rarely |
| `component-client-config.json` | Local infra | MCP client config | Platform (once) | Rarely |

---

## Platform Repository Files

These files live in the platform master repository. They define shared truth
that component repositories align to.

Template: [templates/platform-template/](templates/platform-template/)

### 1. `platform-baseline.md`

The foundational document. Defines platform principles, shared capabilities,
contracts, versioning rules, and JIRA conventions.

- **Created:** Platform phase by the Architect
- **Updated:** When the platform evolves (requires ADR if contradicting existing specs)
- **Used by:** Every phase — Assess reads it to classify changes, Specify reads
  it to confirm refs, Plan reads it to constrain design

Sample: [templates/platform-template/sample-platform-baseline.md](templates/platform-template/sample-platform-baseline.md)

### 2. `ownership/component-ownership-<name>.md`

One file per component. Records what the component owns, what it does NOT own,
and its published and consumed contracts.

- **Created:** Platform phase (write once)
- **Read at:** Every Assess step (rule O-1: confirm ownership before opening any JIRA epic)
- **Updated:** When ownership boundaries change during delivery

```markdown
# Component Ownership: profile-service

## Owns
- Customer profile CRUD
- Email validation rules
- Profile event publishing

## Does NOT Own
- Authentication (auth-service)
- Email delivery (notification-service)

## Contracts
- Publishes: `customer-profile.v2` (event)
- Consumes: `auth-token.v1` (API)
```

### 3. `ownership/dependency-map.md`

One platform file. Records which components must change together (tier 1),
which need monitoring (tier 2), and which adapt independently (tier 3).

- **Created:** Platform phase
- **Read at:** Assess (populate `platform-ref.yaml` impact tiers), Plan (design constraints)
- **Updated:** When platform integrations change

```markdown
# Dependency Map

## Tier 1 — must_change_together
- profile-service ↔ auth-service (shared session contract)

## Tier 2 — watch_for_breakage
- profile-service → notification-service (profile events)

## Tier 3 — adapts_independently
- analytics-service
- reporting-service
```

### 4. `ownership/glossary.md`

Shared glossary of platform terms with "what it is NOT" clauses. Prevents spec
ambiguity before Specify begins.

- **Created:** Platform phase (seeded)
- **Read at:** Specify (rule O-2: all proposal terms must be in the glossary)
- **Updated:** Incrementally during Specify phases

```markdown
# Glossary

## Customer
A registered user with a verified email. NOT a guest or anonymous visitor.

## Eligibility
Whether a customer qualifies for a specific action based on business rules.
NOT a technical permission or auth check.
```

### 5. `jira-conventions.md`

Defines the JIRA hierarchy and issue-link rules for the platform.

- **Created:** Platform phase
- **Updated:** Rarely — only when JIRA structure changes

### 6. `refs-index.md`

Index of all platform refs (principles, capabilities, contracts) with stable IDs.

- **Created:** Platform phase
- **Updated:** When new refs are added

### 7. `adrs/ADR-NNN.md`

Architecture Decision Records. Immutable once accepted.

- **Created:** Any phase when a significant decision is made
- **Template:** [templates/platform-template/adrs/ADR-000-template.md](templates/platform-template/adrs/ADR-000-template.md)

### 8. `contracts/*.md` and `capabilities/*.md`

Shared API/event contract definitions and capability definitions.

- **Created:** Platform phase for initial set
- **Updated:** When contracts or capabilities change (requires versioning per Constitution)
- **Templates:** [templates/platform-template/contracts/](templates/platform-template/contracts/), [templates/platform-template/capabilities/](templates/platform-template/capabilities/)

---

## Component Repository Files

These files live in each component repository. They record platform alignment
and contain the local implementation artifacts.

Template: [templates/component-boilerplate/](templates/component-boilerplate/)

### 1. `platform-ref.yaml`

Records which platform version, refs, and impact tiers this component follows.
This is the key alignment file between platform truth and component work.

- **Created:** Assess phase by the Team Lead
- **Updated:** Plan (refined), Deliver (finalized)
- **Owner:** Team Lead

Key sections: `platform` (id, repo, version), `component` (name, owner_team),
`change` (alignment_type, requires_platform_change), `platform_refs` (principles,
capabilities, contracts with reasons), `alignment_notes` (summary, risks).

Full template with examples: [templates/platform-ref.yaml](templates/platform-ref.yaml)

### 2. `jira-traceability.yaml`

Records the JIRA chain from platform issue to component epic to stories.

- **Created:** Assess phase by the Team Lead
- **Updated:** Plan (stories added), Deliver (PR refs and status updates)
- **Owner:** Team Lead

Key sections: `platform_issue` (key, summary, version), `component_epic`
(key, component, change_package_id), `stories` (key, task_ref, spec_ref, pr_ref,
status), `delivery_rules` (one_story_per_slice, pr_must_reference_story).

Full template with examples: [templates/jira-traceability.yaml](templates/jira-traceability.yaml)

### 3. `openspec/config.yaml`

OpenSpec project configuration. Defines the component name, platform version
context, and local spec rules.

- **Created:** Platform phase (once per component)
- **Updated:** Rarely — only when platform version or rules change

```yaml
project:
  name: "profile-service"
  description: "Manages customer profile data"

context:
  platform_version: "2026.03"
  platform_refs:
    - "capabilities.customer-identity"
    - "contracts.customer-profile.v2"

rules:
  - "keep local specs aligned to platform refs"
  - "do not weaken shared contracts silently"
  - "use reviewable slices and traceable tasks"
```

### 4. `.sdd-spec/changes/<slug>/` — Change Package

Each approved change gets a directory with these artifacts:

```
.sdd-spec/changes/chg-profile-email-validation/
├── proposal.md      ← what and why (Specify phase)
├── specs/           ← delta specs with MUST/SHALL language (Specify phase)
├── design.md        ← architecture decisions and tradeoffs (Plan phase)
└── tasks.md         ← implementation tasks (Plan phase)
```

- **Created:** Specify phase (proposal + specs), Plan phase (design + tasks)
- **Archived:** Deliver phase — merged into `.sdd-spec/specs/` via `/openspec-archive`

---

## Optional Infrastructure Files

These files support the local platform MCP gateway — a read-only tool for
querying platform truth from a developer workstation without hosted infrastructure.

For full detail, see [local-platform-mcp-model.md](local-platform-mcp-model.md).

Template: [templates/platform-mcp-boilerplate/](templates/platform-mcp-boilerplate/)

### 1. `platform-mcp-config.yaml`

Configures the local MCP server: platform repo path, ref mode, spec/contract
directories, cache, and safety settings.

```yaml
server:
  name: platform-truth-mcp
  mode: local-read-only

platform_repo:
  path: /absolute/path/to/platform-repo
  default_ref_mode: pinned
  allow_latest_queries: true

inputs:
  specs_dir: docs/platform
  contracts_dir: docs/contracts
  adrs_dir: docs/adrs

safety:
  read_only: true
  allow_repo_writes: false
```

### 2. `component-client-config.sample.json`

Sample MCP client configuration for component repos to connect to the local
gateway.

---

## Lifecycle by Phase

| Phase | Creates | Updates |
|-------|---------|---------|
| **Platform** | `platform-baseline.md`, ownership files, `glossary.md`, `jira-conventions.md`, `refs-index.md`, `openspec/config.yaml`, MCP configs | — |
| **Assess** | `platform-ref.yaml`, `jira-traceability.yaml` | — |
| **Specify** | `.sdd-spec/changes/<slug>/proposal.md`, delta specs | `glossary.md` (new terms) |
| **Plan** | `.sdd-spec/changes/<slug>/design.md`, `tasks.md`, ADRs | `platform-ref.yaml`, `jira-traceability.yaml` |
| **Deliver** | PR evidence, `verify.md` | `jira-traceability.yaml` (PR refs, status), ownership files (if boundaries changed) |
