package models

import "time"

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
	TrackID     string     `yaml:"track_id,omitempty"`     // Associated track ID
	ClaimedAt   time.Time  `yaml:"claimed_at,omitempty"`   // When the task was claimed
	CompletedAt time.Time  `yaml:"completed_at,omitempty"` // When the task was completed
	Branch      string     `yaml:"branch,omitempty"`       // Git branch when claimed
	Commits     []string   `yaml:"commits,omitempty"`      // Associated git commit hashes
}

type SubTask struct {
	ID          string     `yaml:"id"`
	Title       string     `yaml:"title"`
	Status      TaskStatus `yaml:"status"`
	AssignedTo  string     `yaml:"assigned_to,omitempty"`
}
