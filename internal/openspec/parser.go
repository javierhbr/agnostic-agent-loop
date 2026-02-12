package openspec

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	numberedListRe = regexp.MustCompile(`^\d+\.\s+(.+)$`)
	checkboxRe     = regexp.MustCompile(`^-\s+\[[ x~]\]\s+(.+)$`)
	linkRefRe      = regexp.MustCompile(`\[.*?\]\(\.?/?(tasks/[^\)]+)\)`)
	// Matches the parenthetical "(ver [text](path))" suffix
	verLinkSuffixRe = regexp.MustCompile(`\s*\(ver\s+\[.*?\]\(.*?\)\)\s*$`)
	// Matches a bare markdown link "[text](path)" at end of title
	bareLinkSuffixRe = regexp.MustCompile(`\s*\[.*?\]\(.*?\)\s*$`)
)

// TaskEntry represents a parsed line from tasks.md — title plus optional file reference.
type TaskEntry struct {
	Title   string // cleaned task title
	FileRef string // relative path, e.g. "tasks/01-setup.md" — empty if no reference
}

// TaskDetail holds structured information parsed from an individual task file.
type TaskDetail struct {
	Title         string   // from H1 heading
	Description   string   // content under ## Description
	Prerequisites []string // items under ## Prerequisites
	Acceptance    []string // items under ## Acceptance Criteria
	Notes         string   // content under ## Technical Notes
}

// ParseTasksFile reads a tasks.md file and extracts task titles.
// Supports numbered lists (1. Task) and checkbox lists (- [ ] Task).
func ParseTasksFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open tasks file: %w", err)
	}
	defer f.Close()

	var tasks []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if title := extractTaskTitle(line); title != "" {
			tasks = append(tasks, title)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read tasks file: %w", err)
	}
	if len(tasks) == 0 {
		return nil, fmt.Errorf("no tasks found in %s", path)
	}
	return tasks, nil
}

func extractTaskTitle(line string) string {
	if m := numberedListRe.FindStringSubmatch(line); len(m) == 2 {
		return strings.TrimSpace(m[1])
	}
	if m := checkboxRe.FindStringSubmatch(line); len(m) == 2 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

// ParseTasksFileStructured reads tasks.md and returns TaskEntry objects
// with titles and optional file references extracted from markdown links.
func ParseTasksFileStructured(path string) ([]TaskEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open tasks file: %w", err)
	}
	defer f.Close()

	var entries []TaskEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if entry := extractTaskEntry(line); entry != nil {
			entries = append(entries, *entry)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read tasks file: %w", err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("no tasks found in %s", path)
	}
	return entries, nil
}

// extractTaskEntry parses a single line from tasks.md into a TaskEntry.
// Returns nil if the line is not a task line.
func extractTaskEntry(line string) *TaskEntry {
	var rawTitle string
	if m := numberedListRe.FindStringSubmatch(line); len(m) == 2 {
		rawTitle = strings.TrimSpace(m[1])
	} else if m := checkboxRe.FindStringSubmatch(line); len(m) == 2 {
		rawTitle = strings.TrimSpace(m[1])
	}
	if rawTitle == "" {
		return nil
	}

	entry := &TaskEntry{}

	// Extract file reference if present
	if m := linkRefRe.FindStringSubmatch(rawTitle); len(m) == 2 {
		entry.FileRef = m[1]
	}

	entry.Title = cleanTitleFromLinks(rawTitle)
	return entry
}

// cleanTitleFromLinks strips parenthetical link references and bare markdown
// link syntax from the end of a title string.
func cleanTitleFromLinks(raw string) string {
	// First try "(ver [text](path))" pattern
	cleaned := verLinkSuffixRe.ReplaceAllString(raw, "")
	if cleaned != raw {
		return strings.TrimSpace(cleaned)
	}
	// Then try bare "[text](path)" pattern
	cleaned = bareLinkSuffixRe.ReplaceAllString(raw, "")
	return strings.TrimSpace(cleaned)
}

// HasTasksDir checks whether a tasks/ subdirectory exists next to tasks.md.
func HasTasksDir(tasksFilePath string) bool {
	dir := filepath.Dir(tasksFilePath)
	tasksDir := filepath.Join(dir, "tasks")
	info, err := os.Stat(tasksDir)
	return err == nil && info.IsDir()
}

// ParseTaskDetailFile reads an individual task markdown file and extracts
// structured sections (Description, Prerequisites, Acceptance Criteria, Technical Notes).
func ParseTaskDetailFile(path string) (*TaskDetail, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read task file %s: %w", path, err)
	}

	detail := &TaskDetail{}
	lines := strings.Split(string(data), "\n")

	var currentSection string
	var sectionContent []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// H1 = task title
		if strings.HasPrefix(trimmed, "# ") && !strings.HasPrefix(trimmed, "## ") {
			detail.Title = strings.TrimSpace(strings.TrimPrefix(trimmed, "# "))
			continue
		}

		// H2 = section boundary
		if strings.HasPrefix(trimmed, "## ") {
			flushSection(detail, currentSection, sectionContent)
			currentSection = normalizeSectionName(strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")))
			sectionContent = nil
			continue
		}

		// Accumulate content for current section
		if currentSection != "" {
			sectionContent = append(sectionContent, line)
		}
	}
	// Flush final section
	flushSection(detail, currentSection, sectionContent)

	return detail, nil
}

// normalizeSectionName maps heading text to canonical section names.
func normalizeSectionName(heading string) string {
	lower := strings.ToLower(heading)
	switch {
	case strings.Contains(lower, "description"):
		return "description"
	case strings.Contains(lower, "prerequisite") || strings.Contains(lower, "pre-requisite"):
		return "prerequisites"
	case strings.Contains(lower, "acceptance"):
		return "acceptance"
	case strings.Contains(lower, "technical") || lower == "notes":
		return "notes"
	default:
		return lower
	}
}

// flushSection assigns accumulated content to the appropriate TaskDetail field.
func flushSection(detail *TaskDetail, section string, lines []string) {
	if section == "" || len(lines) == 0 {
		return
	}
	switch section {
	case "description":
		detail.Description = strings.TrimSpace(strings.Join(lines, "\n"))
	case "prerequisites":
		detail.Prerequisites = extractListItems(lines)
	case "acceptance":
		detail.Acceptance = extractListItems(lines)
	case "notes":
		detail.Notes = strings.TrimSpace(strings.Join(lines, "\n"))
	}
}

// extractListItems extracts bullet and checkbox items from lines.
func extractListItems(lines []string) []string {
	var items []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- [") {
			// Checkbox: "- [ ] text" or "- [x] text"
			if _, after, found := strings.Cut(trimmed, "] "); found {
				items = append(items, strings.TrimSpace(after))
			}
		} else if strings.HasPrefix(trimmed, "- ") {
			items = append(items, strings.TrimSpace(trimmed[2:]))
		} else if strings.HasPrefix(trimmed, "* ") {
			items = append(items, strings.TrimSpace(trimmed[2:]))
		}
	}
	return items
}
