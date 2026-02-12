package skills

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/javierbenavides/agentic-agent/internal/config"
	"github.com/javierbenavides/agentic-agent/pkg/models"
)

type Generator struct {
	Registry *SkillRegistry
	Config   *models.Config
}

func NewGenerator() *Generator {
	return &Generator{
		Registry: NewSkillRegistry(),
	}
}

func NewGeneratorWithConfig(cfg *models.Config) *Generator {
	return &Generator{
		Registry: NewSkillRegistry(),
		Config:   cfg,
	}
}

func (g *Generator) Generate(tool string) error {
	skill, err := g.Registry.GetSkill(tool)
	if err != nil {
		return err
	}

	// Load base rules
	baseRules, err := os.ReadFile(".agentic/agent-rules/base.md")
	if err != nil {
		baseRules = []byte("No base rules found.")
	}

	// Load agent-specific rules
	agentRules := g.loadAgentRules(tool)

	// Parse template
	tmplContent, err := templatesFS.ReadFile(skill.TemplatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", skill.TemplatePath, err)
	}

	t, err := template.New(tool).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		BaseRules  string
		AgentRules string
	}{
		BaseRules:  string(baseRules),
		AgentRules: agentRules,
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return err
	}

	// Ensure output dir exists
	outputDir := filepath.Dir(skill.OutputFile)
	if outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return err
		}
	}

	if err := os.WriteFile(skill.OutputFile, buf.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

// loadAgentRules loads per-agent custom rules from .agentic/agent-rules/<tool>.md
// and merges with ExtraRules from config overrides.
func (g *Generator) loadAgentRules(tool string) string {
	var parts []string

	// 1. Load from file: .agentic/agent-rules/<tool>.md
	content, err := os.ReadFile(filepath.Join(".agentic", "agent-rules", tool+".md"))
	if err == nil && len(content) > 0 {
		parts = append(parts, strings.TrimSpace(string(content)))
	}

	// 2. Load from config overrides
	if g.Config != nil {
		agentCfg := config.GetAgentConfig(g.Config, tool)
		if len(agentCfg.ExtraRules) > 0 {
			for _, rule := range agentCfg.ExtraRules {
				parts = append(parts, "- "+rule)
			}
		}
	}

	return strings.Join(parts, "\n")
}

// toolSkillTemplates lists the skill templates that are generated for every agent tool.
var toolSkillTemplates = []struct {
	templateFile string
	outputName   string // relative to tool skill dir (e.g., "prd.md")
}{
	{templateFile: "templates/prd-skill.md", outputName: "prd.md"},
	{templateFile: "templates/ralph-converter-skill.md", outputName: "ralph-converter.md"},
}

// GenerateToolSkills generates prd and ralph-converter skill files for any agent tool.
// Files are placed in the tool's skill directory (e.g., .claude/skills/, .cursor/skills/).
func (g *Generator) GenerateToolSkills(agentName string) error {
	if g.Config == nil {
		return fmt.Errorf("config required for generating tool skills")
	}

	skillDir, ok := ToolSkillDir[agentName]
	if !ok {
		return fmt.Errorf("unsupported tool: %s", agentName)
	}

	data := struct {
		PRDOutputPath string
	}{
		PRDOutputPath: g.Config.Paths.PRDOutputPath,
	}

	for _, skill := range toolSkillTemplates {
		tmplContent, err := templatesFS.ReadFile(skill.templateFile)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", skill.templateFile, err)
		}

		t, err := template.New("skill").Parse(string(tmplContent))
		if err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}

		var buf bytes.Buffer
		if err := t.Execute(&buf, data); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		outputFile := filepath.Join(skillDir, skill.outputName)
		if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", outputFile, err)
		}

		if err := os.WriteFile(outputFile, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write skill file %s: %w", outputFile, err)
		}
	}

	// Gemini also gets slash command TOML files
	if agentName == "gemini" {
		if err := g.generateGeminiCommands(data); err != nil {
			return err
		}
	}

	return nil
}

// generateGeminiCommands generates Gemini CLI slash command TOML files.
func (g *Generator) generateGeminiCommands(data struct{ PRDOutputPath string }) error {
	commands := []struct {
		templateFile string
		outputFile   string
	}{
		{templateFile: "templates/gemini-prd-command.toml", outputFile: ".gemini/commands/prd/gen.toml"},
		{templateFile: "templates/gemini-ralph-command.toml", outputFile: ".gemini/commands/ralph/convert.toml"},
	}

	for _, cmd := range commands {
		tmplContent, err := templatesFS.ReadFile(cmd.templateFile)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", cmd.templateFile, err)
		}

		t, err := template.New("skill").Parse(string(tmplContent))
		if err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}

		var buf bytes.Buffer
		if err := t.Execute(&buf, data); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		if err := os.MkdirAll(filepath.Dir(cmd.outputFile), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", filepath.Dir(cmd.outputFile), err)
		}

		if err := os.WriteFile(cmd.outputFile, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write skill file %s: %w", cmd.outputFile, err)
		}
	}

	return nil
}

// GenerateGeminiSkills generates Gemini CLI slash command files for PRD and Ralph converter.
// Deprecated: Use GenerateToolSkills("gemini") instead, which generates both skill files and commands.
func (g *Generator) GenerateGeminiSkills() error {
	if g.Config == nil {
		return fmt.Errorf("config required for generating Gemini skills")
	}
	return g.generateGeminiCommands(struct{ PRDOutputPath string }{
		PRDOutputPath: g.Config.Paths.PRDOutputPath,
	})
}

// GenerateClaudeCodeSkills generates Claude Code skill files for PRD and Ralph converter.
// Deprecated: Use GenerateToolSkills("claude-code") instead.
func (g *Generator) GenerateClaudeCodeSkills() error {
	return g.GenerateToolSkills("claude-code")
}
