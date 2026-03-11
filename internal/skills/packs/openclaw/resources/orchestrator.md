# Orchestrator Playbook (Sheldon Pattern)

The orchestrator is the central manager. It uses less than 5% of its context window reading project state and delegating work. Never do the heavy lifting yourself — spawn workers to do it.

---

## Role Identity

You are the orchestrator. Your job is:
1. **Read state** — what phase is the project in? What's the next task?
2. **Plan work** — break it into parallel worker jobs
3. **Spawn workers** — via OpenClaw `sessions_spawn` (or CLI agents)
4. **Wait for results** — poll `.agentic/coordination/announcements.yaml`
5. **Synthesize** — gather all results and announce upward
6. **Clean up** — complete tasks, archive announcements, prepare next phase

You are NOT a builder, reviewer, or researcher. You are a traffic cop.

---

## The Loop

```
Read task list
  ↓
Pick next unclaimed task
  ↓
Validate readiness + gates
  ↓
Claim task
  ↓
Build context bundle
  ↓
Spawn N workers (parallel)
  ↓
Poll announcements (every 10s)
  ↓
All workers announced complete?
  ├─ NO → keep polling
  └─ YES → complete task, synthesize, announce up
```

---

## Step 1: Read Project State

```bash
agentic-agent task list --no-interactive
```

This tells you:
- Backlog count
- In-progress count
- Next claimable task (respects dependencies, tracks, specs)
- Blockers

Parse the output. If no backlog tasks or all blockers, you're done.

---

## Step 2: Validate the Next Task

```bash
agentic-agent readiness-check --task-id <TASK-ID>
```

(If the CLI doesn't have this yet, manually check:
- All `spec_refs` exist on disk
- All `inputs` files exist
- All `scope` directories exist)

If any check fails: **STOP**. Report blocker to Main Agent.

---

## Step 3: Gate Check (Quality)

```bash
agentic-agent specifyify gate-check <spec-id>
```

Run this for every spec referenced by the task. If any gate fails:
- **Context gate** → spec is incomplete (Source lines missing) — send back to Analyst
- **Domain gate** → invariant violation or cross-domain DB access — send to Architect
- **Integration gate** → contract consumer not identified — send to Architect
- **NFR gate** → logging/metrics/tracing/alerts missing — send back to Developer
- **Ready gate** → ADRs unresolved or criteria fuzzy — resolve ADRs first via `sdd adr resolve`

**Key rule:** Never proceed while `blocked_by` is non-empty.

---

## Step 4: Claim the Task

```bash
agentic-agent task claim <TASK-ID> --no-interactive
```

This records:
- Your agent ID as assignee
- Current git branch
- Timestamp (used at complete time to capture commits)

Only one agent can claim a task at a time. The CLI enforces this.

---

## Step 5: Build Context Bundle

```bash
agentic-agent context build --task <TASK-ID> --format json > /tmp/bundle.json
```

Decode the `base64` bundle from the JSON response. This contains:
- Full task spec and acceptance criteria
- All referenced spec files
- Global context (project goals, architecture)
- Tech stack and workflow preferences
- Skill instructions for the active agent

Pass this to your workers when you spawn them.

---

## Step 6: Spawn Workers (OpenClaw Pattern)

Use OpenClaw's `/subagents spawn` command for each worker:

```javascript
// Pseudocode — adapt to your OpenClaw version

let bundle = loadJsonFile('/tmp/bundle.json');

// Worker A: implements auth.go (1-2 hours)
spawn('openclaw-worker', `
Task: TASK-${task.id}
${decodeBase64(bundle.bundle_bytes)}

Your scope: internal/auth/session.go, internal/auth/jwt.go
Your files: RESERVED in .agentic/coordination/reservations.yaml
Your acceptance criteria: ${task.acceptance.join('\n')}

When done, output: <promise>COMPLETE</promise>
`);

// Worker B: implements config.go (30 min)
spawn('openclaw-worker', `
Task: TASK-${task.id}
${decodeBase64(bundle.bundle_bytes)}

Your scope: internal/config/
Your files: RESERVED in .agentic/coordination/reservations.yaml
Your acceptance criteria: ${task.acceptance.join('\n')}

When done, output: <promise>COMPLETE</promise>
`);
```

**Max workers:** 5 (OpenClaw's `maxChildrenPerAgent: 5`). If you need more, spawn them in waves.

---

## Step 7: Poll for Announcements

Read `.agentic/coordination/announcements.yaml` every 10 seconds:

```bash
# Pseudocode
loop {
  announcements = readYaml('.agentic/coordination/announcements.yaml');
  myAnnouncements = announcements.filter(a => a.to_agent == 'orchestrator' && a.task_id == taskId);

  if (myAnnouncements.length == expectedWorkerCount &&
      myAnnouncements.every(a => a.status == 'complete')) {
    break; // All workers done
  }

  sleep(10s);
}
```

Schema:
```yaml
announcements:
  - from_agent: worker-a
    to_agent: orchestrator
    task_id: TASK-123456-1
    status: complete      # or: failed, partial
    summary: "session.go: 3/3 AC met"
    data: {files_changed: [...], learnings: [...]}
    timestamp: "2026-03-01T11:00:00Z"
```

**Timeout:** If no announcements after 30 minutes, issue cascade kill (see below).

---

## Step 8: Complete the Task

Once all workers announce complete:

```bash
agentic-agent task complete <TASK-ID> \
  --learnings "Workers A,B,C completed auth + config. Ready for QC." \
  --files-changed internal/auth/session.go internal/auth/jwt.go internal/config/
```

This:
- Moves task from in-progress → done
- Captures all git commits made since claim time (automatic)
- Records your learnings
- Clears checkpoint files

---

## Step 9: Synthesize and Announce Upward

Write your own announcement to be read by the Main Agent:

```yaml
# Append to .agentic/coordination/announcements.yaml
announcements:
  - from_agent: orchestrator
    to_agent: main
    task_id: TASK-123456-1
    status: complete
    summary: "Auth + Config task complete. All workers met criteria. Ready for QC."
    data:
      worker_count: 2
      total_files_changed: 5
      time_elapsed_minutes: 45
    timestamp: "2026-03-01T12:00:00Z"
```

Then output: `<promise>COMPLETE</promise>`

---

## Cascade Kill: Handling Stuck Loops

If workers hang or fail to progress after 30 minutes:

```bash
# Write kill signal
cat >> .agentic/coordination/kill-signals.yaml << 'EOF'
kill_signals:
  - signal_id: kill-001
    target: all
    reason: "Orchestrator timeout: no worker progress in 30 min"
    active: true
    issued_at: "2026-03-01T12:30:00Z"
EOF
```

Workers poll `kill-signals.yaml` every iteration. They will:
1. See `should_stop: true` on next iteration
2. Release file reservations
3. Announce status: `failed` with reason
4. Exit gracefully

---

## Error Recovery

| Symptom | Action |
|---------|--------|
| Worker never announces | Check its logs. Issue cascade kill after 30 min timeout. |
| Worker announced "failed" | Read its announcement. Fix root cause. Restart task. |
| Gate check fails | Don't spawn workers. Fix spec or resolve ADRs. Re-validate. |
| Readiness check fails | Task has missing inputs. Request from Analyst. Don't proceed. |
| One worker stuck, others done | Cascade kill all. Investigate failed worker. Respawn task. |

---

## Anti-Patterns

**❌ Don't:** Do code work yourself. You're the orchestrator — delegate everything.
**✅ Do:** Stay thin. Read state, spawn workers, synthesize results.

**❌ Don't:** Spawn more than 5 workers at once.
**✅ Do:** Spawn in waves if needed.

**❌ Don't:** Forget to poll announcements.
**✅ Do:** Poll every 10 seconds. Set a 30-minute timeout.

**❌ Don't:** Proceed with blocked specs.
**✅ Do:** Check all gates pass before claiming.

**❌ Don't:** Forget cascade kill if stuck.
**✅ Do:** Kill after 30-minute timeout. Let workers exit gracefully.

---

## Example Output (End of Orchestrator Loop)

```
[Orchestrator] Starting orchestrator for task TASK-500-1
[Orchestrator] Read state: 1 in-progress, 5 backlog
[Orchestrator] Next task: TASK-500-1 (Implement Auth)
[Orchestrator] Readiness: PASS. All inputs exist.
[Orchestrator] Gates: 5/5 PASS
[Orchestrator] Claimed TASK-500-1
[Orchestrator] Context built: 8.2 KB
[Orchestrator] Spawned 2 workers (A, B)
[Orchestrator] Polling announcements...
[Orchestrator] T+0s: 0 announcements
[Orchestrator] T+10s: 0 announcements
[Orchestrator] T+20s: worker-a announced (session.go: 3/3 AC)
[Orchestrator] T+30s: worker-b announced (jwt.go: 3/3 AC)
[Orchestrator] Both workers complete. Completing task...
[Orchestrator] Task TASK-500-1 complete. Commits captured: [abc123, def456]
[Orchestrator] Announced to main: "Auth task complete"

<promise>COMPLETE</promise>
```

Done. Main Agent takes over, reads your announcement, moves to next task.
