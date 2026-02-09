package gitops

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Commit represents a git commit.
type Commit struct {
	Hash    string
	Message string
	Author  string
	Time    time.Time
	Files   []string
}

// Tracker provides read-only git operations for task tracking.
type Tracker struct{}

// NewTracker creates a new git tracker.
func NewTracker() *Tracker {
	return &Tracker{}
}

// IsGitRepo returns true if the current directory is inside a git repository.
func (t *Tracker) IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	return cmd.Run() == nil
}

// GetCurrentBranch returns the current git branch name.
func (t *Tracker) GetCurrentBranch() (string, error) {
	out, err := runGit("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// GetCommitsSince returns all commits since the given timestamp.
func (t *Tracker) GetCommitsSince(since time.Time) ([]Commit, error) {
	sinceStr := since.Format("2006-01-02T15:04:05")
	out, err := runGit("log", "--since="+sinceStr, "--format=%H|%s|%an|%aI", "--no-merges")
	if err != nil {
		return nil, err
	}

	return parseCommitLog(out)
}

// GetFilesChangedSince returns all files changed since the given timestamp.
func (t *Tracker) GetFilesChangedSince(since time.Time) ([]string, error) {
	sinceStr := since.Format("2006-01-02T15:04:05")
	out, err := runGit("log", "--since="+sinceStr, "--name-only", "--format=", "--no-merges")
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var files []string
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !seen[line] {
			seen[line] = true
			files = append(files, line)
		}
	}
	return files, nil
}

// GetCommitFiles returns the files changed in a specific commit.
func (t *Tracker) GetCommitFiles(hash string) ([]string, error) {
	out, err := runGit("diff-tree", "--no-commit-id", "--name-only", "--root", "-r", hash)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

func parseCommitLog(output string) ([]Commit, error) {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil, nil
	}

	var commits []Commit
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 4)
		if len(parts) < 4 {
			continue
		}

		commitTime, err := time.Parse(time.RFC3339, parts[3])
		if err != nil {
			commitTime = time.Time{}
		}

		commits = append(commits, Commit{
			Hash:    parts[0],
			Message: parts[1],
			Author:  parts[2],
			Time:    commitTime,
		})
	}
	return commits, nil
}

func runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), stderr.String())
	}
	return stdout.String(), nil
}
