package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/javierbenavides/agentic-agent/internal/config"
	"github.com/javierbenavides/agentic-agent/internal/openspec"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/spf13/cobra"
)

// syncOpenSpecChanges auto-imports tasks from any draft openspec changes
// that have a populated tasks.md. Returns the sync result for optional logging.
func syncOpenSpecChanges(cfg *models.Config) *openspec.SyncResult {
	if cfg.Paths.OpenSpecDir == "" {
		return &openspec.SyncResult{}
	}
	m := openspec.NewManager(cfg.Paths.OpenSpecDir)
	tm := tasks.NewTaskManager(".agentic/tasks")
	result, _ := m.Sync(tm)
	if result != nil && len(result.ChangesImported) > 0 {
		fmt.Printf("Auto-imported %d tasks from %d change(s)\n",
			result.TasksCreated, len(result.ChangesImported))
	}
	return result
}

var openspecCmd = &cobra.Command{
	Use:   "openspec",
	Short: "Manage spec-driven development changes",
	Long:  "Create, import, track, and archive openspec change proposals.",
}

var openspecInitCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Create a new change from a requirements file",
	Long: `Create a new openspec change directory with proposal.md and tasks.md.

The --from flag seeds the proposal with the contents of the requirements file.
The agent then fills in the proposal and writes tasks.md.

Examples:
  agentic-agent openspec init "Auth Feature" --from requirements.md
  agentic-agent openspec init "Payment System" --from docs/payment-spec.md`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := getConfig()

		// Ensure agnostic-agent.yaml exists in the project root
		cwd, _ := os.Getwd()
		configResult, err := openspec.EnsureConfig(cwd, "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not create config: %v\n", err)
		} else if configResult.Created {
			if helpers.ShouldUseInteractiveMode(cmd) {
				fmt.Println(styles.RenderSuccess("Created agnostic-agent.yaml"))
			} else {
				fmt.Println("Created agnostic-agent.yaml")
			}
			// Reload config so paths.openSpecDir is set
			reloaded, loadErr := config.LoadConfig(configResult.Path)
			if loadErr == nil {
				cfg = reloaded
				appConfig = cfg
			}
		}

		m := openspec.NewManager(cfg.Paths.OpenSpecDir)

		fromFile, _ := cmd.Flags().GetString("from")
		if fromFile == "" {
			fmt.Fprintln(os.Stderr, "Error: --from flag is required")
			os.Exit(1)
		}

		change, err := m.Init(args[0], fromFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		changeDir := filepath.Join(cfg.Paths.OpenSpecDir, change.ID)
		if helpers.ShouldUseInteractiveMode(cmd) {
			fmt.Println(styles.RenderSuccess(fmt.Sprintf("Created change: %s", change.ID)))
			fmt.Printf("  %s Proposal:  %s/proposal.md\n", styles.IconArrow, changeDir)
			fmt.Printf("  %s Tasks:     %s/tasks.md\n", styles.IconArrow, changeDir)
			fmt.Printf("  %s Status:    %s\n", styles.IconArrow, change.Status)
			fmt.Printf("\n%s Fill in proposal.md, then write tasks in tasks.md.\n", styles.IconArrow)
			fmt.Printf("%s For complex changes (4+ tasks), write specs in %s/specs/\n", styles.IconArrow, changeDir)
			fmt.Printf("%s Tasks auto-import when you run: agentic-agent task list  or  task claim\n", styles.IconArrow)
		} else {
			fmt.Printf("Created change: %s\n", change.ID)
			fmt.Printf("proposal: %s/proposal.md\n", changeDir)
			fmt.Printf("tasks: %s/tasks.md\n", changeDir)
		}
	},
}

var openspecImportCmd = &cobra.Command{
	Use:        "import <id>",
	Short:      "Import tasks from tasks.md into the backlog (deprecated)",
	Hidden:     true,
	Deprecated: "tasks are now auto-imported when you run 'task list' or 'task claim'.",
	Long: `Parse tasks.md for the given change and create tasks in the backlog.

DEPRECATED: Tasks are now auto-imported when you run 'task list' or 'task claim'.
This command is kept for backward compatibility.

Supports numbered lists (1. Task) and checkbox lists (- [ ] Task).

Examples:
  agentic-agent openspec import auth-feature`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := getConfig()
		m := openspec.NewManager(cfg.Paths.OpenSpecDir)
		tm := tasks.NewTaskManager(".agentic/tasks")

		created, err := m.Import(args[0], tm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			fmt.Println(styles.RenderSuccess(fmt.Sprintf("Imported %d tasks from %s", len(created), args[0])))
			for _, t := range created {
				fmt.Printf("  %s %s: %s\n", styles.IconBullet, t.ID, t.Title)
			}
			fmt.Printf("\n%s Run: agentic-agent task claim %s\n", styles.IconArrow, created[0].ID)
		} else {
			fmt.Printf("Imported %d tasks\n", len(created))
			for _, t := range created {
				fmt.Printf("%s: %s\n", t.ID, t.Title)
			}
		}
	},
}

var openspecListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all changes",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := getConfig()
		m := openspec.NewManager(cfg.Paths.OpenSpecDir)

		changes, err := m.List()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(changes) == 0 {
			fmt.Println("No changes yet. Create one with: agentic-agent openspec init \"My Feature\" --from requirements.md")
			return
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder
			b.WriteString(styles.TitleStyle.Render("OpenSpec Changes") + "\n\n")
			for _, c := range changes {
				icon := changeStatusIcon(c.Status)
				taskCount := fmt.Sprintf("%d tasks", len(c.TaskIDs))
				b.WriteString(fmt.Sprintf("  %s %s  %s  %s\n",
					icon,
					styles.BoldStyle.Render(c.ID),
					changeStatusStyle(c.Status),
					styles.MutedStyle.Render(taskCount),
				))
			}
			fmt.Println(styles.ContainerStyle.Render(b.String()))
		} else {
			for _, c := range changes {
				fmt.Printf("%s  %s  %d tasks\n", c.ID, c.Status, len(c.TaskIDs))
			}
		}
	},
}

var openspecShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show change details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := getConfig()
		m := openspec.NewManager(cfg.Paths.OpenSpecDir)

		change, err := m.Get(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder
			b.WriteString(styles.TitleStyle.Render(change.Name) + "\n\n")
			b.WriteString(fmt.Sprintf("  ID:      %s\n", change.ID))
			b.WriteString(fmt.Sprintf("  Status:  %s\n", changeStatusStyle(change.Status)))
			b.WriteString(fmt.Sprintf("  Source:  %s\n", change.SourceFile))
			b.WriteString(fmt.Sprintf("  Created: %s\n", change.CreatedAt.Format("2006-01-02")))
			if len(change.TaskIDs) > 0 {
				b.WriteString(fmt.Sprintf("\n  %s Tasks:\n", styles.SubtitleStyle.Render("Associated")))
				for _, tid := range change.TaskIDs {
					b.WriteString(fmt.Sprintf("    %s %s\n", styles.IconBullet, tid))
				}
			}
			fmt.Println(styles.ContainerStyle.Render(b.String()))
		} else {
			fmt.Printf("ID: %s\nName: %s\nStatus: %s\nSource: %s\nCreated: %s\n",
				change.ID, change.Name, change.Status, change.SourceFile, change.CreatedAt.Format("2006-01-02"))
			if len(change.TaskIDs) > 0 {
				fmt.Printf("Tasks: %s\n", strings.Join(change.TaskIDs, ", "))
			}
		}
	},
}

var openspecStatusCmd = &cobra.Command{
	Use:   "status <id>",
	Short: "Show change progress",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := getConfig()
		m := openspec.NewManager(cfg.Paths.OpenSpecDir)
		tm := tasks.NewTaskManager(".agentic/tasks")

		change, err := m.Get(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		progress, err := m.Progress(args[0], tm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder
			b.WriteString(styles.TitleStyle.Render("Change: "+change.Name) + "\n\n")
			b.WriteString(fmt.Sprintf("  Status:      %s\n", changeStatusStyle(change.Status)))
			b.WriteString(fmt.Sprintf("  Total:       %d\n", progress.Total))
			b.WriteString(fmt.Sprintf("  %s Done:       %d\n", styles.IconCheckmark, progress.Done))
			b.WriteString(fmt.Sprintf("  %s In Progress: %d\n", styles.IconProgress, progress.InProgress))
			b.WriteString(fmt.Sprintf("  %s Pending:     %d\n", styles.IconPending, progress.Pending))

			if progress.Done == progress.Total && progress.Total > 0 {
				b.WriteString(fmt.Sprintf("\n%s All tasks done! Run: agentic-agent openspec complete %s\n", styles.IconCheckmark, args[0]))
			}
			fmt.Println(styles.ContainerStyle.Render(b.String()))
		} else {
			fmt.Printf("change: %s\nstatus: %s\ntotal: %d\ndone: %d\nin_progress: %d\npending: %d\n",
				change.ID, change.Status, progress.Total, progress.Done, progress.InProgress, progress.Pending)
			for _, tid := range progress.TaskIDs {
				fmt.Printf("task: %s\n", tid)
			}
		}
	},
}

var openspecCompleteCmd = &cobra.Command{
	Use:   "complete <id>",
	Short: "Mark a change as implemented",
	Long: `Validates all tasks are done and writes the IMPLEMENTED marker.

Fails if any tasks are still pending or in progress.

Examples:
  agentic-agent openspec complete auth-feature`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := getConfig()
		m := openspec.NewManager(cfg.Paths.OpenSpecDir)
		tm := tasks.NewTaskManager(".agentic/tasks")

		if err := m.Complete(args[0], tm); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			fmt.Println(styles.RenderSuccess(fmt.Sprintf("Change %s marked as implemented", args[0])))
			fmt.Printf("  %s IMPLEMENTED marker written\n", styles.IconCheckmark)
			fmt.Printf("\n%s Run: agentic-agent openspec archive %s\n", styles.IconArrow, args[0])
		} else {
			fmt.Printf("Change %s implemented\n", args[0])
		}
	},
}

var openspecArchiveCmd = &cobra.Command{
	Use:   "archive <id>",
	Short: "Archive a completed change",
	Long: `Moves a completed change to the _archive directory.

Requires the IMPLEMENTED marker (run 'openspec complete' first).

Examples:
  agentic-agent openspec archive auth-feature`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := getConfig()
		m := openspec.NewManager(cfg.Paths.OpenSpecDir)

		if err := m.Archive(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			fmt.Println(styles.RenderSuccess(fmt.Sprintf("Archived change: %s", args[0])))
		} else {
			fmt.Printf("Archived change: %s\n", args[0])
		}
	},
}

func changeStatusIcon(s openspec.ChangeStatus) string {
	switch s {
	case openspec.StatusDraft:
		return styles.IconPending
	case openspec.StatusImported:
		return styles.IconPending
	case openspec.StatusImplementing:
		return styles.IconProgress
	case openspec.StatusImplemented:
		return styles.IconCheckmark
	case openspec.StatusArchived:
		return styles.IconBullet
	default:
		return styles.IconPending
	}
}

func changeStatusStyle(s openspec.ChangeStatus) string {
	switch s {
	case openspec.StatusDraft:
		return styles.MutedStyle.Render(string(s))
	case openspec.StatusImported:
		return styles.MutedStyle.Render(string(s))
	case openspec.StatusImplementing:
		return styles.SuccessStyle.Render(string(s))
	case openspec.StatusImplemented:
		return styles.SuccessStyle.Render(string(s))
	case openspec.StatusArchived:
		return styles.MutedStyle.Render(string(s))
	default:
		return string(s)
	}
}

func init() {
	openspecInitCmd.Flags().String("from", "", "Requirements file to seed the proposal (required)")

	openspecCmd.AddCommand(openspecInitCmd)
	openspecCmd.AddCommand(openspecImportCmd)
	openspecCmd.AddCommand(openspecListCmd)
	openspecCmd.AddCommand(openspecShowCmd)
	openspecCmd.AddCommand(openspecStatusCmd)
	openspecCmd.AddCommand(openspecCompleteCmd)
	openspecCmd.AddCommand(openspecArchiveCmd)
}
