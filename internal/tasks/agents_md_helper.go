package tasks

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AgentsMdHelper manages AGENTS.md files in project directories
type AgentsMdHelper struct {
	projectRoot string
}

// NewAgentsMdHelper creates a new AGENTS.md helper
func NewAgentsMdHelper(projectRoot string) *AgentsMdHelper {
	return &AgentsMdHelper{
		projectRoot: projectRoot,
	}
}

// UpdateAgentsMd creates or updates AGENTS.md in the specified directory with a learning
func (h *AgentsMdHelper) UpdateAgentsMd(dir string, learning string) error {
	agentsMdPath := filepath.Join(dir, "AGENTS.md")

	// Check if file exists
	exists := true
	if _, err := os.Stat(agentsMdPath); os.IsNotExist(err) {
		exists = false
	}

	if !exists {
		// Create new AGENTS.md
		content := fmt.Sprintf(`# Agent Guidance for %s

## Patterns

- %s
`, filepath.Base(dir), learning)
		return os.WriteFile(agentsMdPath, []byte(content), 0644)
	}

	// Append to existing file
	return h.appendPattern(agentsMdPath, learning)
}

// appendPattern appends a learning to the Patterns section of AGENTS.md
func (h *AgentsMdHelper) appendPattern(agentsMdPath string, learning string) error {
	content, err := os.ReadFile(agentsMdPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")

	// Find the Patterns section
	patternsIdx := -1
	endIdx := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "## Patterns" {
			patternsIdx = i
		} else if patternsIdx != -1 && endIdx == -1 && strings.HasPrefix(trimmed, "##") {
			endIdx = i
			break
		}
	}

	// If Patterns section doesn't exist, create it
	if patternsIdx == -1 {
		newSection := []string{
			"",
			"## Patterns",
			"",
			fmt.Sprintf("- %s", learning),
			"",
		}
		lines = append(lines, newSection...)
	} else {
		// Add to existing section
		insertIdx := patternsIdx + 1
		// Skip any empty lines after the header
		for insertIdx < len(lines) && strings.TrimSpace(lines[insertIdx]) == "" {
			insertIdx++
		}

		if endIdx != -1 {
			// Insert before next section
			lines = append(lines[:endIdx], append([]string{fmt.Sprintf("- %s", learning)}, lines[endIdx:]...)...)
		} else {
			// Append to end of Patterns section
			lines = append(lines[:insertIdx], append([]string{fmt.Sprintf("- %s", learning)}, lines[insertIdx:]...)...)
		}
	}

	// Write back
	return os.WriteFile(agentsMdPath, []byte(strings.Join(lines, "\n")), 0644)
}

// GetModifiedDirectories extracts unique directories from a list of file paths
// Filters to source code directories only (excludes hidden directories, build artifacts, etc.)
func (h *AgentsMdHelper) GetModifiedDirectories(filesChanged []string) []string {
	dirMap := make(map[string]bool)

	for _, file := range filesChanged {
		dir := filepath.Dir(file)

		// Skip if it's the root
		if dir == "." || dir == "/" {
			continue
		}

		// Make path relative to project root if needed
		if filepath.IsAbs(dir) && h.projectRoot != "" {
			relDir, err := filepath.Rel(h.projectRoot, dir)
			if err == nil {
				dir = relDir
			}
		}

		// Skip hidden directories, vendor, node_modules, build artifacts
		if h.shouldSkipDirectory(dir) {
			continue
		}

		// Use the immediate parent directory for source files
		// For deeper structures, use the top-level module directory
		topLevel := h.getTopLevelSourceDir(dir)
		if topLevel != "" {
			dirMap[topLevel] = true
		}
	}

	// Convert map to slice
	var dirs []string
	for dir := range dirMap {
		dirs = append(dirs, dir)
	}

	return dirs
}

// shouldSkipDirectory determines if a directory should be excluded from AGENTS.md tracking
func (h *AgentsMdHelper) shouldSkipDirectory(dir string) bool {
	parts := strings.Split(filepath.Clean(dir), string(filepath.Separator))

	for _, part := range parts {
		// Skip hidden directories
		if strings.HasPrefix(part, ".") {
			return true
		}

		// Skip common non-source directories
		skipDirs := map[string]bool{
			"vendor":       true,
			"node_modules": true,
			"build":        true,
			"dist":         true,
			"bin":          true,
			"target":       true,
			"out":          true,
			"coverage":     true,
		}

		if skipDirs[part] {
			return true
		}
	}

	return false
}

// getTopLevelSourceDir returns the top-level source directory for a file path
// Examples:
//   - internal/tasks/manager.go -> internal/tasks
//   - cmd/agentic-agent/main.go -> cmd/agentic-agent
//   - pkg/models/config.go -> pkg/models
func (h *AgentsMdHelper) getTopLevelSourceDir(dir string) string {
	parts := strings.Split(filepath.Clean(dir), string(filepath.Separator))

	// Look for common top-level source directories
	sourceDirs := map[string]bool{
		"internal": true,
		"cmd":      true,
		"pkg":      true,
		"src":      true,
		"lib":      true,
		"app":      true,
		"api":      true,
	}

	for i, part := range parts {
		if sourceDirs[part] {
			// Return the directory at this level or one level deeper
			if i+1 < len(parts) {
				return filepath.Join(parts[:i+2]...)
			}
			return filepath.Join(parts[:i+1]...)
		}
	}

	// If no common source dir found, return the first directory
	if len(parts) > 0 {
		return parts[0]
	}

	return dir
}

// GetExistingPatterns reads patterns from an existing AGENTS.md file
func (h *AgentsMdHelper) GetExistingPatterns(dir string) ([]string, error) {
	agentsMdPath := filepath.Join(dir, "AGENTS.md")

	if _, err := os.Stat(agentsMdPath); os.IsNotExist(err) {
		return nil, nil
	}

	f, err := os.Open(agentsMdPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var patterns []string
	scanner := bufio.NewScanner(f)
	inPatternsSection := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "## Patterns" {
			inPatternsSection = true
			continue
		}

		// Exit patterns section when we hit another section
		if inPatternsSection && strings.HasPrefix(line, "##") {
			break
		}

		// Collect pattern lines (markdown list items)
		if inPatternsSection && strings.HasPrefix(line, "-") {
			pattern := strings.TrimPrefix(line, "-")
			pattern = strings.TrimSpace(pattern)
			if pattern != "" {
				patterns = append(patterns, pattern)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return patterns, nil
}
