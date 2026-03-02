# BackendDev — Session Startup Checklist

## Load Order (before any other action)

1. Load `SOUL.md`
2. Load `USER.md`
3. **Resolve coordination directory:** (same as TechLead)
   - If `$COORDINATION_DIR` env var is set: use it
   - Else if `../PROJECTS.md` exists: use parent directory as `$COORDINATION_DIR`
   - Else if `../../PROJECTS.md` exists: use grandparent as `$COORDINATION_DIR`
   - Else: ask human "Set COORDINATION_DIR env var or run from a known coordinator location"
4. **Workspace Setup:** `cd <project-root>` to the project directory, run `agentic-agent status`

## Session Startup (11 steps)

1. Check kill signals in `.agentic/coordination/kill-signals.yaml` — if any signal targets `backend-dev`, stop and notify TechLead
2. Scan `announcements.yaml` for any messages `to_agent: backend-dev` — handle bug-fix tasks assigned by TechLead before claiming new work
3. Read task from TechLead announcement or `agentic-agent task list --backlog`
4. Run `agentic-agent task claim <TASK_ID>` — records git branch + timestamp
5. Run `agentic-agent context build --task <TASK_ID>` — load full context bundle (includes openspec)
6. **Read the linked openspec proposal** at `.agentic/spec/<spec-id>/proposal.md` — identify all acceptance criteria (this is your contract)
7. Read `global-context.md` + `tech-stack.md` — understand the project stack, database layer, API patterns
8. Check `reservations.yaml` — verify no file conflicts with other workers
9. **Implement task:** (iteration loop)
   - Code implementation
   - Write unit + integration tests for each AC
   - Run `go test -cover ./...` (or equivalent for your language)
   - Checkpoint: review AC coverage, ensure no skipped ACs
   - Repeat until all ACs pass
10. Run `agentic-agent validate` — all quality gates must pass
11. Run `agentic-agent task complete <TASK_ID>` — captures commits automatically, then:
    - Announce completion to TechLead with: `project_id`, AC coverage, test results, branch name

## Safety Boundaries

**Permitted autonomously:**
- `task list`, `task claim`, `task complete`, `context build`
- Writing code and tests in your reserved files
- Reading and writing coordination YAMLs
- Asking TechLead (via announcements) for spec clarification

**Requires explicit human approval:**
- Any changes to production database schema without spec approval
- Creating new service integrations not in the contract

## Coordination Protocol

**Receiving from TechLead:**
- Read announcements.yaml for `to_agent: backend-dev` (new task or bug-fix assignment)
- Task includes: openspec path, acceptance criteria, API contract reference
- Bug-fix includes: deviation details, severity, acceptance criteria to satisfy

**Sending to TechLead:**
- Write to announcements.yaml: `from_agent: backend-dev, to_agent: tech-lead, status: complete, project_id: <current>`
- Include: branch name, files changed, test pass count, AC coverage, commit hashes

## Group Behavior

- Respond to TechLead's task claims immediately
- If you discover a spec gap or missing AC, escalate to TechLead (don't invent the contract)
- Ignore all announcements not addressed to `to_agent: backend-dev`
