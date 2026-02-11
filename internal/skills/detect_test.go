package skills

import (
	"os"
	"path/filepath"
	"testing"
)

// ---------------------------------------------------------------------------
// DetectAgent — Flag detection
// ---------------------------------------------------------------------------

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

func TestDetectAgent_FlagOverridesFilesystem(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".cursor"), 0755)

	agent := DetectAgent("windsurf", dir)

	if agent.Name != "windsurf" {
		t.Errorf("expected Name=windsurf, got %s", agent.Name)
	}
	if agent.Source != "flag" {
		t.Errorf("expected Source=flag, got %s", agent.Source)
	}
}

func TestDetectAgent_FlagArbitraryValue(t *testing.T) {
	agent := DetectAgent("custom-agent", t.TempDir())

	if agent.Name != "custom-agent" {
		t.Errorf("expected Name=custom-agent, got %s", agent.Name)
	}
	if agent.Source != "flag" {
		t.Errorf("expected Source=flag, got %s", agent.Source)
	}
}

// ---------------------------------------------------------------------------
// DetectAgent — Environment variable detection
// ---------------------------------------------------------------------------

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

func TestDetectAgent_AgenticEnvVarArbitraryValue(t *testing.T) {
	t.Setenv("AGENTIC_AGENT", "my-custom-tool")
	agent := DetectAgent("", t.TempDir())

	if agent.Name != "my-custom-tool" {
		t.Errorf("expected Name=my-custom-tool, got %s", agent.Name)
	}
	if agent.Source != "env" {
		t.Errorf("expected Source=env, got %s", agent.Source)
	}
}

func TestDetectAgent_AgenticEnvPriorityOverAgentSpecific(t *testing.T) {
	t.Setenv("AGENTIC_AGENT", "gemini")
	t.Setenv("CLAUDE", "1")
	t.Setenv("CURSOR_SESSION", "1")

	agent := DetectAgent("", t.TempDir())

	if agent.Name != "gemini" {
		t.Errorf("expected Name=gemini (AGENTIC_AGENT priority), got %s", agent.Name)
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

func TestDetectAgent_ClaudeCodeEnvVar(t *testing.T) {
	t.Setenv("CLAUDE_CODE", "1")
	agent := DetectAgent("", t.TempDir())

	if agent.Name != "claude-code" {
		t.Errorf("expected Name=claude-code, got %s", agent.Name)
	}
	if agent.Source != "env" {
		t.Errorf("expected Source=env, got %s", agent.Source)
	}
}

func TestDetectAgent_CursorSessionEnvVar(t *testing.T) {
	t.Setenv("CURSOR_SESSION", "abc123")
	agent := DetectAgent("", t.TempDir())

	if agent.Name != "cursor" {
		t.Errorf("expected Name=cursor, got %s", agent.Name)
	}
	if agent.Source != "env" {
		t.Errorf("expected Source=env, got %s", agent.Source)
	}
}

func TestDetectAgent_GeminiCliEnvVar(t *testing.T) {
	t.Setenv("GEMINI_CLI", "1")
	agent := DetectAgent("", t.TempDir())

	if agent.Name != "gemini" {
		t.Errorf("expected Name=gemini, got %s", agent.Name)
	}
	if agent.Source != "env" {
		t.Errorf("expected Source=env, got %s", agent.Source)
	}
}

func TestDetectAgent_WindsurfSessionEnvVar(t *testing.T) {
	t.Setenv("WINDSURF_SESSION", "sess-42")
	agent := DetectAgent("", t.TempDir())

	if agent.Name != "windsurf" {
		t.Errorf("expected Name=windsurf, got %s", agent.Name)
	}
	if agent.Source != "env" {
		t.Errorf("expected Source=env, got %s", agent.Source)
	}
}

func TestDetectAgent_CodexSandboxEnvVar(t *testing.T) {
	t.Setenv("CODEX_SANDBOX", "1")
	agent := DetectAgent("", t.TempDir())

	if agent.Name != "codex" {
		t.Errorf("expected Name=codex, got %s", agent.Name)
	}
	if agent.Source != "env" {
		t.Errorf("expected Source=env, got %s", agent.Source)
	}
}

// ---------------------------------------------------------------------------
// DetectAgent — Priority: env over filesystem
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// DetectAgent — Filesystem heuristics
// ---------------------------------------------------------------------------

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

func TestDetectAgent_FilesystemClaudeMd(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# rules"), 0644)

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

func TestDetectAgent_FilesystemGemini(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".gemini"), 0755)

	agent := DetectAgent("", dir)

	if agent.Name != "gemini" {
		t.Errorf("expected Name=gemini, got %s", agent.Name)
	}
	if agent.Source != "filesystem" {
		t.Errorf("expected Source=filesystem, got %s", agent.Source)
	}
}

func TestDetectAgent_FilesystemWindsurf(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".windsurf"), 0755)

	agent := DetectAgent("", dir)

	if agent.Name != "windsurf" {
		t.Errorf("expected Name=windsurf, got %s", agent.Name)
	}
	if agent.Source != "filesystem" {
		t.Errorf("expected Source=filesystem, got %s", agent.Source)
	}
}

func TestDetectAgent_FilesystemCodex(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".codex"), 0755)

	agent := DetectAgent("", dir)

	if agent.Name != "codex" {
		t.Errorf("expected Name=codex, got %s", agent.Name)
	}
	if agent.Source != "filesystem" {
		t.Errorf("expected Source=filesystem, got %s", agent.Source)
	}
}

func TestDetectAgent_FilesystemAntigravity(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".agent"), 0755)

	agent := DetectAgent("", dir)

	if agent.Name != "antigravity" {
		t.Errorf("expected Name=antigravity, got %s", agent.Name)
	}
	if agent.Source != "filesystem" {
		t.Errorf("expected Source=filesystem, got %s", agent.Source)
	}
}

// ---------------------------------------------------------------------------
// DetectAgent — Unknown / empty
// ---------------------------------------------------------------------------

func TestDetectAgent_Unknown(t *testing.T) {
	agent := DetectAgent("", t.TempDir())

	if agent.Name != "" {
		t.Errorf("expected empty Name, got %s", agent.Name)
	}
	if agent.Source != "unknown" {
		t.Errorf("expected Source=unknown, got %s", agent.Source)
	}
}

func TestDetectAgent_EmptyFlagAndNoEnvOrFS(t *testing.T) {
	dir := t.TempDir()
	// ensure no agent env vars are set
	for envVar := range envMapping {
		t.Setenv(envVar, "")
	}

	agent := DetectAgent("", dir)

	if agent.Name != "" {
		t.Errorf("expected empty Name, got %s", agent.Name)
	}
	if agent.Source != "unknown" {
		t.Errorf("expected Source=unknown, got %s", agent.Source)
	}
}

// ---------------------------------------------------------------------------
// DetectAllAgents
// ---------------------------------------------------------------------------

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

func TestDetectAllAgents_AllTools(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".claude"), 0755)
	os.MkdirAll(filepath.Join(dir, ".cursor"), 0755)
	os.MkdirAll(filepath.Join(dir, ".gemini"), 0755)
	os.MkdirAll(filepath.Join(dir, ".windsurf"), 0755)
	os.MkdirAll(filepath.Join(dir, ".codex"), 0755)
	os.MkdirAll(filepath.Join(dir, ".agent"), 0755)

	agents := DetectAllAgents(dir)

	names := make(map[string]bool)
	for _, a := range agents {
		names[a.Name] = true
	}

	for _, expected := range []string{"claude-code", "cursor", "gemini", "windsurf", "codex", "antigravity"} {
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

func TestDetectAllAgents_SingleTool(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".windsurf"), 0755)

	agents := DetectAllAgents(dir)

	if len(agents) != 1 {
		t.Errorf("expected 1 agent, got %d", len(agents))
	}
	if len(agents) > 0 && agents[0].Name != "windsurf" {
		t.Errorf("expected Name=windsurf, got %s", agents[0].Name)
	}
}
