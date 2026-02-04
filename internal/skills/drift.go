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
	baseRules, _ := os.ReadFile(".agentic/agent-rules/base.md")

	for _, skill := range skills {
		// Calculate expected content
		tmplContent, _ := templatesFS.ReadFile(skill.TemplatePath)
		t, _ := template.New(skill.ToolName).Parse(string(tmplContent))

		data := struct{ BaseRules string }{string(baseRules)}
		var buf bytes.Buffer
		t.Execute(&buf, data)
		expected := buf.Bytes()

		// Read actual
		actual, err := os.ReadFile(skill.OutputFile)
		if err != nil {
			if os.IsNotExist(err) {
				drifted = append(drifted, fmt.Sprintf("%s (Missing)", skill.OutputFile))
			} else {
				return nil, err
			}
			continue
		}

		if !bytes.Equal(expected, actual) {
			drifted = append(drifted, fmt.Sprintf("%s (Modified)", skill.OutputFile))
		}
	}

	return drifted, nil
}
