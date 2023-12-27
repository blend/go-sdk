/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package autoflush

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/stats"
)

// Assert Buffer is graceful.
var (
	_ graceful.Graceful = (*Buffer)(nil)
)

func Test_Buffer_MaxLen(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(2)

	var processed int32
	handler := func(_ context.Context, objects []interface{}) error {
		defer wg.Done()
		atomic.AddInt32(&processed, int32(len(objects)))
		return nil
	}

	afb := New(handler,
		OptMaxLen(10),
		OptInterval(time.Hour),
	)

	go func() { _ = afb.Start() }()
	<-afb.NotifyStarted()
	defer func() { _ = afb.Stop() }()

	for x := 0; x < 20; x++ {
		afb.Add(context.TODO(), fmt.Sprintf("foo%d", x))
	}

	wg.Wait()
	assert.Equal(20, processed)
}

func Test_Buffer_Ticker(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(20)
	buffer := New(func(_ context.Context, objects []interface{}) error {
		for range objects {
			wg.Done()
		}
		return nil
	}, OptMaxLen(100), OptInterval(time.Millisecond))

	go func() { _ = buffer.Start() }()
	<-buffer.NotifyStarted()
	defer func() { _ = buffer.Stop() }()

	for x := 0; x < 20; x++ {
		buffer.Add(context.TODO(), fmt.Sprintf("foo%d", x))
	}
	wg.Wait()
	assert.True(true)
}

func Test_Buffer_Stop(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(20)
	handler := func(_ context.Context, objects []interface{}) error {
		for range objects {
			wg.Done()
		}
		return nil
	}

	buffer := New(handler,
		OptMaxLen(10),
		OptInterval(time.Hour),
		OptShutdownGracePeriod(5*time.Second),
	)

	go func() { _ = buffer.Start() }()
	<-buffer.NotifyStarted()

	for x := 0; x < 20; x++ {
		buffer.Add(context.TODO(), fmt.Sprintf("foo%d", x))
	}

	assert.Nil(buffer.Stop())
	assert.True(buffer.Latch.IsStopped())
	wg.Wait()
}

func Test_Buffer_Stats(t *testing.T) {
	assert := assert.New(t)
	todo := context.TODO()

	mockCollector := stats.NewMockCollector(128)

	maxItems := 10
	numFlushes := 5

	var msg struct{}
	var flushBlock []interface{}
	for x := 0; x < maxItems; x++ {
		flushBlock = append(flushBlock, msg)
	}

	afb := New(nil,
		OptMaxLen(maxItems),
		OptStats(mockCollector),
	)
	assert.Equal(maxItems, afb.contents.Capacity())
	afb.flushes = make(chan Flush, maxItems)

	for x := 0; x < numFlushes; x++ {
		afb.AddMany(todo, flushBlock...) // should cause a queued flush
	}

	assert.Equal(numFlushes, len(afb.flushes))
	metrics := mockCollector.AllMetrics()

	assert.AnyCount(metrics, numFlushes, func(v interface{}) bool {
		typed := v.(stats.MockMetric)
		return typed.Name == MetricBufferLength
	}, MetricBufferLength)
	assert.AnyCount(metrics, numFlushes, func(v interface{}) bool {
		typed := v.(stats.MockMetric)
		return typed.Name == MetricFlush && typed.Count == 1
	}, MetricFlush)
	assert.AnyCount(metrics, numFlushes*3, func(v interface{}) bool {
		typed := v.(stats.MockMetric)
		return typed.Name == MetricFlushEnqueueElapsed
	}, MetricFlushEnqueueElapsed)

	assert.AnyCount(metrics, 1, func(v interface{}) bool {
		typed := v.(stats.MockMetric)
		return typed.Name == MetricFlushQueueLength && typed.Gauge == 0.0
	})
	assert.AnyCount(metrics, 1, func(v interface{}) bool {
		typed := v.(stats.MockMetric)
		return typed.Name == MetricFlushQueueLength && typed.Gauge == 1.0
	})
	assert.AnyCount(metrics, 1, func(v interface{}) bool {
		typed := v.(stats.MockMetric)
		return typed.Name == MetricFlushQueueLength && typed.Gauge == 2.0
	})
	assert.AnyCount(metrics, 1, func(v interface{}) bool {
		typed := v.(stats.MockMetric)
		return typed.Name == MetricFlushQueueLength && typed.Gauge == 3.0
	})
	assert.AnyCount(metrics, 1, func(v interface{}) bool {
		typed := v.(stats.MockMetric)
		return typed.Name == MetricFlushQueueLength && typed.Gauge == 4.0
	})
}

func BenchmarkBuffer(b *testing.B) {
	buffer := New(func(_ context.Context, objects []interface{}) error {
		if len(objects) > 128 {
			b.Fail()
		}
		return nil
	}, OptMaxLen(128), OptInterval(500*time.Millisecond))

	go func() { _ = buffer.Start() }()
	<-buffer.NotifyStarted()
	defer func() { _ = buffer.Stop() }()

	for x := 0; x < b.N; x++ {
		for y := 0; y < 1000; y++ {
			buffer.Add(context.TODO(), fmt.Sprintf("asdf%d%d", x, y))
		}
	}
}
