package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/javierbenavides/agentic-agent/internal/config"
	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// EnsureResult holds the outcome of an ensure operation.
type EnsureResult struct {
	Agent          string
	RulesGenerated bool
	RulesFile      string
	AgentRulesSet  bool
	ToolRulesSet   bool
	PacksInstalled []string
	DriftFixed     bool
	Warnings       []string
}

// ensureAgentRules copies AGENT_RULES.md to the project root if it doesn't exist or is outdated.
func ensureAgentRules() (bool, error) {
	targetPath := "AGENT_RULES.md"

	// Check if file exists
	if _, err := os.Stat(targetPath); err == nil {
		return false, nil // File already exists, skip
	}

	// Read template from configs/templates/init/
	templatePath := filepath.Join("configs", "templates", "init", "AGENT_RULES.md")
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return false, fmt.Errorf("failed to read AGENT_RULES.md template: %w", err)
	}

	// Write to project root
	if err := os.WriteFile(targetPath, content, 0644); err != nil {
		return false, fmt.Errorf("failed to write AGENT_RULES.md: %w", err)
	}

	return true, nil
}

// ensureToolRules copies the tool-specific rules file (e.g., CLAUDE.md) to the project root.
func ensureToolRules(agentName string) (bool, error) {
	// Map agent name to rules filename
	rulesFileMap := map[string]string{
		"claude-code": "CLAUDE.md",
		"copilot":     "COPILOT.md",
		"opencode":    "OPENCODE.md",
		"cursor":      "CURSOR.md",
		"windsurf":    "WINDSURF.md",
	}

	rulesFile, ok := rulesFileMap[agentName]
	if !ok {
		return false, nil // No tool-specific rules for this agent
	}

	targetPath := rulesFile

	// Check if file exists
	if _, err := os.Stat(targetPath); err == nil {
		return false, nil // File already exists, skip
	}

	// Read template from configs/templates/init/tool-rules/
	templatePath := filepath.Join("configs", "templates", "init", "tool-rules", rulesFile)
	content, err := os.ReadFile(templatePath)
	if err != nil {
		// If template doesn't exist, it's not an error - just skip
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to read %s template: %w", rulesFile, err)
	}

	// Write to project root
	if err := os.WriteFile(targetPath, content, 0644); err != nil {
		return false, fmt.Errorf("failed to write %s: %w", rulesFile, err)
	}

	return true, nil
}

// Ensure makes sure an agent has all necessary skills and rules.
// It is idempotent and safe to run repeatedly.
func Ensure(agentName string, cfg *models.Config) (*EnsureResult, error) {
	result := &EnsureResult{Agent: agentName}
	gen := NewGeneratorWithConfig(cfg)

	// 0. Ensure AGENT_RULES.md exists in project root
	agentRulesSet, err := ensureAgentRules()
	if err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Failed to set up AGENT_RULES.md: %v", err))
	} else {
		result.AgentRulesSet = agentRulesSet
	}

	// 0.5. Ensure tool-specific rules file exists (e.g., CLAUDE.md)
	toolRulesSet, err := ensureToolRules(agentName)
	if err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Failed to set up tool rules: %v", err))
	} else {
		result.ToolRulesSet = toolRulesSet
	}

	// 1. Check if rules file exists; generate if missing or drifted
	registry := NewSkillRegistry()
	skill, err := registry.GetSkill(agentName)
	if err != nil {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("No template registered for %s, skipping rules generation", agentName))
	} else {
		if _, statErr := os.Stat(skill.OutputFile); os.IsNotExist(statErr) {
			if genErr := gen.Generate(agentName); genErr != nil {
				return nil, fmt.Errorf("failed to generate rules for %s: %w", agentName, genErr)
			}
			result.RulesGenerated = true
			result.RulesFile = skill.OutputFile
		} else {
			// Check for drift
			drifted, _ := gen.CheckDriftFor(agentName)
			if len(drifted) > 0 {
				if genErr := gen.Generate(agentName); genErr != nil {
					result.Warnings = append(result.Warnings,
						fmt.Sprintf("Failed to fix drift for %s: %v", agentName, genErr))
				} else {
					result.DriftFixed = true
					result.RulesFile = skill.OutputFile
				}
			}
		}
	}

	// 2. Install mandatory skill packs + configured skill packs
	installer := NewInstaller()

	// Mandatory packs are always installed regardless of config
	for _, packName := range MandatoryPacks {
		if !installer.IsInstalled(packName, agentName) {
			_, installErr := installer.Install(packName, agentName, false)
			if installErr != nil {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("Failed to install mandatory pack %s: %v", packName, installErr))
			} else {
				result.PacksInstalled = append(result.PacksInstalled, packName)
			}
		}
	}

	// Additional configured packs from agent config
	agentCfg := config.GetAgentConfig(cfg, agentName)
	for _, packName := range agentCfg.SkillPacks {
		if !installer.IsInstalled(packName, agentName) {
			_, installErr := installer.Install(packName, agentName, false)
			if installErr != nil {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("Failed to install pack %s: %v", packName, installErr))
			} else {
				result.PacksInstalled = append(result.PacksInstalled, packName)
			}
		}
	}

	// 3. Generate prd and ralph-converter skill files for this agent
	if gen.Config != nil {
		if _, ok := ToolSkillDir[agentName]; ok {
			if err := gen.GenerateToolSkills(agentName); err != nil {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("Failed to generate skills for %s: %v", agentName, err))
			}
		}
	}

	return result, nil
}

// FormatEnsureResult returns a human-readable summary of what ensure did.
func FormatEnsureResult(r *EnsureResult) string {
	var b strings.Builder

	if r.AgentRulesSet {
		b.WriteString("Set up AGENT_RULES.md\n")
	}
	if r.ToolRulesSet {
		b.WriteString(fmt.Sprintf("Set up %s rules\n", strings.ToUpper(r.Agent)))
	}
	if r.RulesGenerated {
		b.WriteString(fmt.Sprintf("Generated rules: %s\n", r.RulesFile))
	}
	if r.DriftFixed {
		b.WriteString(fmt.Sprintf("Fixed drift: %s\n", r.RulesFile))
	}
	if len(r.PacksInstalled) > 0 {
		b.WriteString(fmt.Sprintf("Installed packs: %s\n", strings.Join(r.PacksInstalled, ", ")))
	}
	if !r.AgentRulesSet && !r.ToolRulesSet && !r.RulesGenerated && !r.DriftFixed && len(r.PacksInstalled) == 0 {
		b.WriteString("Already up to date.\n")
	}
	for _, w := range r.Warnings {
		b.WriteString(fmt.Sprintf("Warning: %s\n", w))
	}

	return b.String()
}
