package tasks

import (
	"fmt"
	"os"
	"strings"

	"github.com/javierbenavides/agentic-agent/internal/specs"
	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// ReadinessCheck represents one check performed on a task.
type ReadinessCheck struct {
	Name    string
	Passed  bool
	Message string
}

// ReadinessResult holds all readiness checks for a task.
type ReadinessResult struct {
	TaskID string
	Ready  bool
	Checks []ReadinessCheck
}

// CanClaimTask performs readiness checks on a task.
// Checks: inputs exist, spec refs resolvable, scope dirs exist (warning only).
func CanClaimTask(task *models.Task, cfg *models.Config) *ReadinessResult {
	result := &ReadinessResult{
		TaskID: task.ID,
		Ready:  true,
	}

	// Check 1: All Inputs files exist on disk
	for _, input := range task.Inputs {
		if _, err := os.Stat(input); os.IsNotExist(err) {
			result.Checks = append(result.Checks, ReadinessCheck{
				Name:    "input-exists",
				Passed:  false,
				Message: fmt.Sprintf("input file %q not found", input),
			})
			result.Ready = false
		} else {
			result.Checks = append(result.Checks, ReadinessCheck{
				Name:    "input-exists",
				Passed:  true,
				Message: fmt.Sprintf("input file %q exists", input),
			})
		}
	}

	// Check 2: All SpecRefs are resolvable
	if len(task.SpecRefs) > 0 {
		resolver := specs.NewResolver(cfg)
		resolved := resolver.ResolveAll(task.SpecRefs)
		for _, r := range resolved {
			if !r.Found {
				result.Checks = append(result.Checks, ReadinessCheck{
					Name:    "spec-resolvable",
					Passed:  false,
					Message: fmt.Sprintf("spec %q not resolvable: %s", r.Ref, r.Error),
				})
				result.Ready = false
			} else {
				result.Checks = append(result.Checks, ReadinessCheck{
					Name:    "spec-resolvable",
					Passed:  true,
					Message: fmt.Sprintf("spec %q resolved at %s", r.Ref, r.Path),
				})
			}
		}
	}

	// Check 3: Scope directories exist (warning only, not blocking)
	for _, dir := range task.Scope {
		if info, err := os.Stat(dir); os.IsNotExist(err) || (err == nil && !info.IsDir()) {
			result.Checks = append(result.Checks, ReadinessCheck{
				Name:    "scope-exists",
				Passed:  false,
				Message: fmt.Sprintf("scope directory %q not found (warning only)", dir),
			})
			// Note: scope missing is a warning, does NOT set Ready=false
		} else {
			result.Checks = append(result.Checks, ReadinessCheck{
				Name:    "scope-exists",
				Passed:  true,
				Message: fmt.Sprintf("scope directory %q exists", dir),
			})
		}
	}

	return result
}

// FormatReadinessResult returns a human-readable string of readiness checks.
func FormatReadinessResult(r *ReadinessResult) string {
	var b strings.Builder
	if r.Ready {
		b.WriteString(fmt.Sprintf("Task %s: READY\n", r.TaskID))
	} else {
		b.WriteString(fmt.Sprintf("Task %s: NOT READY\n", r.TaskID))
	}

	for _, check := range r.Checks {
		icon := "+"
		if !check.Passed {
			icon = "-"
		}
		b.WriteString(fmt.Sprintf("  [%s] %s: %s\n", icon, check.Name, check.Message))
	}

	return b.String()
}
