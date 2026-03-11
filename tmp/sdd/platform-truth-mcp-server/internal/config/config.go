package config

import (
	"fmt"
	"os"
	"path/filepath"

	"platform-truth-mcp-server/internal/yamlish"
)

type Config struct {
	Path               string
	ConfigDir          string
	PlatformRepoAbs    string
	Server             ServerConfig
	PlatformRepo       PlatformRepoConfig
	ComponentAlignment ComponentAlignmentConfig
	Inputs             InputsConfig
	Cache              CacheConfig
	Safety             SafetyConfig
}

type ServerConfig struct {
	Name         string
	Mode         string
	LanguageHint string
}

type PlatformRepoConfig struct {
	Path             string
	DefaultRefMode   string
	AllowLatestQuery bool
}

type ComponentAlignmentConfig struct {
	PlatformRefFile      string
	JiraTraceabilityFile string
}

type InputsConfig struct {
	SpecsDir        string
	ContractsDir    string
	ADRsDir         string
	RefsIndex       string
	JiraMappingFile string
}

type CacheConfig struct {
	Enabled   bool
	Directory string
}

type SafetyConfig struct {
	ReadOnly             bool
	AllowRepoWrites      bool
	AllowComponentWrites bool
}

func Load(path string) (*Config, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve config path: %w", err)
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	parsed, err := yamlish.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	cfg := &Config{
		Path:      absPath,
		ConfigDir: filepath.Dir(absPath),
		Server: ServerConfig{
			Name:         defaultString(yamlish.LookupString(parsed, "server", "name"), "platform-truth-mcp"),
			Mode:         defaultString(yamlish.LookupString(parsed, "server", "mode"), "local-read-only"),
			LanguageHint: defaultString(yamlish.LookupString(parsed, "server", "language_hint"), "go"),
		},
		PlatformRepo: PlatformRepoConfig{
			Path:             yamlish.LookupString(parsed, "platform_repo", "path"),
			DefaultRefMode:   defaultString(yamlish.LookupString(parsed, "platform_repo", "default_ref_mode"), "pinned"),
			AllowLatestQuery: yamlish.LookupBool(parsed, "platform_repo", "allow_latest_queries"),
		},
		ComponentAlignment: ComponentAlignmentConfig{
			PlatformRefFile:      defaultString(yamlish.LookupString(parsed, "component_alignment", "platform_ref_file"), "platform-ref.yaml"),
			JiraTraceabilityFile: defaultString(yamlish.LookupString(parsed, "component_alignment", "jira_traceability_file"), "jira-traceability.yaml"),
		},
		Inputs: InputsConfig{
			SpecsDir:        defaultString(yamlish.LookupString(parsed, "inputs", "specs_dir"), "."),
			ContractsDir:    defaultString(yamlish.LookupString(parsed, "inputs", "contracts_dir"), "."),
			ADRsDir:         defaultString(yamlish.LookupString(parsed, "inputs", "adrs_dir"), "."),
			RefsIndex:       yamlish.LookupString(parsed, "inputs", "refs_index"),
			JiraMappingFile: yamlish.LookupString(parsed, "inputs", "jira_mapping_file"),
		},
		Cache: CacheConfig{
			Enabled:   yamlish.LookupBool(parsed, "cache", "enabled"),
			Directory: defaultString(yamlish.LookupString(parsed, "cache", "directory"), ".cache/platform-mcp"),
		},
		Safety: SafetyConfig{
			ReadOnly:             boolWithDefault(parsed, true, "safety", "read_only"),
			AllowRepoWrites:      yamlish.LookupBool(parsed, "safety", "allow_repo_writes"),
			AllowComponentWrites: yamlish.LookupBool(parsed, "safety", "allow_component_writes"),
		},
	}

	if cfg.PlatformRepo.Path == "" {
		return nil, fmt.Errorf("platform_repo.path is required")
	}
	cfg.PlatformRepoAbs = resolveRelative(cfg.ConfigDir, cfg.PlatformRepo.Path)
	cfg.Inputs.SpecsDir = resolveRelative(cfg.PlatformRepoAbs, cfg.Inputs.SpecsDir)
	cfg.Inputs.ContractsDir = resolveRelative(cfg.PlatformRepoAbs, cfg.Inputs.ContractsDir)
	cfg.Inputs.ADRsDir = resolveRelative(cfg.PlatformRepoAbs, cfg.Inputs.ADRsDir)
	cfg.Inputs.RefsIndex = resolveRelative(cfg.PlatformRepoAbs, cfg.Inputs.RefsIndex)
	cfg.Inputs.JiraMappingFile = resolveRelative(cfg.PlatformRepoAbs, cfg.Inputs.JiraMappingFile)
	cfg.Cache.Directory = resolveRelative(cfg.ConfigDir, cfg.Cache.Directory)
	return cfg, nil
}

func defaultString(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func boolWithDefault(value any, fallback bool, path ...string) bool {
	lookedUp := yamlish.Lookup(value, path...)
	if lookedUp == nil {
		return fallback
	}
	boolean, ok := lookedUp.(bool)
	if !ok {
		return fallback
	}
	return boolean
}

func resolveRelative(base, p string) string {
	if p == "" {
		return ""
	}
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Clean(filepath.Join(base, p))
}
