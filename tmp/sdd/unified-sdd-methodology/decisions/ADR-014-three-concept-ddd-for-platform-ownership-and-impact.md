---
id: "methodology/ADR-014"
title: "Adopt three DDD-derived concepts for platform ownership, dependency impact, and shared language"
status: "proposed"
date: "2026-03-15"
initiative: "unified-sdd-methodology"
author: "javierbenavides + claude"
deciders:
  - "platform-methodology-working-group"
---

## Context

The unified SDD methodology now has a platform master repository that owns
canonical shared truth and multiple component repositories that align to it.

When a change is proposed, the Assess phase must answer three questions quickly
and consistently:

1. Which component owns the thing being changed?
2. Which other components will be affected, and how severely?
3. Are all teams using the same terms when describing the change?

Without a structured answer to these questions, the Assess phase becomes a
judgment call that varies by Team Lead, impact assessment misses affected
components, and specs written by different teams use the same word to mean
different things — causing review cycles, misaligned JIRA epics, and defects
that surface late.

Domain-Driven Design (DDD) offers a well-established vocabulary for answering
these questions. The full DDD model contains many concepts (aggregates, entities,
value objects, domain events, repositories, factories, anti-corruption layers,
and more) that would overwhelm a non-DDD team.

Three concepts from DDD directly map to the three questions above:

| DDD concept | Platform question it answers | Plain name used here |
|-------------|------------------------------|----------------------|
| Bounded Context | Who owns what? | Component ownership boundary |
| Context Map | Who depends on whom and how hard is it to change? | Dependency map |
| Ubiquitous Language | Are we using the same words? | Shared glossary |

The methodology needs a decision on whether to adopt these three concepts, which
artifacts they produce, where those artifacts live, and how agents use them.

## Options Considered

### Option A: Do nothing — rely on informal team knowledge

Teams already know their components. Impact assessment is the Team Lead's
judgment. Terminology drift is caught in review.

**Pros:**
- no new artifacts
- no onboarding cost

**Cons:**
- Assess outcomes vary by Team Lead experience
- missed impact is the most common source of late-stage defects
- "customer" meaning auth-service customer ≠ profile-service customer causes
  spec reviews to loop
- informal knowledge is not readable by agents

### Option B: Adopt the full DDD model

Add bounded contexts, context maps, aggregates, domain events, repositories,
anti-corruption layers, and ubiquitous language to the platform methodology.

**Pros:**
- complete model
- well-documented in the DDD community

**Cons:**
- high onboarding cost for teams with no DDD background
- most concepts answer questions not relevant to this methodology's scope
- risk of the team spending more time on the model than on delivery

### Option C: Adopt three DDD-derived concepts with plain-English names ← CHOSEN

Use only the three DDD concepts that directly solve the three Assess-phase
problems. Rename them to plain English so teams adopt them without DDD training.
Limit each concept to one artifact per platform.

The three concepts and their artifacts:

- **Component ownership boundary** → `ownership/component-ownership-<name>.md`
  (one per component in the platform repository)
- **Dependency map** → `ownership/dependency-map.md`
  (one platform-level file in the platform repository)
- **Shared glossary** → `ownership/glossary.md`
  (one platform-level file in the platform repository)

**Pros:**
- solves the three specific problems without importing the full DDD model
- plain-English names require no DDD background to understand
- three artifacts are easy to maintain and easy for agents to read
- fits directly into the existing platform template directory
- integrates cleanly with `platform-ref.yaml` and `jira-traceability.yaml`

**Cons:**
- teams familiar with DDD will recognise the pattern but use different names;
  the naming table in `platform-ddd-spec.md` bridges this for DDD practitioners
- ownership boundaries require an initial investment during the Platform phase

## Decision

Adopt Option C.

Add the following three artifacts to the platform master repository under
`ownership/`:

```text
platform-repo/
  ownership/
    component-ownership-<name>.md   <- one per component
    dependency-map.md               <- one per platform
    glossary.md                     <- one per platform
```

### Concept 1: Component ownership boundary

Each component gets one `component-ownership-<name>.md` file. It records:

- what the component owns (its authoritative capabilities and data)
- what the component explicitly does NOT own (prevents scope drift)
- which contracts it publishes and consumes

This file is created during the Platform phase and updated when ownership
boundaries change (always requires a platform change package).

During Assess, the Team Lead reads the relevant ownership files to determine
whether a change is `component-only` or crosses an ownership boundary.

### Concept 2: Dependency map

One `dependency-map.md` lives at platform scope. It records:

- which component depends on which other component
- what the dependency is (contract, event, shared data)
- which impact tier applies

Impact tiers:

| Tier | Meaning | JIRA implication |
|------|---------|-----------------|
| `must_change_together` | Shared contract changes require coordinated releases | Coordinate epics, same release window |
| `watch_for_breakage` | A breaking contract change will break the consumer | Open component epic when contract changes |
| `adapts_independently` | Loose coupling; consumer absorbs the change on its own schedule | No forced epic coordination needed |

The impact tier is referenced in `platform-ref.yaml` under the `impact` field.
It makes the JIRA structure decision deterministic: the tier lookup replaces an
informal judgment call.

### Concept 3: Shared glossary

One `glossary.md` lives at platform scope. It records:

- shared terms used in contracts, events, capabilities, and specs
- one definition per term with a clear "what it is NOT" clause
- which capability or contract uses the term

All specs, designs, stories, and PRs must use only glossary terms for
shared concepts. If two teams use the same word with different meanings,
the glossary resolves it before Specify starts.

Referenced in `platform-ref.yaml` under `alignment_notes.glossary_terms_used`.

### Constitution rules added

Three new rules are added to the platform constitution as a result of this ADR:

- **O-1**: Every component must have an ownership boundary file before any
  change that touches shared contracts can be classified in Assess.
- **O-2**: All specs and stories must use glossary terms for shared concepts.
  Terms not in the glossary must be added before Specify is approved.
- **O-3**: The impact tier from the dependency map drives the JIRA structure
  decision. Teams do not override impact tiers without a platform change package.

### Integration with existing artifacts

| Existing artifact | Integration |
|-------------------|-------------|
| `platform-ref.yaml` | Add `ownership.primary_component`, `impact.must_change_together`, `impact.watch_for_breakage`, `impact.adapts_independently`, `alignment_notes.glossary_terms_used` |
| `jira-traceability.yaml` | JIRA issue chain is derived from impact tier lookup |
| `platform-baseline.md` | References the three new artifacts as canonical ownership truth |
| `capability-template.md` | Shared capabilities reference glossary terms |
| `contract-template.md` | Contracts reference glossary terms and appear in dependency map |

## Consequences

- The Platform phase must produce ownership boundary files for all components,
  a dependency map, and a glossary before the methodology is operational.
- Assess gains a deterministic lookup step: read ownership file → read dependency
  map → determine impact tier → determine JIRA structure.
- Specify is blocked on glossary alignment: all shared terms must resolve before
  the proposal is approved.
- Plan traces tasks to `platform-ref.yaml` impact fields, making coordination
  requirements explicit in `design.md` and `tasks.md`.
- Deliver links PRs and archive records to impact-tier-derived JIRA chains,
  closing the traceability loop.
- Agents reading `platform-ref.yaml` can determine impact without asking the
  Team Lead; the Team Lead validates, not originates, impact classification.
- The full rationale and usage guide lives in:
  - `platform-ddd-spec.md`
  - `templates/platform-template/ownership/`
