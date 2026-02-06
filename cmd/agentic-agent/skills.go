package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/javierbenavides/agentic-agent/internal/config"
	"github.com/javierbenavides/agentic-agent/internal/skills"
	"github.com/javierbenavides/agentic-agent/internal/ui/components"
	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/spf13/cobra"
)

var skillsCmd = &cobra.Command{
	Use:   "skills",
	Short: "Manage agent skills",
}

// skillsGenerateModel is a simple Bubble Tea model for tool selection
type skillsGenerateModel struct {
	selector components.SimpleSelect
	gen      *skills.Generator
	done     bool
	success  bool
	message  string
}

func (m skillsGenerateModel) Init() tea.Cmd {
	return nil
}

func (m skillsGenerateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			selected := m.selector.SelectedOption()
			toolName := selected.Value()

			// Generate skills
			if toolName == "all" {
				registry := skills.NewSkillRegistry()
				generated := []string{}
				for _, s := range registry.GetAll() {
					if err := m.gen.Generate(s.ToolName); err == nil {
						generated = append(generated, s.OutputFile)
					}
				}
				m.message = fmt.Sprintf("Generated %d skill file(s)", len(generated))
			} else {
				if err := m.gen.Generate(toolName); err != nil {
					m.message = fmt.Sprintf("Error: %v", err)
					m.success = false
				} else {
					m.message = fmt.Sprintf("Generated skill for %s", toolName)
					m.success = true
				}
			}
			m.done = true
			return m, tea.Quit

		default:
			m.selector = m.selector.Update(msg)
		}
	}

	return m, nil
}

func (m skillsGenerateModel) View() string {
	if m.done {
		if m.success {
			return styles.RenderSuccess(m.message) + "\n"
		}
		return styles.RenderError(m.message) + "\n"
	}

	var b strings.Builder
	b.WriteString(styles.TitleStyle.Render("Generate Skill Files") + "\n\n")
	b.WriteString(m.selector.View() + "\n")
	b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Enter select • Esc cancel") + "\n")

	return styles.ContainerStyle.Render(b.String())
}

var skillsGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate skill files",
	Run: func(cmd *cobra.Command, args []string) {
		tool, _ := cmd.Flags().GetString("tool")
		all, _ := cmd.Flags().GetBool("all")

		gen := skills.NewGenerator()

		// Interactive mode - tool selection
		if helpers.ShouldUseInteractiveMode(cmd) && tool == "" && !all {
			// Get available tools
			registry := skills.NewSkillRegistry()
			allSkills := registry.GetAll()

			// Create options for selection
			options := []components.SelectOption{
				components.NewSelectOption("All Tools", "Generate skills for all available tools", "all"),
			}
			for _, s := range allSkills {
				options = append(options, components.NewSelectOption(
					s.ToolName,
					fmt.Sprintf("Generate skill file: %s", s.OutputFile),
					s.ToolName,
				))
			}

			// Create selector
			selector := components.NewSimpleSelect("Select tool to generate skills for", options)

			// Create a simple model to handle the selection
			model := &skillsGenerateModel{
				selector: selector,
				gen:      gen,
			}

			p := tea.NewProgram(model)
			if _, err := p.Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		// Flag mode or explicit flags
		if all {
			registry := skills.NewSkillRegistry()
			for _, s := range registry.GetAll() {
				if err := gen.Generate(s.ToolName); err != nil {
					fmt.Printf("Error generating %s: %v\n", s.ToolName, err)
				} else {
					fmt.Printf("Generated %s\n", s.OutputFile)
				}
			}
			return
		}

		if tool == "" {
			fmt.Println("Error: --tool or --all required in non-interactive mode")
			fmt.Println("Usage: agentic-agent skills generate --tool <name>")
			fmt.Println("   or: agentic-agent skills generate --all")
			fmt.Println("   or: agentic-agent skills generate  (interactive mode)")
			os.Exit(1)
		}

		if err := gen.Generate(tool); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Generated skill for %s\n", tool)
	},
}

var skillsCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for drift in skill files",
	Run: func(cmd *cobra.Command, args []string) {
		gen := skills.NewGenerator()
		drifted, err := gen.CheckDrift()
		if err != nil {
			fmt.Printf("Error checking drift: %v\n", err)
			os.Exit(1)
		}

		// Interactive mode - styled output
		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder

			b.WriteString(styles.TitleStyle.Render("Skill Drift Check") + "\n\n")

			if len(drifted) > 0 {
				b.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("%s Drift detected in %d file(s):", styles.IconCross, len(drifted))) + "\n\n")
				for _, d := range drifted {
					line := fmt.Sprintf("  %s %s", styles.IconBullet, d)
					b.WriteString(styles.MutedStyle.Render(line) + "\n")
				}
				b.WriteString("\n")
				b.WriteString(styles.HelpStyle.Render("Tip: Use 'agentic-agent skills generate --all' to regenerate skill files") + "\n")

				fmt.Println(styles.ContainerStyle.Render(b.String()))
				os.Exit(1)
			}

			successMsg := fmt.Sprintf("%s No drift detected - all skill files are up to date!", styles.IconCheckmark)
			b.WriteString(styles.SuccessStyle.Render(successMsg) + "\n")

			fmt.Println(styles.ContainerStyle.Render(b.String()))
			return
		}

		// Flag mode - simple text output
		if len(drifted) > 0 {
			fmt.Println("Drift detected in:")
			for _, d := range drifted {
				fmt.Printf("  - %s\n", d)
			}
			os.Exit(1)
		}
		fmt.Println("No drift detected.")
	},
}

var skillsGenerateClaudeCmd = &cobra.Command{
	Use:   "generate-claude-skills",
	Short: "Generate Claude Code skill files (PRD, Ralph converter)",
	Run: func(cmd *cobra.Command, args []string) {
		// Load config
		cfg, err := config.LoadConfig("agnostic-agent.yaml")
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			fmt.Println("Using default paths...")
			// Use defaults if config not found
			cfg = &models.Config{}
			cfg.Paths.PRDOutputPath = ".agentic/tasks/"
		}

		// Create generator with config
		gen := skills.NewGeneratorWithConfig(cfg)

		// Generate Claude Code skills
		if err := gen.GenerateClaudeCodeSkills(); err != nil {
			fmt.Printf("Error generating Claude Code skills: %v\n", err)
			os.Exit(1)
		}

		// Interactive mode - styled output
		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder
			b.WriteString(styles.TitleStyle.Render("Claude Code Skills Generated") + "\n\n")
			b.WriteString(styles.SuccessStyle.Render(fmt.Sprintf("%s Generated skill files:", styles.IconCheckmark)) + "\n")
			b.WriteString(styles.MutedStyle.Render("  • .claude/skills/prd.md") + "\n")
			b.WriteString(styles.MutedStyle.Render("  • .claude/skills/ralph-converter.md") + "\n\n")
			b.WriteString(styles.HelpStyle.Render(fmt.Sprintf("PRD output path: %s", cfg.Paths.PRDOutputPath)) + "\n")
			fmt.Println(styles.ContainerStyle.Render(b.String()))
			return
		}

		// Flag mode - simple output
		fmt.Println("Generated Claude Code skills:")
		fmt.Println("  - .claude/skills/prd.md")
		fmt.Println("  - .claude/skills/ralph-converter.md")
		fmt.Printf("PRD output path: %s\n", cfg.Paths.PRDOutputPath)
	},
}

var skillsGenerateGeminiCmd = &cobra.Command{
	Use:   "generate-gemini-skills",
	Short: "Generate Gemini CLI slash command files (PRD, Ralph converter)",
	Run: func(cmd *cobra.Command, args []string) {
		// Load config
		cfg, err := config.LoadConfig("agnostic-agent.yaml")
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			fmt.Println("Using default paths...")
			cfg = &models.Config{}
			cfg.Paths.PRDOutputPath = ".agentic/tasks/"
		}

		// Create generator with config
		gen := skills.NewGeneratorWithConfig(cfg)

		// Generate Gemini skills
		if err := gen.GenerateGeminiSkills(); err != nil {
			fmt.Printf("Error generating Gemini skills: %v\n", err)
			os.Exit(1)
		}

		// Interactive mode - styled output
		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder
			b.WriteString(styles.TitleStyle.Render("Gemini CLI Skills Generated") + "\n\n")
			b.WriteString(styles.SuccessStyle.Render(fmt.Sprintf("%s Generated skill files:", styles.IconCheckmark)) + "\n")
			b.WriteString(styles.MutedStyle.Render("  • .gemini/commands/prd/gen.toml") + "\n")
			b.WriteString(styles.MutedStyle.Render("  • .gemini/commands/ralph/convert.toml") + "\n\n")
			b.WriteString(styles.HelpStyle.Render(fmt.Sprintf("PRD output path: %s", cfg.Paths.PRDOutputPath)) + "\n")
			fmt.Println(styles.ContainerStyle.Render(b.String()))
			return
		}

		// Flag mode - simple output
		fmt.Println("Generated Gemini CLI skills:")
		fmt.Println("  - .gemini/commands/prd/gen.toml")
		fmt.Println("  - .gemini/commands/ralph/convert.toml")
		fmt.Printf("PRD output path: %s\n", cfg.Paths.PRDOutputPath)
	},
}

func init() {
	skillsGenerateCmd.Flags().String("tool", "", "Tool name (claude-code, cursor, gemini)")
	skillsGenerateCmd.Flags().Bool("all", false, "Generate for all tools")

	skillsCmd.AddCommand(skillsGenerateCmd)
	skillsCmd.AddCommand(skillsCheckCmd)
	skillsCmd.AddCommand(skillsGenerateClaudeCmd)
	skillsCmd.AddCommand(skillsGenerateGeminiCmd)
}
