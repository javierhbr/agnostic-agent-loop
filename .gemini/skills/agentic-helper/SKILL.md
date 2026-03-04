---
name: agentic-helper
description: "CLI guide, executor, and automator. Use when asking about agentic-agent commands, workflow selection, task management, SDD methodology, or how to start/continue/complete work."
---

# skill:agentic-helper

## Does exactly this

Guides you through agentic-agent CLI commands, helps select the right workflow (Tiny/Small/OpenSpec+/Full SDD) based on risk, and automates context generation and validation.

---

## Use this skill when

- You're unsure which workflow to use (Tiny vs Small vs SDD)
- You want to start, claim, or complete a task
- You want to generate context or run gate checks
- You're asking about agentic-agent CLI commands
- You need to understand SDD or OpenSpec methodology

---

## Step 1 — Assess Risk (Workflow Decision Tree)

Ask: **"What are you building and how risky is it?"**

| Signal | Workflow | Setup | Key Commands |
|---|---|---|---|
| Bug fix, < 1 day, 1 file | **Tiny** | 5 min | `task create`, `claim`, `complete` |
| Feature, 1-2 weeks, 1 service | **Small** | 10 min | `openspec init`, tasks |
| Multi-package, monorepo | **OpenSpec+** | 20 min | `openspec init` per package |
| Payment/auth/PII/contract | **Full SDD** | 1+ hour | `sdd start --risk critical` |

**Risk escalators** (each one = increase tier):
- Touches payment, auth, or PII → **Critical**
- Breaks an API/contract → **High or Critical**
- 4+ services involved → **High**
- No rollback strategy → **increase one level**

---

## Step 2 — Execute by Tier

See `resources/workflow-commands.md` for full command sequences for each tier (TINY, SMALL, OpenSpec+, Full SDD) with examples.

---

## Step 3 — Proactive Automation

**Before editing any directory:**
- Check for and read `AGENTS.md` (Read-Before-Write rule)
- If missing: `agentic-agent context generate <DIR>`

**Before completing any task:**
```bash
agentic-agent validate
```

**After architectural changes:**
- Update affected `AGENTS.md` files immediately
- Respect hexagonal boundaries: Core/Domain → nothing; Core/App → Domain only; Infra → Domain+App

---

## Hard Stops (Never Break These)

- Never implement when `blocked_by` is non-empty
- Never skip `task claim` (loses traceability)
- Never edit `.agentic/` YAML files directly
- Never merge without verify.md from Verifier

---

## Error Recovery

| Symptom | Cause | Fix |
|---|---|---|
| Gate fails | Missing Source or empty section | Fix spec, re-run gate-check |
| `blocked_by` non-empty | Unresolved ADR | `sdd adr resolve <id>` |
| Validation fails | Context stale or scope violated | Run context generate, validate again |
| "Which workflow?" | Risk not assessed | Walk the decision tree above |

---

## If you need more detail

→ `resources/workflow-commands.md` — Full command sequences for each tier with examples, error recovery details, and workflow selection flow
