package statsd

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func Test_ClientCount_Buffered(t *testing.T) {
	assert := assert.New(t)

	listener, err := NewUDPListener("127.0.0.1:0")
	assert.Nil(err)

	wg := sync.WaitGroup{}
	wg.Add(5) // 10/2 flushes

	metrics := make(chan Metric, 10)
	mock := &Server{
		Listener: listener,
		Handler: func(ms ...Metric) {
			defer wg.Done()
			for _, m := range ms {
				metrics <- m
			}
		},
	}
	go mock.Start()
	defer mock.Stop()

	client, err := New(
		OptAddr(mock.Listener.LocalAddr().String()),
		OptMaxBufferSize(2),
	)
	assert.Nil(err)

	for x := 0; x < 10; x++ {
		assert.Nil(client.Count(fmt.Sprintf("count%d", x), 10, Tag("env", "dev"), Tag("role", "test"), Tag("index", strconv.Itoa(x))))
	}

	wg.Wait()
	assert.Len(metrics, 10)
}

func Test_ClientGauge_Buffered(t *testing.T) {
	assert := assert.New(t)

	listener, err := NewUDPListener("127.0.0.1:0")
	assert.Nil(err)

	wg := sync.WaitGroup{}
	wg.Add(5)

	metrics := make(chan Metric, 10)
	mock := &Server{
		Listener: listener,
		Handler: func(ms ...Metric) {
			defer wg.Done()
			for _, m := range ms {
				metrics <- m
			}
		},
	}
	go mock.Start()
	defer mock.Stop()

	client, err := New(
		OptAddr(mock.Listener.LocalAddr().String()),
		OptMaxBufferSize(2),
	)
	assert.Nil(err)

	for x := 0; x < 10; x++ {
		assert.Nil(client.Gauge(fmt.Sprintf("gauge%d", x), 10, Tag("env", "dev"), Tag("role", "test"), Tag("index", strconv.Itoa(x))))
	}

	wg.Wait()
	assert.Len(metrics, 10)
}

func Test_ClientTimeInMilliseconds_Buffered(t *testing.T) {
	assert := assert.New(t)

	listener, err := NewUDPListener("127.0.0.1:0")
	assert.Nil(err)

	wg := sync.WaitGroup{}
	wg.Add(5)

	metrics := make(chan Metric, 10)
	mock := &Server{
		Listener: listener,
		Handler: func(ms ...Metric) {
			defer wg.Done()
			for _, m := range ms {
				metrics <- m
			}
		},
	}
	go mock.Start()
	defer mock.Stop()

	client, err := New(
		OptAddr(mock.Listener.LocalAddr().String()),
		OptMaxBufferSize(2),
	)
	assert.Nil(err)

	for x := 0; x < 10; x++ {
		assert.Nil(client.TimeInMilliseconds(fmt.Sprintf("time%d", x), 10*time.Millisecond, Tag("env", "dev"), Tag("role", "test"), Tag("index", strconv.Itoa(x))))
	}

	wg.Wait()
	assert.Len(metrics, 10)
}
