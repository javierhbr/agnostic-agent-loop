# Autopilot vs Ralph Loop - When to Use What

## Quick Comparison

| Feature | Autopilot (CLI) | Ralph Loop (AI Chat) |
|---------|-----------------|----------------------|
| **Execution** | Terminal command | AI chat skill |
| **Visibility** | Terminal output | Interactive conversation |
| **Iteration** | Single pass per task | Loops until complete |
| **Agent** | Optional `--execute-agent` | Always AI-driven |
| **Best For** | Batch processing | Interactive development |
| **Output** | Structured logs | Conversational progress |

## Autopilot (Terminal)

### When to Use

- Processing multiple tasks in batch
- CI/CD pipelines
- Automation scripts
- Non-interactive environments

### Example Usage

```bash
# Run autopilot with agent execution
agentic-agent --agent claude autopilot start --execute-agent --max-iterations 5
```

**Output:**

```
╔═══════════════════════════════════════════════════════════╗
║           Agentic Agent - Autopilot Mode                 ║
╚═══════════════════════════════════════════════════════════╝

🤖 Agent execution: ENABLED (claude-code)
   Tasks will be executed by AI agent automatically
🔄 Max iterations: 5

--- Iteration 1/5 ---
Next task: [TASK-001] Add health check endpoint
Readiness checks:
  ✓ Specs found (1)
  ✓ Scope defined (1 directory)
  ✓ No blockers

Claimed task TASK-001
  Generated context for internal/api
  Context bundle built (12KB)

🤖 Executing claude-code agent...
  ✅ Agent completed (tokens: 3456)
  ✅ All acceptance criteria met!
  ✅ Task TASK-001 completed successfully

--- Iteration 2/5 ---
All tasks complete. Autopilot finished.
```

### Without Agent Execution

```bash
# Prepare tasks without executing
agentic-agent autopilot start --max-iterations 5
```

**Output:**

```
╔═══════════════════════════════════════════════════════════╗
║           Agentic Agent - Autopilot Mode                 ║
╚═══════════════════════════════════════════════════════════╝

📋 Agent execution: DISABLED
   Tasks will be prepared but not executed

💡 Tip: Use --execute-agent to enable AI agent execution
   Or use /ralph-loop in AI chat for interactive iteration
🔄 Max iterations: 5

--- Iteration 1/5 ---
Next task: [TASK-001] Add health check endpoint
Task TASK-001 is ready for agent execution.

--- Iteration 2/5 ---
...
```

## Ralph Loop (AI Chat)

### When to Use

- Interactive development
- Learning and experimentation
- Tasks requiring human judgment
- When you want visibility into iteration process

### Example Usage

**In AI Chat (Claude Code, Cursor, etc.):**

```
You: /ralph-loop

AI: I'm using the ralph-loop skill to iteratively complete this task.

─── Step 1: Get Task ───
$ agentic-agent task continue

Working on: TASK-001 - Add health check endpoint

Acceptance criteria:
  - Tests pass
  - GET /health returns 200
  - Response includes version number

Specs: health-check-spec

─── Step 2: Read Specs ───
Reading .agentic/spec/health-check-spec.md...

Requirements:
- Endpoint: GET /health
- Return JSON with status, version, uptime

─── Iteration 1 ───
Creating health check handler...
Writing tests...

$ go test ./internal/api/...
FAIL: TestHealthCheck not implemented

Status: 0/3 criteria met
Continuing...

─── Iteration 2 ───
Implementing handler function...

$ go test ./internal/api/...
PASS: 5/5 tests

$ curl http://localhost:8080/health
{"status":"ok","version":"1.0.0","uptime":42}

Status: 3/3 criteria met ✅

$ agentic-agent task complete TASK-001

<promise>TASK COMPLETE</promise>
```

## Detecting CLI Command in Chat

If a user pastes the autopilot CLI command into AI chat:

```
User: agentic-agent --agent claude autopilot start --execute-agent --max-iterations 1
```

**AI Response:**

```
I see you've pasted an autopilot CLI command. In AI chat, we use the /ralph-loop
skill instead, which provides the same iterative behavior with better visibility.

Let me run the ralph-loop for you:

─── Step 1: Get Task ───
...
```

## Choosing the Right Approach

### Use Autopilot When

✅ Running in CI/CD pipeline
✅ Processing many tasks automatically
✅ Scripting and automation
✅ Terminal-only environment
✅ Want structured log output

**Example:**
```bash
#!/bin/bash
# CI pipeline script
agentic-agent --agent claude autopilot start \
  --execute-agent \
  --max-iterations 10
```

### Use Ralph Loop When

✅ Working interactively in IDE
✅ Want to see iteration progress
✅ Learning or experimenting
✅ Task requires human judgment
✅ Want conversational feedback

**Example:**
```
# In Claude Code while developing
You: /ralph-loop

# [Interactive iteration with visible progress]
```

## Combining Both Approaches

### Workflow 1: Autopilot Prep + Ralph Loop Execute

```bash
# Terminal: Prepare tasks
agentic-agent autopilot start --max-iterations 5 --dry-run

# Review what would be processed
# Then claim one task manually
agentic-agent task claim TASK-001

# AI Chat: Execute interactively
/ralph-loop
```

### Workflow 2: Ralph Loop Dev + Autopilot Deploy

```
# AI Chat: Develop first task iteratively
/ralph-loop
# [Task TASK-001 complete]

# Terminal: Batch process remaining tasks
agentic-agent --agent claude autopilot start --execute-agent
```

## Configuration

### For Autopilot

```yaml
# agnostic-agent.yaml
active_agent: claude-code

agents:
  claude-code:
    auto_setup: true
```

```bash
# Or via flags
agentic-agent --agent claude autopilot start --execute-agent
```

### For Ralph Loop

No configuration needed! Just:

1. Claim a task
2. Type `/ralph-loop` in AI chat

## Error Handling

### Autopilot

```bash
$ agentic-agent autopilot start --execute-agent

Error: No tasks in backlog
Tip: Create tasks first or run: agentic-agent task list
```

### Ralph Loop

```
You: /ralph-loop

AI: Let me check for active tasks...
$ agentic-agent task continue

No in-progress tasks found.

Would you like me to:
1. List available tasks (agentic-agent task list)
2. Create a new task
3. Claim a specific task ID

Which would you prefer?
```

## Summary

**Autopilot = Automated batch processing**
- Fast
- Structured
- Non-interactive
- Good for CI/CD

**Ralph Loop = Interactive iteration**
- Visible progress
- Conversational
- Educational
- Good for development

Both achieve the same goal (task completion) but with different interaction models. Choose based on your context:
- **Terminal/automation** → Autopilot
- **AI chat/IDE** → Ralph Loop
