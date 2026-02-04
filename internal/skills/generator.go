package skills

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type Generator struct {
	Registry *SkillRegistry
}

func NewGenerator() *Generator {
	return &Generator{
		Registry: NewSkillRegistry(),
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
