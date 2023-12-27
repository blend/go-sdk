/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/blend/go-sdk/statsd"
)

var (
	addr		= flag.String("addr", "127.0.0.1:8125", "The statsd server address")
	dialTimeout	= flag.Duration("dial-timeout", time.Second, "The client dial timeout")
	bufferSize	= flag.Int("buffer-size", 64, "The client buffer size")
	workers		= flag.Int("workers", runtime.NumCPU(), "The number of workers to use")
)

var metrics = []statsd.Metric{
	{Type: statsd.MetricTypeCount, Name: "http.request", Value: "1", Tags: []string{statsd.Tag("env", "test")}},
	{Type: statsd.MetricTypeCount, Name: "error", Value: "1", Tags: []string{statsd.Tag("env", "test")}},
	{Type: statsd.MetricTypeCount, Name: "http.response", Value: "1", Tags: []string{statsd.Tag("env", "test"), statsd.Tag("status_code", "200")}},
	{Type: statsd.MetricTypeTimer, Name: "http.response.elapsed", Value: "500.0", Tags: []string{statsd.Tag("env", "test"), statsd.Tag("status_code", "200")}},
}

func main() {
	c, err := statsd.New(
		statsd.OptAddr(*addr),
		statsd.OptDialTimeout(*dialTimeout),
		statsd.OptMaxBufferSize(*bufferSize),
	)
	if err != nil {
		log.Fatal(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(*workers)
	started := time.Now()
	var sent int32
	for workerID := 0; workerID < *workers; workerID++ {
		go func(id int) {
			defer wg.Done()
			var err error
			for x := 0; x < 1024; x++ {
				for _, m := range metrics {
					switch m.Type {
					case statsd.MetricTypeCount:
						v, _ := m.Int64()
						err = c.Count(m.Name, v, m.Tags...)
					case statsd.MetricTypeGauge:
						v, _ := m.Float64()
						err = c.Gauge(m.Name, v, m.Tags...)
					case statsd.MetricTypeTimer:
						v, _ := m.Duration()
						err = c.TimeInMilliseconds(m.Name, v, m.Tags...)
					}
					if err != nil {
						log.Printf("client error: %v\n", err)
					}
					atomic.AddInt32(&sent, 1)
				}
			}
		}(workerID)
	}
	wg.Wait()

	elapsed := time.Since(started)
	fmt.Printf("sent %d messages in %v (%0.2f m/sec)\n", sent, elapsed, float64(sent)/(float64(elapsed)/float64(time.Second)))
}
