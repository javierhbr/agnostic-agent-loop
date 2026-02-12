package orchestrator

import (
	"context"
	"fmt"
	"os"

	appcontext "github.com/javierbenavides/agentic-agent/internal/context"
	"github.com/javierbenavides/agentic-agent/internal/encoding"
	"github.com/javierbenavides/agentic-agent/internal/openspec"
	"github.com/javierbenavides/agentic-agent/internal/skills"
	"github.com/javierbenavides/agentic-agent/internal/specs"
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/internal/tracks"
	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// AutopilotLoop orchestrates claiming, context building, and task iteration.
type AutopilotLoop struct {
	cfg           *models.Config
	maxIterations int
	stopSignal    string
	dryRun        bool
	taskManager   *tasks.TaskManager
	specResolver  *specs.Resolver
	trackManager  *tracks.Manager
}

// NewAutopilotLoop creates a new autopilot loop.
func NewAutopilotLoop(cfg *models.Config, maxIterations int, stopSignal string, dryRun bool) *AutopilotLoop {
	if maxIterations <= 0 {
		maxIterations = 10
	}
	if stopSignal == "" {
		stopSignal = "<promise>COMPLETE</promise>"
	}
	return &AutopilotLoop{
		cfg:           cfg,
		maxIterations: maxIterations,
		stopSignal:    stopSignal,
		dryRun:        dryRun,
		taskManager:   tasks.NewTaskManager(".agentic/tasks"),
		specResolver:  specs.NewResolver(cfg),
		trackManager:  tracks.NewManager(cfg.Paths.TrackDir),
	}
}

// Run executes the autopilot loop.
func (a *AutopilotLoop) Run(ctx context.Context) error {
	// Ensure agent skills are set up before starting
	if a.cfg.ActiveAgent != "" {
		result, err := skills.Ensure(a.cfg.ActiveAgent, a.cfg)
		if err != nil {
			fmt.Printf("Warning: could not ensure agent skills: %v\n", err)
		} else if result.RulesGenerated || result.DriftFixed || len(result.PacksInstalled) > 0 {
			fmt.Print(skills.FormatEnsureResult(result))
		}
	}

	// Auto-import tasks from draft openspec changes
	if a.cfg.Paths.OpenSpecDir != "" {
		om := openspec.NewManager(a.cfg.Paths.OpenSpecDir)
		syncResult, _ := om.Sync(a.taskManager)
		if syncResult != nil && len(syncResult.ChangesImported) > 0 {
			fmt.Printf("Auto-imported %d tasks from %d change(s)\n",
				syncResult.TasksCreated, len(syncResult.ChangesImported))
		}
	}

	user := os.Getenv("USER")
	if user == "" {
		user = "autopilot"
	}

	for iteration := 1; iteration <= a.maxIterations; iteration++ {
		select {
		case <-ctx.Done():
			fmt.Println("Autopilot cancelled.")
			return ctx.Err()
		default:
		}

		// 1. Find next claimable task
		task, err := a.findNextTask()
		if err != nil {
			return fmt.Errorf("iteration %d: %w", iteration, err)
		}
		if task == nil {
			fmt.Println("All tasks complete. Autopilot finished.")
			return nil
		}

		fmt.Printf("\n--- Iteration %d/%d ---\n", iteration, a.maxIterations)
		fmt.Printf("Next task: [%s] %s\n", task.ID, task.Title)

		// 2. Run readiness checks
		result := tasks.CanClaimTask(task, a.cfg)
		fmt.Print(tasks.FormatReadinessResult(result))

		if a.dryRun {
			fmt.Printf("[DRY RUN] Would claim task %s and generate context\n", task.ID)
			continue
		}

		// 3. Claim task
		if err := a.taskManager.ClaimTaskWithConfig(task.ID, user, a.cfg); err != nil {
			fmt.Printf("Warning: could not claim task %s: %v\n", task.ID, err)
			continue
		}
		fmt.Printf("Claimed task %s\n", task.ID)

		// 4. Generate context for scope dirs
		for _, dir := range task.Scope {
			dirCtx, err := appcontext.GenerateContextWithConfig(dir, a.cfg)
			if err != nil {
				fmt.Printf("  Warning: context generation failed for %s: %v\n", dir, err)
				continue
			}
			dcm := appcontext.NewDirectoryContextManager(dir)
			if err := dcm.SaveContext(dir, dirCtx); err != nil {
				fmt.Printf("  Warning: could not save context for %s: %v\n", dir, err)
				continue
			}
			fmt.Printf("  Generated context for %s\n", dir)
		}

		// 5. Build context bundle (with resolved specs)
		bundle, err := encoding.CreateContextBundle(task.ID, "toon", a.cfg)
		if err != nil {
			fmt.Printf("  Warning: could not build context bundle: %v\n", err)
		} else {
			fmt.Printf("  Context bundle built (%d bytes)\n", len(bundle))
		}

		// 6. Report task ready for agent execution
		fmt.Printf("Task %s is ready for agent execution.\n", task.ID)
	}

	fmt.Printf("Reached max iterations (%d). Stopping autopilot.\n", a.maxIterations)
	return nil
}

// findNextTask finds the next claimable task from the backlog.
// Prefers tasks where readiness checks all pass.
func (a *AutopilotLoop) findNextTask() (*models.Task, error) {
	backlog, err := a.taskManager.LoadTasks("backlog")
	if err != nil {
		return nil, err
	}

	if len(backlog.Tasks) == 0 {
		// Check if there are in-progress tasks still running
		inProgress, err := a.taskManager.LoadTasks("in-progress")
		if err != nil {
			return nil, err
		}
		if len(inProgress.Tasks) > 0 {
			return nil, fmt.Errorf("no backlog tasks but %d still in progress", len(inProgress.Tasks))
		}
		return nil, nil // All done
	}

	// Prefer tasks that are fully ready and whose track (if any) is active
	for _, t := range backlog.Tasks {
		if a.isTaskBlocked(&t) {
			continue
		}
		result := tasks.CanClaimTask(&t, a.cfg)
		if result.Ready {
			return &t, nil
		}
	}

	// Fall back to first unblocked task in backlog
	for _, t := range backlog.Tasks {
		if !a.isTaskBlocked(&t) {
			return &t, nil
		}
	}

	return nil, nil
}

// isTaskBlocked returns true if a task is linked to a track that is not yet active.
func (a *AutopilotLoop) isTaskBlocked(task *models.Task) bool {
	if task.TrackID == "" {
		return false
	}
	track, err := a.trackManager.Get(task.TrackID)
	if err != nil {
		return false // track not found — don't block
	}
	return track.Status == models.TrackStatusIdeation || track.Status == models.TrackStatusPlanning
}
