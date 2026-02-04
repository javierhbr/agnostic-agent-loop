
## 23. Automated Framework Validator

The validator transforms rules into enforceable law. Without it, the framework is merely a convention. With it, the framework becomes a system.

### 23.1 Validator Objectives

The validator must produce a definitive yes/no answer for each of these questions:

1. Does the change respect the task's defined scope?
2. Was `context.md` updated when a logical change occurred?
3. Is the task sufficiently small?
4. Were directories modified that lack a context file?
5. Was the session close protocol followed correctly?

### 23.2 Validator Architecture

```
.agentic/validator
  ├─ index.ts
  ├─ rules/
  │   ├─ task-scope.rule.ts
  │   ├─ context-update.rule.ts
  │   ├─ directory-context.rule.ts
  │   ├─ task-size.rule.ts
  │   └─ session-close.rule.ts
  └─ report.ts
```

Pure rule logic only; no agent-specific code.

### 23.3 Validation Rules

#### Rule 1: Directory Context Required

If a directory contains logic (files with `.ts`, `.js`, `.py`, `.go` extensions) and `context.md` does not exist, the validation fails.

```
IF directory contains .ts|.js|.py|.go
AND context.md does not exist
THEN error
```

#### Rule 2: Context Update on Logical Change (CUOC)

This is the most critical rule.

Practical heuristic (sufficient in 90% of cases):

```
IF files_changed_in_dir > 1
OR public function signatures changed
OR new files added
THEN context.md must be modified
```

Error output:

```
ERROR: Logical changes detected in src/token
but src/token/context.md was not updated.
```

#### Rule 3: Task Scope Enforcement

Compares modified files (from git diff) against the task's allowed directories.

```
IF modified_files NOT IN task.allowed_directories
THEN error
```

#### Rule 4: Task Size Limits

```yaml
limits:
  max_files: 5
  max_directories: 2
```

```
IF files_changed > max_files
THEN error: task too large, must decompose
```

#### Rule 5: Session Close Compliance

When a session ends, both must exist:

- `.agentic/context/session-summary-*.md`
- Updated `rolling-summary.md`

```
IF session_closed
AND no session-summary exists
THEN error
```

### 23.4 Validator Output Format

Machine-readable and human-friendly:

```json
{
  "status": "fail",
  "errors": [
    {
      "rule": "CONTEXT_UPDATE_REQUIRED",
      "directory": "src/token",
      "message": "Logical changes without context.md update"
    }
  ]
}
```

### 23.5 Validator Execution Points

- Before accepting sub-agent output
- Before closing a session
- Within the CI pipeline
- Optionally as a pre-commit hook

---

## 24. Orchestrator State Machine

The Orchestrator writes no code. It manages state, tasks, agents, and context.

### 24.1 State Transitions

```
INIT
 ↓
LOAD_TASK
 ↓
DECOMPOSE
 ↓
ASSIGN_SUBTASK
 ↓
EXECUTE_SUBAGENT
 ↓
VALIDATE_OUTPUT
 ↓
UPDATE_CONTEXT
 ↓
CHECK_LIMITS
 ↓
DECIDE
 ├─ DONE
 ├─ SWITCH_AGENT
 ├─ NEXT_SUBTASK
 └─ CLOSE_SESSION
```

### 24.2 State-by-State Breakdown

#### INIT

Loads global configuration, token limit settings, and reads `global-context.md`.

#### LOAD_TASK

Retrieves the next task from `.agentic/tasks`. Validates that the task is properly defined. If it is too large, transitions to DECOMPOSE.

#### DECOMPOSE

Splits the task into atomic subtasks and persists them.

#### ASSIGN_SUBTASK

Selects an available agent. Constructs a minimal context bundle: task definition, relevant `context.md` files, and referenced specification files.

#### EXECUTE_SUBAGENT

Runs a single unit of work with no prior memory and no access to other tasks.

#### VALIDATE_OUTPUT

Executes the automated validator and checks acceptance criteria. On failure: retry with a different agent or re-decompose.

#### UPDATE_CONTEXT

Updates `context.md` files (where applicable) and `rolling-summary.md`.

#### CHECK_LIMITS

Evaluates token consumption. At the soft limit, triggers summarization. At the hard limit, transitions to CLOSE_SESSION.

#### DECIDE

| Condition | Action |
|-----------|--------|
| All subtasks completed | DONE |
| Subtasks remaining | NEXT_SUBTASK |
| Agent failed validation | SWITCH_AGENT |
| Token hard limit reached | CLOSE_SESSION |

#### CLOSE_SESSION

Follows the mandatory protocol:

1. Generate session summary
2. Update rolling summary
3. List modified `context.md` files
4. Terminate execution

The next session begins at INIT.

### 24.3 Persistent Orchestrator State

```yaml
orchestrator_state:
  current_task: TASK-014
  completed_subtasks:
    - SUB-014-A
  pending_subtasks:
    - SUB-014-B
  active_agent: antigravity
  session_id: 2026-02-04T18:22
```

**This state lives in a file, never in memory.**

---

## 25. Agent Skill Layer: The Missing Operational Bridge

The framework defines *what* agents must do and the *rules* they must follow. But without concrete, tool-native skill files that each agent reads at startup, no tool actually knows how to operate within the Agnostic-Agent system. The Agent Skill Layer closes this gap. It is the operational bridge between abstract rules and executable behavior — the layer that makes any agent, in any environment (IDE or terminal), immediately capable of picking up work and continuing where another agent left off.

### 25.1 The Problem

Every agentic tool consumes persistent instructions through its own native mechanism:

| Tool | Instruction File | Environment |
|------|-----------------|-------------|
| Claude Code | `CLAUDE.md` | Terminal |
| Cursor | `.cursor/rules/*.mdc` | IDE |
| Windsurf | `.windsurfrules` | IDE |
| Cline / Roo Code | `.clinerules` | IDE |
| GitHub Copilot | `.github/copilot-instructions.md` | IDE |
| Aider | `.aider.conf.yml` + conventions | Terminal |
| Continue.dev | `.continue/config.json` + rules | IDE |
| Augment | `.augment/rules.md` | IDE |

The framework's `/.agent-rules/base.md` defines behavioral rules abstractly, but none of these tools will read that file natively. Each tool needs a file in its own expected format and location, containing the full Agnostic-Agent behavioral contract translated into its native language.

### 25.2 The Solution: Generated Skill Files

The Agent Skill Layer consists of **auto-generated, tool-native instruction files** that the CLI produces from a single shared source. Every generated file encodes the identical behavioral contract — only the mechanical format differs.

```
.agentic/
  ├─ agent-rules/
  │   ├─ base.md                  ← canonical behavioral contract (source of truth)
  │   ├─ cloudcode.md             ← mechanical overrides for Claude Code
  │   ├─ cursor.md                ← mechanical overrides for Cursor
  │   ├─ copilot.md               ← mechanical overrides for Copilot
  │   └─ windsurf.md              ← mechanical overrides for Windsurf

# Generated files (in tool-native locations):
CLAUDE.md                         ← generated from base.md + cloudcode.md
.cursor/rules/agnostic-agent.mdc  ← generated from base.md + cursor.md
.windsurfrules                    ← generated from base.md + windsurf.md
.clinerules                       ← generated from base.md + cline.md
.github/copilot-instructions.md  ← generated from base.md + copilot.md
```

**Generation command:**

```bash
aa skills generate --tool claude-code
aa skills generate --tool cursor
aa skills generate --tool all
```

### 25.3 What Every Skill File Must Contain

Regardless of tool, every generated skill file must encode the following sections. These are non-negotiable — a skill file missing any of these is incomplete and invalid.

#### Section A: Identity and Role

```markdown
## Identity

You are an interchangeable execution agent operating within the Agnostic-Agent
(Agnostic-Agent) framework. You are NOT the system.
You are a disposable executor. Your memory is not trusted. Your session
can terminate at any time.
```

#### Section B: Boot Sequence

```markdown
## Boot Sequence (Execute on Every Session Start)

Before performing ANY work, execute these steps in order:

1. Read `.agentic/context/global-context.md` — understand project goals and invariants
2. Read `.agentic/context/rolling-summary.md` — understand current project state
3. Read `.agentic/tasks/in-progress.yaml` — identify assigned or claimable tasks
4. For each directory you will touch, read its `context.md` BEFORE making changes
5. Read the spec files referenced by your assigned task

Do NOT skip any step. Do NOT begin implementation before completing this sequence.
```

#### Section C: Behavioral Rules

```markdown
## Behavioral Rules (Immutable)

- Follow specifications strictly — specs define intent, not code
- Work only on your assigned subtask — never expand scope
- Read `context.md` before modifying any directory — no exceptions
- Do not retain memory beyond your file outputs
- Do not invent intent — if it is not in the spec, do not assume it
- Update `rolling-summary.md` when instructed
- Update `context.md` if you change logic, responsibilities, or patterns in a directory
- Commit code changes and context.md updates together — never separately
- If something feels ambiguous, stop and reduce scope
```

#### Section D: Task Pickup Protocol

```markdown
## Task Pickup Protocol

To begin work:

1. Check `.agentic/tasks/in-progress.yaml` for tasks assigned to you
2. If no task is assigned, check `.agentic/tasks/backlog.yaml` for the
   highest-priority unclaimed task
3. Before claiming a task, verify:
   - The task references valid spec files
   - The task defines clear acceptance criteria
   - The task scope fits within session limits (max 5 files, 2 directories)
4. Claim the task by moving it to `in-progress.yaml` with your agent identifier
5. Build your minimal context bundle: task definition + relevant context.md
   files + referenced specs
6. Begin execution
```

#### Section E: Task Handoff Protocol

```markdown
## Task Handoff Protocol (Session End or Agent Switch)

When your session ends or you are being replaced:

1. Persist all partial outputs to disk — nothing stays in memory
2. Update `context.md` for every directory where you changed logic
3. Generate a factual session summary:
   - Current task and progress percentage
   - Completed subtasks
   - Pending subtasks
   - Key decisions made during this session
   - Open questions or blockers
   - Exact next step to resume
4. Write summary to `.agentic/context/session-summary-<timestamp>.md`
5. Update `.agentic/context/rolling-summary.md`
6. Update task status in `.agentic/tasks/in-progress.yaml`
7. List all modified `context.md` files in the session summary

The next agent (or the same agent in a new session) will resume from these
files alone. No conversation history is passed. No chat context is shared.
```

#### Section F: Context Window Awareness

```markdown
## Context Window Management

You operate under finite token limits. Be aware of context consumption.

- Soft limit: Begin summarizing context. Execute SUMMARIZE_CONTEXT task.
- Hard limit: Stop all work immediately. Execute SESSION_CLOSE protocol.

Never attempt to continue working past the hard limit.
When resuming in a new session, load only:
  - rolling-summary.md
  - Current task definition
  - Relevant context.md files
  - Referenced spec files
```

#### Section G: Forbidden Actions

```markdown
## Forbidden Actions

- ❌ Operating from conversation memory instead of files
- ❌ Modifying code without first reading the directory's context.md
- ❌ Expanding scope beyond the assigned task
- ❌ Implementing large features end-to-end in a single session
- ❌ Keeping logic undocumented (every logical change requires context.md update)
- ❌ Skipping the boot sequence
- ❌ Claiming multiple tasks simultaneously
- ❌ Modifying files outside your task's declared scope
- ❌ Making assumptions about other tasks' internal state
```

#### Section H: Definition of Done

```markdown
## Definition of Done

Your task is complete ONLY when ALL of these are true:

- [ ] Acceptance criteria from the task definition are met
- [ ] Code changes are minimal and stay within declared scope
- [ ] All modified directories have updated context.md files
- [ ] rolling-summary.md reflects the new project state
- [ ] Task status is updated in .agentic/tasks/
- [ ] No undocumented logic remains
```

#### Section I: Context Generation Skill (context.md Lifecycle)

This is the skill that operationalizes the framework's most critical rule: **Context is Code.** Every agent must be capable of analyzing the actual source files in a directory and producing or updating its `context.md` — not from memory, not from conversation history, but from reading the code itself.

```markdown
## Context Generation Skill

You have the ability — and the obligation — to create and update context.md
files by analyzing the actual source code in a directory. This is not
optional documentation. It is a system-level operation equal in importance
to writing the code itself.

### When to Execute

GENERATE a new context.md when:
- You create a new directory that contains logic files
- The validator reports DIRECTORY_CONTEXT_MISSING for a directory
- You are assigned a CONTEXT_GENERATE task by the orchestrator
- You run `aa context generate <directory>`

UPDATE an existing context.md when:
- You modify logic, responsibilities, or patterns in a directory (CUOC protocol)
- You add, remove, or rename public interfaces
- You change dependency relationships between directories
- The validator reports CONTEXT_UPDATE_REQUIRED
- You run `aa context update <directory>`

### How to Analyze a Directory

Follow this exact procedure. Do NOT skip steps.

Step 1: LIST all files in the directory (ignore test files for purpose/responsibility
        analysis, but note their existence under Key Files)
Step 2: READ each source file — focus on:
        - Exported/public function signatures and their parameters
        - Type definitions, interfaces, structs, classes
        - Import statements (these reveal dependency relationships)
        - Error handling patterns
        - Side effects (I/O, network, database, file system)
Step 3: IDENTIFY the directory's role:
        - What does this code DO? (one sentence)
        - What architectural layer does it belong to? (core/infrastructure/adapter)
        - What can it depend on? What MUST it NOT depend on?
        - What patterns does it follow? (stateless, DI, repository, etc.)
Step 4: IDENTIFY constraints:
        - What are the MUST DO rules? (invariants the code enforces)
        - What are the CANNOT DO rules? (boundaries the code must not cross)
Step 5: COMPARE with existing context.md (if updating):
        - What changed? What is new? What was removed?
        - Update ONLY the sections that changed
        - Update the "Updated" timestamp with a brief change description

### Output Format (Strict)

Every context.md you produce must follow this exact structure.
Do NOT add sections. Do NOT remove sections. Do NOT exceed 60 lines / 600 tokens.

# [directory path relative to project root]

## Purpose
[One sentence: what this directory does]

## Responsibilities
[Bullet list: 2-5 specific responsibilities derived from the code]

## Local Architecture
[Bullet list: patterns, approaches, design decisions visible in the code]

## Dependencies
- Depends on: [list directories this code imports from]
- Depended on by: [list directories that import from this code, if known]
- Forbidden: [list directories this code must NEVER import from]

## MUST DO
[Bullet list: invariants and mandatory behaviors observed in the code]

## CANNOT DO
[Bullet list: boundaries and prohibitions — what this code must never do]

## Key Files
[Bullet list: significant files with one-line descriptions]

## Updated
[date] — [brief description of what changed]

### Quality Rules for Generated context.md

- Be FACTUAL — describe what the code does, not what you think it should do
- Be CONCISE — if a section needs more than 3 bullets, you are over-documenting
- Be CURRENT — reflect the code as it IS, not as it was or will be
- NEVER include how-to instructions, examples, or tutorials
- NEVER include historical rationale — that belongs in git history
- NEVER exceed 60 lines. If you cannot fit the description, the directory
  has too many responsibilities and should be split
- Dependencies must be derived from actual import statements, not assumptions
- MUST DO / CANNOT DO rules must be observable in the code or spec, not invented

### Updating vs. Regenerating

When UPDATING an existing context.md:
- Read the existing file FIRST
- Preserve sections that have not changed
- Modify only the sections affected by your code changes
- Always update the "Updated" timestamp

When GENERATING a new context.md:
- Analyze all source files in the directory
- Cross-reference with the nearest parent context.md for architectural role
- Cross-reference with relevant spec files for intended constraints
- Write the complete file from scratch following the output format

### Validation

After writing or updating context.md, self-check:
- [ ] Purpose is one sentence, factual, derived from code
- [ ] Responsibilities match the actual exported functions/types
- [ ] Dependencies match actual import statements
- [ ] MUST DO rules are observable in the code
- [ ] CANNOT DO rules are consistent with the architectural layer
- [ ] File is under 60 lines and 600 tokens
- [ ] Updated timestamp is present with change description
```

### 25.4 Skill File Generation: How It Works

The CLI composes each tool-native file by merging the canonical `base.md` with the tool-specific mechanical overrides. The process is deterministic and repeatable.

```
┌──────────────────────┐     ┌─────────────────────┐
│  base.md             │     │  cloudcode.md        │
│  (behavioral rules)  │  +  │  (mechanical rules)  │
└──────────┬───────────┘     └──────────┬──────────┘
           │                            │
           └────────────┬───────────────┘
                        ▼
              ┌─────────────────┐
              │   CLAUDE.md     │
              │   (generated)   │
              └─────────────────┘
```

The mechanical override files contain ONLY tool-specific instructions:

```markdown
# .agentic/agent-rules/cloudcode.md

## Mechanical Rules (Claude Code Specific)

- Use explicit file diffs when modifying code
- Prefer `claude-code` CLI conventions for file operations
- No inline commentary in code outputs
- Use markdown formatting for all non-code output
- When creating files, write complete files rather than patches
```

```markdown
# .agentic/agent-rules/cursor.md

## Mechanical Rules (Cursor Specific)

- Use Cursor's @file references when reading context
- Apply changes through Cursor's composer interface
- Use inline diff suggestions for small modifications
- Reference context.md files using @context.md notation
```

```markdown
# .agentic/agent-rules/copilot.md

## Mechanical Rules (Copilot Specific)

- Use Copilot Chat's workspace context features
- Apply changes via suggested edits
- Reference files using #file annotations
- Use /explain for verification before /fix operations
```

### 25.5 Context Generation via CLI

The CLI provides commands that trigger agents to execute the Context Generation Skill. These commands can be run manually by a developer, invoked by the orchestrator when the validator detects missing or stale context files, or integrated into CI pipelines as a pre-commit or post-merge check.

**Generate context.md for a single directory:**

```bash
aa context generate src/token
# Agent reads all source files in src/token/
# Agent produces src/token/context.md following the standard template
# Validator checks: under 60 lines, dependencies match imports, purpose is one sentence
```

**Update context.md for a directory after code changes:**

```bash
aa context update src/token
# Agent reads the existing src/token/context.md
# Agent reads all source files to detect what changed
# Agent updates only the affected sections
# Agent appends Updated timestamp
```

**Scan the entire project and generate missing context files:**

```bash
aa context scan
# Walks the source tree
# Identifies every directory containing logic files (.ts, .js, .py, .go)
# Reports which directories have context.md and which do not
# Output:
#   ✔ src/token/context.md           — exists, 42 lines
#   ✔ src/auth/context.md            — exists, 38 lines
#   ✗ src/events/handlers/context.md — MISSING
#   ✗ src/payments/context.md        — MISSING
#   ⚠ src/core/domain/context.md     — STALE (code modified after last context update)

aa context scan --fix
# Same as above, but automatically generates missing context.md files
# and flags stale ones for agent review
```

**Validate all existing context.md files against their directories:**

```bash
aa context validate
# For each context.md, reads the directory's source files and checks:
#   - Do listed dependencies match actual import statements?
#   - Do listed Key Files still exist?
#   - Is the file under 60 lines / 600 tokens?
#   - Is the Updated timestamp present?
# Output:
#   ✔ src/token/context.md           — valid
#   ⚠ src/auth/context.md            — dependency drift: lists src/core but imports src/infra
#   ✗ src/payments/context.md        — Key File 'stripe.adapter.ts' no longer exists
```

**Bulk regenerate all context files (fresh start):**

```bash
aa context regenerate --all
# WARNING: This overwrites all existing context.md files
# Use only during initial framework adoption or major refactors
# Each directory is analyzed from scratch
# All files follow the standard template
```

These commands integrate with the orchestrator's CONTEXT_GENERATE and CONTEXT_UPDATE task types. When the validator detects a DIRECTORY_CONTEXT_MISSING or CONTEXT_UPDATE_REQUIRED violation, the orchestrator can automatically spawn a sub-agent task that executes the Context Generation Skill for the affected directory.

---

## 26. Task Coordination: Parallel Work Across Tools

The framework must support multiple agents working on different tasks simultaneously, potentially on different machines, in different environments (IDE and terminal mixed). This requires an explicit coordination protocol.

### 26.1 Task Locking Mechanism

When an agent claims a task, it must be locked to prevent double-assignment:

```yaml
# .agentic/tasks/in-progress.yaml

- id: TASK-014
  title: Implement token budget controller
  status: in-progress
  locked_by: claude-code
  locked_at: 2026-02-04T14:30:00Z
  environment: terminal
  machine: dev-laptop-01

- id: TASK-015
  title: Build authentication middleware
  status: in-progress
  locked_by: cursor
  locked_at: 2026-02-04T14:32:00Z
  environment: ide
  machine: dev-desktop-02
```

### 26.2 Conflict Prevention Rules

```
RULE: No two agents may lock tasks that modify overlapping directories.

IF task_A.scope.directories INTERSECTS task_B.scope.directories
AND both are in-progress
THEN block the second claim until the first completes.
```

The CLI enforces this:

```bash
aa task claim TASK-015 --agent cursor
# ERROR: TASK-015 modifies src/auth which is locked by TASK-014 (claude-code)

aa task claim TASK-016 --agent cursor
# OK: TASK-016 scope (src/payments) does not overlap any in-progress task
```

### 26.3 Cross-Environment Continuation

The defining capability: an agent in a terminal can start work, and a different agent in an IDE can pick it up — or vice versa.

**Scenario: Terminal → IDE handoff**

```
1. Developer starts TASK-014 with Claude Code in terminal
2. Claude Code reaches session limit → executes SESSION_CLOSE
3. Session summary and rolling summary are persisted
4. Developer opens the project in Cursor
5. Cursor reads its generated skill file (.cursor/rules/agnostic-agent.mdc)
6. Cursor executes boot sequence:
   - Reads global-context.md
   - Reads rolling-summary.md (includes Claude Code's session summary)
   - Reads in-progress.yaml → sees TASK-014 with partial completion
   - Reads relevant context.md files
7. Cursor continues TASK-014 from the exact point Claude Code stopped
```

**Scenario: IDE → IDE parallel work**

```
1. Developer A uses Cursor on TASK-014 (src/token)
2. Developer B uses Windsurf on TASK-016 (src/payments)
3. Both agents read their own skill files (identical behavioral rules)
4. Both agents follow boot sequence independently
5. No directory overlap → no conflicts
6. Both update rolling-summary.md upon completion
7. Git merge handles the rolling-summary convergence
```

### 26.4 Rolling Summary Merge Strategy

When multiple agents complete tasks concurrently, their rolling summary updates may conflict. The resolution protocol:

```
1. Each agent writes its updates to a session-summary file (never conflicts)
2. The rolling-summary.md is updated via append-then-consolidate:
   a. Agent appends its session results to rolling-summary.md
   b. If a merge conflict occurs in Git, the CONSOLIDATE_SUMMARY task runs
   c. Consolidation merges both agents' summaries into a single coherent state
3. The CLI can automate this:
   aa context consolidate
```

---

## 27. Agent Boot Sequence: The Universal Startup Protocol

Every agent, regardless of tool or environment, must execute this exact sequence when beginning a session. This is the protocol that enables seamless continuation.

### 27.1 The Sequence

```
BOOT SEQUENCE (mandatory, in order)

Step 1 │ READ    .agentic/context/global-context.md
        │         → Understand project identity and invariants
        │
Step 2 │ READ    .agentic/context/rolling-summary.md
        │         → Understand current state of all work
        │
Step 3 │ READ    .agentic/tasks/in-progress.yaml
        │         → Find assigned task or identify claimable work
        │
Step 4 │ READ    Task's spec_refs (from task definition)
        │         → Understand what must be built and why
        │
Step 5 │ READ    context.md for each directory in task scope
        │         → Understand local rules and constraints
        │
Step 6 │ READ    Latest session-summary (if resuming a task)
        │         → Understand where the previous agent stopped
        │
Step 7 │ VERIFY  Context bundle is within token budget
        │         → If over soft limit, summarize before starting
        │
Step 8 │ BEGIN   Execution of assigned subtask
```

### 27.2 What the Boot Sequence Replaces

Without this protocol, agents rely on:

- ❌ Conversation history (lost between sessions)
- ❌ Internal memory (unreliable, tool-specific)
- ❌ Developer's verbal briefing (not scalable)
- ❌ Reading the entire codebase (token-wasteful)

With this protocol, agents rely on:

- ✅ Persisted files (survive any session boundary)
- ✅ Minimal, scoped context (token-efficient)
- ✅ Structured state (deterministic, auditable)
- ✅ Self-describing project (no human briefing needed)

### 27.3 Boot Sequence Validation

The CLI can verify an agent is properly initialized:

```bash
aa session validate
# ✔ global-context.md loaded
# ✔ rolling-summary.md loaded
# ✔ Task TASK-014 claimed and locked
# ✔ Spec refs resolved: 03-constraints.md, 04-architecture.md
# ✔ Context files loaded: src/token/context.md
# ✔ Token budget: 34,200 / 120,000 (28%)
# ✔ Boot sequence complete — ready to execute
```

---

## 28. Skill File Registry and Versioning

Skill files are generated artifacts. They must be tracked, versioned, and regenerated when the base rules change.

### 28.1 Registry File

```yaml
# .agentic/agent-rules/skill-registry.yaml

version: "1.3.0"
base_hash: "a4f2c8e1"  # SHA of base.md
generated_at: "2026-02-04T14:00:00Z"

skills:
  - tool: claude-code
    target: CLAUDE.md
    override: cloudcode.md
    hash: "b7d3f1a2"
  - tool: cursor
    target: .cursor/rules/agnostic-agent.mdc
    override: cursor.md
    hash: "c9e4a5b3"
  - tool: windsurf
    target: .windsurfrules
    override: windsurf.md
    hash: "d1f6b7c4"
  - tool: copilot
    target: .github/copilot-instructions.md
    override: copilot.md
    hash: "e3a8c9d5"
  - tool: cline
    target: .clinerules
    override: cline.md
    hash: "f5b0d1e6"
```

### 28.2 Drift Detection

If someone manually edits a generated skill file, it drifts from the canonical source. The CLI detects this:

```bash
aa skills check
# ⚠ DRIFT DETECTED: CLAUDE.md has been manually modified
#   Expected hash: b7d3f1a2
#   Actual hash:   x8k2m4n7
#   Run `aa skills generate --tool claude-code` to regenerate

aa skills check --fix
# Regenerated all skill files from canonical sources
```

### 28.3 Version Bumping

When `base.md` is modified, all skill files become stale:

```bash
aa skills status
# base.md modified since last generation
# 5 skill files need regeneration

aa skills generate --all
# ✔ Regenerated CLAUDE.md
# ✔ Regenerated .cursor/rules/agnostic-agent.mdc
# ✔ Regenerated .windsurfrules
# ✔ Regenerated .clinerules
# ✔ Regenerated .github/copilot-instructions.md
# Updated skill-registry.yaml → version 1.4.0
```

---

## 29. Environment-Specific Adaptations

Different environments (terminal vs. IDE) have different capabilities and constraints. The skill layer must account for these.

### 29.1 Terminal Agents (Claude Code, Aider)

Terminal agents typically have:

- Direct file system access
- Git command execution
- No visual UI for diff review
- Potentially larger context windows
- Session controlled by the developer or automation

Mechanical adaptations:

```markdown
## Terminal-Specific Rules

- Write complete files rather than inline patches
- Use explicit file paths in all references
- Output diffs in unified format for human review
- Use git commands directly for status and commit operations
- Print progress indicators to stdout
```

### 29.2 IDE Agents (Cursor, Windsurf, Copilot, Cline)

IDE agents typically have:

- File access mediated through the IDE
- Visual diff review built in
- Smaller effective context windows (shared with IDE features)
- Session tied to the IDE window lifecycle
- Access to IDE-specific features (symbol search, refactoring)

Mechanical adaptations:

```markdown
## IDE-Specific Rules

- Use the IDE's native file reference system (@file, #file, etc.)
- Apply changes through the IDE's diff/composer interface
- Leverage IDE symbol search for navigation
- Be aware that context window is shared with IDE features
  — budget accordingly (reduce by ~20%)
- Use IDE-native commit interfaces when available
```

### 29.3 Hybrid Workflow Support

A single project may have both terminal and IDE agents working simultaneously on different tasks. The framework supports this natively because:

- All state is in files (both environments see the same files)
- Task locking prevents directory conflicts
- Session summaries are environment-agnostic
- Rolling summary is the universal sync point
- Git is the coordination layer

