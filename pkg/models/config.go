package models

type Config struct {
	Project     ProjectConfig  `yaml:"project"`
	Agents      AgentsConfig   `yaml:"agents"`
	Workflow    WorkflowConfig `yaml:"workflow"`
	Paths       PathsConfig    `yaml:"paths"`
	ActiveAgent string         `yaml:"-"` // Runtime-only: detected agent name
}

type ProjectConfig struct {
	Name    string   `yaml:"name"`
	Version string   `yaml:"version"`
	Roots   []string `yaml:"roots"` // Source roots to scan
}

type AgentsConfig struct {
	Defaults  AgentDefaults `yaml:"defaults"`
	Overrides []AgentConfig `yaml:"overrides"`
}

type AgentDefaults struct {
	MaxTokens int    `yaml:"max_tokens"`
	Model     string `yaml:"model"`
}

type AgentConfig struct {
	Name       string   `yaml:"name"` // e.g., "claude-code", "cursor"
	MaxTokens  int      `yaml:"max_tokens,omitempty"`
	Model      string   `yaml:"model,omitempty"`
	SkillPacks []string `yaml:"skill_packs,omitempty"` // packs to auto-install
	ExtraRules []string `yaml:"extra_rules,omitempty"` // additional rule lines
	AutoSetup  bool     `yaml:"auto_setup,omitempty"`  // auto-generate on init/ensure
}

type WorkflowConfig struct {
	Validators []string `yaml:"validators"`
}

type PathsConfig struct {
	PRDOutputPath    string   `yaml:"prdOutputPath,omitempty"`
	ProgressTextPath string   `yaml:"progressTextPath,omitempty"`
	ProgressYAMLPath string   `yaml:"progressYAMLPath,omitempty"`
	ArchiveDir       string   `yaml:"archiveDir,omitempty"`
	SpecDirs         []string `yaml:"specDirs,omitempty"`
	ContextDirs      []string `yaml:"contextDirs,omitempty"`
	TrackDir         string   `yaml:"trackDir,omitempty"`
	OpenSpecDir      string   `yaml:"openSpecDir,omitempty"`
}
