package dbmetrics

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stats"
)

func TestAddListeners(t *testing.T) {
	assert := assert.New(t)

	log := logger.None()
	AddListeners(nil, nil)
	assert.False(log.HasListener(db.QueryFlag, stats.ListenerNameStats))
	AddListeners(log, stats.NewMockCollector())
	assert.True(log.HasListener(db.QueryFlag, stats.ListenerNameStats))
}
