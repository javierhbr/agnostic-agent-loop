package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/internal/ui/components"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	taskmodels "github.com/javierbenavides/agentic-agent/pkg/models"
)

// TemplateStep represents the current step in template workflow
type TemplateStep int

const (
	TemplateStepSelect TemplateStep = iota
	TemplateStepCustomize
	TemplateStepPreview
	TemplateStepCreating
	TemplateStepComplete
)

// TaskTemplate represents a saved task template
type TaskTemplate struct {
	Name        string
	Description string
	Task        taskmodels.Task
}

// TemplateSelectModel is the model for template selection
type TemplateSelectModel struct {
	step          TemplateStep
	templates     []TaskTemplate
	templateSelect components.SimpleSelect
	title         components.ValidatedInput
	description   components.TextArea
	selectedTemplate *TaskTemplate
	spinner       components.Spinner
	error         string
	taskID        string
	width         int
	height        int
	quitting      bool
}

// NewTemplateSelectModel creates a new template selection model
func NewTemplateSelectModel() TemplateSelectModel {
	// Load templates
	templates := loadBuiltInTemplates()

	// Create template selector
	options := make([]components.SelectOption, len(templates))
	for i, tmpl := range templates {
		options[i] = components.NewSelectOption(
			tmpl.Name,
			tmpl.Description,
			fmt.Sprintf("template-%d", i),
		)
	}

	templateSelect := components.NewSimpleSelect("Select Template", options)

	// Title input
	title := components.NewValidatedInput(
		"Task Title",
		"",
		func(s string) error {
			if len(s) == 0 {
				return fmt.Errorf("title cannot be empty")
			}
			return nil
		},
	)

	// Description
	description := components.NewTextArea("Task Description", "", true)

	spinner := components.NewSpinner("Creating task...")

	return TemplateSelectModel{
		step:           TemplateStepSelect,
		templates:      templates,
		templateSelect: templateSelect,
		title:          title,
		description:    description,
		spinner:        spinner,
	}
}

// loadBuiltInTemplates loads the built-in task templates
func loadBuiltInTemplates() []TaskTemplate {
	return []TaskTemplate{
		{
			Name:        "Feature Implementation",
			Description: "Template for implementing a new feature",
			Task: taskmodels.Task{
				Title:       "Implement [Feature Name]",
				Description: "Add a new feature to the application.\n\nDetails:\n- [Add details here]",
				Acceptance: []string{
					"Feature is implemented according to specification",
					"Unit tests cover the new feature",
					"Integration tests pass",
					"Documentation is updated",
				},
				Scope: []string{
					"internal/",
					"cmd/",
				},
			},
		},
		{
			Name:        "Bug Fix",
			Description: "Template for fixing bugs",
			Task: taskmodels.Task{
				Title:       "Fix: [Bug Description]",
				Description: "Fix the following bug:\n\n**Issue:** [Describe the bug]\n**Expected:** [Expected behavior]\n**Actual:** [Actual behavior]",
				Acceptance: []string{
					"Bug is fixed and verified",
					"Root cause is identified and documented",
					"Regression tests added",
					"No new bugs introduced",
				},
				Scope: []string{},
			},
		},
		{
			Name:        "Refactoring",
			Description: "Template for code refactoring",
			Task: taskmodels.Task{
				Title:       "Refactor: [Component Name]",
				Description: "Refactor code to improve:\n- Code quality\n- Maintainability\n- Performance\n- Test coverage",
				Acceptance: []string{
					"Code is refactored without breaking functionality",
					"All tests pass",
					"Code quality metrics improved",
					"Documentation updated if needed",
				},
				Scope: []string{},
			},
		},
		{
			Name:        "Documentation",
			Description: "Template for documentation tasks",
			Task: taskmodels.Task{
				Title:       "Document: [Component/Feature]",
				Description: "Create or update documentation for:\n\n- [Item 1]\n- [Item 2]",
				Acceptance: []string{
					"Documentation is clear and accurate",
					"Code examples are provided",
					"Common use cases are covered",
					"Documentation is reviewed",
				},
				Scope: []string{
					"docs/",
					"README.md",
				},
			},
		},
		{
			Name:        "Testing",
			Description: "Template for adding tests",
			Task: taskmodels.Task{
				Title:       "Add tests for [Component]",
				Description: "Add comprehensive tests for:\n\n- Unit tests\n- Integration tests\n- Edge cases",
				Acceptance: []string{
					"Test coverage increased",
					"All edge cases covered",
					"Tests are maintainable",
					"CI pipeline passes",
				},
				Scope: []string{},
			},
		},
	}
}

// Init initializes the model
func (m TemplateSelectModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m TemplateSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.step != TemplateStepCreating {
				m.quitting = true
				return m, tea.Quit
			}

		case "enter":
			return m.handleEnter()
		}

	case taskCreateCompleteMsg:
		m.taskID = msg.taskID
		m.step = TemplateStepComplete
		return m, nil

	case taskCreateErrorMsg:
		m.error = msg.err.Error()
		m.step = TemplateStepComplete
		return m, nil
	}

	// Handle step-specific updates
	switch m.step {
	case TemplateStepSelect:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			m.templateSelect = m.templateSelect.Update(keyMsg)
		}

	case TemplateStepCustomize:
		m.title, cmd = m.title.Update(msg)
		return m, cmd

	case TemplateStepCreating:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// handleEnter handles the enter key for each step
func (m TemplateSelectModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case TemplateStepSelect:
		// Get selected template
		selectedIdx := m.templateSelect.SelectedIdx
		if selectedIdx >= 0 && selectedIdx < len(m.templates) {
			m.selectedTemplate = &m.templates[selectedIdx]
			// Pre-fill title with template title
			m.title = components.NewValidatedInput(
				"Task Title",
				m.selectedTemplate.Task.Title,
				func(s string) error {
					if len(s) == 0 {
						return fmt.Errorf("title cannot be empty")
					}
					return nil
				},
			)
			m.description = components.NewTextArea("Task Description", m.selectedTemplate.Task.Description, true)
			m.step = TemplateStepCustomize
			return m, m.title.Focus()
		}

	case TemplateStepCustomize:
		if m.title.IsValid() {
			m.step = TemplateStepPreview
		}

	case TemplateStepPreview:
		m.step = TemplateStepCreating
		return m, tea.Batch(
			m.spinner.Init(),
			m.createTaskFromTemplate(),
		)

	case TemplateStepComplete:
		m.quitting = true
		return m, tea.Quit
	}

	return m, nil
}

// createTaskFromTemplate creates a task from the selected template
func (m *TemplateSelectModel) createTaskFromTemplate() tea.Cmd {
	return func() tea.Msg {
		tm := tasks.NewTaskManager(".agentic/tasks")

		// Create task from template
		task, err := tm.CreateTask(m.title.Value())
		if err != nil {
			return taskCreateErrorMsg{err}
		}

		// Copy template data
		if m.selectedTemplate != nil {
			task.Description = strings.TrimSpace(m.description.Value())
			task.Acceptance = m.selectedTemplate.Task.Acceptance
			task.Scope = m.selectedTemplate.Task.Scope
		}

		// Save task
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

// View renders the model
func (m TemplateSelectModel) View() string {
	if m.quitting && m.step != TemplateStepComplete {
		return styles.MutedStyle.Render("Cancelled.\n")
	}

	switch m.step {
	case TemplateStepSelect:
		return m.renderSelect()
	case TemplateStepCustomize:
		return m.renderCustomize()
	case TemplateStepPreview:
		return m.renderPreview()
	case TemplateStepCreating:
		return m.renderCreating()
	case TemplateStepComplete:
		return m.renderComplete()
	}

	return ""
}

func (m TemplateSelectModel) renderSelect() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Select Task Template") + "\n\n")
	b.WriteString(m.templateSelect.View() + "\n\n")
	b.WriteString(styles.MutedStyle.Render("Choose a template to start with, then customize it for your needs.") + "\n\n")
	b.WriteString(styles.HelpStyle.Render("↑/↓ to navigate • Enter to continue • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m TemplateSelectModel) renderCustomize() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Customize Template") + "\n\n")

	if m.selectedTemplate != nil {
		b.WriteString(styles.SubtitleStyle.Render("Template: "+m.selectedTemplate.Name) + "\n\n")
	}

	b.WriteString(m.title.View() + "\n\n")
	b.WriteString(styles.MutedStyle.Render("You can edit the title and description. Other fields will be copied from the template.") + "\n\n")
	b.WriteString(styles.HelpStyle.Render("Enter to continue • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m TemplateSelectModel) renderPreview() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Task Preview") + "\n\n")

	if m.selectedTemplate != nil {
		// Build preview
		preview := fmt.Sprintf("Title: %s\n", styles.BoldStyle.Render(m.title.Value()))
		preview += fmt.Sprintf("\nDescription:\n%s\n", strings.TrimSpace(m.description.Value()))

		if len(m.selectedTemplate.Task.Acceptance) > 0 {
			preview += "\nAcceptance Criteria:\n"
			for _, criterion := range m.selectedTemplate.Task.Acceptance {
				preview += fmt.Sprintf("  %s %s\n", styles.IconBullet, criterion)
			}
		}

		if len(m.selectedTemplate.Task.Scope) > 0 {
			preview += "\nScope:\n"
			for _, scope := range m.selectedTemplate.Task.Scope {
				preview += fmt.Sprintf("  %s %s\n", styles.IconBullet, scope)
			}
		}

		b.WriteString(styles.CardStyle.Render(preview) + "\n")
	}

	b.WriteString(styles.HelpStyle.Render("Press Enter to create task • Esc to cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m TemplateSelectModel) renderCreating() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Creating Task") + "\n\n")
	b.WriteString(m.spinner.View() + "\n\n")
	b.WriteString(styles.MutedStyle.Render("Creating task from template...") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m TemplateSelectModel) renderComplete() string {
	var b strings.Builder

	if m.error != "" {
		b.WriteString(styles.RenderError("Failed to create task") + "\n\n")
		b.WriteString(styles.MutedStyle.Render(m.error) + "\n\n")
		b.WriteString(styles.HelpStyle.Render("Press Enter to exit") + "\n")
		return styles.ContainerStyle.Render(b.String())
	}

	b.WriteString(styles.RenderSuccess("Task created successfully from template!") + "\n\n")

	// Next steps
	step1 := "1. View task:\n   " + styles.BoldStyle.Render("agentic-agent task show "+m.taskID)
	step2 := "2. Claim task:\n   " + styles.BoldStyle.Render("agentic-agent task claim "+m.taskID)
	step3 := "3. List all tasks:\n   " + styles.BoldStyle.Render("agentic-agent task list")

	b.WriteString(step1 + "\n\n")
	b.WriteString(step2 + "\n\n")
	b.WriteString(step3 + "\n\n")

	b.WriteString(styles.HelpStyle.Render("Press Enter to exit") + "\n")

	return styles.ContainerStyle.Render(b.String())
}
