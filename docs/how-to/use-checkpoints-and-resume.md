# Checkpoint and Resume for Long-Running Tasks

## Overview

Checkpoint/resume functionality allows you to pause and resume long-running task execution without losing progress. This is especially useful when:

- Approaching token limits (200K for Claude)
- Need to pause work and resume later
- Want to save progress after significant milestones
- Experiencing interruptions or errors

## How It Works

```
┌──────────────────────────────────────────────────────────┐
│              CHECKPOINT FLOW                             │
├──────────────────────────────────────────────────────────┤
│                                                           │
│  Iteration 1  →  ⚠️  Not checkpoint yet                  │
│  Iteration 2  →  ⚠️  Not checkpoint yet                  │
│  Iteration 3  →  ⚠️  Not checkpoint yet                  │
│  Iteration 4  →  ⚠️  Not checkpoint yet                  │
│  Iteration 5  →  💾 CHECKPOINT (every 5 iterations)      │
│                                                           │
│  [50% tokens]  →  💾 CHECKPOINT (50% threshold)          │
│  [75% tokens]  →  💾 CHECKPOINT (75% threshold)          │
│  [90% tokens]  →  💾 CHECKPOINT + ⚠️  WARNING            │
│                                                           │
│  Task complete →  🗑️  Clean up checkpoints               │
│                                                           │
└──────────────────────────────────────────────────────────┘
```

## Automatic Checkpointing

Checkpoints are created automatically when:

1. **Every 5 iterations** - Regular progress saves
2. **50% of token limit** - First threshold warning
3. **75% of token limit** - Second threshold warning
4. **90% of token limit** - Final warning before limit

### Token Limits by Agent

| Agent | Token Limit | 50% | 75% | 90% |
|-------|-------------|-----|-----|-----|
| Claude 3.5 | 200,000 | 100K | 150K | 180K |
| GPT-4 | 128,000 | 64K | 96K | 115K |
| Gemini Pro | 1,000,000 | 500K | 750K | 900K |
| Copilot | 32,000 | 16K | 24K | 29K |

## Example Usage

### Automatic Checkpoint During Execution

```bash
$ agentic-agent --agent claude autopilot start --execute-agent

╔═══════════════════════════════════════════════════════════╗
║           Agentic Agent - Autopilot Mode                 ║
╚═══════════════════════════════════════════════════════════╝

🤖 Agent execution: ENABLED (claude-code)
🔄 Max iterations: 10

--- Iteration 1/10 ---
🤖 Executing claude-code agent...
  ✅ Agent completed (tokens: 15000, total: 15000)
  ⚠️  Criteria not met: ["Linter clean"]
  📊 Progress: 2/3 criteria met

--- Iteration 2/10 ---
🤖 Executing claude-code agent...
  ✅ Agent completed (tokens: 12000, total: 27000)
  ⚠️  Criteria not met: ["Linter clean"]
  📊 Progress: 2/3 criteria met

...

--- Iteration 5/10 ---
🤖 Executing claude-code agent...
  ✅ Agent completed (tokens: 18000, total: 89000)
  💾 Checkpoint saved (iteration 5, 66.7% complete)
  ⚠️  Criteria not met: ["Linter clean"]
  📊 Progress: 2/3 criteria met

...

--- Iteration 10/10 ---
🤖 Executing claude-code agent...
  ✅ Agent completed (tokens: 22000, total: 178000)
  ⚠️  Token usage at 89.0% of limit (178000/200000)
  💾 Checkpoint saved (iteration 10, 100.0% complete)
  ✅ All acceptance criteria met!
  ✅ Task TASK-001 completed successfully
```

### Resuming from Checkpoint

If autopilot is interrupted or stopped, it automatically resumes from the last checkpoint:

```bash
$ agentic-agent --agent claude autopilot start --execute-agent

--- Iteration 1/10 ---
Next task: [TASK-001] Add health check endpoint

🤖 Executing claude-code agent...
📌 Resuming from checkpoint (iteration 5, 89000 tokens used)
  ✅ Agent completed (tokens: 15000, total: 104000)
  💾 Checkpoint saved (iteration 6, 66.7% complete)
```

### Approaching Token Limit

```bash
--- Iteration 15/20 ---
🤖 Executing claude-code agent...
  ✅ Agent completed (tokens: 18000, total: 162000)
  ⚠️  Token usage at 81.0% of limit (162000/200000)
  💾 Checkpoint saved (iteration 15, 75.0% complete)

--- Iteration 18/20 ---
🤖 Executing claude-code agent...
  ✅ Agent completed (tokens: 20000, total: 182000)
  ⚠️  Token usage at 91.0% of limit (182000/200000)
  🛑 Approaching token limit! Consider pausing and resuming later.
  💾 Checkpoint saved (iteration 18, 85.0% complete)
```

## Checkpoint Storage

Checkpoints are stored in `.agentic/checkpoints/`:

```
.agentic/checkpoints/
├── TASK-001-001.json          # Iteration 1 checkpoint
├── TASK-001-005.json          # Iteration 5 checkpoint
├── TASK-001-010.json          # Iteration 10 checkpoint
├── TASK-001-latest.json       # Latest checkpoint (symlink/copy)
├── TASK-002-005.json
└── TASK-002-latest.json
```

### Checkpoint Contents

```json
{
  "task_id": "TASK-001",
  "iteration": 5,
  "tokens_used": 89000,
  "created_at": "2026-02-13T15:30:00Z",
  "agent": "claude-code",
  "output": "Implemented health check endpoint...",
  "criteria_met": [
    "Tests pass",
    "GET /health returns 200"
  ],
  "criteria_left": [
    "Response includes version number"
  ],
  "files_modified": [
    "internal/api/health.go",
    "internal/api/health_test.go"
  ],
  "learnings": [],
  "notes": "Iteration 5: 2/3 criteria met"
}
```

## Manual Checkpoint Management

### List Checkpoints

```bash
$ ls -la .agentic/checkpoints/

TASK-001-001.json  (15KB, 2026-02-13 14:00)
TASK-001-005.json  (18KB, 2026-02-13 14:15)
TASK-001-010.json  (22KB, 2026-02-13 14:30)
TASK-001-latest.json -> TASK-001-010.json
```

### View Checkpoint

```bash
$ cat .agentic/checkpoints/TASK-001-latest.json | jq .

{
  "task_id": "TASK-001",
  "iteration": 10,
  "tokens_used": 178000,
  "created_at": "2026-02-13T14:30:00Z",
  "agent": "claude-code",
  "criteria_met": ["Tests pass", "API works", "Linter clean"],
  "criteria_left": [],
  "notes": "Iteration 10: 3/3 criteria met"
}
```

### Delete Checkpoints

Checkpoints are automatically deleted when a task completes successfully. To manually clean up:

```bash
# Delete all checkpoints for a task
rm .agentic/checkpoints/TASK-001-*.json

# Delete specific iteration
rm .agentic/checkpoints/TASK-001-005.json
```

## Integration with Ralph Loop

Ralph loop also supports checkpoints when using the `/ralph-loop` skill:

```
You: /ralph-loop

AI: I'm using the ralph-loop skill...

─── Iteration 5 ───
Implementing fixes...
Running tests...

💾 Checkpoint saved (5th iteration)
  Token usage: 85,000 / 200,000 (42.5%)
  Progress: 2/3 criteria met

─── Iteration 10 ───
...

💾 Checkpoint saved (10th iteration)
  Token usage: 165,000 / 200,000 (82.5%)
  Progress: 3/3 criteria met ✅

<promise>TASK COMPLETE</promise>
```

## Best Practices

### 1. Monitor Token Usage

Keep an eye on token usage warnings:

```
⚠️  Token usage at 81.0% of limit  → Continue carefully
⚠️  Token usage at 91.0% of limit  → Consider pausing
🛑 Approaching token limit!        → Pause and resume
```

### 2. Pause and Resume Strategy

If you hit token limits:

```bash
# 1. Stop current autopilot (Ctrl+C)
^C
Autopilot cancelled.

# 2. Check latest checkpoint
cat .agentic/checkpoints/TASK-001-latest.json | jq '.tokens_used, .criteria_left'

# 3. Resume in a new session
agentic-agent --agent claude autopilot start --execute-agent --max-iterations 5
```

### 3. Checkpoint Cleanup

Completed tasks automatically clean up checkpoints. For abandoned tasks:

```bash
# Check for old checkpoints
find .agentic/checkpoints -name "*.json" -mtime +30

# Clean up manually
rm .agentic/checkpoints/TASK-OLD-*.json
```

### 4. Track Progress

Use checkpoint data to track progress:

```bash
# Get progress for all tasks
for file in .agentic/checkpoints/*-latest.json; do
  task=$(jq -r '.task_id' "$file")
  progress=$(jq -r '.notes' "$file")
  echo "$task: $progress"
done
```

## Troubleshooting

### Problem: Checkpoint Not Loading

**Symptom:** Autopilot doesn't resume from checkpoint

**Solution:**
```bash
# Check if checkpoint file exists
ls -la .agentic/checkpoints/TASK-001-latest.json

# Verify checkpoint is valid JSON
cat .agentic/checkpoints/TASK-001-latest.json | jq .

# Check file permissions
chmod 644 .agentic/checkpoints/TASK-001-latest.json
```

### Problem: Too Many Checkpoints

**Symptom:** `.agentic/checkpoints/` directory is large

**Solution:**
```bash
# Count checkpoints per task
ls .agentic/checkpoints/*.json | cut -d- -f1-2 | sort | uniq -c

# Keep only latest for each task
for task in $(ls .agentic/checkpoints/*-latest.json); do
  task_id=$(basename "$task" | sed 's/-latest.json//')
  # Delete non-latest checkpoints
  find .agentic/checkpoints -name "${task_id}-[0-9]*.json" -type f -delete
done
```

### Problem: Checkpoint Out of Sync

**Symptom:** Iteration count doesn't match

**Solution:**
```bash
# Reset checkpoint for task
rm .agentic/checkpoints/TASK-001-*.json

# Start fresh
agentic-agent --agent claude autopilot start --execute-agent
```

## Advanced Usage

### Configurable Checkpoint Thresholds

You can customize when checkpoints are created by configuring the checkpoint settings in `agnostic-agent.yaml`:

```yaml
# agnostic-agent.yaml
checkpoint:
  # Checkpoint every N iterations (default: 5)
  iteration_interval: 5

  # Checkpoint at these token usage percentages (default: [0.5, 0.75, 0.9])
  # Values are decimal percentages (0.5 = 50%, 0.75 = 75%, 0.9 = 90%)
  token_thresholds: [0.5, 0.75, 0.9]
```

**Examples:**

```yaml
# More frequent checkpoints
checkpoint:
  iteration_interval: 3           # Every 3 iterations instead of 5
  token_thresholds: [0.25, 0.5, 0.75, 0.9]  # More milestones

# Less frequent checkpoints (for agents with higher token limits)
checkpoint:
  iteration_interval: 10
  token_thresholds: [0.8, 0.95]   # Only near the end

# Aggressive checkpointing (for critical work)
checkpoint:
  iteration_interval: 2
  token_thresholds: [0.1, 0.25, 0.5, 0.75, 0.9]
```

**Default behavior** (if not configured):
- Checkpoint every 5 iterations
- Checkpoint at 50%, 75%, 90% of token limit

### Custom Token Limits

For different agents, you can customize token limits (future feature):

```yaml
# agnostic-agent.yaml
agents:
  claude-code:
    token_limit: 200000
  gemini:
    token_limit: 1000000
  gpt-4:
    token_limit: 128000
```

## Summary

**Benefits:**
- ✅ Never lose progress on long-running tasks
- ✅ Pause and resume anytime
- ✅ Automatic token limit monitoring
- ✅ Progress tracking across iterations
- ✅ Clean up on successful completion

**When Checkpoints Are Created:**
- Every 5 iterations
- At 50%, 75%, 90% token usage
- Before hitting token limits

**Storage:**
- Location: `.agentic/checkpoints/`
- Format: JSON
- Auto-cleanup on task completion

Checkpoints enable reliable, resumable execution for even the longest tasks!
