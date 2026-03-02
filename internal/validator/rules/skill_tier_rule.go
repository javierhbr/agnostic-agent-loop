package rules

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/javierbenavides/agentic-agent/internal/validator"
)

const (
	RouterWarnLines   = 70
	RouterFailLines   = 100
	SkillWarnLines    = 130
	SkillFailLines    = 200
	ResourceWarnLines = 500
)

var requiredSkillSections = []string{
	"does exactly this",
	"when to use",
	"if you need more detail",
}

// SkillTierRule enforces the 3-tier layered context model.
// Tier 1: Routers (SKILLS.md) stay slim (<70 lines)
// Tier 2: Skill files (SKILL.md) stay focused (54-130 lines)
// Tier 3: Resources (resources/*.md) hold detail on-demand (150-500 lines)
type SkillTierRule struct {
	PacksDir string // overridable for tests; defaults to "internal/skills/packs"
}

func (r *SkillTierRule) Name() string {
	return "skill-tier"
}

func (r *SkillTierRule) Validate(ctx *validator.ValidationContext) (*validator.RuleResult, error) {
	result := &validator.RuleResult{
		RuleName: r.Name(),
		Status:   "PASS",
		Errors:   []string{},
	}

	// Resolve packs directory
	packsDir := r.PacksDir
	if packsDir == "" {
		packsDir = filepath.Join(ctx.ProjectRoot, "internal", "skills", "packs")
	}

	// Check if packs directory exists; skip silently if not
	if _, err := os.Stat(packsDir); os.IsNotExist(err) {
		return result, nil
	}

	// Check root router once
	rootRouter := filepath.Join(packsDir, "SKILLS.md")
	checkRouter(rootRouter, result)

	// Walk through packs directory
	entries, err := os.ReadDir(packsDir)
	if err != nil {
		return result, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		packDir := filepath.Join(packsDir, entry.Name())

		// Check top-level SKILL.md
		skillPath := filepath.Join(packDir, "SKILL.md")
		if _, err := os.Stat(skillPath); err == nil {
			checkSkillFile(skillPath, result)
		}

		// Check sub-router (e.g., sdd/SKILLS.md)
		subRouter := filepath.Join(packDir, "SKILLS.md")
		if _, err := os.Stat(subRouter); err == nil {
			checkRouter(subRouter, result)
		}

		// Check resources directory
		resourcesDir := filepath.Join(packDir, "resources")
		if resourceEntries, err := os.ReadDir(resourcesDir); err == nil {
			for _, resEntry := range resourceEntries {
				if !resEntry.IsDir() && strings.HasSuffix(resEntry.Name(), ".md") {
					resPath := filepath.Join(resourcesDir, resEntry.Name())
					checkResourceFile(resPath, result)
				}
			}
		}

		// Check subdirectories (e.g., sdd/analyst/)
		subEntries, err := os.ReadDir(packDir)
		if err != nil {
			continue
		}

		for _, subEntry := range subEntries {
			if !subEntry.IsDir() || subEntry.Name() == "resources" {
				continue
			}

			subPackDir := filepath.Join(packDir, subEntry.Name())

			// Check SKILL.md in subdirectory
			subSkillPath := filepath.Join(subPackDir, "SKILL.md")
			if _, err := os.Stat(subSkillPath); err == nil {
				checkSkillFile(subSkillPath, result)
			}

			// Check resources in subdirectory
			subResourcesDir := filepath.Join(subPackDir, "resources")
			if subResEntries, err := os.ReadDir(subResourcesDir); err == nil {
				for _, subResEntry := range subResEntries {
					if !subResEntry.IsDir() && strings.HasSuffix(subResEntry.Name(), ".md") {
						subResPath := filepath.Join(subResourcesDir, subResEntry.Name())
						checkResourceFile(subResPath, result)
					}
				}
			}
		}
	}

	return result, nil
}

func checkRouter(path string, result *validator.RuleResult) {
	lines, err := countLines(path)
	if err != nil {
		addFail(result, fmt.Sprintf("Failed to read router %s: %v", relPath(path), err))
		return
	}

	if lines > RouterFailLines {
		addFail(result, fmt.Sprintf("Router %s: %d lines (fail threshold: %d)", relPath(path), lines, RouterFailLines))
	} else if lines > RouterWarnLines {
		addWarn(result, fmt.Sprintf("Router %s: %d lines (warn threshold: %d)", relPath(path), lines, RouterWarnLines))
	}
}

func checkSkillFile(path string, result *validator.RuleResult) {
	lines, err := countLines(path)
	if err != nil {
		addFail(result, fmt.Sprintf("Failed to read skill %s: %v", relPath(path), err))
		return
	}

	if lines > SkillFailLines {
		addFail(result, fmt.Sprintf("Skill %s: %d lines (fail threshold: %d)", relPath(path), lines, SkillFailLines))
	} else if lines > SkillWarnLines {
		addWarn(result, fmt.Sprintf("Skill %s: %d lines (warn threshold: %d)", relPath(path), lines, SkillWarnLines))
	}

	// Check required sections
	missing, err := hasSections(path, requiredSkillSections)
	if err != nil {
		addFail(result, fmt.Sprintf("Failed to check sections in %s: %v", relPath(path), err))
		return
	}

	if len(missing) > 0 {
		addFail(result, fmt.Sprintf("Skill %s missing required sections: %v", relPath(path), missing))
	}

	// Check resource links
	if err := checkResourceLinks(path, result); err != nil {
		addFail(result, fmt.Sprintf("Failed to check links in %s: %v", relPath(path), err))
	}
}

func checkResourceFile(path string, result *validator.RuleResult) {
	lines, err := countLines(path)
	if err != nil {
		addFail(result, fmt.Sprintf("Failed to read resource %s: %v", relPath(path), err))
		return
	}

	if lines > ResourceWarnLines {
		addWarn(result, fmt.Sprintf("Resource %s: %d lines (warn threshold: %d)", relPath(path), lines, ResourceWarnLines))
	}
}

func countLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func hasSections(path string, required []string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	found := make(map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "##") {
			// Extract heading text
			heading := strings.TrimPrefix(line, "##")
			heading = strings.TrimSpace(heading)
			heading = strings.ToLower(heading)

			// Check for substring match
			for _, req := range required {
				if strings.Contains(heading, req) {
					found[req] = true
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Collect missing sections
	var missing []string
	for _, req := range required {
		if !found[req] {
			missing = append(missing, req)
		}
	}

	return missing, nil
}

func checkResourceLinks(skillPath string, result *validator.RuleResult) error {
	file, err := os.Open(skillPath)
	if err != nil {
		return err
	}
	defer file.Close()

	skillDir := filepath.Dir(skillPath)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		// Look for actual resource links: "→ `resources/..."
		if !strings.Contains(line, "→ `resources/") {
			continue
		}

		// Extract the resource path
		// Format: → `resources/file.md` or → `resources/file.md#anchor`
		parts := strings.Split(line, "→ `")
		if len(parts) < 2 {
			continue
		}

		linkPart := parts[1]

		// Find the closing backtick
		if idx := strings.Index(linkPart, "`"); idx != -1 {
			linkPart = linkPart[:idx]
		}

		// Extract just the file path (remove anchor)
		if idx := strings.Index(linkPart, "#"); idx != -1 {
			linkPart = linkPart[:idx]
		}

		// Build full path
		targetPath := filepath.Join(skillDir, linkPart)

		// Verify file exists
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			addFail(result, fmt.Sprintf("Broken link in %s: %s (file not found)", relPath(skillPath), linkPart))
		}
	}

	return scanner.Err()
}

func relPath(abs string) string {
	// Find "packs/" in the path and strip everything before it
	if idx := strings.Index(abs, "packs/"); idx != -1 {
		return abs[idx:]
	}
	return abs
}

func addFail(result *validator.RuleResult, msg string) {
	result.Status = "FAIL"
	result.Errors = append(result.Errors, msg)
}

func addWarn(result *validator.RuleResult, msg string) {
	if result.Status == "PASS" {
		result.Status = "WARN"
	}
	result.Errors = append(result.Errors, msg)
}
