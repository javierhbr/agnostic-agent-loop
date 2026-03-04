package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/javierbenavides/agentic-agent/internal/config"
	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// EnsureOptions controls how packs are installed during Ensure.
type EnsureOptions struct {
	Global  bool // Install to global user dir (~/.claude/skills/) instead of project-local
	Symlink bool // Create symlinks from destination to canonical copy (~/.agentic/skills/)
}

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

// ensureAgentRules appends AGENT_RULES.md to the project root if it doesn't exist, or appends content if it does.
func ensureAgentRules() (bool, error) {
	targetPath := "AGENT_RULES.md"

	// Read template from configs/templates/init/
	templatePath := filepath.Join("configs", "templates", "init", "AGENT_RULES.md")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return false, fmt.Errorf("failed to read AGENT_RULES.md template: %w", err)
	}

	// Check if file exists
	fileExists := true
	existingContent, err := os.ReadFile(targetPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return false, fmt.Errorf("failed to read existing AGENT_RULES.md: %w", err)
		}
		fileExists = false
	}

	var finalContent []byte
	if fileExists {
		// Append template content to existing file with a separator
		finalContent = append(existingContent, []byte("\n\n---\n\n")...)
		finalContent = append(finalContent, templateContent...)
	} else {
		// Create new file with template content
		finalContent = templateContent
	}

	// Write to project root
	if err := os.WriteFile(targetPath, finalContent, 0644); err != nil {
		return false, fmt.Errorf("failed to write AGENT_RULES.md: %w", err)
	}

	return true, nil
}

// ensureToolRules appends the tool-specific rules file (e.g., CLAUDE.md) to the project root.
// If the file exists, appends new content; otherwise creates it.
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

	// Read template from configs/templates/init/tool-rules/
	templatePath := filepath.Join("configs", "templates", "init", "tool-rules", rulesFile)
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		// If template doesn't exist, it's not an error - just skip
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to read %s template: %w", rulesFile, err)
	}

	// Check if file exists and read existing content
	fileExists := true
	existingContent, err := os.ReadFile(targetPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return false, fmt.Errorf("failed to read existing %s: %w", rulesFile, err)
		}
		fileExists = false
	}

	var finalContent []byte
	if fileExists {
		// Append template content to existing file with a separator
		finalContent = append(existingContent, []byte("\n\n---\n\n")...)
		finalContent = append(finalContent, templateContent...)
	} else {
		// Create new file with template content
		finalContent = templateContent
	}

	// Write to project root
	if err := os.WriteFile(targetPath, finalContent, 0644); err != nil {
		return false, fmt.Errorf("failed to write %s: %w", rulesFile, err)
	}

	return true, nil
}

// Ensure makes sure an agent has all necessary skills and rules.
// It is idempotent and safe to run repeatedly.
func Ensure(agentName string, cfg *models.Config, opts EnsureOptions) (*EnsureResult, error) {
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
		if !installer.IsInstalledAt(packName, agentName, opts.Global) {
			_, installErr := installer.Install(packName, agentName, opts.Global, opts.Symlink)
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
		if !installer.IsInstalledAt(packName, agentName, opts.Global) {
			_, installErr := installer.Install(packName, agentName, opts.Global, opts.Symlink)
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

// FormatEnsureResultCompact returns a compressed single-line summary of what ensure did.
func FormatEnsureResultCompact(r *EnsureResult) string {
	var items []string

	if r.AgentRulesSet {
		items = append(items, "AGENT_RULES")
	}
	if r.ToolRulesSet {
		items = append(items, fmt.Sprintf("%s-rules", strings.ToLower(r.Agent)))
	}
	if r.RulesGenerated {
		items = append(items, "rules")
	}
	if r.DriftFixed {
		items = append(items, "drift-fixed")
	}
	if len(r.PacksInstalled) > 0 {
		items = append(items, fmt.Sprintf("%d-packs", len(r.PacksInstalled)))
	}

	if len(items) == 0 {
		return "✓ up-to-date"
	}

	result := fmt.Sprintf("✓ %s", strings.Join(items, " + "))

	if len(r.Warnings) > 0 {
		result += fmt.Sprintf(" ⚠️ %d warnings", len(r.Warnings))
	}

	return result
}
