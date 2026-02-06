# Agentic Agent

A CLI that gives AI agents the structure they need to ship real features — specs, context, tasks, and validation in one workflow.

---

## The Problem

AI coding agents are powerful but chaotic. Without structure, they:

- **Lose context** — agents don't know what other modules do, what was decided last session, or what constraints exist. They reinvent, contradict, and break things.
- **Sprawl across the codebase** — a "simple feature" touches 20 files across 8 directories. No human can review that. No agent should attempt it.
- **Skip validation** — agents write code but don't verify it meets requirements, follows conventions, or stays within scope.
- **Lock you into one tool** — your workflow shouldn't break when you switch from Claude Code to Cursor to Copilot.

Agentic Agent solves this by giving every AI agent the same structured workflow: specs define what to build, context keeps agents focused, tasks enforce atomic scope, and validation catches drift.

---

## How It Works

```
                         Your Project
                             |
           +-----------------+-----------------+
           |                 |                 |
       Specs             Context            Tasks
   (source of truth)  (per-directory)   (atomic units)
           |                 |                 |
           +--------+--------+--------+--------+
                    |                 |
              Claim + Build       Validate
            (readiness checks)   (scope, context)
                    |                 |
                    +--------+--------+
                             |
                     Complete + Learn
```

The core loop:

```
  init --> write specs --> create tasks --> claim task
   ^                                          |
   |                                    generate context
   |                                          |
   +--- complete <--- validate <--- implement
```

Every command works with **any AI agent** — Claude Code, Cursor, GitHub Copilot, Gemini, or manual development.

---

## Quick Start

```bash
# Install
go install github.com/javierbenavides/agentic-agent/cmd/agentic-agent@latest

# Initialize a project
mkdir my-project && cd my-project
agentic-agent init --name "My Project"

# Create a task
agentic-agent task create \
  --title "Add user authentication" \
  --spec-refs "auth-requirements.md" \
  --outputs "internal/auth/service.go" \
  --acceptance "JWT tokens generated,Expired tokens rejected"

# Work on it
agentic-agent task claim TASK-001       # Readiness checks + claim
agentic-agent context generate internal/auth  # Context for the agent
agentic-agent validate                  # Check scope and context
agentic-agent task complete TASK-001    # Done
```

All commands support **dual mode** — run without arguments for interactive wizards, or pass flags for scripting:

```bash
agentic-agent task create                # Interactive wizard
agentic-agent task create --title "..."  # Flag mode
```

---

## Installation

**Prerequisites:** Go 1.22+, Git

```bash
# Option 1: Go install (recommended)
go install github.com/javierbenavides/agentic-agent/cmd/agentic-agent@latest

# Option 2: Build from source
git clone https://github.com/javierbenavides/agentic-agent.git
cd agentic-agent
go build -o agentic-agent ./cmd/agentic-agent
sudo mv agentic-agent /usr/local/bin/

# Option 3: Development (no install)
go run ./cmd/agentic-agent version
```

---

## Core Concepts

### Specifications — the source of truth

Specs live in `.agentic/spec/` (or Spec Kit / OpenSpec directories). Agents reference them, never serve as the source of truth themselves.

```bash
agentic-agent task create --spec-refs "auth-requirements.md,api-design.md"
```

### Context isolation — agents stay focused

Each directory gets a `context.md` file describing its purpose, dependencies, and constraints. The CLI bundles only the relevant context for each task, so agents don't drown in irrelevant information.

```bash
agentic-agent context generate internal/auth   # Generate context.md
agentic-agent context build --task TASK-001    # Bundle: specs + context + task
```

The context bundle sent to an agent includes:

```
Context Bundle for TASK-001
|-- Global Context (project-wide)
|-- Rolling Summary (recent changes)
|-- Task Details (title, acceptance criteria, scope)
|-- Resolved Specs (full content from spec_refs)
+-- Directory Contexts (context.md from task scope)
```

### Atomic tasks — small, reviewable units

Tasks are capped at **5 files** and **2 directories**. Larger work gets decomposed:

```bash
agentic-agent task decompose TASK-001 \
  "Create JWT service" \
  "Add auth middleware" \
  "Write integration tests"
```

### Readiness checks — verify before starting

When claiming a task, the CLI checks that inputs exist, specs resolve, and scope directories are present:

```
$ agentic-agent task claim TASK-001

Task TASK-001: READY
  [+] spec-resolvable: "auth-requirements.md" resolved at .agentic/spec/auth-requirements.md
  [+] inputs-exist: all input files present
  [!] scope-dir: "internal/auth" does not exist (warning)
Claimed task TASK-001
```

### Agent-agnostic — works with everything

The same project structure works across Claude Code, Cursor, Copilot, Gemini, and Antigravity IDE. Generate tool-specific configs:

```bash
agentic-agent skills generate-claude-skills
agentic-agent skills generate-gemini-skills
```

---

## Spec-Driven Development

Configure `specDirs` to search multiple spec sources in priority order:

```yaml
# agnostic-agent.yaml
paths:
  specDirs:
    - .specify/specs       # Spec Kit
    - openspec/specs        # OpenSpec
    - .agentic/spec         # Native (fallback)
```

When a task references `auth/spec.md`, the resolver searches each directory in order and uses the first match. Specs are read fresh from disk every time — no caching, no sync.

```bash
# List all specs across configured directories
agentic-agent spec list

# Resolve and print a specific spec
agentic-agent spec resolve "auth/spec.md"
```

See [docs/SPEC_DRIVEN_DEVELOPMENT.md](docs/SPEC_DRIVEN_DEVELOPMENT.md) for Spec Kit, OpenSpec, and native spec workflows.

---

## Autopilot Mode

Process backlog tasks sequentially — readiness checks, claim, context generation, and bundling:

```bash
# Preview without making changes
agentic-agent autopilot start --dry-run

# Process up to 5 tasks
agentic-agent autopilot start --max-iterations 5
```

```
--- Iteration 1/5 ---
Next task: [TASK-001] Create JWT token service
Task TASK-001: READY
  [+] spec-resolvable: "auth-requirements.md" resolved
  [+] inputs-exist: all inputs present
Claimed TASK-001
Generated context for internal/auth
Built context bundle (toon format)

--- Iteration 2/5 ---
Next task: [TASK-002] Implement auth middleware
...
```

Autopilot stops when all tasks are processed, `--max-iterations` is reached, or you press Ctrl+C.

---

## CLI Reference

### Project

| Command | Description |
|---------|-------------|
| `init --name "Name"` | Initialize project structure |
| `start` | Interactive setup wizard |
| `version` | Print version |

### Tasks

| Command | Description |
|---------|-------------|
| `task create [flags]` | Create a task (wizard or flags) |
| `task list` | List all tasks by status |
| `task show <id>` | Show task details |
| `task claim <id>` | Claim task with readiness checks |
| `task complete <id>` | Mark task as done |
| `task decompose <id> ...` | Break into subtasks |
| `task from-template` | Create from template (wizard) |
| `task sample-task` | Create a sample task |

**Task create flags:** `--title`, `--description`, `--spec-refs`, `--inputs`, `--outputs`, `--acceptance`

### Context

| Command | Description |
|---------|-------------|
| `context generate <dir>` | Generate context.md for a directory |
| `context scan` | Find directories missing context |
| `context build --task <id>` | Build context bundle with resolved specs |

### Specs

| Command | Description |
|---------|-------------|
| `spec list` | List all specs across configured directories |
| `spec resolve <ref>` | Resolve a spec ref and print content |

### Automation

| Command | Description |
|---------|-------------|
| `autopilot start` | Process backlog tasks sequentially |
| `autopilot start --dry-run` | Preview without changes |
| `run` | Run orchestrator loop |
| `work` | Interactive claim-to-complete workflow |

### Validation and Skills

| Command | Description |
|---------|-------------|
| `validate` | Run all validation rules |
| `validate --format json` | JSON output for CI/CD |
| `skills generate-claude-skills` | Generate Claude Code config |
| `skills generate-gemini-skills` | Generate Gemini config |

### Progress Tracking

| Command | Description |
|---------|-------------|
| `learnings add "..."` | Record a learning |
| `learnings list` | List learnings |
| `learnings show` | Show progress and learnings |
| `token status` | Show token usage stats |

---

## Configuration

The CLI reads `agnostic-agent.yaml` from the current working directory (created by `agentic-agent init`). Use `--config` for a custom path.

```yaml
project:
  name: "My Project"
  version: 0.1.0
  roots:
    - .

agents:
  defaults:
    max_tokens: 4000
    model: claude-3-5-sonnet-20241022

paths:
  specDirs:
    - .specify/specs         # Spec Kit
    - openspec/specs          # OpenSpec
    - .agentic/spec           # Native fallback
  contextDirs:
    - .agentic/context
  prdOutputPath: .agentic/tasks/
  progressTextPath: .agentic/progress.txt
  progressYAMLPath: .agentic/progress.yaml
  archiveDir: .agentic/archive/

workflow:
  validators:
    - context-check
    - task-scope
    - browser-verification
```

---

## Project Structure

```
.agentic/
|-- spec/                # Specifications (source of truth)
|-- context/             # Global and rolling context
|   |-- global-context.md
|   +-- rolling-summary.md
|-- tasks/               # Task lifecycle files
|   |-- backlog.yaml
|   |-- in-progress.yaml
|   +-- done.yaml
+-- agent-rules/         # Tool-specific agent configs
    +-- base.md
agnostic-agent.yaml      # Project configuration
```

Source directories each have a `context.md`:

```
internal/auth/
|-- context.md           # Module context for agents
|-- service.go
+-- service_test.go
```

### Source Code Layout

```
agentic-agent/
|-- cmd/agentic-agent/        # CLI commands (Cobra + Bubble Tea)
|-- internal/
|   |-- config/               # Configuration loading
|   |-- context/              # Context generation
|   |-- encoding/             # Context bundle encoding
|   |-- orchestrator/         # Autopilot and loop orchestration
|   |-- project/              # Project initialization
|   |-- specs/                # Multi-directory spec resolution
|   |-- tasks/                # Task management + readiness checks
|   +-- validator/            # Validation rules
|-- pkg/models/               # Data models
|-- test/
|   |-- bdd/                  # BDD tests (godog/Gherkin)
|   |-- functional/           # Functional CLI tests
|   +-- integration/          # Integration tests
+-- examples/                 # Example projects
```

---

## Testing

```bash
# Run all tests
go test ./...

# With coverage
go test ./... -cover

# HTML coverage report
make coverage-html

# By package
go test ./internal/specs -v           # Spec resolution
go test ./internal/tasks -v           # Task management + readiness
go test ./internal/orchestrator -v    # Autopilot + loop
go test ./internal/validator/rules -v # Validation rules
go test ./test/integration -v         # Integration tests
go test ./test/bdd -v                 # BDD/Gherkin tests
```

---

## Documentation

- [Spec-Driven Development Guide](docs/SPEC_DRIVEN_DEVELOPMENT.md) — Spec Kit, OpenSpec, and native spec workflows
- [CLI Tutorial](docs/guide/CLI_TUTORIAL.md) — Step-by-step scenarios
- [BDD Testing Guide](docs/bdd/BDD_GUIDE.md) — Gherkin feature files and godog
- [Multi-Agent Workflow](examples/multi-agent-workflow/MULTI_AGENT_USE_CASE.md) — Switching between Claude Code, Cursor, Copilot, and Gemini
- [Spec-Driven Workflow Example](examples/spec-driven-workflow/README.md) — End-to-end example with Spec Kit, OpenSpec, and autopilot
- [BDD Infrastructure](test/bdd/README.md) — Test infrastructure overview

---

## Troubleshooting

**`agentic-agent: command not found`** — Add the binary to your PATH or use `go run ./cmd/agentic-agent`.

**Validation fails with "Missing context.md"** — Run `agentic-agent context scan` to find directories missing context, then `agentic-agent context generate <dir>`.

**Task size validation fails** — Decompose with `agentic-agent task decompose TASK-001 "Subtask 1" "Subtask 2"`.

**Task claim fails with "not found in backlog"** — Check task status with `agentic-agent task list`. Tasks can only be claimed from backlog.

---

## Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests for your changes
4. Ensure all tests pass: `go test ./...`
5. Run validation: `agentic-agent validate`
6. Submit a pull request

---

## Built With

- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — Interactive TUI
- [YAML v3](https://github.com/go-yaml/yaml) — YAML parsing
- [godog](https://github.com/cucumber/godog) — BDD testing
- [testify](https://github.com/stretchr/testify) — Test assertions

## License

[Add your license here]

## Support

- GitHub Issues: [Create an issue](https://github.com/javierbenavides/agentic-agent/issues)
