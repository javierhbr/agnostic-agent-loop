package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstaller_Install(t *testing.T) {
	tmpDir := t.TempDir()

	// Override ToolSkillDir for test
	origDir := ToolSkillDir["claude-code"]
	ToolSkillDir["claude-code"] = filepath.Join(tmpDir, ".claude", "skills")
	defer func() { ToolSkillDir["claude-code"] = origDir }()

	installer := NewInstaller()
	result, err := installer.Install("tdd", "claude-code", false)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	if result.PackName != "tdd" {
		t.Errorf("expected PackName 'tdd', got %q", result.PackName)
	}
	if result.Tool != "claude-code" {
		t.Errorf("expected Tool 'claude-code', got %q", result.Tool)
	}
	if len(result.FilesWritten) != 3 {
		t.Errorf("expected 3 files written, got %d", len(result.FilesWritten))
	}

	// Verify files exist and have content
	for _, f := range result.FilesWritten {
		info, err := os.Stat(f)
		if err != nil {
			t.Errorf("file %s does not exist: %v", f, err)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("file %s is empty", f)
		}
	}
}

func TestInstaller_Install_UnknownPack(t *testing.T) {
	installer := NewInstaller()
	_, err := installer.Install("nonexistent", "claude-code", false)
	if err == nil {
		t.Fatal("expected error for unknown pack")
	}
}

func TestInstaller_Install_UnknownTool(t *testing.T) {
	installer := NewInstaller()
	_, err := installer.Install("tdd", "unknown-tool", false)
	if err == nil {
		t.Fatal("expected error for unknown tool")
	}
}

func TestInstaller_IsInstalled(t *testing.T) {
	tmpDir := t.TempDir()

	origDir := ToolSkillDir["claude-code"]
	skillDir := filepath.Join(tmpDir, ".claude", "skills")
	ToolSkillDir["claude-code"] = skillDir
	defer func() { ToolSkillDir["claude-code"] = origDir }()

	installer := NewInstaller()

	// Not installed yet
	if installer.IsInstalled("tdd", "claude-code") {
		t.Error("expected IsInstalled to return false before install")
	}

	// Install
	_, err := installer.Install("tdd", "claude-code", false)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	// Now installed
	if !installer.IsInstalled("tdd", "claude-code") {
		t.Error("expected IsInstalled to return true after install")
	}
}

func TestInstaller_IsInstalledAnywhere(t *testing.T) {
	tmpDir := t.TempDir()

	origDir := ToolSkillDir["gemini"]
	skillDir := filepath.Join(tmpDir, ".gemini", "skills")
	ToolSkillDir["gemini"] = skillDir
	defer func() { ToolSkillDir["gemini"] = origDir }()

	installer := NewInstaller()

	// Install for gemini
	_, err := installer.Install("tdd", "gemini", false)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	tool := installer.IsInstalledAnywhere("tdd")
	if tool != "gemini" {
		t.Errorf("expected IsInstalledAnywhere to find 'gemini', got %q", tool)
	}
}

func TestInstaller_Install_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()

	origDir := ToolSkillDir["claude-code"]
	ToolSkillDir["claude-code"] = filepath.Join(tmpDir, ".claude", "skills")
	defer func() { ToolSkillDir["claude-code"] = origDir }()

	installer := NewInstaller()

	// Install twice — should not error
	_, err := installer.Install("tdd", "claude-code", false)
	if err != nil {
		t.Fatalf("first install failed: %v", err)
	}

	_, err = installer.Install("tdd", "claude-code", false)
	if err != nil {
		t.Fatalf("second install failed: %v", err)
	}
}

func TestInstaller_InstallMulti(t *testing.T) {
	tmpDir := t.TempDir()

	// Override two tool dirs
	origClaude := ToolSkillDir["claude-code"]
	origCursor := ToolSkillDir["cursor"]
	ToolSkillDir["claude-code"] = filepath.Join(tmpDir, ".claude", "skills")
	ToolSkillDir["cursor"] = filepath.Join(tmpDir, ".cursor", "skills")
	defer func() {
		ToolSkillDir["claude-code"] = origClaude
		ToolSkillDir["cursor"] = origCursor
	}()

	installer := NewInstaller()
	results, err := installer.InstallMulti("tdd", []string{"claude-code", "cursor"}, false)
	if err != nil {
		t.Fatalf("InstallMulti failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// Verify each result
	for _, r := range results {
		if r.PackName != "tdd" {
			t.Errorf("expected PackName 'tdd', got %q", r.PackName)
		}
		if len(r.FilesWritten) == 0 {
			t.Errorf("expected files written for tool %s", r.Tool)
		}
		for _, f := range r.FilesWritten {
			if _, err := os.Stat(f); err != nil {
				t.Errorf("file %s does not exist: %v", f, err)
			}
		}
	}

	// Verify both tools report installed
	if !installer.IsInstalled("tdd", "claude-code") {
		t.Error("expected tdd installed for claude-code")
	}
	if !installer.IsInstalled("tdd", "cursor") {
		t.Error("expected tdd installed for cursor")
	}
}

func TestInstaller_InstallMulti_PartialFailure(t *testing.T) {
	tmpDir := t.TempDir()

	origClaude := ToolSkillDir["claude-code"]
	ToolSkillDir["claude-code"] = filepath.Join(tmpDir, ".claude", "skills")
	defer func() { ToolSkillDir["claude-code"] = origClaude }()

	installer := NewInstaller()
	results, err := installer.InstallMulti("tdd", []string{"claude-code", "unknown-tool"}, false)
	if err == nil {
		t.Fatal("expected error for unknown tool")
	}

	// Should have partial results from the first successful install
	if len(results) != 1 {
		t.Errorf("expected 1 partial result, got %d", len(results))
	}
}

func TestInstaller_ListPacks(t *testing.T) {
	installer := NewInstaller()
	packs := installer.ListPacks()
	if len(packs) == 0 {
		t.Fatal("expected at least one pack")
	}
}

func TestInstaller_Install_AgentHelperPack(t *testing.T) {
	tmpDir := t.TempDir()

	// Override both ToolSkillDir and ToolAgentDir for test
	origSkillDir := ToolSkillDir["claude-code"]
	origAgentDir := ToolAgentDir["claude-code"]
	ToolSkillDir["claude-code"] = filepath.Join(tmpDir, ".claude", "skills")
	ToolAgentDir["claude-code"] = filepath.Join(tmpDir, ".claude", "agents")
	defer func() {
		ToolSkillDir["claude-code"] = origSkillDir
		ToolAgentDir["claude-code"] = origAgentDir
	}()

	installer := NewInstaller()
	result, err := installer.Install("agentic-helper", "claude-code", false)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	if result.PackName != "agentic-helper" {
		t.Errorf("expected PackName 'agentic-helper', got %q", result.PackName)
	}
	if result.Tool != "claude-code" {
		t.Errorf("expected Tool 'claude-code', got %q", result.Tool)
	}
	if len(result.FilesWritten) != 2 {
		t.Errorf("expected 2 files written for agentic-helper, got %d", len(result.FilesWritten))
	}

	// Verify AGENT.md is in ToolAgentDir
	agentFile := filepath.Join(ToolAgentDir["claude-code"], "agentic-helper.md")
	if _, err := os.Stat(agentFile); os.IsNotExist(err) {
		t.Errorf("AGENT.md not found at expected location: %s", agentFile)
	} else if info, _ := os.Stat(agentFile); info.Size() == 0 {
		t.Errorf("AGENT.md is empty")
	}

	// Verify SKILL.md is in ToolSkillDir
	skillFile := filepath.Join(ToolSkillDir["claude-code"], "agentic-helper", "SKILL.md")
	if _, err := os.Stat(skillFile); os.IsNotExist(err) {
		t.Errorf("SKILL.md not found at expected location: %s", skillFile)
	} else if info, _ := os.Stat(skillFile); info.Size() == 0 {
		t.Errorf("SKILL.md is empty")
	}
}

func TestInstaller_Install_AgentFile_FallsBackForNonAgentTools(t *testing.T) {
	tmpDir := t.TempDir()

	// Override ToolSkillDir for cursor (no ToolAgentDir entry)
	origSkillDir := ToolSkillDir["cursor"]
	ToolSkillDir["cursor"] = filepath.Join(tmpDir, ".cursor", "skills")
	defer func() { ToolSkillDir["cursor"] = origSkillDir }()

	installer := NewInstaller()
	result, err := installer.Install("agentic-helper", "cursor", false)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	if len(result.FilesWritten) != 2 {
		t.Errorf("expected 2 files written, got %d", len(result.FilesWritten))
	}

	// For cursor (no agent dir), AGENT.md should fall back to ToolSkillDir
	skillDir := ToolSkillDir["cursor"]
	agentFile := filepath.Join(skillDir, "agentic-helper.md")
	if _, err := os.Stat(agentFile); os.IsNotExist(err) {
		t.Errorf("AGENT.md should fall back to ToolSkillDir for cursor, but not found at %s", agentFile)
	}
}

func TestInstaller_IsInstalled_AgentPack(t *testing.T) {
	tmpDir := t.TempDir()

	origSkillDir := ToolSkillDir["claude-code"]
	origAgentDir := ToolAgentDir["claude-code"]
	ToolSkillDir["claude-code"] = filepath.Join(tmpDir, ".claude", "skills")
	ToolAgentDir["claude-code"] = filepath.Join(tmpDir, ".claude", "agents")
	defer func() {
		ToolSkillDir["claude-code"] = origSkillDir
		ToolAgentDir["claude-code"] = origAgentDir
	}()

	installer := NewInstaller()

	// Not installed yet
	if installer.IsInstalled("agentic-helper", "claude-code") {
		t.Error("expected IsInstalled to return false before install")
	}

	// Install
	_, err := installer.Install("agentic-helper", "claude-code", false)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	// Now installed
	if !installer.IsInstalled("agentic-helper", "claude-code") {
		t.Error("expected IsInstalled to return true after install")
	}
}

func TestResolveAgentOutputDir_ClaudeCode(t *testing.T) {
	agentDir, err := resolveAgentOutputDir("claude-code", false)
	if err != nil {
		t.Fatalf("resolveAgentOutputDir failed: %v", err)
	}

	expected := ".claude/agents"
	if agentDir != expected {
		t.Errorf("expected %q, got %q", expected, agentDir)
	}
}

func TestResolveAgentOutputDir_UnsupportedTool(t *testing.T) {
	_, err := resolveAgentOutputDir("cursor", false)
	if err == nil {
		t.Fatal("expected error for unsupported tool cursor")
	}
}

func TestResolveAgentOutputDir_Global(t *testing.T) {
	agentDir, err := resolveAgentOutputDir("claude-code", true)
	if err != nil {
		t.Fatalf("resolveAgentOutputDir failed: %v", err)
	}

	// Should have expanded ~ to home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	expected := filepath.Join(homeDir, ".claude", "agents")
	if agentDir != expected {
		t.Errorf("expected %q, got %q", expected, agentDir)
	}
}
