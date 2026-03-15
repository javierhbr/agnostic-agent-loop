# Commands and Workflows

## Default quick path
Use this path for most teams:

1. `/opsx:propose`
2. `/opsx:apply`
3. `/opsx:archive`

Use `/opsx:explore` before proposing when the problem is still fuzzy.

## Expanded workflow commands
These depend on the project's configured workflow profile.

- `/opsx:new` — scaffold a change without generating everything at once
- `/opsx:continue` — generate the next artifact
- `/opsx:ff` — fast-forward planning artifacts
- `/opsx:verify` — run workflow-specific verification steps
- `/opsx:sync` — sync workflow state when configured
- `/opsx:bulk-archive` — archive multiple completed changes
- `/opsx:onboard` — project onboarding support when available

## How to choose
- Use `explore` for discovery and comparison.
- Use `propose` for the standard all-at-once planning jumpstart.
- Use `new/continue/ff` only when the team wants more explicit control over artifact creation.
- Use `apply` once tasks are good enough to implement.
- Use `archive` when the change is truly finished and the specs should become the new source of truth.
