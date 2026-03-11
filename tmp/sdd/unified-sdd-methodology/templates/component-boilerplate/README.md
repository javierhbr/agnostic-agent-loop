# Component Boilerplate

Use this boilerplate for a component repository that implements local work while aligning to platform truth.

## Recommended structure

```text
component-boilerplate/
  openspec/
    config.yaml
  platform-ref.yaml
  jira-traceability.yaml
  proposal.md
  design.md
  tasks.md
  pr-description.md
  specs/
    example-change/
      spec.md
```

## What belongs here

- local OpenSpec artifacts
- local platform alignment metadata
- JIRA traceability metadata
- local PR and validation notes

## What should not live here

- editable copies of full platform truth
- platform-wide ADRs unless this repo owns them
- delivery work that has no platform version or issue chain when one is required

## Sample

See:

- `how-to-use.md`
- `sample-component-package.md`
- `../../example/component-repo/platform-ref.yaml`
- `../../example/component-repo/jira-traceability.yaml`
- `../../example/component-repo/proposal.md`
- `../../example/component-repo/design.md`
- `../../example/component-repo/tasks.md`
- `../../example/component-repo/specs/validated-email-updates/spec.md`
