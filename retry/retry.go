// Package retry provides functions for executing potentially failing operations with configurable retry policies.
//
// It offers generic functions to execute an operation multiple times with a
// specified backoff period until it succeeds, the maximum number of retries is
// exceeded, or the provided context is done (canceled or timed out).
//
// The generic retry functions allow calling any function that returns a value
// of type T and an error, making them highly versatile.
package retry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrMaxAttemptsExceeded is returned when the function execution has failed all attempts.
var ErrMaxAttemptsExceeded = errors.New("retry: max attempts exceeded")

// ErrNegativeRetries is returned if the retries parameter is negative.
var ErrNegativeRetries = errors.New("retry: retries cannot be negative")

func execute[T any](ctx context.Context, fn func() (T, error), retries int, initialDelay time.Duration, backoffFactor int) (T, error) {
	var zeroT T

	if err := ctx.Err(); err != nil {
		return zeroT, err
	}

	if retries < 0 {
		return zeroT, ErrNegativeRetries
	}

	maxAttempts := 1 + retries
	var lastErr error
	currDelay := initialDelay
	for attempt := range maxAttempts {
		output, err := fn()
		if err == nil {
			return output, nil
		}
		lastErr = err
		if attempt < retries {
			select {
			case <-time.After(currDelay):
			case <-ctx.Done():
				return zeroT, ctx.Err()
			}
			currDelay *= time.Duration(backoffFactor)
		}
	}

	return zeroT, fmt.Errorf("%w (total %d attempts): last error: %w",
		ErrMaxAttemptsExceeded, maxAttempts, lastErr)
}

// Linearly executes the function fn up to 1 + retries times, waiting for a fixed
// initialDelay between attempts. The delay does not increase (backoff factor is 1).
//
// The total number of attempts is 1 (initial execution) + retries.
//
// Parameters:
//
//	ctx: The context used to control cancellation. If the context is canceled
//	  or timed out, the retry process stops immediately.
//	fn: The function to execute. It should return the desired result and an error.
//	  If the error is nil, the retry process stops.
//	retries: The maximum number of times to retry execution (must be >= 0).
//	initialDelay: The constant duration to wait between retries.
//
// Returns:
//
//	The result of fn on success, or an error if all attempts fail or the context is canceled.
//	If all attempts fail, ErrMaxAttemptsExceeded is returned, wrapped with the last error.
func Linearly[T any](ctx context.Context, fn func() (T, error), retries int, initialDelay time.Duration) (T, error) {
	return execute(ctx, fn, retries, initialDelay, 1)
}

// Exponentially executes the function fn up to 1 + retries times, using an
// exponential backoff strategy for delays. The delay between attempts doubles
// with each failure (backoff factor of 2).
//
// The total number of attempts is 1 (initial execution) + retries.
//
// Parameters:
//
//	ctx: The context used to control cancellation. If the context is canceled
//	  or timed out, the retry process stops immediately.
//	fn: The function to execute. It should return the desired result and an error.
//	  If the error is nil, the retry process stops.
//	retries: The maximum number of times to retry execution (must be >= 0).
//	initialDelay: The duration to wait before the first retry. Subsequent delays
//	  are multiplied by 2.
//
// Returns:
//
//	The result of fn on success, or an error if all attempts fail or the context is canceled.
//	If all attempts fail, ErrMaxAttemptsExceeded is returned, wrapped with the last error.
func Exponentially[T any](ctx context.Context, fn func() (T, error), retries int, initialDelay time.Duration) (T, error) {
	return execute(ctx, fn, retries, initialDelay, 2)
}
