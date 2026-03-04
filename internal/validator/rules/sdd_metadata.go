package rules

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/javierbenavides/agentic-agent/internal/validator"
	"gopkg.in/yaml.v3"
)

type SpecMetadataRule struct{}

func (r *SpecMetadataRule) Name() string {
	return "sdd-metadata"
}

func (r *SpecMetadataRule) Validate(ctx *validator.ValidationContext) (*validator.RuleResult, error) {
	result := &validator.RuleResult{
		RuleName: r.Name(),
		Status:   "PASS",
		Errors:   []string{},
	}

	changesDir := ".agentic/openspec/changes"
	if _, err := os.Stat(changesDir); err != nil {
		if os.IsNotExist(err) {
			// Directory doesn't exist, skip
			return result, nil
		}
		return nil, fmt.Errorf("failed to stat changes directory: %w", err)
	}

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

		data, err := os.ReadFile(metadataPath)
		if err != nil {
			if os.IsNotExist(err) {
				result.Status = "FAIL"
				result.Errors = append(result.Errors,
					fmt.Sprintf("Change %s missing metadata.yaml", changeID))
				continue
			}
			return nil, fmt.Errorf("failed to read metadata for %s: %w", changeID, err)
		}

		var metadata map[string]interface{}
		if err := yaml.Unmarshal(data, &metadata); err != nil {
			result.Status = "FAIL"
			result.Errors = append(result.Errors,
				fmt.Sprintf("Change %s has invalid YAML metadata: %v", changeID, err))
			continue
		}

		// Check required fields
		requiredFields := []string{"implements", "context_pack", "blocked_by", "status"}
		for _, field := range requiredFields {
			if _, ok := metadata[field]; !ok {
				result.Status = "FAIL"
				result.Errors = append(result.Errors,
					fmt.Sprintf("Change %s missing metadata field: %s", changeID, field))
			}
		}

		// Check status is one of the valid values
		if statusVal, ok := metadata["status"]; ok {
			status, _ := statusVal.(string)
			validStatuses := []string{"Planned", "Draft", "Approved", "Implementing", "Done", "Paused", "Blocked"}
			found := false
			for _, valid := range validStatuses {
				if status == valid {
					found = true
					break
				}
			}
			if !found && status != "" {
				result.Status = "FAIL"
				result.Errors = append(result.Errors,
					fmt.Sprintf("Change %s has invalid status: %s", changeID, status))
			}
		}
	}

	return result, nil
}
