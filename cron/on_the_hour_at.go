/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import (
	"fmt"
	"time"
)

var (
	_ Schedule     = (*OnTheHourAtUTCSchedule)(nil)
	_ fmt.Stringer = (*OnTheHourAtUTCSchedule)(nil)
)

// EveryHourOnTheHour returns a schedule that fires every 60 minutes on the 00th minute.
func EveryHourOnTheHour() Schedule {
	return OnTheHourAtUTCSchedule{}
}

// EveryHourAtUTC returns a schedule that fires every hour at a given minute.
func EveryHourAtUTC(minute, second int) Schedule {
	return OnTheHourAtUTCSchedule{Minute: minute, Second: second}
}

// OnTheHourAtUTCSchedule is a schedule that fires every hour on the given minute.
type OnTheHourAtUTCSchedule struct {
	Minute int
	Second int
}

// String returns a string representation of the schedule.
func (o OnTheHourAtUTCSchedule) String() string {
	return fmt.Sprintf("on the hour at %v:%v", o.Minute, o.Second)
}

// Next implements the chronometer Schedule api.
func (o OnTheHourAtUTCSchedule) Next(after time.Time) time.Time {
	var returnValue time.Time
	now := Now()
	if after.IsZero() {
		returnValue = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), o.Minute, o.Second, 0, time.UTC)
		if returnValue.Before(now) {
			returnValue = returnValue.Add(time.Hour)
		}
	} else {
		returnValue = time.Date(after.Year(), after.Month(), after.Day(), after.Hour(), o.Minute, o.Second, 0, time.UTC)
		if returnValue.Before(after) {
			returnValue = returnValue.Add(time.Hour)
		}
	}
	return returnValue
}
