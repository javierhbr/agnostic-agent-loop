---
id: "methodology/ADR-002"
title: "Use the change package as the canonical execution unit"
status: "proposed"
date: "2026-03-07"
initiative: "unified-sdd-methodology"
author: "javierbenavides + claude"
deciders:
  - "platform-methodology-working-group"
---

## Context

We need one canonical unit that every intake path can normalize into. The
choice must work across many teams, multiple artifacts, and both platform and
component-level changes.

The main candidates are:

- initiative
- capability spec
- change package

## Options Considered

### Option A: Initiative

Make the initiative the canonical unit.

**Pros:**
- Strong business framing
- Good for portfolio management
- Works well for large, cross-team programs

**Cons:**
- Too large for day-to-day execution
- Poor fit for small or local changes
- Hard to assign implementation-ready tasks directly

### Option B: Capability spec

Make the capability spec the canonical unit.

**Pros:**
- Durable platform source of truth
- Useful for long-lived platform boundaries and standards
- Encourages reuse across initiatives

**Cons:**
- Too stable for many short-term delivery changes
- Weak fit for sequencing, rollout, and task management
- Easy to overload with temporary delivery detail

### Option C: Change package ← CHOSEN

Use the change package as the canonical execution unit.

**Pros:**
- Small enough to execute
- Large enough to contain proposal, deltas, design, tasks, and validation
- Works for both single-team and cross-team changes
- Maps well to agent execution and status tracking

**Cons:**
- Needs links upward to initiatives and capability specs
- Requires discipline to avoid oversized packages

## Decision

Use the change package as the canonical execution unit in the unified
methodology.

Initiatives remain the strategic umbrella. Capability specs remain the durable
platform source of truth. The change package becomes the standard container
for planning, execution, validation, and status tracking.

## Consequences

- Every approved change must have one change package
- Each change package must link to its parent initiative when one exists
- Each change package must link to affected capability specs, components,
  contracts, and ADRs
- Delivery status, tasks, validation, and archive history should attach to the
  change package rather than to the initiative alone
