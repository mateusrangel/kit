package domain

import (
	"context"
	"errors"
	"time"
)

type PaymentGatewayInput struct {
	CardToken string
	Amount    string
}

type PaymentGatewayOutput struct {
	Tid string
}

type PaymentGateway interface {
	ProcessTransaction(ctx context.Context, input *PaymentGatewayInput) (*PaymentGatewayOutput, error)
}

type BrokenPaymentGateway struct{}

func (g *BrokenPaymentGateway) ProcessTransaction(ctx context.Context, input *PaymentGatewayInput) (*PaymentGatewayOutput, error) {
	select {
	// simulates a 1 second round trip request that returned an error
	case <-time.After(1 * time.Second):
		return nil, errors.New("BrokenPaymentGateway: internal server error")
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

type WorkingPaymentGateway struct{}

func (g *WorkingPaymentGateway) ProcessTransaction(ctx context.Context, input *PaymentGatewayInput) (*PaymentGatewayOutput, error) {
	// always works
	return &PaymentGatewayOutput{Tid: "123"}, nil
}
