# Global Context

## Project Overview
A Go REST API demonstrating agent-aware skills. Each AI agent tool receives tailored rules and skill instructions based on its capabilities.

## Architecture
- Go with net/http standard library
- JWT-based authentication
- PostgreSQL for user storage
- Middleware chain: logging, auth, validation

## Conventions
- Package per domain (auth, users, middleware)
- Errors wrapped with `fmt.Errorf("...: %w", err)`
- Tests colocated with source files
