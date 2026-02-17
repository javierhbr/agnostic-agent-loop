# Plugin System

The agnostic-agent CLI supports a plugin system that lets third parties extend functionality without modifying or recompiling the core CLI. Plugins can provide skill packs, validation rules, custom commands, and orchestrator lifecycle hooks.

---

## Design Approach

The system uses a **hybrid model**:

- **File-based** for content (skill packs, templates) — drop markdown files into a directory
- **Subprocess-based** for logic (validators, commands, hooks) — any language, JSON protocol
- **Go SDK** available at `pkg/plugin/` for type-safe plugin development

This mirrors how the CLI already works internally: skill packs are embedded file bundles while validation rules are Go interfaces with structured results. The plugin system externalizes both patterns.

---

## Plugin Location

Plugins are directories containing a `plugin.yaml` manifest:

- **Project-level**: `.agentic/plugins/<name>/`
- **Global (user-level)**: `~/.agentic/plugins/<name>/`

Discovery order: project-level first, then global. Both paths are configurable via `agnostic-agent.yaml`.

---

## Plugin Manifest

Every plugin requires a `plugin.yaml` at its root:

```yaml
apiVersion: v1
kind: plugin
metadata:
  name: team-review
  version: 1.0.0
  description: "Team code review standards and pre-completion checks"
  author: "platform-team"
  minCliVersion: "0.5.0"    # optional: minimum CLI version required

capabilities:
  # File-based: skill packs (same format as built-in packs)
  skillPacks:
    - name: review-checklist
      description: "Code review checklist for team standards"
      entrypoint: packs/review-checklist/SKILL.md
      resources:
        - packs/review-checklist/resources/checklist-template.md

  # Executable: custom CLI commands
  commands:
    - name: lint-check
      description: "Run team lint standards against current scope"
      entrypoint: bin/lint-check
      parent: validate          # becomes: agentic-agent validate lint-check
      flags:
        - name: severity
          type: string
          default: "warning"
          description: "Minimum severity to report"

  # Executable: validation rules
  validators:
    - name: review-standards
      description: "Validate code meets team review standards"
      entrypoint: bin/validate-review

  # Executable: orchestrator lifecycle hooks
  hooks:
    - event: pre-claim
      entrypoint: bin/pre-claim-hook
    - event: post-complete
      entrypoint: bin/post-complete-hook
```

All `entrypoint` paths are relative to the plugin directory. Each capability section is optional — a plugin can provide just skill packs, just validators, or any combination.

---

## Subprocess Protocol

Executable plugins (commands, validators, hooks) communicate via JSON over stdin/stdout. Stderr is forwarded to the CLI's stderr for logging and debugging.

### Input (stdin)

Every invocation receives this JSON structure:

```json
{
  "plugin": {
    "name": "team-review",
    "version": "1.0.0",
    "dir": "/absolute/path/to/.agentic/plugins/team-review"
  },
  "project": {
    "root": "/absolute/path/to/project",
    "name": "my-app",
    "version": "1.0.0"
  },
  "agent": {
    "name": "claude-code",
    "source": "filesystem"
  },
  "config": {
    "paths": {
      "specDirs": [".agentic/spec"],
      "contextDirs": [".agentic/context"],
      "trackDir": ".agentic/tracks"
    },
    "workflow": {
      "validators": ["task-size", "task-scope"]
    }
  },
  "invocation": {
    "type": "validator",
    "name": "review-standards"
  }
}
```

The `invocation` field varies by type:

| Type | Fields |
| --- | --- |
| command | `type: "command"`, `name`, `args: []string`, `flags: {}` |
| validator | `type: "validator"`, `name` |
| hook | `type: "hook"`, `event: "pre-claim"`, `eventData: {}` |

### Output (stdout)

Each invocation type has its own response format.

**Command response:**

```json
{
  "exitCode": 0,
  "output": "Lint check passed: 15 files checked, 0 issues found"
}
```

**Validator response** (matches the internal `RuleResult` struct):

```json
{
  "ruleName": "review-standards",
  "status": "PASS",
  "errors": []
}
```

Status values: `"PASS"`, `"FAIL"`, `"WARN"`.

**Hook response:**

```json
{
  "continue": true,
  "message": "Pre-claim check passed"
}
```

For pre-* hooks, `"continue": false` aborts the operation. Post-* hooks cannot abort.

### Timeouts

All subprocess invocations enforce timeouts to prevent buggy plugins from hanging the CLI:

| Type | Default timeout |
| --- | --- |
| Commands | 30 seconds |
| Validators | 10 seconds |
| Hooks | 10 seconds |

---

## Extension Points

### Skill Packs

Plugin skill packs use the same directory layout as built-in packs. They register into the existing `PackRegistry` and can be installed to any supported agent tool's skill directory (`.claude/skills/`, `.cursor/skills/`, etc.).

No executable needed — just markdown files. The plugin's `packs/` directory mirrors the structure of the built-in `internal/skills/packs/` directory. Each pack has an entrypoint `SKILL.md` and optional resource files.

When the plugin loader discovers a skill pack in a manifest, it constructs a `SkillPack` struct and registers it in the global `PackRegistry`. The installer reads from the plugin's filesystem instead of the embedded FS.

### Validation Rules

Plugin validators implement the same contract as the built-in `ValidationRule` interface, but via subprocess. The CLI creates a `PluginValidationRule` adapter that:

1. Serializes the `ValidationContext` (project root, config) to JSON
2. Executes the plugin binary with JSON on stdin
3. Parses the JSON response into a `RuleResult`
4. Returns the result to the validator framework

Plugin validators run alongside built-in validators during `agentic-agent validate`. They appear in output with a `[plugin:name]` prefix for traceability.

### Custom Commands

Plugin commands appear as subcommands in the CLI tree. A command declaring `parent: validate` becomes `agentic-agent validate lint-check`. Commands with no parent attach to the root command.

The CLI creates a Cobra `Command` for each plugin command, registers declared flags, and delegates execution to the plugin binary via the subprocess protocol.

Plugin commands always run in flag mode (no interactive TUI). The binary handles its own output formatting. If a plugin needs interactive behavior, it implements that within its own binary.

### Orchestrator Hooks

Hooks fire at lifecycle points in the autopilot/orchestrator loop:

| Event | When it fires | Can abort? |
| --- | --- | --- |
| `pre-claim` | Before claiming a task | Yes |
| `post-claim` | After a task is claimed | No |
| `pre-complete` | Before marking a task complete | Yes |
| `post-complete` | After task completion | No |
| `pre-validate` | Before validation runs | Yes |
| `on-iteration` | Each autopilot iteration | No |

Pre-hooks receive the operation details in `eventData` and can prevent the operation by returning `"continue": false`. When any pre-hook aborts, the CLI prints the hook's message and skips the operation.

Multiple hooks for the same event run in registration order (project-level plugins first, then global).

---

## Go SDK

For Go plugin authors, the `pkg/plugin/` package provides type-safe helpers that handle the JSON protocol automatically:

```go
package main

import "github.com/javierbenavides/agentic-agent/pkg/plugin"

type myValidator struct{}

func (v *myValidator) Validate(ctx *plugin.Context) (*plugin.ValidatorResult, error) {
    // Access project info, agent details, config
    root := ctx.Project.Root
    agent := ctx.Agent.Name

    // Custom validation logic
    if someCheck(root) {
        return &plugin.ValidatorResult{
            RuleName: "my-rule",
            Status:   plugin.StatusPass,
        }, nil
    }

    return &plugin.ValidatorResult{
        RuleName: "my-rule",
        Status:   plugin.StatusFail,
        Errors:   []string{"Validation failed: reason"},
    }, nil
}

func main() {
    plugin.RunValidator(&myValidator{})
}
```

The SDK provides three entry points:

- `plugin.RunValidator(handler)` — for validator plugins
- `plugin.RunCommand(handler)` — for command plugins
- `plugin.RunHook(handler)` — for hook plugins

Each reads JSON from stdin, deserializes into a typed `Context` struct, calls the handler, serializes the response, and writes to stdout. Plugin authors compile their code to a binary and reference it in `plugin.yaml`.

---

## Configuration

Plugin discovery and loading is configured in `agnostic-agent.yaml`:

```yaml
plugins:
  # Directories to scan for plugins
  dirs:
    - .agentic/plugins        # project-level (default)
    - ~/.agentic/plugins       # global (default)

  # Allowlist: if specified, only these plugins load
  enabled:
    - team-review
    - custom-lint

  # Blocklist: takes precedence over enabled
  disabled:
    - experimental-plugin
```

If `enabled` is omitted, all discovered plugins load. The `disabled` list always takes precedence.

Plugins are discovered during CLI startup (`PersistentPreRunE`), after config loading and agent detection. Manifests are parsed eagerly but plugin binaries are invoked lazily — only when their trigger fires.

---

## CLI Management Commands

```text
agentic-agent plugin list                  # Show all plugins, status, capabilities
agentic-agent plugin install <path|url>    # Install from local path or git URL
agentic-agent plugin remove <name>         # Remove a plugin
agentic-agent plugin info <name>           # Show plugin manifest details
agentic-agent plugin init <name>           # Scaffold a new plugin directory
```

### `plugin list` output

```text
NAME              VERSION  STATUS    CAPABILITIES
team-review       1.0.0    enabled   validators(1), hooks(2), skill-packs(1)
custom-lint       0.3.0    enabled   commands(1), validators(2)
experimental      0.1.0    disabled  commands(2)
```

### `plugin install`

Installs from a local path (copies directory) or git URL (clones repository) into `.agentic/plugins/`. Use `--global` to install to `~/.agentic/plugins/` instead.

```bash
# From local directory
agentic-agent plugin install ./path/to/my-plugin

# From git repository
agentic-agent plugin install https://github.com/org/agentic-review-plugin.git

# Install globally
agentic-agent plugin install ./my-plugin --global
```

### `plugin init`

Scaffolds a new plugin directory with a template manifest and standard structure:

```bash
agentic-agent plugin init my-custom-rules
```

Creates:

```text
.agentic/plugins/my-custom-rules/
  plugin.yaml          # manifest template
  bin/                  # for executable plugins
  packs/               # for skill pack files
```

---

## Examples

### Validator Plugin in Bash

A plugin that checks for unresolved TODO comments in staged files:

**Directory structure:**

```text
.agentic/plugins/todo-checker/
  plugin.yaml
  bin/
    validate-todos
```

**plugin.yaml:**

```yaml
apiVersion: v1
kind: plugin
metadata:
  name: todo-checker
  version: 1.0.0
  description: "Checks for unresolved TODO comments"

capabilities:
  validators:
    - name: no-todos
      description: "Fail if staged files contain TODO comments"
      entrypoint: bin/validate-todos
```

**bin/validate-todos:**

```bash
#!/usr/bin/env bash
set -euo pipefail

# Read JSON context from stdin
CONTEXT=$(cat)
PROJECT_ROOT=$(echo "$CONTEXT" | jq -r '.project.root')

ERRORS=()

# Check for TODO comments in staged files
TODOS=$(git -C "$PROJECT_ROOT" diff --cached --name-only | xargs grep -l "TODO" 2>/dev/null || true)
if [ -n "$TODOS" ]; then
    while IFS= read -r file; do
        ERRORS+=("Unresolved TODO in $file")
    done <<< "$TODOS"
fi

# Output JSON response
if [ ${#ERRORS[@]} -eq 0 ]; then
    echo '{"ruleName":"no-todos","status":"PASS","errors":[]}'
else
    ERRS=$(printf '"%s",' "${ERRORS[@]}")
    echo "{\"ruleName\":\"no-todos\",\"status\":\"WARN\",\"errors\":[${ERRS%,}]}"
fi
```

### Skill Pack Plugin

A plugin that provides team coding conventions as a skill pack (no executable needed):

**Directory structure:**

```text
.agentic/plugins/team-conventions/
  plugin.yaml
  packs/
    conventions/
      SKILL.md
      resources/
        naming-rules.md
        review-checklist.md
```

**plugin.yaml:**

```yaml
apiVersion: v1
kind: plugin
metadata:
  name: team-conventions
  version: 1.0.0
  description: "Team coding conventions and review standards"

capabilities:
  skillPacks:
    - name: conventions
      description: "Team coding conventions"
      entrypoint: packs/conventions/SKILL.md
      resources:
        - packs/conventions/resources/naming-rules.md
        - packs/conventions/resources/review-checklist.md
```

### Hook Plugin in Python

A pre-completion hook that runs a custom check before a task is marked done:

**bin/pre-complete-hook:**

```python
#!/usr/bin/env python3
import json
import sys
import subprocess

context = json.load(sys.stdin)
project_root = context["project"]["root"]
task_id = context["invocation"].get("eventData", {}).get("taskId", "")

# Check that tests pass before allowing completion
result = subprocess.run(
    ["go", "test", "./..."],
    cwd=project_root,
    capture_output=True,
    text=True
)

if result.returncode == 0:
    json.dump({"continue": True, "message": f"Tests pass for {task_id}"}, sys.stdout)
else:
    json.dump({
        "continue": False,
        "message": f"Tests failing — cannot complete {task_id}: {result.stderr[:200]}"
    }, sys.stdout)
```

### Validator Plugin in Go (using SDK)

```go
package main

import (
    "os"
    "path/filepath"

    "github.com/javierbenavides/agentic-agent/pkg/plugin"
)

type contextValidator struct{}

func (v *contextValidator) Validate(ctx *plugin.Context) (*plugin.ValidatorResult, error) {
    // Check that every directory with .go files has a context.md
    missing := []string{}

    entries, _ := os.ReadDir(filepath.Join(ctx.Project.Root, "internal"))
    for _, entry := range entries {
        if entry.IsDir() {
            ctxPath := filepath.Join(ctx.Project.Root, "internal", entry.Name(), "AGENTS.md")
            if _, err := os.Stat(ctxPath); os.IsNotExist(err) {
                missing = append(missing, "internal/"+entry.Name())
            }
        }
    }

    if len(missing) == 0 {
        return &plugin.ValidatorResult{
            RuleName: "context-coverage",
            Status:   plugin.StatusPass,
        }, nil
    }

    errors := make([]string, len(missing))
    for i, dir := range missing {
        errors[i] = "Missing context.md in " + dir
    }

    return &plugin.ValidatorResult{
        RuleName: "context-coverage",
        Status:   plugin.StatusFail,
        Errors:   errors,
    }, nil
}

func main() {
    plugin.RunValidator(&contextValidator{})
}
```

---

## Implementation Architecture

### Internal Components

The plugin system introduces these internal packages:

| Package | Purpose |
| --- | --- |
| `pkg/models/plugin.go` | Manifest structs (`PluginManifest`, capabilities, etc.) |
| `pkg/plugin/` | Go SDK for plugin authors |
| `internal/plugins/loader.go` | Discovery: scans directories for `plugin.yaml`, parses manifests |
| `internal/plugins/executor.go` | Subprocess protocol: JSON stdin/stdout, timeout enforcement |
| `internal/plugins/registry.go` | Central registry: lookups by capability type, enable/disable |
| `internal/plugins/validator.go` | Adapter: wraps subprocess as `validator.ValidationRule` |
| `internal/plugins/command.go` | Factory: creates Cobra commands from plugin manifest |
| `internal/plugins/hooks.go` | Hook runner: fires events, handles abort logic |
| `cmd/agentic-agent/plugin.go` | CLI commands: list, install, remove, info, init |

### Integration Points

The plugin loader integrates with existing infrastructure at these points:

| Where | What happens |
| --- | --- |
| `PersistentPreRunE` in `root.go` | Plugin discovery and manifest loading (after config + agent detection) |
| `PackRegistry` in `packs.go` | Plugin skill packs registered alongside built-in packs |
| `Validator` in `validate.go` | Plugin validators registered alongside built-in rules |
| `rootCmd` / subcommands | Plugin commands attached to the command tree |
| `AutopilotLoop` in `autopilot.go` | Hook runner fires events at lifecycle points |
| `Config` in `config.go` | `PluginsConfig` added for discovery paths and enable/disable lists |

### Key Design: SkillPackFile FS Field

The existing `SkillPackFile` struct gains an `FS io/fs.FS` field. Built-in packs default to the embedded `packsFS`. Plugin packs set this to `os.DirFS(pluginDir)`. The installer reads from whichever FS is set, making the change backwards-compatible.

---

## Design Decisions

| Decision | Rationale |
| --- | --- |
| Hybrid (files + subprocess) | Skill packs are naturally file-based content; validators and hooks need executable logic. Matching both patterns avoids forcing all plugins into one model. |
| JSON protocol (not gRPC/protobuf) | JSON is supported in every language's standard library. The existing `RuleResult` struct already has `json` tags. Shell scripts can emit JSON with `echo`. No code generation needed. |
| Go SDK alongside JSON protocol | Go authors get type safety and a one-liner entry point. Bash/Python authors use raw JSON. Both compile to the same binary interface. |
| No HashiCorp go-plugin | go-plugin requires the same Go version for host and plugin, does not support cross-compilation, and loads via `plugin.Open()` which is platform-specific. The subprocess model is simpler and universal. |
| Plugin dirs inside `.agentic/` | Follows the established convention: tasks live in `.agentic/tasks/`, specs in `.agentic/spec/`, tracks in `.agentic/tracks/`. Plugins in `.agentic/plugins/` can be version-controlled and shared with the team. |
| Eager discovery, lazy execution | Scanning a few YAML files at startup is negligible overhead. Plugin binaries are only invoked when their trigger fires. This avoids startup latency from executing plugins that may not be needed. |
| No dependency resolution (v1) | Adding a version solver (like Go modules or npm) would be premature. The `minCliVersion` field in the manifest provides basic compatibility checking. If two plugins conflict, the user disables one via the config blocklist. |
| Skill packs + validators first | These build directly on the project's strongest existing patterns (`PackRegistry` and `ValidationRule` interface). They deliver the most immediate value — teams diverge most on standards and conventions. Commands and hooks follow once the protocol is proven. |
| 30s timeout for subprocesses | Prevents a buggy or hanging plugin from blocking the entire CLI. Validators and hooks use a shorter 10s timeout since they should be fast checks. |
| Plugin commands use flag mode only | Interactive TUI behavior requires Bubble Tea models that can only be built in Go within the main binary. Plugin commands handle their own output. This keeps the interface simple. |
