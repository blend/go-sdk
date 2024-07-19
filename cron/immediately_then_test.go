/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func Test_Immediately(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	ts := time.Date(2019, 9, 8, 12, 11, 10, 9, time.UTC)

	is := Immediately()
	its.Equal(StringScheduleImmediately, is.String())
	next := is.Next(ts)
	its.NotEqual(ts, next)
	next = is.Next(ts)
	its.True(next.IsZero())

	thenSchedule := EverySecond()
	is = Immediately().Then(thenSchedule).(*ImmediateSchedule)
	its.Equal(fmt.Sprintf("%s %v", StringScheduleImmediatelyThen, thenSchedule), is.String())

	next = is.Next(ts)
	its.NotEqual(ts, next)
	its.Equal(1, is.didRun)

	next = is.Next(ts)
	its.Equal(ts.Add(time.Second), next)

	// another one to be safe.
	next = is.Next(ts)
	its.Equal(ts.Add(time.Second), next)
	its.Equal(1, is.didRun)
}

func Test_Immediately_new(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	ts := time.Date(2019, 9, 8, 12, 11, 10, 9, time.UTC)

	is := new(ImmediateSchedule)
	its.Equal(StringScheduleImmediately, is.String())
	next := is.Next(ts)
	its.NotEqual(ts, next)
	next = is.Next(ts)
	its.True(next.IsZero())

	thenSchedule := EverySecond()
	is = new(ImmediateSchedule).Then(thenSchedule).(*ImmediateSchedule)
	its.Equal(fmt.Sprintf("%s %v", StringScheduleImmediatelyThen, thenSchedule), is.String())

	next = is.Next(ts)
	its.NotEqual(ts, next)
	its.Equal(1, is.didRun)

	next = is.Next(ts)
	its.Equal(ts.Add(time.Second), next)

	// another one to be safe.
	next = is.Next(ts)
	its.Equal(ts.Add(time.Second), next)
	its.Equal(1, is.didRun)
}

func Test_Immediately_Then(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	s := Immediately().Then(EveryHour())
	its.NotNil(s.Next(Zero))
	now := Now()
	next := s.Next(Now())
	its.True(next.Sub(now) > time.Minute, fmt.Sprintf("%v", next.Sub(now)))
	its.True(next.Sub(now) < (2 * time.Hour))
}

func Test_ImmediateSchedule_parallel(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	ts := time.Date(2019, 9, 8, 12, 11, 10, 9, time.UTC)

	var nextCount int32
	next := ScheduleFunc(func(ts time.Time) time.Time {
		atomic.AddInt32(&nextCount, 1)
		return ts.Add(time.Minute)
	})

	is := Immediately().Then(next)

	start := make(chan struct{})

	now := Now()
	times := make(chan time.Time, 10)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for x := 0; x < 10; x++ {
		go func() {
			defer wg.Done()
			<-start
			out := is.Next(ts)
			times <- out
		}()
	}

	close(start)
	wg.Wait()

	var allTimes []time.Time
	for x := 0; x < 10; x++ {
		allTimes = append(allTimes, <-times)
	}

	its.AnyCount(allTimes, 1, func(v interface{}) bool {
		typed, _ := v.(time.Time)
		return typed.After(now) && typed.Sub(now) < time.Second
	}, "(1) of the times should be '''now''' and the other 9 should have gone through the next schedule")
	its.AnyCount(allTimes, 9, func(v interface{}) bool {
		return v.(time.Time).Equal(ts.Add(time.Minute))
	}, "(1) of the times should be '''now''' and the other 9 should have gone through the next schedule")
}
