# SDD Unified Methodology -- Configuration Files Reference

All configuration files used across the unified SDD methodology, organized by repository scope.

---

## Quick Reference

| File | Scope | Created In | Primary Purpose |
|------|-------|------------|-----------------|
| `platform-baseline.md` | Platform | Platform | Principles, capabilities, contracts, versioning |
| `ownership/component-ownership-<name>.md` | Platform | Platform | Component ownership boundaries |
| `ownership/dependency-map.md` | Platform | Platform | 3-tier impact classification |
| `ownership/glossary.md` | Platform | Platform | Shared terminology |
| `jira-conventions.md` | Platform | Platform | JIRA hierarchy and prefix rules |
| `refs-index.md` | Platform | Platform | Index of all platform refs |
| `adrs/ADR-NNN.md` | Platform | Any phase | Architecture Decision Records |
| `contracts/*.md` | Platform | Platform | Shared API/event contracts |
| `capabilities/*.md` | Platform | Platform | Shared capability definitions |
| `platform-ref.yaml` | Component | Assess | Platform alignment + impact tiers |
| `jira-traceability.yaml` | Component | Assess | JIRA chain from issue to PR |
| `openspec/config.yaml` | Component | Specify | OpenSpec project config |
| `.sdd-spec/changes/<slug>/` | Component | Specify | Change package directory |
| `platform-mcp-config.yaml` | Infra (opt) | Platform | MCP gateway configuration |
| `component-client-config.json` | Infra (opt) | Platform | MCP client configuration |

---

## Platform Repository Files

### 1. platform-baseline.md

Foundational document for the entire platform. Defines principles, capabilities, contracts, versioning strategy, and JIRA conventions.

- **Created by:** Architect, during Platform phase.
- **Updated when:** Platform evolves (new capabilities, version bumps, contract changes).
- **Read by:** All phases. Assess uses it to populate `platform-ref.yaml`.

### 2. ownership/component-ownership-\<name\>.md

One file per component. Declares what the component owns and does not own, which contracts it publishes, and which it consumes.

- **Created by:** Architect, during Platform phase.
- **Updated when:** Ownership boundaries change (new contracts, transferred responsibilities).

### 3. ownership/dependency-map.md

Classifies all cross-component dependencies into three impact tiers:

| Tier | Label | Meaning |
|------|-------|---------|
| 1 | `must_change_together` | Breaking change requires coordinated release |
| 2 | `watch_for_breakage` | May break; requires verification after change |
| 3 | `adapts_independently` | Loosely coupled; independent release safe |

- **Created by:** Architect, during Platform phase.
- **Read in:** Assess phase, to populate `platform-ref.yaml` impact fields.

### 4. ownership/glossary.md

Shared platform terms with explicit "what it is NOT" definitions to prevent ambiguity across teams.

- **Created by:** Architect, during Platform phase.
- **Checked in:** Specify phase, before writing any proposal.

### 5. jira-conventions.md

JIRA hierarchy rules: platform issue prefix, component epic prefix, story prefix, link types.

- **Created by:** Architect or PM, during Platform phase.
- **Read by:** Assess and Plan phases when creating JIRA artifacts.

### 6. refs-index.md

Index of all platform refs (contracts, capabilities, principles). Acts as a lookup table for `platform-ref.yaml` references.

- **Created by:** Architect, during Platform phase.
- **Updated when:** New contracts or capabilities are added.

### 7. adrs/ADR-NNN.md

Architecture Decision Records. Immutable once approved. Each ADR captures context, decision, consequences, and status.

- **Created by:** Architect, during any phase when a significant decision is made.
- **Immutable:** Once status is `accepted`, the record is not modified (superseded by new ADR if needed).

### 8. contracts/\*.md

Shared API and event contracts. Each file defines one contract with schema, versioning, and compatibility guarantees.

- **Created by:** Architect, during Platform phase.
- **Referenced by:** Component `platform-ref.yaml` entries.

### 9. capabilities/\*.md

Shared capability definitions. Each file describes one platform capability, its scope, and its constraints.

- **Created by:** Architect, during Platform phase.
- **Referenced by:** Component `platform-ref.yaml` entries.

---

## Component Repository Files

### 1. platform-ref.yaml

Platform version alignment, impact tiers, and platform refs for a single component.

- **Created in:** Assess phase.
- **Updated in:** Plan and Deliver phases as alignment evolves.
- **Fields:** platform identity, component identity, ownership scope, impact tiers, change package linkage, platform refs with reasons, alignment notes.

```yaml
platform:
  id: acme-platform
  version: "2.4.0"
  baseline_ref: "platform-baseline.md#v2.4"

component:
  name: billing-service
  owner_team: payments

ownership:
  primary_component: billing-service

impact:
  must_change_together:
    - target: invoice-service
      contract: billing-events-v2
      reason: "Shared event schema; breaking change requires joint release"
  watch_for_breakage:
    - target: reporting-service
      contract: billing-read-api
      reason: "Consumes billing queries; field removal may break"
  adapts_independently:
    - target: notification-service
      reason: "Loosely coupled via async events"

change:
  change_package_id: CHG-2024-042
  alignment_type: contract-update

platform_refs:
  contracts:
    - ref: contracts/billing-events-v2.md
      reason: "Publishing updated invoice.created event"
  capabilities:
    - ref: capabilities/multi-currency.md
      reason: "Extending billing to support EUR"

alignment_notes:
  glossary_terms_used:
    - term: settlement
      definition_ref: "ownership/glossary.md#settlement"
```

### 2. jira-traceability.yaml

Full JIRA chain linking platform issues to component epics, stories, and PRs.

- **Created in:** Assess phase.
- **Updated throughout:** Plan, Deliver, and Verify phases as stories progress.

```yaml
platform_issue:
  key: PLAT-1234
  summary: "Add multi-currency support to billing"
  version: "2.4.0"

component_epics:
  - key: BILL-567
    component: billing-service
    change_package_id: CHG-2024-042

stories:
  - key: BILL-568
    task_ref: tasks/currency-model.md
    spec_ref: .sdd-spec/changes/multi-currency/proposal.md
    pr_ref: "https://github.com/acme/billing/pull/142"
    status: merged
  - key: BILL-569
    task_ref: tasks/conversion-api.md
    spec_ref: .sdd-spec/changes/multi-currency/specs/conversion.md
    pr_ref: null
    status: in-progress
```

### 3. openspec/config.yaml

OpenSpec project configuration, typically inherited or seeded from platform conventions.

- **Created in:** Specify phase (or during project init).
- **Contains:** Project name, spec directory paths, template preferences, validation rules.

### 4. .sdd-spec/changes/\<slug\>/

Change package directory. Each slug represents one logical change and contains:

| File | Purpose |
|------|---------|
| `proposal.md` | Problem statement, scope, approach |
| `specs/` | Detailed component specs |
| `design.md` | Architecture and design decisions |
| `tasks.md` | Task breakdown with estimates |
| `archive/` | Previous iterations (if any) |

- **Created in:** Specify phase.
- **Updated in:** Plan and Deliver phases.

---

## Optional Infrastructure Files

### 1. platform-mcp-config.yaml

MCP gateway configuration for platform-level tooling orchestration.

- **Created by:** Platform team during Platform phase.
- **Optional:** Only needed when using MCP-based agent coordination.

### 2. component-client-config.json

MCP client configuration for component-level agent tooling.

- **Created by:** Component team during setup.
- **Optional:** Only needed when using MCP-based agent coordination.

---

## Lifecycle by Phase

Which files are created (C), read (R), or updated (U) in each methodology phase.

| File | Platform | Assess | Specify | Plan | Deliver | Verify |
|------|----------|--------|---------|------|---------|--------|
| `platform-baseline.md` | C | R | R | R | R | R |
| `ownership/component-ownership-*.md` | C | R | R | | | |
| `ownership/dependency-map.md` | C | R | | | | |
| `ownership/glossary.md` | C | | R | | | |
| `jira-conventions.md` | C | R | | R | | |
| `refs-index.md` | C | R | R | | | |
| `adrs/ADR-NNN.md` | C | | C | C | | |
| `contracts/*.md` | C | R | R | | | R |
| `capabilities/*.md` | C | R | R | | | |
| `platform-ref.yaml` | | C | R | U | U | R |
| `jira-traceability.yaml` | | C | | U | U | R |
| `openspec/config.yaml` | | | C | R | R | |
| `.sdd-spec/changes/<slug>/` | | | C | U | U | R |
