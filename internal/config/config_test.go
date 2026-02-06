package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetDefaults(t *testing.T) {
	cfg := &models.Config{}
	SetDefaults(cfg)

	assert.Equal(t, ".agentic/tasks/", cfg.Paths.PRDOutputPath)
	assert.Equal(t, ".agentic/progress.txt", cfg.Paths.ProgressTextPath)
	assert.Equal(t, ".agentic/progress.yaml", cfg.Paths.ProgressYAMLPath)
	assert.Equal(t, ".agentic/archive/", cfg.Paths.ArchiveDir)
	assert.Equal(t, []string{".agentic/spec"}, cfg.Paths.SpecDirs)
	assert.Equal(t, []string{".agentic/context"}, cfg.Paths.ContextDirs)
}

func TestSetDefaults_PreservesExisting(t *testing.T) {
	cfg := &models.Config{}
	cfg.Paths.SpecDirs = []string{"custom/specs"}
	cfg.Paths.ContextDirs = []string{"custom/ctx"}
	cfg.Paths.ArchiveDir = "custom/archive"

	SetDefaults(cfg)

	assert.Equal(t, []string{"custom/specs"}, cfg.Paths.SpecDirs)
	assert.Equal(t, []string{"custom/ctx"}, cfg.Paths.ContextDirs)
	assert.Equal(t, "custom/archive", cfg.Paths.ArchiveDir)
	// But empty fields should still get defaults
	assert.Equal(t, ".agentic/tasks/", cfg.Paths.PRDOutputPath)
}

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "test-config.yaml")

	content := `
project:
  name: test-project
paths:
  specDirs:
    - specs/main
    - specs/shared
  contextDirs:
    - ctx/
`
	require.NoError(t, os.WriteFile(cfgPath, []byte(content), 0644))

	cfg, err := LoadConfig(cfgPath)
	require.NoError(t, err)

	assert.Equal(t, "test-project", cfg.Project.Name)
	assert.Equal(t, []string{"specs/main", "specs/shared"}, cfg.Paths.SpecDirs)
	assert.Equal(t, []string{"ctx/"}, cfg.Paths.ContextDirs)
	// Defaults should be applied to unset fields
	assert.Equal(t, ".agentic/tasks/", cfg.Paths.PRDOutputPath)
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/config.yaml")
	assert.Error(t, err)
}

func TestLoadConfig_EmptyPath(t *testing.T) {
	_, err := LoadConfig("")
	assert.Error(t, err)
}
