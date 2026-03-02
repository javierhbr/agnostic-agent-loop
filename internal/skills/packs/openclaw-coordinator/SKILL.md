# OpenClaw Coordinator Agents

Activate TechLead and ProductLead agents to manage project work across multiple teams.

## Does exactly this

Installs two Claude Code agents configured for autonomous project coordination:

- **TechLead** — Routes technical tasks, spawns builders, runs quality gates
- **ProductLead** — Defines specs, creates requirements, coordinates with TechLead

Both agents:
- Manage work across multiple projects via `$COORDINATION_DIR` environment variable
- Communicate through shared `announcements.yaml` file
- Support project switching with automatic context loading
- Filter announcements by `project_id` for multi-project awareness

## When to use this

- You're running multiple projects in parallel
- You need ProductLead creating specs and TechLead routing technical work
- You want autonomous agents managing project coordination
- You need clear separation between product definition and technical execution

## Installation

```bash
agentic-agent skills install openclaw-coordinator --tool claude-code
```

This installs:
- TechLead agent to `.claude/agents/openclaw-tech-lead.md`
- ProductLead agent to `.claude/agents/openclaw-product-lead.md`

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

After installation, you'll have:
- `.claude/agents/openclaw-tech-lead.md` — TechLead agent
- `.claude/agents/openclaw-product-lead.md` — ProductLead agent
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
