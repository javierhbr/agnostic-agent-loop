package skills

import (
	"os"
	"path/filepath"
)

// DetectedAgent holds the result of agent detection.
type DetectedAgent struct {
	Name   string // e.g., "claude-code", "cursor", "gemini"
	Source string // "flag", "env", "filesystem", "unknown"
}

var envOrder = []struct {
	Var   string
	Agent string
}{
	{"AGENTIC_AGENT", ""},
	{"CLAUDE", "claude-code"},
	{"CLAUDE_CODE", "claude-code"},
	{"CURSOR_SESSION", "cursor"},
	{"GEMINI_CLI", "gemini"},
	{"WINDSURF_SESSION", "windsurf"},
	{"GITHUB_COPILOT", "copilot"},
	{"OPENCODE", "opencode"},
}

// fsMapping maps filesystem paths to agent tool names.
var fsMapping = map[string]string{
	".claude":                         "claude-code",
	"CLAUDE.md":                       "claude-code",
	".cursor":                         "cursor",
	".gemini":                         "gemini",
	".windsurf":                       "windsurf",
	".codex":                          "codex",
	".agent":                          "antigravity",
	".github/copilot-instructions.md": "copilot",
	".opencode":                       "opencode",
}

// DetectAgent tries to determine which agent is running.
// Priority: explicit flag > env var > filesystem.
func DetectAgent(flagValue string, projectRoot string) DetectedAgent {
	// 1. Explicit flag
	if flagValue != "" {
		return DetectedAgent{Name: flagValue, Source: "flag"}
	}

	// 2. Environment variables (ordered: our own convention first)
	if val := os.Getenv("AGENTIC_AGENT"); val != "" {
		return DetectedAgent{Name: val, Source: "env"}
	}
	for _, entry := range envOrder {
		if entry.Var == "AGENTIC_AGENT" {
			continue
		}
		if os.Getenv(entry.Var) != "" {
			return DetectedAgent{Name: entry.Agent, Source: "env"}
		}
	}

	// 3. Filesystem heuristics
	for path, agentName := range fsMapping {
		fullPath := filepath.Join(projectRoot, path)
		if _, err := os.Stat(fullPath); err == nil {
			return DetectedAgent{Name: agentName, Source: "filesystem"}
		}
	}

	// 4. Fallback: CODEX_SANDBOX only if nothing else matched
	if os.Getenv("CODEX_SANDBOX") != "" {
		return DetectedAgent{Name: "codex", Source: "env"}
	}

	return DetectedAgent{Source: "unknown"}
}

// DetectAllAgents returns all agents detected via filesystem heuristics.
func DetectAllAgents(projectRoot string) []DetectedAgent {
	seen := make(map[string]bool)
	var agents []DetectedAgent

	for path, agentName := range fsMapping {
		if seen[agentName] {
			continue
		}
		fullPath := filepath.Join(projectRoot, path)
		if _, err := os.Stat(fullPath); err == nil {
			agents = append(agents, DetectedAgent{Name: agentName, Source: "filesystem"})
			seen[agentName] = true
		}
	}

	return agents
}
