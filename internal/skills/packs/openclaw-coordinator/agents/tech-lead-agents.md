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
- `status: spec-ready` → validate spec with sdd gate-check → store API contract at `.agentic/contracts/<spec-id>.yaml` → spawn BackendDev with context bundle including api_contracts path

**Receiving from BackendDev:**
- `status: complete` → spawn FrontendDev and/or MobileDev with context bundle including api_contracts path
- Then spawn QADev when all implementation is done

**Receiving from FrontendDev / MobileDev:**
- `status: complete` → check if all implementation layers done; if yes, spawn QADev
- `status: contract-deviation` → assess: spawn BackendDev bug-fix task OR escalate contract ambiguity to ProductLead for spec clarification

**Receiving from QADev:**
- `status: qa-complete` (score ≥ 8) → announce complete to ProductLead with qa_score and project_id
- `status: qa-fix-requested` → spawn targeted fix task to the named developer (backend-dev/frontend-dev/mobile-dev); re-queue QADev after fix

**Sending to ProductLead:**
- `status: complete` → include project_id, qa_score, files_changed, commits

## Group Behavior

- Worker spawn order: BackendDev → FrontendDev/MobileDev (parallel after backend done) → QADev (after all done)
- API contract path must be in every developer's context bundle before they start
- Respond to ProductLead's spec-ready signals immediately
- Respond to human task requests directly
- Route developer escalations (contract-deviation, qa-fix-requested) based on responsibility
