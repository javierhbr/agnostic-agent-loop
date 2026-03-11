---
id: "methodology/ADR-013"
title: "Implement the first local platform MCP server as a self-contained Go stdio server"
status: "proposed"
date: "2026-03-08"
initiative: "unified-sdd-methodology"
author: "javierbenavides + codex"
deciders:
  - "platform-methodology-working-group"
---

## Context

ADR-012 established that teams should use a local read-only platform MCP
gateway before hosted infrastructure exists.

The next decision is the first implementation shape.

The first server must:

- run locally on developer machines
- stay small and fast
- avoid runtime drift across teams
- work with stdio-based MCP clients
- remain read-only
- support pinned-version validation against `platform-ref.yaml`

## Options Considered

### Option A: Node.js server

**Pros:**
- fast to prototype
- rich MCP ecosystem

**Cons:**
- higher local runtime drift
- larger support burden across developer machines
- weaker fit for a small, single-binary local tool

### Option B: Python server

**Pros:**
- fast to prototype
- easy local scripting

**Cons:**
- environment drift
- packaging friction
- weaker binary distribution story

### Option C: Rust stdio server

**Pros:**
- small binary
- strong performance
- strong safety guarantees

**Cons:**
- slower initial delivery for this team context
- higher implementation cost for a first internal tool

### Option D: Go stdio server with no external Go dependencies ← CHOSEN

**Pros:**
- small, portable binary
- fast startup
- easy local distribution
- faster initial delivery than Rust
- easy offline build and test

**Cons:**
- requires a small custom implementation of the MCP protocol surface
- less library support than a JavaScript-first implementation

## Decision

Implement the first server as:

- a Go module
- stdio transport
- newline-delimited JSON-RPC messages
- read-only only
- no external Go module dependencies in v1

The first tool set should stay small and focus on:

- platform version lookup
- platform ref discovery
- JIRA-linked metadata lookup
- component alignment validation
- pinned-version drift detection

## Consequences

- the first server can be built and tested offline
- the server is easy for developers to run locally
- the implementation is intentionally narrow and should not become a write API
- later hosted versions can keep the same tool contract while changing the
  deployment model
