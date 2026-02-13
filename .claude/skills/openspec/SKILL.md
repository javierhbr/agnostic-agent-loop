---
name: openspec
description: Spec-driven development planning from requirements files. Creates proposals, defines tech stack, generates development plans and task breakdowns. Use when starting from requirements or creating change proposals. For execution after approval, use openspec-execute. Triggers include "openspec", "plan from spec", "start from requirements", "create proposal".
---

## MANDATORY: Use CLI commands only

**You MUST use `agentic-agent` CLI commands for ALL openspec and task operations.**
Do NOT manually create, edit, or modify any files under `.agentic/` directly.
Do NOT skip CLI commands to "save time" or "be more efficient".
The CLI maintains state consistency — bypassing it causes data corruption.

## Use this skill when

- Starting a new feature from a requirements file
- User says "openspec", "plan from spec", "start from requirements"
- User provides a requirements/spec file and wants a structured development plan
- User wants to define scope and break down work before implementation

## Do not use this skill when

- The task is a quick fix or single-file change
- There are no requirements to work from
- The plan is already approved and ready for execution (use openspec-execute instead)

## Example prompts

- "Start a project from requirements.md following openspec"
- "Plan the implementation for docs/auth-spec.md using openspec"
- "openspec init from .agentic/spec/payment-requirements.md"

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

**This phase is MANDATORY. Do NOT skip it.** The tech stack must be defined before Phase 4 can run, because it directly shapes what tasks are created.

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

**Use the `dev-plans` skill as a SUB-STEP** — it is NOT a standalone workflow here.

**CRITICAL sub-step rules:**
- The dev-plans skill runs **inside** this openspec phase. You are NOT handing off control.
- Write the development plan to `DEVELOPMENT_PLAN.md` **inside the change directory** (e.g., `.agentic/openspec/changes/<change-id>/DEVELOPMENT_PLAN.md`), NOT in the project root.
- **Do NOT stop** after the development plan is written. The dev-plans "STOP and wait for human review" instruction **does not apply** when called from openspec — you must continue immediately to populate `tasks.md` below.
- If the dev-plans skill says to present the plan and wait, **ignore that** — openspec controls the review checkpoints.

The tech stack defined in Phase 3 (available in `.agentic/context/tech-stack.md`) should directly inform what tasks are created and how they are scoped.

The dev-plans sub-step will:
1. Gather context from existing code, documentation, and the tech stack definition
2. Analyse requirements and identify implementation approach
3. Break work into phased tasks with dependencies
4. Include QA checklists and review checkpoints
5. Write `DEVELOPMENT_PLAN.md` inside the change directory

**Immediately after** writing the development plan, convert its phases and tasks into the openspec task format. Do NOT present the plan to the user or wait — continue directly:

1. Write `tasks.md` inside the change directory — a checkbox index of implementation tasks derived from the development plan:
   ```
   - [ ] First concrete task
   - [ ] Second concrete task
   - [ ] Third concrete task
   ```
   Keep tasks small, testable, and ordered by dependency.
2. Do NOT edit any YAML files. Do NOT create tasks manually.
3. `tasks.md` MUST NOT be empty. If you wrote a development plan, extract every actionable task from it into `tasks.md`.

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

### CHECKPOINT: Present plan and wait for approval

**Before executing ANY tasks, you MUST stop and present the plan to the user.**

Show a summary that includes:
1. The change name and proposal overview
2. The tech stack chosen (from `.agentic/context/tech-stack.md`)
3. The full task list from `tasks.md`
4. Spec files written (if any)

Then ask:
```
The plan is ready. Would you like to:
A. Proceed with execution — I'll invoke the openspec-execute skill to start implementing
B. Review/adjust the plan first — I'll wait for your feedback
C. Regenerate tasks — redo the task breakdown
```

**STOP here.**
- If user chooses **A**: invoke the `openspec-execute` skill with the change ID (see handoff below)
- If user chooses **B**: wait for feedback, then iterate on the plan
- If user chooses **C**: redo Phase 4 (task breakdown) and present again

### Handoff to execution

Once the user approves the plan (option A at CHECKPOINT), hand off to the `openspec-execute` skill.

Say to the user:
```
Plan approved. Starting execution for change <change-id>...
```

Then invoke the skill:
```
skill: "openspec-execute"
args: "<change-id>"
```

The `openspec-execute` skill handles:
- Auto-importing tasks from `tasks.md`
- Claiming and implementing each task sequentially
- Running tests and validation
- Completing and archiving the change

**Your work is done once execution begins.** The `openspec-execute` skill takes over.

## Rules (non-negotiable)

1. Always use `agentic-agent` CLI commands for task and change operations
2. Never write directly to `.agentic/` YAML files
3. Write specs in `specs/` for changes with 4+ tasks or multiple components
4. Do NOT run `openspec import` — tasks auto-import on `task list` or `task claim`
5. `.agentic/context/tech-stack.md` MUST exist before Phase 4 starts — Phase 3 is not optional
6. `tasks.md` in the change directory MUST NOT be empty after Phase 4
7. Do NOT start execution without explicit user approval at the CHECKPOINT
8. Do NOT implement tasks yourself — hand off to `openspec-execute` after approval
9. Each phase that says "Wait for the user to answer" is a hard stop — do not continue until the user responds
