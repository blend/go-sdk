package stats

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestNewMultiCollector(t *testing.T) {
	assert := assert.New(t)

	c1 := NewMockCollector()
	c2 := NewMockCollector()

	_, err := NewMultiCollector(c1, c2)
	assert.Nil(err)

	_, err = NewMultiCollector()
	assert.NotNil(err)
}

func TestCount(t *testing.T) {
	assert := assert.New(t)

	assertTags := func(actualTags []string) {
		assert.Len(actualTags, 1)
		assert.Equal("k1:v1", actualTags[0])
	}

	c1 := NewMockCollector()
	c2 := NewMockCollector()

	mc, err := NewMultiCollector(c1, c2)
	go mc.Count("event", 1, "k1:v1")

	assert.Nil(err)
	metric1 := <-c1.Events
	metric2 := <-c2.Events
	assert.Equal("event", metric1.Name)
	assert.Equal(1, metric1.Count)
	assertTags(metric1.Tags)
	assert.Zero(metric1.Gauge)
	assert.Zero(metric1.Histogram)
	assert.Zero(metric1.TimeInMilliseconds)
	assert.Equal(metric1, metric2)
}

func TestIncrement(t *testing.T) {
	assert := assert.New(t)

	c1 := NewMockCollector()
	c2 := NewMockCollector()

	mc, err := NewMultiCollector(c1, c2)
	go mc.Increment("event", "k1:v1")

	assert.Nil(err)
	metric1 := <-c1.Events
	metric2 := <-c2.Events
	assert.Equal("event", metric1.Name)
	assert.Equal(1, metric1.Count)
	assert.Zero(metric1.Gauge)
	assert.Zero(metric1.Histogram)
	assert.Zero(metric1.TimeInMilliseconds)
	assert.Equal(metric1, metric2)
}

func TestGauge(t *testing.T) {
	assert := assert.New(t)
	c1 := NewMockCollector()
	c2 := NewMockCollector()

	mc, err := NewMultiCollector(c1, c2)
	go mc.Gauge("event", .01)

	assert.Nil(err)
	metric1 := <-c1.Events
	metric2 := <-c2.Events
	assert.Equal("event", metric1.Name)
	assert.Equal(.01, metric1.Gauge)
	assert.Zero(metric1.Count)
	assert.Zero(metric1.Histogram)
	assert.Zero(metric1.TimeInMilliseconds)
	assert.Equal(metric1, metric2)
}

func TestHistogram(t *testing.T) {
	assert := assert.New(t)
	c1 := NewMockCollector()
	c2 := NewMockCollector()

	mc, err := NewMultiCollector(c1, c2)
	go mc.Histogram("event", .01)

	assert.Nil(err)
	metric1 := <-c1.Events
	metric2 := <-c2.Events
	assert.Equal("event", metric1.Name)
	assert.Equal(.01, metric1.Histogram)
	assert.Zero(metric1.Count)
	assert.Zero(metric1.Gauge)
	assert.Zero(metric1.TimeInMilliseconds)
	assert.Equal(metric1, metric2)
}

func TestTimeInMilliseconds(t *testing.T) {
	assert := assert.New(t)

	assertTags := func(actualTags []string) {
		assert.Len(actualTags, 1)
		assert.Equal("k1:v1", actualTags[0])
	}

	c1 := NewMockCollector()
	c2 := NewMockCollector()

	mc, err := NewMultiCollector(c1, c2)
	go mc.TimeInMilliseconds("event", time.Second, "k1:v1")

	assert.Nil(err)
	metric1 := <-c1.Events
	metric2 := <-c2.Events
	assert.Equal("event", metric1.Name)
	assert.Equal(1000, metric1.TimeInMilliseconds)
	assertTags(metric1.Tags)
	assert.Zero(metric1.Gauge)
	assert.Zero(metric1.Histogram)
	assert.Zero(metric1.Count)
	assert.Equal(metric1, metric2)
}
