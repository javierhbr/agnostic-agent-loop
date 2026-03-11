# Agentic Helper — Custom Claude Code Agent

**Created:** 2026-02-26

## Overview

The **Agentic Helper** is a custom Claude Code subagent that helps users understand and execute the agentic-agent CLI and SDD methodology. It serves three simultaneous roles:

1. **🎓 Teacher** — Explains *why* steps exist, not just *what* to run
2. **🛠️ Hands-on Executor** — Actively executes workflows: claims tasks, generates context, runs validation
3. **🔄 Proactive Automator** — Checks AGENTS.md, generates context, runs validation without being asked

## Installation

The agent is automatically installed as a mandatory pack when you run:

```bash
agentic-agent skills ensure --agent claude-code
```

Or manually install it:

```bash
agentic-agent skills install agentic-helper --tool claude-code
```

## Where Files Live

- **Claude Code Agent:** `.claude/agents/agentic-helper.md` (frontmatter + system prompt, ~170 lines)
- **Skill for other tools:** `.claude/skills/agentic-helper/SKILL.md` (skill format, ~180 lines)

## How It Works

The agent auto-delegates when you ask Claude Code about:

- `agentic-agent` CLI commands (e.g., "what's the task claim command?")
- SDD or OpenSpec methodology (e.g., "when should I use SDD?")
- Workflow selection (e.g., "I'm building a feature, what workflow should I use?")
- Starting or completing work (e.g., "help me claim a task", "run gate check")
- Understanding the methodology (e.g., "explain SDD", "what are the 5 gates?")

## Core Workflow Decision Tree

The agent guides you through a single decision tree question:

**"What are you building and how risky is it?"**

| Signal | Workflow | Key Commands |
|---|---|---|
| Bug fix, < 1 day, 1 file | **Tiny** | `task create`, `task claim`, `task complete` |
| Feature, 1-2 weeks, 1 service | **Small** | `openspec init`, then tasks |
| Multi-package, monorepo | **OpenSpec+** | `openspec init` per package |
| Payment / auth / PII / breaking contract | **Full SDD** | `sdd start --risk critical` |

## Execution Support

The agent knows every command for each workflow tier:

### Tiny Workflow (Bug Fixes)
```bash
agentic-agent task create --title "Fix: <description>"
agentic-agent task claim <ID>
agentic-agent context generate <DIR>
# ... implement ...
agentic-agent validate
agentic-agent task complete <ID>
```

### Small Workflow (Single Service Features)
```bash
agentic-agent openspec init "<name>" --from <prd-file>
agentic-agent task claim <ID>
agentic-agent context generate <DIR>
# ... implement ...
agentic-agent validate
agentic-agent task complete <ID>
agentic-agent openspec complete <change-id>
```

### Full SDD Workflow (High/Critical Risk)
```bash
agentic-agent specifyify start "<name>" --risk critical
agentic-agent specifyify workflow show <id>
agentic-agent specifyify gate-check <spec-id>
# ... phases 0-4 with gates ...
```

## Proactive Automation

The agent automatically:

1. **Checks AGENTS.md** before you edit any directory (Read-Before-Write rule)
2. **Generates context** if missing: `agentic-agent context generate <DIR>`
3. **Runs validation** before task completion: `agentic-agent validate`
4. **Updates AGENTS.md** after architectural changes
5. **Detects drift** and suggests `agentic-agent skills ensure`

## Teaching Approach

The agent explains the *why* at point of use. Example:

> "We require `task claim` before starting because it records your git branch and timestamp. When you complete later, we capture the commits you made. This gives us traceability and lets us audit scope violations."

## Error Recovery

The agent provides a clear error table:

| Symptom | Cause | Fix |
|---|---|---|
| Gate fails | Missing Source line or empty section | Read gate-check skill, fix spec, re-run |
| `blocked_by` non-empty | Unresolved ADR | `sdd adr list --blocked`, resolve ADR |
| Task won't complete | Validation failing | `agentic-agent validate`, read error, fix |
| Context stale | Forgot `context generate` | Run it now, re-read AGENTS.md |

## Technical Implementation

### Files Created

1. **`internal/skills/packs/agentic-helper/AGENT.md`** — Claude Code subagent with YAML frontmatter
2. **`internal/skills/packs/agentic-helper/SKILL.md`** — Companion skill for non-Claude tools

### Files Modified

1. **`internal/skills/packs.go`**
   - Added `ToolAgentDir` and `ToolGlobalAgentDir` maps (claude-code → `.claude/agents`)
   - Added `IsAgent` field to `SkillPackFile` struct
   - Registered `agentic-helper` pack (2 files: AGENT.md + SKILL.md)
   - Added `agentic-helper` to `MandatoryPacks` list

2. **`internal/skills/installer.go`**
   - Extended `Install()` to route `IsAgent` files to `ToolAgentDir` instead of `ToolSkillDir`
   - Added `resolveAgentOutputDir()` helper function for agent directory resolution

3. **`internal/skills/packs_test.go`**
   - Updated test expectations to include `agentic-helper` (2 files)
   - Updated SDD pack file count (15 files, not 5)
   - Added `superpowers-bridge` to expected packs

## Usage Examples

### Example 1: User Asks About Workflow

```
User: "I want to fix a bug in the payment system. What should I do?"

Claude Code delegates to agentic-helper:
- Runs: agentic-agent status
- Asks: "Is this a critical payment bug? Will it affect transactions?"
- Recommends: Full SDD workflow (high risk)
- Shows: All commands for Phase 0-4
```

### Example 2: User Needs to Execute a Task

```
User: "Help me claim task TASK-001 and get ready to work"

agentic-helper:
- Runs: agentic-agent task claim TASK-001
- Checks: AGENTS.md in the directories you'll edit
- Runs: agentic-agent context generate <DIR>
- Reports: Ready to code
```

### Example 3: User Is About to Complete Work

```
User: "I think I'm done with my changes"

agentic-helper (proactively):
- Runs: agentic-agent validate
- Reports: Validation passed ✓
- Runs: agentic-agent task complete TASK-001
- Reports: Task completed with commits recorded
```

## Agent Configuration

```yaml
---
name: agentic-helper
description: >
  Use when the user asks about agentic-agent CLI commands, SDD methodology,
  OpenSpec, when to use which workflow, how to start/continue/complete a task,
  or says "claim a task", "generate context", "run gate check", "validate my
  work", "explain SDD", "what workflow should I use", "help me start a feature",
  or pastes an agentic-agent command.
tools: Read, Write, Edit, Bash, Glob, Grep, Task
model: sonnet
memory: project
---
```

## Tools Available

- **Read** — Read CLAUDE.md, AGENTS.md, specs, task metadata
- **Write** — Create spec files, verify.md, task metadata
- **Edit** — Update AGENTS.md after changes, fix drift
- **Bash** — Execute `agentic-agent` CLI commands
- **Glob** — Find missing context.md files, audit AGENTS.md coverage
- **Grep** — Search for AGENTS.md, find context drift
- **Task** — Spawn parallel sub-agents for multi-component work

## Memory

The agent uses **project-scoped memory** (`.claude/agent-memory/agentic-helper/`) to remember:

- Task IDs and initiatives across sessions
- Which directories have been context-generated
- Workflow patterns specific to the project
- Lessons learned (command failures, edge cases)

## Next Steps

1. **For users:** Ask Claude Code questions about the CLI or methodology
2. **For teams:** Share this agent across projects via `skills install agentic-helper`
3. **For customization:** Modify the system prompt in `.claude/agents/agentic-helper.md` to add project-specific context

## Key Design Decisions

### Single Agent vs. Multiple Agents

A single agent minimizes coordination friction and keeps context unified. The agent's system prompt is strictly structured to remain scannable even at 170 lines.

### Auto-Installation as Mandatory Pack

Making `agentic-helper` mandatory ensures every Claude Code user gets the benefit of guided methodology. No extra setup required — just `skills ensure`.

### Separate AGENT.md and SKILL.md

- **AGENT.md** is for Claude Code (uses `tools:` and `model:` fields)
- **SKILL.md** is for other tools like Cursor and Gemini (uses "Use this skill when" sections)

This allows the agent to work across all agent tools supported by the CLI.

### ToolAgentDir Map

Claude Code supports agents (subagents). Other tools don't. The `ToolAgentDir` map cleanly separates agent-capable tools (claude-code, opencode) from skill-only tools (cursor, gemini). For skill-only tools, the SKILL.md falls back to `.claude/skills/agentic-helper/`.

## Testing

All tests pass:

```bash
go test ./internal/skills/... -v
# 100+ tests, all passing
# New tests cover agentic-helper pack registration and agent file routing
```

The implementation has been validated:

```bash
# Build succeeds
go build ./cmd/agentic-agent/.

# Agent installs correctly
agentic-agent skills install agentic-helper --tool claude-code
# ✓ Created .claude/agents/agentic-helper.md
# ✓ Created .claude/skills/agentic-helper/SKILL.md

# Auto-installation works
agentic-agent skills ensure --agent claude-code
# ✓ agentic-helper installed (first in mandatory packs list)
```
