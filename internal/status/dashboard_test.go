package status

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/javierbenavides/agentic-agent/internal/config"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDir(t *testing.T) (string, *tasks.TaskManager, *models.Config) {
	t.Helper()
	dir := t.TempDir()
	tasksDir := filepath.Join(dir, ".agentic", "tasks")
	require.NoError(t, os.MkdirAll(tasksDir, 0755))

	// Create empty task lists
	for _, name := range []string{"backlog", "in-progress", "done"} {
		require.NoError(t, os.WriteFile(
			filepath.Join(tasksDir, name+".yaml"),
			[]byte("tasks: []\n"),
			0644,
		))
	}

	tm := tasks.NewTaskManager(tasksDir)
	cfg := &models.Config{}
	config.SetDefaults(cfg)
	cfg.Project.Name = "test-project"

	return dir, tm, cfg
}

func TestGather_EmptyProject(t *testing.T) {
	_, tm, cfg := setupTestDir(t)

	d, err := Gather(tm, cfg)
	require.NoError(t, err)
	assert.Equal(t, "test-project", d.ProjectName)
	assert.Equal(t, 0, d.TotalCount)
	assert.Equal(t, float64(0), d.CompletionPct)
	assert.Nil(t, d.NextReady)
	assert.Empty(t, d.Blockers)
}

func TestGather_WithTasks(t *testing.T) {
	_, tm, cfg := setupTestDir(t)

	// Create tasks via manager
	_, err := tm.CreateTask("Task A")
	require.NoError(t, err)
	_, err = tm.CreateTask("Task B")
	require.NoError(t, err)

	d, err := Gather(tm, cfg)
	require.NoError(t, err)
	assert.Equal(t, 2, d.BacklogCount)
	assert.Equal(t, 0, d.InProgressCount)
	assert.Equal(t, 0, d.DoneCount)
	assert.Equal(t, 2, d.TotalCount)
	assert.Equal(t, float64(0), d.CompletionPct)
	// Both tasks should be ready (no inputs, no spec refs)
	assert.NotNil(t, d.NextReady)
}

func TestGather_CompletionPercentage(t *testing.T) {
	_, tm, cfg := setupTestDir(t)

	// Write tasks directly to avoid timestamp collision from CreateTask
	backlog, err := tm.LoadTasks("backlog")
	require.NoError(t, err)
	backlog.Tasks = append(backlog.Tasks,
		models.Task{ID: "TASK-A", Title: "Task A", Status: models.StatusPending},
		models.Task{ID: "TASK-B", Title: "Task B", Status: models.StatusPending},
	)
	require.NoError(t, tm.SaveTasks("backlog", backlog))

	// Complete one task
	require.NoError(t, tm.MoveTask("TASK-A", "backlog", "done", models.StatusDone))

	d, err := Gather(tm, cfg)
	require.NoError(t, err)
	assert.Equal(t, 1, d.BacklogCount)
	assert.Equal(t, 1, d.DoneCount)
	assert.Equal(t, 50.0, d.CompletionPct)

	// Complete both
	require.NoError(t, tm.MoveTask("TASK-B", "backlog", "done", models.StatusDone))

	d, err = Gather(tm, cfg)
	require.NoError(t, err)
	assert.Equal(t, 100.0, d.CompletionPct)
}

func TestGather_BlockedTasks(t *testing.T) {
	dir, tm, cfg := setupTestDir(t)

	// Create spec dir so resolver works
	specDir := filepath.Join(dir, ".agentic", "spec")
	require.NoError(t, os.MkdirAll(specDir, 0755))
	cfg.Paths.SpecDirs = []string{specDir}

	// Create a task with a missing input
	backlog, err := tm.LoadTasks("backlog")
	require.NoError(t, err)
	backlog.Tasks = append(backlog.Tasks, models.Task{
		ID:     "TASK-blocked",
		Title:  "Blocked task",
		Status: models.StatusPending,
		Inputs: []string{"nonexistent-file.go"},
	})
	require.NoError(t, tm.SaveTasks("backlog", backlog))

	d, err := Gather(tm, cfg)
	require.NoError(t, err)
	assert.Nil(t, d.NextReady)
	assert.NotEmpty(t, d.Blockers)
	assert.Contains(t, d.Blockers[0], "nonexistent-file.go")
}
