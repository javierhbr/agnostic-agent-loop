package openspec

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// configTemplate is the agnostic-agent.yaml content created by EnsureConfig.
// It includes multi-directory spec resolution paths and workflow validators
// suited for spec-driven development.
const configTemplate = `project:
  name: {{PROJECT_NAME}}
  version: 0.1.0
  roots:
    - .

agents:
  defaults:
    max_tokens: 4000
    model: claude-3-5-sonnet-20241022

# Multi-directory spec resolution:
# Specs are searched in order — first match wins.
paths:
  specDirs:
    - .specify/specs       # Spec Kit specs (if using Spec Kit)
    - openspec/specs        # OpenSpec specs (if using OpenSpec)
    - .agentic/spec         # Agentic native specs (fallback)
  openSpecDir: .agentic/openspec/changes  # OpenSpec change lifecycle
  contextDirs:
    - .agentic/context

workflow:
  validators:
    - context-check
    - task-scope
`

// EnsureConfigResult holds the outcome of the config creation attempt.
type EnsureConfigResult struct {
	Created bool
	Path    string
}

// EnsureConfig creates agnostic-agent.yaml in the given directory if it
// doesn't already exist. The projectName is used to set project.name.
// Returns whether the file was created and its path.
func EnsureConfig(dir, projectName string) (*EnsureConfigResult, error) {
	configPath := filepath.Join(dir, "agnostic-agent.yaml")

	if _, err := os.Stat(configPath); err == nil {
		return &EnsureConfigResult{Created: false, Path: configPath}, nil
	}

	name := projectName
	if name == "" {
		name = filepath.Base(dir)
	}

	content := strings.Replace(configTemplate, "{{PROJECT_NAME}}", name, 1)

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write %s: %w", configPath, err)
	}

	return &EnsureConfigResult{Created: true, Path: configPath}, nil
}
