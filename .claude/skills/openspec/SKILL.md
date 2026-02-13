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
- "Continue implementing change auth-feature" (resumes Phase 6)

## Instructions

### Phase 0: Ensure skills are installed

Before starting, verify the CLI and skills are set up:
```bash
agentic-agent skills ensure
```
This auto-installs mandatory skill packs and generates rules files. Safe to run multiple times.

### Phase 0.5: Refine requirements (if starting from scratch)

If the user has a vague idea but no requirements file yet:

1. **Ask the user** if they'd like to refine their idea first:
   ```
   I notice you don't have a requirements file yet. Would you like to:
   A. Brainstorm the idea first (recommended — explores intent and design)
   B. Create a PRD directly (formalizes requirements into a document)
   C. Skip — I'll work from what you've described
   ```
2. If **A**: Use the **brainstorming** skill to explore the idea. Once the design is clear, proceed to B.
3. If **B**: Use the **product-wizard** skill to create a PRD. It will save to `.agentic/spec/prd-<feature>.md`. Then use that file with `openspec init --from .agentic/spec/prd-<feature>.md`.
4. If **C**: Proceed directly to Phase 1.

Always confirm with the user before moving to the next step.

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

### Phase 3: Define the tech stack

Before planning tasks, establish the technical foundation for this change.

1. **Check for an existing tech stack definition** — read `.agentic/context/tech-stack.md`
   - If the file exists and has real content (not just template comments): show it to the user and ask:
     ```
     I found an existing tech stack definition:
     [show contents]

     Is this still accurate for this change? Anything to add or change?
     ```
   - If the file doesn't exist, is empty, or only has template placeholders: proceed to discovery below.

2. **Ask the user about their tech stack** — present structured questions with lettered options:
   ```
   1. What language(s) will this project use?
      A. TypeScript/JavaScript (Node.js)
      B. Go
      C. Python
      D. Other: [please specify]

   2. What frontend framework (if any)?
      A. React (Next.js, Vite, etc.)
      B. Vue (Nuxt, etc.)
      C. Svelte (SvelteKit)
      D. None — backend/CLI only

   3. What data storage?
      A. PostgreSQL
      B. SQLite / local storage
      C. MongoDB
      D. None needed / TBD

   4. Any specific infrastructure or deployment target?
      A. Docker + cloud (AWS/GCP/Azure)
      B. Vercel/Netlify (serverless)
      C. Local only / CLI tool
      D. Other: [please specify]
   ```
   Wait for the user to answer before continuing.

3. **Write or update `.agentic/context/tech-stack.md`** with the user's answers:
   ```markdown
   # Tech Stack

   - **Language**: TypeScript (Node.js 20+)
   - **Frontend**: React 18 with Vite
   - **Storage**: localStorage (offline-first), sync via REST API
   - **Styling**: Tailwind CSS
   - **Testing**: Vitest + Testing Library
   - **Build**: Vite
   ```
   Create the `.agentic/context/` directory if it doesn't exist.

4. **Update `proposal.md`** — add a Tech Stack subsection under Approach with the chosen stack, so it becomes part of the change's permanent record.

### Phase 4: Create the development plan

**Use the `creating-development-plans` skill** to analyse the requirements and generate technical tasks.

The tech stack defined in Phase 3 (available in `.agentic/context/tech-stack.md`) should directly inform what tasks are created and how they are scoped.

The dev-plans skill will:
1. Gather context from existing code, documentation, and the tech stack definition
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
   Keep tasks ordered by dependency.
2. Do NOT edit any YAML files. Do NOT create tasks manually.

#### Task granularity rules

Each task must be **atomic** — one focused unit of work that an agent can complete in a single session.

| Guideline | Example |
|-----------|---------|
| One concern per task | "Add user model" not "Add user model and auth endpoints" |
| Testable in isolation | Task has its own acceptance criteria that can be verified independently |
| Single directory/layer | Prefer "Add API routes" + "Add UI components" over "Add full feature" |
| 10-20 tasks for a medium project | A feature with 3 screens, an API, and storage should produce ~15 tasks |
| 3-5 acceptance criteria per task | More than 5 criteria means the task should be split |

**Split aggressively.** A task that touches both frontend and backend should be two tasks. A task that sets up infrastructure AND implements business logic should be two tasks. When in doubt, split.

**Do NOT bundle** multiple features, layers, or components into one task just to keep the count low.

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

   ## Skills
   - tdd
   ```
3. **Ask the user which skills each task needs.** Present the available skill packs and ask which ones apply:
   ```
   Which skill packs should each task use?
   Available: tdd, api-docs, code-simplification, dev-plans, diataxis, extract-wisdom, openspec

   Task 1 (Set up project structure): [none / list packs]
   Task 2 (Implement storage): [none / list packs]
   ```
   Add the answers as a `## Skills` section in each task detail file. If the user says "none" or skips, omit the section — all installed skills load by default.
4. Update `tasks.md` to be an index that references the detail files:
   ```markdown
   - [ ] Set up project structure (ver [tasks/01-project-setup.md](./tasks/01-project-setup.md))
   - [ ] Implement storage (ver [tasks/02-storage.md](./tasks/02-storage.md))
   ```

**Skip detailed files** for simple changes (1-3 tasks) where `tasks.md` titles are sufficient.

After writing `tasks.md` (and optional detail files), **immediately trigger the auto-import** so tasks appear in the backlog:
```bash
agentic-agent task list
```
This imports tasks from `tasks.md` into `.agentic/tasks/backlog.yaml`. Do NOT skip this step — without it, the task backlog will be empty.

### Phase 5: Write detailed specs (for complex changes)

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

### Phase 6: Execute tasks sequentially

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

### Phase 7: Complete and archive

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
