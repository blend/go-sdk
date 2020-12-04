package dbstats

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

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
	AddListeners(log, stats.NewMockCollector(32))
	assert.True(log.HasListener(db.QueryFlag, stats.ListenerNameStats))
}

func TestAddListenersStats(t *testing.T) {
	assert := assert.New(t)

	log := logger.All(logger.OptOutput(ioutil.Discard))
	collector := stats.NewMockCollector(32)

	AddListeners(log, collector)

	log.TriggerContext(context.Background(), db.NewQueryEvent("select 'ok!'", time.Second))

	qm := <-collector.Metrics
	assert.Equal(MetricNameDBQuery, qm.Name)
	assert.Equal(1, qm.Count)
	assert.NotEmpty(qm.Tags)

	qm = <-collector.Metrics
	assert.Equal(MetricNameDBQueryElapsed, qm.Name)
	assert.Equal(1000, qm.Gauge)
	assert.NotEmpty(qm.Tags)

	qm = <-collector.Metrics
	assert.Equal(MetricNameDBQueryElapsed, qm.Name)
	assert.Equal(1000, qm.TimeInMilliseconds)
	assert.NotEmpty(qm.Tags)
}
