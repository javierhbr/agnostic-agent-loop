package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveSkillRefs_EmptyRefs(t *testing.T) {
	result := ResolveSkillRefs(nil, "claude-code")
	if result != nil {
		t.Errorf("expected nil for empty refs, got %v", result)
	}

	result = ResolveSkillRefs([]string{}, "claude-code")
	if result != nil {
		t.Errorf("expected nil for empty slice, got %v", result)
	}
}

func TestResolveSkillRefs_EmbeddedFallback(t *testing.T) {
	// With no installed skills, should fall back to embedded packs
	result := ResolveSkillRefs([]string{"tdd"}, "nonexistent-agent")

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if !result[0].Found {
		t.Errorf("expected tdd to be found via embedded FS, got error: %s", result[0].Error)
	}
	if result[0].Content == "" {
		t.Error("expected non-empty content for tdd skill")
	}
	if result[0].Ref != "tdd" {
		t.Errorf("expected ref 'tdd', got %q", result[0].Ref)
	}
}

func TestResolveSkillRefs_UnknownPack(t *testing.T) {
	result := ResolveSkillRefs([]string{"nonexistent-pack"}, "claude-code")

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Found {
		t.Error("expected nonexistent-pack to not be found")
	}
	if result[0].Error == "" {
		t.Error("expected non-empty error for unresolved skill")
	}
}

func TestResolveSkillRefs_MultipleRefs(t *testing.T) {
	result := ResolveSkillRefs([]string{"tdd", "code-simplification", "nonexistent"}, "")

	if len(result) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result))
	}

	// tdd should be found (embedded)
	if !result[0].Found {
		t.Errorf("expected tdd to be found, got error: %s", result[0].Error)
	}

	// code-simplification should be found (embedded)
	if !result[1].Found {
		t.Errorf("expected code-simplification to be found, got error: %s", result[1].Error)
	}

	// nonexistent should not be found
	if result[2].Found {
		t.Error("expected nonexistent to not be found")
	}
}

func TestResolveSkillRefs_InstalledSkillTakesPriority(t *testing.T) {
	// Set up a temp directory with installed skill
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	// Create a fake installed skill for a test tool
	skillDir := filepath.Join(tmpDir, ".claude", "skills", "tdd")
	os.MkdirAll(skillDir, 0755)
	customContent := "# Custom Installed TDD Skill\nThis is a custom version."
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(customContent), 0644)

	// Temporarily override ToolSkillDir for testing
	origToolDir := ToolSkillDir["claude-code"]
	ToolSkillDir["claude-code"] = filepath.Join(tmpDir, ".claude", "skills")
	defer func() { ToolSkillDir["claude-code"] = origToolDir }()

	result := ResolveSkillRefs([]string{"tdd"}, "claude-code")

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if !result[0].Found {
		t.Fatalf("expected tdd to be found, got error: %s", result[0].Error)
	}
	// Should pick up the installed version, not embedded
	if result[0].Content != customContent {
		t.Errorf("expected installed content, got embedded content instead")
	}
}

func TestResolveSkillRefs_EmbeddedWithResources(t *testing.T) {
	// diataxis has resources/ subdirectory
	result := ResolveSkillRefs([]string{"diataxis"}, "")

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if !result[0].Found {
		t.Fatalf("expected diataxis to be found, got error: %s", result[0].Error)
	}
	// Should contain content from SKILL.md + resources (joined by separator)
	if !containsSubstring(result[0].Content, "---") {
		t.Error("expected diataxis content to include resource separator")
	}
}

func TestReadEmbeddedSkill_KnownPacks(t *testing.T) {
	packs := []string{"tdd", "api-docs", "code-simplification", "dev-plans", "diataxis", "extract-wisdom"}
	for _, pack := range packs {
		content, err := readEmbeddedSkill(pack)
		if err != nil {
			t.Errorf("failed to read embedded skill %q: %v", pack, err)
			continue
		}
		if content == "" {
			t.Errorf("embedded skill %q returned empty content", pack)
		}
	}
}

func TestReadEmbeddedSkill_Unknown(t *testing.T) {
	_, err := readEmbeddedSkill("does-not-exist")
	if err == nil {
		t.Error("expected error for unknown embedded skill")
	}
}

func containsSubstring(s, sub string) bool {
	return len(s) > 0 && len(sub) > 0 && contains(s, sub)
}

func contains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
