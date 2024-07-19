/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func Test_IntervalSchedule(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	its.Equal(time.Second, EverySecond().Every)
	its.Equal(time.Minute, EveryMinute().Every)
	its.Equal(time.Hour, EveryHour().Every)
	its.Equal(time.Millisecond, Every(time.Millisecond).Every)

	schedule := EveryHour()
	its.Equal("@every 1h0m0s", schedule.String())

	now := time.Now().UTC()
	firstRun := schedule.Next(Zero)
	firstRunDiff := firstRun.Sub(now)
	its.InDelta(float64(firstRunDiff), float64(1*time.Hour), float64(1*time.Second))
	next := schedule.Next(now)
	its.True(next.After(now))
}
