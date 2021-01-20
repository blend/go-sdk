package cron

import (
	"fmt"
	"sync"
	"time"
)

// Interface assertions.
var (
	_ Schedule     = (*TimesSchedule)(nil)
	_ fmt.Stringer = (*TimesSchedule)(nil)
)

// Times returns a new times schedule that returns a given
// next run time from a schedule only a certain number of times.
func Times(times int, schedule Schedule) *TimesSchedule {
	return &TimesSchedule{
		times:    times,
		left:     times,
		schedule: schedule,
	}
}

// TimesSchedule is a schedule that only returns
// a certain number of results.
type TimesSchedule struct {
	sync.Mutex

	times    int
	left     int
	schedule Schedule
}

// Next implements cron.Schedule.
func (ts *TimesSchedule) Next(after time.Time) time.Time {
	ts.Lock()
	defer ts.Unlock()

	if ts.left > 0 {
		ts.left--
		return ts.schedule.Next(after)
	}
	return Zero
}

// String returns a string representation of the schedul.e
func (ts *TimesSchedule) String() string {
	if typed, ok := ts.schedule.(fmt.Stringer); ok {
		return fmt.Sprintf("%s %d/%d times", typed.String(), ts.left, ts.times)
	}
	return fmt.Sprintf("%d/%d times", ts.left, ts.times)
}
