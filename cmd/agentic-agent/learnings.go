package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/javierbenavides/agentic-agent/internal/config"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/spf13/cobra"
)

var learningsCmd = &cobra.Command{
	Use:   "learnings",
	Short: "Manage codebase learnings and patterns",
	Long:  `Manage codebase learnings and patterns discovered during development.`,
}

var addLearningCmd = &cobra.Command{
	Use:   "add [pattern]",
	Short: "Add a codebase pattern to progress.txt",
	Long:  `Add a new codebase pattern or learning to the Codebase Patterns section of progress.txt.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pattern := args[0]

		// Load config to get progress file path
		cfg, err := config.LoadConfig("agnostic-agent.yaml")
		if err != nil {
			// Use defaults
			cfg = &models.Config{}
			cfg.Paths.ProgressTextPath = ".agentic/progress.txt"
		}

		// Create progress writer
		progressWriter := tasks.NewProgressWriter(cfg.Paths.ProgressTextPath, cfg.Paths.ProgressYAMLPath)

		// Add pattern
		if err := progressWriter.AddCodebasePattern(pattern); err != nil {
			fmt.Printf("Error adding pattern: %v\n", err)
			os.Exit(1)
		}

		// Interactive mode - styled output
		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder
			b.WriteString(styles.TitleStyle.Render("Pattern Added") + "\n\n")
			b.WriteString(styles.SuccessStyle.Render(fmt.Sprintf("%s Added to Codebase Patterns:", styles.IconCheckmark)) + "\n")
			b.WriteString(styles.MutedStyle.Render(fmt.Sprintf("  %s", pattern)) + "\n\n")
			b.WriteString(styles.HelpStyle.Render(fmt.Sprintf("Location: %s", cfg.Paths.ProgressTextPath)) + "\n")
			fmt.Println(styles.ContainerStyle.Render(b.String()))
			return
		}

		// Flag mode
		fmt.Printf("Added pattern to %s\n", cfg.Paths.ProgressTextPath)
	},
}

var listPatternsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all codebase patterns",
	Long:  `List all codebase patterns from the Codebase Patterns section of progress.txt.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load config to get progress file path
		cfg, err := config.LoadConfig("agnostic-agent.yaml")
		if err != nil {
			// Use defaults
			cfg = &models.Config{}
			cfg.Paths.ProgressTextPath = ".agentic/progress.txt"
		}

		// Create progress writer
		progressWriter := tasks.NewProgressWriter(cfg.Paths.ProgressTextPath, cfg.Paths.ProgressYAMLPath)

		// Get patterns
		patterns, err := progressWriter.GetCodebasePatterns()
		if err != nil {
			fmt.Printf("Error reading patterns: %v\n", err)
			os.Exit(1)
		}

		if len(patterns) == 0 {
			if helpers.ShouldUseInteractiveMode(cmd) {
				var b strings.Builder
				b.WriteString(styles.TitleStyle.Render("Codebase Patterns") + "\n\n")
				b.WriteString(styles.MutedStyle.Render("No patterns found yet.") + "\n\n")
				b.WriteString(styles.HelpStyle.Render("Add patterns with: agentic-agent learnings add \"your pattern\"") + "\n")
				fmt.Println(styles.ContainerStyle.Render(b.String()))
			} else {
				fmt.Println("No patterns found.")
			}
			return
		}

		// Interactive mode - styled output
		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder
			b.WriteString(styles.TitleStyle.Render("Codebase Patterns") + "\n\n")
			b.WriteString(styles.SuccessStyle.Render(fmt.Sprintf("%d patterns found:", len(patterns))) + "\n\n")
			for i, pattern := range patterns {
				b.WriteString(styles.MutedStyle.Render(fmt.Sprintf("%d. %s", i+1, pattern)) + "\n")
			}
			b.WriteString("\n")
			b.WriteString(styles.HelpStyle.Render(fmt.Sprintf("Source: %s", cfg.Paths.ProgressTextPath)) + "\n")
			fmt.Println(styles.ContainerStyle.Render(b.String()))
			return
		}

		// Flag mode
		fmt.Printf("Codebase Patterns (%d):\n", len(patterns))
		for i, pattern := range patterns {
			fmt.Printf("%d. %s\n", i+1, pattern)
		}
	},
}

var showProgressCmd = &cobra.Command{
	Use:   "show",
	Short: "Show recent progress entries",
	Long:  `Show recent progress entries from progress.yaml.`,
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")

		// Load config to get progress file path
		cfg, err := config.LoadConfig("agnostic-agent.yaml")
		if err != nil {
			// Use defaults
			cfg = &models.Config{}
			cfg.Paths.ProgressYAMLPath = ".agentic/progress.yaml"
		}

		// Create progress writer
		progressWriter := tasks.NewProgressWriter(cfg.Paths.ProgressTextPath, cfg.Paths.ProgressYAMLPath)

		// Get all entries
		entries, err := progressWriter.GetAllEntries()
		if err != nil {
			fmt.Printf("Error reading progress: %v\n", err)
			os.Exit(1)
		}

		if len(entries) == 0 {
			if helpers.ShouldUseInteractiveMode(cmd) {
				var b strings.Builder
				b.WriteString(styles.TitleStyle.Render("Progress Entries") + "\n\n")
				b.WriteString(styles.MutedStyle.Render("No progress entries found yet.") + "\n")
				fmt.Println(styles.ContainerStyle.Render(b.String()))
			} else {
				fmt.Println("No progress entries found.")
			}
			return
		}

		// Limit entries (most recent first)
		if limit > 0 && limit < len(entries) {
			entries = entries[len(entries)-limit:]
		}

		// Interactive mode - styled output
		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder
			b.WriteString(styles.TitleStyle.Render("Recent Progress") + "\n\n")

			for i := len(entries) - 1; i >= 0; i-- {
				entry := entries[i]
				b.WriteString(styles.SuccessStyle.Render(fmt.Sprintf("%s %s", styles.IconCheckmark, entry.StoryID)) + " ")
				b.WriteString(styles.BoldStyle.Render(entry.Title) + "\n")
				b.WriteString(styles.MutedStyle.Render(fmt.Sprintf("  %s", entry.Timestamp.Format("2006-01-02 15:04"))) + "\n")

				if len(entry.FilesChanged) > 0 {
					b.WriteString(styles.MutedStyle.Render(fmt.Sprintf("  Files: %d changed", len(entry.FilesChanged))) + "\n")
				}
				if len(entry.Learnings) > 0 {
					b.WriteString(styles.MutedStyle.Render(fmt.Sprintf("  Learnings: %d recorded", len(entry.Learnings))) + "\n")
				}
				b.WriteString("\n")
			}

			b.WriteString(styles.HelpStyle.Render(fmt.Sprintf("Total: %d entries | Source: %s", len(entries), cfg.Paths.ProgressYAMLPath)) + "\n")
			fmt.Println(styles.ContainerStyle.Render(b.String()))
			return
		}

		// Flag mode
		fmt.Printf("Recent Progress (%d entries):\n\n", len(entries))
		for i := len(entries) - 1; i >= 0; i-- {
			entry := entries[i]
			fmt.Printf("[%s] %s - %s\n", entry.Timestamp.Format("2006-01-02 15:04"), entry.StoryID, entry.Title)
			if len(entry.FilesChanged) > 0 {
				fmt.Printf("  Files changed: %d\n", len(entry.FilesChanged))
			}
			if len(entry.Learnings) > 0 {
				fmt.Printf("  Learnings: %d\n", len(entry.Learnings))
			}
			fmt.Println()
		}
	},
}

func init() {
	showProgressCmd.Flags().Int("limit", 10, "Number of recent entries to show (0 for all)")

	learningsCmd.AddCommand(addLearningCmd)
	learningsCmd.AddCommand(listPatternsCmd)
	learningsCmd.AddCommand(showProgressCmd)
}
