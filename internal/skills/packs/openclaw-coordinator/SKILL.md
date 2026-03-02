# OpenClaw Coordinator Agents

Activate TechLead and ProductLead agents to manage project work across multiple teams.

## Does exactly this

Installs six Claude Code agents configured for autonomous project coordination:

**Coordinator Agents:**
- **TechLead** — Routes technical tasks, spawns developers, runs quality gates
- **ProductLead** — Defines specs, creates requirements, coordinates with TechLead

**Developer Agents (spawned by TechLead):**
- **BackendDev** — Implements APIs, databases, services, business logic
- **FrontendDev** — Builds web UI components, visual design, accessibility
- **MobileDev** — Implements Flutter apps, cross-platform mobile
- **QADev** — Tests implementation, enforces 8/10 quality gate, scores all work

All agents:
- Manage work across multiple projects via `$COORDINATION_DIR` environment variable
- Communicate through shared `announcements.yaml` file
- Support project switching with automatic context loading
- Filter announcements by `project_id` for multi-project awareness
- Coordinate APIs via contract-driven development (ProductLead → TechLead → BackendDev → FrontendDev/MobileDev → QADev)

## When to use this

- You're running multiple projects in parallel
- You need ProductLead creating specs and TechLead routing technical work
- You want autonomous agents managing project coordination
- You need clear separation between product definition and technical execution

## Installation

```bash
agentic-agent skills install openclaw-coordinator --tool claude-code
```

This installs six agents to `.claude/agents/`:
- **TechLead** — `openclaw-tech-lead.md` (coordinator)
- **ProductLead** — `openclaw-product-lead.md` (coordinator)
- **BackendDev** — `openclaw-backend-dev.md` (worker)
- **FrontendDev** — `openclaw-frontend-dev.md` (worker)
- **MobileDev** — `openclaw-mobile-dev.md` (worker)
- **QADev** — `openclaw-qa-dev.md` (quality gatekeeper)

## Setup

1. Set coordination directory:
   ```bash
   mkdir -p ~/my-org/coordinators
   export COORDINATION_DIR=~/my-org/coordinators
   ```

2. Create project registry:
   ```bash
   cp resources/PROJECTS.md $COORDINATION_DIR/PROJECTS.md
   ```

3. Load both agents in Claude Code / Cursor

4. Start coordinating:
   - ProductLead: `agentic-agent openspec init "Feature Name"`
   - TechLead: `agentic-agent task claim <ID>`

## Key Files

After installation, you'll have six agents in `.claude/agents/`:
- `openclaw-tech-lead.md`, `openclaw-product-lead.md` — Coordinators
- `openclaw-backend-dev.md`, `openclaw-frontend-dev.md`, `openclaw-mobile-dev.md`, `openclaw-qa-dev.md` — Workers
- `resources/` — Configuration templates and guides

## Configuration

Both agents use `$COORDINATION_DIR` environment variable to find:
- `PROJECTS.md` — Registry of all known projects
- `active-project.yaml` — Current active workspace (optional)
- `.agentic/coordination/announcements.yaml` — Shared communication

If `$COORDINATION_DIR` is not set, agents auto-discover by searching:
- `../PROJECTS.md` (parent directory)
- `../../PROJECTS.md` (grandparent directory)

## If you need more detail

See `resources/MULTI-PROJECT-GUIDE.md` for:
- Complete setup instructions (3 configuration options)
- Session startup checklist for both agents
- Multi-project coordination patterns
- Troubleshooting and scaling guidance
- Real example timeline showing parallel projects

See `resources/` for:
- PROJECTS.md template (copy to $COORDINATION_DIR)
- Agent configuration references (SOUL.md, AGENTS.md, TOOLS.md)
- Coordination file format (announcements.yaml schema)
