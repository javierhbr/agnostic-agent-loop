# BDD Quick Reference

Quick reference card for working with BDD/ATDD tests in agentic-agent.

## Running Tests

```bash
# Run all BDD tests
make test-bdd

# Run all tests (unit + functional + BDD)
make test-all

# Run specific feature
go test ./test/bdd -v -godog.paths=../../features/tasks/task_lifecycle.feature

# Run specific workflow
go test ./test/bdd -v -godog.paths=../../features/workflows/

# Run with tags
go test ./test/bdd -v -godog.tags=@smoke
go test ./test/bdd -v -godog.tags="@beginner || @intermediate"
go test ./test/bdd -v -godog.tags="~@wip"  # Exclude WIP
```

## Coverage

```bash
# All tests with merged coverage
make coverage-all

# View merged coverage report
open build/coverage/merged-coverage.html

# BDD coverage only
go test ./test/bdd -v -coverprofile=build/build/coverage/bdd.out
go tool cover -html=build/build/coverage/bdd.out
```

## Gherkin Syntax

### Feature Structure

```gherkin
Feature: Short description
  As a <role>
  I want <goal>
  So that <benefit>

  Background:
    Given common setup step
    And another setup step

  @tag1 @tag2
  Scenario: Specific behavior
    Given initial context
    When action occurs
    Then expected outcome
    And additional assertion
```

### Data Tables

```gherkin
When I create the following tasks:
  | title      | priority |
  | Task 1     | high     |
  | Task 2     | low      |
```

### Doc Strings

```gherkin
When I run the following command:
  """
  task create \
    --title "My Task" \
    --description "Details"
  """
```

### Scenario Outlines

```gherkin
Scenario Outline: Validate input
  When I create task with title "<title>"
  Then the result should be "<result>"

  Examples:
    | title        | result  |
    | Valid        | success |
    | <empty>      | error   |
```

## Common Step Patterns

### Setup Steps

```gherkin
Given a clean test environment
And I have initialized a project
And I have initialized project "MyProject"
```

### Action Steps

```gherkin
When I run "task create MyTask"
When I create a task with title "MyTask"
When I claim the task
When I complete the task
When I list all tasks
```

### Assertion Steps

```gherkin
Then the command should succeed
Then the command should fail
Then the error message should contain "not found"
Then the task should be in "backlog" state
Then I should see 2 tasks in backlog
Then the following files should exist:
  | file                  |
  | .agentic/tasks/*.yaml |
```

## Available Tags

- `@smoke` - Critical path tests
- `@wip` - Work in progress (excluded from CI)
- `@slow` - Long-running tests
- `@beginner` - Beginner tutorial
- `@intermediate` - Intermediate tutorial
- `@advanced` - Advanced tutorial

## Step Definition Locations

- **Setup & Commands**: `test/bdd/steps/common_steps.go`
- **Task Operations**: `test/bdd/steps/task_steps.go`
- **Assertions**: `test/bdd/steps/assertion_steps.go`
- **Initialization**: `test/bdd/steps/init_steps.go`
- **Skill Generation**: `test/bdd/steps/skill_steps.go`
- **Agent Detection**: `test/bdd/steps/detection_steps.go`

## ATDD Workflow

1. **Write Feature**
   ```bash
   vim features/tasks/new_feature.feature
   ```

2. **Run Tests** (they fail with undefined steps)
   ```bash
   make test-bdd
   ```

3. **Implement Step Definitions**
   ```bash
   vim test/bdd/steps/task_steps.go
   ```

4. **Implement Feature**
   ```bash
   vim cmd/agentic-agent/task.go
   ```

5. **Tests Pass**
   ```bash
   make test-bdd  # Green!
   ```

6. **Commit**
   ```bash
   git add features/ test/bdd/ cmd/
   git commit -m "feat: Add new feature"
   ```

## Feature File Locations

- `features/init/` - Project initialization
- `features/tasks/` - Task management
- `features/context/` - Context generation
- `features/skills/` - Agent detection and skill generation
- `features/workflows/` - Complete workflows

## Useful Commands

```bash
# List all feature files
find features -name "*.feature"

# Count scenarios
grep -r "Scenario:" features/ | wc -l

# Count steps
grep -r "Given\|When\|Then\|And" features/ | wc -l

# View feature
cat features/workflows/beginner_workflow.feature

# Check step definition usage
grep -r "sc.Step" test/bdd/steps/
```

## CI/CD

### GitHub Actions Workflows

- `.github/workflows/bdd-tests.yml` - BDD-specific tests
- `.github/workflows/test-suite.yml` - Complete suite with matrix

### What Runs on CI

- All BDD tests on every push/PR
- Multi-version Go testing (1.22, 1.23)
- Coverage report generation
- PR comments with coverage %
- Test artifact uploads

## Debugging

### Verbose Output

```bash
# More detail
go test ./test/bdd -v

# Even more detail
go test ./test/bdd -v -godog.format=pretty -test.v
```

### Run Single Scenario

Edit feature file, add `@focus` tag:
```gherkin
@focus
Scenario: Debug this one
  ...
```

Then:
```bash
go test ./test/bdd -v -godog.tags=@focus
```

### Check Step Definitions

```bash
# List all registered steps
go test ./test/bdd -v -godog.format=steps
```

## Best Practices

✅ **DO**:
- Write features before code
- Keep scenarios focused (one behavior)
- Use descriptive scenario names
- Reuse step definitions
- Tag scenarios meaningfully
- Use Background for common setup

❌ **DON'T**:
- Mix multiple behaviors in one scenario
- Use implementation details in steps
- Create duplicate step definitions
- Skip the feature-first workflow
- Commit @wip tagged scenarios

## Common Issues

### Undefined Step

**Error**: `Step is undefined: I do something new`

**Fix**: Add step definition in appropriate file under `test/bdd/steps/`

### Race Condition

**Error**: Tests fail intermittently

**Fix**: Already handled - we use `Concurrency: 1`

### Test Cleanup

**Error**: Temp files remain after test

**Fix**: Register cleanup in suite context:
```go
s.suite.RegisterCleanup(func() {
    os.RemoveAll(tempDir)
})
```

## Help & Documentation

- [BDD Guide](BDD_GUIDE.md) - Complete guide
- [CLI Tutorial](CLI_TUTORIAL.md) - Tutorial with ATDD section
- [Implementation Summary](BDD_IMPLEMENTATION_SUMMARY.md) - Project overview
- [test/bdd/README.md](../test/bdd/README.md) - Infrastructure docs

## Quick Stats

- **Feature Files**: 9
- **Scenarios**: 43
- **Steps**: 255
- **Step Definition Files**: 7
- **Step Definitions**: 70+
- **Pass Rate**: 100%
- **Execution Time**: ~1s

---

**Last Updated**: February 2026
**godog Version**: v0.15.1
