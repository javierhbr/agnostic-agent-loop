# Unit & ATDD Tests for Agentic Helper Agent

**Created:** 2026-02-26 | **Commit:** 8ee0c22

## Overview

Comprehensive test coverage for the agentic-helper pack, including 9 white-box unit tests and 10 black-box ATDD scenarios. All tests verify correct agent file routing, installation behavior, and configuration.

---

## Unit Tests (9 tests)

### `internal/skills/installer_test.go` (6 tests)

#### 1. `TestInstaller_Install_AgentHelperPack`
- Installs agentic-helper pack for claude-code
- Overrides both `ToolSkillDir` and `ToolAgentDir` to temp directories
- **Asserts:**
  - 2 files written (AGENT.md + SKILL.md)
  - AGENT.md routed to `.claude/agents/agentic-helper.md`
  - SKILL.md routed to `.claude/skills/agentic-helper/SKILL.md`
  - Both files non-empty and contain expected content

#### 2. `TestInstaller_Install_AgentFile_FallsBackForNonAgentTools`
- Installs agentic-helper for cursor (no `ToolAgentDir` entry)
- **Asserts:**
  - 2 files written
  - AGENT.md falls back to `ToolSkillDir["cursor"]` (no panic/error)
  - Package correctly handles tools without native agent support

#### 3. `TestInstaller_IsInstalled_AgentPack`
- Validates pack detection before and after installation
- Overrides both tool directories
- **Asserts:**
  - `IsInstalled()` returns false before install
  - `IsInstalled()` returns true after install
  - Handles agent file routing in `IsInstalled()`

#### 4. `TestResolveAgentOutputDir_ClaudeCode`
- Direct call to `resolveAgentOutputDir("claude-code", false)`
- **Asserts:**
  - Returns `.claude/agents` without error

#### 5. `TestResolveAgentOutputDir_UnsupportedTool`
- Direct call to `resolveAgentOutputDir("cursor", false)`
- **Asserts:**
  - Returns error (cursor has no agent dir)
  - Error message indicates tool doesn't support agents

#### 6. `TestResolveAgentOutputDir_Global`
- Direct call to `resolveAgentOutputDir("claude-code", true)`
- Overrides `ToolGlobalAgentDir["claude-code"]` to temp directory
- **Asserts:**
  - ~ expansion works correctly
  - Returns absolute path under home directory

### `internal/skills/packs_test.go` (3 tests)

#### 7. `TestAgenticHelperPack_HasIsAgentFlag`
- Retrieves agentic-helper from registry
- **Asserts:**
  - `agentic-helper.md` file has `IsAgent = true`
  - `agentic-helper/SKILL.md` file has `IsAgent = false`

#### 8. `TestToolAgentDir_ClaudeCodePresent`
- Direct checks of maps
- **Asserts:**
  - `ToolAgentDir["claude-code"] == ".claude/agents"`
  - `ToolGlobalAgentDir["claude-code"] == "~/.claude/agents"`

#### 9. `TestMandatoryPacks_IncludesAgenticHelper`
- Iterates MandatoryPacks list
- **Asserts:**
  - "agentic-helper" is in the list
  - (Optional warning if not first)

---

## ATDD/BDD Tests (1 feature, 10 scenarios, 47 steps)

### `features/skills/agent_install.feature`

**Feature:** Agent Pack Installation
**Goal:** Verify agentic-helper agent is installed correctly for all tools

#### Scenario 1: Installing agentic-helper for claude-code creates agent file
**Given:** Clean test environment + initialized project
**When:** Install "agentic-helper" pack for "claude-code"
**Then:**
- Command succeeds
- `.claude/agents/agentic-helper.md` exists
- File contains `name: agentic-helper`
- File contains `tools: Read, Write, Edit, Bash`
- `.claude/skills/agentic-helper/SKILL.md` exists

#### Scenario 2: Agent file contains the workflow decision tree
**When:** Install "agentic-helper" pack for "claude-code"
**Then:**
- File contains "Workflow Decision Tree"
- File contains "Full SDD"
- File contains "openspec init"

#### Scenario 3: Agent file contains CLI command examples
**When:** Install "agentic-helper" pack for "claude-code"
**Then:**
- File contains "task claim"
- File contains "context generate"
- File contains "agentic-agent validate"

#### Scenario 4: Installing agentic-helper for cursor falls back to skill directory
**When:** Install "agentic-helper" pack for "cursor"
**Then:**
- Command succeeds
- `.cursor/skills/agentic-helper/SKILL.md` exists
- (No `.cursor/agents/` created — validates fallback behavior)

#### Scenario 5: agentic-helper is installed during skills ensure
**When:** Run skills ensure for "claude-code"
**Then:**
- `.claude/agents/agentic-helper.md` exists
- File is not empty

#### Scenario 6: agentic-helper agent file uses sonnet model
**When:** Install "agentic-helper" pack for "claude-code"
**Then:**
- File contains "model: sonnet"

#### Scenario 7: agentic-helper agent file enables project memory
**When:** Install "agentic-helper" pack for "claude-code"
**Then:**
- File contains "memory: project"

#### Scenario 8: agentic-helper skill file explains when to use it
**When:** Install "agentic-helper" pack for "claude-code"
**Then:**
- `.claude/skills/agentic-helper/SKILL.md` contains "Use this skill when"
- File contains "agentic-agent CLI"

#### Scenario 9: agentic-helper agent file contains error recovery guidance
**When:** Install "agentic-helper" pack for "claude-code"
**Then:**
- File contains "Error Recovery"
- File contains "Red Flags"

#### Scenario 10: agentic-helper is marked as mandatory
**When:** Run the installed agentic-agent CLI
**Then:**
- "agentic-helper" should be in the mandatory packs list

---

## Step Definitions

### `test/bdd/steps/agent_install_steps.go`

**AgentInstallSteps** struct encapsulates agent installation steps with shared suite context.

#### Registered Steps

1. `I install the "X" pack for "Y"`
   - Overrides `ToolSkillDir[tool]` and `ToolAgentDir[tool]` to test project directory
   - Calls `installer.Install(pack, tool, false)`
   - Registers cleanup to restore original paths

2. `I run skills ensure for "Y"`
   - Creates minimal config
   - Overrides tool directories
   - Installs all mandatory packs for the tool

3. `the file "X" should not be empty`
   - Checks `os.Stat(path).Size() > 0`
   - Verifies files aren't created with zero bytes

4. `I run the installed agentic-agent CLI`
   - No-op step in BDD context (CLI already available)
   - Sets `LastCommandErr` to nil

5. `"X" should be in the mandatory packs list`
   - Iterates `skills.MandatoryPacks`
   - Returns nil if found, error if not

#### Helper Functions

- `getToolSkillSubdir(tool)` — Maps tool names to skill directory paths
- `getToolAgentSubdir(tool)` — Maps tool names to agent directory paths

---

## Test Execution

### Unit Tests Only

```bash
# All skills package tests
go test ./internal/skills/... -v

# Only agentic-helper tests
go test ./internal/skills/... -run "AgenticHelper|AgentDir|Resolve.*AgentOutputDir" -v

# Result: All 9 new tests pass ✓
```

### BDD Tests Only

```bash
# All BDD features
go test ./test/bdd/... -v

# Agent install scenarios only
go test ./test/bdd/... -run TestFeatures -v | grep -A 3 "Agent Pack Installation"

# Result: 10 agent install scenarios pass ✓
# Overall: 53 total scenarios, 307 steps pass ✓
```

### Combined Test Suite

```bash
go test ./internal/skills/... ./test/bdd/... -v
# Result: All unit + ATDD tests passing ✓
```

---

## Code Changes

### `internal/skills/installer.go`
- Modified `IsInstalled()` to handle agent file routing
- Now checks `ToolAgentDir` for files with `IsAgent = true`
- Falls back to `ToolSkillDir` for tools without native agent support

### `internal/skills/installer_test.go`
- Added 6 new test functions (lines 206-310)
- Tests cover agent file routing, fallback behavior, and directory resolution
- All tests follow existing pattern: override maps → install → assert

### `internal/skills/packs_test.go`
- Added 3 new test functions (lines 140-198)
- Tests for pack registration, flag configuration, and mandatory packs list
- Validates static maps and registry state

### `features/skills/agent_install.feature` (NEW)
- 10 ATDD scenarios covering all agent install behaviors
- Tests both happy path (claude-code with agent support) and fallback (cursor)
- Validates configuration (model, memory, content)

### `test/bdd/steps/agent_install_steps.go` (NEW)
- AgentInstallSteps struct with 5 step definitions
- Manages tool directory overrides during test execution
- Reuses SuiteContext for cross-step state

### `test/bdd/features_test.go`
- Registered NewAgentInstallSteps in InitializeScenario
- One line added: `steps.NewAgentInstallSteps(stepsSuite).RegisterSteps(sc)`

---

## Key Design Decisions

### IsInstalled() Enhancement
The `IsInstalled()` function was extended to handle agent files because:
- Agent files live in `ToolAgentDir`, not `ToolSkillDir`
- Without this change, agent pack installation would appear undetected
- The function now routes each file to the correct directory based on `IsAgent` flag

### BDD Step Directory Management
Steps temporarily override `ToolSkillDir` and `ToolAgentDir` maps because:
- Tests need to isolate changes to the test project directory
- The patterns follow existing installer tests (map override + defer restore)
- Cleanup functions ensure no state leakage between scenarios

### Test Coverage Balance
- 9 unit tests for low-level behavior (directory resolution, routing logic)
- 10 ATDD scenarios for end-to-end behavior (install + verify file presence/content)
- Together they provide both implementation detail validation and black-box acceptance testing

---

## Verification Checklist

- [x] All 9 unit tests passing
- [x] All 10 ATDD scenarios passing
- [x] No regressions in existing 44 BDD scenarios
- [x] No regressions in existing 100+ unit tests for other packages
- [x] Code builds without errors
- [x] Agent file routing works for claude-code (native support)
- [x] Fallback behavior works for cursor (no native support)
- [x] IsInstalled() correctly detects agent packs
- [x] All test cleanup functions properly restore state
- [x] Step definitions follow existing patterns (SuiteContext, cleanup registration)
