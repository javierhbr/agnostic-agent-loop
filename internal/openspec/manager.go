package openspec

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"gopkg.in/yaml.v3"
)

// ChangeStatus represents the lifecycle status of an openspec change.
type ChangeStatus string

const (
	StatusDraft        ChangeStatus = "draft"
	StatusImported     ChangeStatus = "imported"
	StatusImplementing ChangeStatus = "implementing"
	StatusImplemented  ChangeStatus = "implemented"
	StatusArchived     ChangeStatus = "archived"
)

// Change represents an openspec change proposal.
type Change struct {
	ID         string       `yaml:"id"`
	Name       string       `yaml:"name"`
	Status     ChangeStatus `yaml:"status"`
	SourceFile string       `yaml:"source_file"`
	TaskIDs    []string     `yaml:"task_ids,omitempty"`
	CreatedAt  time.Time    `yaml:"created_at"`
}

// Registry is the index file listing all changes.
type Registry struct {
	Changes []Change `yaml:"changes"`
}

// ChangeProgress holds task completion stats for a change.
type ChangeProgress struct {
	Total      int
	Done       int
	InProgress int
	Pending    int
	TaskIDs    []string
}

// Manager handles openspec change CRUD operations.
type Manager struct {
	baseDir string
}

// NewManager creates a change manager for the given base directory.
func NewManager(baseDir string) *Manager {
	return &Manager{baseDir: baseDir}
}

// Init creates a new change directory with proposal.md and tasks.md templates.
// If fromFile is provided, its content is seeded into the proposal.
func (m *Manager) Init(name, fromFile string) (*Change, error) {
	id := toKebabCase(name)
	changeDir := filepath.Join(m.baseDir, id)

	if _, err := os.Stat(changeDir); err == nil {
		return nil, fmt.Errorf("change %q already exists", id)
	}

	if err := os.MkdirAll(changeDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create change directory: %w", err)
	}

	// Read source requirements file
	var requirements string
	if fromFile != "" {
		data, err := os.ReadFile(fromFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read source file %s: %w", fromFile, err)
		}
		requirements = string(data)
	}

	change := Change{
		ID:         id,
		Name:       name,
		Status:     StatusDraft,
		SourceFile: fromFile,
		CreatedAt:  time.Now(),
	}

	// Render and write proposal.md
	tmplData := TemplateData{
		Name:         name,
		SourceFile:   fromFile,
		Requirements: requirements,
	}

	proposalContent, err := renderTemplate("proposal.md.tmpl", tmplData)
	if err != nil {
		return nil, fmt.Errorf("failed to render proposal.md: %w", err)
	}
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write proposal.md: %w", err)
	}

	// Render and write tasks.md
	tasksContent, err := renderTemplate("tasks.md.tmpl", tmplData)
	if err != nil {
		return nil, fmt.Errorf("failed to render tasks.md: %w", err)
	}
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasksContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write tasks.md: %w", err)
	}

	// Write metadata.yaml
	if err := m.writeMetadata(changeDir, &change); err != nil {
		return nil, err
	}

	// Update registry
	if err := m.addToRegistry(&change); err != nil {
		return nil, err
	}

	return &change, nil
}

// Import reads tasks.md for the given change and creates tasks in the backlog.
// When a tasks/ subdirectory exists alongside tasks.md, individual task detail
// files are parsed to populate Description, Acceptance, and Inputs fields.
func (m *Manager) Import(id string, tm *tasks.TaskManager) ([]*models.Task, error) {
	change, err := m.Get(id)
	if err != nil {
		return nil, err
	}

	if change.Status != StatusDraft {
		return nil, fmt.Errorf("change %q has already been imported (status: %s)", id, change.Status)
	}

	tasksPath := filepath.Join(m.baseDir, id, "tasks.md")
	entries, err := ParseTasksFileStructured(tasksPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tasks: %w", err)
	}

	changeDir := filepath.Join(m.baseDir, id)
	hasDir := HasTasksDir(tasksPath)

	var created []*models.Task
	for _, entry := range entries {
		task, err := tm.CreateTask(fmt.Sprintf("[%s] %s", id, entry.Title))
		if err != nil {
			return created, fmt.Errorf("failed to create task %q: %w", entry.Title, err)
		}

		task.ChangeID = id
		task.SpecRefs = []string{filepath.Join(id, "proposal.md")}

		// If tasks/ dir exists and this entry references a detail file, parse it
		if hasDir && entry.FileRef != "" {
			detailPath := filepath.Join(changeDir, entry.FileRef)
			detail, detailErr := ParseTaskDetailFile(detailPath)
			if detailErr != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not parse task detail %s: %v\n", entry.FileRef, detailErr)
			} else {
				applyTaskDetail(task, detail)
				task.SpecRefs = append(task.SpecRefs, filepath.Join(id, entry.FileRef))
			}
		}

		// Save updated task back to backlog
		backlog, err := tm.LoadTasks("backlog")
		if err != nil {
			return created, fmt.Errorf("failed to load backlog: %w", err)
		}
		for i, t := range backlog.Tasks {
			if t.ID == task.ID {
				backlog.Tasks[i] = *task
				break
			}
		}
		if err := tm.SaveTasks("backlog", backlog); err != nil {
			return created, fmt.Errorf("failed to save task: %w", err)
		}

		change.TaskIDs = append(change.TaskIDs, task.ID)
		created = append(created, task)
	}

	// Update change status and task list
	change.Status = StatusImported
	if err := m.updateInRegistry(change); err != nil {
		return created, err
	}
	if err := m.writeMetadata(changeDir, change); err != nil {
		return created, err
	}

	return created, nil
}

// applyTaskDetail maps parsed TaskDetail fields onto a Task model.
func applyTaskDetail(task *models.Task, detail *TaskDetail) {
	if detail.Description != "" {
		task.Description = detail.Description
	}
	if len(detail.Acceptance) > 0 {
		task.Acceptance = detail.Acceptance
	}
	if len(detail.Prerequisites) > 0 {
		task.Inputs = detail.Prerequisites
	}
	if detail.Notes != "" {
		if task.Description != "" {
			task.Description += "\n\n---\nTechnical Notes:\n" + detail.Notes
		} else {
			task.Description = "Technical Notes:\n" + detail.Notes
		}
	}
}

// ScaffoldTaskFiles creates the tasks/ subdirectory and renders a task-detail
// template for each title. Called by agents during Phase 2 for complex changes.
func (m *Manager) ScaffoldTaskFiles(id string, titles []string) error {
	changeDir := filepath.Join(m.baseDir, id)
	tasksSubDir := filepath.Join(changeDir, "tasks")
	if err := os.MkdirAll(tasksSubDir, 0755); err != nil {
		return fmt.Errorf("failed to create tasks/ directory: %w", err)
	}

	for i, title := range titles {
		filename := fmt.Sprintf("%02d-%s.md", i+1, toKebabCase(title))
		tmplData := TemplateData{TaskTitle: title}
		content, err := renderTemplate("task-detail.md.tmpl", tmplData)
		if err != nil {
			return fmt.Errorf("failed to render task template for %q: %w", title, err)
		}
		filePath := filepath.Join(tasksSubDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filePath, err)
		}
	}
	return nil
}

// Get returns a single change by ID.
func (m *Manager) Get(id string) (*Change, error) {
	reg, err := m.loadRegistry()
	if err != nil {
		return nil, err
	}
	for _, c := range reg.Changes {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("change %q not found", id)
}

// List returns all changes from the registry.
func (m *Manager) List() ([]Change, error) {
	reg, err := m.loadRegistry()
	if err != nil {
		return nil, err
	}
	return reg.Changes, nil
}

// Progress returns task completion stats for a change.
func (m *Manager) Progress(id string, tm *tasks.TaskManager) (*ChangeProgress, error) {
	change, err := m.Get(id)
	if err != nil {
		return nil, err
	}

	progress := &ChangeProgress{
		Total:   len(change.TaskIDs),
		TaskIDs: change.TaskIDs,
	}

	for _, listType := range []string{"backlog", "in-progress", "done"} {
		list, err := tm.LoadTasks(listType)
		if err != nil {
			continue
		}
		for _, t := range list.Tasks {
			if t.ChangeID != id {
				continue
			}
			switch listType {
			case "done":
				progress.Done++
			case "in-progress":
				progress.InProgress++
			default:
				progress.Pending++
			}
		}
	}

	return progress, nil
}

// Complete validates all tasks are done and writes the IMPLEMENTED marker.
func (m *Manager) Complete(id string, tm *tasks.TaskManager) error {
	progress, err := m.Progress(id, tm)
	if err != nil {
		return err
	}

	if progress.InProgress > 0 {
		return fmt.Errorf("change %q has %d task(s) still in progress", id, progress.InProgress)
	}
	if progress.Pending > 0 {
		return fmt.Errorf("change %q has %d task(s) still pending", id, progress.Pending)
	}
	if progress.Done == 0 {
		return fmt.Errorf("change %q has no completed tasks", id)
	}

	// Write IMPLEMENTED marker
	markerPath := filepath.Join(m.baseDir, id, "IMPLEMENTED")
	content := fmt.Sprintf("Implementation completed: %s\nTasks completed: %d\n", time.Now().Format(time.RFC3339), progress.Done)
	if err := os.WriteFile(markerPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write IMPLEMENTED marker: %w", err)
	}

	// Update status
	change, _ := m.Get(id)
	change.Status = StatusImplemented
	if err := m.updateInRegistry(change); err != nil {
		return err
	}
	changeDir := filepath.Join(m.baseDir, id)
	return m.writeMetadata(changeDir, change)
}

// Archive moves a completed change to the _archive directory.
func (m *Manager) Archive(id string) error {
	changeDir := filepath.Join(m.baseDir, id)

	// Require IMPLEMENTED marker
	markerPath := filepath.Join(changeDir, "IMPLEMENTED")
	if _, err := os.Stat(markerPath); os.IsNotExist(err) {
		return fmt.Errorf("change %q is not implemented (run 'openspec complete %s' first)", id, id)
	}

	// Update registry status
	change, err := m.Get(id)
	if err != nil {
		return err
	}
	change.Status = StatusArchived
	if err := m.updateInRegistry(change); err != nil {
		return err
	}

	// Move to _archive
	archiveDir := filepath.Join(m.baseDir, "_archive")
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return err
	}
	dest := filepath.Join(archiveDir, id)
	if err := os.Rename(changeDir, dest); err != nil {
		return fmt.Errorf("failed to archive change: %w", err)
	}

	return nil
}

// SyncResult holds the outcome of an auto-sync operation.
type SyncResult struct {
	ChangesImported []string
	TasksCreated    int
	Errors          []string
}

// Sync checks all draft changes for populated tasks.md and auto-imports them.
// This eliminates the need for an explicit "openspec import" step.
// Already-imported changes are skipped via the status guard in Import().
func (m *Manager) Sync(tm *tasks.TaskManager) (*SyncResult, error) {
	result := &SyncResult{}

	reg, err := m.loadRegistry()
	if err != nil {
		return result, nil
	}

	for _, change := range reg.Changes {
		if change.Status != StatusDraft {
			continue
		}

		tasksPath := filepath.Join(m.baseDir, change.ID, "tasks.md")
		titles, err := ParseTasksFile(tasksPath)
		if err != nil || len(titles) == 0 {
			// tasks.md is empty or template-only — not ready for import, skip
			continue
		}

		created, err := m.Import(change.ID, tm)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", change.ID, err))
			continue
		}

		result.ChangesImported = append(result.ChangesImported, change.ID)
		result.TasksCreated += len(created)
	}

	return result, nil
}

// --- internal helpers ---

func (m *Manager) registryPath() string {
	return filepath.Join(m.baseDir, "changes.yaml")
}

func (m *Manager) loadRegistry() (*Registry, error) {
	path := m.registryPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Registry{}, nil
		}
		return nil, err
	}
	var reg Registry
	if err := yaml.Unmarshal(data, &reg); err != nil {
		return nil, err
	}
	return &reg, nil
}

func (m *Manager) saveRegistry(reg *Registry) error {
	if err := os.MkdirAll(m.baseDir, 0755); err != nil {
		return err
	}
	data, err := yaml.Marshal(reg)
	if err != nil {
		return err
	}
	return os.WriteFile(m.registryPath(), data, 0644)
}

func (m *Manager) addToRegistry(change *Change) error {
	reg, err := m.loadRegistry()
	if err != nil {
		return err
	}
	reg.Changes = append(reg.Changes, *change)
	return m.saveRegistry(reg)
}

func (m *Manager) updateInRegistry(change *Change) error {
	reg, err := m.loadRegistry()
	if err != nil {
		return err
	}
	for i := range reg.Changes {
		if reg.Changes[i].ID == change.ID {
			reg.Changes[i] = *change
			return m.saveRegistry(reg)
		}
	}
	return fmt.Errorf("change %q not found in registry", change.ID)
}

func (m *Manager) writeMetadata(changeDir string, change *Change) error {
	data, err := yaml.Marshal(change)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(changeDir, "metadata.yaml"), data, 0644)
}

var nonAlphaNum = regexp.MustCompile(`[^a-z0-9]+`)

func toKebabCase(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = nonAlphaNum.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
