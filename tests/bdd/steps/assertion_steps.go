package steps

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

// AssertionSteps encapsulates generic assertion step definitions
type AssertionSteps struct {
	suite *SuiteContext
}

// NewAssertionSteps creates a new AssertionSteps instance
func NewAssertionSteps(suite *SuiteContext) *AssertionSteps {
	return &AssertionSteps{suite: suite}
}

// RegisterSteps registers all assertion step definitions
func (s *AssertionSteps) RegisterSteps(sc *godog.ScenarioContext) {
	// Command result assertions
	sc.Step(`^the command should succeed$`, s.commandShouldSucceed)
	sc.Step(`^the command should fail$`, s.commandShouldFail)
	sc.Step(`^the error message should contain "([^"]*)"$`, s.errorMessageShouldContain)

	// File system assertions
	sc.Step(`^the following directories should exist:$`, s.directoriesShouldExist)
	sc.Step(`^the following files should exist:$`, s.filesShouldExist)
	sc.Step(`^the project structure should be created$`, s.projectStructureShouldBeCreated)

	// Git assertions
	sc.Step(`^git should be initialized$`, s.gitShouldBeInitialized)
}

// commandShouldSucceed asserts that the last command succeeded
func (s *AssertionSteps) commandShouldSucceed(ctx context.Context) error {
	if s.suite.LastCommandErr != nil {
		return fmt.Errorf("command failed: %w", s.suite.LastCommandErr)
	}
	return nil
}

// commandShouldFail asserts that the last command failed
func (s *AssertionSteps) commandShouldFail(ctx context.Context) error {
	if s.suite.LastCommandErr == nil {
		return fmt.Errorf("expected command to fail, but it succeeded")
	}
	return nil
}

// errorMessageShouldContain asserts that the error message contains the expected text
func (s *AssertionSteps) errorMessageShouldContain(ctx context.Context, expectedText string) error {
	if s.suite.LastCommandErr == nil {
		return fmt.Errorf("no error occurred, cannot check error message")
	}

	errorMsg := s.suite.LastCommandErr.Error()
	if !contains(errorMsg, expectedText) {
		return fmt.Errorf("error message %q does not contain %q", errorMsg, expectedText)
	}

	return nil
}

// directoriesShouldExist asserts that the specified directories exist
func (s *AssertionSteps) directoriesShouldExist(ctx context.Context, table *godog.Table) error {
	for _, row := range table.Rows[1:] { // Skip header row
		dirPath := filepath.Join(s.suite.ProjectDir, row.Cells[0].Value)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", dirPath)
		}
	}
	return nil
}

// filesShouldExist asserts that the specified files exist
func (s *AssertionSteps) filesShouldExist(ctx context.Context, table *godog.Table) error {
	for _, row := range table.Rows[1:] { // Skip header row
		filePath := filepath.Join(s.suite.ProjectDir, row.Cells[0].Value)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", filePath)
		}
	}
	return nil
}

// projectStructureShouldBeCreated asserts that the standard project structure exists
func (s *AssertionSteps) projectStructureShouldBeCreated(ctx context.Context) error {
	requiredDirs := []string{
		".agentic/tasks",
		".agentic/context",
		".agentic/spec",
		".agentic/agent-rules",
	}

	for _, dir := range requiredDirs {
		dirPath := filepath.Join(s.suite.ProjectDir, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			return fmt.Errorf("required directory does not exist: %s", dir)
		}
	}

	requiredFiles := []string{
		"agnostic-agent.yaml",
		".agentic/tasks/backlog.yaml",
		".agentic/tasks/in-progress.yaml",
		".agentic/tasks/done.yaml",
	}

	for _, file := range requiredFiles {
		filePath := filepath.Join(s.suite.ProjectDir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("required file does not exist: %s", file)
		}
	}

	return nil
}

// gitShouldBeInitialized asserts that git is initialized in the project directory
func (s *AssertionSteps) gitShouldBeInitialized(ctx context.Context) error {
	gitDir := filepath.Join(s.suite.ProjectDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("git is not initialized (no .git directory)")
	}
	return nil
}

// contains checks if a string contains a substring (case-insensitive helper)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
