---
name: ralph-converter
description: "Convert PRDs to task YAML format for the agnostic-agent system. Use when you have an existing PRD and need to convert it to YAML tasks. Triggers on: convert this prd, turn this into tasks, create tasks from prd, ralph yaml."
---

# Ralph PRD to YAML Converter

Converts existing PRDs to the YAML task format that agnostic-agent uses for task management.

---

## The Job

Take a PRD (markdown file) from `.agentic/tasks/` and convert it to YAML task entries in `.agentic/tasks/backlog.yaml`.

---

## Output Format

YAML task entries following the Task struct format:

```yaml
- id: "US-001"
  title: "[Story title]"
  description: "As a [user], I want [feature] so that [benefit]"
  status: "backlog"
  scope:
    - "path/to/file1.go"
    - "path/to/file2.go"
  acceptance:
    - "Criterion 1"
    - "Criterion 2"
    - "Typecheck passes"
```

---

## Story Size: The Number One Rule

**Each story must be completable in ONE iteration (one context window).**

The agent spawns a fresh instance per iteration with no memory of previous work. If a story is too big, the LLM runs out of context before finishing and produces broken code.

### Right-sized stories:
- Add a database column and migration
- Add a UI component to an existing page
- Update a server action with new logic
- Add a filter dropdown to a list

### Too big (split these):
- "Build the entire dashboard" - Split into: schema, queries, UI components, filters
- "Add authentication" - Split into: schema, middleware, login UI, session handling
- "Refactor the API" - Split into one story per endpoint or pattern

**Rule of thumb:** If you cannot describe the change in 2-3 sentences, it is too big.

---

## Story Ordering: Dependencies First

Stories execute in priority order. Earlier stories must not depend on later ones.

**Correct order:**
1. Schema/database changes (migrations)
2. Server actions / backend logic
3. UI components that use the backend
4. Dashboard/summary views that aggregate data

**Wrong order:**
1. UI component (depends on schema that does not exist yet)
2. Schema change

---

## Acceptance Criteria: Must Be Verifiable

Each criterion must be something that can be CHECKED, not something vague.

### Good criteria (verifiable):
- "Add `status` column to tasks table with default 'pending'"
- "Filter dropdown has options: All, Active, Completed"
- "Clicking delete shows confirmation dialog"
- "Typecheck passes"
- "Tests pass"

### Bad criteria (vague):
- "Works correctly"
- "User can do X easily"
- "Good UX"
- "Handles edge cases"

### Always include as final criterion:
```
"Typecheck passes"
```

For stories with testable logic, also include:
```
"Tests pass"
```

### For stories that change UI, also include:
```
"Verify in browser using dev-browser skill"
```

Frontend stories are NOT complete until visually verified.

---

## Conversion Rules

1. **Extract user stories from PRD** - Look for sections like "User Stories" or numbered requirements
2. **Generate sequential IDs** - US-001, US-002, etc.
3. **Map fields correctly**:
   - PRD story title → YAML `title`
   - PRD story description → YAML `description`
   - PRD acceptance criteria → YAML `acceptance` array
   - Infer file paths from story context → YAML `scope` array
   - All stories start with `status: "backlog"`
4. **Order by dependency** - Schema first, backend second, UI last
5. **Always add** "Typecheck passes" to every story's acceptance criteria
6. **Validate story size** - If a story seems too large, split it into multiple stories

---

## Splitting Large PRDs

If a PRD has big features, split them:

**Original:**
> "Add user notification system"

**Split into:**
1. US-001: Add notifications table to database
2. US-002: Create notification service for sending notifications
3. US-003: Add notification bell icon to header
4. US-004: Create notification dropdown panel
5. US-005: Add mark-as-read functionality
6. US-006: Add notification preferences page

Each is one focused change that can be completed and verified independently.

---

## Example

**Input PRD** (from `.agentic/tasks/prd-task-status.md`):
```markdown
# PRD: Task Status Feature

## User Stories

### US-001: Add status field to tasks table
**Description:** As a developer, I need to store task status in the database.

**Acceptance Criteria:**
- [ ] Add status column: 'pending' | 'in_progress' | 'done' (default 'pending')
- [ ] Generate and run migration successfully
- [ ] Typecheck passes

### US-002: Display status badge on task cards
**Description:** As a user, I want to see task status at a glance.

**Acceptance Criteria:**
- [ ] Each task card shows colored status badge
- [ ] Badge colors: gray=pending, blue=in_progress, green=done
- [ ] Typecheck passes
- [ ] Verify in browser using dev-browser skill
```

**Output YAML** (appended to `.agentic/tasks/backlog.yaml`):
```yaml
- id: "1738800000"  # Generate using timestamp or UUID
  title: "Add status field to tasks table"
  description: "As a developer, I need to store task status in the database."
  status: "backlog"
  scope:
    - "internal/database/schema.sql"
    - "internal/database/migrations/"
  acceptance:
    - "Add status column: 'pending' | 'in_progress' | 'done' (default 'pending')"
    - "Generate and run migration successfully"
    - "Typecheck passes"

- id: "1738800001"
  title: "Display status badge on task cards"
  description: "As a user, I want to see task status at a glance."
  status: "backlog"
  scope:
    - "internal/ui/components/task_card.go"
    - "internal/ui/components/badge.go"
  acceptance:
    - "Each task card shows colored status badge"
    - "Badge colors: gray=pending, blue=in_progress, green=done"
    - "Typecheck passes"
    - "Verify in browser using dev-browser skill"
```

---

## Inferring Scope (File Paths)

Based on the story description, infer likely file paths:

- **Database/Schema changes** → `internal/database/`, `pkg/models/`, migration files
- **API/Server actions** → `internal/handlers/`, `internal/services/`
- **UI components** → `internal/ui/components/`, `internal/ui/models/`
- **CLI commands** → `cmd/agentic-agent/`
- **Business logic** → `internal/tasks/`, `internal/context/`, `internal/validator/`

If unsure, leave scope empty or use broad directory paths like `internal/ui/`.

---

## Workflow

1. **Read the PRD** from `.agentic/tasks/prd-[feature-name].md`
2. **Extract all user stories** from the PRD
3. **Validate story size** - Split if needed
4. **Order by dependencies** - Database → Backend → UI
5. **Generate YAML entries** with all required fields
6. **Append to backlog** - Add entries to `.agentic/tasks/backlog.yaml`
7. **Verify** - Ensure YAML is valid and all stories accounted for

---

## Checklist Before Saving

Before writing tasks to backlog.yaml, verify:

- [ ] Each story is completable in one iteration (small enough)
- [ ] Stories are ordered by dependency (schema to backend to UI)
- [ ] Every story has "Typecheck passes" as criterion
- [ ] UI stories have "Verify in browser using dev-browser skill" as criterion
- [ ] Acceptance criteria are verifiable (not vague)
- [ ] No story depends on a later story
- [ ] Scope paths are inferred where possible
- [ ] All IDs are unique

---

## Important Notes

- **Always validate YAML syntax** before saving
- **Preserve existing tasks** in backlog.yaml (append, don't replace)
- **Use consistent ID format** (timestamp-based or sequential)
- **Review with user** if story splitting is extensive
