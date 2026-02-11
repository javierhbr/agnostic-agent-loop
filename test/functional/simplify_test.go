package functional

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/javierbenavides/agentic-agent/internal/project"
	"github.com/javierbenavides/agentic-agent/internal/simplify"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSimplifyBundleWithGoFiles tests building a simplify bundle with Go source files
func TestSimplifyBundleWithGoFiles(t *testing.T) {
	tmpDir := SetupTestProject(t)

	err := project.InitProject("SimplifyTest")
	require.NoError(t, err)

	// Create source files
	srcDir := filepath.Join(tmpDir, "internal", "auth")
	os.MkdirAll(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "handler.go"), []byte("package auth\n\nfunc Login() error { return nil }"), 0644)
	os.WriteFile(filepath.Join(srcDir, "middleware.go"), []byte("package auth\n\nfunc Protect() {}"), 0644)
	os.WriteFile(filepath.Join(srcDir, "handler_test.go"), []byte("package auth\n\nimport \"testing\"\n\nfunc TestLogin(t *testing.T) {}"), 0644)

	cfg := &models.Config{}
	bundle, err := simplify.BuildSimplifyBundle([]string{srcDir}, "", cfg)
	require.NoError(t, err)

	// Verify bundle content
	assert.NotEmpty(t, bundle.SkillInstructions, "Should contain code-simplification skill")
	assert.Contains(t, bundle.SkillInstructions, "simplicity", "Should contain simplification principles")

	// Should find 3 .go files
	assert.GreaterOrEqual(t, len(bundle.TargetFiles), 3, "Should find Go source files")

	// Verify all found files are .go
	for _, f := range bundle.TargetFiles {
		assert.Equal(t, ".go", filepath.Ext(f), "All files should be Go source files")
	}

	assert.False(t, bundle.BuiltAt.IsZero())
}

// TestSimplifyBundleMultipleDirs tests building with multiple directories
func TestSimplifyBundleMultipleDirs(t *testing.T) {
	tmpDir := SetupTestProject(t)

	err := project.InitProject("SimplifyMultiDir")
	require.NoError(t, err)

	// Create two source directories
	dir1 := filepath.Join(tmpDir, "internal", "auth")
	dir2 := filepath.Join(tmpDir, "internal", "users")
	os.MkdirAll(dir1, 0755)
	os.MkdirAll(dir2, 0755)
	os.WriteFile(filepath.Join(dir1, "auth.go"), []byte("package auth"), 0644)
	os.WriteFile(filepath.Join(dir2, "users.go"), []byte("package users"), 0644)
	os.WriteFile(filepath.Join(dir2, "repo.go"), []byte("package users"), 0644)

	cfg := &models.Config{}
	bundle, err := simplify.BuildSimplifyBundle([]string{dir1, dir2}, "", cfg)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(bundle.TargetFiles), 3)
	assert.NotEmpty(t, bundle.SkillInstructions)
}

// TestSimplifyBundleEmptyDir tests that empty directories still produce a valid bundle
func TestSimplifyBundleEmptyDir(t *testing.T) {
	tmpDir := SetupTestProject(t)

	err := project.InitProject("SimplifyEmpty")
	require.NoError(t, err)

	emptyDir := filepath.Join(tmpDir, "internal", "empty")
	os.MkdirAll(emptyDir, 0755)

	cfg := &models.Config{}
	bundle, err := simplify.BuildSimplifyBundle([]string{emptyDir}, "", cfg)
	require.NoError(t, err)

	assert.NotEmpty(t, bundle.SkillInstructions, "Skill content should still be present")
	assert.Empty(t, bundle.TargetFiles, "No source files in empty directory")
}

// TestSimplifyBundleNoDirs tests that nil/empty dirs returns an error
func TestSimplifyBundleNoDirs(t *testing.T) {
	cfg := &models.Config{}

	_, err := simplify.BuildSimplifyBundle(nil, "", cfg)
	assert.Error(t, err)

	_, err = simplify.BuildSimplifyBundle([]string{}, "", cfg)
	assert.Error(t, err)
}

// TestSimplifyBundleSkipsHiddenAndVendor tests that hidden dirs and vendor are excluded
func TestSimplifyBundleSkipsHiddenAndVendor(t *testing.T) {
	tmpDir := SetupTestProject(t)

	err := project.InitProject("SimplifySkip")
	require.NoError(t, err)

	srcDir := filepath.Join(tmpDir, "src")
	os.MkdirAll(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "main.go"), []byte("package main"), 0644)

	// Create dirs that should be skipped
	hiddenDir := filepath.Join(srcDir, ".hidden")
	vendorDir := filepath.Join(srcDir, "vendor", "dep")
	nodeDir := filepath.Join(srcDir, "node_modules", "pkg")
	os.MkdirAll(hiddenDir, 0755)
	os.MkdirAll(vendorDir, 0755)
	os.MkdirAll(nodeDir, 0755)
	os.WriteFile(filepath.Join(hiddenDir, "secret.go"), []byte("package hidden"), 0644)
	os.WriteFile(filepath.Join(vendorDir, "lib.go"), []byte("package dep"), 0644)
	os.WriteFile(filepath.Join(nodeDir, "index.js"), []byte("module.exports = {}"), 0644)

	cfg := &models.Config{}
	bundle, err := simplify.BuildSimplifyBundle([]string{srcDir}, "", cfg)
	require.NoError(t, err)

	// Should only find main.go, not files in hidden/vendor/node_modules
	assert.Len(t, bundle.TargetFiles, 1, "Should only find main.go")
	assert.Contains(t, bundle.TargetFiles[0], "main.go")
}

// TestSimplifyBundleWithTechStack tests that tech stack is included when available
func TestSimplifyBundleWithTechStack(t *testing.T) {
	tmpDir := SetupTestProject(t)

	err := project.InitProject("SimplifyTechStack")
	require.NoError(t, err)

	// Write tech stack file
	os.WriteFile(filepath.Join(tmpDir, ".agentic", "context", "tech-stack.md"),
		[]byte("# Tech Stack\n- Go 1.22\n- PostgreSQL 16"), 0644)

	srcDir := filepath.Join(tmpDir, "src")
	os.MkdirAll(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "app.go"), []byte("package main"), 0644)

	cfg := &models.Config{}
	bundle, err := simplify.BuildSimplifyBundle([]string{srcDir}, "", cfg)
	require.NoError(t, err)

	assert.Contains(t, bundle.TechStack, "Go 1.22")
	assert.Contains(t, bundle.TechStack, "PostgreSQL 16")
}
