# Ralph Loop Human Integration Design

## Problem Statement

When a **human** executes `/ralph-loop` in Claude.ai (or similar AI chat interface), there are two separate contexts:

1. **Terminal** — Where the human runs `agentic-agent` CLI commands
2. **AI Chat** — Where Ralph loop executes with file access via MCP

**Current workflow is disconnected:**

```
Terminal (Human)                     Claude.ai (AI)
────────────────                     ──────────────
$ agentic-agent task claim T1
$ agentic-agent task show T1
                                     Human: /ralph-loop "..."
                                     Claude: [reads files via MCP]
                                     Claude: [writes code]
                                     Claude: <promise>COMPLETE</promise>
$ agentic-agent task complete T1
```

**Problems:**
- Ralph doesn't know which task is active in CLI
- Ralph can't update task status automatically
- Human must manually bridge the two contexts
- No automatic commit capture or progress tracking
- Token usage not integrated with CLI tracking

## Design Goals

1. **Bidirectional communication** between CLI and AI chat
2. **Context awareness** — Ralph knows current task, specs, acceptance criteria
3. **Automatic updates** — Task status syncs without manual commands
4. **Progress tracking** — Commits and learnings captured automatically
5. **Token awareness** — Checkpoint when approaching limits

## Solution: MCP Server Integration

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Human's Setup                         │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Terminal                          Claude.ai                 │
│  ────────                          ────────                  │
│  $ agentic-agent                   [AI Chat Interface]       │
│    task claim T1                                             │
│         │                                                     │
│         ├──> Starts MCP Server     ◄────── Connects via MCP  │
│         │    (background process)                            │
│         │                                                     │
│         └──> .agentic/session.json ◄────── Reads via MCP     │
│              (active task state)                             │
│                                                               │
└─────────────────────────────────────────────────────────────┘

Session State Flow:
─────────────────
Terminal writes → .agentic/session.json → MCP server → Claude reads
Claude writes  → MCP server → .agentic/session.json → Terminal reads
```

### Component 1: MCP Server for Agentic Agent

Create an MCP server that exposes agentic-agent operations:

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
import * as path from "path";

const server = new Server(
  {
    name: "agentic-agent-server",
    version: "1.0.0",
  },
  {
    capabilities: {
      tools: {},
      resources: {},
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
          task_id: { type: "string" },
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
    {
      name: "agentic_get_token_status",
      description: "Get current token usage status",
      inputSchema: {
        type: "object",
        properties: {},
      },
    },
    {
      name: "agentic_create_checkpoint",
      description: "Create checkpoint for resumption",
      inputSchema: {
        type: "object",
        properties: {
          task_id: { type: "string" },
          iteration: { type: "number" },
          work_completed: { type: "array", items: { type: "string" } },
          remaining_work: { type: "array", items: { type: "string" } },
        },
        required: ["task_id", "iteration"],
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
          `agentic-agent task show ${args.task_id} --no-interactive --format json`,
          { encoding: "utf-8" }
        );
        return {
          content: [{ type: "text", text: result }],
        };
      }

      case "agentic_task_update_progress": {
        // Update session file with progress
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
        const learningsArg = args.learnings
          ? args.learnings.map((l: string) => `"${l}"`).join(" ")
          : "";

        const result = execSync(
          `agentic-agent task complete ${args.task_id} ${learningsArg ? `--learnings ${learningsArg}` : ""}`,
          { encoding: "utf-8" }
        );

        return {
          content: [{ type: "text", text: result }],
        };
      }

      case "agentic_get_token_status": {
        const result = execSync("agentic-agent token status --format json", {
          encoding: "utf-8",
        });
        return {
          content: [{ type: "text", text: result }],
        };
      }

      case "agentic_create_checkpoint": {
        const result = execSync(
          `agentic-agent token checkpoint ${args.task_id} --iteration ${args.iteration}`,
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
```

### Component 2: Enhanced Task Claim with Session State

```go
// cmd/agentic-agent/task.go

type SessionState struct {
    ActiveTask   *models.Task `json:"activeTask,omitempty"`
    ClaimedAt    time.Time    `json:"claimedAt"`
    Branch       string       `json:"branch"`
    MCPEnabled   bool         `json:"mcpEnabled"`
}

func startMCPServer() error {
    // Check if MCP server script exists
    mcpScript := ".agentic/mcp-server/agentic-agent-server.ts"
    if _, err := os.Stat(mcpScript); os.IsNotExist(err) {
        return fmt.Errorf("MCP server not found. Run: agentic-agent init-mcp")
    }

    // Start MCP server in background
    cmd := exec.Command("npx", "tsx", mcpScript)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Start(); err != nil {
        return fmt.Errorf("failed to start MCP server: %w", err)
    }

    // Save PID for cleanup
    pidFile := ".agentic/mcp-server.pid"
    os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644)

    return nil
}

func writeSessionState(task *models.Task) error {
    state := SessionState{
        ActiveTask:   task,
        ClaimedAt:    time.Now(),
        Branch:       getCurrentGitBranch(),
        MCPEnabled:   true,
    }

    data, err := json.MarshalIndent(state, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(".agentic/session.json", data, 0644)
}

// Enhanced task claim
var taskClaimCmd = &cobra.Command{
    Use:   "claim TASK_ID",
    Short: "Claim task and start MCP server for AI integration",
    Run: func(cmd *cobra.Command, args []string) {
        // ... existing claim logic ...

        // Write session state
        if err := writeSessionState(task); err != nil {
            fmt.Printf("Warning: Failed to write session state: %v\n", err)
        }

        // Check if --with-mcp flag is set
        withMCP, _ := cmd.Flags().GetBool("with-mcp")
        if withMCP {
            if err := startMCPServer(); err != nil {
                fmt.Printf("Warning: Failed to start MCP server: %v\n", err)
            } else {
                fmt.Println("\n✅ MCP Server started")
                fmt.Println("📡 Claude can now interact with agentic-agent via MCP tools")
                fmt.Println("\nAvailable MCP tools:")
                fmt.Println("  - agentic_get_active_task")
                fmt.Println("  - agentic_task_show")
                fmt.Println("  - agentic_task_update_progress")
                fmt.Println("  - agentic_task_complete")
                fmt.Println("  - agentic_get_token_status")
                fmt.Println("  - agentic_create_checkpoint")
            }
        }

        fmt.Printf("\nClaimed task %s\n", task.ID)
        fmt.Printf("\n💡 TIP: In Claude, use /ralph-loop with MCP tools to auto-track progress\n")
    },
}

func init() {
    taskClaimCmd.Flags().Bool("with-mcp", false, "Start MCP server for AI integration")
}
```

### Component 3: Enhanced Ralph Prompt with MCP Integration

When human claims a task with `--with-mcp`, the CLI outputs a ready-to-use prompt:

```bash
$ agentic-agent task claim TASK-123 --with-mcp

Claimed task TASK-123
✅ MCP Server started
📡 Claude can now interact with agentic-agent via MCP tools

──────────────────────────────────────────────────────────────
Copy this prompt to Claude:

I'm working on a task from agentic-agent. Before we start:

1. Use `agentic_get_active_task` to see what I'm working on
2. Use `agentic_task_show` with the task ID to get full details
3. Build a ralph-loop prompt from the acceptance criteria
4. During the loop:
   - Use `agentic_task_update_progress` after each iteration
   - Use `agentic_get_token_status` to check if approaching limits
   - Use `agentic_create_checkpoint` if tokens > 80%
5. When complete, use `agentic_task_complete` with learnings

Let's start by reading the active task.
──────────────────────────────────────────────────────────────
```

### Component 4: Ralph Loop Integration Example

**Human in Claude.ai:**

```
I'm working on a task from agentic-agent. Let me get the active task first.
```

**Claude with MCP:**

```typescript
// Claude calls MCP tool
agentic_get_active_task()
// Returns:
{
  "id": "TASK-123",
  "title": "[todo-app] Configure TypeScript and ESLint",
  "specRefs": ["todo-app/proposal.md", "todo-app/tasks/02-typescript-eslint.md"],
  "acceptance": [
    "tsconfig.json configured with strict mode enabled",
    "ESLint installed and configured for React + TypeScript",
    "No ESLint errors on existing code"
  ]
}

// Claude gets detailed task info
agentic_task_show({ task_id: "TASK-123" })

// Claude builds ralph prompt and executes
/ralph-loop "Implement TASK-123: Configure TypeScript and ESLint.

Read first:
- .agentic/spec/todo-app/proposal.md
- .agentic/spec/todo-app/tasks/02-typescript-eslint.md

Acceptance:
- tsconfig.json with strict mode
- ESLint for React + TypeScript
- No ESLint errors

After EACH iteration, call:
  agentic_task_update_progress({
    task_id: 'TASK-123',
    iteration: N,
    files_modified: [...],
    tests_status: '...'
  })

Check tokens every 3 iterations:
  agentic_get_token_status()

When ALL acceptance criteria pass:
  agentic_task_complete({
    task_id: 'TASK-123',
    learnings: ['...']
  })
  <promise>TASK COMPLETE</promise>
" --max-iterations 10 --completion-promise "TASK COMPLETE"
```

**During Ralph Loop:**

```
Iteration 1: Setup tsconfig.json...
[calls agentic_task_update_progress]

Iteration 2: Install ESLint...
[calls agentic_task_update_progress]

Iteration 3: Configure rules...
[calls agentic_get_token_status]
  → 45,000 tokens used (22.5% of limit) ✅

Iteration 4: Fix errors...
[calls agentic_task_update_progress]

Iteration 5: Verify...
All acceptance criteria pass!
[calls agentic_task_complete with learnings]
<promise>TASK COMPLETE</promise>
```

**Human's terminal automatically shows:**

```
$ # Session state auto-updates
$ cat .agentic/session.json
{
  "activeTask": {
    "id": "TASK-123",
    "status": "done"
  },
  "progress": {
    "iteration": 5,
    "filesModified": ["tsconfig.json", ".eslintrc.js", "package.json"],
    "testsStatus": "all passing",
    "timestamp": "2026-02-13T20:15:30Z"
  }
}

$ agentic-agent status
Task TASK-123: ✅ COMPLETED
Learnings captured: 3
Commits: 5
```

## CLI Commands

### Setup

```bash
# Initialize MCP server
agentic-agent init-mcp

# Install dependencies
cd .agentic/mcp-server && npm install
```

### Workflow

```bash
# 1. Claim task with MCP integration
agentic-agent task claim TASK-123 --with-mcp

# 2. Copy prompt to Claude.ai
# (CLI outputs ready-to-use prompt)

# 3. In Claude, run ralph-loop
# (MCP tools auto-update session state)

# 4. Monitor progress from terminal
agentic-agent session watch

# 5. Stop MCP server when done
agentic-agent session stop
```

## Session Commands

```bash
# Watch session state in real-time
agentic-agent session watch

# Show current session info
agentic-agent session info

# Start MCP server manually
agentic-agent session start

# Stop MCP server
agentic-agent session stop

# Clear session state
agentic-agent session clear
```

## Benefits

1. **Seamless Integration** — Terminal and AI chat stay in sync
2. **Automatic Updates** — No manual `task complete` needed
3. **Progress Visibility** — Human sees updates in real-time
4. **Token Awareness** — Auto-checkpoint when approaching limits
5. **Learnings Captured** — AI can record learnings via MCP
6. **Resume Friendly** — Session state persists for resumption

## Configuration

```yaml
# agnostic-agent.yaml

mcp:
  enabled: true
  server_script: ".agentic/mcp-server/agentic-agent-server.ts"
  session_file: ".agentic/session.json"
  auto_start: true  # Start MCP server on task claim

session:
  watch_interval: 2s  # How often to check for updates
  checkpoint_on_token_limit: true
```

## Implementation Phases

### Phase 1: MCP Server
- [ ] Create TypeScript MCP server
- [ ] Implement basic tools (get_active_task, task_show)
- [ ] Test with Claude.ai

### Phase 2: CLI Integration
- [ ] Add `--with-mcp` flag to task claim
- [ ] Write session state file
- [ ] Start/stop MCP server commands

### Phase 3: Progress Tracking
- [ ] Implement `task_update_progress` tool
- [ ] Add `session watch` command
- [ ] Real-time terminal updates

### Phase 4: Token Awareness
- [ ] Integrate `get_token_status` tool
- [ ] Auto-checkpoint on threshold
- [ ] Resumption prompt generation

## Open Questions

1. **How to handle MCP server crashes?**
   - Auto-restart on failure?
   - Graceful degradation (continue without MCP)?

2. **Multiple concurrent sessions?**
   - One session per project?
   - Session isolation by branch?

3. **Security concerns?**
   - MCP server has full CLI access
   - Should we add permission prompts?

4. **Cross-platform compatibility?**
   - Windows support for MCP server?
   - Alternative to npx/tsx?

## Summary

This design bridges the gap between human's terminal and AI chat by:

- **MCP Server** — Exposes agentic-agent operations to Claude
- **Session State** — Shared file keeps both contexts in sync
- **Auto-Updates** — AI can update task status without human
- **Token Tracking** — Integrated checkpoint/resume flow
- **Ready-to-Use Prompts** — CLI generates optimal ralph prompts

The human gets a **unified workflow** where terminal and AI chat work together seamlessly.
