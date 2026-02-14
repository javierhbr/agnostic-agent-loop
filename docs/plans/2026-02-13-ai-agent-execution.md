# AI Agent Execution Integration Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Integrate AI model execution into autopilot and create MCP server for Ralph loop integration.

**Architecture:** Add agent executor abstraction supporting Claude API, Anthropic SDK, and MCP server communication. Enhance autopilot to invoke agents and check completion. Create MCP server exposing agentic-agent operations for human-driven Ralph loops.

**Tech Stack:** Go 1.21+, Anthropic SDK, MCP protocol (TypeScript), Cobra CLI, existing agentic-agent infrastructure

---

## Task 1: Create Agent Executor Interface

**Files:**
- Create: `internal/agents/executor.go`
- Create: `internal/agents/executor_test.go`
- Create: `pkg/models/agent.go`

**Step 1: Write the failing test**

```go
// internal/agents/executor_test.go
package agents

import (
	"context"
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

func TestExecutor_Execute(t *testing.T) {
	task := &models.Task{
		ID:    "TASK-123",
		Title: "Test task",
		Acceptance: []string{
			"Test passes",
		},
	}

	executor := NewExecutor("mock")
	result, err := executor.Execute(context.Background(), "test prompt", task)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Output == "" {
		t.Error("Expected output, got empty string")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/agents/... -v -run TestExecutor_Execute`
Expected: FAIL with "undefined: NewExecutor"

**Step 3: Create agent models**

```go
// pkg/models/agent.go
package models

type AgentExecutionResult struct {
	Output           string
	Success          bool
	CriteriaMet      []string
	CriteriaFailed   []string
	FilesModified    []string
	ErrorMessage     string
	TokensUsed       int
}

func (r *AgentExecutionResult) AllCriteriaMet() bool {
	return r.Success && len(r.CriteriaFailed) == 0
}
```

**Step 4: Implement executor interface**

```go
// internal/agents/executor.go
package agents

import (
	"context"
	"fmt"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

type Executor interface {
	Execute(ctx context.Context, prompt string, task *models.Task) (*models.AgentExecutionResult, error)
}

type executor struct {
	agentType string
}

func NewExecutor(agentType string) Executor {
	return &executor{
		agentType: agentType,
	}
}

func (e *executor) Execute(ctx context.Context, prompt string, task *models.Task) (*models.AgentExecutionResult, error) {
	// Mock implementation for testing
	if e.agentType == "mock" {
		return &models.AgentExecutionResult{
			Output:  "Mock output",
			Success: true,
			CriteriaMet: task.Acceptance,
			TokensUsed: 1000,
		}, nil
	}

	return nil, fmt.Errorf("unsupported agent type: %s", e.agentType)
}
```

**Step 5: Run test to verify it passes**

Run: `go test ./internal/agents/... -v -run TestExecutor_Execute`
Expected: PASS

**Step 6: Commit**

```bash
git add internal/agents/executor.go internal/agents/executor_test.go pkg/models/agent.go
git commit -m "feat: add agent executor interface with mock implementation"
```

---

## Task 2: Add Claude API Executor

**Files:**
- Modify: `internal/agents/executor.go`
- Create: `internal/agents/claude.go`
- Create: `internal/agents/claude_test.go`
- Modify: `go.mod`

**Step 1: Add Anthropic SDK dependency**

```bash
go get github.com/anthropics/anthropic-sdk-go
```

Run: `go get github.com/anthropics/anthropic-sdk-go`
Expected: Dependency added to go.mod

**Step 2: Write failing test for Claude executor**

```go
// internal/agents/claude_test.go
package agents

import (
	"context"
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

func TestClaudeExecutor_Execute(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	task := &models.Task{
		ID:    "TASK-123",
		Title: "Test task",
		Acceptance: []string{
			"Response contains 'hello'",
		},
	}

	executor := NewClaudeExecutor("test-key", "claude-3-5-sonnet-20241022")
	result, err := executor.Execute(context.Background(), "Say hello", task)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Output == "" {
		t.Error("Expected output, got empty string")
	}
}
```

**Step 3: Run test to verify it fails**

Run: `go test ./internal/agents/... -v -run TestClaudeExecutor_Execute`
Expected: FAIL with "undefined: NewClaudeExecutor"

**Step 4: Implement Claude executor**

```go
// internal/agents/claude.go
package agents

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/javierbenavides/agentic-agent/pkg/models"
)

type ClaudeExecutor struct {
	client *anthropic.Client
	model  string
}

func NewClaudeExecutor(apiKey, model string) *ClaudeExecutor {
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	if model == "" {
		model = "claude-3-5-sonnet-20241022"
	}

	client := anthropic.NewClient(option.WithAPIKey(apiKey))

	return &ClaudeExecutor{
		client: client,
		model:  model,
	}
}

func (c *ClaudeExecutor) Execute(ctx context.Context, prompt string, task *models.Task) (*models.AgentExecutionResult, error) {
	// Build full prompt with task context
	fullPrompt := c.buildPrompt(prompt, task)

	// Call Claude API
	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.F(c.model),
		MaxTokens: anthropic.Int(4096),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(fullPrompt)),
		}),
	})

	if err != nil {
		return nil, fmt.Errorf("claude api error: %w", err)
	}

	// Extract output
	output := ""
	if len(message.Content) > 0 {
		if textBlock, ok := message.Content[0].AsUnion().(anthropic.TextBlock); ok {
			output = textBlock.Text
		}
	}

	// Check acceptance criteria
	criteriaMet, criteriaFailed := c.checkCriteria(output, task.Acceptance)

	return &models.AgentExecutionResult{
		Output:         output,
		Success:        len(criteriaFailed) == 0,
		CriteriaMet:    criteriaMet,
		CriteriaFailed: criteriaFailed,
		TokensUsed:     int(message.Usage.InputTokens + message.Usage.OutputTokens),
	}, nil
}

func (c *ClaudeExecutor) buildPrompt(basePrompt string, task *models.Task) string {
	var b strings.Builder

	b.WriteString(basePrompt)
	b.WriteString("\n\n")

	if len(task.Acceptance) > 0 {
		b.WriteString("## Acceptance Criteria\n")
		for _, criterion := range task.Acceptance {
			b.WriteString(fmt.Sprintf("- %s\n", criterion))
		}
		b.WriteString("\n")
	}

	b.WriteString("When all criteria are met, include in your response: <promise>TASK COMPLETE</promise>\n")

	return b.String()
}

func (c *ClaudeExecutor) checkCriteria(output string, criteria []string) (met []string, failed []string) {
	// Simple string matching for now
	// TODO: More sophisticated criteria checking

	if strings.Contains(output, "<promise>TASK COMPLETE</promise>") {
		// Assume all criteria met if completion signal found
		return criteria, []string{}
	}

	// Otherwise, mark all as failed
	return []string{}, criteria
}
```

**Step 5: Update executor factory**

```go
// internal/agents/executor.go (add to existing file)

func NewExecutor(agentType string) Executor {
	switch agentType {
	case "mock":
		return &executor{agentType: "mock"}
	case "claude-code", "claude":
		return NewClaudeExecutor("", "")
	default:
		return &executor{agentType: agentType}
	}
}
```

**Step 6: Run tests**

Run: `go test ./internal/agents/... -v -short`
Expected: PASS (skips integration test)

**Step 7: Commit**

```bash
git add internal/agents/claude.go internal/agents/claude_test.go internal/agents/executor.go go.mod go.sum
git commit -m "feat: add Claude API executor with Anthropic SDK"
```

---

## Task 3: Enhance Autopilot with Agent Execution

**Files:**
- Modify: `internal/orchestrator/autopilot.go:134-137`
- Modify: `cmd/agentic-agent/autopilot.go:36-40`
- Create: `internal/orchestrator/autopilot_test.go`

**Step 1: Write failing test**

```go
// internal/orchestrator/autopilot_test.go
package orchestrator

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

func TestAutopilotLoop_WithAgentExecution(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir("..")

	// Setup test config
	cfg := &models.Config{
		ActiveAgent: "mock",
		Paths: models.PathsConfig{
			SpecDirs: []string{".agentic/spec"},
		},
	}

	// Create test task
	os.MkdirAll(".agentic/tasks", 0755)
	taskData := `tasks:
  - id: TASK-001
    title: Test task
    status: pending
    acceptance:
      - Task completes successfully
`
	os.WriteFile(".agentic/tasks/backlog.yaml", []byte(taskData), 0644)

	// Run autopilot with agent execution enabled
	loop := NewAutopilotLoop(cfg, 1, "", false)
	loop.enableAgentExecution = true

	err := loop.Run(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify task was completed
	doneData, _ := os.ReadFile(".agentic/tasks/done.yaml")
	if !strings.Contains(string(doneData), "TASK-001") {
		t.Error("Expected task to be in done.yaml")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/orchestrator/... -v -run TestAutopilotLoop_WithAgentExecution`
Expected: FAIL with "enableAgentExecution undefined"

**Step 3: Add agent execution flag to AutopilotLoop**

```go
// internal/orchestrator/autopilot.go (add fields)

type AutopilotLoop struct {
	cfg                   *models.Config
	maxIterations         int
	stopSignal            string
	dryRun                bool
	enableAgentExecution  bool  // NEW
	taskManager           *tasks.TaskManager
	specResolver          *specs.Resolver
	trackManager          *tracks.Manager
}

func NewAutopilotLoop(cfg *models.Config, maxIterations int, stopSignal string, dryRun bool) *AutopilotLoop {
	// ... existing code ...
	return &AutopilotLoop{
		// ... existing fields ...
		enableAgentExecution: false, // NEW - default off
	}
}

// NEW method
func (a *AutopilotLoop) WithAgentExecution(enabled bool) *AutopilotLoop {
	a.enableAgentExecution = enabled
	return a
}
```

**Step 4: Implement agent execution in Run loop**

```go
// internal/orchestrator/autopilot.go (replace lines 134-137)

		// 6. Execute agent or report ready
		if a.enableAgentExecution {
			if err := a.executeAgent(ctx, task, bundle); err != nil {
				fmt.Printf("  Agent execution failed: %v\n", err)
				continue
			}
			fmt.Printf("✓ Task %s completed successfully\n", task.ID)
		} else {
			fmt.Printf("Task %s is ready for agent execution.\n", task.ID)
		}
	}

	// ... rest of loop ...
}

// NEW method
func (a *AutopilotLoop) executeAgent(ctx context.Context, task *models.Task, bundle string) error {
	fmt.Printf("  Executing agent: %s\n", a.cfg.ActiveAgent)

	// Create executor
	executor := agents.NewExecutor(a.cfg.ActiveAgent)

	// Execute
	result, err := executor.Execute(ctx, bundle, task)
	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	fmt.Printf("  Agent output (%d tokens):\n%s\n", result.TokensUsed, result.Output)

	// Check completion
	if !result.AllCriteriaMet() {
		return fmt.Errorf("acceptance criteria not met: %v", result.CriteriaFailed)
	}

	// Complete task
	if err := a.taskManager.CompleteTask(task.ID); err != nil {
		return fmt.Errorf("failed to complete task: %w", err)
	}

	return nil
}
```

**Step 5: Add import statement**

```go
// internal/orchestrator/autopilot.go (add to imports)
import (
	// ... existing imports ...
	"github.com/javierbenavides/agentic-agent/internal/agents"
)
```

**Step 6: Run test to verify it passes**

Run: `go test ./internal/orchestrator/... -v -run TestAutopilotLoop_WithAgentExecution`
Expected: PASS

**Step 7: Add CLI flag**

```go
// cmd/agentic-agent/autopilot.go (modify Run function around line 36)

	Run: func(cmd *cobra.Command, args []string) {
		maxIterations, _ := cmd.Flags().GetInt("max-iterations")
		stopSignal, _ := cmd.Flags().GetString("stop-signal")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		executeAgent, _ := cmd.Flags().GetBool("execute-agent") // NEW

		cfg := getConfig()

		loop := orchestrator.NewAutopilotLoop(cfg, maxIterations, stopSignal, dryRun)
		loop.WithAgentExecution(executeAgent) // NEW

		// ... rest of function ...
	},
```

**Step 8: Register flag**

```go
// cmd/agentic-agent/autopilot.go (add to init function)

func init() {
	autopilotStartCmd.Flags().Int("max-iterations", 10, "Maximum number of tasks to process")
	autopilotStartCmd.Flags().String("stop-signal", "", "Custom stop signal string")
	autopilotStartCmd.Flags().Bool("dry-run", false, "Show what would be processed without making changes")
	autopilotStartCmd.Flags().Bool("execute-agent", false, "Execute AI agent for each task") // NEW

	autopilotCmd.AddCommand(autopilotStartCmd)
}
```

**Step 9: Run integration test**

Run: `go test ./... -v -short`
Expected: PASS

**Step 10: Commit**

```bash
git add internal/orchestrator/autopilot.go internal/orchestrator/autopilot_test.go cmd/agentic-agent/autopilot.go
git commit -m "feat: add agent execution to autopilot with --execute-agent flag"
```

---

## Task 4: Create MCP Server Package

**Files:**
- Create: `.agentic/mcp-server/package.json`
- Create: `.agentic/mcp-server/tsconfig.json`
- Create: `.agentic/mcp-server/agentic-agent-server.ts`

**Step 1: Create package.json**

```json
{
  "name": "agentic-agent-mcp-server",
  "version": "1.0.0",
  "description": "MCP server for agentic-agent CLI integration",
  "main": "agentic-agent-server.ts",
  "scripts": {
    "start": "tsx agentic-agent-server.ts",
    "build": "tsc"
  },
  "dependencies": {
    "@modelcontextprotocol/sdk": "^0.5.0"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "tsx": "^4.7.0",
    "typescript": "^5.3.0"
  }
}
```

**Step 2: Create tsconfig.json**

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "Node16",
    "moduleResolution": "Node16",
    "outDir": "./dist",
    "rootDir": ".",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true
  },
  "include": ["*.ts"],
  "exclude": ["node_modules", "dist"]
}
```

**Step 3: Write MCP server**

```typescript
// .agentic/mcp-server/agentic-agent-server.ts
import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
} from "@modelcontextprotocol/sdk/types.js";
import { execSync } from "child_process";
import * as fs from "fs";

const server = new Server(
  {
    name: "agentic-agent-server",
    version: "1.0.0",
  },
  {
    capabilities: {
      tools: {},
    },
  }
);

// Tool definitions
server.setRequestHandler(ListToolsRequestSchema, async () => ({
  tools: [
    {
      name: "agentic_get_active_task",
      description: "Get the currently active task from the CLI session",
      inputSchema: {
        type: "object",
        properties: {},
      },
    },
    {
      name: "agentic_task_show",
      description: "Get detailed information about a specific task",
      inputSchema: {
        type: "object",
        properties: {
          task_id: { type: "string", description: "Task ID to show" },
        },
        required: ["task_id"],
      },
    },
    {
      name: "agentic_task_update_progress",
      description: "Update task progress during ralph loop",
      inputSchema: {
        type: "object",
        properties: {
          task_id: { type: "string" },
          iteration: { type: "number" },
          files_modified: { type: "array", items: { type: "string" } },
          tests_status: { type: "string" },
          notes: { type: "string" },
        },
        required: ["task_id", "iteration"],
      },
    },
    {
      name: "agentic_task_complete",
      description: "Mark task as complete with learnings",
      inputSchema: {
        type: "object",
        properties: {
          task_id: { type: "string" },
          learnings: { type: "array", items: { type: "string" } },
        },
        required: ["task_id"],
      },
    },
  ],
}));

// Tool implementations
server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const { name, arguments: args } = request.params;

  try {
    switch (name) {
      case "agentic_get_active_task": {
        const sessionPath = ".agentic/session.json";
        if (!fs.existsSync(sessionPath)) {
          return {
            content: [
              {
                type: "text",
                text: JSON.stringify({ error: "No active session" }),
              },
            ],
          };
        }

        const session = JSON.parse(fs.readFileSync(sessionPath, "utf-8"));
        return {
          content: [
            {
              type: "text",
              text: JSON.stringify(session.activeTask || {}),
            },
          ],
        };
      }

      case "agentic_task_show": {
        const result = execSync(
          `agentic-agent task show ${args.task_id} --no-interactive`,
          { encoding: "utf-8" }
        );
        return {
          content: [{ type: "text", text: result }],
        };
      }

      case "agentic_task_update_progress": {
        const sessionPath = ".agentic/session.json";
        const session = fs.existsSync(sessionPath)
          ? JSON.parse(fs.readFileSync(sessionPath, "utf-8"))
          : {};

        session.progress = {
          iteration: args.iteration,
          filesModified: args.files_modified || [],
          testsStatus: args.tests_status || "unknown",
          notes: args.notes || "",
          timestamp: new Date().toISOString(),
        };

        fs.writeFileSync(sessionPath, JSON.stringify(session, null, 2));

        return {
          content: [
            {
              type: "text",
              text: JSON.stringify({ status: "progress updated" }),
            },
          ],
        };
      }

      case "agentic_task_complete": {
        const result = execSync(
          `agentic-agent task complete ${args.task_id}`,
          { encoding: "utf-8" }
        );

        return {
          content: [{ type: "text", text: result }],
        };
      }

      default:
        throw new Error(`Unknown tool: ${name}`);
    }
  } catch (error: any) {
    return {
      content: [
        {
          type: "text",
          text: JSON.stringify({ error: error.message }),
        },
      ],
      isError: true,
    };
  }
});

// Start server
const transport = new StdioServerTransport();
server.connect(transport);

console.error("Agentic Agent MCP Server running on stdio");
```

**Step 4: Test MCP server compiles**

Run: `cd .agentic/mcp-server && npm install && npm run build`
Expected: TypeScript compiles successfully

**Step 5: Commit**

```bash
git add .agentic/mcp-server/
git commit -m "feat: add MCP server for ralph-loop integration"
```

---

## Task 5: Add Session Management Commands

**Files:**
- Create: `cmd/agentic-agent/session.go`
- Modify: `internal/tasks/lock.go:70-94`

**Step 1: Write session state model**

```go
// pkg/models/session.go
package models

import "time"

type SessionState struct {
	ActiveTask *Task     `json:"activeTask,omitempty"`
	ClaimedAt  time.Time `json:"claimedAt"`
	Branch     string    `json:"branch"`
	MCPEnabled bool      `json:"mcpEnabled"`
	Progress   *SessionProgress `json:"progress,omitempty"`
}

type SessionProgress struct {
	Iteration     int       `json:"iteration"`
	FilesModified []string  `json:"filesModified"`
	TestsStatus   string    `json:"testsStatus"`
	Notes         string    `json:"notes"`
	Timestamp     time.Time `json:"timestamp"`
}
```

**Step 2: Create session command**

```go
// cmd/agentic-agent/session.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/spf13/cobra"
)

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Manage AI agent sessions",
}

var sessionStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start MCP server for AI integration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := startMCPServer(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ MCP Server started")
		fmt.Println("📡 AI agents can now interact via MCP tools")
	},
}

var sessionStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop MCP server",
	Run: func(cmd *cobra.Command, args []string) {
		pidFile := ".agentic/mcp-server.pid"
		data, err := os.ReadFile(pidFile)
		if err != nil {
			fmt.Println("No MCP server running")
			return
		}

		var pid int
		fmt.Sscanf(string(data), "%d", &pid)

		process, err := os.FindProcess(pid)
		if err != nil {
			fmt.Printf("Error finding process: %v\n", err)
			os.Exit(1)
		}

		if err := process.Kill(); err != nil {
			fmt.Printf("Error stopping server: %v\n", err)
			os.Exit(1)
		}

		os.Remove(pidFile)
		fmt.Println("✅ MCP Server stopped")
	},
}

var sessionInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show current session information",
	Run: func(cmd *cobra.Command, args []string) {
		sessionPath := ".agentic/session.json"
		if _, err := os.Stat(sessionPath); os.IsNotExist(err) {
			fmt.Println("No active session")
			return
		}

		data, err := os.ReadFile(sessionPath)
		if err != nil {
			fmt.Printf("Error reading session: %v\n", err)
			os.Exit(1)
		}

		var session models.SessionState
		if err := json.Unmarshal(data, &session); err != nil {
			fmt.Printf("Error parsing session: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Session Information:")
		fmt.Printf("  Active Task: %s\n", session.ActiveTask.ID)
		fmt.Printf("  Title: %s\n", session.ActiveTask.Title)
		fmt.Printf("  Claimed At: %s\n", session.ClaimedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Branch: %s\n", session.Branch)
		fmt.Printf("  MCP Enabled: %v\n", session.MCPEnabled)

		if session.Progress != nil {
			fmt.Println("\nProgress:")
			fmt.Printf("  Iteration: %d\n", session.Progress.Iteration)
			fmt.Printf("  Files Modified: %d\n", len(session.Progress.FilesModified))
			fmt.Printf("  Tests Status: %s\n", session.Progress.TestsStatus)
		}
	},
}

func startMCPServer() error {
	mcpScript := ".agentic/mcp-server/agentic-agent-server.ts"
	if _, err := os.Stat(mcpScript); os.IsNotExist(err) {
		return fmt.Errorf("MCP server not found. Run: agentic-agent init-mcp")
	}

	cmd := exec.Command("npx", "tsx", mcpScript)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start MCP server: %w", err)
	}

	pidFile := ".agentic/mcp-server.pid"
	return os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644)
}

func init() {
	sessionCmd.AddCommand(sessionStartCmd)
	sessionCmd.AddCommand(sessionStopCmd)
	sessionCmd.AddCommand(sessionInfoCmd)
	rootCmd.AddCommand(sessionCmd)
}
```

**Step 3: Update task claim to write session**

```go
// internal/tasks/lock.go (modify ClaimTaskWithConfig, add after line 92)

func (tm *TaskManager) ClaimTaskWithConfig(taskID, assignee string, cfg *models.Config) error {
	// ... existing readiness check code ...

	// Claim the task
	if err := tm.ClaimTask(taskID, assignee); err != nil {
		return err
	}

	// Write session state
	task, source, _ := tm.FindTask(taskID)
	if source == "in-progress" && task != nil {
		writeSessionState(task)
	}

	return nil
}

func writeSessionState(task *models.Task) error {
	state := models.SessionState{
		ActiveTask: task,
		ClaimedAt:  time.Now(),
		Branch:     currentGitBranch(),
		MCPEnabled: false,
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(".agentic/session.json", data, 0644)
}
```

**Step 4: Test session commands**

Run: `go build ./cmd/agentic-agent && ./agentic-agent session info`
Expected: "No active session" or session info displayed

**Step 5: Commit**

```bash
git add cmd/agentic-agent/session.go internal/tasks/lock.go pkg/models/session.go
git commit -m "feat: add session management commands for MCP integration"
```

---

## Task 6: Documentation and Examples

**Files:**
- Create: `docs/AI_AGENT_EXECUTION.md`
- Create: `examples/autopilot-with-ai/README.md`
- Modify: `README.md`

**Step 1: Write main documentation**

```markdown
# docs/AI_AGENT_EXECUTION.md

# AI Agent Execution Guide

This guide covers how to execute tasks using AI agents in agentic-agent.

## Autopilot with Agent Execution

### Basic Usage

```bash
# Execute tasks automatically using detected agent
agentic-agent autopilot start --execute-agent --max-iterations 5

# Specify agent explicitly
export AGENTIC_AGENT=claude-code
agentic-agent autopilot start --execute-agent
```

### Configuration

Set your API key:
```bash
export ANTHROPIC_API_KEY=your-key-here
```

### How It Works

1. Autopilot finds next ready task
2. Claims task and generates context
3. Builds context bundle with specs
4. **Executes AI agent** with bundle as prompt
5. Checks acceptance criteria
6. Completes task if all criteria met
7. Moves to next task

## MCP Server for Ralph Loops

### Setup

Install MCP server dependencies:
```bash
cd .agentic/mcp-server
npm install
```

### Usage with Claude.ai

1. **Terminal**: Claim task with MCP
```bash
agentic-agent task claim TASK-123 --with-mcp
```

2. **Claude.ai**: Connect MCP server and use tools
```
Use agentic_get_active_task to see current task.
Then run /ralph-loop with MCP integration.
```

3. **Terminal**: Monitor progress
```bash
agentic-agent session info
```

### Available MCP Tools

- `agentic_get_active_task` - Get current task
- `agentic_task_show` - Get task details
- `agentic_task_update_progress` - Update progress
- `agentic_task_complete` - Complete task

## Examples

See `examples/autopilot-with-ai/` for complete examples.
```

**Step 2: Create example**

```markdown
# examples/autopilot-with-ai/README.md

# Autopilot with AI Agent Execution

Example of using autopilot with AI agent execution.

## Setup

1. Set API key:
```bash
export ANTHROPIC_API_KEY=your-key
```

2. Initialize project:
```bash
agentic-agent init --name "ai-demo"
```

3. Create test tasks:
```bash
agentic-agent task create --title "Test task 1" \
  --acceptance "Task outputs hello world"
agentic-agent task create --title "Test task 2" \
  --acceptance "Task outputs goodbye"
```

## Run Autopilot

```bash
# Dry run first
agentic-agent autopilot start --dry-run --max-iterations 2

# Execute with AI
agentic-agent autopilot start --execute-agent --max-iterations 2
```

## Expected Output

```
--- Iteration 1/2 ---
Next task: [TASK-001] Test task 1
Task TASK-001: READY
Claimed task TASK-001
  Executing agent: claude-code
  Agent output (1234 tokens):
Hello world

✓ Task TASK-001 completed successfully

--- Iteration 2/2 ---
Next task: [TASK-002] Test task 2
...
```
```

**Step 3: Update main README**

Add to README.md in appropriate section:

```markdown
### AI Agent Execution

Autopilot can execute AI agents automatically:

```bash
# Enable agent execution
agentic-agent autopilot start --execute-agent --max-iterations 5

# Agents execute tasks, check criteria, and complete automatically
```

See [AI Agent Execution Guide](docs/AI_AGENT_EXECUTION.md) for details.
```

**Step 4: Commit**

```bash
git add docs/AI_AGENT_EXECUTION.md examples/autopilot-with-ai/README.md README.md
git commit -m "docs: add AI agent execution documentation and examples"
```

---

## Completion Checklist

- [ ] Agent executor interface with mock and Claude implementations
- [ ] Autopilot enhanced with `--execute-agent` flag
- [ ] MCP server for ralph-loop integration
- [ ] Session management commands (start, stop, info)
- [ ] Documentation and examples
- [ ] All tests passing
- [ ] Integration tested with real API

## Final Testing

```bash
# Run all tests
go test ./... -v -short

# Test autopilot dry run
agentic-agent autopilot start --dry-run --max-iterations 1

# Test MCP server compilation
cd .agentic/mcp-server && npm run build

# Test session commands
agentic-agent session info
```
