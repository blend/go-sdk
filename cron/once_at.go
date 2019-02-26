package cron

import (
	"fmt"
	"time"
)

// Interface assertions.
var (
	_ Schedule     = (*OnceAtUTCSchedule)(nil)
	_ fmt.Stringer = (*OnceAtUTCSchedule)(nil)
)

// OnceAtUTC returns a schedule that fires once at a given time.
// It will never fire again unless reloaded.
func OnceAtUTC(t time.Time) Schedule {
	return OnceAtUTCSchedule{Time: t}
}

// OnceAtUTCSchedule is a schedule.
type OnceAtUTCSchedule struct {
	Time time.Time
}

// String returns a string representation of the schedule.
func (oa OnceAtUTCSchedule) String() string {
	return fmt.Sprintf("once at %s", oa.Time.Format(time.RFC3339))
}

// Next returns the next runtime.
func (oa OnceAtUTCSchedule) Next(after time.Time) time.Time {
	if after.IsZero() {
		return oa.Time
	}
	if oa.Time.After(after) {
		return oa.Time
	}
	return Zero
}
