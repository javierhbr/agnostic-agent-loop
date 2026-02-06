package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/javierbenavides/agentic-agent/internal/config"
	"github.com/javierbenavides/agentic-agent/internal/encoding"
	"github.com/javierbenavides/agentic-agent/internal/project"
	"github.com/javierbenavides/agentic-agent/internal/specs"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSpecWorkflow tests the full spec-driven workflow:
// init → create spec files → create task with spec refs → claim with readiness → build context with specs
func TestSpecWorkflow(t *testing.T) {
	tmpDir := setupIntegrationTest(t)

	// Step 1: Initialize project
	err := project.InitProject("SpecWorkflowTest")
	require.NoError(t, err)

	// Step 2: Create config with defaults (the init template uses Go template syntax
	// that isn't valid YAML, so we build config with defaults like getConfig() does)
	cfg := &models.Config{}
	config.SetDefaults(cfg)
	assert.Equal(t, []string{".agentic/spec"}, cfg.Paths.SpecDirs)

	// Step 3: Create spec files in the default spec directory
	specDir := filepath.Join(tmpDir, ".agentic", "spec")
	require.DirExists(t, specDir)

	authSpec := `# Authentication Specification

## Overview
JWT-based authentication for the API.

## Requirements
- RS256 signing
- 15 minute access token TTL
- 7 day refresh token TTL
`
	apiSpec := `# API Specification

## Endpoints
- POST /auth/login
- POST /auth/refresh
- DELETE /auth/logout
`
	require.NoError(t, os.WriteFile(filepath.Join(specDir, "auth.md"), []byte(authSpec), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(specDir, "api.md"), []byte(apiSpec), 0644))

	// Step 4: Verify spec resolver can find them
	resolver := specs.NewResolver(cfg)

	allSpecs, err := resolver.ListSpecs()
	require.NoError(t, err)
	assert.Len(t, allSpecs, 2)

	resolved := resolver.ResolveSpec("auth.md")
	assert.True(t, resolved.Found)
	assert.Contains(t, resolved.Content, "JWT-based authentication")

	content, err := resolver.ReadSpec("api.md")
	require.NoError(t, err)
	assert.Contains(t, content, "POST /auth/login")

	// Missing spec should return not found
	missing := resolver.ResolveSpec("nonexistent.md")
	assert.False(t, missing.Found)

	// Step 5: Create a task with spec refs
	tm := tasks.NewTaskManager(".agentic/tasks")
	task, err := tm.CreateTask("Implement Authentication")
	require.NoError(t, err)

	// Update task with spec refs and inputs
	backlog, err := tm.LoadTasks("backlog")
	require.NoError(t, err)

	backlog.Tasks[0].SpecRefs = []string{"auth.md", "api.md"}
	backlog.Tasks[0].Scope = []string{"internal/auth"}
	backlog.Tasks[0].Acceptance = []string{"JWT tokens work", "All endpoints implemented"}
	err = tm.SaveTasks("backlog", backlog)
	require.NoError(t, err)

	// Step 6: Run readiness checks
	backlog, err = tm.LoadTasks("backlog")
	require.NoError(t, err)
	taskToCheck := backlog.Tasks[0]

	result := tasks.CanClaimTask(&taskToCheck, cfg)
	// Specs should be resolvable, scope dir won't exist (warning only)
	assert.True(t, result.Ready, "Task should be ready since specs exist and scope is warning-only")

	// Verify readiness output format
	formatted := tasks.FormatReadinessResult(result)
	assert.Contains(t, formatted, "READY")
	assert.Contains(t, formatted, "spec-resolvable")

	// Step 7: Claim with readiness checks
	err = tm.ClaimTaskWithConfig(task.ID, "test-agent", cfg)
	require.NoError(t, err)

	inProgress, err := tm.LoadTasks("in-progress")
	require.NoError(t, err)
	assert.Len(t, inProgress.Tasks, 1)
	assert.Equal(t, task.ID, inProgress.Tasks[0].ID)

	// Step 8: Build context bundle — should include resolved specs
	bundle, err := encoding.CreateContextBundle(task.ID, "toon", cfg)
	require.NoError(t, err)
	assert.NotEmpty(t, bundle)

	// The bundle should contain spec content (it's YAML-encoded)
	bundleStr := string(bundle)
	assert.Contains(t, bundleStr, "auth.md")
	assert.Contains(t, bundleStr, "api.md")
	assert.Contains(t, bundleStr, "JWT-based authentication")
}

// TestSpecWorkflow_MultiDir tests spec resolution across multiple directories
func TestSpecWorkflow_MultiDir(t *testing.T) {
	tmpDir := setupIntegrationTest(t)

	err := project.InitProject("MultiDirSpecTest")
	require.NoError(t, err)

	// Create two spec directories
	dir1 := filepath.Join(tmpDir, "specs", "core")
	dir2 := filepath.Join(tmpDir, "specs", "shared")
	require.NoError(t, os.MkdirAll(dir1, 0755))
	require.NoError(t, os.MkdirAll(dir2, 0755))

	require.NoError(t, os.WriteFile(filepath.Join(dir1, "core.md"), []byte("# Core Spec"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir2, "shared.md"), []byte("# Shared Spec"), 0644))
	// Duplicate: dir1 should shadow
	require.NoError(t, os.WriteFile(filepath.Join(dir2, "core.md"), []byte("# Core Spec (shadowed)"), 0644))

	cfg := &models.Config{}
	config.SetDefaults(cfg)
	cfg.Paths.SpecDirs = []string{dir1, dir2}

	resolver := specs.NewResolver(cfg)

	// core.md should come from dir1 (first match wins)
	coreSpec := resolver.ResolveSpec("core.md")
	assert.True(t, coreSpec.Found)
	assert.Equal(t, "# Core Spec", coreSpec.Content)

	// shared.md from dir2
	sharedSpec := resolver.ResolveSpec("shared.md")
	assert.True(t, sharedSpec.Found)
	assert.Equal(t, "# Shared Spec", sharedSpec.Content)

	// ListSpecs should deduplicate
	allSpecs, err := resolver.ListSpecs()
	require.NoError(t, err)
	assert.Len(t, allSpecs, 2, "Should have 2 unique specs (core.md deduplicated)")

	// ResolveAll with mixed results
	results := resolver.ResolveAll([]string{"core.md", "missing.md", "shared.md"})
	assert.Len(t, results, 3)
	assert.True(t, results[0].Found)
	assert.False(t, results[1].Found)
	assert.True(t, results[2].Found)
}

// TestSpecWorkflow_ReadinessBlocking tests that missing inputs block readiness
func TestSpecWorkflow_ReadinessBlocking(t *testing.T) {
	_ = setupIntegrationTest(t)

	err := project.InitProject("ReadinessTest")
	require.NoError(t, err)

	cfg := &models.Config{}
	config.SetDefaults(cfg)

	// Task with missing inputs should NOT be ready
	task := &models.Task{
		ID:     "T-BLOCK",
		Inputs: []string{"/nonexistent/required-file.txt"},
	}

	result := tasks.CanClaimTask(task, cfg)
	assert.False(t, result.Ready, "Task with missing inputs should not be ready")

	// Task with missing specs should NOT be ready
	task2 := &models.Task{
		ID:       "T-BLOCK2",
		SpecRefs: []string{"nonexistent-spec.md"},
	}

	result2 := tasks.CanClaimTask(task2, cfg)
	assert.False(t, result2.Ready, "Task with missing specs should not be ready")

	// Task with only missing scope should still be ready (warning only)
	task3 := &models.Task{
		ID:    "T-WARN",
		Scope: []string{"/nonexistent/dir"},
	}

	result3 := tasks.CanClaimTask(task3, cfg)
	assert.True(t, result3.Ready, "Task with only missing scope should still be ready")
	assert.Len(t, result3.Checks, 1)
	assert.False(t, result3.Checks[0].Passed, "Scope check should fail but not block")
}
