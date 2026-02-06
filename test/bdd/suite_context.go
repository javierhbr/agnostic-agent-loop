package bdd

import (
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// SuiteContext holds shared state across step definitions within a scenario
type SuiteContext struct {
	T              *testing.T
	ProjectDir     string
	CurrentTask    *models.Task
	LastTaskID     string
	LastCommandOut string
	LastCommandErr error
	CleanupFuncs   []func()
}

// NewSuiteContext creates a new test suite context
func NewSuiteContext(t *testing.T) *SuiteContext {
	return &SuiteContext{
		T:            t,
		CleanupFuncs: make([]func(), 0),
	}
}

// Cleanup runs all registered cleanup functions
func (s *SuiteContext) Cleanup() {
	for _, fn := range s.CleanupFuncs {
		fn()
	}
}

// RegisterCleanup adds a cleanup function to be run after the scenario
func (s *SuiteContext) RegisterCleanup(fn func()) {
	s.CleanupFuncs = append(s.CleanupFuncs, fn)
}
