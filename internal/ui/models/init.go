package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/project"
	"github.com/javierbenavides/agentic-agent/internal/skills"
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
	InitStepDescription
	InitStepTechStack
	InitStepWorkflow
	InitStepAgentTools
	InitStepPreview
	InitStepInitializing
	InitStepComplete
)

// InitWizardModel is the model for init wizard
type InitWizardModel struct {
	step            InitWizardStep
	projectName     components.ValidatedInput
	modelSelect     components.SimpleSelect
	validatorSelect components.SimpleSelect
	agentToolSelect components.SimpleSelect
	description     components.TextArea
	techStack       components.TextArea
	workflow        components.TextArea
	spinner         components.Spinner
	error           string
	width           int
	height          int
	quitting        bool
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

	// Project profile fields (all optional)
	description := components.NewTextArea(
		"Project Description",
		"Brief description of your project, its goals, and key features...",
		true,
	)
	techStack := components.NewTextArea(
		"Tech Stack",
		"Languages, frameworks, databases (e.g., Go, React, PostgreSQL)...",
		true,
	)
	workflow := components.NewTextArea(
		"Workflow Preferences",
		"Testing approach, branching strategy, code review process...",
		true,
	)

	// Agent tool selection
	agentToolOptions := []components.SelectOption{
		components.NewSelectOption("Skip", "Don't generate agent rules now", "skip"),
	}
	for _, tool := range skills.SupportedTools() {
		agentToolOptions = append(agentToolOptions, components.NewSelectOption(
			tool,
			fmt.Sprintf("Generate rules for %s", tool),
			tool,
		))
	}
	agentToolSelect := components.NewSimpleSelect("Agent Tool", agentToolOptions)

	spinner := components.NewSpinner("Initializing project...")

	return InitWizardModel{
		step:            InitStepWelcome,
		projectName:     projectName,
		modelSelect:     modelSelect,
		validatorSelect: validatorSelect,
		agentToolSelect: agentToolSelect,
		description:     description,
		techStack:       techStack,
		workflow:        workflow,
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

		case "tab":
			// Tab skips optional profile steps
			switch m.step {
			case InitStepDescription:
				m.step = InitStepTechStack
				return m, m.techStack.Focus()
			case InitStepTechStack:
				m.step = InitStepWorkflow
				return m, m.workflow.Focus()
			case InitStepWorkflow:
				m.step = InitStepAgentTools
				return m, nil
			}

		case "enter":
			// In textarea steps, enter adds a newline (use tab to skip, ctrl+n to advance)
			if m.step == InitStepDescription || m.step == InitStepTechStack || m.step == InitStepWorkflow {
				// Don't intercept enter for textareas - let it add newlines
				break
			}
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

	case InitStepDescription:
		m.description, cmd = m.description.Update(msg)
		return m, cmd

	case InitStepTechStack:
		m.techStack, cmd = m.techStack.Update(msg)
		return m, cmd

	case InitStepWorkflow:
		m.workflow, cmd = m.workflow.Update(msg)
		return m, cmd

	case InitStepAgentTools:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			m.agentToolSelect = m.agentToolSelect.Update(keyMsg)
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
		m.step = InitStepDescription
		return m, m.description.Focus()

	case InitStepAgentTools:
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
		profile := &project.ProjectProfile{
			Description: m.description.Value(),
			TechStack:   m.techStack.Value(),
			Workflow:    m.workflow.Value(),
		}
		// Only pass profile if any field is filled
		if profile.Description == "" && profile.TechStack == "" && profile.Workflow == "" {
			profile = nil
		}
		if err := project.InitProjectWithProfile(m.projectName.Value(), profile); err != nil {
			return projectInitErrorMsg{err}
		}

		// Generate agent skills if a tool was selected
		agentTool := m.agentToolSelect.SelectedOption().Value()
		if agentTool != "" && agentTool != "skip" {
			// Best-effort: don't fail init if skills generation fails
			skills.Ensure(agentTool, nil)
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
	case InitStepDescription:
		return m.renderDescription()
	case InitStepTechStack:
		return m.renderTechStack()
	case InitStepWorkflow:
		return m.renderWorkflow()
	case InitStepAgentTools:
		return m.renderAgentTools()
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
	b.WriteString("  • Project configuration\n")
	b.WriteString("  • Directory structure (.agentic/)\n")
	b.WriteString("  • Specification templates\n")
	b.WriteString("  • Task management files\n")
	b.WriteString("  • AI model configuration\n")
	b.WriteString("  • Project profile (optional)\n\n")

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

	b.WriteString(styles.MutedStyle.Render("Validators enforce best practices:\n"+
		"  • Context files in directories\n"+
		"  • Task scope enforcement\n"+
		"  • Task size limits") + "\n\n")

	b.WriteString(styles.HelpStyle.Render("↑/↓ to navigate • Enter to continue • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m InitWizardModel) renderDescription() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Project Profile") + "\n\n")
	b.WriteString(m.description.View() + "\n")
	b.WriteString(styles.HelpStyle.Render("Tab to skip/next • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m InitWizardModel) renderTechStack() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Project Profile") + "\n\n")
	b.WriteString(m.techStack.View() + "\n")
	b.WriteString(styles.HelpStyle.Render("Tab to skip/next • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m InitWizardModel) renderWorkflow() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Project Profile") + "\n\n")
	b.WriteString(m.workflow.View() + "\n")
	b.WriteString(styles.HelpStyle.Render("Tab to skip/next • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m InitWizardModel) renderAgentTools() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Agent Tool Setup") + "\n\n")
	b.WriteString(styles.MutedStyle.Render("Select the AI agent you'll use with this project.\n") + "\n")
	b.WriteString(styles.MutedStyle.Render("This generates tool-specific rules (e.g., CLAUDE.md, .cursor/rules).\n\n"))
	b.WriteString(m.agentToolSelect.View() + "\n")
	b.WriteString(styles.HelpStyle.Render("↑/↓ to navigate • Enter to continue • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m InitWizardModel) renderPreview() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Confirm Configuration") + "\n\n")

	agentToolLabel := m.agentToolSelect.SelectedOption().Title()
	if m.agentToolSelect.SelectedOption().Value() == "skip" {
		agentToolLabel = "None"
	}

	config := fmt.Sprintf(
		"Project Name: %s\n"+
			"AI Model: %s\n"+
			"Validators: %s\n"+
			"Agent Tool: %s",
		styles.BoldStyle.Render(m.projectName.Value()),
		styles.BoldStyle.Render(m.modelSelect.SelectedOption().Title()),
		styles.BoldStyle.Render(m.validatorSelect.SelectedOption().Title()),
		styles.BoldStyle.Render(agentToolLabel),
	)

	b.WriteString(styles.CardStyle.Render(config) + "\n")

	// Show profile summary if any filled
	if m.description.Value() != "" || m.techStack.Value() != "" || m.workflow.Value() != "" {
		var profile strings.Builder
		if m.description.Value() != "" {
			profile.WriteString("Description: provided\n")
		}
		if m.techStack.Value() != "" {
			profile.WriteString("Tech Stack: provided\n")
		}
		if m.workflow.Value() != "" {
			profile.WriteString("Workflow: provided\n")
		}
		b.WriteString(styles.CardStyle.Render(profile.String()) + "\n")
	}

	b.WriteString(styles.SubtitleStyle.Render("Directory structure to be created:") + "\n\n")
	b.WriteString(styles.MutedStyle.Render(".agentic/\n"+
		"├── spec/              # Specifications\n"+
		"├── context/           # Context summaries\n"+
		"├── tasks/             # Task management\n"+
		"├── tracks/            # Feature/bug tracks\n"+
		"└── agent-rules/       # Tool configs\n"+
		"agnostic-agent.yaml    # Project config") + "\n\n")

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

	b.WriteString(styles.SubtitleStyle.Render("Recommended workflow (confirm each step with your AI agent):") + "\n\n")
	b.WriteString("  1. " + styles.BoldStyle.Render("Brainstorm") + "       — Explore and refine your idea\n")
	b.WriteString("  2. " + styles.BoldStyle.Render("Product PRD") + "      — Formalize into a requirements doc (product-wizard)\n")
	b.WriteString("  3. " + styles.BoldStyle.Render("OpenSpec") + "         — " + styles.BoldStyle.Render("agentic-agent openspec init \"feature\" --from <prd>") + "\n")
	b.WriteString("                       Creates proposal, dev plan, and tasks automatically\n\n")

	b.WriteString(styles.SubtitleStyle.Render("Or start working directly:") + "\n\n")

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
