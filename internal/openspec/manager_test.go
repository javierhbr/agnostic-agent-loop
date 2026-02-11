package openspec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/javierbenavides/agentic-agent/internal/tasks"
)

func TestInit(t *testing.T) {
	dir := t.TempDir()
	changesDir := filepath.Join(dir, "changes")
	m := NewManager(changesDir)

	// Create a source requirements file
	reqFile := filepath.Join(dir, "requirements.md")
	if err := os.WriteFile(reqFile, []byte("# Auth\n- Login\n- Register\n"), 0644); err != nil {
		t.Fatal(err)
	}

	change, err := m.Init("Auth Feature", reqFile)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if change.ID != "auth-feature" {
		t.Errorf("expected id 'auth-feature', got %q", change.ID)
	}
	if change.Status != StatusDraft {
		t.Errorf("expected status 'draft', got %q", change.Status)
	}

	// Check files were created
	changeDir := filepath.Join(changesDir, "auth-feature")
	for _, f := range []string{"proposal.md", "tasks.md", "metadata.yaml"} {
		if _, err := os.Stat(filepath.Join(changeDir, f)); err != nil {
			t.Errorf("expected %s to exist: %v", f, err)
		}
	}
	// Check specs/ dir
	if _, err := os.Stat(filepath.Join(changeDir, "specs")); err != nil {
		t.Error("expected specs/ directory to exist")
	}

	// Check proposal contains requirements
	proposal, err := os.ReadFile(filepath.Join(changeDir, "proposal.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !contains(string(proposal), "Login") {
		t.Error("expected proposal to contain requirements content")
	}

	// Check registry
	reg, err := m.loadRegistry()
	if err != nil {
		t.Fatal(err)
	}
	if len(reg.Changes) != 1 {
		t.Errorf("expected 1 change in registry, got %d", len(reg.Changes))
	}
}

func TestInitDuplicate(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir)

	reqFile := filepath.Join(dir, "req.md")
	os.WriteFile(reqFile, []byte("test"), 0644)

	if _, err := m.Init("My Feature", reqFile); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Init("My Feature", reqFile); err == nil {
		t.Error("expected error for duplicate change")
	}
}

func TestImport(t *testing.T) {
	dir := t.TempDir()
	changesDir := filepath.Join(dir, "changes")
	taskDir := filepath.Join(dir, "tasks")
	os.MkdirAll(taskDir, 0755)

	m := NewManager(changesDir)
	tm := tasks.NewTaskManager(taskDir)

	// Create a source file and init the change
	reqFile := filepath.Join(dir, "req.md")
	os.WriteFile(reqFile, []byte("requirements"), 0644)

	change, err := m.Init("Import Test", reqFile)
	if err != nil {
		t.Fatal(err)
	}

	// Write tasks.md with numbered items
	tasksPath := filepath.Join(changesDir, change.ID, "tasks.md")
	tasksContent := "# Tasks\n\n1. Create user model\n2. Add API endpoints\n3. Write tests\n"
	if err := os.WriteFile(tasksPath, []byte(tasksContent), 0644); err != nil {
		t.Fatal(err)
	}

	created, err := m.Import(change.ID, tm)
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	if len(created) != 3 {
		t.Errorf("expected 3 tasks, got %d", len(created))
	}

	// Check tasks have ChangeID set
	for _, task := range created {
		if task.ChangeID != change.ID {
			t.Errorf("expected ChangeID %q, got %q", change.ID, task.ChangeID)
		}
	}

	// Check change status updated
	updated, err := m.Get(change.ID)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Status != StatusImported {
		t.Errorf("expected status 'imported', got %q", updated.Status)
	}
	if len(updated.TaskIDs) != 3 {
		t.Errorf("expected 3 task IDs, got %d", len(updated.TaskIDs))
	}
}

func TestComplete(t *testing.T) {
	dir := t.TempDir()
	changesDir := filepath.Join(dir, "changes")
	taskDir := filepath.Join(dir, "tasks")
	os.MkdirAll(taskDir, 0755)

	m := NewManager(changesDir)
	tm := tasks.NewTaskManager(taskDir)

	// Init + import
	reqFile := filepath.Join(dir, "req.md")
	os.WriteFile(reqFile, []byte("requirements"), 0644)
	change, _ := m.Init("Complete Test", reqFile)

	tasksPath := filepath.Join(changesDir, change.ID, "tasks.md")
	os.WriteFile(tasksPath, []byte("1. Single task\n"), 0644)
	created, _ := m.Import(change.ID, tm)

	// Should fail — task not done
	if err := m.Complete(change.ID, tm); err == nil {
		t.Error("expected error when tasks are pending")
	}

	// Move task to done
	for _, task := range created {
		tm.MoveTask(task.ID, "backlog", "done", "done")
	}

	// Should succeed now
	if err := m.Complete(change.ID, tm); err != nil {
		t.Fatalf("Complete failed: %v", err)
	}

	// Check IMPLEMENTED marker
	markerPath := filepath.Join(changesDir, change.ID, "IMPLEMENTED")
	if _, err := os.Stat(markerPath); err != nil {
		t.Error("expected IMPLEMENTED marker to exist")
	}

	// Check status
	updated, _ := m.Get(change.ID)
	if updated.Status != StatusImplemented {
		t.Errorf("expected status 'implemented', got %q", updated.Status)
	}
}

func TestArchive(t *testing.T) {
	dir := t.TempDir()
	changesDir := filepath.Join(dir, "changes")
	taskDir := filepath.Join(dir, "tasks")
	os.MkdirAll(taskDir, 0755)

	m := NewManager(changesDir)
	tm := tasks.NewTaskManager(taskDir)

	// Init + import + complete
	reqFile := filepath.Join(dir, "req.md")
	os.WriteFile(reqFile, []byte("requirements"), 0644)
	change, _ := m.Init("Archive Test", reqFile)

	tasksPath := filepath.Join(changesDir, change.ID, "tasks.md")
	os.WriteFile(tasksPath, []byte("1. Do something\n"), 0644)
	created, _ := m.Import(change.ID, tm)
	for _, task := range created {
		tm.MoveTask(task.ID, "backlog", "done", "done")
	}
	m.Complete(change.ID, tm)

	// Should fail without IMPLEMENTED — but we just completed, so test archive
	if err := m.Archive(change.ID); err != nil {
		t.Fatalf("Archive failed: %v", err)
	}

	// Original dir should be gone
	if _, err := os.Stat(filepath.Join(changesDir, change.ID)); !os.IsNotExist(err) {
		t.Error("expected change directory to be moved")
	}

	// Archive dir should exist
	if _, err := os.Stat(filepath.Join(changesDir, "_archive", change.ID)); err != nil {
		t.Error("expected archived directory to exist")
	}
}

func TestArchiveWithoutImplemented(t *testing.T) {
	dir := t.TempDir()
	m := NewManager(dir)

	reqFile := filepath.Join(dir, "req.md")
	os.WriteFile(reqFile, []byte("test"), 0644)
	change, _ := m.Init("No Impl", reqFile)

	if err := m.Archive(change.ID); err == nil {
		t.Error("expected error archiving without IMPLEMENTED marker")
	}
}

func TestParser(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "numbered list",
			content:  "1. First task\n2. Second task\n3. Third task\n",
			expected: []string{"First task", "Second task", "Third task"},
		},
		{
			name:     "checkbox list",
			content:  "- [ ] Alpha\n- [ ] Beta\n- [x] Gamma\n",
			expected: []string{"Alpha", "Beta", "Gamma"},
		},
		{
			name:     "mixed with headers",
			content:  "# Tasks\n\nSome intro text.\n\n1. Real task one\n2. Real task two\n",
			expected: []string{"Real task one", "Real task two"},
		},
		{
			name:     "in-progress checkbox",
			content:  "- [~] Partially done\n- [ ] Not started\n",
			expected: []string{"Partially done", "Not started"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(dir, tt.name+".md")
			if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			got, err := ParseTasksFile(path)
			if err != nil {
				t.Fatalf("ParseTasksFile failed: %v", err)
			}

			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d tasks, got %d: %v", len(tt.expected), len(got), got)
			}
			for i, title := range got {
				if title != tt.expected[i] {
					t.Errorf("task %d: expected %q, got %q", i, tt.expected[i], title)
				}
			}
		})
	}
}

func TestParserEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.md")
	os.WriteFile(path, []byte("# No tasks here\nJust text.\n"), 0644)

	_, err := ParseTasksFile(path)
	if err == nil {
		t.Error("expected error for file with no tasks")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && strings.Contains(s, substr)
}
