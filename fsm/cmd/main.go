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
	var events = []string{"VALIDATATION_SUCCEEDED", "VALIDATION_FAILED", "DISPUTE_WON", "DISPUTE_LOST"}
	m := fsm.New(order.State, states, events)
	orderFSM := &DisputeFSM{Dispute: order, FSM: m}

	transitions := []*fsm.Transition{
		{Event: "VALIDATATION_SUCCEEDED", Src: "RECEIVED", Dst: "PROCESSING"},
		{Event: "VALIDATION_FAILED", Src: "RECEIVED", Dst: "FINISHED", Actions: []fsm.Action{orderFSM.SendWarningMail}},
		{Event: "DISPUTE_WON", Src: "PROCESSING", Dst: "FINISHED"},
		{Event: "DISPUTE_LOST", Src: "PROCESSING", Dst: "FINISHED", Actions: []fsm.Action{orderFSM.SendWarningMail}},
	}

	err := orderFSM.FSM.AddTransitions(transitions)
	if err != nil {
		panic(err)
	}
	return orderFSM
}

func main() {
	orderFSM := NewOrderFSM(New("123abc", "RECEIVED"))
	fmt.Printf("BEFORE: %v\n", orderFSM.FSM.Current())
	_ = orderFSM.FSM.ExecEvent("VALIDATION_FAILED")
	fmt.Printf("AFTER: %v\n", orderFSM.FSM.Current())
	fmt.Println(fsm.Visualize(orderFSM.FSM))
}
