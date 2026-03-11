# Unified SDD Example Package

This directory contains the concrete examples for the unified SDD methodology.

Use it when teams need to see the method applied end to end with:

- one platform capability
- multiple component repositories
- JIRA issue chain
- concrete OpenSpec-style artifacts

## Running example

Shared feature:

- validated customer email updates across the platform

Shared IDs:

- platform capability: `customer-identity`
- shared contract: `contracts.customer-profile.v2`
- platform issue: `PLAT-123`
- profile component epic: `PROF-456`
- auth component epic: `AUTH-234`

Affected components:

- `profile-service`
- `auth-service`
- `notification-service`

## Example flow

```text
[Platform baseline]
        |
        v
[Route]
  PLAT-123 + PROF-456 + AUTH-234
        |
        v
[Specify]
  proposal + delta specs + platform refs
        |
        v
[Plan]
  design + tasks + story mapping
        |
        v
[Deliver]
  stories -> PRs -> verify -> deploy -> archive
```

## Documents in this directory

- `01-platform-phase.md`
- `02-route-phase.md`
- `03-specify-phase.md`
- `04-plan-phase.md`
- `05-deliver-phase.md`
- `06-entry-point-examples.md`
- `07-platform-component-interaction.md`
- `08-platform-plan-to-component-openspec.md`
- `09-local-platform-mcp-usage.md`

## Concrete artifacts in this directory

Platform repository examples:

- `platform-repo/platform-baseline.md`

Component repository examples:

- `component-repo/platform-ref.yaml`
- `component-repo/jira-traceability.yaml`
- `component-repo/proposal.md`
- `component-repo/design.md`
- `component-repo/tasks.md`
- `component-repo/pr-description.md`
- `component-repo/specs/validated-email-updates/spec.md`

## Related templates

If you want reusable starting points instead of filled examples, use:

- `../templates/README.md`
- `../templates/platform-template/README.md`
- `../templates/component-boilerplate/README.md`
- `../templates/platform-mcp-boilerplate/README.md`

## How to use this package

1. read `01-platform-phase.md` through `05-deliver-phase.md` in order
2. use `06-entry-point-examples.md` when deciding how a new request starts
3. use `07-platform-component-interaction.md` when teams need to align local work to platform truth
4. use `08-platform-plan-to-component-openspec.md` when a component team is receiving a platform `Plan` handoff
5. use `09-local-platform-mcp-usage.md` when teams are consuming platform truth through a local MCP server
6. open the artifact examples when you need a concrete target shape
