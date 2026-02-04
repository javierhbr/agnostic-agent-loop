package context

import (
	"os"
	"path/filepath"
	"time"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

type GlobalContextManager struct {
	baseDir string
}

func NewGlobalContextManager(baseDir string) *GlobalContextManager {
	return &GlobalContextManager{baseDir: baseDir}
}

func (gcm *GlobalContextManager) LoadGlobal() (*models.GlobalContext, error) {
	path := filepath.Join(gcm.baseDir, "global-context.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// For MVP, we are not parsing Markdown into struct fields deeply yet,
	// just loading it. If we need structured data, we might need frontmatter/TOON.
	// But models.GlobalContext has structured fields.
	// Real implementation would parse MD sections.
	// For now, let's assume we might just read the raw content or return a struct with basic info.

	return &models.GlobalContext{
		Overview: string(data),
		Updated:  time.Now(), // Mock timestamp for file mod time
	}, nil
}

func (gcm *GlobalContextManager) UpdateGlobal(content string) error {
	path := filepath.Join(gcm.baseDir, "global-context.md")
	return os.WriteFile(path, []byte(content), 0644)
}
