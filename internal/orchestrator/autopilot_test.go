package orchestrator

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/javierbenavides/agentic-agent/internal/config"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func setupAutopilotTestDir(t *testing.T) (string, *models.Config) {
	t.Helper()
	base := t.TempDir()

	tasksDir := filepath.Join(base, ".agentic", "tasks")
	specDir := filepath.Join(base, ".agentic", "spec")
	contextDir := filepath.Join(base, ".agentic", "context")
	require.NoError(t, os.MkdirAll(tasksDir, 0755))
	require.NoError(t, os.MkdirAll(specDir, 0755))
	require.NoError(t, os.MkdirAll(contextDir, 0755))

	cfg := &models.Config{}
	config.SetDefaults(cfg)
	cfg.Paths.SpecDirs = []string{specDir}
	cfg.Paths.ContextDirs = []string{contextDir}

	return base, cfg
}

func writeTasksFile(t *testing.T, dir, listType string, taskList tasks.TaskList) {
	t.Helper()
	data, err := yaml.Marshal(taskList)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, listType+".yaml"), data, 0644))
}

func TestAutopilotLoop_DryRun(t *testing.T) {
	base, cfg := setupAutopilotTestDir(t)
	tasksDir := filepath.Join(base, ".agentic", "tasks")

	// Create a backlog task
	writeTasksFile(t, tasksDir, "backlog", tasks.TaskList{
		Tasks: []models.Task{
			{ID: "T-1", Title: "Test task", Status: models.StatusPending},
		},
	})
	writeTasksFile(t, tasksDir, "in-progress", tasks.TaskList{})

	// Override the task manager base dir
	loop := NewAutopilotLoop(cfg, 3, "", true)
	loop.taskManager = tasks.NewTaskManager(tasksDir)

	err := loop.Run(context.Background())
	assert.NoError(t, err)

	// Task should still be in backlog (dry run doesn't claim)
	tm := tasks.NewTaskManager(tasksDir)
	backlog, _ := tm.LoadTasks("backlog")
	assert.Len(t, backlog.Tasks, 1)
}

func TestAutopilotLoop_AllTasksComplete(t *testing.T) {
	base, cfg := setupAutopilotTestDir(t)
	tasksDir := filepath.Join(base, ".agentic", "tasks")

	// Empty backlog and in-progress
	writeTasksFile(t, tasksDir, "backlog", tasks.TaskList{})
	writeTasksFile(t, tasksDir, "in-progress", tasks.TaskList{})
	writeTasksFile(t, tasksDir, "done", tasks.TaskList{
		Tasks: []models.Task{
			{ID: "T-1", Title: "Completed task", Status: models.StatusDone},
		},
	})

	loop := NewAutopilotLoop(cfg, 5, "", false)
	loop.taskManager = tasks.NewTaskManager(tasksDir)

	err := loop.Run(context.Background())
	assert.NoError(t, err)
}

func TestAutopilotLoop_ContextCancellation(t *testing.T) {
	base, cfg := setupAutopilotTestDir(t)
	tasksDir := filepath.Join(base, ".agentic", "tasks")

	writeTasksFile(t, tasksDir, "backlog", tasks.TaskList{
		Tasks: []models.Task{
			{ID: "T-1", Title: "Task 1", Status: models.StatusPending},
		},
	})
	writeTasksFile(t, tasksDir, "in-progress", tasks.TaskList{})

	loop := NewAutopilotLoop(cfg, 100, "", true)
	loop.taskManager = tasks.NewTaskManager(tasksDir)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := loop.Run(ctx)
	assert.Error(t, err)
}

func TestAutopilotLoop_MaxIterations(t *testing.T) {
	base, cfg := setupAutopilotTestDir(t)
	tasksDir := filepath.Join(base, ".agentic", "tasks")

	// Create multiple tasks that won't be claimed in dry run
	writeTasksFile(t, tasksDir, "backlog", tasks.TaskList{
		Tasks: []models.Task{
			{ID: "T-1", Title: "Task 1", Status: models.StatusPending},
			{ID: "T-2", Title: "Task 2", Status: models.StatusPending},
			{ID: "T-3", Title: "Task 3", Status: models.StatusPending},
		},
	})
	writeTasksFile(t, tasksDir, "in-progress", tasks.TaskList{})

	loop := NewAutopilotLoop(cfg, 2, "", true)
	loop.taskManager = tasks.NewTaskManager(tasksDir)

	err := loop.Run(context.Background())
	assert.NoError(t, err) // Dry run doesn't error on max iterations
}
