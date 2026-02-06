# Use Case: Switching Between Agent Tools

This walkthrough demonstrates the core value of the agentic-agent CLI: **any AI agent tool can pick up exactly where another left off**. All task state, context, and learnings are persisted in the `.agentic/` directory as plain YAML and Markdown files. No proprietary format, no vendor lock-in.

The scenario below follows a developer building a user authentication feature across three different AI agent interfaces — Claude Code CLI in the terminal, the Claude Code VSCode extension, and GitHub Copilot agent mode in VSCode — without losing a single piece of context.

## The Scenario

You are building JWT-based user authentication for a web application. The feature decomposes into three subtasks:

1. **Create JWT token service** — core token generation and validation logic
2. **Implement auth middleware** — HTTP middleware that protects routes
3. **Write integration tests** — end-to-end tests for the auth flow

Each subtask will be completed with a different AI agent tool.

## Prerequisites

- `agentic-agent` CLI built and available (see [Installation](../../README.md#installation))
- VSCode with the [Claude Code extension](https://marketplace.visualstudio.com/items?itemName=Anthropic.claude-code) installed
- GitHub Copilot with agent mode enabled in VSCode
- A Git repository for your project

## Setup

Initialize a new project and create the feature task:

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
  --description "JWT-based auth with token service, middleware, and integration tests" \
  --spec-refs ".agentic/spec/04-architecture.md" \
  --outputs "src/auth/jwt.go,src/auth/middleware.go,tests/auth_test.go" \
  --acceptance "JWT tokens generated on login,Token validation rejects expired tokens,Auth middleware blocks unauthenticated requests,Integration tests pass"
```

```
Created task TASK-1738900000: Implement user authentication
```

### 2. Decompose into subtasks

```bash
agentic-agent task decompose TASK-1738900000 \
  "Create JWT token service" \
  "Implement auth middleware" \
  "Write integration tests"
```

```
Decomposed task TASK-1738900000 with 3 subtasks
```

### 3. Verify the decomposition

```bash
agentic-agent task show TASK-1738900000
```

```
ID: TASK-1738900000
Title: Implement user authentication
Status: pending (backlog)
Subtasks:
  [TASK-1738900000.1] Create JWT token service
  [TASK-1738900000.2] Implement auth middleware
  [TASK-1738900000.3] Write integration tests
```

### 4. Claim the task

```bash
agentic-agent task claim TASK-1738900000
```

```
Claimed task TASK-1738900000
```

The task moves from `.agentic/tasks/backlog.yaml` to `.agentic/tasks/in-progress.yaml` with `assigned_to` set to your `$USER`.

### 5. Generate context for the working directory

```bash
mkdir -p src/auth
agentic-agent context generate src/auth
```

This creates `src/auth/context.md`. Claude Code reads this file before touching any code in that directory, following the rule from `.agentic/agent-rules/base.md`:

> Always read `context.md` before starting work in a directory.

### 6. Implement the JWT token service

Claude Code CLI helps you write `src/auth/jwt.go` with token generation and validation functions. The implementation follows the spec referenced in the task.

### 7. Record learnings

```bash
agentic-agent learnings add "Use RS256 for JWT signing in production, HS256 for dev"
agentic-agent learnings add "Token expiry should be 24h for access tokens, 7d for refresh tokens"
```

These patterns are written to `.agentic/progress.txt` and persist across all agent sessions.

### 8. Update the rolling summary

Edit `.agentic/context/rolling-summary.md` to reflect what was accomplished:

```markdown
## Current State
- JWT token service implemented in src/auth/jwt.go
- Subtask 1 of 3 complete (Create JWT token service)
- Subtasks remaining: auth middleware, integration tests
```

### State on disk after Phase 1

| File | Contents |
|------|----------|
| `.agentic/tasks/in-progress.yaml` | Parent task with 3 subtasks, `assigned_to: your-username` |
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

### 2. Read the rolling summary to catch up

The Claude Code extension reads the shared state files:

```bash
cat .agentic/context/rolling-summary.md
```

This tells the agent what was done in the previous session — the JWT token service is complete.

### 3. Build a context bundle for the task

```bash
agentic-agent context build --task TASK-1738900000
```

This bundles everything the agent needs into a single output:

- The task definition (title, subtasks, acceptance criteria)
- Global context from `.agentic/context/global-context.md`
- Rolling summary from `.agentic/context/rolling-summary.md`
- All directory context files (including `src/auth/context.md`)

The agent now has full context about what was done, what remains, and the project constraints — without needing you to repeat anything.

### 4. Review existing learnings

```bash
agentic-agent learnings list
```

```
Codebase Patterns (2):
1. Use RS256 for JWT signing in production, HS256 for dev
2. Token expiry should be 24h for access tokens, 7d for refresh tokens
```

The agent sees the signing algorithm and expiry decisions from Phase 1 and applies them consistently.

### 5. Implement the auth middleware (Subtask 2)

With the context bundle loaded, the Claude Code extension helps you write `src/auth/middleware.go`. It uses the JWT service from Subtask 1 and follows the patterns recorded in learnings.

### 6. Update context after implementation

```bash
agentic-agent context generate src/auth
```

This regenerates `src/auth/context.md` to reflect the new middleware file.

### 7. Record additional learnings

```bash
agentic-agent learnings add "Middleware should return 401 with WWW-Authenticate header for missing tokens"
agentic-agent learnings add "Use request context to pass authenticated user to handlers"
```

### 8. Validate

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

All validation rules pass regardless of which agent did the work — the validators check the file system, not the agent.

### 9. Update the rolling summary

```markdown
## Current State
- JWT token service implemented in src/auth/jwt.go
- Auth middleware implemented in src/auth/middleware.go
- Subtasks 1 and 2 of 3 complete
- Subtask remaining: integration tests
```

### State on disk after Phase 2

| File | What changed |
|------|-------------|
| `.agentic/context/rolling-summary.md` | Updated with middleware completion |
| `.agentic/progress.txt` | Now has 4 learnings (2 from Phase 1 + 2 from Phase 2) |
| `src/auth/context.md` | Regenerated to include middleware |
| `src/auth/middleware.go` | New file — the auth middleware |

---

## Phase 3 — GitHub Copilot Agent Mode (VSCode)

> **Tool**: GitHub Copilot in agent mode, in the same VSCode
> **Goal**: Implement integration tests (Subtask 3) and complete the feature

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

### 3. Check current task state from the terminal

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
  [TASK-1738900000.3] Write integration tests
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

Copilot sees all four learnings from Phases 1 and 2. These inform the test implementation — the tests should verify RS256 signing, token expiry behavior, 401 responses, and context propagation.

### 5. Read directory context

```bash
cat src/auth/context.md
```

Copilot reads the context file to understand the module structure, dependencies, and constraints before writing tests.

### 6. Write integration tests (Subtask 3)

Copilot helps write `tests/auth_test.go` covering:

- JWT token generation and validation
- Expired token rejection
- Middleware blocking unauthenticated requests
- Authenticated user propagation via request context

Each learning from Phases 1 and 2 maps to a specific test case.

### 7. Generate context for the tests directory

```bash
mkdir -p tests
agentic-agent context generate tests
```

### 8. Record final learnings

```bash
agentic-agent learnings add "Use httptest.NewServer for integration tests against auth middleware"
```

### 9. Complete the task

```bash
agentic-agent task complete TASK-1738900000
```

```
Completed task TASK-1738900000
```

The task moves from `.agentic/tasks/in-progress.yaml` to `.agentic/tasks/done.yaml`.

### 10. Run final validation

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

Validation passes. It does not matter that three different agent tools contributed to this feature.

### State on disk after Phase 3

| File | What changed |
|------|-------------|
| `.agentic/tasks/done.yaml` | Task moved here from in-progress |
| `.agentic/tasks/in-progress.yaml` | Task removed |
| `.agentic/progress.txt` | Now has 5 learnings |
| `tests/context.md` | New — context for the tests directory |
| `tests/auth_test.go` | New — integration tests |

---

## Phase 4 — Back to Claude Code CLI (Terminal)

> **Tool**: Claude Code in the terminal
> **Goal**: Review the completed feature

### 1. Verify the task is done

```bash
agentic-agent task list
```

```
--- DONE ---
[TASK-1738900000] Implement user authentication
  Subtasks:
    [TASK-1738900000.1] Create JWT token service
    [TASK-1738900000.2] Implement auth middleware
    [TASK-1738900000.3] Write integration tests
```

### 2. Review all learnings accumulated across three tools

```bash
agentic-agent learnings list
```

```
Codebase Patterns (5):
1. Use RS256 for JWT signing in production, HS256 for dev
2. Token expiry should be 24h for access tokens, 7d for refresh tokens
3. Middleware should return 401 with WWW-Authenticate header for missing tokens
4. Use request context to pass authenticated user to handlers
5. Use httptest.NewServer for integration tests against auth middleware
```

Five learnings accumulated across three agent tools, all persisted in a single file.

### 3. Final validation

```bash
agentic-agent validate
```

All rules pass. The feature is complete.

---

## What Made This Work

| Principle | How It Applies |
|-----------|---------------|
| **Shared state directory** | `.agentic/` is the single source of truth. Every agent reads and writes the same YAML and Markdown files. |
| **CLI as the bridge** | The `agentic-agent` CLI is the only interface to task state. It works identically in any terminal — standalone, VSCode integrated, or CI/CD. |
| **Context bundling** | `agentic-agent context build --task TASK-ID` packages everything a new agent session needs to continue work. |
| **Persistent learnings** | `agentic-agent learnings add` and `learnings list` carry patterns and decisions across agent boundaries. |
| **Agent-agnostic rules** | `.agentic/agent-rules/base.md` defines behavior rules that apply to every agent equally. |
| **Validation independence** | `agentic-agent validate` checks the file system, not the agent. Work done by any tool is validated the same way. |

## State Flow Across Phases

```
Phase 1 (Claude CLI)          Phase 2 (Claude VSCode)       Phase 3 (Copilot)            Phase 4 (Claude CLI)
─────────────────────          ──────────────────────         ─────────────────            ─────────────────────

task create                    cat rolling-summary.md         task show TASK-...           task list
task decompose                 context build --task ...       learnings list                 -> all in done
task claim                     learnings list                 cat src/auth/context.md      learnings list
context generate src/auth        -> sees 2 learnings           -> sees 4 learnings            -> sees 5 learnings

[implement jwt.go]             [implement middleware.go]      [implement auth_test.go]     validate
                                                                                            -> all pass
learnings add (x2)             context generate src/auth      context generate tests
update rolling-summary         learnings add (x2)             learnings add (x1)
                               validate                       task complete
                               update rolling-summary         validate

Disk state:                    Disk state:                    Disk state:
  in-progress.yaml (task)        in-progress.yaml (task)        done.yaml (task)
  progress.txt (2 learnings)     progress.txt (4 learnings)     progress.txt (5 learnings)
  src/auth/jwt.go                src/auth/middleware.go          tests/auth_test.go
  src/auth/context.md            src/auth/context.md (updated)  tests/context.md
```

## Key Takeaway

The developer never had to explain the project state to a new agent. Each agent read the same `.agentic/` files and continued seamlessly. The `agentic-agent` CLI was the constant across all three tools — the bridge that made agent switching a non-event.
