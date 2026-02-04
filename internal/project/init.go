package project

import (
	"fmt"
)

// InitProject initializes a new Agnostic Agent project.
func InitProject(name string) error {
	fmt.Printf("Initializing project: %s\n", name)

	// 1. Create Directory Structure
	if err := CreateStructure(); err != nil {
		return err
	}

	// 2. Write Default Files
	files := map[string]string{
		"init/agnostic-agent.yaml":        "agnostic-agent.yaml",
		"init/tasks/backlog.yaml":         ".agentic/tasks/backlog.yaml",
		"init/tasks/in-progress.yaml":     ".agentic/tasks/in-progress.yaml",
		"init/tasks/done.yaml":            ".agentic/tasks/done.yaml",
		"init/context/global-context.md":  ".agentic/context/global-context.md",
		"init/context/rolling-summary.md": ".agentic/context/rolling-summary.md",
		"init/context/decisions.md":       ".agentic/context/decisions.md",
		"init/context/assumptions.md":     ".agentic/context/assumptions.md",
		"init/agent-rules/base.md":        ".agentic/agent-rules/base.md",
	}

	for tmpl, dest := range files {
		if err := WriteTemplate(tmpl, dest); err != nil {
			// Warn but maybe don't fail hard? OR fail hard.
			// Ideally we want to fail if we can't write, but we also want to create empty files if templates don't exist?
			// For now, fail hard.
			return fmt.Errorf("failed to write %s: %w", dest, err)
		}
		fmt.Printf("Created %s\n", dest)
	}

	fmt.Println("Project initialized successfully.")
	return nil
}
