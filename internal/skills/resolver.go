package skills

import (
	"fmt"
	"io/fs"
	"strings"
)

// ResolvedSkill holds the result of resolving a single skill pack reference.
type ResolvedSkill struct {
	Ref     string `yaml:"ref" json:"ref"`
	Content string `yaml:"content,omitempty" json:"content,omitempty"`
	Found   bool   `yaml:"found" json:"found"`
	Error   string `yaml:"error,omitempty" json:"error,omitempty"`
}

// ResolveSkillRefs loads skill content for the given refs.
// Resolution order per ref:
//  1. Agent's installed skill directory (e.g., .claude/skills/<pack>/SKILL.md)
//  2. Any installed tool via ReadInstalledSkillFromAnyTool
//  3. Embedded pack content from packsFS
func ResolveSkillRefs(refs []string, agent string) []ResolvedSkill {
	if len(refs) == 0 {
		return nil
	}

	results := make([]ResolvedSkill, 0, len(refs))
	for _, ref := range refs {
		results = append(results, resolveOneSkill(ref, agent))
	}
	return results
}

func resolveOneSkill(ref, agent string) ResolvedSkill {
	// 1. Try the active agent's installed skill directory
	if agent != "" {
		content, err := ReadInstalledSkill(agent, ref)
		if err == nil {
			return ResolvedSkill{Ref: ref, Content: content, Found: true}
		}
	}

	// 2. Try any installed tool
	content, _, err := ReadInstalledSkillFromAnyTool(ref)
	if err == nil {
		return ResolvedSkill{Ref: ref, Content: content, Found: true}
	}

	// 3. Fall back to embedded pack content
	content, err = readEmbeddedSkill(ref)
	if err == nil {
		return ResolvedSkill{Ref: ref, Content: content, Found: true}
	}

	return ResolvedSkill{
		Ref:   ref,
		Found: false,
		Error: fmt.Sprintf("skill pack %q could not be resolved", ref),
	}
}

// readEmbeddedSkill reads the SKILL.md (and any resources) from the embedded packsFS.
func readEmbeddedSkill(packName string) (string, error) {
	mainPath := fmt.Sprintf("packs/%s/SKILL.md", packName)
	mainContent, err := fs.ReadFile(packsFS, mainPath)
	if err != nil {
		return "", fmt.Errorf("embedded skill pack %q not found: %w", packName, err)
	}

	var parts []string
	parts = append(parts, string(mainContent))

	// Also include any resource files in packs/<pack>/resources/
	resourceDir := fmt.Sprintf("packs/%s/resources", packName)
	entries, err := fs.ReadDir(packsFS, resourceDir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			resContent, err := fs.ReadFile(packsFS, fmt.Sprintf("%s/%s", resourceDir, entry.Name()))
			if err == nil {
				parts = append(parts, string(resContent))
			}
		}
	}

	return strings.Join(parts, "\n\n---\n\n"), nil
}
