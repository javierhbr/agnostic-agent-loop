package validator

import "github.com/javierbenavides/agentic-agent/pkg/models"

type ValidationContext struct {
	ProjectRoot string
	Config      *models.Config
}

type RuleResult struct {
	RuleName string   `json:"rule_name"`
	Status   string   `json:"status"` // "PASS", "FAIL", "WARN"
	Errors   []string `json:"errors,omitempty"`
}

type ValidationRule interface {
	Name() string
	Validate(ctx *ValidationContext) (*RuleResult, error)
}

type Validator struct {
	rules []ValidationRule
}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Register(rule ValidationRule) {
	v.rules = append(v.rules, rule)
}

func (v *Validator) Validate(ctx *ValidationContext) ([]*RuleResult, error) {
	var results []*RuleResult
	for _, rule := range v.rules {
		res, err := rule.Validate(ctx)
		if err != nil {
			return nil, err
		}
		results = append(results, res)
	}
	return results, nil
}
