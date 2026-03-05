package skills

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AgentPack represents a collection of Claude Code agent definition files.
type AgentPack struct {
	Name        string      // e.g., "claude-code", "openclaw"
	Description string      // e.g., "Base OpenClaw agents..."
	Files       []AgentFile // Agent .md files to install
}

// AgentFile represents a single agent definition file.
type AgentFile struct {
	SrcPath string // Path in embedded FS
	DstPath string // Filename in .claude/agents/
}

//go:embed packs/claude-code-custom-agent/agents/*
var agentsFS embed.FS

// AgentRegistry provides access to all available agent packs.
var AgentRegistry = map[string]AgentPack{
	"claude-code": {
		Name:        "claude-code",
		Description: "Base Claude Code agents: orchestrator, worker, researcher, reviewer",
		Files: []AgentFile{
			{SrcPath: "packs/claude-code-custom-agent/agents/orchestrator.md", DstPath: "orchestrator.md"},
			{SrcPath: "packs/claude-code-custom-agent/agents/worker.md", DstPath: "worker.md"},
			{SrcPath: "packs/claude-code-custom-agent/agents/researcher.md", DstPath: "researcher.md"},
			{SrcPath: "packs/claude-code-custom-agent/agents/reviewer.md", DstPath: "reviewer.md"},
		},
	},
	"openclaw": {
		Name:        "openclaw",
		Description: "OpenClaw autonomous agents: orchestrator, worker, researcher, reviewer",
		Files: []AgentFile{
			{SrcPath: "packs/claude-code-custom-agent/agents/orchestrator.md", DstPath: "orchestrator.md"},
			{SrcPath: "packs/claude-code-custom-agent/agents/worker.md", DstPath: "worker.md"},
			{SrcPath: "packs/claude-code-custom-agent/agents/researcher.md", DstPath: "researcher.md"},
			{SrcPath: "packs/claude-code-custom-agent/agents/reviewer.md", DstPath: "reviewer.md"},
		},
	},
	"openclaw-coordinator": {
		Name:        "openclaw-coordinator",
		Description: "Coordinator agents: tech-lead, product-lead, backend-dev, frontend-dev, mobile-dev, qa-dev",
		Files: []AgentFile{
			{SrcPath: "packs/claude-code-custom-agent/agents/tech-lead.md", DstPath: "tech-lead.md"},
			{SrcPath: "packs/claude-code-custom-agent/agents/product-lead.md", DstPath: "product-lead.md"},
			{SrcPath: "packs/claude-code-custom-agent/agents/backend-dev.md", DstPath: "backend-dev.md"},
			{SrcPath: "packs/claude-code-custom-agent/agents/frontend-dev.md", DstPath: "frontend-dev.md"},
			{SrcPath: "packs/claude-code-custom-agent/agents/mobile-dev.md", DstPath: "mobile-dev.md"},
			{SrcPath: "packs/claude-code-custom-agent/agents/qa-dev.md", DstPath: "qa-dev.md"},
		},
	},
}

// InstallAgents installs agent files from a pack to .claude/agents/ or ~/.claude/agents/ (if global).
func InstallAgents(packName string, global bool) ([]string, error) {
	pack, exists := AgentRegistry[packName]
	if !exists {
		return nil, fmt.Errorf("agent pack not found: %s", packName)
	}

	// Determine output directory
	outputDir := ".claude/agents"
	if global {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to resolve home directory: %w", err)
		}
		outputDir = filepath.Join(home, ".claude/agents")
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create agent directory %s: %w", outputDir, err)
	}

	var filesWritten []string

	// Write each agent file
	for _, file := range pack.Files {
		content, err := agentsFS.ReadFile(file.SrcPath)
		if err != nil {
			return filesWritten, fmt.Errorf("failed to read embedded agent file %s: %w", file.SrcPath, err)
		}

		destPath := filepath.Join(outputDir, file.DstPath)
		if err := os.WriteFile(destPath, content, 0644); err != nil {
			return filesWritten, fmt.Errorf("failed to write agent file %s: %w", destPath, err)
		}

		filesWritten = append(filesWritten, destPath)
	}

	return filesWritten, nil
}

// ListAgentPacks returns all available agent packs.
func ListAgentPacks() []AgentPack {
	var packs []AgentPack
	for _, pack := range AgentRegistry {
		packs = append(packs, pack)
	}
	return packs
}

// GetAgentPack returns a specific agent pack by name.
func GetAgentPack(name string) (AgentPack, error) {
	pack, exists := AgentRegistry[name]
	if !exists {
		return AgentPack{}, fmt.Errorf("agent pack not found: %s", name)
	}
	return pack, nil
}

// IsAgentInstalled checks whether an agent pack's files are installed in .claude/agents/.
func IsAgentInstalled(packName string, global bool) bool {
	pack, exists := AgentRegistry[packName]
	if !exists {
		return false
	}

	// Determine check directory
	checkDir := ".claude/agents"
	if global {
		home, err := os.UserHomeDir()
		if err != nil {
			return false
		}
		checkDir = filepath.Join(home, ".claude/agents")
	}

	// Check each agent file
	for _, file := range pack.Files {
		path := filepath.Join(checkDir, file.DstPath)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// SupportedAgentTools returns the list of tools that support agent installation.
func SupportedAgentTools() []string {
	return []string{"claude-code"}
}

// GetAgentDirectoryForTool returns the agent directory for a specific tool.
func GetAgentDirectoryForTool(tool string, global bool) (string, error) {
	if tool != "claude-code" {
		return "", fmt.Errorf("tool %s does not support agent installation", tool)
	}

	if global {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to resolve home directory: %w", err)
		}
		return filepath.Join(home, ".claude/agents"), nil
	}

	return ".claude/agents", nil
}

// FormatAgentPackList returns a formatted string listing all agent packs.
func FormatAgentPackList() string {
	var sb strings.Builder
	sb.WriteString("Available Agent Packs:\n\n")

	for _, pack := range ListAgentPacks() {
		fmt.Fprintf(&sb, "  %s\n", pack.Name)
		fmt.Fprintf(&sb, "    %s\n", pack.Description)
		fmt.Fprintf(&sb, "    Agents: %d\n\n", len(pack.Files))
	}

	return sb.String()
}
