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

// Interface assertions.
var (
	_	Schedule	= (*ImmediateSchedule)(nil)
	_	fmt.Stringer	= (*ImmediateSchedule)(nil)
)

// Immediately Returns a schedule that causes a job to run immediately on start,
// with an optional subsequent schedule.
func Immediately() *ImmediateSchedule {
	return &ImmediateSchedule{}
}

// ImmediateSchedule fires immediately with an optional continuation schedule.
type ImmediateSchedule struct {
	didRun	int32
	then	Schedule
}

// String returns a string representation of the schedul.e
func (i *ImmediateSchedule) String() string {
	if i.then != nil {
		return fmt.Sprintf("%s %v", StringScheduleImmediatelyThen, i.then)
	}
	return StringScheduleImmediately
}

// Then allows you to specify a subsequent schedule after the first run.
func (i *ImmediateSchedule) Then(then Schedule) Schedule {
	i.then = then
	return i
}

// Next implements Schedule.
func (i *ImmediateSchedule) Next(after time.Time) time.Time {
	if atomic.CompareAndSwapInt32(&i.didRun, 0, 1) {
		return Now()
	}
	if i.then != nil {
		return i.then.Next(after)
	}
	return Zero
}
