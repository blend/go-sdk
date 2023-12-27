/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package timeutil

import "time"

// EndOfMonth returns the date that represents
// the last day of the month for a given time.
func EndOfMonth(t time.Time) time.Time {
	t2 := t.AddDate(0, 1, 0)                                                // add a month
	t3 := time.Date(t2.Year(), t2.Month(), 1, 00, 00, 00, 00, t.Location()) // move to YY-MM-01 00:00.00
	t4 := t3.Add(-time.Nanosecond)                                          // subtract (1) nanosecond
	return t4
}
