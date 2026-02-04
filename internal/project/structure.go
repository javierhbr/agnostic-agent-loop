package project

import (
	"embed"
	"fmt"
	"os"
)

//go:embed templates/*
var templatesFS embed.FS

const (
	DirAgentic    = ".agentic"
	DirSpec       = ".agentic/spec"
	DirTasks      = ".agentic/tasks"
	DirContext    = ".agentic/context"
	DirAgentRules = ".agentic/agent-rules"
)

var requiredDirs = []string{
	DirAgentic,
	DirSpec,
	DirTasks,
	DirContext,
	DirAgentRules,
}

// CreateStructure creates the .agentic directory structure in the current directory.
func CreateStructure() error {
	for _, dir := range requiredDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

// WriteTemplate writes a template file to the destination path.
func WriteTemplate(tmplPath, destPath string) error {
	content, err := templatesFS.ReadFile("templates/" + tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", tmplPath, err)
	}

	if err := os.WriteFile(destPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", destPath, err)
	}
	return nil
}
