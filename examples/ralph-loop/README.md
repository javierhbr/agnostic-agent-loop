# Ralph Loop Skill - Tool-Agnostic Iterative Task Completion

## Overview

The `/ralph-loop` skill implements Ralph Wiggum methodology: **iterate until the task converges to completion**. It works across all AI chat tools (Claude Code, Cursor, Copilot, Windsurf, etc.).

## Quick Start

### 1. Claim a Task

```bash
# List available tasks
agentic-agent task list

# Claim a task
agentic-agent task claim TASK-001
```

### 2. Start Ralph Loop in AI Chat

```
You: /ralph-loop

AI: I'm using the ralph-loop skill to iteratively complete this task.

[AI reads specs, implements code, verifies criteria, iterates until done]

<promise>TASK COMPLETE</promise>
```

## How It Works

```
┌─────────────────────────────────────────────────────────────┐
│                    RALPH LOOP                                │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Iteration 1:                                                │
│    ├─ Read specs                                             │
│    ├─ Write code                                             │
│    ├─ Run tests         → 2/5 pass                           │
│    └─ Status: Not done yet                                   │
│                                                               │
│  Iteration 2:                                                │
│    ├─ Fix failing tests                                      │
│    ├─ Run tests         → 5/5 pass ✓                         │
│    ├─ Check linter      → 3 errors                           │
│    └─ Status: Not done yet                                   │
│                                                               │
│  Iteration 3:                                                │
│    ├─ Fix linter errors                                      │
│    ├─ Run tests         → 5/5 pass ✓                         │
│    ├─ Check linter      → 0 errors ✓                         │
│    ├─ Verify API        → 200 OK ✓                           │
│    └─ Status: ALL CRITERIA MET!                              │
│                                                               │
│  → Complete task and exit                                    │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

## Example Workflow

### Create a Task

```bash
cat > .agentic/tasks/backlog.yaml <<EOF
tasks:
  - id: TASK-001
    title: "Add health check endpoint"
    description: "Create /health endpoint that returns server status"
    acceptance:
      - "Tests pass"
      - "GET /health returns 200"
      - "Response includes version number"
      - "No linter errors"
    spec_refs:
      - health-check-spec
    scope:
      - internal/api
EOF
```

### Create Spec

```bash
mkdir -p .agentic/spec
cat > .agentic/spec/health-check-spec.md <<EOF
# Health Check Endpoint Specification

## Requirements

- Endpoint: GET /health
- Response: JSON
- Fields:
  - status: "ok" | "error"
  - version: string (from build info)
  - uptime: duration in seconds

## Example Response

\`\`\`json
{
  "status": "ok",
  "version": "1.0.0",
  "uptime": 3600
}
\`\`\`

## Acceptance Criteria

1. Returns 200 OK when server is healthy
2. Includes all required fields
3. Version matches build version
4. Unit tests achieve 100% coverage
EOF
```

### Claim and Run Ralph Loop

```bash
# Claim the task
agentic-agent task claim TASK-001
```

Then in your AI chat tool (Claude Code, Cursor, etc.):

```
You: /ralph-loop
```

The AI will:
1. Load task TASK-001
2. Read the health-check-spec.md
3. Start iterating:
   - Iteration 1: Create basic endpoint
   - Iteration 2: Add tests
   - Iteration 3: Fix test failures
   - Iteration 4: Add version info
   - Iteration 5: Verify all criteria pass
4. Complete the task automatically

## Acceptance Criteria

Tasks must define **clear, verifiable acceptance criteria**:

### ✅ Good Criteria (Specific, Testable)

```yaml
acceptance:
  - "go test ./... passes with 0 failures"
  - "GET /api/users returns 200"
  - "golangci-lint run shows 0 errors"
  - "Response time < 100ms for /api/health"
```

### ❌ Bad Criteria (Vague, Unverifiable)

```yaml
acceptance:
  - "Code works well"           # How do you verify?
  - "Looks good"                # Subjective
  - "No major issues"           # What counts as major?
  - "Ready for production"      # Undefined
```

## Verification Commands

The AI runs actual commands to verify criteria:

```bash
# Tests pass
go test ./... -v

# API returns 200
curl http://localhost:8080/health

# No linter errors
golangci-lint run

# Build succeeds
go build ./cmd/...

# Specific test coverage
go test -cover ./internal/api/...
```

## Convergence Detection

Ralph loop stops when:

### 1. Success (All Criteria Met)

```
─── Iteration 4 ───

Verification:
✓ go test ./... → PASS (10/10 tests)
✓ GET /health → 200 OK
✓ Response includes version
✓ golangci-lint → 0 errors

ALL 4/4 CRITERIA MET! 🎉

$ agentic-agent task complete TASK-001

<promise>TASK COMPLETE</promise>
```

### 2. Max Iterations (10) Reached

```
─── Iteration 10 ───

Status: 3/4 criteria met
⚠ Max iterations reached

Remaining issues:
✗ Response time exceeds 100ms (avg: 150ms)

Recommend:
- Investigate performance bottleneck
- Consider caching or optimization
- Break into separate performance task
```

### 3. Stuck (No Progress)

```
─── Iteration 7 ───
Same error as iterations 5, 6:
  ✗ Tests fail: connection refused

⚠ No progress in last 3 iterations

Possible causes:
- Service not running
- Wrong port configuration
- Test environment issue

Recommend: Manual intervention needed
```

## Tool Support

### Claude Code ✅

```
# Has everything built-in
- ✅ File read/write
- ✅ Bash execution
- ✅ agentic-agent CLI access

Usage: Just type /ralph-loop
```

### Cursor ✅

```
# Full support
- ✅ File operations
- ✅ Terminal access
- ✅ CLI execution

Usage: Type /ralph-loop or @ralph-loop
```

### GitHub Copilot ✅

```
# Works via terminal
- ✅ Can read files
- ✅ Can suggest commands
- ⚠ May need manual command execution

Usage: Paste skill instructions or type /ralph-loop
```

### Windsurf ✅

```
# Full support
- ✅ File access
- ✅ Terminal execution
- ✅ CLI integration

Usage: Type /ralph-loop
```

## Advanced Usage

### Custom Max Iterations

```yaml
# In task definition
metadata:
  max_iterations: 20  # For complex tasks
```

### Multiple Tasks in Sequence

```bash
# Claim first task
agentic-agent task claim TASK-001

# In chat: /ralph-loop
# [Task completes]

# Claim next task
agentic-agent task claim TASK-002

# In chat: /ralph-loop again
# [Task completes]
```

### Resuming After Pause

```bash
# If you paused mid-loop
agentic-agent task continue

# In chat: /ralph-loop
# [Resumes from where it left off]
```

## Best Practices

### 1. Start With Clear Specs

```markdown
# BAD: Vague spec
Add a login feature.

# GOOD: Clear spec
## Login Endpoint Spec

- POST /api/login
- Accept: username (string), password (string)
- Return: JWT token (string)
- Validate credentials against database
- Use bcrypt for password hashing
```

### 2. Define Measurable Criteria

```yaml
# BAD
acceptance:
  - "Works correctly"

# GOOD
acceptance:
  - "POST /api/login returns 200 for valid credentials"
  - "POST /api/login returns 401 for invalid credentials"
  - "JWT token validates with public key"
  - "All tests pass: go test ./internal/auth/..."
```

### 3. Keep Scope Focused

```yaml
# BAD: Too broad
scope:
  - internal/
  - cmd/
  - pkg/

# GOOD: Specific
scope:
  - internal/auth
```

### 4. Provide Examples in Specs

````markdown
## Example Request

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'
```

## Expected Response

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expiresAt": "2024-01-15T10:30:00Z"
}
```
````

## Troubleshooting

### Problem: AI skips verification

**Symptom:** AI says "tests pass" without running them

**Solution:**
- Be explicit: "Run go test and show output"
- Include verification commands in acceptance criteria
- Ask AI to paste actual command output

### Problem: Infinite loop

**Symptom:** Same iteration repeating

**Solution:**
- Check if criteria are achievable
- Review task scope (might be too large)
- Break into smaller tasks

### Problem: Criteria not met after max iterations

**Symptom:** Stuck at iteration 10

**Solution:**
```bash
# Review task
agentic-agent task show TASK-001

# Check what's blocking
# Update criteria if unrealistic
# Or break into sub-tasks
```

### Problem: AI completes without meeting criteria

**Symptom:** Claims done but tests fail

**Solution:**
- Report the issue
- Ask AI to re-verify each criterion
- Be specific: "Show me go test output"

## Integration with Autopilot

Ralph loop works with autopilot:

```bash
# Terminal autopilot (no AI)
agentic-agent autopilot start

# vs.

# Manual ralph loop (AI-driven)
# 1. Claim task manually
# 2. Chat: /ralph-loop
# 3. AI iterates until done
```

**Differences:**

| Feature | Autopilot | Ralph Loop |
|---------|-----------|------------|
| Execution | CLI automated | AI chat iterative |
| Agent | Optional (--execute-agent) | Always (AI in chat) |
| Verification | Readiness checks | Full criteria verification |
| Iteration | One-shot per task | Loop until complete |
| Progress | Batch processing | Real-time visible |

## Examples

See [`examples/ralph-loop/`](./examples/) for:
- Complete task examples
- Spec templates
- Verification scripts
- Multi-iteration scenarios

## Related

- [Autopilot Documentation](../../docs/multi-agent-execution.md)
- [Task Management Guide](../../README.md#task-management)
- [Spec Writing Guide](../../docs/spec-writing.md)
