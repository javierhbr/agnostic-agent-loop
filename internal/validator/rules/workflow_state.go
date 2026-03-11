package rules

import (
	"github.com/javierbenavides/agentic-agent/internal/validator"
)

type WorkflowStateRule struct{}

func (r *WorkflowStateRule) Name() string {
	return "WorkflowState"
}

func (r *WorkflowStateRule) Validate(ctx *validator.ValidationContext) (*validator.RuleResult, error) {
	// Mock implementation
	return &validator.RuleResult{
		RuleName: r.Name(),
		Status:   "PASS",
	}, nil
}
