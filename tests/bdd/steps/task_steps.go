package steps

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/cucumber/godog"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/javierbenavides/agentic-agent/tests/functional"
)

// TaskSteps encapsulates task-related step definitions
type TaskSteps struct {
	suite *SuiteContext
}

// NewTaskSteps creates a new TaskSteps instance
func NewTaskSteps(suite *SuiteContext) *TaskSteps {
	return &TaskSteps{suite: suite}
}

// RegisterSteps registers all task-related step definitions
func (s *TaskSteps) RegisterSteps(sc *godog.ScenarioContext) {
	// Task creation steps
	sc.Step(`^I create a task with title "([^"]*)"$`, s.createTaskWithTitle)
	sc.Step(`^I have created a task with title "([^"]*)"$`, s.createTaskWithTitle)
	sc.Step(`^a task should be created successfully$`, s.taskShouldBeCreated)

	// Task lifecycle steps
	sc.Step(`^I claim the task$`, s.claimCurrentTask)
	sc.Step(`^I complete the task$`, s.completeCurrentTask)
	sc.Step(`^I list all tasks$`, s.listAllTasks)

	// Task state assertions
	sc.Step(`^the task should be in "([^"]*)" state$`, s.assertTaskInState)
	sc.Step(`^the task should appear in the backlog$`, s.assertTaskInBacklog)
	sc.Step(`^the task should move to in-progress$`, s.assertTaskInProgress)
	sc.Step(`^the task should move to done$`, s.assertTaskInDone)
	sc.Step(`^the task status should be "([^"]*)"$`, s.assertTaskStatus)

	// Task counting assertions
	sc.Step(`^I should see (\d+) tasks? in backlog$`, s.assertTaskCountInBacklog)
	sc.Step(`^I should see (\d+) tasks? in progress$`, s.assertTaskCountInProgress)
}

// createTaskWithTitle creates a task with the specified title
func (s *TaskSteps) createTaskWithTitle(ctx context.Context, title string) error {
	tm := tasks.NewTaskManager(filepath.Join(s.suite.ProjectDir, ".agentic/tasks"))
	task, err := tm.CreateTask(title)
	if err != nil {
		s.suite.LastCommandErr = err
		return err
	}

	s.suite.CurrentTask = task
	s.suite.LastTaskID = task.ID
	s.suite.LastCommandErr = nil
	return nil
}

// taskShouldBeCreated asserts that a task was created successfully
func (s *TaskSteps) taskShouldBeCreated(ctx context.Context) error {
	if s.suite.LastTaskID == "" {
		return fmt.Errorf("no task was created")
	}
	return nil
}

// claimCurrentTask claims the current task
func (s *TaskSteps) claimCurrentTask(ctx context.Context) error {
	if s.suite.LastTaskID == "" {
		return fmt.Errorf("no task to claim")
	}

	tm := tasks.NewTaskManager(filepath.Join(s.suite.ProjectDir, ".agentic/tasks"))
	err := tm.ClaimTask(s.suite.LastTaskID, "test-agent")
	s.suite.LastCommandErr = err
	return err
}

// completeCurrentTask completes the current task
func (s *TaskSteps) completeCurrentTask(ctx context.Context) error {
	if s.suite.LastTaskID == "" {
		return fmt.Errorf("no task to complete")
	}

	tm := tasks.NewTaskManager(filepath.Join(s.suite.ProjectDir, ".agentic/tasks"))
	err := tm.MoveTask(s.suite.LastTaskID, "in-progress", "done", models.StatusDone)
	s.suite.LastCommandErr = err
	return err
}

// listAllTasks lists all tasks (stores count for assertions)
func (s *TaskSteps) listAllTasks(ctx context.Context) error {
	tm := tasks.NewTaskManager(filepath.Join(s.suite.ProjectDir, ".agentic/tasks"))

	// Load tasks from all states
	backlog, err := tm.LoadTasks("backlog")
	if err != nil {
		return err
	}

	inProgress, err := tm.LoadTasks("in-progress")
	if err != nil {
		return err
	}

	done, err := tm.LoadTasks("done")
	if err != nil {
		return err
	}

	// Store for later assertions
	s.suite.LastCommandOut = fmt.Sprintf("Backlog: %d, In Progress: %d, Done: %d",
		len(backlog.Tasks), len(inProgress.Tasks), len(done.Tasks))

	return nil
}

// assertTaskInState asserts that the task is in the specified state
func (s *TaskSteps) assertTaskInState(ctx context.Context, state string) error {
	taskFile := filepath.Join(s.suite.ProjectDir, ".agentic/tasks", state+".yaml")
	task := functional.VerifyTaskInFile(s.suite.T, taskFile, s.suite.LastTaskID)
	if task == nil {
		return fmt.Errorf("task %s not found in %s state", s.suite.LastTaskID, state)
	}
	return nil
}

// assertTaskInBacklog asserts that the task is in backlog
func (s *TaskSteps) assertTaskInBacklog(ctx context.Context) error {
	return s.assertTaskInState(ctx, "backlog")
}

// assertTaskInProgress asserts that the task moved to in-progress
func (s *TaskSteps) assertTaskInProgress(ctx context.Context) error {
	return s.assertTaskInState(ctx, "in-progress")
}

// assertTaskInDone asserts that the task moved to done
func (s *TaskSteps) assertTaskInDone(ctx context.Context) error {
	return s.assertTaskInState(ctx, "done")
}

// assertTaskStatus asserts the task has the specified status
func (s *TaskSteps) assertTaskStatus(ctx context.Context, expectedStatus string) error {
	tm := tasks.NewTaskManager(filepath.Join(s.suite.ProjectDir, ".agentic/tasks"))
	task, _, err := tm.FindTask(s.suite.LastTaskID)
	if err != nil {
		return fmt.Errorf("failed to find task: %w", err)
	}

	if string(task.Status) != expectedStatus {
		return fmt.Errorf("expected task status %q, got %q", expectedStatus, task.Status)
	}

	return nil
}

// assertTaskCountInBacklog asserts the number of tasks in backlog
func (s *TaskSteps) assertTaskCountInBacklog(ctx context.Context, expectedCount int) error {
	tm := tasks.NewTaskManager(filepath.Join(s.suite.ProjectDir, ".agentic/tasks"))
	backlog, err := tm.LoadTasks("backlog")
	if err != nil {
		return err
	}

	if len(backlog.Tasks) != expectedCount {
		return fmt.Errorf("expected %d tasks in backlog, got %d", expectedCount, len(backlog.Tasks))
	}

	return nil
}

// assertTaskCountInProgress asserts the number of tasks in progress
func (s *TaskSteps) assertTaskCountInProgress(ctx context.Context, expectedCount int) error {
	tm := tasks.NewTaskManager(filepath.Join(s.suite.ProjectDir, ".agentic/tasks"))
	inProgress, err := tm.LoadTasks("in-progress")
	if err != nil {
		return err
	}

	if len(inProgress.Tasks) != expectedCount {
		return fmt.Errorf("expected %d tasks in progress, got %d", expectedCount, len(inProgress.Tasks))
	}

	return nil
}
