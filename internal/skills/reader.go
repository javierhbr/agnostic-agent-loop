package skills

import (
	"fmt"
	"os"
	"path/filepath"
)

// ReadInstalledSkill reads the SKILL.md content from an installed pack location.
// It looks in the tool's project-level skill directory.
func ReadInstalledSkill(tool, packName string) (string, error) {
	dir, ok := ToolSkillDir[tool]
	if !ok {
		return "", fmt.Errorf("unsupported tool: %s", tool)
	}

	skillPath := filepath.Join(dir, packName, "SKILL.md")
	content, err := os.ReadFile(skillPath)
	if err != nil {
		return "", fmt.Errorf("skill pack %q not found for tool %q at %s: %w", packName, tool, skillPath, err)
	}

	return string(content), nil
}

// ReadInstalledSkillFromAnyTool tries to read a skill pack from any known tool directory.
// Returns the content and the tool name where it was found.
func ReadInstalledSkillFromAnyTool(packName string) (string, string, error) {
	for tool, dir := range ToolSkillDir {
		skillPath := filepath.Join(dir, packName, "SKILL.md")
		content, err := os.ReadFile(skillPath)
		if err == nil {
			return string(content), tool, nil
		}
	}
	return "", "", fmt.Errorf("skill pack %q not installed for any tool", packName)
}
