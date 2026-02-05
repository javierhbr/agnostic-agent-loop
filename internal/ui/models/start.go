package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/project"
	"github.com/javierbenavides/agentic-agent/internal/ui/components"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
)

// StartWizardStep represents the current step in the wizard
type StartWizardStep int

const (
	StepWelcome StartWizardStep = iota
	StepProjectName
	StepModelSelection
	StepConfirmation
	StepInitializing
	StepComplete
	StepNextAction
)

// StartWizardModel is the model for the start wizard
type StartWizardModel struct {
	step         StartWizardStep
	projectName  components.ValidatedInput
	modelSelect  components.SimpleSelect
	nextAction   components.SimpleSelect
	spinner      components.Spinner
	error        string
	width        int
	height       int
	initialized  bool
	quitting     bool
	selectedAction string
}

// NewStartWizardModel creates a new start wizard model
func NewStartWizardModel() StartWizardModel {
	// Project name input with validation
	projectName := components.NewValidatedInput(
		"Project Name",
		"my-awesome-project",
		func(s string) error {
			if len(s) == 0 {
				return fmt.Errorf("project name cannot be empty")
			}
			if len(s) > 100 {
				return fmt.Errorf("project name too long (max 100 characters)")
			}
			if strings.Contains(s, "/") || strings.Contains(s, "\\") {
				return fmt.Errorf("project name cannot contain slashes")
			}
			return nil
		},
	)

	// Model selection
	modelOptions := []components.SelectOption{
		components.NewSelectOption(
			"Claude 3.5 Sonnet (Recommended)",
			"Latest Claude model with excellent reasoning and coding",
			"claude-3-5-sonnet-20241022",
		),
		components.NewSelectOption(
			"GPT-4",
			"OpenAI's powerful language model",
			"gpt-4",
		),
		components.NewSelectOption(
			"Claude 3 Opus",
			"Most capable Claude model for complex tasks",
			"claude-3-opus-20240229",
		),
	}
	modelSelect := components.NewSimpleSelect("AI Model", modelOptions)

	// Next action selection - ALL first-level commands
	nextActionOptions := []components.SelectOption{
		// First-Level Commands (matching agentic-agent --help)
		components.NewSelectOption(
			"🎬 init - Initialize project",
			"Initialize a new agentic project (interactive wizard)",
			"init",
		),
		components.NewSelectOption(
			"📋 task - Manage tasks",
			"Create, list, claim, complete, and manage tasks",
			"task",
		),
		components.NewSelectOption(
			"🔨 work - Interactive workflow",
			"Complete workflow: claim task → work → complete",
			"work",
		),
		components.NewSelectOption(
			"📁 context - Manage context",
			"Generate, scan, and build context files",
			"context",
		),
		components.NewSelectOption(
			"✔️  validate - Run validation",
			"Run all validation rules on the project",
			"validate",
		),
		components.NewSelectOption(
			"🛠️  skills - Manage skills",
			"Generate and check agent skill files",
			"skills",
		),
		components.NewSelectOption(
			"📊 token - Token usage",
			"View and manage token usage statistics",
			"token",
		),
		components.NewSelectOption(
			"🤖 run - Run orchestrator",
			"Run the agent orchestrator for a task",
			"run",
		),
		components.NewSelectOption(
			"ℹ️  version - Version info",
			"Display version information",
			"version",
		),
		components.NewSelectOption(
			"❓ help - Get help",
			"Show help for any command",
			"help",
		),
		components.NewSelectOption(
			"🚪 Exit",
			"Close the CLI",
			"exit",
		),
	}
	nextAction := components.NewSimpleSelect("Select a command:", nextActionOptions)

	spinner := components.NewSpinner("Initializing project...")

	return StartWizardModel{
		step:        StepWelcome,
		projectName: projectName,
		modelSelect: modelSelect,
		nextAction:  nextAction,
		spinner:     spinner,
	}
}

// Init initializes the model
func (m StartWizardModel) Init() tea.Cmd {
	return m.projectName.Focus()
}

// SelectedAction returns the action selected by the user
func (m StartWizardModel) SelectedAction() string {
	return m.selectedAction
}

// Update handles messages
func (m StartWizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.step != StepInitializing {
				m.quitting = true
				return m, tea.Quit
			}

		case "enter":
			return m.handleEnter()
		}

	case projectInitCompleteMsg:
		m.step = StepComplete
		return m, nil

	case projectInitErrorMsg:
		m.error = msg.err.Error()
		m.step = StepComplete
		return m, nil
	}

	// Handle step-specific updates
	switch m.step {
	case StepProjectName:
		m.projectName, cmd = m.projectName.Update(msg)
		return m, cmd

	case StepModelSelection:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			m.modelSelect = m.modelSelect.Update(keyMsg)
		}

	case StepNextAction:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			m.nextAction = m.nextAction.Update(keyMsg)
		}

	case StepInitializing:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// handleEnter handles the enter key press for each step
func (m StartWizardModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case StepWelcome:
		m.step = StepProjectName
		return m, m.projectName.Focus()

	case StepProjectName:
		if m.projectName.IsValid() {
			m.step = StepModelSelection
		}
		return m, nil

	case StepModelSelection:
		m.step = StepConfirmation
		return m, nil

	case StepConfirmation:
		m.step = StepInitializing
		return m, tea.Batch(
			m.spinner.Init(),
			m.initializeProject(),
		)

	case StepComplete:
		if m.error == "" {
			m.step = StepNextAction
			return m, nil
		}
		m.quitting = true
		return m, tea.Quit

	case StepNextAction:
		m.selectedAction = m.nextAction.SelectedOption().Value()
		m.quitting = true
		return m, tea.Quit
	}

	return m, nil
}

// initializeProject initializes the project
func (m *StartWizardModel) initializeProject() tea.Cmd {
	return func() tea.Msg {
		// Initialize the project using the project package
		// Note: For now we use the basic InitProject function
		// In the future, we can enhance it to accept model configuration
		if err := project.InitProject(m.projectName.Value()); err != nil {
			return projectInitErrorMsg{err}
		}

		return projectInitCompleteMsg{}
	}
}

// projectInitCompleteMsg signals project initialization is complete
type projectInitCompleteMsg struct{}

// projectInitErrorMsg signals an error during initialization
type projectInitErrorMsg struct {
	err error
}

// Update handler for initialization messages
func (m StartWizardModel) handleInitMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case projectInitCompleteMsg:
		m.step = StepComplete
		return m, nil

	case projectInitErrorMsg:
		errMsg := msg.(projectInitErrorMsg)
		m.error = errMsg.err.Error()
		m.step = StepComplete
		return m, nil
	}
	return m, nil
}

// View renders the wizard
func (m StartWizardModel) View() string {
	if m.quitting && m.step != StepComplete {
		return styles.MutedStyle.Render("Cancelled.\n")
	}

	switch m.step {
	case StepWelcome:
		return m.renderWelcome()
	case StepProjectName:
		return m.renderProjectName()
	case StepModelSelection:
		return m.renderModelSelection()
	case StepConfirmation:
		return m.renderConfirmation()
	case StepInitializing:
		return m.renderInitializing()
	case StepComplete:
		return m.renderComplete()
	case StepNextAction:
		return m.renderNextAction()
	}

	return ""
}

// renderWelcome renders the welcome screen
func (m StartWizardModel) renderWelcome() string {
	var b strings.Builder

	logo := `
┌─────────────────────────────────────────┐
│                                         │
│     █████╗  ██████╗ ███████╗███╗   ██╗ │
│    ██╔══██╗██╔════╝ ██╔════╝████╗  ██║ │
│    ███████║██║  ███╗█████╗  ██╔██╗ ██║ │
│    ██╔══██║██║   ██║██╔══╝  ██║╚██╗██║ │
│    ██║  ██║╚██████╔╝███████╗██║ ╚████║ │
│    ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝  ╚═══╝ │
│                                         │
│    Agnostic Agent Framework             │
│    Specification-driven AI workflows    │
│                                         │
└─────────────────────────────────────────┘
`

	logoStyle := styles.TitleStyle.Copy().Foreground(styles.Primary)
	b.WriteString(logoStyle.Render(logo) + "\n\n")
	b.WriteString(styles.TitleStyle.Render("Welcome! Let's set up your project.") + "\n\n")
	b.WriteString(styles.HelpStyle.Render("Press Enter to continue or Esc to exit") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

// renderProjectName renders the project name input step
func (m StartWizardModel) renderProjectName() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Project Setup") + "\n\n")
	b.WriteString(m.projectName.View() + "\n")
	b.WriteString(styles.HelpStyle.Render("Press Enter to continue • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

// renderModelSelection renders the model selection step
func (m StartWizardModel) renderModelSelection() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("AI Model Configuration") + "\n\n")
	b.WriteString(m.modelSelect.View() + "\n")
	b.WriteString(styles.HelpStyle.Render("↑/↓ to navigate • Enter to continue • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

// renderConfirmation renders the confirmation step
func (m StartWizardModel) renderConfirmation() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Confirm Project Configuration") + "\n\n")

	// Show configuration summary
	config := styles.BoxStyle.Render(fmt.Sprintf(
		"Project Name: %s\nAI Model: %s\nMax Tokens: 8000",
		styles.BoldStyle.Render(m.projectName.Value()),
		styles.BoldStyle.Render(m.modelSelect.SelectedOption().Title()),
	))

	b.WriteString(config + "\n\n")

	b.WriteString(styles.SubtitleStyle.Render("This will create:") + "\n")
	b.WriteString("  • .agentic/ directory structure\n")
	b.WriteString("  • Specification templates\n")
	b.WriteString("  • Task management files\n")
	b.WriteString("  • Configuration file\n\n")

	b.WriteString(styles.HelpStyle.Render("Press Enter to initialize • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

// renderInitializing renders the initializing step
func (m StartWizardModel) renderInitializing() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Initializing Project") + "\n\n")
	b.WriteString(m.spinner.View() + "\n\n")
	b.WriteString(styles.MutedStyle.Render("Please wait...") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

// renderComplete renders the completion step
func (m StartWizardModel) renderComplete() string {
	var b strings.Builder

	if m.error != "" {
		b.WriteString(styles.RenderError("Failed to initialize project") + "\n\n")
		b.WriteString(styles.MutedStyle.Render(m.error) + "\n\n")
		b.WriteString(styles.HelpStyle.Render("Press Enter to exit") + "\n")
		return styles.ContainerStyle.Render(b.String())
	}

	b.WriteString(styles.RenderSuccess("Project initialized successfully!") + "\n\n")
	b.WriteString(styles.HelpStyle.Render("Press Enter to continue") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

// renderNextAction renders the next action selection
func (m StartWizardModel) renderNextAction() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("What would you like to do next?") + "\n\n")
	b.WriteString(m.nextAction.View() + "\n\n")
	b.WriteString(styles.HelpStyle.Render("↑/↓ to navigate • Enter to continue • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}
