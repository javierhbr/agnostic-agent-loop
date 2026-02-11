package plans

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PlanTaskStatus represents the status of a task in a plan.
type PlanTaskStatus string

const (
	PlanTaskPending    PlanTaskStatus = "pending"
	PlanTaskInProgress PlanTaskStatus = "in_progress"
	PlanTaskDone       PlanTaskStatus = "done"
)

// PlanTask is a single task item in a plan phase.
type PlanTask struct {
	Title  string
	Status PlanTaskStatus
	Line   int // 1-based line number in the file
}

// Phase groups tasks under a named section.
type Phase struct {
	Name  string
	Tasks []PlanTask
}

// Plan represents a parsed plan.md file.
type Plan struct {
	Title  string
	Phases []Phase
	Path   string
}

// ParseFile reads a plan.md file and extracts phases and tasks.
func ParseFile(path string) (*Plan, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open plan: %w", err)
	}
	defer f.Close()

	plan := &Plan{Path: path}
	scanner := bufio.NewScanner(f)
	lineNum := 0
	var currentPhase *Phase

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Parse title (# heading)
		if strings.HasPrefix(trimmed, "# ") && !strings.HasPrefix(trimmed, "## ") {
			plan.Title = strings.TrimPrefix(trimmed, "# ")
			// Strip common prefixes like "Plan: "
			plan.Title = strings.TrimPrefix(plan.Title, "Plan: ")
			continue
		}

		// Parse phase heading (## heading)
		if strings.HasPrefix(trimmed, "## ") {
			phaseName := strings.TrimPrefix(trimmed, "## ")
			plan.Phases = append(plan.Phases, Phase{Name: phaseName})
			currentPhase = &plan.Phases[len(plan.Phases)-1]
			continue
		}

		// Parse task items (- [ ], - [x], - [~])
		if currentPhase != nil {
			if task, ok := parseTaskLine(trimmed, lineNum); ok {
				currentPhase.Tasks = append(currentPhase.Tasks, task)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return plan, nil
}

// parseTaskLine tries to parse a markdown checkbox line.
func parseTaskLine(line string, lineNum int) (PlanTask, bool) {
	if !strings.HasPrefix(line, "- [") {
		return PlanTask{}, false
	}

	// Expect "- [x] ...", "- [ ] ...", or "- [~] ..."
	if len(line) < 6 {
		return PlanTask{}, false
	}

	marker := line[3:4]
	var status PlanTaskStatus
	switch marker {
	case " ":
		status = PlanTaskPending
	case "x", "X":
		status = PlanTaskDone
	case "~":
		status = PlanTaskInProgress
	default:
		return PlanTask{}, false
	}

	if line[4] != ']' {
		return PlanTask{}, false
	}

	title := strings.TrimSpace(line[5:])

	return PlanTask{
		Title:  title,
		Status: status,
		Line:   lineNum,
	}, true
}

// NextTask returns the first pending task across all phases, or nil if all done.
func (p *Plan) NextTask() (*PlanTask, *Phase) {
	for i := range p.Phases {
		for j := range p.Phases[i].Tasks {
			if p.Phases[i].Tasks[j].Status == PlanTaskPending {
				return &p.Phases[i].Tasks[j], &p.Phases[i]
			}
		}
	}
	return nil, nil
}

// Progress returns done/total counts.
func (p *Plan) Progress() (done, total int) {
	for _, phase := range p.Phases {
		for _, task := range phase.Tasks {
			total++
			if task.Status == PlanTaskDone {
				done++
			}
		}
	}
	return
}
