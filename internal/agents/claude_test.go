package agents

import (
	"context"
	"os"
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

func TestClaudeExecutor_Execute(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY not set, skipping integration test")
	}

	task := &models.Task{
		ID:    "TASK-123",
		Title: "Test task",
		Acceptance: []string{
			"Response contains 'hello'",
		},
	}

	executor := NewClaudeExecutor("", "claude-3-5-sonnet-20241022")
	result, err := executor.Execute(context.Background(), "Say hello", task)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Output == "" {
		t.Error("Expected output, got empty string")
	}
}
