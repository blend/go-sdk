package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/blend/go-sdk/statsd"
)

var (
	bindAddr = flag.String("bind-addr", "127.0.0.1:0", "The bind address, defautls to a random local port")
	verbose  = flag.Bool("verbose", false, "If we should print each metric the server receives")
)

func main() {
	flag.Parse()
	log.SetOutput(logger{os.Stdout})

	listener, err := statsd.NewUDPListener(*bindAddr)
	if err != nil {
		log.Fatal(err)
	}

	srv := &statsd.Server{
		Listener: listener,
		Handler:  handleMetrics,
	}

	log.Printf("server listenening on: %s", listener.LocalAddr().String())
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}

var (
	metricRates   = map[string]*rate{}
	metricRatesMu sync.Mutex
)

func printRates() {
	for {
		<-time.Tick(2 * time.Second)

		for key, rate := range metricRates {
			log.Printf("%s: %s\n", key, rate.String())
		}
	}
}

func handleMetrics(ms ...statsd.Metric) {
	metricRatesMu.Lock()
	defer metricRatesMu.Unlock()

	for _, m := range ms {
		if *verbose {
			log.Printf("%#v\n", m)
		}
		_, ok := metricRates[m.Name]
		if !ok {
			metricRates[m.Name] = &rate{
				Count: 1,
				Time:  time.Now(),
			}
			return
		}
		metricRates[m.Name].Count++
	}
}

// rate
type rate struct {
	Count int
	Time  time.Time
}

func (r rate) String() string {
	if r.Count == 0 {
		return "N/A"
	}
	quantum := float64(time.Since(r.Time)) / float64(time.Second)
	rate := float64(r.Count) / quantum
	return fmt.Sprintf("%0.2f/s", rate)
}

type logger struct {
	wr io.Write
}

func (l logger) Write(contents []byte) (int, error) {
	return fmt.Fprint(l.wr, time.Now().UTC().Format(time.RFC3339Nano), " ", string(contents))
}
