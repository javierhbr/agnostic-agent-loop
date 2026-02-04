package token

type LimitStatus string

const (
	StatusOK        LimitStatus = "OK"
	StatusSoftLimit LimitStatus = "SOFT_LIMIT"
	StatusHardLimit LimitStatus = "HARD_LIMIT"
)

// Summarizer logic (limit checking)
type LimitChecker struct {
	MaxTokens int
	SoftLimit int
}

func NewLimitChecker(max int) *LimitChecker {
	return &LimitChecker{
		MaxTokens: max,
		SoftLimit: int(float64(max) * 0.8), // 80%
	}
}

func (lc *LimitChecker) Check(currentUsage int) LimitStatus {
	if currentUsage >= lc.MaxTokens {
		return StatusHardLimit
	}
	if currentUsage >= lc.SoftLimit {
		return StatusSoftLimit
	}
	return StatusOK
}

func (tm *TokenManager) CheckAgentLimit(agent string, maxTokens int) (LimitStatus, int, error) {
	usage, err := tm.LoadUsage()
	if err != nil {
		return StatusOK, 0, err
	}

	current := usage.AgentUsage[agent]
	checker := NewLimitChecker(maxTokens)
	return checker.Check(current), current, nil
}
