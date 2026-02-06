package orchestrator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewArchiver(t *testing.T) {
	archiver := NewArchiver(
		"/tmp/archive",
		"/tmp/progress.txt",
		"/tmp/progress.yaml",
		"/tmp/tasks",
	)

	assert.NotNil(t, archiver)
	assert.Equal(t, "/tmp/archive", archiver.archiveDir)
	assert.Equal(t, "/tmp/progress.txt", archiver.progressText)
	assert.Equal(t, "/tmp/progress.yaml", archiver.progressYAML)
	assert.Equal(t, "/tmp/tasks", archiver.tasksDir)
	assert.Equal(t, "/tmp/.last-branch", archiver.lastBranchFile)
}

func TestArchiver_SaveAndReadLastBranch(t *testing.T) {
	tmpDir := t.TempDir()
	archiver := NewArchiver(
		filepath.Join(tmpDir, "archive"),
		filepath.Join(tmpDir, "progress.txt"),
		filepath.Join(tmpDir, "progress.yaml"),
		filepath.Join(tmpDir, "tasks"),
	)

	// Save branch
	err := archiver.saveLastBranch("feature/test-branch")
	require.NoError(t, err)

	// Read it back
	branch, err := archiver.readLastBranch()
	require.NoError(t, err)
	assert.Equal(t, "feature/test-branch", branch)

	// GetLastBranch should work too
	branch, err = archiver.GetLastBranch()
	require.NoError(t, err)
	assert.Equal(t, "feature/test-branch", branch)
}

func TestArchiver_ArchiveIfBranchChanged_FirstRun(t *testing.T) {
	tmpDir := t.TempDir()
	archiver := NewArchiver(
		filepath.Join(tmpDir, "archive"),
		filepath.Join(tmpDir, "progress.txt"),
		filepath.Join(tmpDir, "progress.yaml"),
		filepath.Join(tmpDir, "tasks"),
	)

	// First run - no previous branch
	err := archiver.ArchiveIfBranchChanged("main")
	require.NoError(t, err)

	// Verify branch was saved
	branch, err := archiver.GetLastBranch()
	require.NoError(t, err)
	assert.Equal(t, "main", branch)

	// No archive should be created
	archives, err := archiver.ListArchives()
	require.NoError(t, err)
	assert.Len(t, archives, 0)
}

func TestArchiver_ArchiveIfBranchChanged_SameBranch(t *testing.T) {
	tmpDir := t.TempDir()
	archiver := NewArchiver(
		filepath.Join(tmpDir, "archive"),
		filepath.Join(tmpDir, "progress.txt"),
		filepath.Join(tmpDir, "progress.yaml"),
		filepath.Join(tmpDir, "tasks"),
	)

	// Set initial branch
	err := archiver.saveLastBranch("feature/auth")
	require.NoError(t, err)

	// Create some progress files
	progressContent := "# Progress Log\n\nSome progress..."
	err = os.WriteFile(archiver.progressText, []byte(progressContent), 0644)
	require.NoError(t, err)

	// Run again with same branch
	err = archiver.ArchiveIfBranchChanged("feature/auth")
	require.NoError(t, err)

	// No archive should be created
	archives, err := archiver.ListArchives()
	require.NoError(t, err)
	assert.Len(t, archives, 0)

	// Progress file should still exist
	assert.FileExists(t, archiver.progressText)
}

func TestArchiver_ArchiveIfBranchChanged_BranchChanged(t *testing.T) {
	tmpDir := t.TempDir()
	archiveDir := filepath.Join(tmpDir, "archive")
	progressText := filepath.Join(tmpDir, "progress.txt")
	progressYAML := filepath.Join(tmpDir, "progress.yaml")
	tasksDir := filepath.Join(tmpDir, "tasks")

	archiver := NewArchiver(archiveDir, progressText, progressYAML, tasksDir)

	// Set initial branch
	err := archiver.saveLastBranch("feature/old-feature")
	require.NoError(t, err)

	// Create progress files
	progressContent := "# Progress Log\n\nOld feature progress..."
	err = os.WriteFile(progressText, []byte(progressContent), 0644)
	require.NoError(t, err)

	yamlContent := "- timestamp: 2026-02-05\n  storyId: US-001"
	err = os.WriteFile(progressYAML, []byte(yamlContent), 0644)
	require.NoError(t, err)

	// Create done tasks
	err = os.MkdirAll(tasksDir, 0755)
	require.NoError(t, err)
	doneTasksPath := filepath.Join(tasksDir, "done.yaml")
	doneContent := "tasks:\n  - id: US-001\n    title: Completed task"
	err = os.WriteFile(doneTasksPath, []byte(doneContent), 0644)
	require.NoError(t, err)

	// Switch to new branch
	err = archiver.ArchiveIfBranchChanged("feature/new-feature")
	require.NoError(t, err)

	// Verify archive was created
	archives, err := archiver.ListArchives()
	require.NoError(t, err)
	require.Len(t, archives, 1)
	assert.Contains(t, archives[0], "old-feature")

	// Verify archived files exist
	archivePath := filepath.Join(archiveDir, archives[0])
	assert.FileExists(t, filepath.Join(archivePath, "progress.txt"))
	assert.FileExists(t, filepath.Join(archivePath, "progress.yaml"))
	assert.FileExists(t, filepath.Join(archivePath, "done.yaml"))

	// Verify archived content
	archivedProgress, err := os.ReadFile(filepath.Join(archivePath, "progress.txt"))
	require.NoError(t, err)
	assert.Contains(t, string(archivedProgress), "Old feature progress")

	// Verify progress files were reset
	resetProgress, err := os.ReadFile(progressText)
	require.NoError(t, err)
	resetProgressStr := string(resetProgress)
	assert.Contains(t, resetProgressStr, "# Progress Log")
	assert.Contains(t, resetProgressStr, "Started:")
	assert.NotContains(t, resetProgressStr, "Old feature")

	// Verify progress.yaml was removed
	assert.NoFileExists(t, progressYAML)

	// Verify branch was updated
	branch, err := archiver.GetLastBranch()
	require.NoError(t, err)
	assert.Equal(t, "feature/new-feature", branch)
}

func TestArchiver_ArchiveIfBranchChanged_WithBranchPrefix(t *testing.T) {
	tmpDir := t.TempDir()
	archiveDir := filepath.Join(tmpDir, "archive")

	archiver := NewArchiver(
		archiveDir,
		filepath.Join(tmpDir, "progress.txt"),
		filepath.Join(tmpDir, "progress.yaml"),
		filepath.Join(tmpDir, "tasks"),
	)

	// Set initial branch with prefix
	err := archiver.saveLastBranch("ralph/feature-one")
	require.NoError(t, err)

	// Create minimal progress file
	err = os.WriteFile(archiver.progressText, []byte("# Progress"), 0644)
	require.NoError(t, err)

	// Switch to new branch
	err = archiver.ArchiveIfBranchChanged("ralph/feature-two")
	require.NoError(t, err)

	// Verify archive uses clean branch name (without prefix)
	archives, err := archiver.ListArchives()
	require.NoError(t, err)
	require.Len(t, archives, 1)
	// Should contain "feature-one" not "ralph/feature-one"
	assert.Contains(t, archives[0], "feature-one")
}

func TestArchiver_CopyFileIfExists(t *testing.T) {
	tmpDir := t.TempDir()
	archiver := NewArchiver(
		filepath.Join(tmpDir, "archive"),
		filepath.Join(tmpDir, "progress.txt"),
		filepath.Join(tmpDir, "progress.yaml"),
		filepath.Join(tmpDir, "tasks"),
	)

	// Test copying existing file
	srcFile := filepath.Join(tmpDir, "source.txt")
	dstFile := filepath.Join(tmpDir, "dest.txt")
	content := "Test content"
	err := os.WriteFile(srcFile, []byte(content), 0644)
	require.NoError(t, err)

	err = archiver.copyFileIfExists(srcFile, dstFile)
	require.NoError(t, err)

	// Verify destination file exists and has same content
	copiedContent, err := os.ReadFile(dstFile)
	require.NoError(t, err)
	assert.Equal(t, content, string(copiedContent))
}

func TestArchiver_CopyFileIfExists_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	archiver := NewArchiver(
		filepath.Join(tmpDir, "archive"),
		filepath.Join(tmpDir, "progress.txt"),
		filepath.Join(tmpDir, "progress.yaml"),
		filepath.Join(tmpDir, "tasks"),
	)

	// Test copying non-existent file (should not error)
	srcFile := filepath.Join(tmpDir, "nonexistent.txt")
	dstFile := filepath.Join(tmpDir, "dest.txt")

	err := archiver.copyFileIfExists(srcFile, dstFile)
	require.NoError(t, err)

	// Destination should not exist
	assert.NoFileExists(t, dstFile)
}

func TestArchiver_ResetProgressFiles(t *testing.T) {
	tmpDir := t.TempDir()
	progressText := filepath.Join(tmpDir, "progress.txt")
	progressYAML := filepath.Join(tmpDir, "progress.yaml")

	archiver := NewArchiver(
		filepath.Join(tmpDir, "archive"),
		progressText,
		progressYAML,
		filepath.Join(tmpDir, "tasks"),
	)

	// Create existing files with old content
	oldProgress := "# Old Progress\n\nOld content..."
	err := os.WriteFile(progressText, []byte(oldProgress), 0644)
	require.NoError(t, err)

	oldYAML := "old: yaml"
	err = os.WriteFile(progressYAML, []byte(oldYAML), 0644)
	require.NoError(t, err)

	// Reset
	err = archiver.resetProgressFiles()
	require.NoError(t, err)

	// Verify progress.txt was reset
	newProgress, err := os.ReadFile(progressText)
	require.NoError(t, err)
	newProgressStr := string(newProgress)
	assert.Contains(t, newProgressStr, "# Progress Log")
	assert.Contains(t, newProgressStr, "Started:")
	assert.NotContains(t, newProgressStr, "Old content")

	// Verify progress.yaml was removed
	assert.NoFileExists(t, progressYAML)
}

func TestArchiver_ListArchives(t *testing.T) {
	tmpDir := t.TempDir()
	archiveDir := filepath.Join(tmpDir, "archive")

	archiver := NewArchiver(
		archiveDir,
		filepath.Join(tmpDir, "progress.txt"),
		filepath.Join(tmpDir, "progress.yaml"),
		filepath.Join(tmpDir, "tasks"),
	)

	// No archives initially
	archives, err := archiver.ListArchives()
	require.NoError(t, err)
	assert.Len(t, archives, 0)

	// Create some archive directories
	err = os.MkdirAll(filepath.Join(archiveDir, "2026-02-05-feature-one"), 0755)
	require.NoError(t, err)
	err = os.MkdirAll(filepath.Join(archiveDir, "2026-02-06-feature-two"), 0755)
	require.NoError(t, err)

	// Create a file (should be ignored)
	err = os.WriteFile(filepath.Join(archiveDir, "readme.txt"), []byte("test"), 0644)
	require.NoError(t, err)

	// List archives
	archives, err = archiver.ListArchives()
	require.NoError(t, err)
	require.Len(t, archives, 2)

	// Verify archive names
	archiveNames := strings.Join(archives, ",")
	assert.Contains(t, archiveNames, "2026-02-05-feature-one")
	assert.Contains(t, archiveNames, "2026-02-06-feature-two")
	assert.NotContains(t, archiveNames, "readme.txt")
}

func TestArchiver_ArchivePreviousRun_MissingFiles(t *testing.T) {
	tmpDir := t.TempDir()
	archiveDir := filepath.Join(tmpDir, "archive")

	archiver := NewArchiver(
		archiveDir,
		filepath.Join(tmpDir, "progress.txt"),
		filepath.Join(tmpDir, "progress.yaml"),
		filepath.Join(tmpDir, "tasks"),
	)

	// Set initial branch
	err := archiver.saveLastBranch("feature/old")
	require.NoError(t, err)

	// Don't create any progress files

	// Switch branch (should succeed even with missing files)
	err = archiver.ArchiveIfBranchChanged("feature/new")
	require.NoError(t, err)

	// Verify archive was created but is empty
	archives, err := archiver.ListArchives()
	require.NoError(t, err)
	require.Len(t, archives, 1)

	archivePath := filepath.Join(archiveDir, archives[0])
	assert.DirExists(t, archivePath)

	// Files should not exist in archive (since they didn't exist originally)
	assert.NoFileExists(t, filepath.Join(archivePath, "progress.txt"))
	assert.NoFileExists(t, filepath.Join(archivePath, "progress.yaml"))
	assert.NoFileExists(t, filepath.Join(archivePath, "done.yaml"))
}
