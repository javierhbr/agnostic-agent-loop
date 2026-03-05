package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/javierbenavides/agentic-agent/internal/skills"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage Claude Code custom agents",
	Long: `Manage Claude Code custom agents (.claude/agents/).

Available agent packs:
  - claude-code: Base agents (orchestrator, worker, researcher, reviewer)
  - openclaw: OpenClaw agents (same as claude-code)
  - openclaw-coordinator: Coordinator agents (tech-lead, product-lead, backend-dev, frontend-dev, mobile-dev, qa-dev)

Usage:
  agentic-agent agent list                      # List available agent packs
  agentic-agent agent install claude-code       # Install agents
  agentic-agent agent install openclaw          # Install openclaw agents
  agentic-agent agent install openclaw --global # Install globally (~/.claude/agents/)
`,
}

var agentInstallCmd = &cobra.Command{
	Use:   "install [pack-name]",
	Short: "Install agent definition files to .claude/agents/",
	Long: `Install agent definition files to .claude/agents/ (project) or ~/.claude/agents/ (global).

Examples:
  agentic-agent agent install claude-code      # Install to .claude/agents/
  agentic-agent agent install openclaw         # Install to .claude/agents/
  agentic-agent agent install openclaw --global # Install to ~/.claude/agents/
`,
	RunE: agentInstallRun,
}

var agentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available agent packs",
	RunE:  agentListRun,
}

// agentInstallRun handles `agentic-agent agent install [pack-name]`.
func agentInstallRun(cmd *cobra.Command, args []string) error {
	globalFlag, _ := cmd.Flags().GetBool("global")

	// Pack name is required
	if len(args) == 0 {
		return fmt.Errorf("pack-name required. Use 'agentic-agent agent list' to see available packs")
	}

	packName := args[0]
	return installAgentPack(packName, globalFlag)
}

// installAgentPack directly installs a pack.
func installAgentPack(packName string, global bool) error {
	filesWritten, err := skills.InstallAgents(packName, global)
	if err != nil {
		return fmt.Errorf("failed to install agents: %w", err)
	}

	if global {
		fmt.Printf("✓ Installed %d agent(s) globally to ~/.claude/agents/\n", len(filesWritten))
	} else {
		fmt.Printf("✓ Installed %d agent(s) to .claude/agents/\n", len(filesWritten))
	}

	for _, f := range filesWritten {
		// Extract just the filename for cleaner output
		name := filepath.Base(f)
		fmt.Printf("  • %s\n", name)
	}

	return nil
}

// agentListRun handles `agentic-agent agent list`.
func agentListRun(cmd *cobra.Command, args []string) error {
	packs := skills.ListAgentPacks()

	fmt.Println("Available Agent Packs:")

	for _, pack := range packs {
		fmt.Printf("  %s\n", pack.Name)
		fmt.Printf("    %s\n", pack.Description)
		fmt.Printf("    %d agents\n\n", len(pack.Files))
	}

	return nil
}

func init() {
	// Register subcommands
	agentInstallCmd.Flags().BoolP("global", "g", false, "Install to ~/.claude/agents/ instead of .claude/agents/")
	agentCmd.AddCommand(agentInstallCmd)
	agentCmd.AddCommand(agentListCmd)

	// Register parent command
	rootCmd.AddCommand(agentCmd)
}
