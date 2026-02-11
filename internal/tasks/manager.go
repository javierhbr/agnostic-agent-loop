package tasks

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/javierbenavides/agentic-agent/pkg/models"
	"gopkg.in/yaml.v3"
)

type TaskList struct {
	Tasks []models.Task `yaml:"tasks"`
}

type TaskManager struct {
	baseDir        string
	progressWriter *ProgressWriter
	agentsMdHelper *AgentsMdHelper
}

func NewTaskManager(baseDir string) *TaskManager {
	return &TaskManager{baseDir: baseDir}
}

func NewTaskManagerWithTracking(baseDir string, progressWriter *ProgressWriter, agentsMdHelper *AgentsMdHelper) *TaskManager {
	return &TaskManager{
		baseDir:        baseDir,
		progressWriter: progressWriter,
		agentsMdHelper: agentsMdHelper,
	}
}

func (tm *TaskManager) LoadTasks(listType string) (*TaskList, error) {
	path := filepath.Join(tm.baseDir, listType+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &TaskList{Tasks: []models.Task{}}, nil
		}
		return nil, err
	}

	var list TaskList
	if err := yaml.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (tm *TaskManager) SaveTasks(listType string, list *TaskList) error {
	path := filepath.Join(tm.baseDir, listType+".yaml")
	data, err := yaml.Marshal(list)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (tm *TaskManager) CreateTask(title string) (*models.Task, error) {
	backlog, err := tm.LoadTasks("backlog")
	if err != nil {
		return nil, err
	}

	// Simple ID generation (can be improved)
	id := fmt.Sprintf("TASK-%d", time.Now().Unix())

	task := models.Task{
		ID:     id,
		Title:  title,
		Status: models.StatusPending,
	}

	backlog.Tasks = append(backlog.Tasks, task)
	if err := tm.SaveTasks("backlog", backlog); err != nil {
		return nil, err
	}

	return &task, nil
}

func (tm *TaskManager) MoveTask(taskID string, fromType, toType string, newStatus models.TaskStatus) error {
	fromList, err := tm.LoadTasks(fromType)
	if err != nil {
		return err
	}

	var taskToMove models.Task
	found := false
	newFromTasks := []models.Task{}

	for _, t := range fromList.Tasks {
		if t.ID == taskID {
			taskToMove = t
			found = true
		} else {
			newFromTasks = append(newFromTasks, t)
		}
	}

	if !found {
		return fmt.Errorf("task %s not found in %s", taskID, fromType)
	}

	fromList.Tasks = newFromTasks
	if err := tm.SaveTasks(fromType, fromList); err != nil {
		return err
	}

	toList, err := tm.LoadTasks(toType)
	if err != nil {
		return err
	}

	taskToMove.Status = newStatus
	toList.Tasks = append(toList.Tasks, taskToMove)
	return tm.SaveTasks(toType, toList)
}

// FindTask searches for a task across all lists (backlog, in-progress, done)
// Returns the task, the source list name, and an error if any
func (tm *TaskManager) FindTask(taskID string) (*models.Task, string, error) {
	sources := []string{"backlog", "in-progress", "done"}

	for _, source := range sources {
		list, err := tm.LoadTasks(source)
		if err != nil {
			return nil, "", fmt.Errorf("error loading %s: %w", source, err)
		}

		for _, task := range list.Tasks {
			if task.ID == taskID {
				return &task, source, nil
			}
			// Also check subtasks
			for _, subtask := range task.SubTasks {
				if subtask.ID == taskID {
					// Convert SubTask to Task for consistent return
					fullTask := models.Task{
						ID:         subtask.ID,
						Title:      subtask.Title,
						Status:     subtask.Status,
						AssignedTo: subtask.AssignedTo,
					}
					return &fullTask, source, nil
				}
			}
		}
	}

	return nil, "", nil // Not found, but no error
}

// CompleteTaskWithTracking marks a task as complete and logs progress.
// If the task has a ClaimedAt timestamp, git commits since that time are auto-captured.
func (tm *TaskManager) CompleteTaskWithTracking(taskID string, learnings []string, filesChanged []string, threadURL string) error {
	// Find the task
	task, source, err := tm.FindTask(taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return fmt.Errorf("task %s not found", taskID)
	}

	// Auto-populate git data if ClaimedAt is set and we're in a git repo
	if !task.ClaimedAt.IsZero() {
		task.CompletedAt = time.Now()
		gitCommits, gitFiles := collectGitData(task.ClaimedAt)
		if len(gitCommits) > 0 {
			task.Commits = gitCommits
		}
		if len(filesChanged) == 0 && len(gitFiles) > 0 {
			filesChanged = gitFiles
		}
	}

	// Move task to done (with updated fields)
	if source != "done" {
		// Update the task in its current list before moving
		tm.updateTaskInList(taskID, source, func(t *models.Task) {
			t.CompletedAt = task.CompletedAt
			t.Commits = task.Commits
		})
		if err := tm.MoveTask(taskID, source, "done", models.StatusDone); err != nil {
			return fmt.Errorf("failed to move task: %w", err)
		}
	}

	// If progress tracking is enabled, log the completion
	if tm.progressWriter != nil {
		entry := ProgressEntry{
			Timestamp:    time.Now(),
			StoryID:      taskID,
			Title:        task.Title,
			FilesChanged: filesChanged,
			Learnings:    learnings,
			ThreadURL:    threadURL,
		}
		if err := tm.progressWriter.AppendEntry(entry); err != nil {
			return fmt.Errorf("failed to write progress: %w", err)
		}
	}

	// If AGENTS.md helper is enabled and there are learnings, prompt for directory updates
	if tm.agentsMdHelper != nil && len(learnings) > 0 && len(filesChanged) > 0 {
		dirs := tm.agentsMdHelper.GetModifiedDirectories(filesChanged)
		_ = dirs
	}

	return nil
}

// updateTaskInList modifies a task in place within a list.
func (tm *TaskManager) updateTaskInList(taskID, listType string, update func(*models.Task)) {
	list, err := tm.LoadTasks(listType)
	if err != nil {
		return
	}
	for i := range list.Tasks {
		if list.Tasks[i].ID == taskID {
			update(&list.Tasks[i])
			tm.SaveTasks(listType, list)
			return
		}
	}
}

// collectGitData gathers commit hashes and changed files since the given time.
func collectGitData(since time.Time) (commits []string, files []string) {
	sinceStr := since.Format("2006-01-02T15:04:05")

	// Get commit hashes
	out, err := execGit("log", "--since="+sinceStr, "--format=%H", "--no-merges")
	if err != nil {
		return nil, nil
	}
	for _, line := range splitLines(out) {
		if line != "" {
			commits = append(commits, line)
		}
	}

	// Get files changed
	out, err = execGit("log", "--since="+sinceStr, "--name-only", "--format=", "--no-merges")
	if err != nil {
		return commits, nil
	}
	seen := make(map[string]bool)
	for _, line := range splitLines(out) {
		if line != "" && !seen[line] {
			seen[line] = true
			files = append(files, line)
		}
	}

	return commits, files
}

func execGit(args ...string) (string, error) {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func splitLines(s string) []string {
	var lines []string
	for _, line := range strings.Split(strings.TrimSpace(s), "\n") {
		lines = append(lines, strings.TrimSpace(line))
	}
	return lines
}
