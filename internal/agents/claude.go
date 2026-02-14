package agents

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/javierbenavides/agentic-agent/pkg/models"
)

type ClaudeExecutor struct {
	client anthropic.Client
	model  string
}

func NewClaudeExecutor(apiKey, model string) *ClaudeExecutor {
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	if model == "" {
		model = "claude-3-5-sonnet-20241022"
	}

	client := anthropic.NewClient(option.WithAPIKey(apiKey))

	return &ClaudeExecutor{
		client: client,
		model:  model,
	}
}

func (c *ClaudeExecutor) Execute(ctx context.Context, prompt string, task *models.Task) (*models.AgentExecutionResult, error) {
	// Build full prompt with task context
	fullPrompt := c.buildPrompt(prompt, task)

	// Call Claude API
	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(c.model),
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(fullPrompt)),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("claude api error: %w", err)
	}

	// Extract output
	output := ""
	if len(message.Content) > 0 {
		// ContentBlockUnion has a Text field for text blocks
		if message.Content[0].Type == "text" {
			output = message.Content[0].Text
		}
	}

	// Check acceptance criteria
	criteriaMet, criteriaFailed := c.checkCriteria(output, task.Acceptance)

	return &models.AgentExecutionResult{
		Output:         output,
		Success:        len(criteriaFailed) == 0,
		CriteriaMet:    criteriaMet,
		CriteriaFailed: criteriaFailed,
		TokensUsed:     int(message.Usage.InputTokens + message.Usage.OutputTokens),
	}, nil
}

func (c *ClaudeExecutor) buildPrompt(basePrompt string, task *models.Task) string {
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

func (c *ClaudeExecutor) checkCriteria(output string, criteria []string) (met []string, failed []string) {
	// Simple string matching for now
	// TODO: More sophisticated criteria checking

	if strings.Contains(output, "<promise>TASK COMPLETE</promise>") {
		// Assume all criteria met if completion signal found
		return criteria, []string{}
	}

	// Otherwise, mark all as failed
	return []string{}, criteria
}
