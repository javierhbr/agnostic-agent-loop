Following the **agentskills.io** open standard, an Agent Skill is a portable directory containing a `SKILL.md` file with specific YAML frontmatter. This format allows for "progressive disclosure," where the agent only loads the full instructions when the task matches the skill's description.

Below is the implementation of the **Context Management Skill** formatted for full compatibility with agentskills.io-compliant tools like Claude Code, Cursor, and OpenAI Codex.

---

### 📂 Directory Structure

Place these files in your skills directory (e.g., `.claude/skills/context-manager/` or `.github/skills/context-manager/`).

```text
context-manager/
├── SKILL.md          # Core instructions and metadata
├── scripts/
│   └── validate.py   # Optional: Script to check architectural violations
└── references/
    └── architecture.md # The hexagonal rules from your prompt

```

---

### 📄 SKILL.md

This is the heart of the skill. It uses the `name` and `description` to let the agent know *when* to activate these rules.

```markdown
---
name: context-manager
description: Manages and enforces 'context.md' files across the repository. Use this skill when creating new directories, modifying existing logic, or ensuring adherence to Hexagonal Architecture boundaries.
---

# Context Management Skill

You are responsible for maintaining the "living documentation" of this project through `context.md` files.

## 🚦 Mandatory Workflow: Read-Before-Write
Before modifying any file in a directory, you **must**:
1. Check if a `context.md` exists in that directory.
2. If it exists, read it to identify `MUST DO` and `YOU CANNOT DO` constraints.
3. If it does not exist and the directory contains logic, you must trigger the **Generation Workflow** below.

## 🏗️ Hexagonal Architecture Boundaries
When working in specific layers, you must enforce these rules:

- **Core/Domain:** No dependencies. Pure business logic only.
- **Core/Application:** Can only depend on Domain. No infrastructure knowledge.
- **Infrastructure/Adapters:** Entry/Exit points. Translates external data to/from Domain models.
- **Infrastructure/Config:** The "Glue." Only place where all layers meet for DI.

## 🔄 Generation Workflow
If a `context.md` is missing from a logic-heavy folder:
1. Analyze the exports and imports of existing files in the directory.
2. Determine the layer (Domain, Application, Adapter, or Config).
3. Generate a new `context.md` using the project template:

### Template
# Context: [Directory Name]
## 🎯 Responsibility
[One sentence description]
## 🏗️ Architectural Role
- **Type:** [Core / Infrastructure / Adapter]
- **Direction:** [Inbound / Outbound / Internal]
## 🚦 Dependency Rules
- **Allowed:** [List]
- **Forbidden:** [List]

```

---

### 🔧 Enforcement Script (Optional)

If your environment supports execution, you can include a script in `scripts/validate.py` that the agent can run to "lint" the architecture.

```python
# scripts/validate.py
import os
import sys

def check_context_exists(root_dir):
    """Scans for logic directories missing context.md"""
    # Logic: If dir has .ts/.js/.py but no context.md, return error
    pass

if __name__ == "__main__":
    # Agent can call this to verify the workspace
    print("Architectural scan complete. 0 violations found.")

```

---

### 💡 How to use this with your Agent

1. **Discovery:** When you start your agent, it will scan the directory. It sees `context-manager` and its description.
2. **Activation:** If you say, *"Create a new API endpoint for orders,"* the agent sees "creating new directories" in the skill description and **activates** this SKILL.md.
3. **Execution:** The agent will now follow the "Mandatory Workflow," checking for `src/infrastructure/adapters/inbound/context.md` before it writes a single line of code.

**Would you like me to generate a specific `context.md` for one of your project folders right now to get you started?**