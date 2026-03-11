package main

import (
	"context"
	"fmt"
	"os"

	"github.com/javierbenavides/agentic-agent/internal/adapters"
	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	"github.com/spf13/cobra"
)

var routeCmd = &cobra.Command{
	Use:   "route [request file or text]",
	Short: "Route a request via BMAD for track selection and sizing",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		request := args[0]
		
		// If it's a file path, we could read it, but for now just pass to adapter
		if _, err := os.Stat(request); err == nil {
			bytes, _ := os.ReadFile(request)
			request = string(bytes)
		}

		ctx := context.Background()
		// Using the MockBMADAdapter for Phase 1/2
		bmad := &adapters.MockBMADAdapter{}
		
		track, err := bmad.Route(ctx, request)
		if err != nil {
			return fmt.Errorf("BMAD routing failed: %w", err)
		}
		
		size, err := bmad.Size(ctx, request)
		if err != nil {
			return fmt.Errorf("BMAD sizing failed: %w", err)
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			fmt.Println(styles.RenderSuccess(fmt.Sprintf("Request routed to: %s (Size: %s)", track, size)))
		} else {
			fmt.Printf("Track: %s\nSize: %s\n", track, size)
		}

		return nil
	},
}

func init() {
	// Global registration happens in root.go, so we need to add routeCmd there later.
}
