package tracks

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/javierbenavides/agentic-agent/internal/plans"
	"github.com/javierbenavides/agentic-agent/internal/project"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"gopkg.in/yaml.v3"
)

// InitOptions holds optional fields for track initialization.
type InitOptions struct {
	Purpose     string
	Constraints string
	Success     string
}

// Registry is the index file listing all tracks.
type Registry struct {
	Tracks []models.Track `yaml:"tracks"`
}

// Manager handles track CRUD operations.
type Manager struct {
	baseDir string
}

// NewManager creates a track manager for the given base directory.
func NewManager(baseDir string) *Manager {
	return &Manager{baseDir: baseDir}
}

// Create creates a new track directory with brainstorm.md, spec.md, plan.md, and metadata.yaml.
// Uses enhanced templates with brainstorming scaffolding. Status starts as "ideation".
func (m *Manager) Create(name string, trackType models.TrackType, opts *InitOptions) (*models.Track, error) {
	id := toKebabCase(name)
	trackDir := filepath.Join(m.baseDir, id)

	if _, err := os.Stat(trackDir); err == nil {
		return nil, fmt.Errorf("track %q already exists", id)
	}

	if err := os.MkdirAll(trackDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create track directory: %w", err)
	}

	track := models.Track{
		ID:             id,
		Name:           name,
		Type:           trackType,
		Status:         models.TrackStatusIdeation,
		CreatedAt:      time.Now(),
		SpecPath:       filepath.Join(id, "spec.md"),
		PlanPath:       filepath.Join(id, "plan.md"),
		BrainstormPath: filepath.Join(id, "brainstorm.md"),
	}

	if opts == nil {
		opts = &InitOptions{}
	}

	tmplData := project.TrackTemplateData{
		Name:        name,
		Type:        string(trackType),
		Purpose:     opts.Purpose,
		Constraints: opts.Constraints,
		Success:     opts.Success,
	}

	// Write brainstorm.md from template
	brainstormContent, err := project.RenderTrackTemplate("brainstorm.md.tmpl", tmplData)
	if err != nil {
		return nil, fmt.Errorf("failed to render brainstorm.md: %w", err)
	}
	if err := os.WriteFile(filepath.Join(trackDir, "brainstorm.md"), []byte(brainstormContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write brainstorm.md: %w", err)
	}

	// Write spec.md from enhanced template
	specContent, err := project.RenderTrackTemplate("spec-enhanced.md.tmpl", tmplData)
	if err != nil {
		return nil, fmt.Errorf("failed to render spec.md: %w", err)
	}
	if err := os.WriteFile(filepath.Join(trackDir, "spec.md"), []byte(specContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write spec.md: %w", err)
	}

	// Write plan.md from template
	planContent, err := project.RenderTrackTemplate("plan-from-spec.md.tmpl", tmplData)
	if err != nil {
		return nil, fmt.Errorf("failed to render plan.md: %w", err)
	}
	if err := os.WriteFile(filepath.Join(trackDir, "plan.md"), []byte(planContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write plan.md: %w", err)
	}

	// Write metadata.yaml
	if err := m.writeMetadata(trackDir, &track); err != nil {
		return nil, err
	}

	// Update registry
	if err := m.addToRegistry(&track); err != nil {
		return nil, err
	}

	return &track, nil
}

// Activate validates a track's spec, generates a plan from it, and
// optionally decomposes the plan into tasks.
func (m *Manager) Activate(id string, decompose bool, taskDir string) ([]*models.Task, error) {
	track, err := m.Get(id)
	if err != nil {
		return nil, err
	}

	specPath := filepath.Join(m.baseDir, track.SpecPath)
	planPath := filepath.Join(m.baseDir, track.PlanPath)

	// Validate spec completeness
	report, err := ValidateSpec(specPath)
	if err != nil {
		return nil, fmt.Errorf("spec validation failed: %w", err)
	}
	if !report.Complete {
		return nil, fmt.Errorf("spec is incomplete, missing sections: %s (run 'track refine %s' to see details)", strings.Join(report.Missing, ", "), id)
	}

	// Generate plan from spec
	if err := plans.GenerateFromSpec(specPath, planPath, track.Name); err != nil {
		return nil, fmt.Errorf("plan generation failed: %w", err)
	}

	// Update status to active
	if err := m.UpdateStatus(id, models.TrackStatusActive); err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	// Optionally decompose into tasks
	if decompose {
		tm := tasks.NewTaskManager(taskDir)
		created, err := tasks.DecomposeFromPlan(planPath, id, tm)
		if err != nil {
			return nil, fmt.Errorf("task decomposition failed: %w", err)
		}
		// Link tasks to track
		for _, t := range created {
			if err := m.AddTask(id, t.ID); err != nil {
				return created, fmt.Errorf("failed to link task %s to track: %w", t.ID, err)
			}
		}
		return created, nil
	}

	return nil, nil
}

// List returns all tracks from the registry.
func (m *Manager) List() ([]models.Track, error) {
	reg, err := m.loadRegistry()
	if err != nil {
		return nil, err
	}
	return reg.Tracks, nil
}

// Get returns a single track by ID.
func (m *Manager) Get(id string) (*models.Track, error) {
	reg, err := m.loadRegistry()
	if err != nil {
		return nil, err
	}
	for _, t := range reg.Tracks {
		if t.ID == id {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("track %q not found", id)
}

// UpdateStatus changes a track's status and persists it.
func (m *Manager) UpdateStatus(id string, status models.TrackStatus) error {
	reg, err := m.loadRegistry()
	if err != nil {
		return err
	}
	found := false
	for i := range reg.Tracks {
		if reg.Tracks[i].ID == id {
			reg.Tracks[i].Status = status
			found = true
			// Update metadata file too
			trackDir := filepath.Join(m.baseDir, id)
			if err := m.writeMetadata(trackDir, &reg.Tracks[i]); err != nil {
				return err
			}
			break
		}
	}
	if !found {
		return fmt.Errorf("track %q not found", id)
	}
	return m.saveRegistry(reg)
}

// AddTask associates a task ID with a track.
func (m *Manager) AddTask(trackID, taskID string) error {
	reg, err := m.loadRegistry()
	if err != nil {
		return err
	}
	for i := range reg.Tracks {
		if reg.Tracks[i].ID == trackID {
			// Avoid duplicates
			for _, existing := range reg.Tracks[i].TaskIDs {
				if existing == taskID {
					return nil
				}
			}
			reg.Tracks[i].TaskIDs = append(reg.Tracks[i].TaskIDs, taskID)
			trackDir := filepath.Join(m.baseDir, trackID)
			if err := m.writeMetadata(trackDir, &reg.Tracks[i]); err != nil {
				return err
			}
			return m.saveRegistry(reg)
		}
	}
	return fmt.Errorf("track %q not found", trackID)
}

// Archive moves a track to the _archive directory and marks it archived.
func (m *Manager) Archive(id string) error {
	trackDir := filepath.Join(m.baseDir, id)
	if _, err := os.Stat(trackDir); os.IsNotExist(err) {
		return fmt.Errorf("track directory %q not found", id)
	}

	// Update registry status before moving the directory
	reg, err := m.loadRegistry()
	if err != nil {
		return err
	}
	found := false
	for i := range reg.Tracks {
		if reg.Tracks[i].ID == id {
			reg.Tracks[i].Status = models.TrackStatusArchived
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("track %q not found in registry", id)
	}
	if err := m.saveRegistry(reg); err != nil {
		return err
	}

	archiveDir := filepath.Join(m.baseDir, "_archive")
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return err
	}

	dest := filepath.Join(archiveDir, id)
	if err := os.Rename(trackDir, dest); err != nil {
		return fmt.Errorf("failed to archive track: %w", err)
	}

	return nil
}

// --- internal helpers ---

func (m *Manager) registryPath() string {
	return filepath.Join(m.baseDir, "tracks.yaml")
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
	data, err := yaml.Marshal(reg)
	if err != nil {
		return err
	}
	return os.WriteFile(m.registryPath(), data, 0644)
}

func (m *Manager) addToRegistry(track *models.Track) error {
	reg, err := m.loadRegistry()
	if err != nil {
		return err
	}
	reg.Tracks = append(reg.Tracks, *track)
	return m.saveRegistry(reg)
}

func (m *Manager) writeMetadata(trackDir string, track *models.Track) error {
	data, err := yaml.Marshal(track)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(trackDir, "metadata.yaml"), data, 0644)
}

var nonAlphaNum = regexp.MustCompile(`[^a-z0-9]+`)

func toKebabCase(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = nonAlphaNum.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
