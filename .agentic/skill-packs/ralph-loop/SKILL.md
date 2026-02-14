---
name: ralph-loop
description: Execute Ralph Wiggum iterative convergence loop for task completion. Use when user says "/ralph-loop", pastes "agentic-agent autopilot" command, or asks to implement a task iteratively until complete.
---

# Ralph Wiggum Loop - Iterative Task Completion

## Overview

Ralph Wiggum methodology: iterate until task converges to completion. Keep trying, refining, and improving until all acceptance criteria are met.

**Announcement required:** "I'm using the ralph-loop skill to iteratively complete this task."

## Detecting User Intent

**Trigger this skill when user:**
1. Types `/ralph-loop` explicitly
2. Pastes CLI command like: `agentic-agent --agent XYZ autopilot start --execute-agent`
3. Says "implement this task iteratively" or similar
4. Asks to "keep trying until it works"

**If user pastes a CLI command in chat:**
```
User: agentic-agent --agent claude autopilot start --execute-agent --max-iterations 1

AI Response:
I see you've pasted an autopilot CLI command. In AI chat, we use the /ralph-loop
skill instead, which provides the same iterative behavior with better visibility.

Let me run the ralph-loop for you:

[Proceeds with ralph-loop methodology]
```

## The Ralph Loop Process

```
┌─────────────────────────────────────────────────────────────┐
│                    RALPH LOOP CYCLE                          │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  1. Get task details        → agentic-agent task continue   │
│  2. Read all specs          → Files from task.SpecRefs      │
│  3. Check acceptance        → Criteria to meet               │
│  4. Implement solution      → Write code                     │
│  5. Verify criteria         → Run tests, check output       │
│  6. Check convergence:                                       │
│     ├─ All criteria met?  → Complete and exit               │
│     └─ Not yet?           → Iterate (back to step 4)        │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

## Step-by-Step Instructions

### Step 1: Get Active Task

```bash
# Get the current in-progress task
agentic-agent task continue
```

**Parse the output to extract:**
- Task ID
- Title
- Description
- Acceptance criteria (CRITICAL - these define "done")
- Spec references
- Scope directories

**If no task claimed:**
```bash
# List available tasks
agentic-agent task list

# Claim a specific task
agentic-agent task claim TASK-XXX
```

### Step 2: Read Context

**Read ALL spec files referenced in the task:**

```bash
# Task will show spec references like:
# Specs:
#   - auth/requirements.md
#   - auth/api-design.md
```

Read each spec file completely. These define WHAT to build.

**Example:**
```bash
# If task says: Specs: auth/requirements.md
# Then read: .agentic/spec/auth/requirements.md
```

### Step 3: Understand Acceptance Criteria

**Acceptance criteria = Definition of "done"**

Example criteria:
```yaml
acceptance:
  - "Tests pass"
  - "API returns 200 for valid credentials"
  - "JWT token is generated and valid"
  - "Linter shows no errors"
```

**You MUST meet ALL criteria before completing.**

### Step 4: Iteration Loop

**For each iteration (max 10):**

#### 4a. Implement/Improve

Based on:
- Spec requirements
- Acceptance criteria
- Previous iteration learnings

Write or modify code in the scope directories.

#### 4b. Verify Criteria

**Check EACH acceptance criterion:**

```bash
# Example verifications:

# Criterion: "Tests pass"
go test ./internal/auth/... -v

# Criterion: "API returns 200"
curl -X POST http://localhost:8080/api/login \
  -d '{"user":"test","pass":"test"}' | jq .

# Criterion: "Linter shows no errors"
golangci-lint run ./internal/auth/...

# Criterion: "Build succeeds"
go build ./cmd/...
```

**Track results:**
```
Iteration 3 Status:
✓ Tests pass (100% coverage)
✓ API returns 200 for valid credentials
✗ JWT token validation failing
✗ Linter shows 2 errors

Continue iterating...
```

#### 4c. Check Convergence

**If ALL criteria met:**
```bash
# Complete the task
agentic-agent task complete TASK-XXX \
  --learnings "Implemented JWT auth with bcrypt" \
  --learnings "Added middleware for token validation"
```

**Output completion signal:**
```
<promise>TASK COMPLETE</promise>
```

**If NOT all criteria met:**
- Identify what's missing
- Make targeted improvements
- Continue to next iteration

### Step 5: Report Progress

**After each iteration, report:**

```
─── Iteration 3/10 ───

Changes made:
- Fixed JWT signing algorithm (HS256 → RS256)
- Added token expiration check in middleware

Verification results:
✓ Tests pass (15/15)
✓ API returns 200 for valid credentials
✓ JWT token is generated and valid
✗ Linter shows 1 error: unused variable

Next steps:
- Remove unused variable
- Re-run linter
```

## Convergence Detection

**Stop iterating when ONE of these occurs:**

1. **Success:** All acceptance criteria met
   ```
   ✓ All 4 criteria passed
   <promise>TASK COMPLETE</promise>
   ```

2. **Max iterations:** Reached 10 iterations without completion
   ```
   ⚠ Reached max iterations (10)
   Task remains in-progress
   Recommend: Review acceptance criteria or break into smaller tasks
   ```

3. **Stuck:** Same errors for 3+ consecutive iterations
   ```
   ⚠ No progress in last 3 iterations
   Recommend: Ask for help or review approach
   ```

## Example: Full Ralph Loop

```
User: /ralph-loop

Claude: I'm using the ralph-loop skill to iteratively complete this task.

─── Step 1: Get Task ───
$ agentic-agent task continue

Working on: TASK-123 - Implement user login

Acceptance criteria:
  - Tests pass
  - API returns 200 for valid credentials
  - JWT token generated
  - No linter errors

Specs:
  - auth/requirements.md
  - auth/api-design.md

Scope: internal/auth

─── Step 2: Read Specs ───
Reading .agentic/spec/auth/requirements.md...
Reading .agentic/spec/auth/api-design.md...

Requirements understood:
- JWT-based authentication
- POST /api/login endpoint
- bcrypt password hashing
- 1-hour token expiration

─── Iteration 1 ───
Implementing:
- Creating auth package
- Adding login handler
- Implementing JWT signing

$ go test ./internal/auth/...
FAIL: TestLoginHandler not found

Status: 0/4 criteria met
Continuing...

─── Iteration 2 ───
Improving:
- Added login handler tests
- Implemented bcrypt hashing
- Added JWT token generation

$ go test ./internal/auth/... -v
PASS: 5/5 tests

$ curl -X POST localhost:8080/api/login -d '{"user":"test","pass":"test"}'
{"token":"eyJhbGc...","status":"ok"}

Status: 2/4 criteria met (tests ✓, API ✓)
Continuing...

─── Iteration 3 ───
Improving:
- Added token validation
- Fixed linter issues

$ golangci-lint run ./internal/auth/...
✓ No issues

$ go test ./internal/auth/...
✓ All tests pass

Verification complete:
✓ Tests pass (5/5)
✓ API returns 200 for valid credentials
✓ JWT token generated and validates
✓ No linter errors

All 4/4 criteria met! 🎉

$ agentic-agent task complete TASK-123 \
  --learnings "Implemented JWT auth with bcrypt hashing" \
  --learnings "Added comprehensive test coverage"

✓ Task TASK-123 completed successfully

<promise>TASK COMPLETE</promise>
```

## Key Principles

### 1. **Be Systematic**
- Always verify ALL criteria each iteration
- Don't assume criteria are met
- Run actual commands, don't simulate

### 2. **Be Incremental**
- Make small, focused improvements each iteration
- Don't try to fix everything at once
- Build on previous iteration's progress

### 3. **Be Honest**
- Report actual verification results
- Don't mark complete unless ALL criteria pass
- Show real command output

### 4. **Be Persistent**
- Keep iterating until convergence
- Learn from each iteration's failures
- Refine approach based on verification results

### 5. **Be Complete**
- Read ALL specs before starting
- Verify ALL acceptance criteria
- Don't skip verification steps

## Commands Reference

```bash
# Task management
agentic-agent task continue              # Get active task
agentic-agent task list                  # List all tasks
agentic-agent task claim <ID>            # Claim a task
agentic-agent task show <ID>             # Show task details
agentic-agent task complete <ID>         # Complete task

# Context
agentic-agent context build --task <ID>  # Build full context

# Validation
agentic-agent validate                   # Run all validators
```

## Tool-Agnostic Design

This skill works in:
- ✅ Claude Code (native file access + Bash)
- ✅ Cursor (file access + terminal)
- ✅ GitHub Copilot (via terminal)
- ✅ Windsurf (file access + terminal)
- ✅ Any AI tool with file and command execution

**Requirements:**
- File read/write capability
- Command execution (Bash/terminal)
- Access to `agentic-agent` CLI

## Anti-Patterns

❌ **DON'T:**
- Skip reading specs
- Assume tests pass without running them
- Mark complete without verifying ALL criteria
- Give up after 1-2 iterations
- Simulate verification results

✅ **DO:**
- Read every spec file
- Run actual verification commands
- Report real output
- Iterate until ALL criteria pass
- Be honest about progress

## Troubleshooting

**Problem: Can't claim task**
```bash
# Check what's in-progress
agentic-agent task list

# Unclaim if needed
agentic-agent task unclaim <ID>

# Then claim
agentic-agent task claim <ID>
```

**Problem: Specs not found**
```bash
# Check spec directories
ls -la .agentic/spec/

# Check config
cat agnostic-agent.yaml | grep spec_dirs
```

**Problem: Stuck in loop (no progress)**
- Review approach - might need different strategy
- Check if acceptance criteria are achievable
- Consider breaking task into smaller sub-tasks
- Ask user for guidance

## Success Criteria

**The ralph-loop is successful when:**
1. ✓ Task claimed and context loaded
2. ✓ All specs read and understood
3. ✓ Code implemented in scope directories
4. ✓ ALL acceptance criteria verified and passing
5. ✓ Task completed with learnings captured
6. ✓ `<promise>TASK COMPLETE</promise>` output

**If max iterations reached without success:**
- Document what criteria remain unmet
- Capture learnings about blockers
- Leave task in-progress for manual review
- Recommend next steps to user
