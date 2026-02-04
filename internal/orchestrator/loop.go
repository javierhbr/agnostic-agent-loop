package orchestrator

import (
	"fmt"
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
