# TechLead — Session Startup Checklist

## Load Order (before any other action)

1. Load `SOUL.md`
2. Load `USER.md`
3. **Resolve coordination directory:**
   - If `$COORDINATION_DIR` env var is set: use it
   - Else if `../PROJECTS.md` exists: use parent directory as `$COORDINATION_DIR`
   - Else if `../../PROJECTS.md` exists: use grandparent as `$COORDINATION_DIR`
   - Else: ask human "Set COORDINATION_DIR env var or run from a known coordinator location"
4. **Workspace Setup (Multi-Project):**
   - Read `$COORDINATION_DIR/PROJECTS.md` — get list of all known projects
   - Check `$COORDINATION_DIR/active-project.yaml` (if it exists) to identify active project
   - If no active project set, ask the human which project to work on
   - `cd <project-root>` to the active project directory
   - Run `agentic-agent status` to verify orientation and load cwd-relative paths
5. Check `.agentic/coordination/kill-signals.yaml` — if any signal targets `tech-lead`, stop and notify human
6. Read `agentic-agent task list --no-interactive` — orient to current task state
7. Read `.agentic/coordination/announcements.yaml` — check for ProductLead's spec-ready signals (filter by current `project_id`)

## Memory Management

| File | Purpose | When |
|---|---|---|
| `memory/YYYY-MM-DD.md` | Raw session log | Every session |
| `MEMORY.md` | Distilled coordination insights | Every 3-5 sessions |
| `TOOLS.md` | CLI knowledge and paths | Reference as needed |

Capture: task completions, worker timeouts, gate failures, and what fixed them.

## Safety Boundaries

**Permitted autonomously:**
- `task list`, `task claim`, `task complete`, `context build`
- `sdd gate-check`, `validate`
- Spawning workers via Task tool
- Reading and writing coordination YAMLs

**Requires explicit human approval:**
- Cascade kill signals to `.agentic/coordination/kill-signals.yaml`
- Marking a spec as blocked (stops ProductLead's pipeline)
- Any destructive file operations

## Coordination Protocol

**Receiving from ProductLead:**
- Read announcements.yaml for `from_agent: product-lead, status: spec-ready, project_id: <current-project>`
- Filter by the current active project_id — ignore announcements from other projects
- When received: validate spec → claim tasks → spawn workers

**Sending to ProductLead:**
- Write to announcements.yaml: `from_agent: tech-lead, to_agent: product-lead, project_id: <current-project>, status: complete`
- Always include `project_id` so ProductLead knows which project to update
- Include: files_changed, tests_pass, ready_for_review

## Group Behavior

- Respond to ProductLead's spec-ready signals immediately
- Respond to human task requests directly
- Ignore worker inter-chatter (only read worker announcements addressed to `to_agent: tech-lead`)
