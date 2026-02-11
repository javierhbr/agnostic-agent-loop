package tracks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir)

	track, err := m.Create("Add Authentication", models.TrackTypeFeature, nil)
	require.NoError(t, err)

	assert.Equal(t, "add-authentication", track.ID)
	assert.Equal(t, "Add Authentication", track.Name)
	assert.Equal(t, models.TrackTypeFeature, track.Type)
	assert.Equal(t, models.TrackStatusIdeation, track.Status)
	assert.NotEmpty(t, track.BrainstormPath)

	// Verify files exist
	assert.FileExists(t, filepath.Join(dir, "add-authentication", "spec.md"))
	assert.FileExists(t, filepath.Join(dir, "add-authentication", "plan.md"))
	assert.FileExists(t, filepath.Join(dir, "add-authentication", "brainstorm.md"))
	assert.FileExists(t, filepath.Join(dir, "add-authentication", "metadata.yaml"))
	assert.FileExists(t, filepath.Join(dir, "tracks.yaml"))
}

func TestCreate_Duplicate(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir)

	_, err := m.Create("My Feature", models.TrackTypeFeature, nil)
	require.NoError(t, err)

	_, err = m.Create("My Feature", models.TrackTypeFeature, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestList(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir)

	_, err := m.Create("Feature A", models.TrackTypeFeature, nil)
	require.NoError(t, err)
	_, err = m.Create("Bug Fix B", models.TrackTypeBug, nil)
	require.NoError(t, err)

	tracks, err := m.List()
	require.NoError(t, err)
	assert.Len(t, tracks, 2)
	assert.Equal(t, "feature-a", tracks[0].ID)
	assert.Equal(t, "bug-fix-b", tracks[1].ID)
}

func TestList_Empty(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir)

	tracks, err := m.List()
	require.NoError(t, err)
	assert.Empty(t, tracks)
}

func TestGet(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir)

	_, err := m.Create("My Track", models.TrackTypeRefactor, nil)
	require.NoError(t, err)

	track, err := m.Get("my-track")
	require.NoError(t, err)
	assert.Equal(t, "My Track", track.Name)
	assert.Equal(t, models.TrackTypeRefactor, track.Type)
}

func TestGet_NotFound(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir)

	_, err := m.Get("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestUpdateStatus(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir)

	_, err := m.Create("Test Track", models.TrackTypeFeature, nil)
	require.NoError(t, err)

	err = m.UpdateStatus("test-track", models.TrackStatusActive)
	require.NoError(t, err)

	track, err := m.Get("test-track")
	require.NoError(t, err)
	assert.Equal(t, models.TrackStatusActive, track.Status)
}

func TestAddTask(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir)

	_, err := m.Create("Test Track", models.TrackTypeFeature, nil)
	require.NoError(t, err)

	require.NoError(t, m.AddTask("test-track", "TASK-001"))
	require.NoError(t, m.AddTask("test-track", "TASK-002"))
	// Duplicate should be a no-op
	require.NoError(t, m.AddTask("test-track", "TASK-001"))

	track, err := m.Get("test-track")
	require.NoError(t, err)
	assert.Equal(t, []string{"TASK-001", "TASK-002"}, track.TaskIDs)
}

func TestArchive(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir)

	_, err := m.Create("Old Feature", models.TrackTypeFeature, nil)
	require.NoError(t, err)

	require.NoError(t, m.Archive("old-feature"))

	// Original dir should be gone
	_, err = os.Stat(filepath.Join(dir, "old-feature"))
	assert.True(t, os.IsNotExist(err))

	// Archive dir should have it
	assert.DirExists(t, filepath.Join(dir, "_archive", "old-feature"))

	// Status should be archived
	track, err := m.Get("old-feature")
	require.NoError(t, err)
	assert.Equal(t, models.TrackStatusArchived, track.Status)
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Add Authentication", "add-authentication"},
		{"  Fix Bug #123  ", "fix-bug-123"},
		{"already-kebab", "already-kebab"},
		{"CamelCaseHere", "camelcasehere"},
		{"multiple   spaces", "multiple-spaces"},
		{"special!@#chars", "special-chars"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, toKebabCase(tt.input))
		})
	}
}
