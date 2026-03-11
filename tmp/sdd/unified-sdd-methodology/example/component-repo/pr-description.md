# PR Description Example

Title: `Validate customer email updates in profile-service`

## Traceability

- platform issue: `PLAT-123`
- component epic: `PROF-456`
- story: `PROF-789`
- change package: `chg-profile-email-validation`

## Scope

- add email format validation before persistence
- add explicit duplicate-email failure handling
- add observability for validation failures

## Platform refs

- `capabilities.customer-identity`
- `contracts.customer-profile.v2`
- `principles.observability`

## Validation

- added API tests for invalid email format
- added API tests for duplicate email attempts
- confirmed event contract remains unchanged

## Reviewer focus

- verify backward compatibility with `contracts.customer-profile.v2`
- verify validation errors remain explicit and recoverable
