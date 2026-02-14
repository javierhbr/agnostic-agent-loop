package checkpoint

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

func TestCheckpointManager_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	checkpoint := &Checkpoint{
		TaskID:       "TASK-001",
		Iteration:    3,
		TokensUsed:   5000,
		CreatedAt:    time.Now(),
		Agent:        "claude-code",
		Output:       "Test output",
		CriteriaMet:  []string{"Tests pass", "API works"},
		CriteriaLeft: []string{"Linter clean"},
	}

	// Save checkpoint
	err := manager.Save(checkpoint)
	if err != nil {
		t.Fatalf("Failed to save checkpoint: %v", err)
	}

	// Load checkpoint
	loaded, err := manager.Load("TASK-001")
	if err != nil {
		t.Fatalf("Failed to load checkpoint: %v", err)
	}

	if loaded == nil {
		t.Fatal("Expected checkpoint, got nil")
	}

	if loaded.TaskID != checkpoint.TaskID {
		t.Errorf("Expected TaskID %s, got %s", checkpoint.TaskID, loaded.TaskID)
	}

	if loaded.Iteration != checkpoint.Iteration {
		t.Errorf("Expected iteration %d, got %d", checkpoint.Iteration, loaded.Iteration)
	}

	if loaded.TokensUsed != checkpoint.TokensUsed {
		t.Errorf("Expected tokens %d, got %d", checkpoint.TokensUsed, loaded.TokensUsed)
	}
}

func TestCheckpointManager_LoadNonexistent(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	checkpoint, err := manager.Load("NONEXISTENT")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if checkpoint != nil {
		t.Error("Expected nil checkpoint for nonexistent task")
	}
}

func TestCheckpointManager_List(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	// Create multiple checkpoints
	for i := 1; i <= 3; i++ {
		checkpoint := &Checkpoint{
			TaskID:     "TASK-001",
			Iteration:  i,
			TokensUsed: i * 1000,
			CreatedAt:  time.Now(),
			Agent:      "claude-code",
		}
		if err := manager.Save(checkpoint); err != nil {
			t.Fatalf("Failed to save checkpoint %d: %v", i, err)
		}
	}

	// List checkpoints
	checkpoints, err := manager.List("TASK-001")
	if err != nil {
		t.Fatalf("Failed to list checkpoints: %v", err)
	}

	if len(checkpoints) != 3 {
		t.Errorf("Expected 3 checkpoints, got %d", len(checkpoints))
	}
}

func TestCheckpointManager_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	checkpoint := &Checkpoint{
		TaskID:    "TASK-001",
		Iteration: 1,
		CreatedAt: time.Now(),
	}

	if err := manager.Save(checkpoint); err != nil {
		t.Fatalf("Failed to save checkpoint: %v", err)
	}

	// Delete checkpoint
	if err := manager.Delete("TASK-001", 1); err != nil {
		t.Fatalf("Failed to delete checkpoint: %v", err)
	}

	// Verify deleted
	loaded, err := manager.LoadIteration("TASK-001", 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if loaded != nil {
		t.Error("Expected checkpoint to be deleted")
	}
}

func TestCheckpointManager_DeleteAll(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	// Create multiple checkpoints
	for i := 1; i <= 3; i++ {
		checkpoint := &Checkpoint{
			TaskID:    "TASK-001",
			Iteration: i,
			CreatedAt: time.Now(),
		}
		if err := manager.Save(checkpoint); err != nil {
			t.Fatalf("Failed to save checkpoint %d: %v", i, err)
		}
	}

	// Delete all
	if err := manager.DeleteAll("TASK-001"); err != nil {
		t.Fatalf("Failed to delete all checkpoints: %v", err)
	}

	// Verify all deleted
	checkpoints, err := manager.List("TASK-001")
	if err != nil {
		t.Fatalf("Failed to list checkpoints: %v", err)
	}

	if len(checkpoints) != 0 {
		t.Errorf("Expected 0 checkpoints, got %d", len(checkpoints))
	}
}

func TestCheckpointManager_ShouldCheckpoint(t *testing.T) {
	manager := NewManager("")

	tests := []struct {
		name        string
		tokensUsed  int
		tokenLimit  int
		iteration   int
		shouldCheck bool
	}{
		{"Every 5 iterations", 1000, 10000, 5, true},
		{"Every 5 iterations - 10", 1000, 10000, 10, true},
		{"Not on iteration 3", 1000, 10000, 3, false},
		{"50% threshold", 5000, 10000, 3, true},
		{"75% threshold", 7500, 10000, 3, true},
		{"90% threshold", 9000, 10000, 3, true},
		{"Below 50%", 4000, 10000, 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.ShouldCheckpoint(tt.tokensUsed, tt.tokenLimit, tt.iteration)
			if result != tt.shouldCheck {
				t.Errorf("Expected %v, got %v", tt.shouldCheck, result)
			}
		})
	}
}

func TestCheckpointManager_GetProgress(t *testing.T) {
	manager := NewManager("")

	checkpoint := &Checkpoint{
		CriteriaMet:  []string{"Test 1", "Test 2"},
		CriteriaLeft: []string{"Test 3", "Test 4"},
	}

	progress := manager.GetProgress(checkpoint, 4)
	expected := 50.0

	if progress != expected {
		t.Errorf("Expected progress %.1f%%, got %.1f%%", expected, progress)
	}
}

func TestCreateFromResult(t *testing.T) {
	result := &models.AgentExecutionResult{
		Output:         "Test output",
		Success:        false,
		CriteriaMet:    []string{"Test 1"},
		CriteriaFailed: []string{"Test 2", "Test 3"},
		FilesModified:  []string{"file1.go", "file2.go"},
		TokensUsed:     2000,
	}

	task := &models.Task{
		ID: "TASK-001",
		Acceptance: []string{
			"Test 1",
			"Test 2",
			"Test 3",
		},
	}

	checkpoint := CreateFromResult("TASK-001", 3, "claude-code", result, task)

	if checkpoint.TaskID != "TASK-001" {
		t.Errorf("Expected TaskID TASK-001, got %s", checkpoint.TaskID)
	}

	if checkpoint.Iteration != 3 {
		t.Errorf("Expected iteration 3, got %d", checkpoint.Iteration)
	}

	if checkpoint.TokensUsed != 2000 {
		t.Errorf("Expected 2000 tokens, got %d", checkpoint.TokensUsed)
	}

	if len(checkpoint.CriteriaMet) != 1 {
		t.Errorf("Expected 1 criteria met, got %d", len(checkpoint.CriteriaMet))
	}

	if len(checkpoint.CriteriaLeft) != 2 {
		t.Errorf("Expected 2 criteria left, got %d", len(checkpoint.CriteriaLeft))
	}
}

func TestCheckpointManager_FileStructure(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	checkpoint := &Checkpoint{
		TaskID:    "TASK-001",
		Iteration: 5,
		CreatedAt: time.Now(),
	}

	if err := manager.Save(checkpoint); err != nil {
		t.Fatalf("Failed to save checkpoint: %v", err)
	}

	// Check iteration file exists
	iterPath := filepath.Join(tmpDir, "TASK-001-005.json")
	if _, err := os.Stat(iterPath); os.IsNotExist(err) {
		t.Errorf("Iteration checkpoint file not created: %s", iterPath)
	}

	// Check latest file exists
	latestPath := filepath.Join(tmpDir, "TASK-001-latest.json")
	if _, err := os.Stat(latestPath); os.IsNotExist(err) {
		t.Errorf("Latest checkpoint file not created: %s", latestPath)
	}
}
