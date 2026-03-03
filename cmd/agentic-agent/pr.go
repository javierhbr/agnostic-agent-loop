package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/javierbenavides/agentic-agent/internal/github"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/pkg/models"
	"github.com/spf13/cobra"
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Manage pull requests (create, review, status)",
	Long: `Manage pull requests integrated with agentic-agent tasks.

Workflow:
1. Developer builds code and commits in worktree
2. agentic-agent pr create --task <ID>     (auto-generates PR from spec)
3. agentic-agent pr review --task <ID> --pr-url <URL>  (spawns reviewer task)
4. Reviewer runs gate-check and validation
5. Developer or orchestrator merges when verdict ≥ 8/10

Examples:
  agentic-agent pr create --task TASK-123
  agentic-agent pr review --task TASK-123 --pr-url https://github.com/owner/repo/pull/42
  agentic-agent pr status --pr-url https://github.com/owner/repo/pull/42`,
}

var prCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a pull request from a completed task",
	Long: `Generate and create a pull request from task spec and commits.

The command will:
1. Load the task context (spec, commits, branch)
2. Extract PR title/body from the spec file
3. Invoke: gh pr create --title --body --base main --head <branch>
4. Record PR URL in task metadata

Example:
  agentic-agent pr create --task TASK-123`,

	RunE: func(cmd *cobra.Command, args []string) error {
		taskID, _ := cmd.Flags().GetString("task")
		if taskID == "" {
			return fmt.Errorf("--task is required")
		}

		tm := tasks.NewTaskManager(".agentic/tasks")

		// Load task from in-progress (must have been claimed and worked on)
		inProgress, err := tm.LoadTasks("in-progress")
		if err != nil {
			return fmt.Errorf("failed to load in-progress tasks: %w", err)
		}

		var task *models.Task
		for i, t := range inProgress.Tasks {
			if t.ID == taskID {
				task = &inProgress.Tasks[i]
				break
			}
		}

		if task == nil {
			return fmt.Errorf("task %s not found in in-progress", taskID)
		}

		// Load spec file
		if len(task.SpecRefs) == 0 {
			return fmt.Errorf("task has no spec references")
		}

		specPath := filepath.Join(".agentic", "specs", task.SpecRefs[0])
		specContent, err := os.ReadFile(specPath)
		if err != nil {
			return fmt.Errorf("failed to read spec: %w", err)
		}

		// Build PR config from spec
		title := github.BuildPRTitle(string(specContent))
		body := github.BuildPRBody(string(specContent), task.Commits, taskID)

		// Default: PR against main from feature branch
		base := "main"
		head := task.Branch
		if head == "" {
			head = fmt.Sprintf("feature/task-%s", taskID)
		}

		fmt.Fprintf(os.Stderr, "Creating PR: %s → %s\n", head, base)

		// Create PR using gh CLI
		pr, err := github.CreatePR(title, body, base, head)
		if err != nil {
			return fmt.Errorf("failed to create PR: %w", err)
		}

		fmt.Fprintf(os.Stderr, "✅ PR created: %s\n", pr.URL)

		// Record PR in task
		task.GithubPR = models.GithubPR{
			URL:       pr.URL,
			Number:    pr.Number,
			CreatedAt: time.Now(),
		}

		// Save updated task
		if err := tm.SaveTasks("in-progress", inProgress); err != nil {
			return fmt.Errorf("failed to update task: %w", err)
		}

		// Output PR URL for scripting
		fmt.Println(pr.URL)
		return nil
	},
}

var prReviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Create a review task for a PR",
	Long: `Spawn a reviewer task to independently review code and spec.

The command will:
1. Create a new task with type=review
2. Link to the original task and PR
3. Output instructions for reviewer

Reviewer will:
- Load context: agentic-agent context build --task <REVIEW-ID>
- Run gates: agentic-agent sdd gate-check <spec-id>
- Run validation: agentic-agent validate
- Check code quality and announce verdict

Example:
  agentic-agent pr review --task TASK-123 --pr-url https://github.com/owner/repo/pull/42`,

	RunE: func(cmd *cobra.Command, args []string) error {
		taskID, _ := cmd.Flags().GetString("task")
		prURL, _ := cmd.Flags().GetString("pr-url")

		if taskID == "" || prURL == "" {
			return fmt.Errorf("--task and --pr-url are required")
		}

		tm := tasks.NewTaskManager(".agentic/tasks")

		// Load original task
		originalTask, source, err := tm.FindTask(taskID)
		if err != nil {
			return fmt.Errorf("failed to load original task: %w", err)
		}

		if source == "" {
			return fmt.Errorf("task %s not found", taskID)
		}

		// Extract PR number from URL
		prNumber, err := github.ExtractPRNumber(prURL)
		if err != nil {
			return fmt.Errorf("invalid PR URL: %w", err)
		}

		// Create review task
		reviewTaskID := fmt.Sprintf("REVIEW-%d", prNumber)
		reviewTask := models.Task{
			ID:       reviewTaskID,
			Title:    fmt.Sprintf("Review PR #%d: %s", prNumber, originalTask.Title),
			Type:     "review",
			Status:   models.StatusPending,
			SpecRefs: originalTask.SpecRefs,
			Scope:    originalTask.Scope,
		}

		// Save to backlog (reviewer will claim it)
		backlog, err := tm.LoadTasks("backlog")
		if err != nil {
			return fmt.Errorf("failed to load backlog: %w", err)
		}

		backlog.Tasks = append(backlog.Tasks, reviewTask)
		if err := tm.SaveTasks("backlog", backlog); err != nil {
			return fmt.Errorf("failed to save review task: %w", err)
		}

		fmt.Fprintf(os.Stderr, "✅ Review task created: %s\n", reviewTaskID)
		fmt.Fprintf(os.Stderr, "\n🔍 Next steps for reviewer:\n")
		fmt.Fprintf(os.Stderr, "  1. agentic-agent task claim %s\n", reviewTaskID)
		fmt.Fprintf(os.Stderr, "  2. agentic-agent context build --task %s\n", reviewTaskID)
		if len(originalTask.SpecRefs) > 0 {
			fmt.Fprintf(os.Stderr, "  3. agentic-agent sdd gate-check %s\n", originalTask.SpecRefs[0])
		}
		fmt.Fprintf(os.Stderr, "  4. agentic-agent validate\n")
		fmt.Fprintf(os.Stderr, "  5. Review code, score quality, announce verdict\n")
		fmt.Fprintf(os.Stderr, "  6. agentic-agent task complete %s\n", reviewTaskID)

		// Output task ID for scripting
		fmt.Println(reviewTaskID)
		return nil
	},
}

var prStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check pull request status",
	Long: `Get current PR state (OPEN, MERGED, CLOSED) and reviewer info.

Example:
  agentic-agent pr status --pr-url https://github.com/owner/repo/pull/42`,

	RunE: func(cmd *cobra.Command, args []string) error {
		prURL, _ := cmd.Flags().GetString("pr-url")
		if prURL == "" {
			return fmt.Errorf("--pr-url is required")
		}

		// Get PR info from GitHub
		pr, err := github.GetPRInfo(prURL)
		if err != nil {
			return fmt.Errorf("failed to get PR info: %w", err)
		}

		// Human-readable output
		fmt.Printf("PR #%d: %s\n", pr.Number, pr.Title)
		fmt.Printf("Status: %s\n", pr.State)
		if len(pr.Reviewers) > 0 {
			fmt.Printf("Reviewers: %v\n", pr.Reviewers)
		}

		// JSON output if requested
		if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
			data, _ := json.MarshalIndent(pr, "", "  ")
			fmt.Println(string(data))
		}

		return nil
	},
}

func init() {
	// pr create flags
	prCreateCmd.Flags().StringP("task", "t", "", "task ID (required)")

	// pr review flags
	prReviewCmd.Flags().StringP("task", "t", "", "original task ID (required)")
	prReviewCmd.Flags().StringP("pr-url", "p", "", "PR URL (required)")

	// pr status flags
	prStatusCmd.Flags().StringP("pr-url", "p", "", "PR URL (required)")
	prStatusCmd.Flags().BoolP("json", "j", false, "output as JSON")

	// Register subcommands
	prCmd.AddCommand(prCreateCmd)
	prCmd.AddCommand(prReviewCmd)
	prCmd.AddCommand(prStatusCmd)

	// Register pr command to root
	rootCmd.AddCommand(prCmd)
}
