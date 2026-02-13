package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/prompts"
	"github.com/javierbenavides/agentic-agent/internal/ui/components"
	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// isPromptsInteractive returns true if the prompts command should use the TUI.
// Unlike other commands, prompts is read-only browsing, so it launches the TUI
// whenever there's a TTY — even when an agent is detected via filesystem.
// Only --no-interactive explicitly disables it.
func isPromptsInteractive(cmd *cobra.Command) bool {
	if noInteractive, _ := cmd.Flags().GetBool("no-interactive"); noInteractive {
		return false
	}
	return term.IsTerminal(int(os.Stdin.Fd()))
}

var promptsCmd = &cobra.Command{
	Use:   "prompts",
	Short: "Browse example prompts and workflow recipes",
	Long: `Browse a built-in library of example prompts organized into three categories:

  Agent Prompts      Ready-to-paste prompts for your AI agent
  CLI Examples       Example CLI commands grouped by workflow
  Workflow Recipes   Step-by-step workflow sequences

Interactive Mode:
  agentic-agent prompts

Flag Mode:
  agentic-agent prompts list [--category <cat>]
  agentic-agent prompts show <slug>`,
	Run: func(cmd *cobra.Command, args []string) {
		if isPromptsInteractive(cmd) {
			runPromptsInteractive()
			return
		}
		// Non-interactive: show the full list instead of just help
		promptsListCmd.Run(cmd, []string{})
	},
}

var promptsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available prompts",
	Run: func(cmd *cobra.Command, args []string) {
		category, _ := cmd.Flags().GetString("category")
		registry := prompts.NewRegistry()

		var items []prompts.Prompt
		if category != "" {
			items = registry.ByCategory(prompts.Category(category))
			if len(items) == 0 {
				fmt.Printf("No prompts found for category: %s\n", category)
				fmt.Printf("Valid categories: agent-prompt, cli-example, workflow-recipe\n")
				os.Exit(1)
			}
		} else {
			items = registry.All()
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder
			b.WriteString(styles.TitleStyle.Render("Example Prompts Library") + "\n\n")
			for _, cat := range prompts.AllCategories() {
				if category != "" && prompts.Category(category) != cat {
					continue
				}
				catPrompts := registry.ByCategory(cat)
				if len(catPrompts) == 0 {
					continue
				}
				b.WriteString(styles.BoldStyle.Render(prompts.CategoryInfo[cat]) + "\n")
				for _, p := range catPrompts {
					b.WriteString(fmt.Sprintf("  %s %-28s %s\n",
						styles.IconBullet,
						styles.BoldStyle.Render(p.Slug),
						styles.MutedStyle.Render(p.Description),
					))
				}
				b.WriteString("\n")
			}
			b.WriteString(styles.HelpStyle.Render("Show details: agentic-agent prompts show <slug>") + "\n")
			fmt.Println(styles.ContainerStyle.Render(b.String()))
		} else {
			for _, p := range items {
				fmt.Printf("%-28s [%-15s] %s\n", p.Slug, p.Category, p.Description)
			}
		}
	},
}

var promptsShowCmd = &cobra.Command{
	Use:   "show <slug>",
	Short: "Show the full content of a prompt",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slug := args[0]
		registry := prompts.NewRegistry()
		prompt := registry.FindBySlug(slug)

		if prompt == nil {
			fmt.Printf("Prompt not found: %s\n", slug)
			fmt.Println("Run 'agentic-agent prompts list' to see available prompts.")
			os.Exit(1)
		}

		fmt.Println(renderPromptDetail(prompt))
	},
}

// renderPromptDetail renders a prompt with clear visual structure.
// Uses horizontal rules instead of a bordered box so copied text stays clean.
func renderPromptDetail(p *prompts.Prompt) string {
	var b strings.Builder
	separator := styles.MutedStyle.Render("────────────────────────────────────────────────────────")

	// Header
	b.WriteString(styles.TitleStyle.Render(p.Title) + "\n")
	b.WriteString(styles.MutedStyle.Render(fmt.Sprintf("Category: %s  |  Slug: %s", prompts.CategoryInfo[p.Category], p.Slug)) + "\n")
	b.WriteString(styles.MutedStyle.Render(p.Description) + "\n\n")

	// Prompt content between horizontal rules (no side borders to copy)
	b.WriteString(styles.BoldStyle.Render("Prompt:") + "\n")
	b.WriteString(separator + "\n\n")
	b.WriteString(p.Content + "\n\n")
	b.WriteString(separator + "\n\n")

	// Usage hint
	b.WriteString(styles.HelpStyle.Render("Copy the text between the lines and paste it into your AI agent.") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

// promptsBrowseModel is the Bubble Tea model for interactive prompt browsing.
type promptsBrowseModel struct {
	step        string // "select-category", "select-prompt", "done"
	categorySel components.SimpleSelect
	promptSel   components.SimpleSelect
	registry    *prompts.Registry
	category    prompts.Category
	showAll     bool
	selected    *prompts.Prompt
	done        bool
	height      int // terminal height for scroll viewport
}

func (m promptsBrowseModel) Init() tea.Cmd { return nil }

func (m promptsBrowseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Reserve lines for title, category label, help text, and padding
		m.height = msg.Height - 6
		m.categorySel.SetMaxVisible(m.height)
		m.promptSel.SetMaxVisible(m.height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q", "esc":
			if m.step == "select-prompt" {
				m.step = "select-category"
				return m, nil
			}
			return m, tea.Quit

		case "enter":
			if m.step == "select-category" {
				catValue := m.categorySel.SelectedOption().Value()
				m.showAll = catValue == "all"
				if !m.showAll {
					m.category = prompts.Category(catValue)
				}

				var catPrompts []prompts.Prompt
				if m.showAll {
					catPrompts = m.registry.All()
				} else {
					catPrompts = m.registry.ByCategory(m.category)
				}

				options := make([]components.SelectOption, 0, len(catPrompts))
				for _, p := range catPrompts {
					options = append(options, components.NewSelectOption(
						p.Title,
						p.Description,
						p.Slug,
					))
				}
				m.promptSel = components.NewSimpleSelect("Select a prompt", options)
				if m.height > 0 {
					m.promptSel.SetMaxVisible(m.height)
				}
				m.step = "select-prompt"

			} else if m.step == "select-prompt" {
				slug := m.promptSel.SelectedOption().Value()
				m.selected = m.registry.FindBySlug(slug)
				m.done = true
				return m, tea.Quit
			}

		default:
			if m.step == "select-category" {
				m.categorySel = m.categorySel.Update(msg)
			} else if m.step == "select-prompt" {
				m.promptSel = m.promptSel.Update(msg)
			}
		}
	}
	return m, nil
}

func (m promptsBrowseModel) View() string {
	if m.done {
		return ""
	}

	var b strings.Builder
	b.WriteString(styles.TitleStyle.Render("Example Prompts Library") + "\n\n")

	if m.step == "select-category" {
		b.WriteString(m.categorySel.View() + "\n")
		b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Enter select • Esc quit") + "\n")
	} else if m.step == "select-prompt" {
		label := "All"
		if !m.showAll {
			label = prompts.CategoryInfo[m.category]
		}
		b.WriteString(styles.MutedStyle.Render(fmt.Sprintf("Category: %s", label)) + "\n\n")
		b.WriteString(m.promptSel.View() + "\n")
		b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Enter select • Esc back") + "\n")
	}

	return styles.ContainerStyle.Render(b.String())
}

func runPromptsInteractive() {
	registry := prompts.NewRegistry()

	categoryOptions := []components.SelectOption{
		components.NewSelectOption("Agent Prompts", "Ready-to-paste prompts for your AI agent", string(prompts.CategoryAgentPrompt)),
		components.NewSelectOption("CLI Examples", "Example CLI commands grouped by workflow", string(prompts.CategoryCLIExample)),
		components.NewSelectOption("Workflow Recipes", "Step-by-step workflow sequences", string(prompts.CategoryWorkflowRecipe)),
		components.NewSelectOption("All", "Browse all prompts", "all"),
	}

	model := &promptsBrowseModel{
		step:        "select-category",
		categorySel: components.NewSimpleSelect("Select a category", categoryOptions),
		registry:    registry,
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(promptsBrowseModel); ok && m.done && m.selected != nil {
		fmt.Println(renderPromptDetail(m.selected))
	}
}

func init() {
	promptsListCmd.Flags().String("category", "", "Filter by category (agent-prompt, cli-example, workflow-recipe)")

	promptsCmd.AddCommand(promptsListCmd)
	promptsCmd.AddCommand(promptsShowCmd)
}
