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

	assert.NotNil(NewJob("test_job", noop))
	assert.True(NewJob("test_job", noop).Enabled())
	assert.Zero(NewJob("test_job", noop).Timeout())
	assert.False(NewJob("test_job", noop).ShouldSkipLoggerListeners())
	assert.False(NewJob("test_job", noop).ShouldSkipLoggerOutput())
	assert.Equal("test_job", NewJob("test_job", noop).Name())
}
