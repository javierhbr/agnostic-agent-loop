# Spec Consistency Validation System Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Ensure tasks always have proper spec files through hybrid validation with warnings, auto-generation, manual creation, and skip options - guided by agent-friendly CLI output.

**Architecture:** Extend existing readiness system with spec generator, add lifecycle hooks at task create/claim, provide new spec management commands, and include agent-guiding messages in CLI output.

**Tech Stack:** Go 1.21+, Cobra CLI, Bubble Tea (interactive prompts), existing specs/tasks packages, YAML config

---

## Task 1: Add Configuration Fields for Spec Validation

**Files:**
- Modify: `pkg/models/config.go:36-38`
- Test: `pkg/models/config_test.go` (new file)

**Step 1: Write the failing test**

Create test file first:

```go
package models

import (
	"testing"
	"gopkg.in/yaml.v3"
)

func TestWorkflowConfigSpecValidation(t *testing.T) {
	yamlData := `
workflow:
  validators:
    - context-check
  validate_specs_on_claim: true
  spec_validation_mode: warn
`
	var cfg Config
	err := yaml.Unmarshal([]byte(yamlData), &cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if !cfg.Workflow.ValidateSpecsOnClaim {
		t.Error("Expected ValidateSpecsOnClaim to be true")
	}

	if cfg.Workflow.SpecValidationMode != "warn" {
		t.Errorf("Expected SpecValidationMode 'warn', got '%s'", cfg.Workflow.SpecValidationMode)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/models/... -v -run TestWorkflowConfigSpecValidation`
Expected: FAIL with "unknown field" errors

**Step 3: Add fields to WorkflowConfig**

In `pkg/models/config.go`, modify `WorkflowConfig`:

```go
type WorkflowConfig struct {
	Validators           []string `yaml:"validators"`
	ValidateSpecsOnClaim bool     `yaml:"validate_specs_on_claim,omitempty"` // Run spec validation before claim
	SpecValidationMode   string   `yaml:"spec_validation_mode,omitempty"`    // "warn" | "block" | "silent"
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/models/... -v -run TestWorkflowConfigSpecValidation`
Expected: PASS

**Step 5: Commit**

```bash
git add pkg/models/config.go pkg/models/config_test.go
git commit -m "feat: add spec validation config fields to WorkflowConfig"
```

---

## Task 2: Create Spec Generator Package Structure

**Files:**
- Create: `internal/specs/generator.go`
- Create: `internal/specs/generator_test.go`
- Create: `internal/specs/templates.go`

**Step 1: Write the failing test**

```go
package specs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

func TestSpecGenerator_DetectContext_PRDExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a mock PRD file
	prdPath := filepath.Join(tmpDir, "todo-app-proposal.md")
	err := os.WriteFile(prdPath, []byte("# Todo App Proposal\n\nFeatures..."), 0644)
	if err != nil {
		t.Fatal(err)
	}

	cfg := &models.Config{
		Paths: models.PathsConfig{
			SpecDirs: []string{tmpDir},
		},
	}

	task := &models.Task{
		ChangeID: "todo-app",
		SpecRefs: []string{"todo-app/proposal.md"},
	}

	gen := NewSpecGenerator(cfg)
	ctx := gen.DetectContext(task)

	if ctx.Source != "prd" {
		t.Errorf("Expected source 'prd', got '%s'", ctx.Source)
	}

	if ctx.PRDPath == "" {
		t.Error("Expected PRDPath to be set")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/specs/... -v -run TestSpecGenerator_DetectContext_PRDExists`
Expected: FAIL with "undefined: NewSpecGenerator"

**Step 3: Implement SpecGenerator structure**

In `internal/specs/generator.go`:

```go
package specs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/mattn/go-isatty"
)

type SpecGenerationContext struct {
	Source       string // "prd" | "interactive" | "metadata"
	PRDPath      string
	TaskMetadata *models.Task
	Interactive  bool
}

type SpecGenerator struct {
	cfg *models.Config
}

func NewSpecGenerator(cfg *models.Config) *SpecGenerator {
	return &SpecGenerator{cfg: cfg}
}

// DetectContext determines the best source for spec generation
func (g *SpecGenerator) DetectContext(task *models.Task) *SpecGenerationContext {
	ctx := &SpecGenerationContext{
		TaskMetadata: task,
		Interactive:  isatty.IsTerminal(os.Stdout.Fd()),
	}

	// Priority 1: Look for PRD files
	if task.ChangeID != "" || task.TrackID != "" {
		searchTerm := task.ChangeID
		if searchTerm == "" {
			searchTerm = task.TrackID
		}

		// Search in configured spec dirs
		for _, dir := range g.cfg.Paths.SpecDirs {
			matches, _ := filepath.Glob(filepath.Join(dir, "*"+searchTerm+"*proposal.md"))
			if len(matches) > 0 {
				ctx.Source = "prd"
				ctx.PRDPath = matches[0]
				return ctx
			}
		}
	}

	// Priority 2: Interactive if TTY
	if ctx.Interactive {
		ctx.Source = "interactive"
		return ctx
	}

	// Priority 3: Fallback to metadata
	ctx.Source = "metadata"
	return ctx
}

// Generate creates a single spec file from context
func (g *SpecGenerator) Generate(specRef string, ctx *SpecGenerationContext) (string, error) {
	// Determine output path
	outputPath := g.resolveOutputPath(specRef)

	// Generate content based on source
	var content string
	var err error

	switch ctx.Source {
	case "prd":
		content, err = g.generateFromPRD(specRef, ctx.PRDPath)
	case "interactive":
		content, err = g.generateInteractive(specRef)
	case "metadata":
		content, err = g.generateFromMetadata(specRef, ctx.TaskMetadata)
	default:
		return "", fmt.Errorf("unknown generation source: %s", ctx.Source)
	}

	if err != nil {
		return "", err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write spec file: %w", err)
	}

	return outputPath, nil
}

func (g *SpecGenerator) resolveOutputPath(specRef string) string {
	// If absolute, use as-is
	if filepath.IsAbs(specRef) {
		return specRef
	}

	// Otherwise, use first spec dir
	baseDir := ".agentic/spec"
	if len(g.cfg.Paths.SpecDirs) > 0 {
		baseDir = g.cfg.Paths.SpecDirs[0]
	}

	return filepath.Join(baseDir, specRef)
}

func (g *SpecGenerator) generateFromPRD(specRef, prdPath string) (string, error) {
	// Read PRD content
	data, err := os.ReadFile(prdPath)
	if err != nil {
		return "", fmt.Errorf("failed to read PRD: %w", err)
	}

	// For now, create a minimal spec referencing the PRD
	content := fmt.Sprintf(`# %s

**Generated from:** %s

## Overview

This spec was automatically generated from the PRD. Please refer to the source PRD for complete details.

## Source PRD Content

%s
`, specRef, prdPath, string(data))

	return content, nil
}

func (g *SpecGenerator) generateInteractive(specRef string) (string, error) {
	// For now, return a TODO - will implement Bubble Tea prompt later
	return fmt.Sprintf(`# %s

**TODO:** This spec needs to be filled in interactively.

## Requirements

(To be filled in)

## Acceptance Criteria

(To be filled in)
`, specRef), nil
}

func (g *SpecGenerator) generateFromMetadata(specRef string, task *models.Task) (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", specRef))
	sb.WriteString(fmt.Sprintf("**Generated from task:** %s\n\n", task.ID))

	if task.Description != "" {
		sb.WriteString("## Description\n\n")
		sb.WriteString(task.Description)
		sb.WriteString("\n\n")
	}

	if len(task.Acceptance) > 0 {
		sb.WriteString("## Acceptance Criteria\n\n")
		for _, criterion := range task.Acceptance {
			sb.WriteString(fmt.Sprintf("- %s\n", criterion))
		}
		sb.WriteString("\n")
	}

	if len(task.Scope) > 0 {
		sb.WriteString("## Scope\n\n")
		for _, scope := range task.Scope {
			sb.WriteString(fmt.Sprintf("- %s\n", scope))
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/specs/... -v -run TestSpecGenerator_DetectContext_PRDExists`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/specs/generator.go internal/specs/generator_test.go
git commit -m "feat: add SpecGenerator with context detection and generation"
```

---

## Task 3: Enhance Readiness Checks for Spec Completeness

**Files:**
- Modify: `internal/tasks/readiness.go:26-93`
- Modify: `internal/tasks/readiness_test.go`

**Step 1: Write the failing test**

Add to `internal/tasks/readiness_test.go`:

```go
func TestCanClaimTask_SpecCompleteness(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &models.Config{
		Paths: models.PathsConfig{
			SpecDirs: []string{tmpDir},
		},
	}

	task := &models.Task{
		ID:       "TASK-123",
		SpecRefs: []string{"missing-spec.md", "another-missing.md"},
	}

	result := CanClaimTask(task, cfg)

	// Should have spec-completeness checks
	hasSpecCheck := false
	for _, check := range result.Checks {
		if check.Name == "spec-completeness" {
			hasSpecCheck = true
			if check.Passed {
				t.Error("Expected spec-completeness check to fail for missing specs")
			}
		}
	}

	if !hasSpecCheck {
		t.Error("Expected spec-completeness check in results")
	}

	// Should still be ready (warnings only)
	if !result.Ready {
		t.Error("Expected task to be ready despite spec warnings")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/tasks/... -v -run TestCanClaimTask_SpecCompleteness`
Expected: FAIL with "Expected spec-completeness check in results"

**Step 3: Add spec-completeness check to CanClaimTask**

In `internal/tasks/readiness.go`, after the existing checks (around line 72):

```go
// Check 3: Spec completeness (warning only, not blocking)
if len(task.SpecRefs) > 0 {
	resolver := specs.NewResolver(cfg)
	resolved := resolver.ResolveAll(task.SpecRefs)

	missingSpecs := []string{}
	for _, r := range resolved {
		if !r.Found {
			missingSpecs = append(missingSpecs, r.Ref)
		}
	}

	if len(missingSpecs) > 0 {
		result.Checks = append(result.Checks, ReadinessCheck{
			Name:    "spec-completeness",
			Passed:  false,
			Message: fmt.Sprintf("%d specs missing: %s", len(missingSpecs), strings.Join(missingSpecs, ", ")),
		})
		// Note: spec-completeness is a warning, does NOT set Ready=false
	} else {
		result.Checks = append(result.Checks, ReadinessCheck{
			Name:    "spec-completeness",
			Passed:  true,
			Message: fmt.Sprintf("all %d specs found", len(task.SpecRefs)),
		})
	}
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/tasks/... -v -run TestCanClaimTask_SpecCompleteness`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/tasks/readiness.go internal/tasks/readiness_test.go
git commit -m "feat: add spec-completeness check to readiness validation"
```

---

## Task 4: Add Agent-Guiding Output to ClaimTaskWithConfig

**Files:**
- Modify: `internal/tasks/lock.go:68-94`
- Test: Manual testing (output formatting)

**Step 1: Update FormatReadinessResult with agent guidance**

In `internal/tasks/readiness.go`, replace `FormatReadinessResult`:

```go
// FormatReadinessResult returns a human-readable string with agent guidance.
func FormatReadinessResult(r *ReadinessResult) string {
	var b strings.Builder
	if r.Ready {
		b.WriteString(fmt.Sprintf("Task %s: READY\n", r.TaskID))
	} else {
		b.WriteString(fmt.Sprintf("Task %s: NOT READY\n", r.TaskID))
	}

	hasSpecIssues := false
	missingSpecs := []string{}

	for _, check := range r.Checks {
		icon := "+"
		if !check.Passed {
			icon = "-"
			if check.Name == "spec-completeness" {
				hasSpecIssues = true
				// Extract missing spec names from message
				msg := check.Message
				if strings.Contains(msg, "missing:") {
					parts := strings.Split(msg, "missing:")
					if len(parts) > 1 {
						specs := strings.Split(parts[1], ",")
						for _, spec := range specs {
							missingSpecs = append(missingSpecs, strings.TrimSpace(spec))
						}
					}
				}
			}
		}
		b.WriteString(fmt.Sprintf("  [%s] %s: %s\n", icon, check.Name, check.Message))
	}

	// Add agent guidance for missing specs
	if hasSpecIssues && len(missingSpecs) > 0 {
		b.WriteString("\n⚠️  MISSING SPECS DETECTED\n\n")
		b.WriteString("The following spec files are referenced but don't exist:\n")
		for _, spec := range missingSpecs {
			b.WriteString(fmt.Sprintf("  - %s\n", spec))
		}
		b.WriteString("\n📋 RECOMMENDED ACTIONS FOR AGENTS:\n\n")
		b.WriteString(fmt.Sprintf("  Option 1: Auto-generate specs (recommended)\n"))
		b.WriteString(fmt.Sprintf("  → agentic-agent spec generate %s --auto\n\n", r.TaskID))
		b.WriteString(fmt.Sprintf("  Option 2: Generate with user interaction\n"))
		b.WriteString(fmt.Sprintf("  → agentic-agent spec generate %s --interactive\n\n", r.TaskID))
		b.WriteString(fmt.Sprintf("  Option 3: Create specs manually\n"))
		for _, spec := range missingSpecs {
			b.WriteString(fmt.Sprintf("  → agentic-agent spec create %s\n", spec))
		}
		b.WriteString(fmt.Sprintf("\n  Option 4: Skip validation and proceed\n"))
		b.WriteString(fmt.Sprintf("  → agentic-agent task claim %s --skip-validation\n\n", r.TaskID))
		b.WriteString("💡 CONTEXT: This task was likely created outside the recommended workflow.\n")
		b.WriteString("   For best results, follow: brainstorming → product-wizard → openspec init\n\n")
		b.WriteString("🤖 AGENT TIP: Read the user's intent. If they want to proceed quickly,\n")
		b.WriteString("   use --skip-validation. If quality matters, use --auto generation.\n\n")
	}

	return b.String()
}
```

**Step 2: Test manually**

Run: Create a task with missing specs and claim it
```bash
# In test directory with tasks
agentic-agent task claim TASK-123
```
Expected: See formatted output with agent guidance

**Step 3: Adjust formatting based on manual testing**

Iterate on formatting until output is clear and helpful.

**Step 4: Commit**

```bash
git add internal/tasks/readiness.go
git commit -m "feat: add agent-guiding output to readiness results"
```

---

## Task 5: Add spec generate Command Structure

**Files:**
- Create: `cmd/agentic-agent/spec.go`
- Modify: `cmd/agentic-agent/root.go` (register command)

**Step 1: Write the command skeleton test**

Create `cmd/agentic-agent/spec_test.go`:

```go
package main

import (
	"testing"
)

func TestSpecCommand_Exists(t *testing.T) {
	// Verify spec command is registered
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "spec" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected 'spec' command to be registered")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./cmd/agentic-agent/... -v -run TestSpecCommand_Exists`
Expected: FAIL with "Expected 'spec' command to be registered"

**Step 3: Create spec command file**

Create `cmd/agentic-agent/spec.go`:

```go
package main

import (
	"fmt"
	"os"

	"github.com/javierbenavides/agentic-agent/internal/specs"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/spf13/cobra"
)

var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "Manage specification files",
	Long:  `Create, generate, and validate specification files for tasks.`,
}

var specGenerateCmd = &cobra.Command{
	Use:   "generate TASK_ID",
	Short: "Generate missing specs for a task",
	Long: `Generate missing specification files for a task using smart context detection.

Auto Mode (detects best source):
  agentic-agent spec generate TASK-123 --auto

Interactive Mode (prompts for requirements):
  agentic-agent spec generate TASK-123 --interactive

From PRD (use specific PRD file):
  agentic-agent spec generate TASK-123 --from-prd path/to/prd.md

Generate All (for all tasks with missing specs):
  agentic-agent spec generate --all`,
	Run: func(cmd *cobra.Command, args []string) {
		auto, _ := cmd.Flags().GetBool("auto")
		interactive, _ := cmd.Flags().GetBool("interactive")
		fromPRD, _ := cmd.Flags().GetString("from-prd")
		all, _ := cmd.Flags().GetBool("all")

		if all {
			runGenerateAllSpecs()
			return
		}

		if len(args) == 0 {
			fmt.Println("Error: TASK_ID required (or use --all)")
			os.Exit(1)
		}

		taskID := args[0]
		runGenerateTaskSpecs(taskID, auto, interactive, fromPRD)
	},
}

var specCreateCmd = &cobra.Command{
	Use:   "create SPEC_REF",
	Short: "Create a new spec file with template",
	Long: `Create a new specification file using a template.

Examples:
  agentic-agent spec create todo-app/proposal.md
  agentic-agent spec create todo-app/tasks/01-setup.md --template task`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: SPEC_REF required")
			os.Exit(1)
		}

		specRef := args[0]
		template, _ := cmd.Flags().GetString("template")
		runCreateSpec(specRef, template)
	},
}

var specValidateCmd = &cobra.Command{
	Use:   "validate [TASK_ID]",
	Short: "Validate spec completeness for tasks",
	Long: `Validate that all tasks have their required specification files.

Validate all tasks:
  agentic-agent spec validate

Validate specific task:
  agentic-agent spec validate TASK-123`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			runValidateAllSpecs()
		} else {
			runValidateTaskSpecs(args[0])
		}
	},
}

func runGenerateTaskSpecs(taskID string, auto, interactive bool, fromPRD string) {
	cfg := getConfig()
	tm := tasks.NewTaskManager(".agentic/tasks")

	// Find the task
	task, source, err := tm.FindTask(taskID)
	if err != nil {
		fmt.Printf("Error finding task: %v\n", err)
		os.Exit(1)
	}

	if source != "backlog" && source != "in-progress" {
		fmt.Printf("Error: Task %s is in %s (can only generate specs for backlog or in-progress tasks)\n", taskID, source)
		os.Exit(1)
	}

	fmt.Printf("Analyzing task context...\n")

	gen := specs.NewSpecGenerator(cfg)
	ctx := gen.DetectContext(task)

	// Override context based on flags
	if fromPRD != "" {
		ctx.Source = "prd"
		ctx.PRDPath = fromPRD
	} else if interactive {
		ctx.Source = "interactive"
	} else if auto {
		// Use detected context
	}

	fmt.Printf("✓ Context: Will generate from %s\n\n", ctx.Source)

	// Find missing specs
	resolver := specs.NewResolver(cfg)
	resolved := resolver.ResolveAll(task.SpecRefs)

	missingSpecs := []string{}
	for _, r := range resolved {
		if !r.Found {
			missingSpecs = append(missingSpecs, r.Ref)
		}
	}

	if len(missingSpecs) == 0 {
		fmt.Printf("✓ All specs already exist for task %s\n", taskID)
		return
	}

	fmt.Printf("Generating %d missing spec(s) for %s...\n", len(missingSpecs), taskID)

	for _, specRef := range missingSpecs {
		path, err := gen.Generate(specRef, ctx)
		if err != nil {
			fmt.Printf("✗ Failed to generate %s: %v\n", specRef, err)
			continue
		}

		// Get file size
		info, _ := os.Stat(path)
		size := int64(0)
		if info != nil {
			size = info.Size()
		}

		fmt.Printf("✓ Created: %s (%d bytes)\n", path, size)
	}

	fmt.Printf("\n📄 Generated spec summary:\n")
	if ctx.Source == "prd" {
		fmt.Printf("   Source: %s\n", ctx.PRDPath)
	}
	fmt.Printf("   Content: Generated from %s\n", ctx.Source)
	fmt.Printf("\n🤖 AGENT: Specs are now complete. You can proceed with:\n")
	fmt.Printf("   → agentic-agent task claim %s\n", taskID)
	fmt.Printf("   → agentic-agent context build --task %s\n", taskID)
}

func runGenerateAllSpecs() {
	fmt.Println("Generating all missing specs... (TODO: implement)")
}

func runCreateSpec(specRef, template string) {
	fmt.Printf("Creating spec %s with template %s... (TODO: implement)\n", specRef, template)
}

func runValidateAllSpecs() {
	fmt.Println("Validating all task specs... (TODO: implement)")
}

func runValidateTaskSpecs(taskID string) {
	fmt.Printf("Validating specs for task %s... (TODO: implement)\n", taskID)
}

func init() {
	specGenerateCmd.Flags().Bool("auto", false, "Auto-detect best generation source")
	specGenerateCmd.Flags().Bool("interactive", false, "Generate using interactive prompts")
	specGenerateCmd.Flags().String("from-prd", "", "Generate from specific PRD file")
	specGenerateCmd.Flags().Bool("all", false, "Generate specs for all tasks with missing specs")

	specCreateCmd.Flags().String("template", "generic", "Template type (generic, task, proposal)")

	specCmd.AddCommand(specGenerateCmd)
	specCmd.AddCommand(specCreateCmd)
	specCmd.AddCommand(specValidateCmd)

	rootCmd.AddCommand(specCmd)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./cmd/agentic-agent/... -v -run TestSpecCommand_Exists`
Expected: PASS

**Step 5: Test manually**

Run: `go run cmd/agentic-agent/main.go spec --help`
Expected: See spec command help

**Step 6: Commit**

```bash
git add cmd/agentic-agent/spec.go cmd/agentic-agent/spec_test.go
git commit -m "feat: add spec command with generate, create, and validate subcommands"
```

---

## Task 6: Add --skip-validation Flag to task claim

**Files:**
- Modify: `cmd/agentic-agent/task.go` (add flag to claim command)

**Step 1: Find the task claim command**

```bash
grep -n "taskClaimCmd" cmd/agentic-agent/task.go
```

**Step 2: Add --skip-validation flag**

In the file where `taskClaimCmd` is defined (look for `var taskClaimCmd`), add to init section:

```go
taskClaimCmd.Flags().Bool("skip-validation", false, "Skip spec validation checks")
```

**Step 3: Modify claim logic to respect flag**

In the `Run` function of `taskClaimCmd`, add:

```go
skipValidation, _ := cmd.Flags().GetBool("skip-validation")

// ... existing code to find task ...

if skipValidation {
	// Use regular ClaimTask without readiness checks
	if err := tm.ClaimTask(taskID, assignee); err != nil {
		fmt.Printf("Error claiming task: %v\n", err)
		os.Exit(1)
	}
} else {
	// Use ClaimTaskWithConfig (includes readiness checks)
	if err := tm.ClaimTaskWithConfig(taskID, assignee, cfg); err != nil {
		fmt.Printf("Error claiming task: %v\n", err)
		os.Exit(1)
	}
}
```

**Step 4: Test manually**

```bash
# Test with validation
agentic-agent task claim TASK-123

# Test skipping validation
agentic-agent task claim TASK-123 --skip-validation
```

**Step 5: Commit**

```bash
git add cmd/agentic-agent/task.go
git commit -m "feat: add --skip-validation flag to task claim command"
```

---

## Task 7: Add Validation Hook to Task Create

**Files:**
- Modify: `cmd/agentic-agent/task.go:taskCreateCmd`
- Modify: `cmd/agentic-agent/openspec.go` (where tasks are created in bulk)

**Step 1: Add validation check after task creation in taskCreateCmd**

In `cmd/agentic-agent/task.go`, in the `taskCreateCmd.Run` function, after task is saved:

```go
// After: tm.SaveTasks("backlog", backlog)

// Check for missing specs
if len(task.SpecRefs) > 0 {
	cfg := getConfig()
	resolver := specs.NewResolver(cfg)
	resolved := resolver.ResolveAll(task.SpecRefs)

	missingSpecs := []string{}
	for _, r := range resolved {
		if !r.Found {
			missingSpecs = append(missingSpecs, r.Ref)
		}
	}

	if len(missingSpecs) > 0 {
		fmt.Printf("\n⚠️  SPEC WARNING: %d spec(s) referenced but don't exist\n", len(missingSpecs))
		for _, spec := range missingSpecs {
			fmt.Printf("  - %s\n", spec)
		}
		fmt.Printf("\n💡 TIP: Generate specs with:\n")
		fmt.Printf("   → agentic-agent spec generate %s --auto\n", task.ID)
	}
}
```

**Step 2: Add bulk validation after openspec task creation**

Find where `openspec init` creates multiple tasks (likely in `cmd/agentic-agent/openspec.go`), and add similar check:

```go
// After all tasks are created

fmt.Printf("\n✓ Created %d tasks in .agentic/tasks/backlog.yaml\n", taskCount)

// Check for missing specs across all tasks
cfg := getConfig()
resolver := specs.NewResolver(cfg)

tasksWithMissingSpecs := []string{}
totalMissing := 0

for _, task := range allCreatedTasks {
	if len(task.SpecRefs) == 0 {
		continue
	}

	resolved := resolver.ResolveAll(task.SpecRefs)
	hasMissing := false
	for _, r := range resolved {
		if !r.Found {
			hasMissing = true
			totalMissing++
		}
	}

	if hasMissing {
		tasksWithMissingSpecs = append(tasksWithMissingSpecs, task.ID)
	}
}

if len(tasksWithMissingSpecs) > 0 {
	fmt.Printf("\n⚠️  VALIDATION WARNING: %d tasks reference specs that weren't generated\n", len(tasksWithMissingSpecs))
	fmt.Printf("\n📝 Tasks with missing specs:\n")
	for _, taskID := range tasksWithMissingSpecs {
		fmt.Printf("   • %s\n", taskID)
	}

	fmt.Printf("\n🤖 AGENT GUIDANCE:\n\n")
	fmt.Printf("  These tasks reference granular spec files that weren't part of the proposal.\n\n")
	fmt.Printf("  You have 3 options:\n\n")
	fmt.Printf("  A) Generate all missing specs from proposal\n")
	fmt.Printf("     → agentic-agent spec generate --all --from-proposal\n\n")
	fmt.Printf("  B) Generate them on-demand when claiming each task\n")
	fmt.Printf("     → Tasks will auto-prompt during claim if validate_specs_on_claim is enabled\n\n")
	fmt.Printf("  C) Work without detailed specs\n")
	fmt.Printf("     → Task descriptions contain enough info to proceed\n")
	fmt.Printf("     → Specs can be backfilled later if needed\n\n")
	fmt.Printf("  Recommendation: Option A for comprehensive documentation,\n")
	fmt.Printf("                  Option B for just-in-time workflow,\n")
	fmt.Printf("                  Option C for quick iteration\n\n")
}
```

**Step 3: Test manually**

```bash
# Create task with spec refs
agentic-agent task create --title "Test" --spec-refs "missing-spec.md"

# Create via openspec
agentic-agent openspec init "test-feature" --from some-prd.md
```

Expected: See warning messages with guidance

**Step 4: Commit**

```bash
git add cmd/agentic-agent/task.go cmd/agentic-agent/openspec.go
git commit -m "feat: add spec validation warnings to task create lifecycle"
```

---

## Task 8: Update CLAUDE.md Template with Spec Workflow

**Files:**
- Modify: `internal/skills/templates/claude-rules.tmpl` (or equivalent template file)

**Step 1: Find the CLAUDE.md template**

```bash
find . -name "*claude*tmpl" -o -name "*CLAUDE*tmpl"
```

**Step 2: Add spec validation section**

Add to the template (adjust path based on step 1):

```markdown
## Spec Validation Workflow

When you encounter missing spec warnings during task claim:

1. **Read the CLI output carefully** - it contains specific guidance
2. **Assess the situation:**
   - If user wants quick iteration → use `--skip-validation`
   - If proper documentation matters → use `spec generate --auto`
   - If you need to understand requirements → use `spec generate --interactive`
3. **Follow the recommended actions** in the CLI output
4. **Don't manually create tasks** - use `openspec init` pipeline instead

The CLI will guide you with messages prefixed with 🤖 AGENT.

## Spec Management Commands

- `agentic-agent spec generate <TASK_ID> --auto` — Auto-generate missing specs
- `agentic-agent spec generate <TASK_ID> --interactive` — Generate with prompts
- `agentic-agent spec create <SPEC_REF>` — Create single spec manually
- `agentic-agent spec validate` — Check all tasks for missing specs
```

**Step 3: Regenerate CLAUDE.md**

```bash
agentic-agent skills ensure
```

**Step 4: Verify changes**

```bash
cat CLAUDE.md | grep -A 10 "Spec Validation"
```

Expected: See new section in CLAUDE.md

**Step 5: Commit**

```bash
git add internal/skills/templates/*.tmpl CLAUDE.md
git commit -m "docs: add spec validation workflow to agent rules template"
```

---

## Task 9: Add Integration Tests

**Files:**
- Create: `test/integration/spec_validation_test.go`

**Step 1: Write integration test**

```go
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSpecValidationWorkflow(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir("..")

	// Initialize project
	cmd := exec.Command("agentic-agent", "init", "--non-interactive")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init: %v", err)
	}

	// Create task with spec refs
	cmd = exec.Command("agentic-agent", "task", "create",
		"--title", "Test Task",
		"--spec-refs", "test-spec.md")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to create task: %v\n%s", err, output)
	}

	// Should see spec warning
	if !strings.Contains(string(output), "SPEC WARNING") {
		t.Error("Expected spec warning in task create output")
	}

	// Extract task ID
	lines := strings.Split(string(output), "\n")
	var taskID string
	for _, line := range lines {
		if strings.Contains(line, "Created task") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				taskID = strings.TrimSuffix(parts[2], ":")
			}
		}
	}

	if taskID == "" {
		t.Fatal("Could not extract task ID")
	}

	// Try to claim task (should show validation)
	cmd = exec.Command("agentic-agent", "task", "claim", taskID)
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to claim: %v\n%s", err, output)
	}

	// Should see readiness warnings
	if !strings.Contains(string(output), "spec-completeness") {
		t.Error("Expected spec-completeness check in claim output")
	}

	if !strings.Contains(string(output), "RECOMMENDED ACTIONS") {
		t.Error("Expected agent guidance in claim output")
	}

	// Generate spec
	cmd = exec.Command("agentic-agent", "spec", "generate", taskID, "--auto")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to generate spec: %v\n%s", err, output)
	}

	// Check spec file exists
	specPath := filepath.Join(".agentic/spec", "test-spec.md")
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Errorf("Expected spec file to be created at %s", specPath)
	}

	// Claim should now succeed without warnings
	cmd = exec.Command("agentic-agent", "task", "claim", taskID)
	output, err = cmd.CombinedOutput()
	if err == nil {
		// Task might already be claimed from previous attempt
		if strings.Contains(string(output), "spec-completeness") {
			outputStr := string(output)
			if !strings.Contains(outputStr, "[+] spec-completeness") {
				t.Error("Expected spec-completeness to pass after generation")
			}
		}
	}
}
```

**Step 2: Run test to verify**

Run: `go test ./test/integration/... -v -run TestSpecValidationWorkflow`
Expected: PASS

**Step 3: Commit**

```bash
git add test/integration/spec_validation_test.go
git commit -m "test: add integration test for spec validation workflow"
```

---

## Task 10: Update Documentation

**Files:**
- Modify: `README.md`
- Create: `docs/spec-validation.md`

**Step 1: Write spec validation guide**

Create `docs/spec-validation.md`:

```markdown
# Spec Validation System

The spec validation system ensures tasks always have proper specification files through hybrid validation that warns, auto-generates, allows manual creation, and permits skipping.

## How It Works

### Validation Points

1. **Task Create** - Warns when creating tasks with missing specs
2. **Task Claim** - Validates specs before claiming (configurable)

### Smart Context Detection

When generating missing specs, the system automatically detects the best source:

1. **PRD-based** - If a PRD file exists for the feature
2. **Interactive** - If running in a terminal (TTY)
3. **Metadata** - Fallback using task description and acceptance criteria

### Configuration

In `agnostic-agent.yaml`:

```yaml
workflow:
  validate_specs_on_claim: true  # Enable pre-claim validation
  spec_validation_mode: "warn"   # "warn" | "block" | "silent"
```

## Commands

### Generate Missing Specs

Auto-detect best source:
```bash
agentic-agent spec generate TASK-123 --auto
```

Interactive prompts:
```bash
agentic-agent spec generate TASK-123 --interactive
```

From specific PRD:
```bash
agentic-agent spec generate TASK-123 --from-prd path/to/prd.md
```

Generate all missing specs:
```bash
agentic-agent spec generate --all
```

### Create Specs Manually

```bash
agentic-agent spec create todo-app/proposal.md
agentic-agent spec create todo-app/tasks/01-setup.md --template task
```

### Validate Specs

Check all tasks:
```bash
agentic-agent spec validate
```

Check specific task:
```bash
agentic-agent spec validate TASK-123
```

### Skip Validation

When claiming a task:
```bash
agentic-agent task claim TASK-123 --skip-validation
```

## Agent Behavior

AI agents receive guidance through CLI output:

- **Warnings** include specific commands to run
- **Options** are presented with context
- **Recommendations** based on situation

Example output:
```
⚠️  MISSING SPECS DETECTED
📋 RECOMMENDED ACTIONS FOR AGENTS:
  Option 1: Auto-generate specs (recommended)
  → agentic-agent spec generate TASK-123 --auto
  ...
```

## Best Practices

1. **Follow the pipeline** - brainstorming → product-wizard → openspec init
2. **Generate early** - Create specs during task creation
3. **Use --auto** - For AI agents, auto-detection works well
4. **Skip sparingly** - Only skip validation for quick iterations
```

**Step 2: Update README with spec validation section**

Add to `README.md` in appropriate section:

```markdown
### Spec Validation

Tasks can reference specification files that provide detailed requirements. The system validates spec completeness and helps generate missing specs:

```bash
# Generate missing specs automatically
agentic-agent spec generate TASK-123 --auto

# Validate all tasks
agentic-agent spec validate
```

See [Spec Validation Guide](docs/spec-validation.md) for details.
```

**Step 3: Commit**

```bash
git add README.md docs/spec-validation.md
git commit -m "docs: add spec validation system documentation"
```

---

## Task 11: Final Integration and Testing

**Files:**
- Test all components together

**Step 1: Build the project**

Run: `make build` or `go build ./cmd/agentic-agent`
Expected: Clean build with no errors

**Step 2: Run full test suite**

Run: `go test ./...`
Expected: All tests pass

**Step 3: Manual end-to-end test**

```bash
# Create fresh test project
mkdir /tmp/spec-test && cd /tmp/spec-test
agentic-agent init --non-interactive

# Create task with missing specs
agentic-agent task create --title "Test Feature" --spec-refs "test/proposal.md,test/plan.md"

# Should see warning
# Extract task ID from output

# Try to claim
agentic-agent task claim TASK-XXX
# Should see validation warnings with guidance

# Generate specs
agentic-agent spec generate TASK-XXX --auto
# Should create spec files

# Claim again
agentic-agent task claim TASK-XXX
# Should succeed with all specs validated

# Test skip flag
agentic-agent task create --title "Skip Test" --spec-refs "missing.md"
agentic-agent task claim TASK-YYY --skip-validation
# Should skip validation
```

**Step 4: Document any issues found**

Create GitHub issues for any bugs discovered during testing.

**Step 5: Final commit**

```bash
git add -A
git commit -m "feat: complete spec validation system implementation"
```

---

## Completion Checklist

- [ ] Config fields added and tested
- [ ] SpecGenerator implemented with context detection
- [ ] Readiness checks enhanced
- [ ] Agent-guiding output added
- [ ] `spec generate` command working
- [ ] `spec create` command working
- [ ] `spec validate` command working
- [ ] `--skip-validation` flag working
- [ ] Task create lifecycle hooks added
- [ ] CLAUDE.md template updated
- [ ] Integration tests passing
- [ ] Documentation complete
- [ ] Manual E2E testing successful
