/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cron

import "time"

// this is used by `Now()`
// it can be overrided in tests etc.
var _nowProvider = time.Now

// Now returns a new timestamp.
func Now() time.Time {
	return _nowProvider().UTC()
}

// Since returns the duration since another timestamp.
func Since(t time.Time) time.Duration {
	return Now().Sub(t)
}

// Min returns the minimum of two times.
func Min(t1, t2 time.Time) time.Time {
	if t1.IsZero() && t2.IsZero() {
		return time.Time{}
	}
	if t1.IsZero() {
		return t2
	}
	if t2.IsZero() {
		return t1
	}
	if t1.Before(t2) {
		return t1
	}
	return t2
}

// Max returns the maximum of two times.
func Max(t1, t2 time.Time) time.Time {
	if t1.Before(t2) {
		return t2
	}
	return t1
}

// FormatTime returns a string for a time.
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// NOTE: time.Zero()? what's that?
var (
	// DaysOfWeek are all the time.Weekday in an array for utility purposes.
	DaysOfWeek	= []time.Weekday{
		time.Sunday,
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
		time.Saturday,
	}

	// WeekDays are the business time.Weekday in an array.
	WeekDays	= []time.Weekday{
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
	}

	// WeekWeekEndDaysDays are the weekend time.Weekday in an array.
	WeekendDays	= []time.Weekday{
		time.Sunday,
		time.Saturday,
	}

	// Epoch is unix epoch saved for utility purposes.
	Epoch	= time.Unix(0, 0)
	// Zero is different than epoch in that it is the "unset" value for a time
	// where Epoch is a valid date. Nominally it is `time.Time{}`.
	Zero	= time.Time{}
)

// NOTE: we have to use shifts here because in their infinite wisdom google didn't make these values powers of two for masking
const (
	// AllDaysMask is a bitmask of all the days of the week.
	AllDaysMask	= 1<<uint(time.Sunday) | 1<<uint(time.Monday) | 1<<uint(time.Tuesday) | 1<<uint(time.Wednesday) | 1<<uint(time.Thursday) | 1<<uint(time.Friday) | 1<<uint(time.Saturday)
	// WeekDaysMask is a bitmask of all the weekdays of the week.
	WeekDaysMask	= 1<<uint(time.Monday) | 1<<uint(time.Tuesday) | 1<<uint(time.Wednesday) | 1<<uint(time.Thursday) | 1<<uint(time.Friday)
	//WeekendDaysMask is a bitmask of the weekend days of the week.
	WeekendDaysMask	= 1<<uint(time.Sunday) | 1<<uint(time.Saturday)
)

// IsWeekDay returns if the day is a monday->friday.
func IsWeekDay(day time.Weekday) bool {
	return !IsWeekendDay(day)
}

// IsWeekendDay returns if the day is a monday->friday.
func IsWeekendDay(day time.Weekday) bool {
	return day == time.Saturday || day == time.Sunday
}

// ConstBool returns a value provider for a constant.
func ConstBool(value bool) func() bool {
	return func() bool {
		return value
	}
}

// ConstInt returns a value provider for a constant.
func ConstInt(value int) func() int {
	return func() int {
		return value
	}
}

// ConstDuration returns a value provider for a constant.
func ConstDuration(value time.Duration) func() time.Duration {
	return func() time.Duration {
		return value
	}
}

// ConstLabels returns a value provider for a constant.
func ConstLabels(labels map[string]string) func() map[string]string {
	return func() map[string]string {
		return labels
	}
}
