---
name: code-documentation-skill
description: "Generate and maintain hierarchical AGENTS.md files from real repository structure. Use when asked to create or update AGENTS.md files, bootstrap directory documentation, audit AGENTS.md coverage, or document code boundaries, entrypoints, and run/test commands."
---

# Code Documentation Skill

Create high-signal AGENTS.md files that help future agent runs be faster and safer.

## Non-negotiable rules

- Never invent files, commands, ownership, or architecture.
- Mention only paths that exist.
- If uncertain, write `Unknown - verify` and point to the file to confirm.
- Update existing AGENTS.md files in place and preserve accurate content.
- Skip generated or vendor directories.

## Working process

1. Discover repository layout
- List top-level directories.
- Identify primary languages and frameworks from manifests (`go.mod`, `package.json`, `pyproject.toml`, etc.).

2. Select targets
- Include source, domain, service, library, infrastructure, and API directories.
- Exclude non-source paths such as `.git`, `node_modules`, `dist`, `build`, `.next`, `.turbo`, `vendor`, `target`, `coverage`, `__pycache__`.

3. Build or update root `AGENTS.md`
- Add a short project overview.
- Add accurate commands to run, test, lint, and build.
- Add a concise architecture map by directory.
- Add global rules and quality gates.
- Explain local AGENTS inheritance (root to leaf).

4. Build or update local `AGENTS.md` per target directory
- Use this structure:
  - `## Purpose`
  - `## Key Entrypoints`
  - `## File Map`
  - `## Important Flows` (if known)
  - `## Local Conventions / Invariants`
  - `## How to Run / Test`
  - `## Do / Don’t`
  - `## Common Tasks` (only tasks that match real directory behavior)

5. Run accuracy checks
- Verify every referenced path exists.
- Verify every command exists (scripts, Make targets, task runners).
- Mark unverifiable items as `Unknown - verify`.
- Keep each AGENTS.md concise.

6. Report output
- List:
  - `Created`
  - `Updated`
  - `Skipped` (with reason)
- Include short next steps only when information is missing.

## Output style

- Prefer bullets and short sentences.
- Keep each file focused on actionable guidance.
- Link to existing docs instead of duplicating large sections.
