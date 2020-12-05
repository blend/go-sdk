package retry

import (
	"time"
)

// DefaultOptions applies defaults to a given options.
func DefaultOptions(o *Options) {
	o.MaxAttempts = DefaultMaxAttempts
	o.DelayProvider = ConstantDelay(DefaultRetryDelay)
	o.ShouldRetryProvider = func(_ error) bool { return true }
}

// Options are the options for the retry.
type Options struct {
	MaxAttempts         uint
	DelayProvider       DelayProvider
	ShouldRetryProvider ShouldRetryProvider
}

// Option mutates retry options.
type Option func(*Options)

// OptMaxAttempts sets the max attempts.
func OptMaxAttempts(maxAttempts uint) Option {
	return func(o *Options) { o.MaxAttempts = maxAttempts }
}

// OptDelayProvider sets the retry delay provider.
func OptDelayProvider(delayProvider DelayProvider) Option {
	return func(o *Options) { o.DelayProvider = delayProvider }
}

// OptConstantDelay sets the retry delay provider.
func OptConstantDelay(d time.Duration) Option {
	return func(o *Options) { o.DelayProvider = ConstantDelay(d) }
}

// OptExponentialBackoff sets the retry delay provider.
func OptExponentialBackoff(d time.Duration) Option {
	return func(o *Options) { o.DelayProvider = ExponentialBackoff(d) }
}

// OptShouldRetryProvider sets the should retry provider.
func OptShouldRetryProvider(provider ShouldRetryProvider) Option {
	return func(o *Options) { o.ShouldRetryProvider = provider }
}
