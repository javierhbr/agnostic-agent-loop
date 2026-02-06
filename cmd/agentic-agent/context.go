package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/context"
	"github.com/javierbenavides/agentic-agent/internal/encoding"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/internal/ui/components"
	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Manage context files",
}

// contextGenerateModel is a Bubble Tea model for directory selection and context generation
type contextGenerateModel struct {
	picker  components.FilePicker
	step    string // "select" or "done"
	success bool
	message string
}

func (m contextGenerateModel) Init() tea.Cmd {
	return nil
}

func (m contextGenerateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			if m.step == "done" {
				return m, tea.Quit
			}
			return m, tea.Quit

		case "enter":
			if m.step == "select" {
				selected := m.picker.GetSelected()
				if len(selected) == 0 {
					m.message = "No directory selected"
					m.success = false
					m.step = "done"
					return m, tea.Quit
				}

				dir := selected[0]

				// Generate context
				ctx, err := context.GenerateContext(dir)
				if err != nil {
					m.message = fmt.Sprintf("Error generating context: %v", err)
					m.success = false
					m.step = "done"
					return m, tea.Quit
				}

				dcm := context.NewDirectoryContextManager(dir)
				if err := dcm.SaveContext(dir, ctx); err != nil {
					m.message = fmt.Sprintf("Error saving context: %v", err)
					m.success = false
					m.step = "done"
					return m, tea.Quit
				}

				m.message = fmt.Sprintf("Generated context for %s", dir)
				m.success = true
				m.step = "done"
				return m, tea.Quit
			}

		default:
			var cmd tea.Cmd
			m.picker, cmd = m.picker.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m contextGenerateModel) View() string {
	if m.step == "done" {
		if m.success {
			return styles.RenderSuccess(m.message) + "\n"
		}
		return styles.RenderError(m.message) + "\n"
	}

	var b strings.Builder
	b.WriteString(styles.TitleStyle.Render("Select Directory for Context Generation") + "\n\n")
	b.WriteString(m.picker.View() + "\n")

	return styles.ContainerStyle.Render(b.String())
}

// contextBuildModel is a Bubble Tea model for task and format selection
type contextBuildModel struct {
	step         string // "select-task", "select-format", "building", "done"
	taskSelector components.SimpleSelect
	tasks        []string
	selectedTask string
	format       string
	success      bool
	message      string
}

func (m *contextBuildModel) Init() tea.Cmd {
	// Load tasks for selection
	tm := tasks.NewTaskManager(".agentic/tasks")

	// Get tasks from all statuses
	var allTasks []string
	for _, status := range []string{"backlog", "in-progress", "done"} {
		taskList, _ := tm.LoadTasks(status)
		for _, t := range taskList.Tasks {
			allTasks = append(allTasks, t.ID)
		}
	}

	m.tasks = allTasks

	// Create task selector
	options := []components.SelectOption{}
	for _, taskID := range allTasks {
		options = append(options, components.NewSelectOption(taskID, "Build context for this task", taskID))
	}

	if len(options) == 0 {
		options = append(options, components.NewSelectOption("No tasks found", "Create a task first", ""))
	}

	m.taskSelector = components.NewSimpleSelect("Select task for context bundle", options)

	return nil
}

func (m *contextBuildModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			if m.step == "select-format" {
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
				m.step = "select-format"
			} else if m.step == "select-format" {
				// Build the context bundle
				m.step = "building"
				bundle, err := encoding.CreateContextBundle(m.selectedTask, m.format, getConfig())
				if err != nil {
					m.message = fmt.Sprintf("Error building bundle: %v", err)
					m.success = false
				} else {
					m.message = string(bundle)
					m.success = true
				}
				m.step = "done"
				return m, tea.Quit
			}

		case "1":
			if m.step == "select-format" {
				m.format = "toon"
			}
		case "2":
			if m.step == "select-format" {
				m.format = "markdown"
			}
		case "3":
			if m.step == "select-format" {
				m.format = "json"
			}

		default:
			if m.step == "select-task" {
				m.taskSelector = m.taskSelector.Update(msg)
			}
		}
	}

	return m, nil
}

func (m *contextBuildModel) View() string {
	if m.step == "done" {
		if m.success {
			return m.message + "\n"
		}
		return styles.RenderError(m.message) + "\n"
	}

	var b strings.Builder

	if m.step == "select-task" {
		b.WriteString(styles.TitleStyle.Render("Build Context Bundle") + "\n\n")
		b.WriteString(m.taskSelector.View() + "\n")
		b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Enter select • q quit") + "\n")
	} else if m.step == "select-format" {
		b.WriteString(styles.TitleStyle.Render("Select Output Format") + "\n\n")
		b.WriteString(fmt.Sprintf("Task: %s\n\n", styles.BoldStyle.Render(m.selectedTask)))

		formats := []struct {
			key  string
			name string
			desc string
		}{
			{"1", "toon", "TOON format (default)"},
			{"2", "markdown", "Markdown format"},
			{"3", "json", "JSON format"},
		}

		for _, f := range formats {
			selected := ""
			if f.name == m.format {
				selected = styles.IconArrow + " "
			}
			line := fmt.Sprintf("%s%s. %s - %s\n",
				selected,
				f.key,
				styles.BoldStyle.Render(f.name),
				styles.MutedStyle.Render(f.desc),
			)
			b.WriteString(line)
		}

		b.WriteString("\n")
		b.WriteString(styles.HelpStyle.Render("1-3 select format • Enter confirm • Esc back") + "\n")
	} else if m.step == "building" {
		b.WriteString(styles.TitleStyle.Render("Building Context Bundle...") + "\n")
	}

	return styles.ContainerStyle.Render(b.String())
}

var contextGenerateCmd = &cobra.Command{
	Use:   "generate [dir]",
	Short: "Generate context.md for a directory",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Interactive mode - use file picker
		if helpers.ShouldUseInteractiveMode(cmd) && len(args) == 0 {
			cwd, _ := os.Getwd()
			picker := components.NewFilePicker("Select directory", cwd, true, false) // directories only, single select

			model := &contextGenerateModel{
				picker: picker,
				step:   "select",
			}

			p := tea.NewProgram(model)
			if _, err := p.Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		// Flag mode - require directory argument
		if len(args) != 1 {
			fmt.Println("Error: directory argument required in non-interactive mode")
			fmt.Println("Usage: agentic-agent context generate <dir>")
			fmt.Println("   or: agentic-agent context generate  (interactive mode)")
			os.Exit(1)
		}

		dir := args[0]
		ctx, err := context.GenerateContext(dir)
		if err != nil {
			fmt.Printf("Error generating context: %v\n", err)
			os.Exit(1)
		}

		dcm := context.NewDirectoryContextManager(dir)
		if err := dcm.SaveContext(dir, ctx); err != nil {
			fmt.Printf("Error saving context: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Generated context for %s\n", dir)
	},
}

var contextScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for directories needing context",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()
		dcm := context.NewDirectoryContextManager(cwd)
		dirs, err := dcm.FindContextDirs(cwd)
		if err != nil {
			fmt.Printf("Error scanning: %v\n", err)
			os.Exit(1)
		}

		// Interactive mode - prettier output with colors
		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder

			b.WriteString(styles.TitleStyle.Render("Context Directory Scan") + "\n\n")

			missingCount := 0
			okCount := 0

			for _, d := range dirs {
				if _, err := os.Stat(d + "/context.md"); os.IsNotExist(err) {
					missingCount++
					line := fmt.Sprintf("%s [%s] %s",
						styles.IconCross,
						styles.ErrorStyle.Render("MISSING"),
						styles.MutedStyle.Render(d),
					)
					b.WriteString(line + "\n")
				} else {
					okCount++
					line := fmt.Sprintf("%s [%s]      %s",
						styles.IconCheckmark,
						styles.SuccessStyle.Render("OK"),
						d,
					)
					b.WriteString(line + "\n")
				}
			}

			b.WriteString("\n")
			summary := fmt.Sprintf("Total: %s directories | %s with context | %s missing context",
				styles.BoldStyle.Render(fmt.Sprintf("%d", len(dirs))),
				styles.SuccessStyle.Render(fmt.Sprintf("%d", okCount)),
				styles.ErrorStyle.Render(fmt.Sprintf("%d", missingCount)),
			)
			b.WriteString(styles.CardStyle.Render(summary) + "\n")

			if missingCount > 0 {
				b.WriteString("\n" + styles.HelpStyle.Render("Tip: Use 'agentic-agent context generate <dir>' to create missing context files") + "\n")
			}

			fmt.Println(styles.ContainerStyle.Render(b.String()))
			return
		}

		// Flag mode - simple text output
		for _, d := range dirs {
			if _, err := os.Stat(d + "/context.md"); os.IsNotExist(err) {
				fmt.Printf("[MISSING] %s\n", d)
			} else {
				fmt.Printf("[OK]      %s\n", d)
			}
		}
	},
}

var contextUpdateCmd = &cobra.Command{
	Use:   "update [dir]",
	Short: "Update context.md for a directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// logic same as generate for now
		contextGenerateCmd.Run(cmd, args)
	},
}

var contextBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build context bundle",
	Run: func(cmd *cobra.Command, args []string) {
		taskID, _ := cmd.Flags().GetString("task")
		format, _ := cmd.Flags().GetString("format")

		// Interactive mode - task and format selection
		if helpers.ShouldUseInteractiveMode(cmd) && taskID == "" {
			model := &contextBuildModel{
				step:   "select-task",
				format: format,
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
			fmt.Println("Usage: agentic-agent context build --task <task-id> [--format toon]")
			fmt.Println("   or: agentic-agent context build  (interactive mode)")
			os.Exit(1)
		}

		bundle, err := encoding.CreateContextBundle(taskID, format, getConfig())
		if err != nil {
			fmt.Printf("Error building bundle: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(bundle))
	},
}

func init() {
	contextBuildCmd.Flags().String("task", "", "Task ID to build context for")
	contextBuildCmd.Flags().String("format", "toon", "Output format")

	contextCmd.AddCommand(contextGenerateCmd)
	contextCmd.AddCommand(contextScanCmd)
	contextCmd.AddCommand(contextUpdateCmd)
	contextCmd.AddCommand(contextBuildCmd)
}
