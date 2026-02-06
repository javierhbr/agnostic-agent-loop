package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/javierbenavides/agentic-agent/internal/specs"
	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	"github.com/spf13/cobra"
)

var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "Manage and resolve specification files",
}

var specResolveCmd = &cobra.Command{
	Use:   "resolve [ref]",
	Short: "Resolve and print a spec file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ref := args[0]
		resolver := specs.NewResolver(getConfig())
		result := resolver.ResolveSpec(ref)

		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder
			b.WriteString(styles.TitleStyle.Render("Spec Resolve") + "\n\n")

			if result.Found {
				b.WriteString(fmt.Sprintf("%s Resolved: %s\n", styles.IconCheckmark, styles.BoldStyle.Render(result.Ref)))
				b.WriteString(fmt.Sprintf("   Path: %s\n\n", styles.MutedStyle.Render(result.Path)))
				b.WriteString(result.Content)
			} else {
				b.WriteString(fmt.Sprintf("%s %s\n", styles.IconCross, styles.ErrorStyle.Render(result.Error)))
			}
			fmt.Println(styles.ContainerStyle.Render(b.String()))
			if !result.Found {
				os.Exit(1)
			}
			return
		}

		if !result.Found {
			fmt.Fprintf(os.Stderr, "Error: %s\n", result.Error)
			os.Exit(1)
		}
		fmt.Println(result.Content)
	},
}

var specListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all specs across configured directories",
	Run: func(cmd *cobra.Command, args []string) {
		resolver := specs.NewResolver(getConfig())
		allSpecs, err := resolver.ListSpecs()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing specs: %v\n", err)
			os.Exit(1)
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder
			b.WriteString(styles.TitleStyle.Render("Spec Files") + "\n\n")

			if len(allSpecs) == 0 {
				b.WriteString(styles.MutedStyle.Render("No spec files found in configured directories.") + "\n")
				b.WriteString(styles.HelpStyle.Render("Directories searched: "+strings.Join(getConfig().Paths.SpecDirs, ", ")) + "\n")
			} else {
				for _, s := range allSpecs {
					status := styles.SuccessStyle.Render("OK")
					if !s.Found {
						status = styles.ErrorStyle.Render("ERR")
					}
					b.WriteString(fmt.Sprintf("  %s [%s] %s\n", styles.IconBullet, status, s.Ref))
					b.WriteString(fmt.Sprintf("       %s\n", styles.MutedStyle.Render(s.Path)))
				}
				b.WriteString(fmt.Sprintf("\nTotal: %s spec file(s)\n", styles.BoldStyle.Render(fmt.Sprintf("%d", len(allSpecs)))))
			}

			fmt.Println(styles.ContainerStyle.Render(b.String()))
			return
		}

		// Flag mode
		if len(allSpecs) == 0 {
			fmt.Println("No spec files found.")
			return
		}
		for _, s := range allSpecs {
			fmt.Printf("%s  %s\n", s.Ref, s.Path)
		}
	},
}

func init() {
	specCmd.AddCommand(specResolveCmd)
	specCmd.AddCommand(specListCmd)
}
