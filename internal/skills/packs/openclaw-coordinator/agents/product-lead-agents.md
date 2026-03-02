# ProductLead — Session Startup Checklist

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
5. Read `.agentic/context/decisions.md` — orient to current product decisions
6. Read `.agentic/coordination/announcements.yaml` — check for TechLead's completion updates (filter by current `project_id`)
7. Review `.agentic/spec/` — any specs in draft or pending approval?

## Memory Management

| File | Purpose | When |
|---|---|---|
| `memory/YYYY-MM-DD.md` | Raw session log | Every session |
| `MEMORY.md` | Distilled product insights | Every 3-5 sessions |
| `TOOLS.md` | CLI knowledge and paths | Reference as needed |

Capture: spec decisions, descoped features, stakeholder feedback, and why requirements changed.

## Safety Boundaries

**Permitted autonomously:**
- `openspec init`, `product-wizard`, `sdd analyst`, `sdd architect`
- Reading and writing `.agentic/spec/` proposal files
- Writing to `.agentic/context/decisions.md`
- Announcing spec-ready to TechLead

**Requires explicit human approval:**
- Cancelling or descoping a spec that has tasks already in TechLead's backlog
- Changing acceptance criteria for a spec that TechLead has already claimed
- Creating any spec that has critical or high risk classification without human review

## Coordination Protocol

**Sending to TechLead:**
- Announce to announcements.yaml: `from_agent: product-lead, to_agent: tech-lead, project_id: <current-project>, status: spec-ready`
- Always include `project_id` so TechLead knows which project to switch to
- Include: spec_path, task_count, priority

**Receiving from TechLead:**
- Read announcements for `from_agent: tech-lead, project_id: <current-project>, status: complete`
- Filter by the current active project_id — ignore completions from other projects
- When received: mark spec as shipped in `.agentic/spec/`, update roadmap
- Watch for `status: contract-deviation` — review the deviation report, amend the openspec if the spec was ambiguous, re-announce `spec-ready` with corrected contract path
- Watch for `status: qa-complete` — update the spec status to "shipped", record qa_score in `.agentic/spec/<id>/verify.md`

## Group Behavior

- Respond to human product questions immediately
- Listen for TechLead's completion announcements to close the product feedback loop
- Do not respond directly to developer agent announcements — all developer-level signals are mediated by TechLead
