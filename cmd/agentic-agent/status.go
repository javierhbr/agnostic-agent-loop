package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/javierbenavides/agentic-agent/internal/status"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show project status dashboard",
	Long:  "Display an overview of project progress: task counts, completion percentage, in-progress work, blockers, and recent activity.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := getConfig()
		tm := tasks.NewTaskManager(".agentic/tasks")

		data, err := status.Gather(tm, cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error gathering status: %v\n", err)
			os.Exit(1)
		}

		format, _ := cmd.Flags().GetString("format")

		if format == "json" {
			out, _ := json.MarshalIndent(data, "", "  ")
			fmt.Println(string(out))
			return
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			renderInteractiveStatus(data)
		} else {
			renderPlainStatus(data)
		}
	},
}

func renderInteractiveStatus(d *status.DashboardData) {
	var b strings.Builder

	// Title
	projectName := d.ProjectName
	if projectName == "" {
		projectName = "Project"
	}
	b.WriteString(styles.TitleStyle.Render(projectName+" Status") + "\n\n")

	// Progress bar
	bar := renderProgressBar(d.CompletionPct, 30)
	b.WriteString(fmt.Sprintf("  %s  %s\n\n",
		bar,
		styles.BoldStyle.Render(fmt.Sprintf("%.0f%% complete", d.CompletionPct)),
	))

	// Task counts
	counts := fmt.Sprintf(
		"  %s %s backlog  %s %s in progress  %s %s done",
		styles.IconPending,
		styles.MutedStyle.Render(fmt.Sprintf("%d", d.BacklogCount)),
		styles.IconProgress,
		styles.WarningStyle.Render(fmt.Sprintf("%d", d.InProgressCount)),
		styles.IconCheckmark,
		styles.SuccessStyle.Render(fmt.Sprintf("%d", d.DoneCount)),
	)
	b.WriteString(counts + "\n\n")

	// In-progress tasks
	if len(d.InProgressTasks) > 0 {
		b.WriteString(styles.SubtitleStyle.Render("In Progress") + "\n")
		for _, t := range d.InProgressTasks {
			assignee := ""
			if t.AssignedTo != "" {
				assignee = styles.MutedStyle.Render(" ("+t.AssignedTo+")")
			}
			b.WriteString(fmt.Sprintf("  %s %s%s\n",
				styles.IconArrow,
				styles.BoldStyle.Render(t.ID+": "+t.Title),
				assignee,
			))
		}
		b.WriteString("\n")
	}

	// Next ready task
	if d.NextReady != nil {
		b.WriteString(styles.SubtitleStyle.Render("Next Up") + "\n")
		b.WriteString(fmt.Sprintf("  %s %s\n\n",
			styles.IconPrompt,
			d.NextReady.ID+": "+d.NextReady.Title,
		))
	}

	// Blockers
	if len(d.Blockers) > 0 {
		b.WriteString(styles.ErrorStyle.Render("Blockers") + "\n")
		for _, blocker := range d.Blockers {
			b.WriteString(fmt.Sprintf("  %s %s\n", styles.IconCross, styles.MutedStyle.Render(blocker)))
		}
		b.WriteString("\n")
	}

	// Recent activity
	if len(d.RecentEntries) > 0 {
		b.WriteString(styles.SubtitleStyle.Render("Recent Activity") + "\n")
		for _, entry := range d.RecentEntries {
			b.WriteString(fmt.Sprintf("  %s %s %s\n",
				styles.IconBullet,
				styles.MutedStyle.Render(entry.Timestamp.Format("Jan 02")),
				entry.Title,
			))
		}
	}

	fmt.Println(styles.ContainerStyle.Render(b.String()))
}

func renderPlainStatus(d *status.DashboardData) {
	projectName := d.ProjectName
	if projectName == "" {
		projectName = "Project"
	}
	fmt.Printf("%s Status\n", projectName)
	fmt.Printf("Progress: %.0f%% (%d/%d tasks done)\n", d.CompletionPct, d.DoneCount, d.TotalCount)
	fmt.Printf("Backlog: %d | In Progress: %d | Done: %d\n\n", d.BacklogCount, d.InProgressCount, d.DoneCount)

	if len(d.InProgressTasks) > 0 {
		fmt.Println("In Progress:")
		for _, t := range d.InProgressTasks {
			fmt.Printf("  -> %s: %s\n", t.ID, t.Title)
		}
		fmt.Println()
	}

	if d.NextReady != nil {
		fmt.Printf("Next Ready: %s: %s\n\n", d.NextReady.ID, d.NextReady.Title)
	}

	if len(d.Blockers) > 0 {
		fmt.Println("Blockers:")
		for _, b := range d.Blockers {
			fmt.Printf("  ! %s\n", b)
		}
		fmt.Println()
	}

	if len(d.RecentEntries) > 0 {
		fmt.Println("Recent Activity:")
		for _, entry := range d.RecentEntries {
			fmt.Printf("  %s  %s\n", entry.Timestamp.Format("Jan 02"), entry.Title)
		}
	}
}

func renderProgressBar(pct float64, width int) string {
	filled := int(pct / 100 * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled
	bar := styles.SuccessStyle.Render(strings.Repeat("█", filled)) +
		styles.MutedStyle.Render(strings.Repeat("░", empty))
	return "[" + bar + "]"
}

func init() {
	statusCmd.Flags().String("format", "text", "Output format (text|json)")
}
