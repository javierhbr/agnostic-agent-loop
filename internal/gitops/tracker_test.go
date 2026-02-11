package gitops

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// Initialize git repo
	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		require.NoError(t, cmd.Run())
	}

	return dir
}

func commitFile(t *testing.T, dir, filename, content, message string) {
	t.Helper()
	path := filepath.Join(dir, filename)
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))

	cmd := exec.Command("git", "add", filename)
	cmd.Dir = dir
	require.NoError(t, cmd.Run())

	cmd = exec.Command("git", "commit", "-m", message)
	cmd.Dir = dir
	require.NoError(t, cmd.Run())
}

func TestIsGitRepo(t *testing.T) {
	dir := setupGitRepo(t)

	// Change to git repo
	origDir, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(origDir)

	tracker := NewTracker()
	assert.True(t, tracker.IsGitRepo())
}

func TestIsGitRepo_NotARepo(t *testing.T) {
	dir := t.TempDir()

	origDir, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(origDir)

	tracker := NewTracker()
	assert.False(t, tracker.IsGitRepo())
}

func TestGetCurrentBranch(t *testing.T) {
	dir := setupGitRepo(t)
	commitFile(t, dir, "init.txt", "init", "initial commit")

	origDir, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(origDir)

	tracker := NewTracker()
	branch, err := tracker.GetCurrentBranch()
	require.NoError(t, err)
	// Default branch could be "main" or "master" depending on git config
	assert.NotEmpty(t, branch)
}

func TestGetCommitsSince(t *testing.T) {
	dir := setupGitRepo(t)

	before := time.Now().Add(-1 * time.Second)
	commitFile(t, dir, "a.txt", "hello", "add a")
	commitFile(t, dir, "b.txt", "world", "add b")

	origDir, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(origDir)

	tracker := NewTracker()
	commits, err := tracker.GetCommitsSince(before)
	require.NoError(t, err)
	assert.Len(t, commits, 2)
	// Most recent first
	assert.Equal(t, "add b", commits[0].Message)
	assert.Equal(t, "add a", commits[1].Message)
}

func TestGetFilesChangedSince(t *testing.T) {
	dir := setupGitRepo(t)

	before := time.Now().Add(-1 * time.Second)
	commitFile(t, dir, "foo.go", "package foo", "add foo")
	commitFile(t, dir, "bar.go", "package bar", "add bar")

	origDir, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(origDir)

	tracker := NewTracker()
	files, err := tracker.GetFilesChangedSince(before)
	require.NoError(t, err)
	assert.Contains(t, files, "foo.go")
	assert.Contains(t, files, "bar.go")
}

func TestGetCommitFiles(t *testing.T) {
	dir := setupGitRepo(t)

	commitFile(t, dir, "first.go", "package first", "add first")

	origDir, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	defer os.Chdir(origDir)

	tracker := NewTracker()
	commits, err := tracker.GetCommitsSince(time.Now().Add(-5 * time.Second))
	require.NoError(t, err)
	require.NotEmpty(t, commits)

	files, err := tracker.GetCommitFiles(commits[0].Hash)
	require.NoError(t, err)
	assert.Contains(t, files, "first.go")
}

func TestParseCommitLog_Empty(t *testing.T) {
	commits, err := parseCommitLog("")
	require.NoError(t, err)
	assert.Nil(t, commits)
}

func TestParseCommitLog_ValidLines(t *testing.T) {
	log := "abc123|fix bug|Alice|2024-01-15T10:30:00+00:00\ndef456|add feature|Bob|2024-01-16T12:00:00+00:00"
	commits, err := parseCommitLog(log)
	require.NoError(t, err)
	assert.Len(t, commits, 2)
	assert.Equal(t, "abc123", commits[0].Hash)
	assert.Equal(t, "fix bug", commits[0].Message)
	assert.Equal(t, "Alice", commits[0].Author)
}
