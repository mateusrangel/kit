package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mateusrangel/kit/examples/internal/domain"
	"github.com/mateusrangel/kit/retry"
)

type PaymentProcessor struct {
	next *PaymentProcessor
	pg   domain.PaymentGateway
}

func NewPaymentProcessor(pg domain.PaymentGateway, next *PaymentProcessor) *PaymentProcessor {
	return &PaymentProcessor{next: next, pg: pg}
}

func (p *PaymentProcessor) ProcessPayment(ctx context.Context, input *domain.PaymentGatewayInput) (*domain.PaymentGatewayOutput, error) {
	output, err := retry.Execute(
		ctx,
		func() (*domain.PaymentGatewayOutput, error) {
			return p.pg.ProcessTransaction(ctx, &domain.PaymentGatewayInput{CardToken: "1234", Amount: "99.00"})
		},
		2,
		2*time.Second,
	)
	if err != nil {
		nextOutput, err := p.next.ProcessPayment(ctx, input)
		if err != nil {
			return nil, err
		}
		return nextOutput, nil
	}
	return output, nil
}

func main() {
	spp := NewPaymentProcessor(&domain.WorkingPaymentGateway{}, nil)
	fpp := NewPaymentProcessor(&domain.BrokenPaymentGateway{}, spp)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := fpp.ProcessPayment(ctx, &domain.PaymentGatewayInput{CardToken: "1234", Amount: "99.00"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(out)
}
