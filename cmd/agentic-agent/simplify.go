package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/javierbenavides/agentic-agent/internal/encoding"
	"github.com/javierbenavides/agentic-agent/internal/simplify"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var simplifyCmd = &cobra.Command{
	Use:   "simplify [dir...]",
	Short: "Generate a code simplification review bundle",
	Long: `Generate a focused context bundle for code simplification review.

Uses the code-simplification skill pack to build a context bundle
containing simplification principles and directory context for the
specified targets.

Examples:
  agentic-agent simplify internal/skills     # Simplify specific directory
  agentic-agent simplify ./cmd ./internal    # Multiple directories
  agentic-agent simplify --task TASK-1001    # Use task scope directories
  agentic-agent simplify . --format json     # Output as JSON`,
	Run: func(cmd *cobra.Command, args []string) {
		taskID, _ := cmd.Flags().GetString("task")
		output, _ := cmd.Flags().GetString("output")
		format, _ := cmd.Flags().GetString("format")

		cfg := getConfig()
		agent := cfg.ActiveAgent

		// Determine target directories
		var dirs []string
		if taskID != "" {
			dirs = dirsFromTask(taskID)
		} else if len(args) > 0 {
			dirs = args
		} else {
			fmt.Println("Error: specify directories or --task")
			fmt.Println("Usage: agentic-agent simplify [dir...] or --task <task-id>")
			os.Exit(1)
		}

		// Build the simplification bundle
		bundle, err := simplify.BuildSimplifyBundle(dirs, agent, cfg)
		if err != nil {
			fmt.Printf("Error building simplify bundle: %v\n", err)
			os.Exit(1)
		}

		// Encode output
		var data []byte
		switch format {
		case "json":
			data, err = json.MarshalIndent(bundle, "", "  ")
		case "yaml":
			data, err = yaml.Marshal(bundle)
		default:
			// toon format — wrap in the toon encoder
			encoder := encoding.NewToonEncoder()
			data, err = encoder.Encode(bundle)
		}
		if err != nil {
			fmt.Printf("Error encoding output: %v\n", err)
			os.Exit(1)
		}

		// Write output
		if output != "" {
			if err := os.WriteFile(output, data, 0644); err != nil {
				fmt.Printf("Error writing to %s: %v\n", output, err)
				os.Exit(1)
			}
			fmt.Printf("Simplification bundle written to %s (%d bytes)\n", output, len(data))
		} else {
			fmt.Print(string(data))
		}
	},
}

// dirsFromTask loads a task and returns its Scope directories.
func dirsFromTask(taskID string) []string {
	tm := tasks.NewTaskManager(".agentic/tasks")

	task, _, err := tm.FindTask(taskID)
	if err != nil || task == nil {
		fmt.Printf("Error: task %s not found: %v\n", taskID, err)
		os.Exit(1)
	}

	if len(task.Scope) == 0 {
		fmt.Printf("Error: task %s has no scope directories\n", taskID)
		os.Exit(1)
	}

	return task.Scope
}

func init() {
	simplifyCmd.Flags().String("task", "", "Task ID (uses task's scope directories)")
	simplifyCmd.Flags().String("output", "", "Output file path (default: stdout)")
	simplifyCmd.Flags().String("format", "toon", "Output format: toon, json, yaml")
}
