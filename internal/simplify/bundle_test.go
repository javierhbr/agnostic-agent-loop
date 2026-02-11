package simplify

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

func TestBuildSimplifyBundle_EmptyDirs(t *testing.T) {
	cfg := &models.Config{}
	_, err := BuildSimplifyBundle(nil, "", cfg)
	if err == nil {
		t.Error("expected error for empty dirs")
	}

	_, err = BuildSimplifyBundle([]string{}, "", cfg)
	if err == nil {
		t.Error("expected error for empty slice")
	}
}

func TestBuildSimplifyBundle_ValidDir(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	// Create a source directory with Go files
	srcDir := filepath.Join(tmpDir, "src")
	os.MkdirAll(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "main.go"), []byte("package main\nfunc main() {}"), 0644)
	os.WriteFile(filepath.Join(srcDir, "helper.go"), []byte("package main\nfunc helper() {}"), 0644)

	cfg := &models.Config{}
	bundle, err := BuildSimplifyBundle([]string{srcDir}, "", cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if bundle.SkillInstructions == "" {
		t.Error("expected non-empty skill instructions")
	}

	if len(bundle.TargetFiles) < 2 {
		t.Errorf("expected at least 2 target files, got %d", len(bundle.TargetFiles))
	}

	if bundle.BuiltAt.IsZero() {
		t.Error("expected non-zero BuiltAt")
	}
}

func TestBuildSimplifyBundle_SkillContent(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	srcDir := filepath.Join(tmpDir, "src")
	os.MkdirAll(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "app.go"), []byte("package app"), 0644)

	cfg := &models.Config{}
	bundle, err := BuildSimplifyBundle([]string{srcDir}, "", cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should contain code-simplification skill content
	if bundle.SkillInstructions == "" {
		t.Error("expected code-simplification skill instructions")
	}
}

func TestScanSourceFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create various file types
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "app.ts"), []byte("const x = 1"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "readme.md"), []byte("# Readme"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "data.json"), []byte("{}"), 0644)

	// Create hidden dir (should be skipped)
	hiddenDir := filepath.Join(tmpDir, ".git")
	os.MkdirAll(hiddenDir, 0755)
	os.WriteFile(filepath.Join(hiddenDir, "config.go"), []byte("package git"), 0644)

	files, err := scanSourceFiles(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should find .go and .ts but not .md, .json, or files in .git
	if len(files) != 2 {
		t.Errorf("expected 2 source files, got %d: %v", len(files), files)
	}
}

func TestScanSourceFiles_SkipsNodeModules(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "index.ts"), []byte("export {}"), 0644)
	nmDir := filepath.Join(tmpDir, "node_modules", "dep")
	os.MkdirAll(nmDir, 0755)
	os.WriteFile(filepath.Join(nmDir, "index.js"), []byte("module.exports = {}"), 0644)

	files, err := scanSourceFiles(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("expected 1 source file (excluding node_modules), got %d: %v", len(files), files)
	}
}
