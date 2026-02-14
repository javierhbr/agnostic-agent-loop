package agents

import (
	"context"
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

func TestExecutor_Execute(t *testing.T) {
	task := &models.Task{
		ID:    "TASK-123",
		Title: "Test task",
		Acceptance: []string{
			"Test passes",
		},
	}

	executor := NewExecutor("mock")
	result, err := executor.Execute(context.Background(), "test prompt", task)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Output == "" {
		t.Error("Expected output, got empty string")
	}
}
