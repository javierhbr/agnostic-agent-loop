package plans

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const samplePlan = `# Plan: Add Authentication

## Phase 1: Data Layer
- [x] Create User model with password hashing
- [ ] Create UserRepository interface
- [ ] Implement PostgreSQL UserRepository

## Phase 2: API Layer
- [ ] Create auth middleware
- [~] Implement login endpoint
- [ ] Implement logout endpoint

## Phase 3: Testing
- [ ] Unit tests for User model
- [ ] Integration tests for auth flow
`

func writePlan(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "plan.md")
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	return path
}

func TestParseFile(t *testing.T) {
	path := writePlan(t, samplePlan)

	plan, err := ParseFile(path)
	require.NoError(t, err)

	assert.Equal(t, "Add Authentication", plan.Title)
	assert.Len(t, plan.Phases, 3)

	// Phase 1
	assert.Equal(t, "Phase 1: Data Layer", plan.Phases[0].Name)
	assert.Len(t, plan.Phases[0].Tasks, 3)
	assert.Equal(t, PlanTaskDone, plan.Phases[0].Tasks[0].Status)
	assert.Equal(t, "Create User model with password hashing", plan.Phases[0].Tasks[0].Title)
	assert.Equal(t, PlanTaskPending, plan.Phases[0].Tasks[1].Status)
	assert.Equal(t, PlanTaskPending, plan.Phases[0].Tasks[2].Status)

	// Phase 2
	assert.Equal(t, "Phase 2: API Layer", plan.Phases[1].Name)
	assert.Len(t, plan.Phases[1].Tasks, 3)
	assert.Equal(t, PlanTaskPending, plan.Phases[1].Tasks[0].Status)
	assert.Equal(t, PlanTaskInProgress, plan.Phases[1].Tasks[1].Status)
	assert.Equal(t, PlanTaskPending, plan.Phases[1].Tasks[2].Status)

	// Phase 3
	assert.Len(t, plan.Phases[2].Tasks, 2)
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/plan.md")
	assert.Error(t, err)
}

func TestNextTask(t *testing.T) {
	path := writePlan(t, samplePlan)
	plan, err := ParseFile(path)
	require.NoError(t, err)

	task, phase := plan.NextTask()
	require.NotNil(t, task)
	assert.Equal(t, "Create UserRepository interface", task.Title)
	assert.Equal(t, "Phase 1: Data Layer", phase.Name)
}

func TestNextTask_AllDone(t *testing.T) {
	content := `# Plan: Done
## Phase 1
- [x] Task A
- [x] Task B
`
	path := writePlan(t, content)
	plan, err := ParseFile(path)
	require.NoError(t, err)

	task, phase := plan.NextTask()
	assert.Nil(t, task)
	assert.Nil(t, phase)
}

func TestProgress(t *testing.T) {
	path := writePlan(t, samplePlan)
	plan, err := ParseFile(path)
	require.NoError(t, err)

	done, total := plan.Progress()
	assert.Equal(t, 1, done)
	assert.Equal(t, 8, total)
}

func TestUpdateTaskStatus(t *testing.T) {
	path := writePlan(t, samplePlan)
	plan, err := ParseFile(path)
	require.NoError(t, err)

	// Mark "Create UserRepository interface" (pending) as done
	task := plan.Phases[0].Tasks[1]
	assert.Equal(t, PlanTaskPending, task.Status)

	err = UpdateTaskStatus(path, task.Line, PlanTaskDone)
	require.NoError(t, err)

	// Re-parse and verify
	plan2, err := ParseFile(path)
	require.NoError(t, err)
	assert.Equal(t, PlanTaskDone, plan2.Phases[0].Tasks[1].Status)
}

func TestUpdateTaskStatus_InProgress(t *testing.T) {
	path := writePlan(t, samplePlan)
	plan, err := ParseFile(path)
	require.NoError(t, err)

	// Mark first pending task as in-progress
	task := plan.Phases[0].Tasks[1]
	err = UpdateTaskStatus(path, task.Line, PlanTaskInProgress)
	require.NoError(t, err)

	plan2, err := ParseFile(path)
	require.NoError(t, err)
	assert.Equal(t, PlanTaskInProgress, plan2.Phases[0].Tasks[1].Status)
}

func TestUpdateTaskStatus_InvalidLine(t *testing.T) {
	path := writePlan(t, samplePlan)

	err := UpdateTaskStatus(path, 999, PlanTaskDone)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "out of range")
}

func TestParseFile_EmptyPlan(t *testing.T) {
	path := writePlan(t, "# Empty Plan\n")
	plan, err := ParseFile(path)
	require.NoError(t, err)
	assert.Equal(t, "Empty Plan", plan.Title)
	assert.Empty(t, plan.Phases)
}
