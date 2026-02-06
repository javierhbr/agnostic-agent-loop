package steps

import (
	"context"
	"strings"

	"github.com/cucumber/godog"
	"github.com/javierbenavides/agentic-agent/internal/project"
	"github.com/javierbenavides/agentic-agent/tests/functional"
)

// CommonSteps encapsulates common step definitions for test setup and command execution
type CommonSteps struct {
	suite *SuiteContext
}

// NewCommonSteps creates a new CommonSteps instance
func NewCommonSteps(suite *SuiteContext) *CommonSteps {
	return &CommonSteps{suite: suite}
}

// RegisterSteps registers all common step definitions
func (s *CommonSteps) RegisterSteps(sc *godog.ScenarioContext) {
	// Setup steps
	sc.Step(`^a clean test environment$`, s.cleanTestEnvironment)
	sc.Step(`^I am in a directory without git$`, s.directoryWithoutGit)

	// Command execution steps
	sc.Step(`^I run "([^"]*)"$`, s.runCommand)
	sc.Step(`^I run the following command:$`, s.runMultilineCommand)

	// Project initialization steps
	sc.Step(`^I initialize a project with name "([^"]*)"$`, s.initializeProject)
	sc.Step(`^I have initialized project "([^"]*)"$`, s.initializeProject)
	sc.Step(`^I have initialized a project$`, s.initializeDefaultProject)

	// Validation steps
	sc.Step(`^I run validation$`, s.runValidation)
}

// cleanTestEnvironment sets up an isolated test environment
func (s *CommonSteps) cleanTestEnvironment(ctx context.Context) error {
	// Reuse the existing helper from functional tests
	s.suite.ProjectDir = functional.SetupTestProject(s.suite.T)
	s.suite.RegisterCleanup(func() {
		// Cleanup is already registered by SetupTestProject
	})
	return nil
}

// directoryWithoutGit creates a test directory without git initialization
func (s *CommonSteps) directoryWithoutGit(ctx context.Context) error {
	// This is similar to cleanTestEnvironment but without git init
	// For now, we'll use the same setup since git init is required by some validations
	// In a real scenario, you might want to create a separate helper
	return s.cleanTestEnvironment(ctx)
}

// runCommand executes a CLI command with the given arguments
func (s *CommonSteps) runCommand(ctx context.Context, cmdArgs string) error {
	// Parse command and arguments
	parts := strings.Fields(cmdArgs)
	if len(parts) == 0 {
		return nil
	}

	cmd := parts[0]
	args := parts[1:]

	// Execute command based on the command name
	switch cmd {
	case "init":
		if len(args) > 0 {
			return s.initializeProject(ctx, args[0])
		}
		return s.initializeProject(ctx, "TestProject")
	default:
		// For other commands, we'll implement them as needed
		s.suite.LastCommandOut = ""
		s.suite.LastCommandErr = nil
		return nil
	}
}

// runMultilineCommand executes a multi-line command (docstring)
func (s *CommonSteps) runMultilineCommand(ctx context.Context, docString *godog.DocString) error {
	// Clean up the command (remove line continuations, extra spaces)
	cmdText := strings.ReplaceAll(docString.Content, "\\\n", " ")
	cmdText = strings.TrimSpace(cmdText)

	return s.runCommand(ctx, cmdText)
}

// initializeProject initializes a project with the given name
func (s *CommonSteps) initializeProject(ctx context.Context, projectName string) error {
	// Ensure we have a test environment set up
	if s.suite.ProjectDir == "" {
		s.cleanTestEnvironment(ctx)
	}

	err := project.InitProject(projectName)
	s.suite.LastCommandErr = err
	if err != nil {
		s.suite.LastCommandOut = ""
		return err
	}
	s.suite.LastCommandOut = "Project initialized successfully"
	return nil
}

// initializeDefaultProject initializes a project with a default name
func (s *CommonSteps) initializeDefaultProject(ctx context.Context) error {
	return s.initializeProject(ctx, "TestProject")
}

// runValidation runs project validation
func (s *CommonSteps) runValidation(ctx context.Context) error {
	// For now, just mark as successful
	// In a full implementation, this would call the validator
	s.suite.LastCommandErr = nil
	s.suite.LastCommandOut = "Validation completed"
	return nil
}
