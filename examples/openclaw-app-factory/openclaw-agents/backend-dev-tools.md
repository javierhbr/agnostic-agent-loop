# BackendDev — Tools & CLI Reference

## Task & Context Management

- `agentic-agent task list` — view backlog and in-progress tasks
- `agentic-agent task claim <ID>` — claim task (records branch + timestamp)
- `agentic-agent context build --task <ID>` — load context bundle (includes openspec + tech-stack)
- `agentic-agent validate` — run all quality validators before completing
- `agentic-agent task complete <ID>` — mark done (auto-captures commits)
- `agentic-agent openspec check <spec-id>` — verify spec readiness before claiming

## Testing & Coverage

- Language-specific test runners: `go test -cover ./...`, `npm run test`, `python -m pytest`, `cargo test`
- Coverage reports: `go tool cover`, `nyc`, `coverage.py`
- Integration testing: `docker-compose up && go test ./integration`, API client tests

## Database Tools

- **PostgreSQL:** `psql -h localhost -U user -d dbname -c "SELECT..."`
- **Redis:** `redis-cli`
- **MongoDB:** `mongosh`
- **DynamoDB:** `aws dynamodb scan --table-name <table>`
- **Migrations:** `migrate -path ./migrations -database "postgres://..." up`

## API Testing & Documentation

- `curl -X GET http://localhost:8080/api/v1/health`
- `httpie` — user-friendly HTTP client: `http GET localhost:8080/api/v1/users`
- API contract reference: read `.agentic/contracts/<spec-id>.yaml` (OpenAPI spec)

## Key Paths

- `.agentic/tasks/` — task manifests (backlog.yaml, in-progress.yaml, done.yaml)
- `.agentic/spec/` — all openspec proposals + acceptance criteria
- `.agentic/context/` — global-context.md, tech-stack.md, decisions.md
- `.agentic/coordination/` — announcements.yaml, kill-signals.yaml, reservations.yaml
- `.agentic/contracts/` — API contracts (endpoints, schemas, auth) stored by TechLead after gate-check

## Announcements Format

When announcing task completion to TechLead:

```yaml
- from_agent: backend-dev
  to_agent: tech-lead
  project_id: proj-001
  status: complete
  summary: "Auth API layer implemented with 3 endpoints"
  data:
    task_id: TASK-042
    branch: feature/auth-api-v1
    files_changed:
      - internal/api/auth/handler.go
      - internal/auth/service.go
      - internal/db/migrations/002_auth_tables.sql
    test_results:
      unit_tests: 24 passed
      integration_tests: 8 passed
      coverage: 87%
    ac_coverage:
      - "POST /auth/register validates email and password" ✅
      - "POST /auth/login returns JWT token" ✅
      - "POST /auth/refresh rotates token pair" ✅
    commits: [abc123, def456, ghi789]
```

## Git Workflow

- Create feature branch: `git checkout -b feature/task-id-short-name` (TechLead handles main)
- Commit messages: `feat: <AC description>` or `fix: <bug description>`
- Push to feature branch: `git push origin feature/<branch>`
- Never force-push or merge directly to main
