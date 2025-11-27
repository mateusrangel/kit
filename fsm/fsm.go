// An FSM is defined by a list of its states, its initial state, and the inputs that trigger each transition.
// A state is a description of the status of a system that is waiting to execute a transition.
package fsm

import (
	"errors"
	"fmt"
	"slices"
)

type State string
type Action func() bool
type StateActionTuple struct {
	NextState State
	Actions   []Action
}
type Event string

type Transition struct {
	Event   string
	Src     string
	Dst     string
	Actions []Action
}

var exists = struct{}{}

type FSM struct {
	state       State
	stateSet    map[State]struct{}
	eventSet    map[Event]struct{}
	transitions map[State]map[Event]*StateActionTuple
}

func New(initial string, states []string, events []string) *FSM {
	machine := &FSM{state: State(initial), transitions: make(map[State]map[Event]*StateActionTuple), stateSet: make(map[State]struct{}), eventSet: make(map[Event]struct{})}
	for _, v := range states {
		machine.addState(State(v))
	}
	for _, v := range events {
		machine.addEvent(Event(v))
	}
	return machine
}

func (m *FSM) addState(state State) {
	m.stateSet[state] = exists
}

func (m *FSM) addEvent(event Event) {
	m.eventSet[event] = exists
}

func (m *FSM) AddTransitions(ts []*Transition) error {
	for _, t := range ts {
		err := m.AddTransition(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *FSM) AddTransition(t *Transition) error {
	evt := Event(t.Event)
	srcState := State(t.Src)
	destState := State(t.Dst)

	_, exists := m.eventSet[evt]
	if !exists {
		return fmt.Errorf("event %s doesn't exist", evt)
	}

	_, exists = m.stateSet[srcState]
	if !exists {
		return fmt.Errorf("source state %s doesn't exist", srcState)
	}

	_, exists = m.stateSet[destState]
	if !exists {
		return fmt.Errorf("destination state %s doesn't exist", destState)
	}

	if _, ok := m.transitions[srcState]; !ok {
		m.transitions[srcState] = make(map[Event]*StateActionTuple)
	}

	m.transitions[srcState][evt] = &StateActionTuple{NextState: destState, Actions: t.Actions}
	return nil
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

func (f *FSM) GetStates() []string {
	states := make([]string, 0)
	for state := range f.stateSet {
		states = append(states, string(state))
	}

	return states
}

func (f *FSM) getSortedStates() []string {
	states := make([]string, 0)
	for state := range f.stateSet {
		states = append(states, string(state))
	}
	slices.Sort(states)
	return states
}

func getSortedStateKeys(transitions map[State]map[Event]*StateActionTuple) []State {
	keys := make([]State, 0, len(transitions))
	for state := range transitions {
		keys = append(keys, state)
	}
	slices.Sort(keys)
	return keys
}

func getSortedEventKeys(eventMap map[Event]*StateActionTuple) []Event {
	keys := make([]Event, 0, len(eventMap))
	for event := range eventMap {
		keys = append(keys, event)
	}
	slices.Sort(keys)
	return keys
}
