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
  track init --> brainstorm --> refine spec --> activate
       |                                          |
       |                          generate plan + decompose tasks
       |                                          |
       +--- complete <--- validate <--- claim + implement
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

# Start a track with brainstorming scaffolding
agentic-agent track init "User Authentication" \
  --purpose "Secure user login" \
  --success "Users can register, login, logout"

# Brainstorm with your AI agent using the generated brainstorm.md
# Then check spec completeness
agentic-agent track refine user-authentication

# Activate: generate plan + decompose into tasks
agentic-agent track activate user-authentication --decompose

# Work on tasks
agentic-agent task claim TASK-001       # Readiness checks + claim
agentic-agent context generate internal/auth  # Context for the agent
agentic-agent validate                  # Check scope and context
agentic-agent task complete TASK-001    # Done
```

All commands support **dual mode** — run without arguments for interactive wizards, or pass flags for scripting:

```bash
agentic-agent track init                        # Interactive wizard
agentic-agent track init "Auth" --purpose "..."  # Flag mode
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
make build
sudo mv build/agentic-agent /usr/local/bin/

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

### Tracks — from idea to implementation

Tracks group a spec, plan, and tasks into a single work unit. The brainstorming workflow guides you from a vague idea to a well-formed spec through structured dialogue with your AI agent.

```bash
# 1. Initialize with brainstorming scaffolding
agentic-agent track init "User Auth" --type feature

# 2. Use brainstorm.md as a dialogue script with your agent
#    Agent asks one question at a time, proposes 2-3 approaches,
#    presents design in sections for validation

# 3. Check spec completeness
agentic-agent track refine user-auth
#   ✓ purpose
#   ✓ constraints
#   ✗ design (missing)

# 4. Activate: validate spec → generate plan → decompose into tasks
agentic-agent track activate user-auth --decompose
```

Track status progresses through: **ideation** → **active** → **done**. Tasks linked to ideation tracks are skipped by autopilot until the track is activated.

Templates are customizable — override the defaults by placing your own in `.agentic/templates/track/`.

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
|-- Resolved Skills (targeted content from skill_refs, or all installed)
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

The same project structure works across Claude Code, Cursor, Copilot, Gemini, Windsurf, Antigravity IDE, and Codex. Generate tool-specific configs or install skill packs:

```bash
agentic-agent skills generate-claude-skills
agentic-agent skills generate-gemini-skills

# Install skill packs for any supported tool
agentic-agent skills install tdd --tool claude-code
agentic-agent skills install tdd --tool cursor --global
agentic-agent skills list
```

### Skill packs — reusable agent skills

Skill packs are tool-agnostic bundles of instructions that install into the correct directory for each AI agent tool:

| Tool | Project directory |
|------|-------------------|
| Claude Code | `.claude/skills/` |
| Cursor | `.cursor/skills/` |
| Gemini CLI | `.gemini/skills/` |
| Windsurf | `.windsurf/skills/` |
| Antigravity | `.agent/skills/` |
| Codex | `.codex/skills/` |

```bash
# Install TDD skill pack for Claude Code (project-level)
agentic-agent skills install tdd --tool claude-code

# Install globally (user-level)
agentic-agent skills install tdd --tool claude-code --global

# Interactive mode — select pack, tool, and scope
agentic-agent skills install
```

### Task-level skill refs — targeted skills per task

Instead of applying all installed skill packs to every task, declare `skill_refs` on individual tasks to include only the relevant skills:

```yaml
# .agentic/tasks/backlog.yaml
tasks:
  - id: "TASK-001"
    title: "Refactor auth middleware"
    skill_refs:
      - code-simplification
      - tdd
    scope:
      - "internal/auth"
```

When a context bundle is built for this task, only the referenced skill packs are included — not all installed packs. Tasks without `skill_refs` fall back to the existing behavior (all installed packs included).

Skill refs resolve through a 3-tier fallback:

1. Agent's installed skill directory (e.g., `.claude/skills/tdd/SKILL.md`)
2. Any installed tool's skill directory
3. Embedded pack content (compiled into the binary — always available)

### Simplify command — targeted code review

The `simplify` command generates a focused context bundle for code simplification review:

```bash
# Review specific directories
agentic-agent simplify internal/auth internal/middleware

# Review directories from a task's scope
agentic-agent simplify --task TASK-001

# Output as JSON
agentic-agent simplify . --format json --output review.json
```

The bundle includes the code-simplification skill instructions, directory context, source file listings, and tech stack information.

### TDD workflow — RED/GREEN/REFACTOR

The `work --follow-tdd` flag decomposes a task into three phased sub-tasks and verifies the TDD skill pack is installed:

```bash
# Install the TDD skill pack first
agentic-agent skills install tdd --tool claude-code

# Work on a task with TDD workflow
agentic-agent work --task TASK-001 --follow-tdd
```

This creates three sub-tasks:

1. **RED** (`TASK-001-red`) — Write failing tests
2. **GREEN** (`TASK-001-green`) — Implement minimal code to pass
3. **REFACTOR** (`TASK-001-refactor`) — Improve code quality

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

### OpenSpec CLI

The `openspec` command group handles the full change lifecycle — from a requirements file to archived implementation:

```bash
# Initialize a change from requirements
agentic-agent openspec init "Auth Feature" --from .agentic/spec/auth-requirements.md

# Fill in proposal.md and tasks.md, then import into backlog
agentic-agent openspec import auth-feature

# Track progress, complete, and archive
agentic-agent openspec status auth-feature
agentic-agent openspec complete auth-feature
agentic-agent openspec archive auth-feature
```

Tell your agent: *"Start a project from requirements.md following openspec"* — the openspec skill handles the full flow automatically.

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

### Tracks

| Command | Description |
|---------|-------------|
| `track init [name]` | Initialize track with brainstorming scaffolding |
| `track refine <id>` | Validate spec completeness |
| `track activate <id>` | Generate plan and tasks from spec |
| `track list` | List all tracks |
| `track show <id>` | Show track details |
| `track archive <id>` | Archive a completed track |

**Track init flags:** `--name`, `--type`, `--purpose`, `--constraints`, `--success`

**Track activate flags:** `--decompose` (auto-create tasks from plan)

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

### OpenSpec

| Command | Description |
|---------|-------------|
| `openspec init <name> --from <file>` | Create change from requirements file |
| `openspec import <id>` | Import tasks.md into backlog |
| `openspec list` | List all changes |
| `openspec show <id>` | Show change details and proposal excerpt |
| `openspec status <id>` | Show task progress for a change |
| `openspec complete <id>` | Validate all tasks done, write IMPLEMENTED marker |
| `openspec archive <id>` | Archive a completed change |

### Automation

| Command | Description |
|---------|-------------|
| `autopilot start` | Process backlog tasks sequentially |
| `autopilot start --dry-run` | Preview without changes |
| `run` | Run orchestrator loop |
| `work` | Interactive claim-to-complete workflow |
| `work --follow-tdd` | TDD workflow: decompose into RED/GREEN/REFACTOR |

### Validation, Skills, and Simplification

| Command | Description |
|---------|-------------|
| `validate` | Run all validation rules |
| `validate --format json` | JSON output for CI/CD |
| `skills generate-claude-skills` | Generate Claude Code config |
| `skills generate-gemini-skills` | Generate Gemini config |
| `skills install [pack]` | Install a skill pack for an agent tool |
| `skills list` | List available skill packs |
| `simplify [dir...]` | Generate code simplification review bundle |
| `simplify --task <id>` | Simplify using task scope directories |

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
  openSpecDir: .agentic/openspec/changes  # OpenSpec change lifecycle
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
|-- tracks/              # Work units (spec + plan + tasks)
|   +-- user-auth/
|       |-- brainstorm.md    # Agent dialogue script
|       |-- spec.md          # Enhanced specification
|       |-- plan.md          # Phased implementation plan
|       +-- metadata.yaml
|-- openspec/            # OpenSpec change lifecycle
|   +-- changes/
|       +-- auth-feature/
|           |-- proposal.md      # Change proposal
|           |-- tasks.md         # Implementation tasks
|           |-- specs/           # Change-specific specs
|           +-- metadata.yaml
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
|   |-- gitops/               # Read-only git integration
|   |-- orchestrator/         # Autopilot and loop orchestration
|   |-- plans/                # Plan parsing + generation from specs
|   |-- project/              # Project init + track templates
|   |-- openspec/             # OpenSpec change lifecycle manager
|   |-- specs/                # Multi-directory spec resolution
|   |-- status/               # Aggregated project status
|   |-- simplify/             # Code simplification bundle builder
|   |-- skills/               # Skill packs, installer, registry, resolver
|   |-- tasks/                # Task management + decomposition
|   |-- tracks/               # Track management + spec validation
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
go test ./internal/openspec -v        # OpenSpec change lifecycle
go test ./internal/specs -v           # Spec resolution
go test ./internal/tasks -v           # Task management + decomposition + TDD
go test ./internal/skills -v          # Skill packs, installer, resolver
go test ./internal/simplify -v        # Code simplification bundles
go test ./internal/tracks -v          # Track management + spec validation
go test ./internal/plans -v           # Plan parsing + generation
go test ./internal/orchestrator -v    # Autopilot + loop
go test ./internal/validator/rules -v # Validation rules
go test ./test/functional -v          # Functional CLI tests
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
