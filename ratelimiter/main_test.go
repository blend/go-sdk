/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ratelimiter

import (
	"time"
)

// Clock is a helper function.
func Clock(t time.Time, offset time.Duration) func() time.Time {
	return func() time.Time {
		return t.Add(offset)
	}
}
