---
id: "methodology/ADR-009"
title: "Create a Platform Contextualizer skill for starting Iteration 1 on existing platforms"
status: "proposed"
date: "2026-03-08"
initiative: "unified-sdd-methodology"
author: "javierbenavides + codex"
deciders:
  - "platform-methodology-working-group"
---

## Context

The Platform phase in Iteration 1 assumes teams can define durable context and
shared rules. On an existing platform with active teams, that work is harder
than it looks.

The platform already has:

- current conventions
- undocumented constraints
- existing artifacts
- team habits
- gaps between how teams think they work and how they actually work

Teams need help reviewing and documenting the current state before they try to
define the target baseline.

## Options Considered

### Option A: Start Platform phase directly from a blank baseline

Begin the Platform phase by writing principles and config without a dedicated
current-state contextualization step.

**Pros:**
- Faster to begin
- Fewer documents

**Cons:**
- High risk of generic rules
- Likely to miss real constraints and team practices
- Easier to produce principles teams will ignore

### Option B: Create a Platform Contextualizer skill ← CHOSEN

Create a starter skill focused on existing-platform discovery, gap analysis,
and durable-context drafting.

**Pros:**
- Fits brownfield reality
- Grounds the Platform phase in evidence
- Helps teams separate current state, gaps, and target baseline
- Makes the Platform phase more practical and actionable

**Cons:**
- Adds one more skill to maintain
- Requires teams to spend time on discovery before standard-setting

## Decision

Create a `Platform Contextualizer` skill under `custom-skills/sdd`.

The skill should help teams:

- review and document the current platform state
- identify gaps and areas for improvement
- establish shared understanding of principles and constraints
- produce actionable outputs for the Platform phase

It should combine:

- BMAD for brownfield context, project-context thinking, and role framing
- OpenSpec for durable context and reusable config
- Speckit for explicit principles and quality guardrails

## Consequences

- The Platform phase gets a practical starting point for existing platforms
- Teams can document current state before proposing durable rules
- The unified methodology should link to this starter skill from the Platform phase
