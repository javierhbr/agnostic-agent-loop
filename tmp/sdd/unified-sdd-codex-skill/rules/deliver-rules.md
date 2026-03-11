# Deliver Rules

Use these rules in the Deliver phase.

## Goal

Execute the plan in controlled slices through code, review, verification,
deploy, and archive.

## Apply these skills throughout delivery

- `bmad-codex-skill` for implementation and review support
- `openspec-codex-skill` for apply and archive
- `speckit-codex-skill` for task discipline and phased delivery

## Default slice flow

1. Build
2. Create PR
3. Review PR
4. Verify
5. Deploy
6. Archive

## Must do

- deliver in slices, not one large batch
- keep task status and artifacts current during execution
- create a PR for each reviewable slice or tightly related slice set
- keep story, epic, and platform issue links current
- collect validation evidence before closure
- capture deploy and rollback notes when relevant
- archive the change package after delivery is complete

## Avoid

- stale tasks or stale specs during execution
- delaying review until everything is finished
- deploying without resolved or tracked review feedback
- treating archive as optional cleanup

## Required outputs

- implemented slice
- PR and review results
- validation evidence
- deploy notes when needed
- archived change package
- final traceability from platform issue to story to PR

## Exit gate

Close the change only when:

- slices are complete or intentionally deferred
- review is complete
- validation evidence exists
- artifacts reflect delivered reality
- archive is done
