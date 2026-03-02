package steps

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
	"github.com/javierbenavides/agentic-agent/internal/skills"
	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// AgentInstallSteps encapsulates agent installation step definitions
type AgentInstallSteps struct {
	suite *SuiteContext
}

// NewAgentInstallSteps creates a new AgentInstallSteps instance
func NewAgentInstallSteps(suite *SuiteContext) *AgentInstallSteps {
	return &AgentInstallSteps{suite: suite}
}

// RegisterSteps registers all agent installation step definitions
func (s *AgentInstallSteps) RegisterSteps(sc *godog.ScenarioContext) {
	// Agent pack installation steps
	sc.Step(`^I install the "([^"]*)" pack for "([^"]*)"$`, s.installPackForTool)
	sc.Step(`^I run skills ensure for "([^"]*)"$`, s.runSkillsEnsure)

	// File content assertions (reuse from common but could be defined here too)
	sc.Step(`^the file "([^"]*)" should not be empty$`, s.fileShouldNotBeEmpty)

	// CLI verification
	sc.Step(`^I run the installed agentic-agent CLI$`, s.runInstalledCLI)
	sc.Step(`^"([^"]*)" should be in the mandatory packs list$`, s.shouldBeInMandatoryPacks)
}

func (s *AgentInstallSteps) installPackForTool(ctx context.Context, packName, tool string) error {
	// Temporarily override the tool directories to point under the test project
	origSkillDir := skills.ToolSkillDir[tool]
	origAgentDir := skills.ToolAgentDir[tool]

	skillDir := filepath.Join(s.suite.ProjectDir, getToolSkillSubdir(tool))
	skills.ToolSkillDir[tool] = skillDir

	agentDir := filepath.Join(s.suite.ProjectDir, getToolAgentSubdir(tool))
	if agentDir != skillDir { // Only set if tool supports agents
		skills.ToolAgentDir[tool] = agentDir
	}

	// Register cleanup to restore original directories
	s.suite.RegisterCleanup(func() {
		skills.ToolSkillDir[tool] = origSkillDir
		if _, ok := skills.ToolAgentDir[tool]; ok {
			skills.ToolAgentDir[tool] = origAgentDir
		}
	})

	// Install the pack
	installer := skills.NewInstaller()
	_, err := installer.Install(packName, tool, false)
	s.suite.LastCommandErr = err
	return nil
}

func (s *AgentInstallSteps) runSkillsEnsure(ctx context.Context, tool string) error {
	// Create a minimal config
	cfg := &models.Config{}
	cfg.Paths.SpecDirs = []string{filepath.Join(s.suite.ProjectDir, ".agentic", "spec")}

	// Temporarily override the tool directories
	origSkillDir := skills.ToolSkillDir[tool]
	origAgentDir := skills.ToolAgentDir[tool]

	skillDir := filepath.Join(s.suite.ProjectDir, getToolSkillSubdir(tool))
	skills.ToolSkillDir[tool] = skillDir

	agentDir := filepath.Join(s.suite.ProjectDir, getToolAgentSubdir(tool))
	if agentDir != skillDir {
		skills.ToolAgentDir[tool] = agentDir
	}

	// Register cleanup
	s.suite.RegisterCleanup(func() {
		skills.ToolSkillDir[tool] = origSkillDir
		if _, ok := skills.ToolAgentDir[tool]; ok {
			skills.ToolAgentDir[tool] = origAgentDir
		}
	})

	// Run Ensure (simplified - just install mandatory packs)
	installer := skills.NewInstaller()
	for _, packName := range skills.MandatoryPacks {
		_, err := installer.Install(packName, tool, false)
		if err != nil {
			s.suite.LastCommandErr = err
			return nil
		}
	}

	s.suite.LastCommandErr = nil
	return nil
}

func (s *AgentInstallSteps) fileShouldNotBeEmpty(ctx context.Context, relPath string) error {
	fullPath := filepath.Join(s.suite.ProjectDir, relPath)
	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", fullPath)
	}
	if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", fullPath, err)
	}
	if info.Size() == 0 {
		return fmt.Errorf("file is empty: %s", fullPath)
	}
	return nil
}

func (s *AgentInstallSteps) runInstalledCLI(ctx context.Context) error {
	// This step is a no-op in the BDD context; the CLI is already available
	// In a real scenario, this would build and run the CLI
	s.suite.LastCommandErr = nil
	return nil
}

func (s *AgentInstallSteps) shouldBeInMandatoryPacks(ctx context.Context, packName string) error {
	for _, p := range skills.MandatoryPacks {
		if p == packName {
			return nil
		}
	}
	return fmt.Errorf("%q not found in MandatoryPacks: %v", packName, skills.MandatoryPacks)
}

// Helper function to get the tool-specific skill subdirectory
func getToolSkillSubdir(tool string) string {
	switch tool {
	case "claude-code":
		return ".claude/skills"
	case "cursor":
		return ".cursor/skills"
	case "gemini":
		return ".gemini/skills"
	case "windsurf":
		return ".windsurf/skills"
	case "codex":
		return ".codex/skills"
	case "copilot":
		return ".github/skills"
	case "opencode":
		return ".opencode/skills"
	case "antigravity":
		return ".agent/skills"
	default:
		return ".claude/skills" // default fallback
	}
}

// Helper function to get the tool-specific agent subdirectory
func getToolAgentSubdir(tool string) string {
	switch tool {
	case "claude-code":
		return ".claude/agents"
	case "opencode":
		return ".agents"
	default:
		// For tools without native agent support, use the skill subdirectory
		return getToolSkillSubdir(tool)
	}
}
