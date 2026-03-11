# Sample Platform Baseline

Platform: `customer-platform`
Version: `2026.03`

## Purpose

This platform baseline governs shared truth for customer identity behavior
across profile, auth, and notification components.

## Shared principles

- customer identity changes must be validated before persistence
- shared contracts must stay backward compatible unless a versioned migration is approved
- identity-related failures must be explicit and observable

## Shared capabilities

- `capabilities.customer-identity`

## Shared contracts

- `contracts.customer-profile.v2`

## Versioning rules

- platform truth is published as `yyyy.mm`
- component repos pin the platform version in `platform-ref.yaml`
- shared contract changes require explicit version or explicit compatibility approval

## JIRA conventions

- one platform issue for shared identity change
- one component epic per affected repo
- one story per reviewable slice

## Alignment rules

- platform repo owns shared truth
- component repos own local implementation truth
- JIRA tracks workflow state and ownership
