package token

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type TokenUsage struct {
	TotalTokens int            `yaml:"total_tokens"`
	AgentUsage  map[string]int `yaml:"agent_usage"`
}

type TokenManager struct {
	baseDir string
}

func NewTokenManager(baseDir string) *TokenManager {
	return &TokenManager{baseDir: baseDir}
}

func (tm *TokenManager) LoadUsage() (*TokenUsage, error) {
	path := filepath.Join(tm.baseDir, "token_usage.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &TokenUsage{AgentUsage: make(map[string]int)}, nil
		}
		return nil, err
	}

	var usage TokenUsage
	if err := yaml.Unmarshal(data, &usage); err != nil {
		return nil, err
	}
	if usage.AgentUsage == nil {
		usage.AgentUsage = make(map[string]int)
	}
	return &usage, nil
}

func (tm *TokenManager) AddUsage(agent string, tokens int) error {
	usage, err := tm.LoadUsage()
	if err != nil {
		return err
	}

	usage.TotalTokens += tokens
	usage.AgentUsage[agent] += tokens

	path := filepath.Join(tm.baseDir, "token_usage.yaml")
	data, err := yaml.Marshal(usage)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
