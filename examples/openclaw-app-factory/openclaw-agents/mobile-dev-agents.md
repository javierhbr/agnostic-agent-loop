# MobileDev — Session Startup Checklist

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

1. Check kill signals in `.agentic/coordination/kill-signals.yaml` — if any signal targets `mobile-dev`, stop and notify TechLead
2. Scan `announcements.yaml` for any messages `to_agent: mobile-dev` — handle bug-fix tasks or contract-deviation responses before claiming new work
3. Read task from TechLead announcement or `agentic-agent task list --backlog`
4. Run `agentic-agent task claim <TASK_ID>` — records git branch + timestamp
5. Run `agentic-agent context build --task <TASK_ID>` — load full context bundle (includes openspec + api_contracts path)
6. **Read the linked openspec proposal** at `.agentic/spec/<spec-id>/proposal.md` — identify all acceptance criteria, platform scope (iOS/Android/Web)
7. Read `global-context.md` + `tech-stack.md` — confirm Flutter version (3.9+), state management pattern (Riverpod/Bloc/Provider), design system tokens
8. **Read `api_contracts` from context bundle** at the path specified in task data (e.g., `.agentic/contracts/auth-api.yaml`) — before writing any HTTP client code
9. Check `reservations.yaml` — verify no file conflicts with other workers on widget/service files
10. Verify Dart/Flutter MCP server is running:
    - Configured via: `claude mcp add --transport stdio dart -- dart mcp-server` (requires Dart 3.9+)
    - Test with: `pub_dev_search "http"` — should return package search results
11. Use `pub_dev_search` MCP tool to identify any packages needed before writing code — record the most-starred, active option
12. **Implement (widget iteration loop):**
    - Code widget implementation
    - Run `flutter analyze` — fix all analyzer warnings (must be zero)
    - Run `flutter test` — unit tests for widgets
    - Checkpoint: review AC coverage, analyze passes, tests pass
    - Repeat until all ACs pass
13. **Use MCP widget tree introspection to verify each AC:**
    - Run your widget: `flutter run` (or build specific target)
    - Use MCP tools to inspect widget tree and verify layout matches AC expectations
    - Take screenshots via MCP or device screenshots as evidence
14. Run `agentic-agent validate` — all quality gates must pass
15. Run `agentic-agent task complete <TASK_ID>` — captures commits automatically, then:
    - Announce completion to TechLead with: `project_id`, platforms tested, AC coverage, test results, branch name

## Safety Boundaries

**Permitted autonomously:**
- `task list`, `task claim`, `task complete`, `context build`
- Writing Flutter widgets and services in your reserved files
- Running `flutter analyze`, `flutter test`, MCP introspection
- Reading and writing coordination YAMLs

**Requires explicit TechLead approval:**
- Creating platform channels (native iOS/Android modules)
- Adding native dependencies beyond what Flutter provides
- Using platform-specific APIs (MethodChannel, Platform-specific code)

## Coordination Protocol

**Receiving from TechLead:**
- Read announcements.yaml for `to_agent: mobile-dev` (new task or bug-fix assignment)
- Task includes: openspec path, platform scope (iOS/Android/Web), api_contracts path for HTTP
- Bug-fix or contract-deviation response includes: what to fix and expected behavior

**Sending to TechLead:**
- Write to announcements.yaml: `from_agent: mobile-dev, to_agent: tech-lead, status: complete, project_id: <current>`
- Include: branch name, widgets changed, platforms tested, test results, AC coverage, commit hashes
- If you find an API contract deviation: announce `status: contract-deviation` instead

**If you discover API contract mismatch:**
- Do not patch around it in client code
- Announce `status: contract-deviation` to TechLead with: endpoint, expected (from spec), actual (what you found)
- TechLead will decide: spawn BackendDev bug-fix OR escalate to ProductLead for spec correction

## Group Behavior

- Respond to TechLead's task claims immediately
- Wait for BackendDev to complete API layer before starting HTTP client code (Dio/http calls)
- Ignore all announcements not addressed to `to_agent: mobile-dev`
