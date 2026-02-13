package main

import (
	"github.com/javierbenavides/agentic-agent/internal/config"
	"github.com/javierbenavides/agentic-agent/internal/skills"
	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/spf13/cobra"
)

var (
	// Version information (set via ldflags at build time)
	Version   = "0.2.16"
	Commit    = "none"
	BuildDate = "unknown"

	cfgFile       string
	agentFlag     string
	appConfig     *models.Config
	detectedAgent skills.DetectedAgent
)

// getConfig returns the loaded config, falling back to defaults if not loaded.
func getConfig() *models.Config {
	if appConfig != nil {
		return appConfig
	}
	cfg := &models.Config{}
	config.SetDefaults(cfg)
	appConfig = cfg
	return appConfig
}

// getAgent returns the detected agent, falling back to unknown.
func getAgent() skills.DetectedAgent {
	return detectedAgent
}

var rootCmd = &cobra.Command{
	Use:   "agentic-agent",
	Short: "Agent-agnostic, specification-driven development CLI",
	Long: `agentic-agent is a CLI tool that enables agent-agnostic, specification-driven
development with token management, context isolation, and multi-tool support.

It orchestrates AI coding agents (Claude Code, Cursor, Windsurf, etc.) through
a unified task and context management system.`,
	Version:      Version,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		path := cfgFile
		if path == "" {
			path = "agnostic-agent.yaml"
		}
		cfg, err := config.LoadConfig(path)
		if err != nil {
			// Config file missing is OK — use defaults
			cfg = &models.Config{}
			config.SetDefaults(cfg)
		}
		appConfig = cfg

		// Detect active agent and propagate to UI helpers
		// so all commands auto-disable interactive TUI when an agent is driving.
		detectedAgent = skills.DetectAgent(agentFlag, ".")
		if detectedAgent.Name != "" {
			appConfig.ActiveAgent = detectedAgent.Name
			helpers.SetActiveAgent(detectedAgent.Name)
		}

		// Auto-ensure mandatory skills for detected agent.
		// Skip for commands that don't need it (version, help, skills, start, init).
		cmdName := cmd.Name()
		skipEnsure := cmdName == "version" || cmdName == "help" || cmdName == "start" ||
			cmdName == "init" || cmd.Parent() != nil && cmd.Parent().Name() == "skills"
		if !skipEnsure && detectedAgent.Name != "" {
			// Silent ensure — only installs missing mandatory packs
			skills.Ensure(detectedAgent.Name, appConfig)
		}

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./agnostic-agent.yaml)")
	rootCmd.PersistentFlags().StringVar(&agentFlag, "agent", "", "agent tool name (claude-code, cursor, gemini, windsurf, antigravity, codex, copilot, opencode)")
	rootCmd.PersistentFlags().Bool("no-interactive", false, "disable interactive mode and use flag-based commands")
	rootCmd.PersistentFlags().BoolP("interactive", "i", false, "force interactive mode (overrides agent detection and flag-based mode)")

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
	rootCmd.AddCommand(specCmd)
	rootCmd.AddCommand(autopilotCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(trackCmd)
	rootCmd.AddCommand(planCmd)
	rootCmd.AddCommand(simplifyCmd)
	rootCmd.AddCommand(openspecCmd)
	rootCmd.AddCommand(promptsCmd)
}
