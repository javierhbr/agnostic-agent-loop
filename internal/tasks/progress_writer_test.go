package tasks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProgressWriter_AppendEntry(t *testing.T) {
	// Setup temp directory
	tmpDir := t.TempDir()
	textPath := filepath.Join(tmpDir, "progress.txt")
	yamlPath := filepath.Join(tmpDir, "progress.yaml")

	pw := NewProgressWriter(textPath, yamlPath)

	entry := ProgressEntry{
		Timestamp:    time.Date(2026, 2, 5, 10, 30, 0, 0, time.UTC),
		StoryID:      "US-001",
		Title:        "Add user authentication",
		FilesChanged: []string{"src/auth/login.go", "src/auth/middleware.go"},
		Learnings:    []string{"Always validate JWT tokens", "Use bcrypt for passwords"},
		ThreadURL:    "https://claude.ai/thread/abc123",
	}

	err := pw.AppendEntry(entry)
	require.NoError(t, err)

	// Verify text file exists
	assert.FileExists(t, textPath)

	// Verify YAML file exists
	assert.FileExists(t, yamlPath)

	// Read and verify text content
	textContent, err := os.ReadFile(textPath)
	require.NoError(t, err)
	textStr := string(textContent)

	assert.Contains(t, textStr, "US-001")
	assert.Contains(t, textStr, "Add user authentication")
	assert.Contains(t, textStr, "src/auth/login.go")
	assert.Contains(t, textStr, "Always validate JWT tokens")
	assert.Contains(t, textStr, "https://claude.ai/thread/abc123")

	// Verify YAML can be read back
	entries, err := pw.GetAllEntries()
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, "US-001", entries[0].StoryID)
	assert.Equal(t, "Add user authentication", entries[0].Title)
	assert.Len(t, entries[0].FilesChanged, 2)
	assert.Len(t, entries[0].Learnings, 2)
}

func TestProgressWriter_AppendMultipleEntries(t *testing.T) {
	tmpDir := t.TempDir()
	textPath := filepath.Join(tmpDir, "progress.txt")
	yamlPath := filepath.Join(tmpDir, "progress.yaml")

	pw := NewProgressWriter(textPath, yamlPath)

	// Add first entry
	entry1 := ProgressEntry{
		Timestamp:    time.Date(2026, 2, 5, 10, 0, 0, 0, time.UTC),
		StoryID:      "US-001",
		Title:        "First task",
		FilesChanged: []string{"file1.go"},
		Learnings:    []string{"Learning 1"},
	}
	err := pw.AppendEntry(entry1)
	require.NoError(t, err)

	// Add second entry
	entry2 := ProgressEntry{
		Timestamp:    time.Date(2026, 2, 5, 11, 0, 0, 0, time.UTC),
		StoryID:      "US-002",
		Title:        "Second task",
		FilesChanged: []string{"file2.go"},
		Learnings:    []string{"Learning 2"},
	}
	err = pw.AppendEntry(entry2)
	require.NoError(t, err)

	// Verify both entries exist
	entries, err := pw.GetAllEntries()
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, "US-001", entries[0].StoryID)
	assert.Equal(t, "US-002", entries[1].StoryID)
}

func TestProgressWriter_GetCodebasePatterns(t *testing.T) {
	tmpDir := t.TempDir()
	textPath := filepath.Join(tmpDir, "progress.txt")
	yamlPath := filepath.Join(tmpDir, "progress.yaml")

	pw := NewProgressWriter(textPath, yamlPath)

	// Initialize with patterns
	content := `# Progress Log

Started: 2026-02-05 10:00:00

## Codebase Patterns
- Always use context.md before editing
- Validate input at boundaries
- Use dependency injection

## 2026-02-05 10:30:00 - US-001
**Task completed**

---
`
	err := os.WriteFile(textPath, []byte(content), 0644)
	require.NoError(t, err)

	patterns, err := pw.GetCodebasePatterns()
	require.NoError(t, err)
	require.Len(t, patterns, 3)
	assert.Equal(t, "Always use context.md before editing", patterns[0])
	assert.Equal(t, "Validate input at boundaries", patterns[1])
	assert.Equal(t, "Use dependency injection", patterns[2])
}

func TestProgressWriter_GetCodebasePatterns_NoFile(t *testing.T) {
	tmpDir := t.TempDir()
	textPath := filepath.Join(tmpDir, "progress.txt")
	yamlPath := filepath.Join(tmpDir, "progress.yaml")

	pw := NewProgressWriter(textPath, yamlPath)

	patterns, err := pw.GetCodebasePatterns()
	require.NoError(t, err)
	assert.Nil(t, patterns)
}

func TestProgressWriter_AddCodebasePattern(t *testing.T) {
	tmpDir := t.TempDir()
	textPath := filepath.Join(tmpDir, "progress.txt")
	yamlPath := filepath.Join(tmpDir, "progress.yaml")

	pw := NewProgressWriter(textPath, yamlPath)

	// Add first pattern (will create file and section)
	err := pw.AddCodebasePattern("Pattern 1")
	require.NoError(t, err)

	// Verify pattern was added
	patterns, err := pw.GetCodebasePatterns()
	require.NoError(t, err)
	require.Len(t, patterns, 1)
	assert.Equal(t, "Pattern 1", patterns[0])

	// Add second pattern
	err = pw.AddCodebasePattern("Pattern 2")
	require.NoError(t, err)

	// Verify both patterns exist
	patterns, err = pw.GetCodebasePatterns()
	require.NoError(t, err)
	require.Len(t, patterns, 2)
	assert.Equal(t, "Pattern 1", patterns[0])
	assert.Equal(t, "Pattern 2", patterns[1])
}

func TestProgressWriter_AddCodebasePattern_WithExistingContent(t *testing.T) {
	tmpDir := t.TempDir()
	textPath := filepath.Join(tmpDir, "progress.txt")
	yamlPath := filepath.Join(tmpDir, "progress.yaml")

	pw := NewProgressWriter(textPath, yamlPath)

	// Create file with existing content
	content := `# Progress Log

Started: 2026-02-05 10:00:00

---

## 2026-02-05 10:30:00 - US-001
**Task completed**

---
`
	err := os.WriteFile(textPath, []byte(content), 0644)
	require.NoError(t, err)

	// Add pattern
	err = pw.AddCodebasePattern("New pattern")
	require.NoError(t, err)

	// Verify pattern was added and existing content preserved
	updatedContent, err := os.ReadFile(textPath)
	require.NoError(t, err)
	contentStr := string(updatedContent)

	assert.Contains(t, contentStr, "## Codebase Patterns")
	assert.Contains(t, contentStr, "- New pattern")
	assert.Contains(t, contentStr, "## 2026-02-05 10:30:00 - US-001")
}

func TestProgressWriter_GetAllEntries(t *testing.T) {
	tmpDir := t.TempDir()
	textPath := filepath.Join(tmpDir, "progress.txt")
	yamlPath := filepath.Join(tmpDir, "progress.yaml")

	pw := NewProgressWriter(textPath, yamlPath)

	// Add entries
	entry1 := ProgressEntry{
		Timestamp: time.Date(2026, 2, 5, 10, 0, 0, 0, time.UTC),
		StoryID:   "US-001",
		Title:     "First task",
	}
	entry2 := ProgressEntry{
		Timestamp: time.Date(2026, 2, 5, 11, 0, 0, 0, time.UTC),
		StoryID:   "US-002",
		Title:     "Second task",
	}

	err := pw.AppendEntry(entry1)
	require.NoError(t, err)
	err = pw.AppendEntry(entry2)
	require.NoError(t, err)

	// Retrieve all entries
	entries, err := pw.GetAllEntries()
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, "US-001", entries[0].StoryID)
	assert.Equal(t, "US-002", entries[1].StoryID)
}

func TestProgressWriter_GetEntriesByFile(t *testing.T) {
	tmpDir := t.TempDir()
	textPath := filepath.Join(tmpDir, "progress.txt")
	yamlPath := filepath.Join(tmpDir, "progress.yaml")

	pw := NewProgressWriter(textPath, yamlPath)

	// Add entries with different files
	entry1 := ProgressEntry{
		Timestamp:    time.Date(2026, 2, 5, 10, 0, 0, 0, time.UTC),
		StoryID:      "US-001",
		Title:        "First task",
		FilesChanged: []string{"src/auth/login.go", "src/auth/middleware.go"},
	}
	entry2 := ProgressEntry{
		Timestamp:    time.Date(2026, 2, 5, 11, 0, 0, 0, time.UTC),
		StoryID:      "US-002",
		Title:        "Second task",
		FilesChanged: []string{"src/users/service.go"},
	}
	entry3 := ProgressEntry{
		Timestamp:    time.Date(2026, 2, 5, 12, 0, 0, 0, time.UTC),
		StoryID:      "US-003",
		Title:        "Third task",
		FilesChanged: []string{"src/auth/login.go", "src/users/service.go"},
	}

	err := pw.AppendEntry(entry1)
	require.NoError(t, err)
	err = pw.AppendEntry(entry2)
	require.NoError(t, err)
	err = pw.AppendEntry(entry3)
	require.NoError(t, err)

	// Get entries that modified src/auth/login.go
	entries, err := pw.GetEntriesByFile("src/auth/login.go")
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, "US-001", entries[0].StoryID)
	assert.Equal(t, "US-003", entries[1].StoryID)

	// Get entries that modified src/users/service.go
	entries, err = pw.GetEntriesByFile("src/users/service.go")
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, "US-002", entries[0].StoryID)
	assert.Equal(t, "US-003", entries[1].StoryID)

	// Get entries for non-existent file
	entries, err = pw.GetEntriesByFile("src/nonexistent.go")
	require.NoError(t, err)
	assert.Len(t, entries, 0)
}

func TestProgressWriter_InitializeTextFile(t *testing.T) {
	tmpDir := t.TempDir()
	textPath := filepath.Join(tmpDir, "progress.txt")
	yamlPath := filepath.Join(tmpDir, "progress.yaml")

	pw := NewProgressWriter(textPath, yamlPath)

	err := pw.initializeTextFile()
	require.NoError(t, err)

	// Verify file exists
	assert.FileExists(t, textPath)

	// Verify header content
	content, err := os.ReadFile(textPath)
	require.NoError(t, err)
	contentStr := string(content)

	assert.Contains(t, contentStr, "# Progress Log")
	assert.Contains(t, contentStr, "Started:")
}

func TestProgressWriter_AppendEntry_WithoutThreadURL(t *testing.T) {
	tmpDir := t.TempDir()
	textPath := filepath.Join(tmpDir, "progress.txt")
	yamlPath := filepath.Join(tmpDir, "progress.yaml")

	pw := NewProgressWriter(textPath, yamlPath)

	entry := ProgressEntry{
		Timestamp:    time.Date(2026, 2, 5, 10, 30, 0, 0, time.UTC),
		StoryID:      "US-001",
		Title:        "Add feature",
		FilesChanged: []string{"file.go"},
		Learnings:    []string{"Learning 1"},
		// ThreadURL intentionally omitted
	}

	err := pw.AppendEntry(entry)
	require.NoError(t, err)

	// Read text content
	textContent, err := os.ReadFile(textPath)
	require.NoError(t, err)
	textStr := string(textContent)

	// Verify entry exists but no Thread line
	assert.Contains(t, textStr, "US-001")
	assert.NotContains(t, textStr, "Thread:")
}

func TestProgressWriter_TextFileFormat(t *testing.T) {
	tmpDir := t.TempDir()
	textPath := filepath.Join(tmpDir, "progress.txt")
	yamlPath := filepath.Join(tmpDir, "progress.yaml")

	pw := NewProgressWriter(textPath, yamlPath)

	entry := ProgressEntry{
		Timestamp:    time.Date(2026, 2, 5, 14, 30, 45, 0, time.UTC),
		StoryID:      "US-001",
		Title:        "Test Task",
		FilesChanged: []string{"file1.go", "file2.go"},
		Learnings:    []string{"Learning A", "Learning B"},
		ThreadURL:    "https://example.com/thread",
	}

	err := pw.AppendEntry(entry)
	require.NoError(t, err)

	content, err := os.ReadFile(textPath)
	require.NoError(t, err)
	lines := strings.Split(string(content), "\n")

	// Verify format
	var foundHeader, foundThread, foundTitle, foundFiles, foundLearnings, foundSeparator bool

	for _, line := range lines {
		if strings.Contains(line, "## 2026-02-05 14:30:45 - US-001") {
			foundHeader = true
		}
		if strings.Contains(line, "Thread: https://example.com/thread") {
			foundThread = true
		}
		if strings.Contains(line, "**Test Task**") {
			foundTitle = true
		}
		if strings.Contains(line, "**Files Changed:**") {
			foundFiles = true
		}
		if strings.Contains(line, "**Learnings for future iterations:**") {
			foundLearnings = true
		}
		if strings.Contains(line, "---") {
			foundSeparator = true
		}
	}

	assert.True(t, foundHeader, "Should have entry header with timestamp and story ID")
	assert.True(t, foundThread, "Should have thread URL")
	assert.True(t, foundTitle, "Should have title in bold")
	assert.True(t, foundFiles, "Should have Files Changed section")
	assert.True(t, foundLearnings, "Should have Learnings section")
	assert.True(t, foundSeparator, "Should have separator line")
}
