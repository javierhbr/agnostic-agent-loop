package agents

import (
	"context"
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

func TestNewExecutor_AllAgents(t *testing.T) {
	tests := []struct {
		name      string
		agentType string
		wantType  string
	}{
		{"Claude", "claude", "*agents.ClaudeExecutor"},
		{"Claude Code", "claude-code", "*agents.ClaudeExecutor"},
		{"Copilot", "copilot", "*agents.CopilotExecutor"},
		{"GitHub Copilot", "github-copilot", "*agents.CopilotExecutor"},
		{"Gemini", "gemini", "*agents.GeminiExecutor"},
		{"Cursor", "cursor", "*agents.CursorExecutor"},
		{"Codex", "codex", "*agents.CodexExecutor"},
		{"OpenAI", "openai", "*agents.CodexExecutor"},
		{"Antigravity", "antigravity", "*agents.AntigravityExecutor"},
		{"OpenCode", "opencode", "*agents.OpenCodeExecutor"},
		{"Mock", "mock", "*agents.executor"},
		{"Unknown", "unknown", "*agents.executor"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewExecutor(tt.agentType)
			if executor == nil {
				t.Error("Expected executor, got nil")
			}
		})
	}
}

func TestCopilotExecutor_Execute(t *testing.T) {
	executor := NewCopilotExecutor("")
	task := &models.Task{
		ID:    "TASK-123",
		Title: "Test task",
		Acceptance: []string{
			"Must complete successfully",
		},
	}

	result, err := executor.Execute(context.Background(), "Test prompt", task)

	// Copilot executor may fail if gh cli not available, but should not panic
	if err != nil {
		t.Logf("Copilot execution failed (expected if gh cli not installed): %v", err)
	}

	if result != nil && result.Output == "" {
		t.Error("Expected output or error, got empty result")
	}
}

func TestGeminiExecutor_Execute(t *testing.T) {
	executor := NewGeminiExecutor("", "")
	task := &models.Task{
		ID:    "TASK-123",
		Title: "Test task",
		Acceptance: []string{
			"Must complete successfully",
		},
	}

	result, err := executor.Execute(context.Background(), "Test prompt", task)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return placeholder until implemented
	if result.Success {
		t.Error("Expected placeholder response to be not successful")
	}
	if result.ErrorMessage == "" {
		t.Error("Expected error message for unimplemented executor")
	}
}

func TestCursorExecutor_Execute(t *testing.T) {
	executor := NewCursorExecutor("")
	task := &models.Task{
		ID:    "TASK-123",
		Title: "Test task",
		Acceptance: []string{
			"Must complete successfully",
		},
	}

	result, err := executor.Execute(context.Background(), "Test prompt", task)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Cursor requires manual execution
	if result.Success {
		t.Error("Expected manual execution to be not automatically successful")
	}
}

func TestCodexExecutor_Execute(t *testing.T) {
	executor := NewCodexExecutor("", "")
	task := &models.Task{
		ID:    "TASK-123",
		Title: "Test task",
		Acceptance: []string{
			"Must complete successfully",
		},
	}

	result, err := executor.Execute(context.Background(), "Test prompt", task)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return placeholder until implemented
	if result.Success {
		t.Error("Expected placeholder response to be not successful")
	}
}

func TestAntigravityExecutor_Execute(t *testing.T) {
	executor := NewAntigravityExecutor("")
	task := &models.Task{
		ID:    "TASK-123",
		Title: "Test task",
		Acceptance: []string{
			"Must complete successfully",
		},
	}

	result, err := executor.Execute(context.Background(), "Test prompt", task)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return placeholder until implemented
	if result.Success {
		t.Error("Expected placeholder response to be not successful")
	}
}

func TestOpenCodeExecutor_Execute(t *testing.T) {
	executor := NewOpenCodeExecutor("")
	task := &models.Task{
		ID:    "TASK-123",
		Title: "Test task",
		Acceptance: []string{
			"Must complete successfully",
		},
	}

	result, err := executor.Execute(context.Background(), "Test prompt", task)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return placeholder until implemented
	if result.Success {
		t.Error("Expected placeholder response to be not successful")
	}
}

func TestCheckCriteria(t *testing.T) {
	tests := []struct {
		name           string
		output         string
		criteria       []string
		wantMet        int
		wantFailed     int
	}{
		{
			name:       "Complete signal present",
			output:     "All done! <promise>TASK COMPLETE</promise>",
			criteria:   []string{"Test 1", "Test 2"},
			wantMet:    2,
			wantFailed: 0,
		},
		{
			name:       "No complete signal",
			output:     "Still working...",
			criteria:   []string{"Test 1", "Test 2"},
			wantMet:    0,
			wantFailed: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			met, failed := checkCriteria(tt.output, tt.criteria)
			if len(met) != tt.wantMet {
				t.Errorf("Expected %d met criteria, got %d", tt.wantMet, len(met))
			}
			if len(failed) != tt.wantFailed {
				t.Errorf("Expected %d failed criteria, got %d", tt.wantFailed, len(failed))
			}
		})
	}
}

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		want  int
	}{
		{"Empty", "", 0},
		{"Short", "test", 1},
		{"Medium", "This is a test message", 5},
		{"Long", "This is a much longer test message that should estimate more tokens", 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := estimateTokens(tt.text)
			if got != tt.want {
				t.Errorf("estimateTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}
