package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/javierbenavides/agentic-agent/internal/sdd"
	"github.com/javierbenavides/agentic-agent/internal/ui/helpers"
	"github.com/javierbenavides/agentic-agent/internal/ui/styles"
	"github.com/spf13/cobra"
)

var sddCmd = &cobra.Command{
	Use:   "sdd",
	Short: "SDD v3.0 workflow management",
}

// sdd start <name> --risk low|medium|high|critical
var sddStartCmd = &cobra.Command{
	Use:   "start <name>",
	Short: "Start a new SDD initiative",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := getConfig()
		name := args[0]

		// Get risk level from flag, config, or interactive prompt
		flagRisk, _ := cmd.Flags().GetString("risk")
		var riskLevel sdd.RiskLevel
		var err error

		if flagRisk != "" {
			// Flag provided
			riskLevel, err = sdd.DetermineRisk(flagRisk, "", false)
		} else if helpers.ShouldUseInteractiveMode(cmd) {
			// Interactive: show risk selection prompt
			riskLevel, err = promptForRisk()
		} else {
			// Non-interactive: use config default
			riskLevel, err = sdd.DetermineRisk("", cfg.SDD.DefaultRisk, false)
		}

		if err != nil {
			return fmt.Errorf("failed to determine risk level: %w", err)
		}

		// Create initiative
		initiativeDir := cfg.SDD.InitiativesDir
		if err := os.MkdirAll(initiativeDir, 0755); err != nil {
			return fmt.Errorf("failed to create initiatives directory: %w", err)
		}

		manager := sdd.NewInitiativeManager(initiativeDir)
		initiative, err := manager.Create(name, riskLevel)
		if err != nil {
			return fmt.Errorf("failed to create initiative: %w", err)
		}

		// Print workflow info
		workflow := sdd.RiskToWorkflow(initiative.Risk)
		agents := sdd.WorkflowAgents(workflow)

		msg := fmt.Sprintf("Initiative '%s' created\n", name)
		msg += fmt.Sprintf("Risk Level: %s\n", riskLevel)
		msg += fmt.Sprintf("Workflow: %s\n", workflow)
		msg += fmt.Sprintf("Agent Sequence: %s\n", formatAgents(agents))
		msg += fmt.Sprintf("Current Agent: %s", agents[0])

		printSuccess(cmd, msg)
		return nil
	},
}

// sdd workflow show <id>
var sddWorkflowShowCmd = &cobra.Command{
	Use:   "workflow show <id>",
	Short: "Show workflow progress for an initiative",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := getConfig()
		id := args[0]

		manager := sdd.NewInitiativeManager(cfg.SDD.InitiativesDir)
		init, err := manager.Get(id)
		if err != nil {
			return fmt.Errorf("failed to get initiative: %w", err)
		}

		agents := sdd.WorkflowAgents(init.Workflow)
		msg := fmt.Sprintf("Initiative: %s\n", init.Name)
		msg += fmt.Sprintf("Risk Level: %s\n", init.Risk)
		msg += fmt.Sprintf("Workflow: %s\n", init.Workflow)
		msg += fmt.Sprintf("Status: %s\n", init.Status)
		msg += fmt.Sprintf("Current Agent: %s\n", init.CurrentAgent)
		msg += fmt.Sprintf("Agent Sequence: %s", formatAgents(agents))

		printSuccess(cmd, msg)
		return nil
	},
}

// sdd agents install [--dir .claude/agents] [--force]
var sddAgentsInstallCmd = &cobra.Command{
	Use:   "agents install",
	Short: "Install SDD agent Markdown files",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := getConfig()
		dir, _ := cmd.Flags().GetString("dir")
		force, _ := cmd.Flags().GetBool("force")

		if dir == "" {
			dir = cfg.SDD.AgentsDir
		}

		written, err := sdd.InstallAgents(dir, force)
		if err != nil {
			return fmt.Errorf("failed to install agents: %w", err)
		}

		msg := fmt.Sprintf("Installed %d agent files to %s", len(written), dir)
		for _, file := range written {
			msg += fmt.Sprintf("\n  - %s", filepath.Base(file))
		}

		printSuccess(cmd, msg)
		return nil
	},
}

// sdd gate-check <spec-id> [--format text|json]
var sddGateCheckCmd = &cobra.Command{
	Use:   "gate-check <spec-id>",
	Short: "Run all 5 SDD gates on a spec",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := getConfig()
		specID := args[0]
		format, _ := cmd.Flags().GetString("format")

		// Find the spec directory
		specDir := filepath.Join(cfg.Paths.OpenSpecDir, specID)
		if _, err := os.Stat(specDir); err != nil {
			return fmt.Errorf("spec directory not found: %s", specDir)
		}

		// Load spec graph to get node info
		graph := &sdd.SpecGraph{Nodes: make(map[string]sdd.SpecGraphNode)}
		if err := graph.Load(cfg.SDD.SpecGraphPath); err == nil {
			// Graph loaded, get the node (optional)
		}

		// Use basic node for now if not in graph
		node, ok := graph.Get(specID)
		if !ok {
			node = sdd.SpecGraphNode{
				ID:     specID,
				Status: sdd.SpecStatusDraft,
			}
		}

		// Run gates
		report, err := sdd.RunGates(specDir, node)
		if err != nil {
			return fmt.Errorf("failed to run gates: %w", err)
		}

		// Format output
		if format == "json" {
			data, _ := json.MarshalIndent(report, "", "  ")
			fmt.Println(string(data))
		} else {
			printGateReport(cmd, report)
		}

		if !report.Passed {
			return fmt.Errorf("gate check failed")
		}

		return nil
	},
}

// sdd sync-graph [--from .agentic/spec-graph.json] [--to graph/index.yaml]
var sddSyncGraphCmd = &cobra.Command{
	Use:   "sync-graph",
	Short: "Sync spec graph to platform repo",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := getConfig()
		fromPath, _ := cmd.Flags().GetString("from")
		toPath, _ := cmd.Flags().GetString("to")

		if fromPath == "" {
			fromPath = cfg.SDD.SpecGraphPath
		}
		if toPath == "" {
			toPath = "graph/index.yaml"
		}

		graph := &sdd.SpecGraph{Nodes: make(map[string]sdd.SpecGraphNode)}
		if err := graph.SyncToRemote(fromPath, toPath); err != nil {
			return fmt.Errorf("failed to sync graph: %w", err)
		}

		msg := fmt.Sprintf("Synced spec graph from %s to %s", fromPath, toPath)
		printSuccess(cmd, msg)
		return nil
	},
}

// sdd adr create --title "..." [--scope global|local]
var sddADRCreateCmd = &cobra.Command{
	Use:   "adr create",
	Short: "Create a new Architecture Decision Record",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := getConfig()
		title, _ := cmd.Flags().GetString("title")
		scope, _ := cmd.Flags().GetString("scope")

		if title == "" {
			return fmt.Errorf("--title is required")
		}
		if scope == "" {
			scope = "local"
		}

		if scope != "global" && scope != "local" {
			return fmt.Errorf("--scope must be 'global' or 'local'")
		}

		manager := sdd.NewADRManager(cfg.SDD.ADRDir)
		adr, err := manager.Create(title, scope)
		if err != nil {
			return fmt.Errorf("failed to create ADR: %w", err)
		}

		msg := fmt.Sprintf("Created ADR: %s\n", adr.ID)
		msg += fmt.Sprintf("Title: %s\n", adr.Title)
		msg += fmt.Sprintf("Scope: %s\n", adr.Scope)
		msg += fmt.Sprintf("File: %s", adr.FilePath)

		printSuccess(cmd, msg)
		return nil
	},
}

// sdd adr resolve <id>
var sddADRResolveCmd = &cobra.Command{
	Use:   "adr resolve <id>",
	Short: "Mark an ADR as approved and unblock dependent specs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := getConfig()
		id := args[0]

		// Load spec graph
		graph := &sdd.SpecGraph{Nodes: make(map[string]sdd.SpecGraphNode)}
		if err := graph.Load(cfg.SDD.SpecGraphPath); err != nil {
			return fmt.Errorf("failed to load spec graph: %w", err)
		}

		// Resolve ADR
		manager := sdd.NewADRManager(cfg.SDD.ADRDir)
		if err := manager.Resolve(id, graph); err != nil {
			return fmt.Errorf("failed to resolve ADR: %w", err)
		}

		// Save updated graph
		if err := graph.Save(cfg.SDD.SpecGraphPath); err != nil {
			return fmt.Errorf("failed to save spec graph: %w", err)
		}

		msg := fmt.Sprintf("Resolved ADR: %s\n", id)
		msg += "Spec graph updated - dependent specs unblocked"

		printSuccess(cmd, msg)
		return nil
	},
}

// sdd adr list [--blocked]
var sddADRListCmd = &cobra.Command{
	Use:   "adr list",
	Short: "List all ADRs",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := getConfig()
		blocked, _ := cmd.Flags().GetBool("blocked")

		manager := sdd.NewADRManager(cfg.SDD.ADRDir)
		var adrs []sdd.ADR
		var err error

		if blocked {
			adrs, err = manager.ListBlocking()
		} else {
			adrs, err = manager.List()
		}

		if err != nil {
			return fmt.Errorf("failed to list ADRs: %w", err)
		}

		if len(adrs) == 0 {
			msg := "No ADRs found"
			if blocked {
				msg = "No blocking ADRs found"
			}
			printSuccess(cmd, msg)
			return nil
		}

		msg := fmt.Sprintf("Found %d ADR(s):\n", len(adrs))
		for _, adr := range adrs {
			msg += fmt.Sprintf("\n%s: %s\n", adr.ID, adr.Title)
			msg += fmt.Sprintf("  Status: %s\n", adr.Status)
			msg += fmt.Sprintf("  Scope: %s", adr.Scope)
			if len(adr.Blocks) > 0 {
				msg += fmt.Sprintf("\n  Blocks: %v", adr.Blocks)
			}
		}

		printSuccess(cmd, msg)
		return nil
	},
}

func initSDDCmd() {
	sddStartCmd.Flags().String("risk", "", "Risk level: low|medium|high|critical")
	sddAgentsInstallCmd.Flags().String("dir", "", "Target directory (default: .claude/agents)")
	sddAgentsInstallCmd.Flags().Bool("force", false, "Overwrite existing files")
	sddGateCheckCmd.Flags().String("format", "text", "Output format: text|json")
	sddSyncGraphCmd.Flags().String("from", "", "Source graph path")
	sddSyncGraphCmd.Flags().String("to", "", "Destination graph path")
	sddADRCreateCmd.Flags().String("title", "", "ADR title (required)")
	sddADRCreateCmd.Flags().String("scope", "local", "Scope: global|local")
	sddADRListCmd.Flags().Bool("blocked", false, "List only blocking ADRs")

	sddCmd.AddCommand(sddStartCmd)
	sddCmd.AddCommand(sddWorkflowShowCmd)
	sddCmd.AddCommand(sddAgentsInstallCmd)
	sddCmd.AddCommand(sddGateCheckCmd)
	sddCmd.AddCommand(sddSyncGraphCmd)
	sddCmd.AddCommand(sddADRCreateCmd)
	sddCmd.AddCommand(sddADRResolveCmd)
	sddCmd.AddCommand(sddADRListCmd)
}

// Helper functions

func printGateReport(cmd *cobra.Command, report *sdd.GateReport) {
	msg := fmt.Sprintf("Gate Report for %s\n", report.SpecID)
	msg += "========================================\n"

	for _, gate := range report.Gates {
		status := "✓ PASS"
		if gate.Status == "FAIL" {
			status = "✗ FAIL"
		}
		msg += fmt.Sprintf("\nGate %d: %s — %s\n", gate.Gate, gate.Name, status)

		if len(gate.Issues) > 0 {
			msg += "Issues:\n"
			for _, issue := range gate.Issues {
				msg += fmt.Sprintf("  - %s\n", issue)
			}
		}

		if len(gate.Remediation) > 0 {
			msg += "Remediation:\n"
			for _, rem := range gate.Remediation {
				msg += fmt.Sprintf("  - %s\n", rem)
			}
		}
	}

	msg += "\n========================================"
	if report.Passed {
		msg += "\nAll gates PASSED ✓"
	} else {
		msg += "\nSome gates FAILED ✗"
	}

	if helpers.ShouldUseInteractiveMode(cmd) {
		fmt.Println(styles.RenderSuccess(msg))
	} else {
		fmt.Println(msg)
	}
}

func formatAgents(agents []string) string {
	result := ""
	for i, agent := range agents {
		if i > 0 {
			result += " → "
		}
		result += agent
	}
	return result
}

func promptForRisk() (sdd.RiskLevel, error) {
	// Stub for now - will be replaced with interactive prompt
	// For testing, return medium as default
	return sdd.RiskMedium, nil
}
