package orchestrator

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestNewLoop(t *testing.T) {
	tm := tasks.NewTaskManager("/tmp/tasks")

	// Test with valid parameters
	loop := NewLoop(5, "<stop>", tm)
	assert.NotNil(t, loop)
	assert.Equal(t, 5, loop.maxIterations)
	assert.Equal(t, "<stop>", loop.stopSignal)
	assert.Equal(t, tm, loop.taskManager)
}

func TestNewLoop_DefaultValues(t *testing.T) {
	tm := tasks.NewTaskManager("/tmp/tasks")

	// Test with empty stop signal
	loop := NewLoop(5, "", tm)
	assert.Equal(t, "<promise>COMPLETE</promise>", loop.stopSignal)

	// Test with zero max iterations
	loop = NewLoop(0, "<stop>", tm)
	assert.Equal(t, 10, loop.maxIterations)

	// Test with negative max iterations
	loop = NewLoop(-5, "<stop>", tm)
	assert.Equal(t, 10, loop.maxIterations)
}

func TestLoop_CheckStopCondition(t *testing.T) {
	tm := tasks.NewTaskManager("/tmp/tasks")
	loop := NewLoop(10, "<promise>COMPLETE</promise>", tm)

	// Test matching stop condition
	assert.True(t, loop.checkStopCondition("Task completed. <promise>COMPLETE</promise> Ready for next."))
	assert.True(t, loop.checkStopCondition("<promise>COMPLETE</promise>"))

	// Test non-matching
	assert.False(t, loop.checkStopCondition("Task in progress"))
	assert.False(t, loop.checkStopCondition("COMPLETE but no promise tags"))
	assert.False(t, loop.checkStopCondition(""))
}

func TestLoop_AllTasksComplete_EmptyBacklog(t *testing.T) {
	tmpDir := t.TempDir()
	tasksDir := filepath.Join(tmpDir, "tasks")
	err := os.MkdirAll(tasksDir, 0755)
	require.NoError(t, err)

	// Create empty task lists
	emptyList := tasks.TaskList{Tasks: []models.Task{}}

	backlogData, err := yaml.Marshal(emptyList)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tasksDir, "backlog.yaml"), backlogData, 0644)
	require.NoError(t, err)

	inProgressData, err := yaml.Marshal(emptyList)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tasksDir, "in-progress.yaml"), inProgressData, 0644)
	require.NoError(t, err)

	tm := tasks.NewTaskManager(tasksDir)
	loop := NewLoop(10, "<stop>", tm)

	// All tasks should be complete
	assert.True(t, loop.allTasksComplete())
}

func TestLoop_AllTasksComplete_TasksRemaining(t *testing.T) {
	tmpDir := t.TempDir()
	tasksDir := filepath.Join(tmpDir, "tasks")
	err := os.MkdirAll(tasksDir, 0755)
	require.NoError(t, err)

	// Create backlog with tasks
	backlogTasks := tasks.TaskList{
		Tasks: []models.Task{
			{ID: "US-001", Title: "Pending task", Status: "backlog"},
		},
	}
	backlogData, err := yaml.Marshal(backlogTasks)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tasksDir, "backlog.yaml"), backlogData, 0644)
	require.NoError(t, err)

	// Empty in-progress
	emptyList := tasks.TaskList{Tasks: []models.Task{}}
	inProgressData, err := yaml.Marshal(emptyList)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tasksDir, "in-progress.yaml"), inProgressData, 0644)
	require.NoError(t, err)

	tm := tasks.NewTaskManager(tasksDir)
	loop := NewLoop(10, "<stop>", tm)

	// Tasks remaining
	assert.False(t, loop.allTasksComplete())
}

func TestLoop_AllTasksComplete_InProgressTasks(t *testing.T) {
	tmpDir := t.TempDir()
	tasksDir := filepath.Join(tmpDir, "tasks")
	err := os.MkdirAll(tasksDir, 0755)
	require.NoError(t, err)

	// Empty backlog
	emptyList := tasks.TaskList{Tasks: []models.Task{}}
	backlogData, err := yaml.Marshal(emptyList)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tasksDir, "backlog.yaml"), backlogData, 0644)
	require.NoError(t, err)

	// In-progress with tasks
	inProgressTasks := tasks.TaskList{
		Tasks: []models.Task{
			{ID: "US-002", Title: "Active task", Status: "in-progress"},
		},
	}
	inProgressData, err := yaml.Marshal(inProgressTasks)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tasksDir, "in-progress.yaml"), inProgressData, 0644)
	require.NoError(t, err)

	tm := tasks.NewTaskManager(tasksDir)
	loop := NewLoop(10, "<stop>", tm)

	// Tasks in progress
	assert.False(t, loop.allTasksComplete())
}

func TestLoop_RunIteration(t *testing.T) {
	tm := tasks.NewTaskManager("/tmp/tasks")
	loop := NewLoop(10, "<stop>", tm)

	ctx := context.Background()
	output, err := loop.runIteration(ctx)

	require.NoError(t, err)
	assert.NotEmpty(t, output)
}

func TestLoop_RunIteration_ContextCanceled(t *testing.T) {
	tm := tasks.NewTaskManager("/tmp/tasks")
	loop := NewLoop(10, "<stop>", tm)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := loop.runIteration(ctx)
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

func TestLoop_RunIteration_ContextTimeout(t *testing.T) {
	tm := tasks.NewTaskManager("/tmp/tasks")
	loop := NewLoop(10, "<stop>", tm)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(2 * time.Millisecond) // Ensure timeout

	_, err := loop.runIteration(ctx)
	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
}

// Note: Testing Run with stop condition would require mocking or integration testing
// since runIteration is a private method that simulates work.
// The checkStopCondition method is tested separately above.

func TestLoop_Run_AllTasksComplete(t *testing.T) {
	tmpDir := t.TempDir()
	tasksDir := filepath.Join(tmpDir, "tasks")
	err := os.MkdirAll(tasksDir, 0755)
	require.NoError(t, err)

	// Create empty task lists
	emptyList := tasks.TaskList{Tasks: []models.Task{}}

	backlogData, err := yaml.Marshal(emptyList)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tasksDir, "backlog.yaml"), backlogData, 0644)
	require.NoError(t, err)

	inProgressData, err := yaml.Marshal(emptyList)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tasksDir, "in-progress.yaml"), inProgressData, 0644)
	require.NoError(t, err)

	tm := tasks.NewTaskManager(tasksDir)
	loop := NewLoop(10, "<stop>", tm)

	ctx := context.Background()
	err = loop.Run(ctx)

	// Should exit successfully
	require.NoError(t, err)
}

func TestLoop_Run_MaxIterationsReached(t *testing.T) {
	tmpDir := t.TempDir()
	tasksDir := filepath.Join(tmpDir, "tasks")
	err := os.MkdirAll(tasksDir, 0755)
	require.NoError(t, err)

	// Create backlog with tasks (so it never completes)
	backlogTasks := tasks.TaskList{
		Tasks: []models.Task{
			{ID: "US-001", Title: "Endless task", Status: "backlog"},
		},
	}
	backlogData, err := yaml.Marshal(backlogTasks)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tasksDir, "backlog.yaml"), backlogData, 0644)
	require.NoError(t, err)

	emptyList := tasks.TaskList{Tasks: []models.Task{}}
	inProgressData, err := yaml.Marshal(emptyList)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tasksDir, "in-progress.yaml"), inProgressData, 0644)
	require.NoError(t, err)

	tm := tasks.NewTaskManager(tasksDir)
	loop := NewLoop(3, "<stop>", tm) // Only 3 iterations

	ctx := context.Background()
	err = loop.Run(ctx)

	// Should error with max iterations message
	require.Error(t, err)
	assert.Contains(t, err.Error(), "reached max iterations")
	assert.Contains(t, err.Error(), "3")
}
