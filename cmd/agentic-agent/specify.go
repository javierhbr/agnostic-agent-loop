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

var specifyCmd = &cobra.Command{
	Use:   "specify [input file or text]",
	Short: "Create an OpenSpec proposal and run a Spec Kit clarity pass",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		
		if _, err := os.Stat(input); err == nil {
			bytes, _ := os.ReadFile(input)
			input = string(bytes)
		}

		ctx := context.Background()
		
		// Using Mock Adapters for Phase 1/2
		openspec := &adapters.MockOpenSpecAdapter{}
		speckit := &adapters.MockSpecKitAdapter{}

		// 1. Specify using OpenSpec
		proposalPath, err := openspec.Specify(ctx, input)
		if err != nil {
			return fmt.Errorf("OpenSpec specification failed: %w", err)
		}

		// 2. Check Clarity using SpecKit
		clarityScore, err := speckit.CheckClarity(ctx, proposalPath)
		if err != nil {
			return fmt.Errorf("SpecKit clarity check failed: %w", err)
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			fmt.Println(styles.RenderSuccess(fmt.Sprintf("Created OpenSpec proposal: %s", proposalPath)))
			fmt.Println(fmt.Sprintf("SpecKit Clarity Score: %d/100", clarityScore))
		} else {
			fmt.Printf("Proposal: %s\nClarity Score: %d\n", proposalPath, clarityScore)
		}

		return nil
	},
}

func init() {
	// Global registration happens in root.go, so we need to add specifyCmd there later.
}
