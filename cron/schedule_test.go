package cron

import (
	"fmt"
	"testing"
	"time"
	"github.com/blend/go-sdk/assert"
)

func TestDailyScheduleEveryDay(t *testing.T) {
	a := assert.New(t)
	schedule := DailyAtUTC(12, 0, 0) //noon
	now := time.Now().UTC()
	beforenoon := time.Date(now.Year(), now.Month(), now.Day(), 11, 0, 0, 0, time.UTC)
	afternoon := time.Date(now.Year(), now.Month(), now.Day(), 13, 0, 0, 0, time.UTC)
	todayAtNoon := schedule.Next(beforenoon)
	tomorrowAtNoon := schedule.Next(afternoon)

	a.True(todayAtNoon.Before(afternoon))
	a.True(tomorrowAtNoon.After(afternoon))
}

func TestDailyScheduleSingleDay(t *testing.T) {
	a := assert.New(t)
	schedule := WeeklyAtUTC(12, 0, 0, time.Monday)               //every monday at noon
	beforenoon := time.Date(2016, 01, 11, 11, 0, 0, 0, time.UTC) //these are both a monday
	afternoon := time.Date(2016, 01, 11, 13, 0, 0, 0, time.UTC)  //these are both a monday

	sundayBeforeNoon := time.Date(2016, 01, 17, 11, 0, 0, 0, time.UTC) //to gut check that it's monday

	todayAtNoon := schedule.Next(beforenoon)
	nextWeekAtNoon := schedule.Next(afternoon)

	a.NonFatal().True(todayAtNoon.Before(afternoon))
	a.NonFatal().True(nextWeekAtNoon.After(afternoon))
	a.NonFatal().True(nextWeekAtNoon.After(sundayBeforeNoon))
	a.NonFatal().Equal(time.Monday, nextWeekAtNoon.Weekday())
}

func TestDayOfWeekFunctions(t *testing.T) {
	assert := assert.New(t)

	for _, wd := range WeekendDays {
		assert.True(IsWeekendDay(wd))
		assert.False(IsWeekDay(wd))
	}

	for _, wd := range WeekDays {
		assert.False(IsWeekendDay(wd))
		assert.True(IsWeekDay(wd))
	}
}

func TestOnTheHourAt(t *testing.T) {
	assert := assert.New(t)

	now := time.Now().UTC()
	schedule := EveryHourAtUTC(40, 00)

	fromNil := schedule.Next(Zero)
	assert.NotNil(fromNil)

	fromNilExpected := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 40, 0, 0, time.UTC)
	if fromNilExpected.Before(now) {
		fromNilExpected = fromNilExpected.Add(time.Hour)
	}
	assert.InTimeDelta(fromNilExpected, fromNil, time.Second)

	fromHalfStart := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 45, 0, 0, time.UTC)
	fromHalfExpected := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 40, 0, 0, time.UTC)

	fromHalf := schedule.Next(fromHalfStart)

	assert.NotNil(fromHalf)
	assert.InTimeDelta(fromHalfExpected, fromHalf, time.Second)
}

func TestImmediatelyThen(t *testing.T) {
	assert := assert.New(t)

	s := Immediately().Then(EveryHour())
	assert.NotNil(s.Next(Zero))
	now := Now()
	next := s.Next(Now())
	assert.True(next.Sub(now) > time.Minute, fmt.Sprintf("%v", next.Sub(now)))
	assert.True(next.Sub(now) < (2 * time.Hour))
}

func TestOnceAt(t *testing.T) {
	assert := assert.New(t)

	fireAt := time.Date(2018, 10, 21, 12, 00, 00, 00, time.UTC)
	before := fireAt.Add(-time.Minute)
	after := fireAt.Add(time.Minute)

	s := OnceAtUTC(fireAt)
	result := s.Next(before)
	assert.Equal(result, fireAt)

	result = s.Next(after)
	assert.True(result.IsZero())
}
