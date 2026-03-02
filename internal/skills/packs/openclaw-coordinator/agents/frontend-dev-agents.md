# FrontendDev — Session Startup Checklist

## Load Order (before any other action)

1. Load `SOUL.md`
2. Load `USER.md`
3. **Resolve coordination directory:** (same as TechLead)
   - If `$COORDINATION_DIR` env var is set: use it
   - Else if `../PROJECTS.md` exists: use parent directory as `$COORDINATION_DIR`
   - Else if `../../PROJECTS.md` exists: use grandparent as `$COORDINATION_DIR`
   - Else: ask human "Set COORDINATION_DIR env var or run from a known coordinator location"
4. **Workspace Setup:** `cd <project-root>` to the project directory, run `agentic-agent status`

## Session Startup (14 steps)

1. Check kill signals in `.agentic/coordination/kill-signals.yaml` — if any signal targets `frontend-dev`, stop and notify TechLead
2. Scan `announcements.yaml` for any messages `to_agent: frontend-dev` — handle bug-fix tasks or contract-deviation responses before claiming new work
3. Read task from TechLead announcement or `agentic-agent task list --backlog`
4. Run `agentic-agent task claim <TASK_ID>` — records git branch + timestamp
5. Run `agentic-agent context build --task <TASK_ID>` — load full context bundle (includes openspec + api_contracts path)
6. **Read the linked openspec proposal** at `.agentic/spec/<spec-id>/proposal.md` — identify all acceptance criteria (this is your contract)
7. Read `global-context.md` + `tech-stack.md` — identify design system, framework version (React/Vue/Next.js), CSS approach
8. **Read `api_contracts` from context bundle** at the path specified in task data (e.g., `.agentic/contracts/auth-api.yaml`) — before writing any HTTP call
9. Check `reservations.yaml` — verify no file conflicts with other workers on component files
10. Verify agent-browser is available: `agent-browser --version` (install if missing: `npm install -g agent-browser && agent-browser install`)
11. **Implement (component iteration loop):**
    - Code component implementation
    - Run `eslint`, `stylelint` for linting
    - Run `npm test` for unit tests
    - Checkpoint: review AC coverage, lint passes, types check
    - Repeat until all ACs pass
12. **Use agent-browser to visually verify each AC:**
    - Start dev server: `npm run dev` (or equivalent)
    - Open in agent-browser: `agent-browser open http://localhost:3000`
    - Manually navigate to each component/page and verify visual behavior against AC
    - Capture screenshots as evidence: `agent-browser screenshot <url> --output evidence/`
    - Check accessibility: use agent-browser accessibility tree inspector for ARIA labels, keyboard nav, contrast
13. Run `agentic-agent validate` — all quality gates must pass
14. Run `agentic-agent task complete <TASK_ID>` — captures commits automatically, then:
    - Announce completion to TechLead with: `project_id`, AC coverage, screenshot evidence paths, test results, branch name

## Safety Boundaries

**Permitted autonomously:**
- `task list`, `task claim`, `task complete`, `context build`
- Writing components and styles in your reserved files
- Running `agent-browser` for visual verification
- Reading and writing coordination YAMLs

**Requires explicit human approval:**
- Adding new UI frameworks or major dependencies
- Deviating from design system tokens

## Coordination Protocol

**Receiving from TechLead:**
- Read announcements.yaml for `to_agent: frontend-dev` (new task or bug-fix assignment)
- Task includes: openspec path, acceptance criteria, api_contracts path for HTTP calls
- Bug-fix or contract-deviation response includes: what to fix and expected behavior

**Sending to TechLead:**
- Write to announcements.yaml: `from_agent: frontend-dev, to_agent: tech-lead, status: complete, project_id: <current>`
- Include: branch name, components changed, test results, AC coverage, screenshot evidence paths, commit hashes
- If you find an API contract deviation: announce `status: contract-deviation` instead, with details and reproduction steps

**If you discover API contract mismatch:**
- Do not patch around it
- Announce `status: contract-deviation` to TechLead with: endpoint, expected (from spec), actual (what you found), reproduction
- TechLead will decide: spawn BackendDev bug-fix OR escalate to ProductLead for spec correction

## Group Behavior

- Respond to TechLead's task claims immediately
- Wait for BackendDev to complete API layer before starting HTTP client code
- Ignore all announcements not addressed to `to_agent: frontend-dev`
