package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/javierbenavides/agentic-agent/internal/tracks"
	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/spf13/cobra"
)

var trackCmd = &cobra.Command{
	Use:   "track",
	Short: "Manage feature, bug, and refactor tracks",
	Long:  "Tracks group a spec, plan, and tasks into a single work unit.",
}

var trackInitCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a new track with brainstorming scaffolding",
	Long: `Create a new track with structured brainstorming templates.

Generates brainstorm.md (agent dialogue script), spec.md (enhanced with
purpose/constraints/success/alternatives sections), and plan.md.

The track starts in "ideation" status. Use 'track refine' to check spec
completeness, then 'track activate' to generate a plan and tasks.

Examples:
  agentic-agent track init "User Authentication"
  agentic-agent track init "Fix Login Bug" --type bug
  agentic-agent track init "Auth System" --purpose "Secure user login" --success "Users can register and login"`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := getConfig()
		m := tracks.NewManager(cfg.Paths.TrackDir)

		var name string
		if len(args) > 0 {
			name = args[0]
		}
		if name == "" {
			name, _ = cmd.Flags().GetString("name")
		}
		if name == "" {
			fmt.Fprintln(os.Stderr, "track name is required: track init \"My Feature\" or --name \"My Feature\"")
			os.Exit(1)
		}

		typeStr, _ := cmd.Flags().GetString("type")
		trackType := models.TrackTypeFeature
		switch typeStr {
		case "bug":
			trackType = models.TrackTypeBug
		case "refactor":
			trackType = models.TrackTypeRefactor
		}

		purpose, _ := cmd.Flags().GetString("purpose")
		constraints, _ := cmd.Flags().GetString("constraints")
		success, _ := cmd.Flags().GetString("success")

		opts := &tracks.InitOptions{
			Purpose:     purpose,
			Constraints: constraints,
			Success:     success,
		}

		track, err := m.Create(name, trackType, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating track: %v\n", err)
			os.Exit(1)
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			fmt.Println(styles.RenderSuccess(fmt.Sprintf("Created track: %s", track.ID)))
			fmt.Printf("  %s Brainstorm: %s/%s\n", styles.IconArrow, cfg.Paths.TrackDir, track.BrainstormPath)
			fmt.Printf("  %s Spec:       %s/%s\n", styles.IconArrow, cfg.Paths.TrackDir, track.SpecPath)
			fmt.Printf("  %s Plan:       %s/%s\n", styles.IconArrow, cfg.Paths.TrackDir, track.PlanPath)
			fmt.Printf("  %s Type:       %s\n", styles.IconArrow, track.Type)
			fmt.Printf("  %s Status:     %s\n", styles.IconArrow, track.Status)
			fmt.Printf("\n%s Use brainstorm.md as a dialogue script with your AI agent.\n", styles.IconArrow)
			fmt.Printf("%s Then run: agentic-agent track refine %s\n", styles.IconArrow, track.ID)
		} else {
			fmt.Printf("Created track: %s\n", track.ID)
		}
	},
}

var trackRefineCmd = &cobra.Command{
	Use:   "refine <id>",
	Short: "Validate spec completeness for a track",
	Long: `Check whether a track's spec.md has all required sections filled in.

Reports which sections are present, missing, or have warnings.
Exits with non-zero status if the spec is incomplete.

Examples:
  agentic-agent track refine user-authentication`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := getConfig()
		m := tracks.NewManager(cfg.Paths.TrackDir)

		track, err := m.Get(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		specPath := filepath.Join(cfg.Paths.TrackDir, track.SpecPath)
		report, err := tracks.ValidateSpec(specPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating spec: %v\n", err)
			os.Exit(1)
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder
			b.WriteString(styles.TitleStyle.Render("Spec Completeness: "+track.Name) + "\n\n")

			for _, name := range report.Present {
				b.WriteString(fmt.Sprintf("  %s %s\n", styles.IconCheckmark, name))
			}
			for _, name := range report.Missing {
				b.WriteString(fmt.Sprintf("  %s %s\n", styles.IconCross, styles.ErrorStyle.Render(name)))
			}
			for _, name := range report.Warnings {
				b.WriteString(fmt.Sprintf("  %s %s (optional)\n", styles.IconPending, styles.MutedStyle.Render(name)))
			}

			if report.Complete {
				b.WriteString(fmt.Sprintf("\n%s Spec is complete. Run: agentic-agent track activate %s\n", styles.IconCheckmark, track.ID))
			} else {
				b.WriteString(fmt.Sprintf("\n%s %d section(s) need work.\n", styles.IconCross, len(report.Missing)))
			}

			fmt.Println(styles.ContainerStyle.Render(b.String()))
		} else {
			for _, name := range report.Present {
				fmt.Printf("OK  %s\n", name)
			}
			for _, name := range report.Missing {
				fmt.Printf("MISSING  %s\n", name)
			}
			for _, name := range report.Warnings {
				fmt.Printf("WARN  %s\n", name)
			}
		}

		if !report.Complete {
			os.Exit(1)
		}
	},
}

var trackActivateCmd = &cobra.Command{
	Use:   "activate <id>",
	Short: "Generate plan and tasks from a validated spec",
	Long: `Validate the track's spec, generate a plan from it, and optionally
decompose the plan into atomic tasks.

Updates the track status from ideation to active.

Examples:
  agentic-agent track activate user-authentication
  agentic-agent track activate user-authentication --decompose`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := getConfig()
		m := tracks.NewManager(cfg.Paths.TrackDir)
		decompose, _ := cmd.Flags().GetBool("decompose")

		created, err := m.Activate(args[0], decompose, ".agentic/tasks")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error activating track: %v\n", err)
			os.Exit(1)
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			fmt.Println(styles.RenderSuccess(fmt.Sprintf("Activated track: %s", args[0])))
			fmt.Printf("  %s Plan generated from spec\n", styles.IconCheckmark)
			fmt.Printf("  %s Status: active\n", styles.IconCheckmark)
			if len(created) > 0 {
				fmt.Printf("  %s Created %d tasks in backlog\n", styles.IconCheckmark, len(created))
				for _, t := range created {
					fmt.Printf("    %s %s: %s\n", styles.IconBullet, t.ID, t.Title)
				}
			}
		} else {
			fmt.Printf("Activated track: %s\n", args[0])
			if len(created) > 0 {
				fmt.Printf("Created %d tasks\n", len(created))
			}
		}
	},
}

var trackListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tracks",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := getConfig()
		m := tracks.NewManager(cfg.Paths.TrackDir)

		trackList, err := m.List()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing tracks: %v\n", err)
			os.Exit(1)
		}

		if len(trackList) == 0 {
			fmt.Println("No tracks yet. Create one with: agentic-agent track init \"My Feature\"")
			return
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder
			b.WriteString(styles.TitleStyle.Render("Tracks") + "\n\n")
			for _, tr := range trackList {
				icon := trackStatusIcon(tr.Status)
				taskCount := fmt.Sprintf("%d tasks", len(tr.TaskIDs))
				b.WriteString(fmt.Sprintf("  %s %s  %s  %s  %s\n",
					icon,
					styles.BoldStyle.Render(tr.ID),
					styles.MutedStyle.Render(string(tr.Type)),
					statusStyle(tr.Status),
					styles.MutedStyle.Render(taskCount),
				))
			}
			fmt.Println(styles.ContainerStyle.Render(b.String()))
		} else {
			for _, tr := range trackList {
				fmt.Printf("%s  %s  %s  %d tasks\n", tr.ID, tr.Type, tr.Status, len(tr.TaskIDs))
			}
		}
	},
}

var trackShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show track details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := getConfig()
		m := tracks.NewManager(cfg.Paths.TrackDir)

		track, err := m.Get(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder
			b.WriteString(styles.TitleStyle.Render(track.Name) + "\n\n")
			b.WriteString(fmt.Sprintf("  ID:      %s\n", track.ID))
			b.WriteString(fmt.Sprintf("  Type:    %s\n", track.Type))
			b.WriteString(fmt.Sprintf("  Status:  %s\n", statusStyle(track.Status)))
			b.WriteString(fmt.Sprintf("  Created: %s\n", track.CreatedAt.Format("2006-01-02")))
			if len(track.TaskIDs) > 0 {
				b.WriteString(fmt.Sprintf("\n  %s Tasks:\n", styles.SubtitleStyle.Render("Associated")))
				for _, tid := range track.TaskIDs {
					b.WriteString(fmt.Sprintf("    %s %s\n", styles.IconBullet, tid))
				}
			}
			fmt.Println(styles.ContainerStyle.Render(b.String()))
		} else {
			fmt.Printf("ID: %s\nName: %s\nType: %s\nStatus: %s\nCreated: %s\n",
				track.ID, track.Name, track.Type, track.Status, track.CreatedAt.Format("2006-01-02"))
			if len(track.TaskIDs) > 0 {
				fmt.Printf("Tasks: %s\n", strings.Join(track.TaskIDs, ", "))
			}
		}
	},
}

var trackArchiveCmd = &cobra.Command{
	Use:   "archive <id>",
	Short: "Archive a track",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := getConfig()
		m := tracks.NewManager(cfg.Paths.TrackDir)

		if err := m.Archive(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error archiving track: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Archived track: %s\n", args[0])
	},
}

func trackStatusIcon(s models.TrackStatus) string {
	switch s {
	case models.TrackStatusIdeation:
		return styles.IconPending
	case models.TrackStatusPlanning:
		return styles.IconPending
	case models.TrackStatusActive:
		return styles.IconProgress
	case models.TrackStatusBlocked:
		return styles.IconCross
	case models.TrackStatusDone:
		return styles.IconCheckmark
	case models.TrackStatusArchived:
		return styles.IconBullet
	default:
		return styles.IconPending
	}
}

func statusStyle(s models.TrackStatus) string {
	switch s {
	case models.TrackStatusIdeation:
		return styles.MutedStyle.Render(string(s))
	case models.TrackStatusPlanning:
		return styles.MutedStyle.Render(string(s))
	case models.TrackStatusActive:
		return styles.SuccessStyle.Render(string(s))
	case models.TrackStatusBlocked:
		return styles.ErrorStyle.Render(string(s))
	case models.TrackStatusDone:
		return styles.SuccessStyle.Render(string(s))
	case models.TrackStatusArchived:
		return styles.MutedStyle.Render(string(s))
	default:
		return string(s)
	}
}

func init() {
	trackInitCmd.Flags().String("name", "", "Track name")
	trackInitCmd.Flags().String("type", "feature", "Track type (feature|bug|refactor)")
	trackInitCmd.Flags().String("purpose", "", "What problem this solves")
	trackInitCmd.Flags().String("constraints", "", "Constraints to work within")
	trackInitCmd.Flags().String("success", "", "Success criteria")

	trackActivateCmd.Flags().Bool("decompose", false, "Decompose plan into tasks")

	trackCmd.AddCommand(trackInitCmd)
	trackCmd.AddCommand(trackRefineCmd)
	trackCmd.AddCommand(trackActivateCmd)
	trackCmd.AddCommand(trackListCmd)
	trackCmd.AddCommand(trackShowCmd)
	trackCmd.AddCommand(trackArchiveCmd)
}
