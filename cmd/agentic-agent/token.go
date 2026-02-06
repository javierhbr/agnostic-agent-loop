package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/javierbenavides/agentic-agent/internal/token"
	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage token usage",
}

var tokenStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show token usage status",
	Run: func(cmd *cobra.Command, args []string) {
		tm := token.NewTokenManager(".agentic")
		usage, err := tm.LoadUsage()
		if err != nil {
			fmt.Printf("Error loading token usage: %v\n", err)
			os.Exit(1)
		}

		// Interactive mode - prettier output
		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder

			b.WriteString(styles.TitleStyle.Render("Token Usage Status") + "\n\n")

			// Total usage
			totalLine := fmt.Sprintf("Total: %s tokens", styles.BoldStyle.Render(fmt.Sprintf("%d", usage.TotalTokens)))
			b.WriteString(styles.CardStyle.Render(totalLine) + "\n\n")

			// Agent breakdown
			if len(usage.AgentUsage) > 0 {
				b.WriteString(styles.SubtitleStyle.Render("Usage by Agent") + "\n\n")

				// Sort agents for consistent display
				agents := make([]string, 0, len(usage.AgentUsage))
				for agent := range usage.AgentUsage {
					agents = append(agents, agent)
				}
				sort.Strings(agents)

				for _, agent := range agents {
					count := usage.AgentUsage[agent]
					percentage := float64(count) / float64(usage.TotalTokens) * 100
					line := fmt.Sprintf("  %s %s: %s (%s)",
						styles.IconBullet,
						agent,
						styles.BoldStyle.Render(fmt.Sprintf("%d tokens", count)),
						styles.MutedStyle.Render(fmt.Sprintf("%.1f%%", percentage)),
					)
					b.WriteString(line + "\n")
				}
			} else {
				b.WriteString(styles.MutedStyle.Render("No agent usage recorded yet.") + "\n")
			}

			fmt.Println(styles.ContainerStyle.Render(b.String()))
			return
		}

		// Flag mode - simple text output
		fmt.Printf("Total Usage: %d tokens\n", usage.TotalTokens)
		fmt.Println("Agent Usage:")
		for agent, count := range usage.AgentUsage {
			fmt.Printf("  - %s: %d\n", agent, count)
		}
	},
}

func init() {
	tokenCmd.AddCommand(tokenStatusCmd)
}
