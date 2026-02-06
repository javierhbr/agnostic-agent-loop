package specs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDirs(t *testing.T) (string, *models.Config) {
	t.Helper()
	base := t.TempDir()

	dir1 := filepath.Join(base, "specs1")
	dir2 := filepath.Join(base, "specs2")
	require.NoError(t, os.MkdirAll(dir1, 0755))
	require.NoError(t, os.MkdirAll(dir2, 0755))

	// Write spec files
	require.NoError(t, os.WriteFile(filepath.Join(dir1, "auth.md"), []byte("# Auth Spec"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir2, "api.md"), []byte("# API Spec"), 0644))
	// Duplicate in dir2 should be shadowed by dir1
	require.NoError(t, os.WriteFile(filepath.Join(dir2, "auth.md"), []byte("# Auth Spec (dir2 shadow)"), 0644))

	cfg := &models.Config{}
	cfg.Paths.SpecDirs = []string{dir1, dir2}

	return base, cfg
}

func TestResolveSpec_FromFirstDir(t *testing.T) {
	_, cfg := setupTestDirs(t)
	r := NewResolver(cfg)

	result := r.ResolveSpec("auth.md")
	assert.True(t, result.Found)
	assert.Equal(t, "# Auth Spec", result.Content)
	assert.Equal(t, "auth.md", result.Ref)
}

func TestResolveSpec_FromSecondDir(t *testing.T) {
	_, cfg := setupTestDirs(t)
	r := NewResolver(cfg)

	result := r.ResolveSpec("api.md")
	assert.True(t, result.Found)
	assert.Equal(t, "# API Spec", result.Content)
}

func TestResolveSpec_DirectPath(t *testing.T) {
	base := t.TempDir()
	specFile := filepath.Join(base, "direct.md")
	require.NoError(t, os.WriteFile(specFile, []byte("# Direct"), 0644))

	r := NewResolver(nil)
	result := r.ResolveSpec(specFile)
	assert.True(t, result.Found)
	assert.Equal(t, "# Direct", result.Content)
}

func TestResolveSpec_NotFound(t *testing.T) {
	_, cfg := setupTestDirs(t)
	r := NewResolver(cfg)

	result := r.ResolveSpec("nonexistent.md")
	assert.False(t, result.Found)
	assert.Contains(t, result.Error, "not found")
}

func TestResolveAll_MixedResults(t *testing.T) {
	_, cfg := setupTestDirs(t)
	r := NewResolver(cfg)

	results := r.ResolveAll([]string{"auth.md", "nonexistent.md", "api.md"})
	require.Len(t, results, 3)
	assert.True(t, results[0].Found)
	assert.False(t, results[1].Found)
	assert.True(t, results[2].Found)
}

func TestReadSpec_Convenience(t *testing.T) {
	_, cfg := setupTestDirs(t)
	r := NewResolver(cfg)

	content, err := r.ReadSpec("auth.md")
	assert.NoError(t, err)
	assert.Equal(t, "# Auth Spec", content)

	_, err = r.ReadSpec("missing.md")
	assert.Error(t, err)
}

func TestListSpecs(t *testing.T) {
	_, cfg := setupTestDirs(t)
	r := NewResolver(cfg)

	specs, err := r.ListSpecs()
	require.NoError(t, err)

	// Should find auth.md (from dir1) and api.md (from dir2)
	// auth.md in dir2 is shadowed
	assert.Len(t, specs, 2)

	names := make(map[string]bool)
	for _, s := range specs {
		names[s.Ref] = true
		assert.True(t, s.Found)
	}
	assert.True(t, names["auth.md"])
	assert.True(t, names["api.md"])
}

func TestListSpecs_EmptyDirs(t *testing.T) {
	base := t.TempDir()
	emptyDir := filepath.Join(base, "empty")
	require.NoError(t, os.MkdirAll(emptyDir, 0755))

	cfg := &models.Config{}
	cfg.Paths.SpecDirs = []string{emptyDir}
	r := NewResolver(cfg)

	specs, err := r.ListSpecs()
	require.NoError(t, err)
	assert.Len(t, specs, 0)
}

func TestListSpecs_NonexistentDir(t *testing.T) {
	cfg := &models.Config{}
	cfg.Paths.SpecDirs = []string{"/nonexistent/dir"}
	r := NewResolver(cfg)

	specs, err := r.ListSpecs()
	require.NoError(t, err)
	assert.Len(t, specs, 0)
}

func TestNewResolver_NilConfig(t *testing.T) {
	r := NewResolver(nil)
	assert.NotNil(t, r)
	assert.Equal(t, []string{".agentic/spec"}, r.specDirs)
}
