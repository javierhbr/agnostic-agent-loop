package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstaller_Install_CreatesCanonicalAndSymlink(t *testing.T) {
	tmpDir := t.TempDir()

	origDir := ToolSkillDir["claude-code"]
	toolSkillDir := filepath.Join(tmpDir, ".claude", "skills")
	ToolSkillDir["claude-code"] = toolSkillDir
	defer func() { ToolSkillDir["claude-code"] = origDir }()

	// Set canonical dir to temp
	canonicalDir := filepath.Join(tmpDir, ".agentic", "skills")

	installer := NewInstallerWithCanonicalDir(canonicalDir)
	result, err := installer.Install("tdd", "claude-code", false)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	// Canonical file should exist as a real file
	canonicalPath := filepath.Join(canonicalDir, "tdd", "SKILL.md")
	info, err := os.Lstat(canonicalPath)
	if err != nil {
		t.Fatalf("canonical file missing: %v", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		t.Error("canonical file should be a real file, not a symlink")
	}

	// Tool dir file should be a symlink
	toolPath := filepath.Join(toolSkillDir, "tdd", "SKILL.md")
	info, err = os.Lstat(toolPath)
	if err != nil {
		t.Fatalf("tool dir file missing: %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("tool dir file should be a symlink")
	}

	// Symlink should point to canonical
	target, _ := os.Readlink(toolPath)
	absCanonical, _ := filepath.Abs(canonicalPath)
	if target != absCanonical {
		t.Errorf("symlink target = %q, want %q", target, absCanonical)
	}

	// Result should list files written
	if len(result.FilesWritten) != 3 {
		t.Errorf("expected 3 files written, got %d", len(result.FilesWritten))
	}
}

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

	// Verify both files are in ToolSkillDir (including AGENT.md which is now treated like any other file)
	for _, filePath := range result.FilesWritten {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("file not found at expected location: %s", filePath)
		} else if info, _ := os.Stat(filePath); info.Size() == 0 {
			t.Errorf("file is empty: %s", filePath)
		}
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

	// All files (including AGENT.md) should be installed to ToolSkillDir
	for _, filePath := range result.FilesWritten {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("file not found at expected location: %s", filePath)
		}
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

	canonicalDir := filepath.Join(tmpDir, ".agentic", "skills")
	installer := NewInstallerWithCanonicalDir(canonicalDir)

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
