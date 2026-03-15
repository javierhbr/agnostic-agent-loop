# Alignment and Traceability Rules

Use these rules across all phases when the platform has a master repository and
components have their own repositories.

## Goal

Keep shared platform truth canonical, keep component work aligned by version,
and keep JIRA traceability linked to the same change chain.

## Must do

- treat the platform master repository as the source of shared platform truth
- treat component repositories as the source of local implementation truth
- record platform version and platform refs in each affected component repository
- keep JIRA as the coordination layer, not the detailed spec store
- create linked platform and component changes when shared truth changes
- map tasks to stories and PRs to stories during planning and delivery

## Avoid

- copying editable platform truth into component repositories
- updating shared truth only in a component repository
- putting the full requirements or design into JIRA descriptions
- opening delivery work with no platform version or issue chain

## Required outputs

- `platform-ref.yaml`
- `jira-traceability.yaml`
- linked platform issue and component epic
- linked tasks, stories, and PRs
