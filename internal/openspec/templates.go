package openspec

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed templates/*
var templatesFS embed.FS

// TemplateData holds the data used to render openspec templates.
type TemplateData struct {
	Name         string
	SourceFile   string
	Requirements string
	TaskTitle    string // used by task-detail.md.tmpl
}

// renderTemplate renders an embedded template by name.
func renderTemplate(tmplName string, data TemplateData) (string, error) {
	content, err := templatesFS.ReadFile("templates/" + tmplName)
	if err != nil {
		return "", fmt.Errorf("template %s not found: %w", tmplName, err)
	}

	tmpl, err := template.New(tmplName).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", tmplName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", tmplName, err)
	}

	return buf.String(), nil
}
