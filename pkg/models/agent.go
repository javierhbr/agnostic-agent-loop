package models

type AgentExecutionResult struct {
	Output           string
	Success          bool
	CriteriaMet      []string
	CriteriaFailed   []string
	FilesModified    []string
	ErrorMessage     string
	TokensUsed       int
}

func (r *AgentExecutionResult) AllCriteriaMet() bool {
	return r.Success && len(r.CriteriaFailed) == 0
}
