# Command Quick Reference - Dual Mode Support

All commands support both **Interactive Mode** (no flags) and **Flag Mode** (with arguments).

## Quick Start

### Interactive Mode (Beginner-Friendly)
Just run the command without any arguments:
```bash
agentic-agent <command>
```
The CLI will guide you through the process with menus, selectors, and prompts.

### Flag Mode (Power User)
Run with arguments for direct execution:
```bash
agentic-agent <command> [args] [flags]
```
Perfect for scripts, automation, and CI/CD.

---

## Command Reference

### 📋 Task Management

| Command | Interactive | Flag Mode |
|---------|------------|-----------|
| **claim** | `agentic-agent task claim` | `agentic-agent task claim <task-id>` |
| **complete** | `agentic-agent task complete` | `agentic-agent task complete <task-id>` |
| **show** | `agentic-agent task show` | `agentic-agent task show <task-id>` |
| **create** | `agentic-agent task create` | `agentic-agent task create --title "..." [options]` |
| **list** | `agentic-agent task list` | `agentic-agent task list --no-interactive` |
| **decompose** | `agentic-agent task decompose` | `agentic-agent task decompose <id> <sub1> <sub2>...` |
| **from-template** | `agentic-agent task from-template` | `agentic-agent task from-template --template <name> --title "..."` |

### 📁 Context Management

| Command | Interactive | Flag Mode |
|---------|------------|-----------|
| **generate** | `agentic-agent context generate` | `agentic-agent context generate <dir>` |
| **update** | `agentic-agent context update` | `agentic-agent context update <dir>` |
| **scan** | `agentic-agent context scan` | `agentic-agent context scan` |
| **build** | `agentic-agent context build` | `agentic-agent context build --task <id> --format <fmt>` |

### 🛠️ Skills & Tools

| Command | Interactive | Flag Mode |
|---------|------------|-----------|
| **generate** | `agentic-agent skills generate` | `agentic-agent skills generate --tool <name>` or `--all` |
| **check** | `agentic-agent skills check` | `agentic-agent skills check` |

### ✅ Validation & Tokens

| Command | Interactive | Flag Mode |
|---------|------------|-----------|
| **validate** | `agentic-agent validate` | `agentic-agent validate --format <text|json>` |
| **token status** | `agentic-agent token status` | `agentic-agent token status --no-interactive` |

### 🔄 Workflows

| Command | Interactive | Flag Mode |
|---------|------------|-----------|
| **work** | `agentic-agent work` | `agentic-agent work --task <id> [--skip-context-gen]` |
| **run** | `agentic-agent run` | `agentic-agent run --task <id>` |

---

## Common Patterns

### Force Flag Mode
Use `--no-interactive` to force flag mode even without arguments:
```bash
agentic-agent task list --no-interactive
```

### Automation & Scripting
All commands work in non-TTY environments (pipes, CI/CD):
```bash
echo "abc123" | agentic-agent task claim --no-interactive
```

### Error Handling
Commands exit with:
- `0` on success
- `1` on error

---

## Tips & Tricks

### 🎯 Interactive Mode Tips
- Use **arrow keys** (↑/↓) or **j/k** for navigation
- Press **Enter** to select/confirm
- Press **Esc** to go back or cancel
- Press **q** or **Ctrl+C** to quit
- Press **Space** to toggle selections (where applicable)
- Press **Tab** to switch between tabs (in list views)

### ⚡ Flag Mode Tips
- Add `--no-interactive` for explicit flag mode
- Use `--help` on any command to see all flags
- Combine with shell tools: `agentic-agent task list --no-interactive | grep "in-progress"`
- Perfect for cron jobs and automation scripts

### 🔍 Discovery
- Run any command without arguments to see the interactive interface
- The UI will show you all available options
- Help text appears at the bottom of each screen

---

## Examples by Use Case

### For Beginners
```bash
# Just run commands and follow the prompts
agentic-agent task create
agentic-agent task claim
agentic-agent work
```

### For Scripts
```bash
#!/bin/bash
# Automated workflow
TASK_ID="abc123"
agentic-agent task claim "$TASK_ID"
agentic-agent context generate ./src
agentic-agent task complete "$TASK_ID"
```

### For CI/CD
```yaml
# GitHub Actions example
- name: Validate project
  run: agentic-agent validate --format json --no-interactive

- name: Check skills drift
  run: agentic-agent skills check --no-interactive
```

### For Power Users
```bash
# Quick operations without prompts
agentic-agent task claim abc123
agentic-agent context generate ./src ./tests
agentic-agent task from-template --template feature --title "New API"
agentic-agent work --task abc123 --skip-context-gen
```

---

## Getting Help

### Command Help
```bash
agentic-agent --help              # General help
agentic-agent task --help         # Task commands help
agentic-agent task create --help  # Specific command help
```

### Interactive Help
Run any command without arguments to see the interactive interface with built-in guidance.

---

## Keyboard Shortcuts Reference

### Universal
- **↑/↓** or **j/k**: Navigate up/down
- **Enter**: Select/Confirm
- **Esc**: Back/Cancel
- **q** or **Ctrl+C**: Quit

### Task Lists
- **Tab**: Next tab
- **Shift+Tab**: Previous tab
- **c**: Quick claim (on task)
- **d**: Quick complete (on task)
- **a**: Actions menu

### File Picker
- **Space**: Toggle selection
- **h**: Toggle hidden files
- **Enter**: Select directory

### Multi-Item Editor
- **Enter**: Add item
- **d**: Delete item (when selected)
- **e**: Edit item (when selected)

---

*For more details, see [DUAL_MODE_IMPLEMENTATION_COMPLETE.md](DUAL_MODE_IMPLEMENTATION_COMPLETE.md)*
