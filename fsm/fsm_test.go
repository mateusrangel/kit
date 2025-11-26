package fsm_test

import (
	"fmt"
	"testing"

	"github.com/mateusrangel/kit/fsm"
)

var saveStateToDb = func() bool { fmt.Println("SAVING STATE TO DB"); return true }
var sendValitationSucceededEmail = func() bool { fmt.Println("SENDING VALIDATION SUCCEEDED EMAIL"); return true }
var states = []string{"RECEIVED", "PROCESSING", "FINISHED"}

func TestMachine_TestTransitions(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		model   any
		states  []string
		initial string
		// Named input parameters for target function.
		trigger       string
		source        string
		dest          string
		triggerToExec string
		wantErr       bool
	}{
		{
			name:          "Valid Transition",
			model:         nil,
			states:        states,
			initial:       "RECEIVED",
			trigger:       "VALIDATED",
			source:        "RECEIVED",
			dest:          "PROCESSING",
			triggerToExec: "VALIDATED",
			wantErr:       false,
		},
		{
			name:          "Invalid Transition",
			model:         nil,
			states:        states,
			initial:       "RECEIVED",
			trigger:       "VALIDATED",
			source:        "RECEIVED",
			dest:          "PROCESSING",
			triggerToExec: "PROCESSED",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := fsm.New(tt.model, tt.states, tt.initial)
			actions := []fsm.Action{saveStateToDb, sendValitationSucceededEmail}
			m.AddTransition(tt.trigger, tt.source, tt.dest, actions)
			gotErr := m.ExecTrigger(tt.triggerToExec)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ExecTrigger() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ExecTrigger() succeeded unexpectedly")
			}
		})
	}
}
