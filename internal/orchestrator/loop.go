package orchestrator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/javierbenavides/agentic-agent/internal/tasks"
)

// RunLoop is the main entry point for the agent's autonomous loop.
// For the MVP, this will just simulate the loop or run one step.
func RunLoop(taskID string) error {
	fmt.Printf("Starting orchestrator for task %s...\n", taskID)

	// 1. Load Task
	tm := tasks.NewTaskManager(".agentic/tasks")
	// Simplified loading... (logic duplicated from context build, should ideally be refactored to a helper)
	// We'll assume the task is in progress for now
	list, err := tm.LoadTasks("in-progress")
	if err != nil {
		return err
	}

	found := false
	for _, t := range list.Tasks {
		if t.ID == taskID {
			found = true
			break
		}
	}
	if !found {
		// Try backlog and move it?
		return fmt.Errorf("task %s not found in in-progress list", taskID)
	}

	// 2. Initialize State Machine
	sm := NewStateMachine(StateIdle)

	// 3. Start Loop
	// In a real agent, this would be an infinite loop waiting for LLM responses.
	// For this CLI tool, we might just print the current state and what needs to happen.

	// Simulate startup
	fmt.Println("State: IDLE")
	if err := sm.HandleEvent(EventTaskStarted); err != nil {
		return err
	}
	fmt.Println("State: PLANNING")
	fmt.Println(">> user should now create implementation_plan.md <<")

	// We could poll or look for file existence here...
	time.Sleep(500 * time.Millisecond)

	return nil
}

// Loop represents an autonomous agent loop with stop conditions
type Loop struct {
	maxIterations int
	stopSignal    string
	taskManager   *tasks.TaskManager
}

// NewLoop creates a new agent loop
func NewLoop(maxIterations int, stopSignal string, taskManager *tasks.TaskManager) *Loop {
	if stopSignal == "" {
		stopSignal = "<promise>COMPLETE</promise>"
	}
	if maxIterations <= 0 {
		maxIterations = 10
	}
	return &Loop{
		maxIterations: maxIterations,
		stopSignal:    stopSignal,
		taskManager:   taskManager,
	}
}

// Run executes the agent loop with stop condition detection
func (l *Loop) Run(ctx context.Context) error {
	for iteration := 1; iteration <= l.maxIterations; iteration++ {
		// Run agent iteration (placeholder - actual implementation would call agent)
		output, err := l.runIteration(ctx)
		if err != nil {
			return fmt.Errorf("iteration %d failed: %w", iteration, err)
		}

		// Check for stop condition
		if l.checkStopCondition(output) {
			fmt.Printf("Stop condition detected at iteration %d\n", iteration)
			return nil
		}

		// Check if all tasks completed
		if l.allTasksComplete() {
			fmt.Printf("All tasks completed at iteration %d\n", iteration)
			return nil
		}

		fmt.Printf("Iteration %d complete. Continuing...\n", iteration)
	}

	return fmt.Errorf("reached max iterations (%d) without completion", l.maxIterations)
}

// runIteration runs a single iteration of the agent loop
func (l *Loop) runIteration(ctx context.Context) (string, error) {
	// Placeholder - in real implementation, this would:
	// 1. Load next task from backlog
	// 2. Spawn agent with task context
	// 3. Capture agent output
	// 4. Return output for stop condition checking

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		// Simulate agent work
		time.Sleep(100 * time.Millisecond)
		return "Agent iteration output", nil
	}
}

// checkStopCondition checks if the output contains the stop signal
func (l *Loop) checkStopCondition(output string) bool {
	return strings.Contains(output, l.stopSignal)
}

// allTasksComplete checks if all tasks are in the done state
func (l *Loop) allTasksComplete() bool {
	// Check backlog and in-progress for any remaining tasks
	backlog, err := l.taskManager.LoadTasks("backlog")
	if err != nil || len(backlog.Tasks) > 0 {
		return false
	}

	inProgress, err := l.taskManager.LoadTasks("in-progress")
	if err != nil || len(inProgress.Tasks) > 0 {
		return false
	}

	return true
}
