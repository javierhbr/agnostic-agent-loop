# API Design Specification

## Architecture
- Go HTTP server using standard library
- SQLite for persistence
- Repository pattern for data access
- Middleware chain: logging -> auth -> handler

## Directory Structure
```
internal/
  auth/       — JWT service, middleware, handlers
  models/     — Domain types (User, Token)
  repository/ — Data access interfaces and implementations
  server/     — HTTP router and middleware
cmd/
  api/        — Main entry point
```

## Error Format
All errors return JSON:
```json
{
  "error": "human-readable message",
  "code": "MACHINE_CODE"
}
```

## Response Format
Success responses wrap data:
```json
{
  "data": { ... }
}
```
