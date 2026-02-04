package context

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type RollingContextManager struct {
	baseDir string
}

func NewRollingContextManager(baseDir string) *RollingContextManager {
	return &RollingContextManager{baseDir: baseDir}
}

func (rcm *RollingContextManager) LoadRolling() (string, error) {
	path := filepath.Join(rcm.baseDir, "rolling-summary.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (rcm *RollingContextManager) AppendEntry(summary string) error {
	path := filepath.Join(rcm.baseDir, "rolling-summary.md")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	entry := fmt.Sprintf("\n## Entry %s\n%s\n", time.Now().Format(time.RFC3339), summary)
	if _, err := f.WriteString(entry); err != nil {
		return err
	}
	return nil
}
