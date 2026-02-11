package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/internal/ui/components"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
)

// TaskCreateStep represents the current step in task creation
type TaskCreateStep int

const (
	TaskStepTitle TaskCreateStep = iota
	TaskStepDescription
	TaskStepSpecRefsConfirm
	TaskStepSpecRefsPicker
	TaskStepScopeConfirm
	TaskStepScopePicker
	TaskStepOutputsConfirm
	TaskStepOutputsPicker
	TaskStepAcceptanceConfirm
	TaskStepAcceptance
	TaskStepPreview
	TaskStepCreating
	TaskStepComplete
)

// TaskCreateModel is the model for task creation wizard
type TaskCreateModel struct {
	step            TaskCreateStep
	title           components.ValidatedInput
	description     components.TextArea
	addSpecRefs     components.Confirm
	specRefsPicker  components.FilePicker
	addScope        components.Confirm
	scopePicker     components.FilePicker
	addOutputs      components.Confirm
	outputsPicker   components.FilePicker
	addAcceptance   components.Confirm
	acceptance      components.MultiItemEditor
	error           string
	taskID          string
	spinner         components.Spinner
	width           int
	height          int
	quitting        bool
	selectedSpecRefs []string
	selectedScope    []string
	selectedOutputs  []string
}

// NewTaskCreateModel creates a new task creation model
func NewTaskCreateModel() TaskCreateModel {
	// Title input with validation
	title := components.NewValidatedInput(
		"Task Title",
		"Implement user authentication",
		func(s string) error {
			if len(s) == 0 {
				return fmt.Errorf("title cannot be empty")
			}
			if len(s) > 200 {
				return fmt.Errorf("title too long (max 200 characters)")
			}
			if strings.Contains(s, "\n") || strings.Contains(s, "\r") {
				return fmt.Errorf("title cannot contain newlines")
			}
			return nil
		},
	)

	// Description textarea
	description := components.NewTextArea(
		"Description",
		"Detailed description of the task...",
		true,
	)

	// Confirmation prompts
	addSpecRefs := components.NewConfirm("Add specification references?", false)
	addScope := components.NewConfirm("Add scope (files/directories)?", false)
	addOutputs := components.NewConfirm("Add expected output files?", false)
	addAcceptance := components.NewConfirm("Add acceptance criteria?", true)

	// Acceptance criteria editor
	acceptance := components.NewMultiItemEditor("Acceptance Criteria")

	spinner := components.NewSpinner("Creating task...")

	return TaskCreateModel{
		step:          TaskStepTitle,
		title:         title,
		description:   description,
		addSpecRefs:   addSpecRefs,
		addScope:      addScope,
		addOutputs:    addOutputs,
		addAcceptance: addAcceptance,
		acceptance:    acceptance,
		spinner:       spinner,
	}
}

// Init initializes the model
func (m TaskCreateModel) Init() tea.Cmd {
	return m.title.Focus()
}

// Update handles messages
func (m TaskCreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.step != TaskStepCreating && !m.acceptance.IsEditing {
				m.quitting = true
				return m, tea.Quit
			}

		case "enter":
			// Handle enter based on current step
			if !m.acceptance.IsEditing {
				return m.handleEnter()
			}

		case "a":
			// Add acceptance criteria
			if m.step == TaskStepAcceptance && !m.acceptance.IsEditing {
				return m, m.acceptance.StartEditing()
			}
		}

	case taskCreateCompleteMsg:
		m.taskID = msg.taskID
		m.step = TaskStepComplete
		return m, nil

	case taskCreateErrorMsg:
		m.error = msg.err.Error()
		m.step = TaskStepComplete
		return m, nil
	}

	// Handle step-specific updates
	switch m.step {
	case TaskStepTitle:
		m.title, cmd = m.title.Update(msg)
		return m, cmd

	case TaskStepDescription:
		m.description, cmd = m.description.Update(msg)
		return m, cmd

	case TaskStepSpecRefsConfirm:
		m.addSpecRefs = m.addSpecRefs.Update(msg)

	case TaskStepSpecRefsPicker:
		m.specRefsPicker, cmd = m.specRefsPicker.Update(msg)
		return m, cmd

	case TaskStepScopeConfirm:
		m.addScope = m.addScope.Update(msg)

	case TaskStepScopePicker:
		m.scopePicker, cmd = m.scopePicker.Update(msg)
		return m, cmd

	case TaskStepOutputsConfirm:
		m.addOutputs = m.addOutputs.Update(msg)

	case TaskStepOutputsPicker:
		m.outputsPicker, cmd = m.outputsPicker.Update(msg)
		return m, cmd

	case TaskStepAcceptanceConfirm:
		m.addAcceptance = m.addAcceptance.Update(msg)

	case TaskStepAcceptance:
		if m.addAcceptance.IsYes() {
			m.acceptance, cmd = m.acceptance.Update(msg)
			return m, cmd
		}

	case TaskStepCreating:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// handleEnter handles the enter key press for each step
func (m TaskCreateModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case TaskStepTitle:
		if m.title.IsValid() {
			m.step = TaskStepDescription
			return m, m.description.Focus()
		}

	case TaskStepDescription:
		m.step = TaskStepSpecRefsConfirm
		m.description.Blur()

	case TaskStepSpecRefsConfirm:
		if m.addSpecRefs.IsYes() {
			// Initialize file picker for spec refs
			m.specRefsPicker = components.NewFilePicker("Select Specification References", ".agentic/spec", false, true)
			m.step = TaskStepSpecRefsPicker
		} else {
			m.step = TaskStepScopeConfirm
		}

	case TaskStepSpecRefsPicker:
		// Save selected spec refs
		m.selectedSpecRefs = m.specRefsPicker.GetSelected()
		m.step = TaskStepScopeConfirm

	case TaskStepScopeConfirm:
		if m.addScope.IsYes() {
			// Initialize file picker for scope (directories/files)
			m.scopePicker = components.NewFilePicker("Select Scope (files/directories)", ".", false, true)
			m.step = TaskStepScopePicker
		} else {
			m.step = TaskStepOutputsConfirm
		}

	case TaskStepScopePicker:
		// Save selected scope
		m.selectedScope = m.scopePicker.GetSelected()
		m.step = TaskStepOutputsConfirm

	case TaskStepOutputsConfirm:
		if m.addOutputs.IsYes() {
			// Initialize file picker for outputs
			m.outputsPicker = components.NewFilePicker("Select Expected Output Files", ".", false, true)
			m.step = TaskStepOutputsPicker
		} else {
			m.step = TaskStepAcceptanceConfirm
		}

	case TaskStepOutputsPicker:
		// Save selected outputs
		m.selectedOutputs = m.outputsPicker.GetSelected()
		m.step = TaskStepAcceptanceConfirm

	case TaskStepAcceptanceConfirm:
		if m.addAcceptance.IsYes() {
			m.step = TaskStepAcceptance
		} else {
			m.step = TaskStepPreview
		}

	case TaskStepAcceptance:
		if !m.acceptance.IsEditing {
			m.step = TaskStepPreview
		}

	case TaskStepPreview:
		m.step = TaskStepCreating
		return m, tea.Batch(
			m.spinner.Init(),
			m.createTask(),
		)

	case TaskStepComplete:
		m.quitting = true
		return m, tea.Quit
	}

	return m, nil
}

// createTask creates the task
func (m *TaskCreateModel) createTask() tea.Cmd {
	return func() tea.Msg {
		tm := tasks.NewTaskManager(".agentic/tasks")
		task, err := tm.CreateTask(m.title.Value())
		if err != nil {
			return taskCreateErrorMsg{err}
		}

		// Update task with additional fields
		if desc := strings.TrimSpace(m.description.Value()); desc != "" {
			task.Description = desc
		}

		if len(m.selectedSpecRefs) > 0 {
			task.SpecRefs = m.selectedSpecRefs
		}

		if len(m.selectedScope) > 0 {
			task.Scope = m.selectedScope
		}

		if len(m.selectedOutputs) > 0 {
			task.Outputs = m.selectedOutputs
		}

		if m.acceptance.HasItems() {
			task.Acceptance = m.acceptance.GetItems()
		}

		// Save updated task
		backlog, err := tm.LoadTasks("backlog")
		if err != nil {
			return taskCreateErrorMsg{err}
		}

		// Update the task in backlog
		for i, t := range backlog.Tasks {
			if t.ID == task.ID {
				backlog.Tasks[i] = *task
				break
			}
		}

		if err := tm.SaveTasks("backlog", backlog); err != nil {
			return taskCreateErrorMsg{err}
		}

		return taskCreateCompleteMsg{taskID: task.ID}
	}
}

// taskCreateCompleteMsg signals task creation is complete
type taskCreateCompleteMsg struct {
	taskID string
}

// taskCreateErrorMsg signals an error during creation
type taskCreateErrorMsg struct {
	err error
}

// View renders the wizard
func (m TaskCreateModel) View() string {
	if m.quitting && m.step != TaskStepComplete {
		return styles.MutedStyle.Render("Cancelled.\n")
	}

	switch m.step {
	case TaskStepTitle:
		return m.renderTitle()
	case TaskStepDescription:
		return m.renderDescription()
	case TaskStepSpecRefsConfirm:
		return m.renderSpecRefsConfirm()
	case TaskStepSpecRefsPicker:
		return m.renderSpecRefsPicker()
	case TaskStepScopeConfirm:
		return m.renderScopeConfirm()
	case TaskStepScopePicker:
		return m.renderScopePicker()
	case TaskStepOutputsConfirm:
		return m.renderOutputsConfirm()
	case TaskStepOutputsPicker:
		return m.renderOutputsPicker()
	case TaskStepAcceptanceConfirm:
		return m.renderAcceptanceConfirm()
	case TaskStepAcceptance:
		return m.renderAcceptance()
	case TaskStepPreview:
		return m.renderPreview()
	case TaskStepCreating:
		return m.renderCreating()
	case TaskStepComplete:
		return m.renderComplete()
	}

	return ""
}

func (m TaskCreateModel) renderTitle() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Create New Task") + "\n\n")
	b.WriteString(m.title.View() + "\n")
	b.WriteString(styles.HelpStyle.Render("Enter to continue • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m TaskCreateModel) renderDescription() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Task Description") + "\n\n")
	b.WriteString(m.description.View() + "\n")
	b.WriteString(styles.HelpStyle.Render("Enter to continue • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m TaskCreateModel) renderSpecRefsConfirm() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Specification References") + "\n\n")
	b.WriteString(m.addSpecRefs.View() + "\n\n")
	b.WriteString(styles.MutedStyle.Render("Select specification files from .agentic/spec/ directory") + "\n\n")
	b.WriteString(styles.HelpStyle.Render("Enter to continue • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m TaskCreateModel) renderSpecRefsPicker() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Select Specification References") + "\n\n")
	b.WriteString(m.specRefsPicker.View() + "\n\n")
	b.WriteString(styles.HelpStyle.Render("Enter when done • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m TaskCreateModel) renderScopeConfirm() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Task Scope") + "\n\n")
	b.WriteString(m.addScope.View() + "\n\n")
	b.WriteString(styles.MutedStyle.Render("Select files or directories that this task will modify") + "\n\n")
	b.WriteString(styles.HelpStyle.Render("Enter to continue • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m TaskCreateModel) renderScopePicker() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Select Task Scope") + "\n\n")
	b.WriteString(m.scopePicker.View() + "\n\n")
	b.WriteString(styles.HelpStyle.Render("Enter when done • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m TaskCreateModel) renderOutputsConfirm() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Expected Outputs") + "\n\n")
	b.WriteString(m.addOutputs.View() + "\n\n")
	b.WriteString(styles.MutedStyle.Render("Select files that will be created or modified by this task") + "\n\n")
	b.WriteString(styles.HelpStyle.Render("Enter to continue • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m TaskCreateModel) renderOutputsPicker() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Select Expected Output Files") + "\n\n")
	b.WriteString(m.outputsPicker.View() + "\n\n")
	b.WriteString(styles.HelpStyle.Render("Enter when done • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m TaskCreateModel) renderAcceptanceConfirm() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Acceptance Criteria") + "\n\n")
	b.WriteString(m.addAcceptance.View() + "\n\n")
	b.WriteString(styles.MutedStyle.Render("Define conditions that must be met for this task to be complete") + "\n\n")
	b.WriteString(styles.HelpStyle.Render("Enter to continue • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m TaskCreateModel) renderAcceptance() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Acceptance Criteria") + "\n\n")
	b.WriteString(m.acceptance.View() + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m TaskCreateModel) renderPreview() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Task Preview") + "\n\n")

	// Build task summary
	summary := fmt.Sprintf("Title: %s\n", styles.BoldStyle.Render(m.title.Value()))

	if desc := strings.TrimSpace(m.description.Value()); desc != "" {
		summary += fmt.Sprintf("Description: %s\n", desc)
	}

	if len(m.selectedSpecRefs) > 0 {
		summary += "\nSpecification References:\n"
		for _, ref := range m.selectedSpecRefs {
			summary += fmt.Sprintf("  %s %s\n", styles.IconBullet, ref)
		}
	}

	if len(m.selectedScope) > 0 {
		summary += "\nScope:\n"
		for _, scope := range m.selectedScope {
			summary += fmt.Sprintf("  %s %s\n", styles.IconBullet, scope)
		}
	}

	if len(m.selectedOutputs) > 0 {
		summary += "\nExpected Outputs:\n"
		for _, output := range m.selectedOutputs {
			summary += fmt.Sprintf("  %s %s\n", styles.IconBullet, output)
		}
	}

	if m.acceptance.HasItems() {
		summary += "\nAcceptance Criteria:\n"
		for _, item := range m.acceptance.GetItems() {
			summary += fmt.Sprintf("  %s %s\n", styles.IconBullet, item)
		}
	}

	b.WriteString(styles.CardStyle.Render(summary) + "\n")
	b.WriteString(styles.HelpStyle.Render("Press Enter to create task • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m TaskCreateModel) renderCreating() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Creating Task") + "\n\n")
	b.WriteString(m.spinner.View() + "\n\n")
	b.WriteString(styles.MutedStyle.Render("Please wait...") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m TaskCreateModel) renderComplete() string {
	var b strings.Builder

	if m.error != "" {
		b.WriteString(styles.RenderError("Failed to create task") + "\n\n")
		b.WriteString(styles.MutedStyle.Render(m.error) + "\n\n")
		b.WriteString(styles.HelpStyle.Render("Press Enter to exit") + "\n")
		return styles.ContainerStyle.Render(b.String())
	}

	b.WriteString(styles.RenderSuccess("Task created successfully!") + "\n\n")

	b.WriteString(styles.SubtitleStyle.Render(fmt.Sprintf("Task ID: %s", m.taskID)) + "\n\n")

	b.WriteString(styles.SubtitleStyle.Render("Next steps:") + "\n\n")
	b.WriteString(fmt.Sprintf("  1. View task: %s\n", styles.BoldStyle.Render("agentic-agent task show "+m.taskID)))
	b.WriteString(fmt.Sprintf("  2. Claim task: %s\n", styles.BoldStyle.Render("agentic-agent task claim "+m.taskID)))
	b.WriteString(fmt.Sprintf("  3. List all tasks: %s\n\n", styles.BoldStyle.Render("agentic-agent task list")))

	b.WriteString(styles.HelpStyle.Render("Press Enter to exit") + "\n")

	return styles.ContainerStyle.Render(b.String())
}
