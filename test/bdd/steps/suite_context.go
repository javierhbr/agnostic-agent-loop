package steps

import (
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// SuiteContext holds shared state across step definitions within a scenario
// This is a re-export of the parent package's SuiteContext for convenience
type SuiteContext struct {
	T              *testing.T
	ProjectDir     string
	CurrentTask    *models.Task
	LastTaskID     string
	LastCommandOut string
	LastCommandErr error
	CleanupFuncs   []func()
}

// RegisterCleanup adds a cleanup function to be run after the scenario
func (s *SuiteContext) RegisterCleanup(fn func()) {
	s.CleanupFuncs = append(s.CleanupFuncs, fn)
}
