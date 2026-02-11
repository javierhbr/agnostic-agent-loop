package openspec

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	numberedListRe = regexp.MustCompile(`^\d+\.\s+(.+)$`)
	checkboxRe     = regexp.MustCompile(`^-\s+\[[ x~]\]\s+(.+)$`)
)

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
