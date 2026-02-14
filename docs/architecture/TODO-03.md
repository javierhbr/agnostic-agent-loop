
## 30. Skill File Content: Complete Reference for Each Tool

Below are the full generated skill file templates for each major tool. These represent what `aa skills generate` produces.

### 30.1 Claude Code (`CLAUDE.md`)

```markdown
# CLAUDE.md — Agnostic-Agent Framework Instructions

## Identity
You are an interchangeable execution agent operating within the Agnostic-Agent
framework. You are a disposable executor. Your memory is not trusted.

## Boot Sequence
1. Read .agentic/context/global-context.md
2. Read .agentic/context/rolling-summary.md
3. Read .agentic/tasks/in-progress.yaml
4. Read context.md for each directory in your task scope
5. Read spec files referenced by your task

## Behavioral Rules
- Follow specs strictly
- Work only on your assigned subtask
- Read context.md before modifying any directory
- Do not expand scope
- Do not retain memory beyond file outputs
- Update context.md when changing logic/responsibilities/patterns
- Commit code and context.md together

## Task Pickup
1. Check .agentic/tasks/in-progress.yaml for assigned tasks
2. If none, check .agentic/tasks/backlog.yaml
3. Claim by moving to in-progress.yaml with agent identifier
4. Build minimal context: task + context.md files + specs

## Session End Protocol
1. Persist all outputs to disk
2. Update context.md for changed directories
3. Write session summary to .agentic/context/session-summary-<timestamp>.md
4. Update .agentic/context/rolling-summary.md
5. Update task status in .agentic/tasks/

## Mechanical Rules (Claude Code Specific)
- Use explicit file diffs
- Prefer claude-code CLI conventions
- Write complete files, not patches
- No inline commentary in code

## Context Generation Skill
When creating or updating context.md for a directory:
1. List all source files: find <dir> -type f -name '*.ts' -o -name '*.go' -o -name '*.py'
2. Read each file — extract exports, imports, types, patterns
3. If context.md exists: read it first, update only changed sections
4. If context.md missing: generate from scratch using standard template
5. Verify: Purpose is one sentence. Max 60 lines. Dependencies match imports.
6. Always include Updated timestamp with change description
Output format: Purpose → Responsibilities → Local Architecture → Dependencies →
MUST DO → CANNOT DO → Key Files → Updated

## Forbidden Actions
- Operating from conversation memory
- Modifying code without reading context.md
- Expanding scope beyond task
- Skipping boot sequence

## Definition of Done
- Acceptance criteria met
- Code changes minimal and scoped
- context.md files updated
- rolling-summary.md reflects new state
```

### 30.2 Cursor (`.cursor/rules/agnostic-agent.mdc`)

```markdown
---
description: Agnostic-Agent behavioral and operational rules
globs:
alwaysApply: true
---

# Agnostic-Agent — Cursor Agent Rules

## Identity
You are an interchangeable execution agent operating within the Agnostic-Agent
framework. You are a disposable executor. Your memory is not trusted.

## Boot Sequence
1. Read @.agentic/context/global-context.md
2. Read @.agentic/context/rolling-summary.md
3. Read @.agentic/tasks/in-progress.yaml
4. Read @context.md for each directory in your task scope
5. Read spec files referenced by your task

## Behavioral Rules
- Follow specs strictly
- Work only on your assigned subtask
- Read context.md before modifying any directory
- Do not expand scope
- Do not retain memory beyond file outputs
- Update context.md when changing logic/responsibilities/patterns
- Commit code and context.md together

## Task Pickup
1. Check @.agentic/tasks/in-progress.yaml for assigned tasks
2. If none, check @.agentic/tasks/backlog.yaml
3. Claim by moving to in-progress.yaml with agent identifier
4. Build minimal context: task + context.md files + specs

## Session End Protocol
1. Persist all outputs to disk
2. Update context.md for changed directories
3. Write session summary to .agentic/context/session-summary-<timestamp>.md
4. Update @.agentic/context/rolling-summary.md
5. Update task status

## Mechanical Rules (Cursor Specific)
- Use @file references when reading context
- Apply changes through Cursor composer
- Use inline diff suggestions for small changes
- Reference context files using @context.md
- Be aware context window is shared with IDE features — budget -20%

## Context Generation Skill
When creating or updating context.md for a directory:
1. Use @file to read all source files in the directory
2. Extract exports, imports, types, and patterns from each file
3. If @context.md exists in directory: read first, update only changed sections
4. If missing: generate from scratch using standard template
5. Verify: Purpose is one sentence. Max 60 lines. Dependencies match imports.
6. Always include Updated timestamp with change description
Output format: Purpose → Responsibilities → Local Architecture → Dependencies →
MUST DO → CANNOT DO → Key Files → Updated

## Forbidden Actions
- Operating from conversation memory
- Modifying code without reading context.md
- Expanding scope beyond task
- Skipping boot sequence

## Definition of Done
- Acceptance criteria met
- Code changes minimal and scoped
- context.md files updated
- rolling-summary.md reflects new state
```

### 30.3 Windsurf (`.windsurfrules`)

```markdown
# Agnostic-Agent — Windsurf Agent Rules

## Identity
You are an interchangeable execution agent operating within the Agnostic-Agent
framework. You are a disposable executor. Your memory is not trusted.

## Boot Sequence
1. Read .agentic/context/global-context.md
2. Read .agentic/context/rolling-summary.md
3. Read .agentic/tasks/in-progress.yaml
4. Read context.md for each directory in your task scope
5. Read spec files referenced by your task

## Behavioral Rules
- Follow specs strictly
- Work only on your assigned subtask
- Read context.md before modifying any directory
- Do not expand scope
- Do not retain memory beyond file outputs
- Update context.md when changing logic/responsibilities/patterns
- Commit code and context.md together

## Task Pickup
1. Check .agentic/tasks/in-progress.yaml for assigned tasks
2. If none, check .agentic/tasks/backlog.yaml
3. Claim by moving to in-progress.yaml with agent identifier

## Session End Protocol
1. Persist all outputs to disk
2. Update context.md for changed directories
3. Write session summary to .agentic/context/session-summary-<timestamp>.md
4. Update .agentic/context/rolling-summary.md
5. Update task status

## Mechanical Rules (Windsurf Specific)
- Use Cascade's flow-based editing interface
- Apply multi-file changes through Cascade flows
- Use Windsurf's contextual awareness features
- Context window shared with IDE — budget -20%

## Context Generation Skill
When creating or updating context.md for a directory:
1. Read all source files in the directory using Cascade's file access
2. Extract exports, imports, types, and patterns from each file
3. If context.md exists: read first, update only changed sections
4. If missing: generate from scratch using standard template
5. Verify: Purpose is one sentence. Max 60 lines. Dependencies match imports.
6. Always include Updated timestamp with change description
Output format: Purpose → Responsibilities → Local Architecture → Dependencies →
MUST DO → CANNOT DO → Key Files → Updated

## Forbidden Actions
- Operating from conversation memory
- Modifying code without reading context.md
- Expanding scope beyond task
- Skipping boot sequence

## Definition of Done
- Acceptance criteria met
- Code changes minimal and scoped
- context.md files updated
- rolling-summary.md reflects new state
```

### 30.4 GitHub Copilot (`.github/copilot-instructions.md`)

```markdown
# Agnostic-Agent — Copilot Agent Rules

## Identity
You are an interchangeable execution agent operating within the Agnostic-Agent
framework. You are a disposable executor. Your memory is not trusted.

## Boot Sequence
1. Read .agentic/context/global-context.md
2. Read .agentic/context/rolling-summary.md
3. Read .agentic/tasks/in-progress.yaml
4. Read context.md for each directory in your task scope
5. Read spec files referenced by your task

## Behavioral Rules
- Follow specs strictly
- Work only on your assigned subtask
- Read context.md before modifying any directory
- Do not expand scope
- Do not retain memory beyond file outputs
- Update context.md when changing logic/responsibilities/patterns
- Commit code and context.md together

## Task Pickup
1. Check .agentic/tasks/in-progress.yaml for assigned tasks
2. If none, check .agentic/tasks/backlog.yaml
3. Claim by moving to in-progress.yaml with agent identifier

## Session End Protocol
1. Persist all outputs to disk
2. Update context.md for changed directories
3. Write session summary to .agentic/context/session-summary-<timestamp>.md
4. Update .agentic/context/rolling-summary.md
5. Update task status

## Mechanical Rules (Copilot Specific)
- Use #file annotations to reference project files
- Use /explain before /fix for verification
- Apply suggestions through Copilot's edit interface
- Use workspace context features for navigation

## Context Generation Skill
When creating or updating context.md for a directory:
1. Use #file references to read all source files in the directory
2. Extract exports, imports, types, and patterns from each file
3. If context.md exists: read first, update only changed sections
4. If missing: generate from scratch using standard template
5. Verify: Purpose is one sentence. Max 60 lines. Dependencies match imports.
6. Always include Updated timestamp with change description
Output format: Purpose → Responsibilities → Local Architecture → Dependencies →
MUST DO → CANNOT DO → Key Files → Updated

## Forbidden Actions
- Operating from conversation memory
- Modifying code without reading context.md
- Expanding scope beyond task
- Skipping boot sequence

## Definition of Done
- Acceptance criteria met
- Code changes minimal and scoped
- context.md files updated
- rolling-summary.md reflects new state
```

### 30.5 Cline / Roo Code (`.clinerules`)

```markdown
# Agnostic-Agent — Cline Agent Rules

## Identity
You are an interchangeable execution agent operating within the Agnostic-Agent
framework. You are a disposable executor. Your memory is not trusted.

## Boot Sequence
1. Read .agentic/context/global-context.md
2. Read .agentic/context/rolling-summary.md
3. Read .agentic/tasks/in-progress.yaml
4. Read context.md for each directory in your task scope
5. Read spec files referenced by your task

## Behavioral Rules
- Follow specs strictly
- Work only on your assigned subtask
- Read context.md before modifying any directory
- Do not expand scope
- Do not retain memory beyond file outputs
- Update context.md when changing logic/responsibilities/patterns
- Commit code and context.md together

## Task Pickup
1. Check .agentic/tasks/in-progress.yaml for assigned tasks
2. If none, check .agentic/tasks/backlog.yaml
3. Claim by moving to in-progress.yaml with agent identifier

## Session End Protocol
1. Persist all outputs to disk
2. Update context.md for changed directories
3. Write session summary to .agentic/context/session-summary-<timestamp>.md
4. Update .agentic/context/rolling-summary.md
5. Update task status

## Mechanical Rules (Cline Specific)
- Use Cline's file access tools for reading context
- Apply changes through Cline's edit workflow
- Use Cline's built-in diff display for review
- Operate within Cline's approval-based execution model

## Context Generation Skill
When creating or updating context.md for a directory:
1. Use Cline's file access tools to read all source files in the directory
2. Extract exports, imports, types, and patterns from each file
3. If context.md exists: read first, update only changed sections
4. If missing: generate from scratch using standard template
5. Verify: Purpose is one sentence. Max 60 lines. Dependencies match imports.
6. Always include Updated timestamp with change description
Output format: Purpose → Responsibilities → Local Architecture → Dependencies →
MUST DO → CANNOT DO → Key Files → Updated

## Forbidden Actions
- Operating from conversation memory
- Modifying code without reading context.md
- Expanding scope beyond task
- Skipping boot sequence

## Definition of Done
- Acceptance criteria met
- Code changes minimal and scoped
- context.md files updated
- rolling-summary.md reflects new state
```

---

## 31. Complete Project Skeleton with Agent Skill Layer

With the Agent Skill Layer integrated, the full project structure becomes:

```
project-root/
│
├─ .agentic/                              # Framework root
│   ├─ spec/                              # Source of truth
│   │   ├─ 00-vision.md
│   │   ├─ 01-requirements.md
│   │   ├─ 02-non-goals.md
│   │   ├─ 03-constraints.md
│   │   ├─ 04-architecture.md
│   │   ├─ 05-domain-model.md
│   │   └─ 06-acceptance-criteria.md
│   │
│   ├─ tasks/                             # Work units
│   │   ├─ backlog.yaml
│   │   ├─ in-progress.yaml
│   │   └─ done.yaml
│   │
│   ├─ context/                           # State and summaries
│   │   ├─ global-context.md
│   │   ├─ rolling-summary.md
│   │   ├─ decisions.md
│   │   ├─ assumptions.md
│   │   └─ session-summary-*.md
│   │
│   ├─ agent-rules/                       # Canonical rules (source of truth)
│   │   ├─ base.md                        # Behavioral contract (shared)
│   │   ├─ cloudcode.md                   # Mechanical: Claude Code
│   │   ├─ cursor.md                      # Mechanical: Cursor
│   │   ├─ windsurf.md                    # Mechanical: Windsurf
│   │   ├─ copilot.md                     # Mechanical: Copilot
│   │   ├─ cline.md                       # Mechanical: Cline
│   │   └─ skill-registry.yaml            # Tracks generated skill versions
│   │
│   ├─ validator/                         # Enforcement rules
│   │   └─ (rule definitions)
│   │
│   └─ orchestrator/                      # State machine config
│       └─ state.yaml
│
├─ CLAUDE.md                              # ← Generated skill (Claude Code)
├─ .cursor/rules/agnostic-agent.mdc       # ← Generated skill (Cursor)
├─ .windsurfrules                         # ← Generated skill (Windsurf)
├─ .clinerules                            # ← Generated skill (Cline)
├─ .github/copilot-instructions.md       # ← Generated skill (Copilot)
│
├─ src/                                   # Implementation
│   ├─ module-a/
│   │   ├─ context.md                     # Per-directory context
│   │   └─ (source files)
│   ├─ module-b/
│   │   ├─ context.md
│   │   └─ (source files)
│   └─ ...
│
└─ agnostic-agent.yaml                              # Framework configuration
```

### 31.1 Framework Configuration (`agnostic-agent.yaml`)

```yaml
# agnostic-agent.yaml — Project-level Agnostic-Agent framework configuration

project:
  name: "my-project"
  version: "0.1.0"

agents:
  cloudcode:
    max_tokens: 200000
    soft_limit: 120000
    hard_limit: 150000
  cursor:
    max_tokens: 100000
    soft_limit: 60000
    hard_limit: 90000
  windsurf:
    max_tokens: 100000
    soft_limit: 60000
    hard_limit: 90000
  copilot:
    max_tokens: 80000
    soft_limit: 50000
    hard_limit: 70000
  cline:
    max_tokens: 100000
    soft_limit: 60000
    hard_limit: 90000

task_limits:
  max_files_touched: 5
  max_directories: 2
  max_output_tokens: 3000

context_md:
  max_lines: 60
  max_tokens: 600

skills:
  auto_generate: true
  tools:
    - claude-code
    - cursor
    - windsurf
    - copilot
    - cline

toon:
  enabled: true
  delimiter: ","          # comma (default), tab, or pipe
  apply_to:
    - tasks
    - orchestrator_state
    - validator_output
    - skill_registry
    - agent_config
  preserve_yaml:
    - spec/*
    - context/*.md
    - agent-rules/*.md
```

---

## 32. TOON Integration: Token-Efficient Serialization for Agent Communication

### 32.1 The Problem: Structural Overhead Consumes Your Token Budget

Every time the CLI builds a context bundle for an agent, it includes structured data: task definitions, orchestrator state, backlog summaries, validator output, token budget configurations, and skill registries. In YAML or JSON, this data carries significant syntactic overhead — braces, brackets, quotes, and repeated key names on every object in every array.

Consider a backlog with 20 tasks. In YAML, each task repeats the same field names (`id`, `title`, `status`, `spec_refs`, `scope`, `acceptance`) 20 times. In JSON, add braces, brackets, and quotes on top. This redundancy is invisible to humans editing the file, but it consumes real tokens when injected into an agent's context window.

In a typical Agnostic-Agent session, the structured data portion of the context bundle consumes 500–1,500 tokens depending on task complexity. For agents with smaller context windows (Copilot at 80K, Cline at 100K), this overhead represents real capacity lost to formatting instead of reasoning. The problem compounds across the task lifecycle: every subtask assignment rebuilds the context bundle, every summarization cycle reloads state, and every agent switch reinjects the same structured data. Across a full task with 5–8 subtasks, the cumulative overhead can reach 5,000–10,000 tokens spent on parsing repeated field names and syntactic punctuation rather than on implementation.

### 32.2 The Solution: TOON as the Wire Format

TOON (Token-Oriented Object Notation) is a compact, human-readable serialization format designed specifically for LLM input. It encodes the same data model as JSON but eliminates redundant syntax: indentation replaces braces, field names are declared once in a header rather than repeated per object, and quoting is applied only when necessary.

TOON's benchmarks across multiple LLMs demonstrate **30–60% fewer tokens** than JSON and **74% comprehension accuracy** versus JSON's 70% — meaning agents not only receive less data, they parse it more reliably.

**The key architectural decision:** TOON is used as the **wire format between the CLI and the agent**, not as the storage format on disk.

```
Human-Editable Files (disk)       Agent Context Bundle (injected)
─────────────────────────         ──────────────────────────────
YAML / Markdown                   TOON (structured data)
                                  + Markdown (prose content)
  ┌──────────────┐
  │ backlog.yaml │──── aa CLI ────▶  TOON-encoded task list
  └──────────────┘    converts      (30-60% fewer tokens)
  ┌──────────────┐
  │ state.yaml   │──── aa CLI ────▶  TOON-encoded orchestrator state
  └──────────────┘    converts
  ┌──────────────┐
  │ context.md   │──── passed ────▶  Markdown (unchanged)
  └──────────────┘    directly
```

Files on disk remain YAML (human-editable, familiar, Git-friendly). The CLI converts only the structured data portions to TOON when assembling the context bundle. Prose content (specs, context.md, rolling-summary.md) passes through unchanged as Markdown.

### 32.3 What Gets TOON-Encoded vs. What Stays Markdown

The distinction is simple: **data tables → TOON, prose documents → Markdown.**

| Content Type | Storage (Disk) | Wire Format (Agent) | Reason |
|---|---|---|---|
| Task backlog | YAML | **TOON** | Uniform arrays of task objects — TOON's optimal case |
| In-progress tasks | YAML | **TOON** | Same structure as backlog |
| Orchestrator state | YAML | **TOON** | Flat key-value + arrays of subtask IDs |
| Validator output | JSON | **TOON** | Uniform array of error objects |
| Skill registry | YAML | **TOON** | Uniform array of tool definitions |
| Agent token config | YAML | **TOON** | Repetitive per-agent settings |
| Spec files | Markdown | **Markdown** | Free-form prose — TOON not applicable |
| context.md files | Markdown | **Markdown** | Prose with rules and constraints |
| rolling-summary.md | Markdown | **Markdown** | Narrative state description |
| Session summaries | Markdown | **Markdown** | Narrative handoff records |
| global-context.md | Markdown | **Markdown** | Project identity and invariants |
| Skill files | Markdown | **Markdown** | Natural language instructions |

### 32.4 Concrete Examples: YAML vs. TOON

#### Example 1: Task Backlog (5 Tasks)

**YAML (on disk — human-editable):**

```yaml
- id: TASK-014
  title: Implement token budget controller
  status: pending
  spec_refs:
    - .agentic/spec/03-constraints.md
    - .agentic/spec/04-architecture.md
  scope:
    directories:
      - src/token
  acceptance:
    - Detects context > threshold
    - Triggers summarization

- id: TASK-015
  title: Build authentication middleware
  status: pending
  spec_refs:
    - .agentic/spec/05-domain-model.md
  scope:
    directories:
      - src/auth
  acceptance:
    - Validates JWT tokens
    - Role-based access enforced

- id: TASK-016
  title: Create payment adapter
  status: pending
  spec_refs:
    - .agentic/spec/04-architecture.md
  scope:
    directories:
      - src/payments
  acceptance:
    - Stripe integration functional
    - Idempotent transactions

- id: TASK-017
  title: Implement event bus
  status: pending
  spec_refs:
    - .agentic/spec/04-architecture.md
  scope:
    directories:
      - src/events
  acceptance:
    - Publish/subscribe operational
    - Dead letter queue configured

- id: TASK-018
  title: Build CLI interface
  status: pending
  spec_refs:
    - .agentic/spec/01-requirements.md
  scope:
    directories:
      - src/cli
  acceptance:
    - All commands documented
    - Help text complete
```

**TOON (injected into agent context — generated by CLI):**

```toon
tasks[5]{id,title,status,spec_refs,scope_dirs,acceptance}:
  TASK-014,Implement token budget controller,pending,.agentic/spec/03-constraints.md;.agentic/spec/04-architecture.md,src/token,Detects context > threshold;Triggers summarization
  TASK-015,Build authentication middleware,pending,.agentic/spec/05-domain-model.md,src/auth,Validates JWT tokens;Role-based access enforced
  TASK-016,Create payment adapter,pending,.agentic/spec/04-architecture.md,src/payments,Stripe integration functional;Idempotent transactions
  TASK-017,Implement event bus,pending,.agentic/spec/04-architecture.md,src/events,Publish/subscribe operational;Dead letter queue configured
  TASK-018,Build CLI interface,pending,.agentic/spec/01-requirements.md,src/cli,All commands documented;Help text complete
```

The YAML version repeats `id`, `title`, `status`, `spec_refs`, `scope`, `directories`, and `acceptance` five times each. TOON declares these fields once in the header and streams pure data rows. For this five-task example, the token reduction is approximately **55%**.

#### Example 2: Orchestrator State

**YAML (on disk):**

```yaml
orchestrator_state:
  current_task: TASK-014
  completed_subtasks:
    - SUB-014-A
    - SUB-014-B
  pending_subtasks:
    - SUB-014-C
    - SUB-014-D
  active_agent: cursor
  session_id: 2026-02-04T18:22
  token_usage: 67400
  token_limit: 100000
```

**TOON (injected into agent context):**

```toon
orchestrator_state:
  current_task: TASK-014
  active_agent: cursor
  session_id: 2026-02-04T18:22
  token_usage: 67400
  token_limit: 100000
  completed_subtasks[2]:
    SUB-014-A,SUB-014-B
  pending_subtasks[2]:
    SUB-014-C,SUB-014-D
```

Modest savings on flat objects, but the explicit array length markers (`[2]`) provide the agent with structural guardrails that improve parsing reliability.

#### Example 3: Validator Output (Multiple Errors)

**JSON (produced by validator):**

```json
{
  "status": "fail",
  "errors": [
    {
      "rule": "CONTEXT_UPDATE_REQUIRED",
      "directory": "src/token",
      "message": "Logical changes without context.md update"
    },
    {
      "rule": "TASK_SCOPE_EXCEEDED",
      "directory": "src/auth",
      "message": "Modified files outside task scope"
    },
    {
      "rule": "DIRECTORY_CONTEXT_MISSING",
      "directory": "src/events/handlers",
      "message": "Directory contains logic but no context.md"
    }
  ]
}
```

**TOON (injected into agent context):**

```toon
status: fail
errors[3]{rule,directory,message}:
  CONTEXT_UPDATE_REQUIRED,src/token,Logical changes without context.md update
  TASK_SCOPE_EXCEEDED,src/auth,Modified files outside task scope
  DIRECTORY_CONTEXT_MISSING,src/events/handlers,Directory contains logic but no context.md
```

Every brace, bracket, quote, colon, and repeated key name is eliminated. The agent receives the same information in roughly **50% fewer tokens**.

### 32.5 CLI Integration: How the Conversion Works

The CLI handles TOON conversion transparently. The developer never writes TOON manually.

```
aa context build --task TASK-014 --agent cursor

Pipeline:
1. Read agnostic-agent.yaml → determine TOON settings
2. Read backlog.yaml → convert to TOON
3. Read in-progress.yaml → convert to TOON
4. Read orchestrator state.yaml → convert to TOON
5. Read relevant context.md files → pass through as Markdown
6. Read relevant spec files → pass through as Markdown
7. Read rolling-summary.md → pass through as Markdown
8. Assemble context bundle:
   ┌─────────────────────────────────┐
   │  [Markdown] global-context.md   │
   │  [Markdown] rolling-summary.md  │
   │  [TOON]     current task        │
   │  [TOON]     orchestrator state  │
   │  [Markdown] spec refs           │
   │  [Markdown] directory contexts  │
   │  [TOON]     token budget        │
   └─────────────────────────────────┘
9. Validate total tokens < soft_limit
10. Inject into agent
```

**Explicit CLI commands:**

```bash
# Preview the TOON-encoded context bundle (for debugging)
aa context preview --task TASK-014 --agent cursor --format toon

# Compare token counts between YAML and TOON
aa context compare --task TASK-014
# backlog.yaml:     YAML: 847 tokens → TOON: 382 tokens (savings: 55%)
# state.yaml:       YAML: 156 tokens → TOON: 112 tokens (savings: 28%)
# validator.json:   JSON: 234 tokens → TOON: 108 tokens (savings: 54%)
# Total savings:    465 tokens (38% of structured data overhead)

# Force YAML injection (disable TOON for debugging)
aa context build --task TASK-014 --agent cursor --no-toon
```

### 32.6 Go Implementation

The Agnostic-Agent CLI is built in Go (see architecture discussion). TOON has a mature Go library available:

```go
import toon "github.com/toon-format/toon-go"

// Convert task backlog from internal Go struct to TOON
func EncodeTasks(tasks []Task) (string, error) {
    data := make([]map[string]any, len(tasks))
    for i, t := range tasks {
        data[i] = map[string]any{
            "id":          t.ID,
            "title":       t.Title,
            "status":      t.Status,
            "spec_refs":   strings.Join(t.SpecRefs, ";"),
            "scope_dirs":  strings.Join(t.ScopeDirs, ";"),
            "acceptance":  strings.Join(t.Acceptance, ";"),
        }
    }
    wrapper := map[string]any{"tasks": data}
    return toon.Encode(wrapper)
}
```

The library handles tabular detection automatically — when all objects in an array share the same fields, TOON encodes them as a compact table with a single header row.

### 32.7 Agent Comprehension: Why TOON Works Better Than JSON for LLMs

TOON's advantage is not merely fewer tokens. The format is structurally clearer for language models:

**Explicit schema declaration.** The header `tasks[5]{id,title,status}:` tells the agent immediately: there are 5 tasks, each has exactly 3 fields. The agent knows the structure before reading any data. With JSON, the agent must infer structure from repeated patterns.

**Reduced parsing noise.** Every `{`, `}`, `[`, `]`, `"`, and `,` in JSON is a token that carries syntactic meaning but zero semantic value. TOON removes most of this noise, letting the agent focus on actual data content.

**Tabular clarity.** For uniform arrays (the dominant data structure in Agnostic-Agent), TOON's tabular format maps directly to how agents internally represent structured data. Benchmark results show 73.9% accuracy for TOON vs 69.7% for JSON across multiple models.

**Deterministic output.** TOON sorts keys alphabetically and uses consistent formatting, producing identical output for identical data. This predictability helps agents compare states across sessions.

### 32.8 When NOT to Use TOON

TOON is not a universal replacement. The framework explicitly preserves Markdown for:

- **Specification documents** — free-form prose with headings, paragraphs, lists. TOON cannot represent narrative content.
- **Context files** — `AGENTS.md` contains rules, constraints, and reasoning in natural language.
- **Summaries** — `rolling-summary.md` and session summaries are narrative descriptions of project state.
- **Skill files** — `CLAUDE.md`, `.cursorrules`, etc. are natural language behavioral instructions.
- **Deeply nested, non-uniform data** — if a data structure has irregular nesting or mixed-type fields, TOON may not offer savings over compact JSON.

The rule is straightforward: **if it's a uniform array of objects or a flat configuration block, encode it as TOON. If it's prose or deeply nested, leave it as Markdown or JSON.**

### 32.9 Configuration

TOON behavior is controlled in `agnostic-agent.yaml`:

```yaml
toon:
  enabled: true
  delimiter: ","          # comma (default), tab, or pipe
  length_markers: true    # explicit [N] array lengths
  apply_to:               # which structured files get TOON-encoded
    - tasks
    - orchestrator_state
    - validator_output
    - skill_registry
    - agent_config
  preserve_yaml:          # these always stay as Markdown/YAML
    - spec/*
    - context/*.md
    - agent-rules/*.md
```

Setting `toon.enabled: false` disables TOON entirely and falls back to YAML injection for all structured data. This is useful for debugging or when working with agents that have not been tested with TOON input.

### 32.10 Token Savings Across a Full Task Lifecycle

To quantify the cumulative impact, consider a typical task (TASK-014) with 4 subtasks, executed across 2 agent sessions with one agent switch:

| Event | Structured Data Injected | YAML Tokens | TOON Tokens | Saved |
|---|---|---|---|---|
| Session 1: Boot + Task Load | backlog + state + config | 1,140 | 580 | 560 |
| Session 1: Subtask 1 assign | subtask def + state | 320 | 185 | 135 |
| Session 1: Subtask 2 assign | subtask def + state | 340 | 195 | 145 |
| Session 1: Summarize (soft limit) | state snapshot | 280 | 160 | 120 |
| Agent Switch: Context rebuild | summary + backlog + state | 980 | 510 | 470 |
| Session 2: Boot + Resume | backlog + state + config | 1,060 | 545 | 515 |
| Session 2: Subtask 3 assign | subtask def + state | 330 | 190 | 140 |
| Session 2: Subtask 4 assign | subtask def + state | 350 | 200 | 150 |
| Session 2: Validator output | error/pass report | 180 | 85 | 95 |
| **Total** | | **4,980** | **2,650** | **2,330** |

**Result: 2,330 tokens saved across one task lifecycle — a 47% reduction in structural overhead.** For agents with 80K–100K context windows, this is equivalent to approximately 50 additional lines of source code the agent can reason about, or two extra context.md files it can hold in memory simultaneously.

For a project with 50 tasks, the cumulative savings approach **100,000+ tokens** — a full Copilot context window worth of capacity recovered from formatting overhead.

---

## 33. Summary: What This Framework Delivers

With the complete Agnostic-Agent system in place — including the Agent Skill Layer, Context Generation Skill, and TOON-optimized communication — you have:

- ✔ Executable, enforceable rules
- ✔ Automated compliance enforcement
- ✔ Deterministic orchestration
- ✔ Full context isolation
- ✔ Genuine agent switching capability
- ✔ Complete resilience to token limits
- ✔ Tool-native skill files for every major agent
- ✔ Seamless continuation across IDE and terminal environments
- ✔ Parallel task execution with conflict prevention
- ✔ Universal boot sequence for instant agent onboarding
- ✔ Skill versioning and drift detection
- ✔ Automated context.md generation and maintenance from source code analysis
- ✔ TOON-optimized structured data serialization (30–60% token savings)
- ✔ Transparent CLI-managed format conversion (YAML on disk, TOON in context)

This is framework-level infrastructure, not a workflow suggestion. It is agent-agnostic, context-isolated, token-managed, and designed to function reliably whether agents are AI models, human developers, or a mix of both — working in any combination of IDEs, terminals, and machines.

The Agnostic-Agent CLI (`aa`) is the single entry point: it generates skill files, builds TOON-optimized context bundles, validates compliance, generates and validates context.md files, coordinates tasks, detects drift, and manages the entire lifecycle from task creation through session close.

The essence: **Git + Specification-Driven Development + Agent Orchestration + Agent Skill Layer + TOON-Optimized Communication** — genuine AI-native development infrastructure where any agent can pick up any task and continue the work, regardless of which tool started it, with maximum token efficiency.

> *Agents come and go. Tools change. Environments differ. Files remember. Context explains. Specs decide. Skills instruct. TOON compresses.*