package rules

import (
	"github.com/javierbenavides/agentic-agent/internal/validator"
)

type PlatformAlignmentRule struct{}

func (r *PlatformAlignmentRule) Name() string {
	return "PlatformAlignment"
}

func (r *PlatformAlignmentRule) Validate(ctx *validator.ValidationContext) (*validator.RuleResult, error) {
	// Mock implementation
	return &validator.RuleResult{
		RuleName: r.Name(),
		Status:   "PASS",
	}, nil
}
