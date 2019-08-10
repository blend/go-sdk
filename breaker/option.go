package breaker

import (
	"context"
	"time"
)

// Option is a mutator for a breaker.
type Option func(*Breaker) error

// OptHalfOpenMaxActions sets the HalfOpenMaxActions.
func OptHalfOpenMaxActions(maxActions int64) Option {
	return func(b *Breaker) error {
		b.HalfOpenMaxActions = maxActions
		return nil
	}
}

// OptClosedExpiryInterval sets the ClosedExpiryInterval.
func OptClosedExpiryInterval(interval time.Duration) Option {
	return func(b *Breaker) error {
		b.ClosedExpiryInterval = interval
		return nil
	}
}

// OptOpenExpiryInterval sets the OpenExpiryInterval.
func OptOpenExpiryInterval(interval time.Duration) Option {
	return func(b *Breaker) error {
		b.OpenExpiryInterval = interval
		return nil
	}
}

// OptConfig sets the breaker based on a config.
func OptConfig(cfg Config) Option {
	return func(b *Breaker) error {
		b.HalfOpenMaxActions = cfg.HalfOpenMaxActions
		b.ClosedExpiryInterval = cfg.ClosedExpiryInterval
		b.OpenExpiryInterval = cfg.OpenExpiryInterval
		return nil
	}
}

// OptOpenAction sets the open action on the breaker.
func OptOpenAction(action func(context.Context)) Option {
	return func(b *Breaker) error {
		b.OpenAction = action
		return nil
	}
}

// OptOnStateChange sets the OnFaiilure handler on the breaker.
func OptOnStateChange(handler func(ctx context.Context, from, to State, generation int64)) Option {
	return func(b *Breaker) error {
		b.OnStateChange = handler
		return nil
	}
}

// OptShouldCloseProvider sets the ShouldCloseProvider provider on the breaker.
func OptShouldCloseProvider(provider func(ctx context.Context, counts Counts) bool) Option {
	return func(b *Breaker) error {
		b.ShouldCloseProvider = provider
		return nil
	}
}

// OptNowProvider sets the now provider on the breaker.
func OptNowProvider(provider func() time.Time) Option {
	return func(b *Breaker) error {
		b.NowProvider = provider
		return nil
	}
}
