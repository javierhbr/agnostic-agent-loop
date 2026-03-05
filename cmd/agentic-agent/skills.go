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
		// Handle -i as shortcut for --interactive even if passed as argument
		if len(args) > 0 && args[0] == "-i" {
			cmd.Flags().Set("interactive", "true")
			args = args[1:] // Remove -i from args
		}

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
				// UI failed - check if -i was explicitly requested and we have auto-detected agent
				interactive, _ := cmd.Flags().GetBool("interactive")
				if !interactive {
					// User didn't request -i, so this is a genuine error
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				// -i was requested but UI failed - just continue without showing the UI
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
	step         string // "select-pack", "select-scope", "select-method", "select-agents", "done"
	packSel      components.SimpleSelect
	scopeSel     components.SimpleSelect
	methodSel    components.SimpleSelect
	agentSel     components.MultiSelect
	installer    *skills.Installer
	packName     string
	agents       []string
	global       bool
	symlink      bool
	results      []*skills.InstallResult
	done         bool
	success      bool
	message      string
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
				// Build scope selector
				m.scopeSel = components.NewSimpleSelect("Install scope", []components.SelectOption{
					components.NewSelectOption("Project-level", "Install to the current project directory", "local"),
					components.NewSelectOption("Global (user-level)", "Install to your home directory for all projects", "global"),
				})
				m.step = "select-scope"
			} else if m.step == "select-scope" {
				m.global = m.scopeSel.SelectedOption().Value() == "global"
				// Build method selector
				m.methodSel = components.NewSimpleSelect("Installation method", []components.SelectOption{
					components.NewSelectOption("Copy", "Write files directly to destination", "copy"),
					components.NewSelectOption("Symlink", "Create symlinks to canonical copy (~/.agentic/skills/)", "symlink"),
				})
				m.step = "select-method"
			} else if m.step == "select-method" {
				m.symlink = m.methodSel.SelectedOption().Value() == "symlink"

				// If global + symlink, prompt for which agents; otherwise use all supported tools
				if m.global && m.symlink {
					agentOptions := []components.SelectOption{}
					for _, agent := range skills.SupportedTools() {
						agentOptions = append(agentOptions, components.NewSelectOption(
							agent,
							fmt.Sprintf("Install to %s", agent),
							agent,
						))
					}
					m.agentSel = components.NewMultiSelect(
						"Select agents to install for (Space to toggle, Enter to confirm)",
						agentOptions,
					)
					m.step = "select-agents"
				} else {
					// For non-global or non-symlink, use all supported tools
					m.agents = skills.SupportedTools()
					m.performInstallation()
					m.done = true
					return m, tea.Quit
				}
			} else if m.step == "select-agents" {
				m.agents = m.agentSel.SelectedValues()
				if len(m.agents) == 0 {
					m.message = "No agents selected"
					m.success = false
					m.done = true
					return m, tea.Quit
				}
				m.performInstallation()
				m.done = true
				return m, tea.Quit
			}

		default:
			switch m.step {
			case "select-pack":
				m.packSel = m.packSel.Update(msg)
			case "select-scope":
				m.scopeSel = m.scopeSel.Update(msg)
			case "select-method":
				m.methodSel = m.methodSel.Update(msg)
			case "select-agents":
				m.agentSel = m.agentSel.Update(msg)
			}
		}
	}

	return m, nil
}

func (m *skillsInstallModel) performInstallation() {
	results, err := m.installer.InstallMulti(m.packName, m.agents, m.global)
	m.results = results
	if err != nil {
		m.message = fmt.Sprintf("Error: %v", err)
		m.success = false
	} else {
		var agentNames []string
		for _, r := range results {
			agentNames = append(agentNames, r.Tool)
		}
		scope := "project"
		if m.global {
			scope = "global"
		}
		method := "copy"
		if m.symlink {
			method = "symlink"
		}
		m.message = fmt.Sprintf("Installed pack %q (%s, %s) for %s", m.packName, scope, method, strings.Join(agentNames, ", "))
		m.success = true
	}
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

	switch m.step {
	case "select-pack":
		b.WriteString(m.packSel.View() + "\n")
		b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Enter select • Esc cancel") + "\n")
	case "select-scope":
		b.WriteString(fmt.Sprintf("Pack: %s\n\n", styles.BoldStyle.Render(m.packName)))
		b.WriteString(m.scopeSel.View() + "\n")
		b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Enter select • Esc cancel") + "\n")
	case "select-method":
		scopeStr := "project"
		if m.global {
			scopeStr = "global"
		}
		b.WriteString(fmt.Sprintf("Pack: %s  Scope: %s\n\n",
			styles.BoldStyle.Render(m.packName),
			styles.BoldStyle.Render(scopeStr),
		))
		b.WriteString(m.methodSel.View() + "\n")
		b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Enter select • Esc cancel") + "\n")
	case "select-agents":
		scopeStr := "project"
		if m.global {
			scopeStr = "global"
		}
		methodStr := "copy"
		if m.symlink {
			methodStr = "symlink"
		}
		b.WriteString(fmt.Sprintf("Pack: %s  Scope: %s  Method: %s\n\n",
			styles.BoldStyle.Render(m.packName),
			styles.BoldStyle.Render(scopeStr),
			styles.BoldStyle.Render(methodStr),
		))
		b.WriteString(m.agentSel.View() + "\n")
		b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Space toggle • Enter confirm • Esc cancel") + "\n")
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
		// Handle -i as shortcut for --interactive even if passed as argument
		forceInteractive := false
		if len(args) > 0 && args[0] == "-i" {
			cmd.Flags().Set("interactive", "true")
			forceInteractive = true
			args = args[1:] // Remove -i from args
		}

		tool, _ := cmd.Flags().GetString("tool")
		global, _ := cmd.Flags().GetBool("global")
		listFlag, _ := cmd.Flags().GetBool("list")

		installer := skills.NewInstaller()

		// List mode
		if listFlag {
			packs := installer.ListPacks()
			if forceInteractive || helpers.ShouldUseInteractiveMode(cmd) {
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
		if (forceInteractive || helpers.ShouldUseInteractiveMode(cmd)) && len(args) == 0 && tool == "" {
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
			finalModel, err := p.Run()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			// Cast the returned model back to skillsInstallModel
			installModel, ok := finalModel.(skillsInstallModel)
			if !ok {
				fmt.Println("Error: unexpected model type")
				os.Exit(1)
			}

			if installModel.done {
				if installModel.success {
					var b strings.Builder
					b.WriteString(styles.TitleStyle.Render("Skill Pack Installed") + "\n\n")
					for _, r := range installModel.results {
						b.WriteString(styles.SuccessStyle.Render(fmt.Sprintf("%s %s", styles.IconCheckmark, r.Tool)) + "\n")
						for _, f := range r.FilesWritten {
							b.WriteString(styles.MutedStyle.Render(fmt.Sprintf("  %s %s", styles.IconBullet, f)) + "\n")
						}
					}
					fmt.Println(styles.ContainerStyle.Render(b.String()))
				} else {
					fmt.Println(styles.RenderError(installModel.message))
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
		globalFlag, _ := cmd.Flags().GetBool("global")
		symlinkFlag, _ := cmd.Flags().GetBool("symlink")
		cfg := getConfig()
		agent := getAgent()

		// First, determine installation options (destination + method)
		global := globalFlag
		symlink := symlinkFlag

		// Check if flags were explicitly set
		globalChanged := cmd.Flags().Changed("global")
		symlinkChanged := cmd.Flags().Changed("symlink")
		noInteractive, _ := cmd.Flags().GetBool("no-interactive")

		// If options not explicitly set and in interactive mode, prompt via Pretty UI
		var selectedAgents []string
		if !globalChanged && !symlinkChanged && !noInteractive && helpers.ShouldUseInteractiveMode(cmd) {
			optsModel := &skillsEnsureOptionsModel{
				step: "scope",
				scopeSel: components.NewSimpleSelect("Installation destination", []components.SelectOption{
					components.NewSelectOption("Project-level", "Install to the current project directory", "local"),
					components.NewSelectOption("Global (user-level)", "Install to your home directory for all projects", "global"),
				}),
			}
			p := tea.NewProgram(optsModel, tea.WithAltScreen())
			finalModel, err := p.Run()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			if optsModel, ok := finalModel.(*skillsEnsureOptionsModel); ok && optsModel.done {
				global = optsModel.global
				symlink = optsModel.symlink
				selectedAgents = optsModel.agents
			} else {
				return
			}
		}

		opts := skills.EnsureOptions{Global: global, Symlink: symlink}

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
		} else if len(selectedAgents) > 0 {
			// Use agents selected from the interactive modal (destination + method + agents flow)
			// Priority: interactive selection > auto-detected agent
			for _, agentName := range selectedAgents {
				agents = append(agents, skills.DetectedAgent{Name: agentName, Source: "interactive"})
			}
		} else if agent.Name != "" {
			// Agent auto-detected (or explicitly provided via --agent flag)
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
				model := &skillsEnsureModel{selector: selector, cfg: cfg, opts: opts}
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
				result, err := skills.Ensure(a.Name, cfg, opts)
				if err != nil {
					fmt.Printf("  %-20s %s\n", a.Name, styles.ErrorStyle.Render(fmt.Sprintf("Error: %v", err)))
				} else {
					fmt.Printf("  %-20s %s\n", a.Name, styles.SuccessStyle.Render(skills.FormatEnsureResultCompact(result)))
				}
			} else {
				fmt.Printf("Ensuring skills for %s...\n", a.Name)
				result, err := skills.Ensure(a.Name, cfg, opts)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				fmt.Print(skills.FormatEnsureResult(result))
			}
		}
	},
}

// skillsEnsureOptionsModel is a Bubble Tea model for interactive selection of install options (scope + method + agents + confirm)
type skillsEnsureOptionsModel struct {
	step      string                  // "scope", "method", "agents", or "confirm"
	scopeSel  components.SimpleSelect // For selecting local vs global
	methodSel components.SimpleSelect // For selecting copy vs symlink
	agentSel  components.MultiSelect  // For selecting which agents
	confirmSel components.SimpleSelect // For confirming the operation
	global    bool                    // Result: global or local
	symlink   bool                    // Result: symlink or copy
	agents    []string                // Result: selected agents
	done      bool
}

func (m *skillsEnsureOptionsModel) Init() tea.Cmd {
	return nil
}

func (m *skillsEnsureOptionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.step == "scope" {
				// Move to method selection
				m.global = m.scopeSel.SelectedOption().Value() == "global"
				methodOptions := []components.SelectOption{
					components.NewSelectOption("Copy", "Write files directly", "copy"),
					components.NewSelectOption("Symlink", "Create symlinks to canonical copy (~/.agentic/skills/)", "symlink"),
				}
				m.methodSel = components.NewSimpleSelect("Installation method", methodOptions)
				m.step = "method"
			} else if m.step == "method" {
				// Move to agent selection
				m.symlink = m.methodSel.SelectedOption().Value() == "symlink"
				agentOptions := []components.SelectOption{}
				for _, agent := range skills.SupportedTools() {
					agentOptions = append(agentOptions, components.NewSelectOption(
						agent,
						fmt.Sprintf("Ensure skills for %s", agent),
						agent,
					))
				}
				m.agentSel = components.NewMultiSelect(
					"Select agents to ensure (Space to toggle, Enter to confirm)",
					agentOptions,
				)
				m.step = "agents"
			} else if m.step == "agents" {
				// Move to confirmation
				m.agents = m.agentSel.SelectedValues()
				if len(m.agents) == 0 {
					m.agents = skills.SupportedTools() // Default to all if none selected
				}
				confirmOptions := []components.SelectOption{
					components.NewSelectOption("Yes", "Proceed with the operation", "confirm"),
					components.NewSelectOption("No", "Cancel and go back", "cancel"),
				}
				m.confirmSel = components.NewSimpleSelect("Confirm operation", confirmOptions)
				m.step = "confirm"
			} else if m.step == "confirm" {
				// Final confirmation
				if m.confirmSel.SelectedOption().Value() == "confirm" {
					m.done = true
					return m, tea.Quit
				}
				// Go back to agent selection
				m.step = "agents"
			}
		default:
			switch m.step {
			case "scope":
				m.scopeSel = m.scopeSel.Update(msg)
			case "method":
				m.methodSel = m.methodSel.Update(msg)
			case "agents":
				m.agentSel = m.agentSel.Update(msg)
			case "confirm":
				m.confirmSel = m.confirmSel.Update(msg)
			}
		}
	}
	return m, nil
}

func (m *skillsEnsureOptionsModel) View() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Ensure Agent Skills") + "\n\n")

	if m.step == "scope" {
		b.WriteString(m.scopeSel.View() + "\n")
		b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Enter confirm • Ctrl+C cancel") + "\n")
	} else if m.step == "method" {
		scopeStr := "project"
		if m.global {
			scopeStr = "global"
		}
		b.WriteString(fmt.Sprintf("Destination: %s\n\n", styles.BoldStyle.Render(scopeStr)))
		b.WriteString(m.methodSel.View() + "\n")
		b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Enter confirm • Ctrl+C cancel") + "\n")
	} else if m.step == "agents" {
		scopeStr := "project"
		if m.global {
			scopeStr = "global"
		}
		methodStr := "copy"
		if m.symlink {
			methodStr = "symlink"
		}
		b.WriteString(fmt.Sprintf("Destination: %s  Method: %s\n\n",
			styles.BoldStyle.Render(scopeStr),
			styles.BoldStyle.Render(methodStr),
		))
		b.WriteString(m.agentSel.View() + "\n")
		b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Space toggle • Enter confirm • Ctrl+C cancel") + "\n")
	} else if m.step == "confirm" {
		scopeStr := "project"
		if m.global {
			scopeStr = "global"
		}
		methodStr := "copy"
		if m.symlink {
			methodStr = "symlink"
		}
		b.WriteString(fmt.Sprintf("Destination: %s  Method: %s\n", styles.BoldStyle.Render(scopeStr), styles.BoldStyle.Render(methodStr)))
		b.WriteString(fmt.Sprintf("Agents: %s\n\n", styles.BoldStyle.Render(fmt.Sprintf("%d selected", len(m.agents)))))
		b.WriteString(m.confirmSel.View() + "\n")
		b.WriteString(styles.HelpStyle.Render("↑/↓ navigate • Enter confirm • Ctrl+C cancel") + "\n")
	}

	return styles.ContainerStyle.Render(b.String())
}

// skillsEnsureModel is a Bubble Tea model for interactive multi-agent selection in ensure
type skillsEnsureModel struct {
	selector components.MultiSelect
	cfg      *models.Config
	opts     skills.EnsureOptions
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
				result, err := skills.Ensure(tool, m.cfg, m.opts)
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
	skillsEnsureCmd.Flags().Bool("global", false, "Install to global user directory (~/.claude/skills/) instead of project-level")
	skillsEnsureCmd.Flags().Bool("symlink", false, "Create symlinks from destination to canonical copy at ~/.agentic/skills/")

	skillsCmd.AddCommand(skillsGenerateCmd)
	skillsCmd.AddCommand(skillsCheckCmd)
	skillsCmd.AddCommand(skillsGenerateClaudeCmd)
	skillsCmd.AddCommand(skillsGenerateGeminiCmd)
	skillsCmd.AddCommand(skillsInstallCmd)
	skillsCmd.AddCommand(skillsListCmd)
	skillsCmd.AddCommand(skillsEnsureCmd)
}
