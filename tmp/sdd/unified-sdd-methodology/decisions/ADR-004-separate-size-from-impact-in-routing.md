---
id: "methodology/ADR-004"
title: "Separate change size from change impact in workflow routing"
status: "proposed"
date: "2026-03-08"
initiative: "unified-sdd-methodology"
author: "javierbenavides + claude"
deciders:
  - "platform-methodology-working-group"
---

## Context

The current workflow draft uses a single classification step for "change size
and impact." That label is easy to read, but it mixes two different concerns:

- delivery effort
- business and technical risk

If we treat them as the same thing, teams may over-plan low-risk large changes
or under-control small but high-risk changes.

## Options Considered

### Option A: One combined classification only

Use one label and one score for both size and impact.

**Pros:**
- Simple to explain
- Fast to apply

**Cons:**
- Hides important differences
- Leads to inconsistent planning depth
- Leads to inconsistent validation and release controls

### Option B: Separate size and impact ← CHOSEN

Classify each change on two axes:

- size = delivery effort and coordination depth
- impact = risk, blast radius, and control depth

**Pros:**
- Better routing decisions
- Better fit for both small high-risk and large low-risk changes
- Clearer ownership for planning and validation

**Cons:**
- Slightly more process at intake
- Requires teams to learn two terms instead of one

## Decision

Separate change size from change impact in the routing model.

The methodology keeps one core workflow for all changes:

- intake
- classify
- change package
- proposal and specs
- clarify
- design
- tasks
- implement
- verify
- archive

Size changes the workflow depth. Impact changes the governance and validation
depth.

## Consequences

- Small changes use compact planning artifacts
- Medium changes use the standard path
- Large changes use phased planning and delivery
- Low-impact changes use normal controls
- High-impact changes require stronger review, validation, and rollout controls
- The proposal should define default thresholds for size and impact
