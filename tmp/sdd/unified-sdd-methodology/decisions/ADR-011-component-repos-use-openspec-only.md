---
id: "methodology/ADR-011"
title: "Use OpenSpec only inside component repositories"
status: "proposed"
date: "2026-03-08"
initiative: "unified-sdd-methodology"
author: "javierbenavides + codex"
deciders:
  - "platform-methodology-working-group"
---

## Context

The methodology now distinguishes between:

- platform-level work that sets shared rules, contracts, and planning boundaries
- component-level work that turns those boundaries into local implementation

The current documentation still mixes BMAD, OpenSpec, and Speckit inside some
component-level examples. That creates two problems:

- teams do not know which skill should own the local component change package
- component repositories end up with a mixed authoring model instead of one
  simple, repeatable workflow

The platform side benefits from combining the three skills because it must:

- frame context and routing
- define shared guardrails
- reason about cross-team impact

The component side has a different job. It must:

- pin platform alignment
- define local behavior
- define local design and tasks
- execute and archive the local change package

That is already the natural shape of OpenSpec.

## Options Considered

### Option A: Keep all three skills active inside component repositories

Use BMAD, OpenSpec, and Speckit directly inside component repositories during
Specify, Plan, and Deliver.

**Pros:**
- maximum flexibility
- preserves the original combined-skill model everywhere

**Cons:**
- too much method mixing at the local repo level
- harder to teach and enforce
- component artifacts can drift between skill styles

### Option B: Let each component team choose its own dominant skill

Allow each component team to prefer BMAD, OpenSpec, or Speckit for local work.

**Pros:**
- high local autonomy
- low short-term process resistance

**Cons:**
- inconsistent artifacts across teams
- weak traceability
- harder cross-team review and onboarding

### Option C: Use the combined skill model at platform level, but OpenSpec only inside component repositories ← CHOSEN

Use BMAD, OpenSpec, and Speckit together for platform context, routing,
governance, and shared planning. Once the work enters a component repository,
use OpenSpec only for the local change package and local execution artifacts.

**Pros:**
- one simple component-level operating model
- strong alignment with the existing OpenSpec artifact chain
- easier team adoption and review consistency
- clean separation between platform governance and local implementation

**Cons:**
- requires explicit documentation of the handoff from platform plan to
  component OpenSpec
- reduces local freedom to mix skill styles inside the component repo

## Decision

Adopt this rule:

- platform-level work may use BMAD, OpenSpec, and Speckit together
- component repositories use OpenSpec only for local SDD artifacts and local
  workflow execution

Inside a component repository, OpenSpec owns:

- `platform-ref.yaml`
- `jira-traceability.yaml`
- `proposal.md`
- component delta specs
- `design.md`
- `tasks.md`
- apply, PR traceability updates, and archive-ready artifact closure

BMAD and Speckit stay upstream as platform-side support for:

- context framing
- routing depth
- shared quality rules
- cross-team planning discipline

They do not become the local authoring model inside the component repository.

## Consequences

- the documentation must show a clear handoff from platform `Plan` to component
  `OpenSpec`
- worked examples must use OpenSpec-only prompts for component roles
- the methodology must state that the component repo is the local OpenSpec
  execution boundary
- agent guides must not tell component teams to mix BMAD or Speckit into the
  local component change package
