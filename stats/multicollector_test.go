/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

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

	c1 := NewMockCollector(32)
	c2 := NewMockCollector(32)

	mc := MultiCollector{c1, c2}

	err := mc.Count("event", 1, "k1:v1")
	assert.Nil(err)
	metric1 := <-c1.Metrics
	metric2 := <-c2.Metrics
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
	metric1 = <-c1.Metrics
	assert.Zero(metric1.Gauge)
	assert.Zero(metric1.Histogram)
	assert.Zero(metric1.TimeInMilliseconds)
}

func TestIncrement(t *testing.T) {
	assert := assert.New(t)

	c1 := NewMockCollector(32)
	c2 := NewMockCollector(32)

	var err error
	mc := MultiCollector{c1, c2}
	err = mc.Increment("event", "k1:v1")
	assert.Nil(err)

	metric1 := <-c1.Metrics
	metric2 := <-c2.Metrics
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
	metric1 = <-c1.Metrics
	assert.Zero(metric1.Gauge)
	assert.Zero(metric1.Histogram)
	assert.Zero(metric1.TimeInMilliseconds)
}

func TestGauge(t *testing.T) {
	assert := assert.New(t)
	c1 := NewMockCollector(32)
	c2 := NewMockCollector(32)

	var err error
	mc := MultiCollector{c1, c2}
	err = mc.Gauge("event", .01)
	assert.Nil(err)

	metric1 := <-c1.Metrics
	metric2 := <-c2.Metrics
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
	metric1 = <-c1.Metrics
	assert.Zero(metric1.Count)
	assert.Zero(metric1.Histogram)
	assert.Zero(metric1.TimeInMilliseconds)
}

func TestHistogram(t *testing.T) {
	assert := assert.New(t)
	c1 := NewMockCollector(32)
	c2 := NewMockCollector(32)

	var err error
	mc := MultiCollector{c1, c2}
	err = mc.Histogram("event", .01)
	assert.Nil(err)

	metric1 := <-c1.Metrics
	metric2 := <-c2.Metrics
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
	metric1 = <-c1.Metrics
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

	c1 := NewMockCollector(32)
	c2 := NewMockCollector(32)

	var err error
	mc := MultiCollector{c1, c2}
	err = mc.TimeInMilliseconds("event", time.Second, "k1:v1")
	assert.Nil(err)
	metric1 := <-c1.Metrics
	metric2 := <-c2.Metrics
	assert.Equal("event", metric1.Name)
	assert.Equal(1000, metric1.TimeInMilliseconds)
	assertTags(metric1.Tags)
	assert.Equal(metric1, metric2)

	mc = MultiCollector{c1, c2}
	c1.Errors <- fmt.Errorf("error")
	err = mc.TimeInMilliseconds("event", time.Second, "k1:v1")
	assert.NotNil(err)
	assert.Equal("error", err.Error())
	metric1 = <-c1.Metrics
	assert.Zero(metric1.Gauge)
	assert.Zero(metric1.Histogram)
	assert.Zero(metric1.Count)
}

func TestFlush(t *testing.T) {
	assert := assert.New(t)

	c1 := NewMockCollector(32)
	c2 := NewMockCollector(32)

	var err error
	mc := MultiCollector{c1, c2}

	err = mc.Flush()
	assert.Nil(err)

	expectedError := fmt.Errorf("err")
	c2.FlushErrors <- expectedError
	err = mc.Flush()
	assert.Equal(expectedError.Error(), err.Error())
}

func TestClose(t *testing.T) {
	assert := assert.New(t)

	c1 := NewMockCollector(32)
	c2 := NewMockCollector(32)

	var err error
	mc := MultiCollector{c1, c2}

	err = mc.Close()
	assert.Nil(err)

	expectedError := fmt.Errorf("err")
	c2.CloseErrors <- expectedError
	err = mc.Close()
	assert.Equal(expectedError.Error(), err.Error())
}

func TestSendEvent(t *testing.T) {
	assert := assert.New(t)

	c1 := NewMockCollector(32)
	c2 := NewMockCollector(32)

	var err error
	mc := MultiCollector{c1, c2}

	err = mc.SendEvent(Event{})
	assert.Nil(err)
}

func TestCreateEvent(t *testing.T) {
	assert := assert.New(t)

	c1 := NewMockCollector(32)
	c2 := NewMockCollector(32)

	mc := MultiCollector{c1, c2}

	event := mc.CreateEvent("title", "text")
	assert.Equal(event.Title, "title")
	assert.Equal(event.Text, "text")
}
