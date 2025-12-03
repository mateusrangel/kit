package main

import (
	"fmt"

	"github.com/mateusrangel/kit/fsm"
)

// Const States
const (
	StateReceived    = "RECEIVED"
	StateCreateClaim = "CREATE_CLAIM"
	StateProcessing  = "PROCESSING"
	StateFinished    = "FINISHED"
)

// Events
const (
	EventValidationSucceded = "VALIDATION_SUCCEEDED"
	EventValidationFailed   = "VALIDATION_FAILED"
	EventClaimCreated       = "CLAIM_CREATED"
	EventDisputeWon         = "DISPUTE_WON"
	EventDispotLost         = "DISPUTE_LOST"
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
	fmt.Printf("EMAIL: DISPUTE %s STATE WAS TRANSITIONED TO %s\n", o.Dispute.Id, o.FSM.Current())
	return true
}

func NewDisputeService(d *Dispute) *DisputeService {
	disputeService := &DisputeService{Dispute: d}
	var states = []string{StateReceived, StateCreateClaim, StateProcessing, StateFinished}
	var events = []string{EventValidationSucceded, EventValidationFailed, EventClaimCreated, EventDisputeWon, EventDispotLost}
	transitions := []*fsm.Transition{
		{Event: EventValidationSucceded, Src: StateReceived, Dst: StateCreateClaim},
		{Event: EventValidationFailed, Src: StateReceived, Dst: StateFinished, Actions: []fsm.Action{disputeService.SendWarningMail}},
		{Event: EventClaimCreated, Src: StateCreateClaim, Dst: StateProcessing},
		{Event: EventDisputeWon, Src: StateProcessing, Dst: StateFinished},
		{Event: EventDispotLost, Src: StateProcessing, Dst: StateFinished, Actions: []fsm.Action{disputeService.SendWarningMail}},
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
	_ = orderFSM.FSM.ExecEvent("VALIDATION_SUCCEEDED")
	fmt.Printf("AFTER: %v\n", orderFSM.FSM.Current())

	fmt.Printf("BEFORE: %v\n", orderFSM.FSM.Current())
	_ = orderFSM.FSM.ExecEvent("CLAIM_CREATED")
	fmt.Printf("AFTER: %v\n", orderFSM.FSM.Current())

	fmt.Println(fsm.Visualize(orderFSM.FSM))
}
