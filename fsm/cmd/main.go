package main

import (
	"fmt"

	"github.com/mateusrangel/kit/fsm"
)

type Dispute struct {
	Id    string
	State string
}

func NewDispute(id string, state string) *Dispute {
	return &Dispute{Id: id, State: state}
}

type DisputeService struct {
	Dispute *Dispute
	FSM     *fsm.FSM
}

func (o *DisputeService) SendWarningMail() bool {
	fmt.Println("EMAIL: DISPUTE STATE WAS TRANSITIONED TO", o.FSM.Current())
	return true
}

func NewDisputeService(d *Dispute) *DisputeService {
	disputeService := &DisputeService{Dispute: d}
	var states = []string{"RECEIVED", "PROCESSING", "FINISHED"}
	var events = []string{"VALIDATATION_SUCCEEDED", "VALIDATION_FAILED", "DISPUTE_WON", "DISPUTE_LOST"}
	transitions := []*fsm.Transition{
		{Event: "VALIDATATION_SUCCEEDED", Src: "RECEIVED", Dst: "PROCESSING"},
		{Event: "VALIDATION_FAILED", Src: "RECEIVED", Dst: "FINISHED", Actions: []fsm.Action{disputeService.SendWarningMail}},
		{Event: "DISPUTE_WON", Src: "PROCESSING", Dst: "FINISHED"},
		{Event: "DISPUTE_LOST", Src: "PROCESSING", Dst: "FINISHED", Actions: []fsm.Action{disputeService.SendWarningMail}},
	}
	m, err := fsm.New(d.State, states, events, transitions)
	if err != nil {
		panic(err)
	}
	disputeService.FSM = m
	return disputeService
}

func main() {
	orderFSM := NewDisputeService(NewDispute("123abc", "RECEIVED"))
	fmt.Printf("BEFORE: %v\n", orderFSM.FSM.Current())
	_ = orderFSM.FSM.ExecEvent("VALIDATION_FAILED")
	fmt.Printf("AFTER: %v\n", orderFSM.FSM.Current())
	fmt.Println(fsm.Visualize(orderFSM.FSM))
}
