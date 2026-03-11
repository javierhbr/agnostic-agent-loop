# Platform Template

Use this template for the master platform repository.

It is the starting point for canonical platform truth.

## Recommended structure

```text
platform-template/
  platform-baseline.md
  refs-index.md
  jira-conventions.md
  capabilities/
    capability-template.md
  contracts/
    contract-template.md
  adrs/
    ADR-000-template.md
```

## What belongs here

- platform principles and guardrails
- shared capabilities
- shared contracts and schemas
- platform-level ADRs
- JIRA hierarchy conventions
- versioning and durable ref rules

## What should not live here

- component-only implementation detail
- local task lists
- component PR descriptions
- copied component specs

## Sample

See:

- `how-to-use.md`
- `sample-platform-baseline.md`
- `../../example/platform-repo/platform-baseline.md`
- `../../example/07-platform-component-interaction.md`
