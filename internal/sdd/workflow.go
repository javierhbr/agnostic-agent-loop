package sdd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// DetermineRisk resolves the effective risk level from multiple sources.
// Priority: non-empty flag > config default > error.
// If interactive is true and both flag and default are empty, returns empty string
// (caller should handle interactive prompt).
func DetermineRisk(flagValue string, defaultRisk string, interactive bool) (RiskLevel, error) {
	// Flag takes precedence
	if flagValue != "" {
		switch flagValue {
		case "low":
			return RiskLow, nil
		case "medium":
			return RiskMedium, nil
		case "high":
			return RiskHigh, nil
		case "critical":
			return RiskCritical, nil
		default:
			return "", fmt.Errorf("invalid risk level: %s", flagValue)
		}
	}

	// Config default takes precedence over interactive
	if defaultRisk != "" {
		switch defaultRisk {
		case "low":
			return RiskLow, nil
		case "medium":
			return RiskMedium, nil
		case "high":
			return RiskHigh, nil
		case "critical":
			return RiskCritical, nil
		default:
			return "", fmt.Errorf("invalid default risk level in config: %s", defaultRisk)
		}
	}

	// If interactive mode and nothing provided, return error for caller to handle prompt
	if interactive {
		return "", fmt.Errorf("no risk level provided (will be prompted)")
	}

	// Non-interactive with no flag or default is an error
	return "", fmt.Errorf("provide --risk or set sdd.default_risk in config")
}

// InitiativeManager manages SDD initiatives and their lifecycle.
type InitiativeManager struct {
	baseDir string
}

// NewInitiativeManager creates a new initiative manager for a given directory.
func NewInitiativeManager(baseDir string) *InitiativeManager {
	return &InitiativeManager{baseDir: baseDir}
}

// Create creates a new initiative with the given name and risk level.
func (m *InitiativeManager) Create(name string, risk RiskLevel) (*Initiative, error) {
	// Ensure directory exists
	if err := os.MkdirAll(m.baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create initiatives directory: %w", err)
	}

	// Generate a kebab-case ID from the name
	id := toKebab(name)

	// Create the initiative
	init := &Initiative{
		ID:        id,
		Name:      name,
		Risk:      risk,
		Workflow:  RiskToWorkflow(risk),
		Status:    "Planned",
		CreatedAt: time.Now(),
	}

	// Set the current agent to the first agent in the workflow
	agents := WorkflowAgents(init.Workflow)
	if len(agents) > 0 {
		init.CurrentAgent = agents[0]
	}

	// Write to file
	filePath := filepath.Join(m.baseDir, id+".yaml")
	data, err := yaml.Marshal(init)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal initiative: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write initiative file: %w", err)
	}

	return init, nil
}

// Get retrieves an initiative by ID or name.
func (m *InitiativeManager) Get(id string) (*Initiative, error) {
	// Try to read as a file first
	filePath := filepath.Join(m.baseDir, id+".yaml")
	return m.readInitiative(filePath)
}

// List returns all initiatives in the directory.
func (m *InitiativeManager) List() ([]Initiative, error) {
	var initiatives []Initiative

	files, err := os.ReadDir(m.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return initiatives, nil // Return empty list if directory doesn't exist
		}
		return nil, fmt.Errorf("failed to read initiatives directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !isYAMLFile(file.Name()) {
			continue
		}

		filePath := filepath.Join(m.baseDir, file.Name())
		init, err := m.readInitiative(filePath)
		if err != nil {
			continue // Skip files that can't be parsed
		}

		initiatives = append(initiatives, *init)
	}

	return initiatives, nil
}

// Advance moves an initiative to the next agent in the workflow sequence.
func (m *InitiativeManager) Advance(id string) error {
	init, err := m.Get(id)
	if err != nil {
		return err
	}

	// Find current agent index
	agents := WorkflowAgents(init.Workflow)
	currentIdx := -1
	for i, agent := range agents {
		if agent == init.CurrentAgent {
			currentIdx = i
			break
		}
	}

	if currentIdx < 0 {
		return fmt.Errorf("current agent %s not found in workflow", init.CurrentAgent)
	}

	// Move to next agent
	if currentIdx+1 < len(agents) {
		init.CurrentAgent = agents[currentIdx+1]
	} else {
		// All agents complete, mark as done
		init.Status = "Done"
	}

	// Write updated initiative
	filePath := filepath.Join(m.baseDir, init.ID+".yaml")
	data, err := yaml.Marshal(init)
	if err != nil {
		return fmt.Errorf("failed to marshal initiative: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write initiative file: %w", err)
	}

	return nil
}

// Helper functions

func (m *InitiativeManager) readInitiative(filePath string) (*Initiative, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var init Initiative
	if err := yaml.Unmarshal(data, &init); err != nil {
		return nil, fmt.Errorf("failed to parse initiative YAML: %w", err)
	}

	return &init, nil
}

func isYAMLFile(name string) bool {
	return len(name) > 5 && name[len(name)-5:] == ".yaml" || name[len(name)-4:] == ".yml"
}
