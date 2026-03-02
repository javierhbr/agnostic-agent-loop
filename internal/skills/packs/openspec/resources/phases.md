# OpenSpec Phases — Detailed Reference

This document contains the full detail for each phase of the OpenSpec workflow. Reference this when executing a phase — the slim SKILL.md will direct you here.

---

## Phase 0.5: Refine Requirements into Formal Spec Document

**This phase is mandatory.** Do NOT skip to Phase 1 without a clear, written requirements document.

If the user has a vague idea but no formal requirements file yet:

### Step 1: Present the refinement options

Present with this exact prompt:

```
✓ REQUIRED: Before proceeding, we need a formal requirements document.

You have three paths:

A. Brainstorm first (RECOMMENDED)
   - Use the brainstorming skill to explore intent, design, and scope
   - Results in clear, documented ideas
   - Then proceed to create a PRD

B. Create a PRD directly
   - Use the product-wizard skill to formalize requirements
   - Saves to `.agentic/spec/prd-<feature>.md`
   - Becomes your requirements for Phase 1

C. Use existing written requirements
   - You already have a spec/requirements file
   - Provide the file path

Which would you like to do? [A / B / C]
```

### Step 2: Complete one of these workflows based on user selection

**If A (Brainstorm):**
- Invoke the **brainstorming** skill
- Wait for brainstorm output to stabilize
- Ask: "Are you satisfied with the brainstorm results?"
- Once confirmed, proceed to B below

**If B (Create PRD):**
- Invoke the **product-wizard** skill
- The skill will save to `.agentic/spec/prd-<feature>.md`
- Wait for PRD completion
- Ask user: "Is the PRD complete and accurate?"
- Confirm user approval before proceeding

**If C (Use existing):**
- Verify the file exists and is readable
- Proceed to Phase 1

### Step 3: Confirm completion

```
✓ Phase 0.5 COMPLETE
Requirements file: [path]
Ready to proceed to Phase 1: Understand the input
```

**You must have a formal requirements document before proceeding to Phase 1. Do NOT continue without one.**

---

## Phase 1: Understand the Input Deeply

**This phase is mandatory.** You must fully understand and document the requirements before Phase 2.

### Step 1: Read and analyze the requirements file completely

Thoroughly review the entire file (requirements, PRD, spec, etc.). Create detailed notes on:
- **Problem statement**: What problem does this solve?
- **Main features/components**: What are the major pieces?
- **End state**: What should be built and delivered?
- **Constraints**: Technical, business, or timeline constraints?
- **Dependencies**: What systems or integrations are needed?

### Step 2: Ask 3-5 clarifying questions

Present these questions with lettered options for quick answers:

```
Based on the requirements, I have clarifying questions:

1. What is the scope for the FIRST VERSION (v1)?
   A. Minimal viable — core functionality only
   B. Full feature set as described
   C. Backend/API only (UI later)
   D. Other: [please specify]

2. Who is the PRIMARY USER/AUDIENCE for this?
   A. Internal team/users
   B. External customers
   C. Public/open-source community
   D. Other: [please specify]

3. If multiple features exist, what is the PRIORITY ORDER?
   A. Features are equally important
   B. [List user's prioritized order from requirements]
   C. Other: [please specify]

4. What are the KNOWN RISKS or BLOCKERS?
   A. No known blockers
   B. [List any identified risks]
   C. Other: [please specify]

5. Is there any specific TIMELINE or DEADLINE?
   A. No deadline — iterative development
   B. [Specific date/milestone]
   C. Other: [please specify]
```

### Step 3: Wait for user answers and document them

- Capture all answers exactly as provided
- If any answer is unclear, ask a follow-up question
- Do NOT proceed until you have clear answers to all questions

### Step 4: Confirm Phase 1 completion

```
✓ Phase 1 COMPLETE - Understanding Confirmed

Requirements file: [path]
Problem: [1-sentence summary]
Main components: [list]
Key constraints: [list]
User priorities: [documented]

Ready to proceed to Phase 2: Initialize the change
```

**Do NOT skip clarifying questions. A clear understanding now prevents rework later.**

---

## Phase 3: Define the Tech Stack

**This phase is mandatory.** Do NOT proceed to Phase 4 without a documented, user-approved tech stack.

Before planning tasks, you must establish and document the technical foundation for this change.

### Step 1: Check for an existing tech stack definition

- Read `.agentic/context/tech-stack.md` if it exists
- If it exists and contains detailed, non-template content:
  ```
  I found an existing tech stack definition:
  [show full contents]

  Is this still accurate for THIS specific change?
  Any changes needed? [Yes / No]
  ```
  - If user says **Yes**, apply changes and proceed to step 3
  - If user says **No**, proceed to step 2 below
- If the file doesn't exist, is empty, or only has placeholders: proceed to step 2

### Step 2: Ask the user to define their tech stack

Present these questions and wait for complete answers:

```
✓ TECH STACK DEFINITION

I need to understand the technical foundation for this change.
Please answer each question:

1. What LANGUAGE(S) will this project use?
   A. TypeScript/JavaScript (Node.js 20+)
   B. Go
   C. Python
   D. Other: [please specify]

2. What FRONTEND FRAMEWORK (if any)?
   A. React (Next.js, Vite, etc.)
   B. Vue (Nuxt, etc.)
   C. Svelte (SvelteKit)
   D. Vanilla JS/HTML
   E. None — backend/CLI only

3. What DATA STORAGE/DATABASE?
   A. PostgreSQL (relational, self-hosted)
   B. SQLite / local storage
   C. MongoDB (document store)
   D. None needed / in-memory only
   E. Other: [please specify]

4. INFRASTRUCTURE & DEPLOYMENT target?
   A. Docker + cloud (AWS/GCP/Azure)
   B. Vercel/Netlify (serverless frontend)
   C. Local only / CLI tool
   D. Other: [please specify]

5. Any OTHER tech decisions (testing, styling, build tools)?
   A. Use standard recommendations (Vitest, Tailwind, Vite)
   B. Custom setup: [please describe]
```

- Wait for complete answers to all questions
- If answers are vague, ask follow-up clarifications
- Get explicit user confirmation before proceeding

### Step 3: Document the tech stack decision

Create/update `.agentic/context/tech-stack.md` with the chosen stack:

```markdown
# Tech Stack for [change-name]

- **Language**: [chosen language]
- **Frontend Framework**: [framework or "None"]
- **Database**: [database choice]
- **Storage/ORM**: [if applicable]
- **API Layer**: [REST/GraphQL/gRPC]
- **Styling**: [Tailwind, CSS Modules, etc.]
- **Testing**: [Vitest, Jest, etc.]
- **Build Tool**: [Vite, webpack, tsc, etc.]
- **Deployment**: [target and method]
- **Key Libraries**: [any critical dependencies]
```

### Step 4: Update proposal.md with the tech stack

- Edit the `proposal.md` file generated in Phase 2
- Add a **Tech Stack** subsection under the Approach section
- Include the full stack details from tech-stack.md
- This becomes the permanent technical record for this change

### Step 5: Confirm Phase 3 completion

```
✓ Phase 3 COMPLETE - Tech Stack Defined

Tech stack file: `.agentic/context/tech-stack.md`
Proposal updated: [change-name]/proposal.md

Ready to proceed to Phase 4: Create the development plan
```

**The tech stack is the foundation for all tasks. Get it right before planning.**

---

## Phase 4: Create the Development Plan and Task Breakdown

**This phase is mandatory.** Do NOT proceed to Phase 5+ without a complete, documented development plan and task list.

### Step 1: Use the development planning process

- **Recommended**: Invoke the **superpowers:writing-plans** skill (if available) for structured planning with human review checkpoints
- **Fallback**: Use the in-house `dev-plans` skill or create the plan manually
- Provide the skill with:
  - Requirements from Phase 1 (with user's clarifications)
  - Tech stack from Phase 3 (from `.agentic/context/tech-stack.md`)
  - User's scope/priority answers from Phase 1
- The skill will output a detailed development plan

### Step 2: Analyze the development plan output

- Review all phases, tasks, and dependencies identified
- Understand the implementation sequence
- Confirm the scope matches Phase 1 goals

### Step 3: Create `tasks.md` within the change directory

**CRITICAL: Do NOT skip this step. This is what agents will execute.**

Create a `tasks.md` file in the change directory with:

```markdown
# Implementation Tasks for [change-name]

## Task Breakdown (from development plan)

1. [Task 1: description]
2. [Task 2: description]
3. [Task 3: description]
4. [Task 4: description]
... (continue for all tasks)

## Task Granularity Rules Applied
- [ ] Each task is atomic (one focused unit of work)
- [ ] Each task is testable in isolation
- [ ] Tasks ordered by dependencies
- [ ] No task bundles frontend + backend (split as two tasks)
- [ ] No task bundles infrastructure + business logic (split as two tasks)
- [ ] 3-5 acceptance criteria per task maximum
```

**Ordering rules:**
- List all foundational tasks first (setup, infrastructure)
- Then feature/feature-related tasks
- Then QA/testing tasks
- Keep tasks ordered by dependency

### Step 4: For changes with 4 or more tasks, create detailed task files

Create a `tasks/` directory inside the change directory. For each task, create a numbered markdown file (e.g., `tasks/01-setup.md`):

```markdown
# Task Title

## Description
[What this task accomplishes and why it matters]

## Prerequisites
- [Previous task(s) that must complete first]

## Acceptance Criteria
- [ ] Criteria 1
- [ ] Criteria 2
- [ ] Criteria 3 (max 5)

## Technical Notes
[Implementation hints, architectural decisions, gotchas]

## Related Specs
[Links to relevant spec files if they exist]
```

Then update `tasks.md` to reference detail files:

```markdown
- [ ] Task 1 (see `tasks/01-task-name.md`)
- [ ] Task 2 (see `tasks/02-task-name.md`)
```

### Step 5: Trigger task auto-import into the backlog

```bash
agentic-agent task list
```

**CRITICAL: Do NOT skip this.** Without running this command, tasks won't appear in the backlog and agents won't have work to do.

### Step 6: Verify task import succeeded

- Run: `agentic-agent task list`
- Confirm all tasks appear in the backlog
- If any are missing, edit `tasks.md` and run `agentic-agent task list` again

### Step 7: Confirm Phase 4 completion

```
✓ Phase 4 COMPLETE - Development Plan Created

Tasks file: [change-name]/tasks.md
Task count: [number]
Tasks imported: Yes
All tasks visible in backlog: Yes

Ready to proceed to Phase 5: Write detailed specs (if 4+ tasks)
or Phase 6: Execute tasks sequentially
```

**Without a complete task breakdown and import, agents have no work to execute. This phase is not optional.**

### Task Granularity Rules

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

### Detailed Task Files (for Complex Changes)

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
   - [ ] Set up project structure (see `tasks/01-project-setup.md`)
   - [ ] Implement storage (see `tasks/02-storage.md`)
   ```

**Skip detailed files** for simple changes (1-3 tasks) where `tasks.md` titles are sufficient.

After writing `tasks.md` (and optional detail files), **immediately trigger the auto-import** so tasks appear in the backlog:

```bash
agentic-agent task list
```

This imports tasks from `tasks.md` into `.agentic/tasks/backlog.yaml`. Do NOT skip this step — without it, the task backlog will be empty.

---

## Phase 5: Write Detailed Specs (for Complex Changes)

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

---

## Phase 6: Execute Tasks Sequentially

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
4. **Execute with optional Superpowers skills** (if available):
   - **superpowers:executing-plans** — automated execution with verification checkpoints
   - **superpowers:test-driven-development** — enforce TDD (write tests first, then code)
   - **superpowers:using-git-worktrees** — create isolated workspace for the task
   - Or use in-house `run-with-ralph` skill for iterative convergence
5. Implement the work
6. Run tests — verify the task works
7. **Complete via CLI**: `agentic-agent task complete <task-id>`

**Never skip ahead.** Complete and verify each task before starting the next.
**Never modify `.agentic/tasks/*.yaml` files directly.**

---

## Phase 7: Final Verification and Completion

When all tasks are done:

1. **Optional: Run final verification with Superpowers (if available):**
   - **superpowers:verification-before-completion** — hard gate requiring evidence before marking done
   - Or use in-house `sdd/verifier` skill for comprehensive verification

2. **Complete and archive via CLI:**
   ```bash
   agentic-agent openspec complete <change-id>
   agentic-agent openspec archive <change-id>
   ```

3. **Report the summary to the user.**
