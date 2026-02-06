package main

import (
	"github.com/spf13/cobra"
)

var (
	// Version information (set via ldflags at build time)
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"

	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "agentic-agent",
	Short: "Agent-agnostic, specification-driven development CLI",
	Long: `agentic-agent is a CLI tool that enables agent-agnostic, specification-driven
development with token management, context isolation, and multi-tool support.

It orchestrates AI coding agents (Claude Code, Cursor, Windsurf, etc.) through
a unified task and context management system.`,
	SilenceUsage: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./agnostic-agent.yaml)")
	rootCmd.PersistentFlags().Bool("no-interactive", false, "disable interactive mode and use flag-based commands")

	// Register commands
	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(taskCmd)
	rootCmd.AddCommand(contextCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(skillsCmd)
	rootCmd.AddCommand(tokenCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(learningsCmd)
}
