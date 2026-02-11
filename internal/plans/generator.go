package plans

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GenerateFromSpec reads a spec.md file and generates a phased plan.md.
// It extracts checkbox items from Requirements and Acceptance Criteria sections
// and organizes them into Setup, Implementation, and Validation phases.
func GenerateFromSpec(specPath, planPath, trackName string) error {
	requirements, acceptance, err := extractSpecItems(specPath)
	if err != nil {
		return err
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("# Plan: %s\n\n", trackName))

	b.WriteString("## Phase 1: Setup\n\n")
	b.WriteString("- [ ] Define approach and validate design\n")
	b.WriteString("- [ ] Set up scaffolding and dependencies\n")

	b.WriteString("\n## Phase 2: Implementation\n\n")
	if len(requirements) > 0 {
		for _, req := range requirements {
			b.WriteString(fmt.Sprintf("- [ ] %s\n", req))
		}
	} else {
		b.WriteString("- [ ] Core implementation\n")
		b.WriteString("- [ ] Integration with existing code\n")
	}

	b.WriteString("\n## Phase 3: Validation\n\n")
	if len(acceptance) > 0 {
		for _, ac := range acceptance {
			b.WriteString(fmt.Sprintf("- [ ] Verify: %s\n", ac))
		}
	} else {
		b.WriteString("- [ ] Write tests\n")
		b.WriteString("- [ ] Code review\n")
	}
	b.WriteString("- [ ] Update documentation\n")

	return os.WriteFile(planPath, []byte(b.String()), 0644)
}

// extractSpecItems reads a spec file and pulls out checkbox items
// from the Requirements and Acceptance Criteria sections.
func extractSpecItems(specPath string) (requirements, acceptance []string, err error) {
	f, err := os.Open(specPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open spec: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var currentSection string

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "## ") {
			switch trimmed {
			case "## Requirements":
				currentSection = "requirements"
			case "## Acceptance Criteria":
				currentSection = "acceptance"
			default:
				currentSection = ""
			}
			continue
		}

		if currentSection != "" && strings.HasPrefix(trimmed, "- [") {
			title := extractCheckboxTitle(trimmed)
			if title != "" && !isPlaceholderItem(title) {
				switch currentSection {
				case "requirements":
					requirements = append(requirements, title)
				case "acceptance":
					acceptance = append(acceptance, title)
				}
			}
		}
	}

	return requirements, acceptance, scanner.Err()
}

func extractCheckboxTitle(line string) string {
	// "- [ ] Some text" or "- [x] Some text"
	if len(line) < 6 {
		return ""
	}
	if line[4] != ']' {
		return ""
	}
	return strings.TrimSpace(line[5:])
}

func isPlaceholderItem(title string) bool {
	lower := strings.ToLower(title)
	placeholders := []string{"requirement 1", "requirement 2", "criterion 1", "criterion 2"}
	for _, p := range placeholders {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}
