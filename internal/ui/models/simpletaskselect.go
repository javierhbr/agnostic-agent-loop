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

// SimpleTaskSelectAction represents the action to perform
type SimpleTaskSelectAction string

const (
	ActionClaim    SimpleTaskSelectAction = "claim"
	ActionComplete SimpleTaskSelectAction = "complete"
	ActionShow     SimpleTaskSelectAction = "show"
)

// actionCompleteMsg is sent when the action completes
type actionCompleteMsg struct {
	success bool
	message string
}

// SimpleTaskSelectModel is a simplified task selector for single actions
type SimpleTaskSelectModel struct {
	action       SimpleTaskSelectAction
	tasks        []taskmodels.Task
	cursorPos    int
	offset       int
	height       int
	width        int
	selectedTask *taskmodels.Task
	showDetails  bool
	confirm      components.Confirm
	showConfirm  bool
	done         bool
	success      bool
	message      string
	error        string
	taskManager  *tasks.TaskManager
}

// NewSimpleTaskSelectModel creates a new simple task selector
func NewSimpleTaskSelectModel(action SimpleTaskSelectAction, filterStatus string) SimpleTaskSelectModel {
	tm := tasks.NewTaskManager(".agentic/tasks")

	// Load tasks based on filter
	taskList, _ := tm.LoadTasks(filterStatus)

	// Create confirmation prompt
	var confirmMsg string
	switch action {
	case ActionClaim:
		confirmMsg = "Claim this task?"
	case ActionComplete:
		confirmMsg = "Mark this task as complete?"
	case ActionShow:
		confirmMsg = "View task details?"
	}
	confirm := components.NewConfirm(confirmMsg, true)

	return SimpleTaskSelectModel{
		action:      action,
		tasks:       taskList.Tasks,
		cursorPos:   0,
		offset:      0,
		height:      10,
		confirm:     confirm,
		taskManager: tm,
	}
}

// Init initializes the model
func (m SimpleTaskSelectModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m SimpleTaskSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case actionCompleteMsg:
		m.done = true
		m.success = msg.success
		if msg.success {
			m.message = msg.message
		} else {
			m.error = msg.message
		}
		return m, tea.Quit

	case tea.KeyMsg:
		// Confirmation has priority
		if m.showConfirm {
			switch msg.String() {
			case "ctrl+c", "esc", "n":
				return m, tea.Quit

			case "enter", "y":
				m.showConfirm = false
				return m, m.executeAction()

			default:
				m.confirm = m.confirm.Update(msg)
			}
			return m, nil
		}

		// Details view
		if m.showDetails {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit

			case "esc":
				m.showDetails = false

			case "enter", "y":
				m.showDetails = false
				m.showConfirm = true
			}
			return m, nil
		}

		// Normal navigation
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursorPos > 0 {
				m.cursorPos--
				if m.cursorPos < m.offset {
					m.offset = m.cursorPos
				}
			}

		case "down", "j":
			if m.cursorPos < len(m.tasks)-1 {
				m.cursorPos++
				if m.cursorPos >= m.offset+m.height {
					m.offset = m.cursorPos - m.height + 1
				}
			}

		case "enter":
			// Select task and show details
			if m.cursorPos < len(m.tasks) {
				m.selectedTask = &m.tasks[m.cursorPos]
				m.showDetails = true
			}

		case "y":
			// Quick confirm without details
			if m.cursorPos < len(m.tasks) {
				m.selectedTask = &m.tasks[m.cursorPos]
				m.showConfirm = true
			}
		}
	}

	return m, cmd
}

// executeAction performs the selected action
func (m SimpleTaskSelectModel) executeAction() tea.Cmd {
	return func() tea.Msg {
		if m.selectedTask == nil {
			return actionCompleteMsg{
				success: false,
				message: "No task selected",
			}
		}

		var err error
		var successMsg string

		switch m.action {
		case ActionClaim:
			user := os.Getenv("USER")
			if user == "" {
				user = "unknown-agent"
			}
			err = m.taskManager.ClaimTask(m.selectedTask.ID, user)
			successMsg = fmt.Sprintf("Task %s claimed successfully!", m.selectedTask.ID)

		case ActionComplete:
			err = m.taskManager.MoveTask(m.selectedTask.ID, "in-progress", "done", taskmodels.StatusDone)
			successMsg = fmt.Sprintf("Task %s completed successfully!", m.selectedTask.ID)

		case ActionShow:
			// Show is handled in the view, not an action
			successMsg = "Showing task details"
		}

		if err != nil {
			return actionCompleteMsg{
				success: false,
				message: fmt.Sprintf("Error: %v", err),
			}
		}

		return actionCompleteMsg{
			success: true,
			message: successMsg,
		}
	}
}

// View renders the model
func (m SimpleTaskSelectModel) View() string {
	if m.done {
		if m.success {
			return styles.RenderSuccess(m.message) + "\n"
		}
		return styles.RenderError(m.error) + "\n"
	}

	if m.showConfirm {
		return m.renderConfirm()
	}

	if m.showDetails {
		return m.renderDetails()
	}

	return m.renderTaskList()
}

// renderTaskList renders the task selection list
func (m SimpleTaskSelectModel) renderTaskList() string {
	var b strings.Builder

	// Title based on action
	var title string
	switch m.action {
	case ActionClaim:
		title = "Select Task to Claim"
	case ActionComplete:
		title = "Select Task to Complete"
	case ActionShow:
		title = "Select Task to View"
	}
	b.WriteString(styles.TitleStyle.Render(title) + "\n\n")

	// Task list
	if len(m.tasks) == 0 {
		b.WriteString(styles.MutedStyle.Render("No tasks available.") + "\n\n")
	} else {
		b.WriteString(m.renderTasks() + "\n\n")
	}

	// Help text
	helpParts := []string{
		"↑/↓ navigate",
		"Enter details",
		"y quick confirm",
		"q quit",
	}
	b.WriteString(styles.HelpStyle.Render(strings.Join(helpParts, " • ")))

	return styles.ContainerStyle.Render(b.String())
}

// renderTasks renders the task list
func (m SimpleTaskSelectModel) renderTasks() string {
	var b strings.Builder

	visibleStart := m.offset
	visibleEnd := m.offset + m.height
	if visibleEnd > len(m.tasks) {
		visibleEnd = len(m.tasks)
	}

	for i := visibleStart; i < visibleEnd; i++ {
		task := m.tasks[i]
		cursor := "  "
		if i == m.cursorPos {
			cursor = styles.IconArrow + " "
		}

		style := styles.ListItemStyle
		if i == m.cursorPos {
			style = styles.SelectedItemStyle
		}

		// Truncate title if too long
		title := task.Title
		if len(title) > 60 {
			title = title[:57] + "..."
		}

		line := fmt.Sprintf("%s%s  %s", cursor, task.ID, title)
		b.WriteString(style.Render(line) + "\n")
	}

	return b.String()
}

// renderDetails renders the task details view
func (m SimpleTaskSelectModel) renderDetails() string {
	if m.selectedTask == nil {
		return ""
	}

	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Task Details") + "\n\n")

	// Task card
	card := fmt.Sprintf("ID: %s\n", styles.BoldStyle.Render(m.selectedTask.ID))
	card += fmt.Sprintf("Title: %s\n", styles.BoldStyle.Render(m.selectedTask.Title))
	card += fmt.Sprintf("Status: %s\n", m.selectedTask.Status)

	if m.selectedTask.Description != "" {
		card += fmt.Sprintf("\nDescription:\n%s\n", m.selectedTask.Description)
	}

	if len(m.selectedTask.SpecRefs) > 0 {
		card += "\nSpecification References:\n"
		for _, ref := range m.selectedTask.SpecRefs {
			card += fmt.Sprintf("  %s %s\n", styles.IconBullet, ref)
		}
	}

	if len(m.selectedTask.Scope) > 0 {
		card += "\nScope:\n"
		for _, scope := range m.selectedTask.Scope {
			card += fmt.Sprintf("  %s %s\n", styles.IconBullet, scope)
		}
	}

	if len(m.selectedTask.Outputs) > 0 {
		card += "\nExpected Outputs:\n"
		for _, output := range m.selectedTask.Outputs {
			card += fmt.Sprintf("  %s %s\n", styles.IconBullet, output)
		}
	}

	if len(m.selectedTask.Acceptance) > 0 {
		card += "\nAcceptance Criteria:\n"
		for _, criterion := range m.selectedTask.Acceptance {
			card += fmt.Sprintf("  %s %s\n", styles.IconBullet, criterion)
		}
	}

	b.WriteString(styles.CardStyle.Render(card) + "\n")

	var actionText string
	switch m.action {
	case ActionClaim:
		actionText = "Enter to claim"
	case ActionComplete:
		actionText = "Enter to complete"
	case ActionShow:
		actionText = "Enter to confirm"
	}

	b.WriteString(styles.HelpStyle.Render(actionText+" • Esc back") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

// renderConfirm renders the confirmation prompt
func (m SimpleTaskSelectModel) renderConfirm() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Confirm Action") + "\n\n")

	if m.selectedTask != nil {
		taskInfo := fmt.Sprintf("Task: %s\n%s", m.selectedTask.ID, m.selectedTask.Title)
		b.WriteString(styles.CardStyle.Render(taskInfo) + "\n\n")
	}

	b.WriteString(m.confirm.View() + "\n")

	return styles.ContainerStyle.Render(b.String())
}
