# How to Use the Platform Template

Use this template when you want to start or refresh the master platform
repository.

## Quick start

1. copy `platform-template/` into the target platform repo or docs area
2. replace the placeholders in `platform-baseline.md`
3. define the stable refs in `refs-index.md`
4. add the first shared capability under `capabilities/`
5. add the first shared contract under `contracts/`
6. define the JIRA conventions in `jira-conventions.md`
7. create ADRs under `adrs/` when a shared technical decision is needed

## Minimum files to fill first

Start with these files:

- `platform-baseline.md`
- `refs-index.md`
- `jira-conventions.md`

Then add:

- one capability file
- one contract file

## What the first version should answer

The first platform baseline should answer:

- what shared truth this platform owns
- which capability refs are stable
- which shared contracts matter first
- how components pin platform version and refs
- how JIRA issues, epics, and stories are linked

## Small sample

See:

- `sample-platform-baseline.md`
- `../../example/platform-repo/platform-baseline.md`
