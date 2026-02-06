package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/orchestrator"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/internal/ui/components"
	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the agent orchestrator for a task",
	Run: func(cmd *cobra.Command, args []string) {
		taskID, _ := cmd.Flags().GetString("task")

		// Interactive mode - task selection
		if helpers.ShouldUseInteractiveMode(cmd) && taskID == "" {
			model := &runOrchestratorModel{
				step: "select-task",
			}

			p := tea.NewProgram(model)
			if _, err := p.Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		// Flag mode - require task ID
		if taskID == "" {
			fmt.Println("Error: --task required in non-interactive mode")
			fmt.Println("Usage: agentic-agent run --task <task-id>")
			fmt.Println("   or: agentic-agent run  (interactive mode)")
			os.Exit(1)
		}

		if err := orchestrator.RunLoop(taskID); err != nil {
			fmt.Printf("Error running orchestrator: %v\n", err)
			os.Exit(1)
		}
	},
}

// runOrchestratorModel is a Bubble Tea model for running the orchestrator
type runOrchestratorModel struct {
	step         string // "select-task", "confirm", "running", "done"
	taskSelector components.SimpleSelect
	tasks        []string
	selectedTask string
	success      bool
	message      string
}

func (m *runOrchestratorModel) Init() tea.Cmd {
	// Load tasks for selection
	tm := tasks.NewTaskManager(".agentic/tasks")

	// Get tasks from backlog and in-progress
	var allTasks []string
	for _, status := range []string{"backlog", "in-progress"} {
		taskList, _ := tm.LoadTasks(status)
		for _, t := range taskList.Tasks {
			allTasks = append(allTasks, t.ID)
		}
	}

	m.tasks = allTasks

	// Create task selector
	options := []components.SelectOption{}
	for _, taskID := range allTasks {
		options = append(options, components.NewSelectOption(taskID, "Run orchestrator for this task", taskID))
	}

	if len(options) == 0 {
		options = append(options, components.NewSelectOption("No tasks found", "Create a task first", ""))
	}

	m.taskSelector = components.NewSimpleSelect("Select task for orchestrator", options)

	return nil
}

func (m *runOrchestratorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			if m.step == "confirm" {
				m.step = "select-task"
			} else {
				return m, tea.Quit
			}

		case "enter":
			if m.step == "select-task" {
				m.selectedTask = m.taskSelector.SelectedOption().Value()
				if m.selectedTask == "" {
					m.message = "No task selected"
					m.success = false
					m.step = "done"
					return m, tea.Quit
				}
				m.step = "confirm"
			} else if m.step == "confirm" {
				// Run the orchestrator
				m.step = "running"
				if err := orchestrator.RunLoop(m.selectedTask); err != nil {
					m.message = fmt.Sprintf("Error: %v", err)
					m.success = false
				} else {
					m.message = fmt.Sprintf("Orchestrator completed for task %s", m.selectedTask)
					m.success = true
				}
				m.step = "done"
				return m, tea.Quit
			}

		case "n":
			if m.step == "confirm" {
				return m, tea.Quit
			}

		default:
			if m.step == "select-task" {
				m.taskSelector = m.taskSelector.Update(msg)
			}
		}
	}

	return m, nil
}

func (m *runOrchestratorModel) View() string {
	if m.step == "done" {
		if m.success {
			return styles.RenderSuccess(m.message) + "\n"
		}
		return styles.RenderError(m.message) + "\n"
	}

	var b strings.Builder

	switch m.step {
	case "select-task":
		b.WriteString(styles.TitleStyle.Render("Run Agent Orchestrator") + "\n\n")
		b.WriteString(m.taskSelector.View() + "\n")
		b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Enter select • q quit") + "\n")

	case "confirm":
		b.WriteString(styles.TitleStyle.Render("Confirm Orchestrator Run") + "\n\n")
		b.WriteString(fmt.Sprintf("Task: %s\n\n", styles.BoldStyle.Render(m.selectedTask)))
		b.WriteString("This will run the agent orchestrator loop for this task.\n\n")
		b.WriteString(styles.HelpStyle.Render("Enter confirm • n/Esc cancel") + "\n")

	case "running":
		b.WriteString(styles.TitleStyle.Render("Running Orchestrator...") + "\n")
		b.WriteString("This may take a while...\n")
	}

	return styles.ContainerStyle.Render(b.String())
}

func init() {
	runCmd.Flags().String("task", "", "Task ID to run")
}
