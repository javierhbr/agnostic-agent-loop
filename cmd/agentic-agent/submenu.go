package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/ui/components"
	"github.com/javierbenavides/agentic-agent/internal/ui/models"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
)

// submenuModel handles interactive submenus for commands with multiple subcommands
type submenuModel struct {
	selector       components.SimpleSelect
	selectedAction string
	done           bool
}

func (m submenuModel) Init() tea.Cmd {
	return nil
}

func (m submenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.done = true
			m.selectedAction = "exit"
			return m, tea.Quit

		case "enter":
			m.selectedAction = m.selector.SelectedOption().Value()
			m.done = true
			return m, tea.Quit

		default:
			m.selector = m.selector.Update(msg)
		}
	}

	return m, nil
}

func (m submenuModel) View() string {
	if m.done {
		return ""
	}

	var result string
	result += m.selector.View() + "\n"
	result += styles.HelpStyle.Render("↑/↓ navigate • Enter select • Esc cancel") + "\n"

	return styles.ContainerStyle.Render(result)
}

// runTaskSubmenu shows the task command submenu
func runTaskSubmenu() {
	options := []components.SelectOption{
		components.NewSelectOption(
			"📋 Create a task",
			"Create a new task with interactive wizard",
			"create",
		),
		components.NewSelectOption(
			"📝 Create from template",
			"Use pre-built templates (feature, bug-fix, etc.)",
			"from-template",
		),
		components.NewSelectOption(
			"📊 List tasks",
			"Browse tasks in Backlog, In Progress, and Done",
			"list",
		),
		components.NewSelectOption(
			"🎯 Claim a task",
			"Select and claim a task to work on",
			"claim",
		),
		components.NewSelectOption(
			"✅ Complete a task",
			"Mark an in-progress task as done",
			"complete",
		),
		components.NewSelectOption(
			"🔍 Show task details",
			"View detailed information about a task",
			"show",
		),
		components.NewSelectOption(
			"🔢 Decompose task",
			"Break a task into subtasks",
			"decompose",
		),
		components.NewSelectOption(
			"📝 Sample task",
			"Create a sample task with example data",
			"sample",
		),
		components.NewSelectOption(
			"🚪 Back",
			"Return to main menu",
			"exit",
		),
	}

	selector := components.NewSimpleSelect("task - Manage tasks", options)
	model := &submenuModel{
		selector: selector,
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(submenuModel); ok {
		switch m.selectedAction {
		case "create":
			runInteractiveTaskCreate()
		case "from-template":
			runTemplateWorkflow()
		case "list":
			runInteractiveTaskList()
		case "claim":
			claimModel := models.NewSimpleTaskSelectModel(models.ActionClaim, "backlog")
			cp := tea.NewProgram(claimModel, tea.WithAltScreen())
			if fm, err := cp.Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else if m, ok := fm.(models.SimpleTaskSelectModel); ok && m.Done() {
				fmt.Println(m.ResultMessage())
			}
		case "complete":
			completeModel := models.NewSimpleTaskSelectModel(models.ActionComplete, "in-progress")
			cp := tea.NewProgram(completeModel, tea.WithAltScreen())
			if fm, err := cp.Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else if m, ok := fm.(models.SimpleTaskSelectModel); ok && m.Done() {
				fmt.Println(m.ResultMessage())
			}
		case "show":
			showModel := models.NewTaskSelectModel()
			sp := tea.NewProgram(showModel, tea.WithAltScreen())
			if _, err := sp.Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case "decompose":
			taskDecomposeCmd.Run(taskDecomposeCmd, []string{})
		case "sample":
			taskSampleCmd.Run(taskSampleCmd, []string{})
		case "exit":
			// Do nothing, return to main menu
		}
	}
}

// runContextSubmenu shows the context command submenu
func runContextSubmenu() {
	options := []components.SelectOption{
		components.NewSelectOption(
			"📁 Generate context",
			"Generate context.md for a directory",
			"generate",
		),
		components.NewSelectOption(
			"🔎 Scan for context",
			"Find directories missing context files",
			"scan",
		),
		components.NewSelectOption(
			"📦 Build context bundle",
			"Create context bundle for a task",
			"build",
		),
		components.NewSelectOption(
			"🔄 Update context",
			"Update existing context.md file",
			"update",
		),
		components.NewSelectOption(
			"🚪 Back",
			"Return to main menu",
			"exit",
		),
	}

	selector := components.NewSimpleSelect("context - Manage context files", options)
	model := &submenuModel{
		selector: selector,
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(submenuModel); ok {
		switch m.selectedAction {
		case "generate":
			contextGenerateCmd.Run(contextGenerateCmd, []string{})
		case "scan":
			contextScanCmd.Run(contextScanCmd, []string{})
		case "build":
			contextBuildCmd.Run(contextBuildCmd, []string{})
		case "update":
			contextUpdateCmd.Run(contextUpdateCmd, []string{})
		case "exit":
			// Do nothing, return to main menu
		}
	}
}

// runSkillsSubmenu shows the skills command submenu
func runSkillsSubmenu() {
	options := []components.SelectOption{
		components.NewSelectOption(
			"🛠️  Generate skills",
			"Generate agent skill files for your tools",
			"generate",
		),
		components.NewSelectOption(
			"📦 Install skill pack",
			"Install a bundle of skill files for one or more tools",
			"install",
		),
		components.NewSelectOption(
			"📋 List skill packs",
			"Show all available skill packs",
			"list",
		),
		components.NewSelectOption(
			"🔄 Check skill drift",
			"Check if skill files need regeneration",
			"check",
		),
		components.NewSelectOption(
			"✅ Ensure skills",
			"Ensure skills are up to date for an agent",
			"ensure",
		),
		components.NewSelectOption(
			"🚪 Back",
			"Return to main menu",
			"exit",
		),
	}

	selector := components.NewSimpleSelect("skills - Manage agent skills", options)
	model := &submenuModel{
		selector: selector,
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(submenuModel); ok {
		switch m.selectedAction {
		case "generate":
			skillsGenerateCmd.Run(skillsGenerateCmd, []string{})
		case "install":
			skillsInstallCmd.Run(skillsInstallCmd, []string{})
		case "list":
			skillsListCmd.Run(skillsListCmd, []string{})
		case "check":
			skillsCheckCmd.Run(skillsCheckCmd, []string{})
		case "ensure":
			skillsEnsureCmd.Run(skillsEnsureCmd, []string{})
		case "exit":
			// Do nothing, return to main menu
		}
	}
}

// Commands are already declared in their respective files:
// - task.go: taskDecomposeCmd, taskSampleCmd
// - context.go: contextGenerateCmd, contextScanCmd, contextBuildCmd, contextUpdateCmd
// - skills.go: skillsGenerateCmd, skillsInstallCmd, skillsListCmd, skillsCheckCmd, skillsEnsureCmd
// - validate.go: validateCmd
// - token.go: tokenStatusCmd
// - run.go: runCmd
