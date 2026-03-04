package github

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// PRInfo contains GitHub PR metadata
type PRInfo struct {
	URL       string   `json:"url"`
	Number    int      `json:"number"`
	State     string   `json:"state"` // OPEN, MERGED, CLOSED
	Title     string   `json:"title"`
	Reviewers []string `json:"reviewers"`
}

// CreatePR creates a new pull request using gh CLI.
// Invokes: gh pr create --title --body --base --head --json url,number,state,title
func CreatePR(title, body, base, head string) (*PRInfo, error) {
	cmd := exec.Command("gh", "pr", "create",
		"--title", title,
		"--body", body,
		"--base", base,
		"--head", head,
		"--json", "url,number,state,title")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gh pr create failed: %w", err)
	}

	var pr PRInfo
	if err := json.Unmarshal(output, &pr); err != nil {
		return nil, fmt.Errorf("failed to parse gh output: %w", err)
	}

	return &pr, nil
}

// GetPRInfo retrieves current PR status using gh CLI.
// Invokes: gh pr view <pr-url> --json url,number,state,title,reviewers
func GetPRInfo(prURL string) (*PRInfo, error) {
	cmd := exec.Command("gh", "pr", "view", prURL,
		"--json", "url,number,state,title,reviewers")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gh pr view failed: %w", err)
	}

	var pr PRInfo
	if err := json.Unmarshal(output, &pr); err != nil {
		return nil, fmt.Errorf("failed to parse gh output: %w", err)
	}

	return &pr, nil
}

// RequestReview requests a review from a reviewer using gh CLI.
// Invokes: gh pr review <pr-url> --request-review <reviewer>
func RequestReview(prURL, reviewer string) error {
	cmd := exec.Command("gh", "pr", "review", prURL,
		"--request-review", reviewer)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gh pr review failed: %w", err)
	}

	return nil
}

// MergePR merges a PR with the specified strategy using gh CLI.
// Invokes: gh pr merge <pr-url> --<strategy>
// strategy: "squash", "merge", or "rebase"
func MergePR(prURL, strategy string) error {
	if strategy == "" {
		strategy = "squash"
	}

	cmd := exec.Command("gh", "pr", "merge", prURL, "--"+strategy)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gh pr merge failed: %w", err)
	}

	return nil
}
