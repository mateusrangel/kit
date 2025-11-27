package main

import (
	"fmt"

	"github.com/mateusrangel/kit/fsm"
)

type Dispute struct {
	Id    string
	State string
}

func New(id string, state string) *Dispute {
	return &Dispute{Id: id, State: state}
}

type DisputeFSM struct {
	Dispute *Dispute
	FSM     *fsm.FSM
}

func (o *DisputeFSM) SendWarningMail() bool {
	fmt.Println("EMAIL: DISPUTE STATE WAS TRANSITIONED TO", o.FSM.Current())
	return true
}

func NewOrderFSM(order *Dispute) *DisputeFSM {
	var states = []string{"RECEIVED", "PROCESSING", "FINISHED"}
	m := fsm.New(order.State, states)
	orderFSM := &DisputeFSM{Dispute: order, FSM: m}
	orderFSM.FSM.AddTransition("VALIDATATION_SUCCEEDED", "RECEIVED", "PROCESSING", []fsm.Action{})
	orderFSM.FSM.AddTransition("VALIDATION_FAILED", "RECEIVED", "FINISHED", []fsm.Action{orderFSM.SendWarningMail})
	orderFSM.FSM.AddTransition("DISPUTE_WON", "PROCESSING", "FINISHED", []fsm.Action{})
	orderFSM.FSM.AddTransition("DISPUTE_LOST", "PROCESSING", "FINISHED", []fsm.Action{orderFSM.SendWarningMail})
	return orderFSM
}

func main() {
	orderFSM := NewOrderFSM(New("123abc", "RECEIVED"))
	fmt.Printf("BEFORE: %v\n", orderFSM.FSM.Current())
	_ = orderFSM.FSM.ExecEvent("VALIDATION_FAILED")
	fmt.Printf("AFTER: %v\n", orderFSM.FSM.Current())
	fmt.Println(fsm.Visualize(orderFSM.FSM))
}
