package stats

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/blend/go-sdk/assert"
)

var runtimeMetrics = []string{
	"go.runtime.mem.num_gc",
	"go.runtime.mem.num_forced_gc",
	"go.runtime.mem.pause_total_ns",
	"go.runtime.mem.frees",
	"go.runtime.mem.mallocs",
	"go.runtime.num_cpu",
	"go.runtime.num_goroutine",
	"go.runtime.mem.alloc",
	"go.runtime.mem.gc_sys",
	"go.runtime.mem.other_sys",
	"go.runtime.mem.heap_alloc",
	"go.runtime.mem.heap_idle",
	"go.runtime.mem.heap_inuse",
	"go.runtime.mem.heap_objects",
	"go.runtime.mem.heap_sys",
	"go.runtime.mem.stack_inuse",
	"go.runtime.mem.stack_sys",
	"go.runtime.mem.sys",
	"go.runtime.mem.total_alloc",
}

func TestRuntimeCollect(t *testing.T) {
	assert := assert.New(t)

	var previous, current runtime.MemStats
	runtime.ReadMemStats(&previous)

	collector := NewMockCollector()
	for i := 0; i < len(runtimeMetrics); i++ {
		go func() { collector.Errors <- fmt.Errorf("error") }()
	}
	go runtimeCollect(collector, &previous, &current)

	for _, metricName := range runtimeMetrics {
		metric := <-collector.Events
		assert.Equal(metricName, metric.Name)
		assert.Zero(metric.Count)
		assert.Zero(metric.Histogram)
		assert.Zero(metric.TimeInMilliseconds)
	}
}
