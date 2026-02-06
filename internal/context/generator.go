package context

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// GenerateContextWithConfig analyzes the directory using the provided config.
// If cfg is nil, behaves identically to GenerateContext.
func GenerateContextWithConfig(dir string, cfg *models.Config) (*models.DirectoryContext, error) {
	// For now, config is reserved for future enhancements (e.g. custom scan rules).
	// Delegate to the core logic.
	return generateContextCore(dir)
}

// GenerateContext analyzes the directory and produces a DirectoryContext.
// This is a simplified "Expert System" or "Skill" implementation.
func GenerateContext(dir string) (*models.DirectoryContext, error) {
	return generateContextCore(dir)
}

func generateContextCore(dir string) (*models.DirectoryContext, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	var imports []string
	var exports []string

	for _, e := range entries {
		if !e.IsDir() {
			ext := filepath.Ext(e.Name())
			if ext == ".go" || ext == ".ts" || ext == ".js" || ext == ".py" {
				files = append(files, e.Name())
				// Analyze file content
				content, _ := os.ReadFile(filepath.Join(dir, e.Name()))
				sContent := string(content)

				// Very basic regex analysis
				if ext == ".go" {
					importRe := regexp.MustCompile(`"(.*?)"`)
					matches := importRe.FindAllStringSubmatch(sContent, -1)
					for _, m := range matches {
						if len(m) > 1 {
							imports = append(imports, m[1])
						}
					}

					funcRe := regexp.MustCompile(`func ([A-Z][a-zA-Z0-9_]*)`)
					matchesFunc := funcRe.FindAllStringSubmatch(sContent, -1)
					for _, m := range matchesFunc {
						if len(m) > 1 {
							exports = append(exports, m[1])
						}
					}
				}
			}
		}
	}

	// Deduplicate
	imports = unique(imports)
	exports = unique(exports)

	purpose := fmt.Sprintf("Contains %d source files. Implements functionality related to %s.", len(files), filepath.Base(dir))

	return &models.DirectoryContext{
		Path:             dir,
		Purpose:          purpose,
		Responsibilities: []string{fmt.Sprintf("Exported symbols: %s", strings.Join(limit(exports, 5), ", "))},
		Dependencies:     limit(imports, 10),
		KeyFiles:         files,
	}, nil
}

func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func limit(slice []string, n int) []string {
	if len(slice) > n {
		return slice[:n]
	}
	return slice
}
