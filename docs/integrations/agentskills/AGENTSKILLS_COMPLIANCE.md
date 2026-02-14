# agentskills.io Compliance Guide

## Overview

The Ralph PDR integration skills are designed to be compatible with the [agentskills.io](https://agentskills.io) open standard, enabling portability across Claude Code, Cursor, Windsurf, and other AI coding assistants.

## Current Compliance Status

### ✅ Already Compliant

Our generated skills follow the core agentskills.io format:

```markdown
---
name: skill-name
description: "When to use this skill..."
---

# Skill Title

Instructions and workflow...
```

**Generated Skills:**
- `.claude/skills/prd.md` - PRD Generator
- `.claude/skills/ralph-converter.md` - PRD to YAML Converter

### 🔄 Enhanced Compliance (Optional)

For full agentskills.io compatibility with advanced features, skills can be organized into directories:

## Recommended Directory Structure

```
.claude/skills/
├── prd-generator/
│   ├── SKILL.md                    # Main skill file
│   ├── examples/
│   │   └── sample-prd.md          # Example output
│   └── templates/
│       └── story-template.md      # User story template
│
├── ralph-converter/
│   ├── SKILL.md                    # Main skill file
│   ├── scripts/
│   │   └── validate-stories.js    # Optional validator
│   └── examples/
│       ├── input-prd.md           # Example input
│       └── output-tasks.yaml      # Example output
│
└── context-manager/
    ├── SKILL.md                    # Context.md enforcement
    ├── scripts/
    │   └── validate-context.py    # Architectural validator
    └── references/
        └── architecture.md         # Hexagonal rules
```

## Enhanced SKILL.md Format

### Extended Frontmatter (Optional)

```yaml
---
name: prd-generator
description: "Generate Product Requirements Documents using Ralph PDR methodology"
version: 1.0.0
author: agnostic-agent
triggers:
  - create a prd
  - write prd for
  - plan this feature
  - spec out
dependencies:
  - context-manager  # Suggested skill to use together
output_location: ".agentic/tasks/"
compatible_tools:
  - claude-code
  - cursor
  - windsurf
  - copilot
---
```

### Skill Body Sections

```markdown
# PRD Generator Skill

## 🎯 Purpose
[One-sentence description]

## 🚀 Activation
This skill activates when you:
- Say "create a PRD"
- Start planning a new feature
- Need to spec out requirements

## 📋 Workflow
1. Step one
2. Step two
3. Step three

## 🔧 Configuration
- Output path: Configurable via agnostic-agent.yaml
- Default: `.agentic/tasks/`

## 📚 Examples
See `examples/` directory for sample PRDs.

## 🤝 Integration
Works with:
- context-manager skill (reads context.md)
- ralph-converter skill (converts PRD → tasks)

## 📖 References
- [Ralph PDR Methodology](../docs/RALPH_PDR_WORKFLOW.md)
- [Task Management](../docs/TASK_MANAGEMENT.md)
```

## Creating agentskills.io-Compliant Skills

### Step 1: Generate Base Skills

```bash
# Generate skills in flat format (current)
agentic-agent skills generate-claude-skills

# Output:
# .claude/skills/prd.md
# .claude/skills/ralph-converter.md
```

### Step 2: Enhance to Directory Format (Optional)

Manually restructure for advanced features:

```bash
# Create directory structure
mkdir -p .claude/skills/prd-generator/{examples,templates}
mkdir -p .claude/skills/ralph-converter/{examples,scripts}

# Move skill files
mv .claude/skills/prd.md .claude/skills/prd-generator/SKILL.md
mv .claude/skills/ralph-converter.md .claude/skills/ralph-converter/SKILL.md

# Add examples and scripts
cp docs/examples/sample-prd.md .claude/skills/prd-generator/examples/
```

### Step 3: Add Supporting Files

**Example: `.claude/skills/prd-generator/examples/sample-prd.md`**

```markdown
# PRD: Task Priority System

## Introduction
Add priority levels to tasks...

## User Stories

### US-001: Add priority field to database
**Description:** As a developer, I need to store task priority.

**Acceptance Criteria:**
- [ ] Add priority column: 'high' | 'medium' | 'low'
- [ ] Typecheck passes
```

**Example: `.claude/skills/ralph-converter/scripts/validate-stories.js`**

```javascript
// Optional: Validate story structure
function validateStory(story) {
  if (!story.acceptanceCriteria.includes('Typecheck passes')) {
    console.warn(`Story ${story.id} missing typecheck criteria`);
  }
}
```

## Integration with Context Manager Skill

Following the agentskills.io pattern from your example, here's how to integrate:

### Context Manager Skill

**File: `.claude/skills/context-manager/SKILL.md`**

```markdown
---
name: context-manager
description: "Manages and enforces 'context.md' files across the repository. Activates when creating new directories, modifying existing logic, or ensuring adherence to Hexagonal Architecture boundaries."
triggers:
  - creating new directory
  - modifying existing code
  - adding new module
compatible_with:
  - prd-generator
  - ralph-converter
---

# Context Management Skill

You are responsible for maintaining the "living documentation" through `AGENTS.md` files.

## 🚦 Mandatory Workflow: Read-Before-Write

Before modifying any file:
1. Check if `AGENTS.md` exists in that directory
2. If exists, read it for `MUST DO` and `YOU CANNOT DO` constraints
3. If missing and directory contains logic, trigger Generation Workflow

## 🏗️ Hexagonal Architecture Boundaries

- **Core/Domain:** No dependencies. Pure business logic only.
- **Core/Application:** Can only depend on Domain.
- **Infrastructure/Adapters:** Entry/Exit points.
- **Infrastructure/Config:** The "Glue" for DI.

## 🔄 Generation Workflow

If `AGENTS.md` is missing:
1. Analyze exports and imports of existing files
2. Determine layer (Domain, Application, Adapter, Config)
3. Generate using template:

\`\`\`markdown
# Context: [Directory Name]

## 🎯 Responsibility
[One sentence description]

## 🏗️ Architectural Role
- **Type:** [Core / Infrastructure / Adapter]
- **Direction:** [Inbound / Outbound / Internal]

## 🚦 Dependency Rules
- **Allowed:** [List]
- **Forbidden:** [List]
\`\`\`

## 🤝 Integration with Other Skills

- **PRD Generator**: Reads context.md to understand constraints
- **Ralph Converter**: Uses context to infer file paths
- **Progress Tracking**: Updates AGENTS.md with learnings
```

## Skill Discovery & Activation

### How Agents Find Skills

1. **Directory Scan**: Agent scans `.claude/skills/` or configured skill directories
2. **Metadata Parse**: Reads frontmatter (`name`, `description`, `triggers`)
3. **Activation Match**: When user input matches triggers or description
4. **Progressive Loading**: Loads full instructions only when activated

### Example Activation Flow

```
User: "I want to create a PRD for user authentication"

Agent Decision Tree:
1. Scans skills directory
2. Finds "prd-generator" skill
3. Matches trigger: "create a prd"
4. Loads SKILL.md fully
5. Sees integration with "context-manager"
6. Loads context-manager as well
7. Follows combined workflow:
   - Check context.md in target directories
   - Generate PRD with architecture constraints
   - Ensure acceptance criteria include browser verification
```

## Cross-Tool Compatibility

### Claude Code
```bash
# Skills directory
~/.claude/skills/

# Project-specific skills
<project>/.claude/skills/
```

### Cursor
```bash
# Skills directory
<project>/.cursor/skills/
```

### Gemini
```bash
# Skills directory
<project>/.gemini/skills/
```

### Windsurf
```bash
# Skills directory
<project>/.windsurf/skills/
```

### Codex
```bash
# Skills directory
<project>/.codex/skills/
```

### Antigravity (generic)
```bash
# Skills directory
<project>/.agent/skills/
```

### Universal Pattern

Place skills in project root for portability:

```bash
<project>/.agentskills/
├── prd-generator/
├── ralph-converter/
└── context-manager/
```

Update tool configs to point to `.agentskills/`.

## Publishing to agentskills.io

Skills can be published to the agentskills.io registry:

```bash
# Package skill
cd .claude/skills/prd-generator
tar -czf prd-generator-v1.0.0.tar.gz .

# Publish (hypothetical - depends on registry API)
agentskills publish prd-generator-v1.0.0.tar.gz \
  --author agnostic-agent \
  --license MIT
```

## Summary

**Current State:**
- ✅ Skills are agentskills.io compliant (flat format)
- ✅ YAML frontmatter with name and description
- ✅ Clear activation triggers
- ✅ Portable across tools

**Optional Enhancements:**
- Directory structure for advanced features
- Example files and validators
- Extended frontmatter metadata
- Integration with context-manager skill
- Publishing to skill registry

The current implementation works perfectly with Claude Code and other tools. The enhancements are **optional** for teams wanting advanced features like validation scripts and extensive examples.
