package skills

import (
	"embed"
	"fmt"
)

//go:embed templates/*
var templatesFS embed.FS

type SkillDefinition struct {
	ToolName     string
	TemplatePath string
	OutputFile   string
}

type SkillRegistry struct {
	skills map[string]SkillDefinition
}

func NewSkillRegistry() *SkillRegistry {
	r := &SkillRegistry{
		skills: make(map[string]SkillDefinition),
	}

	// Register supported tools
	r.Register("claude-code", "templates/claude.tmpl", "CLAUDE.md")
	r.Register("cursor", "templates/cursor.tmpl", ".cursor/rules/agnostic-agent.mdc")
	r.Register("gemini", "templates/gemini.tmpl", ".gemini/GEMINI.md")

	return r
}

func (r *SkillRegistry) Register(tool, tmpl, output string) {
	r.skills[tool] = SkillDefinition{
		ToolName:     tool,
		TemplatePath: tmpl,
		OutputFile:   output,
	}
}

func (r *SkillRegistry) GetSkill(tool string) (SkillDefinition, error) {
	s, ok := r.skills[tool]
	if !ok {
		return SkillDefinition{}, fmt.Errorf("unknown tool: %s", tool)
	}
	return s, nil
}

func (r *SkillRegistry) GetAll() []SkillDefinition {
	var list []SkillDefinition
	for _, s := range r.skills {
		list = append(list, s)
	}
	return list
}
