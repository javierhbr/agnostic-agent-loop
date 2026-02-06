# Spec Kit: Authentication Feature

## Constitution Reference
This spec follows the project constitution defined in `.specify/memory/constitution.md`.

## Feature: User Authentication
Add JWT-based authentication to the REST API.

## Requirements
1. **Registration**: Users register with email + password. Passwords hashed with bcrypt.
2. **Login**: Returns access token (15min) and refresh token (7d).
3. **Token Refresh**: Exchange valid refresh token for new access token.
4. **Logout**: Invalidate refresh token server-side.

## Scenarios

### Happy Path
- User registers -> receives confirmation
- User logs in -> receives access + refresh tokens
- User accesses protected endpoint -> succeeds with valid token
- User refreshes token -> receives new access token

### Error Cases
- Register with existing email -> 409 Conflict
- Login with wrong password -> 401 Unauthorized
- Access with expired token -> 401 Unauthorized
- Refresh with revoked token -> 401 Unauthorized

## Implementation Notes
- Use RS256 for JWT signing
- Store refresh tokens in database for revocation
- Rate limit login to 5 attempts/minute per IP
