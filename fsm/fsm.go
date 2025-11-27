package fsm

import "errors"

type State string
type Action func() bool
type StateActionTuple struct {
	NextState State
	Actions   []Action
}
type Trigger string

type Stateful interface {
	GetState() string
	SetState(newState string)
}
type Machine[T Stateful] struct {
	Model       T
	transitions map[State]map[Trigger]*StateActionTuple
	States      []State
}

func New[T Stateful](model T, states []string) *Machine[T] {
	stateSlice := make([]State, len(states))
	for i, s := range states {
		stateSlice[i] = State(s)
	}
	return &Machine[T]{transitions: make(map[State]map[Trigger]*StateActionTuple), Model: model, States: stateSlice}
}

func (m *Machine[T]) AddTransition(trigger string, source string, dest string, actions []Action) {
	srcState := State(source)
	trig := Trigger(trigger)
	destState := State(dest)

	if _, ok := m.transitions[srcState]; !ok {
		m.transitions[srcState] = make(map[Trigger]*StateActionTuple)
	}

	m.transitions[srcState][trig] = &StateActionTuple{NextState: destState, Actions: actions}
}

func (m *Machine[T]) ExecTrigger(trigger string) error {
	stateActionTuple, ok := m.transitions[State(m.Model.GetState())][Trigger(trigger)]
	if !ok {
		return errors.New("invalid Transition")
	}
	m.Model.SetState(string(stateActionTuple.NextState))
	for _, action := range stateActionTuple.Actions {
		action()
	}
	return nil
}
