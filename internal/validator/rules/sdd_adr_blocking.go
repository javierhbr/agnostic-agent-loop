package rules

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/javierbenavides/agentic-agent/internal/sdd"
	"github.com/javierbenavides/agentic-agent/internal/validator"
)

type ADRBlockingRule struct{}

func (r *ADRBlockingRule) Name() string {
	return "sdd-adr-blocking"
}

func (r *ADRBlockingRule) Validate(ctx *validator.ValidationContext) (*validator.RuleResult, error) {
	result := &validator.RuleResult{
		RuleName: r.Name(),
		Status:   "PASS",
		Errors:   []string{},
	}

	specGraphPath := ".agentic/spec-graph.json"
	data, err := os.ReadFile(specGraphPath)
	if err != nil {
		if os.IsNotExist(err) {
			return result, nil // Graph doesn't exist, skip
		}
		return nil, fmt.Errorf("failed to read spec graph: %w", err)
	}

	var nodes map[string]sdd.SpecGraphNode
	if err := json.Unmarshal(data, &nodes); err != nil {
		return nil, fmt.Errorf("failed to parse spec graph: %w", err)
	}

	// Check for specs that are blocked but in Implementing or Done state
	for id, node := range nodes {
		if len(node.BlockedBy) > 0 {
			// Has blockers
			if node.Status == sdd.SpecStatusImplementing || node.Status == sdd.SpecStatusDone {
				result.Status = "FAIL"
				result.Errors = append(result.Errors,
					fmt.Sprintf("Spec %s is %s but still blocked by: %v", id, node.Status, node.BlockedBy))
			}
		}
	}

	return result, nil
}
