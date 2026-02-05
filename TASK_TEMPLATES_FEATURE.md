# Task Templates Feature

## Overview

Two new commands have been added to make task creation faster and easier:

1. **`task sample-task`** - Quick sample task creation
2. **`task from-template`** - Interactive template-based task creation

## Commands

### 1. Create Sample Task

```bash
agentic-agent task sample-task
```

**What it does:**
- Creates a pre-configured sample task instantly
- Perfect for testing the workflow or learning how tasks work
- Includes example description, acceptance criteria, and scope

**Sample task includes:**
- Title: "Implement user authentication system"
- Complete description with implementation details
- 5 acceptance criteria
- Suggested scope directories

**Use cases:**
- Testing the task management workflow
- Learning how tasks are structured
- Quick demo of the system
- Development and testing

### 2. Create from Template

```bash
agentic-agent task from-template
```

**What it does:**
- Launches an interactive wizard to create tasks from built-in templates
- Shows 5 professional templates covering common development tasks
- Allows customization before creation

**Available Templates:**

1. **Feature Implementation**
   - For adding new features
   - Includes implementation checklist
   - Pre-configured acceptance criteria

2. **Bug Fix**
   - Structured bug fix template
   - Issue/Expected/Actual format
   - Regression test reminders

3. **Refactoring**
   - Code quality improvement template
   - Maintainability focus
   - Quality metrics tracking

4. **Documentation**
   - Documentation task template
   - Code examples structure
   - Review checklist

5. **Testing**
   - Test coverage template
   - Unit/Integration/Edge cases
   - CI pipeline integration

**Workflow:**

1. **Select Template**
   - Browse available templates with ↑/↓ keys
   - See description for each template
   - Press Enter to select

2. **Customize**
   - Edit the task title
   - Modify description as needed
   - Template fields (acceptance criteria, scope) are pre-filled

3. **Preview**
   - Review the complete task
   - See all fields before creation
   - Confirm or go back

4. **Create**
   - Task is created in backlog
   - Ready to claim and start working

## Integration with Start Command

The `agentic-agent start` command now includes an interactive menu after project initialization:

```
What would you like to do next?
→ Create your first task
  Start working
  View all tasks
  Exit
```

Selecting "Create your first task" launches the task creation wizard, which can use templates via the `from-template` option.

## Examples

### Quick Sample for Testing
```bash
# Create a sample task instantly
agentic-agent task sample-task

# Output:
# ✓ Sample task created successfully!
#
# Task ID: TASK-001
# Title: Implement user authentication system
#
# Next steps:
# 1. View task: agentic-agent task show TASK-001
# 2. Claim task: agentic-agent task claim TASK-001
# 3. List all tasks: agentic-agent task list
```

### Create Bug Fix from Template
```bash
# Launch template wizard
agentic-agent task from-template

# 1. Select "Bug Fix" template
# 2. Edit title: "Fix: Login fails with special characters"
# 3. Edit description with bug details
# 4. Preview and create

# Task is created with:
# - Structured bug format
# - Pre-filled acceptance criteria
# - Regression test reminders
```

### Create Feature from Template
```bash
# Launch template wizard
agentic-agent task from-template

# 1. Select "Feature Implementation" template
# 2. Edit title: "Implement dark mode toggle"
# 3. Add feature details
# 4. Preview and create

# Task is created with:
# - Implementation checklist
# - Test coverage requirements
# - Documentation reminders
```

## Benefits

### For Junior Developers
- **No memorization required** - Templates guide structure
- **Best practices built-in** - Follow proven patterns
- **Quick start** - Create professional tasks in seconds
- **Learning tool** - See how tasks should be structured

### For All Users
- **Consistency** - All tasks follow similar structure
- **Speed** - Faster than creating from scratch
- **Quality** - Pre-filled acceptance criteria ensure completeness
- **Flexibility** - Fully customizable before creation

## Technical Details

### Template Structure

Templates are defined in `internal/ui/models/template.go` with:

```go
type TaskTemplate struct {
    Name        string              // Template name
    Description string              // Template description
    Task        taskmodels.Task     // Pre-configured task
}
```

### Built-in Templates

Currently 5 built-in templates are available. These are loaded from `loadBuiltInTemplates()` function.

### Future Enhancements

Potential future features:
- Save custom templates
- Import/export templates
- Team template sharing
- Template categories
- Template variables/placeholders

## Commands Summary

| Command | Description | Use Case |
|---------|-------------|----------|
| `task create` | Interactive wizard or flags | Custom task creation |
| `task sample-task` | Instant sample task | Testing, learning, demos |
| `task from-template` | Template-based creation | Professional structured tasks |
| `task list` | View all tasks | Task management |
| `task show <id>` | View task details | Review task information |
| `task claim <id>` | Claim a task | Start working |

## Integration with Workflow

```
┌─────────────────────────────────────┐
│   agentic-agent start               │
│   (Initialize project)              │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│   What would you like to do next?   │
│   → Create your first task          │
└──────────────┬──────────────────────┘
               │
               ├──► task create (custom)
               ├──► task sample-task (instant)
               └──► task from-template (wizard)
                              │
                              ▼
                    ┌───────────────────┐
                    │  Select Template  │
                    └─────────┬─────────┘
                              │
                              ▼
                    ┌───────────────────┐
                    │    Customize      │
                    └─────────┬─────────┘
                              │
                              ▼
                    ┌───────────────────┐
                    │     Preview       │
                    └─────────┬─────────┘
                              │
                              ▼
                    ┌───────────────────┐
                    │   Task Created!   │
                    └───────────────────┘
```

## Conclusion

These new features make task creation:
- **Faster** - Sample tasks in seconds
- **Easier** - No need to remember structure
- **Better** - Professional templates with best practices
- **More accessible** - Perfect for junior developers

The combination of `sample-task` for quick testing and `from-template` for structured creation provides a complete solution for all task creation needs.
