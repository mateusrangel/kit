package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mateusrangel/kit/retry"
)

var functionThatAlwaysWorks = func() (any, error) { return true, nil }

var functionThatWorksOnTheThirdTry = func() func() (any, error) {
	var callCount int

	return func() (any, error) {
		callCount++
		if callCount < 3 {
			return nil, errors.New("simulated initial failure")
		}

		return "worked on third try", nil
	}
}

func TestLinearly(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		fn      func() (any, error)
		retries int
		backoff time.Duration
		want    any
		wantErr bool
	}{
		{
			name:    "Function that always works",
			ctx:     context.Background(),
			fn:      functionThatAlwaysWorks,
			retries: 1,
			backoff: 1 * time.Nanosecond,
			want:    true,
			wantErr: false,
		},
		{
			name:    "Function that never works",
			ctx:     context.Background(),
			fn:      func() (any, error) { return nil, errors.New("error") },
			retries: 1,
			backoff: 1 * time.Nanosecond,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Function that works on the trird try",
			ctx:     context.Background(),
			fn:      functionThatWorksOnTheThirdTry(),
			retries: 2,
			backoff: 1 * time.Nanosecond,
			want:    "worked on third try",
			wantErr: false,
		},
		{
			name:    "ErrNegativeRetries",
			ctx:     context.Background(),
			fn:      functionThatAlwaysWorks,
			retries: -1,
			backoff: 1 * time.Nanosecond,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "max attempts exceeded",
			ctx:     context.Background(),
			fn:      functionThatWorksOnTheThirdTry(),
			retries: 1,
			backoff: 1 * time.Nanosecond,
			want:    nil,
			wantErr: true,
		},
		{
			name: "context was already timed out",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				return ctx
			}(),
			fn:      functionThatAlwaysWorks,
			retries: 2,
			backoff: 1 * time.Hour,
			want:    nil,
			wantErr: true,
		},
		{
			name: "context time out mid retry",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
				return ctx
			}(),
			fn:      functionThatWorksOnTheThirdTry(),
			retries: 2,
			backoff: 1 * time.Hour,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := retry.Linearly(tt.ctx, tt.fn, tt.retries, tt.backoff)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Linearly() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Linearly() succeeded unexpectedly")
			}

			if got != tt.want {
				t.Errorf("Linearly() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExponentially(t *testing.T) {

	tests := []struct {
		name    string
		ctx     context.Context
		fn      func() (any, error)
		retries int
		backoff time.Duration
		want    any
		wantErr bool
	}{
		{
			name:    "Function that always works",
			ctx:     context.Background(),
			fn:      functionThatAlwaysWorks,
			retries: 1,
			backoff: 1 * time.Nanosecond,
			want:    true,
			wantErr: false,
		},
		{
			name:    "Function that never works",
			ctx:     context.Background(),
			fn:      func() (any, error) { return nil, errors.New("error") },
			retries: 1,
			backoff: 1 * time.Nanosecond,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Function that works on the trird try",
			ctx:     context.Background(),
			fn:      functionThatWorksOnTheThirdTry(),
			retries: 2,
			backoff: 1 * time.Nanosecond,
			want:    "worked on third try",
			wantErr: false,
		},
		{
			name:    "ErrNegativeRetries",
			ctx:     context.Background(),
			fn:      functionThatAlwaysWorks,
			retries: -1,
			backoff: 1 * time.Nanosecond,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "max attempts exceeded",
			ctx:     context.Background(),
			fn:      functionThatWorksOnTheThirdTry(),
			retries: 1,
			backoff: 1 * time.Nanosecond,
			want:    nil,
			wantErr: true,
		},
		{
			name: "context was already timed out",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				return ctx
			}(),
			fn:      functionThatAlwaysWorks,
			retries: 2,
			backoff: 1 * time.Hour,
			want:    nil,
			wantErr: true,
		},
		{
			name: "context time out mid retry",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
				return ctx
			}(),
			fn:      functionThatWorksOnTheThirdTry(),
			retries: 2,
			backoff: 1 * time.Hour,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := retry.Exponentially(tt.ctx, tt.fn, tt.retries, tt.backoff)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Exponentially() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Exponentially() succeeded unexpectedly")
			}

			if got != tt.want {
				t.Errorf("Exponentially() = %v, want %v", got, tt.want)
			}
		})
	}
}
