package checkpoint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// Checkpoint represents a saved state during task execution
type Checkpoint struct {
	TaskID        string    `json:"task_id"`
	Iteration     int       `json:"iteration"`
	TokensUsed    int       `json:"tokens_used"`
	CreatedAt     time.Time `json:"created_at"`
	Agent         string    `json:"agent"`
	Output        string    `json:"output"`
	CriteriaMet   []string  `json:"criteria_met"`
	CriteriaLeft  []string  `json:"criteria_left"`
	FilesModified []string  `json:"files_modified"`
	Learnings     []string  `json:"learnings"`
	Notes         string    `json:"notes"`
}

// Manager handles checkpoint creation and retrieval
type Manager struct {
	checkpointDir string
}

// NewManager creates a new checkpoint manager
func NewManager(checkpointDir string) *Manager {
	if checkpointDir == "" {
		checkpointDir = ".agentic/checkpoints"
	}
	return &Manager{checkpointDir: checkpointDir}
}

// Save creates a checkpoint for the current task state
func (m *Manager) Save(checkpoint *Checkpoint) error {
	if err := os.MkdirAll(m.checkpointDir, 0755); err != nil {
		return fmt.Errorf("failed to create checkpoint directory: %w", err)
	}

	filename := m.getCheckpointPath(checkpoint.TaskID, checkpoint.Iteration)
	data, err := json.MarshalIndent(checkpoint, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal checkpoint: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write checkpoint: %w", err)
	}

	// Also save as latest checkpoint for this task
	latestPath := m.getLatestCheckpointPath(checkpoint.TaskID)
	if err := os.WriteFile(latestPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write latest checkpoint: %w", err)
	}

	return nil
}

// Load retrieves the latest checkpoint for a task
func (m *Manager) Load(taskID string) (*Checkpoint, error) {
	latestPath := m.getLatestCheckpointPath(taskID)

	data, err := os.ReadFile(latestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No checkpoint exists
		}
		return nil, fmt.Errorf("failed to read checkpoint: %w", err)
	}

	var checkpoint Checkpoint
	if err := json.Unmarshal(data, &checkpoint); err != nil {
		return nil, fmt.Errorf("failed to unmarshal checkpoint: %w", err)
	}

	return &checkpoint, nil
}

// LoadIteration retrieves a specific iteration checkpoint
func (m *Manager) LoadIteration(taskID string, iteration int) (*Checkpoint, error) {
	path := m.getCheckpointPath(taskID, iteration)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read checkpoint: %w", err)
	}

	var checkpoint Checkpoint
	if err := json.Unmarshal(data, &checkpoint); err != nil {
		return nil, fmt.Errorf("failed to unmarshal checkpoint: %w", err)
	}

	return &checkpoint, nil
}

// List returns all checkpoints for a task
func (m *Manager) List(taskID string) ([]Checkpoint, error) {
	pattern := filepath.Join(m.checkpointDir, fmt.Sprintf("%s-*.json", taskID))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to list checkpoints: %w", err)
	}

	var checkpoints []Checkpoint
	for _, match := range matches {
		// Skip the latest symlink/file
		if filepath.Base(match) == fmt.Sprintf("%s-latest.json", taskID) {
			continue
		}

		data, err := os.ReadFile(match)
		if err != nil {
			continue
		}

		var checkpoint Checkpoint
		if err := json.Unmarshal(data, &checkpoint); err != nil {
			continue
		}

		checkpoints = append(checkpoints, checkpoint)
	}

	return checkpoints, nil
}

// Delete removes a checkpoint
func (m *Manager) Delete(taskID string, iteration int) error {
	path := m.getCheckpointPath(taskID, iteration)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete checkpoint: %w", err)
	}
	return nil
}

// DeleteAll removes all checkpoints for a task
func (m *Manager) DeleteAll(taskID string) error {
	pattern := filepath.Join(m.checkpointDir, fmt.Sprintf("%s-*.json", taskID))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to list checkpoints: %w", err)
	}

	for _, match := range matches {
		if err := os.Remove(match); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete checkpoint %s: %w", match, err)
		}
	}

	return nil
}

// ShouldCheckpoint determines if a checkpoint should be created based on token usage
func (m *Manager) ShouldCheckpoint(tokensUsed int, tokenLimit int, iteration int) bool {
	return m.ShouldCheckpointWithThresholds(tokensUsed, tokenLimit, iteration, 5, []float64{0.5, 0.75, 0.9})
}

// ShouldCheckpointWithThresholds determines if a checkpoint should be created with custom settings
func (m *Manager) ShouldCheckpointWithThresholds(tokensUsed int, tokenLimit int, iteration int, iterationInterval int, thresholds []float64) bool {
	// Checkpoint every N iterations
	if iterationInterval > 0 && iteration%iterationInterval == 0 {
		return true
	}

	// Checkpoint at specified percentage thresholds
	if len(thresholds) > 0 {
		percentage := float64(tokensUsed) / float64(tokenLimit)

		for _, threshold := range thresholds {
			// Check if we're at the threshold (within 5% range)
			if percentage >= threshold && percentage < threshold+0.05 {
				return true
			}
		}
	}

	return false
}

// GetProgress calculates progress percentage based on criteria met
func (m *Manager) GetProgress(checkpoint *Checkpoint, totalCriteria int) float64 {
	if totalCriteria == 0 {
		return 0
	}
	return float64(len(checkpoint.CriteriaMet)) / float64(totalCriteria) * 100
}

// Helper functions

func (m *Manager) getCheckpointPath(taskID string, iteration int) string {
	filename := fmt.Sprintf("%s-%03d.json", taskID, iteration)
	return filepath.Join(m.checkpointDir, filename)
}

func (m *Manager) getLatestCheckpointPath(taskID string) string {
	filename := fmt.Sprintf("%s-latest.json", taskID)
	return filepath.Join(m.checkpointDir, filename)
}

// CreateFromResult creates a checkpoint from an agent execution result
func CreateFromResult(taskID string, iteration int, agent string, result *models.AgentExecutionResult, task *models.Task) *Checkpoint {
	return &Checkpoint{
		TaskID:        taskID,
		Iteration:     iteration,
		TokensUsed:    result.TokensUsed,
		CreatedAt:     time.Now(),
		Agent:         agent,
		Output:        result.Output,
		CriteriaMet:   result.CriteriaMet,
		CriteriaLeft:  result.CriteriaFailed,
		FilesModified: result.FilesModified,
		Notes:         fmt.Sprintf("Iteration %d: %d/%d criteria met", iteration, len(result.CriteriaMet), len(task.Acceptance)),
	}
}
