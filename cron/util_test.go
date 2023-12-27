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

func Test_IsWeekendDay_IsWeekDay(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	for _, wd := range WeekendDays {
		its.True(IsWeekendDay(wd))
		its.False(IsWeekDay(wd))
	}

	for _, wd := range WeekDays {
		its.False(IsWeekendDay(wd))
		its.True(IsWeekDay(wd))
	}
}

func Test_MinMax(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	a := time.Date(2018, 10, 21, 12, 0, 0, 0, time.UTC)
	b := time.Date(2018, 10, 20, 12, 0, 0, 0, time.UTC)

	its.True(Min(time.Time{}, time.Time{}).IsZero())
	its.Equal(a, Min(a, time.Time{}))
	its.Equal(b, Min(time.Time{}, b))
	its.Equal(b, Min(a, b))
	its.Equal(b, Min(b, a))

	its.Equal(a, Max(a, b))
	its.Equal(a, Max(b, a))
}

func Test_Const(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	its.Equal(true, ConstBool(true)())
	its.Equal(false, ConstBool(false)())

	its.Equal(123, ConstInt(123)())
	its.Equal(6*time.Hour, ConstDuration(6*time.Hour)())
	its.Equal("foo", ConstLabels(map[string]string{
		"bar": "buzz",
		"moo": "foo",
	})()["moo"])
}
