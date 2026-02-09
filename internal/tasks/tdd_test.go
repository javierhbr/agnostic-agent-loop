package tasks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
	"gopkg.in/yaml.v3"
)

func TestDecomposeForTDD(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a task in in-progress
	task := models.Task{
		ID:     "TASK-100",
		Title:  "Implement user login",
		Status: models.StatusInProgress,
	}

	list := &TaskList{Tasks: []models.Task{task}}
	data, err := yaml.Marshal(list)
	if err != nil {
		t.Fatalf("failed to marshal task list: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "in-progress.yaml"), data, 0644); err != nil {
		t.Fatalf("failed to write in-progress.yaml: %v", err)
	}

	// Create empty backlog and done
	emptyList, _ := yaml.Marshal(&TaskList{Tasks: []models.Task{}})
	os.WriteFile(filepath.Join(tmpDir, "backlog.yaml"), emptyList, 0644)
	os.WriteFile(filepath.Join(tmpDir, "done.yaml"), emptyList, 0644)

	tm := NewTaskManager(tmpDir)

	subtasks, err := DecomposeForTDD(tm, "TASK-100")
	if err != nil {
		t.Fatalf("DecomposeForTDD failed: %v", err)
	}

	if len(subtasks) != 3 {
		t.Fatalf("expected 3 sub-tasks, got %d", len(subtasks))
	}

	// Verify RED phase
	if subtasks[0].ID != "TASK-100-red" {
		t.Errorf("expected sub-task ID 'TASK-100-red', got %q", subtasks[0].ID)
	}
	if subtasks[0].Status != models.StatusPending {
		t.Errorf("expected pending status, got %q", subtasks[0].Status)
	}

	// Verify GREEN phase
	if subtasks[1].ID != "TASK-100-green" {
		t.Errorf("expected sub-task ID 'TASK-100-green', got %q", subtasks[1].ID)
	}

	// Verify REFACTOR phase
	if subtasks[2].ID != "TASK-100-refactor" {
		t.Errorf("expected sub-task ID 'TASK-100-refactor', got %q", subtasks[2].ID)
	}

	// Verify sub-tasks were saved to the parent task
	updatedTask, _, err := tm.FindTask("TASK-100")
	if err != nil {
		t.Fatalf("failed to find updated task: %v", err)
	}
	if len(updatedTask.SubTasks) != 3 {
		t.Errorf("expected parent task to have 3 sub-tasks, got %d", len(updatedTask.SubTasks))
	}
}

func TestDecomposeForTDD_TaskNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	emptyList, _ := yaml.Marshal(&TaskList{Tasks: []models.Task{}})
	os.WriteFile(filepath.Join(tmpDir, "backlog.yaml"), emptyList, 0644)
	os.WriteFile(filepath.Join(tmpDir, "in-progress.yaml"), emptyList, 0644)
	os.WriteFile(filepath.Join(tmpDir, "done.yaml"), emptyList, 0644)

	tm := NewTaskManager(tmpDir)

	_, err := DecomposeForTDD(tm, "NONEXISTENT")
	if err == nil {
		t.Fatal("expected error for nonexistent task")
	}
}
