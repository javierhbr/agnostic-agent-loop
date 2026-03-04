package sdd

import (
	"fmt"
	"time"
)

// RiskLevel represents the risk classification of a change.
type RiskLevel string

const (
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

// WorkflowType represents the development workflow type determined by risk level.
type WorkflowType string

const (
	WorkflowQuick    WorkflowType = "quick"
	WorkflowStandard WorkflowType = "standard"
	WorkflowFull     WorkflowType = "full"
)

// SpecStatus represents the state of a spec in its lifecycle.
type SpecStatus string

const (
	SpecStatusPlanned      SpecStatus = "Planned"
	SpecStatusDraft        SpecStatus = "Draft"
	SpecStatusApproved     SpecStatus = "Approved"
	SpecStatusImplementing SpecStatus = "Implementing"
	SpecStatusDone         SpecStatus = "Done"
	SpecStatusPaused       SpecStatus = "Paused"
	SpecStatusBlocked      SpecStatus = "Blocked"
)

// ADRStatus represents the state of an Architecture Decision Record.
type ADRStatus string

const (
	ADRStatusProposed  ADRStatus = "Proposed"
	ADRStatusInReview  ADRStatus = "InReview"
	ADRStatusApproved  ADRStatus = "Approved"
	ADRStatusRejected  ADRStatus = "Rejected"
)

// GateResult represents the outcome of a single gate check.
type GateResult struct {
	Gate         int      `json:"gate" yaml:"gate"`
	Name         string   `json:"name" yaml:"name"`
	Status       string   `json:"status" yaml:"status"` // "PASS" | "FAIL"
	Issues       []string `json:"issues,omitempty" yaml:"issues,omitempty"`
	Remediation  []string `json:"remediation,omitempty" yaml:"remediation,omitempty"`
}

// GateReport contains the results of running all five gates on a spec.
type GateReport struct {
	SpecID string       `json:"spec_id" yaml:"spec_id"`
	Gates  []GateResult `json:"gates" yaml:"gates"`
	Passed bool         `json:"passed" yaml:"passed"`
}

// SpecGraphNode represents a single artifact in the Spec Graph.
type SpecGraphNode struct {
	ID                  string        `json:"id" yaml:"id"`
	Implements          string        `json:"implements,omitempty" yaml:"implements,omitempty"`
	DependsOn           []string      `json:"depends_on,omitempty" yaml:"depends_on,omitempty"`
	Affects             []string      `json:"affects,omitempty" yaml:"affects,omitempty"`
	ContextPack         string        `json:"context_pack,omitempty" yaml:"context_pack,omitempty"`
	BlockedBy           []string      `json:"blocked_by,omitempty" yaml:"blocked_by,omitempty"`
	Status              SpecStatus    `json:"status" yaml:"status"`
	ContractsReferenced []string      `json:"contracts_referenced,omitempty" yaml:"contracts_referenced,omitempty"`
	UpdatedAt           time.Time     `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
}

// ADR represents an Architecture Decision Record.
type ADR struct {
	ID        string    `yaml:"id"`
	Title     string    `yaml:"title"`
	Status    ADRStatus `yaml:"status"`
	Scope     string    `yaml:"scope"` // "global" | "local"
	Owner     string    `yaml:"owner,omitempty"`
	Blocks    []string  `yaml:"blocks,omitempty"` // spec IDs blocked by this ADR
	CreatedAt time.Time `yaml:"created_at"`
	ResolvedAt time.Time `yaml:"resolved_at,omitempty"`
	FilePath  string    `yaml:"file_path"`
}

// Initiative represents an SDD initiative with its workflow state.
type Initiative struct {
	ID        string       `yaml:"id"`
	Name      string       `yaml:"name"`
	Risk      RiskLevel    `yaml:"risk"`
	Workflow  WorkflowType `yaml:"workflow"`
	CurrentAgent string     `yaml:"current_agent"`
	Status    string       `yaml:"status"`
	SpecIDs   []string     `yaml:"spec_ids,omitempty"`
	CreatedAt time.Time    `yaml:"created_at"`
}

// ValidateTransition checks if a spec status transition is valid.
func ValidateTransition(from, to SpecStatus) error {
	validTransitions := map[SpecStatus][]SpecStatus{
		SpecStatusPlanned: {SpecStatusDraft},
		SpecStatusDraft: {SpecStatusApproved, SpecStatusPaused},
		SpecStatusApproved: {SpecStatusImplementing, SpecStatusPaused, SpecStatusBlocked},
		SpecStatusImplementing: {SpecStatusDone, SpecStatusPaused, SpecStatusBlocked},
		SpecStatusDone: {},
		SpecStatusPaused: {SpecStatusDraft, SpecStatusApproved},
		SpecStatusBlocked: {SpecStatusDraft, SpecStatusApproved},
	}

	allowed, ok := validTransitions[from]
	if !ok {
		return fmt.Errorf("invalid source status: %s", from)
	}

	for _, valid := range allowed {
		if valid == to {
			return nil
		}
	}

	return fmt.Errorf("invalid transition from %s to %s", from, to)
}

// RiskToWorkflow maps a risk level to its corresponding workflow type.
func RiskToWorkflow(r RiskLevel) WorkflowType {
	switch r {
	case RiskLow:
		return WorkflowQuick
	case RiskMedium:
		return WorkflowStandard
	case RiskHigh, RiskCritical:
		return WorkflowFull
	default:
		return WorkflowStandard
	}
}

// WorkflowAgents returns the ordered list of agents for a given workflow type.
func WorkflowAgents(wf WorkflowType) []string {
	switch wf {
	case WorkflowQuick:
		return []string{"developer", "verifier"}
	case WorkflowStandard:
		return []string{"architect", "developer", "verifier"}
	case WorkflowFull:
		return []string{"analyst", "architect", "developer", "verifier"}
	default:
		return []string{"developer", "verifier"}
	}
}
