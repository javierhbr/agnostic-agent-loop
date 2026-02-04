# Agnostic-Agent Framework Implementation Plan

## Overview

Build the complete Agnostic-Agent framework as specified in TODO-01.md, TODO-02.md, and TODO-03.md. This is a Go-based CLI tool (`agentic-agent`) that enables agent-agnostic, specification-driven development with token management, context isolation, and multi-tool support.

---

## Technology Stack

| Component | Technology | Rationale |
|-----------|------------|-----------|
| CLI | Go 1.22+ | Spec mandates Go; excellent for CLI tools |
| CLI Framework | Cobra | Industry standard for Go CLIs |
| YAML Parsing | gopkg.in/yaml.v3 | Native YAML support |
| TOON Encoding | github.com/toon-format/toon-go | Specified in docs |
| File Watching | fsnotify | For validator hooks |
| Testing | Go testing + testify | Standard Go testing |
| Templates | text/template | For skill file generation |

---

## Project Structure

```
agnostic-agent-loop/
├── cmd/
│   └── agentic-agent/
│       └── main.go                 # CLI entry point
├── internal/
│   ├── config/
│   │   └── config.go               # agnostic-agent.yaml parsing
│   ├── project/
│   │   ├── init.go                 # Project initialization
│   │   └── structure.go            # Directory structure management
│   ├── tasks/
│   │   ├── task.go                 # Task model
│   │   ├── manager.go              # Task lifecycle (backlog→in-progress→done)
│   │   ├── decomposer.go           # Task splitting logic
│   │   └── lock.go                 # Task locking for parallel work
│   ├── context/
│   │   ├── global.go               # Global context management
│   │   ├── rolling.go              # Rolling summary management
│   │   ├── session.go              # Session summary management
│   │   ├── directory.go            # Per-directory context.md
│   │   └── generator.go            # Context generation from source analysis
│   ├── validator/
│   │   ├── validator.go            # Main validator orchestration
│   │   ├── rules/
│   │   │   ├── directory_context.go    # Rule 1: context.md required
│   │   │   ├── context_update.go       # Rule 2: CUOC enforcement
│   │   │   ├── task_scope.go           # Rule 3: Scope enforcement
│   │   │   ├── task_size.go            # Rule 4: Size limits
│   │   │   └── session_close.go        # Rule 5: Session close compliance
│   │   └── report.go               # Validation report generation
│   ├── orchestrator/
│   │   ├── state.go                # State machine model
│   │   ├── machine.go              # State transitions
│   │   └── executor.go             # Subtask execution coordination
│   ├── token/
│   │   ├── budget.go               # Token budget tracking
│   │   ├── counter.go              # Token counting (tiktoken)
│   │   └── summarizer.go           # Auto-summarization trigger
│   ├── skills/
│   │   ├── generator.go            # Skill file generation
│   │   ├── registry.go             # Skill registry & versioning
│   │   ├── drift.go                # Drift detection
│   │   └── templates/
│   │       ├── base.tmpl           # Base behavioral rules
│   │       ├── claude.tmpl         # Claude Code specific
│   │       ├── cursor.tmpl         # Cursor specific
│   │       ├── windsurf.tmpl       # Windsurf specific
│   │       ├── copilot.tmpl        # GitHub Copilot specific
│   │       └── cline.tmpl          # Cline specific
│   ├── encoding/
│   │   ├── toon.go                 # TOON encoder/decoder wrapper
│   │   └── bundle.go               # Context bundle assembly
│   └── adapters/
│       ├── adapter.go              # AgentExecutor interface
│       ├── cloudcode.go            # Claude Code adapter
│       ├── cursor.go               # Cursor adapter
│       └── generic.go              # Generic adapter for new tools
├── pkg/
│   └── models/
│       ├── task.go                 # Task struct definitions
│       ├── context.go              # Context structs
│       └── config.go               # Config structs
├── templates/
│   └── init/                       # Project initialization templates
│       ├── spec/                   # Spec file templates
│       ├── tasks/                  # Task file templates
│       └── context/                # Context file templates
├── go.mod
├── go.sum
└── Makefile
```

---

## Implementation Phases

### Phase 1: Foundation (Core Infrastructure)
**Goal:** Establish CLI skeleton and core data models

**Tasks:**
1. Initialize Go module and dependencies
2. Set up Cobra CLI with base command structure
3. Define core data models (Task, Context, Config)
4. Implement YAML file parsing utilities
5. Create `agentic-agent version` and `agentic-agent help` commands

**Files to create:**
- `cmd/agentic-agent/main.go`
- `cmd/agentic-agent/root.go`
- `pkg/models/*.go`
- `internal/config/config.go`
- `go.mod`, `Makefile`

**Dependencies:** None (foundation layer)

---

### Phase 2: Project Initialization
**Goal:** Create `.agentic/` directory structure

**Tasks:**
1. Implement `agentic-agent init` command
2. Create directory structure generator
3. Create template files for:
   - `.agentic/spec/` (7 spec documents)
   - `.agentic/tasks/` (backlog.yaml, in-progress.yaml, done.yaml)
   - `.agentic/context/` (global-context.md, rolling-summary.md, decisions.md, assumptions.md)
   - `.agentic/agent-rules/` (base.md + tool overrides)
   - `agnostic-agent.yaml` (framework config)
4. Validate initialization success

**Commands:**
```bash
agentic-agent init                    # Initialize in current directory
agentic-agent init --name "project"   # Initialize with project name
```

**Files to create:**
- `internal/project/init.go`
- `internal/project/structure.go`
- `templates/init/**/*`

**Dependencies:** Phase 1

---

### Phase 3: Task Management
**Goal:** Implement task lifecycle management

**Tasks:**
1. Implement task YAML parsing/writing
2. Create task state transitions (pending→in-progress→done)
3. Implement task claiming with locking
4. Implement task decomposition logic
5. Validate task constraints (max 5 files, 2 dirs)

**Commands:**
```bash
agentic-agent task list               # List all tasks
agentic-agent task show TASK-014      # Show task details
agentic-agent task claim TASK-014     # Claim a task (with lock)
agentic-agent task complete TASK-014  # Mark task complete
agentic-agent task decompose TASK-014 # Split large task
agentic-agent task create --title "..." --spec-refs "..."
```

**Files to create:**
- `internal/tasks/task.go`
- `internal/tasks/manager.go`
- `internal/tasks/decomposer.go`
- `internal/tasks/lock.go`

**Dependencies:** Phase 1, Phase 2

---

### Phase 4: Context Management
**Goal:** Implement context file management and generation

**Tasks:**
1. Implement global context management
2. Implement rolling summary updates
3. Implement session summary creation
4. Implement per-directory context.md management
5. **Context Generation Skill**: Analyze source files to generate context.md
   - List source files in directory
   - Extract exports, imports, types, patterns
   - Generate context.md following strict template
   - Validate: max 60 lines, 600 tokens

**Commands:**
```bash
agentic-agent context generate <dir>       # Generate context.md from source
agentic-agent context update <dir>         # Update existing context.md
agentic-agent context scan                 # Find all dirs needing context.md
agentic-agent context scan --fix           # Auto-generate missing context.md
agentic-agent context validate             # Validate all context.md files
agentic-agent context regenerate --all     # Regenerate all (destructive)
```

**Files to create:**
- `internal/context/global.go`
- `internal/context/rolling.go`
- `internal/context/session.go`
- `internal/context/directory.go`
- `internal/context/generator.go`

**Dependencies:** Phase 1, Phase 2

---

### Phase 5: Validator
**Goal:** Implement 5 validation rules

**Rules to implement:**
1. **Directory Context Required**: Dirs with .ts/.js/.py/.go must have context.md
2. **Context Update on Change (CUOC)**: Logic changes require context.md update
3. **Task Scope Enforcement**: Modified files must be in task's allowed directories
4. **Task Size Limits**: Max 5 files, 2 directories per task
5. **Session Close Compliance**: Session close requires summary files

**Commands:**
```bash
agentic-agent validate                # Run all validation rules
agentic-agent validate --rule cuoc    # Run specific rule
agentic-agent validate --report json  # Output JSON report
```

**Files to create:**
- `internal/validator/validator.go`
- `internal/validator/rules/directory_context.go`
- `internal/validator/rules/context_update.go`
- `internal/validator/rules/task_scope.go`
- `internal/validator/rules/task_size.go`
- `internal/validator/rules/session_close.go`
- `internal/validator/report.go`

**Dependencies:** Phase 1, Phase 3, Phase 4

---

### Phase 6: Agent Skill Layer
**Goal:** Generate tool-native skill files

**Tasks:**
1. Create base behavioral rules template
2. Create tool-specific mechanical rule templates
3. Implement skill file generator (base + override → tool file)
4. Implement skill registry with versioning
5. Implement drift detection

**Supported tools:**
- Claude Code → `CLAUDE.md`
- Cursor → `.cursor/rules/agnostic-agent.mdc`
- Windsurf → `.windsurfrules`
- GitHub Copilot → `.github/copilot-instructions.md`
- Cline → `.clinerules`

**Commands:**
```bash
agentic-agent skills generate --tool claude-code
agentic-agent skills generate --tool cursor
agentic-agent skills generate --all
agentic-agent skills check              # Check for drift
agentic-agent skills check --fix        # Regenerate drifted files
agentic-agent skills status             # Show version info
```

**Files to create:**
- `internal/skills/generator.go`
- `internal/skills/registry.go`
- `internal/skills/drift.go`
- `internal/skills/templates/*.tmpl`

**Dependencies:** Phase 1, Phase 2

---

### Phase 7: Token Management
**Goal:** Implement token budgeting and auto-summarization

**Tasks:**
1. Implement token counter (using tiktoken-go or similar)
2. Implement per-agent budget tracking
3. Implement soft limit detection → trigger summarization
4. Implement hard limit detection → trigger session close
5. Implement automatic summarization task

**Commands:**
```bash
agentic-agent token status              # Show current token usage
agentic-agent token budget --agent cursor  # Show agent-specific budget
agentic-agent context summarize         # Manual summarization trigger
```

**Files to create:**
- `internal/token/budget.go`
- `internal/token/counter.go`
- `internal/token/summarizer.go`

**Dependencies:** Phase 1, Phase 4

---

### Phase 8: TOON Encoding
**Goal:** Implement YAML→TOON conversion for agent context bundles

**Tasks:**
1. Integrate toon-go library
2. Implement YAML→TOON conversion for:
   - Task definitions
   - Orchestrator state
   - Validator output
   - Skill registry
   - Agent config
3. Implement context bundle assembly (TOON + Markdown)
4. Implement token savings calculation

**Commands:**
```bash
agentic-agent context build --task TASK-014 --agent cursor
agentic-agent context preview --task TASK-014 --format toon
agentic-agent context compare --task TASK-014  # YAML vs TOON token comparison
```

**Files to create:**
- `internal/encoding/toon.go`
- `internal/encoding/bundle.go`

**Dependencies:** Phase 1, Phase 4, Phase 7

---

### Phase 9: Orchestrator State Machine
**Goal:** Implement the 13-state orchestrator

**States:**
```
INIT → LOAD_TASK → DECOMPOSE → ASSIGN_SUBTASK → EXECUTE_SUBAGENT
→ VALIDATE_OUTPUT → UPDATE_CONTEXT → CHECK_LIMITS → DECIDE
→ (DONE | SWITCH_AGENT | NEXT_SUBTASK | CLOSE_SESSION)
```

**Tasks:**
1. Implement state machine model
2. Implement state transitions
3. Implement persistent state (`.agentic/orchestrator/state.yaml`)
4. Implement agent assignment logic
5. Implement session close protocol

**Commands:**
```bash
agentic-agent orchestrate start TASK-014    # Start orchestrating a task
agentic-agent orchestrate status            # Show orchestrator state
agentic-agent session validate              # Validate boot sequence
agentic-agent session close                 # Execute session close protocol
```

**Files to create:**
- `internal/orchestrator/state.go`
- `internal/orchestrator/machine.go`
- `internal/orchestrator/executor.go`

**Dependencies:** Phase 3, Phase 4, Phase 5, Phase 7

---

### Phase 10: Integration & Polish
**Goal:** Full integration testing and documentation

**Tasks:**
1. End-to-end integration tests
2. Error handling improvements
3. CLI help text and documentation
4. Example project with all components working
5. CI/CD setup (GitHub Actions)

**Files to create:**
- `tests/integration/*.go`
- `docs/README.md`
- `.github/workflows/ci.yaml`
- `examples/sample-project/`

**Dependencies:** All previous phases

---

## Critical Implementation Details

### Context Generation Algorithm (Phase 4)

```go
func GenerateContext(dir string) (*DirectoryContext, error) {
    // Step 1: List source files
    files := listSourceFiles(dir) // .ts, .js, .py, .go

    // Step 2: Read and analyze each file
    analysis := &DirectoryAnalysis{}
    for _, f := range files {
        exports := extractExports(f)
        imports := extractImports(f)
        types := extractTypes(f)
        patterns := detectPatterns(f)
        analysis.Merge(exports, imports, types, patterns)
    }

    // Step 3: Identify role
    role := identifyArchitecturalRole(analysis)

    // Step 4: Generate context
    ctx := &DirectoryContext{
        Path: dir,
        Purpose: derivePurpose(analysis),
        Responsibilities: deriveResponsibilities(analysis),
        LocalArchitecture: role.Patterns,
        Dependencies: analysis.Imports,
        MustDo: deriveInvariants(analysis),
        CannotDo: deriveBoundaries(role),
        KeyFiles: analysis.SignificantFiles,
        Updated: time.Now(),
    }

    // Step 5: Validate constraints
    if lines := countLines(ctx); lines > 60 {
        return nil, errors.New("context exceeds 60 lines")
    }

    return ctx, nil
}
```

### Validator Rule Interface (Phase 5)

```go
type ValidationRule interface {
    Name() string
    Validate(ctx *ValidationContext) (*RuleResult, error)
}

type RuleResult struct {
    Rule    string   `json:"rule"`
    Status  string   `json:"status"` // "pass" or "fail"
    Errors  []RuleError `json:"errors,omitempty"`
}

type RuleError struct {
    Directory string `json:"directory"`
    Message   string `json:"message"`
}
```

### Orchestrator State Transitions (Phase 9)

```go
type OrchestratorState string

const (
    StateInit           OrchestratorState = "INIT"
    StateLoadTask       OrchestratorState = "LOAD_TASK"
    StateDecompose      OrchestratorState = "DECOMPOSE"
    StateAssignSubtask  OrchestratorState = "ASSIGN_SUBTASK"
    StateExecuteSubagent OrchestratorState = "EXECUTE_SUBAGENT"
    StateValidateOutput OrchestratorState = "VALIDATE_OUTPUT"
    StateUpdateContext  OrchestratorState = "UPDATE_CONTEXT"
    StateCheckLimits    OrchestratorState = "CHECK_LIMITS"
    StateDecide         OrchestratorState = "DECIDE"
    StateDone           OrchestratorState = "DONE"
    StateSwitchAgent    OrchestratorState = "SWITCH_AGENT"
    StateNextSubtask    OrchestratorState = "NEXT_SUBTASK"
    StateCloseSession   OrchestratorState = "CLOSE_SESSION"
)
```

---

## Verification Strategy

### Per-Phase Testing

| Phase | Verification Method |
|-------|---------------------|
| 1. Foundation | `go build ./...`, `agentic-agent version` works |
| 2. Init | `agentic-agent init` creates correct structure |
| 3. Tasks | Task CRUD operations work, locking works |
| 4. Context | Context generation produces valid files |
| 5. Validator | All 5 rules detect violations correctly |
| 6. Skills | Generated files match tool formats |
| 7. Token | Token counting matches expected values |
| 8. TOON | YAML→TOON conversion is reversible |
| 9. Orchestrator | State machine transitions correctly |
| 10. Integration | Full workflow end-to-end |

### End-to-End Test Scenario

```bash
# 1. Initialize project
agentic-agent init --name "test-project"

# 2. Create a task
agentic-agent task create --title "Add auth module" --scope "src/auth"

# 3. Generate skill files
agentic-agent skills generate --all

# 4. Claim and work on task
agentic-agent task claim TASK-001

# 5. After code changes, generate context
agentic-agent context generate src/auth

# 6. Validate changes
agentic-agent validate

# 7. Complete task
agentic-agent task complete TASK-001

# 8. Build context bundle for agent
agentic-agent context build --task TASK-001 --agent cursor --format toon
```

---

## Estimated Component Sizes

| Component | Estimated Lines of Code |
|-----------|------------------------|
| CLI/Commands | ~500 |
| Config | ~200 |
| Project Init | ~400 |
| Task Management | ~600 |
| Context Management | ~800 |
| Context Generator | ~500 |
| Validator | ~700 |
| Skills Generator | ~600 |
| Token Management | ~400 |
| TOON Encoding | ~300 |
| Orchestrator | ~800 |
| Models | ~300 |
| **Total** | **~6,100** |

---

## Summary

This plan implements the complete Agnostic-Agent framework in 10 phases, starting with foundational infrastructure and building up to the full orchestrator. Each phase has clear deliverables, dependencies, and verification methods. The Go-based CLI (`agentic-agent`) will provide all commands specified in the TODO documents.
