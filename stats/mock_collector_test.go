package stats

import (
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

	collector := NewMockCollector()
	collector.AddDefaultTag("k1", "v1")
	collector.AddDefaultTag("k2", "v2")

	tags := collector.DefaultTags()
	assert.Len(tags, 2)
	assert.Equal("k1:v1", tags[0])
	assert.Equal("k2:v2", tags[1])

	go collector.Count("event", 10, "k3:v3")
	mockMetric := <-collector.Events
	assertTags(mockMetric.Tags)

	go collector.Increment("event", "k3:v3")
	mockMetric = <-collector.Events
	assertTags(mockMetric.Tags)

	go collector.Gauge("event", 0.1, "k3:v3")
	mockMetric = <-collector.Events
	assertTags(mockMetric.Tags)

	go collector.Histogram("event", 0.1, "k3:v3")
	mockMetric = <-collector.Events
	assertTags(mockMetric.Tags)

	go collector.TimeInMilliseconds("event", time.Second, "k3:v3")
	mockMetric = <-collector.Events
	assertTags(mockMetric.Tags)
}

func TestMockCollectorCount(t *testing.T) {
	assert := assert.New(t)

	collector := NewMockCollector()

	go collector.Count("event", 10)

	mockMetric := <-collector.Events
	assert.Equal("event", mockMetric.Name)
	assert.Equal(10, mockMetric.Count)
	assert.Zero(mockMetric.Gauge)
	assert.Zero(mockMetric.Histogram)
	assert.Zero(mockMetric.TimeInMilliseconds)
}

func TestMockCollectorIncrement(t *testing.T) {
	assert := assert.New(t)

	collector := NewMockCollector()

	go collector.Increment("event", "")

	mockMetric := <-collector.Events
	assert.Equal("event", mockMetric.Name)
	assert.Equal(1, mockMetric.Count)
	assert.Zero(mockMetric.Gauge)
	assert.Zero(mockMetric.Histogram)
	assert.Zero(mockMetric.TimeInMilliseconds)
}

func TestMockCollectorGauge(t *testing.T) {
	assert := assert.New(t)

	collector := NewMockCollector()

	go collector.Gauge("event", 0.1)

	mockMetric := <-collector.Events
	assert.Equal("event", mockMetric.Name)
	assert.Equal(0.1, mockMetric.Gauge)
	assert.Zero(mockMetric.Count)
	assert.Zero(mockMetric.Histogram)
	assert.Zero(mockMetric.TimeInMilliseconds)
}

func TestMockCollectorHistogram(t *testing.T) {
	assert := assert.New(t)

	collector := NewMockCollector()

	go collector.Histogram("event", 0.1)

	mockMetric := <-collector.Events
	assert.Equal("event", mockMetric.Name)
	assert.Equal(0.1, mockMetric.Histogram)
	assert.Zero(mockMetric.Gauge)
	assert.Zero(mockMetric.Count)
	assert.Zero(mockMetric.TimeInMilliseconds)
}

func TestMockCollectorTimeInMilliseconds(t *testing.T) {
	assert := assert.New(t)

	collector := NewMockCollector()

	go collector.TimeInMilliseconds("event", time.Second)

	mockMetric := <-collector.Events
	assert.Equal("event", mockMetric.Name)
	assert.Equal(1000, mockMetric.TimeInMilliseconds)
	assert.Zero(mockMetric.Gauge)
	assert.Zero(mockMetric.Count)
	assert.Zero(mockMetric.Histogram)
}
