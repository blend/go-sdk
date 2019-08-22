package stats

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

func TestAddWebListeners(t *testing.T) {
	assert := assert.New(t)

	log := logger.None()
	AddWebListeners(nil, nil)
	assert.False(log.HasListener(webutil.HTTPResponse, ListenerNameStats))
	AddWebListeners(log, NewMockCollector())
	assert.True(log.HasListener(webutil.HTTPResponse, ListenerNameStats))
}

func TestAddQueryListeners(t *testing.T) {
	assert := assert.New(t)

	log := logger.None()
	AddQueryListeners(nil, nil)
	assert.False(log.HasListener(db.QueryFlag, ListenerNameStats))
	AddQueryListeners(log, NewMockCollector())
	assert.True(log.HasListener(db.QueryFlag, ListenerNameStats))
}

func TestAddErrorListeners(t *testing.T) {
	assert := assert.New(t)

	log := logger.None()
	AddErrorListeners(nil, nil)
	assert.False(log.HasListener(logger.Warning, ListenerNameStats))
	assert.False(log.HasListener(logger.Error, ListenerNameStats))
	assert.False(log.HasListener(logger.Fatal, ListenerNameStats))
	AddErrorListeners(log, NewMockCollector())
	assert.True(log.HasListener(logger.Warning, ListenerNameStats))
	assert.True(log.HasListener(logger.Error, ListenerNameStats))
	assert.True(log.HasListener(logger.Fatal, ListenerNameStats))
}
