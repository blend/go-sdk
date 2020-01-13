package cron

import (
	"testing"
	"time"
	"github.com/blend/go-sdk/assert"
)

func TestImmediateSchedule(t *testing.T) {
	assert := assert.New(t)

	ts := time.Date(2019, 9, 8, 12, 11, 10, 9, time.UTC)

	is := new(ImmediateSchedule)
	assert.Equal("immediately, once", is.String())
	next := is.Next(ts)
	assert.NotEqual(ts, next)
	next = is.Next(ts)
	assert.True(next.IsZero())

	is = new(ImmediateSchedule).Then(EverySecond()).(*ImmediateSchedule)
	assert.Equal("immediately, then every 1s", is.String())

	next = is.Next(ts)
	assert.NotEqual(ts, next)

	next = is.Next(ts)
	assert.Equal(ts.Add(time.Second), next)
}
