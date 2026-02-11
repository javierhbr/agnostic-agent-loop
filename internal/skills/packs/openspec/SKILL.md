---
name: openspec
description: Spec-driven development from requirements files. Scaffolds changes, imports tasks, and tracks implementation. Use when starting from requirements, implementing specs, or managing change proposals. Triggers include "openspec", "implement from spec", "start from requirements", "apply requirements".
---

## Use this skill when

- Starting a new feature from a requirements file
- User says "openspec", "implement from spec", "start from requirements"
- User provides a requirements/spec file and wants structured implementation

## Do not use this skill when

- The task is a quick fix or single-file change
- There are no requirements to work from

## Example prompts

- "Start a project from requirements.md following openspec"
- "Implement the features described in docs/auth-spec.md using openspec"
- "openspec init from .agentic/spec/payment-requirements.md"
- "Continue implementing change auth-feature" (resumes Phase 4)

## Instructions

### Phase 1: Initialize from requirements

1. Read the requirements file the user provided
2. Derive a short name for the change (e.g., "auth-feature", "payment-system")
3. Run:
   ```bash
   agentic-agent openspec init "<change-name>" --from <requirements-file>
   ```
4. Read the generated proposal at the path printed by the command

### Phase 2: Fill in the proposal and tasks

1. Edit `proposal.md` — fill in the Problem, Approach, Scope, and Acceptance Criteria sections based on the requirements
2. Edit `tasks.md` — write a numbered list of implementation tasks:
   ```
   1. First concrete task
   2. Second concrete task
   3. Third concrete task
   ```
   Keep tasks small, testable, and ordered by dependency.

### Phase 3: Import tasks into the backlog

Run:
```bash
agentic-agent openspec import <change-id>
```

This creates tasks in the backlog. Use TodoWrite to track them.

### Phase 4: Execute tasks sequentially

For each task (check progress with `agentic-agent openspec status <change-id>`):

1. `agentic-agent task claim <task-id>`
2. Implement the work
3. Run tests — verify the task works
4. `agentic-agent task complete <task-id>`

**Never skip ahead.** Complete and verify each task before starting the next.

### Phase 5: Complete and archive

When all tasks are done:

```bash
agentic-agent openspec complete <change-id>
agentic-agent openspec archive <change-id>
```

Report the summary to the user.
