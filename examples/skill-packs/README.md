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
- Reference skills per-task with `skill_refs`
- Run targeted code simplification with `simplify`

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

  • atdd                   Acceptance Test-Driven Development from openspec task criteria (1 file)
  • tdd                    Test-Driven Development with red-green-refactor workflow (3 files)
  • api-docs               Generate comprehensive API documentation from code (1 file)
  • code-simplification    Review and refactor code for simplicity and maintainability (1 file)
  • dev-plans              Create structured development plans with phased task breakdowns (1 file)
  • diataxis               Write documentation using the Diataxis framework (3 files)
  • extract-wisdom         Extract insights and actionable takeaways from text sources (1 file)
  • openspec               Spec-driven development from requirements files (1 file)
  • product-wizard         Generate robust, production-grade PRDs (5 files)
  • run-with-ralph         Execute openspec tasks using Ralph Wiggum iterative loops (1 file)
```

Each pack contains skill files that teach the AI agent a specific workflow or capability.

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
| copilot     | `.github/skills/`    |
| opencode    | `.opencode/skills/`  |

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

## 8. Task-Level Skill Refs

Instead of including all installed packs in every context bundle, tasks can declare which skill packs they need:

```yaml
tasks:
  - id: "TASK-001"
    title: "Refactor auth middleware"
    skill_refs:
      - code-simplification
      - tdd
    scope:
      - "internal/auth"
```

When a task has `skill_refs`, the context bundle includes **only** those packs. Without `skill_refs`, all installed packs are included (backwards compatible).

Skill refs resolve through a 3-tier fallback:

1. Agent's installed directory (e.g., `.claude/skills/tdd/SKILL.md`)
2. Any other tool's installed directory
3. Embedded content (compiled into the binary)

This means `skill_refs` always resolve, even if the pack isn't installed locally.

See [agent-aware-skills/README.md](../agent-aware-skills/README.md) for detailed examples.

---

## 9. Simplify Command

The `simplify` command generates a focused review bundle using the `code-simplification` skill pack:

```bash
# Review specific directories
./agentic-agent simplify internal/auth

# Review a task's scope directories
./agentic-agent simplify --task TASK-001

# Output as JSON or YAML
./agentic-agent simplify internal/auth --format json
./agentic-agent simplify internal/auth --output review.yaml --format yaml
```

The bundle contains the code-simplification skill instructions, directory context, source file listings, and tech stack info.

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
| Simplify directories  | `simplify internal/auth`                         |
| Simplify task scope   | `simplify --task TASK-001`                       |
