/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package timeutil

import "time"

// BeginningOfMonth returns the date that represents
// the last day of the month for a given time.
func BeginningOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 01, 00, 00, 00, 00, t.Location()) // move to YY-MM-01 00:00.00
}
