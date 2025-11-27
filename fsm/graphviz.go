package fsm

import (
	"bytes"
	"fmt"
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

func writeHeaderLine(buf *bytes.Buffer) {
	buf.WriteString(`digraph fsm {`)
	buf.WriteString("\n")
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
