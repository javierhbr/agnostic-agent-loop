package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"platform-truth-mcp-server/internal/config"
)

func TestValidateComponentAlignmentAndJiraChain(t *testing.T) {
	tmp := t.TempDir()
	platformDir := filepath.Join(tmp, "platform")
	componentDir := filepath.Join(tmp, "component")
	if err := os.MkdirAll(platformDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(componentDir, 0o755); err != nil {
		t.Fatal(err)
	}

	mustWrite(t, filepath.Join(platformDir, "platform-baseline.md"), `Platform: customer-platform
Version: 2026.03
- principles.api-versioning
- contracts.customer-profile.v2
- capabilities.customer-identity
- PLAT-123
`)
	mustWrite(t, filepath.Join(componentDir, "platform-ref.yaml"), `platform:
  id: "customer-platform"
  version: "2026.03"
component:
  name: "profile-service"
change:
  change_package_id: "chg-profile-email-validation"
  alignment_type: "shared-change"
  requires_platform_change: true
  platform_change_package_id: "plat-email-validation"
platform_refs:
  principles:
    - id: "principles.api-versioning"
      reason: "Required"
  capabilities:
    - id: "capabilities.customer-identity"
      reason: "Required"
  contracts:
    - id: "contracts.customer-profile.v2"
      reason: "Required"
`)
	mustWrite(t, filepath.Join(componentDir, "jira-traceability.yaml"), `platform_issue:
  key: "PLAT-123"
  platform_version: "2026.03"
component_epic:
  key: "PROF-456"
  linked_platform_issue: "PLAT-123"
  change_package_id: "chg-profile-email-validation"
stories:
  - key: "PROF-789"
    summary: "Task one"
`)
	cfg := &config.Config{
		PlatformRepoAbs: platformDir,
		PlatformRepo: config.PlatformRepoConfig{
			DefaultRefMode: "pinned",
		},
		ComponentAlignment: config.ComponentAlignmentConfig{
			PlatformRefFile:      "platform-ref.yaml",
			JiraTraceabilityFile: "jira-traceability.yaml",
		},
		Inputs: config.InputsConfig{
			SpecsDir:        platformDir,
			ContractsDir:    platformDir,
			ADRsDir:         platformDir,
			RefsIndex:       filepath.Join(platformDir, "platform-baseline.md"),
			JiraMappingFile: filepath.Join(platformDir, "platform-baseline.md"),
		},
		Server: config.ServerConfig{Name: "platform-truth-mcp"},
	}

	svc, err := New(cfg, "0.1.0")
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	alignment, err := svc.Call(context.Background(), "validate_component_alignment", map[string]any{"componentRepoPath": componentDir})
	if err != nil {
		t.Fatalf("validate_component_alignment error: %v", err)
	}
	structured := alignment["structuredContent"].(map[string]any)
	if valid, _ := structured["valid"].(bool); !valid {
		t.Fatalf("expected alignment to be valid, got %v", structured)
	}

	jira, err := svc.Call(context.Background(), "validate_component_jira_chain", map[string]any{"componentRepoPath": componentDir})
	if err != nil {
		t.Fatalf("validate_component_jira_chain error: %v", err)
	}
	jiraStructured := jira["structuredContent"].(map[string]any)
	if valid, _ := jiraStructured["valid"].(bool); !valid {
		t.Fatalf("expected jira chain to be valid, got %v", jiraStructured)
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
