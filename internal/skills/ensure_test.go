package skills

import (
	"os"
	"path/filepath"
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

func TestFormatEnsureResult_UpToDate(t *testing.T) {
	r := &EnsureResult{Agent: "claude-code"}
	out := FormatEnsureResult(r)
	if out != "Already up to date.\n" {
		t.Errorf("unexpected output: %s", out)
	}
}
