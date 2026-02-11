package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/javierbenavides/agentic-agent/internal/plans"
	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Manage phased implementation plans",
	Long:  "Parse, display, and update markdown-based implementation plans with checkbox status markers.",
}

var planShowCmd = &cobra.Command{
	Use:   "show <path>",
	Short: "Display plan with progress",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := resolvePlanPath(cmd, args)

		plan, err := plans.ParseFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		done, total := plan.Progress()

		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder
			title := plan.Title
			if title == "" {
				title = "Plan"
			}
			b.WriteString(styles.TitleStyle.Render(title) + "\n\n")

			pct := float64(0)
			if total > 0 {
				pct = float64(done) / float64(total) * 100
			}
			b.WriteString(fmt.Sprintf("  Progress: %d/%d tasks (%.0f%%)\n\n", done, total, pct))

			for _, phase := range plan.Phases {
				b.WriteString("  " + styles.SubtitleStyle.Render(phase.Name) + "\n")
				for _, task := range phase.Tasks {
					icon := styles.IconPending
					var titleStr string
					switch task.Status {
					case plans.PlanTaskDone:
						icon = styles.IconCheckmark
						titleStr = styles.MutedStyle.Render(task.Title)
					case plans.PlanTaskInProgress:
						icon = styles.IconProgress
						titleStr = styles.WarningStyle.Render(task.Title)
					default:
						titleStr = task.Title
					}
					b.WriteString(fmt.Sprintf("    %s %s\n", icon, titleStr))
				}
				b.WriteString("\n")
			}
			fmt.Println(styles.ContainerStyle.Render(b.String()))
		} else {
			fmt.Printf("%s (%d/%d done)\n\n", plan.Title, done, total)
			for _, phase := range plan.Phases {
				fmt.Printf("## %s\n", phase.Name)
				for _, task := range phase.Tasks {
					marker := " "
					if task.Status == plans.PlanTaskDone {
						marker = "x"
					} else if task.Status == plans.PlanTaskInProgress {
						marker = "~"
					}
					fmt.Printf("  - [%s] %s\n", marker, task.Title)
				}
				fmt.Println()
			}
		}
	},
}

var planNextCmd = &cobra.Command{
	Use:   "next [path]",
	Short: "Show the next pending task",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := resolvePlanPath(cmd, args)

		plan, err := plans.ParseFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		task, phase := plan.NextTask()
		if task == nil {
			fmt.Println("All tasks complete!")
			return
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			fmt.Printf("  %s %s\n", styles.SubtitleStyle.Render(phase.Name), styles.IconArrow)
			fmt.Printf("  %s %s\n", styles.IconPrompt, styles.BoldStyle.Render(task.Title))
		} else {
			fmt.Printf("%s: %s\n", phase.Name, task.Title)
		}
	},
}

var planMarkCmd = &cobra.Command{
	Use:   "mark <path> <line> <status>",
	Short: "Update a task's status by line number",
	Long:  "Mark a plan task as pending, in_progress, or done. Use 'plan show' to see line numbers.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		lineNum, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid line number: %s\n", args[1])
			os.Exit(1)
		}

		var status plans.PlanTaskStatus
		switch args[2] {
		case "done", "x":
			status = plans.PlanTaskDone
		case "pending", " ":
			status = plans.PlanTaskPending
		case "in_progress", "~", "wip":
			status = plans.PlanTaskInProgress
		default:
			fmt.Fprintf(os.Stderr, "Invalid status: %s (use: done, pending, in_progress)\n", args[2])
			os.Exit(1)
		}

		if err := plans.UpdateTaskStatus(path, lineNum, status); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Updated line %d to %s\n", lineNum, status)
	},
}

func resolvePlanPath(cmd *cobra.Command, args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	// Try to find plan.md in the current track or default location
	track, _ := cmd.Flags().GetString("track")
	if track != "" {
		cfg := getConfig()
		return filepath.Join(cfg.Paths.TrackDir, track, "plan.md")
	}
	return "plan.md"
}

func init() {
	planShowCmd.Flags().String("track", "", "Track ID to show plan for")
	planNextCmd.Flags().String("track", "", "Track ID to show next task for")

	planCmd.AddCommand(planShowCmd)
	planCmd.AddCommand(planNextCmd)
	planCmd.AddCommand(planMarkCmd)
}
