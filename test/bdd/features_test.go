package bdd

import (
	"context"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/javierbenavides/agentic-agent/tests/bdd/steps"
)

var opts = godog.Options{
	Format:      "pretty",
	Paths:       []string{"../../features"},
	Tags:        "~@wip", // Exclude work-in-progress scenarios
	Concurrency: 1,       // Run scenarios sequentially to avoid race conditions
}

func TestFeatures(t *testing.T) {
	// Enable more parallelism in CI environments
	if os.Getenv("CI") == "true" {
		opts.Concurrency = 8
	}

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

// InitializeScenario registers all step definitions for each scenario
func InitializeScenario(sc *godog.ScenarioContext) {
	// Create suite context for this scenario
	t := &testing.T{}
	suite := NewSuiteContext(t)

	// Convert to steps.SuiteContext for step definitions
	stepsSuite := &steps.SuiteContext{
		T:            suite.T,
		ProjectDir:   "",
		CurrentTask:  nil,
		LastTaskID:   "",
		LastCommandOut: "",
		LastCommandErr: nil,
		CleanupFuncs: suite.CleanupFuncs,
	}

	// Register all step definitions
	steps.NewCommonSteps(stepsSuite).RegisterSteps(sc)
	steps.NewTaskSteps(stepsSuite).RegisterSteps(sc)
	steps.NewInitSteps(stepsSuite).RegisterSteps(sc)
	steps.NewAssertionSteps(stepsSuite).RegisterSteps(sc)

	// Cleanup after each scenario
	sc.After(func(ctx context.Context, _ *godog.Scenario, err error) (context.Context, error) {
		for _, fn := range stepsSuite.CleanupFuncs {
			fn()
		}
		return ctx, nil
	})
}
