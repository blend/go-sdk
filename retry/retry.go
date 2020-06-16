package retry

import (
	"context"
	"time"

	"github.com/blend/go-sdk/ex"
)

// Options are the options for the retry.
type Options struct {
	MaxAttempts         int
	DelayProvider       DelayProvider
	ShouldRetryProvider ShouldRetryProvider
}

// Option mutates retry options.
type Option func(*Options)

// OptMaxAttempts sets the max attempts.
func OptMaxAttempts(maxAttempts int) Option {
	return func(ro *Options) { ro.MaxAttempts = maxAttempts }
}

// OptDelayProvider sets the retry delay provider.
func OptDelayProvider(delayProvider DelayProvider) Option {
	return func(ro *Options) { ro.DelayProvider = delayProvider }
}

// OptConstantDelay sets the retry delay provider.
func OptConstantDelay(d time.Duration) Option {
	return func(ro *Options) { ro.DelayProvider = ConstantDelay(d) }
}

// OptExponentialDelay sets the retry delay provider.
func OptExponentialDelay(d time.Duration) Option {
	return func(ro *Options) { ro.DelayProvider = ExponentialBackoff(d) }
}

// OptShouldRetryProvider sets the should retry provider.
func OptShouldRetryProvider(provider ShouldRetryProvider) Option {
	return func(ro *Options) { ro.ShouldRetryProvider = provider }
}

// Retry calls a function and
func Retry(ctx context.Context, action Action, opts ...Option) (res interface{}, err error) {
	options := Options{
		MaxAttempts:         DefaultMaxAttempts,
		DelayProvider:       ConstantDelay(DefaultRetryDelay),
		ShouldRetryProvider: func(_ error) bool { return true },
	}
	for _, opt := range opts {
		opt(&options)
	}

	for attempt := 0; attempt < options.MaxAttempts; attempt++ {
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
