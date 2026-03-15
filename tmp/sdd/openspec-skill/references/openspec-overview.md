# OpenSpec Overview

## What OpenSpec is for
OpenSpec adds a lightweight spec layer before implementation so people and AI assistants can agree on what to build before code changes start.

## Core idea
A change is represented as a folder containing planning artifacts. The usual artifact chain is:

`proposal -> specs -> design -> tasks -> implementation`

This is not a one-way waterfall. OpenSpec expects teams to revise earlier artifacts as they learn.

## Default artifacts
- `proposal.md` — why the change exists, scope, approach, and risk framing
- `specs/` — delta specs that describe behavior changes
- `design.md` — technical approach and architecture decisions
- `tasks.md` — ordered implementation checklist

## Delta specs
Delta specs describe how current behavior should change.
Typical change sections:
- ADDED
- MODIFIED
- REMOVED

Archive merges accepted deltas into the main specs and moves the completed change into archive history.

## Profiles and workflow shape
The default `core` profile centers on:
- propose
- explore
- apply
- archive

Expanded workflows can add commands like `new`, `continue`, `ff`, `verify`, `sync`, `bulk-archive`, and `onboard`.
