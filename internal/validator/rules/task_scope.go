package rules

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/internal/validator"
)

type TaskScopeRule struct{}

func (r *TaskScopeRule) Name() string {
	return "task-scope"
}

func (r *TaskScopeRule) Validate(ctx *validator.ValidationContext) (*validator.RuleResult, error) {
	result := &validator.RuleResult{
		RuleName: r.Name(),
		Status:   "PASS",
		Errors:   []string{},
	}

	// Check if we're in a git repository
	if !isGitRepo() {
		// Not in a git repo, skip validation (can't check modified files)
		return result, nil
	}

	// Get modified files from git
	modifiedFiles, err := getModifiedFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to get modified files: %w", err)
	}

	if len(modifiedFiles) == 0 {
		// No modified files, nothing to validate
		return result, nil
	}

	// Load in-progress tasks
	tm := tasks.NewTaskManager(".agentic/tasks")
	inProgress, err := tm.LoadTasks("in-progress")
	if err != nil {
		return nil, fmt.Errorf("failed to load in-progress tasks: %w", err)
	}

	if len(inProgress.Tasks) == 0 {
		// No tasks in progress but files are modified
		result.Status = "FAIL"
		result.Errors = append(result.Errors,
			fmt.Sprintf("Modified files detected but no task is in progress. Claim a task before making changes."))
		return result, nil
	}

	// Check each modified file against in-progress task scopes
	for _, modFile := range modifiedFiles {
		fileInScope := false

		for _, task := range inProgress.Tasks {
			if len(task.Scope) == 0 {
				// Task has no scope restrictions (allowed for now)
				fileInScope = true
				break
			}

			// Check if modified file is within task scope
			for _, scopePath := range task.Scope {
				if isFileInScope(modFile, scopePath) {
					fileInScope = true
					break
				}
			}

			if fileInScope {
				break
			}
		}

		if !fileInScope {
			result.Status = "FAIL"
			result.Errors = append(result.Errors,
				fmt.Sprintf("File '%s' modified but not in scope of any in-progress task", modFile))
		}
	}

	return result, nil
}

// isGitRepo checks if the current directory is a git repository
func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// getModifiedFiles returns a list of modified files from git status
func getModifiedFiles() ([]string, error) {
	// Get both staged and unstaged modifications
	cmd := exec.Command("git", "diff", "--name-only", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		// Try just unstaged changes
		cmd = exec.Command("git", "diff", "--name-only")
		output, err = cmd.Output()
		if err != nil {
			return nil, err
		}
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	result := make([]string, 0, len(files))
	for _, f := range files {
		if f != "" {
			result = append(result, f)
		}
	}

	return result, nil
}

// isFileInScope checks if a file path is within the given scope path
func isFileInScope(filePath, scopePath string) bool {
	// Normalize paths
	filePath = filepath.Clean(filePath)
	scopePath = filepath.Clean(scopePath)

	// If scope is a specific file, check exact match
	if strings.HasSuffix(scopePath, filepath.Ext(scopePath)) {
		return filePath == scopePath
	}

	// If scope is a directory, check if file is within it
	// Make sure scopePath ends without separator for consistent matching
	scopePath = strings.TrimSuffix(scopePath, string(filepath.Separator))

	// Check if file is in the directory or subdirectory
	return strings.HasPrefix(filePath, scopePath+string(filepath.Separator)) || filePath == scopePath
}
