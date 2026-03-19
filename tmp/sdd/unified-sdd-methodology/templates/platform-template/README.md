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
  ownership/
    component-ownership-template.md   <- one per component
    dependency-map-template.md        <- one per platform
    glossary-template.md              <- one per platform
```

## What belongs here

- platform principles and guardrails
- shared capabilities
- shared contracts and schemas
- platform-level ADRs
- JIRA hierarchy conventions
- versioning and durable ref rules
- component ownership boundaries
- dependency map with impact tiers
- shared glossary

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
