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
