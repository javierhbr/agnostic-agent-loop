package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/ui/models"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Interactive project setup wizard",
	Long: `Start an interactive wizard to set up a new Agnostic Agent project.

This command provides a friendly, step-by-step guide through:
- Project naming and configuration
- AI model selection
- Directory structure creation
- Initial setup

Perfect for first-time users or quick project initialization.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create and run the start wizard
		model := models.NewStartWizardModel()
		p := tea.NewProgram(model, tea.WithAltScreen())

		finalModel, err := p.Run()
		if err != nil {
			fmt.Printf("Error running start wizard: %v\n", err)
			os.Exit(1)
		}

		// Get the selected action
		if wizardModel, ok := finalModel.(models.StartWizardModel); ok {
			switch wizardModel.SelectedAction() {
			// First-Level Commands
			case "init":
				runInteractiveInit()

			case "task":
				// Show task submenu
				runTaskSubmenu()

			case "work":
				runWorkWorkflow()

			case "context":
				// Show context submenu
				runContextSubmenu()

			case "validate":
				validateCmd.Run(nil, []string{})

			case "skills":
				// Show skills submenu
				runSkillsSubmenu()

			case "token":
				tokenStatusCmd.Run(nil, []string{})

			case "run":
				runCmd.Run(nil, []string{})

			case "version":
				fmt.Printf("agentic-agent %s\n", Version)
				fmt.Printf("  Commit:     %s\n", Commit)
				fmt.Printf("  Build Date: %s\n", BuildDate)

			case "help":
				rootCmd.Help()

			case "exit":
				// Do nothing, just exit
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
