package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"

	"github.com/javierbenavides/agentic-agent/internal/validator"
)

type VerifyMdRule struct{}

func (r *VerifyMdRule) Name() string {
	return "sdd-verify-md"
}

func (r *VerifyMdRule) Validate(ctx *validator.ValidationContext) (*validator.RuleResult, error) {
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
				continue
			}
			return nil, fmt.Errorf("failed to read metadata for %s: %w", changeID, err)
		}

		var metadata map[string]interface{}
		if err := yaml.Unmarshal(data, &metadata); err != nil {
			continue
		}

		// Check status
		statusVal, ok := metadata["status"]
		if !ok {
			continue
		}
		status, _ := statusVal.(string)

		// If status is Done, verify.md must exist
		if status == "Done" {
			verifyPath := filepath.Join(changesDir, changeID, "verify.md")
			if _, err := os.Stat(verifyPath); err != nil {
				if os.IsNotExist(err) {
					result.Status = "FAIL"
					result.Errors = append(result.Errors,
						fmt.Sprintf("Change %s is Done but missing verify.md", changeID))
				} else {
					return nil, fmt.Errorf("failed to stat verify.md for %s: %w", changeID, err)
				}
			}
		}
	}

	return result, nil
}
