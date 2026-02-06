# Agentic Agent Framework

A specification-driven framework for building agent-agnostic AI workflows with strong context isolation and task management.

## Overview

The Agentic Agent framework provides a structured approach to AI-assisted development by:

- **Specification-Driven Development**: Agents reference specifications, never serve as source of truth
- **Context Isolation**: Per-directory context files keep agents focused on relevant information
- **Task Management**: Structured workflow for creating, claiming, and completing atomic tasks
- **Agent-Agnostic Design**: Works with any AI agent tool (Claude Code, Cursor, GitHub Copilot, etc.)
- **Validation Rules**: Enforces best practices for context management and task constraints

## Features

- ✅ **Interactive CLI** - Firebase-inspired wizards for junior developers
- ✅ Project initialization with structured directory layout
- ✅ Task lifecycle management (backlog → in-progress → done)
- ✅ Context generation and management
- ✅ Validation rules for quality assurance
- ✅ Task size constraints (max 5 files, 2 directories per task)
- ✅ Specification references and acceptance criteria
- ✅ Subtask decomposition
- ✅ Skills generation for AI agent tools
- ✅ Dual-mode support (interactive wizards + traditional flags)
- ✅ **Ralph PDR Integration** - Plan-Do-Review methodology with PRD generation
- ✅ **Progress Tracking** - Dual-format progress logs (text + YAML)
- ✅ **Learnings Management** - Codebase patterns and directory-specific guidance
- ✅ **Browser Verification** - Validation for UI changes
- ✅ **agentskills.io Compliance** - Portable skills across tools
- ✅ **BDD/ATDD Testing** - Gherkin feature files with godog framework for living documentation

## Installation

### Prerequisites

- Go 1.22 or higher
- Git (for task scope validation)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/javierbenavides/agentic-agent.git
cd agentic-agent

# Build the binary
go build -o agentic-agent ./cmd/agentic-agent

# (Optional) Install globally
sudo mv agentic-agent /usr/local/bin/

# Verify installation
agentic-agent version
```

### Using Go Install

```bash
go install github.com/javierbenavides/agentic-agent/cmd/agentic-agent@latest
```

### Development & Testing Without Installation

If you want to test the binary without installing it globally, you have several options:

#### Option 1: Temporary Alias (Current Session Only)

```bash
# Build the binary
go build -o agentic-agent ./cmd/agentic-agent

# Create alias for current session
alias agentic-agent='./agentic-agent'

# Use it
agentic-agent version
agentic-agent task list
```

#### Option 2: Persistent Alias (Recommended for Development)

```bash
# Build the binary
go build -o agentic-agent ./cmd/agentic-agent

# Add to your shell config (~/.bashrc or ~/.zshrc)
echo "alias agentic-agent='$(pwd)/agentic-agent'" >> ~/.zshrc

# Reload shell config
source ~/.zshrc

# Now the alias persists across sessions
agentic-agent version
```

#### Option 3: Shell Function (More Flexible)

```bash
# Add to ~/.bashrc or ~/.zshrc
agentic-agent() {
    /full/path/to/agnostic-agent-loop/agentic-agent "$@"
}

# Or export for use
export -f agentic-agent

# Reload config
source ~/.zshrc
```

#### Option 4: Add to PATH

```bash
# Build the binary
go build -o agentic-agent ./cmd/agentic-agent

# Add project directory to PATH (temporary)
export PATH="$PATH:$(pwd)"

# Or add permanently to ~/.zshrc or ~/.bashrc
echo "export PATH=\"\$PATH:$(pwd)\"" >> ~/.zshrc
source ~/.zshrc

# Use from anywhere
agentic-agent version
```

#### Option 5: Use Go Run (No Build Needed)

```bash
# Run directly without building
go run ./cmd/agentic-agent version
go run ./cmd/agentic-agent task list
go run ./cmd/agentic-agent init --name "Test"

# Create an alias for convenience
alias agentic-agent='go run ./cmd/agentic-agent'
```

#### Option 6: Install to GOPATH/bin

```bash
# Install to $GOPATH/bin (usually already in PATH)
go install ./cmd/agentic-agent

# If $GOPATH/bin is in your PATH, this works immediately
agentic-agent version

# Check if GOPATH/bin is in PATH
echo $PATH | grep "$(go env GOPATH)/bin"

# If not, add it
echo "export PATH=\"\$PATH:$(go env GOPATH)/bin\"" >> ~/.zshrc
source ~/.zshrc
```

**Recommended for local development**: Use Option 2 (persistent alias) or Option 6 (go install) for the best experience.

## Quick Start

### 1. Initialize a New Project

#### Option A: Interactive Wizard (Recommended for Beginners)

```bash
# Create and navigate to your project directory
mkdir my-project
cd my-project

# Launch the interactive setup wizard
agentic-agent start
```

This launches a friendly, step-by-step wizard that guides you through:

- Project naming and validation
- AI model selection (Claude, GPT-4, etc.)
- Directory structure setup
- Initial configuration

## 🎯 Dual-Mode Support

**All commands support both interactive and flag-based modes!**

### Interactive Mode (Beginner-Friendly)

Run any command without arguments to get a guided experience:

```bash
agentic-agent task claim          # Shows task list to select from
agentic-agent context generate    # Shows directory picker
agentic-agent task create         # Step-by-step wizard
```

### Flag Mode (Power User / Automation)

Use traditional flags for scripts and automation:

```bash
agentic-agent task claim abc123
agentic-agent context generate ./src
agentic-agent task create --title "Fix bug" --description "Details"
```

**📖 See [COMMAND_QUICK_REFERENCE.md](COMMAND_QUICK_REFERENCE.md) for complete command reference**

---

#### Option B: Traditional Command Line

```bash
# Create and navigate to your project directory
mkdir my-project
cd my-project

# Initialize the agentic framework
agentic-agent init --name "My Project"
```

Both options create the following structure:

```
.agentic/
├── spec/              # Specifications and architecture docs
├── context/           # Global and rolling context summaries
├── tasks/             # Task management files
│   ├── backlog.yaml
│   ├── in-progress.yaml
│   └── done.yaml
└── agent-rules/       # Tool-specific agent configurations
agnostic-agent.yaml    # Project configuration
```

### 2. Create Your First Task

#### Option A: Interactive Task Creation (Recommended)

```bash
# Launch interactive task wizard
agentic-agent task create
```

The wizard guides you through:

- Title input with validation
- Description (multi-line, optional)
- Acceptance criteria editor (press 'a' to add items)
- Preview before creation
- Success confirmation with next steps

#### Option B: Quick Sample Task (Perfect for Testing)

```bash
# Create a pre-configured sample task instantly
agentic-agent task sample-task
```

Creates a sample task with:

- Example title and description
- Pre-filled acceptance criteria
- Suggested scope

Perfect for testing the workflow or learning the system!

#### Option C: Create from Template (Best Practices Built-in)

```bash
# Launch template selection wizard
agentic-agent task from-template
```

Choose from professional templates:

- **Feature Implementation** - For adding new features
- **Bug Fix** - Structured bug fix workflow
- **Refactoring** - Code quality improvements
- **Documentation** - Documentation tasks
- **Testing** - Test coverage tasks

Each template includes:

- Pre-configured acceptance criteria
- Suggested structure
- Best practices

#### Option D: Flag-Based Task Creation

```bash
# Create a simple task
agentic-agent task create --title "Implement user authentication"

# Create a task with all fields
agentic-agent task create \
  --title "Implement JWT authentication" \
  --description "Add JWT-based authentication with token validation" \
  --spec-refs ".agentic/spec/04-architecture.md,.agentic/spec/05-domain-model.md" \
  --inputs ".agentic/context/rolling-summary.md" \
  --outputs "src/auth/jwt.go,tests/auth_test.go" \
  --acceptance "JWT tokens generated,Token validation works,All tests pass"
```

### 3. Work on a Task

```bash
# List all tasks
agentic-agent task list

# View task details
agentic-agent task show TASK-001

# Claim a task (marks it as in-progress)
agentic-agent task claim TASK-001

# Generate context for the directory you'll be working in
agentic-agent context generate src/auth

# Do your work...

# Complete the task
agentic-agent task complete TASK-001
```

### 4. Validate Your Work

```bash
# Run all validation rules
agentic-agent validate

# Get JSON output
agentic-agent validate --format json
```

## Interactive CLI Mode

The Agentic Agent CLI offers two modes of operation:

### Interactive Mode (Default for Commands Without Flags)

When you run commands without any flags, the CLI automatically launches interactive wizards:

```bash
# These launch interactive wizards
agentic-agent start          # Project setup wizard
agentic-agent init          # Interactive initialization
agentic-agent task create    # Task creation with file pickers
agentic-agent task list      # Task management with tabs
agentic-agent work          # Complete workflow (claim → work → complete)
```

**Features:**

- **Guided Wizards** - Step-by-step guidance with real-time validation
- **File Pickers** - Browse and select files/directories with multi-select
- **Task Management** - Tabbed interface (Backlog/In Progress/Done) with quick actions
- **Progress Tracking** - Interactive checklists for acceptance criteria
- **Complete Workflows** - End-to-end task completion flow
- **AI Model Selection** - Choose models with descriptions
- **Keyboard Navigation** - Full keyboard support (↑/↓/j/k/Tab/Space/Enter/Esc)
- **Beautiful UI** - Colors, icons, and styled components
- **Perfect for Junior Developers** - No need to memorize commands or flags

### Flag Mode (Traditional CLI)

When you provide flags, the CLI uses traditional command-line mode:

```bash
# These use traditional flag-based mode
agentic-agent init --name "My Project"
agentic-agent task create --title "My Task" --description "Details"
```

**Features:**

- Fast, scriptable commands
- CI/CD friendly
- Backward compatible with existing workflows
- Power user efficiency

### Forcing Non-Interactive Mode

Use the `--no-interactive` flag to force flag-based mode even when no flags are provided:

```bash
agentic-agent task create --no-interactive --title "Automated Task"
```

This is useful in scripts or automation where you want to ensure non-interactive behavior.

---

## Usage Guide

### Task Management

#### Creating Tasks

Tasks are the fundamental unit of work. Each task should be atomic and focused.

```bash
# Basic task
agentic-agent task create --title "Fix login bug"

# Task with full metadata
agentic-agent task create \
  --title "Implement password reset" \
  --description "Add email-based password reset flow" \
  --spec-refs ".agentic/spec/auth.md" \
  --inputs ".agentic/context/rolling-summary.md" \
  --outputs "src/auth/reset.go,src/templates/reset.html" \
  --acceptance "Reset email sent,Token validation works,Password updated"
```

#### Task Constraints

The framework enforces task size limits:
- **Maximum 5 files** per task
- **Maximum 2 directories** per task

If your task is larger, decompose it into subtasks:

```bash
# Decompose a large task
agentic-agent task decompose TASK-001 \
  "Create JWT service" \
  "Add middleware" \
  "Write tests"
```

#### Task Workflow

```bash
# 1. List available tasks
agentic-agent task list

# 2. View task details
agentic-agent task show TASK-001

# 3. Claim the task
agentic-agent task claim TASK-001

# 4. Work on the task (generate context, write code, etc.)
agentic-agent context generate src/module

# 5. Complete the task
agentic-agent task complete TASK-001

# 6. Validate your work
agentic-agent validate
```

### Context Management

Context files provide focused information to AI agents about specific parts of your codebase.

#### Generating Context

```bash
# Generate context for a directory
agentic-agent context generate src/auth

# Scan for directories missing context
agentic-agent context scan
```

#### Context Structure

Each `context.md` file should include:

```markdown
# Module Context

## Purpose
Brief description of what this module does

## Responsibilities
- Key responsibility 1
- Key responsibility 2

## Dependencies
- External dependency 1
- Internal dependency 2

## Must Do
- Constraint or requirement 1
- Constraint or requirement 2

## Cannot Do
- Prohibition 1
- Prohibition 2
```

### Validation Rules

The framework includes several validation rules:

1. **Directory Context Required**: Every directory with source code must have a `context.md` file
2. **Context Update on Change**: `context.md` must be updated when source files change
3. **Task Scope Enforcement**: Modified files must be within in-progress task scope
4. **Task Size Limits**: Tasks cannot exceed 5 files or 2 directories

```bash
# Run all validations
agentic-agent validate

# View validation results in JSON
agentic-agent validate --format json
```

### Skills Generation

Generate tool-specific configuration files for AI agents:

```bash
# Generate Claude Code skill file
agentic-agent skills generate --tool claude-code

# Generate for other tools (coming soon)
agentic-agent skills generate --tool cursor
agentic-agent skills generate --tool copilot
```

## Project Structure

### .agentic/spec/

Specification files that serve as the source of truth:

- `01-overview.md` - Project overview
- `02-goals.md` - Project goals and objectives
- `03-constraints.md` - Technical constraints
- `04-architecture.md` - System architecture
- `05-domain-model.md` - Domain models and entities
- `06-data-flow.md` - Data flow and integrations

### .agentic/context/

Context summaries for agents:

- `global-context.md` - Immutable project-wide context
- `rolling-summary.md` - Updated after each session with key changes

### .agentic/tasks/

Task management files (YAML format):

- `backlog.yaml` - Pending tasks
- `in-progress.yaml` - Tasks currently being worked on
- `done.yaml` - Completed tasks

### Source Directories

Each source directory should have:

```
src/module/
├── context.md        # Module-specific context
├── *.go              # Source files
└── *_test.go         # Test files
```

## Configuration

### How the CLI Finds `agnostic-agent.yaml`

The CLI expects `agnostic-agent.yaml` to exist in the **current working directory** where you run the command. It does **not** search parent directories or perform any file discovery — it simply reads `./agnostic-agent.yaml` relative to your shell's working directory.

This file is created automatically when you run `agentic-agent init`. Several commands (`learnings`, `skills generate`) load this file for path configuration. If the file is not found, those commands fall back to sensible defaults.

You can also pass a custom path via the `--config` flag:

```bash
agentic-agent --config /path/to/my-config.yaml task list
```

### Configuration Reference

The `agnostic-agent.yaml` file in your project root contains project configuration:

```yaml
project:
  name: "My Project"
  version: 0.1.0
  roots:
    - .

agents:
  defaults:
    max_tokens: 4000
    model: claude-3-5-sonnet-20241022

# Path configuration (used by learnings, skills, archiver)
paths:
  prdOutputPath: .agentic/tasks/
  progressTextPath: .agentic/progress.txt
  progressYAMLPath: .agentic/progress.yaml
  archiveDir: .agentic/archive/

workflow:
  validators:
    - context-check
    - task-scope
    - browser-verification
```

## Best Practices

### 1. Keep Tasks Atomic

Break large features into small, focused tasks:

```bash
# Bad: One large task
"Implement entire user management system"

# Good: Multiple focused tasks
"Create user model and repository"
"Implement user registration endpoint"
"Add user authentication"
"Write user management tests"
```

### 2. Maintain Context Files

Update `context.md` whenever you change a module's logic:

```bash
# After making changes
agentic-agent context generate src/auth

# Review and edit the generated context
vim src/auth/context.md
```

### 3. Reference Specifications

Always link tasks to relevant specifications:

```bash
agentic-agent task create \
  --title "Implement rate limiting" \
  --spec-refs ".agentic/spec/03-constraints.md,.agentic/spec/04-architecture.md"
```

### 4. Use Acceptance Criteria

Define clear success criteria for each task:

```bash
agentic-agent task create \
  --title "Add email validation" \
  --acceptance "Valid emails accepted,Invalid emails rejected,Error messages shown,Tests cover edge cases"
```

### 5. Run Validation Regularly

```bash
# Before committing
agentic-agent validate

# In CI/CD pipeline
agentic-agent validate --format json
```

## Examples

### Example 1: Starting a New Feature

```bash
# 1. Create the feature task
agentic-agent task create \
  --title "Add password reset feature" \
  --spec-refs ".agentic/spec/auth.md"

# 2. Decompose into subtasks
agentic-agent task decompose TASK-001 \
  "Create password reset request endpoint" \
  "Add email notification service" \
  "Create reset token validation" \
  "Add password update endpoint" \
  "Write integration tests"

# 3. Claim first subtask
agentic-agent task claim TASK-001.1

# 4. Generate context for relevant directories
agentic-agent context generate src/auth
agentic-agent context generate src/notifications

# 5. Work on the subtask...

# 6. Complete and validate
agentic-agent task complete TASK-001.1
agentic-agent validate
```

### Example 2: Bug Fix Workflow

```bash
# 1. Create bug fix task
agentic-agent task create \
  --title "Fix null pointer in login handler" \
  --outputs "src/auth/login.go,tests/auth/login_test.go" \
  --acceptance "NPE fixed,Edge case covered,Test added"

# 2. Claim the task
agentic-agent task claim TASK-042

# 3. Read relevant context
cat src/auth/context.md

# 4. Fix the bug and add tests...

# 5. Update context if logic changed
agentic-agent context generate src/auth

# 6. Validate and complete
agentic-agent validate
agentic-agent task complete TASK-042
```

## Testing

The framework includes comprehensive tests covering all functionality including the Ralph PDR integration.

### Run All Tests

```bash
# Run all tests in the project
go test ./...

# Run all tests with verbose output
go test ./... -v

# Run all tests with coverage report
go test ./... -cover
```

### Run Tests by Package

```bash
# Task management tests (including progress writer)
go test ./internal/tasks -v

# Validator tests (including browser verification)
go test ./internal/validator/rules -v

# Orchestrator tests (loop and archiver)
go test ./internal/orchestrator -v

# Integration tests
go test ./test/integration -v

# Model tests
go test ./pkg/models -v
```

### Run Specific Tests

```bash
# Run a specific test function
go test ./internal/tasks -run TestProgressWriter_AppendEntry

# Run all tests matching a pattern
go test ./internal/validator/rules -run TestBrowserVerification

# Run tests with detailed output
go test ./internal/orchestrator -run TestArchiver -v
```

### View Test Coverage

#### Quick Coverage Check

```bash
# Show coverage by package
go test ./... -cover
```

#### Comprehensive Coverage Reports

We provide multiple ways to generate detailed coverage reports:

**Option 1: Using Make (Recommended)**

```bash
# Generate HTML coverage report and open in browser
make coverage-html

# Show coverage by function
make coverage-func

# Show coverage summary by package
make coverage-summary

# Clean coverage files
make clean-coverage
```

**Option 2: Using Coverage Script**

```bash
# Generate comprehensive coverage report
./scripts/coverage-report.sh

# Generate report and open HTML in browser
./scripts/coverage-report.sh --open
```

**Option 3: Manual Go Commands**

```bash
# Generate detailed coverage report
go test ./... -coverprofile=build/coverage/coverage.out -covermode=count
go tool cover -html=build/coverage/coverage.out -o coverage/coverage.html

# View coverage by function
go tool cover -func=build/coverage/coverage.out

# Count passing tests
go test ./... -v 2>&1 | grep -c "^--- PASS:"
```

#### Coverage Report Features

The coverage report tools provide:
- **Total coverage percentage** with color-coded output
- **Package-by-package breakdown** sorted by coverage
- **Packages without tests** highlighted
- **HTML visualization** of line-by-line coverage
- **Function-level coverage** details
- **Coverage badge** for documentation
- **Threshold validation** (fails CI if below 50%)

#### Understanding Coverage Output

Coverage levels are color-coded:
- 🟢 **Green (70%+)**: Good coverage
- 🟡 **Yellow (40-69%)**: Needs improvement
- 🔴 **Red (<40%)**: Insufficient coverage

### Ralph PDR Integration Tests

The Ralph PDR integration includes 43 new tests:

```bash
# Progress writer tests (11 tests)
go test ./internal/tasks -run TestProgressWriter -v

# Browser verification validator tests (13 tests)
go test ./internal/validator/rules -run TestBrowserVerification -v

# Archiver tests (10 tests)
go test ./internal/orchestrator -run TestArchiver -v

# Loop orchestrator tests (9 tests)
go test ./internal/orchestrator -run TestLoop -v
```

### Test Results

The framework has been fully validated with:
- **87 unit and integration tests** (total)
- **43 tests** for Ralph PDR integration
- **44 tests** for core framework
- All tests passing ✅
- Build compiles successfully ✅

### Quick Verification

```bash
# Verify all Ralph PDR tests pass
go test ./internal/tasks ./internal/validator/rules ./internal/orchestrator -v

# Run full test suite and show summary
go test ./... -v 2>&1 | tail -20
```

See [VALIDATION_REPORT.md](VALIDATION_REPORT.md) for detailed test results.

## Development

### Project Structure

```
agentic-agent/
├── cmd/
│   └── agentic-agent/      # CLI commands
├── internal/
│   ├── config/             # Configuration management
│   ├── context/            # Context generation
│   ├── project/            # Project initialization
│   ├── tasks/              # Task management
│   └── validator/          # Validation rules
├── pkg/
│   └── models/             # Data models
└── tests/
    └── integration/        # Integration tests
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests for your changes
4. Ensure all tests pass (`go test ./...`)
5. Run validation (`agentic-agent validate`)
6. Submit a pull request

## Troubleshooting

### Common Issues

**Issue**: `agentic-agent: command not found`

**Solution**: Ensure the binary is in your PATH or use the full path to the binary.

```bash
# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH="$PATH:/path/to/agentic-agent"

# Or use full path
/path/to/agentic-agent version
```

---

**Issue**: Validation fails with "Missing context.md"

**Solution**: Generate context for all source directories:
```bash
# Find directories needing context
agentic-agent context scan

# Generate context for specific directory
agentic-agent context generate src/module
```

---

**Issue**: Task size validation fails

**Solution**: Decompose large tasks into smaller subtasks:
```bash
agentic-agent task decompose TASK-001 \
  "Subtask 1" \
  "Subtask 2" \
  "Subtask 3"
```

---

**Issue**: Git integration not working

**Solution**: Ensure you're in a git repository:
```bash
git init  # If not already a git repo
git add .
git commit -m "Initial commit"
```

---

**Issue**: Task claim fails with "not found in backlog"

**Solution**: Check task status and location:
```bash
# Show task details
agentic-agent task show TASK-001

# List all tasks
agentic-agent task list

# Tasks can only be claimed from backlog
```

## CLI Commands Reference

### Project Commands

```bash
agentic-agent init --name "Project Name"    # Initialize new project
agentic-agent version                       # Show version info
```

### Task Commands

```bash
agentic-agent task create [flags]          # Create new task
agentic-agent task list                    # List all tasks
agentic-agent task show TASK-ID            # Show task details
agentic-agent task claim TASK-ID           # Claim task (move to in-progress)
agentic-agent task complete TASK-ID        # Complete task (move to done)
agentic-agent task decompose TASK-ID ...   # Break task into subtasks
```

Task create flags:
- `--title` - Task title (required)
- `--description` - Detailed description
- `--spec-refs` - Comma-separated spec file references
- `--inputs` - Comma-separated input files
- `--outputs` - Comma-separated output files
- `--acceptance` - Comma-separated acceptance criteria

### Context Commands

```bash
agentic-agent context generate DIR         # Generate context for directory
agentic-agent context scan                 # Find directories missing context
```

### Validation Commands

```bash
agentic-agent validate                     # Run all validation rules
agentic-agent validate --format json       # Output as JSON
```

### Skills Commands

```bash
agentic-agent skills generate --tool TOOL  # Generate tool-specific config
```

Supported tools:
- `claude-code` - Claude Code configuration

## Architecture

The framework follows these principles:

### 1. Specification-Driven Development

Specifications in `.agentic/spec/` are the single source of truth. Agents reference these specs, never serve as the source of truth themselves.

### 2. Context Isolation

Each directory maintains its own `context.md` file with focused information. This prevents agents from being overwhelmed with irrelevant context.

### 3. Atomic Tasks

Tasks are kept small and focused (max 5 files, 2 directories). Large features are decomposed into subtasks.

### 4. Agent-Agnostic Design

The framework works with any AI agent tool through generated configuration files and adapters.

### 5. Validation First

Validation rules enforce best practices and catch issues before they become problems.

## Documentation

### User Guides
- [docs/CLI_TUTORIAL.md](docs/CLI_TUTORIAL.md) - Step-by-step CLI tutorial with scenarios
- [docs/BDD_GUIDE.md](docs/BDD_GUIDE.md) - Complete guide to BDD/ATDD testing
- [examples/multi-agent-workflow/MULTI_AGENT_USE_CASE.md](examples/multi-agent-workflow/MULTI_AGENT_USE_CASE.md) - Switching between Claude Code CLI, VSCode extension, and Copilot

### Technical Documentation
- [VALIDATION_REPORT.md](VALIDATION_REPORT.md) - Detailed validation and test results
- [CLAUDE.md](CLAUDE.md) - Claude-specific agent rules
- [docs/COVERAGE.md](docs/COVERAGE.md) - Comprehensive test coverage guide
- [COVERAGE_QUICK_REF.md](COVERAGE_QUICK_REF.md) - Coverage quick reference

### Specifications
- Specification files in `.agentic/spec/` - Project specifications

### Testing
- [test/bdd/README.md](test/bdd/README.md) - BDD infrastructure overview
- [features/](features/) - Gherkin feature files (executable specifications)

## Roadmap

### Implemented (Phases 1-3)
- ✅ Project initialization
- ✅ Task management (CRUD, lifecycle)
- ✅ Context generation
- ✅ Validation rules
- ✅ Task constraints enforcement
- ✅ Skills generation (Claude Code)
- ✅ Comprehensive test suite

### Planned (Phase 4+)
- 🔲 Enhanced context generation with AST parsing
- 🔲 Multi-language context support (TypeScript, Python, Java)
- 🔲 Token limit enforcement
- 🔲 Session management and summaries
- 🔲 More tool adapters (Cursor, Copilot)
- 🔲 Web UI for task management
- 🔲 Integration with CI/CD pipelines

## License

[Add your license here]

## Support

For issues, questions, or contributions:
- GitHub Issues: [Create an issue](https://github.com/javierbenavides/agentic-agent/issues)
- Email: [your-email@example.com]

## Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [YAML v3](https://github.com/go-yaml/yaml) - YAML parsing
- [testify](https://github.com/stretchr/testify) - Testing framework

---

**Status**: Phases 1-3 complete and validated ✅
**Test Coverage**: 66.4% on critical packages
**Tests**: 59/59 passing
**Ready for**: Phase 4 development
