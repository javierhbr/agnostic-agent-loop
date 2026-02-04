package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/javierbenavides/agentic-agent/internal/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestProject(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "validator-test-*")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})
	return tmpDir
}

func TestDirectoryContextRule_Name(t *testing.T) {
	rule := &DirectoryContextRule{}
	assert.Equal(t, "context-required", rule.Name())
}

func TestDirectoryContextRule_NoSourceFiles(t *testing.T) {
	tmpDir := setupTestProject(t)

	// Create directory with no source files
	emptyDir := filepath.Join(tmpDir, "empty")
	err := os.MkdirAll(emptyDir, 0755)
	require.NoError(t, err)

	// Create a text file (not source code)
	err = os.WriteFile(filepath.Join(emptyDir, "readme.txt"), []byte("text"), 0644)
	require.NoError(t, err)

	rule := &DirectoryContextRule{}
	ctx := &validator.ValidationContext{ProjectRoot: tmpDir}

	result, err := rule.Validate(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "PASS", result.Status)
	assert.Empty(t, result.Errors)
}

func TestDirectoryContextRule_SourceFilesWithContext(t *testing.T) {
	tmpDir := setupTestProject(t)

	// Create src directory with source files and context.md
	srcDir := filepath.Join(tmpDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(srcDir, "main.go"), []byte("package main"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(srcDir, "context.md"), []byte("# Context\nPurpose: Main entry"), 0644)
	require.NoError(t, err)

	rule := &DirectoryContextRule{}
	ctx := &validator.ValidationContext{ProjectRoot: tmpDir}

	result, err := rule.Validate(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "PASS", result.Status)
	assert.Empty(t, result.Errors)
}

func TestDirectoryContextRule_SourceFilesWithoutContext(t *testing.T) {
	tmpDir := setupTestProject(t)

	// Create src directory with source files but NO context.md
	srcDir := filepath.Join(tmpDir, "src", "auth")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(srcDir, "auth.go"), []byte("package auth"), 0644)
	require.NoError(t, err)

	rule := &DirectoryContextRule{}
	ctx := &validator.ValidationContext{ProjectRoot: tmpDir}

	result, err := rule.Validate(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "FAIL", result.Status)
	assert.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0], "src/auth")
	assert.Contains(t, result.Errors[0], "Missing context.md")
}

func TestDirectoryContextRule_MultipleDirectories(t *testing.T) {
	tmpDir := setupTestProject(t)

	// Create multiple directories
	authDir := filepath.Join(tmpDir, "src", "auth")
	coreDir := filepath.Join(tmpDir, "src", "core")

	err := os.MkdirAll(authDir, 0755)
	require.NoError(t, err)
	err = os.MkdirAll(coreDir, 0755)
	require.NoError(t, err)

	// auth has context.md
	err = os.WriteFile(filepath.Join(authDir, "auth.go"), []byte("package auth"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(authDir, "context.md"), []byte("# Auth"), 0644)
	require.NoError(t, err)

	// core does NOT have context.md
	err = os.WriteFile(filepath.Join(coreDir, "core.go"), []byte("package core"), 0644)
	require.NoError(t, err)

	rule := &DirectoryContextRule{}
	ctx := &validator.ValidationContext{ProjectRoot: tmpDir}

	result, err := rule.Validate(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "FAIL", result.Status)
	assert.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0], "src/core")
	assert.NotContains(t, result.Errors[0], "src/auth")
}

func TestDirectoryContextRule_NestedDirectories(t *testing.T) {
	tmpDir := setupTestProject(t)

	// Create nested structure
	moduleDir := filepath.Join(tmpDir, "src", "module", "submodule")
	err := os.MkdirAll(moduleDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(moduleDir, "code.ts"), []byte("export const x = 1;"), 0644)
	require.NoError(t, err)

	rule := &DirectoryContextRule{}
	ctx := &validator.ValidationContext{ProjectRoot: tmpDir}

	result, err := rule.Validate(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "FAIL", result.Status)
	assert.Contains(t, result.Errors[0], "src/module/submodule")
}

func TestDirectoryContextRule_DifferentFileTypes(t *testing.T) {
	tmpDir := setupTestProject(t)

	testCases := []struct {
		name           string
		filename       string
		shouldRequire  bool
	}{
		{"Go file", "main.go", true},
		{"TypeScript", "app.ts", true},
		{"TypeScript React", "component.tsx", true},
		{"JavaScript", "script.js", true},
		{"JavaScript React", "view.jsx", true},
		{"Python", "script.py", true},
		{"Markdown", "README.md", false},
		{"Text", "notes.txt", false},
		{"JSON", "config.json", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := filepath.Join(tmpDir, "test-"+tc.name)
			err := os.MkdirAll(dir, 0755)
			require.NoError(t, err)

			err = os.WriteFile(filepath.Join(dir, tc.filename), []byte("content"), 0644)
			require.NoError(t, err)

			rule := &DirectoryContextRule{}
			ctx := &validator.ValidationContext{ProjectRoot: tmpDir}

			result, err := rule.Validate(ctx)
			assert.NoError(t, err)

			if tc.shouldRequire {
				assert.Equal(t, "FAIL", result.Status, "Expected FAIL for %s", tc.filename)
			} else {
				// Might pass or fail depending on other directories, but shouldn't fail for this specific dir
				if result.Status == "FAIL" {
					for _, errMsg := range result.Errors {
						assert.NotContains(t, errMsg, "test-"+tc.name)
					}
				}
			}
		})
	}
}

func TestDirectoryContextRule_HiddenDirectories(t *testing.T) {
	tmpDir := setupTestProject(t)

	// Create hidden directory (should be skipped)
	hiddenDir := filepath.Join(tmpDir, ".hidden")
	err := os.MkdirAll(hiddenDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(hiddenDir, "code.go"), []byte("package hidden"), 0644)
	require.NoError(t, err)

	rule := &DirectoryContextRule{}
	ctx := &validator.ValidationContext{ProjectRoot: tmpDir}

	result, err := rule.Validate(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "PASS", result.Status)
	assert.Empty(t, result.Errors)
}

func TestDirectoryContextRule_ExcludedDirectories(t *testing.T) {
	tmpDir := setupTestProject(t)

	// Only these are actually excluded by the implementation
	excludedDirs := []string{"node_modules", "vendor", ".git"}

	for _, dirName := range excludedDirs {
		dir := filepath.Join(tmpDir, dirName)
		err := os.MkdirAll(dir, 0755)
		require.NoError(t, err)

		err = os.WriteFile(filepath.Join(dir, "code.js"), []byte("code"), 0644)
		require.NoError(t, err)
	}

	rule := &DirectoryContextRule{}
	ctx := &validator.ValidationContext{ProjectRoot: tmpDir}

	result, err := rule.Validate(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "PASS", result.Status)
	assert.Empty(t, result.Errors)
}
