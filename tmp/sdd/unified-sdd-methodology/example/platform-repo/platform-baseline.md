# Platform Baseline Example

Platform: `customer-platform`
Version: `2026.03`
Capability: `customer-identity`

## Shared principles

- customer identity changes must be validated before persistence
- shared contracts must remain backward compatible unless a versioned migration is approved
- identity-related events must emit standard observability fields
- customer-facing validation failures must be explicit and recoverable

## Shared refs

- `capabilities.customer-identity`
- `contracts.customer-profile.v2`
- `principles.api-versioning`
- `principles.observability`
- `principles.backward-compatibility`

## JIRA conventions

- one platform issue for shared capability change
- one component epic per affected repo
- one story per reviewable slice
- every PR should reference at least one story key

## Alignment rule

- platform repo owns shared truth
- component repos pin the platform version and refs they align to
- local OpenSpec artifacts stay in the component repo
