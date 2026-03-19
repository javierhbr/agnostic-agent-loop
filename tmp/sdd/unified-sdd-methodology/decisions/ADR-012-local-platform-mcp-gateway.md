---
id: "methodology/ADR-012"
title: "Adopt a local read-only platform MCP gateway for developer access to platform truth"
status: "proposed"
date: "2026-03-08"
initiative: "unified-sdd-methodology"
author: "javierbenavides + claude"
deciders:
  - "platform-methodology-working-group"
---

## Context

The methodology now assumes:

- one canonical platform repository for shared truth
- multiple component repositories for local OpenSpec work
- multiple teams that need local, fast access to platform rules, refs,
  contracts, and issue mapping

A hosted MCP server is not available yet. Teams still need a practical way to:

- query platform refs and contracts locally
- validate component alignment against the pinned platform version
- inspect shared planning decisions during component planning and delivery
- include JIRA context without making JIRA the spec source of truth

Without a local access layer, teams will tend to:

- browse docs manually and inconsistently
- use stale links or tribal knowledge
- copy platform context into component repositories
- answer alignment questions differently across developers

## Options Considered

### Option A: No MCP layer, rely on manual docs and repo browsing

**Pros:**
- no extra tooling
- lowest setup cost

**Cons:**
- weak consistency
- poor discoverability
- repeated manual work
- no standard local validation path

### Option B: Copy or mirror platform truth into component repositories

**Pros:**
- easy local access inside each repo

**Cons:**
- high drift risk
- duplicates the source of truth
- makes local copies look authoritative

### Option C: Local read-only MCP gateway backed by a developer's local platform clone ← CHOSEN

Each developer keeps a local clone of the platform repository. A small local,
read-only MCP server runs against that clone and exposes query and validation
capabilities to component work.

**Pros:**
- works without hosted infrastructure
- preserves one platform source of truth
- fast local access for every developer
- supports future migration to a hosted MCP server with the same tool surface

**Cons:**
- requires local setup and update discipline
- needs version-aware behavior so components do not validate against the wrong
  platform baseline

## Decision

Adopt a local read-only platform MCP gateway model.

Use this operating rule:

- developers may keep the platform clone updated locally
- component work must default to the platform version pinned in
  `platform-ref.yaml`, not blindly use the latest platform branch state
- the MCP server is a read-only gateway to platform truth, not a second source
  of truth and not a write API for either repository

The local MCP server should read both:

- platform specs, refs, contracts, ADRs, and version markers
- JIRA mapping and issue references that are already represented in local
  platform artifacts or synced metadata

## Consequences

- Platform phase should define the local MCP usage model and configuration
- component teams can query platform truth locally without copying it into the
  component repo
- Plan and Deliver can validate against the pinned platform version through the
  local MCP server
- the first MCP version should stay small, local, and read-only
- the same tool contract should be reusable later if the MCP server is hosted
