/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cron

import (
	"fmt"
	"sync/atomic"
	"time"
)

// Delay returns a composite schedule that delays
// a given schedule by a given duration.
func Delay(d time.Duration, then Schedule) *DelaySchedule {
	return &DelaySchedule{
		delay:	d,
		then:	then,
	}
}

var (
	_	Schedule	= (*DelaySchedule)(nil)
	_	fmt.Stringer	= (*DelaySchedule)(nil)
)

// DelaySchedule wraps a schedule with a delay.
type DelaySchedule struct {
	didRun	int32
	delay	time.Duration
	then	Schedule
}

// Next implements Schedule.
func (ds *DelaySchedule) Next(after time.Time) time.Time {
	if atomic.CompareAndSwapInt32(&ds.didRun, 0, 1) {
		return ds.then.Next(after).Add(ds.delay)
	}
	return ds.then.Next(after)
}

// String implements a string schedule.
func (ds *DelaySchedule) String() string {
	return fmt.Sprintf("%s %v %v", StringScheduleDelay, ds.delay, ds.then)
}
