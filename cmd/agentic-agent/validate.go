package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	"github.com/javierbenavides/agentic-agent/internal/validator"
	"github.com/javierbenavides/agentic-agent/internal/validator/rules"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Run validation rules",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()

		v := validator.NewValidator()
		v.Register(&rules.DirectoryContextRule{})
		v.Register(&rules.ContextUpdateRule{})
		v.Register(&rules.TaskScopeRule{})
		v.Register(&rules.TaskSizeRule{})
		v.Register(&rules.BrowserVerificationRule{})

		ctx := &validator.ValidationContext{
			ProjectRoot: cwd,
		}

		results, err := v.Validate(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating: %v\n", err)
			os.Exit(1)
		}

		format, _ := cmd.Flags().GetString("format")

		// Interactive mode - styled output
		if helpers.ShouldUseInteractiveMode(cmd) && format == "text" {
			var b strings.Builder

			b.WriteString(styles.TitleStyle.Render("Validation Results") + "\n\n")

			passCount := 0
			warnCount := 0
			failCount := 0

			for _, res := range results {
				var icon string
				var statusText string

				switch res.Status {
				case "PASS":
					icon = styles.IconCheckmark
					statusText = styles.SuccessStyle.Render(res.Status)
					passCount++
				case "WARN":
					icon = "⚠"
					statusText = styles.WarningStyle.Render(res.Status)
					warnCount++
				case "FAIL":
					icon = styles.IconCross
					statusText = styles.ErrorStyle.Render(res.Status)
					failCount++
				default:
					icon = styles.IconPending
					statusText = styles.MutedStyle.Render(res.Status)
				}

				ruleLine := fmt.Sprintf("%s %s %s",
					icon,
					statusText,
					styles.BoldStyle.Render(res.RuleName),
				)
				b.WriteString(ruleLine + "\n")

				for _, errMsg := range res.Errors {
					b.WriteString(fmt.Sprintf("    %s %s\n", styles.IconBullet, styles.MutedStyle.Render(errMsg)))
				}
			}

			b.WriteString("\n")
			summary := fmt.Sprintf("Summary: %s passed | %s warnings | %s failed",
				styles.SuccessStyle.Render(fmt.Sprintf("%d", passCount)),
				styles.WarningStyle.Render(fmt.Sprintf("%d", warnCount)),
				styles.ErrorStyle.Render(fmt.Sprintf("%d", failCount)),
			)
			b.WriteString(styles.CardStyle.Render(summary) + "\n")

			fmt.Println(styles.ContainerStyle.Render(b.String()))

			if failCount > 0 {
				os.Exit(1)
			}
			return
		}

		// Flag mode or JSON format - use existing report
		validator.PrintReport(results, format)
	},
}

func init() {
	validateCmd.Flags().String("format", "text", "Output format (text|json)")
	// Register validateCmd in root.go via this init?
	// No, standard pattern in this codebase is to have root.go add it.
}
