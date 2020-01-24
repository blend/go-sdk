package cron

import (
	"fmt"
	"sync"
	"time"
)

// Interface assertions.
var (
	_ Schedule     = (*ImmediateSchedule)(nil)
	_ fmt.Stringer = (*ImmediateSchedule)(nil)
)

// Immediately Returns a schedule that casues a job to run immediately on start,
// with an optional subsequent schedule.
func Immediately() *ImmediateSchedule {
	return &ImmediateSchedule{}
}

// ImmediateSchedule fires immediately with an optional continuation schedule.
type ImmediateSchedule struct {
	sync.Mutex
	didRun bool
	then   Schedule
}

// String returns a string representation of the schedul.e
func (i *ImmediateSchedule) String() string {
	if i.then != nil {
		return fmt.Sprintf("immediately, then %v", i.then)
	}
	return "immediately, once"
}

// Then allows you to specify a subsequent schedule after the first run.
func (i *ImmediateSchedule) Then(then Schedule) Schedule {
	i.then = then
	return i
}

// Next implements Schedule.
func (i *ImmediateSchedule) Next(after time.Time) time.Time {
	i.Lock()
	defer i.Unlock()

	if !i.didRun {
		i.didRun = true
		return Now()
	}
	if i.then != nil {
		return i.then.Next(after)
	}
	return Zero
}
