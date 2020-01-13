package sentry

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestAddListeners(t *testing.T) {
	assert := assert.New(t)

	AddListeners(nil, Config{})

	log := logger.None()
	AddListeners(log, Config{})
	assert.False(log.HasListeners(logger.Error))
	assert.False(log.HasListeners(logger.Fatal))

	AddListeners(log, Config{DSN: "http://foo@example.org/1"})
	assert.True(log.HasListeners(logger.Error))
	assert.True(log.HasListeners(logger.Fatal))

	assert.True(log.HasListener(logger.Error, ListenerName))
	assert.True(log.HasListener(logger.Fatal, ListenerName))
}
