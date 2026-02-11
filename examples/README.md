# Examples

Step-by-step demonstrations of the Agentic Agent framework.

## Available Examples

### `test-sandbox/` - End-to-End Walkthrough

Build a project from scratch: idea elaboration, project init, PRD generation, task decomposition, claim-to-complete workflow, and tool switching.

**Start here** if you're new. See [SAMPLE.README.md](test-sandbox/SAMPLE.README.md).

### `track-workflow/` - Tracks: Idea to Implementation

Walk through the full track lifecycle: brainstorm with an AI agent, refine the spec, activate to generate a plan and tasks, work through tasks, and complete the track.

Demonstrates: `track init`, `track refine`, `track activate --decompose`, `plan show`, `plan next`, `plan mark`, `status`.

See [README.md](track-workflow/README.md).

### `skill-packs/` - Skill Pack Installation and Usage

Install reusable skill bundles for any AI agent tool. Shows project-level and global installation, listing packs, and verifying drift.

Demonstrates: `skills install`, `skills list`, `skills check`, multi-tool support.

See [README.md](skill-packs/README.md).

### `tdd/` - Test-Driven Development Workflow

Use `work --follow-tdd` to decompose a task into RED/GREEN/REFACTOR sub-tasks. The TDD skill pack provides phase-specific instructions to the AI agent.

See [README.md](tdd/README.md).

### `multi-agent-workflow/` - Multi-Agent Tool Switching

Bounce between Claude Code CLI, Claude Code VSCode, GitHub Copilot, and Antigravity IDE with Gemini on the same project. Demonstrates non-linear back-and-forth switching, shared state, and cross-tool bug discovery.

See [MULTI_AGENT_USE_CASE.md](multi-agent-workflow/MULTI_AGENT_USE_CASE.md).

### `spec-driven-workflow/` - Spec Kit and OpenSpec Integration

Multi-directory spec resolution with Spec Kit, OpenSpec, and native specs. Shows how tasks reference specs across directories and how autopilot processes them sequentially.

See [README.md](spec-driven-workflow/README.md).

### `agent-aware-skills/` - Agent Detection, Setup, and Per-Agent Rules

Automatically detect which AI agent is running (Claude Code, Cursor, Gemini, etc.), ensure its skills and rules are installed, and tailor instructions per tool. Covers per-agent config overrides, custom rules files, scoped drift checks, auto-ensure in init/run/autopilot, task-level `skill_refs`, and the `simplify` command.

Demonstrates: `skills ensure`, `--agent` flag, `AGENTIC_AGENT` env var, `.agentic/agent-rules/`, per-agent `skill_packs` and `extra_rules`, task-level `skill_refs`, `simplify`.

See [README.md](agent-aware-skills/README.md).

---

## Running Any Example

### Initialize a project

```bash
cd examples/<example-dir>
agentic-agent init --name "My Project"
```

### Work with tracks

```bash
# Start a track with brainstorming scaffolding
agentic-agent track init "My Feature" --type feature

# Check spec completeness
agentic-agent track refine my-feature

# Activate: generate plan + decompose into tasks
agentic-agent track activate my-feature --decompose
```

### Work with tasks

```bash
agentic-agent task list                 # List all tasks
agentic-agent task claim TASK-001       # Claim with readiness checks
agentic-agent task complete TASK-001    # Mark done
```

### Work with plans

```bash
agentic-agent plan show --track my-feature   # View plan progress
agentic-agent plan next --track my-feature   # See next pending step
agentic-agent plan mark plan.md 12 done      # Mark step done by line
```

### Generate context

```bash
agentic-agent context generate internal/auth   # Directory context
agentic-agent context build --task TASK-001    # Full context bundle
```

### Work with specs

```bash
agentic-agent spec list                     # All specs across directories
agentic-agent spec resolve auth-spec.md     # Resolve and print
```

### Install skill packs

```bash
agentic-agent skills list                             # Available packs
agentic-agent skills install tdd --tool claude-code   # Install for a tool
agentic-agent skills check                            # Detect drift
```

### Ensure agent skills

```bash
agentic-agent skills ensure                           # Auto-detect agent
agentic-agent skills ensure --agent claude-code       # Explicit agent
agentic-agent skills ensure --all                     # All detected agents
agentic-agent skills check --agent cursor             # Scoped drift check
```

### Run code simplification review

```bash
agentic-agent simplify internal/auth              # Review specific directories
agentic-agent simplify --task TASK-001            # Review task scope directories
agentic-agent simplify . --format json            # JSON output
agentic-agent simplify . --output review.yaml     # Write to file
```

### Run autopilot

```bash
agentic-agent autopilot start --dry-run          # Preview
agentic-agent autopilot start --max-iterations 5  # Process tasks
```

### Check project status

```bash
agentic-agent status                # Dashboard with progress bar
agentic-agent status --format json  # Machine-readable output
```

---

## Example Use Cases

- **[test-sandbox](test-sandbox/SAMPLE.README.md)** — Full workflow, tool switching, parallel agents
- **[track-workflow](track-workflow/README.md)** — Brainstorming, spec refinement, plan generation
- **[skill-packs](skill-packs/README.md)** — Multi-tool skill installation, drift detection
- **[tdd](tdd/README.md)** — RED/GREEN/REFACTOR decomposition
- **[multi-agent-workflow](multi-agent-workflow/MULTI_AGENT_USE_CASE.md)** — 4 tools, 6 phases, cross-tool bug fix
- **[spec-driven-workflow](spec-driven-workflow/README.md)** — Spec Kit, OpenSpec, autopilot mode
- **[agent-aware-skills](agent-aware-skills/README.md)** — Agent detection, per-agent rules, `skills ensure`, `skill_refs`, `simplify`

---

## Configuration Reference

`agentic-agent init` generates `agnostic-agent.yaml` in the project root. Here is a fully annotated sample:

```yaml
# ── Project metadata ────────────────────────────────────────────────
project:
  name: "my-project"           # Project name shown in status dashboard
  version: 0.1.0               # Semantic version (informational)
  roots:                        # Source roots to scan for context
    - .

# ── Agent defaults ──────────────────────────────────────────────────
agents:
  defaults:
    max_tokens: 4000            # Token budget for context bundles
    model: claude-3-5-sonnet-20241022
  overrides:                    # Per-tool overrides (optional)
    - name: cursor
      max_tokens: 8000
    - name: claude-code
      max_tokens: 8000
      skill_packs:              # Packs auto-installed by `skills ensure`
        - tdd
      extra_rules:              # Additional lines injected into rules file
        - "Run tests before completing tasks"
      auto_setup: true          # Generate rules during `init`

# ── Paths ───────────────────────────────────────────────────────────
# All paths are relative to the project root.
paths:
  # Spec resolution: searched in order, first match wins.
  specDirs:
    - .specify/specs            # Spec Kit (if using Spec Kit)
    - openspec/specs            # OpenSpec (if using OpenSpec)
    - .agentic/spec             # Native agentic specs (default)

  contextDirs:
    - .agentic/context          # Global context files

  trackDir: .agentic/tracks     # Track work units (spec + plan + tasks)
  prdOutputPath: .agentic/tasks/       # Where PRD converter writes tasks (task YAML supports skill_refs)
  progressTextPath: .agentic/progress.txt
  progressYAMLPath: .agentic/progress.yaml
  archiveDir: .agentic/archive/        # Archived tracks and tasks

# ── Workflow ────────────────────────────────────────────────────────
workflow:
  validators:                   # Validation rules run by `agentic-agent validate`
    - context-check             # Verify context.md exists in scope dirs
    - task-scope                # Enforce max 5 files / 2 dirs per task
    - browser-verification      # (optional) Browser-based checks
```

### Minimal config

If you only need the basics, most fields have sensible defaults:

```yaml
project:
  name: "my-project"

paths:
  specDirs:
    - .agentic/spec

workflow:
  validators:
    - context-check
    - task-scope
```

Omitted fields use these defaults:

| Field              | Default                       |
|--------------------|-------------------------------|
| `project.version`  | `0.1.0`                       |
| `project.roots`    | `[.]`                         |
| `agents.defaults`  | 4000 tokens, Sonnet           |
| `paths.specDirs`   | `[.agentic/spec]`             |
| `paths.contextDirs`| `[.agentic/context]`          |
| `paths.trackDir`   | `.agentic/tracks`             |
| `paths.archiveDir` | `.agentic/archive/`           |

---

## Creating Your Own Example

```bash
mkdir examples/my-example && cd examples/my-example
agentic-agent init --name "My Example"
```

Structure:

```text
examples/my-example/
|-- README.md                    # Walkthrough
|-- agnostic-agent.yaml          # Configuration
|-- .agentic/
|   |-- tasks/                   # backlog, in-progress, done
|   |-- spec/                    # Specification files
|   |-- context/                 # global-context, rolling-summary
|   |-- tracks/                  # Track work units
|   +-- agent-rules/             # base.md
+-- src/                         # Your code
```

## Related Documentation

- [Main README](../README.md) - Project overview and CLI reference
- [Spec-Driven Development Guide](../docs/SPEC_DRIVEN_DEVELOPMENT.md) - Spec Kit and OpenSpec workflows
- [CLI Tutorial](../docs/guide/CLI_TUTORIAL.md) - Command-line usage
- [BDD Guide](../docs/bdd/BDD_GUIDE.md) - Testing with Gherkin
