package tasks

import (
	"fmt"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

func (tm *TaskManager) DecomposeTask(taskID string, subtasks []string) error {
	// Task must be in In-Progress to decompose? Or any state?
	// Usually In-Progress or Pending. Let's look in In-Progress first, then Backlog.

	listType := "in-progress"
	list, err := tm.LoadTasks(listType)
	if err != nil {
		return err
	}

	foundIdx := -1
	for i, t := range list.Tasks {
		if t.ID == taskID {
			foundIdx = i
			break
		}
	}

	if foundIdx == -1 {
		// Check backlog
		listType = "backlog"
		list, err = tm.LoadTasks(listType)
		if err != nil {
			return err
		}
		for i, t := range list.Tasks {
			if t.ID == taskID {
				foundIdx = i
				break
			}
		}
	}

	if foundIdx == -1 {
		return fmt.Errorf("task %s not found", taskID)
	}

	// Add subtasks
	for i, title := range subtasks {
		subID := fmt.Sprintf("%s.%d", taskID, i+1)
		list.Tasks[foundIdx].SubTasks = append(list.Tasks[foundIdx].SubTasks, models.SubTask{
			ID:     subID,
			Title:  title,
			Status: models.StatusPending,
		})
	}

	return tm.SaveTasks(listType, list)
}
