---
name: openspec
description: Spec-driven development from requirements files. Scaffolds changes, writes specs, and tracks implementation. Use when starting from requirements, implementing specs, or managing change proposals. Triggers include "openspec", "implement from spec", "start from requirements", "apply requirements".
---

## MANDATORY: Use CLI commands only

**You MUST use `agentic-agent` CLI commands for ALL openspec and task operations.**
Do NOT manually create, edit, or modify any files under `.agentic/` directly.
Do NOT skip CLI commands to "save time" or "be more efficient".
The CLI maintains state consistency — bypassing it causes data corruption.

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
- "Continue implementing change auth-feature" (resumes Phase 5)

## Instructions

### Phase 0: Ensure skills are installed

Before starting, verify the CLI and skills are set up:
```bash
agentic-agent skills ensure
```
This auto-installs mandatory skill packs and generates rules files. Safe to run multiple times.

### Phase 1: Understand the input

1. Read the requirements file the user provided
2. **Analyse the scope** — identify what the input is asking for:
   - What problem does it solve?
   - What are the main components/features?
   - What's the expected end state?
   - What constraints or dependencies exist?
3. **Ask clarifying questions** — before proceeding, ask the user 3-5 targeted questions to fill gaps:
   - Scope boundaries (what's in/out for v1?)
   - Target audience or users
   - Technical constraints (stack, existing systems, integrations)
   - Priority order if multiple features are described
   - Any known risks or blockers

   Format questions with lettered options for quick answers:
   ```
   1. What is the scope for the first version?
      A. Minimal viable — core functionality only
      B. Full feature set as described
      C. Backend/API only (UI later)
      D. Other: [please specify]
   ```
   Wait for the user to answer before continuing.

### Phase 2: Initialize the change

1. Derive a short name for the change (e.g., "auth-feature", "payment-system")
2. **Run the CLI command** (do NOT create files manually):
   ```bash
   agentic-agent openspec init "<change-name>" --from <requirements-file>
   ```
   This command also creates `agnostic-agent.yaml` if it doesn't exist yet, with the correct paths for spec-driven development (specDirs, openSpecDir, contextDirs, workflow validators).
3. Read the generated proposal at the path printed by the command
4. Edit `proposal.md` — fill in the Problem, Approach, Scope, and Acceptance Criteria sections using the requirements and the user's answers from Phase 1

### Phase 3: Create the development plan

**Use the `creating-development-plans` skill** to analyse the requirements and generate technical tasks.

The dev-plans skill will:
1. Gather context from existing code and documentation
2. Analyse requirements and identify implementation approach
3. Break work into phased tasks with dependencies
4. Include QA checklists and review checkpoints

After the development plan is created, convert its phases and tasks into the openspec format:

1. Write `tasks.md` — a numbered list or checkbox index of implementation tasks derived from the development plan:
   ```
   1. First concrete task
   2. Second concrete task
   3. Third concrete task
   ```
   Keep tasks small, testable, and ordered by dependency.
2. Do NOT edit any YAML files. Do NOT create tasks manually.

#### Detailed task files (for complex changes)

For changes with **4 or more tasks**, create individual task detail files to give agents focused context:

1. Create a `tasks/` subdirectory inside the change directory
2. For each task, create a markdown file (e.g., `tasks/01-project-setup.md`):
   ```markdown
   # Set up project structure

   ## Description
   What this task accomplishes specifically.

   ## Prerequisites
   - None (first task)

   ## Acceptance Criteria
   - [ ] Package.json created
   - [ ] Dev server runs without errors

   ## Technical Notes
   Implementation hints and architecture decisions.
   ```
3. Update `tasks.md` to be an index that references the detail files:
   ```markdown
   - [ ] Set up project structure (ver [tasks/01-project-setup.md](./tasks/01-project-setup.md))
   - [ ] Implement storage (ver [tasks/02-storage.md](./tasks/02-storage.md))
   ```

**Skip detailed files** for simple changes (1-3 tasks) where `tasks.md` titles are sufficient.

### Phase 4: Write detailed specs (for complex changes)

For changes with **4 or more tasks** or **multiple components**, write specification files in the `specs/` directory:

1. Create spec files for each major component or concern:
   - `specs/data-model.md` — database schemas, data structures, relationships
   - `specs/api-design.md` — endpoints, request/response formats, authentication
   - `specs/architecture.md` — component diagram, data flow, integration points
   - `specs/ui-design.md` — wireframes, component hierarchy, user flows
2. Each spec file should include:
   - **Overview** — what this spec covers
   - **Design decisions** — why this approach was chosen
   - **Details** — concrete schemas, endpoints, components, etc.
   - **Dependencies** — what other specs or systems this relates to
3. Reference spec files from `proposal.md` in the Approach section

**Skip this phase** for simple changes (1-3 tasks) where `proposal.md` covers everything.

### Phase 5: Execute tasks sequentially

Tasks are **auto-imported** into the backlog when you run `task list` or `task claim`.
No manual import step is needed.

When tasks have detail files, the imported tasks include:
- Description from the task file
- Acceptance criteria (viewable via `agentic-agent task show <id>`)
- Prerequisites mapped as task inputs
- Technical notes in the description

For each task (check progress with `agentic-agent openspec status <change-id>`):

1. **List tasks**: `agentic-agent task list` (triggers auto-import if needed)
2. **Review details**: `agentic-agent task show <task-id>` (see description, acceptance criteria)
3. **Claim via CLI**: `agentic-agent task claim <task-id>`
4. Implement the work
5. Run tests — verify the task works
6. **Complete via CLI**: `agentic-agent task complete <task-id>`

**Never skip ahead.** Complete and verify each task before starting the next.
**Never modify `.agentic/tasks/*.yaml` files directly.**

### Phase 6: Complete and archive

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
5. Write specs in `specs/` for changes with 4+ tasks or multiple components
6. Do NOT run `openspec import` — tasks auto-import on `task list` or `task claim`
