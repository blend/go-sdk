/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package timeutil

import "time"

// BeginningOfMonth returns the date that represents
// the last day of the month for a given time.
func BeginningOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 01, 00, 00, 00, 00, t.Location())	// move to YY-MM-01 00:00.00
}
