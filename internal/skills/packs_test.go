package skills

import (
	"testing"
)

func TestNewPackRegistry_HasTDDPack(t *testing.T) {
	r := NewPackRegistry()

	pack, err := r.GetPack("tdd")
	if err != nil {
		t.Fatalf("expected tdd pack to exist: %v", err)
	}

	if pack.Name != "tdd" {
		t.Errorf("expected name 'tdd', got %q", pack.Name)
	}

	if len(pack.Files) != 3 {
		t.Errorf("expected 3 files in tdd pack, got %d", len(pack.Files))
	}
}

func TestNewPackRegistry_GetPackUnknown(t *testing.T) {
	r := NewPackRegistry()

	_, err := r.GetPack("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown pack")
	}
}

func TestNewPackRegistry_GetAll(t *testing.T) {
	r := NewPackRegistry()

	packs := r.GetAll()
	if len(packs) == 0 {
		t.Fatal("expected at least one pack")
	}

	found := false
	for _, p := range packs {
		if p.Name == "tdd" {
			found = true
		}
	}
	if !found {
		t.Error("expected tdd pack in GetAll results")
	}
}

func TestPacksFS_AllFilesReadable(t *testing.T) {
	r := NewPackRegistry()

	for _, pack := range r.GetAll() {
		for _, f := range pack.Files {
			content, err := packsFS.ReadFile(f.SrcPath)
			if err != nil {
				t.Errorf("pack %q: failed to read embedded file %s: %v", pack.Name, f.SrcPath, err)
			}
			if len(content) == 0 {
				t.Errorf("pack %q: embedded file %s is empty", pack.Name, f.SrcPath)
			}
		}
	}
}

func TestNewPackRegistry_AllPacksRegistered(t *testing.T) {
	r := NewPackRegistry()

	expected := []struct {
		name      string
		fileCount int
	}{
		{"tdd", 3},
		{"api-docs", 1},
		{"code-simplification", 1},
		{"dev-plans", 1},
		{"diataxis", 3},
		{"extract-wisdom", 1},
		{"openspec", 1},
	}

	packs := r.GetAll()
	if len(packs) != len(expected) {
		t.Errorf("expected %d packs, got %d", len(expected), len(packs))
	}

	for _, exp := range expected {
		pack, err := r.GetPack(exp.name)
		if err != nil {
			t.Errorf("expected pack %q to exist: %v", exp.name, err)
			continue
		}
		if len(pack.Files) != exp.fileCount {
			t.Errorf("pack %q: expected %d files, got %d", exp.name, exp.fileCount, len(pack.Files))
		}
	}
}

func TestToolSkillDir_AllToolsPresent(t *testing.T) {
	expected := []string{"claude-code", "cursor", "gemini", "windsurf", "antigravity", "codex", "copilot", "opencode"}
	for _, tool := range expected {
		if _, ok := ToolSkillDir[tool]; !ok {
			t.Errorf("ToolSkillDir missing tool %q", tool)
		}
		if _, ok := ToolGlobalSkillDir[tool]; !ok {
			t.Errorf("ToolGlobalSkillDir missing tool %q", tool)
		}
	}
}

func TestToolSkillDir_CorrectPaths(t *testing.T) {
	cases := map[string]string{
		"claude-code": ".claude/skills",
		"antigravity": ".agent/skills",
		"codex":       ".codex/skills",
	}
	for tool, expected := range cases {
		if got := ToolSkillDir[tool]; got != expected {
			t.Errorf("ToolSkillDir[%q] = %q, want %q", tool, got, expected)
		}
	}
}

func TestSupportedTools(t *testing.T) {
	tools := SupportedTools()
	if len(tools) != len(ToolSkillDir) {
		t.Errorf("SupportedTools returned %d tools, expected %d", len(tools), len(ToolSkillDir))
	}
}
