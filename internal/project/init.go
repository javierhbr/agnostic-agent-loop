package project

import (
	"fmt"
	"os"
	"strings"
)

// ProjectProfile holds optional profile data gathered during init.
type ProjectProfile struct {
	Description string
	TechStack   string
	Workflow    string
}

// InitProject initializes a new Agnostic Agent project.
func InitProject(name string) error {
	return InitProjectWithProfile(name, nil)
}

// InitProjectWithProfile initializes a project with optional profile data.
func InitProjectWithProfile(name string, profile *ProjectProfile) error {
	fmt.Printf("Initializing project: %s\n", name)

	// 1. Create Directory Structure
	if err := CreateStructure(); err != nil {
		return err
	}

	// 2. Write Default Files
	files := map[string]string{
		"init/agnostic-agent.yaml":            "agnostic-agent.yaml",
		"init/tasks/backlog.yaml":             ".agentic/tasks/backlog.yaml",
		"init/tasks/in-progress.yaml":         ".agentic/tasks/in-progress.yaml",
		"init/tasks/done.yaml":                ".agentic/tasks/done.yaml",
		"init/context/global-context.md":      ".agentic/context/global-context.md",
		"init/context/rolling-summary.md":     ".agentic/context/rolling-summary.md",
		"init/context/decisions.md":           ".agentic/context/decisions.md",
		"init/context/assumptions.md":         ".agentic/context/assumptions.md",
		"init/context/tech-stack.md":          ".agentic/context/tech-stack.md",
		"init/context/workflow-preferences.md": ".agentic/context/workflow-preferences.md",
		"init/agent-rules/base.md":            ".agentic/agent-rules/base.md",
	}

	for tmpl, dest := range files {
		if err := WriteTemplate(tmpl, dest); err != nil {
			return fmt.Errorf("failed to write %s: %w", dest, err)
		}
		fmt.Printf("Created %s\n", dest)
	}

	// 3. Enrich context files with profile data if provided
	if profile != nil {
		if err := writeProfileContext(name, profile); err != nil {
			return fmt.Errorf("failed to write profile context: %w", err)
		}
	}

	fmt.Println("Project initialized successfully.")
	return nil
}

func writeProfileContext(name string, profile *ProjectProfile) error {
	// Enrich global-context.md with description
	if profile.Description != "" {
		content := fmt.Sprintf("# Global Context\n\n## Project Overview\n\n%s\n\n## Goals\n\n## Guidelines\n", strings.TrimSpace(profile.Description))
		if err := os.WriteFile(".agentic/context/global-context.md", []byte(content), 0644); err != nil {
			return err
		}
	}

	// Enrich tech-stack.md
	if profile.TechStack != "" {
		content := fmt.Sprintf("# Tech Stack\n\n%s\n", strings.TrimSpace(profile.TechStack))
		if err := os.WriteFile(".agentic/context/tech-stack.md", []byte(content), 0644); err != nil {
			return err
		}
	}

	// Enrich workflow-preferences.md
	if profile.Workflow != "" {
		content := fmt.Sprintf("# Workflow Preferences\n\n%s\n", strings.TrimSpace(profile.Workflow))
		if err := os.WriteFile(".agentic/context/workflow-preferences.md", []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}
