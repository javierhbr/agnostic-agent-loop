package tasks

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ProgressEntry represents a single task completion entry
type ProgressEntry struct {
	Timestamp    time.Time `yaml:"timestamp"`
	StoryID      string    `yaml:"storyId"`
	Title        string    `yaml:"title"`
	FilesChanged []string  `yaml:"filesChanged"`
	Learnings    []string  `yaml:"learnings"`
	ThreadURL    string    `yaml:"threadUrl,omitempty"`
}

// ProgressWriter handles writing to both progress.txt and progress.yaml
type ProgressWriter struct {
	textPath string
	yamlPath string
}

// NewProgressWriter creates a new progress writer with the specified paths
func NewProgressWriter(textPath, yamlPath string) *ProgressWriter {
	return &ProgressWriter{
		textPath: textPath,
		yamlPath: yamlPath,
	}
}

// AppendEntry appends a progress entry to both text and YAML files
func (pw *ProgressWriter) AppendEntry(entry ProgressEntry) error {
	// Ensure directories exist
	if err := pw.ensureDirectories(); err != nil {
		return err
	}

	// Write to progress.txt
	if err := pw.appendToTextFile(entry); err != nil {
		return fmt.Errorf("failed to append to text file: %w", err)
	}

	// Write to progress.yaml
	if err := pw.appendToYAMLFile(entry); err != nil {
		return fmt.Errorf("failed to append to YAML file: %w", err)
	}

	return nil
}

// appendToTextFile appends an entry to progress.txt in markdown format
func (pw *ProgressWriter) appendToTextFile(entry ProgressEntry) error {
	// Initialize file if it doesn't exist
	if err := pw.initializeTextFile(); err != nil {
		return err
	}

	f, err := os.OpenFile(pw.textPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Format entry
	var content strings.Builder
	content.WriteString(fmt.Sprintf("\n## %s - %s\n", entry.Timestamp.Format("2006-01-02 15:04:05"), entry.StoryID))
	if entry.ThreadURL != "" {
		content.WriteString(fmt.Sprintf("Thread: %s\n", entry.ThreadURL))
	}
	content.WriteString(fmt.Sprintf("**%s**\n\n", entry.Title))

	if len(entry.FilesChanged) > 0 {
		content.WriteString("**Files Changed:**\n")
		for _, file := range entry.FilesChanged {
			content.WriteString(fmt.Sprintf("- %s\n", file))
		}
		content.WriteString("\n")
	}

	if len(entry.Learnings) > 0 {
		content.WriteString("**Learnings for future iterations:**\n")
		for _, learning := range entry.Learnings {
			content.WriteString(fmt.Sprintf("- %s\n", learning))
		}
	}

	content.WriteString("---\n")

	_, err = f.WriteString(content.String())
	return err
}

// appendToYAMLFile appends an entry to progress.yaml
func (pw *ProgressWriter) appendToYAMLFile(entry ProgressEntry) error {
	// Read existing entries
	var entries []ProgressEntry

	if _, err := os.Stat(pw.yamlPath); err == nil {
		data, err := os.ReadFile(pw.yamlPath)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(data, &entries); err != nil {
			return err
		}
	}

	// Append new entry
	entries = append(entries, entry)

	// Write back
	data, err := yaml.Marshal(entries)
	if err != nil {
		return err
	}

	return os.WriteFile(pw.yamlPath, data, 0644)
}

// GetCodebasePatterns retrieves the consolidated patterns from progress.txt
func (pw *ProgressWriter) GetCodebasePatterns() ([]string, error) {
	if _, err := os.Stat(pw.textPath); os.IsNotExist(err) {
		return nil, nil
	}

	f, err := os.Open(pw.textPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var patterns []string
	scanner := bufio.NewScanner(f)
	inPatternsSection := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "## Codebase Patterns" {
			inPatternsSection = true
			continue
		}

		// Exit patterns section when we hit another section
		if inPatternsSection && strings.HasPrefix(line, "##") {
			break
		}

		// Collect pattern lines (markdown list items starting with "- ")
		if inPatternsSection && strings.HasPrefix(line, "- ") {
			pattern := strings.TrimPrefix(line, "- ")
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

// AddCodebasePattern adds a pattern to the Codebase Patterns section
func (pw *ProgressWriter) AddCodebasePattern(pattern string) error {
	// Initialize file if needed
	if err := pw.initializeTextFile(); err != nil {
		return err
	}

	// Read existing content
	content, err := os.ReadFile(pw.textPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")

	// Find the Codebase Patterns section
	patternsIdx := -1
	endIdx := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "## Codebase Patterns" {
			patternsIdx = i
		} else if patternsIdx != -1 && endIdx == -1 && strings.HasPrefix(trimmed, "##") {
			endIdx = i
			break
		}
	}

	// If Codebase Patterns section doesn't exist, create it after header
	if patternsIdx == -1 {
		// Insert after the header (first 3 lines)
		headerEnd := 3
		if len(lines) < headerEnd {
			headerEnd = len(lines)
		}

		newSection := []string{
			"",
			"## Codebase Patterns",
			fmt.Sprintf("- %s", pattern),
			"",
		}

		lines = append(lines[:headerEnd], append(newSection, lines[headerEnd:]...)...)
	} else {
		// Add to existing section
		insertIdx := patternsIdx + 1
		if endIdx != -1 {
			insertIdx = endIdx
		} else {
			insertIdx = len(lines)
		}

		lines = append(lines[:insertIdx], append([]string{fmt.Sprintf("- %s", pattern)}, lines[insertIdx:]...)...)
	}

	// Write back
	return os.WriteFile(pw.textPath, []byte(strings.Join(lines, "\n")), 0644)
}

// initializeTextFile creates progress.txt with a header if it doesn't exist
func (pw *ProgressWriter) initializeTextFile() error {
	if _, err := os.Stat(pw.textPath); os.IsNotExist(err) {
		header := `# Progress Log

Started: ` + time.Now().Format("2006-01-02 15:04:05") + `

---
`
		return os.WriteFile(pw.textPath, []byte(header), 0644)
	}
	return nil
}

// ensureDirectories creates necessary directories for the progress files
func (pw *ProgressWriter) ensureDirectories() error {
	for _, path := range []string{pw.textPath, pw.yamlPath} {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// GetAllEntries retrieves all progress entries from progress.yaml
func (pw *ProgressWriter) GetAllEntries() ([]ProgressEntry, error) {
	if _, err := os.Stat(pw.yamlPath); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := os.ReadFile(pw.yamlPath)
	if err != nil {
		return nil, err
	}

	var entries []ProgressEntry
	if err := yaml.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// GetEntriesByFile retrieves all entries that modified a specific file
func (pw *ProgressWriter) GetEntriesByFile(filePath string) ([]ProgressEntry, error) {
	entries, err := pw.GetAllEntries()
	if err != nil {
		return nil, err
	}

	var matching []ProgressEntry
	for _, entry := range entries {
		for _, file := range entry.FilesChanged {
			if file == filePath || strings.Contains(file, filePath) {
				matching = append(matching, entry)
				break
			}
		}
	}

	return matching, nil
}
