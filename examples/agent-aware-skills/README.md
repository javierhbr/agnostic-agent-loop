# Agent-Aware Skills: Detection, Setup, and Per-Agent Rules

Automatically detect which AI agent is running, ensure its skills and rules are installed, and tailor instructions per tool.

---

## What You'll Learn

- Detect the active agent via flag, environment variable, or filesystem
- Use `skills ensure` to idempotently set up rules, skill packs, and tool-specific files
- Configure per-agent overrides (token budgets, skill packs, extra rules)
- Write custom per-agent rules in `.agentic/agent-rules/`
- Scope drift checks to a single agent with `skills check --agent`
- See how `init`, `run`, and autopilot auto-ensure skills
- Use `skill_refs` on tasks for targeted skill inclusion
- Run `simplify` to generate a code simplification review bundle

---

## 0. Setup

```bash
# From the project root
go build -o examples/agent-aware-skills/agentic-agent ./cmd/agentic-agent
cd examples/agent-aware-skills
```

This example ships with a pre-configured `agnostic-agent.yaml` and `.agentic/` directory. No `init` needed.

---

## 1. Agent Detection

The CLI detects which agent is running using three strategies (in priority order):

| Priority | Strategy | Example |
|----------|----------|---------|
| 1 | `--agent` flag | `--agent claude-code` |
| 2 | Environment variable | `AGENTIC_AGENT=cursor` or `CLAUDE=1` |
| 3 | Filesystem heuristic | `.claude/` directory or `CLAUDE.md` file exists |

### Explicit flag

```bash
./agentic-agent skills ensure --agent claude-code
```

### Environment variable

```bash
# Generic (works for any agent)
AGENTIC_AGENT=cursor ./agentic-agent skills ensure

# Agent-specific variables
CLAUDE=1 ./agentic-agent skills ensure
CURSOR_SESSION=1 ./agentic-agent skills check
GEMINI_CLI=1 ./agentic-agent skills ensure
```

### Filesystem heuristic

When no flag or env var is set, the CLI checks the project root for agent-specific paths:

| Path | Detected Agent |
|------|---------------|
| `.claude/` or `CLAUDE.md` | claude-code |
| `.cursor/` | cursor |
| `.gemini/` | gemini |
| `.windsurf/` | windsurf |
| `.codex/` | codex |
| `.agent/` | antigravity |

```bash
# Create a .cursor/ directory and detection kicks in automatically
mkdir .cursor
./agentic-agent skills check
# → Checks drift for cursor only
```

---

## 2. Ensure Skills (Idempotent Setup)

`skills ensure` is the single command that makes sure an agent has everything it needs. It is safe to run repeatedly.

### What it does

1. Generates the agent's rules file if missing (e.g., `CLAUDE.md`, `.cursor/rules/agnostic-agent.mdc`)
2. Fixes drift if the rules file is outdated
3. Installs configured skill packs from `agnostic-agent.yaml`
4. Generates tool-specific skill files (Claude PRD/Ralph, Gemini slash commands)

### For a single agent

```bash
./agentic-agent skills ensure --agent claude-code
```

Output:

```text
Ensuring skills for claude-code

Generated rules: CLAUDE.md
Installed packs: tdd
```

### For all detected agents

```bash
./agentic-agent skills ensure --all
```

### Run it again (idempotent)

```bash
./agentic-agent skills ensure --agent claude-code
```

Output:

```text
Ensuring skills for claude-code

Already up to date.
```

### Interactive mode

```bash
./agentic-agent skills ensure
```

When no agent is detected, the TUI prompts you to select one:

```text
Ensure Agent Skills

  ❯ claude-code    Ensure skills for claude-code
    cursor          Ensure skills for cursor
    gemini          Ensure skills for gemini
    windsurf        Ensure skills for windsurf
    codex           Ensure skills for codex

  ↑/↓ navigate • Enter select • Esc cancel
```

---

## 3. Per-Agent Configuration

The `agnostic-agent.yaml` in this example shows per-agent overrides:

```yaml
agents:
  defaults:
    max_tokens: 4000
    model: claude-3-5-sonnet-20241022
  overrides:
    - name: claude-code
      max_tokens: 8000
      skill_packs:
        - tdd
        - api-docs
        - code-simplification
      extra_rules:
        - "Always run `go test ./...` before completing a task"
        - "Use table-driven tests for all new test functions"
      auto_setup: true

    - name: cursor
      max_tokens: 12000
      skill_packs:
        - tdd
        - dev-plans
        - diataxis
      extra_rules:
        - "Prefer small, focused commits"

    - name: gemini
      skill_packs:
        - tdd
        - extract-wisdom

    - name: windsurf
      auto_setup: true
```

### Fields

| Field | Purpose |
|-------|---------|
| `max_tokens` | Token budget for context bundles |
| `skill_packs` | Packs to install via `skills ensure` |
| `extra_rules` | Additional lines injected into the agent's rules file |
| `auto_setup` | Automatically generate rules during `init` |

When `skills ensure` runs for `claude-code`, it:
1. Generates `CLAUDE.md` with base rules + claude-code.md custom rules + extra_rules from config
2. Installs the `tdd` skill pack to `.claude/skills/tdd/`
3. Generates `.claude/skills/prd.md` and `.claude/skills/ralph-converter.md`

---

## 4. Custom Per-Agent Rules

Each agent can have a dedicated rules file in `.agentic/agent-rules/<tool>.md`. These are injected into the generated rules file via the `{{ .AgentRules }}` template variable.

### File convention

```text
.agentic/agent-rules/
├── base.md           # Shared rules for all agents
├── claude-code.md    # Claude Code-specific rules
├── cursor.md         # Cursor-specific rules
└── gemini.md         # Gemini-specific rules
```

### How it works

When generating `CLAUDE.md`, the template includes:

```
## Base Rules
{{ .BaseRules }}

## Claude-Specific Rules
{{ .AgentRules }}
```

Where `AgentRules` is the concatenation of:
1. Content from `.agentic/agent-rules/claude-code.md`
2. Lines from `extra_rules` in `agnostic-agent.yaml`

### Example

This project includes `.agentic/agent-rules/claude-code.md`:

```markdown
- When starting a task, run `agentic-agent task claim <TASK_ID>`.
- Before editing files in a directory, run `agentic-agent context generate <DIR>`.
- After completing work, run `agentic-agent task complete <TASK_ID>`.
- Use the Bash tool for `go test` and `go build` commands.
- Prefer editing existing files over creating new ones.
- When writing tests, use table-driven test patterns.
```

Generate and inspect the result:

```bash
./agentic-agent skills ensure --agent claude-code
cat CLAUDE.md
```

The generated `CLAUDE.md` will contain the base rules plus the claude-code-specific rules above.

---

## 5. Scoped Drift Checks

After upgrading the CLI, rules files may become outdated. Check drift for a specific agent:

```bash
# Check only claude-code
./agentic-agent skills check --agent claude-code
```

Output (no drift):

```text
Skill Drift Check

  ✓ No drift detected - all skill files are up to date!
```

Output (drift found):

```text
Skill Drift Check

  ✗ Drift detected in 1 file(s):

  • CLAUDE.md

  Tip: Use 'agentic-agent skills generate --all' to regenerate skill files
```

Fix drift automatically:

```bash
./agentic-agent skills ensure --agent claude-code
# → Fixed drift: CLAUDE.md
```

Without `--agent`, `skills check` checks all registered tools (which may report "Missing" for tools you don't use).

---

## 6. Auto-Ensure in Init, Run, and Autopilot

The skills system integrates with three key workflows.

### During `init`

The interactive wizard includes an "Agent Tools" step where you select which tools to set up:

```bash
./agentic-agent init --name "my-project"
```

The wizard walks through:
1. Project name
2. Workflow preferences
3. **Agent tools** (select claude-code, cursor, gemini, etc.)
4. Preview and confirm

Or use the flag:

```bash
./agentic-agent init --name "my-project" --agent claude-code
```

### During `run`

Before orchestrating a task, the `run` command checks if the active agent's rules file exists:

```bash
./agentic-agent run --task TASK-1001 --agent claude-code
```

If `CLAUDE.md` is missing, it auto-generates:

```text
Skills not found for claude-code. Generating CLAUDE.md...
Generated rules: CLAUDE.md
```

### During autopilot

The autopilot loop calls `skills.Ensure()` at startup, so the agent always has current rules and packs before processing tasks.

---

## 7. Task-Level Skill Refs

Tasks can declare `skill_refs` to include specific skill packs in their context bundle, rather than loading all installed packs. This works like `spec_refs` but for skill packs.

### Task YAML

```yaml
tasks:
  - id: "TASK-1001"
    title: "Implement user authentication endpoint"
    spec_refs:
      - auth-spec.md
    skill_refs:          # Only these skills are included in the context bundle
      - tdd
      - api-docs
    scope:
      - "internal/auth"

  - id: "TASK-1002"
    title: "Add input validation middleware"
    skill_refs:
      - code-simplification    # Only simplification principles for this task
    scope:
      - "internal/middleware"
```

### Targeted vs default mode

When `skill_refs` is present on a task, the context bundle includes **only** those skill packs (targeted mode). When absent, **all** installed skill packs are included (existing behavior).

The agent's generated rules file (e.g., `CLAUDE.md`) is always included regardless.

### Resolution order

Each skill ref is resolved in order:

1. Agent's installed skill directory (e.g., `.claude/skills/code-simplification/SKILL.md`)
2. Any other installed tool's directory
3. Embedded pack content (compiled into the CLI binary)

This means skill refs always resolve, even for packs not explicitly installed.

### Build a context bundle with skill refs

```bash
AGENTIC_AGENT=claude-code ./agentic-agent context build --task TASK-1001
```

The bundle's `skill_instructions` field will contain only the TDD and API docs skill content (from `skill_refs`), plus the agent's base rules.

---

## 8. Simplify Command

The `simplify` command generates a focused context bundle for code simplification review using the `code-simplification` skill pack.

### Simplify specific directories

```bash
./agentic-agent simplify internal/auth internal/middleware
```

### Simplify a task's scope

```bash
./agentic-agent simplify --task TASK-1002
```

This loads the task's `scope` directories and builds a simplification bundle.

### Output formats

```bash
# Default (toon format)
./agentic-agent simplify internal/auth

# JSON output
./agentic-agent simplify internal/auth --format json

# YAML output
./agentic-agent simplify internal/auth --format yaml

# Write to file
./agentic-agent simplify internal/auth --output review.yaml --format yaml
```

### What the bundle contains

| Field                | Content                                                      |
|----------------------|--------------------------------------------------------------|
| `skill_instructions` | The code-simplification SKILL.md content                     |
| `directories`        | Generated context for each target directory                  |
| `target_files`       | Source files found in the target directories                 |
| `tech_stack`         | Tech stack info (if `.agentic/context/tech-stack.md` exists) |

---

## 9. Context Bundles with Skill Instructions

When building a context bundle for a task, the CLI includes agent-specific skill instructions:

```bash
AGENTIC_AGENT=claude-code ./agentic-agent context build --task TASK-1001
```

The bundle's `skill_instructions` field contains:
1. The agent's generated rules file (e.g., `CLAUDE.md` content)
2. Skill pack content — either from `skill_refs` (if set) or all installed `SKILL.md` files

This means the AI agent receives its own rules and installed skill instructions as part of the context payload, without needing to read separate files.

---

## 10. Multi-Agent Project

A single project can serve multiple agents. Each gets its own rules file and skill directory:

```bash
# Set up for all agents in one shot
./agentic-agent skills ensure --all
```

Result:

```text
examples/agent-aware-skills/
├── CLAUDE.md                              # Claude Code rules
├── .cursor/rules/agnostic-agent.mdc       # Cursor rules
├── .gemini/GEMINI.md                      # Gemini rules
├── .windsurf/rules/agnostic-agent.md      # Windsurf rules
├── .codex/CODEX.md                        # Codex rules
├── .claude/skills/tdd/                    # TDD pack for Claude
├── .cursor/skills/tdd/                    # TDD pack for Cursor
├── .gemini/skills/tdd/                    # TDD pack for Gemini
└── agnostic-agent.yaml                    # Shared config
```

Each agent reads only its own files. The shared `agnostic-agent.yaml` drives what gets generated for each.

---

## Quick Reference

| Action | Command |
|--------|---------|
| Ensure (auto-detect) | `skills ensure` |
| Ensure (explicit) | `skills ensure --agent claude-code` |
| Ensure (all agents) | `skills ensure --all` |
| Check drift (scoped) | `skills check --agent cursor` |
| Check drift (all) | `skills check` |
| Init with agent | `init --name "proj" --agent claude-code` |
| Run with agent | `run --task TASK-001 --agent gemini` |
| Set via env var | `AGENTIC_AGENT=cursor agentic-agent run --task TASK-001` |
| Simplify directories | `simplify internal/auth internal/middleware` |
| Simplify task scope | `simplify --task TASK-1002` |
| Simplify to file | `simplify internal/auth --output review.yaml --format yaml` |

## Detection Priority

```text
--agent flag  ──→  AGENTIC_AGENT env  ──→  Agent-specific env  ──→  Filesystem
  (highest)                                (CLAUDE, CURSOR_SESSION)    (lowest)
```

## Agent Rules Resolution

```text
.agentic/agent-rules/base.md         ──→  {{ .BaseRules }} in template
.agentic/agent-rules/<tool>.md        ──→  {{ .AgentRules }} in template
agnostic-agent.yaml extra_rules       ──→  appended to {{ .AgentRules }}
```
