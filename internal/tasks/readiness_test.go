package tasks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanClaimTask_AllReady(t *testing.T) {
	base := t.TempDir()

	// Create input files
	inputFile := filepath.Join(base, "input.txt")
	require.NoError(t, os.WriteFile(inputFile, []byte("data"), 0644))

	// Create spec dir and file
	specDir := filepath.Join(base, "specs")
	require.NoError(t, os.MkdirAll(specDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(specDir, "auth.md"), []byte("# Auth"), 0644))

	// Create scope dir
	scopeDir := filepath.Join(base, "src")
	require.NoError(t, os.MkdirAll(scopeDir, 0755))

	cfg := &models.Config{}
	cfg.Paths.SpecDirs = []string{specDir}

	task := &models.Task{
		ID:       "T-1",
		Inputs:   []string{inputFile},
		SpecRefs: []string{"auth.md"},
		Scope:    []string{scopeDir},
	}

	result := CanClaimTask(task, cfg)
	assert.True(t, result.Ready)
	assert.Len(t, result.Checks, 3)
	for _, check := range result.Checks {
		assert.True(t, check.Passed)
	}
}

func TestCanClaimTask_MissingInput(t *testing.T) {
	task := &models.Task{
		ID:     "T-2",
		Inputs: []string{"/nonexistent/file.txt"},
	}

	result := CanClaimTask(task, nil)
	assert.False(t, result.Ready)
	assert.Len(t, result.Checks, 1)
	assert.False(t, result.Checks[0].Passed)
	assert.Equal(t, "input-exists", result.Checks[0].Name)
}

func TestCanClaimTask_MissingSpec(t *testing.T) {
	cfg := &models.Config{}
	cfg.Paths.SpecDirs = []string{"/nonexistent/specs"}

	task := &models.Task{
		ID:       "T-3",
		SpecRefs: []string{"missing.md"},
	}

	result := CanClaimTask(task, cfg)
	assert.False(t, result.Ready)
}

func TestCanClaimTask_MissingScopeIsWarningOnly(t *testing.T) {
	task := &models.Task{
		ID:    "T-4",
		Scope: []string{"/nonexistent/dir"},
	}

	result := CanClaimTask(task, nil)
	// Missing scope is a warning, NOT a blocker
	assert.True(t, result.Ready)
	assert.Len(t, result.Checks, 1)
	assert.False(t, result.Checks[0].Passed)
}

func TestCanClaimTask_EmptyTask(t *testing.T) {
	task := &models.Task{ID: "T-5"}
	result := CanClaimTask(task, nil)
	assert.True(t, result.Ready)
	assert.Len(t, result.Checks, 0)
}

func TestFormatReadinessResult(t *testing.T) {
	result := &ReadinessResult{
		TaskID: "T-6",
		Ready:  false,
		Checks: []ReadinessCheck{
			{Name: "input-exists", Passed: true, Message: "ok"},
			{Name: "spec-resolvable", Passed: false, Message: "missing"},
		},
	}

	output := FormatReadinessResult(result)
	assert.Contains(t, output, "NOT READY")
	assert.Contains(t, output, "[+] input-exists")
	assert.Contains(t, output, "[-] spec-resolvable")
}
