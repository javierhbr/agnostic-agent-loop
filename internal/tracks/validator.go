package tracks

import (
	"fmt"
	"os"
	"strings"
)

// RequiredSection represents a section expected in a spec document.
type RequiredSection struct {
	Name    string
	Header  string // markdown header to search for (e.g. "## Purpose")
	Warning bool   // if true, missing section is a warning, not a blocker
}

// DefaultRequiredSections defines the sections a complete spec should have.
var DefaultRequiredSections = []RequiredSection{
	{Name: "purpose", Header: "## Purpose"},
	{Name: "constraints", Header: "## Constraints"},
	{Name: "success", Header: "## Success Criteria"},
	{Name: "alternatives", Header: "## Alternatives Considered", Warning: true},
	{Name: "design", Header: "## Design"},
	{Name: "requirements", Header: "## Requirements"},
	{Name: "acceptance", Header: "## Acceptance Criteria"},
}

// ValidationReport describes the completeness of a spec.
type ValidationReport struct {
	Complete bool
	Present  []string
	Missing  []string
	Warnings []string
}

// ValidateSpec checks a spec file for required sections with content.
func ValidateSpec(specPath string) (*ValidationReport, error) {
	data, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec: %w", err)
	}

	return ValidateSpecContent(string(data)), nil
}

// ValidateSpecContent checks spec content string for required sections.
func ValidateSpecContent(content string) *ValidationReport {
	report := &ValidationReport{Complete: true}

	for _, section := range DefaultRequiredSections {
		if sectionHasContent(content, section.Header) {
			report.Present = append(report.Present, section.Name)
		} else if section.Warning {
			report.Warnings = append(report.Warnings, section.Name)
		} else {
			report.Missing = append(report.Missing, section.Name)
			report.Complete = false
		}
	}

	return report
}

// sectionHasContent checks if a markdown section header exists and has
// non-placeholder content below it.
func sectionHasContent(content, header string) bool {
	idx := strings.Index(content, header)
	if idx == -1 {
		return false
	}

	// Get content after the header line
	after := content[idx+len(header):]
	// Find the next section header or end of file
	nextHeader := strings.Index(after, "\n## ")
	var sectionBody string
	if nextHeader == -1 {
		sectionBody = after
	} else {
		sectionBody = after[:nextHeader]
	}

	// Strip comments and whitespace
	lines := strings.Split(sectionBody, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "<!--") || strings.HasPrefix(trimmed, "-->") {
			continue
		}
		// Check it's not just a placeholder
		if isPlaceholder(trimmed) {
			continue
		}
		return true
	}

	return false
}

// isPlaceholder detects common placeholder text that doesn't count as real content.
func isPlaceholder(line string) bool {
	placeholders := []string{
		"requirement 1", "requirement 2",
		"criterion 1", "criterion 2",
		"describe what", "high-level structure",
		"what pieces", "how does information",
		"what can go wrong", "how do we validate",
		"why are we building", "what must we work within",
		"what does \"done\"", "what approaches were explored",
	}
	lower := strings.ToLower(line)
	for _, p := range placeholders {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}
