package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

func setupEnsureTestDir(t *testing.T) (string, func()) {
	t.Helper()
	dir := t.TempDir()
	origDir, _ := os.Getwd()

	os.Chdir(dir)
	os.MkdirAll(".agentic/agent-rules", 0755)
	os.WriteFile(".agentic/agent-rules/base.md", []byte("- Base rule 1\n"), 0644)

	return dir, func() { os.Chdir(origDir) }
}

func TestEnsure_GeneratesRulesWhenMissing(t *testing.T) {
	_, cleanup := setupEnsureTestDir(t)
	defer cleanup()

	cfg := &models.Config{}
	result, err := Ensure("claude-code", cfg)
	if err != nil {
		t.Fatalf("Ensure failed: %v", err)
	}

	if !result.RulesGenerated {
		t.Error("expected RulesGenerated=true")
	}
	if result.RulesFile != "CLAUDE.md" {
		t.Errorf("expected RulesFile=CLAUDE.md, got %s", result.RulesFile)
	}

	// Verify file exists
	if _, err := os.Stat("CLAUDE.md"); os.IsNotExist(err) {
		t.Error("CLAUDE.md was not created")
	}
}

func TestEnsure_Idempotent(t *testing.T) {
	_, cleanup := setupEnsureTestDir(t)
	defer cleanup()

	cfg := &models.Config{}

	// First run generates
	result1, err := Ensure("claude-code", cfg)
	if err != nil {
		t.Fatalf("First ensure failed: %v", err)
	}
	if !result1.RulesGenerated {
		t.Error("first run should generate rules")
	}

	// Second run is a no-op
	result2, err := Ensure("claude-code", cfg)
	if err != nil {
		t.Fatalf("Second ensure failed: %v", err)
	}
	if result2.RulesGenerated {
		t.Error("second run should not regenerate")
	}
	if result2.DriftFixed {
		t.Error("second run should not fix drift")
	}
}

func TestEnsure_FixesDrift(t *testing.T) {
	_, cleanup := setupEnsureTestDir(t)
	defer cleanup()

	cfg := &models.Config{}

	// Generate first
	Ensure("claude-code", cfg)

	// Manually modify the file to create drift
	os.WriteFile("CLAUDE.md", []byte("modified content"), 0644)

	// Ensure should fix drift
	result, err := Ensure("claude-code", cfg)
	if err != nil {
		t.Fatalf("Ensure failed: %v", err)
	}
	if !result.DriftFixed {
		t.Error("expected DriftFixed=true")
	}
}

func TestEnsure_InstallsConfiguredPacks(t *testing.T) {
	_, cleanup := setupEnsureTestDir(t)
	defer cleanup()

	cfg := &models.Config{
		Agents: models.AgentsConfig{
			Overrides: []models.AgentConfig{
				{
					Name:       "claude-code",
					SkillPacks: []string{"tdd"},
				},
			},
		},
	}

	result, err := Ensure("claude-code", cfg)
	if err != nil {
		t.Fatalf("Ensure failed: %v", err)
	}

	if len(result.PacksInstalled) == 0 {
		t.Error("expected tdd pack to be installed")
	}

	// Verify TDD skill file exists
	skillPath := filepath.Join(".claude", "skills", "tdd", "SKILL.md")
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		t.Errorf("expected %s to exist", skillPath)
	}
}

func TestEnsure_UnknownTool(t *testing.T) {
	_, cleanup := setupEnsureTestDir(t)
	defer cleanup()

	cfg := &models.Config{}
	result, err := Ensure("unknown-tool", cfg)
	if err != nil {
		t.Fatalf("Ensure should not fail for unknown tool: %v", err)
	}

	if len(result.Warnings) == 0 {
		t.Error("expected warnings for unknown tool")
	}
}

func TestEnsure_CursorTool(t *testing.T) {
	_, cleanup := setupEnsureTestDir(t)
	defer cleanup()

	cfg := &models.Config{}
	result, err := Ensure("cursor", cfg)
	if err != nil {
		t.Fatalf("Ensure failed: %v", err)
	}

	if !result.RulesGenerated {
		t.Error("expected RulesGenerated=true for cursor")
	}

	if _, err := os.Stat(".cursor/rules/agnostic-agent.mdc"); os.IsNotExist(err) {
		t.Error("cursor rules file was not created")
	}
}

func TestEnsure_GeneratesPrdAndRalphForAllTools(t *testing.T) {
	tools := []struct {
		name     string
		skillDir string
	}{
		{"claude-code", ".claude/skills"},
		{"cursor", ".cursor/skills"},
		{"gemini", ".gemini/skills"},
		{"windsurf", ".windsurf/skills"},
		{"codex", ".codex/skills"},
		{"copilot", ".github/skills"},
		{"opencode", ".opencode/skills"},
		{"antigravity", ".agent/skills"},
	}

	for _, tc := range tools {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			origDir, _ := os.Getwd()
			os.Chdir(dir)
			defer os.Chdir(origDir)

			os.MkdirAll(".agentic/agent-rules", 0755)
			os.WriteFile(".agentic/agent-rules/base.md", []byte("- Base rule\n"), 0644)

			// Override ToolSkillDir for this tool to use temp dir
			origSkillDir := ToolSkillDir[tc.name]
			ToolSkillDir[tc.name] = tc.skillDir
			defer func() { ToolSkillDir[tc.name] = origSkillDir }()

			cfg := &models.Config{}
			cfg.Paths.PRDOutputPath = ".agentic/tasks/"

			_, err := Ensure(tc.name, cfg)
			if err != nil {
				t.Fatalf("Ensure failed for %s: %v", tc.name, err)
			}

			// Verify prd.md exists in the tool's skill dir
			prdPath := filepath.Join(tc.skillDir, "prd.md")
			if _, err := os.Stat(prdPath); os.IsNotExist(err) {
				t.Errorf("expected %s to exist for %s", prdPath, tc.name)
			}

			// Verify ralph-converter.md exists
			ralphPath := filepath.Join(tc.skillDir, "ralph-converter.md")
			if _, err := os.Stat(ralphPath); os.IsNotExist(err) {
				t.Errorf("expected %s to exist for %s", ralphPath, tc.name)
			}

			// Verify template variable was rendered
			content, _ := os.ReadFile(prdPath)
			if !strings.Contains(string(content), ".agentic/tasks/") {
				t.Errorf("prd.md should contain rendered PRDOutputPath")
			}
			if strings.Contains(string(content), "{{.PRDOutputPath}}") {
				t.Errorf("prd.md should not contain unrendered template variable")
			}
		})
	}
}

func TestEnsure_GeminiAlsoGetsCommands(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	os.MkdirAll(".agentic/agent-rules", 0755)
	os.WriteFile(".agentic/agent-rules/base.md", []byte("- Base rule\n"), 0644)

	cfg := &models.Config{}
	cfg.Paths.PRDOutputPath = ".agentic/tasks/"

	_, err := Ensure("gemini", cfg)
	if err != nil {
		t.Fatalf("Ensure failed: %v", err)
	}

	// Gemini should have skill files AND command TOML files
	for _, path := range []string{
		".gemini/skills/prd.md",
		".gemini/skills/ralph-converter.md",
		".gemini/commands/prd/gen.toml",
		".gemini/commands/ralph/convert.toml",
	} {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected %s to exist for gemini", path)
		}
	}
}

func TestGenerateToolSkills_UnsupportedTool(t *testing.T) {
	cfg := &models.Config{}
	cfg.Paths.PRDOutputPath = ".agentic/tasks/"
	gen := NewGeneratorWithConfig(cfg)

	err := gen.GenerateToolSkills("nonexistent-tool")
	if err == nil {
		t.Fatal("expected error for unsupported tool")
	}
	if !strings.Contains(err.Error(), "unsupported tool") {
		t.Errorf("expected 'unsupported tool' in error, got: %s", err.Error())
	}
}

func TestGenerateToolSkills_RequiresConfig(t *testing.T) {
	gen := NewGenerator() // no config
	err := gen.GenerateToolSkills("claude-code")
	if err == nil {
		t.Fatal("expected error without config")
	}
	if !strings.Contains(err.Error(), "config required") {
		t.Errorf("expected 'config required' in error, got: %s", err.Error())
	}
}

func TestFormatEnsureResult_UpToDate(t *testing.T) {
	r := &EnsureResult{Agent: "claude-code"}
	out := FormatEnsureResult(r)
	if out != "Already up to date.\n" {
		t.Errorf("unexpected output: %s", out)
	}
}
