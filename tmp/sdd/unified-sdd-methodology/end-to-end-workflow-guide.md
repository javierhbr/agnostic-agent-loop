# End-to-End Workflow Guide

This document walks through the complete unified SDD methodology step by step.
For each step it shows the role, the tool, the skill, and a concrete sample
from the running example: **validated customer email updates**.

---

## Running example

A product team requests that customers be able to update their email address
with proper validation and clear failure feedback. This change affects multiple
services and a shared platform contract.

Shared IDs used throughout:

| ID | Meaning |
|---|---|
| `PLAT-123` | Platform issue — shared change |
| `PROF-456` | Component epic — profile-service |
| `AUTH-234` | Component epic — auth-service |
| `PROF-789` | Story — email validation task 1 |
| `PROF-790` | Story — email validation task 2 |
| `PROF-791` | Story — verification and release |
| `chg-profile-email-validation` | OpenSpec change package ID |
| `contracts.customer-profile.v2` | Shared platform contract |

---

## Methodology overview

```text
  ITERATION 1                               ITERATION 2
  ┌─────────────────────────────────┐       ┌───────────────────────────────────────────────┐
  │                                 │       │                                               │
  │  [Platform]──>[Assess]──>[Specify]──────>[Plan]──>[Deliver]                            │
  │                                 │       │              │                                │
  └─────────────────────────────────┘       │         ┌────┴──────────────────────┐        │
                                            │         │ Build > PR > Review       │        │
                                            │         │ Verify > Deploy > Archive │        │
                                            │         └───────────────────────────┘        │
                                            └───────────────────────────────────────────────┘
```

**The hard boundary rule:**

```text
  ┌──────────────────────────────────────────────────────────────┐
  │  PLATFORM-SIDE WORK  (cross-cutting)                         │
  │  Speckit  +  OpenSpec  +  BMAD  — all three tools allowed   │
  └──────────────────────────────────────────────────────────────┘
                          │
                          │  handoff
                          ▼
  ┌──────────────────────────────────────────────────────────────┐
  │  COMPONENT REPO  (local work)                                │
  │  OpenSpec ONLY — no Speckit, no BMAD                        │
  └──────────────────────────────────────────────────────────────┘
```

---

## Phase 1: Platform

**Goal:** Create shared context and durable rules before any change-level work starts.

**Owner:** Architect | **Support:** Product, Team Lead

```text
  Architect + Product + Team Lead
              │
              ▼
  ┌───────────────────────────────────────────────────────────┐
  │                     PLATFORM PHASE                        │
  │                                                           │
  │  Step 1.1  BMAD      ──► Brownfield context review       │
  │               │                                          │
  │               ▼                                          │
  │  Step 1.2  Speckit   ──► Constitution + principles       │
  │               │                                          │
  │               ▼                                          │
  │  Step 1.3  OpenSpec  ──► config.yaml (durable context)   │
  │               │                                          │
  │               ▼                                          │
  │  Step 1.4  OpenSpec  ──► Versioning + JIRA conventions   │
  │               │                                          │
  │               ▼                                          │
  │  Step 1.5  Architect ──► Ownership artifacts             │
  │                          component-ownership-<name>.md   │
  │                          dependency-map.md               │
  │                          glossary.md                     │
  └───────────────────────────────────────────────────────────┘
              │
              ▼
        Platform baseline ready
        (constitution · config · role map · versioning ·
         ownership boundaries · dependency map · glossary)
```

---

### Step 1.1 — Document the current platform state

**Role:** Architect
**Tool:** BMAD
**Skill:** `sdd-bmad` (brownfield + project-context playbook)
**When:** Before any rules are written. Use when the platform already exists.

```text
  Existing platform repos
  (conventions · contracts · team habits · partial docs)
              │
              ▼  sdd-bmad  brownfield inspection
  ┌───────────────────────────────────┐
  │  BMAD brownfield output           │
  │  · integration points             │
  │  · existing contracts             │
  │  · architectural constraints      │
  │  · team conventions               │
  └───────────────────────────────────┘
```

**What happens:**

The Architect uses BMAD's brownfield inspection mode to document the current
state of the platform: existing conventions, architectural constraints, repo
boundaries, and integration points across services.

**Sample prompt:**

> "Using the BMAD skill, review the current platform as a brownfield system.
> Identify the architectural constraints for customer identity flows, the
> integration points across profile-service, auth-service, and
> notification-service, and the existing conventions around contracts and API
> versioning."

**Sample output:**

```
Brownfield context — customer-platform

Existing conventions:
- API versioning enforced via contract semver (contracts.customer-profile.v2)
- Observability: all service writes emit structured logs and metrics
- Auth and profile share a customer identity root

Integration points:
- profile-service owns the profile write path
- auth-service owns duplicate-email enforcement
- notification-service triggers on profile update events

Constraints:
- contracts.customer-profile.v2 must not be bumped without a shared migration plan
- API changes require backwards compatibility for at least one release cycle
```

---

### Step 1.2 — Draft the platform constitution

**Role:** Architect
**Tool:** Speckit
**Skill:** `sdd-speckit` (constitution command + `rules/constitution-rules.md`)
**When:** After the brownfield review. First formal output of Platform phase.

```text
  BMAD brownfield context
              │
              ▼  sdd-speckit  constitution
  ┌─────────────────────────────────────────┐
  │  Constitution                           │
  │  · Validation rules   (MUST / SHALL)    │
  │  · Contract rules     (MUST / SHALL)    │
  │  · Observability rules                  │
  │  · API versioning rules                 │
  │  · Security / quality bar               │
  └─────────────────────────────────────────┘
              │
              ▼
        Explicit · testable · durable
```

**What happens:**

The Architect uses Speckit to turn the brownfield observations into explicit,
testable platform principles. These become the governing rules every future
change must follow.

**Sample prompt:**

> "Using the Speckit skill, convert the customer identity platform observations
> into explicit, testable rules for validation, contracts, observability,
> security, and API versioning."

**Sample output — draft constitution (excerpt):**

```markdown
# Platform Constitution — customer-platform v1

## Validation
- All profile write operations MUST validate input before persistence.
- Invalid inputs MUST return structured error responses with an error code.

## Contracts
- Shared contracts MUST NOT be modified without a versioned migration plan.
- Components consuming a shared contract MUST pin the contract version in
  platform-ref.yaml.

## Observability
- All validation failures MUST emit a structured log entry with service name,
  operation name, and failure reason.
- All service writes MUST emit a metrics counter with success/failure labels.

## API versioning
- API changes MUST maintain backwards compatibility for one release cycle.
- Breaking changes MUST go through an ADR and a shared migration plan.
```

---

### Step 1.3 — Encode durable project config

**Role:** Architect
**Tool:** OpenSpec
**Skill:** `sdd-openspec` (config encoding + `rules/project-config-template.yaml`)
**When:** After the constitution is approved.

```text
  Constitution + versioning decisions
              │
              ▼  sdd-openspec  config
  ┌─────────────────────────────────────────┐
  │  openspec/config.yaml                   │
  │  · platform ID + version                │
  │  · JIRA conventions                     │
  │  · required artifact rules              │
  │  · stable context inherited by all      │
  │    component change packages            │
  └─────────────────────────────────────────┘
              │
              ▼
        Reused by every future change package
```

**What happens:**

The Architect uses OpenSpec to encode the stable context into a reusable project
config. This config is inherited by every component change package.

**Sample prompt:**

> "Using the OpenSpec skill, turn the platform constitution, platform versioning
> model, and JIRA hierarchy conventions into a reusable openspec/config.yaml
> that component teams can inherit."

**Sample output — `openspec/config.yaml` (excerpt):**

```yaml
platform:
  id: "customer-platform"
  current_version: "2026.03"
  baseline_ref: "platform-baseline/2026.03"

jira_conventions:
  platform_issue_prefix: "PLAT"
  component_epic_prefix: "<TEAM>"
  story_prefix: "<TEAM>"
  link_type: "is part of"

artifact_rules:
  proposal_required: true
  delta_specs_required: true
  platform_ref_required: true
  jira_traceability_required: true
  design_required_for: ["medium", "large", "architecture-heavy"]
  tasks_required: true
  archive_required: true
```

---

### Step 1.4 — Agree on platform versioning and JIRA conventions

**Role:** Team Lead (with Architect and Product)
**Tool:** OpenSpec
**Skill:** `sdd-openspec`
**When:** After config is drafted. Validate with team leads before locking.

**Sample prompt:**

> "Using the OpenSpec skill, capture the durable workflow expectations for
> platform-ref.yaml, jira-traceability.yaml, and component-level OpenSpec
> change packages."

---

### Step 1.5 — Write ownership artifacts

**Role:** Architect
**Tool:** Platform template (`ownership/` directory)
**Skill:** Reference [platform-ddd-spec.md](platform-ddd-spec.md)
**When:** After versioning and JIRA conventions are agreed. One-time setup.

```text
  Constitution + brownfield context + contract inventory
              │
              ▼  Architect writes three files
  ┌──────────────────────────────────────────────────────────────┐
  │  ownership/component-ownership-profile-service.md            │
  │    owns:      email update flow, profile update event        │
  │    does NOT:  authentication, sessions, email delivery       │
  │                                                              │
  │  ownership/component-ownership-auth-service.md              │
  │    owns:      login sessions, duplicate-email enforcement    │
  │    does NOT:  profile data shape, event publishing           │
  │                                                              │
  │  ownership/dependency-map.md                                 │
  │    profile ◄──tier 1──► auth   (must change together)       │
  │    profile ──tier 2──► notif   (watch for breakage)         │
  │    profile ──tier 3──► audit   (adapts independently)       │
  │                                                              │
  │  ownership/glossary.md                                       │
  │    email update  — defined                                   │
  │    validation    — defined                                   │
  │    duplicate email — defined                                 │
  └──────────────────────────────────────────────────────────────┘
              │
              ▼
        Assess phase can now look up ownership and impact
        instead of deciding them from memory each time
```

**What happens:**

The Architect writes one `component-ownership-<name>.md` per component,
records all component relationships in `dependency-map.md` with their impact
tier, and seeds `glossary.md` from the constitution and brownfield review.

**Sample prompt:**

> "Using the platform-ddd-spec, write the component ownership file for
> profile-service, add its relationships to the dependency map, and seed
> the shared glossary with the terms email update, validation, and
> duplicate email."

**Exit gate for Platform:**

- Constitution is written and approved by the team.
- Durable project config is encoded in `openspec/config.yaml`.
- Versioning model, JIRA hierarchy, and role map are agreed.
- Teams know which service repos are in scope.
- One `component-ownership-<name>.md` exists for every component in scope.
- `dependency-map.md` records all cross-component relationships with tiers.
- `glossary.md` is seeded with terms used in the constitution and contracts.

---

## Phase 2: Assess

**Goal:** Turn an incoming request into one scoped change package with a clear path forward.

**Owner:** Team Lead | **Support:** Product, Architect, Engineering Manager

```text
  Incoming request
  (initiative · requirement · team proposal)
              │
              ▼
  ┌───────────────────────────────────────────────────────────┐
  │                      ASSESS PHASE                         │
  │                                                           │
  │  Step 2.1  BMAD      ──► Classify: size · type · path    │
  │               │          + ownership lookup              │
  │               ├─ greenfield / brownfield?                 │
  │               ├─ small / medium / large?                  │
  │               ├─ Quick Flow / Standard / Full?            │
  │               └─ read component-ownership-<name>.md      │
  │                   → which component owns this change?     │
  │               │                                          │
  │               ▼                                          │
  │  Step 2.2  OpenSpec  ──► Open change package             │
  │               │           + dependency map lookup        │
  │               │           platform-ref.yaml              │
  │               │             ownership.primary_component  │
  │               │             impact.must_change_together  │
  │               │             impact.watch_for_breakage    │
  │               │             impact.adapts_independently  │
  │               │           jira-traceability.yaml         │
  │               │             (JIRA chain from tier)       │
  │               │                                          │
  │               ▼                                          │
  │  Step 2.3  BMAD      ──► Architect impact review         │
  │                          (when contracts involved)        │
  └───────────────────────────────────────────────────────────┘
              │
              ▼
        Assessed change package ready
        (scope · classification · ownership · impact tiers ·
         platform refs · issue chain)
```

---

### Step 2.1 — Classify the request

**Role:** Team Lead
**Tool:** BMAD
**Skill:** `sdd-bmad` (routing agent + `rules/track-selection-rules.md`)
**When:** First action when any new request arrives.

```text
  New request
      │
      ▼  sdd-bmad  track-selection-rules
      │
      ├─► Is it greenfield or brownfield?
      │         │
      │         ├─ brownfield ──► anchor in existing conventions
      │         └─ greenfield ──► new patterns allowed
      │
      ├─► How large?
      │         │
      │         ├─ small  (1–5 stories)  ──► Quick Flow
      │         ├─ medium (5–15 stories) ──► Standard / PRD-first
      │         └─ large  (15+ stories)  ──► Full / architecture-heavy
      │
      └─► Component-only, platform-only, or shared?
                │
                ├─ component-only ──► one epic, local OpenSpec
                ├─ platform-only  ──► platform change package only
                └─ shared         ──► PLAT issue + component epics
```

**What happens:**

The Team Lead uses BMAD to classify the request. Before opening any JIRA issue,
they read `ownership/component-ownership-<name>.md` to confirm which component
owns the primary change. The classification answers: is this greenfield or
brownfield? How large is it? Does it need Quick Flow, a PRD, or a full
architecture review? Is it component-only or shared? Which component is the
primary owner?

**Sample prompt:**

> "Using the BMAD skill, classify this validated-email-update request by size,
> impact, and architecture depth. Tell me whether it is platform-only,
> component-only, or shared. Recommend a path. Also check
> component-ownership-profile-service.md and component-ownership-auth-service.md
> to confirm which component owns the email update flow."

**Sample output:**

```
Classification: validated customer email updates

Type: brownfield
Size: medium (2–3 components, 2–4 stories per component)
Impact: high (shared contract, cross-team coordination)

Ownership lookup (component-ownership-profile-service.md):
  profile-service owns: email update flow ✓
  auth-service owns: duplicate-email enforcement ✓
  Primary owner: profile-service
  Auth-service has a dependency — check dependency-map.md for tier

Path: Standard (PRD-first)
Reason: the change touches a shared contract and two component teams.
        Architecture review is needed before component specs begin.

Classification: shared change
- platform issue required (PLAT-123)
- component epics required (PROF-456, AUTH-234)
- contracts.customer-profile.v2 is in scope
```

---

### Step 2.2 — Open the change package

**Role:** Team Lead
**Tool:** OpenSpec
**Skill:** `sdd-openspec` (`/opsx:propose` + `rules/artifact-rules.md`)
**When:** Immediately after classification.

```text
  BMAD classification result
              │
              ▼  sdd-openspec  /opsx:propose
  ┌──────────────────────────────────────────────┐
  │  Change package  chg-profile-email-validation │
  │                                              │
  │  platform-ref.yaml                           │
  │  ├─ platform version: 2026.03                │
  │  ├─ contracts.customer-profile.v2            │
  │  └─ capabilities.customer-identity           │
  │                                              │
  │  jira-traceability.yaml                      │
  │  ├─ PLAT-123  (platform issue)               │
  │  ├─ PROF-456  (profile-service epic)         │
  │  └─ AUTH-234  (auth-service epic)            │
  └──────────────────────────────────────────────┘
              │
              ▼
        Canonical execution container exists
```

**What happens:**

The Team Lead uses OpenSpec to open the change package. This creates the
canonical execution container that all downstream artifacts will live inside.

**Sample prompt:**

> "Using the OpenSpec skill, open the change package for validated customer
> email updates. Identify the affected component repos and create the initial
> platform-ref.yaml and jira-traceability.yaml."

**Sample prompt:**

> "Using the OpenSpec skill, open the change package for validated customer
> email updates. Identify the affected component repos and create the initial
> platform-ref.yaml and jira-traceability.yaml. Also read dependency-map.md
> and populate the impact fields for profile-service."

**Sample output — `platform-ref.yaml`:**

```yaml
platform:
  id: "customer-platform"
  version: "2026.03"
  baseline_ref: "platform-baseline/2026.03"

component:
  name: "profile-service"
  owner_team: "identity-team"

ownership:
  primary_component: "profile-service"
  # verified against ownership/component-ownership-profile-service.md

impact:
  must_change_together: [auth-service]          # tier 1 — from dependency-map.md
  watch_for_breakage:   [notification-service]  # tier 2 — from dependency-map.md
  adapts_independently: [audit-service]         # tier 3 — from dependency-map.md

change:
  change_package_id: "chg-profile-email-validation"
  alignment_type: "shared-change"
  requires_platform_change: true
  platform_change_package_id: "plat-email-validation"

platform_refs:
  contracts:
    - id: "contracts.customer-profile.v2"
      reason: "Profile service reads and writes the shared customer profile contract."
  capabilities:
    - id: "capabilities.customer-identity"
      reason: "The feature changes identity-related customer behavior."
```

**Sample output — `jira-traceability.yaml`:**

```yaml
platform_issue:
  key: "PLAT-123"
  summary: "Platform: Validated customer email updates"

component_epics:
  - key: "PROF-456"
    component: "profile-service"
    summary: "Profile: Validated email update flow"
  - key: "AUTH-234"
    component: "auth-service"
    summary: "Auth: Duplicate email enforcement alignment"
```

---

### Step 2.3 — Architect impact review (when contracts are involved)

**Role:** Architect
**Tool:** BMAD
**Skill:** `sdd-bmad` (architecture playbook)
**When:** When the BMAD classification shows shared contract impact.

**Sample prompt:**

> "Using the BMAD skill, assess whether validated customer email updates
> require a shared platform change package, a contract update, or only local
> component changes. Identify the main architectural dependencies and
> cross-team risks for PLAT-123."

**Sample output:**

```
Architecture assessment — PLAT-123

Decision: shared change required
- contracts.customer-profile.v2 must be reviewed but not bumped in v1
- AUTH-234 must stay aligned with PROF-456 for duplicate email checks
- notification-service may need a follow-up epic

Risks:
- Existing consumers may assume weaker validation than the new flow enforces.
- Auth duplicate-email logic must not diverge from profile validation behavior.
```

**Exit gate for Assess:**

- Change package exists with one clear owner.
- Ownership verified against `component-ownership-<name>.md`.
- Size and impact are classified separately (medium size, high impact).
- `platform-ref.yaml` pins the platform version, affected refs, ownership,
  and impact tiers (populated from `dependency-map.md`).
- `jira-traceability.yaml` records PLAT-123, PROF-456, AUTH-234 — the issue
  chain is derived from the impact tier (tier 1 → AUTH-234 epic opened now).
- Next artifact and owner are clear: `proposal.md` owned by Product.

---

## Phase 3: Specify

**Goal:** Define the required behavior before planning starts.

**Owner:** Product | **Support:** Team Lead, Architect, Developers

**Rule: once work enters the component repo, use OpenSpec only.**

```text
  Assessed change package
  (PLAT-123 · PROF-456 · AUTH-234)
              │
              ▼
  ┌───────────────────────────────────────────────────────────┐
  │                      SPECIFY PHASE                        │
  │                                                           │
  │  Step 3.1  OpenSpec ──► proposal.md                      │
  │               │         (problem · goals · non-goals ·   │
  │               │          platform refs · acceptance)      │
  │               │                                          │
  │               ▼                                          │
  │  Step 3.2  OpenSpec ──► delta specs                      │
  │               │         ADDED / MODIFIED / REMOVED       │
  │               │         + edge cases from Developers      │
  │               │                                          │
  │               ▼                                          │
  │  Step 3.3  OpenSpec ──► Scope + alignment check          │
  │                         (Team Lead review)               │
  └───────────────────────────────────────────────────────────┘
              │
              ▼
        Approved spec package ready
        (proposal · specs · platform refs confirmed · JIRA updated)
```

---

### Step 3.1 — Write the proposal

**Role:** Product
**Tool:** OpenSpec
**Skill:** `sdd-openspec` (`proposal.md` + `artifact-rules.md`)
**When:** First action in Specify.

**What happens:**

Product reads `ownership/glossary.md` before writing `proposal.md`. Every term
used in goals and acceptance criteria must be in the glossary. Missing terms are
added before the proposal is approved (rule O-2). Then Product uses OpenSpec to
write `proposal.md`.

**Sample prompt:**

> "Using the OpenSpec skill, create proposal.md for validated customer email
> updates. Before writing, check glossary.md to confirm that 'email update',
> 'validation', and 'duplicate email' are defined. Define the problem statement,
> goals, non-goals, affected platform refs, and acceptance summary using only
> glossary terms."

**Sample output — `proposal.md`:**

```markdown
# Proposal: Validated Customer Email Updates

Change package: `chg-profile-email-validation`
Platform issue: `PLAT-123`
Component epic: `PROF-456`

## Problem

Customers can attempt to update their email address today, but validation and
failure handling are inconsistent across profile and auth flows.

## Goals

- validate new email addresses before persistence
- return clear user-facing validation failures
- keep shared customer profile behavior aligned across services

## Non-goals

- redesign account recovery
- migrate all historical profile records in this change
- replace the shared customer profile contract version

## Affected platform refs

- `capabilities.customer-identity`
- `contracts.customer-profile.v2`
- `principles.observability`

## Acceptance summary

- invalid email formats are rejected before persistence
- duplicate-email attempts return explicit failure behavior
- successful updates emit the expected profile event fields
```

---

### Step 3.2 — Write delta specs

**Role:** Product + Developer
**Tool:** OpenSpec
**Skill:** `sdd-openspec` (delta specs + `artifact-rules.md`)
**When:** After proposal is approved.

```text
  Approved proposal.md
              │
              ▼  sdd-openspec  delta specs
  ┌─────────────────────────────────────────────────┐
  │  specs/validated-email-updates/spec.md          │
  │                                                 │
  │  ## MODIFIED                                    │
  │  POST /profile/{id}/email                       │
  │    MUST validate format before write            │
  │    MUST reject duplicates → HTTP 409            │
  │                                                 │
  │  ## ADDED                                       │
  │  Validation failure response shape              │
  │                                                 │
  │  ## REMOVED                                     │
  │  Silent acceptance of malformed input           │
  │                                                 │
  │  ## Edge cases  (Developer review)              │
  │  Empty string → 400 EMAIL_REQUIRED              │
  │  Unicode RFC 5322 valid → accept                │
  └─────────────────────────────────────────────────┘
```

**What happens:**

Product writes delta specs using `ADDED`, `MODIFIED`, and `REMOVED` sections.
Developers review for testability and edge-case coverage.

**Sample prompt (Product):**

> "Using the OpenSpec skill, create delta specs for profile-service and
> auth-service. Make the shared contract references explicit. Use ADDED,
> MODIFIED, and REMOVED sections."

**Sample output — `specs/validated-email-updates/spec.md` (excerpt):**

```markdown
# Delta Spec: Validated Email Updates — profile-service

## MODIFIED

- `POST /profile/{id}/email`
  - MUST validate email format before write
  - MUST reject duplicate emails with HTTP 409 and error code EMAIL_ALREADY_EXISTS
  - MUST emit a structured log entry on validation failure

## ADDED

- Validation failure response shape:
  ```json
  {
    "error": "EMAIL_ALREADY_EXISTS",
    "message": "This email address is already associated with another account."
  }
  ```

## REMOVED

- Silent acceptance of malformed email input (previously allowed to pass)
```

**Sample prompt (Developer — edge case review):**

> "Using the OpenSpec skill, review the approved user stories and turn them
> into executable acceptance expectations, including edge cases for duplicate
> emails, invalid formats, and downstream sync failures."

**Sample output (Developer addition to spec):**

```markdown
## Edge cases

- Empty string inputs MUST be rejected with HTTP 400 and error code EMAIL_REQUIRED
- Unicode email addresses that pass RFC 5322 validation MUST be accepted
- Downstream profile event emission failure MUST NOT leave the record in an inconsistent state
```

---

### Step 3.3 — Scope and alignment check

**Role:** Team Lead
**Tool:** OpenSpec
**Skill:** `sdd-openspec`
**When:** Before approving the spec for planning.

**Sample prompt:**

> "Using the OpenSpec skill, review this spec package for scope control,
> platform alignment, and readiness to move into planning. Confirm that the
> component specs reference the correct platform version, contracts, and JIRA
> issue chain."

**Exit gate for Specify:**

- `proposal.md` is approved with explicit goals and non-goals.
- All terms in the proposal and delta specs are in `glossary.md`.
- `platform-ref.yaml` includes `alignment_notes.glossary_terms_used`.
- Delta specs use concrete behaviors with `ADDED / MODIFIED / REMOVED`.
- `platform-ref.yaml` confirms the correct contract and capability refs.
- `jira-traceability.yaml` is up to date with PLAT-123, PROF-456.
- Team agrees planning can start without guessing.

---

## Phase 4: Plan

**Goal:** Convert the approved spec into a technical execution plan.

**Owner:** Architect | **Support:** Team Lead, Developers, Product

**Rule: once work is inside the component repo, use OpenSpec only.**

```text
  Platform Plan handoff
  (platform version · refs · shared decisions · rollout constraints)
              │
              ▼
  ┌───────────────────────────────────────────────────────────┐
  │                       PLAN PHASE                          │
  │                                                           │
  │  Step 4.1  OpenSpec ──► design.md                        │
  │               │         (architecture · decisions ·       │
  │               │          dependencies · slices)           │
  │               │                                          │
  │               ▼                                          │
  │  Step 4.2  OpenSpec ──► tasks.md                         │
  │               │         Task 1 ──► PROF-789              │
  │               │         Task 2 ──► PROF-790              │
  │               │         Task 3 ──► PROF-791              │
  │               │                                          │
  │               ▼                                          │
  │  Step 4.3  OpenSpec ──► Developer feasibility review     │
  └───────────────────────────────────────────────────────────┘
              │
              ▼
        Implementation-ready plan
        (design · tasks · story mapping · dependencies visible)
```

---

### Step 4.1 — Write the design

**Role:** Architect
**Tool:** OpenSpec
**Skill:** `sdd-openspec` (`design.md` + `artifact-rules.md`)
**When:** First action in Plan, after receiving the platform handoff.

```text
  Platform handoff                    Approved spec
  (version · refs · constraints)      (proposal · delta specs)
              │                               │
              └───────────────┬───────────────┘
                              │
                              ▼  sdd-openspec  design.md
  ┌──────────────────────────────────────────────────────┐
  │  design.md                                           │
  │                                                      │
  │  platform: 2026.03                                   │
  │  contract: contracts.customer-profile.v2             │
  │                                                      │
  │  [API layer]                                         │
  │      │   validate format + duplicate                 │
  │      ▼                                               │
  │  [profile-service write path]                        │
  │      │   preserve contract shape                     │
  │      ▼                                               │
  │  [profile update event]                              │
  │      │   emit observability metrics                  │
  │      ▼                                               │
  │  [AUTH-234 dependency] ◄── must stay aligned         │
  └──────────────────────────────────────────────────────┘
```

**What happens:**

The Architect reads `platform-ref.yaml` impact tiers before designing. Tier 1
entries (`must_change_together`) become hard constraints in `design.md`. Tier 2
entries (`watch_for_breakage`) become rollout risks. Then the Architect turns
the platform plan handoff into a local `design.md` inside the change package.

**Sample prompt:**

> "Using the OpenSpec skill, turn the platform plan handoff for validated
> customer email updates into design.md for profile-service. Read
> platform-ref.yaml impact fields first: auth-service is tier 1 (must change
> together) — make that a hard constraint in the design. notification-service
> is tier 2 (watch for breakage) — make that a rollout risk. Include shared
> contract implications, service boundaries, and rollout constraints."

**Sample output — `design.md`:**

```markdown
# Design: Profile Service Email Validation

Change package: `chg-profile-email-validation`
Epic: `PROF-456`

## Design summary

Profile service adds a validation layer before persisting a customer email
update and before emitting the shared profile update event.

## Platform alignment

- platform version: `2026.03`
- capability: `capabilities.customer-identity`
- contract: `contracts.customer-profile.v2`

## Main decisions

- perform format validation in profile-service before write operations
- preserve contracts.customer-profile.v2 — avoid contract version bump in this slice
- emit explicit validation-failure metrics and logs for observability

## Dependencies

- AUTH-234 — tier 1 (must change together): auth-service duplicate-email
  checks must ship in the same release as this change. Hard constraint.
- notification-service — tier 2 (watch for breakage): monitor after deploy
  to confirm profile event fields are still consumed correctly. Rollout risk.
- audit-service — tier 3 (adapts independently): no coordination required.

## Delivery slices

1. request validation and API failure handling
2. event mapping and test updates
3. verification and release readiness
```

---

### Step 4.2 — Write tasks

**Role:** Architect + Team Lead
**Tool:** OpenSpec
**Skill:** `sdd-openspec` (`tasks.md` + `artifact-rules.md`)
**When:** After `design.md` is approved.

```text
  design.md  (delivery slices decided)
              │
              ▼  sdd-openspec  tasks.md
  ┌──────────────────────────────────────────────────────┐
  │  tasks.md  (PROF-456)                                │
  │                                                      │
  │  [ ] Task 1 ──► PROF-789                            │
  │       validation: API rejects bad + dup inputs       │
  │            │                                         │
  │            ▼                                         │
  │  [ ] Task 2 ──► PROF-790                            │
  │       validation: event compatible with contract v2  │
  │            │                                         │
  │            ▼                                         │
  │  [ ] Task 3 ──► PROF-791                            │
  │       validation: logs · metrics · rollout notes     │
  └──────────────────────────────────────────────────────┘
              │
              ▼
        Each task = one reviewable PR
```

**Sample prompt:**

> "Using the OpenSpec skill, draft tasks.md for PROF-456. Map each task to a
> story key and include explicit validation criteria."

**Sample output — `tasks.md`:**

```markdown
# Tasks: Profile Service Email Validation

Change package: `chg-profile-email-validation`
Epic: `PROF-456`

- [ ] Task 1 — Add email format and duplicate checks in profile-service API
  - story: PROF-789
  - validation: API rejects invalid and duplicate email inputs with explicit errors

- [ ] Task 2 — Update profile event mapping and regression tests
  - story: PROF-790
  - validation: shared profile event remains compatible with contracts.customer-profile.v2

- [ ] Task 3 — Capture verification evidence and release notes
  - story: PROF-791
  - validation: logs, metrics, and rollout notes are recorded before deploy
```

---

### Step 4.3 — Developer feasibility review

**Role:** Developer
**Tool:** OpenSpec
**Skill:** `sdd-openspec`
**When:** Before the plan is approved for execution.

**Sample prompt:**

> "Using the OpenSpec skill, review the tasks for validated customer email
> updates. Confirm that each task is executable, testable, and small enough
> for a reviewable PR."

**Exit gate for Plan:**

- `design.md` explains why each decision was made.
- `tasks.md` maps every task to a story key with a validation step.
- Each task is small enough for one reviewable PR.
- Dependencies (AUTH-234 alignment) are visible and assigned.
- Team agrees the work can be delivered in controlled slices.

---

## Phase 5: Deliver

**Goal:** Execute the plan through Build → Create PR → Review PR → Verify → Deploy → Archive.

**Owner:** Team Lead | **Support:** Developers, QA, Architect, Product

**Rule: once work is inside the component repo, use OpenSpec only.**

```text
  tasks.md + delivery slices
              │
              ▼
  ┌───────────────────────────────────────────────────────────────────┐
  │                        DELIVER PHASE                              │
  │                    (repeat per task slice)                        │
  │                                                                   │
  │  Step 5.1  Build                                                  │
  │    Developer: task claim ──► implement ──► validate ──► complete  │
  │               │                                                   │
  │               ▼                                                   │
  │  Step 5.2  Create PR                                              │
  │    Developer: PR description links change package + story + tasks │
  │               │                                                   │
  │               ▼                                                   │
  │  Step 5.3  Review PR                                              │
  │    Architect: design drift · contract safety                      │
  │    Team Lead: scope · traceability                                │
  │               │                                                   │
  │               ▼   [feedback resolved]                             │
  │  Step 5.4  Verify                                                 │
  │    QA + Dev: record evidence inside change package                │
  │               │                                                   │
  │               ▼   [evidence recorded]                             │
  │  Step 5.5  Deploy                                                 │
  │    Team Lead: timing · AUTH-234 aligned · rollback ready          │
  │               │                                                   │
  │               ▼   [change deployed]                               │
  │  Step 5.6  Archive                                                │
  │    /openspec-archive ──► delta specs merged ──► package closed    │
  └───────────────────────────────────────────────────────────────────┘
              │
              ▼
        Change complete + archived
        (PLAT-123 → PROF-456 → stories → PRs → verified → deployed)
```

---

### Step 5.1 — Build (slice 1)

**Role:** Developer
**Tool:** OpenSpec
**Skill:** `sdd-openspec` (task execution + CLI tracking)
**When:** For each task in `tasks.md`, in order.

```text
  tasks.md  Task 1 (PROF-789)
              │
              ▼  agentic-agent CLI
  task list
      │
      ▼
  task claim PROF-789  ──► branch locked · trace started
      │
      ▼
  implement code change
      │
      ▼
  agentic-agent validate  ──► gate checks pass?
      │                                │
      │  yes                           │  no ──► fix · retry
      ▼                                │
  task complete PROF-789 ◄─────────────┘
      │
      ▼
  Ready to create PR
```

**What happens:**

The Developer claims Task 1 via the `agentic-agent` CLI, implements the email
format and duplicate checks, and updates task state.

**CLI sequence:**

```bash
agentic-agent task list
agentic-agent task claim PROF-789
# ... implement the code change ...
agentic-agent validate
agentic-agent task complete PROF-789
```

**Sample prompt:**

> "Using the OpenSpec skill, implement the current PROF-789 task slice
> according to the approved local spec. Keep the change narrow enough for one
> reviewable PR."

---

### Step 5.2 — Create PR

**Role:** Developer
**Tool:** OpenSpec
**Skill:** `sdd-openspec` (PR description + `artifact-rules.md`)
**When:** After each task slice is implemented and validated locally.

**Sample prompt:**

> "Using the OpenSpec skill, write the PR description for the PROF-789 slice.
> Reference the change package, affected tasks, story key, and validation
> performed."

**Sample output — PR description:**

```markdown
## Change package: chg-profile-email-validation
## Epic: PROF-456 | Story: PROF-789
## Platform version: 2026.03

### What this PR does

Adds email format validation and duplicate-email rejection to the
profile-service update-email API endpoint.

### Platform alignment

- contracts.customer-profile.v2: preserved, no contract version bump
- principles.observability: validation failures now emit structured logs

### Validation performed

- Unit tests: format rejection, duplicate rejection, valid input acceptance
- Contract compatibility: profile event shape unchanged
- Metrics: validation-failure counter verified in local run

### Tasks completed

- [x] Task 1 — Add email format and duplicate checks in profile-service API
  - story: PROF-789

### Related

- PLAT-123 — platform issue
- PROF-456 — component epic
- AUTH-234 — parallel auth-service alignment (separate PR)
```

---

### Step 5.3 — Review PR

**Role:** Architect (design integrity) + Team Lead (scope + traceability)
**Tool:** OpenSpec
**Skill:** `sdd-openspec`
**When:** After PR is opened.

**Sample prompt (Architect):**

> "Using the OpenSpec skill, review this PR for design drift, contract safety,
> and alignment with the approved platform refs. Confirm that tier 1 dependency
> (auth-service) is addressed in this slice or explicitly staged."

**Sample prompt (Team Lead):**

> "Using the OpenSpec skill, confirm that this PR references the change
> package, affected tasks, story keys, and validation performed, and that it
> is scoped to one reviewable slice. Check the PR description for tier 1 and
> tier 2 dependency verification notes."

---

### Step 5.4 — Verify

**Role:** QA / Validation + Developer
**Tool:** OpenSpec
**Skill:** `sdd-openspec`
**When:** After PR is reviewed and before deploy decision.

**What happens:**

Verification evidence is recorded inside the change package. This includes
test results, observability checks, and rollback readiness notes.

**Sample prompt:**

> "Using the OpenSpec skill, record the verification evidence for PROF-789
> inside the change package. Include test pass results, metrics confirmation,
> and any rollout or rollback notes."

---

### Step 5.5 — Deploy

**Role:** Team Lead
**Tool:** OpenSpec
**Skill:** `sdd-openspec`
**When:** After verification passes and review feedback is resolved.

**What happens:**

Team Lead coordinates deploy timing, dependencies (AUTH-234 must be aligned
before deploy), and rollback readiness. Product confirms business acceptance.

**Exit criteria before deploy:**

- All tasks for this slice are complete.
- PR review feedback is resolved or explicitly deferred.
- Verification evidence is recorded in the change package.
- AUTH-234 alignment is confirmed or explicitly staged.

---

### Step 5.6 — Archive

**Role:** Developer + Team Lead
**Tool:** OpenSpec
**Skill:** `sdd-openspec` (`/openspec-archive` command)
**When:** After the change is deployed.

```text
  Change deployed + verified
              │
              ▼  /openspec-archive
  ┌────────────────────────────────────────────────────────┐
  │  Archive                                               │
  │                                                        │
  │  .sdd-spec/changes/chg-profile-email-validation/       │
  │      delta specs ──► merged into .sdd-spec/specs/      │
  │      tasks.md    ──► all complete                      │
  │      PR links    ──► recorded                          │
  │      evidence    ──► recorded                          │
  │                                                        │
  │  jira-traceability.yaml                                │
  │      PLAT-123 → PROF-456 → PROF-789/790/791 → PRs      │
  │                                                        │
  │  Change package status: CLOSED                         │
  └────────────────────────────────────────────────────────┘
              │
              ▼
        New component truth promoted
        (delta specs are now canonical specs)
```

**What happens:**

Archive merges the local delta specs into the component's main spec directory,
closes the change package, and records the delivery as history.

**CLI command:**

```bash
/openspec-archive
```

**Sample prompt:**

> "Using the OpenSpec skill, archive the change package for
> chg-profile-email-validation. Merge the delta specs into the main spec
> directory, update jira-traceability.yaml with final story and PR links,
> and mark the change as complete."

**Exit gate for Deliver (and the full workflow):**

- All planned tasks are complete or intentionally deferred with a logged reason.
- PRs are reviewed and feedback is resolved.
- Verification evidence is recorded.
- Tier 1 dependency verification is recorded in the archive (auth-service aligned).
- Tier 2 consumer check is noted (notification-service checked after deploy).
- Deploy decisions and rollback notes are captured.
- If any ownership boundary or dependency tier changed during delivery, the
  relevant `component-ownership-<name>.md` or `dependency-map.md` is flagged
  for update in the platform repo.
- Change package is archived.
- `jira-traceability.yaml` reflects the final PLAT-123 → PROF-456 → story → PR chain.

---

## Complete workflow summary

```text
 ┌──────────────────────────────────────────────────────────────────────────────────┐
 │  PHASE 1: PLATFORM          Owner: Architect                                     │
 │  ├─ 1.1  BMAD      sdd-bmad      Review brownfield platform state               │
 │  ├─ 1.2  Speckit   sdd-speckit   Draft platform constitution                    │
 │  ├─ 1.3  OpenSpec  sdd-openspec  Encode durable project config                  │
 │  ├─ 1.4  OpenSpec  sdd-openspec  Agree versioning + JIRA conventions            │
 │  └─ 1.5  Architect ownership/   Write ownership boundaries · dependency         │
 │                                  map · shared glossary                           │
 │              │                                                                   │
 │              ▼  constitution · config · role map · ownership · dep map · glossary│
 ├──────────────────────────────────────────────────────────────────────────────────┤
 │  PHASE 2: ASSESS            Owner: Team Lead                                     │
 │  ├─ 2.1  BMAD      sdd-bmad      Classify: size · type · path                  │
 │  │                               + read component-ownership for primary owner   │
 │  ├─ 2.2  OpenSpec  sdd-openspec  Open change package + read dependency-map      │
 │  │                               → populate impact tiers in platform-ref.yaml   │
 │  └─ 2.3  BMAD      sdd-bmad      Architect impact review (shared changes)       │
 │              │                                                                   │
 │              ▼  change package · ownership · impact tiers · issue chain         │
 ├──────────────────────────────────────────────────────────────────────────────────┤
 │  PHASE 3: SPECIFY           Owner: Product          [OpenSpec ONLY in comp repo] │
 │  ├─ 3.1  OpenSpec  sdd-openspec  Check glossary.md → write proposal.md          │
 │  │                               (all terms must be in glossary before approval) │
 │  ├─ 3.2  OpenSpec  sdd-openspec  Write delta specs (ADDED / MODIFIED / REMOVED) │
 │  └─ 3.3  OpenSpec  sdd-openspec  Scope + alignment check (Team Lead)            │
 │              │                                                                   │
 │              ▼  approved spec · glossary aligned · confirmed platform refs       │
 ├──────────────────────────────────────────────────────────────────────────────────┤
 │  PHASE 4: PLAN              Owner: Architect        [OpenSpec ONLY in comp repo] │
 │  ├─ 4.1  OpenSpec  sdd-openspec  Read impact tiers → write design.md            │
 │  │                               (tier 1 = hard constraint · tier 2 = risk)     │
 │  ├─ 4.2  OpenSpec  sdd-openspec  Write tasks.md (story-mapped · validation gate) │
 │  └─ 4.3  OpenSpec  sdd-openspec  Developer feasibility review                   │
 │              │                                                                   │
 │              ▼  design (with tier constraints) · tasks · story mapping · slices  │
 ├──────────────────────────────────────────────────────────────────────────────────┤
 │  PHASE 5: DELIVER           Owner: Team Lead        [OpenSpec ONLY in comp repo] │
 │  ├─ 5.1  OpenSpec  sdd-openspec  Build — claim · implement · validate · complete │
 │  ├─ 5.2  OpenSpec  sdd-openspec  Create PR — link change pkg · tasks · stories  │
 │  ├─ 5.3  OpenSpec  sdd-openspec  Review PR — design integrity · scope check     │
 │  ├─ 5.4  OpenSpec  sdd-openspec  Verify — record evidence in change package     │
 │  ├─ 5.5  OpenSpec  sdd-openspec  Deploy — timing · dependencies · rollback      │
 │  └─ 5.6  OpenSpec  sdd-openspec  Archive — merge delta specs · close package    │
 │              │                                                                   │
 │              ▼  change complete · archived · truth promoted                     │
 └──────────────────────────────────────────────────────────────────────────────────┘
```

---

## Adoption sequence

Teams do not need to run all five phases at once. The methodology is designed
for two-iteration adoption:

**Iteration 1** (learn how to define work well):

```text
  ┌─────────────────────────────────────────────────────┐
  │  ITERATION 1                                        │
  │                                                     │
  │  [Platform] ──► [Assess] ──► [Specify]              │
  │                                  │                  │
  │                                  ▼                  │
  │                           STOP · EVALUATE           │
  │                                                     │
  │  Are specs clear enough to plan without guessing?   │
  │    yes ──► move to Iteration 2                      │
  │    no  ──► fix constitution · config · templates    │
  └─────────────────────────────────────────────────────┘
```

**Iteration 2** (learn how to deliver work well):

```text
  ┌─────────────────────────────────────────────────────┐
  │  ITERATION 2                                        │
  │                                                     │
  │  [Plan] ──► [Deliver]                               │
  │                  │                                  │
  │             Build ──► PR ──► Review                 │
  │             Verify ──► Deploy ──► Archive           │
  └─────────────────────────────────────────────────────┘
```

Add Plan and Deliver once Specify is consistently producing strong specs.

---

## Traceability chain

Every artifact in the workflow traces back to a shared root and forward to an
archived result. The diagram below shows the full traceability chain for the
running example.

```text
  Platform repo (canonical shared truth)
  ┌──────────────────────────────────────┐
  │  platform-baseline/2026.03           │
  │  capabilities.customer-identity      │
  │  contracts.customer-profile.v2       │
  │  principles.observability            │
  └──────────────────────────────────────┘
              │  publishes version + refs
              ▼
  PLAT-123  (platform issue)
              │
        ┌─────┴──────┐
        │            │
        ▼            ▼
  PROF-456        AUTH-234
  (profile-svc)   (auth-svc)
        │
        │  component repo  (local truth — OpenSpec only)
        ▼
  change package: chg-profile-email-validation
  ├─ platform-ref.yaml      (pins 2026.03 + contract refs)
  ├─ jira-traceability.yaml (links PLAT-123 → PROF-456)
  ├─ proposal.md            (goals · non-goals · acceptance)
  ├─ specs/validated-email-updates/spec.md  (delta: ADDED/MODIFIED/REMOVED)
  ├─ design.md              (architecture · decisions · slices)
  ├─ tasks.md
  │   ├─ Task 1 ──► PROF-789 ──► PR #1  (validated · merged)
  │   ├─ Task 2 ──► PROF-790 ──► PR #2  (validated · merged)
  │   └─ Task 3 ──► PROF-791 ──► PR #3  (validated · merged)
  └─ [archived]
       delta specs ──► .sdd-spec/specs/ (new component truth)
```

---

## Key rules to keep visible

1. **Constitution first.** No spec-driven change should start without a
   platform constitution defining the quality bar.

2. **Classify before specifying.** BMAD classification happens in Assess.
   Teams should not start writing specs without knowing whether the change is
   platform-level, component-level, or shared.

3. **Component repos use OpenSpec only.** BMAD and Speckit are permitted on
   the platform side. Once work enters a component repo, only OpenSpec artifacts
   are used.

4. **Never skip task claim.** The `agentic-agent task claim <ID>` command is
   mandatory. It locks the branch and starts the trace. Skipping it breaks
   traceability.

5. **Archive is required.** Archive is not optional cleanup. It is the step
   that promotes the delta specs into the component's new canonical truth.

6. **Never close a change without verification evidence.** Evidence must be
   inside the change package before deploy, not added afterwards.

7. **Ownership before epic.** Before opening any JIRA epic, confirm the
   primary owner using `component-ownership-<name>.md`. The wrong epic owner
   is one of the most common sources of misaligned delivery.

8. **Impact tiers drive JIRA structure.** Read `dependency-map.md` at Assess
   and populate `platform-ref.yaml` impact fields. Tier 1 components must have
   coordinated epics opened immediately. Never override tiers without a platform
   change package.

9. **Glossary before proposal.** All terms used in a proposal or delta spec
   must be in `glossary.md` before the proposal is approved. If a term is
   missing, add it first. Do not invent synonyms for existing glossary terms
   inside a spec.
