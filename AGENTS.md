## Maintenance Rules

- If a directory lacks AGENTS.md, create one following the repo template
- Keep AGENTS.md files small and accurate
- Update AGENTS.md when new patterns are introduced

## Local Skill Rules

- Keep repo-local skills in `.agents/skills/<skill-name>/SKILL.md`
- Use lowercase-hyphen skill names
- Keep compatibility files in `.agents/skills/*.md` pointer-only
- Mirror skills that Codex should auto-discover in `.codex/skills/<skill-name>/SKILL.md`
- When adding or renaming a skill, update `.agents/AGENTS.md` and `.agents/skills/AGENTS.md`
