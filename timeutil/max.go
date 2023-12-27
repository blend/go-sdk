/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package timeutil

import "time"

// Max returns the earliest (min) time in a list of times.
func Max(times ...time.Time) time.Time {
	if len(times) == 0 {
		return time.Time{}
	}

	end := times[0]
	for _, t := range times[1:] {
		if t.After(end) {
			end = t
		}
	}
	return end
}
