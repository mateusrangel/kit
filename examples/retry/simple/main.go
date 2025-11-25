package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mateusrangel/kit/examples/internal/domain"
	"github.com/mateusrangel/kit/retry"
)

func main() {
	pg := &domain.BrokenPaymentGateway{}
	// Run within avaiable time
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := retry.Execute(ctx, func() (*domain.PaymentGatewayOutput, error) {
		return pg.ProcessTransaction(ctx, &domain.PaymentGatewayInput{CardToken: "1234", Amount: "99.00"})
	}, 2, 2*time.Second)
	if err != nil {
		fmt.Println(err)
	}

	//timeouts
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = retry.Execute(ctx, func() (*domain.PaymentGatewayOutput, error) {
		return pg.ProcessTransaction(ctx, &domain.PaymentGatewayInput{CardToken: "1234", Amount: "99.00"})
	}, 2, 2*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
}
