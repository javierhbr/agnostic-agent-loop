package tasks

import (
	"fmt"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// TDD phase identifiers used in sub-task IDs.
const (
	TDDPhaseRed      = "red"
	TDDPhaseGreen    = "green"
	TDDPhaseRefactor = "refactor"
)

// DecomposeForTDD creates RED, GREEN, and REFACTOR sub-tasks on the parent task.
// The parent task is updated in its current list with the new sub-tasks.
func DecomposeForTDD(tm *TaskManager, parentTaskID string) ([]models.SubTask, error) {
	task, source, err := tm.FindTask(parentTaskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("task %s not found", parentTaskID)
	}

	subtasks := []models.SubTask{
		{
			ID:     fmt.Sprintf("%s-%s", parentTaskID, TDDPhaseRed),
			Title:  fmt.Sprintf("[RED] Write failing tests for: %s", task.Title),
			Status: models.StatusPending,
		},
		{
			ID:     fmt.Sprintf("%s-%s", parentTaskID, TDDPhaseGreen),
			Title:  fmt.Sprintf("[GREEN] Implement minimal code to pass tests for: %s", task.Title),
			Status: models.StatusPending,
		},
		{
			ID:     fmt.Sprintf("%s-%s", parentTaskID, TDDPhaseRefactor),
			Title:  fmt.Sprintf("[REFACTOR] Refactor implementation for: %s", task.Title),
			Status: models.StatusPending,
		},
	}

	// Update the parent task with sub-tasks
	tm.updateTaskInList(parentTaskID, source, func(t *models.Task) {
		t.SubTasks = subtasks
	})

	return subtasks, nil
}
