package specs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// ResolvedSpec holds the result of resolving a single spec reference.
type ResolvedSpec struct {
	Ref     string `yaml:"ref" json:"ref"`
	Path    string `yaml:"path,omitempty" json:"path,omitempty"`
	Content string `yaml:"content,omitempty" json:"content,omitempty"`
	Found   bool   `yaml:"found" json:"found"`
	Error   string `yaml:"error,omitempty" json:"error,omitempty"`
}

// Resolver resolves spec references against configured directories.
type Resolver struct {
	specDirs []string
}

// NewResolver creates a Resolver from the given config.
func NewResolver(cfg *models.Config) *Resolver {
	dirs := []string{".agentic/spec"}
	if cfg != nil && len(cfg.Paths.SpecDirs) > 0 {
		dirs = cfg.Paths.SpecDirs
	}
	return &Resolver{specDirs: dirs}
}

// ResolveSpec resolves a single spec reference.
// Resolution order: absolute/relative path first, then search each specDir.
func (r *Resolver) ResolveSpec(ref string) *ResolvedSpec {
	// Try as absolute or relative path first
	if filepath.IsAbs(ref) {
		return r.tryPath(ref, ref)
	}

	// Try relative to cwd
	if info, err := os.Stat(ref); err == nil && !info.IsDir() {
		abs, _ := filepath.Abs(ref)
		return r.readSpec(ref, abs)
	}

	// Search configured spec directories
	for _, dir := range r.specDirs {
		candidate := filepath.Join(dir, ref)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			abs, _ := filepath.Abs(candidate)
			return r.readSpec(ref, abs)
		}
	}

	return &ResolvedSpec{
		Ref:   ref,
		Found: false,
		Error: fmt.Sprintf("spec %q not found in any configured directory", ref),
	}
}

// ResolveAll resolves a batch of spec references.
func (r *Resolver) ResolveAll(refs []string) []*ResolvedSpec {
	results := make([]*ResolvedSpec, 0, len(refs))
	for _, ref := range refs {
		results = append(results, r.ResolveSpec(ref))
	}
	return results
}

// ReadSpec is a convenience method that resolves and returns content or an error.
func (r *Resolver) ReadSpec(ref string) (string, error) {
	resolved := r.ResolveSpec(ref)
	if !resolved.Found {
		return "", fmt.Errorf("%s", resolved.Error)
	}
	return resolved.Content, nil
}

// ListSpecs lists all spec files found across configured directories.
func (r *Resolver) ListSpecs() ([]*ResolvedSpec, error) {
	var results []*ResolvedSpec
	seen := make(map[string]bool)

	for _, dir := range r.specDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("failed to read spec dir %s: %w", dir, err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if seen[name] {
				continue // first directory wins
			}
			seen[name] = true

			abs, _ := filepath.Abs(filepath.Join(dir, name))
			results = append(results, r.readSpec(name, abs))
		}
	}

	return results, nil
}

func (r *Resolver) tryPath(ref, path string) *ResolvedSpec {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return &ResolvedSpec{
			Ref:   ref,
			Found: false,
			Error: fmt.Sprintf("spec %q not found at path %s", ref, path),
		}
	}
	return r.readSpec(ref, path)
}

func (r *Resolver) readSpec(ref, absPath string) *ResolvedSpec {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return &ResolvedSpec{
			Ref:   ref,
			Path:  absPath,
			Found: false,
			Error: fmt.Sprintf("failed to read spec %q: %v", ref, err),
		}
	}
	return &ResolvedSpec{
		Ref:     ref,
		Path:    absPath,
		Content: string(data),
		Found:   true,
	}
}
