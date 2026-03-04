package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	"github.com/spf13/cobra"
)

type platformSpecMetadata struct {
	ID           string   `yaml:"id"`
	Name         string   `yaml:"name"`
	Initiative   string   `yaml:"initiative"`
	ContextPack  string   `yaml:"context_pack"`
	Status       string   `yaml:"status"`
	BlockedBy    []string `yaml:"blocked_by,omitempty"`
	Contracts    []string `yaml:"contracts,omitempty"`
	ChangeType   string   `yaml:"change_type,omitempty"`
	SpecFilePath string   `yaml:"spec_file"`
}

type initiativeRecord struct {
	ID       string `yaml:"id"`
	Priority string `yaml:"priority"`
	Status   string `yaml:"status,omitempty"`
}

var platformCmd = &cobra.Command{
	Use:   "platform",
	Short: "Manage platform-level SpecKit artifacts (Platform Repo tier)",
}

var platformInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Bootstrap platform repo structure",
	RunE: func(cmd *cobra.Command, args []string) error {
		dirs := []string{
			"constitution",
			"initiatives",
			"platform-specs",
			"contracts",
			"adr",
			filepath.Join(".agentic", "context", "packs"),
		}
		for _, d := range dirs {
			if err := os.MkdirAll(d, 0755); err != nil {
				return err
			}
		}
		graphPath := "spec-graph.json"
		if _, err := os.Stat(graphPath); os.IsNotExist(err) {
			if err := os.WriteFile(graphPath, []byte("{\"nodes\":[],\"edges\":[]}"), 0644); err != nil {
				return err
			}
		}
		printSuccess(cmd, "Platform repo initialized")
		return nil
	},
}

var platformAddFeatureCmd = &cobra.Command{
	Use:   "add-feature",
	Short: "Create a new Platform Spec",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		initiative, _ := cmd.Flags().GetString("initiative")
		contextPack, _ := cmd.Flags().GetString("context-pack")
		fanout, _ := cmd.Flags().GetString("fanout")

		if name == "" || initiative == "" || contextPack == "" {
			return fmt.Errorf("--name, --initiative, and --context-pack are required")
		}

		id := toKebab(name)
		specDir := filepath.Join("platform-specs", id)
		if err := os.MkdirAll(specDir, 0755); err != nil {
			return err
		}

		meta := platformSpecMetadata{
			ID:           id,
			Name:         name,
			Initiative:   initiative,
			ContextPack:  contextPack,
			Status:       "Draft",
			ChangeType:   "platform",
			SpecFilePath: filepath.Join(specDir, "spec.md"),
		}
		metaBytes, _ := yaml.Marshal(meta)
		if err := os.WriteFile(filepath.Join(specDir, "metadata.yaml"), metaBytes, 0644); err != nil {
			return err
		}

		specContent := fmt.Sprintf(`# Platform Spec: %s

## Metadata
- ID: %s
- Initiative: %s
- Context Pack: %s
- Status: Draft

## Problem

## Responsibilities

## Contracts

## ADRs

`, name, id, initiative, contextPack)
		if err := os.WriteFile(meta.SpecFilePath, []byte(specContent), 0644); err != nil {
			return err
		}

		if fanout != "" {
			components := splitCSV(fanout)
			if len(components) > 0 {
				var b strings.Builder
				b.WriteString("tasks:\n")
				for _, c := range components {
					b.WriteString(fmt.Sprintf("  - component_repo: %s\n    platform_spec_id: %s\n    context_pack_version: %s\n    contract_change: \"no\"\n    blocked_by: []\n", c, id, contextPack))
				}
				if err := os.WriteFile(filepath.Join(specDir, "fanout.yaml"), []byte(b.String()), 0644); err != nil {
					return err
				}
			}
		}

		printSuccess(cmd, fmt.Sprintf("Platform spec created: %s", id))
		return nil
	},
}

var platformChangeFeatureCmd = &cobra.Command{
	Use:   "change-feature",
	Short: "Update metadata of an existing Platform Spec",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetString("id")
		field, _ := cmd.Flags().GetString("field")
		value, _ := cmd.Flags().GetString("value")
		if id == "" || field == "" {
			return fmt.Errorf("--id and --field are required")
		}
		metaPath := filepath.Join("platform-specs", id, "metadata.yaml")
		metaBytes, err := os.ReadFile(metaPath)
		if err != nil {
			return err
		}
		var meta platformSpecMetadata
		if err := yaml.Unmarshal(metaBytes, &meta); err != nil {
			return err
		}

		switch field {
		case "status":
			meta.Status = value
		case "context_pack":
			meta.ContextPack = value
		case "blocked_by":
			meta.BlockedBy = splitCSV(value)
		default:
			return fmt.Errorf("unsupported field: %s", field)
		}

		metaBytes, _ = yaml.Marshal(meta)
		if err := os.WriteFile(metaPath, metaBytes, 0644); err != nil {
			return err
		}

		printSuccess(cmd, fmt.Sprintf("Updated %s.%s", id, field))
		return nil
	},
}

var platformChangePriorityCmd = &cobra.Command{
	Use:   "change-priority",
	Short: "Set initiative priority/status",
	RunE: func(cmd *cobra.Command, args []string) error {
		initiative, _ := cmd.Flags().GetString("initiative")
		priority, _ := cmd.Flags().GetString("priority")
		status, _ := cmd.Flags().GetString("status")
		if initiative == "" || priority == "" {
			return fmt.Errorf("--initiative and --priority are required")
		}

		rec := initiativeRecord{ID: initiative, Priority: priority, Status: status}
		data, _ := yaml.Marshal(rec)
		path := filepath.Join("initiatives", fmt.Sprintf("%s.yaml", initiative))
		if err := os.WriteFile(path, data, 0644); err != nil {
			return err
		}
		printSuccess(cmd, fmt.Sprintf("Initiative %s priority set to %s", initiative, priority))
		return nil
	},
}

func init() {
	platformAddFeatureCmd.Flags().String("name", "", "Platform feature name")
	platformAddFeatureCmd.Flags().String("initiative", "", "Initiative/Epic ID (e.g., ECO-123)")
	platformAddFeatureCmd.Flags().String("context-pack", "", "Context Pack version (cp-vX)")
	platformAddFeatureCmd.Flags().String("fanout", "", "Comma-separated component repos for fan-out tasks")

	platformChangeFeatureCmd.Flags().String("id", "", "Platform spec ID")
	platformChangeFeatureCmd.Flags().String("field", "", "Field to change (status|context_pack|blocked_by)")
	platformChangeFeatureCmd.Flags().String("value", "", "New value")

	platformChangePriorityCmd.Flags().String("initiative", "", "Initiative ID")
	platformChangePriorityCmd.Flags().String("priority", "", "Priority/phase (Planned|Discovery|Draft|Approved|Implementing|Done|Paused|Blocked)")
	platformChangePriorityCmd.Flags().String("status", "", "Optional status note")

	platformCmd.AddCommand(platformInitCmd)
	platformCmd.AddCommand(platformAddFeatureCmd)
	platformCmd.AddCommand(platformChangeFeatureCmd)
	platformCmd.AddCommand(platformChangePriorityCmd)
}

func printSuccess(cmd *cobra.Command, msg string) {
	if helpers.ShouldUseInteractiveMode(cmd) {
		fmt.Println(styles.RenderSuccess(msg))
	} else {
		fmt.Println(msg)
	}
}

func toKebab(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
