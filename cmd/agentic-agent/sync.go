package main

import (
	"context"
	"fmt"

	"github.com/javierbenavides/agentic-agent/internal/adapters"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Bidirectional state sync between local OpenSpec and Delivery Planner",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		
		// Load local tasks
		tm := tasks.NewTaskManager(".agentic/tasks")
		localList, err := tm.LoadTasks("backlog")
		if err != nil {
			return fmt.Errorf("failed to load local tasks: %w", err)
		}
		
		// Using Mock PlannerAdapter for Phase 2
		planner := &adapters.MockPlannerAdapter{}

		err = planner.Sync(ctx, localList.Tasks)
		if err != nil {
			return fmt.Errorf("Planner sync failed: %w", err)
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			fmt.Println(styles.RenderSuccess("Successfully synchronized local tasks with Delivery Planner"))
		} else {
			fmt.Println("Sync complete.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
