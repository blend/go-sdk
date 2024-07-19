/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package retry

import (
	"time"
)

// Option mutates retry options.
type Option func(*Retrier)

// OptMaxAttempts sets the max attempts.
func OptMaxAttempts(maxAttempts uint) Option {
	return func(o *Retrier) { o.MaxAttempts = maxAttempts }
}

// OptDelayProvider sets the retry delay provider.
func OptDelayProvider(delayProvider DelayProvider) Option {
	return func(o *Retrier) { o.DelayProvider = delayProvider }
}

// OptConstantDelay sets the retry delay provider.
func OptConstantDelay(d time.Duration) Option {
	return func(o *Retrier) { o.DelayProvider = ConstantDelay(d) }
}

// OptExponentialBackoff sets the retry delay provider.
func OptExponentialBackoff(d time.Duration) Option {
	return func(o *Retrier) { o.DelayProvider = ExponentialBackoff(d) }
}

// OptShouldRetryProvider sets the should retry provider.
func OptShouldRetryProvider(provider ShouldRetryProvider) Option {
	return func(o *Retrier) { o.ShouldRetryProvider = provider }
}
