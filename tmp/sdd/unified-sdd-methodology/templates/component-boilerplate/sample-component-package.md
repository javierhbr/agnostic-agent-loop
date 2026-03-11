# Sample Component Package

Component: `profile-service`
Platform version: `2026.03`
Platform issue: `PLAT-123`
Component epic: `PROF-456`
Change package: `chg-profile-email-validation`

## Core files

### `platform-ref.yaml`

- platform ref: `capabilities.customer-identity`
- contract ref: `contracts.customer-profile.v2`
- alignment type: `shared-change`

### `jira-traceability.yaml`

- platform issue: `PLAT-123`
- component epic: `PROF-456`
- story examples:
  - `PROF-789`
  - `PROF-790`

### `proposal.md`

Problem:

- email validation and failure handling are inconsistent today

Goals:

- validate before persistence
- keep explicit customer-facing errors
- preserve shared contract compatibility

### `design.md`

Main decisions:

- validate format before write
- keep contract version `v2`
- emit validation metrics and logs

### `tasks.md`

- Task 1 -> `PROF-789`
- Task 2 -> `PROF-790`

### `pr-description.md`

Traceability:

- platform issue
- component epic
- story
- change package

## Sample outcome

This component package shows the minimum chain teams should keep aligned:

```text
[Platform version 2026.03]
        |
        v
[platform-ref.yaml]
        |
        v
[proposal.md + spec.md]
        |
        v
[design.md + tasks.md]
        |
        v
[story] -> [PR]
```
