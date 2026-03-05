# AGENTS.md - .agents

## Purpose

- Store repo-local agent instructions and reusable skill sources
- Keep local agent behavior versioned with the codebase

## Structure

- `skills/` contains canonical local skills
- Each skill should live at `skills/<skill-name>/SKILL.md`

## Rules

- Keep instructions concise and grounded in real repo files
- Use only `name` and `description` in SKILL.md frontmatter
- Mirror skills needed by Codex auto-discovery under `.codex/skills/`
