/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stats

import (
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestMockCollectorDefaultTags(t *testing.T) {
	assert := assert.New(t)

	assertTags := func(actualTags []string) {
		assert.Len(actualTags, 3)
		assert.Equal("k1:v1", actualTags[0])
		assert.Equal("k2:v2", actualTags[1])
		assert.Equal("k3:v3", actualTags[2])
	}

	collector := NewMockCollector(32)
	collector.AddDefaultTags(Tag("k1", "v1"))
	collector.AddDefaultTags(Tag("k2", "v2"))

	tags := collector.DefaultTags()
	assert.Len(tags, 2)
	assert.Equal("k1:v1", tags[0])
	assert.Equal("k2:v2", tags[1])

	assert.Nil(collector.Count("event", 10, "k3:v3"))
	mockMetric := <-collector.Metrics
	assertTags(mockMetric.Tags)

	assert.Nil(collector.Increment("event", "k3:v3"))
	mockMetric = <-collector.Metrics
	assertTags(mockMetric.Tags)

	assert.Nil(collector.Gauge("event", 0.1, "k3:v3"))
	mockMetric = <-collector.Metrics
	assertTags(mockMetric.Tags)

	assert.Nil(collector.Histogram("event", 0.1, "k3:v3"))
	mockMetric = <-collector.Metrics
	assertTags(mockMetric.Tags)

	assert.Nil(collector.TimeInMilliseconds("event", time.Second, "k3:v3"))
	mockMetric = <-collector.Metrics
	assertTags(mockMetric.Tags)
}

func TestMockCollectorCount(t *testing.T) {
	assert := assert.New(t)

	collector := NewMockCollector(32)

	assert.Nil(collector.Count("event", 10))

	mockMetric := <-collector.Metrics
	assert.Equal("event", mockMetric.Name)
	assert.Equal(10, mockMetric.Count)
	assert.Zero(mockMetric.Gauge)
	assert.Zero(mockMetric.Histogram)
	assert.Zero(mockMetric.TimeInMilliseconds)
}

func TestMockCollectorIncrement(t *testing.T) {
	assert := assert.New(t)

	collector := NewMockCollector(32)

	assert.Nil(collector.Increment("event", ""))

	mockMetric := <-collector.Metrics
	assert.Equal("event", mockMetric.Name)
	assert.Equal(1, mockMetric.Count)
	assert.Zero(mockMetric.Gauge)
	assert.Zero(mockMetric.Histogram)
	assert.Zero(mockMetric.TimeInMilliseconds)
}

func TestMockCollectorGauge(t *testing.T) {
	assert := assert.New(t)

	collector := NewMockCollector(32)

	assert.Nil(collector.Gauge("event", 0.1))

	mockMetric := <-collector.Metrics
	assert.Equal("event", mockMetric.Name)
	assert.Equal(0.1, mockMetric.Gauge)
	assert.Zero(mockMetric.Count)
	assert.Zero(mockMetric.Histogram)
	assert.Zero(mockMetric.TimeInMilliseconds)
}

func TestMockCollectorHistogram(t *testing.T) {
	assert := assert.New(t)

	collector := NewMockCollector(32)

	assert.Nil(collector.Histogram("event", 0.1))

	mockMetric := <-collector.Metrics
	assert.Equal("event", mockMetric.Name)
	assert.Equal(0.1, mockMetric.Histogram)
	assert.Zero(mockMetric.Gauge)
	assert.Zero(mockMetric.Count)
	assert.Zero(mockMetric.TimeInMilliseconds)
}

func TestMockCollectorTimeInMilliseconds(t *testing.T) {
	assert := assert.New(t)

	collector := NewMockCollector(32)

	assert.Nil(collector.TimeInMilliseconds("event", time.Second))

	mockMetric := <-collector.Metrics
	assert.Equal("event", mockMetric.Name)
	assert.Equal(1000, mockMetric.TimeInMilliseconds)
	assert.Zero(mockMetric.Gauge)
	assert.Zero(mockMetric.Count)
	assert.Zero(mockMetric.Histogram)
}

func TestMockCollectorFlush(t *testing.T) {
	assert := assert.New(t)

	collector := NewMockCollector(32)

	err := collector.Flush()
	assert.Nil(err)

	expectedErr := fmt.Errorf("err")
	collector.FlushErrors <- expectedErr
	err = collector.Flush()
	assert.Equal(expectedErr.Error(), err.Error())
}

func TestMockCollectorClose(t *testing.T) {
	assert := assert.New(t)

	collector := NewMockCollector(32)

	err := collector.Close()
	assert.Nil(err)

	expectedErr := fmt.Errorf("err")
	collector.CloseErrors <- expectedErr
	err = collector.Close()
	assert.Equal(expectedErr.Error(), err.Error())
}
