# Spec Validation Workflow Example

This example demonstrates the spec validation system that ensures tasks always have proper specification files through hybrid validation with warnings, auto-generation, manual creation, and skip options.

## Overview

The spec validation system helps maintain documentation quality by:

1. **Detecting missing specs** when tasks are created or claimed
2. **Auto-generating specs** from PRDs, interactive prompts, or task metadata
3. **Guiding AI agents** with clear recommendations in CLI output
4. **Allowing flexibility** with skip options for quick iteration

## Scenario: Building a Todo App

Let's walk through creating a todo app feature to see spec validation in action.

### Step 1: Create Tasks Without Following the Pipeline

```bash
# Initialize project
agentic-agent init --name "todo-app-demo"

# Create a task directly (skipping brainstorming → PRD → openspec pipeline)
agentic-agent task create \
  --title "[todo-app] Project setup: scaffold repo, package.json, Vite, Tailwind" \
  --description "Scaffold the repository, package.json, Vite + React template, Tailwind, ESLint, Prettier, and basic npm scripts." \
  --spec-refs "todo-app/proposal.md,todo-app/tasks/01-project-setup.md" \
  --acceptance "`npm install` completes,`npm run dev` starts without errors,Lint and format scripts run"
```

**Output:**
```
Created task TASK-991363-1: [todo-app] Project setup: scaffold repo, package.json, Vite, Tailwind
  Spec refs: todo-app/proposal.md, todo-app/tasks/01-project-setup.md
  Acceptance criteria: 3 items

⚠️  SPEC WARNING: 2 spec(s) referenced but don't exist
  - todo-app/proposal.md
  - todo-app/tasks/01-project-setup.md

💡 TIP: Generate specs with:
   → agentic-agent spec generate TASK-991363-1 --auto
```

### Step 2: Try to Claim the Task

```bash
agentic-agent task claim TASK-991363-1
```

**Output:**
```
Task TASK-991363-1: READY (with warnings)
  [+] input-exists: input file "None" valid (no inputs required)
  [-] spec-completeness: 2 specs missing: todo-app/proposal.md, todo-app/tasks/01-project-setup.md

⚠️  MISSING SPECS DETECTED

The following spec files are referenced but don't exist:
  - todo-app/proposal.md
  - todo-app/tasks/01-project-setup.md

📋 RECOMMENDED ACTIONS FOR AGENTS:

  Option 1: Auto-generate specs (recommended)
  → agentic-agent spec generate TASK-991363-1 --auto

  Option 2: Generate with user interaction
  → agentic-agent spec generate TASK-991363-1 --interactive

  Option 3: Create specs manually
  → agentic-agent spec create todo-app/proposal.md
  → agentic-agent spec create todo-app/tasks/01-project-setup.md

  Option 4: Skip validation and proceed
  → agentic-agent task claim TASK-991363-1 --skip-validation

💡 CONTEXT: This task was likely created outside the recommended workflow.
   For best results, follow: brainstorming → product-wizard → openspec init

🤖 AGENT TIP: Read the user's intent. If they want to proceed quickly,
   use --skip-validation. If quality matters, use --auto generation.

Claimed task TASK-991363-1
```

### Step 3: Auto-Generate Missing Specs

```bash
agentic-agent spec generate TASK-991363-1 --auto
```

**Output:**
```
Analyzing task context...
✓ Context: Will generate from metadata

Generating 2 missing spec(s) for TASK-991363-1...
✓ Created: .agentic/spec/todo-app/proposal.md (492 bytes)
✓ Created: .agentic/spec/todo-app/tasks/01-project-setup.md (318 bytes)

📄 Generated spec summary:
   Content: Generated from metadata

🤖 AGENT: Specs are now complete. You can proceed with:
   → agentic-agent task claim TASK-991363-1
   → agentic-agent context build --task TASK-991363-1
```

**Generated spec example** (`.agentic/spec/todo-app/proposal.md`):
```markdown
# todo-app/proposal.md

**Generated from task:** TASK-991363-1

## Description

Scaffold the repository, package.json, Vite + React template, Tailwind, ESLint, Prettier, and basic npm scripts.

## Acceptance Criteria

- `npm install` completes
- `npm run dev` starts without errors
- Lint and format scripts run
```

### Step 4: Claim Task Again (Now With Complete Specs)

```bash
agentic-agent task claim TASK-991363-1
```

**Output:**
```
Task TASK-991363-1: READY
  [+] input-exists: input file "None" valid (no inputs required)
  [+] spec-completeness: all 2 specs found
Claimed task TASK-991363-1
```

## Different Generation Modes

### Interactive Generation

For tasks where you need to provide detailed requirements:

```bash
agentic-agent spec generate TASK-991363-1 --interactive
```

This will prompt you for:
- Purpose and goals
- Requirements
- Constraints
- Acceptance criteria

### PRD-Based Generation

If you have a PRD file:

```bash
# Create PRD first
cat > .agentic/spec/todo-app-prd.md <<EOF
# Todo App PRD

## Overview
Build a progressive web app for task management...

## Features
- Task creation with titles, descriptions, due dates
- Color coding and priority levels
- Recurring tasks support
...
EOF

# Generate specs from PRD
agentic-agent spec generate TASK-991363-1 --from-prd .agentic/spec/todo-app-prd.md
```

### Manual Creation

For complete control:

```bash
agentic-agent spec create todo-app/proposal.md --template proposal
```

This creates a template you can fill in manually.

## Skip Validation (Quick Iteration)

When you need to move fast and specs aren't critical:

```bash
agentic-agent task claim TASK-991363-1 --skip-validation
```

This bypasses all spec validation checks.

## Configuration

Enable/disable spec validation in `agnostic-agent.yaml`:

```yaml
workflow:
  validate_specs_on_claim: true   # Enable validation before claim
  spec_validation_mode: "warn"    # "warn" | "block" | "silent"
```

**Modes:**
- `warn` (default) — Show warnings but allow claiming
- `block` — Prevent claiming until specs exist
- `silent` — No validation, completely disabled

## Best Practices

### For AI Agents

1. **Read CLI output carefully** — It contains specific guidance on what to do
2. **Choose based on context:**
   - Quick iteration → `--skip-validation`
   - Quality documentation → `--auto` generation
   - Need requirements clarity → `--interactive` generation
3. **Follow recommended pipeline** — brainstorming → product-wizard → openspec init

### For Manual Development

1. **Use the full pipeline** when starting new features
2. **Generate specs early** during task creation
3. **Validate regularly** with `agentic-agent spec validate`
4. **Keep specs up to date** as requirements change

## Bulk Operations

### Generate Specs for All Tasks

```bash
# After creating multiple tasks via openspec init
agentic-agent openspec init "todo-app" --from todo-app-prd.md

# Generate all missing specs at once
agentic-agent spec generate --all
```

### Validate All Tasks

```bash
# Check which tasks have missing specs
agentic-agent spec validate

# Output:
# ✓ TASK-001: All specs present (2/2)
# ✗ TASK-002: Missing specs (0/3)
#   - todo-app/proposal.md
#   - todo-app/tasks/02-ui.md
#   - todo-app/tasks/03-storage.md
# ✓ TASK-003: All specs present (1/1)
```

## Integration with Workflows

### With Autopilot

```bash
# Enable validation before claiming
echo "workflow:
  validate_specs_on_claim: true" >> agnostic-agent.yaml

# Run autopilot — it will validate before each claim
agentic-agent autopilot start --max-iterations 5
```

### With OpenSpec

```bash
# Full pipeline with spec generation
agentic-agent openspec init "feature-name" --from requirements.md
agentic-agent spec generate --all  # Generate missing task specs
agentic-agent openspec import feature-name
agentic-agent task claim TASK-XXX  # Validated automatically
```

### With Tracks

```bash
# Track workflow includes spec validation
agentic-agent track init "User Auth" --type feature
# ... brainstorm and refine spec ...
agentic-agent track activate user-auth --decompose
# Tasks created with spec refs automatically validated
```

## Troubleshooting

**Q: Validation shows missing specs but they exist**
```bash
# Check configured spec directories
cat agnostic-agent.yaml | grep -A 5 "specDirs"

# Verify spec path matches exactly
agentic-agent spec list
```

**Q: Auto-generation creates wrong content**
```bash
# Use interactive mode instead
agentic-agent spec generate TASK-123 --interactive

# Or generate from specific PRD
agentic-agent spec generate TASK-123 --from-prd path/to/prd.md
```

**Q: Want to disable validation temporarily**
```bash
# Use skip flag
agentic-agent task claim TASK-123 --skip-validation

# Or change config
workflow:
  spec_validation_mode: "silent"
```

## Complete Example: Todo App with Full Pipeline

```bash
# 1. Initialize project
mkdir todo-app-demo && cd todo-app-demo
agentic-agent init --name "todo-app"

# 2. Create PRD (manually or with product-wizard skill)
cat > .agentic/spec/todo-app-prd.md <<EOF
# Todo App PRD
## Overview
Progressive web app for task management with offline support.
## Features
- Create/edit/delete tasks
- Color coding and priorities
- Recurring tasks
- Offline sync
EOF

# 3. Generate proposal and tasks from PRD
agentic-agent openspec init "todo-app" --from .agentic/spec/todo-app-prd.md

# Output shows tasks with missing detailed specs:
# ⚠️  VALIDATION WARNING: 5 tasks reference specs that weren't generated
#
# 🤖 AGENT GUIDANCE:
#   Option A) Generate all missing specs from proposal
#   → agentic-agent spec generate --all --from-proposal

# 4. Generate all missing specs
agentic-agent spec generate --all

# 5. Import tasks and start working
agentic-agent openspec import todo-app
agentic-agent task claim TASK-001  # ✓ All specs validated
```

## Summary

The spec validation system ensures documentation quality while maintaining flexibility:

- **Automatic detection** of missing specs at create and claim time
- **Smart generation** from PRDs, prompts, or task metadata
- **Agent-friendly output** with clear recommendations
- **Configurable behavior** from strict to permissive
- **Multiple workflows** supporting different team preferences

This keeps AI agents aligned with requirements while allowing quick iteration when needed.
