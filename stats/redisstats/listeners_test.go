/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package redisstats

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/redis"
	"github.com/blend/go-sdk/stats"
)

func TestAddListeners(t *testing.T) {
	assert := assert.New(t)

	log := logger.None()
	AddListeners(nil, nil)
	assert.False(log.HasListener(redis.Flag, stats.ListenerNameStats))
	AddListeners(log, stats.NewMockCollector(32))
	assert.True(log.HasListener(redis.Flag, stats.ListenerNameStats))
}

func TestAddListenersStats(t *testing.T) {
	assert := assert.New(t)

	log := logger.All(logger.OptOutput(io.Discard))
	defer log.Close()
	collector := stats.NewMockCollector(32)

	AddListeners(log, collector)

	log.TriggerContext(context.Background(), redis.NewEvent("GET", []string{"foo"}, 250*time.Millisecond,
		redis.OptEventNetwork(redis.DefaultNetwork),
		redis.OptEventAddr(redis.DefaultAddr),
		redis.OptEventAuthUser("system"),
	))

	m := <-collector.Metrics
	assert.Equal(MetricName, m.Name)
	assert.Equal(1, m.Count)
	assert.NotEmpty(m.Tags)

	m = <-collector.Metrics
	assert.Equal(MetricNameElapsedLast, m.Name)
	assert.Equal(250, m.Gauge)
	assert.NotEmpty(m.Tags)

	m = <-collector.Metrics
	assert.Equal(MetricNameElapsed, m.Name)
	assert.Equal(250, m.Histogram)
	assert.NotEmpty(m.Tags)
}
