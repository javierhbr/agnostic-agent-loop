package models

import "time"

type TrackType string

const (
	TrackTypeFeature  TrackType = "feature"
	TrackTypeBug      TrackType = "bug"
	TrackTypeRefactor TrackType = "refactor"
)

type TrackStatus string

const (
	TrackStatusIdeation TrackStatus = "ideation"
	TrackStatusPlanning TrackStatus = "planning"
	TrackStatusActive   TrackStatus = "active"
	TrackStatusBlocked  TrackStatus = "blocked"
	TrackStatusDone     TrackStatus = "done"
	TrackStatusArchived TrackStatus = "archived"
)

type Track struct {
	ID        string      `yaml:"id"`
	Name      string      `yaml:"name"`
	Type      TrackType   `yaml:"type"`
	Status    TrackStatus `yaml:"status"`
	TaskIDs   []string    `yaml:"task_ids,omitempty"`
	CreatedAt time.Time   `yaml:"created_at"`
	SpecPath       string      `yaml:"spec_path,omitempty"`
	PlanPath       string      `yaml:"plan_path,omitempty"`
	BrainstormPath string      `yaml:"brainstorm_path,omitempty"`
}
