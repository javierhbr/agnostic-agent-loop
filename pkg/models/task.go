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
	SubTasks    []SubTask  `yaml:"subtasks,omitempty"`
}

type SubTask struct {
	ID          string     `yaml:"id"`
	Title       string     `yaml:"title"`
	Status      TaskStatus `yaml:"status"`
	AssignedTo  string     `yaml:"assigned_to,omitempty"`
}
