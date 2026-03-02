package skills

import (
	"fmt"
	"os"
	"path/filepath"
)

// DetectSuperpowers checks if the Superpowers plugin is installed
// by looking for Superpowers skill files in both global and project directories.
func DetectSuperpowers() bool {
	// Check global Superpowers skills directory
	home, err := os.UserHomeDir()
	if err == nil {
		globalSuperpowersDir := filepath.Join(home, ".claude", "skills")
		if isSuperpowersInstalled(globalSuperpowersDir) {
			return true
		}
	}

	// Check project-level Superpowers skills directory
	if isSuperpowersInstalled(".claude/skills") {
		return true
	}

	return false
}

// isSuperpowersInstalled checks if a directory contains Superpowers skill files.
// Superpowers skills follow the pattern: superpowers:*
func isSuperpowersInstalled(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}

	superpowersSkills := []string{
		"superpowers:brainstorming",
		"superpowers:writing-plans",
		"superpowers:executing-plans",
		"superpowers:using-git-worktrees",
		"superpowers:test-driven-development",
		"superpowers:systematic-debugging",
		"superpowers:verification-before-completion",
	}

	for _, entry := range entries {
		name := entry.Name()
		for _, skill := range superpowersSkills {
			// Check for directory or file with Superpowers skill pattern
			if name == skill || name == skill+".md" {
				return true
			}
		}
	}

	return false
}

// SuperpowersInstallInstructions returns installation instructions for the Superpowers plugin.
func SuperpowersInstallInstructions() string {
	return `⚠️  Superpowers plugin not detected.

For enhanced TDD, debugging, and planning workflows, install the Superpowers plugin:

1. In Claude Code, run: /help → marketplace
2. Search for "Superpowers"
3. Install and restart Claude Code

Or visit: https://github.com/anthropics/superpowers

The following Superpowers skills will be available after installation:
  • superpowers:brainstorming — Explore user intent and requirements
  • superpowers:writing-plans — Create detailed implementation plans
  • superpowers:executing-plans — Execute plans with checkpoints
  • superpowers:using-git-worktrees — Create isolated workspaces
  • superpowers:test-driven-development — TDD workflow (red-green-refactor)
  • superpowers:systematic-debugging — 4-phase debugging protocol
  • superpowers:verification-before-completion — Verification gates

Continuing without Superpowers... (feature integration available but optional)`
}

// WarnIfSuperpowersNotFound prints a warning if Superpowers is not installed.
// Use this in commands where Superpowers integration would be beneficial.
func WarnIfSuperpowersNotFound() {
	if !DetectSuperpowers() {
		fmt.Println(SuperpowersInstallInstructions())
		fmt.Println()
	}
}
