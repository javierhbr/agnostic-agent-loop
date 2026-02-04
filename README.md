# Agnostic Agent Framework (Go Implementation)

A specialized CLI tool designed to facilitate **Agentic Specification-Driven Development (ASDD)**. It acts as the "brain" and "manager" for AI coding agents (like Claude Code, Cursor, Windsurf), ensuring they follow strict protocols for task management, context maintenance, and architectural integrity.

## Features

- **Project Initialization**: Bootstraps a standard directory structure (`.agentic/`).
- **Task Management**: CLI-based CRUD for tasks with strict lifecycle (pending -> in-progress -> done).
- **Context Management**: Enforces `context.md` files in source directories and maintains a global `rolling-summary.md`.
- **Validation**: Rules engine to enforce architectural constraints (e.g., "No code editing without context update").
- **Skill Generation**: Generates tool-specific configuration files (`CLAUDE.md`, `.cursor/rules/...`) to align AI agents with framework rules.
- **Orchestration**: State machine to guide agents through the Planning -> Execution -> Verification loop.
- **Token Management**: Basic accounting of token usage per agent.

## Installation

```bash
go build -o agentic-agent ./cmd/agentic-agent
mv agentic-agent /usr/local/bin/ # Optional
```

## Quick Start Guide

### 1. Initialize a Project
Run this in your project root:
```bash
agentic-agent init --name "My Project"
```

### 2. Create a Task
```bash
agentic-agent task create --title "Implement Login Feature"
```

### 3. Start Working (Claim Task)
This moves the task to "in-progress" and assigns it to the current agent/user.
```bash
agentic-agent task claim TASK-123456
```

### 4. Generate/Update Context
Before editing code in a directory, generate the context:
```bash
agentic-agent context generate src/auth
```
If you change code, update it again:
```bash
agentic-agent context update src/auth
```

### 5. Generate Agent Skills
Configure your AI assistant to follow these rules:
```bash
# For Claude Code
agentic-agent skills generate --tool claude-code

# For Cursor
agentic-agent skills generate --tool cursor
```

### 6. Run the Orchestrator
To simulate the agent loop for a task:
```bash
agentic-agent run --task TASK-123456
```

### 7. Validate project
Check if you are following all rules:
```bash
agentic-agent validate
```

### 8. Check Token Usage
```bash
agentic-agent token status
```
