package main

import (
	"fmt"

	"github.com/mateusrangel/kit/fsm"
)

type Order struct {
	Id    string
	State string
}

func New(id string, state string) *Order {
	return &Order{Id: id, State: state}
}

type OrderFSM struct {
	Order *Order
	FSM   *fsm.FSM
}

func NewOrderFSM(order *Order) *OrderFSM {
	var states = []string{"RECEIVED", "PROCESSING", "FINISHED"}
	m := fsm.New(order.State, states)
	m.AddTransition("VALIDATATION_SUCCEEDED", "RECEIVED", "PROCESSING", []fsm.Action{})
	m.AddTransition("VALIDATION_FAILED", "RECEIVED", "FINISHED", []fsm.Action{})
	m.AddTransition("DISPUTE_WON", "PROCESSING", "FINISHED", []fsm.Action{})
	m.AddTransition("DISPUTE_LOST", "PROCESSING", "FINISHED", []fsm.Action{})

	return &OrderFSM{Order: order, FSM: m}
}

func main() {
	orderFSM := NewOrderFSM(New("123abc", "RECEIVED"))
	fmt.Println(orderFSM.FSM.AvailableTransitions())
	fmt.Printf("BEFORE: %v\n", orderFSM.FSM.Current())
	fmt.Println(orderFSM.FSM.Can("VALIDATATION_SUCCEEDED"))
	fmt.Println(orderFSM.FSM.Can("DELIVERED"))
	_ = orderFSM.FSM.ExecEvent("VALIDATATION_SUCCEEDED")
	fmt.Printf("AFTER: %v\n", orderFSM.FSM.Current())
	fmt.Println(orderFSM.FSM.Can("VALIDATATION_SUCCEEDED"))
	fmt.Println(orderFSM.FSM.Can("DELIVERED"))
	fmt.Println(fsm.Visualize(orderFSM.FSM))
}
