---
id: "methodology/ADR-008"
title: "Use role-specific agents on top of the unified SDD skill package"
status: "proposed"
date: "2026-03-08"
initiative: "unified-sdd-methodology"
author: "javierbenavides + codex"
deciders:
  - "platform-methodology-working-group"
---

## Context

The unified SDD skill package provides one shared workflow, but platform teams
still work through role-specific responsibilities.

Architects, Team Leads, Product partners, and Developers all need access to
the same three source skills. They should not use them in the same way.

## Options Considered

### Option A: One generic orchestrator only

Use only the unified SDD skill and expect every role to adapt it manually.

**Pros:**
- Simple package structure
- One place to maintain

**Cons:**
- Less practical for teams
- Harder to onboard role-by-role
- Prompts and skill usage remain too generic

### Option B: One agent per role on top of the unified skill ← CHOSEN

Create role-specific agents that all reuse the same unified methodology and
source skills.

**Pros:**
- Matches how platform teams actually work
- Keeps one shared method with role-specific behavior
- Makes prompts, responsibilities, and handoffs clearer

**Cons:**
- More files to maintain
- Requires the role agents to stay aligned with the unified skill

## Decision

Create one agent per role on top of the unified SDD skill package.

The initial role agents are:

- Architect
- Team Lead
- Product
- Developer

Each agent uses the same three source skills:

- BMAD
- OpenSpec
- Speckit

But each role applies them differently according to its phase ownership and
responsibilities.

## Consequences

- The unified skill remains the shared operating model
- Role agents become the practical entry points for team usage
- Each role agent must define:
  - mission
  - primary phases
  - skill emphasis
  - role interactions
  - outputs
  - prompt examples
