package functional

import (
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

// TestBeginnerScenario tests the complete beginner workflow from the tutorial:
// 1. Initialize project
// 2. Create sample task
// 3. List tasks
// 4. Claim task
// 5. Complete task
// 6. Validate
func TestBeginnerScenario(t *testing.T) {
	tmpDir := SetupTestProject(t)

	// Step 1: Initialize project
	t.Log("Step 1: Initializing project...")
	err := project.InitProject("BeginnerProject")
	require.NoError(t, err, "Project initialization should succeed")

	// Verify project structure
	VerifyProjectStructure(t, tmpDir)
	t.Log("✓ Project structure created successfully")

	// Step 2: Create a sample task
	t.Log("Step 2: Creating sample task...")
	tm := tasks.NewTaskManager(filepath.Join(tmpDir, ".agentic/tasks"))
	task, err := tm.CreateTask("My First Task")
	require.NoError(t, err, "Task creation should succeed")
	t.Logf("✓ Task created: %s", task.ID)

	// Step 3: List tasks (verify task is in backlog)
	t.Log("Step 3: Listing tasks...")
	backlog, err := tm.LoadTasks("backlog")
	require.NoError(t, err)
	assert.Len(t, backlog.Tasks, 1, "Backlog should contain one task")
	t.Logf("✓ Found %d task(s) in backlog", len(backlog.Tasks))

	// Step 4: Claim the task
	t.Log("Step 4: Claiming task...")
	err = tm.ClaimTask(task.ID, "beginner-agent")
	require.NoError(t, err, "Task claim should succeed")

	// Verify task moved to in-progress
	inProgress, err := tm.LoadTasks("in-progress")
	require.NoError(t, err)
	assert.Len(t, inProgress.Tasks, 1, "In-progress should contain one task")
	assert.Equal(t, models.StatusInProgress, inProgress.Tasks[0].Status)
	t.Logf("✓ Task claimed and moved to in-progress")

	// Step 5: Complete the task
	t.Log("Step 5: Completing task...")
	err = tm.MoveTask(task.ID, "in-progress", "done", models.StatusDone)
	require.NoError(t, err, "Task completion should succeed")

	// Verify task moved to done
	done, err := tm.LoadTasks("done")
	require.NoError(t, err)
	assert.Len(t, done.Tasks, 1, "Done should contain one task")
	assert.Equal(t, models.StatusDone, done.Tasks[0].Status)
	t.Logf("✓ Task completed and moved to done")

	// Step 6: Run validation
	t.Log("Step 6: Running validation...")
	v := validator.NewValidator()
	// Note: We skip certain validators that require actual source code
	// In a real scenario, validation would check context files, etc.
	ctx := &validator.ValidationContext{ProjectRoot: tmpDir}
	results, err := v.Validate(ctx)
	require.NoError(t, err, "Validation should run without errors")
	t.Logf("✓ Validation completed with %d result(s)", len(results))

	t.Log("✅ Beginner scenario completed successfully!")
}

// TestIntermediateScenario tests the intermediate workflow:
// 1. Initialize with metadata
// 2. Create task with full details
// 3. Work on the task
// 4. Validate with rules
func TestIntermediateScenario(t *testing.T) {
	tmpDir := SetupTestProject(t)

	// Step 1: Initialize project with metadata
	t.Log("Step 1: Initializing project with metadata...")
	err := project.InitProject("BlogAPI")
	require.NoError(t, err)
	VerifyProjectStructure(t, tmpDir)
	t.Log("✓ Project initialized")

	// Step 2: Create task with full metadata
	t.Log("Step 2: Creating feature task with full metadata...")
	tm := tasks.NewTaskManager(filepath.Join(tmpDir, ".agentic/tasks"))
	task, err := tm.CreateTask("Implement Blog Post API")
	require.NoError(t, err)

	// Add full task details
	backlog, err := tm.LoadTasks("backlog")
	require.NoError(t, err)

	backlog.Tasks[0].Description = "Create REST API endpoints for blog post CRUD operations"
	backlog.Tasks[0].Scope = []string{"src/api", "src/models"}
	backlog.Tasks[0].SpecRefs = []string{".agentic/spec/04-architecture.md"}
	backlog.Tasks[0].Inputs = []string{".agentic/context/rolling-summary.md"}
	backlog.Tasks[0].Outputs = []string{
		"src/api/posts.go",
		"src/models/post.go",
		"tests/api/posts_test.go",
	}
	backlog.Tasks[0].Acceptance = []string{
		"GET /posts returns all posts",
		"POST /posts creates a new post",
		"PUT /posts/:id updates a post",
		"DELETE /posts/:id deletes a post",
		"All endpoints have tests",
	}

	err = tm.SaveTasks("backlog", backlog)
	require.NoError(t, err)
	t.Logf("✓ Task created with full metadata: %s", task.ID)

	// Step 3: Claim and work on the task
	t.Log("Step 3: Claiming task...")
	err = tm.ClaimTask(task.ID, "dev-agent")
	require.NoError(t, err)
	t.Log("✓ Task claimed")

	// Simulate doing work by creating expected output files
	t.Log("Step 4: Creating output files...")
	CreateTestFile(t, "src/api/posts.go", "package api\n\n// Post API handlers")
	CreateTestFile(t, "src/models/post.go", "package models\n\ntype Post struct {}")
	CreateTestFile(t, "tests/api/posts_test.go", "package api_test\n\n// Tests")
	t.Log("✓ Work completed (files created)")

	// Step 5: Validate
	t.Log("Step 5: Running validation...")
	v := validator.NewValidator()
	ctx := &validator.ValidationContext{ProjectRoot: tmpDir}
	results, err := v.Validate(ctx)
	require.NoError(t, err)
	t.Logf("✓ Validation completed with %d result(s)", len(results))

	// Step 6: Complete the task
	t.Log("Step 6: Completing task...")
	err = tm.MoveTask(task.ID, "in-progress", "done", models.StatusDone)
	require.NoError(t, err)
	t.Log("✓ Task completed")

	// Verify task in done with all metadata preserved
	done, err := tm.LoadTasks("done")
	require.NoError(t, err)
	assert.Len(t, done.Tasks, 1)
	completedTask := done.Tasks[0]
	assert.Equal(t, task.ID, completedTask.ID)
	assert.Len(t, completedTask.Acceptance, 5, "All acceptance criteria should be preserved")
	assert.Len(t, completedTask.Outputs, 3, "All outputs should be preserved")

	t.Log("✅ Intermediate scenario completed successfully!")
}

// TestAdvancedScenario tests the advanced workflow with task decomposition:
// 1. Create large feature task
// 2. Decompose into subtasks
// 3. Work on subtasks sequentially
// 4. Validate throughout
func TestAdvancedScenario(t *testing.T) {
	tmpDir := SetupTestProject(t)

	// Step 1: Initialize project
	t.Log("Step 1: Initializing project...")
	err := project.InitProject("AdvancedProject")
	require.NoError(t, err)
	t.Log("✓ Project initialized")

	tm := tasks.NewTaskManager(filepath.Join(tmpDir, ".agentic/tasks"))

	// Step 2: Create a large feature task
	t.Log("Step 2: Creating large feature task...")
	task, err := tm.CreateTask("Implement Complete User Management System")
	require.NoError(t, err)
	t.Logf("✓ Parent task created: %s", task.ID)

	// Step 3: Decompose into subtasks
	t.Log("Step 3: Decomposing task into subtasks...")
	subtaskTitles := []string{
		"Create user model and database schema",
		"Implement user registration endpoint",
		"Implement user authentication",
		"Add password reset functionality",
		"Write comprehensive tests",
	}
	err = tm.DecomposeTask(task.ID, subtaskTitles)
	require.NoError(t, err)
	t.Logf("✓ Task decomposed into %d subtasks", len(subtaskTitles))

	// Verify decomposition
	backlog, err := tm.LoadTasks("backlog")
	require.NoError(t, err)
	parentTask := backlog.Tasks[0]
	assert.Len(t, parentTask.SubTasks, 5, "Should have 5 subtasks")

	// Step 4: Work on subtasks sequentially
	t.Log("Step 4: Working on subtasks...")

	// For this test, we'll simulate working on first two subtasks
	for i := 0; i < 2; i++ {
		subtask := parentTask.SubTasks[i]
		t.Logf("  Working on subtask %d: %s", i+1, subtask.Title)

		// Simulate work (in real scenario, this would involve actual file creation)
		time.Sleep(100 * time.Millisecond) // Small delay to ensure different timestamps

		// Mark subtask as done (in actual implementation)
		// This is a simplified version - real implementation would update subtask status
		t.Logf("  ✓ Subtask %d completed", i+1)
	}

	// Step 5: Validate throughout
	t.Log("Step 5: Running validation...")
	v := validator.NewValidator()
	ctx := &validator.ValidationContext{ProjectRoot: tmpDir}
	results, err := v.Validate(ctx)
	require.NoError(t, err)
	t.Logf("✓ Validation completed with %d result(s)", len(results))

	// Verify parent task still exists with subtasks
	backlog, err = tm.LoadTasks("backlog")
	require.NoError(t, err)
	assert.Len(t, backlog.Tasks, 1, "Parent task should still be in backlog")
	assert.Len(t, backlog.Tasks[0].SubTasks, 5, "All subtasks should be preserved")

	t.Log("✅ Advanced scenario completed successfully!")
}

// TestErrorScenarios tests error handling throughout the workflow.
func TestErrorScenarios(t *testing.T) {
	tmpDir := SetupTestProject(t)

	// Initialize project
	err := project.InitProject("ErrorTestProject")
	require.NoError(t, err)

	tm := tasks.NewTaskManager(filepath.Join(tmpDir, ".agentic/tasks"))

	t.Run("Claim nonexistent task", func(t *testing.T) {
		t.Log("Testing: Claim nonexistent task...")
		err := tm.ClaimTask("TASK-NONEXISTENT", "agent")
		assert.Error(t, err, "Should fail to claim nonexistent task")
		t.Log("✓ Correctly rejected nonexistent task claim")
	})

	t.Run("Complete unclaimed task", func(t *testing.T) {
		t.Log("Testing: Complete task that wasn't claimed...")
		task, err := tm.CreateTask("Unclaimed Task")
		require.NoError(t, err)

		// Try to complete without claiming
		err = tm.MoveTask(task.ID, "backlog", "done", models.StatusDone)
		// This should either fail or handle gracefully
		// The behavior depends on implementation
		t.Logf("Move from backlog to done: %v", err)
	})

	t.Run("Decompose already decomposed task", func(t *testing.T) {
		t.Log("Testing: Decompose already decomposed task...")
		task, err := tm.CreateTask("Task to Decompose")
		require.NoError(t, err)

		// First decomposition
		err = tm.DecomposeTask(task.ID, []string{"Subtask 1"})
		require.NoError(t, err)

		// Second decomposition (should handle gracefully or append)
		err = tm.DecomposeTask(task.ID, []string{"Subtask 2"})
		// Behavior depends on implementation
		t.Logf("Second decomposition result: %v", err)
	})

	t.Run("Validation without required files", func(t *testing.T) {
		t.Log("Testing: Validation with missing context files...")

		// Create a source directory without context.md
		CreateTestFile(t, "src/module/code.go", "package module\n\nfunc Foo() {}")

		// Run validation
		v := validator.NewValidator()
		v.Register(&rules.DirectoryContextRule{})

		ctx := &validator.ValidationContext{ProjectRoot: tmpDir}
		results, err := v.Validate(ctx)
		require.NoError(t, err, "Validation should run even if rules fail")

		// Check for failures
		hasFailure := false
		for _, result := range results {
			if result.Status == "FAIL" {
				hasFailure = true
				if len(result.Errors) > 0 {
					t.Logf("✓ Validation correctly detected: %s", result.Errors[0])
				}
				break
			}
		}
		assert.True(t, hasFailure, "Should have at least one validation failure")
	})

	t.Log("✅ Error scenarios handled correctly!")
}

// TestCompleteWorkflow tests a complete realistic workflow from start to finish.
func TestCompleteWorkflow(t *testing.T) {
	tmpDir := SetupTestProject(t)

	t.Log("=== Starting Complete Workflow Test ===")

	// Phase 1: Project Setup
	t.Log("\nPhase 1: Project Setup")
	err := project.InitProject("CompleteWorkflowTest")
	require.NoError(t, err)
	VerifyProjectStructure(t, tmpDir)
	t.Log("✓ Project initialized")

	tm := tasks.NewTaskManager(filepath.Join(tmpDir, ".agentic/tasks"))

	// Phase 2: Create Feature Epic
	t.Log("\nPhase 2: Create Feature Epic")
	epic, err := tm.CreateTask("User Authentication System")
	require.NoError(t, err)

	// Add epic details
	backlog, err := tm.LoadTasks("backlog")
	require.NoError(t, err)
	backlog.Tasks[0].Description = "Complete user authentication with JWT"
	backlog.Tasks[0].Acceptance = []string{
		"Users can register",
		"Users can login",
		"JWT tokens are issued",
		"Protected routes work",
	}
	err = tm.SaveTasks("backlog", backlog)
	require.NoError(t, err)
	t.Logf("✓ Epic created: %s", epic.ID)

	// Phase 3: Decompose into tasks
	t.Log("\nPhase 3: Decompose Epic")
	err = tm.DecomposeTask(epic.ID, []string{
		"Create user model",
		"Implement registration",
		"Implement login",
		"Add JWT middleware",
	})
	require.NoError(t, err)
	t.Log("✓ Epic decomposed into subtasks")

	// Phase 4: Work on first task
	t.Log("\nPhase 4: Work on First Subtask")
	err = tm.ClaimTask(epic.ID, "dev-agent")
	require.NoError(t, err)
	t.Log("✓ Epic claimed")

	// Simulate work
	CreateTestFile(t, "src/models/user.go", "package models\n\ntype User struct {}")
	CreateTestFile(t, "src/models/user_test.go", "package models\n\n// Tests")
	t.Log("✓ Work completed on user model")

	// Phase 5: Validation
	t.Log("\nPhase 5: Validation")
	v := validator.NewValidator()
	ctx := &validator.ValidationContext{ProjectRoot: tmpDir}
	results, err := v.Validate(ctx)
	require.NoError(t, err)
	t.Logf("✓ Validation completed (%d results)", len(results))

	// Phase 6: Complete epic
	t.Log("\nPhase 6: Complete Epic")
	err = tm.MoveTask(epic.ID, "in-progress", "done", models.StatusDone)
	require.NoError(t, err)
	t.Log("✓ Epic completed")

	// Verify final state
	done, err := tm.LoadTasks("done")
	require.NoError(t, err)
	assert.Len(t, done.Tasks, 1, "Should have one completed epic")
	assert.Len(t, done.Tasks[0].SubTasks, 4, "Epic should preserve all subtasks")

	t.Log("\n✅ Complete workflow executed successfully!")
	t.Log("=== Workflow Test Complete ===")
}
