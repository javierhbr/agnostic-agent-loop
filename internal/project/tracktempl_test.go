package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderTrackTemplate_Embedded(t *testing.T) {
	data := TrackTemplateData{
		Name:        "User Auth",
		Type:        "feature",
		Purpose:     "Allow login",
		Constraints: "Use JWT",
		Success:     "Users can log in",
	}

	content, err := RenderTrackTemplate("brainstorm.md.tmpl", data)
	require.NoError(t, err)

	assert.Contains(t, content, "# Brainstorming: User Auth")
	assert.Contains(t, content, "Allow login")
	assert.Contains(t, content, "Use JWT")
	assert.Contains(t, content, "Users can log in")
	assert.Contains(t, content, "Phase 1: Understanding")
}

func TestRenderTrackTemplate_SpecEnhanced(t *testing.T) {
	data := TrackTemplateData{
		Name:    "Auth System",
		Purpose: "Secure access",
	}

	content, err := RenderTrackTemplate("spec-enhanced.md.tmpl", data)
	require.NoError(t, err)

	assert.Contains(t, content, "# Specification: Auth System")
	assert.Contains(t, content, "Secure access")
	assert.Contains(t, content, "## Purpose")
	assert.Contains(t, content, "## Constraints")
	assert.Contains(t, content, "## Success Criteria")
	assert.Contains(t, content, "## Alternatives Considered")
	assert.Contains(t, content, "## Design")
}

func TestRenderTrackTemplate_PlanFromSpec(t *testing.T) {
	data := TrackTemplateData{Name: "My Feature"}

	content, err := RenderTrackTemplate("plan-from-spec.md.tmpl", data)
	require.NoError(t, err)

	assert.Contains(t, content, "# Plan: My Feature")
	assert.Contains(t, content, "## Phase 1: Setup")
	assert.Contains(t, content, "## Phase 2: Implementation")
	assert.Contains(t, content, "## Phase 3: Validation")
}

func TestRenderTrackTemplate_UserOverride(t *testing.T) {
	// Create user override directory
	overrideDir := filepath.Join(".agentic", "templates", "track")
	require.NoError(t, os.MkdirAll(overrideDir, 0755))
	defer os.RemoveAll(".agentic")

	customTmpl := "# Custom Brainstorm: {{.Name}}\nMy custom template."
	require.NoError(t, os.WriteFile(filepath.Join(overrideDir, "brainstorm.md.tmpl"), []byte(customTmpl), 0644))

	data := TrackTemplateData{Name: "Override Test"}
	content, err := RenderTrackTemplate("brainstorm.md.tmpl", data)
	require.NoError(t, err)

	assert.Contains(t, content, "# Custom Brainstorm: Override Test")
	assert.Contains(t, content, "My custom template.")
	assert.NotContains(t, content, "Phase 1")
}

func TestRenderTrackTemplate_NotFound(t *testing.T) {
	_, err := RenderTrackTemplate("nonexistent.tmpl", TrackTemplateData{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
