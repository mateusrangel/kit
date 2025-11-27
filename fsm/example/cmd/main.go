package main

import (
	"fmt"

	"github.com/mateusrangel/kit/fsm/example/internal/domain/dispute"
	"github.com/mateusrangel/kit/fsm/example/internal/domain/disputefsm"
)

func main() {
	orderFSM, err := disputefsm.NewDisputeFSM(dispute.New("123abc", "RECEIVED"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("BEFORE: %v\n", orderFSM.FSM.Current())
	_ = orderFSM.FSM.ExecEvent("VALIDATION_FAILED")
	fmt.Printf("AFTER: %v\n", orderFSM.FSM.Current())
	fmt.Println(orderFSM.FSM.Visualize())
}
