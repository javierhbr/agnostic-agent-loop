# Run with Ralph: Iterative Task Implementation

Execute openspec tasks using Ralph Wiggum iterative loops. One task per loop, acceptance criteria as the completion signal.

---

## What You'll Learn

- Claim openspec tasks and extract context for Ralph prompts
- Build effective Ralph prompts from `task show` output
- Run `/ralph-loop` with convergence guarantees
- Complete tasks with automatic git tracking
- Process an entire openspec change task by task

---

## 0. Prerequisites

- An openspec change with imported tasks (see [spec-driven-workflow](../spec-driven-workflow/README.md))
- Ralph Wiggum plugin installed (provides `/ralph-loop` and `/cancel-ralph`)

---

## 1. List Available Tasks

```bash
agentic-agent task list --no-interactive
```

```text
--- BACKLOG ---
[TASK-936281-1] [todo-pwa] Set up project structure with Vite, React, TypeScript
[TASK-936281-2] [todo-pwa] Configure TypeScript and ESLint
[TASK-936281-3] [todo-pwa] Set up Tailwind with dark mode configuration
...
```

---

## 2. Claim a Task (BEFORE Ralph)

```bash
agentic-agent task claim TASK-936281-2
```

This records your git branch and timestamp. **Never skip this** — `task complete` needs the claim timestamp to capture commits.

---

## 3. Get Task Details

```bash
agentic-agent task show TASK-936281-2 --no-interactive
```

```text
ID: TASK-936281-2
Title: [todo-pwa] Configure TypeScript and ESLint
Description: Set up TypeScript with strict mode and ESLint with React/TypeScript rules...
Spec Refs:
  - todo-pwa/proposal.md
  - todo-pwa/tasks/02-typescript-eslint.md
Acceptance Criteria:
  - tsconfig.json configured with strict mode enabled
  - ESLint installed and configured for React + TypeScript
  - No ESLint errors on existing code
  - Path aliases configured (@/ for src/)
  - VS Code settings recommended for team consistency
```

Extract three things:
1. **Spec Refs** → files Ralph reads each iteration
2. **Acceptance Criteria** → Ralph's completion conditions
3. **Description** → implementation hints

---

## 4. Build the Ralph Prompt

Template — fill from `task show` output:

```text
You are implementing TASK-936281-2: Configure TypeScript and ESLint.

## Context — read these files first:
- .agentic/openspec/changes/todo-pwa/proposal.md
- .agentic/openspec/changes/todo-pwa/tasks/02-typescript-eslint.md
- .agentic/context/tech-stack.md

## Acceptance Criteria — ALL must pass:
- tsconfig.json configured with strict mode enabled
- ESLint installed and configured for React + TypeScript
- No ESLint errors on existing code
- Path aliases configured (@/ for src/)
- VS Code settings recommended for team consistency

## On each iteration:
1. Read the task detail file for technical notes
2. Implement what's missing
3. Run npx tsc --noEmit && npx eslint . to verify
4. If ALL criteria pass, output <promise>TASK COMPLETE</promise>
```

---

## 5. Launch Ralph

```
/ralph-loop "YOUR PROMPT FROM STEP 4" --max-iterations 10 --completion-promise "TASK COMPLETE"
```

Ralph will:
1. Read the spec files
2. Implement changes
3. Run verification commands
4. Check acceptance criteria
5. Loop until all pass, then output `<promise>TASK COMPLETE</promise>`

---

## 6. Complete the Task

```bash
agentic-agent task complete TASK-936281-2
agentic-agent openspec status todo-pwa
```

```text
Change: todo-pwa
  Total: 22  Done: 1  In Progress: 0  Pending: 21
```

---

## 7. Repeat for Next Task

```bash
agentic-agent task claim TASK-936281-3
agentic-agent task show TASK-936281-3 --no-interactive
# Build ralph prompt from output...
/ralph-loop "..." --max-iterations 10 --completion-promise "TASK COMPLETE"
agentic-agent task complete TASK-936281-3
```

---

## Full Change Lifecycle

```text
agentic-agent task list
    ↓
For each task:
    ├── task claim TASK-ID
    ├── task show TASK-ID → extract criteria + spec refs
    ├── /ralph-loop "prompt" --max-iterations 10
    ├── task complete TASK-ID
    └── openspec status CHANGE-ID
    ↓
agentic-agent openspec complete CHANGE-ID
agentic-agent openspec archive CHANGE-ID
```

---

## Combining with ATDD

For maximum reliability, write acceptance tests (ATDD) before launching Ralph:

```bash
# 1. Claim task
agentic-agent task claim TASK-001

# 2. Write acceptance tests from criteria (ATDD RED phase)
# ... agent writes failing tests ...

# 3. Launch Ralph with tests as verification
/ralph-loop "Implement TASK-001. Acceptance tests in tests/auth.test.ts.
Run npm test each iteration.
When all pass: <promise>TASK COMPLETE</promise>" --max-iterations 10 --completion-promise "TASK COMPLETE"

# 4. Complete
agentic-agent task complete TASK-001
```

---

## Quick Reference

| Step | Command |
|------|---------|
| List tasks | `task list --no-interactive` |
| Claim | `task claim TASK-ID` |
| Show details | `task show TASK-ID --no-interactive` |
| Launch Ralph | `/ralph-loop "prompt" --max-iterations 10 --completion-promise "TASK COMPLETE"` |
| Cancel Ralph | `/cancel-ralph` |
| Complete task | `task complete TASK-ID` |
| Check progress | `openspec status CHANGE-ID` |
| Finish change | `openspec complete CHANGE-ID` |

## Critical Rules

| Rule | Why |
|------|-----|
| One task per loop | Multiple tasks = no convergence |
| Always claim before ralph | Git tracking needs the timestamp |
| Include spec refs in prompt | Ralph loses context between iterations |
| Include `<promise>` tag instruction | Without it, Ralph loops forever |
| Set `--max-iterations` | Safety net (10 is a good default) |
| `task complete` after ralph | Captures commits, updates progress |
