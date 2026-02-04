package tasks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "agentic-test-*")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})
	return tmpDir
}

func TestNewTaskManager(t *testing.T) {
	tm := NewTaskManager("/test/path")
	assert.NotNil(t, tm)
	assert.Equal(t, "/test/path", tm.baseDir)
}

func TestLoadTasks_EmptyFile(t *testing.T) {
	tmpDir := setupTestDir(t)

	// Create empty backlog file
	err := os.WriteFile(filepath.Join(tmpDir, "backlog.yaml"), []byte("tasks: []\n"), 0644)
	require.NoError(t, err)

	tm := NewTaskManager(tmpDir)
	list, err := tm.LoadTasks("backlog")

	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Empty(t, list.Tasks)
}

func TestLoadTasks_NonExistent(t *testing.T) {
	tmpDir := setupTestDir(t)

	tm := NewTaskManager(tmpDir)
	list, err := tm.LoadTasks("nonexistent")

	// Should return empty list, not error
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Empty(t, list.Tasks)
}

func TestLoadTasks_MalformedYAML(t *testing.T) {
	tmpDir := setupTestDir(t)

	// Create malformed YAML file
	err := os.WriteFile(filepath.Join(tmpDir, "backlog.yaml"), []byte("invalid: [\n"), 0644)
	require.NoError(t, err)

	tm := NewTaskManager(tmpDir)
	_, err = tm.LoadTasks("backlog")

	assert.Error(t, err)
}

func TestLoadTasks_ValidTasks(t *testing.T) {
	tmpDir := setupTestDir(t)

	// Create valid YAML with tasks
	yamlContent := `tasks:
  - id: TASK-001
    title: Test Task
    status: pending
    spec_refs:
      - spec1.md
      - spec2.md
    acceptance:
      - Criterion 1
      - Criterion 2
`
	err := os.WriteFile(filepath.Join(tmpDir, "backlog.yaml"), []byte(yamlContent), 0644)
	require.NoError(t, err)

	tm := NewTaskManager(tmpDir)
	list, err := tm.LoadTasks("backlog")

	assert.NoError(t, err)
	assert.Len(t, list.Tasks, 1)
	assert.Equal(t, "TASK-001", list.Tasks[0].ID)
	assert.Equal(t, "Test Task", list.Tasks[0].Title)
	assert.Equal(t, models.StatusPending, list.Tasks[0].Status)
	assert.Len(t, list.Tasks[0].SpecRefs, 2)
	assert.Len(t, list.Tasks[0].Acceptance, 2)
}

func TestSaveTasks_Success(t *testing.T) {
	tmpDir := setupTestDir(t)

	tm := NewTaskManager(tmpDir)
	list := &TaskList{
		Tasks: []models.Task{
			{
				ID:     "TASK-001",
				Title:  "Test Task",
				Status: models.StatusPending,
			},
		},
	}

	err := tm.SaveTasks("backlog", list)
	assert.NoError(t, err)

	// Verify file exists and contains correct data
	data, err := os.ReadFile(filepath.Join(tmpDir, "backlog.yaml"))
	assert.NoError(t, err)
	assert.Contains(t, string(data), "TASK-001")
	assert.Contains(t, string(data), "Test Task")
}

func TestCreateTask_Success(t *testing.T) {
	tmpDir := setupTestDir(t)

	// Create empty backlog
	err := os.WriteFile(filepath.Join(tmpDir, "backlog.yaml"), []byte("tasks: []\n"), 0644)
	require.NoError(t, err)

	tm := NewTaskManager(tmpDir)
	task, err := tm.CreateTask("New Task")

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.NotEmpty(t, task.ID)
	assert.Equal(t, "New Task", task.Title)
	assert.Equal(t, models.StatusPending, task.Status)

	// Verify task was added to backlog
	list, err := tm.LoadTasks("backlog")
	assert.NoError(t, err)
	assert.Len(t, list.Tasks, 1)
	assert.Equal(t, task.ID, list.Tasks[0].ID)
}

func TestCreateTask_EmptyTitle(t *testing.T) {
	tmpDir := setupTestDir(t)

	err := os.WriteFile(filepath.Join(tmpDir, "backlog.yaml"), []byte("tasks: []\n"), 0644)
	require.NoError(t, err)

	tm := NewTaskManager(tmpDir)
	task, err := tm.CreateTask("")

	// Should still create task (validation is in CLI layer)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Empty(t, task.Title)
}

func TestMoveTask_BacklogToInProgress(t *testing.T) {
	tmpDir := setupTestDir(t)

	// Setup backlog with task
	backlog := &TaskList{
		Tasks: []models.Task{
			{ID: "TASK-001", Title: "Test", Status: models.StatusPending},
		},
	}
	tm := NewTaskManager(tmpDir)
	err := tm.SaveTasks("backlog", backlog)
	require.NoError(t, err)

	// Create empty in-progress
	err = tm.SaveTasks("in-progress", &TaskList{Tasks: []models.Task{}})
	require.NoError(t, err)

	// Move task
	err = tm.MoveTask("TASK-001", "backlog", "in-progress", models.StatusInProgress)
	assert.NoError(t, err)

	// Verify task moved
	backlogList, err := tm.LoadTasks("backlog")
	assert.NoError(t, err)
	assert.Empty(t, backlogList.Tasks)

	inProgressList, err := tm.LoadTasks("in-progress")
	assert.NoError(t, err)
	assert.Len(t, inProgressList.Tasks, 1)
	assert.Equal(t, "TASK-001", inProgressList.Tasks[0].ID)
	assert.Equal(t, models.StatusInProgress, inProgressList.Tasks[0].Status)
}

func TestMoveTask_TaskNotFound(t *testing.T) {
	tmpDir := setupTestDir(t)

	tm := NewTaskManager(tmpDir)
	err := tm.SaveTasks("backlog", &TaskList{Tasks: []models.Task{}})
	require.NoError(t, err)

	err = tm.MoveTask("NONEXISTENT", "backlog", "in-progress", models.StatusInProgress)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestMoveTask_InProgressToDone(t *testing.T) {
	tmpDir := setupTestDir(t)

	tm := NewTaskManager(tmpDir)
	inProgress := &TaskList{
		Tasks: []models.Task{
			{ID: "TASK-001", Title: "Test", Status: models.StatusInProgress, AssignedTo: "user"},
		},
	}
	err := tm.SaveTasks("in-progress", inProgress)
	require.NoError(t, err)

	err = tm.SaveTasks("done", &TaskList{Tasks: []models.Task{}})
	require.NoError(t, err)

	err = tm.MoveTask("TASK-001", "in-progress", "done", models.StatusDone)
	assert.NoError(t, err)

	// Verify move
	doneList, err := tm.LoadTasks("done")
	assert.NoError(t, err)
	assert.Len(t, doneList.Tasks, 1)
	assert.Equal(t, models.StatusDone, doneList.Tasks[0].Status)
	assert.Equal(t, "user", doneList.Tasks[0].AssignedTo) // Should preserve assigned user
}

func TestFindTask_InBacklog(t *testing.T) {
	tmpDir := setupTestDir(t)

	tm := NewTaskManager(tmpDir)
	backlog := &TaskList{
		Tasks: []models.Task{
			{ID: "TASK-001", Title: "Test Task", Status: models.StatusPending},
		},
	}
	err := tm.SaveTasks("backlog", backlog)
	require.NoError(t, err)

	task, source, err := tm.FindTask("TASK-001")

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "TASK-001", task.ID)
	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, "backlog", source)
}

func TestFindTask_InInProgress(t *testing.T) {
	tmpDir := setupTestDir(t)

	tm := NewTaskManager(tmpDir)
	inProgress := &TaskList{
		Tasks: []models.Task{
			{ID: "TASK-002", Title: "In Progress Task", Status: models.StatusInProgress},
		},
	}
	err := tm.SaveTasks("in-progress", inProgress)
	require.NoError(t, err)

	task, source, err := tm.FindTask("TASK-002")

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "in-progress", source)
}

func TestFindTask_NotFound(t *testing.T) {
	tmpDir := setupTestDir(t)

	tm := NewTaskManager(tmpDir)
	err := tm.SaveTasks("backlog", &TaskList{Tasks: []models.Task{}})
	require.NoError(t, err)

	task, source, err := tm.FindTask("NONEXISTENT")

	assert.NoError(t, err) // No error, just not found
	assert.Nil(t, task)
	assert.Empty(t, source)
}

func TestFindTask_WithSubtasks(t *testing.T) {
	tmpDir := setupTestDir(t)

	tm := NewTaskManager(tmpDir)
	backlog := &TaskList{
		Tasks: []models.Task{
			{
				ID:     "TASK-001",
				Title:  "Parent Task",
				Status: models.StatusPending,
				SubTasks: []models.SubTask{
					{ID: "TASK-001.1", Title: "Subtask 1", Status: models.StatusPending},
					{ID: "TASK-001.2", Title: "Subtask 2", Status: models.StatusDone},
				},
			},
		},
	}
	err := tm.SaveTasks("backlog", backlog)
	require.NoError(t, err)

	// Find parent task
	task, source, err := tm.FindTask("TASK-001")
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "TASK-001", task.ID)
	assert.Len(t, task.SubTasks, 2)

	// Find subtask
	subtask, source, err := tm.FindTask("TASK-001.1")
	assert.NoError(t, err)
	assert.NotNil(t, subtask)
	assert.Equal(t, "TASK-001.1", subtask.ID)
	assert.Equal(t, "Subtask 1", subtask.Title)
	assert.Equal(t, "backlog", source)
}
