package skills

import (
	"fmt"
	"os"
	"strings"

	"github.com/javierbenavides/agentic-agent/internal/config"
	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// EnsureResult holds the outcome of an ensure operation.
type EnsureResult struct {
	Agent          string
	RulesGenerated bool
	RulesFile      string
	PacksInstalled []string
	DriftFixed     bool
	Warnings       []string
}

// Ensure makes sure an agent has all necessary skills and rules.
// It is idempotent and safe to run repeatedly.
func Ensure(agentName string, cfg *models.Config) (*EnsureResult, error) {
	result := &EnsureResult{Agent: agentName}
	gen := NewGeneratorWithConfig(cfg)

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

	// 2. Install configured skill packs
	agentCfg := config.GetAgentConfig(cfg, agentName)
	if len(agentCfg.SkillPacks) > 0 {
		installer := NewInstaller()
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
	}

	// 3. Generate tool-specific skill files
	switch agentName {
	case "claude-code":
		if gen.Config != nil {
			if err := gen.GenerateClaudeCodeSkills(); err != nil {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("Failed to generate Claude Code skills: %v", err))
			}
		}
	case "gemini":
		if gen.Config != nil {
			if err := gen.GenerateGeminiSkills(); err != nil {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("Failed to generate Gemini skills: %v", err))
			}
		}
	}

	return result, nil
}

// FormatEnsureResult returns a human-readable summary of what ensure did.
func FormatEnsureResult(r *EnsureResult) string {
	var b strings.Builder

	if r.RulesGenerated {
		b.WriteString(fmt.Sprintf("Generated rules: %s\n", r.RulesFile))
	}
	if r.DriftFixed {
		b.WriteString(fmt.Sprintf("Fixed drift: %s\n", r.RulesFile))
	}
	if len(r.PacksInstalled) > 0 {
		b.WriteString(fmt.Sprintf("Installed packs: %s\n", strings.Join(r.PacksInstalled, ", ")))
	}
	if !r.RulesGenerated && !r.DriftFixed && len(r.PacksInstalled) == 0 {
		b.WriteString("Already up to date.\n")
	}
	for _, w := range r.Warnings {
		b.WriteString(fmt.Sprintf("Warning: %s\n", w))
	}

	return b.String()
}
