package models

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in-progress"
	StatusDone       TaskStatus = "done"
)

type Task struct {
	ID          string     `yaml:"id"`
	Title       string     `yaml:"title"`
	Description string     `yaml:"description"`
	Status      TaskStatus `yaml:"status"`
	AssignedTo  string     `yaml:"assigned_to,omitempty"`
	Scope       []string   `yaml:"scope,omitempty"`
	SpecRefs    []string   `yaml:"spec_refs,omitempty"`    // Specification file references
	Inputs      []string   `yaml:"inputs,omitempty"`       // Required input files
	Outputs     []string   `yaml:"outputs,omitempty"`      // Expected output files
	Acceptance  []string   `yaml:"acceptance,omitempty"`   // Acceptance criteria
	SubTasks    []SubTask  `yaml:"subtasks,omitempty"`
}

type SubTask struct {
	ID          string     `yaml:"id"`
	Title       string     `yaml:"title"`
	Status      TaskStatus `yaml:"status"`
	AssignedTo  string     `yaml:"assigned_to,omitempty"`
}
