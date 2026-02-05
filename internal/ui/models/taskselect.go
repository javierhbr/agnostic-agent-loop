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

// TaskTab represents the task status tab
type TaskTab int

const (
	TabBacklog TaskTab = iota
	TabInProgress
	TabDone
)

// TaskSelectModel is the model for task selection menu
type TaskSelectModel struct {
	currentTab    TaskTab
	backlogTasks  []taskmodels.Task
	inProgressTasks []taskmodels.Task
	doneTasks     []taskmodels.Task
	cursorPos     int
	offset        int
	height        int
	width         int
	selectedTask  *taskmodels.Task
	showDetails   bool
	actionMenu    components.SimpleSelect
	showActionMenu bool
	error         string
	message       string
	quitting      bool
	taskManager   *tasks.TaskManager
}

// NewTaskSelectModel creates a new task selection model
func NewTaskSelectModel() TaskSelectModel {
	tm := tasks.NewTaskManager(".agentic/tasks")

	// Load tasks from all categories
	backlog, _ := tm.LoadTasks("backlog")
	inProgress, _ := tm.LoadTasks("in-progress")
	done, _ := tm.LoadTasks("done")

	// Action menu options
	actionOptions := []components.SelectOption{
		components.NewSelectOption("Claim Task", "Move task to in-progress and assign to you", "claim"),
		components.NewSelectOption("Complete Task", "Mark task as complete and move to done", "complete"),
		components.NewSelectOption("Show Details", "View full task details", "show"),
		components.NewSelectOption("Cancel", "Return to task list", "cancel"),
	}
	actionMenu := components.NewSimpleSelect("Action", actionOptions)

	return TaskSelectModel{
		currentTab:      TabBacklog,
		backlogTasks:    backlog.Tasks,
		inProgressTasks: inProgress.Tasks,
		doneTasks:       done.Tasks,
		cursorPos:       0,
		offset:          0,
		height:          15,
		actionMenu:      actionMenu,
		taskManager:     tm,
	}
}

// Init initializes the model
func (m TaskSelectModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m TaskSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		// Action menu has priority
		if m.showActionMenu {
			switch msg.String() {
			case "ctrl+c", "esc":
				m.showActionMenu = false
				return m, nil

			case "enter":
				return m.handleAction()

			default:
				m.actionMenu = m.actionMenu.Update(msg)
			}
			return m, nil
		}

		// Normal navigation
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "esc":
			if m.showDetails {
				m.showDetails = false
			} else {
				m.quitting = true
				return m, tea.Quit
			}

		case "tab":
			m.nextTab()
			m.cursorPos = 0
			m.offset = 0

		case "shift+tab":
			m.prevTab()
			m.cursorPos = 0
			m.offset = 0

		case "up", "k":
			if m.cursorPos > 0 {
				m.cursorPos--
				if m.cursorPos < m.offset {
					m.offset = m.cursorPos
				}
			}

		case "down", "j":
			tasks := m.getCurrentTasks()
			if m.cursorPos < len(tasks)-1 {
				m.cursorPos++
				if m.cursorPos >= m.offset+m.height {
					m.offset = m.cursorPos - m.height + 1
				}
			}

		case "enter":
			// Select task and show details
			tasks := m.getCurrentTasks()
			if m.cursorPos < len(tasks) {
				m.selectedTask = &tasks[m.cursorPos]
				m.showDetails = true
			}

		case "c":
			// Quick claim
			tasks := m.getCurrentTasks()
			if m.cursorPos < len(tasks) {
				m.selectedTask = &tasks[m.cursorPos]
				return m.claimTask()
			}

		case "d":
			// Quick complete
			tasks := m.getCurrentTasks()
			if m.cursorPos < len(tasks) {
				m.selectedTask = &tasks[m.cursorPos]
				return m.completeTask()
			}

		case "a":
			// Show action menu
			tasks := m.getCurrentTasks()
			if m.cursorPos < len(tasks) {
				m.selectedTask = &tasks[m.cursorPos]
				m.showActionMenu = true
			}
		}
	}

	return m, cmd
}

// nextTab switches to the next tab
func (m *TaskSelectModel) nextTab() {
	m.currentTab = (m.currentTab + 1) % 3
}

// prevTab switches to the previous tab
func (m *TaskSelectModel) prevTab() {
	if m.currentTab == 0 {
		m.currentTab = 2
	} else {
		m.currentTab--
	}
}

// getCurrentTasks returns the tasks for the current tab
func (m TaskSelectModel) getCurrentTasks() []taskmodels.Task {
	switch m.currentTab {
	case TabBacklog:
		return m.backlogTasks
	case TabInProgress:
		return m.inProgressTasks
	case TabDone:
		return m.doneTasks
	default:
		return []taskmodels.Task{}
	}
}

// handleAction handles the selected action from action menu
func (m TaskSelectModel) handleAction() (tea.Model, tea.Cmd) {
	action := m.actionMenu.SelectedOption().Value()

	switch action {
	case "claim":
		return m.claimTask()
	case "complete":
		return m.completeTask()
	case "show":
		m.showActionMenu = false
		m.showDetails = true
		return m, nil
	case "cancel":
		m.showActionMenu = false
		return m, nil
	}

	return m, nil
}

// claimTask claims the selected task
func (m TaskSelectModel) claimTask() (tea.Model, tea.Cmd) {
	if m.selectedTask == nil {
		return m, nil
	}

	// Get current user
	user := "current-user" // TODO: Get from environment or config

	// Move task from backlog to in-progress
	if err := m.taskManager.ClaimTask(m.selectedTask.ID, user); err != nil {
		m.error = err.Error()
		m.showActionMenu = false
		return m, nil
	}

	m.message = fmt.Sprintf("Task %s claimed successfully!", m.selectedTask.ID)
	m.showActionMenu = false
	m.showDetails = false

	// Reload tasks
	backlog, _ := m.taskManager.LoadTasks("backlog")
	inProgress, _ := m.taskManager.LoadTasks("in-progress")
	m.backlogTasks = backlog.Tasks
	m.inProgressTasks = inProgress.Tasks

	// Switch to in-progress tab
	m.currentTab = TabInProgress
	m.cursorPos = 0
	m.offset = 0

	return m, nil
}

// completeTask completes the selected task
func (m TaskSelectModel) completeTask() (tea.Model, tea.Cmd) {
	if m.selectedTask == nil {
		return m, nil
	}

	// Move task from in-progress to done
	if err := m.taskManager.MoveTask(m.selectedTask.ID, "in-progress", "done", taskmodels.StatusDone); err != nil {
		m.error = err.Error()
		m.showActionMenu = false
		return m, nil
	}

	m.message = fmt.Sprintf("Task %s completed successfully!", m.selectedTask.ID)
	m.showActionMenu = false
	m.showDetails = false

	// Reload tasks
	inProgress, _ := m.taskManager.LoadTasks("in-progress")
	done, _ := m.taskManager.LoadTasks("done")
	m.inProgressTasks = inProgress.Tasks
	m.doneTasks = done.Tasks

	// Switch to done tab
	m.currentTab = TabDone
	m.cursorPos = 0
	m.offset = 0

	return m, nil
}

// View renders the task selection menu
func (m TaskSelectModel) View() string {
	if m.quitting {
		return ""
	}

	if m.showActionMenu {
		return m.renderActionMenu()
	}

	if m.showDetails {
		return m.renderDetails()
	}

	return m.renderTaskList()
}

// renderTaskList renders the main task list with tabs
func (m TaskSelectModel) renderTaskList() string {
	var b strings.Builder

	// Title
	b.WriteString(styles.TitleStyle.Render("Task Manager") + "\n\n")

	// Show message if any
	if m.message != "" {
		b.WriteString(styles.RenderSuccess(m.message) + "\n\n")
	}
	if m.error != "" {
		b.WriteString(styles.RenderError(m.error) + "\n\n")
	}

	// Tabs
	b.WriteString(m.renderTabs() + "\n\n")

	// Task list
	tasks := m.getCurrentTasks()
	if len(tasks) == 0 {
		b.WriteString(styles.MutedStyle.Render("No tasks in this category.") + "\n\n")
	} else {
		b.WriteString(m.renderTasks(tasks) + "\n\n")
	}

	// Help text
	helpParts := []string{
		"↑/↓ navigate",
		"Tab switch",
		"Enter details",
		"c claim",
		"d complete",
		"a actions",
		"q quit",
	}
	b.WriteString(styles.HelpStyle.Render(strings.Join(helpParts, " • ")))

	return styles.ContainerStyle.Render(b.String())
}

// renderTabs renders the tab bar
func (m TaskSelectModel) renderTabs() string {
	tabs := []string{"Backlog", "In Progress", "Done"}
	counts := []int{len(m.backlogTasks), len(m.inProgressTasks), len(m.doneTasks)}

	var parts []string
	for i, tab := range tabs {
		count := fmt.Sprintf("(%d)", counts[i])
		tabText := fmt.Sprintf("%s %s", tab, count)

		if TaskTab(i) == m.currentTab {
			parts = append(parts, styles.SelectedTabStyle.Render(tabText))
		} else {
			parts = append(parts, styles.TabStyle.Render(tabText))
		}
	}

	return strings.Join(parts, " ")
}

// renderTasks renders the task list
func (m TaskSelectModel) renderTasks(tasks []taskmodels.Task) string {
	var b strings.Builder

	visibleStart := m.offset
	visibleEnd := m.offset + m.height
	if visibleEnd > len(tasks) {
		visibleEnd = len(tasks)
	}

	for i := visibleStart; i < visibleEnd; i++ {
		task := tasks[i]
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

	return b.String()
}

// renderDetails renders the task details view
func (m TaskSelectModel) renderDetails() string {
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
	b.WriteString(styles.HelpStyle.Render("c claim • d complete • a actions • Esc back") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

// renderActionMenu renders the action menu
func (m TaskSelectModel) renderActionMenu() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Task Actions") + "\n\n")

	if m.selectedTask != nil {
		b.WriteString(styles.SubtitleStyle.Render(fmt.Sprintf("Task: %s", m.selectedTask.Title)) + "\n\n")
	}

	b.WriteString(m.actionMenu.View() + "\n")
	b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Enter select • Esc cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}
