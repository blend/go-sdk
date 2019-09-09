package cron

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestIntervalSchedule(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(time.Second, EverySecond().Every)
	assert.Equal(time.Minute, EveryMinute().Every)
	assert.Equal(time.Hour, EveryHour().Every)
	assert.Equal(time.Millisecond, Every(time.Millisecond).Every)

	schedule := EveryHour()
	assert.Equal("every 1h0m0s", schedule.String())

	now := time.Now().UTC()
	firstRun := schedule.Next(Zero)
	firstRunDiff := firstRun.Sub(now)
	assert.InDelta(float64(firstRunDiff), float64(1*time.Hour), float64(1*time.Second))
	next := schedule.Next(now)
	assert.True(next.After(now))

	delay := EveryDelayed(time.Hour, time.Second)
	assert.Equal("every 1h0m0s with an initial delay of 1s", delay.String())

	now = time.Now().UTC()
	firstRun = delay.Next(Zero)
	firstRunDiff = firstRun.Sub(now)
	assert.InDelta(float64(firstRunDiff), float64(1*time.Hour), float64(2*time.Second))
	next = schedule.Next(now)
	assert.True(next.After(now))
}
