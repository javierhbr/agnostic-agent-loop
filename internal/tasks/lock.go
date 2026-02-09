package tasks

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// ClaimTask claims a task from backlog, recording claim time and git branch.
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
		return err
	}

	task.Status = models.StatusInProgress
	task.AssignedTo = assignee
	task.ClaimedAt = time.Now()
	task.Branch = currentGitBranch()
	inProgress.Tasks = append(inProgress.Tasks, task)

	return tm.SaveTasks("in-progress", inProgress)
}

// currentGitBranch returns the current git branch, or empty string if not in a repo.
func currentGitBranch() string {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// ClaimTaskWithConfig claims a task after running readiness checks.
// Readiness failures are printed as warnings but do not block the claim.
func (tm *TaskManager) ClaimTaskWithConfig(taskID, assignee string, cfg *models.Config) error {
	// Find the task to run readiness checks before claiming
	backlog, err := tm.LoadTasks("backlog")
	if err != nil {
		return err
	}

	var task *models.Task
	for _, t := range backlog.Tasks {
		if t.ID == taskID {
			task = &t
			break
		}
	}

	if task != nil {
		result := CanClaimTask(task, cfg)
		if len(result.Checks) > 0 {
			fmt.Print(FormatReadinessResult(result))
		}
	}

	return tm.ClaimTask(taskID, assignee)
}
