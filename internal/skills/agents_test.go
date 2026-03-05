package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListAgentPacks(t *testing.T) {
	packs := ListAgentPacks()

	if len(packs) != 3 {
		t.Fatalf("expected 3 packs, got %d", len(packs))
	}

	// Check that all expected packs are present
	packNames := make(map[string]bool)
	for _, pack := range packs {
		packNames[pack.Name] = true
	}

	expectedPacks := []string{"claude-code", "openclaw", "openclaw-coordinator"}
	for _, expected := range expectedPacks {
		if !packNames[expected] {
			t.Errorf("expected pack %s not found", expected)
		}
	}
}

func TestGetAgentPack(t *testing.T) {
	tests := []struct {
		name      string
		packName  string
		wantError bool
		wantFiles int
	}{
		{"claude-code", "claude-code", false, 4},
		{"openclaw", "openclaw", false, 4},
		{"openclaw-coordinator", "openclaw-coordinator", false, 6},
		{"invalid", "invalid-pack", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pack, err := GetAgentPack(tt.packName)

			if tt.wantError && err == nil {
				t.Errorf("expected error, got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.wantError && len(pack.Files) != tt.wantFiles {
				t.Errorf("expected %d files, got %d", tt.wantFiles, len(pack.Files))
			}
		})
	}
}

func TestInstallAgents(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	// Change to temp directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}
	defer os.Chdir(origWd)

	tests := []struct {
		name      string
		packName  string
		global    bool
		wantError bool
		wantFiles int
	}{
		{"claude-code project", "claude-code", false, false, 4},
		{"openclaw-coordinator project", "openclaw-coordinator", false, false, 6},
		{"invalid pack", "invalid-pack", false, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up from previous test
			os.RemoveAll(".claude")

			filesWritten, err := InstallAgents(tt.packName, tt.global)

			if tt.wantError && err == nil {
				t.Errorf("expected error, got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.wantError && len(filesWritten) != tt.wantFiles {
				t.Errorf("expected %d files written, got %d", tt.wantFiles, len(filesWritten))
			}

			// Verify files were actually created
			if !tt.wantError {
				for _, filePath := range filesWritten {
					if _, err := os.Stat(filePath); os.IsNotExist(err) {
						t.Errorf("expected file %s not found", filePath)
					}
				}
			}
		})
	}
}

func TestIsAgentInstalled(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}
	defer os.Chdir(origWd)

	// Before installation, agents should not be installed
	if IsAgentInstalled("claude-code", false) {
		t.Errorf("expected agents not installed before install")
	}

	// Install agents
	_, err = InstallAgents("claude-code", false)
	if err != nil {
		t.Fatalf("failed to install agents: %v", err)
	}

	// After installation, agents should be installed
	if !IsAgentInstalled("claude-code", false) {
		t.Errorf("expected agents installed after install")
	}

	// Non-existent pack should return false
	if IsAgentInstalled("invalid-pack", false) {
		t.Errorf("expected non-existent pack to return false")
	}
}

func TestAgentFileContent(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}
	defer os.Chdir(origWd)

	// Install claude-code agents
	filesWritten, err := InstallAgents("claude-code", false)
	if err != nil {
		t.Fatalf("failed to install agents: %v", err)
	}

	if len(filesWritten) != 4 {
		t.Fatalf("expected 4 files, got %d", len(filesWritten))
	}

	// Verify each agent file has correct frontmatter
	expectedAgents := map[string]string{
		"orchestrator.md": "name: orchestrator",
		"worker.md":       "name: worker",
		"researcher.md":   "name: researcher",
		"reviewer.md":     "name: reviewer",
	}

	for file, expectedContent := range expectedAgents {
		filePath := filepath.Join(".claude/agents", file)
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("failed to read file %s: %v", filePath, err)
			continue
		}

		contentStr := string(content)
		if len(contentStr) == 0 {
			t.Errorf("file %s is empty", filePath)
			continue
		}

		// Check that the file contains the expected name field
		if !strings.Contains(contentStr, expectedContent) {
			t.Errorf("file %s missing expected content: %s", filePath, expectedContent)
		}

		// Check that the file has YAML frontmatter
		if !strings.Contains(contentStr, "---") {
			t.Errorf("file %s missing YAML frontmatter", filePath)
		}

		// Check that the file has the required frontmatter fields
		requiredFields := []string{"name:", "description:", "tools:", "model:", "memory:"}
		for _, field := range requiredFields {
			if !strings.Contains(contentStr, field) {
				t.Errorf("file %s missing required field: %s", filePath, field)
			}
		}
	}
}

func TestCoordinatorAgents(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change to temp directory: %v", err)
	}
	defer os.Chdir(origWd)

	// Install coordinator agents
	filesWritten, err := InstallAgents("openclaw-coordinator", false)
	if err != nil {
		t.Fatalf("failed to install coordinator agents: %v", err)
	}

	if len(filesWritten) != 6 {
		t.Fatalf("expected 6 coordinator files, got %d", len(filesWritten))
	}

	expectedAgents := []string{
		"tech-lead.md",
		"product-lead.md",
		"backend-dev.md",
		"frontend-dev.md",
		"mobile-dev.md",
		"qa-dev.md",
	}

	for _, agent := range expectedAgents {
		filePath := filepath.Join(".claude/agents", agent)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("expected agent file %s not found", filePath)
		}
	}
}

func TestGetAgentDirectoryForTool(t *testing.T) {
	tests := []struct {
		name      string
		tool      string
		global    bool
		wantPath  string
		wantError bool
	}{
		{"claude-code project", "claude-code", false, ".claude/agents", false},
		{"claude-code global", "claude-code", true, "", false}, // Will contain home dir
		{"unsupported tool", "unsupported", false, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := GetAgentDirectoryForTool(tt.tool, tt.global)

			if tt.wantError && err == nil {
				t.Errorf("expected error, got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.wantError && tt.wantPath != "" && path != tt.wantPath {
				t.Errorf("expected path %s, got %s", tt.wantPath, path)
			}
		})
	}
}

func TestSupportedAgentTools(t *testing.T) {
	tools := SupportedAgentTools()

	if len(tools) != 1 {
		t.Fatalf("expected 1 supported tool, got %d", len(tools))
	}

	if tools[0] != "claude-code" {
		t.Errorf("expected 'claude-code', got '%s'", tools[0])
	}
}

