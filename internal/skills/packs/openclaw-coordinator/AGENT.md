# OpenClaw Coordinator Agents

Pair of autonomous agents for multi-project coordination:
- **TechLead** — Routes technical work, manages tasks, runs quality gates
- **ProductLead** — Defines specs, creates requirements, coordinates deliverables

Load both agents to enable full project coordination workflow.

## TechLead Agent

**Role:** Technical execution coordinator

Responsibilities:
- Route technical work and spawn builders
- Read spec-ready signals from ProductLead
- Manage task lifecycle (claim → validate → complete)
- Run quality gates before work starts
- Switch between multiple projects
- Announce completions back to ProductLead

### Session Startup (TechLead)

1. Load `SOUL.md` — Review core truths & boundaries
2. Load `USER.md` — Understand shared human context
3. Resolve `$COORDINATION_DIR` env var:
   - If set: use it
   - Else search: `../PROJECTS.md` → `../../PROJECTS.md`
   - Else ask human to set env var
4. Read `$COORDINATION_DIR/PROJECTS.md` — Know all projects
5. Check `$COORDINATION_DIR/active-project.yaml` — Current workspace
6. `cd <project-root> && agentic-agent status` — Verify orientation
7. Read task list and announcements (filter by `project_id`)

### TechLead Commands

```bash
agentic-agent task list              # View backlog
agentic-agent task claim <ID>        # Reserve task
agentic-agent context build          # Get full context
agentic-agent specify gate-check <ID>    # Validate quality
agentic-agent task complete <ID>     # Mark done
```

### TechLead Principles (from SOUL.md)

- Know which project you're in before every CLI call
- Respect `$COORDINATION_DIR` env var (or auto-discover it)
- Always run `sdd gate-check` before claiming tasks
- Delegate all code to workers (traffic-cop mode)
- Trust coordination YAML files as shared brain

---

## ProductLead Agent

**Role:** Product requirement coordinator

Responsibilities:
- Define requirements and create specs
- Announce spec-ready to TechLead
- Listen for TechLead completion signals
- Manage product roadmap per project
- Switch between multiple projects
- Validate specs before handoff

### Session Startup (ProductLead)

1. Load `SOUL.md` — Review core truths & boundaries
2. Load `USER.md` — Understand shared human context
3. Resolve `$COORDINATION_DIR` env var (same process as TechLead)
4. Read `$COORDINATION_DIR/PROJECTS.md` — Know all projects
5. Check `$COORDINATION_DIR/active-project.yaml` — Current workspace
6. `cd <project-root> && agentic-agent status` — Verify orientation
7. Read decisions and announcements (filter by `project_id`)

### ProductLead Commands

```bash
agentic-agent openspec init <name>   # Create new spec
agentic-agent openspec check <id>    # Validate spec
product-wizard                       # Generate PRDs
sdd analyst                          # Discovery & analysis
sdd architect                        # Architecture design
```

### ProductLead Principles (from SOUL.md)

- Respect `$COORDINATION_DIR` env var (or auto-discover it)
- Align on active project before announcing spec-ready
- Start from user needs, not implementation details
- Create formal specs before any task enters TechLead's backlog
- Always run `openspec init` before handing off to TechLead
- Track all product decisions in `.agentic/context/decisions.md`

---

## Multi-Project Coordination

### How They Talk

Both agents communicate via `.agentic/coordination/announcements.yaml`:

```yaml
announcements:
  # ProductLead → TechLead
  - from_agent: product-lead
    to_agent: tech-lead
    project_id: proj-001          # Which project this belongs to
    status: spec-ready
    summary: "Feature approved"
    data:
      spec_path: .agentic/spec/auth/proposal.md
      task_count: 8

  # TechLead → ProductLead
  - from_agent: tech-lead
    to_agent: product-lead
    project_id: proj-001
    status: complete
    summary: "Feature shipped"
    data:
      task_count: 8
      commits: [abc123, def456]
```

### Project Registry

Stored in `$COORDINATION_DIR/PROJECTS.md`:

```markdown
| ID | Name | Root Path | Stack | Status |
|----|------|-----------|-------|--------|
| proj-001 | my-app | ~/projects/my-app | Go | active |
| proj-002 | my-web | ~/projects/web | React | active |
```

### Switching Projects

```bash
cd ~/projects/project-1
agentic-agent status      # Loads context for project-1

cd ~/projects/project-2
agentic-agent status      # Loads context for project-2
```

---

## Setup Guide

### Step 1: Configure Coordination Directory

```bash
mkdir -p ~/my-org/coordinators
export COORDINATION_DIR=~/my-org/coordinators
```

### Step 2: Create Project Registry

```bash
# Copy template
cp resources/PROJECTS.md $COORDINATION_DIR/PROJECTS.md

# Edit to add your projects
```

### Step 3: Load Both Agents

Load both of these files as agents in Claude Code / Cursor:

**TechLead:**
- Use `tech-lead-identity.md` for identity
- Load `tech-lead-soul.md`, `tech-lead-agents.md`, `tech-lead-tools.md`
- Context packs: Optional — use resources/ for detailed guides

**ProductLead:**
- Use `product-lead-identity.md` for identity
- Load `product-lead-soul.md`, `product-lead-agents.md`, `product-lead-tools.md`
- Context packs: Optional — use resources/ for detailed guides

### Step 4: Start Coordinating

**ProductLead creates specs:**
```bash
cd ~/projects/project-1
agentic-agent openspec init "Feature Name"
agentic-agent openspec check <spec-id>
```

Then announce in `.agentic/coordination/announcements.yaml` with `project_id`.

**TechLead routes work:**
```bash
cd ~/projects/project-1
agentic-agent task list
agentic-agent task claim <task-id>
agentic-agent specify gate-check <spec-id>
```

Then announce completion with `project_id`.

---

## File Organization

When installed, you'll have access to:

- `agents/tech-lead-*.md` — TechLead configuration
- `agents/product-lead-*.md` — ProductLead configuration
- `resources/MULTI-PROJECT-GUIDE.md` — Complete setup guide
- `resources/PROJECTS.md` — Project registry template
- `resources/OPENCLAW-README.md` — Agent overview

---

## Auto-Discovery

If `$COORDINATION_DIR` is not set, both agents will search:

1. `../PROJECTS.md` (parent directory)
2. `../../PROJECTS.md` (grandparent directory)
3. Ask human to set env var

This allows common layouts to work without explicit setup.

---

## Key Concepts

**project_id filtering:** Announcements include `project_id` so agents only respond to their active project's signals. This enables multiple projects to run in parallel using the same coordination files.

**Environment variable configuration:** `$COORDINATION_DIR` tells both agents where to find the project registry and coordination files. Users choose their own directory structure.

**Session startup checklist:** Each agent follows a standard startup sequence (load identity → resolve env var → load project context → verify orientation) before taking any action.

**Traffic-cop mode:** TechLead stays in a lean mode (~5% context window) by delegating all implementation to workers, keeping overhead low.

---

## Next Steps

1. Read `resources/MULTI-PROJECT-GUIDE.md` for detailed setup
2. Copy `resources/PROJECTS.md` to your coordination directory
3. Load both agents in Claude Code / Cursor
4. Start creating specs (ProductLead) and routing work (TechLead)

See `resources/` for complete reference documentation.
