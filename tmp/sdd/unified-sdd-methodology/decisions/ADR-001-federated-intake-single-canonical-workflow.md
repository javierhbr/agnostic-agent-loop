---
id: "methodology/ADR-001"
title: "Allow federated intake with a single canonical workflow"
status: "proposed"
date: "2026-03-07"
initiative: "unified-sdd-methodology"
author: "javierbenavides + codex"
deciders:
  - "platform-methodology-working-group"
---

## Context

The platform serves multiple teams and artifacts. New work can begin from
different places:

- a platform initiative
- a product requirement
- a component or team proposal

If we force one intake path, we create friction for teams and slow small
changes. If we allow each entry point to define its own workflow, teams will
create duplicate artifacts, drift in quality, and lose traceability.

## Options Considered

### Option A: Platform-first intake only

Require every change to start as a platform initiative.

**Pros:**
- Strong central governance
- Clear portfolio visibility
- Easier executive reporting

**Cons:**
- Slows local or low-risk work
- Creates a bottleneck around central approval
- Pushes teams to work around the process

### Option B: Team-first intake only

Allow each team to start from its own proposal and link upward later.

**Pros:**
- Fast for local delivery
- Low coordination overhead for small changes

**Cons:**
- Weak cross-team alignment
- Duplicate specs and conflicting assumptions
- Harder to govern contracts and shared standards

### Option C: Federated intake with normalization ← CHOSEN

Allow work to start from any approved entry point, but normalize every
change into one canonical workflow before implementation.

**Pros:**
- Preserves team speed
- Supports platform and product planning
- Creates one traceable path from intent to implementation
- Fits multi-team, multi-artifact delivery

**Cons:**
- Requires a routing layer
- Needs clear ownership for normalization and status changes

## Decision

Use a federated intake model with a single canonical workflow.

The workflow may start from a platform initiative, a product requirement, or
a component proposal, but all approved work must be normalized into the same
delivery model before planning and implementation proceed.

## Consequences

- We need an intake router that classifies scope, risk, and affected teams
- We need one canonical execution unit for every approved change
- Entry artifacts remain valid, but none of them become the sole source of truth
- Traceability must link intake artifacts to specs, plans, tasks, contracts,
  and implementation
