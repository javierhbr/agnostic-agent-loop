package steps

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
	"github.com/javierbenavides/agentic-agent/internal/skills"
)

// DetectionSteps encapsulates agent-detection step definitions.
type DetectionSteps struct {
	suite      *SuiteContext
	projectDir string
	envVars    map[string]string // env vars to set before detection
	detected   skills.DetectedAgent
	allAgents  []skills.DetectedAgent
}

// NewDetectionSteps creates a new DetectionSteps instance.
func NewDetectionSteps(suite *SuiteContext) *DetectionSteps {
	return &DetectionSteps{
		suite:   suite,
		envVars: make(map[string]string),
	}
}

// RegisterSteps registers all agent-detection step definitions.
func (s *DetectionSteps) RegisterSteps(sc *godog.ScenarioContext) {
	// Setup steps
	sc.Step(`^the environment variable "([^"]*)" is set to "([^"]*)"$`, s.setEnvVar)
	sc.Step(`^a project root with directory "([^"]*)"$`, s.createProjectDirectory)
	sc.Step(`^a project root with file "([^"]*)"$`, s.createProjectFile)
	sc.Step(`^a project root with (directory|file) "([^"]*)"$`, s.createProjectMarker)
	sc.Step(`^an empty project root$`, s.createEmptyProjectRoot)

	// Action steps
	sc.Step(`^I detect the agent with flag "([^"]*)"$`, s.detectWithFlag)
	sc.Step(`^I detect the agent without a flag$`, s.detectWithoutFlag)
	sc.Step(`^I detect all agents$`, s.detectAllAgents)

	// Assertion steps
	sc.Step(`^the detected agent name should be "([^"]*)"$`, s.agentNameShouldBe)
	sc.Step(`^the detected agent source should be "([^"]*)"$`, s.agentSourceShouldBe)
	sc.Step(`^the detected agents should include "([^"]*)"$`, s.agentsShouldInclude)
	sc.Step(`^the detected agent count should be (\d+)$`, s.agentCountShouldBe)
}

// --- Setup ---

func (s *DetectionSteps) setEnvVar(_ context.Context, key, value string) error {
	s.envVars[key] = value
	return nil
}

func (s *DetectionSteps) ensureProjectDir() {
	if s.projectDir == "" {
		dir, err := os.MkdirTemp("", "detect-test-*")
		if err != nil {
			return
		}
		s.projectDir = dir
		s.suite.RegisterCleanup(func() { os.RemoveAll(s.projectDir) })
	}
}

func (s *DetectionSteps) createProjectDirectory(_ context.Context, dirName string) error {
	s.ensureProjectDir()
	return os.MkdirAll(filepath.Join(s.projectDir, dirName), 0755)
}

func (s *DetectionSteps) createProjectFile(_ context.Context, fileName string) error {
	s.ensureProjectDir()
	return os.WriteFile(filepath.Join(s.projectDir, fileName), []byte("# marker"), 0644)
}

func (s *DetectionSteps) createProjectMarker(_ context.Context, markerType, name string) error {
	s.ensureProjectDir()
	if markerType == "directory" {
		return os.MkdirAll(filepath.Join(s.projectDir, name), 0755)
	}
	return os.WriteFile(filepath.Join(s.projectDir, name), []byte("# marker"), 0644)
}

func (s *DetectionSteps) createEmptyProjectRoot(_ context.Context) error {
	s.ensureProjectDir()
	return nil
}

// --- Actions ---

func (s *DetectionSteps) applyEnv() (restore func()) {
	originals := make(map[string]string)
	for key := range s.envVars {
		originals[key] = os.Getenv(key)
	}
	for key, val := range s.envVars {
		os.Setenv(key, val)
	}
	return func() {
		for key := range s.envVars {
			if orig, ok := originals[key]; ok && orig != "" {
				os.Setenv(key, orig)
			} else {
				os.Unsetenv(key)
			}
		}
	}
}

func (s *DetectionSteps) detectWithFlag(_ context.Context, flag string) error {
	s.ensureProjectDir()
	restore := s.applyEnv()
	defer restore()

	s.detected = skills.DetectAgent(flag, s.projectDir)
	return nil
}

func (s *DetectionSteps) detectWithoutFlag(_ context.Context) error {
	s.ensureProjectDir()
	restore := s.applyEnv()
	defer restore()

	s.detected = skills.DetectAgent("", s.projectDir)
	return nil
}

func (s *DetectionSteps) detectAllAgents(_ context.Context) error {
	s.ensureProjectDir()
	s.allAgents = skills.DetectAllAgents(s.projectDir)
	return nil
}

// --- Assertions ---

func (s *DetectionSteps) agentNameShouldBe(_ context.Context, expected string) error {
	if s.detected.Name != expected {
		return fmt.Errorf("expected agent name %q, got %q", expected, s.detected.Name)
	}
	return nil
}

func (s *DetectionSteps) agentSourceShouldBe(_ context.Context, expected string) error {
	if s.detected.Source != expected {
		return fmt.Errorf("expected agent source %q, got %q", expected, s.detected.Source)
	}
	return nil
}

func (s *DetectionSteps) agentsShouldInclude(_ context.Context, name string) error {
	for _, a := range s.allAgents {
		if a.Name == name {
			return nil
		}
	}
	return fmt.Errorf("expected detected agents to include %q, got %v", name, s.allAgents)
}

func (s *DetectionSteps) agentCountShouldBe(_ context.Context, expected int) error {
	if len(s.allAgents) != expected {
		return fmt.Errorf("expected %d agents, got %d: %v", expected, len(s.allAgents), s.allAgents)
	}
	return nil
}
