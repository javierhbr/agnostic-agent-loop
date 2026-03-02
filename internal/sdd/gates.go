package sdd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// RunGates runs all five SDD gates on a spec and returns a comprehensive report.
func RunGates(specDir string, node SpecGraphNode) (*GateReport, error) {
	report := &GateReport{
		SpecID: node.ID,
		Gates:  make([]GateResult, 5),
		Passed: true,
	}

	// Run each gate in sequence
	report.Gates[0] = runGate1ContextCompleteness(specDir, node)
	report.Gates[1] = runGate2DomainValidity(specDir, node)
	report.Gates[2] = runGate3IntegrationSafety(specDir, node)
	report.Gates[3] = runGate4NFRCompliance(specDir, node)
	report.Gates[4] = runGate5ReadyToImplement(specDir, node)

	// Check if all gates passed
	for _, gate := range report.Gates {
		if gate.Status == "FAIL" {
			report.Passed = false
			break
		}
	}

	return report, nil
}

// runGate1ContextCompleteness checks that all required metadata is present.
func runGate1ContextCompleteness(specDir string, node SpecGraphNode) GateResult {
	result := GateResult{
		Gate:   1,
		Name:   "Context Completeness",
		Status: "PASS",
		Issues: []string{},
		Remediation: []string{},
	}

	metadataPath := filepath.Join(specDir, "metadata.yaml")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		result.Status = "FAIL"
		result.Issues = append(result.Issues, "metadata.yaml not found")
		result.Remediation = append(result.Remediation, "Create metadata.yaml with required fields")
		return result
	}

	var metadata map[string]interface{}
	if err := yaml.Unmarshal(data, &metadata); err != nil {
		result.Status = "FAIL"
		result.Issues = append(result.Issues, "metadata.yaml is invalid YAML")
		result.Remediation = append(result.Remediation, "Fix YAML syntax in metadata.yaml")
		return result
	}

	// Check required fields
	requiredFields := []string{"implements", "context_pack", "blocked_by", "status"}
	for _, field := range requiredFields {
		if _, ok := metadata[field]; !ok {
			result.Status = "FAIL"
			result.Issues = append(result.Issues, fmt.Sprintf("metadata.%s is missing", field))
		}
	}

	if result.Status == "FAIL" {
		result.Remediation = append(result.Remediation,
			"Add all required metadata fields: implements, context_pack, blocked_by, status")
		result.Remediation = append(result.Remediation,
			"See platform-spec/component-spec SKILL.md for template")
	}

	// Check constitution exists at platform repo root (best effort)
	if _, err := os.Stat("constitution/policies.md"); err == nil {
		// constitution exists, good
	} else {
		// Not a failure, just a warning condition, but we don't issue it here
	}

	return result
}

// runGate2DomainValidity checks for cross-domain violations and invariant references.
func runGate2DomainValidity(specDir string, node SpecGraphNode) GateResult {
	result := GateResult{
		Gate:   2,
		Name:   "Domain Validity",
		Status: "PASS",
		Issues: []string{},
		Remediation: []string{},
	}

	// Read spec.md or proposal.md
	specPath := filepath.Join(specDir, "spec.md")
	if _, err := os.Stat(specPath); err != nil {
		specPath = filepath.Join(specDir, "proposal.md")
	}

	data, err := os.ReadFile(specPath)
	if err != nil {
		result.Status = "WARN" // Not a FAIL, just a warning
		result.Issues = append(result.Issues, "Could not find spec.md or proposal.md to validate domain rules")
		result.Remediation = append(result.Remediation, "Ensure spec.md or proposal.md exists in the change directory")
		return result
	}

	specContent := string(data)

	// Check for cross-domain DB access patterns (anti-pattern: "SELECT FROM other_service_db")
	if strings.Contains(specContent, "SELECT FROM") && strings.Contains(specContent, "other") {
		result.Status = "FAIL"
		result.Issues = append(result.Issues, "Spec contains cross-domain DB access pattern")
		result.Remediation = append(result.Remediation,
			"Use domain MCP contracts (NATS events, HTTP APIs) instead of direct DB access")
	}

	// Check for invariant section
	if !strings.Contains(specContent, "Invariant") && !strings.Contains(specContent, "invariant") {
		// Not a failure, but good to have
	}

	return result
}

// runGate3IntegrationSafety checks for consumer identification and compatibility plans.
func runGate3IntegrationSafety(specDir string, node SpecGraphNode) GateResult {
	result := GateResult{
		Gate:   3,
		Name:   "Integration Safety",
		Status: "PASS",
		Issues: []string{},
		Remediation: []string{},
	}

	// If contracts_referenced is empty and no contract change declared, no check needed
	if len(node.ContractsReferenced) == 0 {
		return result
	}

	// If contracts are referenced, check that consumers are identified
	specPath := filepath.Join(specDir, "spec.md")
	if _, err := os.Stat(specPath); err != nil {
		specPath = filepath.Join(specDir, "proposal.md")
	}

	data, err := os.ReadFile(specPath)
	if err != nil {
		result.Status = "WARN"
		result.Issues = append(result.Issues, "Could not read spec to verify consumer identification")
		return result
	}

	specContent := string(data)

	// Check for "Consumers:" or "consumer" section
	if !strings.Contains(specContent, "Consumer") && !strings.Contains(specContent, "consumer") {
		result.Status = "FAIL"
		result.Issues = append(result.Issues,
			"Contracts are referenced but no consumer identification found")
		result.Remediation = append(result.Remediation,
			"Add a 'Consumers' section listing all services that depend on the changed contracts")
	}

	// Check for compatibility plan if breaking change
	if strings.Contains(specContent, "break") || strings.Contains(specContent, "Breaking") {
		if !strings.Contains(specContent, "compatibility") && !strings.Contains(specContent, "Compatibility") {
			result.Status = "FAIL"
			result.Issues = append(result.Issues, "Breaking change declared but no compatibility plan")
			result.Remediation = append(result.Remediation,
				"Add a 'Compatibility Plan' section with dual-publish strategy or versioning approach")
		}
	}

	return result
}

// runGate4NFRCompliance checks for non-functional requirements declarations.
func runGate4NFRCompliance(specDir string, node SpecGraphNode) GateResult {
	result := GateResult{
		Gate:   4,
		Name:   "NFR Compliance",
		Status: "PASS",
		Issues: []string{},
		Remediation: []string{},
	}

	specPath := filepath.Join(specDir, "spec.md")
	if _, err := os.Stat(specPath); err != nil {
		specPath = filepath.Join(specDir, "proposal.md")
	}

	data, err := os.ReadFile(specPath)
	if err != nil {
		result.Status = "WARN"
		result.Issues = append(result.Issues, "Could not read spec to verify NFR compliance")
		return result
	}

	specContent := string(data)

	// Check for NFR sections
	nfrChecks := map[string]string{
		"Logging": "log",
		"Metrics": "metric",
		"Tracing": "trace",
		"PII": "pii",
		"Performance": "performance",
	}

	missingNFRs := []string{}
	for nfrName, searchTerm := range nfrChecks {
		if !strings.Contains(strings.ToLower(specContent), searchTerm) {
			missingNFRs = append(missingNFRs, nfrName)
		}
	}

	if len(missingNFRs) > 0 {
		result.Status = "FAIL"
		result.Issues = append(result.Issues,
			fmt.Sprintf("Missing NFR sections: %s", strings.Join(missingNFRs, ", ")))
		result.Remediation = append(result.Remediation,
			"Add NFR sections for: Logging, Metrics, Tracing, PII Handling, Performance Targets")
	}

	return result
}

// runGate5ReadyToImplement checks for no open blockers and testable acceptance criteria.
func runGate5ReadyToImplement(specDir string, node SpecGraphNode) GateResult {
	result := GateResult{
		Gate:   5,
		Name:   "Ready to Implement",
		Status: "PASS",
		Issues: []string{},
		Remediation: []string{},
	}

	// Check if blocked_by is non-empty
	if len(node.BlockedBy) > 0 {
		result.Status = "FAIL"
		result.Issues = append(result.Issues,
			fmt.Sprintf("Spec is blocked by ADR(s): %s", strings.Join(node.BlockedBy, ", ")))
		result.Remediation = append(result.Remediation,
			"Resolve blocking ADR(s) before proceeding with implementation")
		result.Remediation = append(result.Remediation,
			"Use: agentic-agent sdd adr resolve <ADR-ID>")
		return result
	}

	// Read spec and check for acceptance criteria
	specPath := filepath.Join(specDir, "spec.md")
	if _, err := os.Stat(specPath); err != nil {
		specPath = filepath.Join(specDir, "proposal.md")
	}

	data, err := os.ReadFile(specPath)
	if err != nil {
		result.Status = "WARN"
		result.Issues = append(result.Issues, "Could not read spec to verify acceptance criteria")
		return result
	}

	specContent := string(data)

	// Check for AC section and GWT format
	if !strings.Contains(specContent, "Acceptance Criteria") &&
		!strings.Contains(specContent, "acceptance criteria") {
		result.Status = "FAIL"
		result.Issues = append(result.Issues, "No Acceptance Criteria section found")
		result.Remediation = append(result.Remediation,
			"Add 'Acceptance Criteria' section with minimum 3 ACs in Given/When/Then format")
	}

	// Check for GWT pattern (very basic heuristic)
	gwtCount := strings.Count(specContent, "Given ") +
		strings.Count(specContent, "When ") +
		strings.Count(specContent, "Then ")

	if gwtCount < 3 {
		result.Status = "FAIL"
		result.Issues = append(result.Issues,
			"Acceptance criteria not in Given/When/Then format or less than 3 ACs")
		result.Remediation = append(result.Remediation,
			"Write each AC in GWT format. Example: Given [context], When [action], Then [outcome]")
	}

	return result
}
