# TDD Workflow: RED / GREEN / REFACTOR

Use `work --follow-tdd` to decompose a task into three phased sub-tasks and guide your AI agent through strict test-driven development.

---

## What You'll Learn

- Install the TDD skill pack for your AI agent tool
- Use `--follow-tdd` to decompose a task into RED/GREEN/REFACTOR phases
- Work through each phase with agent-specific skill files

---

## 0. Setup

```bash
# From the project root
go build -o examples/tdd/agentic-agent ./cmd/agentic-agent
cd examples/tdd

# Initialize a project
rm -rf .agentic agnostic-agent.yaml
./agentic-agent init --name "tdd-demo"
```

---

## 1. Install the TDD Skill Pack

The skill pack gives your AI agent phase-specific instructions (what to do in RED, GREEN, REFACTOR).

```bash
# Install for Claude Code (project-level)
./agentic-agent skills install tdd --tool claude-code
```

Output:

```text
Skill Pack Installed

  ✓ Installed "tdd" for claude-code

  • .claude/skills/tdd/red.md
  • .claude/skills/tdd/green.md
  • .claude/skills/tdd/refactor.md
```

Install for other tools as needed:

```bash
./agentic-agent skills install tdd --tool cursor
./agentic-agent skills install tdd --tool gemini
```

Verify:

```bash
./agentic-agent skills list
```

---

## 2. Create a Task

```bash
./agentic-agent task create \
  --title "Add email validation utility" \
  --description "Validate email format, check for disposable domains, normalize" \
  --spec-refs "validation-spec.md"
```

Note the task ID from the output (e.g., `TASK-1738000001`).

---

## 3. Run with TDD

```bash
./agentic-agent work --task TASK-1738000001 --follow-tdd
```

Output:

```text
✓ Claimed task TASK-1738000001
TDD workflow enabled (skill pack found for claude-code)
Created TDD sub-tasks:
  - TASK-1738000001-red: Write failing tests for email validation
  - TASK-1738000001-green: Implement minimal code to pass tests
  - TASK-1738000001-refactor: Improve code quality

TDD sub-tasks created. Work through RED → GREEN → REFACTOR phases.
```

---

## 4. Phase 1: RED — Write Failing Tests

Claim the RED sub-task and tell your agent to write tests only:

```bash
./agentic-agent task claim TASK-1738000001-red
```

The agent reads `.claude/skills/tdd/red.md` which instructs it to:

- Write test cases covering expected behavior and edge cases
- Run the tests and confirm they fail
- Do NOT write any implementation code

```bash
# After writing tests
./agentic-agent task complete TASK-1738000001-red
```

---

## 5. Phase 2: GREEN — Minimal Implementation

```bash
./agentic-agent task claim TASK-1738000001-green
```

The agent reads `.claude/skills/tdd/green.md` which instructs it to:

- Write the minimum code to make all tests pass
- Do NOT optimize, refactor, or add features beyond what tests require
- Run tests and confirm they all pass

```bash
./agentic-agent task complete TASK-1738000001-green
```

---

## 6. Phase 3: REFACTOR — Improve Quality

```bash
./agentic-agent task claim TASK-1738000001-refactor
```

The agent reads `.claude/skills/tdd/refactor.md` which instructs it to:

- Improve naming, structure, and readability
- Extract helpers if needed
- Run tests after every change to confirm nothing breaks

```bash
./agentic-agent task complete TASK-1738000001-refactor
```

---

## 7. Verify Completion

```bash
./agentic-agent status
./agentic-agent task list
```

---

## Supported Tools

| Tool        | Skill Directory          |
|-------------|--------------------------|
| claude-code | `.claude/skills/tdd/`    |
| cursor      | `.cursor/skills/tdd/`    |
| gemini      | `.gemini/skills/tdd/`    |
| windsurf    | `.windsurf/skills/tdd/`  |
| antigravity | `.agent/skills/tdd/`     |
| codex       | `.codex/skills/tdd/`     |

## Quick Reference

| Action           | Command                                          |
|------------------|--------------------------------------------------|
| Install TDD pack | `skills install tdd --tool claude-code`          |
| Install globally | `skills install tdd --tool claude-code --global` |
| Run TDD workflow | `work --task <id> --follow-tdd`                  |
| List packs       | `skills list`                                    |
