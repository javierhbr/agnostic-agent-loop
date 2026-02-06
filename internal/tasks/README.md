# Task Management System

This package implements the core task lifecycle management for the Agentic Agent framework.

## Overview

The task management system handles the complete lifecycle of development tasks, from creation in the backlog through completion and archiving. It provides:

- **Task State Management**: Tracks tasks through their lifecycle (backlog → in-progress → done)
- **Task Decomposition**: Breaks large tasks into smaller, manageable pieces
- **Progress Tracking**: Dual-format progress tracking (human-readable text + machine-readable YAML)
- **Task Locking**: Prevents conflicts when multiple agents work in parallel
- **Learnings Tracking**: Records patterns and insights discovered during task completion

## Key Components

### [`manager.go`](manager.go)
Main task manager coordinating the task lifecycle.

**Key Functions:**
- `ClaimTask(id string)` - Move task from backlog to in-progress
- `CompleteTask(id string)` - Move task from in-progress to done
- `ListTasks(state string)` - List tasks by state
- `DecomposeTask(id string)` - Split large task into subtasks

### [`lock.go`](lock.go)
Task locking mechanism for parallel work.

Prevents race conditions when multiple AI agents or developers work simultaneously. Uses file-based locking to ensure only one agent can claim a task at a time.

### [`decomposer.go`](decomposer.go)
Task decomposition logic.

Analyzes tasks and suggests breakdowns when they exceed complexity thresholds:
- Max 5 files per task
- Max 2 directories per task
- Logical grouping by feature/concern

### [`progress_writer.go`](progress_writer.go)
Dual-format progress tracking.

Writes task progress in two formats:
1. **Text format** (`.agentic/progress.txt`) - Human-readable, easy to scan
2. **YAML format** (`.agentic/progress.yaml`) - Machine-readable, structured data

### [`agents_md_helper.go`](agents_md_helper.go)
AGENTS.md file generation.

Maintains an `AGENTS.md` file with learnings and patterns discovered during development. Helps agents learn from past work and apply consistent patterns.

## Task Lifecycle

```
┌─────────────┐
│   Backlog   │  Tasks waiting to be started
└─────────────┘
      │
      │ claim
      ▼
┌─────────────┐
│ In-Progress │  Tasks currently being worked on
└─────────────┘
      │
      │ complete
      ▼
┌─────────────┐
│    Done     │  Completed tasks
└─────────────┘
      │
      │ archive
      ▼
┌─────────────┐
│   Archive   │  Historical record of completed work
└─────────────┘
```

## Task State Files

Tasks are stored in YAML files under `.agentic/tasks/`:

- **`backlog.yaml`** - Tasks not yet started
- **`in-progress.yaml`** - Tasks currently being worked on
- **`done.yaml`** - Completed tasks (before archiving)

### Task Format

```yaml
tasks:
  - id: "task-001"
    title: "Add user authentication"
    description: "Implement JWT-based authentication"
    status: "in-progress"
    files:
      - src/auth/jwt.go
      - src/middleware/auth.go
    created_at: "2024-01-15T10:30:00Z"
    claimed_at: "2024-01-15T11:00:00Z"
```

## Validation Rules

The task system enforces several quality rules:

1. **Size Limits**:
   - Max 5 files per task
   - Max 2 directories per task
   - Prevents overly complex tasks

2. **Scope Validation**:
   - Tasks must be within git-tracked files
   - Prevents scope creep
   - Ensures focused work

3. **Context Requirements**:
   - Each directory touched must have a `context.md`
   - Ensures adequate context for AI agents

## Usage Examples

### Claim a Task

```bash
agentic-agent task claim task-001
```

This moves the task from `backlog.yaml` to `in-progress.yaml` and creates a lock file.

### Complete a Task

```bash
agentic-agent task complete task-001
```

This moves the task from `in-progress.yaml` to `done.yaml`, releases the lock, and updates progress tracking.

### List Tasks

```bash
# List all backlog tasks
agentic-agent task list backlog

# List in-progress tasks
agentic-agent task list in-progress

# List completed tasks
agentic-agent task list done
```

### Decompose a Large Task

```bash
agentic-agent task decompose task-001
```

Analyzes the task and suggests a breakdown if it's too large or complex.

## Integration with Other Components

### Context System
Tasks integrate with the [context system](../context/README.md) to ensure each directory has adequate context before work begins.

### Validator
The [validator](../validator/README.md) checks tasks against rules before allowing completion.

### Skills Generator
Task patterns inform [skill generation](../skills/README.md) to create reusable AI agent skills.

## Testing

See test files in this directory:
- `manager_test.go` - Task manager unit tests
- `decomposer_test.go` - Task decomposition tests
- `progress_writer_test.go` - Progress tracking tests

Run tests:
```bash
go test ./internal/tasks/...
```

## See Also

- [Context System](../context/README.md) - Directory-level context management
- [Validator](../validator/README.md) - Task validation rules
- [Orchestrator](../orchestrator/README.md) - Agent loop coordination
- [CLI Tutorial](../../docs/guide/CLI_TUTORIAL.md) - User-facing task commands
