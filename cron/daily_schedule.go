package cron

import (
	"fmt"
	"strings"
	"time"
)

var (
	_ Schedule     = (*DailySchedule)(nil)
	_ fmt.Stringer = (*DailySchedule)(nil)
)

// WeeklyAtUTC returns a schedule that fires on every of the given days at the given time by hour, minute and second in UTC.
func WeeklyAtUTC(hour, minute, second int, days ...time.Weekday) Schedule {
	dayOfWeekMask := uint(0)
	for _, day := range days {
		dayOfWeekMask = dayOfWeekMask | 1<<uint(day)
	}

	return &DailySchedule{DayOfWeekMask: dayOfWeekMask, TimeOfDayUTC: time.Date(0, 0, 0, hour, minute, second, 0, time.UTC)}
}

// DailyAtUTC returns a schedule that fires every day at the given hour, minute and second in UTC.
func DailyAtUTC(hour, minute, second int) Schedule {
	return &DailySchedule{DayOfWeekMask: AllDaysMask, TimeOfDayUTC: time.Date(0, 0, 0, hour, minute, second, 0, time.UTC)}
}

// WeekdaysAtUTC returns a schedule that fires every week day at the given hour, minute and second in UTC>
func WeekdaysAtUTC(hour, minute, second int) Schedule {
	return &DailySchedule{DayOfWeekMask: WeekDaysMask, TimeOfDayUTC: time.Date(0, 0, 0, hour, minute, second, 0, time.UTC)}
}

// WeekendsAtUTC returns a schedule that fires every weekend day at the given hour, minut and second.
func WeekendsAtUTC(hour, minute, second int) Schedule {
	return &DailySchedule{DayOfWeekMask: WeekendDaysMask, TimeOfDayUTC: time.Date(0, 0, 0, hour, minute, second, 0, time.UTC)}
}

// DailySchedule is a schedule that fires every day that satisfies the DayOfWeekMask at the given TimeOfDayUTC.
type DailySchedule struct {
	DayOfWeekMask uint
	TimeOfDayUTC  time.Time
}

func (ds DailySchedule) String() string {
	if ds.DayOfWeekMask > 0 {
		var days []string
		for _, d := range DaysOfWeek {
			if ds.checkDayOfWeekMask(d) {
				days = append(days, d.String())
			}
		}
		return fmt.Sprintf("%s on %s each week", ds.TimeOfDayUTC.Format(time.RFC3339), strings.Join(days, ", "))
	}
	return fmt.Sprintf("%s every day", ds.TimeOfDayUTC.Format(time.RFC3339))
}

func (ds DailySchedule) checkDayOfWeekMask(day time.Weekday) bool {
	trialDayMask := uint(1 << uint(day))
	bitwiseResult := (ds.DayOfWeekMask & trialDayMask)
	return bitwiseResult > uint(0)
}

// Next implements Schedule.
func (ds DailySchedule) Next(after time.Time) time.Time {
	if after.IsZero() {
		after = Now()
	}

	todayInstance := time.Date(after.Year(), after.Month(), after.Day(), ds.TimeOfDayUTC.Hour(), ds.TimeOfDayUTC.Minute(), ds.TimeOfDayUTC.Second(), 0, time.UTC)
	for day := 0; day < 8; day++ {
		next := todayInstance.AddDate(0, 0, day) //the first run here it should be adding nothing, i.e. returning todayInstance ...

		if ds.checkDayOfWeekMask(next.Weekday()) && next.After(after) { //we're on a day ...
			return next
		}
	}

	return Zero
}
