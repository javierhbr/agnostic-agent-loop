# Unified SDD Methodology: 5-Phase Reference

> Reference document for the five-phase Spec-Driven Development lifecycle.
> Each phase has a single owner, defined entry/exit criteria, and explicit artifact outputs.

---

## Phase 1: Platform

| Field            | Value                                                        |
|------------------|--------------------------------------------------------------|
| Owner            | Architect                                                    |
| Support          | Product, Team Lead                                           |
| Goal             | Create durable shared context before any change-level work   |
| Entry criteria   | New platform or platform evolution needed                    |

**Activities:**

1. BMAD brownfield review -- assess existing platform state
2. Speckit constitution -- encode platform identity, boundaries, and principles
3. OpenSpec config.yaml -- machine-readable platform configuration
4. Versioning and JIRA conventions -- agree on release strategy and issue taxonomy
5. Ownership artifacts -- assign component stewardship

**Artifacts produced:**

- `constitution` -- platform identity, principles, boundaries
- `openspec/config.yaml` -- encoded platform configuration
- `component-ownership-<name>.md` -- per-component ownership record
- `dependency-map.md` -- inter-component dependency graph
- `glossary.md` -- canonical term definitions

**Tools:** BMAD + OpenSpec + Speckit (all three required)

**Exit criteria:**

- Constitution approved by stakeholders
- Config encoded and validated
- Ownership artifacts written for every component
- Versioning conventions agreed and documented

---

## Phase 2: Assess

| Field            | Value                                                        |
|------------------|--------------------------------------------------------------|
| Owner            | Team Lead                                                    |
| Support          | Product, Architect, EM                                       |
| Goal             | Turn incoming request into one scoped change package         |
| Entry criteria   | New request arrives (initiative, requirement, team proposal) |

**Activities:**

1. BMAD classification -- greenfield/brownfield, size, path
2. Read `component-ownership` -- verify ownership for affected areas
3. Open change package -- initialize the structured change container
4. Populate `platform-ref.yaml` -- map impact tiers from `dependency-map.md`
5. Populate `jira-traceability.yaml` -- record upstream and downstream issue links
6. Architect impact review -- required when contracts or shared interfaces are involved

**Artifacts produced:**

- Change package (directory with standard structure)
- `platform-ref.yaml` -- impact tier mapping
- `jira-traceability.yaml` -- JIRA chain record

**Tools:** BMAD + OpenSpec (Speckit optional for clarification)

**Exit criteria:**

- Change package exists with clear owner
- Ownership verified against component-ownership records
- Size and impact classified
- `platform-ref.yaml` populated with impact tiers
- JIRA chain recorded in `jira-traceability.yaml`

---

## Phase 3: Specify

| Field            | Value                                                        |
|------------------|--------------------------------------------------------------|
| Owner            | Product                                                      |
| Support          | Team Lead, Architect, Developers                             |
| Goal             | Define required behavior before planning                     |
| Entry criteria   | Assessed change package exists                               |

**Activities:**

1. Check `glossary.md` -- ensure all terms are defined before writing specs
2. Write `proposal.md` -- describe the change in business and technical terms
3. Write delta specs -- ADDED, MODIFIED, REMOVED behavioral changes
4. Scope and alignment check -- confirm change stays within assessed boundaries

**Artifacts produced:**

- `proposal.md` -- change description and rationale
- Delta specs -- behavioral diffs (ADDED/MODIFIED/REMOVED)
- Confirmed platform refs

**Tools:** OpenSpec ONLY (used in component repos)

**Exit criteria:**

- Proposal approved by stakeholders
- All terms present in glossary
- Delta specs contain concrete, testable behaviors
- `platform-ref.yaml` confirmed against actual impact
- JIRA updated with spec links

---

## Phase 4: Plan

| Field            | Value                                                        |
|------------------|--------------------------------------------------------------|
| Owner            | Architect                                                    |
| Support          | Team Lead, Developers, Product                               |
| Goal             | Convert approved spec into technical execution plan          |
| Entry criteria   | Approved spec package                                        |

**Activities:**

1. Read `platform-ref.yaml` impact tiers -- tier 1 = hard constraints, tier 2 = risks
2. Write `design.md` -- document decisions, trade-offs, tier-based constraints
3. Write `tasks.md` -- decompose into deliverable units, map to JIRA story keys
4. Developer feasibility review -- confirm tasks are realistic and sized correctly

**Artifacts produced:**

- `design.md` -- technical design with decision rationale
- `tasks.md` -- ordered task list mapped to story keys
- ADRs (Architecture Decision Records) -- when significant trade-offs exist

**Tools:** OpenSpec ONLY (used in component repos)

**Exit criteria:**

- Design explains all significant decisions
- Tasks mapped to JIRA stories
- Each task corresponds to one reviewable PR
- Dependencies between tasks are visible and ordered

---

## Phase 5: Deliver

| Field            | Value                                                        |
|------------------|--------------------------------------------------------------|
| Owner            | Team Lead                                                    |
| Support          | Developers, QA, Architect, Product                           |
| Goal             | Execute plan through controlled slices                       |
| Entry criteria   | Approved plan with tasks                                     |

**Internal slices (per task):**

1. **Build** -- implement the code change in an isolated worktree
2. **Create PR** -- open pull request with spec references
3. **Review PR** -- independent code review against spec behaviors
4. **Verify** -- run tests, check verification evidence
5. **Deploy** -- release to target environment
6. **Archive** -- close change package, store final artifacts

**Artifacts produced:**

- Implemented code
- Reviewed PRs with approval records
- Verification evidence (test results, coverage, manual checks)
- Archived change package

**Tools:** OpenSpec ONLY (used in component repos)

**Exit criteria:**

- All tasks complete
- PRs reviewed and merged
- Verification evidence recorded for every task
- Tier 1 dependencies verified against `platform-ref.yaml`
- Change package archived
- JIRA chain complete (all issues closed or linked)

---

## Core Boundary Rules

These rules apply across all five phases without exception.

1. **Platform artifacts live in the platform repo.** Constitutions, dependency maps, glossaries, and ownership records are never duplicated into component repos.
2. **Component artifacts live in component repos.** Proposals, delta specs, designs, and tasks belong to the repo where the code changes.
3. **OpenSpec is the only spec tool in component repos.** Speckit and BMAD operate at platform level only (Phase 1 and Phase 2).
4. **One change package = one scoped change.** Never bundle unrelated changes. If scope grows, split into separate change packages.
5. **Traceability is mandatory.** Every artifact links back to its parent: task to story, story to spec, spec to change package, change package to platform ref.

---

## Two-Iteration Adoption Overview

Teams adopt the methodology in two iterations to avoid overload.

**Iteration 1 -- Foundation (Phases 1-3):**

- Establish platform context (constitution, config, ownership)
- Practice assessment workflow on real requests
- Write first specs using OpenSpec in component repos
- Goal: team can produce a well-scoped, well-specified change package

**Iteration 2 -- Execution (Phases 4-5):**

- Add planning discipline (design.md, tasks.md, ADRs)
- Integrate delivery slices with CI/CD pipeline
- Establish verification and archival habits
- Goal: team delivers traced, verified changes end-to-end

---

## Phase-by-Phase Operating Table

```
Phase   | Owner     | Tools               | Key Input              | Key Output
--------|-----------|----------------------|------------------------|---------------------------
1 Platf | Architect | BMAD+OpenSpec+Speckit| Platform need          | Constitution, config, maps
2 Assss | Team Lead | BMAD+OpenSpec        | Incoming request       | Change package, refs, JIRA
3 Specf | Product   | OpenSpec only        | Change package         | Proposal, delta specs
4 Plan  | Architect | OpenSpec only        | Approved spec          | Design, tasks, ADRs
5 Delivr| Team Lead | OpenSpec only        | Approved plan          | Code, PRs, evidence
```

---

## Traceability Chain

The following diagram shows the artifact chain from platform down to deployed code. Every link is mandatory -- a missing link means the chain is broken and the change cannot proceed.

```
Platform Constitution
    |
    v
openspec/config.yaml  +  dependency-map.md  +  component-ownership-<name>.md
    |
    v
Change Package (assessed, owned, classified)
    |
    +--- platform-ref.yaml (impact tiers from dependency-map)
    +--- jira-traceability.yaml (issue chain)
    |
    v
proposal.md  -->  delta specs (ADDED / MODIFIED / REMOVED)
    |
    v
design.md  -->  tasks.md  -->  ADRs (when needed)
    |
    v
Task 1 -----> Build --> PR --> Review --> Verify --> Deploy
Task 2 -----> Build --> PR --> Review --> Verify --> Deploy
Task N -----> Build --> PR --> Review --> Verify --> Deploy
    |
    v
Archived Change Package (all evidence, all links intact)
```

Each arrow represents a traceable link. Auditors or agents can walk the chain in either direction: from a deployed PR back to the platform constitution, or from a platform change forward to every affected task.

---

## Quick Reference: Who Owns What

```
Artifact                     | Created in Phase | Owner
-----------------------------|------------------|----------
Constitution                 | 1 Platform       | Architect
config.yaml                  | 1 Platform       | Architect
component-ownership-<name>   | 1 Platform       | Architect
dependency-map.md            | 1 Platform       | Architect
glossary.md                  | 1 Platform       | Architect
Change package               | 2 Assess         | Team Lead
platform-ref.yaml            | 2 Assess         | Team Lead
jira-traceability.yaml       | 2 Assess         | Team Lead
proposal.md                  | 3 Specify        | Product
Delta specs                  | 3 Specify        | Product
design.md                    | 4 Plan           | Architect
tasks.md                     | 4 Plan           | Architect
ADRs                         | 4 Plan           | Architect
Code + PRs                   | 5 Deliver        | Team Lead
Verification evidence        | 5 Deliver        | Team Lead
Archived change package      | 5 Deliver        | Team Lead
```
