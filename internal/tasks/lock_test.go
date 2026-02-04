package tasks

import (
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClaimTask_Success(t *testing.T) {
	tmpDir := setupTestDir(t)

	tm := NewTaskManager(tmpDir)

	// Setup backlog with task
	backlog := &TaskList{
		Tasks: []models.Task{
			{ID: "TASK-001", Title: "Test Task", Status: models.StatusPending},
		},
	}
	err := tm.SaveTasks("backlog", backlog)
	require.NoError(t, err)

	// Create empty in-progress
	err = tm.SaveTasks("in-progress", &TaskList{Tasks: []models.Task{}})
	require.NoError(t, err)

	// Claim task
	err = tm.ClaimTask("TASK-001", "test-user")
	assert.NoError(t, err)

	// Verify task moved to in-progress with assignee
	inProgress, err := tm.LoadTasks("in-progress")
	assert.NoError(t, err)
	assert.Len(t, inProgress.Tasks, 1)
	assert.Equal(t, "TASK-001", inProgress.Tasks[0].ID)
	assert.Equal(t, models.StatusInProgress, inProgress.Tasks[0].Status)
	assert.Equal(t, "test-user", inProgress.Tasks[0].AssignedTo)

	// Verify removed from backlog
	backlogAfter, err := tm.LoadTasks("backlog")
	assert.NoError(t, err)
	assert.Empty(t, backlogAfter.Tasks)
}

func TestClaimTask_NotInBacklog(t *testing.T) {
	tmpDir := setupTestDir(t)

	tm := NewTaskManager(tmpDir)

	// Setup empty backlog
	err := tm.SaveTasks("backlog", &TaskList{Tasks: []models.Task{}})
	require.NoError(t, err)

	// Try to claim non-existent task
	err = tm.ClaimTask("NONEXISTENT", "test-user")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found in backlog")
}

func TestClaimTask_AlreadyInProgress(t *testing.T) {
	tmpDir := setupTestDir(t)

	tm := NewTaskManager(tmpDir)

	// Setup task already in progress
	inProgress := &TaskList{
		Tasks: []models.Task{
			{ID: "TASK-001", Title: "Test Task", Status: models.StatusInProgress, AssignedTo: "user1"},
		},
	}
	err := tm.SaveTasks("in-progress", inProgress)
	require.NoError(t, err)

	// Setup empty backlog
	err = tm.SaveTasks("backlog", &TaskList{Tasks: []models.Task{}})
	require.NoError(t, err)

	// Try to claim (should fail since not in backlog)
	err = tm.ClaimTask("TASK-001", "user2")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found in backlog")

	// Verify original assignment preserved
	inProgressAfter, err := tm.LoadTasks("in-progress")
	assert.NoError(t, err)
	assert.Len(t, inProgressAfter.Tasks, 1)
	assert.Equal(t, "user1", inProgressAfter.Tasks[0].AssignedTo)
}

func TestClaimTask_WithEmptyAssignee(t *testing.T) {
	tmpDir := setupTestDir(t)

	tm := NewTaskManager(tmpDir)

	backlog := &TaskList{
		Tasks: []models.Task{
			{ID: "TASK-001", Title: "Test", Status: models.StatusPending},
		},
	}
	err := tm.SaveTasks("backlog", backlog)
	require.NoError(t, err)

	err = tm.SaveTasks("in-progress", &TaskList{Tasks: []models.Task{}})
	require.NoError(t, err)

	// Claim with empty assignee
	err = tm.ClaimTask("TASK-001", "")
	assert.NoError(t, err)

	// Verify task moved
	inProgress, err := tm.LoadTasks("in-progress")
	assert.NoError(t, err)
	assert.Len(t, inProgress.Tasks, 1)
	assert.Empty(t, inProgress.Tasks[0].AssignedTo)
}

func TestClaimTask_PreservesTaskFields(t *testing.T) {
	tmpDir := setupTestDir(t)

	tm := NewTaskManager(tmpDir)

	// Setup task with all fields populated
	backlog := &TaskList{
		Tasks: []models.Task{
			{
				ID:          "TASK-001",
				Title:       "Complex Task",
				Description: "This is a detailed description",
				Status:      models.StatusPending,
				Scope:       []string{"src/module1", "src/module2"},
				SpecRefs:    []string{".agentic/spec/01.md"},
				Inputs:      []string{"input.txt"},
				Outputs:     []string{"output.txt"},
				Acceptance:  []string{"Works correctly", "Has tests"},
				SubTasks: []models.SubTask{
					{ID: "TASK-001.1", Title: "Sub 1", Status: models.StatusPending},
				},
			},
		},
	}
	err := tm.SaveTasks("backlog", backlog)
	require.NoError(t, err)

	err = tm.SaveTasks("in-progress", &TaskList{Tasks: []models.Task{}})
	require.NoError(t, err)

	// Claim task
	err = tm.ClaimTask("TASK-001", "test-user")
	assert.NoError(t, err)

	// Verify all fields preserved
	inProgress, err := tm.LoadTasks("in-progress")
	assert.NoError(t, err)
	assert.Len(t, inProgress.Tasks, 1)

	task := inProgress.Tasks[0]
	assert.Equal(t, "TASK-001", task.ID)
	assert.Equal(t, "Complex Task", task.Title)
	assert.Equal(t, "This is a detailed description", task.Description)
	assert.Equal(t, models.StatusInProgress, task.Status)
	assert.Equal(t, "test-user", task.AssignedTo)
	assert.Equal(t, []string{"src/module1", "src/module2"}, task.Scope)
	assert.Equal(t, []string{".agentic/spec/01.md"}, task.SpecRefs)
	assert.Equal(t, []string{"input.txt"}, task.Inputs)
	assert.Equal(t, []string{"output.txt"}, task.Outputs)
	assert.Equal(t, []string{"Works correctly", "Has tests"}, task.Acceptance)
	assert.Len(t, task.SubTasks, 1)
	assert.Equal(t, "TASK-001.1", task.SubTasks[0].ID)
}
