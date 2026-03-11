---
id: "methodology/ADR-010"
title: "Adopt canonical platform truth with versioned component alignment and JIRA-linked execution"
status: "proposed"
date: "2026-03-08"
initiative: "unified-sdd-methodology"
author: "javierbenavides + codex"
deciders:
  - "platform-methodology-working-group"
---

## Context

The methodology now spans:

- one platform-level source of shared rules and contracts
- multiple component repositories with local implementation detail
- multiple teams delivering through JIRA epics, stories, pull requests, and releases

Without an explicit interaction model, teams will drift in at least one of
these ways:

- copy platform rules into component repositories and let them diverge
- put too much implementation detail into the platform truth
- use JIRA as a duplicate requirements store
- lose traceability between platform-level intent and component-level delivery

The methodology needs a clear relationship between:

- platform specs
- component specs
- JIRA hierarchy
- local OpenSpec artifacts

## Options Considered

### Option A: Full sync of platform specs into every component repository

Copy or mirror platform specs into every component repository and update them
locally as part of the component workflow.

**Pros:**
- easy local access
- no external reference lookup

**Cons:**
- high drift risk
- unclear source of truth
- duplicated maintenance
- component teams can silently redefine platform truth

### Option B: Manual links only

Keep a platform repository and component repositories separate, but rely mostly
on URLs, wiki links, and issue descriptions for alignment.

**Pros:**
- simple to start
- little tooling required

**Cons:**
- weak traceability
- no durable version pinning
- easy for teams to forget which platform baseline they followed

### Option C: Canonical platform truth with versioned component alignment and JIRA-linked execution ← CHOSEN

Use one master platform repository as the canonical source of shared truth.
Each component repository keeps local OpenSpec artifacts and explicitly pins
the platform version and references it aligns to. JIRA tracks delivery through
platform-level issues, component epics, and stories.

**Pros:**
- one upstream source of shared truth
- clean separation of platform and component responsibilities
- clear versioned traceability
- supports multiple repositories and teams without forcing a monorepo
- fits the change package model and the phased methodology

**Cons:**
- requires simple alignment metadata
- requires teams to keep JIRA links and spec refs current
- needs clear rules for when a component change also requires a platform change

## Decision

Adopt this model:

- the platform master repository is the canonical source of shared platform truth
- component repositories own local OpenSpec artifacts and implementation detail
- each component repository records the platform version and references it aligns to
- JIRA is the coordination and tracking layer, not the spec source of truth

Use the following hierarchy where possible:

- platform initiative or platform epic for shared platform-level outcomes
- component epic per affected component or repository
- story per reviewable delivery slice or task group

If JIRA does not support a hierarchy above epic, use:

- one platform epic
- linked component epics
- stories under each component epic

Shared changes must update:

- the platform truth in the master repository
- the affected component specs in the component repositories
- the linked JIRA issues

## Consequences

- Platform phase must define the master repo baseline, versioning approach, and shared ref IDs
- Route must classify changes as platform-only, component-only, or shared
- Specify must keep component specs aligned to platform refs and open linked platform deltas when shared truth changes
- Plan must trace tasks to stories and keep platform refs visible in design and task artifacts
- Deliver must link PRs, stories, verification, and archive records back to the same aligned change

The methodology should provide:

- a detailed explanation of the model
- templates for `platform-ref.yaml` and `jira-traceability.yaml`
- skill and rule updates so agents apply the model consistently
