/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package breaker

import (
	"time"
)

// Option is a mutator for a breaker.
type Option func(*Breaker)

// OptOpenFailureThreshold sets the OpenFailureThreshold.
func OptOpenFailureThreshold(openFailureThreshold int64) Option {
	return func(b *Breaker) {
		b.OpenFailureThreshold = openFailureThreshold
	}
}

// OptHalfOpenMaxActions sets the HalfOpenMaxActions.
func OptHalfOpenMaxActions(maxActions int64) Option {
	return func(b *Breaker) {
		b.HalfOpenMaxActions = maxActions
	}
}

// OptClosedExpiryInterval sets the ClosedExpiryInterval.
func OptClosedExpiryInterval(interval time.Duration) Option {
	return func(b *Breaker) {
		b.ClosedExpiryInterval = interval
	}
}

// OptOpenExpiryInterval sets the OpenExpiryInterval.
func OptOpenExpiryInterval(interval time.Duration) Option {
	return func(b *Breaker) {
		b.OpenExpiryInterval = interval
	}
}

// OptConfig sets the breaker based on a config.
func OptConfig(cfg Config) Option {
	return func(b *Breaker) {
		b.HalfOpenMaxActions = cfg.HalfOpenMaxActions
		b.ClosedExpiryInterval = cfg.ClosedExpiryInterval
		b.OpenExpiryInterval = cfg.OpenExpiryInterval
	}
}

// OptOpenAction sets the open action on the breaker.
//
// The "Open" action is called when the breaker opens,
// that is, when it no longer allows calls.
func OptOpenAction(action Actioner) Option {
	return func(b *Breaker) {
		b.OpenAction = action
	}
}

// OptOnStateChange sets the OnStateChange handler on the breaker.
func OptOnStateChange(handler OnStateChangeHandler) Option {
	return func(b *Breaker) {
		b.OnStateChange = handler
	}
}

// OptShouldOpenProvider sets the ShouldCloseProvider provider on the breaker.
func OptShouldOpenProvider(provider ShouldOpenProvider) Option {
	return func(b *Breaker) {
		b.ShouldOpenProvider = provider
	}
}

// OptNowProvider sets the now provider on the breaker.
func OptNowProvider(provider NowProvider) Option {
	return func(b *Breaker) {
		b.NowProvider = provider
	}
}
