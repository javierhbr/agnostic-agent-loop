package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/javierbenavides/agentic-agent/internal/orchestrator"
	"github.com/spf13/cobra"
)

var autopilotCmd = &cobra.Command{
	Use:   "autopilot",
	Short: "Autopilot mode for automated task processing",
}

var autopilotStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the autopilot loop",
	Long: `Start the autopilot loop to automatically process tasks.

The autopilot will:
1. Find the next claimable task from the backlog
2. Run readiness checks (specs, inputs, scope)
3. Claim the task
4. Generate context for scope directories
5. Build a context bundle with resolved specs
6. Report the task as ready for agent execution

Flags:
  --max-iterations  Maximum number of tasks to process (default 10)
  --stop-signal     Custom stop signal string
  --dry-run         Show what would be processed without making changes`,
	Run: func(cmd *cobra.Command, args []string) {
		maxIterations, _ := cmd.Flags().GetInt("max-iterations")
		stopSignal, _ := cmd.Flags().GetString("stop-signal")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		cfg := getConfig()

		loop := orchestrator.NewAutopilotLoop(cfg, maxIterations, stopSignal, dryRun)

		// Set up context with Ctrl+C cancellation
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigCh
			fmt.Println("\nReceived interrupt signal. Stopping autopilot...")
			cancel()
		}()

		if err := loop.Run(ctx); err != nil {
			fmt.Printf("Autopilot error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	autopilotStartCmd.Flags().Int("max-iterations", 10, "Maximum number of tasks to process")
	autopilotStartCmd.Flags().String("stop-signal", "", "Custom stop signal string")
	autopilotStartCmd.Flags().Bool("dry-run", false, "Show what would be processed without making changes")

	autopilotCmd.AddCommand(autopilotStartCmd)
}
