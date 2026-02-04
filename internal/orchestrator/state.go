package orchestrator

type State string

const (
	StateIdle         State = "IDLE"
	StatePlanning     State = "PLANNING"
	StateExecution    State = "EXECUTION"
	StateVerification State = "VERIFICATION"
	StateDone         State = "DONE"
)

type Event string

const (
	EventTaskStarted      Event = "TASK_STARTED"
	EventPlanApproved     Event = "PLAN_APPROVED"
	EventWorkCompleted    Event = "WORK_COMPLETED"
	EventVerificationPass Event = "VERIFICATION_PASS"
	EventVerificationFail Event = "VERIFICATION_FAIL"
)

type StateMachine struct {
	CurrentState State
}

func NewStateMachine(initial State) *StateMachine {
	if initial == "" {
		initial = StateIdle
	}
	return &StateMachine{CurrentState: initial}
}
