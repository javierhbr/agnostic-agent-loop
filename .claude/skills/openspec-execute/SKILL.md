---
name: openspec-execute
description: Execute implementation tasks from an approved openspec change. Use after openspec planning is complete and user has approved the plan. Triggers include "execute openspec", "implement change", "work on change", "continue implementing".
---

## MANDATORY: Use CLI commands only

**You MUST use `agentic-agent` CLI commands for ALL openspec and task operations.**
Do NOT manually create, edit, or modify any files under `.agentic/` directly.
Do NOT skip CLI commands to "save time" or "be more efficient".
The CLI maintains state consistency — bypassing it causes data corruption.

## Use this skill when

- User approves an openspec plan and wants to start implementation
- User says "execute openspec", "implement change", "work on <change-id>"
- User wants to continue working on an in-progress openspec change
- The planning phase (openspec skill) is complete and user said "proceed"

## Do not use this skill when

- Planning/proposal phase is not complete (use openspec skill instead)
- No approved change exists to work from
- User is still discussing requirements or scope
- The task is a quick fix unrelated to an openspec change

## Example prompts

- "Execute openspec auth-feature"
- "Implement the tasks for payment-system"
- "Continue working on change notifications"
- "Start implementing" (after openspec planning completes)

## Instructions

### Phase 0: Identify the change to work on

1. **Check for explicit change ID** in the user's prompt
   - Look for patterns: "execute <id>", "implement <id>", "work on <id>", "<id>"
   - If found, verify it exists and proceed to Phase 1

2. **If no ID provided**, list eligible changes:
   ```bash
   agentic-agent openspec list
   ```
   Look for changes with status "imported" or "implementing" (not "draft" or "archived").

3. **Handle results**:
   - If exactly one eligible change: use it automatically, confirm with user
   - If multiple eligible changes: ask user which one to work on
     ```
     Found multiple changes ready for execution:
     - auth-feature (imported, 5 tasks)
     - payment-system (implementing, 3 tasks remaining)

     Which change would you like to work on?
     ```
   - If none found: inform user and suggest next steps
     ```
     No changes are ready for execution.
     Create a change first: use the openspec skill with a requirements file.
     ```

4. **Validate change is ready**: Run `agentic-agent openspec show <change-id>`
   - If status is "draft": the plan is not complete. Tell the user to run the openspec skill first.
   - If status is "imported" or "implementing": proceed to Phase 1.
   - If status is "implemented" or "archived": inform user this change is already done.

### Phase 1: Execute tasks sequentially

Tasks are **auto-imported** into the backlog when you run `task list` or `task claim`.
No manual import step is needed.

When tasks have detail files, the imported tasks include:
- Description from the task file
- Acceptance criteria (viewable via `agentic-agent task show <id>`)
- Prerequisites mapped as task inputs
- Technical notes in the description

For each task (check progress with `agentic-agent openspec status <change-id>`):

1. **Check for in-progress tasks first**: `agentic-agent task continue` (resumes an already-claimed task, or use `task continue <task-id>` for a specific one)
   - If a task is already in-progress, continue from where it was left off
   - If the task is still pending, `task continue` auto-claims it
2. **If no in-progress task**, list and claim the next one:
   - `agentic-agent task list` (triggers auto-import if needed)
   - `agentic-agent task claim <task-id>`
3. **Review details**: `agentic-agent task show <task-id>` (see description, acceptance criteria)
4. **Read task detail file** if it exists (in `.agentic/openspec/changes/<change-id>/tasks/`)
5. **Read relevant specs** from `.agentic/openspec/changes/<change-id>/specs/` if they exist
6. **Read tech stack** from `.agentic/context/tech-stack.md` for technology guidance
7. Implement the work
8. Run tests — verify the task works
9. **Complete via CLI**: `agentic-agent task complete <task-id>`

**Never skip ahead.** Complete and verify each task before starting the next.
**Never modify `.agentic/tasks/*.yaml` files directly.**

### Phase 2: Complete and archive

When all tasks are done:

```bash
agentic-agent openspec complete <change-id>
agentic-agent openspec archive <change-id>
```

Report the summary to the user.

## Rules (non-negotiable)

1. Always use `agentic-agent` CLI commands for task and change operations
2. Never write directly to `.agentic/` YAML files
3. Always claim a task before working on it
4. Always complete a task after finishing it
5. Execute tasks in order — do not skip ahead to later tasks
6. Never modify task status manually — use `task claim` and `task complete`
7. If a task fails or is blocked, report to user and ask how to proceed
8. Read task detail files (in `tasks/` subdirectory) for implementation guidance
9. Reference specs in `specs/` directory for technical details
10. Each task must pass its acceptance criteria before being marked complete
