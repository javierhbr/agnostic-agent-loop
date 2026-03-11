---
id: "methodology/ADR-006"
title: "Roll out the 5-phase workflow in two iterations"
status: "proposed"
date: "2026-03-08"
initiative: "unified-sdd-methodology"
author: "javierbenavides + codex"
deciders:
  - "platform-methodology-working-group"
---

## Context

The 5-phase workflow is simpler than the earlier draft, but adopting all five
phases at once still asks teams to change intake, specification, planning,
delivery coordination, and release behavior in one move.

That is more change than most teams can absorb safely.

## Options Considered

### Option A: Roll out all 5 phases at once

Adopt Platform, Route, Specify, Plan, and Deliver in a single release.

**Pros:**
- Fastest path to the full target model
- Less transition documentation

**Cons:**
- High adoption risk
- Harder to learn what is working
- More cross-team coordination required up front

### Option B: Roll out in two iterations ← CHOSEN

Adopt the workflow in two steps:

- Iteration 1: Platform, Route, Specify
- Iteration 2: Plan, Deliver

**Pros:**
- Lower adoption risk
- Teams learn the front half of the workflow first
- Better quality of inputs before planning and delivery change
- Easier to pilot and refine

**Cons:**
- Slower path to the full operating model
- Temporary split between front-half and back-half maturity

## Decision

Roll out the methodology in two iterations.

Iteration 1 focuses on the first three phases:

- Platform
- Route
- Specify

Iteration 2 focuses on the remaining two phases:

- Plan
- Deliver

Deliver may split into Build and Deploy in a future evolution, but not during
the first two adoption iterations.

## Consequences

- The proposal should show both the target 5-phase model and the two-step
  adoption path
- Early pilots should validate platform context, routing, and specification
  quality before changing delivery behavior
- Teams can improve planning and delivery later using better upstream inputs
