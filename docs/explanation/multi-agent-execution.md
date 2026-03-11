# Multi-Agent Execution

## Overview

The `agentic-agent` CLI now supports executing tasks with multiple AI agents:

- **Claude** (Anthropic API) - Production ready
- **Copilot** (GitHub CLI) - Experimental
- **Gemini** (Google) - Placeholder
- **Cursor** - Manual execution
- **Codex** (OpenAI) - Placeholder
- **Antigravity** - Placeholder
- **OpenCode** - Placeholder

## Quick Start

### 1. Configure Your Agent

Set your active agent in `agnostic-agent.yaml`:

```yaml
active_agent: claude-code

agents:
  claude-code:
    skill_packs:
      - atdd
      - dev-plans
    extra_rules:
      - Use TDD for all implementations
    auto_setup: true
```

Or use environment variable:

```bash
export AGENTIC_AGENT=claude-code
```

Or use CLI flag:

```bash
agentic-agent --agent claude autopilot start --execute-agent
```

### 2. Set API Keys

```bash
# For Claude
export ANTHROPIC_API_KEY="sk-ant-..."

# For Gemini
export GEMINI_API_KEY="..."

# For Codex/OpenAI
export OPENAI_API_KEY="sk-..."
```

### 3. Run Autopilot with Agent Execution

```bash
# Start autopilot with agent execution
agentic-agent autopilot start --execute-agent

# With specific agent
agentic-agent --agent claude autopilot start --execute-agent --max-iterations 5

# Dry run first
agentic-agent autopilot start --execute-agent --dry-run
```

## Agent Implementations

### Claude (Production Ready)

Uses Anthropic API to execute tasks with Claude models.

**Features:**
- ✅ Full API integration
- ✅ Token tracking
- ✅ Acceptance criteria validation
- ✅ Auto-completion on success

**Configuration:**

```yaml
agents:
  claude-code:
    model: claude-3-5-sonnet-20241022  # Optional, defaults to latest
```

**Environment:**

```bash
export ANTHROPIC_API_KEY="sk-ant-..."
```

**Example:**

```bash
# Claim task
agentic-agent task claim TASK-123

# Execute with Claude
agentic-agent --agent claude autopilot start --execute-agent --max-iterations 1
```

### Copilot (Experimental)

Uses GitHub CLI (`gh copilot`) to execute tasks.

**Features:**
- ⚠️ Requires `gh` CLI installed
- ⚠️ Limited to command-line suggestions
- ⚠️ Manual verification needed

**Requirements:**

```bash
# Install GitHub CLI
brew install gh

# Authenticate
gh auth login

# Install copilot extension
gh extension install github/gh-copilot
```

**Example:**

```bash
agentic-agent --agent copilot autopilot start --execute-agent
```

### Gemini (Placeholder)

Google Gemini integration - not yet implemented.

**Status:** Returns placeholder response until Google Gemini API is integrated.

**Planned Features:**
- Google Gemini API integration
- 1M token context window
- Multi-modal support

### Cursor (Manual Execution)

Cursor IDE integration - writes prompt files for manual execution.

**Features:**
- ✅ Generates prompt files
- ⚠️ Requires manual execution in Cursor IDE
- ⚠️ No automatic completion

**Workflow:**

```bash
# Start with Cursor agent
agentic-agent --agent cursor autopilot start --execute-agent

# Opens: .agentic/cursor-iteration.md
# → Open this file in Cursor IDE
# → Execute the prompt
# → Return to complete task manually
```

### Codex (Placeholder)

OpenAI Codex/GPT integration - not yet implemented.

**Status:** Returns placeholder response until OpenAI API is integrated.

**Planned Features:**
- OpenAI API integration
- GPT-4 support
- Code generation focus

### Antigravity (Placeholder)

Antigravity AI integration - not yet implemented.

**Status:** Returns placeholder response until integration is added.

### OpenCode (Placeholder)

OpenCode AI integration - not yet implemented.

**Status:** Returns placeholder response until integration is added.

## Architecture

### Executor Interface

All agents implement the `Executor` interface:

```go
type Executor interface {
    Execute(ctx context.Context, prompt string, task *models.Task) (*models.AgentExecutionResult, error)
}
```

### Agent Factory

The factory pattern selects the appropriate executor:

```go
executor := agents.NewExecutor(agentType)
result, err := executor.Execute(ctx, prompt, task)
```

### Result Structure

```go
type AgentExecutionResult struct {
    Output           string
    Success          bool
    CriteriaMet      []string
    CriteriaFailed   []string
    FilesModified    []string
    ErrorMessage     string
    TokensUsed       int
}
```

## Autopilot Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    Autopilot Loop                            │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  1. Find next claimable task from backlog                   │
│  2. Run readiness checks (specs, scope, etc.)               │
│  3. Claim task (record branch + timestamp)                  │
│  4. Generate context for scope directories                  │
│  5. Build context bundle with specs                         │
│  6. Execute agent (if --execute-agent enabled)              │
│     ├─ Build prompt with acceptance criteria                │
│     ├─ Call agent API/CLI                                   │
│     ├─ Check for completion signal                          │
│     └─ Auto-complete if success                             │
│  7. Loop to next task                                        │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

## Acceptance Criteria

Tasks can define acceptance criteria:

```yaml
tasks:
  - id: TASK-123
    title: "Implement user login"
    acceptance:
      - "Tests pass"
      - "API returns 200 for valid credentials"
      - "JWT token is generated"
```

Agents must output the completion signal when all criteria are met:

```
<promise>TASK COMPLETE</promise>
```

## Examples

### Example 1: Single Task with Claude

```bash
# Create task
cat > .agentic/tasks/backlog.yaml <<EOF
tasks:
  - id: TASK-001
    title: "Add hello world function"
    description: "Create a function that returns 'Hello, World!'"
    acceptance:
      - "Function exists in utils.go"
      - "Test passes"
    scope:
      - internal/utils
EOF

# Execute
agentic-agent --agent claude autopilot start --execute-agent --max-iterations 1

# Output:
# 🤖 Executing claude-code agent...
# ✅ Agent completed (tokens: 1234)
# ✅ All acceptance criteria met!
# ✅ Task TASK-001 completed successfully
```

### Example 2: Multiple Tasks in Sequence

```bash
# Start autopilot to process up to 5 tasks
agentic-agent --agent claude autopilot start --execute-agent --max-iterations 5

# Autopilot will:
# - Claim TASK-001
# - Execute Claude
# - Complete if successful
# - Move to TASK-002
# - Repeat until max iterations or no tasks left
```

### Example 3: Dry Run First

```bash
# See what would happen without executing
agentic-agent autopilot start --execute-agent --dry-run

# Output:
# [DRY RUN] Would claim task TASK-001 and generate context
# [DRY RUN] Would execute claude-code agent
# [DRY RUN] Would auto-complete on success
```

## Token Management

Each agent tracks token usage:

```go
result.TokensUsed  // Total tokens for this execution
```

**Token limits by agent:**

| Agent        | Context Limit | Notes                    |
|------------- |-------------- |------------------------- |
| Claude 3.5   | 200K tokens   | Production ready         |
| Gemini Pro   | 1M tokens     | Not yet implemented      |
| GPT-4        | 128K tokens   | Not yet implemented      |
| Copilot      | 32K tokens    | Via gh CLI               |

## Error Handling

If agent execution fails:

```bash
# Agent will report error but continue to next task
agentic-agent autopilot start --execute-agent --max-iterations 10

# Output:
# ⚠️  Agent execution error: API rate limit exceeded
# Task TASK-001 remains in-progress
# Moving to next task...
```

Tasks that fail remain in `in-progress` state and can be retried later.

## Testing

```bash
# Run all agent tests
go test ./internal/agents/... -v

# Run with integration tests (requires API keys)
ANTHROPIC_API_KEY=sk-ant-... go test ./internal/agents/... -v

# Run specific agent test
go test ./internal/agents/... -v -run TestClaudeExecutor
```

## Adding New Agents

To add a new agent:

1. **Create executor implementation:**

```go
// internal/agents/myagent.go

type MyAgentExecutor struct {
    apiKey string
    model  string
}

func NewMyAgentExecutor(apiKey, model string) *MyAgentExecutor {
    return &MyAgentExecutor{apiKey: apiKey, model: model}
}

func (e *MyAgentExecutor) Execute(ctx context.Context, prompt string, task *models.Task) (*models.AgentExecutionResult, error) {
    // Build prompt with acceptance criteria
    fullPrompt := buildStandardPrompt(prompt, task)

    // Call your agent's API
    output, err := callMyAgentAPI(fullPrompt)
    if err != nil {
        return nil, err
    }

    // Check criteria
    met, failed := checkCriteria(output, task.Acceptance)

    return &models.AgentExecutionResult{
        Output:         output,
        Success:        len(failed) == 0,
        CriteriaMet:    met,
        CriteriaFailed: failed,
        TokensUsed:     estimateTokens(output),
    }, nil
}
```

2. **Register in factory:**

```go
// internal/agents/executor.go

func NewExecutor(agentType string) Executor {
    switch agentType {
    case "myagent":
        return NewMyAgentExecutor("", "")
    // ... other cases
    }
}
```

3. **Add tests:**

```go
// internal/agents/myagent_test.go

func TestMyAgentExecutor_Execute(t *testing.T) {
    executor := NewMyAgentExecutor("test-key", "")
    // ... test implementation
}
```

4. **Update configuration:**

```yaml
# agnostic-agent.yaml

agents:
  myagent:
    skill_packs:
      - dev-plans
    auto_setup: true
```

## Next Steps

- [ ] Implement Gemini API integration
- [ ] Implement OpenAI/Codex API integration
- [ ] Add streaming support for real-time output
- [ ] Add token budget management
- [ ] Add checkpoint/resume for long-running tasks
- [ ] Add MCP server for chat AI integration
- [ ] Add session state persistence

## Related Documentation

- [Autopilot Design](./designs/cli-ralph-loop.md)
- [Agent Token Awareness](./designs/agent-token-awareness.md)
- [Ralph Loop Human Integration](./designs/ralph-loop-human-integration.md)
