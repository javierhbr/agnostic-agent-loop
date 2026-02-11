package simplify

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	appcontext "github.com/javierbenavides/agentic-agent/internal/context"
	"github.com/javierbenavides/agentic-agent/internal/skills"
	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// SimplifyBundle is a focused context bundle for code simplification review.
type SimplifyBundle struct {
	SkillInstructions string                    `yaml:"skill_instructions" json:"skill_instructions"`
	Directories       []*models.DirectoryContext `yaml:"directories" json:"directories"`
	TargetFiles       []string                  `yaml:"target_files,omitempty" json:"target_files,omitempty"`
	TechStack         string                    `yaml:"tech_stack,omitempty" json:"tech_stack,omitempty"`
	BuiltAt           time.Time                 `yaml:"built_at" json:"built_at"`
}

// BuildSimplifyBundle creates a simplification-oriented context bundle for the given directories.
func BuildSimplifyBundle(dirs []string, agent string, cfg *models.Config) (*SimplifyBundle, error) {
	if len(dirs) == 0 {
		return nil, fmt.Errorf("at least one directory is required")
	}

	// Resolve code-simplification skill
	resolved := skills.ResolveSkillRefs([]string{"code-simplification"}, agent)
	if len(resolved) == 0 || !resolved[0].Found {
		errMsg := "unknown error"
		if len(resolved) > 0 {
			errMsg = resolved[0].Error
		}
		return nil, fmt.Errorf("code-simplification skill pack not available: %s", errMsg)
	}

	// Generate/load directory contexts
	var dirContexts []*models.DirectoryContext
	for _, dir := range dirs {
		dirCtx, err := appcontext.GenerateContextWithConfig(dir, cfg)
		if err != nil {
			// Try loading existing context as fallback
			dcm := appcontext.NewDirectoryContextManager(dir)
			existing, loadErr := dcm.LoadContext(dir)
			if loadErr != nil {
				return nil, fmt.Errorf("failed to generate or load context for %s: %w", dir, err)
			}
			dirContexts = append(dirContexts, existing)
			continue
		}
		dirContexts = append(dirContexts, dirCtx)
	}

	// Scan target files
	var targetFiles []string
	for _, dir := range dirs {
		files, err := scanSourceFiles(dir)
		if err != nil {
			continue
		}
		targetFiles = append(targetFiles, files...)
	}

	// Load tech stack if available
	techStack, _ := os.ReadFile(".agentic/context/tech-stack.md")

	return &SimplifyBundle{
		SkillInstructions: resolved[0].Content,
		Directories:       dirContexts,
		TargetFiles:       targetFiles,
		TechStack:         string(techStack),
		BuiltAt:           time.Now(),
	}, nil
}

// scanSourceFiles returns source file paths in a directory, skipping common non-source files.
func scanSourceFiles(dir string) ([]string, error) {
	var files []string
	sourceExts := map[string]bool{
		".go": true, ".js": true, ".ts": true, ".tsx": true, ".jsx": true,
		".py": true, ".rb": true, ".rs": true, ".java": true, ".kt": true,
		".c": true, ".cpp": true, ".h": true, ".hpp": true,
		".cs": true, ".swift": true, ".dart": true,
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		// Skip hidden directories and common non-source dirs
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" || name == "__pycache__" {
				return filepath.SkipDir
			}
			return nil
		}
		ext := filepath.Ext(path)
		if sourceExts[ext] {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
