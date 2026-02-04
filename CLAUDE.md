# CLAUDE.md - Agnostic Agent Rules

## Base Rules
# Base Agent Rules

- Always read `context.md` before starting work in a directory.
- Update `context.md` if you change the logic/architecture.
- Keep tasks small.


## Claude-Specific Rules
- Use `agentic-agent` CLI for all task and context operations.
- When starting a task, run `agentic-agent task claim <TASK_ID>`.
- Before editing files in a directory, run `agentic-agent context generate <DIR>`.
- After completing work, run `agentic-agent task complete <TASK_ID>`.
