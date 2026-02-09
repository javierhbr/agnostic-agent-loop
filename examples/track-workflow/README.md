# Track Workflow: Idea to Implementation

End-to-end walkthrough of the track lifecycle — from a vague idea to a completed feature with validated spec, phased plan, and atomic tasks.

---

## What You'll Learn

- Initialize a track with brainstorming scaffolding
- Use brainstorm.md as a dialogue script with your AI agent
- Validate spec completeness with `track refine`
- Activate a track to generate a plan and decompose tasks
- Use `plan show`, `plan next`, and `plan mark` to track progress
- Monitor project status with the `status` dashboard
- Complete and archive the track

---

## 0. Setup

```bash
# From the project root
go build -o examples/track-workflow/agentic-agent ./cmd/agentic-agent
cd examples/track-workflow

# Start fresh
rm -rf .agentic agnostic-agent.yaml

# Initialize the project
./agentic-agent init --name "track-demo"
```

---

## 1. Initialize a Track

A track groups a spec, plan, and tasks into a single work unit. Start with a name and optional metadata:

```bash
# Flag mode
./agentic-agent track init "Search Feature" \
  --type feature \
  --purpose "Let users search books by title, author, or genre" \
  --success "Users can search and get results within 200ms"

# Or interactive mode (no flags)
./agentic-agent track init
```

Output:

```text
✓ Created track: search-feature
  → Brainstorm: .agentic/tracks/search-feature/brainstorm.md
  → Spec:       .agentic/tracks/search-feature/spec.md
  → Plan:       .agentic/tracks/search-feature/plan.md
  → Type:       feature
  → Status:     ideation

→ Use brainstorm.md as a dialogue script with your AI agent.
→ Then run: agentic-agent track refine search-feature
```

Three files are scaffolded:

| File             | Purpose                                                                |
|------------------|------------------------------------------------------------------------|
| `brainstorm.md`  | Dialogue script — guides the agent through structured questions        |
| `spec.md`        | Enhanced specification with purpose, constraints, design, requirements |
| `plan.md`        | Empty, generated during activation                                     |

---

## 2. Brainstorm with Your AI Agent

Open the project in your AI agent (Claude Code, Cursor, etc.) and point it at the brainstorm file:

```text
Read .agentic/tracks/search-feature/brainstorm.md and walk me
through the brainstorming phases. Ask one question at a time.
```

The brainstorm template has four phases:

**Phase 1: Understanding** — Purpose, constraints, success criteria, existing context.

**Phase 2: Exploring Approaches** — The agent proposes 2-3 approaches with trade-offs, then you pick one.

**Phase 3: Design Presentation** — Architecture, components, data flow, error handling, testing. Presented in 200-300 word sections for validation.

**Phase 4: Finalization** — Agent writes the validated design into spec.md.

As the conversation progresses, the agent fills in the "Answers Captured" section of brainstorm.md.

---

## 3. Refine the Spec

After brainstorming, check that the spec has all required sections:

```bash
./agentic-agent track refine search-feature
```

Output:

```text
Spec Completeness: Search Feature

  ✓ purpose
  ✓ constraints
  ✓ success criteria
  ✗ design (missing)
  ✗ requirements (missing)
  ○ alternatives considered (optional)

  ✗ 2 section(s) need work.
```

Fill in the missing sections (with your agent or manually), then re-run until complete:

```bash
./agentic-agent track refine search-feature
# ✓ Spec is complete. Run: agentic-agent track activate search-feature
```

---

## 4. Activate the Track

Activation validates the spec, generates a phased plan from it, and optionally decomposes the plan into tasks:

```bash
# Generate plan + create tasks in one step
./agentic-agent track activate search-feature --decompose
```

Output:

```text
✓ Activated track: search-feature
  ✓ Plan generated from spec
  ✓ Status: active
  ✓ Created 4 tasks in backlog
    • TASK-1738000001: Set up search index schema
    • TASK-1738000002: Implement search query parser
    • TASK-1738000003: Add search HTTP endpoint
    • TASK-1738000004: Write search integration tests
```

Without `--decompose`, only the plan is generated. You can decompose later or create tasks manually.

---

## 5. View the Plan

The plan is a markdown file with phased implementation steps using checkbox markers:

```bash
# Show plan with progress indicators
./agentic-agent plan show --track search-feature
```

Output:

```text
Search Feature Plan (0/8 done)

  Phase 1: Foundation
    ○ Set up search index schema
    ○ Define search query types

  Phase 2: Core Logic
    ○ Implement full-text search
    ○ Add fuzzy matching
    ○ Implement result ranking

  Phase 3: Integration
    ○ Add HTTP endpoint
    ○ Wire up middleware
    ○ Write integration tests
```

Find the next step to work on:

```bash
./agentic-agent plan next --track search-feature
# Phase 1: Foundation →
# ❯ Set up search index schema
```

Mark steps as done or in-progress:

```bash
# Mark line 5 as in-progress
./agentic-agent plan mark .agentic/tracks/search-feature/plan.md 5 wip

# Mark line 5 as done
./agentic-agent plan mark .agentic/tracks/search-feature/plan.md 5 done
```

---

## 6. Work Through Tasks

Claim and complete tasks from the backlog:

```bash
# Check project status first
./agentic-agent status
```

Output:

```text
track-demo Status

  [████████░░░░░░░░░░░░░░░░░░░░░░]  0% complete

  ○ 4 backlog  ◐ 0 in progress  ✓ 0 done

  Next Up
  ❯ TASK-1738000001: Set up search index schema
```

Work on the first task:

```bash
# Claim the task
./agentic-agent task claim TASK-1738000001

# Generate context for the target directory
./agentic-agent context generate internal/search

# Show task details
./agentic-agent task show TASK-1738000001

# ... implement the code ...

# Validate
./agentic-agent validate

# Complete the task
./agentic-agent task complete TASK-1738000001
```

Or use the interactive workflow:

```bash
./agentic-agent work
```

The TUI guides you through: select task, claim, generate context, implement, validate, complete.

Check status again:

```bash
./agentic-agent status
```

```text
track-demo Status

  [████████░░░░░░░░░░░░░░░░░░░░░░]  25% complete

  ○ 3 backlog  ◐ 0 in progress  ✓ 1 done

  Next Up
  ❯ TASK-1738000002: Implement search query parser

  Recent Activity
  • Feb 09  Completed TASK-1738000001
```

---

## 7. Complete and Archive the Track

Once all tasks are done:

```bash
# Verify nothing is left
./agentic-agent task list
./agentic-agent status

# View the track
./agentic-agent track show search-feature

# Archive the track
./agentic-agent track archive search-feature
```

---

## Quick Reference

| Step                 | Command                                                              |
|----------------------|----------------------------------------------------------------------|
| Create track         | `track init "Name" --type feature --purpose "..." --success "..."`   |
| List tracks          | `track list`                                                         |
| Show track           | `track show <id>`                                                    |
| Check spec           | `track refine <id>`                                                  |
| Activate + decompose | `track activate <id> --decompose`                                    |
| View plan progress   | `plan show --track <id>`                                             |
| Next plan step       | `plan next --track <id>`                                             |
| Mark plan step       | `plan mark <path> <line> done`                                       |
| Project status       | `status`                                                             |
| Archive track        | `track archive <id>`                                                 |

## Lifecycle

```text
track init ──→ brainstorm ──→ track refine ──→ track activate --decompose
    │                                                │
    │                              plan show / plan next / plan mark
    │                                                │
    +─── track archive ←── status ←── task claim + implement + complete
```
