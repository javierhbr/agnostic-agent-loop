package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/javierbenavides/agentic-agent/internal/validator"
)

type ContextUpdateRule struct{}

func (r *ContextUpdateRule) Name() string {
	return "context-freshness"
}

func (r *ContextUpdateRule) Validate(ctx *validator.ValidationContext) (*validator.RuleResult, error) {
	var failures []string

	err := filepath.Walk(ctx.ProjectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Skip logic same as above
			if strings.HasPrefix(info.Name(), ".") || info.Name() == "vendor" {
				if info.Name() != "." {
					return filepath.SkipDir
				}
			}

			contextPath := filepath.Join(path, "context.md")
			contextInfo, err := os.Stat(contextPath)
			if err == nil {
				// Check source files vs context mod time
				entries, _ := os.ReadDir(path)
				for _, e := range entries {
					if !e.IsDir() && e.Name() != "context.md" {
						ext := filepath.Ext(e.Name())
						if ext == ".go" || ext == ".ts" || ext == ".js" {
							fileInfo, _ := e.Info()
							if fileInfo.ModTime().After(contextInfo.ModTime()) {
								failures = append(failures, fmt.Sprintf("Stale context in %s (newer file: %s)", path, e.Name()))
								break // Report once per dir
							}
						}
					}
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
