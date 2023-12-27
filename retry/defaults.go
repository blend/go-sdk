/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package retry

// Defaults applies defaults to a given options.
func Defaults(r *Retrier) {
	r.MaxAttempts = DefaultMaxAttempts
	r.DelayProvider = ConstantDelay(DefaultRetryDelay)
	r.ShouldRetryProvider = func(_ error) bool { return true }
}
