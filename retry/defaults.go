/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package retry

// Defaults applies defaults to a given options.
func Defaults(r *Retrier) {
	r.MaxAttempts = DefaultMaxAttempts
	r.DelayProvider = ConstantDelay(DefaultRetryDelay)
	r.ShouldRetryProvider = func(_ error) bool { return true }
}
