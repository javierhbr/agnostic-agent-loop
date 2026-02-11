package functional

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/javierbenavides/agentic-agent/internal/project"
	"github.com/javierbenavides/agentic-agent/internal/skills"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// TestSkillRefsYAMLRoundTrip tests that skill_refs survive YAML marshal/unmarshal
func TestSkillRefsYAMLRoundTrip(t *testing.T) {
	task := models.Task{
		ID:        "T-001",
		Title:     "Test task",
		Status:    models.StatusPending,
		SpecRefs:  []string{"auth-spec.md"},
		SkillRefs: []string{"tdd", "code-simplification"},
		Scope:     []string{"internal/auth"},
	}

	// Marshal
	data, err := yaml.Marshal(task)
	require.NoError(t, err)
	assert.Contains(t, string(data), "skill_refs:")
	assert.Contains(t, string(data), "- tdd")
	assert.Contains(t, string(data), "- code-simplification")

	// Unmarshal
	var parsed models.Task
	err = yaml.Unmarshal(data, &parsed)
	require.NoError(t, err)
	assert.Equal(t, []string{"tdd", "code-simplification"}, parsed.SkillRefs)
	assert.Equal(t, []string{"auth-spec.md"}, parsed.SpecRefs)
}

// TestSkillRefsYAMLOmitEmpty tests that skill_refs is omitted when empty
func TestSkillRefsYAMLOmitEmpty(t *testing.T) {
	task := models.Task{
		ID:     "T-002",
		Title:  "Task without skills",
		Status: models.StatusPending,
	}

	data, err := yaml.Marshal(task)
	require.NoError(t, err)
	assert.NotContains(t, string(data), "skill_refs")
}

// TestSkillRefsPersistThroughTaskLifecycle tests skill_refs persist through claim and completion
func TestSkillRefsPersistThroughTaskLifecycle(t *testing.T) {
	SetupTestProject(t)

	err := project.InitProject("SkillRefsLifecycle")
	require.NoError(t, err)

	tm := tasks.NewTaskManager(".agentic/tasks")

	// Create task and add skill_refs
	task, err := tm.CreateTask("Task with skills")
	require.NoError(t, err)

	backlog, err := tm.LoadTasks("backlog")
	require.NoError(t, err)
	backlog.Tasks[0].SkillRefs = []string{"tdd", "api-docs"}
	err = tm.SaveTasks("backlog", backlog)
	require.NoError(t, err)

	// Claim
	err = tm.ClaimTask(task.ID, "agent")
	require.NoError(t, err)

	inProgress, err := tm.LoadTasks("in-progress")
	require.NoError(t, err)
	require.Len(t, inProgress.Tasks, 1)
	assert.Equal(t, []string{"tdd", "api-docs"}, inProgress.Tasks[0].SkillRefs, "skill_refs should survive claim")

	// Complete
	err = tm.MoveTask(task.ID, "in-progress", "done", models.StatusDone)
	require.NoError(t, err)

	done, err := tm.LoadTasks("done")
	require.NoError(t, err)
	require.Len(t, done.Tasks, 1)
	assert.Equal(t, []string{"tdd", "api-docs"}, done.Tasks[0].SkillRefs, "skill_refs should survive completion")
}

// TestResolveSkillRefsAllPacks tests that all registered packs resolve via embedded FS
func TestResolveSkillRefsAllPacks(t *testing.T) {
	packs := []string{"tdd", "api-docs", "code-simplification", "dev-plans", "diataxis", "extract-wisdom"}

	for _, pack := range packs {
		t.Run(pack, func(t *testing.T) {
			resolved := skills.ResolveSkillRefs([]string{pack}, "")
			require.Len(t, resolved, 1)
			assert.True(t, resolved[0].Found, "pack %s should resolve via embedded FS", pack)
			assert.NotEmpty(t, resolved[0].Content, "pack %s should have content", pack)
		})
	}
}

// TestResolveSkillRefsFromInstalledDir tests priority over embedded content
func TestResolveSkillRefsFromInstalledDir(t *testing.T) {
	tmpDir := SetupTestProject(t)

	// Install a custom skill for claude-code
	skillDir := filepath.Join(tmpDir, ".claude", "skills", "tdd")
	os.MkdirAll(skillDir, 0755)
	customContent := "# Custom Installed TDD\nOverrides embedded."
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(customContent), 0644)

	origDir := skills.ToolSkillDir["claude-code"]
	skills.ToolSkillDir["claude-code"] = filepath.Join(tmpDir, ".claude", "skills")
	defer func() { skills.ToolSkillDir["claude-code"] = origDir }()

	resolved := skills.ResolveSkillRefs([]string{"tdd"}, "claude-code")
	require.Len(t, resolved, 1)
	assert.True(t, resolved[0].Found)
	assert.Equal(t, customContent, resolved[0].Content)
}

// TestResolveSkillRefsMixed tests a mix of found and not-found refs
func TestResolveSkillRefsMixed(t *testing.T) {
	resolved := skills.ResolveSkillRefs([]string{"tdd", "nonexistent", "code-simplification"}, "")
	require.Len(t, resolved, 3)

	assert.True(t, resolved[0].Found)
	assert.Equal(t, "tdd", resolved[0].Ref)

	assert.False(t, resolved[1].Found)
	assert.Equal(t, "nonexistent", resolved[1].Ref)
	assert.NotEmpty(t, resolved[1].Error)

	assert.True(t, resolved[2].Found)
	assert.Equal(t, "code-simplification", resolved[2].Ref)
}
