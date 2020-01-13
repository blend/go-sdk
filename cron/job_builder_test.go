package cron

import (
	"context"
	"testing"
	"github.com/blend/go-sdk/assert"
)

func noop(_ context.Context) error {
	return nil
}

func TestJobBuilder(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(NewJob())
	assert.False(NewJob().Disabled())
	assert.Zero(NewJob().Timeout())
	assert.False(NewJob().ShouldSkipLoggerListeners())
	assert.False(NewJob().ShouldSkipLoggerOutput())
	assert.Equal("test_job", NewJob(OptJobName("test_job")).Name())
}
