package statsd

import (
	"sync"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_ClientCount(t *testing.T) {
	assert := assert.New(t)

	listener, err := NewUDPListener("127.0.0.1:0")
	assert.Nil(err)

	wg := sync.WaitGroup{}
	wg.Add(4)

	metrics := make(chan Metric, 4)
	mock := &Server{
		Listener: listener,
		Handler: func(m Metric) {
			defer wg.Done()
			metrics <- m
		},
	}
	go mock.Start()
	defer mock.Stop()

	client, err := New(OptAddr(mock.Listener.LocalAddr().String()))
	assert.Nil(err)

	assert.Nil(client.Count("count0", 10, Tag("env", "dev"), Tag("role", "test")))
	assert.Nil(client.Count("count0", 5, Tag("env", "sandbox"), Tag("role", "test")))
	assert.Nil(client.Count("count1", 10, Tag("env", "dev"), Tag("role", "test")))
	assert.Nil(client.Count("count1", 5, Tag("env", "sandbox"), Tag("role", "test")))

	wg.Wait()

	assert.Len(metrics, 4)

	metric := <-metrics
	assert.Equal("count0", metric.Name)
	assert.Equal("10", metric.Value)
	assert.Equal("c", metric.Type)
	assert.Len(metric.Tags, 2)

	metric = <-metrics
	assert.Equal("count0", metric.Name)
	assert.Equal("5", metric.Value)
	assert.Equal("c", metric.Type)
	assert.Len(metric.Tags, 2)

	metric = <-metrics
	assert.Equal("count1", metric.Name)
	assert.Equal("10", metric.Value)
	assert.Equal("c", metric.Type)
	assert.Len(metric.Tags, 2)

	metric = <-metrics
	assert.Equal("count1", metric.Name)
	assert.Equal("5", metric.Value)
	assert.Equal("c", metric.Type)
	assert.Len(metric.Tags, 2)
}
