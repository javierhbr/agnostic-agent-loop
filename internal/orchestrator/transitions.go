package orchestrator

import "fmt"

func (sm *StateMachine) HandleEvent(event Event) error {
	switch sm.CurrentState {
	case StateIdle:
		if event == EventTaskStarted {
			sm.CurrentState = StatePlanning
			return nil
		}
	case StatePlanning:
		if event == EventPlanApproved {
			sm.CurrentState = StateExecution
			return nil
		}
	case StateExecution:
		if event == EventWorkCompleted {
			sm.CurrentState = StateVerification
			return nil
		}
	case StateVerification:
		if event == EventVerificationPass {
			sm.CurrentState = StateDone
			return nil
		}
		if event == EventVerificationFail {
			sm.CurrentState = StateExecution // Back to fix
			return nil
		}
	case StateDone:
		return fmt.Errorf("cannot transition from DONE state")
	}

	return fmt.Errorf("invalid transition from %s with event %s", sm.CurrentState, event)
}
