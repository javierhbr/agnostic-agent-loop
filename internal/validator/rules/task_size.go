package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/internal/validator"
)

const (
	MaxFilesPerTask       = 5
	MaxDirectoriesPerTask = 2
)

type TaskSizeRule struct{}

func (r *TaskSizeRule) Name() string {
	return "task-size"
}

func (r *TaskSizeRule) Validate(ctx *validator.ValidationContext) (*validator.RuleResult, error) {
	result := &validator.RuleResult{
		RuleName: r.Name(),
		Status:   "PASS",
		Errors:   []string{},
	}

	// Load in-progress tasks
	tm := tasks.NewTaskManager(".agentic/tasks")
	inProgress, err := tm.LoadTasks("in-progress")
	if err != nil {
		return nil, fmt.Errorf("failed to load in-progress tasks: %w", err)
	}

	// Check each task's size
	for _, task := range inProgress.Tasks {
		if len(task.Scope) == 0 {
			// No scope defined, skip validation
			continue
		}

		// Count files and directories in scope
		fileCount := 0
		dirSet := make(map[string]bool)

		for _, scopePath := range task.Scope {
			// Normalize path
			cleanPath := filepath.Clean(scopePath)

			// Check if it's a file or directory
			info, err := os.Stat(cleanPath)
			if err != nil {
				if os.IsNotExist(err) {
					// Path doesn't exist yet (might be an output), count it anyway
					if strings.Contains(filepath.Base(cleanPath), ".") {
						// Has extension, likely a file
						fileCount++
						dirSet[filepath.Dir(cleanPath)] = true
					} else {
						// No extension, likely a directory
						dirSet[cleanPath] = true
					}
					continue
				}
				return nil, fmt.Errorf("failed to stat %s: %w", cleanPath, err)
			}

			if info.IsDir() {
				// Count all files in directory
				err := filepath.Walk(cleanPath, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() && isSourceFile(path) {
						fileCount++
					}
					return nil
				})
				if err != nil {
					return nil, fmt.Errorf("failed to walk directory %s: %w", cleanPath, err)
				}
				dirSet[cleanPath] = true
			} else {
				// It's a file
				fileCount++
				dirSet[filepath.Dir(cleanPath)] = true
			}
		}

		dirCount := len(dirSet)

		// Check limits
		if fileCount > MaxFilesPerTask {
			result.Status = "FAIL"
			result.Errors = append(result.Errors,
				fmt.Sprintf("Task %s exceeds file limit: %d files (max %d). Consider decomposing into subtasks.",
					task.ID, fileCount, MaxFilesPerTask))
		}

		if dirCount > MaxDirectoriesPerTask {
			result.Status = "FAIL"
			result.Errors = append(result.Errors,
				fmt.Sprintf("Task %s exceeds directory limit: %d directories (max %d). Consider decomposing into subtasks.",
					task.ID, dirCount, MaxDirectoriesPerTask))
		}
	}

	return result, nil
}

// isSourceFile checks if a file is a source code file
func isSourceFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	sourceExts := []string{".go", ".ts", ".tsx", ".js", ".jsx", ".py", ".java", ".rs", ".c", ".cpp", ".h", ".hpp"}
	for _, se := range sourceExts {
		if ext == se {
			return true
		}
	}
	return false
}
