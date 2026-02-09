package plans

import (
	"fmt"
	"os"
	"strings"
)

// statusMarker returns the checkbox character for a given status.
func statusMarker(s PlanTaskStatus) string {
	switch s {
	case PlanTaskDone:
		return "x"
	case PlanTaskInProgress:
		return "~"
	default:
		return " "
	}
}

// UpdateTaskStatus updates a task's status marker in the plan file by line number.
func UpdateTaskStatus(path string, lineNum int, newStatus PlanTaskStatus) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read plan: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	if lineNum < 1 || lineNum > len(lines) {
		return fmt.Errorf("line %d out of range (file has %d lines)", lineNum, len(lines))
	}

	idx := lineNum - 1
	line := lines[idx]
	trimmed := strings.TrimSpace(line)

	// Verify it's a checkbox line
	if !strings.HasPrefix(trimmed, "- [") || len(trimmed) < 6 || trimmed[4] != ']' {
		return fmt.Errorf("line %d is not a checkbox task", lineNum)
	}

	// Preserve leading whitespace
	indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
	rest := trimmed[5:] // everything after "- [x]"
	lines[idx] = fmt.Sprintf("%s- [%s]%s", indent, statusMarker(newStatus), rest)

	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}
