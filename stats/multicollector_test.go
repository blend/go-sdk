package stats

import (
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestCount(t *testing.T) {
	assert := assert.New(t)

	assertTags := func(actualTags []string) {
		assert.Len(actualTags, 1)
		assert.Equal("k1:v1", actualTags[0])
	}

	c1 := NewMockCollector()
	c2 := NewMockCollector()

	mc := MultiCollector{c1, c2}

	err := mc.Count("event", 1, "k1:v1")
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

	mc = MultiCollector{c1}
	c1.Errors <- fmt.Errorf("error")
	err = mc.Count("event", 1, "k1:v1")
	assert.NotNil(err)
	assert.Equal("error", err.Error())
	metric1 = <-c1.Events
	assert.Zero(metric1.Gauge)
	assert.Zero(metric1.Histogram)
	assert.Zero(metric1.TimeInMilliseconds)
}

func TestIncrement(t *testing.T) {
	assert := assert.New(t)

	c1 := NewMockCollector()
	c2 := NewMockCollector()

	var err error
	mc := MultiCollector{c1, c2}
	err = mc.Increment("event", "k1:v1")
	assert.Nil(err)

	metric1 := <-c1.Events
	metric2 := <-c2.Events
	assert.Equal("event", metric1.Name)
	assert.Equal(1, metric1.Count)
	assert.Zero(metric1.Gauge)
	assert.Zero(metric1.Histogram)
	assert.Zero(metric1.TimeInMilliseconds)
	assert.Equal(metric1, metric2)

	mc = MultiCollector{c1}

	c1.Errors <- fmt.Errorf("error")
	err = mc.Increment("event", "k1:v1")
	assert.NotNil(err)
	assert.Equal("error", err.Error())
	metric1 = <-c1.Events
	assert.Zero(metric1.Gauge)
	assert.Zero(metric1.Histogram)
	assert.Zero(metric1.TimeInMilliseconds)
}

func TestGauge(t *testing.T) {
	assert := assert.New(t)
	c1 := NewMockCollector()
	c2 := NewMockCollector()

	var err error
	mc := MultiCollector{c1, c2}
	err = mc.Gauge("event", .01)
	assert.Nil(err)

	metric1 := <-c1.Events
	metric2 := <-c2.Events
	assert.Equal("event", metric1.Name)
	assert.Equal(.01, metric1.Gauge)
	assert.Zero(metric1.Count)
	assert.Zero(metric1.Histogram)
	assert.Zero(metric1.TimeInMilliseconds)
	assert.Equal(metric1, metric2)

	mc = MultiCollector{c1}

	c1.Errors <- fmt.Errorf("error")
	err = mc.Gauge("event", .01)
	assert.NotNil(err)
	assert.Equal("error", err.Error())
	metric1 = <-c1.Events
	assert.Zero(metric1.Count)
	assert.Zero(metric1.Histogram)
	assert.Zero(metric1.TimeInMilliseconds)
}

func TestHistogram(t *testing.T) {
	assert := assert.New(t)
	c1 := NewMockCollector()
	c2 := NewMockCollector()

	var err error
	mc := MultiCollector{c1, c2}
	err = mc.Histogram("event", .01)
	assert.Nil(err)

	metric1 := <-c1.Events
	metric2 := <-c2.Events
	assert.Equal("event", metric1.Name)
	assert.Equal(.01, metric1.Histogram)
	assert.Zero(metric1.Count)
	assert.Zero(metric1.Gauge)
	assert.Zero(metric1.TimeInMilliseconds)
	assert.Equal(metric1, metric2)

	mc = MultiCollector{c1, c2}

	c1.Errors <- fmt.Errorf("error")
	err = mc.Histogram("event", .01)
	assert.NotNil(err)
	assert.Equal("error", err.Error())
	metric1 = <-c1.Events
	assert.Zero(metric1.Count)
	assert.Zero(metric1.Gauge)
	assert.Zero(metric1.TimeInMilliseconds)
}

func TestTimeInMilliseconds(t *testing.T) {
	assert := assert.New(t)

	assertTags := func(actualTags []string) {
		assert.Len(actualTags, 1)
		assert.Equal("k1:v1", actualTags[0])
	}

	c1 := NewMockCollector()
	c2 := NewMockCollector()

	var err error
	mc := MultiCollector{c1, c2}
	err = mc.TimeInMilliseconds("event", time.Second, "k1:v1")
	assert.Nil(err)
	metric1 := <-c1.Events
	metric2 := <-c2.Events
	assert.Equal("event", metric1.Name)
	assert.Equal(1000, metric1.TimeInMilliseconds)
	assertTags(metric1.Tags)
	assert.Equal(metric1, metric2)

	mc = MultiCollector{c1, c2}
	c1.Errors <- fmt.Errorf("error")
	err = mc.TimeInMilliseconds("event", time.Second, "k1:v1")
	assert.NotNil(err)
	assert.Equal("error", err.Error())
	metric1 = <-c1.Events
	assert.Zero(metric1.Gauge)
	assert.Zero(metric1.Histogram)
	assert.Zero(metric1.Count)
}
