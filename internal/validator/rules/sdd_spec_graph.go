package rules

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/javierbenavides/agentic-agent/internal/sdd"
	"github.com/javierbenavides/agentic-agent/internal/validator"
)

type SpecGraphRule struct{}

func (r *SpecGraphRule) Name() string {
	return "sdd-spec-graph"
}

func (r *SpecGraphRule) Validate(ctx *validator.ValidationContext) (*validator.RuleResult, error) {
	result := &validator.RuleResult{
		RuleName: r.Name(),
		Status:   "PASS",
		Errors:   []string{},
	}

	changesDir := ".agentic/openspec/changes"
	if _, err := os.Stat(changesDir); err != nil {
		if os.IsNotExist(err) {
			return result, nil
		}
		return nil, fmt.Errorf("failed to stat changes directory: %w", err)
	}

	// Load spec graph
	specGraphPath := ".agentic/spec-graph.json"
	data, err := os.ReadFile(specGraphPath)
	if err != nil {
		if os.IsNotExist(err) {
			// No graph yet, that's OK for early stages
			return result, nil
		}
		return nil, fmt.Errorf("failed to read spec graph: %w", err)
	}

	var nodes map[string]sdd.SpecGraphNode
	if err := json.Unmarshal(data, &nodes); err != nil {
		return nil, fmt.Errorf("failed to parse spec graph: %w", err)
	}

	// Check that Approved, Implementing, and Done specs are in the graph
	entries, err := os.ReadDir(changesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read changes directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		changeID := entry.Name()
		metadataPath := filepath.Join(changesDir, changeID, "metadata.yaml")

		// Read status from metadata
		metaData, err := os.ReadFile(metadataPath)
		if err != nil {
			continue
		}

		var metadata map[string]interface{}
		if err := json.Unmarshal(metaData, &metadata); err != nil {
			continue
		}

		statusVal, ok := metadata["status"]
		if !ok {
			continue
		}
		status, _ := statusVal.(string)

		// Check if spec is in graph if status requires it
		if status == "Approved" || status == "Implementing" || status == "Done" {
			if _, inGraph := nodes[changeID]; !inGraph {
				result.Status = "FAIL"
				result.Errors = append(result.Errors,
					fmt.Sprintf("Spec %s with status %s not found in spec-graph.json", changeID, status))
			}
		}
	}

	return result, nil
}
