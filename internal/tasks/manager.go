package tasks

import (
	"fmt"
	"os"
	"path/filepath"
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

// CompleteTaskWithTracking marks a task as complete and logs progress
func (tm *TaskManager) CompleteTaskWithTracking(taskID string, learnings []string, filesChanged []string, threadURL string) error {
	// Find the task
	task, source, err := tm.FindTask(taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return fmt.Errorf("task %s not found", taskID)
	}

	// Move task to done
	if source != "done" {
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
		// Note: In a CLI context, this would prompt the user interactively
		// For now, we'll just prepare the infrastructure
		// The actual prompting will be done in the CLI command handler
		_ = dirs
	}

	return nil
}
