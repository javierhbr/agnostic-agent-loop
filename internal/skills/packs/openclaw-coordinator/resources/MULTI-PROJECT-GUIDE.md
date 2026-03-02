# Multi-Project Workspace Guide

TechLead and ProductLead can now manage multiple projects in parallel. Configuration is flexible via `$COORDINATION_DIR` env var.

---

## Quick Start: 3 Projects

### 1. Choose Where to Store Coordinator Files

Decide where you want to keep `PROJECTS.md`, `active-project.yaml`, `USER.md`, etc:

```bash
# Option A: Separate coordinator directory (recommended)
mkdir -p ~/my-org/coordinators
export COORDINATION_DIR=~/my-org/coordinators

# Option B: Custom location
export COORDINATION_DIR=~/.app-factory/metadata
export COORDINATION_DIR=/opt/projects/shared

# Option C: If you don't set env var, agents auto-discover:
# They search: ../PROJECTS.md → ../../PROJECTS.md → ask you to set env var
```

### 2. Set Up Project Registry

Create or edit `$COORDINATION_DIR/PROJECTS.md` and add your projects:

```markdown
| ID | Name | Root Path | Stack | Status |
|----|------|-----------|-------|--------|
| proj-trailmate | Trailmate | ~/projects/trailmate | Swift/Go | active |
| proj-fithub | FitHub | ~/projects/fithub | Swift/Node | active |
| proj-web | Fitness Web | ~/projects/web | React/Python | active |
```

### 3. Session Startup (Both Agents)

Both agents follow this pattern at session start:

```
1. Load SOUL.md & USER.md (identity + context)
2. Resolve $COORDINATION_DIR:
   - If env var set: use it
   - Else search: ../PROJECTS.md → ../../PROJECTS.md
   - Else: ask human to set env var
3. Read $COORDINATION_DIR/PROJECTS.md (know all projects)
4. Check $COORDINATION_DIR/active-project.yaml (know current workspace)
5. cd <project-root> && agentic-agent status (orient to that project)
6. Start work (ProductLead creates specs, TechLead builds them)
```

### 4. Key Pattern: project_id in Announcements

Every announcement includes `project_id` so agents know which workspace to cd to:

```yaml
# ProductLead announces in Trailmate
- from_agent: product-lead
  to_agent: tech-lead
  project_id: proj-trailmate    # ← This tells TechLead which cd to use
  status: spec-ready
  summary: "GPS Recording approved"

# TechLead reads it, sees proj-trailmate, does:
cd ~/projects/trailmate
agentic-agent task claim TASK-ID
```

### 5. Filtering by project_id

When agents read announcements, they filter:

**TechLead:**
```
Read .agentic/coordination/announcements.yaml
Filter: from_agent == "product-lead" AND project_id == "proj-trailmate"
Ignore: entries with other project_ids
```

**ProductLead:**
```
Read .agentic/coordination/announcements.yaml
Filter: from_agent == "tech-lead" AND project_id == "proj-fithub"
Ignore: entries with other project_ids
```

---

## Workspace Switching

### Manual Switch (Any Moment)

```bash
# First time: set env var
export COORDINATION_DIR=~/my-org/coordinators

# Current: working in Trailmate
agentic-agent status  # shows Trailmate context

# Switch to FitHub
cd ~/projects/fithub
agentic-agent status  # shows FitHub context

# Update active project (optional)
cat > $COORDINATION_DIR/active-project.yaml << 'EOF'
project_id: proj-fithub
root: ~/projects/fithub
name: FitHub
switched_at: 2026-03-01T11:00:00Z
EOF
```

### Automated Switch (via Announcement)

When TechLead sees a spec-ready announcement with `project_id: proj-web`:

```bash
# Agent logic (assuming COORDINATION_DIR is set or auto-discovered):
export COORDINATION_DIR=~/my-org/coordinators  # or already set

# Look up proj-web root in $COORDINATION_DIR/PROJECTS.md
cd ~/projects/web          # ← determined by project_id lookup
agentic-agent status       # verify orientation
agentic-agent task list    # see tasks for proj-web
agentic-agent task claim TASK-ID

# Update active-project.yaml
cat > $COORDINATION_DIR/active-project.yaml << 'EOF'
project_id: proj-web
root: ~/projects/web
name: Fitness Web
switched_at: [current timestamp]
EOF
```

---

## Core Principles

### ProductLead

1. **Read PROJECTS.md** — Know which projects exist
2. **Verify active project** — Where should I create this spec?
3. **Include project_id** — Every announcement must say which project
4. **Announce with confidence** — TechLead will switch to your project

### TechLead

1. **Read PROJECTS.md** — Know which projects exist
2. **Filter announcements** — Only respond to my active project_id
3. **Switch projects explicitly** — `cd <root> && agentic-agent status`
4. **Announce completion** — Include project_id so ProductLead knows where you shipped

---

## Files

- **PROJECTS.md** — Human-editable registry of all projects
- **active-project.yaml** — Optional tracking of current workspace
- **USER.md** — Shared user context (both agents read)
- **tech-lead/TOOLS.md** — Workspace management commands
- **tech-lead/AGENTS.md** — Workspace setup in load order
- **tech-lead/SOUL.md** — "Know which project you're in"
- **product-lead/TOOLS.md** — Workspace management + announcement format
- **product-lead/AGENTS.md** — Workspace setup in load order
- **product-lead/SOUL.md** — "Align on active project before announcing"

---

## Example: Two Specs in Parallel

```
Timeline:

T=0:  ProductLead reads PROJECTS.md
      → knows: Trailmate, FitHub, Web

T=1:  ProductLead cd ~/projects/trailmate
      → creates GPS spec
      → announces: project_id=proj-trailmate, status=spec-ready

T=2:  TechLead reads announcement
      → sees project_id=proj-trailmate
      → cd ~/projects/trailmate
      → claims TASK-500-1
      → spawns 2 builders

T=3:  ProductLead cd ~/projects/fithub
      → creates Social spec
      → announces: project_id=proj-fithub, status=spec-ready

T=4:  TechLead finishes Trailmate builders
      → reads announcements
      → filters by project_id=proj-trailmate (own active project)
      → announces: status=complete

T=5:  ProductLead reads completion announcement
      → sees project_id=proj-trailmate
      → knows Trailmate shipped
      → closes Trailmate spec
      → but still monitoring proj-fithub for TechLead...

T=6:  TechLead finishes current Trailmate tasks
      → reads announcements
      → sees project_id=proj-fithub (new spec from ProductLead)
      → cd ~/projects/fithub
      → claims TASK-600-1
      → spawns builders for Social spec
      → Both projects now being built in parallel!
```

---

## Troubleshooting

| Problem | Cause | Solution |
|---------|-------|----------|
| TechLead reads wrong project's tasks | `active-project.yaml` stale | Run `agentic-agent status` to verify cwd |
| ProductLead forgets project_id in announcement | Human error | Always include `project_id: proj-XXX` in YAML |
| Workers in wrong project directory | Agent didn't cd before claiming | Verify: `agentic-agent status` shows correct project |
| Announcements mixed up between projects | Single file, multiple projects | Filter by `project_id` when reading |

---

## Scaling

### 1 Project
- Just ignore project_id filtering
- Single-project mode works fine

### 2-3 Projects
- Use this guide
- Both agents manage all projects
- Announcements keep them in sync

### 4+ Projects
- Same pattern, just more entries in PROJECTS.md
- Consider adding priority/urgency to announcements
- Add rate-limiting if too many specs created at once

---

## For Developers

### Adding a New Project

1. Add row to PROJECTS.md:
   ```
   | proj-newapp | NewApp | ~/projects/newapp | Stack | active |
   ```

2. Create project directory:
   ```bash
   mkdir ~/projects/newapp
   cd ~/projects/newapp
   agentic-agent init --name "NewApp"
   ```

3. Both agents auto-discover on next session startup

### Archiving a Project

1. Update PROJECTS.md:
   ```
   | proj-oldapp | OldApp | ~/projects/oldapp | Stack | archived |
   ```

2. Both agents will skip `archived` projects (you can add that filter)

### Monitoring All Projects

```bash
# From .openclaw/
for id in proj-trailmate proj-fithub proj-web; do
  echo "=== $id ==="
  grep "project_id: $id" .agentic/coordination/announcements.yaml | tail -3
done
```

---

**See also:** [OpenClaw App Factory Example](../examples/openclaw-app-factory/README.md) for full multi-project orchestration patterns.
