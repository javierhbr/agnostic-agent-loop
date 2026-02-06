package rules

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/internal/validator"
)

type BrowserVerificationRule struct{}

func (r *BrowserVerificationRule) Name() string {
	return "browser-verification"
}

func (r *BrowserVerificationRule) Validate(ctx *validator.ValidationContext) (*validator.RuleResult, error) {
	var failures []string

	// Load in-progress tasks
	tm := tasks.NewTaskManager(filepath.Join(ctx.ProjectRoot, ".agentic/tasks"))
	inProgressList, err := tm.LoadTasks("in-progress")
	if err != nil {
		return nil, fmt.Errorf("failed to load in-progress tasks: %w", err)
	}

	// Check each in-progress task
	for _, task := range inProgressList.Tasks {
		// Check if task modifies UI files
		hasUIChanges := false
		for _, file := range task.Scope {
			if r.isUIFile(file) {
				hasUIChanges = true
				break
			}
		}

		if hasUIChanges {
			// Check acceptance criteria for browser verification
			hasBrowserCriteria := false
			for _, criterion := range task.Acceptance {
				lowerCriterion := strings.ToLower(criterion)
				if strings.Contains(lowerCriterion, "verify in browser") ||
					strings.Contains(lowerCriterion, "browser verification") ||
					strings.Contains(lowerCriterion, "visual verification") ||
					strings.Contains(lowerCriterion, "test in browser") {
					hasBrowserCriteria = true
					break
				}
			}

			if !hasBrowserCriteria {
				failures = append(failures, fmt.Sprintf(
					"Task %s (%s) modifies UI files but lacks browser verification in acceptance criteria",
					task.ID, task.Title,
				))
			}
		}
	}

	status := "PASS"
	if len(failures) > 0 {
		status = "FAIL"
	}

	return &validator.RuleResult{
		RuleName: r.Name(),
		Status:   status,
		Errors:   failures,
	}, nil
}

// isUIFile checks if a file path is a UI/frontend file
func (r *BrowserVerificationRule) isUIFile(filePath string) bool {
	ext := filepath.Ext(filePath)

	// Common UI file extensions
	uiExtensions := map[string]bool{
		".tsx":   true,
		".jsx":   true,
		".vue":   true,
		".svelte": true,
		".html":  true,
		".css":   true,
		".scss":  true,
		".sass":  true,
		".less":  true,
	}

	if uiExtensions[ext] {
		return true
	}

	// Check for .ts/.js files in UI directories
	if ext == ".ts" || ext == ".js" {
		lowerPath := strings.ToLower(filePath)
		uiDirs := []string{
			"/components/",
			"/ui/",
			"/views/",
			"/pages/",
			"/layouts/",
			"/templates/",
			"/widgets/",
		}

		for _, dir := range uiDirs {
			if strings.Contains(lowerPath, dir) {
				return true
			}
		}
	}

	return false
}

// hasBrowserVerification checks if acceptance criteria includes browser verification
func (r *BrowserVerificationRule) hasBrowserVerification(acceptanceCriteria []string) bool {
	keywords := []string{
		"verify in browser",
		"browser verification",
		"visual verification",
		"test in browser",
	}

	for _, criterion := range acceptanceCriteria {
		lowerCriterion := strings.ToLower(criterion)
		for _, keyword := range keywords {
			if strings.Contains(lowerCriterion, keyword) {
				return true
			}
		}
	}

	return false
}
