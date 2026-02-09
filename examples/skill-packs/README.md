# Skill Packs: Reusable Agent Skills

Install tool-agnostic skill bundles so that any AI agent tool gets the right instructions for your workflow.

---

## What You'll Learn

- List available skill packs
- Install a pack for a specific AI agent tool
- Install globally (user-level) vs project-level
- Use interactive mode to select pack and tool
- Detect and fix skill drift
- Understand where skill files land for each tool

---

## 0. Setup

```bash
# From the project root
go build -o examples/skill-packs/agentic-agent ./cmd/agentic-agent
cd examples/skill-packs

# Initialize a project
rm -rf .agentic agnostic-agent.yaml
./agentic-agent init --name "skill-demo"
```

---

## 1. List Available Skill Packs

```bash
./agentic-agent skills list
```

Output:

```text
Available Skill Packs

  • tdd         Test-Driven Development (RED/GREEN/REFACTOR) (3 files)
```

Each pack contains multiple skill files that teach the AI agent a specific workflow.

---

## 2. Install a Skill Pack (Flag Mode)

Install the TDD pack for Claude Code:

```bash
./agentic-agent skills install tdd --tool claude-code
```

Output:

```text
Skill Pack Installed

  ✓ Installed "tdd" for claude-code

  • .claude/skills/tdd/red.md
  • .claude/skills/tdd/green.md
  • .claude/skills/tdd/refactor.md
```

The files land in the tool's skill directory:

| Tool        | Skill Directory      |
|-------------|----------------------|
| claude-code | `.claude/skills/`    |
| cursor      | `.cursor/skills/`    |
| gemini      | `.gemini/skills/`    |
| windsurf    | `.windsurf/skills/`  |
| antigravity | `.agent/skills/`     |
| codex       | `.codex/skills/`     |

---

## 3. Install for Multiple Tools

```bash
# Install for Cursor
./agentic-agent skills install tdd --tool cursor

# Install for Gemini
./agentic-agent skills install tdd --tool gemini

# Verify the files
ls .cursor/skills/tdd/
ls .gemini/skills/tdd/
```

---

## 4. Install Globally (User-Level)

Project-level installs put skills in the current directory. Global installs put them in your home directory so every project gets the skills:

```bash
# Global install for Claude Code
./agentic-agent skills install tdd --tool claude-code --global
```

This writes to `~/.claude/skills/tdd/` instead of `.claude/skills/tdd/`.

---

## 5. Interactive Mode

Run without arguments for a guided wizard:

```bash
./agentic-agent skills install
```

The TUI walks you through:

1. Select a skill pack from the list
2. Select the target AI agent tool
3. Confirm installation

---

## 6. Check for Drift

After upgrading the CLI, skill files may be outdated. Check for drift:

```bash
./agentic-agent skills check
```

If everything is current:

```text
Skill Drift Check

  ✓ No drift detected - all skill files are up to date!
```

If files are outdated:

```text
Skill Drift Check

  ✗ Drift detected in 2 file(s):

  • .claude/skills/tdd/red.md
  • .claude/skills/tdd/green.md

  Tip: Use 'agentic-agent skills generate --all' to regenerate skill files
```

Fix drift by regenerating:

```bash
./agentic-agent skills generate --all
```

---

## 7. Generate Tool-Specific Config Files

Beyond skill packs, generate PRD and converter skills for specific tools:

```bash
# Claude Code: generates .claude/skills/prd.md and .claude/skills/ralph-converter.md
./agentic-agent skills generate-claude-skills

# Gemini CLI: generates .gemini/commands/prd/gen.toml and .gemini/commands/ralph/convert.toml
./agentic-agent skills generate-gemini-skills

# Generic: select tool interactively
./agentic-agent skills generate
```

---

## Quick Reference

| Action                | Command                                          |
|-----------------------|--------------------------------------------------|
| List packs            | `skills list`                                    |
| Install (flag mode)   | `skills install tdd --tool claude-code`          |
| Install (global)      | `skills install tdd --tool cursor --global`      |
| Install (interactive) | `skills install`                                 |
| Check drift           | `skills check`                                   |
| Regenerate all        | `skills generate --all`                          |
| Claude Code skills    | `skills generate-claude-skills`                  |
| Gemini skills         | `skills generate-gemini-skills`                  |
