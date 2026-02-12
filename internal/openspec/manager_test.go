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
	// specs/ dir is NOT created upfront — agents create it on demand (Phase 3)

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

	// Check all task IDs are unique
	seenIDs := make(map[string]bool)
	for _, task := range created {
		if seenIDs[task.ID] {
			t.Errorf("duplicate task ID: %s", task.ID)
		}
		seenIDs[task.ID] = true

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

	// Re-import should fail (idempotency guard)
	_, err = m.Import(change.ID, tm)
	if err == nil {
		t.Error("expected error on re-import")
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

func TestSync(t *testing.T) {
	dir := t.TempDir()
	changesDir := filepath.Join(dir, "changes")
	taskDir := filepath.Join(dir, "tasks")
	os.MkdirAll(taskDir, 0755)

	m := NewManager(changesDir)
	tm := tasks.NewTaskManager(taskDir)

	// Init a change
	reqFile := filepath.Join(dir, "req.md")
	os.WriteFile(reqFile, []byte("requirements"), 0644)
	change, err := m.Init("Sync Test", reqFile)
	if err != nil {
		t.Fatal(err)
	}

	// Write tasks.md with real tasks
	tasksPath := filepath.Join(changesDir, change.ID, "tasks.md")
	os.WriteFile(tasksPath, []byte("1. Create models\n2. Add routes\n3. Write tests\n"), 0644)

	// Sync should auto-import
	result, err := m.Sync(tm)
	if err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	if len(result.ChangesImported) != 1 {
		t.Errorf("expected 1 change imported, got %d", len(result.ChangesImported))
	}
	if result.TasksCreated != 3 {
		t.Errorf("expected 3 tasks created, got %d", result.TasksCreated)
	}

	// Verify tasks are in backlog
	backlog, err := tm.LoadTasks("backlog")
	if err != nil {
		t.Fatal(err)
	}
	if len(backlog.Tasks) != 3 {
		t.Errorf("expected 3 tasks in backlog, got %d", len(backlog.Tasks))
	}

	// Verify change status updated
	updated, _ := m.Get(change.ID)
	if updated.Status != StatusImported {
		t.Errorf("expected status 'imported', got %q", updated.Status)
	}
}

func TestSyncSkipsEmptyTasks(t *testing.T) {
	dir := t.TempDir()
	changesDir := filepath.Join(dir, "changes")
	taskDir := filepath.Join(dir, "tasks")
	os.MkdirAll(taskDir, 0755)

	m := NewManager(changesDir)
	tm := tasks.NewTaskManager(taskDir)

	// Init a change (tasks.md will be template-only)
	reqFile := filepath.Join(dir, "req.md")
	os.WriteFile(reqFile, []byte("requirements"), 0644)
	_, err := m.Init("Empty Tasks", reqFile)
	if err != nil {
		t.Fatal(err)
	}

	// Sync should be a no-op since tasks.md has no parseable tasks
	result, err := m.Sync(tm)
	if err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	if len(result.ChangesImported) != 0 {
		t.Errorf("expected 0 changes imported, got %d", len(result.ChangesImported))
	}
	if result.TasksCreated != 0 {
		t.Errorf("expected 0 tasks created, got %d", result.TasksCreated)
	}

	// Change should still be draft
	change, _ := m.Get("empty-tasks")
	if change.Status != StatusDraft {
		t.Errorf("expected status 'draft', got %q", change.Status)
	}
}

func TestSyncIdempotent(t *testing.T) {
	dir := t.TempDir()
	changesDir := filepath.Join(dir, "changes")
	taskDir := filepath.Join(dir, "tasks")
	os.MkdirAll(taskDir, 0755)

	m := NewManager(changesDir)
	tm := tasks.NewTaskManager(taskDir)

	// Init and populate tasks
	reqFile := filepath.Join(dir, "req.md")
	os.WriteFile(reqFile, []byte("requirements"), 0644)
	change, _ := m.Init("Idempotent Test", reqFile)

	tasksPath := filepath.Join(changesDir, change.ID, "tasks.md")
	os.WriteFile(tasksPath, []byte("1. Task one\n2. Task two\n"), 0644)

	// First sync imports
	result1, _ := m.Sync(tm)
	if result1.TasksCreated != 2 {
		t.Fatalf("expected 2 tasks on first sync, got %d", result1.TasksCreated)
	}

	// Second sync is a no-op (change is no longer draft)
	result2, err := m.Sync(tm)
	if err != nil {
		t.Fatalf("Second sync failed: %v", err)
	}
	if len(result2.ChangesImported) != 0 {
		t.Errorf("expected 0 changes on second sync, got %d", len(result2.ChangesImported))
	}
	if result2.TasksCreated != 0 {
		t.Errorf("expected 0 tasks on second sync, got %d", result2.TasksCreated)
	}

	// Backlog should still have exactly 2 tasks
	backlog, _ := tm.LoadTasks("backlog")
	if len(backlog.Tasks) != 2 {
		t.Errorf("expected 2 tasks in backlog after double sync, got %d", len(backlog.Tasks))
	}
}

func TestSyncMultipleChanges(t *testing.T) {
	dir := t.TempDir()
	changesDir := filepath.Join(dir, "changes")
	taskDir := filepath.Join(dir, "tasks")
	os.MkdirAll(taskDir, 0755)

	m := NewManager(changesDir)
	tm := tasks.NewTaskManager(taskDir)

	reqFile := filepath.Join(dir, "req.md")
	os.WriteFile(reqFile, []byte("requirements"), 0644)

	// Init two changes
	change1, _ := m.Init("Feature Alpha", reqFile)
	change2, _ := m.Init("Feature Beta", reqFile)

	// Populate tasks.md for change1 only
	tasksPath1 := filepath.Join(changesDir, change1.ID, "tasks.md")
	os.WriteFile(tasksPath1, []byte("1. Alpha task one\n2. Alpha task two\n"), 0644)

	// Leave change2 tasks.md as template (no parseable tasks)

	// Sync should only import change1
	result, err := m.Sync(tm)
	if err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	if len(result.ChangesImported) != 1 {
		t.Errorf("expected 1 change imported, got %d", len(result.ChangesImported))
	}
	if result.ChangesImported[0] != change1.ID {
		t.Errorf("expected change %q imported, got %q", change1.ID, result.ChangesImported[0])
	}
	if result.TasksCreated != 2 {
		t.Errorf("expected 2 tasks created, got %d", result.TasksCreated)
	}

	// change1 should be imported, change2 should still be draft
	updated1, _ := m.Get(change1.ID)
	if updated1.Status != StatusImported {
		t.Errorf("change1: expected 'imported', got %q", updated1.Status)
	}
	updated2, _ := m.Get(change2.ID)
	if updated2.Status != StatusDraft {
		t.Errorf("change2: expected 'draft', got %q", updated2.Status)
	}
}

func TestSyncNoRegistry(t *testing.T) {
	dir := t.TempDir()
	changesDir := filepath.Join(dir, "no-such-dir")
	taskDir := filepath.Join(dir, "tasks")
	os.MkdirAll(taskDir, 0755)

	m := NewManager(changesDir)
	tm := tasks.NewTaskManager(taskDir)

	// Sync on non-existent dir should not error
	result, err := m.Sync(tm)
	if err != nil {
		t.Fatalf("Sync should not fail on missing dir: %v", err)
	}
	if len(result.ChangesImported) != 0 {
		t.Errorf("expected 0 changes, got %d", len(result.ChangesImported))
	}
}

func TestEnsureConfig(t *testing.T) {
	dir := t.TempDir()

	result, err := EnsureConfig(dir, "my-project")
	if err != nil {
		t.Fatalf("EnsureConfig failed: %v", err)
	}
	if !result.Created {
		t.Error("expected config to be created")
	}

	// Verify file exists and has correct content
	data, err := os.ReadFile(filepath.Join(dir, "agnostic-agent.yaml"))
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	content := string(data)

	// Check key fields
	if !strings.Contains(content, "name: my-project") {
		t.Error("expected project name in config")
	}
	if !strings.Contains(content, "openSpecDir: .agentic/openspec/changes") {
		t.Error("expected openSpecDir in config")
	}
	if !strings.Contains(content, ".specify/specs") {
		t.Error("expected .specify/specs in specDirs")
	}
	if !strings.Contains(content, "openspec/specs") {
		t.Error("expected openspec/specs in specDirs")
	}
	if !strings.Contains(content, ".agentic/spec") {
		t.Error("expected .agentic/spec in specDirs")
	}
	if !strings.Contains(content, "context-check") {
		t.Error("expected context-check validator")
	}
	if !strings.Contains(content, "task-scope") {
		t.Error("expected task-scope validator")
	}
}

func TestEnsureConfigSkipsExisting(t *testing.T) {
	dir := t.TempDir()

	// Create an existing config
	existingContent := "project:\n  name: existing\n"
	os.WriteFile(filepath.Join(dir, "agnostic-agent.yaml"), []byte(existingContent), 0644)

	result, err := EnsureConfig(dir, "new-project")
	if err != nil {
		t.Fatalf("EnsureConfig failed: %v", err)
	}
	if result.Created {
		t.Error("expected config NOT to be created when one already exists")
	}

	// Verify original content is preserved
	data, _ := os.ReadFile(filepath.Join(dir, "agnostic-agent.yaml"))
	if string(data) != existingContent {
		t.Error("expected existing config to be preserved")
	}
}

func TestEnsureConfigDefaultsProjectName(t *testing.T) {
	dir := t.TempDir()

	result, err := EnsureConfig(dir, "")
	if err != nil {
		t.Fatalf("EnsureConfig failed: %v", err)
	}
	if !result.Created {
		t.Error("expected config to be created")
	}

	// Should use directory basename as project name
	data, _ := os.ReadFile(filepath.Join(dir, "agnostic-agent.yaml"))
	dirName := filepath.Base(dir)
	if !strings.Contains(string(data), "name: "+dirName) {
		t.Errorf("expected project name %q from dir basename, got:\n%s", dirName, string(data))
	}
}

// --- Structured parser tests ---

func TestParserStructured(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name     string
		content  string
		expected []TaskEntry
	}{
		{
			name:    "numbered list without refs",
			content: "1. First task\n2. Second task\n",
			expected: []TaskEntry{
				{Title: "First task", FileRef: ""},
				{Title: "Second task", FileRef: ""},
			},
		},
		{
			name:    "checkbox list with ver refs",
			content: "- [ ] Setup (ver [tasks/01-setup.md](./tasks/01-setup.md))\n- [ ] Build (ver [tasks/02-build.md](./tasks/02-build.md))\n",
			expected: []TaskEntry{
				{Title: "Setup", FileRef: "tasks/01-setup.md"},
				{Title: "Build", FileRef: "tasks/02-build.md"},
			},
		},
		{
			name:    "mixed - some with refs some without",
			content: "- [ ] Has ref (ver [tasks/01-foo.md](./tasks/01-foo.md))\n- [ ] No ref here\n",
			expected: []TaskEntry{
				{Title: "Has ref", FileRef: "tasks/01-foo.md"},
				{Title: "No ref here", FileRef: ""},
			},
		},
		{
			name:    "numbered list with refs",
			content: "1. Setup project (ver [tasks/01-setup.md](tasks/01-setup.md))\n2. Add API\n",
			expected: []TaskEntry{
				{Title: "Setup project", FileRef: "tasks/01-setup.md"},
				{Title: "Add API", FileRef: ""},
			},
		},
		{
			name:    "headers and intro text ignored",
			content: "# Tasks\n\nSome intro.\n\n- [ ] Real task (ver [tasks/01-real.md](./tasks/01-real.md))\n",
			expected: []TaskEntry{
				{Title: "Real task", FileRef: "tasks/01-real.md"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(dir, tt.name+".md")
			if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			got, err := ParseTasksFileStructured(path)
			if err != nil {
				t.Fatalf("ParseTasksFileStructured failed: %v", err)
			}

			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d entries, got %d: %v", len(tt.expected), len(got), got)
			}
			for i, entry := range got {
				if entry.Title != tt.expected[i].Title {
					t.Errorf("entry %d title: expected %q, got %q", i, tt.expected[i].Title, entry.Title)
				}
				if entry.FileRef != tt.expected[i].FileRef {
					t.Errorf("entry %d fileref: expected %q, got %q", i, tt.expected[i].FileRef, entry.FileRef)
				}
			}
		})
	}
}

func TestParseTaskDetailFile(t *testing.T) {
	dir := t.TempDir()

	t.Run("full file with all sections", func(t *testing.T) {
		content := `# Set up project structure

## Description
Create the initial Vite/React project with skeleton UI components.

## Prerequisites
- None (first task)

## Acceptance Criteria
- [ ] Package.json created with all dependencies
- [ ] Dev server runs without errors
- [x] README exists

## Technical Notes
Use Vite for instant HMR. React 18 with Suspense.
`
		path := filepath.Join(dir, "full.md")
		os.WriteFile(path, []byte(content), 0644)

		detail, err := ParseTaskDetailFile(path)
		if err != nil {
			t.Fatalf("ParseTaskDetailFile failed: %v", err)
		}

		if detail.Title != "Set up project structure" {
			t.Errorf("title: expected %q, got %q", "Set up project structure", detail.Title)
		}
		if !strings.Contains(detail.Description, "Vite/React") {
			t.Errorf("description should contain 'Vite/React', got %q", detail.Description)
		}
		if len(detail.Prerequisites) != 1 || detail.Prerequisites[0] != "None (first task)" {
			t.Errorf("prerequisites: expected [None (first task)], got %v", detail.Prerequisites)
		}
		if len(detail.Acceptance) != 3 {
			t.Errorf("acceptance: expected 3 items, got %d: %v", len(detail.Acceptance), detail.Acceptance)
		}
		if detail.Acceptance[0] != "Package.json created with all dependencies" {
			t.Errorf("acceptance[0]: expected 'Package.json created with all dependencies', got %q", detail.Acceptance[0])
		}
		if !strings.Contains(detail.Notes, "Vite for instant HMR") {
			t.Errorf("notes should contain 'Vite for instant HMR', got %q", detail.Notes)
		}
	})

	t.Run("empty sections", func(t *testing.T) {
		content := `# Minimal task

## Description

## Acceptance Criteria
`
		path := filepath.Join(dir, "empty-sections.md")
		os.WriteFile(path, []byte(content), 0644)

		detail, err := ParseTaskDetailFile(path)
		if err != nil {
			t.Fatalf("ParseTaskDetailFile failed: %v", err)
		}

		if detail.Title != "Minimal task" {
			t.Errorf("title: expected 'Minimal task', got %q", detail.Title)
		}
		if detail.Description != "" {
			t.Errorf("description should be empty, got %q", detail.Description)
		}
		if len(detail.Acceptance) != 0 {
			t.Errorf("acceptance should be empty, got %v", detail.Acceptance)
		}
	})

	t.Run("unknown sections ignored", func(t *testing.T) {
		content := `# Task with extras

## Description
The real description.

## Random Section
This should be ignored.

## Acceptance Criteria
- [ ] It works
`
		path := filepath.Join(dir, "unknown-sections.md")
		os.WriteFile(path, []byte(content), 0644)

		detail, err := ParseTaskDetailFile(path)
		if err != nil {
			t.Fatalf("ParseTaskDetailFile failed: %v", err)
		}

		if !strings.Contains(detail.Description, "real description") {
			t.Errorf("description: expected 'real description', got %q", detail.Description)
		}
		if len(detail.Acceptance) != 1 || detail.Acceptance[0] != "It works" {
			t.Errorf("acceptance: expected [It works], got %v", detail.Acceptance)
		}
	})

	t.Run("bullet variants", func(t *testing.T) {
		content := `# Bullet test

## Prerequisites
- Dash item
* Star item
`
		path := filepath.Join(dir, "bullets.md")
		os.WriteFile(path, []byte(content), 0644)

		detail, err := ParseTaskDetailFile(path)
		if err != nil {
			t.Fatalf("ParseTaskDetailFile failed: %v", err)
		}

		if len(detail.Prerequisites) != 2 {
			t.Fatalf("prerequisites: expected 2, got %d: %v", len(detail.Prerequisites), detail.Prerequisites)
		}
		if detail.Prerequisites[0] != "Dash item" {
			t.Errorf("prerequisites[0]: expected 'Dash item', got %q", detail.Prerequisites[0])
		}
		if detail.Prerequisites[1] != "Star item" {
			t.Errorf("prerequisites[1]: expected 'Star item', got %q", detail.Prerequisites[1])
		}
	})
}

func TestHasTasksDir(t *testing.T) {
	dir := t.TempDir()
	tasksFile := filepath.Join(dir, "tasks.md")
	os.WriteFile(tasksFile, []byte("1. task\n"), 0644)

	// No tasks/ dir yet
	if HasTasksDir(tasksFile) {
		t.Error("expected HasTasksDir=false when tasks/ doesn't exist")
	}

	// Create tasks/ dir
	os.MkdirAll(filepath.Join(dir, "tasks"), 0755)
	if !HasTasksDir(tasksFile) {
		t.Error("expected HasTasksDir=true when tasks/ exists")
	}
}

// --- Import with task detail files ---

func TestImportWithTasksDir(t *testing.T) {
	dir := t.TempDir()
	changesDir := filepath.Join(dir, "changes")
	taskDir := filepath.Join(dir, "tasks")
	os.MkdirAll(taskDir, 0755)

	m := NewManager(changesDir)
	tm := tasks.NewTaskManager(taskDir)

	// Init change
	reqFile := filepath.Join(dir, "req.md")
	os.WriteFile(reqFile, []byte("requirements"), 0644)
	change, err := m.Init("Detail Test", reqFile)
	if err != nil {
		t.Fatal(err)
	}

	changeDir := filepath.Join(changesDir, change.ID)

	// Create tasks/ subdirectory with a detail file
	tasksSubDir := filepath.Join(changeDir, "tasks")
	os.MkdirAll(tasksSubDir, 0755)

	os.WriteFile(filepath.Join(tasksSubDir, "01-setup.md"), []byte(`# Project Setup

## Description
Create the initial project structure with all dependencies.

## Prerequisites
- Node.js installed

## Acceptance Criteria
- [ ] Package.json exists
- [ ] Dev server runs

## Technical Notes
Use Vite for HMR.
`), 0644)

	// Write tasks.md index with reference
	tasksPath := filepath.Join(changeDir, "tasks.md")
	os.WriteFile(tasksPath, []byte(
		"- [ ] Project Setup (ver [tasks/01-setup.md](./tasks/01-setup.md))\n",
	), 0644)

	created, err := m.Import(change.ID, tm)
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	if len(created) != 1 {
		t.Fatalf("expected 1 task, got %d", len(created))
	}

	task := created[0]

	// Check title
	if !strings.Contains(task.Title, "Project Setup") {
		t.Errorf("title should contain 'Project Setup', got %q", task.Title)
	}

	// Check description populated
	if !strings.Contains(task.Description, "initial project structure") {
		t.Errorf("description should contain 'initial project structure', got %q", task.Description)
	}

	// Check technical notes appended to description
	if !strings.Contains(task.Description, "Vite for HMR") {
		t.Errorf("description should contain technical notes, got %q", task.Description)
	}

	// Check acceptance criteria
	if len(task.Acceptance) != 2 {
		t.Errorf("expected 2 acceptance criteria, got %d: %v", len(task.Acceptance), task.Acceptance)
	}

	// Check inputs (prerequisites)
	if len(task.Inputs) != 1 || task.Inputs[0] != "Node.js installed" {
		t.Errorf("expected inputs [Node.js installed], got %v", task.Inputs)
	}

	// Check spec refs include both proposal and detail file
	if len(task.SpecRefs) != 2 {
		t.Errorf("expected 2 spec refs, got %d: %v", len(task.SpecRefs), task.SpecRefs)
	}
}

func TestImportMixedWithAndWithoutDetailFiles(t *testing.T) {
	dir := t.TempDir()
	changesDir := filepath.Join(dir, "changes")
	taskDir := filepath.Join(dir, "tasks")
	os.MkdirAll(taskDir, 0755)

	m := NewManager(changesDir)
	tm := tasks.NewTaskManager(taskDir)

	reqFile := filepath.Join(dir, "req.md")
	os.WriteFile(reqFile, []byte("requirements"), 0644)
	change, _ := m.Init("Mixed Test", reqFile)

	changeDir := filepath.Join(changesDir, change.ID)
	tasksSubDir := filepath.Join(changeDir, "tasks")
	os.MkdirAll(tasksSubDir, 0755)

	// Only create detail file for first task
	os.WriteFile(filepath.Join(tasksSubDir, "01-detailed.md"), []byte(`# Detailed task

## Description
This task has full details.

## Acceptance Criteria
- [ ] Everything works
`), 0644)

	// tasks.md: first has ref, second doesn't
	tasksPath := filepath.Join(changeDir, "tasks.md")
	os.WriteFile(tasksPath, []byte(
		"- [ ] Detailed task (ver [tasks/01-detailed.md](./tasks/01-detailed.md))\n"+
			"- [ ] Simple task without details\n",
	), 0644)

	created, err := m.Import(change.ID, tm)
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	if len(created) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(created))
	}

	// First task has details
	if created[0].Description == "" {
		t.Error("first task should have description populated")
	}
	if len(created[0].Acceptance) != 1 {
		t.Errorf("first task should have 1 acceptance criterion, got %d", len(created[0].Acceptance))
	}

	// Second task has no details (title only)
	if created[1].Description != "" {
		t.Errorf("second task should have empty description, got %q", created[1].Description)
	}
	if len(created[1].Acceptance) != 0 {
		t.Errorf("second task should have no acceptance criteria, got %v", created[1].Acceptance)
	}
}

func TestImportWithMissingDetailFile(t *testing.T) {
	dir := t.TempDir()
	changesDir := filepath.Join(dir, "changes")
	taskDir := filepath.Join(dir, "tasks")
	os.MkdirAll(taskDir, 0755)

	m := NewManager(changesDir)
	tm := tasks.NewTaskManager(taskDir)

	reqFile := filepath.Join(dir, "req.md")
	os.WriteFile(reqFile, []byte("requirements"), 0644)
	change, _ := m.Init("Missing File Test", reqFile)

	changeDir := filepath.Join(changesDir, change.ID)
	os.MkdirAll(filepath.Join(changeDir, "tasks"), 0755)

	// tasks.md references a file that doesn't exist
	tasksPath := filepath.Join(changeDir, "tasks.md")
	os.WriteFile(tasksPath, []byte(
		"- [ ] Ghost task (ver [tasks/01-ghost.md](./tasks/01-ghost.md))\n",
	), 0644)

	// Import should succeed with a warning, not fail
	created, err := m.Import(change.ID, tm)
	if err != nil {
		t.Fatalf("Import should not fail on missing detail file: %v", err)
	}

	if len(created) != 1 {
		t.Fatalf("expected 1 task, got %d", len(created))
	}

	// Task created with title only
	if !strings.Contains(created[0].Title, "Ghost task") {
		t.Errorf("title should contain 'Ghost task', got %q", created[0].Title)
	}
	if created[0].Description != "" {
		t.Errorf("description should be empty, got %q", created[0].Description)
	}
	// Only proposal.md in spec refs (detail file not added since it failed)
	if len(created[0].SpecRefs) != 1 {
		t.Errorf("expected 1 spec ref (proposal only), got %d: %v", len(created[0].SpecRefs), created[0].SpecRefs)
	}
}

func TestImportWithTasksDirButNoReferences(t *testing.T) {
	dir := t.TempDir()
	changesDir := filepath.Join(dir, "changes")
	taskDir := filepath.Join(dir, "tasks")
	os.MkdirAll(taskDir, 0755)

	m := NewManager(changesDir)
	tm := tasks.NewTaskManager(taskDir)

	reqFile := filepath.Join(dir, "req.md")
	os.WriteFile(reqFile, []byte("requirements"), 0644)
	change, _ := m.Init("No Refs Test", reqFile)

	changeDir := filepath.Join(changesDir, change.ID)
	os.MkdirAll(filepath.Join(changeDir, "tasks"), 0755)

	// tasks.md has simple numbered list (no references)
	tasksPath := filepath.Join(changeDir, "tasks.md")
	os.WriteFile(tasksPath, []byte("1. Simple task one\n2. Simple task two\n"), 0644)

	created, err := m.Import(change.ID, tm)
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	if len(created) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(created))
	}

	// Tasks should be title-only (no detail parsing)
	for i, task := range created {
		if task.Description != "" {
			t.Errorf("task %d should have empty description, got %q", i, task.Description)
		}
		if len(task.SpecRefs) != 1 {
			t.Errorf("task %d should have 1 spec ref, got %d", i, len(task.SpecRefs))
		}
	}
}

func TestScaffoldTaskFiles(t *testing.T) {
	dir := t.TempDir()
	changesDir := filepath.Join(dir, "changes")
	m := NewManager(changesDir)

	reqFile := filepath.Join(dir, "req.md")
	os.WriteFile(reqFile, []byte("requirements"), 0644)
	change, err := m.Init("Scaffold Test", reqFile)
	if err != nil {
		t.Fatal(err)
	}

	titles := []string{"Create user model", "Add API endpoints", "Write tests"}
	if err := m.ScaffoldTaskFiles(change.ID, titles); err != nil {
		t.Fatalf("ScaffoldTaskFiles failed: %v", err)
	}

	// Verify files exist
	tasksSubDir := filepath.Join(changesDir, change.ID, "tasks")
	expectedFiles := []string{
		"01-create-user-model.md",
		"02-add-api-endpoints.md",
		"03-write-tests.md",
	}
	for _, f := range expectedFiles {
		path := filepath.Join(tasksSubDir, f)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("expected %s to exist: %v", f, err)
			continue
		}
		content := string(data)
		if !strings.Contains(content, "## Description") {
			t.Errorf("%s should contain '## Description'", f)
		}
		if !strings.Contains(content, "## Acceptance Criteria") {
			t.Errorf("%s should contain '## Acceptance Criteria'", f)
		}
	}

	// First file should have the title in H1
	data, _ := os.ReadFile(filepath.Join(tasksSubDir, "01-create-user-model.md"))
	if !strings.Contains(string(data), "# Create user model") {
		t.Error("first file should contain '# Create user model'")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && strings.Contains(s, substr)
}
