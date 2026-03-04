# OpenClaw Quick Start (5 Minutes)

Get a working autonomous app factory running in 5 minutes.

---

## 1. Setup (2 min)

```bash
# Create project
mkdir -p ~/trailmate && cd ~/trailmate

# Initialize
agentic-agent init --name "Trailmate" --no-interactive

# Install OpenClaw skill pack
agentic-agent skills install openclaw --tool claude-code --no-interactive
```

## 2. Create First Task (1 min)

```bash
agentic-agent task create \
  --title "Research: Offline Fitness Tracker" \
  --description "Find market opportunity for offline GPS fitness app." \
  --scope docs \
  --acceptance "One-pager written"
```

Copy the task ID (e.g., `TASK-1234567-1`).

## 3. Claim Task (1 min)

```bash
agentic-agent task claim TASK-1234567-1 --no-interactive
agentic-agent context build --task TASK-1234567-1 --format json > /tmp/context.json
```

## 4. Spawn Researcher Worker (1 min)

```bash
/subagents spawn openclaw-researcher "
Task: TASK-1234567-1

Scan Reddit, X, App Store for fitness app complaints.
Identify market gap.

Write one-pager to docs/research/pitch.md:
- Problem statement
- Target audience (specific niche)
- Core features (minimum viable)
- Market demand (search volume, discussions)
- Revenue model

When complete: <promise>COMPLETE</promise>
"
```

## 5. Wait & Complete

Watch announcements:
```bash
watch -n 5 'cat .agentic/coordination/announcements.yaml | tail -10'
```

When researcher announces complete:
```bash
agentic-agent task complete TASK-1234567-1
```

Done! You just ran your first autonomous worker.

---

## Next: Try Orchestrator Mode

Spawn the full orchestrator to manage multiple phases:

```bash
/subagents spawn openclaw-orchestrator "
Read: .agentic/project-state.md

For each task in backlog:
1. agentic-agent task claim <ID>
2. agentic-agent context build --task <ID>
3. Spawn 2-4 workers (based on scope)
4. Poll .agentic/coordination/announcements.yaml
5. When all workers complete: agentic-agent task complete <ID>
6. Repeat

Max 10 iterations.
Output: <promise>COMPLETE</promise>
"
```

---

## Key Files Created

```
.agentic/
├── tasks/
│   ├── backlog.yaml       ← Your tasks
│   ├── in-progress.yaml
│   └── done.yaml
├── coordination/
│   ├── reservations.yaml  ← File locks (workers write here)
│   ├── announcements.yaml ← Results (workers write here)
│   └── kill-signals.yaml  ← Stop signals (orchestrator writes here)
├── checkpoints/
│   └── TASK-ID-*.json     ← Worker progress (context compaction)
└── context/
    ├── global-context.md
    └── tech-stack.md
```

---

## Verifying It Works

```bash
# Check task moved to in-progress
agentic-agent task list

# Check announcements from worker
cat .agentic/coordination/announcements.yaml

# Check checkpoint saved
ls -la .agentic/checkpoints/

# Check task moved to done
agentic-agent task list
```

---

## Common First Errors

| Error | Fix |
|-------|-----|
| `task not found` | Verify task ID matches `TASK-*` format from create output |
| `context build failed` | Run `agentic-agent context generate docs` first |
| `no announcements` | Worker may still be running. Wait 10s and re-check. |
| `permission denied .agentic/` | Run from project root where you ran `init` |

---

## Next Steps

1. Read the full example: `examples/openclaw-app-factory/README.md`
2. Study playbooks: `.claude/skills/openclaw/resources/*.md`
3. Scale to 3+ workers in parallel
4. Add reviewer phase (quality gate)
5. Run 24/7 on cron

You now have a working autonomous agent factory! 🚀
