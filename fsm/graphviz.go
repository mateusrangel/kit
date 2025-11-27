package fsm

import (
	"bytes"
	"fmt"
	"slices"
)

// Visualize outputs a visualization of a FSM in Graphviz format.
func Visualize(fsm *FSM) string {
	var buf bytes.Buffer

	sortedSrcStates := getSortedStateKeys(fsm.transitions)

	writeHeaderLine(&buf)
	writeTransitions(&buf, fsm.transitions, sortedSrcStates)
	writeStates(&buf, string(fsm.state), fsm.getStates())
	writeFooter(&buf)

	return buf.String()
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

func writeTransitions(buf *bytes.Buffer, transitions map[State]map[Event]*StateActionTuple, sortedSrcStates []State) {

	for _, state := range sortedSrcStates {
		eventMap := transitions[state]

		sortedEvents := getSortedEventKeys(eventMap)
		for _, event := range sortedEvents {
			stateActionTuple := eventMap[event]
			v := event
			fmt.Fprintf(buf, `    "%s" -> "%s" [ label = "%s" ];`, state, string(stateActionTuple.NextState), v)
			buf.WriteString("\n")
		}
	}

	buf.WriteString("\n")
}

func writeHeaderLine(buf *bytes.Buffer) {
	buf.WriteString(`digraph fsm {`)
	buf.WriteString("\n")
}

func writeStates(buf *bytes.Buffer, current string, sortedStateKeys []string) {
	for _, k := range sortedStateKeys {
		if k == current {
			fmt.Fprintf(buf, `    "%s" [color = "red"];`, k)
		} else {
			fmt.Fprintf(buf, `    "%s";`, k)
		}
		buf.WriteString("\n")
	}
}

func writeFooter(buf *bytes.Buffer) {
	fmt.Fprintln(buf, "}")
}
