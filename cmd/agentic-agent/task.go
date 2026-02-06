package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/internal/ui/components"
	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	uimodels "github.com/javierbenavides/agentic-agent/internal/ui/models"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks",
}

var taskCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task (interactive wizard or flags)",
	Long: `Create a new task using either interactive mode or traditional flags.

Interactive Mode (no flags):
  agentic-agent task create

  Launches a step-by-step wizard for creating tasks with:
  - Title input with validation
  - Optional description (multi-line)
  - Acceptance criteria editor
  - Preview before creation

Flag Mode (with flags):
  agentic-agent task create --title "My Task" --description "Details"

  Traditional command-line mode with all options as flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if we should use interactive mode
		if shouldUseInteractiveTaskCreate(cmd) {
			runInteractiveTaskCreate()
			return
		}

		// Traditional flag-based mode
		title, _ := cmd.Flags().GetString("title")
		if title == "" {
			fmt.Println("Error: --title is required")
			os.Exit(1)
		}

		// Validate title
		if err := validateTaskTitle(title); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		description, _ := cmd.Flags().GetString("description")
		scopeStr, _ := cmd.Flags().GetString("scope")
		specRefsStr, _ := cmd.Flags().GetString("spec-refs")
		inputsStr, _ := cmd.Flags().GetString("inputs")
		outputsStr, _ := cmd.Flags().GetString("outputs")
		acceptanceStr, _ := cmd.Flags().GetString("acceptance")

		tm := tasks.NewTaskManager(".agentic/tasks")
		task, err := tm.CreateTask(title)
		if err != nil {
			fmt.Printf("Error creating task: %v\n", err)
			os.Exit(1)
		}

		// Update task with additional fields
		if description != "" {
			task.Description = description
		}
		if scopeStr != "" {
			task.Scope = parseCommaSeparated(scopeStr)
		}
		if specRefsStr != "" {
			task.SpecRefs = parseCommaSeparated(specRefsStr)
		}
		if inputsStr != "" {
			task.Inputs = parseCommaSeparated(inputsStr)
		}
		if outputsStr != "" {
			task.Outputs = parseCommaSeparated(outputsStr)
		}
		if acceptanceStr != "" {
			task.Acceptance = parseCommaSeparated(acceptanceStr)
		}

		// Save updated task
		backlog, err := tm.LoadTasks("backlog")
		if err != nil {
			fmt.Printf("Error loading backlog: %v\n", err)
			os.Exit(1)
		}

		// Update the task in backlog
		for i, t := range backlog.Tasks {
			if t.ID == task.ID {
				backlog.Tasks[i] = *task
				break
			}
		}

		if err := tm.SaveTasks("backlog", backlog); err != nil {
			fmt.Printf("Error saving task: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Created task %s: %s\n", task.ID, task.Title)
		if len(task.SpecRefs) > 0 {
			fmt.Printf("  Spec refs: %s\n", strings.Join(task.SpecRefs, ", "))
		}
		if len(task.Acceptance) > 0 {
			fmt.Printf("  Acceptance criteria: %d items\n", len(task.Acceptance))
		}
	},
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks (interactive menu or simple list)",
	Long: `List tasks with either an interactive menu or simple text output.

Interactive Mode (no flags):
  agentic-agent task list

  Launches an interactive menu with tabs for Backlog/In Progress/Done tasks.
  Supports navigation, claiming, completing, and viewing task details.

Flag Mode (with flags or --no-interactive):
  agentic-agent task list --no-interactive

  Simple text output listing all tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if we should use interactive mode
		if helpers.ShouldUseInteractiveMode(cmd) {
			runInteractiveTaskList()
			return
		}

		// Traditional text-based list
		tm := tasks.NewTaskManager(".agentic/tasks")

		sources := []string{"backlog", "in-progress", "done"}
		for _, src := range sources {
			list, err := tm.LoadTasks(src)
			if err != nil {
				fmt.Printf("Error loading %s: %v\n", src, err)
				continue
			}
			if len(list.Tasks) > 0 {
				fmt.Printf("\n--- %s ---\n", strings.ToUpper(src))
				for _, t := range list.Tasks {
					assignee := ""
					if t.AssignedTo != "" {
						assignee = fmt.Sprintf(" (@%s)", t.AssignedTo)
					}
					fmt.Printf("[%s] %s%s\n", t.ID, t.Title, assignee)
					for _, st := range t.SubTasks {
						fmt.Printf("  - [%s] %s\n", st.ID, st.Title)
					}
				}
			}
		}
	},
}

var taskClaimCmd = &cobra.Command{
	Use:   "claim [task-id]",
	Short: "Claim a task",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Interactive mode if no args
		if helpers.ShouldUseInteractiveMode(cmd) && len(args) == 0 {
			model := uimodels.NewSimpleTaskSelectModel(uimodels.ActionClaim, "backlog")
			p := tea.NewProgram(model)
			if _, err := p.Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		// Flag mode - require task ID
		if len(args) != 1 {
			fmt.Println("Error: task-id argument required in non-interactive mode")
			fmt.Println("Usage: agentic-agent task claim <task-id>")
			fmt.Println("   or: agentic-agent task claim  (interactive mode)")
			os.Exit(1)
		}

		taskID := args[0]
		user := os.Getenv("USER")
		if user == "" {
			user = "unknown-agent"
		}

		tm := tasks.NewTaskManager(".agentic/tasks")
		if err := tm.ClaimTask(taskID, user); err != nil {
			fmt.Printf("Error claiming task: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Claimed task %s\n", taskID)
	},
}

var taskCompleteCmd = &cobra.Command{
	Use:   "complete [task-id]",
	Short: "Mark a task as done",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Interactive mode if no args
		if helpers.ShouldUseInteractiveMode(cmd) && len(args) == 0 {
			model := uimodels.NewSimpleTaskSelectModel(uimodels.ActionComplete, "in-progress")
			p := tea.NewProgram(model)
			if _, err := p.Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		// Flag mode - require task ID
		if len(args) != 1 {
			fmt.Println("Error: task-id argument required in non-interactive mode")
			fmt.Println("Usage: agentic-agent task complete <task-id>")
			fmt.Println("   or: agentic-agent task complete  (interactive mode)")
			os.Exit(1)
		}

		taskID := args[0]
		tm := tasks.NewTaskManager(".agentic/tasks")

		// Assuming from in-progress to done
		if err := tm.MoveTask(taskID, "in-progress", "done", models.StatusDone); err != nil {
			fmt.Printf("Error completing task: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Completed task %s\n", taskID)
	},
}

var taskDecomposeCmd = &cobra.Command{
	Use:   "decompose [task-id] [subtask1] [subtask2] ...",
	Short: "Decompose a task into subtasks",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// Interactive mode if no args
		if helpers.ShouldUseInteractiveMode(cmd) && len(args) == 0 {
			model := &taskDecomposeModel{
				step: "select-task",
			}

			p := tea.NewProgram(model)
			if _, err := p.Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		// Flag mode - require task ID and at least one subtask
		if len(args) < 2 {
			fmt.Println("Error: task-id and at least one subtask required in non-interactive mode")
			fmt.Println("Usage: agentic-agent task decompose <task-id> <subtask1> [subtask2] ...")
			fmt.Println("   or: agentic-agent task decompose  (interactive mode)")
			os.Exit(1)
		}

		taskID := args[0]
		subtasks := args[1:]

		tm := tasks.NewTaskManager(".agentic/tasks")
		if err := tm.DecomposeTask(taskID, subtasks); err != nil {
			fmt.Printf("Error decomposing task: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Decomposed task %s with %d subtasks\n", taskID, len(subtasks))
	},
}

// taskDecomposeModel is a Bubble Tea model for task decomposition
type taskDecomposeModel struct {
	step            string // "select-task", "edit-subtasks", "confirm", "done"
	taskSelector    uimodels.SimpleTaskSelectModel
	selectedTask    *models.Task
	subtaskEditor   components.MultiItemEditor
	success         bool
	message         string
	taskManager     *tasks.TaskManager
}

func (m *taskDecomposeModel) Init() tea.Cmd {
	// Initialize with task selector showing all tasks
	m.taskSelector = uimodels.NewSimpleTaskSelectModel(uimodels.ActionShow, "backlog")
	m.taskManager = tasks.NewTaskManager(".agentic/tasks")
	return nil
}

func (m *taskDecomposeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			if m.step == "edit-subtasks" {
				m.step = "select-task"
			} else if m.step == "confirm" {
				m.step = "edit-subtasks"
			} else {
				return m, tea.Quit
			}

		case "enter":
			if m.step == "select-task" {
				// Get selected task from selector
				// For now, we'll need to access the internal state
				m.step = "edit-subtasks"
				m.subtaskEditor = components.NewMultiItemEditor("Subtasks")
			} else if m.step == "edit-subtasks" {
				items := m.subtaskEditor.GetItems()
				if len(items) == 0 {
					m.message = "At least one subtask is required"
					return m, nil
				}
				m.step = "confirm"
			} else if m.step == "confirm" {
				// Execute decomposition
				subtasks := m.subtaskEditor.GetItems()
				if len(subtasks) == 0 {
					m.message = "No subtasks to create"
					m.success = false
					m.step = "done"
					return m, tea.Quit
				}

				// For now, use a placeholder task ID since we need better integration
				// In a real implementation, we'd get this from the task selector
				m.message = fmt.Sprintf("Would decompose task into %d subtasks", len(subtasks))
				m.success = true
				m.step = "done"
				return m, tea.Quit
			}

		default:
			if m.step == "edit-subtasks" {
				var cmd tea.Cmd
				m.subtaskEditor, cmd = m.subtaskEditor.Update(msg)
				return m, cmd
			}
		}
	}

	return m, nil
}

func (m *taskDecomposeModel) View() string {
	if m.step == "done" {
		if m.success {
			return styles.RenderSuccess(m.message) + "\n"
		}
		return styles.RenderError(m.message) + "\n"
	}

	var b strings.Builder

	switch m.step {
	case "select-task":
		b.WriteString(styles.TitleStyle.Render("Select Task to Decompose") + "\n\n")
		b.WriteString("Note: Task selection integration coming soon.\n")
		b.WriteString("For now, use flag mode: agentic-agent task decompose <task-id> <subtask1> <subtask2>\n\n")
		b.WriteString(styles.HelpStyle.Render("Press Enter to continue to subtask editor demo • q quit") + "\n")

	case "edit-subtasks":
		b.WriteString(styles.TitleStyle.Render("Add Subtasks") + "\n\n")
		b.WriteString(m.subtaskEditor.View() + "\n")

	case "confirm":
		b.WriteString(styles.TitleStyle.Render("Confirm Decomposition") + "\n\n")
		b.WriteString("Subtasks to create:\n\n")
		for i, subtask := range m.subtaskEditor.GetItems() {
			b.WriteString(fmt.Sprintf("  %d. %s\n", i+1, subtask))
		}
		b.WriteString("\n")
		b.WriteString(styles.HelpStyle.Render("Enter confirm • Esc back • q cancel") + "\n")
	}

	return styles.ContainerStyle.Render(b.String())
}

var taskShowCmd = &cobra.Command{
	Use:   "show [task-id]",
	Short: "Display detailed information about a task",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Interactive mode if no args - show task list with details view
		if helpers.ShouldUseInteractiveMode(cmd) && len(args) == 0 {
			// Use the full task list interface which already has details view
			model := uimodels.NewTaskSelectModel()
			p := tea.NewProgram(model)
			if _, err := p.Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		// Flag mode - require task ID
		if len(args) != 1 {
			fmt.Println("Error: task-id argument required in non-interactive mode")
			fmt.Println("Usage: agentic-agent task show <task-id>")
			fmt.Println("   or: agentic-agent task show  (interactive mode)")
			os.Exit(1)
		}

		taskID := args[0]
		tm := tasks.NewTaskManager(".agentic/tasks")

		task, source, err := tm.FindTask(taskID)
		if err != nil {
			fmt.Printf("Error finding task: %v\n", err)
			os.Exit(1)
		}
		if task == nil {
			fmt.Printf("Task %s not found\n", taskID)
			os.Exit(1)
		}

		// Display task details
		fmt.Printf("ID: %s\n", task.ID)
		fmt.Printf("Title: %s\n", task.Title)
		fmt.Printf("Status: %s (%s)\n", task.Status, source)

		if task.Description != "" {
			fmt.Printf("Description: %s\n", task.Description)
		}

		if task.AssignedTo != "" {
			fmt.Printf("Assigned To: %s\n", task.AssignedTo)
		}

		if len(task.Scope) > 0 {
			fmt.Printf("Scope:\n")
			for _, s := range task.Scope {
				fmt.Printf("  - %s\n", s)
			}
		}

		if len(task.SpecRefs) > 0 {
			fmt.Printf("Spec Refs:\n")
			for _, ref := range task.SpecRefs {
				fmt.Printf("  - %s\n", ref)
			}
		}

		if len(task.Inputs) > 0 {
			fmt.Printf("Inputs:\n")
			for _, input := range task.Inputs {
				fmt.Printf("  - %s\n", input)
			}
		}

		if len(task.Outputs) > 0 {
			fmt.Printf("Outputs:\n")
			for _, output := range task.Outputs {
				fmt.Printf("  - %s\n", output)
			}
		}

		if len(task.Acceptance) > 0 {
			fmt.Printf("Acceptance Criteria:\n")
			for _, criterion := range task.Acceptance {
				fmt.Printf("  - %s\n", criterion)
			}
		}

		if len(task.SubTasks) > 0 {
			fmt.Printf("Subtasks:\n")
			for _, st := range task.SubTasks {
				statusIcon := "○"
				if st.Status == models.StatusDone {
					statusIcon = "✓"
				} else if st.Status == models.StatusInProgress {
					statusIcon = "◐"
				}
				assignee := ""
				if st.AssignedTo != "" {
					assignee = fmt.Sprintf(" (@%s)", st.AssignedTo)
				}
				fmt.Printf("  %s [%s] %s%s\n", statusIcon, st.ID, st.Title, assignee)
			}
		}
	},
}

// shouldUseInteractiveTaskCreate checks if task create should run in interactive mode
func shouldUseInteractiveTaskCreate(cmd *cobra.Command) bool {
	return helpers.ShouldUseInteractiveMode(cmd)
}

// runInteractiveTaskCreate runs the interactive task creation wizard
func runInteractiveTaskCreate() {
	model := uimodels.NewTaskCreateModel()
	p := tea.NewProgram(model)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running task creation wizard: %v\n", err)
		os.Exit(1)
	}
}

// runInteractiveTaskList runs the interactive task list menu
func runInteractiveTaskList() {
	model := uimodels.NewTaskSelectModel()
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running task list: %v\n", err)
		os.Exit(1)
	}
}

// validateTaskTitle validates task title input
func validateTaskTitle(title string) error {
	if len(title) == 0 {
		return fmt.Errorf("title cannot be empty")
	}
	if len(title) > 200 {
		return fmt.Errorf("title too long (max 200 characters)")
	}
	// Check for invalid characters that could cause issues
	if strings.Contains(title, "\n") || strings.Contains(title, "\r") {
		return fmt.Errorf("title cannot contain newlines")
	}
	return nil
}

// parseCommaSeparated splits a comma-separated string into a slice
func parseCommaSeparated(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

var taskSampleCmd = &cobra.Command{
	Use:   "sample-task",
	Short: "Create a sample task with pre-filled data",
	Long: `Create a sample task to demonstrate the task management workflow.

This command creates a pre-configured task with example data including:
- A descriptive title
- Sample description
- Example acceptance criteria
- Suggested scope

Perfect for testing the workflow or learning how tasks work.`,
	Run: func(cmd *cobra.Command, args []string) {
		tm := tasks.NewTaskManager(".agentic/tasks")

		// Create sample task
		task, err := tm.CreateTask("Implement user authentication system")
		if err != nil {
			fmt.Printf("Error creating sample task: %v\n", err)
			os.Exit(1)
		}

		// Fill with sample data
		task.Description = `Add a complete user authentication system with login, logout, and session management.

This task involves:
- Creating authentication endpoints
- Implementing JWT token generation
- Adding session management
- Creating login/logout UI components`

		task.Acceptance = []string{
			"User can login with email and password",
			"JWT tokens are generated on successful login",
			"Session persists across page refreshes",
			"User can logout and clear session",
			"Tests cover authentication flow",
		}

		task.Scope = []string{
			"internal/auth",
			"cmd/api/handlers",
			"web/components/auth",
		}

		// Save updated task
		backlog, err := tm.LoadTasks("backlog")
		if err != nil {
			fmt.Printf("Error loading backlog: %v\n", err)
			os.Exit(1)
		}

		// Update the task in backlog
		for i, t := range backlog.Tasks {
			if t.ID == task.ID {
				backlog.Tasks[i] = *task
				break
			}
		}

		if err := tm.SaveTasks("backlog", backlog); err != nil {
			fmt.Printf("Error saving task: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Sample task created successfully!\n\n")
		fmt.Printf("Task ID: %s\n", task.ID)
		fmt.Printf("Title: %s\n\n", task.Title)
		fmt.Printf("Next steps:\n")
		fmt.Printf("1. View task: agentic-agent task show %s\n", task.ID)
		fmt.Printf("2. Claim task: agentic-agent task claim %s\n", task.ID)
		fmt.Printf("3. List all tasks: agentic-agent task list\n")
	},
}

var taskFromTemplateCmd = &cobra.Command{
	Use:   "from-template",
	Short: "Create a task from a saved template",
	Long: `Create a new task from a saved template.

Interactive Mode:
  agentic-agent task from-template

  Launches an interactive wizard that:
  - Shows available built-in templates
  - Lets you select a template
  - Allows customization before creation

Flag Mode:
  agentic-agent task from-template --template <template-name> --title "Task title" [options]

Built-in templates include:
  - feature: Feature Implementation
  - bug-fix: Bug Fix
  - refactoring: Refactoring
  - documentation: Documentation
  - testing: Testing`,
	Run: func(cmd *cobra.Command, args []string) {
		templateName, _ := cmd.Flags().GetString("template")

		// Interactive mode if no template specified
		if helpers.ShouldUseInteractiveMode(cmd) && templateName == "" {
			runTemplateWorkflow()
			return
		}

		// Flag mode - require template and title
		if templateName == "" {
			fmt.Println("Error: --template required in non-interactive mode")
			fmt.Println("Usage: agentic-agent task from-template --template <name> --title \"Task title\"")
			fmt.Println("   or: agentic-agent task from-template  (interactive mode)")
			fmt.Println("\nAvailable templates: feature, bug-fix, refactoring, documentation, testing")
			os.Exit(1)
		}

		title, _ := cmd.Flags().GetString("title")
		if title == "" {
			fmt.Println("Error: --title required when using --template flag")
			os.Exit(1)
		}

		// Get other optional flags
		description, _ := cmd.Flags().GetString("description")
		scope, _ := cmd.Flags().GetString("scope")
		specRefs, _ := cmd.Flags().GetString("spec-refs")
		inputs, _ := cmd.Flags().GetString("inputs")
		outputs, _ := cmd.Flags().GetString("outputs")
		acceptance, _ := cmd.Flags().GetString("acceptance")

		// Create task manager
		tm := tasks.NewTaskManager(".agentic/tasks")

		// Create new task
		task, err := tm.CreateTask(title)
		if err != nil {
			fmt.Printf("Error creating task: %v\n", err)
			os.Exit(1)
		}

		// Apply template defaults
		switch templateName {
		case "feature":
			if description == "" {
				description = "Implement new feature"
			}
		case "bug-fix":
			if description == "" {
				description = "Fix bug in the system"
			}
		case "refactoring":
			if description == "" {
				description = "Refactor existing code"
			}
		case "documentation":
			if description == "" {
				description = "Add or update documentation"
			}
		case "testing":
			if description == "" {
				description = "Add tests"
			}
		default:
			fmt.Printf("Error: unknown template '%s'\n", templateName)
			fmt.Println("Available templates: feature, bug-fix, refactoring, documentation, testing")
			os.Exit(1)
		}

		// Apply user-provided values
		task.Description = description

		if scope != "" {
			task.Scope = strings.Split(scope, ",")
		}
		if specRefs != "" {
			task.SpecRefs = strings.Split(specRefs, ",")
		}
		if inputs != "" {
			task.Inputs = strings.Split(inputs, ",")
		}
		if outputs != "" {
			task.Outputs = strings.Split(outputs, ",")
		}
		if acceptance != "" {
			task.Acceptance = strings.Split(acceptance, ",")
		}

		// Save task
		backlog, err := tm.LoadTasks("backlog")
		if err != nil {
			fmt.Printf("Error loading backlog: %v\n", err)
			os.Exit(1)
		}

		for i, t := range backlog.Tasks {
			if t.ID == task.ID {
				backlog.Tasks[i] = *task
				break
			}
		}

		if err := tm.SaveTasks("backlog", backlog); err != nil {
			fmt.Printf("Error saving task: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Created task %s from template '%s'\n", task.ID, templateName)
	},
}

// runTemplateWorkflow runs the template selection workflow
func runTemplateWorkflow() {
	model := uimodels.NewTemplateSelectModel()
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running template workflow: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	taskCreateCmd.Flags().String("title", "", "Task title (required)")
	taskCreateCmd.Flags().String("description", "", "Task description")
	taskCreateCmd.Flags().String("scope", "", "Comma-separated list of directories/files in scope")
	taskCreateCmd.Flags().String("spec-refs", "", "Comma-separated specification file references")
	taskCreateCmd.Flags().String("inputs", "", "Comma-separated required input files")
	taskCreateCmd.Flags().String("outputs", "", "Comma-separated expected output files")
	taskCreateCmd.Flags().String("acceptance", "", "Comma-separated acceptance criteria")

	// from-template flags
	taskFromTemplateCmd.Flags().String("template", "", "Template name (feature, bug-fix, refactoring, documentation, testing)")
	taskFromTemplateCmd.Flags().String("title", "", "Task title (required with --template)")
	taskFromTemplateCmd.Flags().String("description", "", "Task description (optional)")
	taskFromTemplateCmd.Flags().String("scope", "", "Comma-separated list of directories/files in scope")
	taskFromTemplateCmd.Flags().String("spec-refs", "", "Comma-separated specification file references")
	taskFromTemplateCmd.Flags().String("inputs", "", "Comma-separated required input files")
	taskFromTemplateCmd.Flags().String("outputs", "", "Comma-separated expected output files")
	taskFromTemplateCmd.Flags().String("acceptance", "", "Comma-separated acceptance criteria")

	taskCmd.AddCommand(taskCreateCmd)
	taskCmd.AddCommand(taskSampleCmd)
	taskCmd.AddCommand(taskFromTemplateCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskShowCmd)
	taskCmd.AddCommand(taskClaimCmd)
	taskCmd.AddCommand(taskCompleteCmd)
	taskCmd.AddCommand(taskDecomposeCmd)
}
