# Artifact Authoring Guide

## proposal.md
Good proposals answer:
- what problem exists
- why now
- what changes
- what is explicitly out of scope
- what risks or rollback concerns exist

## Delta specs
Good delta specs:
- describe behavior, not implementation
- make the changed requirement easy to compare with current reality
- use specific scenarios
- keep ADDED, MODIFIED, and REMOVED sections explicit

## design.md
Good design documents:
- explain the chosen approach
- note tradeoffs and rejected alternatives
- identify affected systems and interfaces
- call out rollout, migration, observability, or failure handling when relevant

## tasks.md
Good tasks are:
- implementation-ready
- ordered logically
- small enough to complete and verify
- traceable to requirements and design

## Writing style for all artifacts
- prefer precise language over aspirational wording
- avoid ambiguity like “handle better” or “support properly”
- keep artifacts easy for another agent to continue from
