# Canonical Platform Truth + Versioned Component Alignment

## Purpose

This document explains how the unified SDD methodology should work when the
organization has:

- one master repository for platform-level truth
- multiple component repositories for local implementation
- multiple teams delivering through JIRA

The goal is to keep one durable source of shared truth while still letting
component teams move at implementation speed inside their own repositories.

This model is the recommended way to connect:

- platform specs
- component specs
- change packages
- JIRA epics and stories
- pull requests and delivery evidence

## Core idea

Use three layers with different jobs:

```text
TRUTH LAYER
[Platform master repo]
  shared truth
        |
        v
[Component repo]
  local truth

TRACKING LAYER
[JIRA]
  platform issue -> component epic -> stories

EXECUTION LAYER
[Git + PRs + CI/CD]
  code, review, verification, deploy, archive
```

The model is simple:

- the platform repository owns shared truth
- component repositories own local implementation truth
- JIRA owns workflow status and coordination
- a local read-only platform MCP gateway may expose that platform truth to
  developers without changing the source-of-truth model

Do not let one of those layers silently replace another.

## Definitions

### Canonical platform truth

The platform repository is the only authoritative place for:

- platform principles and guardrails
- shared capabilities
- shared contracts, APIs, events, and schemas
- cross-team ADRs
- shared ownership boundaries
- cross-cutting quality, security, and compliance rules

This repository changes more slowly and should not contain component-level
implementation detail unless that detail is itself part of a shared contract.

### Versioned component alignment

Each component repository records which platform version and platform refs it
aligns to.

This means the component says:

- which platform baseline it follows
- which platform rules or contracts it depends on
- whether the current change is local-only or also changes shared truth

That alignment should be explicit in a durable artifact such as
`platform-ref.yaml`.

### JIRA-linked execution

JIRA tracks the operational chain of delivery:

- platform-level issue or initiative
- component epic per repository or component
- story per reviewable slice or task group

JIRA should summarize and link. It should not become the detailed spec store.

## Source-of-truth rules

These are the non-negotiables:

- the platform repository is upstream truth for shared platform behavior
- the component repository is upstream truth for local component behavior and local implementation detail
- JIRA is upstream truth for workflow state, ownership, and delivery coordination
- a component repository may extend platform rules, but it should not silently weaken or contradict them
- if shared truth changes, the change must be visible in both the platform repo and the affected component repos
- if only local component behavior changes, the component repo may move without changing platform truth

## Repository model

### Platform master repository responsibilities

The platform repository should normally contain:

- platform constitution or principles
- shared capability definitions
- shared contracts or schemas
- platform-level change packages when shared truth changes
- version markers or release tags for platform truth
- ownership and dependency language used by teams

It should also define:

- the versioning approach for platform truth
- the naming scheme for durable refs
- the default JIRA hierarchy conventions

### Component repository responsibilities

Each component repository should normally contain:

- local OpenSpec change packages
- local component specs and delta specs
- `design.md`
- `tasks.md`
- local PR and verification evidence
- `platform-ref.yaml`
- `jira-traceability.yaml`

The component repository is where implementation is planned and executed. It
should not hold a copied editable version of the full platform truth.

### Optional local platform MCP gateway

When hosted infrastructure is not available, developers may run a small local,
read-only MCP server against their local platform clone.

Use it to:

- query platform refs and contracts locally
- validate component alignment against the pinned platform version
- inspect JIRA-linked metadata that is already part of the local artifact chain

Reference:

- [local-platform-mcp-model.md](local-platform-mcp-model.md)

## JIRA model

### Recommended hierarchy

If JIRA supports levels above epic, use this:

```text
[Platform initiative / capability]
        |
        +--> [Component epic A]
        |         |
        |         +--> [Story A1]
        |         +--> [Story A2]
        |
        +--> [Component epic B]
                  |
                  +--> [Story B1]
                  +--> [Story B2]
```

### Fallback hierarchy when epic-under-epic is not available

Standard JIRA often does not support true epic-under-epic modeling.

If that is your situation, use this:

```text
[Platform epic]
        |
        +--> linked [Component epic A]
        |               |
        |               +--> [Story A1]
        |               +--> [Story A2]
        |
        +--> linked [Component epic B]
                        |
                        +--> [Story B1]
                        +--> [Story B2]
```

Use explicit issue links such as:

- implements
- depends on
- blocked by
- relates to

### JIRA issue responsibilities

- `Platform issue`
  - represents the platform-level outcome, shared capability, or shared truth change
  - links to the platform spec or platform change package

- `Component epic`
  - represents one component or repository implementation stream
  - links to the component repo change package and platform refs

- `Story`
  - represents one reviewable slice or task group
  - links to a task ID, PR, and verification evidence

## Artifact model

### `platform-ref.yaml`

This file answers:

- which platform version does this component follow
- which platform refs apply
- does this change require a platform change too

This file belongs in the component repository.

### `jira-traceability.yaml`

This file answers:

- which platform issue owns the larger outcome
- which component epic owns local delivery
- which stories map to tasks and slices

This file belongs in the component repository.

### OpenSpec artifacts

Use OpenSpec inside the component repository for:

- `proposal.md`
- delta specs
- `design.md`
- `tasks.md`

Those artifacts are the detailed implementation-facing truth for the component.

## Component repository skill rule

Use this boundary rule throughout the methodology:

- platform-level work may use BMAD, OpenSpec, and Speckit together
- once the work enters a component repository, use OpenSpec only

This means the component repository uses OpenSpec for:

- `platform-ref.yaml`
- `jira-traceability.yaml`
- `proposal.md`
- component delta specs
- `design.md`
- `tasks.md`
- delivery traceability updates through PR and archive

Do not mix BMAD or Speckit into the local component change package. Use them
upstream on the platform side when teams are still framing, routing, or
governing shared work.

## Platform Plan to component OpenSpec handoff

The platform `Plan` phase must hand the component team enough information to
start local OpenSpec work without guessing.

Required handoff inputs:

- platform version
- platform refs
- shared design and contract decisions
- migration, rollout, or rollback constraints
- platform issue and component epic chain

The handoff flow is:

```text
[Platform Plan]
  shared boundaries
        |
        v
[Component repo alignment]
  platform-ref.yaml
  jira-traceability.yaml
        |
        v
[Component OpenSpec]
  proposal.md -> spec.md -> design.md -> tasks.md
        |
        v
[PR + verification + archive]
```

The platform plan sets shared boundaries. The component repository converts
those boundaries into local executable OpenSpec artifacts.

Worked example:

- [example/08-platform-plan-to-component-openspec.md](example/08-platform-plan-to-component-openspec.md)

## Phase-by-phase workflow inside the methodology

### Platform

The Platform phase defines the baseline in the platform repository.

What happens here:

- define the canonical platform truth location
- define the platform versioning scheme
- define the durable ref naming scheme
- define which artifacts belong in platform vs component scope
- define the JIRA hierarchy and issue-link conventions

Expected outputs:

- platform baseline
- platform versioning approach
- shared ref IDs
- JIRA conventions for platform issue, component epic, and stories
- component alignment templates

Platform does not yet break work into local component detail. It defines the
shared rules that later changes must follow.

### Assess

Assess is where the issue chain and scope chain are created.

The Team Lead should classify the incoming work as one of these:

- `platform-only`
- `component-only`
- `shared platform + component`
- `multi-component shared change`

Assess must answer:

- which platform refs are affected
- which component repositories are affected
- whether a platform change package is needed
- which JIRA issue chain is needed

Expected outputs:

- change classification
- linked platform issue when required
- component epic
- initial `platform-ref.yaml`
- initial `jira-traceability.yaml`

### Specify

Specify is where the component-level detail is defined.

Use the component repository to create:

- `proposal.md`
- delta specs
- clarification notes
- checklist results

During Specify, the component spec must reference:

- platform version
- platform refs
- whether the change is local-only or shared

If shared truth changes, Specify must also open or update the linked platform
change package in the platform repository.

This is the key rule:

- local behavior change -> component spec delta only
- shared truth change -> platform delta + component delta

### Plan

Plan turns the approved component spec into executable work.

The design and tasks should show:

- which platform refs constrain the design
- which component tasks implement the change
- which stories correspond to which tasks
- which repositories or teams must coordinate

Good planning creates a clean chain:

```text
platform ref
   -> component spec delta
      -> design decision
         -> task
            -> story
               -> PR
```

Expected outputs:

- `design.md`
- `tasks.md`
- final task-to-story mapping
- updated `platform-ref.yaml`
- updated `jira-traceability.yaml`

### Deliver

Deliver executes the planned slices inside the component repository.

Each slice should normally have:

- one story or tightly related story set
- one PR or tightly related PR set
- explicit validation evidence

The delivery chain should look like this:

```text
[Platform issue]
        |
        v
[Component epic]
        |
        v
[Story]
        |
        v
[PR]
        |
        v
[Verification]
        |
        v
[Archive]
```

At archive time:

- the component change package is closed
- JIRA links are complete
- the final platform alignment is visible
- if shared truth changed, the platform repo is updated too

## Change scenarios

### Scenario 1: Local-only component change

Example:

- small UI validation change inside one component
- no shared contract or platform rule changes

Expected pattern:

- no platform change package
- one component epic
- component stories from `tasks.md`
- component repo updates only
- `platform-ref.yaml` still records which platform version was followed

### Scenario 2: Shared contract change

Example:

- API or event schema changes across components

Expected pattern:

- platform issue exists
- platform change package exists
- component epics exist for each affected repository
- component specs reference the new or modified contract
- migration or rollout sequencing is visible in planning and delivery

### Scenario 3: Platform rule adoption

Example:

- platform updates a logging or security rule
- components must align over time

Expected pattern:

- platform repo updates shared truth first
- platform issue tracks adoption
- each component epic is an adoption stream
- each component repo updates `platform-ref.yaml` when aligned to the new version

## Common mistakes to avoid

- copying platform specs into component repositories
- editing shared truth only in component repositories
- using JIRA descriptions as the full spec
- creating component stories with no link to tasks or platform refs
- creating platform issues with no linked component delivery stream
- letting PRs reference only epics and not the actual implementing story

## Operational rules

- every component change must know its platform version
- every shared change must know its platform issue and component epics
- every story should point to a task or slice
- every PR should point to at least one story
- archive should close the traceability chain, not just the code change

## Recommended templates

Use these templates:

- [templates/README.md](templates/README.md)
- [templates/platform-template/README.md](templates/platform-template/README.md)
- [templates/component-boilerplate/README.md](templates/component-boilerplate/README.md)
- [templates/platform-mcp-boilerplate/README.md](templates/platform-mcp-boilerplate/README.md)
- [templates/platform-ref.yaml](templates/platform-ref.yaml)
- [templates/jira-traceability.yaml](templates/jira-traceability.yaml)

## Summary

This model gives you:

- one source of shared truth at platform scope
- one source of detailed execution truth at component scope
- one tracking chain in JIRA
- clean traceability from platform intent to component implementation

In one sentence:

- platform truth stays canonical upstream
- component repositories align to it by version
- OpenSpec breaks work down locally
- JIRA tracks the delivery chain from platform issue to story to PR
