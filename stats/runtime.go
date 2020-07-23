package stats

import (
	"runtime"
	"time"
)

// Runtime reports golang vm runtime stats.
func Runtime(collector Collector) {
	if collector == nil {
		return
	}

	go func() {
		var previous, current runtime.MemStats
		runtimeCollect(collector, &previous, &current)
		for {
			<-time.After(250 * time.Millisecond)
			runtimeCollect(collector, &previous, &current)
		}
	}()
}

func runtimeCollect(collector Collector, previous, current *runtime.MemStats) {
	runtime.ReadMemStats(current)

	// these depend on the previous values
	_ = collector.Gauge("go.runtime.mem.num_gc", float64(current.NumGC-previous.NumGC))
	_ = collector.Gauge("go.runtime.mem.num_forced_gc", float64(current.NumForcedGC-previous.NumForcedGC))
	_ = collector.Gauge("go.runtime.mem.pause_total_ns", float64(current.PauseTotalNs-previous.PauseTotalNs))
	_ = collector.Gauge("go.runtime.mem.frees", float64(current.Frees-previous.Frees))
	_ = collector.Gauge("go.runtime.mem.mallocs", float64(current.Mallocs-previous.Mallocs))

	// these are mostly points in time.
	_ = collector.Gauge("go.runtime.num_cpu", float64(runtime.NumCPU()))
	_ = collector.Gauge("go.runtime.num_goroutine", float64(runtime.NumGoroutine()))

	_ = collector.Gauge("go.runtime.mem.alloc", float64(current.Alloc))

	_ = collector.Gauge("go.runtime.mem.gc_sys", float64(current.GCSys))
	_ = collector.Gauge("go.runtime.mem.other_sys", float64(current.OtherSys))

	_ = collector.Gauge("go.runtime.mem.heap_alloc", float64(current.HeapAlloc))
	_ = collector.Gauge("go.runtime.mem.heap_idle", float64(current.HeapIdle))
	_ = collector.Gauge("go.runtime.mem.heap_inuse", float64(current.HeapInuse))
	_ = collector.Gauge("go.runtime.mem.heap_objects", float64(current.HeapObjects))
	_ = collector.Gauge("go.runtime.mem.heap_sys", float64(current.HeapSys))

	_ = collector.Gauge("go.runtime.mem.stack_inuse", float64(current.StackInuse))
	_ = collector.Gauge("go.runtime.mem.stack_sys", float64(current.StackSys))
	_ = collector.Gauge("go.runtime.mem.sys", float64(current.Sys))
	_ = collector.Gauge("go.runtime.mem.total_alloc", float64(current.TotalAlloc))

	// rotate the results ...
	*previous = *current
}
