---
name: context-manager
description: System rules for context management and architectural boundaries. These are mandatory rules, not optional techniques. See AGENT_RULES.md for complete documentation.
---

# Context Management Rules

⚠️ **This is now a system-level rule, not an optional skill.**

**For complete documentation, see [AGENT_RULES.md](../../../../AGENT_RULES.md) in the repository root.**

---

## Quick Reference

### Mandatory: Read-Before-Write

Before modifying ANY file in a directory:

1. Check if `context.md` exists in that directory
2. If it exists, **read it** - identify dependencies, architectural role, constraints
3. If it does NOT exist and the directory has source files, generate context
4. Only after reading context may you edit files
5. Update `context.md` after architectural changes

### Red Flags — STOP and Fix

- ❌ Editing without reading `context.md` first
- ❌ Creating directories without generating context
- ❌ Importing forbidden dependencies
- ❌ "Too small to need context" — no task is too small
- ❌ "I'll update context later" — later means never

### Hexagonal Architecture Boundaries

| Layer | Can Depend On | Cannot Depend On |
|-------|--------------|-------------------|
| Core/Domain | Nothing | Application, Infrastructure, Config |
| Core/Application | Domain only | Infrastructure, Config |
| Infrastructure/Adapters | Domain, Application | Other adapters directly |
| Infrastructure/Config | All layers | — |

---

## Full Documentation

**👉 See [AGENT_RULES.md](../../../../AGENT_RULES.md) for:**

- Complete context management workflow
- Code quality rules
- Task management rules
- Documentation guidelines
- Communication protocols
- Emergency procedures

**These rules apply to all agents regardless of tool (Claude, OpenCode, Copilot).**

---

## Tool-Specific Commands

- **Claude**: See `CLAUDE.md` for `agentic-agent` CLI commands
- **OpenCode**: See `OPENCODE.md` for tool integration
- **Copilot**: See `COPILOT.md` for GitHub-specific features
