package cron

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/ref"

	"github.com/blend/go-sdk/assert"
)

func TestJobConfig(t *testing.T) {
	assert := assert.New(t)

	var jc JobConfig
	assert.Equal(DefaultTimeout, jc.TimeoutOrDefault())
	assert.Equal(DefaultShutdownGracePeriod, jc.ShutdownGracePeriodOrDefault())
	assert.Equal(DefaultHistoryDisabled, jc.HistoryDisabledOrDefault())
	assert.Equal(0, jc.HistoryMaxCountOrDefault())
	assert.Equal(0, jc.HistoryMaxAgeOrDefault())

	jc.Timeout = time.Second
	jc.ShutdownGracePeriod = time.Minute
	jc.HistoryDisabled = ref.Bool(true)
	jc.HistoryMaxCount = 5
	jc.HistoryMaxAge = time.Hour

	assert.Equal(jc.Timeout, jc.TimeoutOrDefault())
	assert.Equal(jc.ShutdownGracePeriod, jc.ShutdownGracePeriodOrDefault())
	assert.Equal(*jc.HistoryDisabled, jc.HistoryDisabledOrDefault())
	assert.Equal(jc.HistoryMaxCount, jc.HistoryMaxCountOrDefault())
	assert.Equal(jc.HistoryMaxAge, jc.HistoryMaxAgeOrDefault())
}
