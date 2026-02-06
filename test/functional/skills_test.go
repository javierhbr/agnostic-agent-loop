package functional

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/javierbenavides/agentic-agent/internal/project"
	"github.com/javierbenavides/agentic-agent/internal/skills"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSkillsGenerateGemini tests generating the Gemini rules file (.gemini/GEMINI.md).
func TestSkillsGenerateGemini(t *testing.T) {
	tmpDir := SetupTestProject(t)

	// Initialize project (creates .agentic/agent-rules/base.md)
	err := project.InitProject("GeminiTestProject")
	require.NoError(t, err)

	gen := skills.NewGenerator()
	err = gen.Generate("gemini")
	require.NoError(t, err, "Gemini skill generation should succeed")

	// Verify output file exists
	geminiPath := filepath.Join(tmpDir, ".gemini/GEMINI.md")
	assert.FileExists(t, geminiPath, ".gemini/GEMINI.md should be created")

	// Verify content
	content, err := os.ReadFile(geminiPath)
	require.NoError(t, err)
	contentStr := string(content)

	assert.Contains(t, contentStr, "# GEMINI.md - Agnostic Agent Rules")
	assert.Contains(t, contentStr, "## Base Rules")
	assert.Contains(t, contentStr, "## Gemini-Specific Rules")
	assert.Contains(t, contentStr, "agentic-agent task claim")
	assert.Contains(t, contentStr, "agentic-agent context generate")
	assert.Contains(t, contentStr, "agentic-agent task complete")
}

// TestSkillsGenerateGeminiSlashCommands tests generating Gemini slash command TOML files.
func TestSkillsGenerateGeminiSlashCommands(t *testing.T) {
	tmpDir := SetupTestProject(t)

	// Initialize project
	err := project.InitProject("GeminiCommandsProject")
	require.NoError(t, err)

	// Create generator with config
	cfg := &models.Config{}
	cfg.Paths.PRDOutputPath = ".agentic/tasks/"
	gen := skills.NewGeneratorWithConfig(cfg)

	err = gen.GenerateGeminiSkills()
	require.NoError(t, err, "Gemini slash command generation should succeed")

	// Verify PRD command file
	prdPath := filepath.Join(tmpDir, ".gemini/commands/prd/gen.toml")
	assert.FileExists(t, prdPath, "PRD command file should be created")

	prdContent, err := os.ReadFile(prdPath)
	require.NoError(t, err)
	prdStr := string(prdContent)

	assert.Contains(t, prdStr, `description = "Generate a Product Requirements Document`)
	assert.Contains(t, prdStr, ".agentic/tasks/prd-[feature-name].md")
	assert.Contains(t, prdStr, "$PROMPT")

	// Verify Ralph converter command file
	ralphPath := filepath.Join(tmpDir, ".gemini/commands/ralph/convert.toml")
	assert.FileExists(t, ralphPath, "Ralph command file should be created")

	ralphContent, err := os.ReadFile(ralphPath)
	require.NoError(t, err)
	ralphStr := string(ralphContent)

	assert.Contains(t, ralphStr, `description = "Convert a PRD to YAML task format`)
	assert.Contains(t, ralphStr, ".agentic/tasks/")
	assert.Contains(t, ralphStr, "backlog.yaml")
}

// TestSkillsGenerateGeminiSlashCommandsCustomPath tests Gemini commands with custom PRD path.
func TestSkillsGenerateGeminiSlashCommandsCustomPath(t *testing.T) {
	SetupTestProject(t)

	// Initialize project
	err := project.InitProject("CustomPathProject")
	require.NoError(t, err)

	// Use custom PRD output path
	cfg := &models.Config{}
	cfg.Paths.PRDOutputPath = "docs/prds/"
	gen := skills.NewGeneratorWithConfig(cfg)

	err = gen.GenerateGeminiSkills()
	require.NoError(t, err)

	// Verify custom path appears in generated files
	prdContent, err := os.ReadFile(".gemini/commands/prd/gen.toml")
	require.NoError(t, err)
	assert.Contains(t, string(prdContent), "docs/prds/prd-[feature-name].md")

	ralphContent, err := os.ReadFile(".gemini/commands/ralph/convert.toml")
	require.NoError(t, err)
	assert.Contains(t, string(ralphContent), "docs/prds/")
}

// TestSkillsGenerateGeminiWithoutConfig tests that Gemini slash commands fail without config.
func TestSkillsGenerateGeminiWithoutConfig(t *testing.T) {
	SetupTestProject(t)

	err := project.InitProject("NoConfigProject")
	require.NoError(t, err)

	gen := skills.NewGenerator() // No config
	err = gen.GenerateGeminiSkills()
	assert.Error(t, err, "Should fail without config")
	assert.Contains(t, err.Error(), "config required")
}

// TestSkillsRegistryIncludesGemini verifies Gemini is registered in the skill registry.
func TestSkillsRegistryIncludesGemini(t *testing.T) {
	registry := skills.NewSkillRegistry()

	// Verify Gemini is registered
	skill, err := registry.GetSkill("gemini")
	require.NoError(t, err, "Gemini should be a registered skill")
	assert.Equal(t, "gemini", skill.ToolName)
	assert.Equal(t, "templates/gemini.tmpl", skill.TemplatePath)
	assert.Equal(t, ".gemini/GEMINI.md", skill.OutputFile)

	// Verify all three tools are registered
	all := registry.GetAll()
	toolNames := make([]string, 0, len(all))
	for _, s := range all {
		toolNames = append(toolNames, s.ToolName)
	}
	assert.Contains(t, toolNames, "claude-code")
	assert.Contains(t, toolNames, "cursor")
	assert.Contains(t, toolNames, "gemini")
}

// TestSkillsGenerateAll verifies --all generates Gemini alongside other tools.
func TestSkillsGenerateAll(t *testing.T) {
	tmpDir := SetupTestProject(t)

	err := project.InitProject("AllToolsProject")
	require.NoError(t, err)

	gen := skills.NewGenerator()
	registry := skills.NewSkillRegistry()

	generated := []string{}
	for _, s := range registry.GetAll() {
		if err := gen.Generate(s.ToolName); err == nil {
			generated = append(generated, s.OutputFile)
		}
	}

	// Verify Gemini was generated
	assert.Contains(t, generated, ".gemini/GEMINI.md", "Gemini should be in generated list")

	// Verify file exists
	geminiPath := filepath.Join(tmpDir, ".gemini/GEMINI.md")
	assert.FileExists(t, geminiPath)
}

// TestSkillsDriftCheckGemini tests drift detection for Gemini files.
func TestSkillsDriftCheckGemini(t *testing.T) {
	SetupTestProject(t)

	err := project.InitProject("DriftProject")
	require.NoError(t, err)

	gen := skills.NewGenerator()

	// Generate Gemini skill
	err = gen.Generate("gemini")
	require.NoError(t, err)

	// No drift initially
	drifted, err := gen.CheckDrift()
	require.NoError(t, err)

	geminiDrifted := false
	for _, d := range drifted {
		if strings.Contains(d, "GEMINI.md") {
			geminiDrifted = true
		}
	}
	assert.False(t, geminiDrifted, "No drift expected right after generation")

	// Modify the file to create drift
	err = os.WriteFile(".gemini/GEMINI.md", []byte("modified content"), 0644)
	require.NoError(t, err)

	// Should detect drift now
	drifted, err = gen.CheckDrift()
	require.NoError(t, err)

	geminiDrifted = false
	for _, d := range drifted {
		if strings.Contains(d, "GEMINI.md") {
			geminiDrifted = true
		}
	}
	assert.True(t, geminiDrifted, "Drift should be detected after modification")
}

// TestSkillsDriftCheckGeminiMissing tests drift detection when Gemini file is missing.
func TestSkillsDriftCheckGeminiMissing(t *testing.T) {
	SetupTestProject(t)

	err := project.InitProject("MissingDriftProject")
	require.NoError(t, err)

	gen := skills.NewGenerator()

	// Don't generate Gemini - it should show as missing
	drifted, err := gen.CheckDrift()
	require.NoError(t, err)

	geminiMissing := false
	for _, d := range drifted {
		if strings.Contains(d, "GEMINI.md") && strings.Contains(d, "Missing") {
			geminiMissing = true
		}
	}
	assert.True(t, geminiMissing, "Missing Gemini file should be detected as drift")
}

// TestSkillsGeminiBaseRulesIncluded tests that base rules are injected into Gemini template.
func TestSkillsGeminiBaseRulesIncluded(t *testing.T) {
	SetupTestProject(t)

	err := project.InitProject("BaseRulesProject")
	require.NoError(t, err)

	// Write custom base rules
	err = os.WriteFile(".agentic/agent-rules/base.md", []byte("- Custom rule: always test first\n- Custom rule: keep it simple\n"), 0644)
	require.NoError(t, err)

	gen := skills.NewGenerator()
	err = gen.Generate("gemini")
	require.NoError(t, err)

	content, err := os.ReadFile(".gemini/GEMINI.md")
	require.NoError(t, err)
	contentStr := string(content)

	assert.Contains(t, contentStr, "Custom rule: always test first")
	assert.Contains(t, contentStr, "Custom rule: keep it simple")
}

// TestSkillsGeminiDirectoryCreation tests that .gemini directory is created if missing.
func TestSkillsGeminiDirectoryCreation(t *testing.T) {
	tmpDir := SetupTestProject(t)

	err := project.InitProject("DirCreationProject")
	require.NoError(t, err)

	// Ensure .gemini doesn't exist yet
	geminiDir := filepath.Join(tmpDir, ".gemini")
	_, err = os.Stat(geminiDir)
	assert.True(t, os.IsNotExist(err), ".gemini should not exist before generation")

	gen := skills.NewGenerator()
	err = gen.Generate("gemini")
	require.NoError(t, err)

	// Verify directory was created
	assert.DirExists(t, geminiDir, ".gemini directory should be created")
}

// TestSkillsGeminiCommandDirectoryStructure tests the nested directory structure for commands.
func TestSkillsGeminiCommandDirectoryStructure(t *testing.T) {
	tmpDir := SetupTestProject(t)

	err := project.InitProject("CommandDirProject")
	require.NoError(t, err)

	cfg := &models.Config{}
	cfg.Paths.PRDOutputPath = ".agentic/tasks/"
	gen := skills.NewGeneratorWithConfig(cfg)

	err = gen.GenerateGeminiSkills()
	require.NoError(t, err)

	// Verify the nested directory structure
	assert.DirExists(t, filepath.Join(tmpDir, ".gemini/commands/prd"))
	assert.DirExists(t, filepath.Join(tmpDir, ".gemini/commands/ralph"))
	assert.FileExists(t, filepath.Join(tmpDir, ".gemini/commands/prd/gen.toml"))
	assert.FileExists(t, filepath.Join(tmpDir, ".gemini/commands/ralph/convert.toml"))
}
