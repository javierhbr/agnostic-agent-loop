package config

import (
	"fmt"
	"os"

	"github.com/javierbenavides/agentic-agent/pkg/models"
	"gopkg.in/yaml.v3"
)

// LoadConfig reads the configuration from the specified path.
func LoadConfig(path string) (*models.Config, error) {
	if path == "" {
		return nil, fmt.Errorf("config path cannot be empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg models.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	SetDefaults(&cfg)

	return &cfg, nil
}

// SetDefaults applies default values to the configuration.
func SetDefaults(cfg *models.Config) {
	// Set default paths if not configured
	if cfg.Paths.PRDOutputPath == "" {
		cfg.Paths.PRDOutputPath = ".agentic/tasks/"
	}
	if cfg.Paths.ProgressTextPath == "" {
		cfg.Paths.ProgressTextPath = ".agentic/progress.txt"
	}
	if cfg.Paths.ProgressYAMLPath == "" {
		cfg.Paths.ProgressYAMLPath = ".agentic/progress.yaml"
	}
	if cfg.Paths.ArchiveDir == "" {
		cfg.Paths.ArchiveDir = ".agentic/archive/"
	}
	if len(cfg.Paths.SpecDirs) == 0 {
		cfg.Paths.SpecDirs = []string{".agentic/spec"}
	}
	if len(cfg.Paths.ContextDirs) == 0 {
		cfg.Paths.ContextDirs = []string{".agentic/context"}
	}
}
