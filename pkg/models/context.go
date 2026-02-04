package models

import "time"

type ContextType string

const (
	ContextTypeDirectory ContextType = "directory"
	ContextTypeGlobal    ContextType = "global"
	ContextTypeSession   ContextType = "session"
)

type DirectoryContext struct {
	Path              string    `yaml:"path" json:"path"`
	Purpose           string    `yaml:"purpose" json:"purpose"`
	Responsibilities  []string  `yaml:"responsibilities" json:"responsibilities"`
	LocalArchitecture []string  `yaml:"local_architecture" json:"local_architecture"`
	Dependencies      []string  `yaml:"dependencies" json:"dependencies"`
	MustDo            []string  `yaml:"must_do" json:"must_do"`
	CannotDo          []string  `yaml:"cannot_do" json:"cannot_do"`
	KeyFiles          []string  `yaml:"key_files" json:"key_files"`
	Updated           time.Time `yaml:"updated" json:"updated"`
}

type GlobalContext struct {
	ProjectName string    `yaml:"project_name" json:"project_name"`
	Overview    string    `yaml:"overview" json:"overview"`
	Goals       []string  `yaml:"goals" json:"goals"`
	Guidelines  []string  `yaml:"guidelines" json:"guidelines"`
	Updated     time.Time `yaml:"updated" json:"updated"`
}
