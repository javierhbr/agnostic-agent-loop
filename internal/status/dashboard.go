package status

import (
	"github.com/javierbenavides/agentic-agent/internal/tasks"
	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// DashboardData holds aggregated project status information.
type DashboardData struct {
	ProjectName     string
	BacklogCount    int
	InProgressCount int
	DoneCount       int
	TotalCount      int
	CompletionPct   float64
	InProgressTasks []models.Task
	BacklogTasks    []models.Task
	NextReady       *models.Task
	Blockers        []string
	RecentEntries   []tasks.ProgressEntry
}

// Gather collects status data from the task manager and config.
func Gather(tm *tasks.TaskManager, cfg *models.Config) (*DashboardData, error) {
	backlog, err := tm.LoadTasks("backlog")
	if err != nil {
		return nil, err
	}
	inProgress, err := tm.LoadTasks("in-progress")
	if err != nil {
		return nil, err
	}
	done, err := tm.LoadTasks("done")
	if err != nil {
		return nil, err
	}

	d := &DashboardData{
		ProjectName:     cfg.Project.Name,
		BacklogCount:    len(backlog.Tasks),
		InProgressCount: len(inProgress.Tasks),
		DoneCount:       len(done.Tasks),
		InProgressTasks: inProgress.Tasks,
		BacklogTasks:    backlog.Tasks,
	}
	d.TotalCount = d.BacklogCount + d.InProgressCount + d.DoneCount
	if d.TotalCount > 0 {
		d.CompletionPct = float64(d.DoneCount) / float64(d.TotalCount) * 100
	}

	// Find next ready task and collect blockers
	for i := range backlog.Tasks {
		t := &backlog.Tasks[i]
		result := tasks.CanClaimTask(t, cfg)
		if result.Ready {
			if d.NextReady == nil {
				d.NextReady = t
			}
		} else {
			for _, check := range result.Checks {
				if !check.Passed {
					d.Blockers = append(d.Blockers, t.ID+": "+check.Message)
				}
			}
		}
	}

	// Load recent progress entries
	if cfg.Paths.ProgressYAMLPath != "" {
		pw := tasks.NewProgressWriter(cfg.Paths.ProgressTextPath, cfg.Paths.ProgressYAMLPath)
		entries, err := pw.GetAllEntries()
		if err == nil && len(entries) > 0 {
			// Show last 5 entries
			start := 0
			if len(entries) > 5 {
				start = len(entries) - 5
			}
			d.RecentEntries = entries[start:]
		}
	}

	return d, nil
}
