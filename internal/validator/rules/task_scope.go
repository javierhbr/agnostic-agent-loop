package rules

import (
	"github.com/javierbenavides/agentic-agent/internal/validator"
)

type TaskScopeRule struct{}

func (r *TaskScopeRule) Name() string {
	return "task-scope"
}

func (r *TaskScopeRule) Validate(ctx *validator.ValidationContext) (*validator.RuleResult, error) {
	// Mock implementation
	// Real one would check git status vs active task scope
	return &validator.RuleResult{
		RuleName: r.Name(),
		Status:   "PASS",
		Errors:   []string{},
	}, nil
}
