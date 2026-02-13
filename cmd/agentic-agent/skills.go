package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
	Run: func(cmd *cobra.Command, args []string) {
		if helpers.ShouldUseInteractiveMode(cmd) {
			runSkillsSubmenu()
			return
		}
		cmd.Help()
	},
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

			p := tea.NewProgram(model, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			if model.done {
				if model.success {
					fmt.Println(styles.RenderSuccess(model.message))
				} else {
					fmt.Println(styles.RenderError(model.message))
				}
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
		cfg := getConfig()
		gen := skills.NewGeneratorWithConfig(cfg)

		// If --agent is set (global flag), check only that tool
		agent := getAgent()
		var drifted []string
		var err error
		if agent.Name != "" {
			drifted, err = gen.CheckDriftFor(agent.Name)
		} else {
			drifted, err = gen.CheckDrift()
		}
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
		cfg := getConfig()

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
		cfg := getConfig()

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

// skillsInstallModel is a Bubble Tea model for interactive pack installation
type skillsInstallModel struct {
	step      string // "select-pack", "select-tool", "select-scope", "done"
	packSel   components.SimpleSelect
	toolSel   components.MultiSelect
	scopeSel  components.SimpleSelect
	installer *skills.Installer
	packName  string
	tools     []string
	results   []*skills.InstallResult
	done      bool
	success   bool
	message   string
}

func (m skillsInstallModel) Init() tea.Cmd {
	return nil
}

func (m skillsInstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			if m.step == "select-pack" {
				m.packName = m.packSel.SelectedOption().Value()
				// Build tool multi-selector
				toolOptions := []components.SelectOption{}
				for _, t := range skills.SupportedTools() {
					dir := skills.ToolSkillDir[t]
					toolOptions = append(toolOptions, components.NewSelectOption(
						t,
						fmt.Sprintf("Install to %s/", dir),
						t,
					))
				}
				m.toolSel = components.NewMultiSelect("Select target tools (Space to toggle, Enter to confirm)", toolOptions)
				m.step = "select-tool"
			} else if m.step == "select-tool" {
				m.tools = m.toolSel.SelectedValues()
				if len(m.tools) == 0 {
					m.message = "No tools selected"
					m.success = false
					m.done = true
					return m, tea.Quit
				}
				// Build scope selector
				m.scopeSel = components.NewSimpleSelect("Install scope", []components.SelectOption{
					components.NewSelectOption("Project-level", "Install to the current project directory", "local"),
					components.NewSelectOption("Global (user-level)", "Install to your home directory for all projects", "global"),
				})
				m.step = "select-scope"
			} else if m.step == "select-scope" {
				global := m.scopeSel.SelectedOption().Value() == "global"
				results, err := m.installer.InstallMulti(m.packName, m.tools, global)
				m.results = results
				if err != nil {
					m.message = fmt.Sprintf("Error: %v", err)
					m.success = false
				} else {
					var toolNames []string
					for _, r := range results {
						toolNames = append(toolNames, r.Tool)
					}
					scope := "project"
					if global {
						scope = "global"
					}
					m.message = fmt.Sprintf("Installed pack %q (%s) for %s", m.packName, scope, strings.Join(toolNames, ", "))
					m.success = true
				}
				m.done = true
				return m, tea.Quit
			}

		default:
			if m.step == "select-pack" {
				m.packSel = m.packSel.Update(msg)
			} else if m.step == "select-tool" {
				m.toolSel = m.toolSel.Update(msg)
			} else if m.step == "select-scope" {
				m.scopeSel = m.scopeSel.Update(msg)
			}
		}
	}

	return m, nil
}

func (m skillsInstallModel) View() string {
	if m.done {
		if m.success {
			var b strings.Builder
			b.WriteString(styles.TitleStyle.Render("Skill Pack Installed") + "\n\n")
			for _, r := range m.results {
				b.WriteString(styles.SuccessStyle.Render(fmt.Sprintf("%s %s", styles.IconCheckmark, r.Tool)) + "\n")
				for _, f := range r.FilesWritten {
					b.WriteString(styles.MutedStyle.Render(fmt.Sprintf("  %s %s", styles.IconBullet, f)) + "\n")
				}
			}
			return styles.ContainerStyle.Render(b.String())
		}
		return styles.RenderError(m.message) + "\n"
	}

	var b strings.Builder
	b.WriteString(styles.TitleStyle.Render("Install Skill Pack") + "\n\n")

	if m.step == "select-pack" {
		b.WriteString(m.packSel.View() + "\n")
		b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Enter select • Esc cancel") + "\n")
	} else if m.step == "select-tool" {
		b.WriteString(fmt.Sprintf("Pack: %s\n\n", styles.BoldStyle.Render(m.packName)))
		b.WriteString(m.toolSel.View() + "\n")
		b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Space toggle • Enter confirm • Esc cancel") + "\n")
	} else if m.step == "select-scope" {
		b.WriteString(fmt.Sprintf("Pack: %s  Tools: %s\n\n",
			styles.BoldStyle.Render(m.packName),
			styles.BoldStyle.Render(strings.Join(m.tools, ", ")),
		))
		b.WriteString(m.scopeSel.View() + "\n")
		b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Enter select • Esc cancel") + "\n")
	}

	return styles.ContainerStyle.Render(b.String())
}

var skillsInstallCmd = &cobra.Command{
	Use:   "install [pack-name]",
	Short: "Install a skill pack",
	Long: `Install a bundle of related skill files for one or more agent tools.

Interactive Mode:
  agentic-agent skills install

Flag Mode (single tool):
  agentic-agent skills install tdd --tool claude-code

Flag Mode (multiple tools):
  agentic-agent skills install tdd --tool claude-code,cursor,gemini

Global install:
  agentic-agent skills install tdd --tool antigravity --global

List available packs:
  agentic-agent skills install --list`,
	Run: func(cmd *cobra.Command, args []string) {
		tool, _ := cmd.Flags().GetString("tool")
		global, _ := cmd.Flags().GetBool("global")
		listFlag, _ := cmd.Flags().GetBool("list")

		installer := skills.NewInstaller()

		// List mode
		if listFlag {
			packs := installer.ListPacks()
			if helpers.ShouldUseInteractiveMode(cmd) {
				var b strings.Builder
				b.WriteString(styles.TitleStyle.Render("Available Skill Packs") + "\n\n")
				for _, p := range packs {
					b.WriteString(fmt.Sprintf("  %s %s  %s (%d files)\n",
						styles.IconBullet,
						styles.BoldStyle.Render(p.Name),
						styles.MutedStyle.Render(p.Description),
						len(p.Files),
					))
				}
				fmt.Println(styles.ContainerStyle.Render(b.String()))
			} else {
				for _, p := range packs {
					fmt.Printf("%-15s %s (%d files)\n", p.Name, p.Description, len(p.Files))
				}
			}
			return
		}

		// Interactive mode
		if helpers.ShouldUseInteractiveMode(cmd) && len(args) == 0 && tool == "" {
			packs := installer.ListPacks()
			options := []components.SelectOption{}
			for _, p := range packs {
				options = append(options, components.NewSelectOption(
					p.Name,
					fmt.Sprintf("%s (%d files)", p.Description, len(p.Files)),
					p.Name,
				))
			}

			model := &skillsInstallModel{
				step:      "select-pack",
				packSel:   components.NewSimpleSelect("Select skill pack to install", options),
				installer: installer,
			}

			p := tea.NewProgram(model, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			if model.done {
				if model.success {
					var b strings.Builder
					b.WriteString(styles.TitleStyle.Render("Skill Pack Installed") + "\n\n")
					for _, r := range model.results {
						b.WriteString(styles.SuccessStyle.Render(fmt.Sprintf("%s %s", styles.IconCheckmark, r.Tool)) + "\n")
						for _, f := range r.FilesWritten {
							b.WriteString(styles.MutedStyle.Render(fmt.Sprintf("  %s %s", styles.IconBullet, f)) + "\n")
						}
					}
					fmt.Println(styles.ContainerStyle.Render(b.String()))
				} else {
					fmt.Println(styles.RenderError(model.message))
				}
			}
			return
		}

		// Flag mode
		if len(args) == 0 {
			fmt.Println("Error: pack name required")
			fmt.Println("Usage: agentic-agent skills install <pack-name> --tool <tool>")
			fmt.Println("       agentic-agent skills install <pack-name> --tool tool1,tool2")
			fmt.Println("   or: agentic-agent skills install --list")
			fmt.Println("   or: agentic-agent skills install  (interactive mode)")
			os.Exit(1)
		}

		if tool == "" {
			fmt.Println("Error: --tool required in flag mode")
			fmt.Printf("Supported tools: %s\n", strings.Join(skills.SupportedTools(), ", "))
			os.Exit(1)
		}

		// Split comma-separated tools
		tools := strings.Split(tool, ",")
		for i := range tools {
			tools[i] = strings.TrimSpace(tools[i])
		}

		results, err := installer.InstallMulti(args[0], tools, global)
		if err != nil {
			// Print any partial successes before the error
			for _, r := range results {
				fmt.Printf("Installed pack %q for %s\n", r.PackName, r.Tool)
			}
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder
			b.WriteString(styles.TitleStyle.Render("Skill Pack Installed") + "\n\n")
			for _, r := range results {
				b.WriteString(styles.SuccessStyle.Render(fmt.Sprintf("%s Installed %q for %s", styles.IconCheckmark, r.PackName, r.Tool)) + "\n")
				for _, f := range r.FilesWritten {
					b.WriteString(styles.MutedStyle.Render(fmt.Sprintf("  %s %s", styles.IconBullet, f)) + "\n")
				}
				b.WriteString("\n")
			}
			fmt.Println(styles.ContainerStyle.Render(b.String()))
		} else {
			for _, r := range results {
				fmt.Printf("Installed pack %q for %s:\n", r.PackName, r.Tool)
				for _, f := range r.FilesWritten {
					fmt.Printf("  - %s\n", f)
				}
			}
		}
	},
}

var skillsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available skill packs",
	Run: func(cmd *cobra.Command, args []string) {
		installer := skills.NewInstaller()
		packs := installer.ListPacks()

		if helpers.ShouldUseInteractiveMode(cmd) {
			var b strings.Builder
			b.WriteString(styles.TitleStyle.Render("Available Skill Packs") + "\n\n")
			for _, p := range packs {
				b.WriteString(fmt.Sprintf("  %s %s  %s (%d files)\n",
					styles.IconBullet,
					styles.BoldStyle.Render(p.Name),
					styles.MutedStyle.Render(p.Description),
					len(p.Files),
				))
			}
			b.WriteString("\n")
			b.WriteString(styles.HelpStyle.Render("Install: agentic-agent skills install <pack> --tool <tool>") + "\n")
			fmt.Println(styles.ContainerStyle.Render(b.String()))
		} else {
			for _, p := range packs {
				fmt.Printf("%-15s %s (%d files)\n", p.Name, p.Description, len(p.Files))
			}
		}
	},
}

var skillsEnsureCmd = &cobra.Command{
	Use:   "ensure",
	Short: "Ensure all skills and rules are set up for the active agent",
	Long: `Ensure that the detected (or specified) agent has all necessary
skill files, rules, and packs installed. This is idempotent and safe to run repeatedly.

Steps performed:
1. Detect which agent(s) are active (via --agent flag, env var, or filesystem)
2. Generate rules file if missing or drifted (e.g., CLAUDE.md)
3. Install configured skill packs if not already installed
4. Report status

Usage:
  agentic-agent skills ensure                       # Auto-detect agent
  agentic-agent skills ensure --agent claude-code    # Explicit agent
  agentic-agent skills ensure --all                  # All detected agents`,
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")
		cfg := getConfig()
		agent := getAgent()

		var agents []skills.DetectedAgent
		if all {
			agents = skills.DetectAllAgents(".")
			if len(agents) == 0 {
				// Fall back to all registered tools
				registry := skills.NewSkillRegistry()
				for _, s := range registry.GetAll() {
					agents = append(agents, skills.DetectedAgent{Name: s.ToolName, Source: "registry"})
				}
			}
		} else if agent.Name != "" {
			agents = []skills.DetectedAgent{agent}
		} else {
			// Interactive mode: prompt for agent selection
			if helpers.ShouldUseInteractiveMode(cmd) {
				registry := skills.NewSkillRegistry()
				allSkills := registry.GetAll()
				options := []components.SelectOption{}
				for _, s := range allSkills {
					options = append(options, components.NewSelectOption(
						s.ToolName,
						fmt.Sprintf("Ensure skills for %s", s.ToolName),
						s.ToolName,
					))
				}

				selector := components.NewMultiSelect("Select agent tools to ensure", options)
				model := &skillsEnsureModel{selector: selector, cfg: cfg}
				p := tea.NewProgram(model, tea.WithAltScreen())
				if _, err := p.Run(); err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				if model.done && model.message != "" {
					if model.success {
						fmt.Println(styles.RenderSuccess(model.message))
					} else {
						fmt.Println(styles.RenderError(model.message))
					}
				}
				return
			}

			fmt.Println("No agent detected. Use --agent <name> or --all")
			fmt.Printf("Supported: %s\n", strings.Join(skills.SupportedTools(), ", "))
			os.Exit(1)
		}

		for _, a := range agents {
			if helpers.ShouldUseInteractiveMode(cmd) {
				var b strings.Builder
				b.WriteString(styles.TitleStyle.Render(fmt.Sprintf("Ensuring skills for %s", a.Name)) + "\n\n")

				result, err := skills.Ensure(a.Name, cfg)
				if err != nil {
					b.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("Error: %v", err)) + "\n")
				} else {
					b.WriteString(styles.SuccessStyle.Render(skills.FormatEnsureResult(result)))
				}
				fmt.Println(styles.ContainerStyle.Render(b.String()))
			} else {
				fmt.Printf("Ensuring skills for %s...\n", a.Name)
				result, err := skills.Ensure(a.Name, cfg)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				fmt.Print(skills.FormatEnsureResult(result))
			}
		}
	},
}

// skillsEnsureModel is a Bubble Tea model for interactive multi-agent selection in ensure
type skillsEnsureModel struct {
	selector components.MultiSelect
	cfg      *models.Config
	done     bool
	message  string
	success  bool
}

func (m *skillsEnsureModel) Init() tea.Cmd {
	return nil
}

func (m *skillsEnsureModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Reserve lines for title (2), help text (1), container padding (4)
		m.selector.SetMaxVisible(msg.Height - 7)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			selected := m.selector.SelectedValues()
			if len(selected) == 0 {
				m.message = "No agents selected"
				m.success = false
				m.done = true
				return m, tea.Quit
			}
			var msgs []string
			m.success = true
			for _, tool := range selected {
				result, err := skills.Ensure(tool, m.cfg)
				if err != nil {
					msgs = append(msgs, fmt.Sprintf("%s: Error: %v", tool, err))
					m.success = false
				} else {
					msgs = append(msgs, fmt.Sprintf("%s: %s", tool, skills.FormatEnsureResult(result)))
				}
			}
			m.message = strings.Join(msgs, "")
			m.done = true
			return m, tea.Quit
		default:
			m.selector = m.selector.Update(msg)
		}
	}
	return m, nil
}

func (m *skillsEnsureModel) View() string {
	if m.done {
		if m.success {
			return styles.RenderSuccess(m.message) + "\n"
		}
		return styles.RenderError(m.message) + "\n"
	}

	var b strings.Builder
	b.WriteString(styles.TitleStyle.Render("Ensure Agent Skills") + "\n\n")
	b.WriteString(m.selector.View() + "\n")
	b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Space toggle • Enter confirm • Esc cancel") + "\n")
	return styles.ContainerStyle.Render(b.String())
}

func init() {
	skillsGenerateCmd.Flags().String("tool", "", "Tool name (claude-code, cursor, gemini, windsurf, codex, copilot, opencode)")
	skillsGenerateCmd.Flags().Bool("all", false, "Generate for all tools")

	skillsInstallCmd.Flags().String("tool", "", "Target agent tool (claude-code, cursor, gemini, windsurf, antigravity, codex, copilot, opencode)")
	skillsInstallCmd.Flags().Bool("global", false, "Install to user-level directory instead of project-level")
	skillsInstallCmd.Flags().Bool("list", false, "List available skill packs")

	skillsEnsureCmd.Flags().Bool("all", false, "Ensure for all detected agents")

	skillsCmd.AddCommand(skillsGenerateCmd)
	skillsCmd.AddCommand(skillsCheckCmd)
	skillsCmd.AddCommand(skillsGenerateClaudeCmd)
	skillsCmd.AddCommand(skillsGenerateGeminiCmd)
	skillsCmd.AddCommand(skillsInstallCmd)
	skillsCmd.AddCommand(skillsListCmd)
	skillsCmd.AddCommand(skillsEnsureCmd)
}
