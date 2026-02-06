# Spec-Driven Workflow: Spec Kit, OpenSpec & Autopilot

This example demonstrates how to use specification files from Spec Kit, OpenSpec, or native agentic specs — all resolved automatically through multi-directory configuration. It also shows autopilot mode for sequential task processing.

---

## How It Works

The key insight: **you don't import specs — you reference them.** Configure `specDirs` in `agnostic-agent.yaml` to point at your spec directories. The resolver searches them in order and uses the first match.

```yaml
# agnostic-agent.yaml
paths:
  specDirs:
    - .specify/specs       # Spec Kit
    - openspec/specs        # OpenSpec
    - .agentic/spec         # Agentic native (fallback)
```

When a task references `auth-requirements.md`, the resolver searches:
1. `.specify/specs/auth-requirements.md`
2. `openspec/specs/auth-requirements.md`
3. `.agentic/spec/auth-requirements.md` — found here

---

## Part A: Spec Kit Integration

### 1. Plan with Spec Kit

Spec Kit generates specs in `.specify/specs/`. After planning:

```
.specify/
└── specs/
    └── 001-auth/
        └── spec.md       # Feature spec with requirements and scenarios
```

### 2. Create tasks referencing Spec Kit specs

```bash
agentic-agent task create \
  --title "Implement JWT token service" \
  --spec-refs "001-auth/spec.md" \
  --outputs "internal/auth/jwt.go" \
  --acceptance "Token generation works,Validation rejects expired tokens"
```

The spec ref `001-auth/spec.md` resolves to `.specify/specs/001-auth/spec.md` because `.specify/specs` is first in `specDirs`.

### 3. Verify spec resolution

```bash
# List all specs across configured directories
agentic-agent spec list
```

```
auth-requirements.md  /path/to/.agentic/spec/auth-requirements.md
api-design.md         /path/to/.agentic/spec/api-design.md
001-auth/spec.md      /path/to/.specify/specs/001-auth/spec.md
auth/spec.md          /path/to/openspec/specs/auth/spec.md
```

```bash
# Resolve and read a specific spec
agentic-agent spec resolve "001-auth/spec.md"
```

Outputs the full content of the Spec Kit spec.

### 4. Claim task with readiness checks

```bash
agentic-agent task claim SPEC-001
```

```
Task SPEC-001: READY
  [+] spec-resolvable: spec "auth-requirements.md" resolved at .agentic/spec/auth-requirements.md
  [+] spec-resolvable: spec "api-design.md" resolved at .agentic/spec/api-design.md
Claimed task SPEC-001
```

Readiness checks verify that all referenced specs exist before claiming. If a spec is missing, the claim still proceeds but you see a warning.

### 5. Build context bundle with specs

```bash
agentic-agent context build --task SPEC-001
```

The output includes a `specs:` section with the full content of each resolved spec. The agent receives everything it needs — task definition, project context, and specification content — in a single bundle.

### 6. Implement and complete

```bash
# Work on the task...
agentic-agent validate
agentic-agent task complete SPEC-001
```

---

## Part B: OpenSpec Integration

### 1. Plan with OpenSpec

OpenSpec generates specs in `openspec/specs/`. After running `/opsx:new add-auth` and `/opsx:ff`:

```
openspec/
└── specs/
    └── auth/
        └── spec.md       # OpenSpec spec with proposal, design, tasks
```

### 2. Add OpenSpec directory to config

Already configured in `agnostic-agent.yaml`:

```yaml
paths:
  specDirs:
    - .specify/specs       # Spec Kit
    - openspec/specs        # OpenSpec ← specs found here
    - .agentic/spec         # Fallback
```

### 3. Create tasks referencing OpenSpec specs

```bash
agentic-agent task create \
  --title "Create User model and repository" \
  --spec-refs "auth/spec.md" \
  --outputs "internal/models/user.go,internal/repository/user_repo.go"
```

The ref `auth/spec.md` resolves to `openspec/specs/auth/spec.md`.

### 4. Execute the same way

```bash
# Verify the spec resolves
agentic-agent spec resolve "auth/spec.md"

# Claim with readiness checks
agentic-agent task claim TASK-001

# Build context (includes OpenSpec content)
agentic-agent context build --task TASK-001

# Work, validate, complete
agentic-agent validate
agentic-agent task complete TASK-001
```

### 5. Verify with OpenSpec after completion

```
/opsx:verify    # OpenSpec validates against its specs
/opsx:archive   # Archive the completed change
```

For the full spec-driven development guide, see [docs/SPEC_DRIVEN_DEVELOPMENT.md](../../docs/SPEC_DRIVEN_DEVELOPMENT.md).

---

## Part C: Autopilot Mode

Autopilot processes backlog tasks sequentially: readiness check, claim, generate context, build bundle.

### 1. Preview with dry run

```bash
agentic-agent autopilot start --dry-run
```

```
--- Iteration 1/10 ---
Next task: [SPEC-001] Create JWT token service
Task SPEC-001: READY
  [+] spec-resolvable: spec "auth-requirements.md" resolved at .agentic/spec/auth-requirements.md
  [+] spec-resolvable: spec "api-design.md" resolved at .agentic/spec/api-design.md
[DRY RUN] Would claim task SPEC-001 and generate context

--- Iteration 2/10 ---
Next task: [SPEC-002] Implement auth middleware
Task SPEC-002: READY
  [+] spec-resolvable: spec "auth-requirements.md" resolved at .agentic/spec/auth-requirements.md
[DRY RUN] Would claim task SPEC-002 and generate context
```

### 2. Run autopilot

```bash
# Process up to 3 tasks
agentic-agent autopilot start --max-iterations 3
```

Per iteration, autopilot:
1. Finds the next claimable task (prefers tasks where all readiness checks pass)
2. Prints readiness check results
3. Claims the task
4. Generates context for each scope directory
5. Builds a context bundle with resolved specs
6. Reports the task as ready for agent execution

### 3. Stop conditions

Autopilot stops when:
- All backlog tasks are processed
- `--max-iterations` limit is reached
- You press Ctrl+C

---

## Directory Structure

```
spec-driven-workflow/
├── README.md                          # This file
├── agnostic-agent.yaml                # Multi-dir spec config
├── .agentic/
│   ├── spec/
│   │   ├── auth-requirements.md       # Native spec (auth requirements)
│   │   └── api-design.md             # Native spec (API design)
│   ├── tasks/
│   │   ├── backlog.yaml              # 2 tasks referencing specs
│   │   ├── in-progress.yaml
│   │   └── done.yaml
│   └── context/
│       ├── global-context.md
│       └── rolling-summary.md
├── .specify/
│   └── specs/
│       └── 001-auth/
│           └── spec.md               # Spec Kit sample spec
└── openspec/
    └── specs/
        └── auth/
            └── spec.md               # OpenSpec sample spec
```

## Quick Reference

| Command | Purpose |
|---------|---------|
| `spec list` | List all specs across all configured directories |
| `spec resolve <ref>` | Resolve a spec ref and print its content |
| `task create --spec-refs "..."` | Create a task that references specifications |
| `task claim <id>` | Claim task (runs readiness checks first) |
| `context build --task <id>` | Build context bundle including resolved specs |
| `autopilot start --dry-run` | Preview autopilot without making changes |
| `autopilot start` | Run autopilot to process backlog tasks |

## Related Documentation

- [Spec-Driven Development Guide](../../docs/SPEC_DRIVEN_DEVELOPMENT.md)
