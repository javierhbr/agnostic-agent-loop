# Agentic Agent CLI Tutorial

A comprehensive, step-by-step guide to using the agentic-agent CLI tool for specification-driven AI-assisted development.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Tutorial Structure](#tutorial-structure)
- [Scenario 1: Beginner - Your First Project](#scenario-1-beginner---your-first-project)
- [Scenario 2: Intermediate - Building a Blog API](#scenario-2-intermediate---building-a-blog-api)
- [Scenario 3: Advanced - Complex Feature Decomposition](#scenario-3-advanced---complex-feature-decomposition)
- [Testing Workflow](#testing-workflow)
- [ATDD/BDD Workflow](#atddbdd-workflow)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)
- [Quick Reference](#quick-reference)

## Overview

This tutorial guides you through three progressive scenarios:

1. **Beginner**: Learn the basics by creating your first project and completing a simple task
2. **Intermediate**: Build a blog API with proper specifications and context management
3. **Advanced**: Decompose a complex feature into subtasks and manage them systematically

Each scenario includes:
- Step-by-step CLI commands
- Expected outputs
- Automated test verification
- Best practices and tips

## Prerequisites

Before starting this tutorial, ensure you have:

- **Go 1.22+** installed
- **Git** installed and configured
- **agentic-agent CLI** built and available in your PATH
- Basic understanding of command-line interfaces
- Familiarity with YAML file format

### Installation Verification

Verify your installation:

```bash
# Check agentic-agent is installed
agentic-agent version

# Expected output:
# Agentic Agent Framework v0.1.0
```

If you get "command not found", see the [README.md](../README.md#installation) for installation instructions.

## Tutorial Structure

### Learning Objectives

By completing this tutorial, you will learn how to:

✅ Initialize new agentic projects
✅ Create and manage tasks with the CLI
✅ Use both interactive and flag-based modes
✅ Generate and maintain context files
✅ Decompose large tasks into manageable subtasks
✅ Run validation rules
✅ Execute and verify tests
✅ Follow specification-driven development practices

### Tutorial Approach

- **Hands-on**: You'll execute actual commands and see real outputs
- **Progressive**: Each scenario builds on the previous one
- **Tested**: All commands are verified by automated functional tests
- **Practical**: Based on real-world development workflows

---

## Scenario 1: Beginner - Your First Project

**Goal**: Learn the basic workflow by creating a project, managing a simple task, and running tests.

**Time**: ~10 minutes

**What You'll Learn**:
- Project initialization
- Task creation and lifecycle
- Basic validation
- Running tests

### Step 1: Create Your Project Directory

```bash
# Create a new directory for your project
mkdir my-first-project
cd my-first-project
```

### Step 2: Initialize the Project

You can initialize in two ways:

#### Option A: Interactive Mode (Recommended for Beginners)

```bash
agentic-agent start
```

This launches an interactive wizard that guides you through:
- Project naming (with validation)
- AI model selection
- Directory structure setup

Follow the prompts and answer the questions.

#### Option B: Flag Mode (Quick Setup)

```bash
agentic-agent init --name "My First Project"
```

**Expected Output**:
```
Initializing project: My First Project
Created agnostic-agent.yaml
Created .agentic/tasks/backlog.yaml
Created .agentic/tasks/in-progress.yaml
Created .agentic/tasks/done.yaml
Created .agentic/context/global-context.md
Created .agentic/context/rolling-summary.md
Created .agentic/context/decisions.md
Created .agentic/context/assumptions.md
Created .agentic/agent-rules/base.md
Project initialized successfully.
```

### Step 3: Explore the Project Structure

```bash
# List the created structure
tree .agentic  # or use ls -R .agentic if tree is not installed

# View the project structure
ls -la
```

**You Should See**:
```
.agentic/
├── tasks/
│   ├── backlog.yaml        # Pending tasks
│   ├── in-progress.yaml    # Active tasks
│   └── done.yaml           # Completed tasks
├── context/
│   ├── global-context.md   # Project-wide context
│   ├── rolling-summary.md  # Session summaries
│   ├── decisions.md        # Design decisions
│   └── assumptions.md      # Project assumptions
├── spec/                   # Specification files (optional)
└── agent-rules/            # Tool-specific agent configs
    └── base.md
agnostic-agent.yaml         # Project configuration
```

### Step 4: Create Your First Task

#### Option A: Interactive Mode

```bash
agentic-agent task create
```

Follow the wizard to:
- Enter a title (e.g., "Set up project README")
- Add a description (optional)
- Add acceptance criteria (press 'a' to add items)

#### Option B: Quick Sample Task

```bash
agentic-agent task sample-task
```

This creates a pre-configured sample task instantly.

#### Option C: Flag Mode

```bash
agentic-agent task create \
  --title "Set up project README" \
  --description "Create a comprehensive README file" \
  --acceptance "README includes project overview,README includes installation steps,README includes usage examples"
```

**Expected Output**:
```
✓ Task created successfully!

ID: TASK-1234567890
Title: Set up project README
Status: pending

Next steps:
1. View details: agentic-agent task show TASK-1234567890
2. Claim task: agentic-agent task claim TASK-1234567890
```

### Step 5: List All Tasks

```bash
agentic-agent task list
```

#### Interactive Mode Output:
Shows a beautiful tabbed interface with:
- **Backlog** tab: Your pending tasks
- **In Progress** tab: Currently active tasks
- **Done** tab: Completed tasks

Use arrow keys or j/k to navigate, Tab to switch tabs, Enter to select, q to quit.

#### Flag Mode Output (with `--no-interactive`):
```
=== Backlog ===
ID: TASK-1234567890
Title: Set up project README
Status: pending

=== In Progress ===
(empty)

=== Done ===
(empty)
```

### Step 6: View Task Details

```bash
# Replace TASK-ID with your actual task ID from Step 4
agentic-agent task show TASK-1234567890
```

**Output**:
```
Task Details
─────────────────────────────────────
ID:          TASK-1234567890
Title:       Set up project README
Status:      pending
Description: Create a comprehensive README file

Acceptance Criteria:
 ✓ README includes project overview
 ✓ README includes installation steps
 ✓ README includes usage examples
```

### Step 7: Claim the Task

```bash
agentic-agent task claim TASK-1234567890
```

**Output**:
```
✓ Task claimed successfully!

Task TASK-1234567890 moved from backlog to in-progress
Assigned to: [your system username]
```

**What Happened**:
- Task moved from `backlog.yaml` to `in-progress.yaml`
- Status changed from `pending` to `in-progress`
- Task assigned to you

### Step 8: Do the Work

In a real scenario, you would now:
1. Create the README.md file
2. Write the content
3. Test it

For this tutorial, let's create a simple README:

```bash
cat > README.md << 'EOF'
# My First Project

## Overview
This is my first project using the agentic-agent framework.

## Installation
See the main README for installation instructions.

## Usage
Run `agentic-agent --help` to see available commands.
EOF
```

### Step 9: Complete the Task

```bash
agentic-agent task complete TASK-1234567890
```

**Output**:
```
✓ Task completed successfully!

Task TASK-1234567890 moved from in-progress to done
Status: done
```

### Step 10: Validate Your Work

```bash
agentic-agent validate
```

**Output**:
```
Running validation rules...

✓ directory-context: PASS
✓ task-scope: PASS
✓ task-size: PASS

All validation rules passed!
```

### Step 11: Run the Tests

```bash
# Run all tests
make test

# Or run specific functional tests
go test ./tests/functional -run TestBeginnerScenario -v
```

**Expected**: All tests pass, including the automated beginner scenario test.

### Beginner Scenario Summary

Congratulations! You've completed the beginner scenario. You learned:

✅ How to initialize a project
✅ How to create tasks
✅ The task lifecycle: backlog → in-progress → done
✅ How to use both interactive and flag modes
✅ Basic validation
✅ Running tests

**Next**: Move to the Intermediate scenario to learn about specifications and context management.

---

## Scenario 2: Intermediate - Building a Blog API

**Goal**: Build a blog API feature using specification-driven development with proper context management.

**Time**: ~20 minutes

**What You'll Learn**:
- Creating tasks with full metadata
- Using specification references
- Generating context files
- Validation rules
- Template-based task creation

### Step 1: Create the Blog API Project

```bash
mkdir blog-api
cd blog-api

agentic-agent init --name "Blog API"
```

### Step 2: Create a Specification

Specifications are the source of truth. Create one for your API:

```bash
cat > .agentic/spec/01-blog-api.md << 'EOF'
# Blog API Specification

## Overview
A RESTful API for managing blog posts.

## Endpoints

### GET /posts
Returns all blog posts.

### POST /posts
Creates a new blog post.
- Request: { title, content, author }
- Response: { id, title, content, author, created_at }

### GET /posts/:id
Returns a single blog post by ID.

### PUT /posts/:id
Updates an existing blog post.

### DELETE /posts/:id
Deletes a blog post.

## Data Model

```go
type Post struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Author    string    `json:"author"`
    CreatedAt time.Time `json:"created_at"`
}
```

## Acceptance Criteria
- All CRUD operations work
- Proper error handling
- Input validation
- Comprehensive tests
EOF
```

### Step 3: Create Task from Template

```bash
# Interactive: Choose "Feature Implementation" template
agentic-agent task from-template

# Or use flags
agentic-agent task create \
  --title "Implement Blog Post API" \
  --description "Create REST API endpoints for blog post CRUD operations" \
  --spec-refs ".agentic/spec/01-blog-api.md" \
  --outputs "src/api/posts.go,src/models/post.go,tests/api/posts_test.go" \
  --acceptance "GET /posts returns all posts,POST /posts creates a new post,PUT /posts/:id updates a post,DELETE /posts/:id deletes a post,All endpoints have tests"
```

**Output**:
```
✓ Task created successfully!

ID: TASK-1234567891
Title: Implement Blog Post API
Spec Refs: .agentic/spec/01-blog-api.md
Expected Outputs: 3 files
```

### Step 4: Create Directory Structure

```bash
mkdir -p src/api
mkdir -p src/models
mkdir -p tests/api
```

### Step 5: Generate Context for Each Directory

Context files guide AI agents about what each directory does and its constraints.

```bash
# Generate context for API directory
agentic-agent context generate src/api
```

This creates `src/api/context.md`. Edit it:

```markdown
# API Module Context

## Purpose
HTTP handlers for blog API endpoints

## Responsibilities
- Handle HTTP requests
- Validate input
- Call business logic
- Return JSON responses

## Dependencies
- src/models (for Post struct)
- net/http (HTTP handling)

## Must Do
- Validate all inputs
- Return proper HTTP status codes
- Handle errors gracefully

## Cannot Do
- Direct database access (use repository pattern)
- Business logic (delegate to services)
```

Repeat for other directories:

```bash
agentic-agent context generate src/models
agentic-agent context generate tests/api
```

### Step 6: Claim the Task

```bash
agentic-agent task claim TASK-1234567891
```

### Step 7: Implement the Code

Create the model:

```bash
cat > src/models/post.go << 'EOF'
package models

import "time"

type Post struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Author    string    `json:"author"`
    CreatedAt time.Time `json:"created_at"`
}

// Validate checks if the post has required fields
func (p *Post) Validate() error {
    if p.Title == "" {
        return errors.New("title is required")
    }
    if p.Content == "" {
        return errors.New("content is required")
    }
    if p.Author == "" {
        return errors.New("author is required")
    }
    return nil
}
EOF
```

Create the API handlers:

```bash
cat > src/api/posts.go << 'EOF'
package api

import (
    "encoding/json"
    "net/http"
    "github.com/yourproject/blog-api/src/models"
)

type PostsHandler struct {
    // Repository would go here
}

func (h *PostsHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
    // Implementation
    posts := []models.Post{}
    json.NewEncoder(w).Encode(posts)
}

func (h *PostsHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
    var post models.Post
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := post.Validate(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Save post
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(post)
}
EOF
```

Create tests:

```bash
cat > tests/api/posts_test.go << 'EOF'
package api_test

import (
    "testing"
    "net/http/httptest"
    "github.com/yourproject/blog-api/src/api"
)

func TestGetPosts(t *testing.T) {
    handler := &api.PostsHandler{}
    req := httptest.NewRequest("GET", "/posts", nil)
    w := httptest.NewRecorder()

    handler.GetPosts(w, req)

    if w.Code != 200 {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
}
EOF
```

### Step 8: Validate Your Work

```bash
agentic-agent validate
```

**Expected Validation Checks**:
- ✓ All directories have context.md files
- ✓ Modified files are within task scope
- ✓ Task size is within limits

### Step 9: Build Context Bundle for Task

```bash
agentic-agent context build --task TASK-1234567891 --format markdown
```

This creates a focused context bundle for the task, combining:
- Global context
- Rolling summary
- Relevant directory contexts
- Task specifications

### Step 10: Complete the Task

```bash
agentic-agent task complete TASK-1234567891
```

### Step 11: Run Tests

```bash
# Run all tests
make test

# Run intermediate scenario test
go test ./tests/functional -run TestIntermediateScenario -v

# Run your API tests
go test ./tests/api -v
```

### Intermediate Scenario Summary

Excellent work! You've completed the intermediate scenario. You learned:

✅ Specification-driven development
✅ Context file generation and maintenance
✅ Task metadata (spec refs, outputs, acceptance criteria)
✅ Template-based task creation
✅ Context bundling
✅ Comprehensive validation

**Next**: Move to the Advanced scenario to learn about task decomposition.

---

## Scenario 3: Advanced - Complex Feature Decomposition

**Goal**: Manage a complex feature by decomposing it into subtasks and working through them systematically.

**Time**: ~30 minutes

**What You'll Learn**:
- Breaking large features into subtasks
- Managing subtask dependencies
- Progress tracking
- Comprehensive testing workflow

### Step 1: Create the Project

```bash
mkdir user-management
cd user-management

agentic-agent init --name "User Management System"
```

### Step 2: Create a Large Feature Task

```bash
agentic-agent task create \
  --title "Implement Complete User Authentication System" \
  --description "Full-featured authentication with registration, login, JWT, and password reset" \
  --acceptance "Users can register,Users can login,JWT tokens are issued,Protected routes work,Password reset via email,All features have tests"
```

Let's say the task ID is `TASK-1234567892`.

### Step 3: Decompose into Subtasks

Large tasks should be broken down into manageable pieces:

```bash
agentic-agent task decompose TASK-1234567892 \
  "Create user model and database schema" \
  "Implement user registration endpoint" \
  "Implement login with JWT generation" \
  "Add JWT validation middleware" \
  "Implement password reset flow" \
  "Write comprehensive integration tests"
```

**Output**:
```
✓ Task decomposed successfully!

Parent Task: TASK-1234567892
Subtasks created:
  1. TASK-1234567892.1: Create user model and database schema
  2. TASK-1234567892.2: Implement user registration endpoint
  3. TASK-1234567892.3: Implement login with JWT generation
  4. TASK-1234567892.4: Add JWT validation middleware
  5. TASK-1234567892.5: Implement password reset flow
  6. TASK-1234567892.6: Write comprehensive integration tests
```

### Step 4: View the Parent Task

```bash
agentic-agent task show TASK-1234567892
```

**Output**:
```
Task Details
─────────────────────────────────────
ID:    TASK-1234567892
Title: Implement Complete User Authentication System

Subtasks:
 ☐ TASK-1234567892.1: Create user model and database schema
 ☐ TASK-1234567892.2: Implement user registration endpoint
 ☐ TASK-1234567892.3: Implement login with JWT generation
 ☐ TASK-1234567892.4: Add JWT validation middleware
 ☐ TASK-1234567892.5: Implement password reset flow
 ☐ TASK-1234567892.6: Write comprehensive integration tests
```

### Step 5: Work on Subtask 1 - User Model

```bash
# Claim the parent task (or individual subtasks)
agentic-agent task claim TASK-1234567892

# Create directory structure
mkdir -p src/models
mkdir -p src/db

# Generate context
agentic-agent context generate src/models
agentic-agent context generate src/db
```

Edit context files to describe the modules.

Create the user model:

```bash
cat > src/models/user.go << 'EOF'
package models

import (
    "time"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID           string    `json:"id" db:"id"`
    Email        string    `json:"email" db:"email"`
    PasswordHash string    `json:"-" db:"password_hash"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func (u *User) SetPassword(password string) error {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.PasswordHash = string(hash)
    return nil
}

func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
    return err == nil
}
EOF
```

Mark subtask 1 as done (in actual implementation, you'd update the subtask status).

### Step 6: Work on Subtask 2 - Registration

```bash
mkdir -p src/auth
agentic-agent context generate src/auth
```

Create registration endpoint:

```bash
cat > src/auth/register.go << 'EOF'
package auth

import (
    "encoding/json"
    "net/http"
    "github.com/yourproject/user-mgmt/src/models"
)

type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    user := &models.User{Email: req.Email}
    if err := user.SetPassword(req.Password); err != nil {
        http.Error(w, "Failed to hash password", http.StatusInternalServerError)
        return
    }

    // Save user to database

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
EOF
```

### Step 7: Validate Continuously

As you work on each subtask:

```bash
agentic-agent validate --format json
```

This helps catch issues early:
- Missing context files
- Files outside task scope
- Task size violations

### Step 8: Track Progress with Learnings

Document patterns you discover:

```bash
agentic-agent learnings add "Use bcrypt for password hashing with DefaultCost"
agentic-agent learnings add "Always validate email format before user creation"
agentic-agent learnings add "JWT tokens should expire after 24 hours"

# View learnings
agentic-agent learnings list
```

### Step 9: Run Tests for Each Subtask

```bash
# Create tests for the user model
cat > src/models/user_test.go << 'EOF'
package models

import "testing"

func TestUserSetPassword(t *testing.T) {
    user := &User{Email: "test@example.com"}
    err := user.SetPassword("password123")
    if err != nil {
        t.Fatalf("Failed to set password: %v", err)
    }

    if !user.CheckPassword("password123") {
        t.Error("Password verification failed")
    }

    if user.CheckPassword("wrongpassword") {
        t.Error("Wrong password should not verify")
    }
}
EOF

# Run the test
go test ./src/models -v
```

### Step 10: Complete the Epic

After all subtasks are done:

```bash
agentic-agent task complete TASK-1234567892
```

### Step 11: Run Full Test Suite

```bash
# Run all tests
make test

# Run with coverage
make coverage-html

# Run advanced scenario test
go test ./tests/functional -run TestAdvancedScenario -v

# Run complete workflow test
go test ./tests/functional -run TestCompleteWorkflow -v
```

### Advanced Scenario Summary

Outstanding! You've completed the advanced scenario. You learned:

✅ Task decomposition strategies
✅ Managing complex multi-step features
✅ Progress tracking with learnings
✅ Continuous validation
✅ Comprehensive testing workflow
✅ Context generation for all modules

---

## Testing Workflow

This section covers how to integrate testing throughout your development process.

### Test Levels

The agentic-agent framework supports multiple testing levels:

1. **Unit Tests** - Test individual functions and modules
2. **Integration Tests** - Test complete workflows
3. **Functional Tests** - Test CLI commands and scenarios

### Running Tests

#### All Tests

```bash
# Run all tests in the project
make test

# Run with verbose output
make test-verbose
```

#### Package-Specific Tests

```bash
# Test task management
go test ./internal/tasks -v

# Test validation rules
go test ./internal/validator/rules -v

# Test functional scenarios
go test ./tests/functional -v
```

#### Specific Test Functions

```bash
# Run beginner scenario test
go test ./tests/functional -run TestBeginnerScenario -v

# Run intermediate scenario test
go test ./tests/functional -run TestIntermediateScenario -v

# Run advanced scenario test
go test ./tests/functional -run TestAdvancedScenario -v
```

### Test Coverage

#### Generate Coverage Reports

```bash
# Quick coverage check
make test-coverage

# Generate HTML coverage report
make coverage-html

# Show coverage by function
make coverage-func

# Show coverage summary
make coverage-summary
```

#### Understanding Coverage

Coverage reports show:
- **Green (70%+)**: Good coverage
- **Yellow (40-69%)**: Needs improvement
- **Red (<40%)**: Insufficient coverage

### Test-Driven Development (TDD)

Follow this workflow for TDD:

1. **Write the test first**:
```bash
cat > src/models/post_test.go << 'EOF'
package models

import "testing"

func TestPostValidation(t *testing.T) {
    post := &Post{}
    err := post.Validate()
    if err == nil {
        t.Error("Expected validation error for empty post")
    }
}
EOF
```

2. **Run the test (it should fail)**:
```bash
go test ./src/models -v
```

3. **Implement the feature**:
```bash
# Add Validate() method to Post struct
```

4. **Run the test again (it should pass)**:
```bash
go test ./src/models -v
```

5. **Refactor if needed**

### CI/CD Integration

The functional tests are designed for CI/CD pipelines:

```bash
# Non-interactive mode with JSON output
agentic-agent validate --format json --no-interactive

# Exit codes
# 0 = success
# 1 = failure
```

Example GitHub Actions workflow:

```yaml
name: Test
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.22'
      - name: Run tests
        run: make test
      - name: Run functional tests
        run: go test ./tests/functional -v
      - name: Validate
        run: ./agentic-agent validate --format json
```

---

## Best Practices

### When to Use Interactive vs Flag Mode

**Use Interactive Mode When**:
- You're learning the tool
- You're unsure of exact syntax
- You want to browse and select from options
- You're creating tasks with complex metadata

**Use Flag Mode When**:
- You're writing scripts
- You're in a CI/CD pipeline
- You know exactly what you want
- You're batch processing

### Task Sizing Guidelines

**Keep tasks small and focused**:

✅ **Good Task Size**:
- 1-5 files
- 1-2 directories
- Completable in one session
- Single, clear responsibility

❌ **Too Large**:
- 10+ files
- 5+ directories
- Spans multiple concerns
- Takes days to complete

**Solution**: Use `task decompose` to break large tasks into subtasks.

### Context Maintenance

**Update context files when**:
- You change a module's responsibility
- You add new dependencies
- You modify constraints
- You discover new patterns

```bash
# Regenerate context after changes
agentic-agent context generate src/module
```

### Specification-Driven Development

**Always start with specifications**:

1. Write the spec first
2. Create tasks that reference the spec
3. Implement according to spec
4. Update spec if requirements change

**Don't let agents be the source of truth**
**Specs are the source of truth**

### Validation Habits

Run validation frequently:

```bash
# Before committing
agentic-agent validate

# After major changes
agentic-agent validate --format json

# In pre-commit hooks
#!/bin/bash
agentic-agent validate || exit 1
```

---

## Troubleshooting

### Common Issues

#### Issue: "command not found: agentic-agent"

**Solution**:
```bash
# Build the binary
go build -o agentic-agent ./cmd/agentic-agent

# Add to PATH
export PATH="$PATH:$(pwd)"

# Or install globally
sudo mv agentic-agent /usr/local/bin/
```

#### Issue: "task not found in backlog"

**Solution**:
The task might be in a different state. Check all lists:

```bash
agentic-agent task list
# or
agentic-agent task show TASK-ID
```

Tasks can only be claimed from the backlog.

#### Issue: Validation fails with "Missing context.md"

**Solution**:
Generate context for all source directories:

```bash
# Find directories needing context
agentic-agent context scan

# Generate for specific directory
agentic-agent context generate src/module
```

#### Issue: Task size validation fails

**Solution**:
Decompose the task:

```bash
agentic-agent task decompose TASK-ID \
  "Subtask 1" \
  "Subtask 2" \
  "Subtask 3"
```

#### Issue: Git integration not working

**Solution**:
Ensure you're in a git repository:

```bash
git init
git add .
git commit -m "Initial commit"
```

#### Issue: Tests fail with timeout

**Solution**:
Some tests create multiple tasks with time-based IDs. If tests fail randomly, it might be timing-related. The functional tests include appropriate delays.

---

## Quick Reference

### Essential Commands

```bash
# Project
agentic-agent init --name "Project"
agentic-agent start                    # Interactive wizard

# Tasks
agentic-agent task create              # Interactive
agentic-agent task create --title "Task"
agentic-agent task list
agentic-agent task show TASK-ID
agentic-agent task claim TASK-ID
agentic-agent task complete TASK-ID
agentic-agent task decompose TASK-ID "Sub1" "Sub2"
agentic-agent task sample-task
agentic-agent task from-template

# Context
agentic-agent context generate DIR
agentic-agent context scan
agentic-agent context build --task TASK-ID

# Validation
agentic-agent validate
agentic-agent validate --format json

# Testing
make test
make test-verbose
make coverage-html
go test ./tests/functional -v
```

### Keyboard Shortcuts (Interactive Mode)

- **↑/↓ or j/k**: Navigate
- **Tab**: Switch tabs
- **Space**: Select/deselect
- **Enter**: Confirm
- **Esc or q**: Quit
- **a**: Add item (in lists)

### Directory Structure

```
.agentic/
├── tasks/              # Task YAML files
├── context/            # Context summaries
├── spec/               # Specifications
└── agent-rules/        # Agent configs

src/                    # Source code
└── */context.md        # Directory contexts

tests/                  # Test files
├── integration/        # Integration tests
└── functional/         # CLI functional tests
```

### Validation Rules

1. **directory-context**: Every source directory must have context.md
2. **context-update**: Context must be updated when files change
3. **task-scope**: Modified files must be in task scope
4. **task-size**: Max 5 files, 2 directories per task

---

## ATDD/BDD Workflow

**Goal**: Learn how to use Acceptance Test-Driven Development (ATDD) with Gherkin features to drive your development process.

**Time**: ~15 minutes

**What You'll Learn**:
- Writing executable specifications with Gherkin
- Running BDD tests
- Test-first development workflow
- Living documentation

### What is ATDD/BDD?

**Acceptance Test-Driven Development (ATDD)** is a development approach where you:
1. Write acceptance tests **before** implementing features
2. Use plain language that stakeholders can understand
3. Let tests serve as both specification and validation

**Behavior-Driven Development (BDD)** uses Gherkin syntax (Given/When/Then) to describe behavior from the user's perspective.

### Step 1: Understanding Feature Files

Feature files live in the `features/` directory and use Gherkin syntax:

```gherkin
Feature: Task Creation
  As a developer using the CLI
  I want to create tasks with metadata
  So that I have clear requirements

  Scenario: Create a simple task
    Given a clean test environment
    And I have initialized a project
    When I create a task with title "Implement login"
    Then the command should succeed
    And the task should appear in the backlog
```

### Step 2: Run Existing BDD Tests

The project includes comprehensive BDD tests:

```bash
# Run all BDD tests
make test-bdd

# Expected output:
# Running BDD feature tests...
# Feature: Project Initialization
#   Scenario: Initialize a new project successfully ✓
#   Scenario: Initialize project in non-git directory ✓
#
# Feature: Beginner Workflow
#   Scenario: Complete beginner workflow from CLI tutorial ✓
#
# 12 scenarios (12 passed)
# 107 steps (107 passed)
```

### Step 3: Run Tests for Specific Features

You can run specific feature files or scenarios:

```bash
# Run only workflow features
go test ./tests/bdd -v -godog.paths=../../features/workflows/

# Run only beginner workflow
go test ./tests/bdd -v -godog.paths=../../features/workflows/beginner_workflow.feature

# Run tests tagged as @smoke
go test ./tests/bdd -v -godog.tags=@smoke
```

### Step 4: The ATDD Development Cycle

Follow this workflow when adding new features:

#### 4.1 Write Feature First (Red Phase)

Create a feature file describing desired behavior:

```gherkin
# features/tasks/task_priority.feature
Feature: Task Priority Management
  As a developer managing multiple tasks
  I want to set task priorities
  So that I can focus on important work

  Scenario: Set task priority to high
    Given I have created a task with title "Critical Bug Fix"
    When I set the task priority to "high"
    Then the task priority should be "high"
    And the task should appear first in backlog listings
```

#### 4.2 Run Tests (They Fail)

```bash
make test-bdd

# Output shows undefined steps:
# Step definition missing for:
#   I set the task priority to "high"
```

#### 4.3 Implement Step Definitions (Yellow Phase)

Add step definitions in `tests/bdd/steps/task_steps.go`:

```go
func (s *TaskSteps) RegisterSteps(sc *godog.ScenarioContext) {
    // ... existing steps ...
    sc.Step(`^I set the task priority to "([^"]*)"$`, s.setTaskPriority)
    sc.Step(`^the task priority should be "([^"]*)"$`, s.assertTaskPriority)
}

func (s *TaskSteps) setTaskPriority(ctx context.Context, priority string) error {
    // Call actual CLI command or internal API
    return nil
}

func (s *TaskSteps) assertTaskPriority(ctx context.Context, expected string) error {
    // Verify priority was set correctly
    return nil
}
```

#### 4.4 Implement Feature (Green Phase)

Now implement the actual feature in your code:

```go
// cmd/agentic-agent/task.go
var setPriorityCmd = &cobra.Command{
    Use:   "set-priority [task-id] [priority]",
    Short: "Set task priority",
    Run:   setTaskPriority,
}

func setTaskPriority(cmd *cobra.Command, args []string) {
    // Implementation...
}
```

#### 4.5 Tests Pass

```bash
make test-bdd

# Output:
# Feature: Task Priority Management
#   Scenario: Set task priority to high ✓
#
# All tests pass!
```

#### 4.6 Refactor with Confidence

Refactor code knowing tests will catch regressions:

```bash
make test-all  # Unit + Functional + BDD all pass
```

### Step 5: Exploring Existing Features

Review existing feature files to learn patterns:

```bash
# List all feature files
find features -name "*.feature"

# View beginner workflow
cat features/workflows/beginner_workflow.feature

# View error handling examples
cat features/tasks/error_handling.feature
```

### Step 6: Running Coverage with BDD Tests

See how BDD tests contribute to overall coverage:

```bash
# Run all tests with coverage
make coverage-all

# Output shows:
# - Unit test coverage
# - Functional test coverage
# - BDD test coverage
# - Merged total coverage
#
# Reports generated:
# - coverage/merged-coverage.html (all tests)
# - coverage/bdd-coverage.html (BDD only)
```

### Benefits of ATDD/BDD

**Living Documentation**: Feature files are always up-to-date because:
- They're executable tests
- They fail if out of sync with code
- They're readable by non-developers

**Clear Requirements**: Before writing any code, you know:
- Exactly what behavior is expected
- What acceptance criteria must be met
- How to verify the feature works

**Better Design**: Writing tests first:
- Forces thinking about user experience
- Identifies edge cases early
- Results in simpler, more testable code

**Collaboration**: Gherkin syntax enables:
- Product owners to review specifications
- Developers to implement with clarity
- QA to understand test coverage

### Common Gherkin Patterns

**Background** (runs before each scenario):
```gherkin
Background:
  Given a clean test environment
  And I have initialized a project
```

**Data Tables**:
```gherkin
When I add the following acceptance criteria:
  | criterion                  |
  | All endpoints have tests   |
  | API returns proper status  |
```

**Scenario Outlines** (parameterized tests):
```gherkin
Scenario Outline: Validate title length
  When I create a task with title "<title>"
  Then the command should <result>

  Examples:
    | title        | result  |
    | Valid        | succeed |
    |              | fail    |
```

**Tags** (organize and filter):
```gherkin
@smoke @critical
Scenario: Critical user workflow
  ...
```

### ATDD Best Practices

1. **Write features before code**: Define behavior first
2. **Keep scenarios focused**: One scenario tests one behavior
3. **Use descriptive names**: Make intent clear
4. **Avoid implementation details**: Focus on behavior, not internals
5. **Reuse step definitions**: Create generic, reusable steps
6. **Tag strategically**: Use tags to organize test suites

### Further Reading

For comprehensive guidance on ATDD/BDD:
- [BDD Guide](BDD_GUIDE.md) - Complete guide to BDD with examples
- [Feature Files](../features/) - All existing feature specifications
- [Step Definitions](../tests/bdd/steps/) - Implementation patterns

---

## Next Steps

After completing this tutorial:

1. **Practice**: Create a real project using these patterns
2. **Explore**: Try different templates and workflows
3. **Customize**: Adjust the framework to your needs
4. **Contribute**: Share your patterns and improvements
5. **Read Documentation**:
   - [README.md](../README.md) - Full feature list
   - [COVERAGE.md](COVERAGE.md) - Test coverage details
   - [CLAUDE.md](../CLAUDE.md) - Claude-specific rules

## Automated Verification

All scenarios in this tutorial are verified by automated functional tests:

```bash
# Test beginner scenario
go test ./tests/functional -run TestBeginnerScenario -v

# Test intermediate scenario
go test ./tests/functional -run TestIntermediateScenario -v

# Test advanced scenario
go test ./tests/functional -run TestAdvancedScenario -v

# Test complete workflow
go test ./tests/functional -run TestCompleteWorkflow -v
```

These tests ensure the tutorial stays accurate and commands work as documented.

---

**Happy building with agentic-agent!** 🚀

For questions or issues, visit: https://github.com/javierbenavides/agentic-agent/issues
