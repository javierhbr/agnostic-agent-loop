package agents

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// CopilotExecutor executes tasks using GitHub Copilot
type CopilotExecutor struct {
	model string
}

func NewCopilotExecutor(model string) *CopilotExecutor {
	if model == "" {
		model = "gpt-4"
	}
	return &CopilotExecutor{model: model}
}

func (e *CopilotExecutor) Execute(ctx context.Context, prompt string, task *models.Task) (*models.AgentExecutionResult, error) {
	fullPrompt := e.buildPrompt(prompt, task)

	// Write prompt to temp file for copilot CLI
	tmpFile := ".agentic/copilot-iteration.md"
	os.MkdirAll(".agentic", 0755)
	if err := os.WriteFile(tmpFile, []byte(fullPrompt), 0644); err != nil {
		return nil, fmt.Errorf("failed to write prompt: %w", err)
	}

	// Execute via gh copilot (if available)
	cmd := exec.CommandContext(ctx, "gh", "copilot", "suggest", "-t", "shell", fullPrompt)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("copilot execution error: %w", err)
	}

	outputStr := string(output)
	criteriaMet, criteriaFailed := checkCriteria(outputStr, task.Acceptance)

	return &models.AgentExecutionResult{
		Output:         outputStr,
		Success:        len(criteriaFailed) == 0,
		CriteriaMet:    criteriaMet,
		CriteriaFailed: criteriaFailed,
		TokensUsed:     estimateTokens(fullPrompt + outputStr),
	}, nil
}

func (e *CopilotExecutor) buildPrompt(basePrompt string, task *models.Task) string {
	var b strings.Builder
	b.WriteString(basePrompt)
	b.WriteString("\n\n")

	if len(task.Acceptance) > 0 {
		b.WriteString("## Acceptance Criteria\n")
		for _, criterion := range task.Acceptance {
			b.WriteString(fmt.Sprintf("- %s\n", criterion))
		}
		b.WriteString("\n")
	}

	b.WriteString("When all criteria are met, include in your response: <promise>TASK COMPLETE</promise>\n")
	return b.String()
}

// GeminiExecutor executes tasks using Google Gemini
type GeminiExecutor struct {
	apiKey string
	model  string
}

func NewGeminiExecutor(apiKey, model string) *GeminiExecutor {
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	if model == "" {
		model = "gemini-pro"
	}
	return &GeminiExecutor{apiKey: apiKey, model: model}
}

func (e *GeminiExecutor) Execute(ctx context.Context, prompt string, task *models.Task) (*models.AgentExecutionResult, error) {
	// TODO: Implement Gemini API integration
	// For now, return placeholder
	return &models.AgentExecutionResult{
		Output:         "Gemini integration pending",
		Success:        false,
		CriteriaFailed: task.Acceptance,
		ErrorMessage:   "Gemini API integration not yet implemented",
		TokensUsed:     0,
	}, nil
}

// CursorExecutor executes tasks using Cursor
type CursorExecutor struct {
	model string
}

func NewCursorExecutor(model string) *CursorExecutor {
	if model == "" {
		model = "gpt-4"
	}
	return &CursorExecutor{model: model}
}

func (e *CursorExecutor) Execute(ctx context.Context, prompt string, task *models.Task) (*models.AgentExecutionResult, error) {
	fullPrompt := buildStandardPrompt(prompt, task)

	// Write prompt to temp file for cursor
	tmpFile := ".agentic/cursor-iteration.md"
	os.MkdirAll(".agentic", 0755)
	if err := os.WriteFile(tmpFile, []byte(fullPrompt), 0644); err != nil {
		return nil, fmt.Errorf("failed to write prompt: %w", err)
	}

	// Cursor doesn't have a CLI, so we create instruction file
	return &models.AgentExecutionResult{
		Output:         fmt.Sprintf("Prompt written to %s. Open in Cursor to execute.", tmpFile),
		Success:        false,
		CriteriaFailed: task.Acceptance,
		ErrorMessage:   "Manual execution required in Cursor IDE",
		TokensUsed:     0,
	}, nil
}

// CodexExecutor executes tasks using OpenAI Codex
type CodexExecutor struct {
	apiKey string
	model  string
}

func NewCodexExecutor(apiKey, model string) *CodexExecutor {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if model == "" {
		model = "gpt-4"
	}
	return &CodexExecutor{apiKey: apiKey, model: model}
}

func (e *CodexExecutor) Execute(ctx context.Context, prompt string, task *models.Task) (*models.AgentExecutionResult, error) {
	// TODO: Implement OpenAI API integration
	return &models.AgentExecutionResult{
		Output:         "Codex integration pending",
		Success:        false,
		CriteriaFailed: task.Acceptance,
		ErrorMessage:   "OpenAI API integration not yet implemented",
		TokensUsed:     0,
	}, nil
}

// AntigravityExecutor executes tasks using Antigravity
type AntigravityExecutor struct {
	model string
}

func NewAntigravityExecutor(model string) *AntigravityExecutor {
	return &AntigravityExecutor{model: model}
}

func (e *AntigravityExecutor) Execute(ctx context.Context, prompt string, task *models.Task) (*models.AgentExecutionResult, error) {
	// TODO: Implement Antigravity integration
	return &models.AgentExecutionResult{
		Output:         "Antigravity integration pending",
		Success:        false,
		CriteriaFailed: task.Acceptance,
		ErrorMessage:   "Antigravity integration not yet implemented",
		TokensUsed:     0,
	}, nil
}

// OpenCodeExecutor executes tasks using OpenCode
type OpenCodeExecutor struct {
	model string
}

func NewOpenCodeExecutor(model string) *OpenCodeExecutor {
	return &OpenCodeExecutor{model: model}
}

func (e *OpenCodeExecutor) Execute(ctx context.Context, prompt string, task *models.Task) (*models.AgentExecutionResult, error) {
	// TODO: Implement OpenCode integration
	return &models.AgentExecutionResult{
		Output:         "OpenCode integration pending",
		Success:        false,
		CriteriaFailed: task.Acceptance,
		ErrorMessage:   "OpenCode integration not yet implemented",
		TokensUsed:     0,
	}, nil
}

// Helper functions

func buildStandardPrompt(basePrompt string, task *models.Task) string {
	var b strings.Builder
	b.WriteString(basePrompt)
	b.WriteString("\n\n")

	if len(task.Acceptance) > 0 {
		b.WriteString("## Acceptance Criteria\n")
		for _, criterion := range task.Acceptance {
			b.WriteString(fmt.Sprintf("- %s\n", criterion))
		}
		b.WriteString("\n")
	}

	b.WriteString("When all criteria are met, include in your response: <promise>TASK COMPLETE</promise>\n")
	return b.String()
}

func checkCriteria(output string, criteria []string) (met []string, failed []string) {
	if strings.Contains(output, "<promise>TASK COMPLETE</promise>") {
		return criteria, []string{}
	}
	return []string{}, criteria
}

func estimateTokens(text string) int {
	// Rough estimate: ~4 characters per token
	return len(text) / 4
}
