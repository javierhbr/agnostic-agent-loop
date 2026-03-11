---
name: openclaw
description: >
  Use when building or running an OpenClaw autonomous agent factory, when asked to
  "set up openclaw", "create openclaw agents", "implement orchestrator pattern",
  "build autonomous agent", "spawn workers", "announce results", or when configuring
  multi-agent pipelines using agnostic-agent-loop CLI.
---

# OpenClaw Autonomous Agent Factory

Build a multi-agent system using agnostic-agent-loop and OpenClaw's orchestrator pattern.

## Quick Decision Tree

**Are you the Orchestrator (Sheldon)?**
- Read state, decide what task comes next
- Claim task, spawn workers to parallelized work
- Poll workers for completion announcements
- Use: `task list`, `task claim`, `context build`, `task complete`, `.agentic/coordination/`

**Are you a Worker?**
- Do one assigned task
- Watch for cascade kill signals
- Reserve files you edit
- Announce when done
- Use: `context build`, `task complete`, file reservation YAML, kill-signal YAML

**Are you a Researcher (Shan)?**
- Scan for opportunities
- Create one-pager pitches
- Use: `task create`, task `complete`, announcement YAML

**Are you a Reviewer (independent)?**
- Verify code and specs
- Run gates and validation
- Report findings
- Use: `sdd gate-check`, `validate`, announcement YAML

---

## Role Playbooks

- **Orchestrator**: `openclaw/resources/orchestrator.md` — read state → spawn workers → synthesize → announce
- **Worker**: `openclaw/resources/worker.md` — kill-check → reserve files → work → announce
- **Researcher**: `openclaw/resources/researcher.md` — scan → pitch → announce
- **Reviewer**: `openclaw/resources/reviewer.md` — gate-check → validate → report

---

## Shared Coordination Files

All agents read/write these YAML files directly (see `openclaw/resources/coordination.md`):

| File | Used By | Purpose |
|------|---------|---------|
| `.agentic/coordination/reservations.yaml` | Workers | Soft file locks (TTL 10 min) |
| `.agentic/coordination/announcements.yaml` | Workers + Orchestrator | Result propagation (child → parent) |
| `.agentic/coordination/kill-signals.yaml` | Main + Orchestrator | Cascade stop signals |

---

## Critical Rules

| Rule | Why |
|------|-----|
| Always `task claim` before starting work | Records git branch + timestamp for traceability |
| Always check kill-signals on iteration start | Enables graceful shutdown |
| Always release file reservations after edits | Prevents deadlock on next run |
| Always output `<promise>COMPLETE</promise>` when done | Signals loop termination |
| Reviewer never uses Builder's model | Prevents bias and corner-cutting |
| Orchestrator stays thin (<5% context) | Fat orchestrator fails under load |

---

## Anti-Patterns

| ❌ Don't | ✅ Do |
|---------|-------|
| Orchestrator does heavy lifting | Delegate to workers, stay in read-only mode |
| Skip `task claim` | Always claim first — enables git tracking |
| Forget kill-signal check | Check `.agentic/coordination/kill-signals.yaml` each iteration |
| Edit same file concurrently | Use file reservation schema (TTL 10 min, release after) |
| Use same model to build and review | Assign different model to Reviewer agent |
| No stop signal | Always end output with `<promise>COMPLETE</promise>` |

---

## Key CLI Commands by Role

### Orchestrator
```bash
agentic-agent task list --no-interactive           # What's next?
agentic-agent validate                             # Gate checks
agentic-agent specify gate-check <spec-id>             # Quality check
agentic-agent task claim <ID> --no-interactive     # Reserve task
agentic-agent context build --task <ID>            # Gather context
agentic-agent task complete <ID>                   # Mark done
```

### Worker
```bash
agentic-agent context build --task <ID>            # Get full context
agentic-agent task complete <ID>                   # Mark done
# Plus direct YAML read/write for reservations, kill-signals, announcements
```

### Researcher
```bash
agentic-agent task create --title "Research: <topic>"    # Log finding
agentic-agent task complete <ID>                         # Mark done
```

### Reviewer
```bash
agentic-agent specify gate-check <spec-id>             # Quality check
agentic-agent validate                             # Validation rules
```

---

## Next Steps

1. **Decide your role** — are you orchestrator, worker, researcher, or reviewer?
2. **Read the role playbook** — e.g., `openclaw/resources/orchestrator.md`
3. **Understand coordination** — read `openclaw/resources/coordination.md` for YAML schemas
4. **Start the loop** — follow step-by-step instructions in your playbook
5. **Always end with** `<promise>COMPLETE</promise>` so OpenClaw knows when you're done
