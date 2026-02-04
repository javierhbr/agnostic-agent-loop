package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/javierbenavides/agentic-agent/internal/project"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/internal/validator"
	"github.com/javierbenavides/agentic-agent/internal/validator/rules"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupIntegrationTest(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "integration-test-*")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	// Change to test directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(tmpDir)
	require.NoError(t, err)
	t.Cleanup(func() {
		os.Chdir(originalDir)
	})

	return tmpDir
}

// TestHappyPath tests the complete workflow:
// init → create task → claim task → complete task
func TestHappyPath(t *testing.T) {
	tmpDir := setupIntegrationTest(t)

	// Step 1: Initialize project
	err := project.InitProject("IntegrationTest")
	require.NoError(t, err)

	// Verify structure created
	assert.DirExists(t, filepath.Join(tmpDir, ".agentic"))
	assert.DirExists(t, filepath.Join(tmpDir, ".agentic", "tasks"))
	assert.DirExists(t, filepath.Join(tmpDir, ".agentic", "context"))
	assert.DirExists(t, filepath.Join(tmpDir, ".agentic", "spec"))
	assert.DirExists(t, filepath.Join(tmpDir, ".agentic", "agent-rules"))
	assert.FileExists(t, filepath.Join(tmpDir, "agnostic-agent.yaml"))

	// Step 2: Create a task
	tm := tasks.NewTaskManager(".agentic/tasks")
	task, err := tm.CreateTask("Test Feature Implementation")
	require.NoError(t, err)
	assert.NotEmpty(t, task.ID)
	assert.Equal(t, "Test Feature Implementation", task.Title)
	assert.Equal(t, models.StatusPending, task.Status)

	// Verify task in backlog
	backlog, err := tm.LoadTasks("backlog")
	require.NoError(t, err)
	assert.Len(t, backlog.Tasks, 1)
	assert.Equal(t, task.ID, backlog.Tasks[0].ID)

	// Step 3: Claim the task
	err = tm.ClaimTask(task.ID, "test-agent")
	require.NoError(t, err)

	// Verify task moved to in-progress
	inProgress, err := tm.LoadTasks("in-progress")
	require.NoError(t, err)
	assert.Len(t, inProgress.Tasks, 1)
	assert.Equal(t, task.ID, inProgress.Tasks[0].ID)
	assert.Equal(t, models.StatusInProgress, inProgress.Tasks[0].Status)
	assert.Equal(t, "test-agent", inProgress.Tasks[0].AssignedTo)

	// Verify removed from backlog
	backlog, err = tm.LoadTasks("backlog")
	require.NoError(t, err)
	assert.Empty(t, backlog.Tasks)

	// Step 4: Complete the task
	err = tm.MoveTask(task.ID, "in-progress", "done", models.StatusDone)
	require.NoError(t, err)

	// Verify task moved to done
	done, err := tm.LoadTasks("done")
	require.NoError(t, err)
	assert.Len(t, done.Tasks, 1)
	assert.Equal(t, task.ID, done.Tasks[0].ID)
	assert.Equal(t, models.StatusDone, done.Tasks[0].Status)
	assert.Equal(t, "test-agent", done.Tasks[0].AssignedTo) // Should preserve

	// Verify removed from in-progress
	inProgress, err = tm.LoadTasks("in-progress")
	require.NoError(t, err)
	assert.Empty(t, inProgress.Tasks)
}

// TestTaskWithFullFields tests creating a task with all fields populated
func TestTaskWithFullFields(t *testing.T) {
	_ = setupIntegrationTest(t)

	err := project.InitProject("FullFieldsTest")
	require.NoError(t, err)

	tm := tasks.NewTaskManager(".agentic/tasks")

	// Create basic task
	_, err = tm.CreateTask("Implement Authentication")
	require.NoError(t, err)

	// Load and update with all fields
	backlog, err := tm.LoadTasks("backlog")
	require.NoError(t, err)

	backlog.Tasks[0].Description = "Implement JWT-based authentication"
	backlog.Tasks[0].Scope = []string{"src/auth", "src/middleware"}
	backlog.Tasks[0].SpecRefs = []string{".agentic/spec/04-architecture.md", ".agentic/spec/05-domain-model.md"}
	backlog.Tasks[0].Inputs = []string{".agentic/context/rolling-summary.md"}
	backlog.Tasks[0].Outputs = []string{"src/auth/jwt.go", "src/auth/middleware.go", "tests/auth_test.go"}
	backlog.Tasks[0].Acceptance = []string{
		"JWT tokens can be generated",
		"Token validation works",
		"Middleware protects routes",
		"All tests pass",
	}

	err = tm.SaveTasks("backlog", backlog)
	require.NoError(t, err)

	// Verify all fields persisted
	reloaded, err := tm.LoadTasks("backlog")
	require.NoError(t, err)
	assert.Len(t, reloaded.Tasks, 1)

	persistedTask := reloaded.Tasks[0]
	assert.Equal(t, "Implement Authentication", persistedTask.Title)
	assert.Equal(t, "Implement JWT-based authentication", persistedTask.Description)
	assert.Equal(t, []string{"src/auth", "src/middleware"}, persistedTask.Scope)
	assert.Len(t, persistedTask.SpecRefs, 2)
	assert.Len(t, persistedTask.Inputs, 1)
	assert.Len(t, persistedTask.Outputs, 3)
	assert.Len(t, persistedTask.Acceptance, 4)
}

// TestValidationWorkflow tests the validation rules
func TestValidationWorkflow(t *testing.T) {
	tmpDir := setupIntegrationTest(t)

	err := project.InitProject("ValidationTest")
	require.NoError(t, err)

	// Create a source directory without context.md
	srcDir := filepath.Join(tmpDir, "src", "module")
	err = os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(srcDir, "code.go"), []byte("package module\n\nfunc Foo() {}"), 0644)
	require.NoError(t, err)

	// Run validation - should fail
	v := validator.NewValidator()
	v.Register(&rules.DirectoryContextRule{})

	ctx := &validator.ValidationContext{ProjectRoot: tmpDir}
	results, err := v.Validate(ctx)
	require.NoError(t, err)

	// Should have at least one failure
	hasFailure := false
	for _, result := range results {
		if result.Status == "FAIL" {
			hasFailure = true
			break
		}
	}
	assert.True(t, hasFailure, "Expected validation to fail for missing context.md")

	// Add context.md
	contextContent := `# Module Context

## Purpose
Core business logic module

## Responsibilities
- Implements domain models
- Handles business rules

## Dependencies
None

## Must Do
- Keep logic pure
- No external dependencies

## Cannot Do
- Direct database access
- HTTP calls
`
	err = os.WriteFile(filepath.Join(srcDir, "context.md"), []byte(contextContent), 0644)
	require.NoError(t, err)

	// Run validation again - should pass
	results, err = v.Validate(ctx)
	require.NoError(t, err)

	// Check if directory context rule passed
	for _, result := range results {
		if result.RuleName == "directory-context" {
			assert.Equal(t, "PASS", result.Status)
		}
	}
}

// TestTaskDecomposition tests breaking tasks into subtasks
func TestTaskDecomposition(t *testing.T) {
	_ = setupIntegrationTest(t)

	err := project.InitProject("DecompositionTest")
	require.NoError(t, err)

	tm := tasks.NewTaskManager(".agentic/tasks")

	// Create parent task
	task, err := tm.CreateTask("Implement User Management")
	require.NoError(t, err)

	// Decompose into subtasks
	subtaskTitles := []string{
		"Create user model",
		"Implement registration",
		"Implement login",
		"Add password reset",
	}

	err = tm.DecomposeTask(task.ID, subtaskTitles)
	require.NoError(t, err)

	// Verify subtasks added
	backlog, err := tm.LoadTasks("backlog")
	require.NoError(t, err)
	assert.Len(t, backlog.Tasks, 1)

	parentTask := backlog.Tasks[0]
	assert.Len(t, parentTask.SubTasks, 4)

	for i, subtask := range parentTask.SubTasks {
		assert.Equal(t, subtaskTitles[i], subtask.Title)
		assert.Equal(t, models.StatusPending, subtask.Status)
		assert.Contains(t, subtask.ID, task.ID) // Subtask ID should reference parent
	}
}

// TestFindTaskAcrossLists tests finding tasks in different states
func TestFindTaskAcrossLists(t *testing.T) {
	_ = setupIntegrationTest(t)

	err := project.InitProject("FindTest")
	require.NoError(t, err)

	tm := tasks.NewTaskManager(".agentic/tasks")

	// Create tasks in different states
	// Note: CreateTask uses time.Now().Unix() for IDs (seconds), so we need delays
	task1, err := tm.CreateTask("Pending Task")
	require.NoError(t, err)

	// Delay to ensure unique Unix timestamp-based IDs
	time.Sleep(1100 * time.Millisecond)

	task2, err := tm.CreateTask("In Progress Task")
	require.NoError(t, err)
	err = tm.ClaimTask(task2.ID, "agent1")
	require.NoError(t, err)

	time.Sleep(1100 * time.Millisecond)

	task3, err := tm.CreateTask("Done Task")
	require.NoError(t, err)
	err = tm.ClaimTask(task3.ID, "agent2")
	require.NoError(t, err)
	err = tm.MoveTask(task3.ID, "in-progress", "done", models.StatusDone)
	require.NoError(t, err)

	// Find each task
	found1, source1, err := tm.FindTask(task1.ID)
	require.NoError(t, err)
	assert.NotNil(t, found1)
	assert.Equal(t, "backlog", source1)
	assert.Equal(t, "Pending Task", found1.Title)

	found2, source2, err := tm.FindTask(task2.ID)
	require.NoError(t, err)
	assert.NotNil(t, found2)
	assert.Equal(t, "in-progress", source2)
	assert.Equal(t, "In Progress Task", found2.Title)

	found3, source3, err := tm.FindTask(task3.ID)
	require.NoError(t, err)
	assert.NotNil(t, found3)
	assert.Equal(t, "done", source3)
	assert.Equal(t, "Done Task", found3.Title)

	// Try to find non-existent task
	notFound, _, err := tm.FindTask("NONEXISTENT")
	require.NoError(t, err)
	assert.Nil(t, notFound)
}
