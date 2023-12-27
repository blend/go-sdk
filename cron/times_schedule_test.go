/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cron

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestTimesSchedule(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	ts := time.Date(2020, 01, 23, 12, 11, 10, 9, time.UTC)
	schedule := Times(5, Every(time.Millisecond))
	its.Equal(ts.Add(time.Millisecond), schedule.Next(ts))
	its.Equal(ts.Add(time.Millisecond), schedule.Next(ts))
	its.Equal(ts.Add(time.Millisecond), schedule.Next(ts))
	its.Equal(ts.Add(time.Millisecond), schedule.Next(ts))
	its.Equal(ts.Add(time.Millisecond), schedule.Next(ts))
	its.True(schedule.Next(ts).IsZero())
}
