package agents

import (
	"context"
	"fmt"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

type Executor interface {
	Execute(ctx context.Context, prompt string, task *models.Task) (*models.AgentExecutionResult, error)
}

type executor struct {
	agentType string
}

func NewExecutor(agentType string) Executor {
	switch agentType {
	case "mock":
		return &executor{agentType: "mock"}
	case "claude-code", "claude":
		return NewClaudeExecutor("", "")
	case "copilot", "github-copilot":
		return NewCopilotExecutor("")
	case "gemini":
		return NewGeminiExecutor("", "")
	case "cursor":
		return NewCursorExecutor("")
	case "codex", "openai":
		return NewCodexExecutor("", "")
	case "antigravity":
		return NewAntigravityExecutor("")
	case "opencode":
		return NewOpenCodeExecutor("")
	default:
		return &executor{agentType: agentType}
	}
}

func (e *executor) Execute(ctx context.Context, prompt string, task *models.Task) (*models.AgentExecutionResult, error) {
	// Mock implementation for testing
	if e.agentType == "mock" {
		return &models.AgentExecutionResult{
			Output:  "Mock output",
			Success: true,
			CriteriaMet: task.Acceptance,
			TokensUsed: 1000,
		}, nil
	}

	return nil, fmt.Errorf("unsupported agent type: %s", e.agentType)
}
