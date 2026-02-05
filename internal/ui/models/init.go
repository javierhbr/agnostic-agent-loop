package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/project"
	"github.com/javierbenavides/agentic-agent/internal/ui/components"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
)

// InitWizardStep represents the current step in init wizard
type InitWizardStep int

const (
	InitStepWelcome InitWizardStep = iota
	InitStepProjectName
	InitStepModel
	InitStepValidators
	InitStepPreview
	InitStepInitializing
	InitStepComplete
)

// InitWizardModel is the model for init wizard
type InitWizardModel struct {
	step           InitWizardStep
	projectName    components.ValidatedInput
	modelSelect    components.SimpleSelect
	validatorSelect components.SimpleSelect
	spinner        components.Spinner
	error          string
	width          int
	height         int
	quitting       bool
}

// NewInitWizardModel creates a new init wizard model
func NewInitWizardModel() InitWizardModel {
	// Project name input
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
			return nil
		},
	)

	// Model selection
	modelOptions := []components.SelectOption{
		components.NewSelectOption(
			"Claude 3.5 Sonnet (Recommended)",
			"Latest Claude model - excellent balance of speed and capability",
			"claude-3-5-sonnet-20241022",
		),
		components.NewSelectOption(
			"GPT-4 Turbo",
			"OpenAI's powerful model with great reasoning",
			"gpt-4-turbo-preview",
		),
		components.NewSelectOption(
			"Claude 3 Opus",
			"Most capable Claude model for complex tasks",
			"claude-3-opus-20240229",
		),
	}
	modelSelect := components.NewSimpleSelect("AI Model", modelOptions)

	// Validator selection
	validatorOptions := []components.SelectOption{
		components.NewSelectOption(
			"All Validators (Recommended)",
			"Enable all quality checks and best practices",
			"all",
		),
		components.NewSelectOption(
			"Essential Only",
			"Only critical validators (context, scope, size)",
			"essential",
		),
		components.NewSelectOption(
			"None",
			"Disable validators (not recommended)",
			"none",
		),
	}
	validatorSelect := components.NewSimpleSelect("Validation Rules", validatorOptions)

	spinner := components.NewSpinner("Initializing project...")

	return InitWizardModel{
		step:            InitStepWelcome,
		projectName:     projectName,
		modelSelect:     modelSelect,
		validatorSelect: validatorSelect,
		spinner:         spinner,
	}
}

// Init initializes the model
func (m InitWizardModel) Init() tea.Cmd {
	return m.projectName.Focus()
}

// Update handles messages
func (m InitWizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "esc":
			if m.step != InitStepInitializing {
				m.quitting = true
				return m, tea.Quit
			}

		case "enter":
			return m.handleEnter()
		}

	case projectInitCompleteMsg:
		m.step = InitStepComplete
		return m, nil

	case projectInitErrorMsg:
		m.error = msg.err.Error()
		m.step = InitStepComplete
		return m, nil
	}

	// Handle step-specific updates
	switch m.step {
	case InitStepProjectName:
		m.projectName, cmd = m.projectName.Update(msg)
		return m, cmd

	case InitStepModel:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			m.modelSelect = m.modelSelect.Update(keyMsg)
		}

	case InitStepValidators:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			m.validatorSelect = m.validatorSelect.Update(keyMsg)
		}

	case InitStepInitializing:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// handleEnter handles the enter key press for each step
func (m InitWizardModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case InitStepWelcome:
		m.step = InitStepProjectName
		return m, m.projectName.Focus()

	case InitStepProjectName:
		if m.projectName.IsValid() {
			m.step = InitStepModel
		}
		return m, nil

	case InitStepModel:
		m.step = InitStepValidators
		return m, nil

	case InitStepValidators:
		m.step = InitStepPreview
		return m, nil

	case InitStepPreview:
		m.step = InitStepInitializing
		return m, tea.Batch(
			m.spinner.Init(),
			m.initializeProject(),
		)

	case InitStepComplete:
		m.quitting = true
		return m, tea.Quit
	}

	return m, nil
}

// initializeProject initializes the project
func (m *InitWizardModel) initializeProject() tea.Cmd {
	return func() tea.Msg {
		if err := project.InitProject(m.projectName.Value()); err != nil {
			return projectInitErrorMsg{err}
		}
		return projectInitCompleteMsg{}
	}
}

// View renders the wizard
func (m InitWizardModel) View() string {
	if m.quitting && m.step != InitStepComplete {
		return styles.MutedStyle.Render("Cancelled.\n")
	}

	switch m.step {
	case InitStepWelcome:
		return m.renderWelcome()
	case InitStepProjectName:
		return m.renderProjectName()
	case InitStepModel:
		return m.renderModel()
	case InitStepValidators:
		return m.renderValidators()
	case InitStepPreview:
		return m.renderPreview()
	case InitStepInitializing:
		return m.renderInitializing()
	case InitStepComplete:
		return m.renderComplete()
	}

	return ""
}

func (m InitWizardModel) renderWelcome() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Initialize Agentic Agent Project") + "\n\n")
	b.WriteString("This wizard will guide you through setting up a new project with:\n\n")
	b.WriteString(styles.ListItemStyle.Render("  • Project configuration\n"))
	b.WriteString(styles.ListItemStyle.Render("  • Directory structure (.agentic/)\n"))
	b.WriteString(styles.ListItemStyle.Render("  • Specification templates\n"))
	b.WriteString(styles.ListItemStyle.Render("  • Task management files\n"))
	b.WriteString(styles.ListItemStyle.Render("  • AI model configuration\n\n"))

	b.WriteString(styles.HelpStyle.Render("Press Enter to continue • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m InitWizardModel) renderProjectName() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Project Configuration") + "\n\n")
	b.WriteString(m.projectName.View() + "\n")
	b.WriteString(styles.HelpStyle.Render("Enter to continue • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m InitWizardModel) renderModel() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("AI Model Selection") + "\n\n")
	b.WriteString(m.modelSelect.View() + "\n")
	b.WriteString(styles.HelpStyle.Render("↑/↓ to navigate • Enter to continue • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m InitWizardModel) renderValidators() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Validation Rules") + "\n\n")
	b.WriteString(m.validatorSelect.View() + "\n\n")

	b.WriteString(styles.MutedStyle.Render("Validators enforce best practices:\n"))
	b.WriteString(styles.MutedStyle.Render("  • Context files in directories\n"))
	b.WriteString(styles.MutedStyle.Render("  • Task scope enforcement\n"))
	b.WriteString(styles.MutedStyle.Render("  • Task size limits\n\n"))

	b.WriteString(styles.HelpStyle.Render("↑/↓ to navigate • Enter to continue • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m InitWizardModel) renderPreview() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Confirm Configuration") + "\n\n")

	config := fmt.Sprintf(
		"Project Name: %s\n"+
			"AI Model: %s\n"+
			"Validators: %s",
		styles.BoldStyle.Render(m.projectName.Value()),
		styles.BoldStyle.Render(m.modelSelect.SelectedOption().Title()),
		styles.BoldStyle.Render(m.validatorSelect.SelectedOption().Title()),
	)

	b.WriteString(styles.CardStyle.Render(config) + "\n")

	b.WriteString(styles.SubtitleStyle.Render("Directory structure to be created:") + "\n\n")
	b.WriteString(styles.MutedStyle.Render(".agentic/\n"))
	b.WriteString(styles.MutedStyle.Render("├── spec/              # Specifications\n"))
	b.WriteString(styles.MutedStyle.Render("├── context/           # Context summaries\n"))
	b.WriteString(styles.MutedStyle.Render("├── tasks/             # Task management\n"))
	b.WriteString(styles.MutedStyle.Render("└── agent-rules/       # Tool configs\n"))
	b.WriteString(styles.MutedStyle.Render("agnostic-agent.yaml    # Project config\n\n"))

	b.WriteString(styles.HelpStyle.Render("Press Enter to initialize • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m InitWizardModel) renderInitializing() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Initializing Project") + "\n\n")
	b.WriteString(m.spinner.View() + "\n\n")
	b.WriteString(styles.MutedStyle.Render("Creating directory structure and files...") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m InitWizardModel) renderComplete() string {
	var b strings.Builder

	if m.error != "" {
		b.WriteString(styles.RenderError("Failed to initialize project") + "\n\n")
		b.WriteString(styles.MutedStyle.Render(m.error) + "\n\n")
		b.WriteString(styles.HelpStyle.Render("Press Enter to exit") + "\n")
		return styles.ContainerStyle.Render(b.String())
	}

	b.WriteString(styles.RenderSuccess("Project initialized successfully!") + "\n\n")

	b.WriteString(styles.SubtitleStyle.Render("Next steps:") + "\n\n")

	// Format each step on its own line with consistent spacing
	step1 := "1. Create your first task:\n   " + styles.BoldStyle.Render("agentic-agent task create")
	step2 := "2. Start working:\n   " + styles.BoldStyle.Render("agentic-agent work")
	step3 := "3. View all tasks:\n   " + styles.BoldStyle.Render("agentic-agent task list")

	b.WriteString(step1 + "\n\n")
	b.WriteString(step2 + "\n\n")
	b.WriteString(step3 + "\n\n")

	b.WriteString(styles.HelpStyle.Render("Press Enter to exit") + "\n")

	return styles.ContainerStyle.Render(b.String())
}
