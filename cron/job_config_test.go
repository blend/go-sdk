package cron

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestJobConfig(t *testing.T) {
	assert := assert.New(t)

	var jc JobConfig
	assert.Equal(DefaultTimeout, jc.TimeoutOrDefault())
	assert.Equal(DefaultShutdownGracePeriod, jc.ShutdownGracePeriodOrDefault())

	jc.Timeout = time.Second
	jc.ShutdownGracePeriod = time.Minute

	assert.Equal(jc.Timeout, jc.TimeoutOrDefault())
	assert.Equal(jc.ShutdownGracePeriod, jc.ShutdownGracePeriodOrDefault())
}
