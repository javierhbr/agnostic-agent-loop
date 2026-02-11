package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/javierbenavides/agentic-agent/internal/config"
	"github.com/javierbenavides/agentic-agent/internal/encoding"
	"github.com/javierbenavides/agentic-agent/internal/project"
	"github.com/javierbenavides/agentic-agent/internal/skills"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSkillRefsWorkflow tests the full skill_refs workflow:
// init → create task with skill_refs → claim → build context bundle → verify targeted skills
func TestSkillRefsWorkflow(t *testing.T) {
	_ = setupIntegrationTest(t)

	// Step 1: Initialize project
	err := project.InitProject("SkillRefsTest")
	require.NoError(t, err)

	cfg := &models.Config{}
	config.SetDefaults(cfg)

	// Step 2: Create a task with skill_refs
	tm := tasks.NewTaskManager(".agentic/tasks")
	task, err := tm.CreateTask("Refactor auth middleware")
	require.NoError(t, err)

	// Update task with skill_refs
	backlog, err := tm.LoadTasks("backlog")
	require.NoError(t, err)

	backlog.Tasks[0].SkillRefs = []string{"code-simplification", "tdd"}
	backlog.Tasks[0].Scope = []string{"internal/auth"}
	err = tm.SaveTasks("backlog", backlog)
	require.NoError(t, err)

	// Step 3: Verify skill_refs persisted
	reloaded, err := tm.LoadTasks("backlog")
	require.NoError(t, err)
	assert.Equal(t, []string{"code-simplification", "tdd"}, reloaded.Tasks[0].SkillRefs)

	// Step 4: Claim the task
	err = tm.ClaimTask(task.ID, "test-agent")
	require.NoError(t, err)

	// Step 5: Verify skill_refs survive claim
	inProgress, err := tm.LoadTasks("in-progress")
	require.NoError(t, err)
	require.Len(t, inProgress.Tasks, 1)
	assert.Equal(t, []string{"code-simplification", "tdd"}, inProgress.Tasks[0].SkillRefs)

	// Step 6: Build context bundle — should use skill_refs for targeted content
	bundle, err := encoding.CreateContextBundle(task.ID, "toon", cfg)
	require.NoError(t, err)
	assert.NotEmpty(t, bundle)
}

// TestSkillRefsWorkflow_NoRefs tests that tasks without skill_refs still work (backwards compatible)
func TestSkillRefsWorkflow_NoRefs(t *testing.T) {
	_ = setupIntegrationTest(t)

	err := project.InitProject("NoSkillRefsTest")
	require.NoError(t, err)

	cfg := &models.Config{}
	config.SetDefaults(cfg)

	tm := tasks.NewTaskManager(".agentic/tasks")
	task, err := tm.CreateTask("Simple task without skills")
	require.NoError(t, err)

	err = tm.ClaimTask(task.ID, "test-agent")
	require.NoError(t, err)

	// Build bundle — should work without skill_refs
	bundle, err := encoding.CreateContextBundle(task.ID, "toon", cfg)
	require.NoError(t, err)
	assert.NotEmpty(t, bundle)
}

// TestSkillRefsWorkflow_Resolution tests skill ref resolution across tiers
func TestSkillRefsWorkflow_Resolution(t *testing.T) {
	_ = setupIntegrationTest(t)

	// Test embedded fallback (no installed skills, just embedded packs)
	resolved := skills.ResolveSkillRefs([]string{"tdd", "code-simplification"}, "nonexistent-agent")
	require.Len(t, resolved, 2)

	// Both should resolve via embedded FS
	assert.True(t, resolved[0].Found, "tdd should resolve from embedded FS")
	assert.True(t, resolved[1].Found, "code-simplification should resolve from embedded FS")

	// Content should be non-empty
	assert.NotEmpty(t, resolved[0].Content)
	assert.NotEmpty(t, resolved[1].Content)

	// Unknown pack should fail
	unknown := skills.ResolveSkillRefs([]string{"nonexistent-pack"}, "")
	require.Len(t, unknown, 1)
	assert.False(t, unknown[0].Found)
	assert.NotEmpty(t, unknown[0].Error)
}

// TestSkillRefsWorkflow_InstalledPriority tests that installed skills take priority over embedded
func TestSkillRefsWorkflow_InstalledPriority(t *testing.T) {
	tmpDir := setupIntegrationTest(t)

	// Install a custom version of tdd skill for claude-code
	skillDir := filepath.Join(tmpDir, ".claude", "skills", "tdd")
	require.NoError(t, os.MkdirAll(skillDir, 0755))
	customContent := "# Custom TDD Skill\nThis is a custom installed version for testing."
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(customContent), 0644))

	// Override tool skill dir for this test
	origDir := skills.ToolSkillDir["claude-code"]
	skills.ToolSkillDir["claude-code"] = filepath.Join(tmpDir, ".claude", "skills")
	defer func() { skills.ToolSkillDir["claude-code"] = origDir }()

	resolved := skills.ResolveSkillRefs([]string{"tdd"}, "claude-code")
	require.Len(t, resolved, 1)
	assert.True(t, resolved[0].Found)
	assert.Equal(t, customContent, resolved[0].Content, "Should use installed version, not embedded")
}
