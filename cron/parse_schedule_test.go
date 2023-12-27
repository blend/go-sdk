/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cron

import (
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/stringutil"
)

type stringScheduleTestCase struct {
	Input		string
	ExpectedErr	error
	ExpectedNow	bool
	Expected	time.Time
	After		time.Time
}

func TestParseSchedule(t *testing.T) {
	assert := assert.New(t)

	testCases := []stringScheduleTestCase{
		{Input: "", ExpectedErr: ErrStringScheduleInvalid},
		{Input: stringutil.Random(stringutil.Letters, 10), ExpectedErr: ErrStringScheduleInvalid},
		{Input: "*/1 * * * * * *", After: time.Date(2018, 12, 29, 13, 12, 11, 10, time.UTC), Expected: time.Date(2018, 12, 29, 13, 12, 12, 0, time.UTC)},
		{Input: "*/5 * * * * * *", After: time.Date(2018, 12, 29, 13, 12, 11, 0, time.UTC), Expected: time.Date(2018, 12, 29, 13, 12, 15, 0, time.UTC)},
		{Input: "* 2 1 * * 1-6 *", After: time.Date(2019, 01, 01, 12, 0, 0, 0, time.UTC), Expected: time.Date(2019, 01, 02, 01, 02, 0, 0, time.UTC)},
		{Input: "* 2 1 * * MON-FRI *", After: time.Date(2019, 01, 01, 12, 0, 0, 0, time.UTC), Expected: time.Date(2019, 01, 02, 01, 02, 0, 0, time.UTC)},
		{Input: "* 9 10 * * SUN-TUE *", After: time.Date(2019, 01, 02, 12, 0, 0, 0, time.UTC), Expected: time.Date(2019, 01, 06, 10, 9, 0, 0, time.UTC)},
		{Input: "0 0 0 * * 0 *", After: time.Date(2019, 01, 02, 12, 3, 4, 5, time.UTC), Expected: time.Date(2019, 01, 06, 0, 0, 0, 0, time.UTC)},							// every week at midnight sat/sun
		{Input: "0 0 0 * * 0 *", After: time.Date(2019, 01, 06, 0, 0, 0, 1, time.UTC), Expected: time.Date(2019, 01, 13, 0, 0, 0, 0, time.UTC)},							// every week at midnight sat/sun (on almost exactly the same time)
		{Input: "0 0 0 * * * *", After: time.Date(2019, 01, 02, 12, 3, 4, 5, time.UTC), Expected: time.Date(2019, 01, 03, 0, 0, 0, 0, time.UTC)},							// every day at midnight
		{Input: "0 0 * * * * *", After: time.Date(2019, 01, 02, 12, 3, 4, 5, time.UTC), Expected: time.Date(2019, 01, 02, 13, 0, 0, 0, time.UTC)},							// every hour on the hour
		{Input: "0 * * * *", After: time.Date(2019, 01, 02, 12, 3, 4, 5, time.UTC), Expected: time.Date(2019, 01, 02, 13, 0, 0, 0, time.UTC)},								// every hour on the hour (5 field)
		{Input: "0 0 * * * *", After: time.Date(2019, 01, 02, 12, 3, 4, 5, time.UTC), Expected: time.Date(2019, 01, 02, 13, 0, 0, 0, time.UTC)},							// every hour on the hour (6 field)
		{Input: "@never", After: time.Date(2019, 01, 02, 12, 3, 4, 5, time.UTC), Expected: time.Time{}},												// never shorthand
		{Input: "@daily", After: time.Date(2019, 01, 02, 12, 3, 4, 5, time.UTC), Expected: time.Date(2019, 01, 03, 0, 0, 0, 0, time.UTC)},								// daily shorthand
		{Input: "@hourly", After: time.Date(2019, 01, 02, 12, 3, 4, 5, time.UTC), Expected: time.Date(2019, 01, 02, 13, 0, 0, 0, time.UTC)},								// hourly shorthand
		{Input: "@every not-a-value", ExpectedErr: ErrStringScheduleInvalid},																// every
		{Input: "@every 500ms", After: time.Date(2019, 01, 02, 12, 3, 4, 0, time.UTC), Expected: time.Date(2019, 01, 02, 12, 3, 4, int(500*time.Millisecond), time.UTC)},				// every
		{Input: "@once-at not-a-value", ExpectedErr: ErrStringScheduleInvalid},																// every
		{Input: "@once-at 2019-01-02T13:14:15.555Z", After: time.Date(2019, 01, 02, 13, 14, 14, 0, time.UTC), Expected: time.Date(2019, 01, 02, 13, 14, 15, int(555*time.Millisecond), time.UTC)},	// every
		{Input: "@immediately", After: time.Date(2019, 01, 02, 12, 3, 4, 0, time.UTC), ExpectedNow: true},												// immediately then every
		{Input: "@immediately-then @every 500ms", After: time.Date(2019, 01, 02, 12, 3, 4, 0, time.UTC), ExpectedNow: true},										// immediately then every
	}

	for _, tc := range testCases {
		parsed, err := ParseSchedule(tc.Input)
		if tc.ExpectedErr != nil {
			assert.NotNil(err)
			assert.True(ex.Is(err, tc.ExpectedErr))
		} else if tc.ExpectedNow {
			assert.Nil(err)
			next := parsed.Next(tc.After)
			assert.InTimeDelta(Now(), next, time.Second, fmt.Sprintf("%s parsed as %v\nexpected to be near-ish to now", tc.Input, parsed))
		} else {
			assert.Nil(err)
			next := parsed.Next(tc.After)
			assert.Equal(tc.Expected, next, fmt.Sprintf("%s parsed as %v\n%v vs. %v", tc.Input, parsed, tc.Expected.Format(time.RFC3339), next.Format(time.RFC3339)))
		}
	}
}

func TestStringScheduleEvery(t *testing.T) {
	assert := assert.New(t)

	schedule, err := ParseSchedule("*/1 * * * * * *")
	assert.Nil(err)

	last := time.Date(2019, 01, 29, 0, 0, 0, 0, time.UTC)
	for x := 0; x < 60*60*6; x++ {
		last = schedule.Next(last)
	}
	assert.Equal(time.Date(2019, 01, 29, 6, 0, 0, 0, time.UTC), last)
}

func TestMapKeysToArray(t *testing.T) {
	assert := assert.New(t)

	assert.Equal([]int{1, 2, 3}, mapKeysToArray(map[int]bool{
		3:	true,
		1:	true,
		2:	true,
	}))
	assert.Empty(mapKeysToArray(nil))
	assert.Empty(mapKeysToArray(map[int]bool{}))
}

func Test_ParseString_immediately(t *testing.T) {
	its := assert.New(t)

	var parsed Schedule
	var err error
	var input string
	var now, next time.Time
	after := time.Date(2018, 12, 29, 13, 12, 11, 10, time.UTC)

	now = Now()
	input = "@immediately-then */1 * * * * *"
	parsed, err = ParseSchedule(input)
	its.Nil(err)
	next = parsed.Next(after)
	its.NotInTimeDelta(after, next, time.Second)	// should be now
	its.InTimeDelta(now, next, time.Second)		// should be now
	next = parsed.Next(after)			// should kick in real schedule
	its.InTimeDelta(time.Date(2018, 12, 29, 13, 12, 12, 10, time.UTC), next, time.Millisecond)

	input = "@immediately-then bogus"
	parsed, err = ParseSchedule(input)
	its.True(ex.Is(err, ErrStringScheduleInvalid))
	its.Nil(parsed)

	now = Now()
	input = "@immediately-then @every 500ms"
	parsed, err = ParseSchedule(input)
	its.Nil(err)
	next = parsed.Next(after)
	its.NotInTimeDelta(after, next, time.Second)	// should be now
	its.InTimeDelta(now, next, time.Second)		// should be now

	next = parsed.Next(after)	// should kick in real schedule
	its.InTimeDelta(time.Date(2018, 12, 29, 13, 12, 11, 10+int(500*time.Millisecond), time.UTC), next, time.Millisecond)
}
