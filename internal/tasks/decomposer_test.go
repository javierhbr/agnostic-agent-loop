package tasks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecomposeFromPlan(t *testing.T) {
	dir := t.TempDir()
	taskDir := filepath.Join(dir, "tasks")
	require.NoError(t, os.MkdirAll(taskDir, 0755))

	// Initialize empty backlog
	require.NoError(t, os.WriteFile(filepath.Join(taskDir, "backlog.yaml"), []byte("tasks: []\n"), 0644))

	planPath := filepath.Join(dir, "plan.md")
	planContent := `# Plan: Auth

## Phase 1: Setup

- [ ] Define approach
- [ ] Set up scaffolding

## Phase 2: Implementation

- [ ] Login endpoint
- [ ] Logout endpoint

## Phase 3: Validation

- [ ] Write tests
`
	require.NoError(t, os.WriteFile(planPath, []byte(planContent), 0644))

	tm := NewTaskManager(taskDir)
	created, err := DecomposeFromPlan(planPath, "auth-track", tm)
	require.NoError(t, err)

	assert.Len(t, created, 3, "expected 3 tasks (one per phase)")

	// Each task should be linked to the track
	for _, task := range created {
		assert.Equal(t, "auth-track", task.TrackID)
		assert.NotEmpty(t, task.Acceptance)
	}

	// Check titles contain track ID and phase name
	assert.Contains(t, created[0].Title, "[auth-track]")
	assert.Contains(t, created[0].Title, "Setup")
	assert.Contains(t, created[1].Title, "Implementation")
	assert.Contains(t, created[2].Title, "Validation")

	// Verify tasks are saved in backlog
	backlog, err := tm.LoadTasks("backlog")
	require.NoError(t, err)
	assert.Len(t, backlog.Tasks, 3)
}

func TestDecomposeFromPlan_EmptyPhase(t *testing.T) {
	dir := t.TempDir()
	taskDir := filepath.Join(dir, "tasks")
	require.NoError(t, os.MkdirAll(taskDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(taskDir, "backlog.yaml"), []byte("tasks: []\n"), 0644))

	planPath := filepath.Join(dir, "plan.md")
	planContent := `# Plan: Minimal

## Phase 1: Setup

## Phase 2: Implementation

- [ ] Do the thing
`
	require.NoError(t, os.WriteFile(planPath, []byte(planContent), 0644))

	tm := NewTaskManager(taskDir)
	created, err := DecomposeFromPlan(planPath, "minimal-track", tm)
	require.NoError(t, err)

	// Only Phase 2 has tasks, so only 1 task should be created
	assert.Len(t, created, 1)
	assert.Contains(t, created[0].Title, "Implementation")
}
