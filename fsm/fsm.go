// An FSM is defined by a list of its states, its initial state, and the inputs that trigger each transition.
// A state is a description of the status of a system that is waiting to execute a transition.
package fsm

import (
	"errors"
)

type State string
type Action func() bool
type StateActionTuple struct {
	NextState State
	Actions   []Action
}
type Event string

type FSM struct {
	state       State
	States      []State
	transitions map[State]map[Event]*StateActionTuple
}

type Transition struct {
	Trigger string
	Src     string
	Dst     string
	Actions []Action
}

func New(initial string, states []string) *FSM {
	stateSlice := make([]State, len(states))
	for i, s := range states {
		stateSlice[i] = State(s)
	}
	return &FSM{state: State(initial), transitions: make(map[State]map[Event]*StateActionTuple), States: stateSlice}
}

func (m *FSM) AddTransition(trigger string, source string, dest string, actions []Action) {
	srcState := State(source)
	trig := Event(trigger)
	destState := State(dest)

	if _, ok := m.transitions[srcState]; !ok {
		m.transitions[srcState] = make(map[Event]*StateActionTuple)
	}

	m.transitions[srcState][trig] = &StateActionTuple{NextState: destState, Actions: actions}
}

func (m *FSM) ExecEvent(event string) error {
	stateActionTuple, ok := m.transitions[State(m.state)][Event(event)]
	if !ok {
		return errors.New("invalid Transition")
	}
	m.state = stateActionTuple.NextState
	for _, action := range stateActionTuple.Actions {
		action()
	}
	return nil
}

func (f *FSM) AvailableTransitions() []string {
	avaiableTransactions := make([]string, 0, len(f.transitions[f.state]))
	for k := range f.transitions[f.state] {
		avaiableTransactions = append(avaiableTransactions, string(k))
	}
	return avaiableTransactions
}

func (f *FSM) Can(event string) bool {
	_, ok := f.transitions[f.state][Event(event)]
	return ok
}

func (f *FSM) Current() string {
	return string(f.state)
}

func (f *FSM) getStates() []string {
	states := make([]string, 0)
	for _, state := range f.States {
		states = append(states, string(state))
	}
	return states
}
