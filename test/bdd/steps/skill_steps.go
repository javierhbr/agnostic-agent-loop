package steps

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
	"github.com/javierbenavides/agentic-agent/internal/skills"
	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// SkillSteps encapsulates skill-related step definitions
type SkillSteps struct {
	suite       *SuiteContext
	driftResult []string
}

// NewSkillSteps creates a new SkillSteps instance
func NewSkillSteps(suite *SuiteContext) *SkillSteps {
	return &SkillSteps{suite: suite}
}

// RegisterSteps registers all skill-related step definitions
func (s *SkillSteps) RegisterSteps(sc *godog.ScenarioContext) {
	// Skill generation steps
	sc.Step(`^I generate skills for "([^"]*)"$`, s.generateSkillForTool)
	sc.Step(`^I generate skills for all tools$`, s.generateSkillsForAll)
	sc.Step(`^I generate Gemini slash commands$`, s.generateGeminiSlashCommands)
	sc.Step(`^I generate Gemini slash commands with PRD path "([^"]*)"$`, s.generateGeminiSlashCommandsWithPath)

	// Drift detection steps
	sc.Step(`^I check for skill drift$`, s.checkSkillDrift)
	sc.Step(`^drift should be detected for "([^"]*)"$`, s.driftShouldBeDetectedFor)

	// File content and existence assertions
	sc.Step(`^the file "([^"]*)" should exist$`, s.fileShouldExist)
	sc.Step(`^the file "([^"]*)" should contain "([^"]*)"$`, s.fileShouldContain)
	sc.Step(`^I modify the file "([^"]*)" with "([^"]*)"$`, s.modifyFile)

	// Registry assertions
	sc.Step(`^the skill registry should contain "([^"]*)"$`, s.registryShouldContain)

	// Custom base rules
	sc.Step(`^I have custom base rules:$`, s.setCustomBaseRules)
}

// generateSkillForTool generates a skill file for the specified tool
func (s *SkillSteps) generateSkillForTool(ctx context.Context, tool string) error {
	gen := skills.NewGenerator()
	err := gen.Generate(tool)
	s.suite.LastCommandErr = err
	return nil
}

// generateSkillsForAll generates skill files for all registered tools
func (s *SkillSteps) generateSkillsForAll(ctx context.Context) error {
	gen := skills.NewGenerator()
	registry := skills.NewSkillRegistry()

	var lastErr error
	for _, skill := range registry.GetAll() {
		if err := gen.Generate(skill.ToolName); err != nil {
			lastErr = err
		}
	}
	s.suite.LastCommandErr = lastErr
	return nil
}

// generateGeminiSlashCommands generates Gemini slash command TOML files with default config
func (s *SkillSteps) generateGeminiSlashCommands(ctx context.Context) error {
	cfg := &models.Config{}
	cfg.Paths.PRDOutputPath = ".agentic/tasks/"
	gen := skills.NewGeneratorWithConfig(cfg)

	err := gen.GenerateGeminiSkills()
	s.suite.LastCommandErr = err
	return nil
}

// generateGeminiSlashCommandsWithPath generates Gemini slash commands with a custom PRD path
func (s *SkillSteps) generateGeminiSlashCommandsWithPath(ctx context.Context, prdPath string) error {
	cfg := &models.Config{}
	cfg.Paths.PRDOutputPath = prdPath
	gen := skills.NewGeneratorWithConfig(cfg)

	err := gen.GenerateGeminiSkills()
	s.suite.LastCommandErr = err
	return nil
}

// checkSkillDrift checks for drift in all registered skill files
func (s *SkillSteps) checkSkillDrift(ctx context.Context) error {
	gen := skills.NewGenerator()
	drifted, err := gen.CheckDrift()
	if err != nil {
		s.suite.LastCommandErr = err
		return nil
	}
	s.driftResult = drifted
	s.suite.LastCommandErr = nil
	return nil
}

// driftShouldBeDetectedFor asserts that drift was detected for a specific file
func (s *SkillSteps) driftShouldBeDetectedFor(ctx context.Context, filePath string) error {
	for _, d := range s.driftResult {
		if strings.Contains(d, filePath) {
			return nil
		}
	}
	return fmt.Errorf("no drift detected for %s, drifted files: %v", filePath, s.driftResult)
}

// fileShouldExist asserts that a file exists relative to the project directory
func (s *SkillSteps) fileShouldExist(ctx context.Context, relPath string) error {
	fullPath := filepath.Join(s.suite.ProjectDir, relPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", fullPath)
	}
	return nil
}

// fileShouldContain asserts that a file contains expected text
func (s *SkillSteps) fileShouldContain(ctx context.Context, relPath, expectedText string) error {
	fullPath := filepath.Join(s.suite.ProjectDir, relPath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", fullPath, err)
	}

	if !strings.Contains(string(content), expectedText) {
		return fmt.Errorf("file %s does not contain %q", relPath, expectedText)
	}
	return nil
}

// modifyFile overwrites a file with the given content
func (s *SkillSteps) modifyFile(ctx context.Context, relPath, content string) error {
	fullPath := filepath.Join(s.suite.ProjectDir, relPath)
	return os.WriteFile(fullPath, []byte(content), 0644)
}

// registryShouldContain asserts that the skill registry contains the specified tool
func (s *SkillSteps) registryShouldContain(ctx context.Context, toolName string) error {
	registry := skills.NewSkillRegistry()
	_, err := registry.GetSkill(toolName)
	if err != nil {
		return fmt.Errorf("skill registry does not contain %q: %w", toolName, err)
	}
	return nil
}

// setCustomBaseRules writes custom base rules from a table
func (s *SkillSteps) setCustomBaseRules(ctx context.Context, table *godog.Table) error {
	var rules []string
	for _, row := range table.Rows[1:] { // Skip header
		rules = append(rules, "- "+row.Cells[0].Value)
	}

	content := "# Base Agent Rules\n\n" + strings.Join(rules, "\n") + "\n"
	rulesPath := filepath.Join(s.suite.ProjectDir, ".agentic/agent-rules/base.md")
	return os.WriteFile(rulesPath, []byte(content), 0644)
}
