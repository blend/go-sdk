/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func Test_DailyAtUTC(t *testing.T) {
	t.Parallel()
	its := assert.New(t)
	schedule := DailyAtUTC(12, 0, 0)	//noon
	now := time.Now().UTC()
	beforenoon := time.Date(now.Year(), now.Month(), now.Day(), 11, 0, 0, 0, time.UTC)
	afternoon := time.Date(now.Year(), now.Month(), now.Day(), 13, 0, 0, 0, time.UTC)
	todayAtNoon := schedule.Next(beforenoon)
	tomorrowAtNoon := schedule.Next(afternoon)

	its.True(todayAtNoon.Before(afternoon))
	its.True(tomorrowAtNoon.After(afternoon))
}

func Test_WeeklyAtUTC(t *testing.T) {
	t.Parallel()
	its := assert.New(t)
	schedule := WeeklyAtUTC(12, 0, 0, time.Monday)			//every monday at noon
	beforenoon := time.Date(2016, 01, 11, 11, 0, 0, 0, time.UTC)	//these are both a monday
	afternoon := time.Date(2016, 01, 11, 13, 0, 0, 0, time.UTC)	//these are both a monday

	sundayBeforeNoon := time.Date(2016, 01, 17, 11, 0, 0, 0, time.UTC)	//to gut check that it's monday

	todayAtNoon := schedule.Next(beforenoon)
	nextWeekAtNoon := schedule.Next(afternoon)

	its.NonFatal().True(todayAtNoon.Before(afternoon))
	its.NonFatal().True(nextWeekAtNoon.After(afternoon))
	its.NonFatal().True(nextWeekAtNoon.After(sundayBeforeNoon))
	its.NonFatal().Equal(time.Monday, nextWeekAtNoon.Weekday())
}

func Test_OnTheHourAt(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	now := time.Now().UTC()
	schedule := EveryHourAtUTC(40, 00)

	fromNil := schedule.Next(Zero)
	its.NotNil(fromNil)

	fromNilExpected := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 40, 0, 0, time.UTC)
	if fromNilExpected.Before(now) {
		fromNilExpected = fromNilExpected.Add(time.Hour)
	}
	its.InTimeDelta(fromNilExpected, fromNil, time.Second)

	fromHalfStart := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 45, 0, 0, time.UTC)
	fromHalfExpected := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 40, 0, 0, time.UTC)

	fromHalf := schedule.Next(fromHalfStart)

	its.NotNil(fromHalf)
	its.InTimeDelta(fromHalfExpected, fromHalf, time.Second)
}
