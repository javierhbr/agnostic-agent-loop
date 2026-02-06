# Authentication Requirements

## Overview
JWT-based authentication for the REST API.

## Requirements
1. Users authenticate with email and password
2. Access tokens use RS256 signing, expire in 15 minutes
3. Refresh tokens expire in 7 days
4. Protected endpoints return 401 for missing or expired tokens

## Endpoints
- `POST /auth/register` — Create account
- `POST /auth/login` — Get access + refresh tokens
- `POST /auth/refresh` — Exchange refresh token for new access token
- `DELETE /auth/logout` — Invalidate refresh token

## Security
- Passwords hashed with bcrypt (cost factor 12)
- Tokens stored server-side for revocation support
- Rate limiting on login endpoint (5 attempts per minute)
