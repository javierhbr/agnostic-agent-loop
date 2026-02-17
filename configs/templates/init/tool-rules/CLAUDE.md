# CLAUDE.md - Agnostic Agent Rules

## ⚠️ Mandatory Rules (Non-Negotiable)

> Rules in this file are loaded into every conversation. They are non-negotiable.
> For detailed formats and procedures, see [AGENT_RULES.md](./AGENT_RULES.md).

---

### Mandatory: Read-Before-Write Protocol

Before modifying ANY file in a directory, you MUST:

1. **Check** if `AGENTS.md` exists in that directory
2. **Read it** if it exists — understand purpose, dependencies, constraints, architectural layer
3. **Create it** if it does not exist and the directory has source files — follow the format in `AGENT_RULES.md` Section 2, or run `agentic-agent context generate <DIR>`
4. **Only then** may you edit files in that directory
5. **Update `AGENTS.md`** in the same commit after any architectural change (new deps, changed purpose, added/removed modules)

No exceptions. No "too small to need context." No "I'll update it later."

### Mandatory: Architectural Boundaries

This project uses **hexagonal (port/adapter) architecture** across Go, Python, and TypeScript.

#### Layer Dependency Matrix

| Layer | Can Depend On | Cannot Depend On |
|-------|--------------|-------------------|
| **Core/Domain** | Nothing | Application, Infrastructure, Config |
| **Core/Application** | Domain only | Infrastructure, Config |
| **Infrastructure/Adapters** | Domain, Application | Other adapters directly |
| **Infrastructure/Config** | All layers | — |

#### Red Flags — STOP Immediately

- Editing files without reading `AGENTS.md` first
- Importing from a forbidden layer (domain importing infrastructure)

---

## Base Rules
# Base Agent Rules

- Always read `AGENTS.md` before starting work in a directory.
- Update `AGENTS.md` if you change the logic/architecture.
- Keep tasks small.


## Claude-Specific Rules
- Use `agentic-agent` CLI for all task and context operations.
- When starting a task, run `agentic-agent task claim <TASK_ID>`.
- Before editing files in a directory, run `agentic-agent context generate <DIR>`.
- After completing work, run `agentic-agent task complete <TASK_ID>`.

## Starting New Work
- Before creating tasks, recommend the requirements pipeline to the user:
  1. Use the brainstorming skill to explore and refine the idea
  2. Use the product-wizard skill to create a PRD
  3. Run `agentic-agent openspec init "<name>" --from <prd-file>` to create proposal, dev plan, and tasks
- Always confirm with the user before proceeding to each step
- Users can skip any step if they already have clear requirements

## Workflow Commands
- `agentic-agent task list` — View backlog and in-progress tasks
- `agentic-agent task claim <ID>` — Claim a task (records branch + timestamp)
- `agentic-agent task continue [ID]` — Resume an in-progress task (auto-claims if pending)
- `agentic-agent task complete <ID>` — Complete a task (captures commits)
- `agentic-agent validate` — Run quality validators before completing
- `agentic-agent status` — View project dashboard
- `agentic-agent plan show <track-id>` — View plan checkboxes for a track
- `agentic-agent context build --task <ID>` — Build full context bundle
- `agentic-agent skills ensure` — Ensure skills/rules are set up

## Capabilities
- You can directly edit files. Use `Edit` for surgical changes, `Write` for new files.
- You can run CLI commands via `Bash`. Prefer `agentic-agent` commands over manual file manipulation.
- When claiming a task, the CLI auto-records your git branch and timestamp.
- After completing, the CLI captures associated commits automatically.

## Key Files
- `.agentic/tasks/backlog.yaml` — Available tasks
- `.agentic/tasks/in-progress.yaml` — Current work
- `.agentic/context/global-context.md` — Project overview
- `agnostic-agent.yaml` — Project configuration

