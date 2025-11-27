package fsm

import (
	"bytes"
	"fmt"
)

// Visualize outputs a visualization of a FSM in Graphviz format.
func Visualize(fsm *FSM) string {
	var buf bytes.Buffer

	writeHeaderLine(&buf)
	writeTransitions(&buf, fsm.transitions)
	writeStates(&buf, string(fsm.state), fsm.getStates())
	writeFooter(&buf)

	return buf.String()
}

func writeHeaderLine(buf *bytes.Buffer) {
	buf.WriteString(`digraph fsm {`)
	buf.WriteString("\n")
}

func writeTransitions(buf *bytes.Buffer, transitions map[State]map[Event]*StateActionTuple) {
	for state, eventMap := range transitions {
		for event, stateActionTuple := range eventMap {
			v := event
			buf.WriteString(fmt.Sprintf(`    "%s" -> "%s" [ label = "%s" ];`, state, string(stateActionTuple.NextState), v))
			buf.WriteString("\n")
		}

	}

	buf.WriteString("\n")
}

func writeStates(buf *bytes.Buffer, current string, sortedStateKeys []string) {
	for _, k := range sortedStateKeys {
		if k == current {
			buf.WriteString(fmt.Sprintf(`    "%s" [color = "red"];`, k))
		} else {
			buf.WriteString(fmt.Sprintf(`    "%s";`, k))
		}
		buf.WriteString("\n")
	}
}

func writeFooter(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintln("}"))
}
