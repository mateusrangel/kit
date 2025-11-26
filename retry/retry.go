// Package retry provides a generic functions to execute a function multiple times
// with a specified backoff period until it succeeds, a maximum number of
// retries is exceeded, or the provided context is done (canceled or timed out).
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

// Execute attempts to run the function fn up to (1 + retries) times.
//
// If fn returns success (non-nil output and nil error), its output is returned immediately.
// If fn returns an error, the function pauses for the backoff duration before the next
// attempt, provided the maximum number of attempts has not been reached.
//
// The execution halts immediately if the context is already done (canceled/expired) upon entry, or if the context signals done during a backoff period, returning ctx.Err().
//
// If retries is negative, Execute returns ErrNegativeRetries.
// If all attempts fail, it returns an error wrapping ErrMaxAttemptsExceeded,
// including the total number of attempts and the last error encountered.
func Execute[T any](ctx context.Context, fn func() (T, error), retries int, backoff time.Duration) (T, error) {
	var zeroT T

	if err := ctx.Err(); err != nil {
		return zeroT, err
	}

	if retries < 0 {
		return zeroT, ErrNegativeRetries
	}

	maxAttempts := 1 + retries
	var lastErr error

	for attempt := range maxAttempts {
		output, err := fn()
		if err == nil {
			return output, nil
		}
		lastErr = err
		if attempt < retries {
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return zeroT, ctx.Err()
			}
		}
	}

	return zeroT, fmt.Errorf("%w (total %d attempts): last error: %w",
		ErrMaxAttemptsExceeded, maxAttempts, lastErr)
}
