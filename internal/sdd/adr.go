package sdd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ADRManager manages ADR lifecycle and persistence.
type ADRManager struct {
	baseDir string
}

// NewADRManager creates a new ADR manager for a given base directory.
func NewADRManager(baseDir string) *ADRManager {
	return &ADRManager{baseDir: baseDir}
}

// Create creates a new ADR file with the given title and scope.
// Returns the created ADR with an auto-generated ID.
func (m *ADRManager) Create(title, scope string) (*ADR, error) {
	// Ensure directory exists
	if err := os.MkdirAll(m.baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create ADR directory: %w", err)
	}

	// Find the next ADR number
	nextID := m.getNextADRID()
	adrID := fmt.Sprintf("ADR-%03d", nextID)

	// Create kebab-case filename from title
	kebabTitle := toKebab(title)
	filename := fmt.Sprintf("%s-%s.md", adrID, kebabTitle)
	filePath := filepath.Join(m.baseDir, filename)

	// Create the ADR object
	adr := &ADR{
		ID:        adrID,
		Title:     title,
		Status:    ADRStatusProposed,
		Scope:     scope,
		CreatedAt: time.Now(),
		FilePath:  filePath,
	}

	// Write the ADR template to file
	content := m.adrTemplate(adr)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write ADR file: %w", err)
	}

	return adr, nil
}

// Get retrieves an ADR by ID.
func (m *ADRManager) Get(id string) (*ADR, error) {
	// List all ADR files and find the one matching the ID
	files, err := os.ReadDir(m.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("ADR directory does not exist")
		}
		return nil, fmt.Errorf("failed to read ADR directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		if strings.HasPrefix(file.Name(), id+"-") {
			filePath := filepath.Join(m.baseDir, file.Name())
			return m.parseADRFile(filePath)
		}
	}

	return nil, fmt.Errorf("ADR %s not found", id)
}

// List returns all ADRs in the directory.
func (m *ADRManager) List() ([]ADR, error) {
	var adrs []ADR

	files, err := os.ReadDir(m.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return adrs, nil // Return empty list if directory doesn't exist
		}
		return nil, fmt.Errorf("failed to read ADR directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(m.baseDir, file.Name())
		adr, err := m.parseADRFile(filePath)
		if err != nil {
			continue // Skip files that can't be parsed
		}

		adrs = append(adrs, *adr)
	}

	return adrs, nil
}

// ListBlocking returns all ADRs with non-empty Blocks.
func (m *ADRManager) ListBlocking() ([]ADR, error) {
	all, err := m.List()
	if err != nil {
		return nil, err
	}

	var blocking []ADR
	for _, adr := range all {
		if len(adr.Blocks) > 0 {
			blocking = append(blocking, adr)
		}
	}

	return blocking, nil
}

// Resolve marks an ADR as approved and clears blocked_by on dependent specs in the graph.
func (m *ADRManager) Resolve(id string, graph *SpecGraph) error {
	adr, err := m.Get(id)
	if err != nil {
		return err
	}

	adr.Status = ADRStatusApproved
	adr.ResolvedAt = time.Now()

	// Update the ADR file with new status
	content := m.adrTemplate(adr)
	if err := os.WriteFile(adr.FilePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to update ADR file: %w", err)
	}

	// Clear this ADR from blocked_by on all dependent specs
	for nodeID, node := range graph.Nodes {
		// Remove this ADR ID from BlockedBy
		var newBlockedBy []string
		for _, blockerID := range node.BlockedBy {
			if blockerID != id {
				newBlockedBy = append(newBlockedBy, blockerID)
			}
		}
		if len(newBlockedBy) != len(node.BlockedBy) {
			// This node was blocked by the resolved ADR
			node.BlockedBy = newBlockedBy
			if len(newBlockedBy) == 0 {
				// No longer blocked, can move back to prior state
				node.Status = SpecStatusApproved
			}
			graph.Upsert(node)
		}
		graph.Nodes[nodeID] = node
	}

	return nil
}

// Helper functions

func (m *ADRManager) getNextADRID() int {
	files, err := os.ReadDir(m.baseDir)
	if err != nil {
		return 1
	}

	maxID := 0
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "ADR-") && strings.HasSuffix(file.Name(), ".md") {
			// Extract the number
			parts := strings.Split(file.Name(), "-")
			if len(parts) > 1 {
				var numStr string
				for _, ch := range parts[1] {
					if ch >= '0' && ch <= '9' {
						numStr += string(ch)
					} else {
						break
					}
				}
				if numStr != "" {
					var id int
					fmt.Sscanf(numStr, "%d", &id)
					if id > maxID {
						maxID = id
					}
				}
			}
		}
	}

	return maxID + 1
}

func (m *ADRManager) parseADRFile(filePath string) (*ADR, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	content := string(data)
	adr := &ADR{FilePath: filePath}

	// Extract ID from filename
	filename := filepath.Base(filePath)
	if strings.HasPrefix(filename, "ADR-") {
		parts := strings.Split(filename, "-")
		if len(parts) >= 2 {
			var numStr string
			for _, ch := range parts[1] {
				if ch >= '0' && ch <= '9' {
					numStr += string(ch)
				} else {
					break
				}
			}
			if numStr != "" {
				adr.ID = fmt.Sprintf("ADR-%s", numStr)
			}
		}
	}

	// Extract title from first H1 heading
	titleRegex := regexp.MustCompile(`^# (.*?)$`)
	if matches := titleRegex.FindStringSubmatch(content); len(matches) > 1 {
		adr.Title = matches[1]
	}

	// Extract status from ## Status section
	statusRegex := regexp.MustCompile(`(?i)## Status\s*\n([^\n]+)`)
	if matches := statusRegex.FindStringSubmatch(content); len(matches) > 1 {
		status := strings.TrimSpace(matches[1])
		switch status {
		case "Proposed":
			adr.Status = ADRStatusProposed
		case "In Review":
			adr.Status = ADRStatusInReview
		case "Approved":
			adr.Status = ADRStatusApproved
		case "Rejected":
			adr.Status = ADRStatusRejected
		}
	}

	// Extract scope if present
	scopeRegex := regexp.MustCompile(`(?i)Scope:\s*([^\n]+)`)
	if matches := scopeRegex.FindStringSubmatch(content); len(matches) > 1 {
		adr.Scope = strings.TrimSpace(matches[1])
	}

	return adr, nil
}

func (m *ADRManager) adrTemplate(adr *ADR) string {
	return fmt.Sprintf(`# %s: %s

## Status
%s

## Date
%s

## Context
[What is the situation? Why is a decision needed now? What constraints exist?]

## Decision Drivers
- [Driver 1]
- [Driver 2]
- [Driver 3]

## Options Considered

### Option A: [Name]
[Description]

Pros:
- [Pro 1]

Cons:
- [Con 1]

### Option B: [Name]
[Description]

Pros:
- [Pro 1]

Cons:
- [Con 1]

## Decision
[Leave as PENDING until approved]

## Consequences
[Fill in after decision is made]

## Owner
%s

## Blocks
%s

## Scope
%s

## References
[Link to relevant context, if applicable]
`,
		adr.ID,
		adr.Title,
		string(adr.Status),
		adr.CreatedAt.Format("2006-01-02"),
		adr.Owner,
		formatBlocks(adr.Blocks),
		adr.Scope,
	)
}

func formatBlocks(blocks []string) string {
	if len(blocks) == 0 {
		return "- (none)"
	}
	var lines []string
	for _, block := range blocks {
		lines = append(lines, fmt.Sprintf("- %s", block))
	}
	return strings.Join(lines, "\n")
}

func toKebab(s string) string {
	// Convert to lowercase and replace spaces and underscores with dashes
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	// Remove any non-alphanumeric characters except dashes
	var result []rune
	for _, ch := range s {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-' {
			result = append(result, ch)
		}
	}
	s = string(result)
	// Remove leading/trailing/multiple dashes
	s = strings.Trim(s, "-")
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return s
}
