package tasks

import (
	"fmt"

	"github.com/javierbenavides/agentic-agent/internal/plans"
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

// DecomposeFromPlan parses a plan.md and creates one task per phase,
// linking each task to the given track.
func DecomposeFromPlan(planPath, trackID string, tm *TaskManager) ([]*models.Task, error) {
	plan, err := plans.ParseFile(planPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse plan: %w", err)
	}

	if len(plan.Phases) == 0 {
		return nil, fmt.Errorf("plan has no phases")
	}

	var created []*models.Task
	for _, phase := range plan.Phases {
		if len(phase.Tasks) == 0 {
			continue
		}

		title := fmt.Sprintf("[%s] %s", trackID, phase.Name)
		task, err := tm.CreateTask(title)
		if err != nil {
			return created, fmt.Errorf("failed to create task for phase %q: %w", phase.Name, err)
		}

		desc := fmt.Sprintf("Phase from track %s:\n", trackID)
		var acceptance []string
		for _, pt := range phase.Tasks {
			desc += fmt.Sprintf("- %s\n", pt.Title)
			acceptance = append(acceptance, pt.Title)
		}

		task.Description = desc
		task.Acceptance = acceptance
		task.TrackID = trackID

		// Save updated task back to backlog
		backlog, err := tm.LoadTasks("backlog")
		if err != nil {
			return created, fmt.Errorf("failed to load backlog: %w", err)
		}
		for i, t := range backlog.Tasks {
			if t.ID == task.ID {
				backlog.Tasks[i] = *task
				break
			}
		}
		if err := tm.SaveTasks("backlog", backlog); err != nil {
			return created, fmt.Errorf("failed to save task: %w", err)
		}

		created = append(created, task)
	}

	return created, nil
}
