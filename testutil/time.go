/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package testutil

import (
	"time"
)

type assertions interface {
	Nil(interface{}, ...interface{}) bool
}

// NowRounded returns the current time, but rounded to a given precision and
// then placed into a timezone given by a location name from the IANA Time Zone
// database (or "Local").
//
// This is useful in situations where timestamps are written to and then read
// back from a foreign system, like a database. For example, a `TIMESTAMP WITH TIME ZONE`
// column in `postgres` will truncate to microsecond precision and will return
// a "bare" timezone even if the timezone written was UTC.
func NowRounded(it assertions, locationName string, precision time.Duration) time.Time {
	loc, err := time.LoadLocation(locationName)
	it.Nil(err)
	// Round to the nearest `precision` (e.g. microsecond) to ensure accuracy
	// across Go / PostgreSQL boundaries and across different platforms.
	return time.Now().UTC().Truncate(precision).In(loc)
}
