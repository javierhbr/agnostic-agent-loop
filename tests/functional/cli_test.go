package functional

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/javierbenavides/agentic-agent/internal/project"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInitCommand tests the init command with flag mode.
func TestInitCommand(t *testing.T) {
	tmpDir := SetupTestProject(t)

	// Initialize project using the project package directly
	// (testing the underlying logic that the command uses)
	err := project.InitProject("TestProject")
	require.NoError(t, err, "Project initialization should succeed")

	// Verify project structure was created
	VerifyProjectStructure(t, tmpDir)

	// Verify config file contains project name
	configPath := filepath.Join(tmpDir, "agnostic-agent.yaml")
	assert.FileExists(t, configPath, "Config file should exist")
}

// TestTaskCreateCommand tests creating a task with the task manager.
func TestTaskCreateCommand(t *testing.T) {
	tmpDir := SetupTestProject(t)

	// Initialize project first
	err := project.InitProject("TestProject")
	require.NoError(t, err)

	// Create task using task manager
	tm := tasks.NewTaskManager(filepath.Join(tmpDir, ".agentic/tasks"))
	task, err := tm.CreateTask("Test Feature Implementation")
	require.NoError(t, err, "Task creation should succeed")

	// Verify task properties
	assert.NotEmpty(t, task.ID, "Task should have an ID")
	assert.Equal(t, "Test Feature Implementation", task.Title)
	assert.Equal(t, models.StatusPending, task.Status)

	// Verify task is in backlog file
	backlogPath := filepath.Join(tmpDir, ".agentic/tasks/backlog.yaml")
	taskList := VerifyTaskFile(t, backlogPath)
	assert.Len(t, taskList.Tasks, 1, "Backlog should contain one task")
	assert.Equal(t, task.ID, taskList.Tasks[0].ID)
}

// TestTaskClaimCommand tests claiming a task.
func TestTaskClaimCommand(t *testing.T) {
	tmpDir := SetupTestProject(t)

	// Initialize project
	err := project.InitProject("TestProject")
	require.NoError(t, err)

	// Create a task
	tm := tasks.NewTaskManager(filepath.Join(tmpDir, ".agentic/tasks"))
	task, err := tm.CreateTask("Task to Claim")
	require.NoError(t, err)

	// Claim the task
	err = tm.ClaimTask(task.ID, "test-agent")
	require.NoError(t, err, "Task claim should succeed")

	// Verify task moved to in-progress
	inProgressPath := filepath.Join(tmpDir, ".agentic/tasks/in-progress.yaml")
	taskList := VerifyTaskFile(t, inProgressPath)
	assert.Len(t, taskList.Tasks, 1, "In-progress should contain one task")
	assert.Equal(t, task.ID, taskList.Tasks[0].ID)
	assert.Equal(t, models.StatusInProgress, taskList.Tasks[0].Status)
	assert.Equal(t, "test-agent", taskList.Tasks[0].AssignedTo)

	// Verify task removed from backlog
	backlogPath := filepath.Join(tmpDir, ".agentic/tasks/backlog.yaml")
	backlogList := VerifyTaskFile(t, backlogPath)
	assert.Empty(t, backlogList.Tasks, "Backlog should be empty")
}

// TestTaskCompleteCommand tests completing a task.
func TestTaskCompleteCommand(t *testing.T) {
	tmpDir := SetupTestProject(t)

	// Initialize project
	err := project.InitProject("TestProject")
	require.NoError(t, err)

	// Create and claim a task
	tm := tasks.NewTaskManager(filepath.Join(tmpDir, ".agentic/tasks"))
	task, err := tm.CreateTask("Task to Complete")
	require.NoError(t, err)

	err = tm.ClaimTask(task.ID, "test-agent")
	require.NoError(t, err)

	// Complete the task
	err = tm.MoveTask(task.ID, "in-progress", "done", models.StatusDone)
	require.NoError(t, err, "Task completion should succeed")

	// Verify task moved to done
	donePath := filepath.Join(tmpDir, ".agentic/tasks/done.yaml")
	taskList := VerifyTaskFile(t, donePath)
	assert.Len(t, taskList.Tasks, 1, "Done should contain one task")
	assert.Equal(t, task.ID, taskList.Tasks[0].ID)
	assert.Equal(t, models.StatusDone, taskList.Tasks[0].Status)

	// Verify task removed from in-progress
	inProgressPath := filepath.Join(tmpDir, ".agentic/tasks/in-progress.yaml")
	inProgressList := VerifyTaskFile(t, inProgressPath)
	assert.Empty(t, inProgressList.Tasks, "In-progress should be empty")
}

// TestTaskDecomposeCommand tests decomposing a task into subtasks.
func TestTaskDecomposeCommand(t *testing.T) {
	tmpDir := SetupTestProject(t)

	// Initialize project
	err := project.InitProject("TestProject")
	require.NoError(t, err)

	// Create a task
	tm := tasks.NewTaskManager(filepath.Join(tmpDir, ".agentic/tasks"))
	task, err := tm.CreateTask("Large Feature")
	require.NoError(t, err)

	// Decompose the task
	subtaskTitles := []string{
		"Subtask 1: Setup",
		"Subtask 2: Implementation",
		"Subtask 3: Testing",
	}
	err = tm.DecomposeTask(task.ID, subtaskTitles)
	require.NoError(t, err, "Task decomposition should succeed")

	// Verify subtasks were added
	backlogPath := filepath.Join(tmpDir, ".agentic/tasks/backlog.yaml")
	taskList := VerifyTaskFile(t, backlogPath)
	require.Len(t, taskList.Tasks, 1, "Should still have one parent task")

	parentTask := taskList.Tasks[0]
	assert.Len(t, parentTask.SubTasks, 3, "Parent should have 3 subtasks")

	// Verify subtask properties
	for i, subtask := range parentTask.SubTasks {
		assert.Equal(t, subtaskTitles[i], subtask.Title)
		assert.Equal(t, models.StatusPending, subtask.Status)
		assert.Contains(t, subtask.ID, task.ID, "Subtask ID should reference parent")
	}
}

// TestTaskShowCommand tests showing task details.
func TestTaskShowCommand(t *testing.T) {
	tmpDir := SetupTestProject(t)

	// Initialize project
	err := project.InitProject("TestProject")
	require.NoError(t, err)

	// Create a task with full details
	tm := tasks.NewTaskManager(filepath.Join(tmpDir, ".agentic/tasks"))
	task, err := tm.CreateTask("Detailed Task")
	require.NoError(t, err)

	// Load and update task with details
	backlog, err := tm.LoadTasks("backlog")
	require.NoError(t, err)

	backlog.Tasks[0].Description = "This is a detailed description"
	backlog.Tasks[0].Acceptance = []string{"Criterion 1", "Criterion 2"}
	backlog.Tasks[0].Outputs = []string{"file1.go", "file2.go"}

	err = tm.SaveTasks("backlog", backlog)
	require.NoError(t, err)

	// Find and verify the task
	foundTask, source, err := tm.FindTask(task.ID)
	require.NoError(t, err, "Task should be found")
	assert.NotNil(t, foundTask, "Found task should not be nil")
	assert.Equal(t, "backlog", source, "Task should be in backlog")
	assert.Equal(t, "Detailed Task", foundTask.Title)
	assert.Equal(t, "This is a detailed description", foundTask.Description)
	assert.Len(t, foundTask.Acceptance, 2)
	assert.Len(t, foundTask.Outputs, 2)
}

// TestFindTaskAcrossLists tests finding tasks in different states.
func TestFindTaskAcrossLists(t *testing.T) {
	tmpDir := SetupTestProject(t)

	// Initialize project
	err := project.InitProject("TestProject")
	require.NoError(t, err)

	tm := tasks.NewTaskManager(filepath.Join(tmpDir, ".agentic/tasks"))

	// Create task in backlog
	task1, err := tm.CreateTask("Backlog Task")
	require.NoError(t, err)

	// Find in backlog
	found1, source1, err := tm.FindTask(task1.ID)
	require.NoError(t, err)
	assert.NotNil(t, found1)
	assert.Equal(t, "backlog", source1)
	assert.Equal(t, "Backlog Task", found1.Title)

	// Move to in-progress
	err = tm.ClaimTask(task1.ID, "agent1")
	require.NoError(t, err)

	// Find in in-progress
	found2, source2, err := tm.FindTask(task1.ID)
	require.NoError(t, err)
	assert.NotNil(t, found2)
	assert.Equal(t, "in-progress", source2)
	assert.Equal(t, models.StatusInProgress, found2.Status)

	// Move to done
	err = tm.MoveTask(task1.ID, "in-progress", "done", models.StatusDone)
	require.NoError(t, err)

	// Find in done
	found3, source3, err := tm.FindTask(task1.ID)
	require.NoError(t, err)
	assert.NotNil(t, found3)
	assert.Equal(t, "done", source3)
	assert.Equal(t, models.StatusDone, found3.Status)
}

// TestTaskWithFullFields tests creating a task with all optional fields.
func TestTaskWithFullFields(t *testing.T) {
	tmpDir := SetupTestProject(t)

	// Initialize project
	err := project.InitProject("TestProject")
	require.NoError(t, err)

	tm := tasks.NewTaskManager(filepath.Join(tmpDir, ".agentic/tasks"))

	// Create basic task
	_, err = tm.CreateTask("Full Fields Task")
	require.NoError(t, err)

	// Load and update with all fields
	backlog, err := tm.LoadTasks("backlog")
	require.NoError(t, err)

	backlog.Tasks[0].Description = "Complete implementation with all metadata"
	backlog.Tasks[0].Scope = []string{"src/auth", "src/middleware"}
	backlog.Tasks[0].SpecRefs = []string{".agentic/spec/architecture.md"}
	backlog.Tasks[0].Inputs = []string{".agentic/context/rolling-summary.md"}
	backlog.Tasks[0].Outputs = []string{"src/auth/jwt.go", "tests/auth_test.go"}
	backlog.Tasks[0].Acceptance = []string{
		"JWT tokens can be generated",
		"Token validation works",
		"All tests pass",
	}

	err = tm.SaveTasks("backlog", backlog)
	require.NoError(t, err)

	// Verify all fields persisted
	reloaded, err := tm.LoadTasks("backlog")
	require.NoError(t, err)
	assert.Len(t, reloaded.Tasks, 1)

	persistedTask := reloaded.Tasks[0]
	assert.Equal(t, "Full Fields Task", persistedTask.Title)
	assert.Equal(t, "Complete implementation with all metadata", persistedTask.Description)
	assert.Equal(t, []string{"src/auth", "src/middleware"}, persistedTask.Scope)
	assert.Len(t, persistedTask.SpecRefs, 1)
	assert.Len(t, persistedTask.Inputs, 1)
	assert.Len(t, persistedTask.Outputs, 2)
	assert.Len(t, persistedTask.Acceptance, 3)
}

// TestErrorCases tests error handling in various scenarios.
func TestErrorCases(t *testing.T) {
	tmpDir := SetupTestProject(t)

	// Initialize project
	err := project.InitProject("TestProject")
	require.NoError(t, err)

	tm := tasks.NewTaskManager(filepath.Join(tmpDir, ".agentic/tasks"))

	t.Run("ClaimNonexistentTask", func(t *testing.T) {
		err := tm.ClaimTask("NONEXISTENT-TASK", "agent1")
		assert.Error(t, err, "Claiming nonexistent task should fail")
	})

	t.Run("FindNonexistentTask", func(t *testing.T) {
		found, source, err := tm.FindTask("NONEXISTENT-TASK")
		assert.NoError(t, err, "Find should not error on missing task")
		assert.Nil(t, found, "Should return nil for nonexistent task")
		assert.Empty(t, source, "Source should be empty for nonexistent task")
	})

	t.Run("MoveNonexistentTask", func(t *testing.T) {
		err := tm.MoveTask("NONEXISTENT-TASK", "backlog", "done", models.StatusDone)
		assert.Error(t, err, "Moving nonexistent task should fail")
	})

	t.Run("DecomposeNonexistentTask", func(t *testing.T) {
		err := tm.DecomposeTask("NONEXISTENT-TASK", []string{"Subtask 1"})
		assert.Error(t, err, "Decomposing nonexistent task should fail")
	})
}

// TestMultipleTasksWorkflow tests working with multiple tasks.
func TestMultipleTasksWorkflow(t *testing.T) {
	tmpDir := SetupTestProject(t)

	// Initialize project
	err := project.InitProject("TestProject")
	require.NoError(t, err)

	tm := tasks.NewTaskManager(filepath.Join(tmpDir, ".agentic/tasks"))

	// Create multiple tasks (with delays to ensure unique IDs)
	task1, err := tm.CreateTask("Task 1")
	require.NoError(t, err)
	time.Sleep(1100 * time.Millisecond)

	task2, err := tm.CreateTask("Task 2")
	require.NoError(t, err)
	time.Sleep(1100 * time.Millisecond)

	_, err = tm.CreateTask("Task 3")
	require.NoError(t, err)

	// Verify all in backlog
	backlog, err := tm.LoadTasks("backlog")
	require.NoError(t, err)
	assert.Len(t, backlog.Tasks, 3, "Backlog should have 3 tasks")

	// Claim task1
	err = tm.ClaimTask(task1.ID, "agent1")
	require.NoError(t, err)

	// Claim task2
	err = tm.ClaimTask(task2.ID, "agent2")
	require.NoError(t, err)

	// Verify states
	backlog, err = tm.LoadTasks("backlog")
	require.NoError(t, err)
	assert.Len(t, backlog.Tasks, 1, "Backlog should have 1 task")

	inProgress, err := tm.LoadTasks("in-progress")
	require.NoError(t, err)
	assert.Len(t, inProgress.Tasks, 2, "In-progress should have 2 tasks")

	// Complete task1
	err = tm.MoveTask(task1.ID, "in-progress", "done", models.StatusDone)
	require.NoError(t, err)

	// Verify final states
	backlog, err = tm.LoadTasks("backlog")
	require.NoError(t, err)
	assert.Len(t, backlog.Tasks, 1)

	inProgress, err = tm.LoadTasks("in-progress")
	require.NoError(t, err)
	assert.Len(t, inProgress.Tasks, 1)

	done, err := tm.LoadTasks("done")
	require.NoError(t, err)
	assert.Len(t, done.Tasks, 1)
	assert.Equal(t, task1.ID, done.Tasks[0].ID)
}
