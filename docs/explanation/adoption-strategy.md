# Two-Iteration Adoption Strategy

## Why Two Iterations

The unified SDD methodology is adopted in two iterations, not one.
Teams first learn how to define work well (Iteration 1), then apply
those stronger inputs to planning and delivery (Iteration 2).

The reasoning is straightforward: specs must be consistently strong
before delivery workflow is layered on top. If a team cannot write
clear proposals and delta specs, no amount of planning discipline
will compensate. Garbage in, garbage out.

This decision is documented in ADR-006.

---

## Target Model

The full methodology has five phases:

```
[Platform] --> [Assess] --> [Specify] --> [Plan] --> [Deliver]
```

Teams do not adopt all five at once. They adopt in two iterations,
with a hard evaluation gate between them.

---

## Iteration 1: Platform + Assess + Specify

**Goal:** Learn to define work well.

### Phases Covered

- **Platform.** Architect creates constitution, config, and ownership
  artifacts. These form the shared context every other phase depends on.
- **Assess.** Team Lead classifies incoming requests and opens change
  packages. Classification uses size and impact as separate dimensions.
- **Specify.** Product writes proposals and delta specs inside the
  change package.

### What Teams Learn

- How to create durable shared context (constitution, config).
- How to classify changes by size and impact separately, rather than
  conflating them into a single "priority" number.
- How to write clear proposals using glossary-aligned terms so that
  reviewers do not have to guess at meaning.
- How to write delta specs using ADDED / MODIFIED / REMOVED markers
  so that the change surface is explicit.

### Stop and Evaluate

At the end of Iteration 1, the team asks one question:

> Are specs clear enough to plan without guessing?

- If yes: move to Iteration 2.
- If no: fix the constitution, config, or templates. Do not proceed
  until specs are consistently strong.

```
Iteration 1:

[Platform] --> [Assess] --> [Specify]  STOP  --> EVALUATE
                                        |
                              specs clear?
                             /            \
                           yes             no
                            |               |
                     Iteration 2      fix inputs,
                                      repeat
```

---

## Iteration 2: Plan + Deliver

**Goal:** Learn to deliver work well.

### Phases Covered

- **Plan.** Architect writes design.md and tasks.md using impact tiers
  derived from the delta spec. Impact tiers constrain the design:
  higher impact means more review, more verification evidence, and
  smaller delivery slices.
- **Deliver.** Team Lead coordinates the delivery sequence:
  Build, PR, Review, Verify, Deploy, Archive.

### What Teams Learn

- How to turn specs into executable plans with explicit tasks.
- How to use impact tiers as design constraints rather than optional
  metadata.
- How to create reviewable delivery slices -- one task, one PR, one
  review cycle.
- How to archive and promote delta specs to canonical truth once
  the change is deployed.

### Success Criteria

- Each task produces exactly one reviewable PR.
- Verification evidence is inside the change package, not scattered
  across Slack threads or meeting notes.
- The archive step promotes delta specs into the canonical spec,
  so the canonical spec always reflects the deployed system.
- The JIRA chain is complete: platform issue to stories to PRs,
  with no orphaned artifacts.

```
Iteration 2:

[Plan] --> [Deliver]
             |
             +-- Build
             +-- PR
             +-- Review
             +-- Verify
             +-- Deploy
             +-- Archive
```

---

## Evolution Path

The five-phase model is the current target (v1). It will evolve
only when complexity justifies it.

| Version | Phases | Trigger to Split |
|---------|--------|------------------|
| v1 | Platform, Assess, Specify, Plan, Deliver | Current model |
| v2 | Platform, Assess, Specify, Plan, Build, Deploy | Split Deliver when release coordination grows complex enough to justify a separate phase |

The rule is simple: only split a phase when the coordination cost
of keeping it together exceeds the coordination cost of managing
two phases. For most teams, Deliver as a single phase is sufficient
for a long time.

---

## Summary

```
FULL TARGET MODEL:
  [Platform] --> [Assess] --> [Specify] --> [Plan] --> [Deliver]

ITERATION 1 (define work well):
  [Platform] --> [Assess] --> [Specify]  ||  EVALUATE

ITERATION 2 (deliver work well):
                                         [Plan] --> [Deliver]
```

Two iterations. One evaluation gate. Specs before plans.
