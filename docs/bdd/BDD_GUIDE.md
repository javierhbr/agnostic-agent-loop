# BDD Guide for agentic-agent CLI

## Table of Contents

1. [Introduction](#introduction)
2. [What is ATDD/BDD?](#what-is-atddbdd)
3. [Why We Use BDD](#why-we-use-bdd)
4. [Getting Started](#getting-started)
5. [Writing Gherkin Features](#writing-gherkin-features)
6. [Step Definitions](#step-definitions)
7. [Running BDD Tests](#running-bdd-tests)
8. [ATDD Workflow](#atdd-workflow)
9. [Best Practices](#best-practices)
10. [Examples](#examples)

---

## Introduction

This guide explains how to use Behavior-Driven Development (BDD) and Acceptance Test-Driven Development (ATDD) for testing the agentic-agent CLI application. We use Gherkin syntax for writing executable specifications and the godog framework for Go.

## What is ATDD/BDD?

**Acceptance Test-Driven Development (ATDD)** is a development methodology where:
- You write acceptance tests **before** implementing features
- Tests are written in plain language that stakeholders understand
- Tests serve as both specifications and validation

**Behavior-Driven Development (BDD)** is a flavor of ATDD that:
- Uses **Gherkin** syntax (Given/When/Then)
- Focuses on behavior from the user's perspective
- Creates living documentation that stays in sync with code

### Gherkin Syntax

Gherkin uses a simple structure:

```gherkin
Feature: Short description of the feature
  As a <role>
  I want <goal>
  So that <benefit>

  Scenario: Specific example of behavior
    Given <initial context>
    When <action or event>
    Then <expected outcome>
```

## Why We Use BDD

### 1. **Living Documentation**
Feature files serve as up-to-date documentation:
- Always accurate (tests fail if out of sync)
- Readable by non-developers
- Examples of how to use the CLI

### 2. **Clear Requirements**
Features define acceptance criteria upfront:
- No ambiguity about what "done" means
- Stakeholders can review before implementation
- Test-first mindset

### 3. **Regression Protection**
Automated tests prevent breaking changes:
- Run on every commit via CI/CD
- Catch bugs before they reach production
- Safe refactoring

### 4. **Better Design**
Writing tests first improves design:
- Forces thinking about user experience
- Identifies edge cases early
- Simpler, more testable code

## Getting Started

### Prerequisites

```bash
# Go 1.22 or higher
go version

# godog framework
go install github.com/cucumber/godog/cmd/godog@latest
```

### Installation

Initialize the BDD infrastructure:

```bash
make bdd-init
```

This creates:
- `features/` - Gherkin feature files
- `test/bdd/` - Step definitions and test infrastructure

### Project Structure

```
agnostic-agent-loop/
├── features/                      # Gherkin feature files
│   ├── init/
│   │   └── project_initialization.feature
│   ├── tasks/
│   │   ├── task_lifecycle.feature
│   │   └── error_handling.feature
│   ├── context/
│   │   └── context_generation.feature
│   └── workflows/
│       ├── beginner_workflow.feature
│       ├── intermediate_workflow.feature
│       └── advanced_workflow.feature
├── test/bdd/                     # BDD test infrastructure
│   ├── features_test.go          # Main test runner
│   ├── suite_context.go          # Shared test context
│   └── steps/                    # Step definitions
│       ├── common_steps.go
│       ├── task_steps.go
│       ├── assertion_steps.go
│       └── ...
└── Makefile
```

## Writing Gherkin Features

### Anatomy of a Feature File

```gherkin
Feature: Task Creation
  As a developer using the CLI
  I want to create tasks with metadata
  So that I have clear requirements to implement

  Background:
    Given a clean test environment
    And I have initialized a project

  @smoke @task-creation
  Scenario: Create a simple task
    When I create a task with title "Implement login"
    Then the command should succeed
    And the task should appear in the backlog

  @task-metadata
  Scenario: Create task with acceptance criteria
    When I create a task with the following details:
      | field       | value              |
      | title       | Implement API      |
      | description | Create REST API    |
    And I add the following acceptance criteria:
      | criterion                  |
      | GET /users returns users   |
      | POST /users creates user   |
    Then the task should have 2 acceptance criteria
```

### Gherkin Keywords

- **Feature**: High-level description of functionality
- **Background**: Steps that run before every scenario
- **Scenario**: Specific example of behavior
- **Given**: Set up initial context
- **When**: Perform action
- **Then**: Assert expected outcome
- **And**: Continue previous step type
- **But**: Negative continuation
- **@tags**: Organize and filter scenarios

### Data Tables

Use tables for structured data:

```gherkin
When I create the following files:
  | file                     | content              |
  | src/api/users.go        | package api          |
  | tests/api/users_test.go | package api_test     |
```

### Doc Strings

Use doc strings for multi-line text:

```gherkin
When I run the following command:
  """
  task create \
    --title "Build API" \
    --description "Create endpoints"
  """
```

## Step Definitions

Step definitions connect Gherkin steps to Go code.

### Example Step Definition

**Feature file:**
```gherkin
When I create a task with title "My Task"
Then the task should be in "backlog" state
```

**Step definition:**
```go
package steps

import (
    "context"
    "fmt"
    "github.com/cucumber/godog"
    "github.com/javierbenavides/agentic-agent/internal/tasks"
)

type TaskSteps struct {
    suite *SuiteContext
}

func (s *TaskSteps) RegisterSteps(sc *godog.ScenarioContext) {
    sc.Step(`^I create a task with title "([^"]*)"$`, s.createTaskWithTitle)
    sc.Step(`^the task should be in "([^"]*)" state$`, s.assertTaskInState)
}

func (s *TaskSteps) createTaskWithTitle(ctx context.Context, title string) error {
    tm := tasks.NewTaskManager(s.suite.ProjectDir + "/.agentic/tasks")
    task, err := tm.CreateTask(title)
    if err != nil {
        return fmt.Errorf("failed to create task: %w", err)
    }
    s.suite.CurrentTask = task
    s.suite.LastTaskID = task.ID
    return nil
}

func (s *TaskSteps) assertTaskInState(ctx context.Context, state string) error {
    // Implementation...
    return nil
}
```

### Step Definition Patterns

**Simple string capture:**
```go
sc.Step(`^I run "([^"]*)"$`, s.runCommand)
// Matches: I run "task list"
```

**Number capture:**
```go
sc.Step(`^I should see (\d+) tasks? in backlog$`, s.assertTaskCount)
// Matches: I should see 1 task in backlog
// Matches: I should see 5 tasks in backlog
```

**Table data:**
```go
sc.Step(`^I add the following outputs:$`, s.addOutputs)

func (s *TaskSteps) addOutputs(ctx context.Context, table *godog.Table) error {
    for _, row := range table.Rows[1:] { // Skip header
        output := row.Cells[0].Value
        // Process output...
    }
    return nil
}
```

**Doc string:**
```go
sc.Step(`^I run the following command:$`, s.runMultilineCommand)

func (s *CommonSteps) runMultilineCommand(ctx context.Context, docString *godog.DocString) error {
    cmdText := docString.Content
    // Execute command...
    return nil
}
```

## Running BDD Tests

### Run All BDD Tests

```bash
make test-bdd
```

### Run Specific Features

```bash
# Run only workflow features
go test ./test/bdd -v -godog.paths=../../features/workflows/

# Run specific feature file
go test ./test/bdd -v -godog.paths=../../features/tasks/task_lifecycle.feature
```

### Run Tagged Scenarios

```bash
# Run only @smoke tests
go test ./test/bdd -v -godog.tags=@smoke

# Run @smoke but exclude @wip (work in progress)
go test ./test/bdd -v -godog.tags="@smoke && ~@wip"

# Run beginner or intermediate workflows
go test ./test/bdd -v -godog.tags="@beginner || @intermediate"
```

### Run All Tests

```bash
# Unit + Functional + BDD
make test-all
```

### Coverage Reports

```bash
# Run all tests with coverage
make coverage-all

# Opens merged HTML coverage report showing:
# - Unit test coverage
# - Functional test coverage
# - BDD test coverage
```

## ATDD Workflow

Follow this workflow when developing new features:

### 1. Write Feature First (Specification)

Create a feature file before any implementation:

```gherkin
# features/tasks/task_archiving.feature
Feature: Task Archiving
  As a developer with many completed tasks
  I want to archive old tasks
  So that my task lists remain manageable

  Scenario: Archive tasks older than 30 days
    Given I have completed tasks older than 30 days
    When I run "task archive --older-than=30d"
    Then tasks should be moved to archived.yaml
    And the backlog should not contain archived tasks
```

### 2. Run Tests (They Fail)

```bash
make test-bdd
```

Output:
```
Step definition missing for:
  I have completed tasks older than 30 days
```

### 3. Implement Step Definitions

```go
// test/bdd/steps/task_steps.go
func (s *TaskSteps) RegisterSteps(sc *godog.ScenarioContext) {
    sc.Step(`^I have completed tasks older than (\d+) days$`, s.createOldCompletedTasks)
    sc.Step(`^tasks should be moved to archived\.yaml$`, s.assertTasksArchived)
}

func (s *TaskSteps) createOldCompletedTasks(ctx context.Context, days int) error {
    // Create test data...
    return nil
}

func (s *TaskSteps) assertTasksArchived(ctx context.Context) error {
    // Verify archiving...
    return nil
}
```

### 4. Implement Feature

Implement the actual CLI command:

```go
// cmd/agentic-agent/task.go
var archiveCmd = &cobra.Command{
    Use:   "archive",
    Short: "Archive old completed tasks",
    Run:   archiveTask,
}

func archiveTask(cmd *cobra.Command, args []string) {
    // Implementation...
}
```

### 5. Tests Pass

```bash
make test-bdd
```

Output:
```
Feature: Task Archiving
  Scenario: Archive tasks older than 30 days ✓

1 scenario (1 passed)
4 steps (4 passed)
```

### 6. Refactor with Confidence

Refactor knowing tests will catch regressions:

```bash
make test-all  # All tests still pass
```

### 7. Commit with Living Documentation

```bash
git add features/tasks/task_archiving.feature
git add cmd/agentic-agent/task.go test/bdd/steps/task_steps.go
git commit -m "feat: Add task archiving

Implemented task archiving feature as specified
in task_archiving.feature. Tasks older than a
specified number of days can now be archived."
```

## Best Practices

### 1. **One Scenario, One Behavior**

**Good:**
```gherkin
Scenario: Create task successfully
  When I create a task with title "My Task"
  Then the task should appear in the backlog
```

**Bad:**
```gherkin
Scenario: Create and claim and complete task
  When I create a task with title "My Task"
  And I claim the task
  And I complete the task
  Then the task should be done
```

### 2. **Use Background for Common Setup**

**Good:**
```gherkin
Background:
  Given a clean test environment
  And I have initialized a project

Scenario: Create task
  When I create a task...

Scenario: List tasks
  When I list all tasks...
```

**Bad:**
```gherkin
Scenario: Create task
  Given a clean test environment
  And I have initialized a project
  When I create a task...

Scenario: List tasks
  Given a clean test environment
  And I have initialized a project
  When I list all tasks...
```

### 3. **Descriptive Scenario Names**

**Good:**
```gherkin
Scenario: Cannot claim task that doesn't exist
Scenario: Task preserves acceptance criteria after decomposition
```

**Bad:**
```gherkin
Scenario: Test error
Scenario: Check task
```

### 4. **Use Tags Meaningfully**

```gherkin
@smoke          # Critical paths
@wip            # Work in progress (exclude from CI)
@slow           # Long-running tests
@integration    # Integration tests
@unit           # Unit-level BDD tests
```

### 5. **Keep Steps Reusable**

Create generic step definitions:

```go
// Good - reusable
sc.Step(`^the command should succeed$`, s.commandShouldSucceed)
sc.Step(`^the command should fail$`, s.commandShouldFail)

// Bad - too specific
sc.Step(`^the task create command should succeed$`, s.taskCreateShouldSucceed)
sc.Step(`^the task claim command should fail$`, s.taskClaimShouldFail)
```

### 6. **Avoid Implementation Details**

**Good (behavior-focused):**
```gherkin
When I create a task with title "Build API"
Then the task should be claimable
```

**Bad (implementation-focused):**
```gherkin
When I write "Build API" to backlog.yaml
Then the YAML file should contain the task ID
```

### 7. **Use Examples for Multiple Cases**

```gherkin
Scenario Outline: Validate task title length
  When I create a task with title "<title>"
  Then the command should <result>

  Examples:
    | title                           | result  |
    | Good title                      | succeed |
    |                                 | fail    |
    | Very long title that exceeds... | fail    |
```

## Examples

### Example 1: Simple CRUD Scenario

```gherkin
Feature: Task CRUD Operations

  Scenario: Complete task lifecycle
    Given a clean test environment
    And I have initialized a project

    # Create
    When I create a task with title "Implement feature"
    Then the task should appear in the backlog

    # Read/List
    When I list all tasks
    Then I should see 1 task in backlog

    # Update
    When I claim the task
    Then the task should move to in-progress

    # Delete (complete)
    When I complete the task
    Then the task should move to done
```

### Example 2: Error Handling

```gherkin
Feature: Validation and Error Handling

  Scenario: Prevent claiming nonexistent task
    Given a clean test environment
    When I try to claim task "TASK-NONEXISTENT"
    Then the command should fail
    And the error message should contain "not found"

  Scenario: Require task title
    When I try to create a task with empty title
    Then the command should fail
    And the error message should contain "title"
```

### Example 3: Complex Workflow

```gherkin
Feature: Task Decomposition

  @advanced
  Scenario: Decompose large task into subtasks
    Given I have created a task with title "Build User Management"

    When I decompose the task into the following subtasks:
      | subtask                          |
      | Create user model                |
      | Implement registration endpoint  |
      | Add authentication               |

    Then the task should have 3 subtasks
    And the task should remain in backlog

    When I work on subtask 1
    Then subtask 1 should be marked complete
    And the parent task should still be in backlog
```

### Example 4: Data-Driven Testing

```gherkin
Feature: Task Metadata Validation

  Scenario Outline: Validate task metadata fields
    Given I create a task with title "Test Task"
    When I set <field> to "<value>"
    Then validation should <result>

    Examples:
      | field       | value          | result  |
      | description | Valid text     | pass    |
      | description |                | pass    |
      | priority    | high           | pass    |
      | priority    | invalid        | fail    |
      | due_date    | 2024-12-31     | pass    |
      | due_date    | invalid-date   | fail    |
```

## CI/CD Integration

The BDD tests run automatically on:
- Every push to main/develop branches
- Every pull request
- Multiple Go versions (1.22, 1.23)

See [`.github/workflows/bdd-tests.yml`](.github/workflows/bdd-tests.yml) for configuration.

### Coverage Reports

After each CI run:
- Coverage reports are generated
- Artifacts are uploaded (30-day retention)
- PR comments show coverage percentage
- Test summaries appear in GitHub Actions

## Troubleshooting

### Undefined Step

**Error:**
```
Step is undefined:
  I do something new
```

**Solution:**
Add step definition in appropriate `test/bdd/steps/*.go` file.

### Race Conditions

**Error:**
Tests fail intermittently with file system errors.

**Solution:**
We run tests sequentially (`Concurrency: 1` in `features_test.go`) to avoid this.

### Test Cleanup

**Error:**
Test leaves files in temp directory.

**Solution:**
Ensure cleanup functions are registered:
```go
s.suite.RegisterCleanup(func() {
    os.RemoveAll(s.suite.ProjectDir)
})
```

## Additional Resources

- [Cucumber Documentation](https://cucumber.io/docs/guides/)
- [Godog GitHub](https://github.com/cucumber/godog)
- [Gherkin Reference](https://cucumber.io/docs/gherkin/reference/)
- [BDD Best Practices](https://cucumber.io/docs/bdd/)

---

**Questions?** Check [test/bdd/README.md](test/bdd/README.md) or review existing feature files in `features/`.
