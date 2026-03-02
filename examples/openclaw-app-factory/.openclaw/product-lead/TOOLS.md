# ProductLead's Tools & Environment

## CLI

- **Binary:** `agentic-agent` (on PATH)
- **Config:** `agnostic-agent.yaml` at project root

## Key Directories

- `.agentic/spec/` — change proposals, tech stacks, task files per change
- `.agentic/context/` — global-context.md, tech-stack.md, decisions.md
- `.agentic/tasks/` — backlog (read-only for ProductLead; TechLead manages this)

## Command Cheat Sheet

### Spec Lifecycle

```bash
agentic-agent openspec init "<name>" --from <requirements-file>   # Create new change
agentic-agent openspec check <id>                                  # Validate spec readiness
```

### Product Skills (trigger in conversation)

```
product-wizard      → Generate PRD from raw idea
openspec            → Full spec-driven development lifecycle
sdd analyst         → Discovery and requirements analysis
sdd architect       → Architecture and component design
```

### Context & Decisions

```bash
# Update decisions log
echo "Decision: [title] — [rationale]" >> .agentic/context/decisions.md
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
cd <project-root> && agentic-agent openspec init "<name>"
cd <project-root> && agentic-agent openspec check <id>
```

## Coordination YAML Path

- Announcements: `.agentic/coordination/announcements.yaml`

## Spec-Ready Announcement Format

```yaml
announcements:
  - from_agent: product-lead
    to_agent: tech-lead
    project_id: proj-001              # Which workspace this belongs to
    task_id: <spec-id>
    status: spec-ready
    summary: "<spec name> approved. Tasks created in backlog."
    data:
      spec_path: .agentic/spec/<id>/proposal.md
      task_count: <n>
      priority: high|medium|low
```

**Always include `project_id` so TechLead knows which project to switch to.**
