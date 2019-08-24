package httpmetrics

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/webutil"
)

func TestAddListeners(t *testing.T) {
	assert := assert.New(t)

	log := logger.None()
	AddListeners(nil, nil)
	assert.False(log.HasListener(webutil.HTTPRequest, stats.ListenerNameStats))
	assert.False(log.HasListener(webutil.HTTPResponse, stats.ListenerNameStats))
	AddListeners(log, stats.NewMockCollector())
	assert.True(log.HasListener(webutil.HTTPRequest, stats.ListenerNameStats))
	assert.True(log.HasListener(webutil.HTTPResponse, stats.ListenerNameStats))
}
