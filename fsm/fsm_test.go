package fsm_test

import (
	"fmt"
	"testing"

	"github.com/mateusrangel/kit/fsm"
)

type Order struct {
	Id      string
	Machine *fsm.Machine[*Order]
}

func New(id string) *Order {
	order := &Order{Id: id}
	var states = []string{"RECEIVED", "PROCESSING", "FINISHED"}

	m := fsm.New(order, states, states[0])
	actions := []fsm.Action{m.Model.WhoAmi}
	m.AddTransition("VALIDATED", "RECEIVED", "PROCESSING", actions)
	order.Machine = m
	return order
}
func (o *Order) WhoAmi() bool {
	fmt.Println(o.Id)
	return true
}

func TestMachine_TestTransitions(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		model any
		// Named input parameters for target function.
		triggerToExec string
		wantErr       bool
	}{
		{
			name:          "Valid Transition",
			model:         nil,
			triggerToExec: "VALIDATED",
			wantErr:       false,
		},
		{
			name:          "Invalid Transition",
			model:         nil,
			triggerToExec: "PROCESSED",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := New("123abc")

			gotErr := order.Machine.ExecTrigger(tt.triggerToExec)
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
