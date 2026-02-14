# OPENCODE.md - OpenCode-Specific Instructions

> **📋 For comprehensive agent-agnostic rules, see [AGENT_RULES.md](./AGENT_RULES.md)**
> 
> This file contains OpenCode-specific commands and workflow integration.

---

## Quick Reference: Core Rules

⚠️ **Read [AGENT_RULES.md](./AGENT_RULES.md) for complete context management, architectural boundaries, and quality rules.**

Essential reminders:
- Always read `AGENTS.md` before editing files in a directory
- Update `AGENTS.md` when you change architecture or dependencies
- Never violate architectural boundaries defined in context files
- Keep tasks small and focused

---

## OpenCode-Specific Commands

Use `agentic-agent` CLI for all task and context operations.

### Task Management

```bash
# View available tasks
agentic-agent task list

# Claim a task (records branch + timestamp)
agentic-agent task claim <TASK_ID>

# Resume an in-progress task (auto-claims if pending)
agentic-agent task continue [ID]

# Complete a task (captures commits)
agentic-agent task complete <ID>
```

### Context Operations

```bash
# Generate context for a directory
agentic-agent context generate <DIR>

# Scan for missing context files
agentic-agent context scan

# Build full context bundle for a task
agentic-agent context build --task <ID>
```

### Validation & Status

```bash
# Run quality validators before completing
agentic-agent validate

# View project dashboard
agentic-agent status

# View plan checkboxes for a track
agentic-agent plan show <track-id>
```

### Skills Management

```bash
# Ensure skills/rules are set up
agentic-agent skills ensure
```

---

## Starting New Work

Before creating tasks, recommend the requirements pipeline to the user:

1. Use the brainstorming skill to explore and refine the idea
2. Use the product-wizard skill to create a PRD
3. Run `agentic-agent openspec init "<name>" --from <prd-file>` to create proposal, dev plan, and tasks
4. Always confirm with the user before proceeding to each step

Users can skip any step if they already have clear requirements.

---

## Capabilities

- You can directly edit files
- You can run CLI commands
- When claiming a task, the CLI auto-records your git branch and timestamp
- After completing, the CLI captures associated commits automatically

---

## Key Files

| File | Purpose |
|------|---------|
| `.agentic/tasks/backlog.yaml` | Available tasks |
| `.agentic/tasks/in-progress.yaml` | Current work |
| `.agentic/context/global-context.md` | Project overview |
| `agnostic-agent.yaml` | Project configuration |
| `AGENTS.md` (per directory) | Architectural rules and constraints |

---

For complete rules and guidelines, see [AGENT_RULES.md](./AGENT_RULES.md).
