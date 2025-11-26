package fsm

import "errors"

type State string
type Action func() bool
type StateActions struct {
	State   State
	Actions []Action
}
type Trigger string

type Machine[T any] struct {
	currState   State
	transitions map[State]map[Trigger]*StateActions
	Initial     State
	Model       T
	States      []State
}

func New[T any](model T, states []string, initial string) *Machine[T] {
	stateSlice := make([]State, len(states))
	for i, s := range states {
		stateSlice[i] = State(s)
	}
	return &Machine[T]{currState: State(initial), transitions: make(map[State]map[Trigger]*StateActions), Model: model, States: stateSlice, Initial: State(initial)}
}

func (m *Machine[T]) AddTransition(trigger string, source string, dest string, actions []Action) {
	srcState := State(source)
	trig := Trigger(trigger)
	destState := State(dest)

	if _, ok := m.transitions[srcState]; !ok {
		m.transitions[srcState] = make(map[Trigger]*StateActions)
	}

	m.transitions[srcState][trig] = &StateActions{State: destState, Actions: actions}
}

func (m *Machine[T]) ExecTrigger(trigger string) error {
	stateActions, ok := m.transitions[m.currState][Trigger(trigger)]
	if !ok {
		return errors.New("invalid Transition")
	}
	m.currState = stateActions.State
	for _, action := range stateActions.Actions {
		action()
	}
	return nil
}
