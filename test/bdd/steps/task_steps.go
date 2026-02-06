package steps

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/cucumber/godog"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/javierbenavides/agentic-agent/test/functional"
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
	sc.Step(`^the task should not exist in backlog$`, s.assertTaskNotInBacklog)
	sc.Step(`^the task should not exist in "([^"]*)" state$`, s.assertTaskNotInState)

	// Task counting assertions
	sc.Step(`^I should see (\d+) tasks? in backlog$`, s.assertTaskCountInBacklog)
	sc.Step(`^I should see (\d+) tasks? in progress$`, s.assertTaskCountInProgress)

	// Task metadata steps
	sc.Step(`^I set task description to "([^"]*)"$`, s.setTaskDescription)
	sc.Step(`^I add the following acceptance criteria:$`, s.addAcceptanceCriteria)
	sc.Step(`^I add the following outputs:$`, s.addOutputs)
	sc.Step(`^the task should have (\d+) acceptance criteria$`, s.assertAcceptanceCriteriaCount)
	sc.Step(`^the task should preserve all acceptance criteria$`, s.assertAcceptanceCriteriaPreserved)
	sc.Step(`^the task should preserve all outputs$`, s.assertOutputsPreserved)

	// File creation steps
	sc.Step(`^I create the following files:$`, s.createFiles)

	// Task decomposition steps
	sc.Step(`^I decompose the task into the following subtasks:$`, s.decomposeTask)
	sc.Step(`^the task should have (\d+) subtasks$`, s.assertSubtaskCount)
	sc.Step(`^the task should remain in backlog$`, s.assertTaskInBacklog)

	// Error handling steps
	sc.Step(`^I try to claim task "([^"]*)"$`, s.tryClaimTaskByID)
	sc.Step(`^I try to complete task without claiming$`, s.tryCompleteWithoutClaim)
	sc.Step(`^I try to create a task with empty title$`, s.tryCreateEmptyTask)
	sc.Step(`^the task should still be in backlog$`, s.assertTaskInBacklog)
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

// setTaskDescription sets the description for the current task
func (s *TaskSteps) setTaskDescription(ctx context.Context, description string) error {
	if s.suite.CurrentTask == nil {
		return fmt.Errorf("no current task to set description")
	}
	s.suite.CurrentTask.Description = description
	return nil
}

// addAcceptanceCriteria adds acceptance criteria to the current task
func (s *TaskSteps) addAcceptanceCriteria(ctx context.Context, table *godog.Table) error {
	if s.suite.CurrentTask == nil {
		return fmt.Errorf("no current task to add acceptance criteria")
	}

	for _, row := range table.Rows[1:] { // Skip header
		criterion := row.Cells[0].Value
		s.suite.CurrentTask.Acceptance = append(s.suite.CurrentTask.Acceptance, criterion)
	}

	// Save the task with updated acceptance criteria
	tm := tasks.NewTaskManager(filepath.Join(s.suite.ProjectDir, ".agentic/tasks"))
	backlog, err := tm.LoadTasks("backlog")
	if err != nil {
		return err
	}

	// Find and update the task
	for i, task := range backlog.Tasks {
		if task.ID == s.suite.CurrentTask.ID {
			backlog.Tasks[i].Acceptance = s.suite.CurrentTask.Acceptance
			backlog.Tasks[i].Description = s.suite.CurrentTask.Description
			break
		}
	}

	return tm.SaveTasks("backlog", backlog)
}

// addOutputs adds expected outputs to the current task
func (s *TaskSteps) addOutputs(ctx context.Context, table *godog.Table) error {
	if s.suite.CurrentTask == nil {
		return fmt.Errorf("no current task to add outputs")
	}

	for _, row := range table.Rows[1:] { // Skip header
		output := row.Cells[0].Value
		s.suite.CurrentTask.Outputs = append(s.suite.CurrentTask.Outputs, output)
	}

	// Save the task with updated outputs
	tm := tasks.NewTaskManager(filepath.Join(s.suite.ProjectDir, ".agentic/tasks"))
	backlog, err := tm.LoadTasks("backlog")
	if err != nil {
		return err
	}

	// Find and update the task
	for i, task := range backlog.Tasks {
		if task.ID == s.suite.CurrentTask.ID {
			backlog.Tasks[i].Outputs = s.suite.CurrentTask.Outputs
			break
		}
	}

	return tm.SaveTasks("backlog", backlog)
}

// assertAcceptanceCriteriaCount asserts the number of acceptance criteria
func (s *TaskSteps) assertAcceptanceCriteriaCount(ctx context.Context, expectedCount int) error {
	if s.suite.CurrentTask == nil {
		return fmt.Errorf("no current task")
	}

	if len(s.suite.CurrentTask.Acceptance) != expectedCount {
		return fmt.Errorf("expected %d acceptance criteria, got %d", expectedCount, len(s.suite.CurrentTask.Acceptance))
	}

	return nil
}

// assertAcceptanceCriteriaPreserved asserts that acceptance criteria are preserved
func (s *TaskSteps) assertAcceptanceCriteriaPreserved(ctx context.Context) error {
	tm := tasks.NewTaskManager(filepath.Join(s.suite.ProjectDir, ".agentic/tasks"))
	task, _, err := tm.FindTask(s.suite.LastTaskID)
	if err != nil {
		return fmt.Errorf("failed to find task: %w", err)
	}

	if len(task.Acceptance) == 0 {
		return fmt.Errorf("acceptance criteria not preserved")
	}

	return nil
}

// assertOutputsPreserved asserts that outputs are preserved
func (s *TaskSteps) assertOutputsPreserved(ctx context.Context) error {
	tm := tasks.NewTaskManager(filepath.Join(s.suite.ProjectDir, ".agentic/tasks"))
	task, _, err := tm.FindTask(s.suite.LastTaskID)
	if err != nil {
		return fmt.Errorf("failed to find task: %w", err)
	}

	if len(task.Outputs) == 0 {
		return fmt.Errorf("outputs not preserved")
	}

	return nil
}

// createFiles creates test files with specified content
func (s *TaskSteps) createFiles(ctx context.Context, table *godog.Table) error {
	for _, row := range table.Rows[1:] { // Skip header
		filePath := filepath.Join(s.suite.ProjectDir, row.Cells[0].Value)
		content := row.Cells[1].Value

		functional.CreateTestFile(s.suite.T, filePath, content)
	}

	return nil
}

// decomposeTask decomposes the current task into subtasks
func (s *TaskSteps) decomposeTask(ctx context.Context, table *godog.Table) error {
	if s.suite.LastTaskID == "" {
		return fmt.Errorf("no task to decompose")
	}

	subtaskTitles := make([]string, 0)
	for _, row := range table.Rows[1:] { // Skip header
		subtaskTitles = append(subtaskTitles, row.Cells[0].Value)
	}

	tm := tasks.NewTaskManager(filepath.Join(s.suite.ProjectDir, ".agentic/tasks"))
	err := tm.DecomposeTask(s.suite.LastTaskID, subtaskTitles)
	s.suite.LastCommandErr = err
	return err
}

// assertSubtaskCount asserts the number of subtasks
func (s *TaskSteps) assertSubtaskCount(ctx context.Context, expectedCount int) error {
	tm := tasks.NewTaskManager(filepath.Join(s.suite.ProjectDir, ".agentic/tasks"))
	task, _, err := tm.FindTask(s.suite.LastTaskID)
	if err != nil {
		return fmt.Errorf("failed to find task: %w", err)
	}

	if len(task.SubTasks) != expectedCount {
		return fmt.Errorf("expected %d subtasks, got %d", expectedCount, len(task.SubTasks))
	}

	return nil
}

// assertTaskNotInBacklog asserts that the task is not in backlog
func (s *TaskSteps) assertTaskNotInBacklog(ctx context.Context) error {
	return s.assertTaskNotInState(ctx, "backlog")
}

// assertTaskNotInState asserts that the task is not in the specified state
func (s *TaskSteps) assertTaskNotInState(ctx context.Context, state string) error {
	tm := tasks.NewTaskManager(filepath.Join(s.suite.ProjectDir, ".agentic/tasks"))
	taskList, err := tm.LoadTasks(state)
	if err != nil {
		return err
	}

	// Check that task is NOT in this state
	for _, task := range taskList.Tasks {
		if task.ID == s.suite.LastTaskID {
			return fmt.Errorf("task %s should not be in %s state", s.suite.LastTaskID, state)
		}
	}

	return nil
}

// tryClaimTaskByID attempts to claim a task by ID (may fail)
func (s *TaskSteps) tryClaimTaskByID(ctx context.Context, taskID string) error {
	tm := tasks.NewTaskManager(filepath.Join(s.suite.ProjectDir, ".agentic/tasks"))
	err := tm.ClaimTask(taskID, "test-agent")
	s.suite.LastCommandErr = err
	// Don't return error - we want to test failures
	return nil
}

// tryCompleteWithoutClaim attempts to complete a task without claiming it
func (s *TaskSteps) tryCompleteWithoutClaim(ctx context.Context) error {
	if s.suite.LastTaskID == "" {
		return fmt.Errorf("no task to complete")
	}

	tm := tasks.NewTaskManager(filepath.Join(s.suite.ProjectDir, ".agentic/tasks"))
	err := tm.MoveTask(s.suite.LastTaskID, "backlog", "done", models.StatusDone)
	s.suite.LastCommandErr = err
	// Don't return error - we want to test failures
	return nil
}

// tryCreateEmptyTask attempts to create a task with empty title
func (s *TaskSteps) tryCreateEmptyTask(ctx context.Context) error {
	tm := tasks.NewTaskManager(filepath.Join(s.suite.ProjectDir, ".agentic/tasks"))
	task, err := tm.CreateTask("")
	s.suite.LastCommandErr = err
	if task != nil {
		s.suite.LastTaskID = task.ID
	}
	// Don't return error - we want to test failures
	return nil
}
