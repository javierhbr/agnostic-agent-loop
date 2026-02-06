package skills

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

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
		// Fallback if not initialized? Or error.
		baseRules = []byte("No base rules found.")
	}

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
		BaseRules string
	}{
		BaseRules: string(baseRules),
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

// GenerateGeminiSkills generates Gemini CLI slash command files for PRD and Ralph converter
func (g *Generator) GenerateGeminiSkills() error {
	if g.Config == nil {
		return fmt.Errorf("config required for generating Gemini skills")
	}

	skills := []struct {
		templateFile string
		outputFile   string
	}{
		{
			templateFile: "templates/gemini-prd-command.toml",
			outputFile:   ".gemini/commands/prd/gen.toml",
		},
		{
			templateFile: "templates/gemini-ralph-command.toml",
			outputFile:   ".gemini/commands/ralph/convert.toml",
		},
	}

	data := struct {
		PRDOutputPath string
	}{
		PRDOutputPath: g.Config.Paths.PRDOutputPath,
	}

	for _, skill := range skills {
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

		outputDir := filepath.Dir(skill.outputFile)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
		}

		if err := os.WriteFile(skill.outputFile, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write skill file %s: %w", skill.outputFile, err)
		}
	}

	return nil
}

// GenerateClaudeCodeSkills generates Claude Code skill files for PRD and Ralph converter
func (g *Generator) GenerateClaudeCodeSkills() error {
	if g.Config == nil {
		return fmt.Errorf("config required for generating Claude Code skills")
	}

	// Define skill templates to generate
	skills := []struct {
		templateFile string
		outputFile   string
	}{
		{
			templateFile: "templates/prd-skill.md",
			outputFile:   ".claude/skills/prd.md",
		},
		{
			templateFile: "templates/ralph-converter-skill.md",
			outputFile:   ".claude/skills/ralph-converter.md",
		},
	}

	// Template data with config paths
	data := struct {
		PRDOutputPath string
	}{
		PRDOutputPath: g.Config.Paths.PRDOutputPath,
	}

	for _, skill := range skills {
		// Read template
		tmplContent, err := templatesFS.ReadFile(skill.templateFile)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", skill.templateFile, err)
		}

		// Parse and execute template
		t, err := template.New("skill").Parse(string(tmplContent))
		if err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}

		var buf bytes.Buffer
		if err := t.Execute(&buf, data); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		// Ensure output directory exists
		outputDir := filepath.Dir(skill.outputFile)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
		}

		// Write skill file
		if err := os.WriteFile(skill.outputFile, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write skill file %s: %w", skill.outputFile, err)
		}
	}

	return nil
}
