package helpers

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

// ShouldUseInteractiveMode determines if the command should run in interactive mode
// based on flags and terminal state.
//
// Interactive mode is enabled when:
// - No command-specific flags are provided
// - The input is a TTY (not a pipe or file)
// - --no-interactive flag is not set
func ShouldUseInteractiveMode(cmd *cobra.Command) bool {
	// Check if explicitly disabled
	if noInteractive, _ := cmd.Flags().GetBool("no-interactive"); noInteractive {
		return false
	}

	// Check if any command-specific flags provided
	hasFlags := false
	cmd.Flags().Visit(func(f *pflag.Flag) {
		// Ignore global flags
		if f.Name != "config" && f.Name != "no-interactive" {
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
