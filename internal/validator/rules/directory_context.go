package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/javierbenavides/agentic-agent/internal/validator"
)

type DirectoryContextRule struct{}

func (r *DirectoryContextRule) Name() string {
	return "context-required"
}

func (r *DirectoryContextRule) Validate(ctx *validator.ValidationContext) (*validator.RuleResult, error) {
	var failures []string

	err := filepath.Walk(ctx.ProjectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Skip hidden and system dirs
			if strings.HasPrefix(info.Name(), ".") || info.Name() == "vendor" || info.Name() == "node_modules" {
				if info.Name() != "." { // Don't skip root
					return filepath.SkipDir
				}
			}

			// Check for source files
			hasSource := false
			entries, _ := os.ReadDir(path)
			for _, e := range entries {
				if !e.IsDir() {
					ext := filepath.Ext(e.Name())
					if ext == ".go" || ext == ".ts" || ext == ".js" || ext == ".py" {
						hasSource = true
						break
					}
				}
			}

			if hasSource {
				// Check for context.md
				if _, err := os.Stat(filepath.Join(path, "context.md")); os.IsNotExist(err) {
					failures = append(failures, fmt.Sprintf("Missing context.md in %s", path))
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	status := "PASS"
	if len(failures) > 0 {
		status = "FAIL"
	}

	return &validator.RuleResult{
		RuleName: r.Name(),
		Status:   status,
		Errors:   failures,
	}, nil
}
