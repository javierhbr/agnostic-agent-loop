package helpers

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

// ShouldUseInteractiveMode determines if the command should run in interactive mode
// based on flags, terminal state, and whether an AI agent is driving the CLI.
//
// Interactive mode is enabled when:
// - -i/--interactive flag is set (overrides everything else), OR
// - No command-specific flags are provided AND the input is a TTY AND no agent detected
//
// Interactive mode is disabled when:
// - --no-interactive flag is set
// - An AI agent is detected (--agent flag, AGENTIC_AGENT env, or filesystem detection)
// - Command-specific flags are provided
// - stdin is not a TTY
//
// When an AI agent is driving the CLI, the full-screen Bubble Tea TUI would
// flood the terminal renderer with escape sequences and alt-screen buffers,
// causing VSCode and other IDEs to freeze. Agents should always get plain text output.
func ShouldUseInteractiveMode(cmd *cobra.Command) bool {
	// Check if explicitly forced on — overrides agent detection and flag-based mode
	if interactive, _ := cmd.Flags().GetBool("interactive"); interactive {
		return true
	}

	// Check if explicitly disabled
	if noInteractive, _ := cmd.Flags().GetBool("no-interactive"); noInteractive {
		return false
	}

	// Disable interactive mode when an AI agent is driving the CLI.
	// Agents can't interact with Bubble Tea TUIs, and alt-screen mode
	// crashes VSCode's terminal renderer.
	if isAgentDriven(cmd) {
		return false
	}

	// Check if any command-specific flags provided
	hasFlags := false
	cmd.Flags().Visit(func(f *pflag.Flag) {
		// Ignore global flags
		if f.Name != "config" && f.Name != "no-interactive" && f.Name != "interactive" && f.Name != "agent" {
			hasFlags = true
		}
	})

	// If flags provided, use traditional mode
	if hasFlags {
		return false
	}

	// No flags + TTY = interactive mode
	return isTerminal()
}

// activeAgent holds the agent name detected by PersistentPreRunE.
// Set via SetActiveAgent() after filesystem/env/flag detection completes.
var activeAgent string

// SetActiveAgent records the detected agent name so all commands
// automatically disable interactive TUI mode when an agent is driving the CLI.
// Called from root.go's PersistentPreRunE after agent detection.
func SetActiveAgent(name string) {
	activeAgent = name
}

// agentEnvVars lists environment variables that indicate an AI agent is running.
// Matches the env vars in skills/detect.go.
var agentEnvVars = []string{
	"AGENTIC_AGENT",
	"CLAUDE",
	"CLAUDE_CODE",
	"CURSOR_SESSION",
	"GEMINI_CLI",
	"WINDSURF_SESSION",
	"CODEX_SANDBOX",
	"GITHUB_COPILOT",
}

// isAgentDriven returns true when an AI agent is controlling the CLI.
// Checks: ActiveAgent (filesystem detection) > --agent flag > env vars.
func isAgentDriven(cmd *cobra.Command) bool {
	// Check ActiveAgent set by PersistentPreRunE (covers filesystem detection)
	if activeAgent != "" {
		return true
	}

	// Check --agent flag
	if agent, _ := cmd.Flags().GetString("agent"); agent != "" {
		return true
	}

	// Check all known agent environment variables
	for _, envVar := range agentEnvVars {
		if os.Getenv(envVar) != "" {
			return true
		}
	}

	return false
}

// isTerminal checks if stdin is connected to a terminal (TTY)
func isTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

// IsCI checks if we're running in a CI environment
func IsCI() bool {
	ciVars := []string{"CI", "CONTINUOUS_INTEGRATION", "GITHUB_ACTIONS", "GITLAB_CI", "CIRCLECI"}
	for _, v := range ciVars {
		if os.Getenv(v) != "" {
			return true
		}
	}
	return false
}
