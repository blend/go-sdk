package retry

import (
	"context"
	"time"

	"github.com/blend/go-sdk/ex"
)

// Retry calls a function and retries on error or if a should retry provider
// is set, based on the should retry result.
func Retry(ctx context.Context, action Action, opts ...Option) (interface{}, error) {
	var options Options
	DefaultOptions(&options)
	for _, opt := range opts {
		opt(&options)
	}
	return options.Retry(ctx, action)
}

// Retry calls a function and retries on error or if a should retry provider
// is set, based on the should retry result.
func (options Options) Retry(ctx context.Context, action Action) (res interface{}, err error) {
	var attempt uint
	for attempt = 0; attempt < options.MaxAttempts; attempt++ {
		func() {
			defer func() {
				if r := recover(); r != nil {
					err = ex.New(r)
				}
			}()
			res, err = action(ctx)
		}()
		if err == nil {
			return
		}
		if !options.ShouldRetryProvider(err) {
			return
		}

		select {
		case <-ctx.Done():
			return nil, context.Canceled
		case <-time.After(options.DelayProvider(ctx, attempt)):
		}
	}
	return
}
