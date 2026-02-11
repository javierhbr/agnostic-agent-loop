package skills

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

func (g *Generator) CheckDrift() ([]string, error) {
	var drifted []string

	skills := g.Registry.GetAll()

	for _, skill := range skills {
		d, err := g.checkDriftForSkill(skill)
		if err != nil {
			return nil, err
		}
		if d != "" {
			drifted = append(drifted, d)
		}
	}

	return drifted, nil
}

// CheckDriftFor checks drift for a specific tool only.
func (g *Generator) CheckDriftFor(tool string) ([]string, error) {
	skill, err := g.Registry.GetSkill(tool)
	if err != nil {
		return nil, err
	}

	d, err := g.checkDriftForSkill(skill)
	if err != nil {
		return nil, err
	}
	if d != "" {
		return []string{d}, nil
	}
	return nil, nil
}

func (g *Generator) checkDriftForSkill(skill SkillDefinition) (string, error) {
	baseRules, _ := os.ReadFile(".agentic/agent-rules/base.md")
	agentRules := g.loadAgentRules(skill.ToolName)

	tmplContent, _ := templatesFS.ReadFile(skill.TemplatePath)
	t, _ := template.New(skill.ToolName).Parse(string(tmplContent))

	data := struct {
		BaseRules  string
		AgentRules string
	}{string(baseRules), agentRules}

	var buf bytes.Buffer
	t.Execute(&buf, data)
	expected := buf.Bytes()

	actual, err := os.ReadFile(skill.OutputFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Sprintf("%s (Missing)", skill.OutputFile), nil
		}
		return "", err
	}

	if !bytes.Equal(expected, actual) {
		return fmt.Sprintf("%s (Modified)", skill.OutputFile), nil
	}

	return "", nil
}
