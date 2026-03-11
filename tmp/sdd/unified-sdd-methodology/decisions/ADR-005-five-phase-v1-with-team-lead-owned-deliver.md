---
id: "methodology/ADR-005"
title: "Adopt a 5-phase v1 workflow with Team Lead owning Deliver"
status: "proposed"
date: "2026-03-08"
initiative: "unified-sdd-methodology"
author: "javierbenavides + codex"
deciders:
  - "platform-methodology-working-group"
---

## Context

The current draft explains the methodology mostly through artifacts and stage
steps. That is accurate, but it is harder to teach from platform strategy down
to engineering execution and deploy.

We need a simpler operating model that:

- starts at the platform level
- shows how work moves into development and deploy
- explains how agents interact with teams
- shows which of the three SDD skills apply in each phase
- shows which rule sets govern each phase

## Options Considered

### Option A: Keep the detailed artifact-first flow

Keep proposal, specs, clarify, design, tasks, implement, verify, and archive
as the main visible workflow.

**Pros:**
- Very explicit
- Maps directly to the underlying artifacts

**Cons:**
- Harder to teach across teams
- Makes phased delivery harder to see
- Puts too much process detail at the top level

### Option B: Move directly to 6 phases

Use Platform, Route, Specify, Plan, Build, and Deploy.

**Pros:**
- Very strong separation between engineering and release
- Good long-term target for mature programs

**Cons:**
- More ceremony for initial adoption
- Adds a phase boundary many teams may not need yet

### Option C: Use 5 phases in v1 ← CHOSEN

Use Platform, Route, Specify, Plan, and Deliver. Keep build, verify, deploy,
and archive inside Deliver as internal slices.

**Pros:**
- Easier to adopt
- Still preserves phased implementation
- Maps cleanly to BMAD, OpenSpec, and Speckit
- Leaves room to split Deliver later if release complexity grows

**Cons:**
- Build and deploy are less explicit at the top level
- Requires a clear Deliver owner

## Decision

Adopt a 5-phase v1 workflow:

- Platform
- Route
- Specify
- Plan
- Deliver

In v1, Deliver is owned by the Team Lead.

Deliver contains four internal slices:

- build
- verify
- deploy
- archive

The methodology may evolve later to 6 phases by splitting Deliver into Build
and Deploy when release complexity justifies that change.

## Consequences

- The team proposal should explain the methodology through phases first
- Artifacts remain important, but they live inside phases
- Each phase must show agent interaction, skill usage, and applied rules
- Team Lead becomes the accountable coordinator from implementation through
  release and archive in v1
