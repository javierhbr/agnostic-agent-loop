---
id: "methodology/ADR-007"
title: "Make pull request creation and review explicit inside Deliver"
status: "proposed"
date: "2026-03-08"
initiative: "unified-sdd-methodology"
author: "javierbenavides + claude"
deciders:
  - "platform-methodology-working-group"
---

## Context

The current 5-phase model describes Deliver as build, verify, deploy, and
archive. That is directionally correct, but it leaves pull request creation
and review implicit.

Across the three source skills, review is already important:

- BMAD includes QA / review and code-review as part of the workflow
- OpenSpec prefers incremental, reviewable change sets
- Speckit requires tasks to be independently reviewable

If pull request creation and review remain implicit, teams may treat them as
local process details instead of a required quality and coordination step.

## Options Considered

### Option A: Keep pull request work implicit

Leave PR creation and review inside local team practices.

**Pros:**
- Less visible process in the methodology
- Flexible for each team

**Cons:**
- Weakens review discipline
- Makes delivery less consistent across teams
- Hides an important coordination and quality gate

### Option B: Make pull request creation and review explicit ← CHOSEN

Add PR creation and PR review as explicit steps inside Deliver.

**Pros:**
- Makes review a standard part of phased delivery
- Aligns with BMAD review roles
- Fits OpenSpec's preference for reviewable change sets
- Reinforces Speckit's small, verifiable tasks

**Cons:**
- Adds visible process to the Deliver phase
- Requires clearer role expectations for reviewers

## Decision

Make pull request creation and review explicit steps inside Deliver.

The default Deliver flow becomes:

- Build
- Create PR
- Review PR
- Verify
- Deploy
- Archive

The PR should normally map to one delivery slice or a small set of tightly
related tasks.

## Consequences

- Plan should produce slices that are small enough to review safely
- Deliver must define PR owner, reviewers, and review expectations
- Review feedback should update code and artifacts before deploy
- Archive should happen only after PR review and delivery closure are complete
