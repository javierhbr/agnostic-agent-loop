# Ralph PDR Workflow Guide

## Overview

This guide explains how to use the Ralph Plan-Do-Review (PDR) methodology integrated into the agnostic-agent framework.

## The PDR Cycle

### 📋 Plan: Create PRD

Use the PRD skill to create a structured Product Requirements Document:

```bash
# Generate Claude Code skills first (one-time setup)
agentic-agent skills generate-claude-skills

# Now use the /prd skill in Claude Code
# This will create .agentic/tasks/prd-[feature-name].md
```

**PRD Structure:**
- Introduction/Overview
- Goals (measurable objectives)
- User Stories (with acceptance criteria)
- Functional Requirements
- Non-Goals (scope boundaries)
- Success Metrics

**Key Features:**
- **Lettered Q&A**: Answer with "1A, 2C, 3B" for quick iteration
- **Story Sizing**: Each story should be completable in one context window
- **Browser Verification**: UI stories must include "Verify in browser" criteria
- **Dependency Ordering**: Database → Backend → UI

### 🔄 Convert: PRD to Tasks

Use the Ralph converter skill to transform PRD into YAML tasks:

```bash
# Use the /ralph-converter skill in Claude Code
# This reads from .agentic/tasks/prd-[feature].md
# And creates entries in .agentic/tasks/backlog.yaml
```

**Conversion Features:**
- Maps user stories to Task YAML format
- Validates story size and dependencies
- Infers file paths from story descriptions
- Enforces acceptance criteria best practices

### ⚙️ Do: Implement with Context

Before implementing, **always read context.md** files:

```bash
# Generate context for directories you'll modify
agentic-agent context generate internal/tasks

# Claim a task
agentic-agent task claim TASK-123

# Work on the task (read context.md first!)
# - Check .agentic/spec/ for specifications
# - Read directory context.md for patterns
# - Follow the rules in CLAUDE.md or agent-specific rules
```

**Key Practices:**
1. **Context First**: Read `context.md` before editing any directory
2. **Scope Awareness**: Only modify files in task.Scope
3. **Acceptance Criteria**: Complete ALL criteria before marking done
4. **Browser Verification**: For UI changes, verify visually (not just typecheck)

### ✅ Review: Track Progress & Learnings

When completing a task, track your learnings:

```bash
# Complete task with learnings
agentic-agent task complete TASK-123

# Manually add codebase patterns
agentic-agent learnings add "Always validate UI files in task scope have browser verification"

# View progress
agentic-agent learnings show --limit 10
agentic-agent learnings list
```

**Progress Tracking:**
- **progress.txt**: Human-readable log with learnings
- **progress.yaml**: Machine-readable metadata
- **Codebase Patterns**: Consolidated learnings at top of progress.txt

## Integration with Context Files

### context.md vs AGENTS.md vs progress.txt

| File | Purpose | Scope | Updated When |
|------|---------|-------|--------------|
| **context.md** | Architectural rules, responsibilities | Per-directory | Architecture changes |
| **AGENTS.md** | Agent-specific patterns, gotchas | Per-directory | Task completion (optional) |
| **progress.txt** | Codebase-wide learnings, progress log | Global | Task completion |

### Example Workflow

1. **Read context.md** - Understand directory architecture
   ```bash
   cat internal/tasks/context.md
   ```

2. **Check AGENTS.md** (if exists) - Learn directory-specific patterns
   ```bash
   cat internal/tasks/AGENTS.md
   ```

3. **Review Codebase Patterns** - Apply global learnings
   ```bash
   agentic-agent learnings list
   ```

4. **Implement** - Follow all guidelines

5. **Update learnings** - Record new discoveries
   ```bash
   agentic-agent learnings add "TaskManager requires progress writer for tracking"
   ```

6. **Update AGENTS.md** (if significant) - Add directory-specific pattern
   ```
   # In internal/tasks/AGENTS.md
   ## Patterns
   - Always use NewTaskManagerWithTracking when progress logging is needed
   - Progress writer paths come from config.Paths
   ```

## Validation Rules

### Browser Verification Rule

Ensures UI tasks have proper verification:

```yaml
# In agnostic-agent.yaml
workflow:
  validators:
    - context-check      # Ensures context.md exists
    - task-scope         # Validates files in scope
    - browser-verification  # NEW: Checks UI tasks have browser criteria
```

**Triggered when:**
- Task scope includes UI files (`.tsx`, `.jsx`, `.vue`, `.html`, `.css`, etc.)
- Task scope includes files in UI directories (`/components/`, `/ui/`, `/pages/`)

**Passes when:**
- Acceptance criteria includes "verify in browser", "browser verification", or "visual verification"

**Example:**
```yaml
# ❌ FAIL - UI task without browser verification
- id: "US-002"
  title: "Add button component"
  scope:
    - "internal/ui/components/button.tsx"
  acceptance:
    - "Button renders with correct styling"
    - "Typecheck passes"

# ✅ PASS - UI task with browser verification
- id: "US-002"
  title: "Add button component"
  scope:
    - "internal/ui/components/button.tsx"
  acceptance:
    - "Button renders with correct styling"
    - "Typecheck passes"
    - "Verify in browser using dev-browser skill"
```

## Orchestrator & Archiving

### Autonomous Loop

For autonomous multi-task execution:

```go
import "github.com/javierbenavides/agentic-agent/internal/orchestrator"

// Create loop with stop conditions
loop := orchestrator.NewLoop(
    10,  // max iterations
    "<promise>COMPLETE</promise>",  // stop signal
    taskManager,
)

// Run loop
err := loop.Run(context.Background())
```

**Features:**
- Runs multiple iterations automatically
- Detects stop signal in agent output
- Checks if all tasks completed
- Returns after max iterations or completion

### Branch-Based Archiving

Automatically archive progress when switching features:

```go
import "github.com/javierbenavides/agentic-agent/internal/orchestrator"

// Create archiver
archiver := orchestrator.NewArchiver(
    ".agentic/archive/",
    ".agentic/progress.txt",
    ".agentic/progress.yaml",
    ".agentic/tasks/",
)

// Archive if branch changed (typically at start of new feature)
err := archiver.ArchiveIfBranchChanged("feature/new-auth-flow")
```

**Archive Format:**
```
.agentic/archive/
└── 2026-02-05-feature-previous-feature/
    ├── progress.txt
    ├── progress.yaml
    └── done.yaml
```

**Automatic Reset:**
- Creates fresh progress.txt with new header
- Removes progress.yaml (recreated on first entry)
- Preserves completed tasks in archive

## Best Practices

### 1. Story Sizing
- **Right-sized**: "Add status column to tasks table with migration"
- **Too large**: "Build entire dashboard" (split into: schema, queries, UI, filters)
- **Rule of thumb**: If you can't describe it in 2-3 sentences, it's too big

### 2. Acceptance Criteria
- **Good (verifiable)**: "Filter dropdown has options: All, Active, Completed"
- **Bad (vague)**: "Works correctly"
- **Always include**: "Typecheck passes"
- **For UI**: "Verify in browser using dev-browser skill"

### 3. Dependency Order
- ✅ **Correct**: Schema → Backend → UI → Dashboard
- ❌ **Wrong**: UI component → Schema (UI depends on schema that doesn't exist yet)

### 4. Progress Logging
- Log **what changed** (files modified)
- Log **why it matters** (learnings for future work)
- Log **what to avoid** (gotchas encountered)
- Include **thread URL** if working in Claude Code threads

### 5. Codebase Patterns
Add to Codebase Patterns section when you discover:
- Reusable architectural patterns
- Non-obvious requirements
- Common gotchas
- Testing approaches
- Configuration patterns

**Don't add:**
- Story-specific details
- Temporary debugging notes
- Information already in context.md

## Example: Complete PDR Cycle

### 1. Plan
```markdown
# PRD: Task Status Feature

## User Stories

### US-001: Add status field to database
**Description:** As a developer, I need to store task status.

**Acceptance Criteria:**
- [ ] Add status column: 'pending' | 'in_progress' | 'done'
- [ ] Generate and run migration
- [ ] Typecheck passes

### US-002: Display status badge on cards
**Description:** As a user, I want to see task status at a glance.

**Acceptance Criteria:**
- [ ] Each task card shows colored status badge
- [ ] Badge colors: gray=pending, blue=in_progress, green=done
- [ ] Typecheck passes
- [ ] Verify in browser using dev-browser skill
```

### 2. Convert to Tasks
```yaml
# In .agentic/tasks/backlog.yaml
- id: "1738800000"
  title: "Add status field to database"
  description: "As a developer, I need to store task status."
  status: "backlog"
  scope:
    - "internal/database/schema.sql"
    - "internal/database/migrations/"
  acceptance:
    - "Add status column: 'pending' | 'in_progress' | 'done'"
    - "Generate and run migration"
    - "Typecheck passes"

- id: "1738800001"
  title: "Display status badge on cards"
  description: "As a user, I want to see task status at a glance."
  status: "backlog"
  scope:
    - "internal/ui/components/task_card.tsx"
  acceptance:
    - "Each task card shows colored status badge"
    - "Badge colors: gray=pending, blue=in_progress, green=done"
    - "Typecheck passes"
    - "Verify in browser using dev-browser skill"
```

### 3. Do (Implement)
```bash
# Claim task
agentic-agent task claim 1738800000

# Read context
cat internal/database/context.md

# Implement (following context.md rules)
# - Add migration file
# - Update schema
# - Run typecheck

# Verify acceptance criteria
```

### 4. Review (Log Progress)
```bash
# Complete with learnings
agentic-agent task complete 1738800000

# Add global pattern
agentic-agent learnings add "Use IF NOT EXISTS for all migrations"

# View progress
agentic-agent learnings show
```

### 5. Repeat for US-002
```bash
# Claim UI task
agentic-agent task claim 1738800001

# Read context
cat internal/ui/components/context.md

# Check for AGENTS.md patterns
cat internal/ui/components/AGENTS.md

# Implement
# - Add badge component
# - Update task card
# - Run typecheck
# - VERIFY IN BROWSER (required!)

# Complete
agentic-agent task complete 1738800001
agentic-agent learnings add "Badge component reuses existing styles from theme"
```

## Summary

The Ralph PDR integration provides:
- ✅ Structured PRD creation with validated user stories
- ✅ Automatic PRD to task conversion
- ✅ Progress tracking with dual formats (text + YAML)
- ✅ Learnings consolidation (Codebase Patterns)
- ✅ Directory-specific pattern tracking (AGENTS.md)
- ✅ Browser verification for UI changes
- ✅ Orchestrator with stop conditions
- ✅ Branch-based archiving

Combined with the existing `context.md` pattern, this creates a complete knowledge management system for agent-driven development.
