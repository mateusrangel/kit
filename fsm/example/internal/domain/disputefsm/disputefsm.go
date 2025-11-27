package disputefsm

import (
	"fmt"

	"github.com/mateusrangel/kit/fsm"
	"github.com/mateusrangel/kit/fsm/example/internal/domain/dispute"
)

// --- 1. Constants and Types ---

// DisputeStates defines the possible states in the FSM.
const (
	StateReceived   fsm.State = "RECEIVED"
	StateProcessing fsm.State = "PROCESSING"
	StateFinished   fsm.State = "FINISHED"
)

// DisputeEvents defines the triggers for state transitions.
const (
	EventValidationSucceeded fsm.Event = "VALIDATION_SUCCEEDED" // Corrected typo
	EventValidationFailed    fsm.Event = "VALIDATION_FAILED"
	EventDisputeWon          fsm.Event = "DISPUTE_WON"
	EventDisputeLost         fsm.Event = "DISPUTE_LOST"
)

var allStates = []fsm.State{
	StateReceived,
	StateProcessing,
	StateFinished,
}

var allEvents = []fsm.Event{
	EventValidationSucceeded,
	EventValidationFailed,
	EventDisputeWon,
	EventDisputeLost,
}

// DisputeMachine wraps the core domain object and the FSM instance.
type DisputeMachine struct {
	Dispute *dispute.Dispute
	FSM     *fsm.FSM
}

// --- 2. Action Handlers ---

// SendWarningMail is an FSM action triggered by certain events.
func (dm *DisputeMachine) SendWarningMail() bool {
	// Note: We use dm.FSM.Current() to get the destination state after transition
	fmt.Printf("EMAIL SENT: Dispute ID %s state transitioned to %s\n", dm.Dispute.Id, dm.FSM.Current())
	return true
}

// --- 3. Constructor and Initialization ---

// DefineTransitions sets up all valid transitions for the Dispute FSM.
func (dm *DisputeMachine) DefineTransitions() []*fsm.Transition {
	return []*fsm.Transition{
		{
			Event: EventValidationSucceeded,
			Src:   StateReceived,
			Dst:   StateProcessing,
		},
		{
			Event:   EventValidationFailed,
			Src:     StateReceived,
			Dst:     StateFinished,
			Actions: []fsm.Action{dm.SendWarningMail}, // Action triggered on failure
		},
		{
			Event: EventDisputeWon,
			Src:   StateProcessing,
			Dst:   StateFinished,
		},
		{
			Event:   EventDisputeLost,
			Src:     StateProcessing,
			Dst:     StateFinished,
			Actions: []fsm.Action{dm.SendWarningMail}, // Action triggered on loss
		},
	}
}

// NewDisputeFSM creates and initializes a new FSM wrapper for a Dispute.
// It returns the initialized machine or an error if transition setup fails.
func NewDisputeFSM(disputeData *dispute.Dispute) (*DisputeMachine, error) {
	// 1. Initialize core FSM
	m := fsm.New(disputeData.State, allStates, allEvents) // assuming fsm.New only needs states

	// 2. Wrap FSM and domain object
	machine := &DisputeMachine{
		Dispute: disputeData,
		FSM:     m,
	}

	// 3. Add transitions
	transitions := machine.DefineTransitions()
	err := machine.FSM.AddTransitions(transitions)
	if err != nil {
		// Return a wrapped error instead of panicking
		return nil, fmt.Errorf("failed to add FSM transitions: %w", err)
	}

	return machine, nil
}
