# Spec-Driven Development Guide

Complete guide for using Spec Kit, OpenSpec, and native specs with the agentic-agent CLI.

---

## Overview

Agentic-agent supports **spec-driven development** where specifications drive task creation, context generation, and execution. Three spec sources are supported:

| Framework | Directory | Strengths |
|-----------|-----------|-----------|
| **Spec Kit** | `.specify/specs/` | AI-assisted planning, constitution, structured specs |
| **OpenSpec** | `openspec/specs/` | Lightweight delta specs, fluid iteration, change isolation |
| **Native** | `.agentic/spec/` | Simple markdown specs, no external tools needed |

All three work through the same mechanism: **configure paths, reference specs, read on-demand**.

---

## Configuration

### Multi-Directory Spec Resolution

Configure `specDirs` in `agnostic-agent.yaml` to search multiple spec directories:

```yaml
# agnostic-agent.yaml
project:
  name: "My Project"
  version: 0.1.0
  roots:
    - .

paths:
  specDirs:
    - .specify/specs       # Spec Kit (priority 1)
    - openspec/specs        # OpenSpec (priority 2)
    - .agentic/spec         # Native fallback (priority 3)
  contextDirs:
    - .agentic/context

workflow:
  validators:
    - context-check
    - task-scope
```

When a task references `auth/spec.md`, the resolver searches each directory in order and uses the first match. If no config file exists, the default is `.agentic/spec`.

### Spec Resolution Rules

1. **Absolute/relative path** — used directly if the file exists
2. **Search specDirs** — each directory checked in configured order, first match wins
3. **No caching** — specs are read fresh from disk every time they're needed

This means specs are always up-to-date. If Spec Kit or OpenSpec updates a spec, the next read picks up the change automatically.

---

## Context System

### Three Mechanisms

**1. Agent Rules** — `.agentic/agent-rules/base.md` instructs agents to read `AGENTS.md` before working in any directory. Created by `agentic-agent init`.

**2. Context Bundles** — `agentic-agent context build --task TASK-001` aggregates everything the agent needs:

```
Context Bundle:
├── Global Context (.agentic/context/global-context.md)
├── Rolling Summary (.agentic/context/rolling-summary.md)
├── Task Details (from backlog.yaml)
├── Resolved Specs (from task spec_refs)
└── Directory Contexts (context.md from task scope)
```

**3. Validation** — `agentic-agent validate` checks that `AGENTS.md` exists in each source directory and flags when files change without a context update.

### Bundle Formats

```bash
agentic-agent context build --task TASK-001 --format toon      # Compact
agentic-agent context build --task TASK-001 --format markdown   # Readable
agentic-agent context build --task TASK-001 --format json       # Structured
```

### Context Generation

```bash
# Generate context for a directory
agentic-agent context generate internal/models

# Scan for missing context files
agentic-agent context scan
```

---

## Task Management

### Creating Tasks with Spec References

```bash
agentic-agent task create \
  --title "Create User model" \
  --spec-refs "001-auth/spec.md,001-auth/plan.md" \
  --outputs "internal/models/user.go" \
  --acceptance "User struct defined,Password hashing with bcrypt"
```

The `--spec-refs` values are resolved through the configured `specDirs` at claim time and bundle time.

### Task Ordering and Dependencies

Tasks execute in **FIFO order** from `backlog.yaml`. Dependencies are expressed through `inputs` and `outputs`:

```yaml
# backlog.yaml
tasks:
  - id: TASK-001
    title: Create User model
    outputs:
      - internal/models/user.go

  - id: TASK-002
    title: Create UserRepository
    inputs:
      - internal/models/user.go    # Depends on TASK-001
    outputs:
      - internal/repository/user_repo.go
```

### Readiness Checks

When claiming a task, readiness checks run automatically:

```bash
agentic-agent task claim TASK-001
```

```
Task TASK-001: READY
  [+] spec-resolvable: spec "001-auth/spec.md" resolved at .specify/specs/001-auth/spec.md
  [+] inputs-exist: all input files present
  [!] scope-dir: directory "internal/models" does not exist (warning only)
Claimed task TASK-001
```

Checks performed:
- **inputs-exist** (blocking) — all input files must exist
- **spec-resolvable** (blocking) — all spec refs must resolve
- **scope-dir** (warning) — scope directories checked but missing ones don't block

### Spec Commands

```bash
# List all specs across configured directories
agentic-agent spec list

# Resolve a spec ref and print its content
agentic-agent spec resolve "001-auth/spec.md"
```

---

## Workflows

### Spec Kit Workflow

```bash
# 1. Plan with Spec Kit
/speckit.constitution
/speckit.specify Build authentication system
/speckit.plan Use Go, PostgreSQL, JWT
/speckit.tasks

# 2. Create tasks referencing Spec Kit specs
agentic-agent task create \
  --title "Create User model" \
  --spec-refs "001-auth/spec.md" \
  --outputs "internal/models/user.go"

# 3. Verify spec resolution
agentic-agent spec list
agentic-agent spec resolve "001-auth/spec.md"

# 4. Execute
agentic-agent task claim TASK-001
agentic-agent context generate internal/models
# ... implement ...
agentic-agent validate
agentic-agent task complete TASK-001
```

### OpenSpec CLI Workflow

The `openspec` command group handles the full change lifecycle from a requirements file:

```bash
# 1. Initialize a change from requirements
agentic-agent openspec init "Add Authentication" --from docs/auth-requirements.md

# 2. Fill in proposal.md and tasks.md (the agent does this)

# 3. Import tasks into the backlog
agentic-agent openspec import add-authentication

# 4. Execute tasks sequentially
agentic-agent openspec status add-authentication   # Check progress
agentic-agent task claim TASK-001
# ... implement ...
agentic-agent task complete TASK-001

# 5. Complete and archive when all tasks are done
agentic-agent openspec complete add-authentication
agentic-agent openspec archive add-authentication
```

### OpenSpec Manual Workflow

You can also reference OpenSpec specs directly with `--spec-refs`:

```bash
# 1. Plan with OpenSpec
/opsx:new add-authentication
/opsx:ff

# 2. Create tasks referencing OpenSpec specs
agentic-agent task create \
  --title "Create User model" \
  --spec-refs "auth/spec.md" \
  --outputs "internal/models/user.go"

# 3. Execute
agentic-agent task claim TASK-001
agentic-agent context generate internal/models
# ... implement ...
agentic-agent validate
agentic-agent task complete TASK-001

# 4. Verify and archive with OpenSpec
/opsx:verify
/opsx:archive
```

### Native Spec Workflow

```bash
# 1. Create specs manually in .agentic/spec/
# 2. Create tasks
agentic-agent task create \
  --title "Implement auth middleware" \
  --spec-refs "auth-requirements.md"

# 3. Execute
agentic-agent task claim TASK-001
agentic-agent context build --task TASK-001
# ... implement ...
agentic-agent task complete TASK-001
```

### Ralph PDR Integration

Ralph Plan-Do-Review methodology works alongside any spec source:

```bash
# PLAN — use Spec Kit, OpenSpec, or native specs
# DO — execute with context and validation
agentic-agent task claim TASK-001
agentic-agent context generate internal/auth
# ... implement ...
agentic-agent validate
# REVIEW — track progress and learnings
agentic-agent task complete TASK-001
agentic-agent learnings add "User model requires bcrypt for password hashing"
agentic-agent learnings show
```

---

## Autopilot Mode

Autopilot processes backlog tasks sequentially: readiness check, claim, generate context, build bundle.

### Dry Run

```bash
agentic-agent autopilot start --dry-run
```

```
--- Iteration 1/10 ---
Next task: [TASK-001] Create User model
Task TASK-001: READY
  [+] spec-resolvable: spec "001-auth/spec.md" resolved
  [+] inputs-exist: all inputs present
[DRY RUN] Would claim task TASK-001 and generate context

--- Iteration 2/10 ---
Next task: [TASK-002] Create UserRepository
...
```

### Run

```bash
agentic-agent autopilot start --max-iterations 5
```

Per iteration, autopilot:
1. Finds the next claimable task (prefers tasks where all readiness checks pass)
2. Prints readiness check results
3. Claims the task
4. Generates context for each scope directory
5. Builds a context bundle with resolved specs
6. Reports the task as ready for agent execution

### Stop Conditions

- All backlog tasks processed
- `--max-iterations` limit reached
- Ctrl+C

---

## Directory Structures

### Spec Kit Project

```
project/
├── .specify/
│   ├── memory/
│   │   └── constitution.md         # Project principles
│   └── specs/
│       └── 001-auth/
│           ├── spec.md             # Feature spec
│           ├── plan.md             # Technical plan
│           └── tasks.md            # Task breakdown
├── .agentic/
│   ├── spec/                       # Native specs (fallback)
│   ├── tasks/
│   │   ├── backlog.yaml
│   │   ├── in-progress.yaml
│   │   └── done.yaml
│   ├── context/
│   │   ├── global-context.md
│   │   └── rolling-summary.md
│   └── agent-rules/
│       └── base.md
└── agnostic-agent.yaml
```

### OpenSpec Project

```
project/
├── openspec/
│   └── specs/
│       └── auth/
│           └── spec.md             # Source-of-truth specs (manual)
├── .agentic/
│   ├── openspec/
│   │   └── changes/                # Managed by `openspec` commands
│   │       ├── changes.yaml        # Change registry
│   │       ├── add-authentication/
│   │       │   ├── proposal.md     # Seeded from requirements
│   │       │   ├── tasks.md        # Numbered task list
│   │       │   ├── metadata.yaml   # Change metadata
│   │       │   ├── specs/          # Delta specs
│   │       │   └── IMPLEMENTED     # Completion marker
│   │       └── _archive/           # Archived changes
│   ├── spec/                       # Native specs (fallback)
│   ├── tasks/
│   ├── context/
│   └── agent-rules/
└── agnostic-agent.yaml
```

---

## Framework Comparison

### Complementary Strengths

| Capability | Spec Kit | OpenSpec | Agentic Agent |
|------------|----------|----------|---------------|
| AI-assisted spec generation | Yes | Yes | No |
| Per-directory context | No | No | Yes |
| Task size enforcement | No | No | Yes |
| Delta specs (changes only) | No | Yes | No |
| Change isolation | No | Yes | No |
| Automated validation | No | Verify command | Yes |
| Interactive CLI | No | No | Yes |
| Progress tracking | No | No | Yes (Ralph) |
| Fluid iteration | Rigid phases | Yes | Flexible |

### When to Use What

- **Spec Kit** — AI-assisted planning for new projects, comprehensive specs
- **OpenSpec** — Lightweight change management, quick features, exploration
- **Agentic Agent** — Task execution, context management, validation
- **All three** — Spec Kit/OpenSpec for planning, Agentic for execution

---

## CLI Quick Reference

### Task & Spec Commands

| Command | Purpose |
|---------|---------|
| `task create --spec-refs "..."` | Create a task referencing specs |
| `task claim <id>` | Claim task (runs readiness checks) |
| `task complete <id>` | Mark task as done (captures commits) |
| `task list` | List tasks across backlog, in-progress, done |
| `task show <id>` | Display task details |
| `task decompose <id>` | Break task into subtasks |
| `spec list` | List all specs across configured directories |
| `spec resolve <ref>` | Resolve a spec ref and print its content |

### Track & Plan Commands

| Command | Purpose |
|---------|---------|
| `track create` | Create a new feature/bug track with brainstorm, spec, plan files |
| `track list` | List all tracks |
| `track show <id>` | Display track details |
| `track archive <id>` | Archive a completed track |
| `plan show` | Display parsed plan.md with phase/task breakdown |
| `plan next` | Show the next pending task in the plan |
| `plan mark <line> <status>` | Update a task's checkbox status in plan.md |

### Context & Validation Commands

| Command | Purpose |
|---------|---------|
| `context generate <dir>` | Generate context.md for a directory |
| `context build --task <id>` | Build context bundle with resolved specs |
| `context scan` | Find directories missing context.md |
| `validate` | Run validation rules |

### Skills & Status Commands

| Command | Purpose |
|---------|---------|
| `skills ensure` | Idempotent setup: generate rules, install packs, fix drift |
| `skills detect` | Show detected agent (flag/env/filesystem) |
| `skills list` | List available skill packs |
| `skills install <pack>` | Install a skill pack for the active agent |
| `status` | Show project dashboard (task counts, blockers, next ready) |

### OpenSpec Commands

| Command | Purpose |
|---------|---------|
| `openspec init <name> --from <file>` | Create change from requirements file |
| `openspec import <id>` | Import tasks.md into backlog |
| `openspec status <id>` | Show change progress (tasks done/total) |
| `openspec complete <id>` | Validate all tasks done, write IMPLEMENTED marker |
| `openspec archive <id>` | Move completed change to archive |
| `openspec list` | List all changes |
| `openspec show <id>` | Show change details |

### Autopilot & Learnings Commands

| Command | Purpose |
|---------|---------|
| `autopilot start` | Process backlog tasks sequentially |
| `autopilot start --dry-run` | Preview autopilot without changes |
| `learnings add "..."` | Record a learning (Ralph) |
| `learnings show` | View progress and learnings |

---

## Future Enhancements

The following are not yet implemented:

- **Explicit task dependencies** — `depends_on` field for task-to-task ordering
- **Priority levels** — `priority: high/medium/low` for task sorting
- **Agent spawning** — Autopilot spawns agents (Claude API, etc.) automatically
- **OpenSpec bidirectional sync** — `openspec sync` to update OpenSpec tasks.md from agentic completion status
- **Parallel task execution** — Run independent tasks concurrently
- **Cost tracking** — Budget limits for autopilot runs
