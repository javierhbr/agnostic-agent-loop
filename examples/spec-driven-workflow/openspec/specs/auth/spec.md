# OpenSpec: Authentication

## Change: add-user-authentication

### Proposal
Add JWT-based user authentication to support protected API endpoints.

### Requirements
- Email/password registration with bcrypt hashing
- JWT access tokens (RS256, 15 min TTL)
- Refresh tokens (7 day TTL, server-side storage)
- Auth middleware for protected routes
- Rate limiting on login endpoint

### Design
- Repository pattern for user/token storage
- Middleware chain: logging -> rate-limit -> auth -> handler
- Tokens stored in SQLite for revocation support

### Tasks
- [ ] Create User model and repository
- [ ] Implement JWT token service
- [ ] Add auth middleware
- [ ] Create registration endpoint
- [ ] Create login endpoint
- [ ] Create token refresh endpoint
- [ ] Write integration tests

### Verification Criteria
- All endpoints return correct status codes
- Expired tokens are rejected
- Revoked refresh tokens cannot be reused
- Rate limiting blocks excessive login attempts
