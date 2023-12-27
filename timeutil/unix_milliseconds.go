/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package timeutil

import "time"

// UnixMilliseconds returns the time in unix (seconds) format
// with a floating point remainder for subsecond fraction.
func UnixMilliseconds(t time.Time) float64 {
	nanosPerSecond := float64(time.Second) / float64(time.Nanosecond)
	return float64(t.UnixNano()) / nanosPerSecond
}
