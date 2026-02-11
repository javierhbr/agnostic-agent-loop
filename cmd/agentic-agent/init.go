package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/project"
	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	uimodels "github.com/javierbenavides/agentic-agent/internal/ui/models"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new agentic project (interactive wizard or flags)",
	Long: `Initialize a new Agentic Agent project with either interactive mode or traditional flags.

Interactive Mode (no flags):
  agentic-agent init

  Launches a step-by-step wizard for project setup with:
  - Project name configuration
  - AI model selection
  - Validator preferences
  - Preview before initialization

Flag Mode (with flags):
  agentic-agent init --name "My Project"

  Traditional command-line mode with project name as a flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if we should use interactive mode
		if helpers.ShouldUseInteractiveMode(cmd) {
			runInteractiveInit()
			return
		}

		// Traditional flag-based mode
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			name = "my-agentic-project"
		}

		desc, _ := cmd.Flags().GetString("description")
		tech, _ := cmd.Flags().GetString("tech-stack")
		wf, _ := cmd.Flags().GetString("workflow")

		var profile *project.ProjectProfile
		if desc != "" || tech != "" || wf != "" {
			profile = &project.ProjectProfile{
				Description: desc,
				TechStack:   tech,
				Workflow:    wf,
			}
		}

		if err := project.InitProjectWithProfile(name, profile); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// runInteractiveInit runs the interactive init wizard
func runInteractiveInit() {
	model := uimodels.NewInitWizardModel()
	p := tea.NewProgram(model)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running init wizard: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	initCmd.Flags().String("name", "", "Name of the project")
	initCmd.Flags().String("description", "", "Project description (populates global context)")
	initCmd.Flags().String("tech-stack", "", "Tech stack summary (e.g., 'Go, React, PostgreSQL')")
	initCmd.Flags().String("workflow", "", "Workflow preferences (e.g., 'TDD, trunk-based, PR reviews')")
}
