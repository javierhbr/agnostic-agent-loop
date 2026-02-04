package context

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

type DirectoryContextManager struct {
	rootDir string
}

func NewDirectoryContextManager(rootDir string) *DirectoryContextManager {
	return &DirectoryContextManager{rootDir: rootDir}
}

func (dcm *DirectoryContextManager) LoadContext(dir string) (*models.DirectoryContext, error) {
	path := filepath.Join(dir, "context.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Mock parsing for MVP
	return &models.DirectoryContext{
		Path:    dir,
		Purpose: string(data), // Just dump all content here for now
	}, nil
}

func (dcm *DirectoryContextManager) SaveContext(dir string, ctx *models.DirectoryContext) error {
	content := fmt.Sprintf("# Context for %s\n\n## Purpose\n%s\n\n## Responsibilities\n%s\n\n## Dependencies\n%s\n",
		ctx.Path, ctx.Purpose, strings.Join(ctx.Responsibilities, "\n- "), strings.Join(ctx.Dependencies, "\n- "))

	path := filepath.Join(dir, "context.md")
	return os.WriteFile(path, []byte(content), 0644)
}

func (dcm *DirectoryContextManager) FindContextDirs(root string) ([]string, error) {
	var dirs []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() != ".git" && info.Name() != ".agentic" {
			// Check if it has source files
			hasSource := false
			entries, _ := os.ReadDir(path)
			for _, e := range entries {
				if !e.IsDir() && (strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), ".ts") || strings.HasSuffix(e.Name(), ".js")) {
					hasSource = true
					break
				}
			}
			if hasSource {
				dirs = append(dirs, path)
			}
		}
		return nil
	})
	return dirs, err
}
