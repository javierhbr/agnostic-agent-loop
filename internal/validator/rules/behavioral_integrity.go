package rules

import (
	"github.com/javierbenavides/agentic-agent/internal/validator"
)

type BehavioralIntegrityRule struct{}

func (r *BehavioralIntegrityRule) Name() string {
	return "BehavioralIntegrity"
}

func (r *BehavioralIntegrityRule) Validate(ctx *validator.ValidationContext) (*validator.RuleResult, error) {
	// Mock implementation
	return &validator.RuleResult{
		RuleName: r.Name(),
		Status:   "PASS",
	}, nil
}
