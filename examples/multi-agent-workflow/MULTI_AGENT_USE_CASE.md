# Use Case: Switching Between Agent Tools

This walkthrough demonstrates the core value of the agentic-agent CLI: **any AI agent tool can pick up exactly where another left off**. All task state, context, and learnings are persisted in the `.agentic/` directory as plain YAML and Markdown files. No proprietary format, no vendor lock-in.

The scenario below follows a developer building a user authentication feature while bouncing between four different AI agent interfaces — Claude Code CLI, Claude Code VSCode extension, GitHub Copilot, and Antigravity IDE with Gemini. The developer switches tools freely, even going back to tools already used, and never loses context.

## The Scenario

You are building JWT-based user authentication for a web application. The feature decomposes into four subtasks:

1. **Create JWT token service** — core token generation and validation logic
2. **Implement auth middleware** — HTTP middleware that protects routes
3. **Add user registration endpoint** — registration handler with password hashing
4. **Write integration tests** — end-to-end tests for the auth flow

The work happens across six phases using four different tools, with deliberate back-and-forth switching:

| Phase | Tool | Work |
|-------|------|------|
| 1 | Claude Code CLI (terminal) | Setup, task creation, JWT service |
| 2 | Claude Code VSCode Extension | Auth middleware |
| 3 | GitHub Copilot (VSCode) | Registration endpoint |
| 4 | Antigravity IDE + Gemini | Integration tests — discovers a bug |
| 5 | Claude Code VSCode Extension | **Returns** to fix the bug Gemini found |
| 6 | Claude Code CLI (terminal) | **Returns** for final review and completion |

## Prerequisites

- `agentic-agent` CLI built and available (see [Installation](../../README.md#installation))
- VSCode with the [Claude Code extension](https://marketplace.visualstudio.com/items?itemName=Anthropic.claude-code) installed
- GitHub Copilot with agent mode enabled in VSCode
- Antigravity IDE with Gemini integration
- A Git repository for your project

## Setup

Initialize a new project:

```bash
mkdir my-auth-project && cd my-auth-project
git init

# Initialize the agentic framework
agentic-agent init --name "My Auth Project"
```

This creates:

```
my-auth-project/
├── agnostic-agent.yaml          # Project config (read from current directory)
└── .agentic/
    ├── tasks/
    │   ├── backlog.yaml
    │   ├── in-progress.yaml
    │   └── done.yaml
    ├── spec/                    # Specification files referenced by tasks
    ├── context/
    │   ├── global-context.md
    │   ├── rolling-summary.md
    │   ├── decisions.md
    │   └── assumptions.md
    └── agent-rules/
        └── base.md
```

---

## Phase 1 — Claude Code CLI (Terminal)

> **Tool**: Claude Code in the terminal
> **Goal**: Create the feature task, decompose it, and implement the JWT token service (Subtask 1)

### 1. Create the parent feature task

```bash
agentic-agent task create \
  --title "Implement user authentication" \
  --description "JWT-based auth with token service, middleware, registration, and tests" \
  --spec-refs ".agentic/spec/04-architecture.md" \
  --outputs "src/auth/jwt.go,src/auth/middleware.go,src/auth/register.go,tests/auth_test.go" \
  --acceptance "JWT tokens generated on login,Token validation rejects expired tokens,Auth middleware blocks unauthenticated requests,User registration with password hashing,Integration tests pass"
```

```
Created task TASK-1738900000: Implement user authentication
```

Verify the spec reference resolves:

```bash
agentic-agent spec resolve .agentic/spec/04-architecture.md
```

### 2. Decompose into subtasks

```bash
agentic-agent task decompose TASK-1738900000 \
  "Create JWT token service" \
  "Implement auth middleware" \
  "Add user registration endpoint" \
  "Write integration tests"
```

```
Decomposed task TASK-1738900000 with 4 subtasks
```

### 3. Claim the task

```bash
agentic-agent task claim TASK-1738900000
```

```
Claimed task TASK-1738900000
```

The task moves from `.agentic/tasks/backlog.yaml` to `.agentic/tasks/in-progress.yaml` with `assigned_to` set to your `$USER`.

### 4. Generate context for the working directory

```bash
mkdir -p src/auth
agentic-agent context generate src/auth
```

This creates `src/auth/context.md`. Claude Code reads this file before touching any code in that directory, following the rule from `.agentic/agent-rules/base.md`:

> Always read `context.md` before starting work in a directory.

### 5. Implement the JWT token service (Subtask 1)

Claude Code CLI helps you write `src/auth/jwt.go` with token generation and validation functions. The implementation follows the spec referenced in the task.

### 6. Record learnings

```bash
agentic-agent learnings add "Use RS256 for JWT signing in production, HS256 for dev"
agentic-agent learnings add "Token expiry should be 24h for access tokens, 7d for refresh tokens"
```

These patterns are written to `.agentic/progress.txt` and persist across all agent sessions.

### 7. Update the rolling summary

Edit `.agentic/context/rolling-summary.md` to reflect what was accomplished:

```markdown
## Current State
- JWT token service implemented in src/auth/jwt.go
- Subtask 1 of 4 complete (Create JWT token service)
- Subtasks remaining: auth middleware, registration endpoint, integration tests
```

### State on disk after Phase 1

| File | Contents |
|------|----------|
| `.agentic/tasks/in-progress.yaml` | Parent task with 4 subtasks, `assigned_to: your-username` |
| `.agentic/context/rolling-summary.md` | Updated with JWT service completion notes |
| `.agentic/progress.txt` | 2 learnings recorded |
| `src/auth/context.md` | Directory context for `src/auth/` |
| `src/auth/jwt.go` | The implemented JWT service |

---

## Phase 2 — Claude Code VSCode Extension

> **Tool**: Claude Code extension panel in VSCode
> **Goal**: Catch up on Phase 1, then implement auth middleware (Subtask 2)

### 1. Open VSCode and the Claude Code panel

Open your project in VSCode. The Claude Code extension has terminal access, so it can run all the same CLI commands.

### 2. Build a context bundle to catch up

```bash
agentic-agent context build --task TASK-1738900000
```

This bundles everything the agent needs into a single output:

- The task definition (title, subtasks, acceptance criteria)
- **Resolved specs** from the task's `spec_refs` (e.g., `.agentic/spec/04-architecture.md` content)
- Global context from `.agentic/context/global-context.md`
- Rolling summary from `.agentic/context/rolling-summary.md`
- All directory context files (including `src/auth/context.md`)

The agent now has full context about what was done, what remains, and the project constraints — without needing you to repeat anything.

### 3. Review existing learnings

```bash
agentic-agent learnings list
```

```
Codebase Patterns (2):
1. Use RS256 for JWT signing in production, HS256 for dev
2. Token expiry should be 24h for access tokens, 7d for refresh tokens
```

The agent sees the signing algorithm and expiry decisions from Phase 1 and applies them consistently.

### 4. Implement the auth middleware (Subtask 2)

With the context bundle loaded, the Claude Code extension helps you write `src/auth/middleware.go`. It uses the JWT service from Subtask 1 and follows the patterns recorded in learnings.

### 5. Update context after implementation

```bash
agentic-agent context generate src/auth
```

This regenerates `src/auth/context.md` to reflect the new middleware file.

### 6. Record additional learnings

```bash
agentic-agent learnings add "Middleware should return 401 with WWW-Authenticate header for missing tokens"
agentic-agent learnings add "Use request context to pass authenticated user to handlers"
```

### 7. Update the rolling summary

```markdown
## Current State
- JWT token service implemented in src/auth/jwt.go
- Auth middleware implemented in src/auth/middleware.go
- Subtasks 1 and 2 of 4 complete
- Subtasks remaining: registration endpoint, integration tests
```

### State on disk after Phase 2

| File | What changed |
|------|-------------|
| `.agentic/context/rolling-summary.md` | Updated with middleware completion |
| `.agentic/progress.txt` | Now has 4 learnings (2 new from this phase) |
| `src/auth/context.md` | Regenerated to include middleware |
| `src/auth/middleware.go` | New file — the auth middleware |

---

## Phase 3 — GitHub Copilot Agent Mode (VSCode)

> **Tool**: GitHub Copilot in agent mode, in the same VSCode
> **Goal**: Implement user registration endpoint (Subtask 3)

### 1. Switch to Copilot agent mode

Still in VSCode, switch from Claude Code to GitHub Copilot in agent mode. Copilot does not have native CLI integration, so you run CLI commands in the VSCode terminal yourself while Copilot assists with code.

### 2. Point Copilot to the base rules

Tell Copilot to read `.agentic/agent-rules/base.md`:

```
- Always read context.md before starting work in a directory.
- Update context.md if you change the logic/architecture.
- Keep tasks small.
```

These rules apply identically to every agent tool.

### 3. Check current task state

```bash
agentic-agent task show TASK-1738900000
```

```
ID: TASK-1738900000
Title: Implement user authentication
Status: in-progress
Assigned To: your-username
Subtasks:
  [TASK-1738900000.1] Create JWT token service
  [TASK-1738900000.2] Implement auth middleware
  [TASK-1738900000.3] Add user registration endpoint
  [TASK-1738900000.4] Write integration tests
```

### 4. Review all learnings from prior phases

```bash
agentic-agent learnings list
```

```
Codebase Patterns (4):
1. Use RS256 for JWT signing in production, HS256 for dev
2. Token expiry should be 24h for access tokens, 7d for refresh tokens
3. Middleware should return 401 with WWW-Authenticate header for missing tokens
4. Use request context to pass authenticated user to handlers
```

Copilot sees all four learnings from Phases 1 and 2. These inform the registration endpoint — it needs to generate tokens using RS256 and set the correct expiry.

### 5. Read directory context

```bash
cat src/auth/context.md
```

Copilot reads the context file to understand the existing module structure, the JWT service API, and the middleware behavior.

### 6. Implement user registration (Subtask 3)

Copilot helps write `src/auth/register.go` with:

- Input validation for email and password
- Password hashing with bcrypt
- User storage
- JWT token generation on successful registration (using the existing `jwt.go` service)

### 7. Update context and record learnings

```bash
agentic-agent context generate src/auth
agentic-agent learnings add "Use bcrypt cost factor 12 for password hashing"
```

### 8. Update the rolling summary

```markdown
## Current State
- JWT token service implemented in src/auth/jwt.go
- Auth middleware implemented in src/auth/middleware.go
- Registration endpoint implemented in src/auth/register.go
- Subtasks 1, 2, and 3 of 4 complete
- Subtask remaining: integration tests
```

### State on disk after Phase 3

| File | What changed |
|------|-------------|
| `.agentic/context/rolling-summary.md` | Updated with registration completion |
| `.agentic/progress.txt` | Now has 5 learnings |
| `src/auth/context.md` | Regenerated to include registration endpoint |
| `src/auth/register.go` | New file — user registration handler |

---

## Phase 4 — Antigravity IDE with Gemini

> **Tool**: Antigravity IDE with Gemini
> **Goal**: Write integration tests (Subtask 4) — and discover a bug along the way

### 1. Open the project in Antigravity IDE

Close VSCode (or leave it open — it doesn't matter). Open the project in Antigravity IDE. Gemini is available as the AI assistant, and the IDE has a built-in terminal.

### 2. Catch up using the context bundle

In the Antigravity terminal:

```bash
agentic-agent context build --task TASK-1738900000
```

Gemini reads the bundled output. It now knows:

- Three subtasks are complete (JWT service, middleware, registration)
- One subtask remains (integration tests)
- Five learnings have been recorded by three different agent tools
- The `src/auth/` directory has three Go files

### 3. Review all learnings

```bash
agentic-agent learnings list
```

```
Codebase Patterns (5):
1. Use RS256 for JWT signing in production, HS256 for dev
2. Token expiry should be 24h for access tokens, 7d for refresh tokens
3. Middleware should return 401 with WWW-Authenticate header for missing tokens
4. Use request context to pass authenticated user to handlers
5. Use bcrypt cost factor 12 for password hashing
```

Gemini sees every pattern accumulated across three prior tools. Each learning becomes a test case.

### 4. Read the auth module context

```bash
cat src/auth/context.md
```

### 5. Write integration tests (Subtask 4)

Gemini helps write `tests/auth_test.go` covering:

- JWT token generation and validation (learning 1, 2)
- Expired token rejection (learning 2)
- Middleware blocking unauthenticated requests (learning 3)
- Authenticated user propagation via request context (learning 4)
- Registration with password hashing (learning 5)
- Registration with invalid input

### 6. Discover a bug in middleware

While writing the test for expired tokens, Gemini notices that `middleware.go` does not check the `exp` claim properly — it parses the token but does not reject expired ones. The test fails:

```
--- FAIL: TestMiddleware_ExpiredToken
    Expected: 401 Unauthorized
    Got: 200 OK
```

### 7. Record the bug as a learning

```bash
agentic-agent learnings add "BUG: middleware.go does not validate token expiry - needs exp claim check"
```

### 8. Generate context for the tests directory

```bash
mkdir -p tests
agentic-agent context generate tests
```

### 9. Update the rolling summary

```markdown
## Current State
- JWT token service implemented in src/auth/jwt.go
- Auth middleware implemented in src/auth/middleware.go (BUG: missing exp validation)
- Registration endpoint implemented in src/auth/register.go
- Integration tests written in tests/auth_test.go (1 failing: expired token test)
- Subtasks 1-4 attempted, but middleware bug blocks completion
- BUG FOUND: middleware does not check token expiry claim
```

### State on disk after Phase 4

| File | What changed |
|------|-------------|
| `.agentic/context/rolling-summary.md` | Updated with bug report |
| `.agentic/progress.txt` | Now has 6 learnings (including the bug) |
| `tests/context.md` | New — context for the tests directory |
| `tests/auth_test.go` | New — integration tests (1 failing) |

---

## Phase 5 — Back to Claude Code VSCode Extension

> **Tool**: Claude Code extension panel in VSCode (second visit)
> **Goal**: Fix the middleware bug that Gemini discovered in Phase 4

This is where the back-and-forth pays off. A bug was found in Antigravity IDE, and now we return to the Claude Code VSCode extension to fix it.

### 1. Reopen VSCode and the Claude Code panel

Open the project in VSCode again. Claude Code extension has no memory of the Phase 2 session — but it doesn't need one. Everything is on disk.

### 2. Read the rolling summary

```bash
cat .agentic/context/rolling-summary.md
```

The Claude Code extension immediately sees:

```
BUG FOUND: middleware does not check token expiry claim
```

### 3. Review the bug learning

```bash
agentic-agent learnings list
```

```
Codebase Patterns (6):
1. Use RS256 for JWT signing in production, HS256 for dev
2. Token expiry should be 24h for access tokens, 7d for refresh tokens
3. Middleware should return 401 with WWW-Authenticate header for missing tokens
4. Use request context to pass authenticated user to handlers
5. Use bcrypt cost factor 12 for password hashing
6. BUG: middleware.go does not validate token expiry - needs exp claim check
```

Learning #6 was recorded by Gemini in Antigravity IDE. Claude Code in VSCode reads it and knows exactly what to fix.

### 4. Fix the middleware bug

Claude Code helps you update `src/auth/middleware.go` to add expiry validation:

```go
// Check token expiry
if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
    w.Header().Set("WWW-Authenticate", "Bearer error=\"invalid_token\"")
    http.Error(w, "Token expired", http.StatusUnauthorized)
    return
}
```

### 5. Run the failing test

```bash
go test ./tests/ -run TestMiddleware_ExpiredToken -v
```

```
--- PASS: TestMiddleware_ExpiredToken (0.01s)
PASS
```

The test that Gemini wrote in Antigravity now passes in VSCode.

### 6. Update context and record the fix

```bash
agentic-agent context generate src/auth
agentic-agent learnings add "Fixed: middleware now validates exp claim and returns 401 for expired tokens"
```

### 7. Validate

```bash
agentic-agent validate
```

```
PASS  directory-context
PASS  context-update
PASS  task-scope
PASS  task-size
PASS  browser-verification

Summary: 5 passed | 0 warnings | 0 failed
```

### 8. Update the rolling summary

```markdown
## Current State
- JWT token service implemented in src/auth/jwt.go
- Auth middleware implemented in src/auth/middleware.go (bug fixed: exp validation added)
- Registration endpoint implemented in src/auth/register.go
- Integration tests passing in tests/auth_test.go
- All 4 subtasks complete, all tests passing
- Ready for final review and task completion
```

### State on disk after Phase 5

| File | What changed |
|------|-------------|
| `.agentic/context/rolling-summary.md` | Updated — bug marked as fixed |
| `.agentic/progress.txt` | Now has 7 learnings |
| `src/auth/middleware.go` | Updated — expiry validation added |
| `src/auth/context.md` | Regenerated to reflect the fix |

---

## Phase 6 — Back to Claude Code CLI (Terminal)

> **Tool**: Claude Code in the terminal (second visit)
> **Goal**: Final review, complete the task, validate everything

### 1. Review the task state

```bash
agentic-agent task show TASK-1738900000
```

```
ID: TASK-1738900000
Title: Implement user authentication
Status: in-progress
Assigned To: your-username
Subtasks:
  [TASK-1738900000.1] Create JWT token service
  [TASK-1738900000.2] Implement auth middleware
  [TASK-1738900000.3] Add user registration endpoint
  [TASK-1738900000.4] Write integration tests
```

### 2. Run all tests

```bash
go test ./... -v
```

```
--- PASS: TestJWT_GenerateToken (0.01s)
--- PASS: TestJWT_ValidateToken (0.01s)
--- PASS: TestJWT_ExpiredToken (0.01s)
--- PASS: TestMiddleware_ValidToken (0.02s)
--- PASS: TestMiddleware_ExpiredToken (0.01s)
--- PASS: TestMiddleware_MissingToken (0.00s)
--- PASS: TestRegister_ValidUser (0.15s)
--- PASS: TestRegister_InvalidEmail (0.00s)
--- PASS: TestRegister_WeakPassword (0.00s)
PASS
```

All tests pass. Code written across four different tools works together.

### 3. Review all learnings accumulated across four tools

```bash
agentic-agent learnings list
```

```
Codebase Patterns (7):
1. Use RS256 for JWT signing in production, HS256 for dev
2. Token expiry should be 24h for access tokens, 7d for refresh tokens
3. Middleware should return 401 with WWW-Authenticate header for missing tokens
4. Use request context to pass authenticated user to handlers
5. Use bcrypt cost factor 12 for password hashing
6. BUG: middleware.go does not validate token expiry - needs exp claim check
7. Fixed: middleware now validates exp claim and returns 401 for expired tokens
```

Seven learnings from four tools. The bug lifecycle is visible: learning #6 (found by Gemini in Antigravity) and learning #7 (fixed by Claude Code in VSCode).

### 4. Complete the task

```bash
agentic-agent task complete TASK-1738900000
```

```
Completed task TASK-1738900000
```

The task moves from `.agentic/tasks/in-progress.yaml` to `.agentic/tasks/done.yaml`.

### 5. Final validation

```bash
agentic-agent validate
```

```
PASS  directory-context
PASS  context-update
PASS  task-scope
PASS  task-size
PASS  browser-verification

Summary: 5 passed | 0 warnings | 0 failed
```

The feature is complete. Four tools contributed. Zero context was lost.

---

## What Made This Work

| Principle | How It Applies |
|-----------|---------------|
| **Shared state directory** | `.agentic/` is the single source of truth. Every agent reads and writes the same YAML and Markdown files — Claude Code, Copilot, and Gemini all see the same state. |
| **CLI as the bridge** | The `agentic-agent` CLI is the only interface to task state. It works identically in every terminal — standalone, VSCode integrated, or Antigravity IDE. |
| **Context bundling** | `agentic-agent context build --task TASK-ID` packages everything a new agent session needs to continue work without re-explanation. |
| **Persistent learnings** | `agentic-agent learnings add` and `learnings list` carry patterns, decisions, and even bug reports across agent boundaries. |
| **Agent-agnostic rules** | `.agentic/agent-rules/base.md` defines behavior rules that apply to every agent equally. |
| **Spec resolution** | `spec resolve` and `context build` resolve spec references from configured directories. Specs travel with the project — any agent reads the same requirements. |
| **Validation independence** | `agentic-agent validate` checks the file system, not the agent. Work done by any tool is validated the same way. |
| **Back-and-forth is free** | Returning to a previously used tool costs nothing — the agent reads the current state from disk and continues. No "session" to restore. |

## State Flow Across All Phases

```
Phase 1               Phase 2               Phase 3               Phase 4               Phase 5               Phase 6
Claude CLI             Claude VSCode         Copilot (VSCode)      Antigravity+Gemini    Claude VSCode         Claude CLI
──────────             ─────────────         ────────────────      ──────────────────    ─────────────         ──────────

task create            context build         task show             context build         cat rolling-summary   task show
task decompose         learnings list        learnings list        learnings list        learnings list        go test ./...
task claim               -> 2 learnings       -> 4 learnings       -> 5 learnings         -> 6 learnings      learnings list
context generate                                                                                                -> 7 learnings

[jwt.go]               [middleware.go]       [register.go]         [auth_test.go]        [fix middleware.go]   task complete
                                                                    BUG FOUND!                                  validate
learnings add (x2)     context generate      context generate      context generate      context generate
update summary         learnings add (x2)    learnings add (x1)    learnings add (x1)    learnings add (x1)
                       update summary        update summary        update summary        validate
                                                                                         update summary

Learnings: 2           Learnings: 4          Learnings: 5          Learnings: 6          Learnings: 7          DONE
```

## Key Takeaway

The developer never had to explain the project state to a new agent. Each agent — regardless of which tool it ran in — read the same `.agentic/` files and continued seamlessly. A bug found by Gemini in Antigravity IDE was fixed by Claude Code in VSCode using a learning recorded in the shared progress file. The `agentic-agent` CLI was the constant across all four tools — the bridge that made agent switching a non-event.
