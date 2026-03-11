# How to Use the Component Boilerplate

Use this boilerplate when you want a component repository to align to platform
truth and still keep local OpenSpec artifacts in the component repo.

## Quick start

1. copy `component-boilerplate/` into the target component repo
2. fill `openspec/config.yaml` with the local component context
3. fill `platform-ref.yaml` with the platform version and refs
4. fill `jira-traceability.yaml` with the issue, epic, and story chain
5. write `proposal.md`
6. write the first spec delta under `specs/`
7. add `design.md` and `tasks.md` when the change reaches planning
8. use `pr-description.md` as the PR traceability pattern

## Minimum files to fill first

Start with these files:

- `platform-ref.yaml`
- `jira-traceability.yaml`
- `proposal.md`
- `specs/example-change/spec.md`

Then add:

- `design.md`
- `tasks.md`
- `pr-description.md`

## What the first version should answer

The first component setup should answer:

- which platform version the component follows
- which platform refs constrain the change
- which JIRA issue chain owns the work
- what local behavior is changing
- how tasks map to stories and PRs

## Small sample

See:

- `sample-component-package.md`
- `../../example/component-repo/platform-ref.yaml`
- `../../example/component-repo/jira-traceability.yaml`
