/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package timeutil

import "time"

// MinMax returns the minimum and maximum times in a given range.
func MinMax(times ...time.Time) (min time.Time, max time.Time) {
	if len(times) == 0 {
		return
	}
	min = times[0]
	max = times[0]

	for index := 1; index < len(times); index++ {
		if times[index].Before(min) {
			min = times[index]
		}
		if times[index].After(max) {
			max = times[index]
		}
	}
	return
}
