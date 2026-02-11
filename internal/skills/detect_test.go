package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectAgent_FlagTakesPriority(t *testing.T) {
	t.Setenv("AGENTIC_AGENT", "gemini")
	agent := DetectAgent("claude-code", t.TempDir())

	if agent.Name != "claude-code" {
		t.Errorf("expected Name=claude-code, got %s", agent.Name)
	}
	if agent.Source != "flag" {
		t.Errorf("expected Source=flag, got %s", agent.Source)
	}
}

func TestDetectAgent_AgenticEnvVar(t *testing.T) {
	t.Setenv("AGENTIC_AGENT", "cursor")
	agent := DetectAgent("", t.TempDir())

	if agent.Name != "cursor" {
		t.Errorf("expected Name=cursor, got %s", agent.Name)
	}
	if agent.Source != "env" {
		t.Errorf("expected Source=env, got %s", agent.Source)
	}
}

func TestDetectAgent_ClaudeEnvVar(t *testing.T) {
	t.Setenv("CLAUDE", "1")
	agent := DetectAgent("", t.TempDir())

	if agent.Name != "claude-code" {
		t.Errorf("expected Name=claude-code, got %s", agent.Name)
	}
	if agent.Source != "env" {
		t.Errorf("expected Source=env, got %s", agent.Source)
	}
}

func TestDetectAgent_FilesystemClaude(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".claude"), 0755)

	agent := DetectAgent("", dir)

	if agent.Name != "claude-code" {
		t.Errorf("expected Name=claude-code, got %s", agent.Name)
	}
	if agent.Source != "filesystem" {
		t.Errorf("expected Source=filesystem, got %s", agent.Source)
	}
}

func TestDetectAgent_FilesystemCursor(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".cursor"), 0755)

	agent := DetectAgent("", dir)

	if agent.Name != "cursor" {
		t.Errorf("expected Name=cursor, got %s", agent.Name)
	}
	if agent.Source != "filesystem" {
		t.Errorf("expected Source=filesystem, got %s", agent.Source)
	}
}

func TestDetectAgent_Unknown(t *testing.T) {
	agent := DetectAgent("", t.TempDir())

	if agent.Name != "" {
		t.Errorf("expected empty Name, got %s", agent.Name)
	}
	if agent.Source != "unknown" {
		t.Errorf("expected Source=unknown, got %s", agent.Source)
	}
}

func TestDetectAgent_EnvPriorityOverFilesystem(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".cursor"), 0755)
	t.Setenv("CLAUDE", "1")

	agent := DetectAgent("", dir)

	if agent.Name != "claude-code" {
		t.Errorf("expected Name=claude-code (env priority), got %s", agent.Name)
	}
	if agent.Source != "env" {
		t.Errorf("expected Source=env, got %s", agent.Source)
	}
}

func TestDetectAllAgents_MultipleTools(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".claude"), 0755)
	os.MkdirAll(filepath.Join(dir, ".cursor"), 0755)
	os.MkdirAll(filepath.Join(dir, ".gemini"), 0755)

	agents := DetectAllAgents(dir)

	if len(agents) < 3 {
		t.Errorf("expected at least 3 agents, got %d", len(agents))
	}

	names := make(map[string]bool)
	for _, a := range agents {
		names[a.Name] = true
		if a.Source != "filesystem" {
			t.Errorf("expected Source=filesystem for all, got %s", a.Source)
		}
	}

	for _, expected := range []string{"claude-code", "cursor", "gemini"} {
		if !names[expected] {
			t.Errorf("expected %s in detected agents", expected)
		}
	}
}

func TestDetectAllAgents_Empty(t *testing.T) {
	agents := DetectAllAgents(t.TempDir())

	if len(agents) != 0 {
		t.Errorf("expected 0 agents, got %d", len(agents))
	}
}

func TestDetectAllAgents_NoDuplicates(t *testing.T) {
	dir := t.TempDir()
	// Both .claude/ and CLAUDE.md map to claude-code
	os.MkdirAll(filepath.Join(dir, ".claude"), 0755)
	os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("test"), 0644)

	agents := DetectAllAgents(dir)

	claudeCount := 0
	for _, a := range agents {
		if a.Name == "claude-code" {
			claudeCount++
		}
	}

	if claudeCount > 1 {
		t.Errorf("expected at most 1 claude-code, got %d", claudeCount)
	}
}
