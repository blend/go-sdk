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
	addr    = flag.String("addr", "127.0.0.1:8125", "The statsd server address")
	workers = flag.Int("workers", runtime.NumCPU(), "The number of workers to use")
)

var metrics = []statsd.Metric{
	{Type: "c", Name: "http.request", Value: "1", Tags: []string{statsd.Tag("env", "test")}},
	{Type: "c", Name: "error", Value: "1", Tags: []string{statsd.Tag("env", "test")}},
	{Type: "c", Name: "http.response", Value: "1", Tags: []string{statsd.Tag("env", "test"), statsd.Tag("status_code", "200")}},
	{Type: "ms", Name: "http.response.elapsed", Value: "500.0", Tags: []string{statsd.Tag("env", "test"), statsd.Tag("status_code", "200")}},
}

func main() {
	c, err := statsd.New(
		statsd.OptAddr(*addr),
		statsd.OptDialTimeout(250*time.Millisecond),
		statsd.OptMaxBufferSize(64),
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
					case "c":
						v, _ := m.Int64()
						err = c.Count(m.Name, v, m.Tags...)
					case "g":
						v, _ := m.Float64()
						err = c.Gauge(m.Name, v, m.Tags...)
					case "ms":
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
