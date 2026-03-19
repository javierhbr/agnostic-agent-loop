# Component Ownership Boundary

Component: `<component-name>`
Platform version: `<yyyy.mm>`

## What this component owns

List the things this component is the single source of truth for.

- `<capability or data domain 1>`
- `<capability or data domain 2>`
- `<capability or data domain 3>`

## What this component does NOT own

List things closely related but owned by another component. This prevents scope drift.

- `<related thing>` — owned by `<other-component>`
- `<related thing>` — owned by `<other-component>`

## Shared contracts this component publishes

List contracts or events this component exposes to others.

- `<contract-ref>` — `<brief description>`

## Shared contracts this component consumes

List contracts or events this component depends on from others.

- `<contract-ref>` — from `<other-component>`

## Notes

Any boundary decisions that need explanation for the team.

---

## How to use this file

- Fill in one file per component when the Platform phase defines the baseline.
- Update when ownership boundaries change (requires platform change package).
- Reference this file in `platform-ref.yaml` under `ownership.primary_component`.
- The Team Lead reads this during Assess to classify whether a change is local-only
  or crosses an ownership boundary.
