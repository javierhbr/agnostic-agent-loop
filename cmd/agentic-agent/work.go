package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/context"
	"github.com/javierbenavides/agentic-agent/internal/skills"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	uimodels "github.com/javierbenavides/agentic-agent/internal/ui/models"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/spf13/cobra"
)

var workCmd = &cobra.Command{
	Use:   "work",
	Short: "Interactive workflow: claim task → work → complete",
	Long: `Complete interactive workflow for working on tasks.

This command guides you through the entire workflow:
1. Select a task from the backlog
2. Claim the task (moves to in-progress)
3. Optionally generate context for scope directories
4. View task details with acceptance criteria
5. Mark as complete when done
6. Run validation
7. Move to done

Interactive Mode:
  agentic-agent work

Flag Mode:
  agentic-agent work --task <task-id> [--skip-context-gen] [--follow-tdd]`,
	Run: func(cmd *cobra.Command, args []string) {
		taskID, _ := cmd.Flags().GetString("task")
		skipContextGen, _ := cmd.Flags().GetBool("skip-context-gen")
		followTDD, _ := cmd.Flags().GetBool("follow-tdd")

		// Interactive mode if no task specified
		if helpers.ShouldUseInteractiveMode(cmd) && taskID == "" {
			runWorkWorkflow()
			return
		}

		// Flag mode - require task ID
		if taskID == "" {
			fmt.Println("Error: --task required in non-interactive mode")
			fmt.Println("Usage: agentic-agent work --task <task-id> [--skip-context-gen]")
			fmt.Println("   or: agentic-agent work  (interactive mode)")
			os.Exit(1)
		}

		// Execute workflow programmatically
		tm := tasks.NewTaskManager(".agentic/tasks")

		// 1. Find and claim the task
		task, source, err := tm.FindTask(taskID)
		if err != nil {
			fmt.Printf("Error finding task: %v\n", err)
			os.Exit(1)
		}
		if task == nil {
			fmt.Printf("Task %s not found\n", taskID)
			os.Exit(1)
		}

		// Get current user
		user := os.Getenv("USER")
		if user == "" {
			user = "unknown-agent"
		}

		// Claim task if it's in backlog
		if source == "backlog" {
			if err := tm.ClaimTask(taskID, user); err != nil {
				fmt.Printf("Error claiming task: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✓ Claimed task %s\n", taskID)
		} else {
			fmt.Printf("Task %s is already in %s\n", taskID, source)
		}

		// 2. TDD workflow: decompose task and verify skill pack installed
		if followTDD {
			installer := skills.NewInstaller()
			tool := installer.IsInstalledAnywhere("tdd")
			if tool == "" {
				fmt.Println("Error: TDD skill pack not installed.")
				fmt.Println("Install it first: agentic-agent skills install tdd --tool <tool>")
				os.Exit(1)
			}

			subtasks, err := tasks.DecomposeForTDD(tm, taskID)
			if err != nil {
				fmt.Printf("Error decomposing task for TDD: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("TDD workflow enabled (skill pack found for %s)\n", tool)
			fmt.Println("Created TDD sub-tasks:")
			for _, st := range subtasks {
				fmt.Printf("  - %s: %s\n", st.ID, st.Title)
			}
			fmt.Println()
		}

		// 3. Generate context for scope directories if not skipped
		if !skipContextGen && len(task.Scope) > 0 {
			fmt.Println("Generating context for scope directories...")
			for _, dir := range task.Scope {
				ctx, err := context.GenerateContextWithConfig(dir, getConfig())
				if err != nil {
					fmt.Printf("  Warning: Could not generate context for %s: %v\n", dir, err)
					continue
				}
				dcm := context.NewDirectoryContextManager(dir)
				if err := dcm.SaveContext(dir, ctx); err != nil {
					fmt.Printf("  Warning: Could not save context for %s: %v\n", dir, err)
					continue
				}
				fmt.Printf("  ✓ Generated context for %s\n", dir)
			}
		}

		// 4. In TDD mode, don't auto-complete — user works through sub-tasks
		if followTDD {
			fmt.Println("TDD sub-tasks created. Work through RED → GREEN → REFACTOR phases.")
			fmt.Println("\nWorkflow ready!")
			return
		}

		// 5. Complete the task
		if err := tm.MoveTask(taskID, "in-progress", "done", models.StatusDone); err != nil {
			fmt.Printf("Error completing task: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Completed task %s\n", taskID)
		fmt.Println("\nWorkflow complete!")
	},
}

// runWorkWorkflow runs the complete work workflow
func runWorkWorkflow() {
	model := uimodels.NewWorkflowModel()
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running work workflow: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	workCmd.Flags().String("task", "", "Task ID to work on")
	workCmd.Flags().Bool("skip-context-gen", false, "Skip context generation for scope directories")
	workCmd.Flags().Bool("follow-tdd", false, "Follow TDD workflow: decompose task into RED/GREEN/REFACTOR phases")

	rootCmd.AddCommand(workCmd)
}
