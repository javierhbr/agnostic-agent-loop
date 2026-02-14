# Agnostic-Agent Framework

**A Comprehensive Architecture for Agent-Agnostic, Context-Isolated, Token-Managed AI Development — with Agent Skill Layer and TOON-Optimized Agent Communication**

*February 2026 — Reference Document (v3: includes Agent Skill Layer + TOON Integration)*

---

## 1. Foundational Principle: The Key That Unlocks Everything

This is genuine agentic development infrastructure, not superficial prompt engineering. The architecture presented here is practical, tool-independent, and resilient against both tool changes and context window constraints. The **Agnostic-Agent** framework uses Specification-Driven Development (SDD) as its structural backbone.

### 1.1 The Central Insight

**Agents must never serve as the source of truth.** The state of all work resides outside the agent, stored in versioned, structured artifacts.

When this principle is correctly applied:

- Migrating between tools (Cloud Code, Antigravity, Copilot, or any other agent) becomes straightforward
- Token consumption ceases to be a bottleneck
- Each agent becomes a replaceable executor rather than an indispensable part of the system

---

## 2. Architecture Layers

### 2.1 Layer 1: Project Brain (Source of Truth)

The foundation is a Git repository containing mandatory (not optional) documents:

```
.agentic/spec
  ├─ 00-vision.md
  ├─ 01-requirements.md
  ├─ 02-non-goals.md
  ├─ 03-constraints.md
  ├─ 04-architecture.md
  ├─ 05-domain-model.md
  └─ 06-acceptance-criteria.md

.agentic/tasks
  ├─ backlog.yaml
  ├─ in-progress.yaml
  └─ done.yaml

.agentic/context
  ├─ rolling-summary.md
  ├─ decisions.md
  └─ assumptions.md
```

**Golden rule:** If something is not recorded here, it does not exist.

---

## 3. Specification-Driven Development (SDD) Formalized

Every task is entirely independent of any specific agent. Here is an example sub-task definition (`tasks/backlog.yaml`):

```yaml
- id: TASK-014
  title: Implement token budget controller
  spec_refs:
    - .agentic/spec/03-constraints.md
    - .agentic/spec/04-architecture.md
  inputs:
    - .agentic/context/rolling-summary.md
  outputs:
    - src/token_controller.ts
    - tests/token_controller.test.ts
  acceptance:
    - Detects context > threshold
    - Triggers summarization
    - Persists summary
  status: pending
```

Notice the complete absence of any agent-specific references. There is no mention of Claude, Copilot, Antigravity, or any other particular tool. The task definition is purely functional:

- ❌ "use Claude"
- ❌ "use Copilot"
- ❌ "use Antigravity"

---

## 4. Layer 2: Agent Adapter Layer

Each tool connects through an adapter rather than embedding its own proprietary logic.

### 4.1 Shared Interface (Conceptual Model)

```typescript
interface AgentExecutor {
  name: string
  maxContextTokens: number
  execute(task: Task, context: ProjectContext): Result
}
```

### 4.2 Adapter Mapping

| Tool | Adapter File |
|------|-------------|
| Cloud Code CLI | `cloudcode.adapter.ts` |
| Antigravity | `antigravity.adapter.ts` |
| Copilot Agent | `copilot.adapter.ts` |

Every adapter performs the same operations: it reads identical input files, produces identical output formats, and updates the rolling summary following the same protocol.

---

## 5. Context Control and Token Budgeting

### 5.1 Fundamental Rule

**Never feed the complete context to an agent.** Instead, always provide only:

1. `rolling-summary.md`
2. Relevant specification files
3. The current task definition
4. Related source code (minimum viable scope)

Additionally, all structured data (task definitions, orchestrator state, validator output, configuration) should be serialized using TOON (Token-Oriented Object Notation) when injected into the agent's context window. This reduces structural overhead by 30–60% compared to JSON/YAML, freeing tokens for reasoning and implementation. See **Section 32: TOON Integration** for the complete strategy.

### 5.2 Token Monitor

Set thresholds per agent:

```yaml
agents:
  cloudcode:
    max: 200000
    soft_limit: 120000
  antigravity:
    max: 100000
    soft_limit: 60000
  copilot:
    max: 80000
    soft_limit: 50000
```

### 5.3 Automatic Summarization Trigger

When the soft limit is reached:

1. Pause execution
2. Run an automatic task: `SUMMARIZE_CONTEXT`
3. Store result in `.agentic/context/rolling-summary.md`
4. Discard stale context
5. Resume execution

The summarization prompt (agent-agnostic):

```
Summarize the current project state into:
- What is implemented
- What is pending
- Key decisions
- Open risks
- Next steps

Output must be factual, concise, and implementation-focused.
```

This ensures each agent inherits a clean state rather than accumulated noise.

---

## 6. Painless Agent Switching

Switching does not depend on the agent itself; it depends on the Task Engine.

### 6.1 Switching Logic

```
IF
  agent fails acceptance
  OR token soft-limit exceeded
  OR confidence score < threshold
THEN
  switch agent
```

### 6.2 What Gets Passed During a Switch

When transitioning to a different agent:

- ❌ **Do NOT pass:** the previous conversation or chat history
- ✅ **DO pass:** rolling-summary, specification files, task definition, and any partial outputs

---

## 7. Multi-Agent Quality Assurance

Agents can be assigned distinct roles while operating against the same specification:

1. Cloud Code handles implementation
2. Copilot reviews output against acceptance criteria
3. Antigravity performs optimization or refactoring

Each agent executes a different function, but all share the same specification as their contract.

---

## 8. End-to-End Execution Flow

```
1. Load task from backlog
2. Select agent (default: Cloud Code)
3. Build minimal context bundle
4. Execute task
5. Monitor tokens
   ├─ if soft-limit → summarize
6. Validate acceptance criteria
   ├─ pass → move to done
   └─ fail → switch agent
7. Update rolling-summary
```

---

## 9. Why This Approach Scales (and Others Do Not)

- ✔ Compatible with any LLM or agent tool
- ✔ No reliance on internal agent memory
- ✔ Token usage is actively managed
- ✔ Fully debuggable
- ✔ Fully auditable
- ✔ Supports humans in the loop

The essence of this framework is: **Git + Specification-Driven Development + Agent Orchestration** — not just a chat interface with extra features.

---

## 10. Hard Session Close and Fresh Session Protocol

### 10.1 Non-Negotiable Rule

**When the limit is reached, all work stops.** The session closes, state is persisted, and work resumes in an entirely new session.

### 10.2 Hard Limit Configuration

```yaml
session:
  hard_limit:
    cloudcode: 150000
    antigravity: 90000
    copilot: 70000
```

### 10.3 Mandatory Session Close Protocol

#### Step 1: Automatic Task — SESSION_CLOSE_SUMMARY

**Input:**
- `rolling-summary.md`
- Current task
- Active subtasks
- Recent decisions

**Output:**
- `.agentic/context/session-summary-<timestamp>.md`
- Updated `rolling-summary.md`

Standardized prompt (agent-agnostic):

```
You must close the current session.

Produce a factual summary containing:
1. Current task and progress
2. Completed subtasks
3. Pending subtasks
4. Key decisions and constraints
5. Open questions or risks
6. Exact next step to resume

No speculation. No repetition. No commentary.
```

#### Step 2: Persistence

```
.agentic/context
  ├─ rolling-summary.md        ← live state
  └─ session-summary-2026-02-04.md
```

#### Step 3: Terminate the Session

No further agent calls are permitted. A new session begins with a fresh, minimal context bundle.

---

## 11. Sub-Agents: Genuine Context Isolation

### 11.1 Core Rule

**The primary agent never implements directly.** The primary agent only orchestrates.

### 11.2 Agent Role Taxonomy

| Type | Role | Context Scope |
|------|------|--------------|
| Orchestrator | Plans and delegates | Minimal |
| Sub-agent | Executes a specific task | Ultra-reduced |
| Reviewer | Verifies acceptance criteria | Ultra-reduced |

### 11.3 Example Sub-Agent Task

```yaml
- id: SUB-023-A
  parent: TASK-023
  scope:
    directories:
      - src/token
  input_files:
    - src/token/context.md
    - .agentic/spec/03-constraints.md
  output:
    - src/token/budget.ts
```

The sub-agent has no visibility into:

- Conversation history
- Other directories
- Other tasks

It delivers its output, the orchestrator validates it, and the primary context remains uncontaminated.

---

## 12. Per-Directory Context (context.md) — A Critical Element

The `AGENTS.md` file placed in each directory is an essential practice for maintainability. It functions as **living documentation** that communicates the _intent_ of the folder, the architectural rules that govern it, and the reasoning behind the code it contains.

### 12.1 Mandatory Structure

Every significant directory must contain its own `AGENTS.md` file. This applies across the entire source tree:

```
src/
  ├─ core/
  │   ├─ domain/
  │   │   └─ context.md
  │   └─ application/
  │       └─ context.md
  ├─ infrastructure/
  │   ├─ adapters/
  │   │   ├─ inbound/
  │   │   │   └─ context.md
  │   │   └─ outbound/
  │   │       └─ context.md
  │   └─ config/
  │       └─ context.md
  └─ ...
```

Every directory containing logic gets a `AGENTS.md`. No exceptions.

### 12.2 Hexagonal Architecture: Context Files by Layer

The following definitions ensure that developers (and agents) understand the hexagonal boundaries within the codebase. Each layer has its own `AGENTS.md` with specific rules about what it does, what it depends on, and what it is forbidden from touching.

#### 12.2.1 `src/core/domain/context.md` — Pure Business Logic

> **Key Rule:** No dependencies on any other folder. No frameworks (except language primitives).

- **Purpose:** Contains the enterprise-wide business rules and state.
- **Contents:** Definitions of Domain Entities and Value Objects.
- **Constraint:** Logic here must be platform-agnostic. If the project migrates from Go to Java or Node, the logic in this folder should remain conceptually identical.

#### 12.2.2 `src/core/application/context.md` — Orchestration (Use Cases)

> **Key Rule:** Can only depend on the `domain`.

- **Purpose:** Coordinate the flow of data to and from the domain entities.
- **Contents:**
  - **Inbound Ports:** Interfaces defining what the application can do.
  - **Outbound Ports:** Interfaces defining what the application needs from the outside (Persistence, Messaging).
  - **Services:** Implementation of the Use Cases.

#### 12.2.3 `src/infrastructure/adapters/inbound/context.md` — The "Driving" Side

> **Key Rule:** Calls the Application Ports. Cannot call the Database directly.

- **Purpose:** Entry points into the application.
- **Contents:** REST controllers, GraphQL resolvers, CLI command handlers, or Message Queue listeners.
- **Responsibility:** Translating external requests (JSON, Protobuf) into Domain Models.

#### 12.2.4 `src/infrastructure/adapters/outbound/context.md` — The "Driven" Side

> **Key Rule:** Implements the Interfaces defined in `core/application/ports/outbound`.

- **Purpose:** Technical implementation of external needs.
- **Contents:** SQL queries, API clients for third-party services (Stripe, Twilio), and File System logic.
- **Constraint:** If the project swaps Postgres for MongoDB, only the code in a sub-folder here should change.

#### 12.2.5 `src/infrastructure/config/context.md` — The "Glue"

> **Key Rule:** The only place where all layers meet.

- **Purpose:** Dependency Injection and bootstrapping.
- **Contents:** Wire-up logic where you instantiate a `PostgresRepository` and inject it into an `OrderService`.
- **Responsibility:** Reading environment variables and configuring the application lifecycle.

### 12.3 Standard Format (Agent-Agnostic)

Here is a general-purpose example showing the format for any directory:

```markdown
# Directory: src/auth

## Purpose
Authentication and authorization logic.

## Responsibilities
- Token validation
- Role-based access

## Patterns Used
- Stateless services
- Dependency injection

## MUST DO
- Validate tokens on every request
- Log auth failures

## YOU CANNOT DO
- Store secrets in code
- Access database directly

## Interfaces
- auth.service.ts

## Related Specs
- .agentic/spec/05-domain-model.md
```

### 12.4 Reusable Template for `AGENTS.md`

Use this template for every directory to maintain consistency across the entire project:

```markdown
# Context: [Directory Name]

## 🎯 Responsibility
Describe what this folder does in one sentence.

## 🏗️ Architectural Role
- **Type:** [Core / Infrastructure / Adapter]
- **Direction:** [Inbound / Outbound / Internal]

## 🚦 Dependency Rules
- **Allowed to depend on:** [List folders]
- **Forbidden to depend on:** [List folders]

## 📝 Usage Notes
Briefly explain how to add a new file here (e.g., "When adding a new REST
endpoint, ensure the DTO is mapped to a Domain Model before calling the service").
```

### 12.5 Enforcement

**Enforced rule:** If any agent modifies a directory without first reading its `AGENTS.md`, that constitutes a system-level bug.

**Automated generation:** When the validator detects a directory containing logic files but no `AGENTS.md`, the CLI can generate one automatically using the **Context Generation Skill** (see Section 25.3, Section I). Agents execute this skill by analyzing the actual source files in the directory — reading exports, imports, types, and patterns — and producing a context file that reflects the code as it is. Run `aa context scan --fix` to detect and generate all missing context files across the project, or `aa context generate <directory>` for a specific directory.

---

## 13. Agent Rules and Skills Configuration

### 13.1 The Problem Being Solved

Each tool has its own workflow, user experience, and command set. However, the underlying behavioral logic must remain identical across all of them.

### 13.2 Solution: `/.agent-rules` Directory

```
.agentic/agent-rules
  ├─ base.md
  ├─ cloudcode.md
  ├─ antigravity.md
  ├─ copilot.md
  └─ cursor.md
```

#### `base.md` (Required for All Agents)

```
You are an interchangeable execution agent.

Rules:
- Follow specs strictly
- Work only on assigned subtask
- Read context.md before modifying any directory
- Do not expand scope
- Do not retain memory beyond outputs
- Always update rolling-summary if instructed
```

#### Tool-Specific Files (Only Mechanical Differences)

Example for `cloudcode.md`:

```
Use Cloud Code CLI conventions.
Prefer explicit file diffs.
No inline commentary.
```

Antigravity, Gemini, Cursor, and others follow the same pattern: only mechanical and interface differences, never behavioral deviations.

---

## 14. Forced Task Decomposition

### 14.1 Strict Rule

**If a task cannot fit within a single clean session, it is invalid** and must be split.

### 14.2 Maximum Task Size

```yaml
task_limits:
  max_files_touched: 5
  max_directories: 2
  max_output_tokens: 3000
```

### 14.3 Automatic Decomposition

The Orchestrator can execute a `TASK_SPLIT` operation with this prompt:

```
Split this task into the smallest possible independent subtasks.
Each subtask must be executable in isolation.
No subtask may require knowledge of another subtask's internal state.
```

---

## 15. Context is Code: The New Fundamental Rule

The `AGENTS.md` file is a first-class part of the system, not secondary documentation. Modifying logic without updating the corresponding context produces a corrupted state.

### 15.1 Mandatory Update Rule (Enforced)

Whenever an agent alters the logic, responsibilities, or patterns within a directory, it must update that directory's `AGENTS.md` within the same iteration.

This applies to:

- Architectural changes
- New responsibilities
- Pattern modifications
- Removal of previous rules
- New constraints

### 15.2 Exemptions

Updates are not required for:

- Trivial changes (typos, renames with no logical impact)
- Purely mechanical refactors with no semantic change

⚠️ **When in doubt, update.**

### 15.3 Context File Size Constraints

```yaml
context_md:
  max_lines: 60
  max_tokens: 600
```

**Containment rule:** If a context file starts growing, simplify it rather than expanding it.

### 15.4 Optimized context.md Format

```markdown
# src/token

## Purpose
Token budget control and session lifecycle.

## Responsibilities
- Track token usage
- Trigger summarization
- Close sessions on hard limit

## Local Architecture
- Stateless controller
- Pure functions
- No I/O side effects

## MUST DO
- Update token counters on every call
- Persist summaries before closing session

## YOU CANNOT DO
- Hold session state in memory
- Call agents directly

## Key Files
- budget.ts
- summarizer.ts

## Updated
2026-02-04 – Added hard session close logic
```

The **Updated** section is essential for audit traceability.

---

## 16. Context Update on Change (CUOC) Protocol

Every relevant change must follow this mini-protocol:

1. Modify the code
2. Evaluate: did the logic or responsibilities change?
3. If yes, execute the **Context Generation Skill** (Section 25.3, Section I) to update `AGENTS.md`
4. Commit both files together

The agent must read the directory's source files to determine what changed — not rely on its memory of what it thinks changed. The Context Generation Skill procedure ensures the update is factual and derived from code, not from the agent's internal state.

**Prohibited actions:**

- 🚫 Changing code without updating context
- 🚫 Changing context without changing code
- 🚫 Updating context from memory instead of reading the actual source files

### 16.1 Automated Validation

The Orchestrator should enforce:

```
IF files_changed in directory
AND change_type == logical
AND context.md not modified
THEN fail task
```

This can be implemented via git diff heuristics, explicit agent instructions, or CI hooks.

---

## 17. Mandatory Comments Before Session Close

Before executing `SESSION_CLOSE`, the agent must leave an explicit record of all context updates:

```markdown
## Context Updates
- src/token/context.md updated to reflect session close logic
- src/agent/context.md updated with sub-agent orchestration rules
```

This record goes into the `session-summary-*.md` file or at the end of `rolling-summary.md`.

---

## 18. Global Context: The Project's Mental Boot Loader

### 18.1 File Location

```
.agentic/context/global-context.md
```

### 18.2 Content (No Longer Than One Page)

```markdown
# Project Global Context

## Goal
Agnostic-Agent: agent-agnostic, specification-driven development framework.

## Core Rules
- Agents are interchangeable
- State lives in files
- Context is code
- Tasks must be atomic

## Non-Negotiables
- Update context.md on logic changes
- Hard session close on token limit
- Sub-agents for execution, orchestrator only coordinates

## Directory Context Rule
Every directory with logic must have context.md.

## Source of Truth
Specs > Contexts > Tasks > Code
```

This file serves as the mental initialization script for any new agent entering the project.

---

## 19. Directory Context: Lightweight, Not Heavy Documentation

Per-directory context replaces most traditional documentation, but it must remain small, factual, and operational.

### 19.1 Containment Rules

```yaml
directory_context:
  include_only:
    - what it does
    - rules
    - constraints
  exclude:
    - how to implement
    - extensive examples
    - historical rationale
```

Historical information belongs in Git commit history and session summaries, not in context files.

---

## 20. Sub-Agents Are Also Bound by Context Rules

Sub-agents must update `AGENTS.md` whenever they change logic, even though they:

- Work in isolation
- Receive minimal context
- Exist for a short duration

The Orchestrator validates this compliance before accepting any sub-agent output.

---

## 21. Fully Integrated Execution Flow

```
1.  Load task
2.  Decompose into subtasks
3.  Assign subtask to sub-agent
4.  Sub-agent reads context.md
5.  Sub-agent modifies code
6.  Sub-agent updates context.md (if needed)
7.  Orchestrator validates:
    - acceptance criteria
    - context consistency
8.  Update rolling-summary
9.  Check token limits
    ├─ soft → summarize
    └─ hard → session close
10. Persist session summary
11. Start new session
```

---

## 22. Foundational README: The Operational Contract

This README is not promotional material. It is the operational contract that every agent (or human) must read before contributing to the project. It is designed to be agent-agnostic, concise, strict, and execution-oriented.

### 22.1 Purpose

This project is constructed using the Agnostic-Agent framework — an agent-agnostic, specification-driven development system. Agents are replaceable executors. State, memory, and decisions reside in files, never in agent conversations.

If you are an AI agent or a human contributor, you must follow the rules below.

### 22.2 Core Principles

1. **Agents are disposable** — No agent memory is trusted. Sessions can terminate at any moment.
2. **Context is Code** — `AGENTS.md` files are system components. Modifying logic without updating context is invalid.
3. **Specifications govern everything** — Priority order: Specs > Contexts > Tasks > Code. Code alone never defines intent.
4. **Tasks must be atomic** — Oversized tasks must be decomposed. Execution occurs through sub-tasks and sub-agents.

### 22.3 Project Structure

```
.agentic/spec          → Source of truth (requirements, constraints, architecture)
.agentic/tasks         → Atomic tasks and subtasks
.agentic/context       → Global and rolling context summaries
.agentic/agent-rules   → Agent-specific execution adapters
/src                   → Implementation (each directory has its own context.md)
```

### 22.4 Global Context

Before performing any action, every agent must read:

```
.agentic/context/global-context.md
```

This file defines project goals, non-negotiable rules, and system invariants.

### 22.5 Directory Context (context.md)

Every directory that contains logic must have a `AGENTS.md`.

**What context.md is:**
- A local contract
- A rapid mental model
- A constraint enforcer

**What it is NOT:**
- Not a tutorial
- Not a changelog
- Not a design document

**Binding rule:** If you alter the logic, responsibilities, or patterns of a directory, you must update its `AGENTS.md` in the same iteration. Failure to do so invalidates the task.

### 22.6 Tasks and Subtasks

All work flows through tasks defined in `.agentic/tasks`.

**Task rules:**
- Tasks must be small and isolated
- Tasks must reference specs and contexts
- Tasks must define acceptance criteria
- Oversized tasks must be decomposed into subtasks

**Sub-agents:**
- Sub-agents execute implementation
- Orchestrator agents only coordinate
- Sub-agents receive minimal context
- Sub-agents are also obligated to update `AGENTS.md` when logic changes

### 22.7 Agent Rules (Agent-Agnostic)

All agents must adhere to `.agentic/agent-rules/base.md`.

Tool-specific rules reside in their own files:

```
.agentic/agent-rules/cloudcode.md
.agentic/agent-rules/antigravity.md
.agentic/agent-rules/copilot.md
.agentic/agent-rules/cursor.md
```

Only mechanical differences between tools are permitted. Behavioral rules are uniform across all agents.

### 22.8 Context and Token Management

Context windows are finite and actively enforced:

- **Soft limit** → summarize context
- **Hard limit** → close session immediately

#### Session Close Protocol

1. Generate a factual session summary
2. Update `rolling-summary.md`
3. Explicitly list all updated `AGENTS.md` files
4. Terminate the session

New sessions resume exclusively from persisted files.

### 22.9 Prohibited Actions

- ❌ Operating from conversation memory
- ❌ Modifying code without consulting `AGENTS.md`
- ❌ Expanding scope beyond the assigned task
- ❌ Allowing a single agent to implement large features end-to-end
- ❌ Maintaining undocumented logic in code

### 22.10 Definition of Done

A task is complete only when:

- Acceptance criteria are satisfied
- Code changes are minimal and properly scoped
- Relevant `AGENTS.md` files are updated
- The rolling context reflects the new state

### 22.11 Mental Model

> *Agents come and go. Files remember. Context explains. Specs decide.*

### 22.12 Instructions for AI Agents

You are an executor, not the system. Follow the files. Do not invent intent. Do not retain memory. Do not expand scope.

### 22.13 Design Intent

This framework is engineered to:

- Survive token limits
- Survive agent switching
- Scale across tools and humans
- Remain auditable and deterministic

If something feels ambiguous, stop and reduce scope.
