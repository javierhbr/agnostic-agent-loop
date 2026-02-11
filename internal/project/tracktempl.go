package project

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// TrackTemplateData holds the data used to render track templates.
type TrackTemplateData struct {
	Name        string
	Type        string
	Purpose     string
	Constraints string
	Success     string
}

// RenderTrackTemplate renders a track template by name (e.g. "brainstorm.md.tmpl").
// Resolution order: user override in .agentic/templates/track/ → embedded default.
func RenderTrackTemplate(tmplName string, data TrackTemplateData) (string, error) {
	content, err := loadTrackTemplate(tmplName)
	if err != nil {
		return "", err
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

// loadTrackTemplate loads a template, checking user override first.
func loadTrackTemplate(tmplName string) ([]byte, error) {
	// Check user override
	userPath := filepath.Join(".agentic", "templates", "track", tmplName)
	if data, err := os.ReadFile(userPath); err == nil {
		return data, nil
	}

	// Fall back to embedded
	embeddedPath := "templates/track/" + tmplName
	data, err := templatesFS.ReadFile(embeddedPath)
	if err != nil {
		return nil, fmt.Errorf("template %s not found (checked %s and embedded): %w", tmplName, userPath, err)
	}

	return data, nil
}
