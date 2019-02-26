package cron

import (
	"fmt"
	"time"
)

var (
	_ Schedule     = (*IntervalSchedule)(nil)
	_ fmt.Stringer = (*IntervalSchedule)(nil)
)

// EverySecond returns a schedule that fires every second.
func EverySecond() Schedule {
	return IntervalSchedule{Every: 1 * time.Second}
}

// EveryMinute returns a schedule that fires every minute.
func EveryMinute() Schedule {
	return IntervalSchedule{Every: 1 * time.Minute}
}

// EveryHour returns a schedule that fire every hour.
func EveryHour() Schedule {
	return IntervalSchedule{Every: 1 * time.Hour}
}

// Every returns a schedule that fires every given interval.
func Every(interval time.Duration) Schedule {
	return IntervalSchedule{Every: interval}
}

// IntervalSchedule is as chedule that fires every given interval with an optional start delay.
type IntervalSchedule struct {
	Every      time.Duration
	StartDelay time.Duration
}

// String returns a string representation of the schedule.
func (i IntervalSchedule) String() string {
	if i.StartDelay > 0 {
		return fmt.Sprintf("every %v with an initial delay of %v", i.Every, i.StartDelay)
	}
	return fmt.Sprintf("every %v", i.Every)
}

// Next implements Schedule.
func (i IntervalSchedule) Next(after time.Time) time.Time {
	if after.IsZero() {
		if i.StartDelay > 0 {
			return Now().Add(i.StartDelay).Add(i.Every)
		}
		return Now().Add(i.Every)
	}
	return after.Add(i.Every)
}
