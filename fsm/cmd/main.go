package main

import (
	"fmt"

	"github.com/mateusrangel/kit/fsm"
)

type Order struct {
	Id    string
	State string
}

func (o *Order) GetState() string {
	return o.State
}

func (o *Order) SetState(newState string) {
	o.State = newState
}

func New(id string, status string) *Order {
	return &Order{Id: id, State: status}
}

type OrderFSM struct {
	*fsm.Machine[*Order]
}

func SendMail() bool {
	fmt.Println("SENDING E-MAIL")
	return true
}

func NewOrderFSM(order *Order) *OrderFSM {
	var states = []string{"RECEIVED", "PROCESSING", "FINISHED"}
	m := fsm.New(order, states)
	m.AddTransition("VALIDATED", "RECEIVED", "PROCESSING", []fsm.Action{SendMail})
	m.AddTransition("DELIVERED", "PROCESSING", "FINISHED", []fsm.Action{SendMail})
	return &OrderFSM{m}
}

func main() {
	order := New("123abc", "RECEIVED")
	machine := NewOrderFSM(order)
	fmt.Printf("BEFORE: %v\n", order)
	_ = machine.ExecTrigger("VALIDATED")
	fmt.Printf("AFTER: %v\n", order)
}
