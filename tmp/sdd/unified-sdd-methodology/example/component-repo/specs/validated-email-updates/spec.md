# Spec Delta: Validated Email Updates

## ADDED Email Validation

The system must validate customer email updates before persistence.

### Scenario: invalid email format

- given a customer submits an invalid email format
- when profile-service receives the update request
- then the request is rejected before persistence
- and the response returns an explicit validation failure

### Scenario: duplicate email address

- given a customer submits an email already in use
- when profile-service processes the update
- then the request is rejected
- and the failure reason is explicit and observable

## MODIFIED Profile Event Contract

The system must preserve compatibility with `contracts.customer-profile.v2` while ensuring successful email updates emit the expected fields.

### Scenario: successful update

- given a valid and unique email address
- when the update succeeds
- then the customer profile event is emitted using the existing contract version
- and logs and metrics record the successful path
