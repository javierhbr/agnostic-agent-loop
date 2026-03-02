# TechLead's Tools & Environment

## CLI

- **Binary:** `agentic-agent` (on PATH)
- **Config:** `agnostic-agent.yaml` at project root

## Key Directories

- `.agentic/tasks/` — backlog.yaml, in-progress.yaml, done.yaml
- `.agentic/coordination/` — reservations.yaml, announcements.yaml, kill-signals.yaml
- `.agentic/spec/` — change proposals and specs (ProductLead manages these)
- `.agentic/context/` — global-context.md, tech-stack.md

## Command Cheat Sheet

### Task Management

```bash
agentic-agent task list --no-interactive        # View all tasks
agentic-agent task claim <ID> --no-interactive  # Reserve task
agentic-agent context build --task <ID>         # Get full context bundle
agentic-agent task complete <ID>                # Mark done
```

### Quality Gates

```bash
agentic-agent validate                          # Run all validators
agentic-agent sdd gate-check <spec-id>          # Run 5 SDD gates
```

### OpenSpec

```bash
agentic-agent openspec check <id>               # Spec readiness check
```

## Workspace Management

The `agentic-agent` CLI is single-project and cwd-relative. Coordination across projects is configurable via environment variable.

### Coordination Directory

Set the location where both agents find project registry and coordination files:

```bash
# Option 1: Set env var (recommended)
export COORDINATION_DIR=~/app-factory/coordinators
export COORDINATION_DIR=~/.agents/coordination
export COORDINATION_DIR=/opt/projects/metadata

# Option 2: Auto-discovery (if COORDINATION_DIR not set)
# Agent searches: ../PROJECTS.md → ../../PROJECTS.md → asks you to set env var
```

### Project Registry

Located in `$COORDINATION_DIR/`:

- `PROJECTS.md` — All known projects (read at session start)
- `active-project.yaml` — Current workspace (optional, for tracking)
- `USER.md` — Shared user context (both agents read)

### Switching Between Projects

```bash
# Step 1: cd to the target project root
cd ~/app-factory/trailmate

# Step 2: Verify orientation
agentic-agent status

# Step 3: Update active-project.yaml (optional)
cat > $COORDINATION_DIR/active-project.yaml << 'EOF'
project_id: proj-trailmate
root: ~/app-factory/trailmate
name: Trailmate
switched_at: 2026-03-01T10:00:00Z
EOF
```

### Running CLI Commands in Any Project

```bash
cd <project-root> && agentic-agent task list
cd <project-root> && agentic-agent context build --task <ID>
cd <project-root> && agentic-agent sdd gate-check <spec-id>
```

## Coordination YAML Paths

- Announcements: `.agentic/coordination/announcements.yaml` (shared across all projects)
- Kill signals: `.agentic/coordination/kill-signals.yaml`

### Announcement Format (with Multi-Project Support)

```yaml
announcements:
  - from_agent: product-lead
    to_agent: tech-lead
    project_id: proj-001              # Which workspace this belongs to
    task_id: spec-auth-login-001
    status: spec-ready
    summary: "Auth login spec approved. Tasks created in backlog."
    data:
      spec_path: .agentic/spec/auth-login/proposal.md
      task_count: 5
      priority: high
```

**Always include `project_id` so you know which project an announcement belongs to.**

## Worker Spawn Format

```
Task: TASK-ID
[decoded context bundle]
Scope: path/to/file.go
Accept criteria: [from task spec]
When done: <promise>COMPLETE</promise>
```
