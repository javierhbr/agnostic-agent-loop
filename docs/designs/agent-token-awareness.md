# Agent-Aware Token Management for Ralph Loops

## Problem Statement

When AI agents (Claude Code, Cursor, Copilot, etc.) execute `/ralph-loop`, they need token awareness to:

1. **Avoid context exhaustion** — Running out of tokens mid-task
2. **Create checkpoints** — Save state before hitting limits
3. **Resume intelligently** — Pick up where they left off
4. **Provide user guidance** — Tell users when sessions need refreshing

**Current limitation:** The `internal/token/counter.go` is a basic estimator. AI agents have no visibility into their actual context usage or limits.

## Solution: Agent-Aware Token Tracking

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Ralph Loop Iteration                     │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ 1. Read specs + context                               │  │
│  │ 2. Execute iteration                                  │  │
│  │ 3. Check acceptance criteria                          │  │
│  │ 4. ✨ NEW: Track token usage                          │  │
│  │ 5. ✨ NEW: Check if approaching agent limit           │  │
│  └───────────────────────────────────────────────────────┘  │
│                           │                                  │
│                           ▼                                  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │        Approaching Limit? (e.g., 80% of max)          │  │
│  └───────────────────────────────────────────────────────┘  │
│          │                                  │                │
│          │ No                               │ Yes            │
│          ▼                                  ▼                │
│   Continue loop            ┌─────────────────────────────┐  │
│                            │ Create checkpoint:          │  │
│                            │ - Save current progress     │  │
│                            │ - Write resumption prompt   │  │
│                            │ - Output guidance to user   │  │
│                            └─────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Components

#### 1. Agent Registry with Token Limits

```go
// internal/token/limits.go

type AgentLimits struct {
    Name         string
    MaxTokens    int  // Total context window
    SafeLimit    int  // When to warn (e.g., 80% of max)
    OutputTokens int  // Reserved for output
}

var AgentLimitRegistry = map[string]AgentLimits{
    "claude-code": {
        Name:         "claude-code",
        MaxTokens:    200000,  // Sonnet 3.5/4
        SafeLimit:    160000,  // 80%
        OutputTokens: 8000,
    },
    "cursor": {
        Name:         "cursor",
        MaxTokens:    128000,  // Depends on model
        SafeLimit:    102400,
        OutputTokens: 4096,
    },
    "copilot": {
        Name:         "copilot",
        MaxTokens:    32000,   // Standard model
        SafeLimit:    25600,
        OutputTokens: 2048,
    },
    "gemini": {
        Name:         "gemini",
        MaxTokens:    1000000, // Gemini 1.5 Pro
        SafeLimit:    800000,
        OutputTokens: 8192,
    },
    "windsurf": {
        Name:         "windsurf",
        MaxTokens:    200000,
        SafeLimit:    160000,
        OutputTokens: 8000,
    },
}

func GetAgentLimits(agentName string) (AgentLimits, bool) {
    limits, ok := AgentLimitRegistry[agentName]
    return limits, ok
}
```

#### 2. Session Token Tracker

```go
// internal/token/session.go

type SessionTracker struct {
    agentName     string
    limits        AgentLimits
    currentUsage  int
    checkpointDir string
}

func NewSessionTracker(agentName, checkpointDir string) (*SessionTracker, error) {
    limits, ok := GetAgentLimits(agentName)
    if !ok {
        // Default fallback
        limits = AgentLimits{
            Name:         agentName,
            MaxTokens:    100000,
            SafeLimit:    80000,
            OutputTokens: 4000,
        }
    }

    return &SessionTracker{
        agentName:     agentName,
        limits:        limits,
        checkpointDir: checkpointDir,
    }, nil
}

func (st *SessionTracker) AddUsage(tokens int) {
    st.currentUsage += tokens
}

func (st *SessionTracker) IsApproachingLimit() bool {
    return st.currentUsage >= st.limits.SafeLimit
}

func (st *SessionTracker) GetUsagePercent() float64 {
    return float64(st.currentUsage) / float64(st.limits.MaxTokens) * 100
}

func (st *SessionTracker) GetRemainingTokens() int {
    return st.limits.MaxTokens - st.currentUsage - st.limits.OutputTokens
}
```

#### 3. Checkpoint Manager

```go
// internal/token/checkpoint.go

type Checkpoint struct {
    Timestamp       time.Time
    TaskID          string
    Iteration       int
    TokensUsed      int
    WorkCompleted   []string  // What was done so far
    RemainingWork   []string  // What's left to do
    FilesModified   []string
    TestsStatus     string
    ResumptionHints string    // How to resume
}

type CheckpointManager struct {
    baseDir string
}

func (cm *CheckpointManager) Save(cp *Checkpoint) (string, error) {
    filename := fmt.Sprintf("checkpoint-%s-%d.yaml", cp.TaskID, cp.Iteration)
    path := filepath.Join(cm.baseDir, filename)

    data, err := yaml.Marshal(cp)
    if err != nil {
        return "", err
    }

    if err := os.WriteFile(path, data, 0644); err != nil {
        return "", err
    }

    return path, nil
}

func (cm *CheckpointManager) GenerateResumptionPrompt(cp *Checkpoint) string {
    var b strings.Builder

    b.WriteString(fmt.Sprintf("# Resuming Task %s (Iteration %d)\n\n", cp.TaskID, cp.Iteration))
    b.WriteString("## Context\n")
    b.WriteString(fmt.Sprintf("Previous session used %d tokens (%.1f%% of limit).\n",
        cp.TokensUsed, float64(cp.TokensUsed)/200000*100))
    b.WriteString("This is a fresh session. Review the work completed so far and continue.\n\n")

    b.WriteString("## Work Completed\n")
    for _, item := range cp.WorkCompleted {
        b.WriteString(fmt.Sprintf("- ✓ %s\n", item))
    }
    b.WriteString("\n")

    b.WriteString("## Remaining Work\n")
    for _, item := range cp.RemainingWork {
        b.WriteString(fmt.Sprintf("- [ ] %s\n", item))
    }
    b.WriteString("\n")

    b.WriteString("## Files Modified\n")
    for _, file := range cp.FilesModified {
        b.WriteString(fmt.Sprintf("- %s\n", file))
    }
    b.WriteString("\n")

    b.WriteString("## Tests Status\n")
    b.WriteString(cp.TestsStatus)
    b.WriteString("\n\n")

    b.WriteString("## Next Steps\n")
    b.WriteString(cp.ResumptionHints)
    b.WriteString("\n")

    return b.String()
}
```

#### 4. Integration with Ralph Loop

```go
// internal/orchestrator/loop.go (enhanced)

type Loop struct {
    // ... existing fields ...
    sessionTracker *token.SessionTracker
    checkpointMgr  *token.CheckpointManager
    taskID         string
}

func (l *Loop) Run(ctx context.Context) error {
    // Detect agent
    agent := skills.DetectAgent()

    // Initialize session tracker
    tracker, err := token.NewSessionTracker(agent.Name, ".agentic/checkpoints")
    if err != nil {
        return err
    }
    l.sessionTracker = tracker
    l.checkpointMgr = token.NewCheckpointManager(".agentic/checkpoints")

    for iteration := 1; iteration <= l.maxIterations; iteration++ {
        // Run iteration
        output, err := l.runIteration(ctx)

        // Estimate token usage for this iteration
        iterationTokens := token.CountTokens(output)
        l.sessionTracker.AddUsage(iterationTokens)

        // Check if approaching limit
        if l.sessionTracker.IsApproachingLimit() {
            return l.createCheckpointAndPause(iteration)
        }

        // Check for completion
        if l.checkStopCondition(output) {
            return nil
        }
    }

    return fmt.Errorf("reached max iterations without completion")
}

func (l *Loop) createCheckpointAndPause(iteration int) error {
    // Gather current state
    cp := &token.Checkpoint{
        Timestamp:     time.Now(),
        TaskID:        l.taskID,
        Iteration:     iteration,
        TokensUsed:    l.sessionTracker.currentUsage,
        WorkCompleted: l.extractCompletedWork(),
        RemainingWork: l.extractRemainingWork(),
        FilesModified: l.getModifiedFiles(),
        TestsStatus:   l.getCurrentTestStatus(),
    }

    // Save checkpoint
    checkpointPath, err := l.checkpointMgr.Save(cp)
    if err != nil {
        return err
    }

    // Generate resumption prompt
    resumptionPrompt := l.checkpointMgr.GenerateResumptionPrompt(cp)
    resumptionPath := filepath.Join(".agentic/checkpoints",
        fmt.Sprintf("resume-%s.md", l.taskID))
    os.WriteFile(resumptionPath, []byte(resumptionPrompt), 0644)

    // Output guidance to agent/user
    l.outputCheckpointGuidance(checkpointPath, resumptionPath)

    return fmt.Errorf("paused: approaching token limit")
}

func (l *Loop) outputCheckpointGuidance(checkpointPath, resumptionPath string) {
    fmt.Println("\n⚠️  TOKEN LIMIT APPROACHING\n")
    fmt.Printf("Used: %d tokens (%.1f%% of limit)\n",
        l.sessionTracker.currentUsage,
        l.sessionTracker.GetUsagePercent())
    fmt.Printf("Remaining: %d tokens\n\n", l.sessionTracker.GetRemainingTokens())

    fmt.Println("📄 Checkpoint created:")
    fmt.Printf("  → %s\n\n", checkpointPath)

    fmt.Println("🔄 To resume in a fresh session:\n")
    fmt.Println("  Option 1 (Recommended for AI agents):")
    fmt.Println("    1. Start a new session/conversation")
    fmt.Printf("    2. Read: %s\n", resumptionPath)
    fmt.Println("    3. Continue with the resumption prompt\n")

    fmt.Println("  Option 2 (Manual):")
    fmt.Println("    1. Review checkpoint file")
    fmt.Println("    2. Close this session")
    fmt.Println("    3. Open fresh session with context from checkpoint\n")

    fmt.Println("🤖 AGENT GUIDANCE:")
    fmt.Println("  I've reached ~80% of my token limit. To continue this task:")
    fmt.Println("  - Tell the user to start a new session")
    fmt.Println("  - Provide them with the resumption file path")
    fmt.Println("  - They can paste the resumption prompt to continue work")
}
```

### CLI Integration

#### New Commands

```bash
# View current session token usage
agentic-agent token status

# List checkpoints for a task
agentic-agent token checkpoints TASK-123

# Load resumption prompt
agentic-agent token resume TASK-123

# Configure agent token limits
agentic-agent token limits --agent claude-code --max 200000 --safe 160000
```

#### Example Output

```
$ agentic-agent token status

Token Usage (claude-code session)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Current Task: TASK-936281-2
Iteration: 7/10

Usage:  165,432 / 200,000 tokens (82.7%) ⚠️
Safe Limit: 160,000 tokens (exceeded)
Remaining:  26,568 tokens

Status: 🟡 APPROACHING LIMIT

Recommendation:
  → Create checkpoint and resume in fresh session
  → Run: agentic-agent token checkpoint TASK-936281-2
```

### Configuration

```yaml
# agnostic-agent.yaml

tokens:
  tracking_enabled: true
  checkpoint_threshold: 0.80  # 80% of limit
  auto_checkpoint: true        # Auto-create on threshold

agents:
  defaults:
    max_tokens: 100000
    safe_limit: 80000

  overrides:
    - name: claude-code
      max_tokens: 200000
      safe_limit: 160000
      output_tokens: 8000

    - name: gemini
      max_tokens: 1000000
      safe_limit: 800000
      output_tokens: 8192
```

## Usage Scenarios

### Scenario 1: Agent Reaches Limit Mid-Task

```bash
# Agent executing ralph loop
/ralph-loop "Implement TASK-123..." --max-iterations 15

# After iteration 8:
⚠️  TOKEN LIMIT APPROACHING
Used: 164,000 tokens (82.0% of limit)

📄 Checkpoint created:
  → .agentic/checkpoints/checkpoint-TASK-123-8.yaml

🔄 To resume in a fresh session:
  Read: .agentic/checkpoints/resume-TASK-123.md

🤖 AGENT GUIDANCE:
  I've reached my token limit. Please:
  1. Start a new session
  2. Use command: agentic-agent token resume TASK-123
  3. I'll continue from iteration 9
```

### Scenario 2: User Manually Resumes

```bash
# New session starts
$ agentic-agent token resume TASK-123

# Resuming Task TASK-123 (Iteration 8)

## Context
Previous session used 164,000 tokens (82.0% of limit).
This is a fresh session. Review the work completed so far and continue.

## Work Completed
- ✓ Set up TypeScript with strict mode
- ✓ Configured ESLint for React
- ✓ Added path aliases (@/ → src/)
- ✓ Created VS Code settings

## Remaining Work
- [ ] Fix remaining ESLint errors in components/
- [ ] Add ESLint pre-commit hook
- [ ] Update documentation

## Files Modified
- tsconfig.json
- .eslintrc.js
- .vscode/settings.json
- package.json

## Tests Status
All tests passing. 3 ESLint warnings remain in src/components/TaskList.tsx

## Next Steps
1. Fix ESLint warnings in TaskList.tsx
2. Add husky pre-commit hook for ESLint
3. Verify all acceptance criteria met
4. Output <promise>TASK COMPLETE</promise>
```

### Scenario 3: Preventive Checkpoint

```bash
# Before starting long task, check if session is fresh
$ agentic-agent token status

Usage: 120,000 / 200,000 tokens (60.0%)
Status: 🟢 HEALTHY

# User decides to checkpoint anyway
$ agentic-agent token checkpoint TASK-124 --force

Checkpoint created: .agentic/checkpoints/checkpoint-TASK-124-0.yaml

Note: This is a preventive checkpoint. Session still has capacity.
```

## Benefits

1. **No Context Loss** — Checkpoints preserve all progress
2. **Agent-Aware** — Knows limits for Claude, Copilot, Gemini, etc.
3. **User Transparency** — Clear guidance on when/how to resume
4. **Automatic** — Triggers at 80% by default (configurable)
5. **Resumption Friendly** — Generated prompts make it easy to continue

## Implementation Phases

### Phase 1: Core Token Tracking
- [ ] Create `internal/token/limits.go` with agent registry
- [ ] Create `internal/token/session.go` for tracking
- [ ] Integrate with orchestrator loop
- [ ] Add basic CLI: `agentic-agent token status`

### Phase 2: Checkpointing
- [ ] Create `internal/token/checkpoint.go`
- [ ] Implement checkpoint creation logic
- [ ] Generate resumption prompts
- [ ] Add CLI: `agentic-agent token checkpoint TASK-ID`

### Phase 3: Auto-Recovery
- [ ] Detect existing checkpoints on task claim
- [ ] Offer to resume from checkpoint
- [ ] Auto-load last checkpoint state
- [ ] Add CLI: `agentic-agent token resume TASK-ID`

### Phase 4: Advanced Features
- [ ] Token usage analytics
- [ ] Per-task token budgets
- [ ] Multi-session task chains
- [ ] Integration with progress tracking

## Open Questions

1. **How to estimate tokens accurately?**
   - Current: 4 chars = 1 token (rough)
   - Better: Use tiktoken library for OpenAI models
   - Best: Agent-specific estimators

2. **Should we support mid-iteration checkpoints?**
   - Pro: More granular recovery
   - Con: Complexity in state management

3. **How to handle git state across sessions?**
   - Checkpoint includes branch name
   - Resume verifies correct branch
   - Warn if uncommitted changes exist

4. **What about parallel agent sessions?**
   - Each agent gets own session tracker
   - Checkpoints are agent-scoped
   - Resume from correct agent type

## Summary

This design adds token-awareness to the ralph loop workflow, enabling AI agents to:

- **Monitor** their context usage in real-time
- **Pause** automatically when approaching limits
- **Checkpoint** progress for resumption
- **Resume** seamlessly in fresh sessions
- **Guide** users on session management

The system is agent-agnostic (works with Claude, Copilot, Gemini, etc.) and provides clear user feedback on when/how to create new sessions.
