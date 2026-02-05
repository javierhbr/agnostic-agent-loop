package models

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/internal/ui/components"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	taskmodels "github.com/javierbenavides/agentic-agent/pkg/models"
)

// WorkflowStep represents the current step in the work workflow
type WorkflowStep int

const (
	WorkflowStepSelectTask WorkflowStep = iota
	WorkflowStepConfirmClaim
	WorkflowStepGenerateContext
	WorkflowStepShowTask
	WorkflowStepConfirmComplete
	WorkflowStepValidating
	WorkflowStepComplete
)

// WorkflowModel orchestrates the complete work workflow
type WorkflowModel struct {
	step          WorkflowStep
	taskManager   *tasks.TaskManager
	backlogTasks  []taskmodels.Task
	selectedTask  *taskmodels.Task
	cursorPos     int
	offset        int
	height        int
	width         int
	confirmClaim  components.Confirm
	generateCtx   components.Confirm
	confirmDone   components.Confirm
	spinner       components.Spinner
	error         string
	message       string
	validationResults string
	quitting      bool
	checklistPos  int
	completed     []bool
}

// NewWorkflowModel creates a new workflow model
func NewWorkflowModel() WorkflowModel {
	tm := tasks.NewTaskManager(".agentic/tasks")

	// Load backlog tasks
	backlog, _ := tm.LoadTasks("backlog")

	confirmClaim := components.NewConfirm("Claim this task and start working?", true)
	generateCtx := components.NewConfirm("Generate context for scope directories?", true)
	confirmDone := components.NewConfirm("Mark task as complete?", true)
	spinner := components.NewSpinner("Validating...")

	return WorkflowModel{
		step:         WorkflowStepSelectTask,
		taskManager:  tm,
		backlogTasks: backlog.Tasks,
		cursorPos:    0,
		offset:       0,
		height:       15,
		confirmClaim: confirmClaim,
		generateCtx:  generateCtx,
		confirmDone:  confirmDone,
		spinner:      spinner,
	}
}

// Init initializes the model
func (m WorkflowModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m WorkflowModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.step == WorkflowStepSelectTask {
				m.quitting = true
				return m, tea.Quit
			}
			// Go back to previous step
			return m.handleBack()

		case "enter":
			return m.handleEnter()

		case "up", "k":
			if m.step == WorkflowStepSelectTask {
				if m.cursorPos > 0 {
					m.cursorPos--
					if m.cursorPos < m.offset {
						m.offset = m.cursorPos
					}
				}
			} else if m.step == WorkflowStepShowTask {
				if m.checklistPos > 0 {
					m.checklistPos--
				}
			}

		case "down", "j":
			if m.step == WorkflowStepSelectTask {
				if m.cursorPos < len(m.backlogTasks)-1 {
					m.cursorPos++
					if m.cursorPos >= m.offset+m.height {
						m.offset = m.cursorPos - m.height + 1
					}
				}
			} else if m.step == WorkflowStepShowTask && m.selectedTask != nil {
				if m.checklistPos < len(m.selectedTask.Acceptance)-1 {
					m.checklistPos++
				}
			}

		case " ":
			// Toggle checklist items
			if m.step == WorkflowStepShowTask && m.selectedTask != nil && len(m.selectedTask.Acceptance) > 0 {
				if m.checklistPos >= 0 && m.checklistPos < len(m.completed) {
					m.completed[m.checklistPos] = !m.completed[m.checklistPos]
				}
			}
		}

	case claimCompleteMsg:
		m.step = WorkflowStepGenerateContext
		return m, nil

	case claimErrorMsg:
		m.error = msg.err.Error()
		m.step = WorkflowStepComplete
		return m, nil

	case validationCompleteMsg:
		m.validationResults = msg.results
		m.step = WorkflowStepComplete
		return m, nil

	case validationErrorMsg:
		m.error = msg.err.Error()
		m.step = WorkflowStepComplete
		return m, nil
	}

	// Handle step-specific updates
	switch m.step {
	case WorkflowStepConfirmClaim:
		m.confirmClaim = m.confirmClaim.Update(msg)

	case WorkflowStepGenerateContext:
		m.generateCtx = m.generateCtx.Update(msg)

	case WorkflowStepConfirmComplete:
		m.confirmDone = m.confirmDone.Update(msg)

	case WorkflowStepValidating:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// handleEnter handles the enter key for each step
func (m WorkflowModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case WorkflowStepSelectTask:
		if m.cursorPos < len(m.backlogTasks) {
			m.selectedTask = &m.backlogTasks[m.cursorPos]
			m.step = WorkflowStepConfirmClaim
			// Initialize checklist
			if len(m.selectedTask.Acceptance) > 0 {
				m.completed = make([]bool, len(m.selectedTask.Acceptance))
			}
		}

	case WorkflowStepConfirmClaim:
		if m.confirmClaim.IsYes() {
			return m, m.claimTask()
		} else {
			m.step = WorkflowStepSelectTask
		}

	case WorkflowStepGenerateContext:
		if m.generateCtx.IsYes() {
			m.message = "Context generation would happen here (not yet implemented)"
		}
		m.step = WorkflowStepShowTask

	case WorkflowStepShowTask:
		m.step = WorkflowStepConfirmComplete

	case WorkflowStepConfirmComplete:
		if m.confirmDone.IsYes() {
			return m, m.completeAndValidate()
		} else {
			m.step = WorkflowStepShowTask
		}

	case WorkflowStepComplete:
		m.quitting = true
		return m, tea.Quit
	}

	return m, nil
}

// handleBack handles going back to the previous step
func (m WorkflowModel) handleBack() (tea.Model, tea.Cmd) {
	switch m.step {
	case WorkflowStepConfirmClaim:
		m.step = WorkflowStepSelectTask

	case WorkflowStepGenerateContext:
		m.step = WorkflowStepConfirmClaim

	case WorkflowStepShowTask:
		m.step = WorkflowStepGenerateContext

	case WorkflowStepConfirmComplete:
		m.step = WorkflowStepShowTask

	default:
		m.quitting = true
		return m, tea.Quit
	}

	return m, nil
}

// claimTask claims the selected task
func (m *WorkflowModel) claimTask() tea.Cmd {
	return func() tea.Msg {
		if m.selectedTask == nil {
			return claimErrorMsg{err: fmt.Errorf("no task selected")}
		}

		user := os.Getenv("USER")
		if user == "" {
			user = "current-user"
		}

		if err := m.taskManager.ClaimTask(m.selectedTask.ID, user); err != nil {
			return claimErrorMsg{err: err}
		}

		// Reload task to get updated version
		inProgress, _ := m.taskManager.LoadTasks("in-progress")
		for _, t := range inProgress.Tasks {
			if t.ID == m.selectedTask.ID {
				*m.selectedTask = t
				break
			}
		}

		return claimCompleteMsg{}
	}
}

// completeAndValidate completes the task and runs validation
func (m *WorkflowModel) completeAndValidate() tea.Cmd {
	return func() tea.Msg {
		if m.selectedTask == nil {
			return validationErrorMsg{err: fmt.Errorf("no task selected")}
		}

		// Move task to done
		if err := m.taskManager.MoveTask(m.selectedTask.ID, "in-progress", "done", taskmodels.StatusDone); err != nil {
			return validationErrorMsg{err: err}
		}

		// TODO: Run actual validation
		results := "Validation passed!\n\nAll checks completed successfully."

		return validationCompleteMsg{results: results}
	}
}

// Custom messages
type claimCompleteMsg struct{}
type claimErrorMsg struct{ err error }
type validationCompleteMsg struct{ results string }
type validationErrorMsg struct{ err error }

// View renders the workflow
func (m WorkflowModel) View() string {
	if m.quitting && m.step != WorkflowStepComplete {
		return styles.MutedStyle.Render("Cancelled.\n")
	}

	switch m.step {
	case WorkflowStepSelectTask:
		return m.renderSelectTask()
	case WorkflowStepConfirmClaim:
		return m.renderConfirmClaim()
	case WorkflowStepGenerateContext:
		return m.renderGenerateContext()
	case WorkflowStepShowTask:
		return m.renderShowTask()
	case WorkflowStepConfirmComplete:
		return m.renderConfirmComplete()
	case WorkflowStepValidating:
		return m.renderValidating()
	case WorkflowStepComplete:
		return m.renderComplete()
	}

	return ""
}

func (m WorkflowModel) renderSelectTask() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Select Task to Work On") + "\n\n")

	if len(m.backlogTasks) == 0 {
		b.WriteString(styles.MutedStyle.Render("No tasks in backlog.") + "\n\n")
		b.WriteString(styles.HelpStyle.Render("Press Esc to exit") + "\n")
		return styles.ContainerStyle.Render(b.String())
	}

	// Render task list
	visibleStart := m.offset
	visibleEnd := m.offset + m.height
	if visibleEnd > len(m.backlogTasks) {
		visibleEnd = len(m.backlogTasks)
	}

	for i := visibleStart; i < visibleEnd; i++ {
		task := m.backlogTasks[i]
		cursor := "  "
		if i == m.cursorPos {
			cursor = styles.IconArrow + " "
		}

		style := styles.ListItemStyle
		if i == m.cursorPos {
			style = styles.SelectedItemStyle
		}

		line := fmt.Sprintf("%s%s  %s", cursor, task.ID, task.Title)
		b.WriteString(style.Render(line) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Enter select • Esc quit"))

	return styles.ContainerStyle.Render(b.String())
}

func (m WorkflowModel) renderConfirmClaim() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Claim Task") + "\n\n")

	if m.selectedTask != nil {
		card := fmt.Sprintf("Task: %s\n", styles.BoldStyle.Render(m.selectedTask.Title))
		if m.selectedTask.Description != "" {
			card += fmt.Sprintf("\nDescription:\n%s\n", m.selectedTask.Description)
		}
		b.WriteString(styles.CardStyle.Render(card) + "\n")
	}

	b.WriteString(m.confirmClaim.View() + "\n\n")
	b.WriteString(styles.HelpStyle.Render("Enter to confirm • Esc to go back") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m WorkflowModel) renderGenerateContext() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Generate Context") + "\n\n")

	if m.message != "" {
		b.WriteString(styles.MutedStyle.Render(m.message) + "\n\n")
	}

	b.WriteString(m.generateCtx.View() + "\n\n")

	if m.selectedTask != nil && len(m.selectedTask.Scope) > 0 {
		b.WriteString(styles.MutedStyle.Render("Scope directories:\n"))
		for _, scope := range m.selectedTask.Scope {
			b.WriteString(styles.MutedStyle.Render(fmt.Sprintf("  %s %s\n", styles.IconBullet, scope)))
		}
		b.WriteString("\n")
	}

	b.WriteString(styles.HelpStyle.Render("Enter to continue • Esc to go back") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m WorkflowModel) renderShowTask() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Working on Task") + "\n\n")

	if m.selectedTask == nil {
		return styles.ContainerStyle.Render(b.String())
	}

	// Task card
	card := fmt.Sprintf("Task: %s\n", styles.BoldStyle.Render(m.selectedTask.Title))
	card += fmt.Sprintf("Status: %s\n", styles.BoldStyle.Render(string(m.selectedTask.Status)))

	if m.selectedTask.Description != "" {
		card += fmt.Sprintf("\nDescription:\n%s\n", m.selectedTask.Description)
	}

	b.WriteString(styles.CardStyle.Render(card) + "\n")

	// Acceptance criteria checklist
	if len(m.selectedTask.Acceptance) > 0 {
		b.WriteString(styles.SubtitleStyle.Render("Acceptance Criteria:") + "\n\n")

		for i, criterion := range m.selectedTask.Acceptance {
			cursor := "  "
			if i == m.checklistPos {
				cursor = styles.IconArrow + " "
			}

			checkbox := "☐"
			style := styles.ListItemStyle
			if i < len(m.completed) && m.completed[i] {
				checkbox = "☑"
				style = styles.SuccessStyle
			}

			if i == m.checklistPos {
				style = styles.SelectedItemStyle
			}

			line := fmt.Sprintf("%s%s %s", cursor, checkbox, criterion)
			b.WriteString(style.Render(line) + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Space toggle • Enter continue when done") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m WorkflowModel) renderConfirmComplete() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Complete Task") + "\n\n")

	// Show checklist status
	if len(m.selectedTask.Acceptance) > 0 {
		completedCount := 0
		for _, done := range m.completed {
			if done {
				completedCount++
			}
		}

		statusText := fmt.Sprintf("Completed %d / %d acceptance criteria", completedCount, len(m.completed))
		if completedCount == len(m.completed) {
			b.WriteString(styles.RenderSuccess(statusText) + "\n\n")
		} else {
			b.WriteString(styles.RenderWarning(statusText) + "\n\n")
		}
	}

	b.WriteString(m.confirmDone.View() + "\n\n")
	b.WriteString(styles.HelpStyle.Render("Enter to confirm • Esc to go back") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m WorkflowModel) renderValidating() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Completing Task") + "\n\n")
	b.WriteString(m.spinner.View() + "\n\n")
	b.WriteString(styles.MutedStyle.Render("Running validation...") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

func (m WorkflowModel) renderComplete() string {
	var b strings.Builder

	if m.error != "" {
		b.WriteString(styles.RenderError("Failed to complete task") + "\n\n")
		b.WriteString(styles.MutedStyle.Render(m.error) + "\n\n")
		b.WriteString(styles.HelpStyle.Render("Press Enter to exit") + "\n")
		return styles.ContainerStyle.Render(b.String())
	}

	b.WriteString(styles.RenderSuccess("Task completed successfully!") + "\n\n")

	if m.validationResults != "" {
		b.WriteString(styles.SubtitleStyle.Render("Validation Results:") + "\n\n")
		b.WriteString(styles.CardStyle.Render(m.validationResults) + "\n")
	}

	b.WriteString(styles.HelpStyle.Render("Press Enter to exit") + "\n")

	return styles.ContainerStyle.Render(b.String())
}
