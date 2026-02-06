package steps

import (
	"github.com/cucumber/godog"
)

// InitSteps encapsulates initialization-related step definitions
type InitSteps struct {
	suite *SuiteContext
}

// NewInitSteps creates a new InitSteps instance
func NewInitSteps(suite *SuiteContext) *InitSteps {
	return &InitSteps{suite: suite}
}

// RegisterSteps registers all initialization-related step definitions
func (s *InitSteps) RegisterSteps(sc *godog.ScenarioContext) {
	// Init steps are primarily handled by CommonSteps
	// This file exists for organization and future expansion
	// Specific init-related assertions can be added here
}
