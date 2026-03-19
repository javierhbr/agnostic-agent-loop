---
id: "methodology/ADR-003"
title: "Use stage-based role accountability for each change package"
status: "proposed"
date: "2026-03-07"
initiative: "unified-sdd-methodology"
author: "javierbenavides + claude"
deciders:
  - "platform-methodology-working-group"
---

## Context

The platform spans 5 teams, about 20 experiences, and about 60 customer
flows. A unified methodology needs clear ownership at each stage of a change
package without creating a process that teams will ignore.

We need a responsibility model that works for:

- intake and prioritization
- planning and design
- implementation and validation
- cross-team coordination

## Options Considered

### Option A: One end-to-end owner for the whole change

Assign one person to own the entire change package from intake to archive.

**Pros:**
- Simple to explain
- Clear escalation point

**Cons:**
- Unrealistic for cross-team work
- Blurs business, design, and engineering responsibilities
- Creates bottlenecks around one role

### Option B: Full RACI matrix for every artifact

Define responsible, accountable, consulted, and informed roles for each
artifact and task.

**Pros:**
- Very explicit
- Strong audit trail

**Cons:**
- Heavy to maintain
- Hard to apply in day-to-day delivery
- Too much ceremony for small changes

### Option C: Stage-based accountability ← CHOSEN

Assign one primary accountable role to each workflow stage and list the key
supporting roles for that stage.

**Pros:**
- Clear ownership without heavy process
- Fits both small and large changes
- Matches BMAD role-based handoffs
- Easier for teams and agents to follow

**Cons:**
- Requires discipline at handoff points
- Some edge cases still need local agreement

## Decision

Use a stage-based accountability model for the unified methodology.

Each change package stage has one primary accountable role. Supporting roles
contribute based on the stage. The default roles are:

- Product
- Team Lead
- Architect
- Engineering Manager
- Developers
- QA / Validation

## Consequences

- Every stage in the methodology must show a primary owner
- Team leads coordinate the flow of the change package across stages
- Product owns business intent and acceptance scope
- Architects own design integrity, contracts, and ADR-quality decisions
- Developers own implementation and test execution
- Engineering managers own staffing, prioritization support, and escalation
- QA / validation owns quality evidence and release readiness checks
