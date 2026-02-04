package tasks

import (
	"fmt"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// AcquireLock attempts to lock a task by checking if it is already assigned/locked.
// For MVP, we will assume "In Progress" means locked by the assignee.
// True file locking can be added later if concurrent access is a major issue.
// This function mainly updates the assigned_to field.
func (tm *TaskManager) ClaimTask(taskID string, assignee string) error {
	// Find task in backlog
	backlog, err := tm.LoadTasks("backlog")
	if err != nil {
		return err
	}

	var task models.Task
	found := false
	for _, t := range backlog.Tasks {
		if t.ID == taskID {
			task = t
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task %s not found in backlog (can only claim pending tasks)", taskID)
	}

	// Move to In Progress
	// We need to slightly modify MoveTask to allow updating properties during move,
	// Or just do it manually here.

	// Better: Use MoveTask then update? OR update then save?
	// MoveTask does remove and add.

	// Let's implement specific logic here reusing Load/Save to be safe.

	newBacklogTasks := []models.Task{}
	for _, t := range backlog.Tasks {
		if t.ID != taskID {
			newBacklogTasks = append(newBacklogTasks, t)
		}
	}
	backlog.Tasks = newBacklogTasks
	if err := tm.SaveTasks("backlog", backlog); err != nil {
		return err
	}

	inProgress, err := tm.LoadTasks("in-progress")
	if err != nil {
		// Try to rollback? Not implementing full rollback for MVP
		return err
	}

	task.Status = models.StatusInProgress
	task.AssignedTo = assignee
	inProgress.Tasks = append(inProgress.Tasks, task)

	return tm.SaveTasks("in-progress", inProgress)
}
