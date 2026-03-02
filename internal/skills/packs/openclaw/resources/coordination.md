# Coordination: Shared State Files

The orchestrator and workers coordinate via three YAML files in `.agentic/coordination/`. These are the "shared blackboard" that enables multi-agent work without a central message queue.

---

## File Locations

All three files are in:
```
.agentic/coordination/
├── reservations.yaml      # Soft file locks (TTL 10 min)
├── announcements.yaml     # Child → Parent result propagation
├── kill-signals.yaml      # Cascade stop control
└── .lock                  # Internal lock file (skip this)
```

Create the directory if it doesn't exist:
```bash
mkdir -p .agentic/coordination
```

---

## reservations.yaml

**Purpose:** Prevent concurrent writes to the same file by different workers.

**Schema:**
```yaml
reservations:
  - reservation_id: res-worker-a-001
    file_path: internal/auth/session.go
    owner: openclaw-worker-1
    task_id: TASK-500-1
    created_at: "2026-03-01T10:00:00Z"
    expires_at: "2026-03-01T10:10:00Z"
```

**Fields:**
- `reservation_id` — unique identifier (e.g., `res-{worker-id}-{timestamp}`)
- `file_path` — relative path to file being reserved (e.g., `src/auth.go`)
- `owner` — agent ID that holds this reservation (e.g., `worker-1`)
- `task_id` — which task this file belongs to
- `created_at` — when the reservation was made (ISO 8601)
- `expires_at` — when the reservation expires (TTL 10 minutes)

**Worker Responsibility:**

Before editing a file:
1. Read `.agentic/coordination/reservations.yaml`
2. Check if any existing reservation matches `file_path` and `owner != myId`
3. If yes, and reservation has not expired: WAIT or ask orchestrator to cascade kill
4. If expired or doesn't exist: append a new entry with `expires_at = now + 10min`
5. Edit the file
6. When done, remove your reservation entry

**Orchestrator Responsibility:**

Before spawning workers for a task:
1. Verify no two workers are assigned overlapping file paths (build `scope` with no intersection)
2. Check that no existing active reservations block this task
3. Inform workers of their reserved files via context bundle

**Example Conflict Resolution:**

Worker-A reserves `src/auth.go` until 10:10.
Worker-B tries to reserve the same file at 10:05.
Worker-B sees the conflict, reports to orchestrator:
```
"Cannot reserve src/auth.go. Held by worker-A until 2026-03-01T10:10:00Z"
```
Orchestrator choices:
- Wait for Worker-A's TTL to expire (OK for quick tasks)
- Cascade kill Worker-A and re-spawn with exclusive access
- Redesign scope: give Worker-B a different file

**TTL Semantics:**

Once `expires_at` passes, the reservation is "stale." Workers can safely ignore or delete stale reservations. The 10-minute TTL ensures that even if a worker crashes without releasing, the file becomes available again automatically.

---

## announcements.yaml

**Purpose:** Workers report results back to the orchestrator (and orchestrator reports to main agent).

**Schema:**
```yaml
announcements:
  - announcement_id: ann-worker-a-001
    from_agent: openclaw-worker-1
    to_agent: orchestrator
    task_id: TASK-500-1
    status: complete
    summary: "Implemented session.go and jwt.go. All 4 AC pass."
    data:
      files_changed: ["internal/auth/session.go", "internal/auth/jwt.go"]
      iterations: 2
      learnings: ["Redis for session", "JWT expiry: 24h"]
    timestamp: "2026-03-01T11:00:00Z"
```

**Fields:**
- `announcement_id` — unique ID (e.g., `ann-{from}-{timestamp}`)
- `from_agent` — who made this announcement (e.g., `worker-1`, `orchestrator`)
- `to_agent` — who reads this (e.g., `orchestrator`, `main`)
- `task_id` — which task this is about
- `status` — `complete` | `failed` | `partial`
- `summary` — human-readable result (1-2 sentences)
- `data` — JSON object with arbitrary details (files, learnings, metrics)
- `timestamp` — when announced (ISO 8601)

**Worker Responsibility:**

When done:
1. Generate unique `announcement_id`
2. Append entry to `.agentic/coordination/announcements.yaml`
3. Set `from_agent` = your agent ID
4. Set `to_agent` = `orchestrator`
5. Set `status` = one of: `complete` (all AC pass) | `failed` (critical error) | `partial` (some AC pass)
6. Write human-readable `summary`
7. Include `data` with files changed, iterations, learnings

**Orchestrator Responsibility:**

Poll for announcements every 10 seconds:
```bash
# Pseudocode
loop {
  announcements = readYaml('.agentic/coordination/announcements.yaml');
  myAnnouncements = announcements
    .filter(a => a.to_agent == 'orchestrator' && a.task_id == currentTaskId);

  if (myAnnouncements.length == expectedWorkerCount &&
      myAnnouncements.every(a => a.status == 'complete')) {
    // All workers done, all passed
    break;
  }

  sleep(10s);
  if (elapsedTime > 1800s) {  // 30 min timeout
    cascadeKill('timeout');
    break;
  }
}
```

After all workers announce:
1. Synthesize their results (count passed AC, list files changed, summarize learnings)
2. Complete the task via CLI
3. Append your own announcement:
   ```yaml
   - from_agent: orchestrator
     to_agent: main
     task_id: TASK-500-1
     status: complete
     summary: "Task complete. 2 workers, 10 files changed, ready for QC."
     data: {workers: ["worker-1", "worker-2"], ...}
   ```
4. Output `<promise>COMPLETE</promise>` so main agent resumes

**Status Legend:**

| Status | Meaning | Orchestrator Action |
|--------|---------|---|
| `complete` | All acceptance criteria pass | Proceed normally |
| `failed` | Critical error; cannot continue | Cascade kill remaining workers, investigate failure, respawn task |
| `partial` | Some criteria pass, some fail | Rare. Possible if orchestrator cascade kills mid-task. Record for audit. |

---

## kill-signals.yaml

**Purpose:** Orchestrator can abort all workers gracefully if looping hangs or runs out of iterations.

**Schema:**
```yaml
kill_signals:
  - signal_id: kill-orchestrator-001
    orchestrator_id: openclaw-orchestrator
    target: all
    reason: "Task timeout: no worker progress in 30 minutes"
    active: true
    issued_at: "2026-03-01T12:30:00Z"
```

**Fields:**
- `signal_id` — unique ID (e.g., `kill-{orchestrator}-{timestamp}`)
- `orchestrator_id` — which orchestrator issued this
- `target` — `all` | `<specific-agent-id>` (default: all)
- `reason` — human-readable explanation (for debugging)
- `active` — boolean (false = signal has been cleared)
- `issued_at` — when the kill was issued

**Orchestrator Responsibility:**

If a task times out (no announcements after 30 minutes):
1. Generate `signal_id`
2. Append entry to `.agentic/coordination/kill-signals.yaml`
3. Set `target` = `all` (or specific worker IDs if you want to be surgical)
4. Set `active` = `true`
5. Write `reason` explaining why
6. Workers will detect this on their next iteration and exit

After workers exit and task is dealt with, you can clear the signal:
```yaml
# Append or update to mark inactive
- signal_id: kill-orchestrator-001
  ...
  active: false
```

**Worker Responsibility:**

At the start of every iteration:
```bash
# Pseudocode
killSignals = readYaml('.agentic/coordination/kill-signals.yaml');
for signal in killSignals {
  if (signal.active &&
      (signal.target == 'all' || signal.target == myAgentId)) {
    announce(status: 'failed', summary: 'Killed by cascade: ' + signal.reason);
    return;  // Exit gracefully
  }
}
```

If you see an active kill signal directed at you or "all":
1. Release all file reservations immediately
2. Announce status: `failed` with reason: "Cascade kill: [reason from signal]"
3. Complete the task with status: `abandoned`
4. Exit (output `<promise>COMPLETE</promise>`)

You do NOT delete kill signals. The orchestrator owns them.

---

## Manual File Editing

These YAML files are meant to be read/written by agents using standard file tools.

To manually check state:
```bash
cat .agentic/coordination/reservations.yaml
cat .agentic/coordination/announcements.yaml
cat .agentic/coordination/kill-signals.yaml
```

To manually fix a stuck reservation (emergency only):
```bash
# Edit directly (not recommended)
nano .agentic/coordination/reservations.yaml
# Remove the stuck entry
```

**Warning:** Manual edits can cause race conditions if agents are concurrently writing. Only edit if you're certain no agents are running.

---

## Concurrency Safety

The `.lock` file in `.agentic/coordination/.lock` is an implementation detail. Ignore it. It's used internally by the CLI (if we ever add coordination support) to prevent race conditions during concurrent reads/writes.

For now, since agents use standard file I/O (Read/Write/Edit tools), expect eventual consistency. If two agents write to the same YAML file simultaneously, the last write wins. Avoid this by:
1. Workers reserve files before editing
2. Orchestrator waits for all announcements before completing
3. Keep YAML files small (< 100 entries each)

---

## Schema Validation

No automatic validation is performed. Agents are responsible for writing valid YAML. If a file becomes corrupt:
1. Orchestrator should detect no announcements for 30+ min
2. Orchestrator cascade kills
3. Workers clean up and exit
4. Fix the coordination files manually
5. Restart the orchestrator

---

## Cleanup (End of Task)

After a task completes:
1. Workers' file reservations auto-expire after 10 min (don't need cleanup)
2. Announcements accumulate (no auto-cleanup) — you may want to archive old announcements manually if the file grows large
3. Kill signals marked `inactive = false` should be left alone

If coordination files grow large (1000+ entries), consider:
- Archiving old announcements to `.agentic/coordination/_archive/announcements-{date}.yaml`
- Purging expired reservations and inactive kill signals

---

## Common Mistakes

| ❌ Mistake | ✅ Fix |
|-----------|--------|
| Forget to release file reservation | Reservation TTL will auto-expire after 10 min. Still, release ASAP. |
| Announce before all work is done | Only announce when ALL acceptance criteria pass. |
| Workers edit same file without reservation | Use reservations. Check for conflicts before editing. |
| Orchestrator doesn't poll announcements | Poll every 10 sec. Set 30-min timeout. |
| Kill signal never checked | Check at start of EVERY iteration. |
| Manual edit to YAML while agents running | Don't. Only edit when all agents are stopped. |

---

## Example: Full Coordination Flow

```
[T+0] Orchestrator claims TASK-500-1
[T+0] Orchestrator spawns Worker-A, Worker-B
[T+0] Worker-A reserves src/auth/session.go (expires T+10m)
[T+0] Worker-B reserves src/auth/jwt.go (expires T+10m)
[T+5] Worker-A: AC-1,2 pass. Checkpoint saved. Still working.
[T+10] Worker-A: AC-1,2,3,4 all pass! Release reservation. Announce(complete).
[T+11] Worker-B: AC-1,2 pass. Checkpoint saved. Still working.
[T+15] Worker-B: AC-1,2,3,4 all pass! Release reservation. Announce(complete).
[T+15] Orchestrator polls. Sees 2 announcements (A, B), both complete. Synthesizes.
[T+16] Orchestrator completes TASK-500-1 via CLI.
[T+16] Orchestrator announces to main(complete).
[T+16] Both Worker-A and Worker-B exit with <promise>COMPLETE</promise>
```

Clean, predictable, decentralized coordination.
