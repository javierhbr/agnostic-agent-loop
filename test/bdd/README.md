# BDD Testing with Godog

This directory contains Behavior-Driven Development (BDD) tests using Gherkin feature files and the godog framework.

## Structure

```
tests/bdd/
├── features_test.go         # Main test runner
├── suite_context.go         # Shared test context
└── steps/                   # Step definition packages
    ├── common_steps.go      # Setup & command execution
    ├── task_steps.go        # Task operations
    ├── init_steps.go        # Initialization
    ├── assertion_steps.go   # Generic assertions
    └── suite_context.go     # Context wrapper

features/                    # Gherkin feature files
├── init/
│   └── project_initialization.feature
├── tasks/
├── context/
├── validation/
└── workflows/
```

## Running Tests

```bash
# Run BDD tests only
make test-bdd

# Run BDD tests with verbose output
make test-bdd-verbose

# Run all tests (unit + functional + BDD)
make test-all

# Run with coverage
make coverage-all
```

## Writing New Features

### 1. Create a Feature File

Create a `.feature` file in the appropriate directory under `features/`:

```gherkin
Feature: Task Creation
  As a developer
  I want to create tasks
  So that I can track work

  Scenario: Create a simple task
    Given a clean test environment
    And I have initialized a project
    When I create a task with title "My Task"
    Then a task should be created successfully
    And the task should appear in the backlog
```

### 2. Implement Step Definitions

If needed, add step definitions in `tests/bdd/steps/`:

```go
func (s *TaskSteps) RegisterSteps(sc *godog.ScenarioContext) {
    sc.Step(`^I create a task with title "([^"]*)"$`, s.createTaskWithTitle)
}

func (s *TaskSteps) createTaskWithTitle(ctx context.Context, title string) error {
    // Implementation
    return nil
}
```

### 3. Run Tests

```bash
make test-bdd
```

## Available Step Definitions

### Common Steps (Setup & Execution)
- `Given a clean test environment`
- `When I run "init ProjectName"`
- `When I initialize a project with name "ProjectName"`

### Task Steps
- `When I create a task with title "Title"`
- `When I claim the task`
- `When I complete the task`
- `When I list all tasks`
- `Then the task should appear in the backlog`
- `Then the task should be in "state" state`
- `Then I should see N task(s) in backlog`

### Assertion Steps
- `Then the command should succeed`
- `Then the command should fail`
- `Then the error message should contain "text"`
- `Then the following directories should exist:`
- `Then the following files should exist:`
- `Then the project structure should be created`
- `Then git should be initialized`

## Reusing Existing Helpers

The BDD step definitions reuse helpers from `tests/functional/helpers.go`:
- `SetupTestProject()` - Creates isolated test environment
- `VerifyTaskFile()` - Parses and validates task YAML
- `VerifyTaskInFile()` - Finds specific task in file
- `VerifyProjectStructure()` - Validates project structure

## ATDD Workflow

Follow this workflow for new features:

1. **Write Feature First** (Specification)
   ```gherkin
   Scenario: Archive old tasks
     Given I have completed tasks older than 30 days
     When I run "task archive --older-than=30d"
     Then tasks should be moved to archived.yaml
   ```

2. **Run Tests** (They fail - undefined steps)
   ```bash
   make test-bdd
   ```

3. **Implement Step Definitions**

4. **Implement Feature** in CLI commands

5. **Tests Pass**
   ```bash
   make test-bdd  # Green!
   ```

## Current Test Results

```
✅ 2 scenarios (2 passed)
✅ 10 steps (10 passed)
✅ 52ms execution time
```

## Next Steps

See `/Users/javierbenavides/.claude/plans/distributed-hugging-swan.md` for the complete implementation plan including:
- Converting existing workflow scenarios to Gherkin
- Adding task lifecycle features
- Error handling scenarios
- CI/CD integration
- Comprehensive documentation
