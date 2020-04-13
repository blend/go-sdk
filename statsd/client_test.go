package statsd

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

type noOpWriteCloser struct {
	io.Writer
}

// Close is a no-op.
func (n noOpWriteCloser) Close() error { return nil }

func Test_Client_Options(t *testing.T) {
	assert := assert.New(t)

	c := new(Client)
	assert.Empty(c.Addr)
	assert.Nil(OptAddr("192.168.1.1:0")(c))
	assert.Equal("192.168.1.1:0", c.Addr)

	assert.Zero(c.DialTimeout)
	assert.Nil(OptDialTimeout(time.Second)(c))
	assert.Equal(time.Second, c.DialTimeout)

	assert.Zero(c.MaxPacketSize)
	assert.Nil(OptMaxPacketSize(1024)(c))
	assert.Equal(1024, c.MaxPacketSize)

	assert.Zero(c.MaxBufferSize)
	assert.Nil(OptMaxBufferSize(512)(c))
	assert.Equal(512, c.MaxBufferSize)

	cfg := Config{
		Addr:          "127.0.0.1:0",
		DialTimeout:   500 * time.Millisecond,
		MaxPacketSize: 1024,
		MaxBufferSize: 64,
		DefaultTags: map[string]string{
			"foo": "bar",
			"env": "sandbox",
		},
		SampleRate: 0.8,
	}

	configClient := new(Client)
	assert.Nil(configClient.SampleProvider)
	assert.Nil(OptConfig(cfg)(configClient))

	assert.Equal("127.0.0.1:0", configClient.Addr)
	assert.Equal(500*time.Millisecond, configClient.DialTimeout)
	assert.Equal(1024, configClient.MaxPacketSize)
	assert.Equal(64, configClient.MaxBufferSize)

	assert.Any(configClient.DefaultTags(), func(v interface{}) bool { return v.(string) == "foo:bar" })
	assert.Any(configClient.DefaultTags(), func(v interface{}) bool { return v.(string) == "env:sandbox" })

	assert.NotNil(configClient.SampleProvider)
}

func Test_Client_AddDefaultTag(t *testing.T) {
	assert := assert.New(t)

	c := new(Client)
	assert.Empty(c.defaultTags)
	c.AddDefaultTag("foo", "bar")
	assert.Equal([]string{"foo:bar"}, c.defaultTags)
}

func Test_ClientCount_Sampling(t *testing.T) {
	assert := assert.New(t)

	buffer := new(bytes.Buffer)

	client := &Client{
		SampleProvider: func() bool {
			return rand.Float64() < 0.5
		},
		conn: noOpWriteCloser{buffer},
	}

	for x := 0; x < 512; x++ {
		assert.Nil(client.Count("sampling test", int64(x)))
	}

	contents := strings.Split(buffer.String(), "\n")
	assert.True(len(contents) > 200, len(contents))
	assert.True(len(contents) < 300, len(contents))
}

func Test_ClientGauge_Sampling(t *testing.T) {
	assert := assert.New(t)

	buffer := new(bytes.Buffer)

	client := &Client{
		SampleProvider: func() bool {
			return rand.Float64() < 0.5
		},
		conn: noOpWriteCloser{buffer},
	}

	for x := 0; x < 512; x++ {
		assert.Nil(client.Gauge("sampling test", float64(x)))
	}

	contents := strings.Split(buffer.String(), "\n")
	assert.True(len(contents) > 200, len(contents))
	assert.True(len(contents) < 300, len(contents))
}

func Test_ClientCount_Unbuffered(t *testing.T) {
	assert := assert.New(t)

	listener, err := NewUDPListener("127.0.0.1:0")
	assert.Nil(err)

	wg := sync.WaitGroup{}
	wg.Add(10)

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
		OptMaxBufferSize(0),
	)
	assert.Nil(err)

	for x := 0; x < 10; x++ {
		assert.Nil(client.Count(fmt.Sprintf("count%d", x), 10+int64(x), Tag("env", "dev"), Tag("role", "test"), Tag("index", strconv.Itoa(x))))
	}

	wg.Wait()
	assert.Len(metrics, 10)

	m := <-metrics
	assert.Equal("c", m.Type)
	assert.Equal("count0", m.Name)
	assert.Equal("10", m.Value)
	assert.Equal([]string{"env:dev", "role:test", "index:0"}, m.Tags)

	m = <-metrics
	assert.Equal("c", m.Type)
	assert.Equal("count1", m.Name)
	assert.Equal("11", m.Value)
	assert.Equal([]string{"env:dev", "role:test", "index:1"}, m.Tags)
}

func Test_ClientGauge_Unbuffered(t *testing.T) {
	assert := assert.New(t)

	listener, err := NewUDPListener("127.0.0.1:0")
	assert.Nil(err)

	wg := sync.WaitGroup{}
	wg.Add(10)

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
		OptMaxBufferSize(0),
	)
	assert.Nil(err)

	for x := 0; x < 10; x++ {
		assert.Nil(client.Gauge(fmt.Sprintf("gauge%d", x), 10+float64(x), Tag("env", "dev"), Tag("role", "test"), Tag("index", strconv.Itoa(x))))
	}

	wg.Wait()
	assert.Len(metrics, 10)

	m := <-metrics
	assert.Equal("g", m.Type)
	assert.Equal("gauge0", m.Name)
	assert.Equal("10", m.Value)
	assert.Equal([]string{"env:dev", "role:test", "index:0"}, m.Tags)

	m = <-metrics
	assert.Equal("g", m.Type)
	assert.Equal("gauge1", m.Name)
	assert.Equal("11", m.Value)
	assert.Equal([]string{"env:dev", "role:test", "index:1"}, m.Tags)
}

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
