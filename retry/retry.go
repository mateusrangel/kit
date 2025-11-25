package retry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrMaxAttemptsExceeded error = errors.New("retry: max attempts exceeded")
var ErrNegativeRetries error = errors.New("retry: retries cannot be negative")

func Execute[T any](ctx context.Context, fn func() (T, error), retries int, backoff time.Duration) (T, error) {
	var zeroT T
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
