/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package statsd

import (
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/stats"
)

// Error classes.
const (
	ErrAddrUnset		ex.Class	= "statsd client address unset"
	ErrMaxPacketSize	ex.Class	= "statsd max packet size exceeded"
	ErrSampleRateInvalid	ex.Class	= "statsd invalid sample rate"
)

var (
	_ stats.Collector = (*Client)(nil)
)

// New creates a new statsd client and opens
// the underlying UDP connection.
func New(opts ...ClientOpt) (*Client, error) {
	client := Client{
		DialTimeout:	DefaultDialTimeout,
		MaxPacketSize:	DefaultMaxPacketSize,
		MaxBufferSize:	DefaultMaxBufferSize,
	}

	var err error
	for _, opt := range opts {
		if err = opt(&client); err != nil {
			return nil, err
		}
	}
	if client.Addr == "" {
		return nil, ex.New(ErrAddrUnset)
	}
	client.conn, err = net.DialTimeout("udp", client.Addr, client.DialTimeout)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// OptAddr sets the client address.
func OptAddr(addr string) ClientOpt {
	return func(c *Client) error {
		c.Addr = addr
		return nil
	}
}

// OptDialTimeout sets the client dial timeout.
func OptDialTimeout(timeout time.Duration) ClientOpt {
	return func(c *Client) error {
		c.DialTimeout = timeout
		return nil
	}
}

// OptMaxPacketSize sets the client max dial size.
func OptMaxPacketSize(sizeBytes int) ClientOpt {
	return func(c *Client) error {
		c.MaxPacketSize = sizeBytes
		return nil
	}
}

// OptMaxBufferSize sets the client max buffer size in messages.
func OptMaxBufferSize(count int) ClientOpt {
	return func(c *Client) error {
		c.MaxBufferSize = count
		return nil
	}
}

// OptConfig sets fields on a client from a given config.
func OptConfig(cfg Config) ClientOpt {
	return func(c *Client) error {
		c.Addr = cfg.Addr
		c.DialTimeout = cfg.DialTimeout
		c.MaxPacketSize = cfg.MaxPacketSize
		c.MaxBufferSize = cfg.MaxBufferSize
		for key, value := range cfg.DefaultTags {
			c.AddDefaultTags(stats.Tag(key, value))
		}
		return OptSampleRate(cfg.SampleRate)(c)
	}
}

// OptSampleRate sets the sample rate on the client or the percent of packets to send
// on the interval [0,1.0).
// A value of `0.0` will drop all packets, a value of `1.0` will send all packets.
func OptSampleRate(rate float64) ClientOpt {
	return func(c *Client) error {
		if rate < 0 || rate > 1.0 {
			return ex.New(ErrSampleRateInvalid, ex.OptMessagef("rate: %0.2f", rate))
		}
		if rate == 1.0 {	// unset on 100%
			c.SampleProvider = nil
		} else {
			c.SampleProvider = func() bool {
				return rand.Float64() <= rate
			}
		}
		return nil
	}
}

// ClientOpt is an option for a client.
type ClientOpt func(*Client) error

// Client is a statsd client.
type Client struct {
	Addr		string
	DialTimeout	time.Duration
	MaxPacketSize	int
	SampleProvider	func() bool
	MaxBufferSize	int

	defaultTags	[]string

	conn	io.WriteCloser
	connMu	sync.Mutex

	bufferMu	sync.Mutex
	buffer		[]byte
	bufferCount	int
}

// AddDefaultTag adds a new default tag.
func (c *Client) AddDefaultTag(name, value string) {
	c.defaultTags = append(c.defaultTags, stats.Tag(name, value))
}

// AddDefaultTags adds default tags.
func (c *Client) AddDefaultTags(tags ...string) {
	c.defaultTags = append(c.defaultTags, tags...)
}

// DefaultTags returns the default tags.
func (c *Client) DefaultTags() []string {
	return c.defaultTags
}

// Count sends a count message.
func (c *Client) Count(name string, value int64, tags ...string) error {
	return c.sendInt(MetricTypeCount, name, value, tags...)
}

// Increment sends a count message with a value of (1).
func (c *Client) Increment(name string, tags ...string) error {
	return c.sendInt(MetricTypeCount, name, 1, tags...)
}

// Gauge sends a point in time value.
func (c *Client) Gauge(name string, value float64, tags ...string) error {
	return c.sendFloat(MetricTypeGauge, name, value, tags...)
}

// TimeInMilliseconds sends a gauge method with a given value represented in milliseconds.
func (c *Client) TimeInMilliseconds(name string, value time.Duration, tags ...string) error {
	return c.sendFloat(MetricTypeTimer, name, float64(value)/float64(time.Millisecond), tags...)
}

// Histogram is an no-op for raw statsd.
func (c *Client) Histogram(name string, value float64, tags ...string) error {
	return c.sendFloat(MetricTypeHistogram, name, value, tags...)
}

// Distribution is an no-op for raw statsd.
func (c *Client) Distribution(name string, value float64, tags ...string) error {
	return c.sendFloat(MetricTypeDistribution, name, value, tags...)
}

// Flush is a no-op.
func (c *Client) Flush() error {
	c.bufferMu.Lock()
	defer c.bufferMu.Unlock()
	return c.flushBuffer()
}

// Close closes the underlying connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) sendInt(metricType, name string, value int64, tags ...string) error {
	if !c.shouldSend() {
		return nil
	}
	if c.MaxBufferSize == 0 {
		return c.send(c.appendInt(nil, metricType, name, value, tags...))
	}
	c.bufferMu.Lock()
	defer c.bufferMu.Unlock()

	c.bufferCount++
	c.buffer = c.appendInt(c.buffer, metricType, name, value, tags...)
	if c.bufferCount < c.MaxBufferSize {
		c.buffer = c.appendMetricSeparator(c.buffer)
		return nil
	}
	return c.flushBuffer()
}

func (c *Client) sendFloat(metricType, name string, value float64, tags ...string) error {
	if !c.shouldSend() {
		return nil
	}
	if c.MaxBufferSize == 0 {
		return c.send(c.appendFloat(nil, metricType, name, value, tags...))
	}
	c.bufferMu.Lock()
	defer c.bufferMu.Unlock()

	c.bufferCount++
	c.buffer = c.appendFloat(c.buffer, metricType, name, value, tags...)
	if c.bufferCount < c.MaxBufferSize {
		c.buffer = c.appendMetricSeparator(c.buffer)
		return nil
	}
	return c.flushBuffer()
}

func (c *Client) appendInt(data []byte, metricType, name string, value int64, tags ...string) []byte {
	data = append(data, []byte(name)...)
	data = append(data, ':')
	data = strconv.AppendInt(data, value, 10)
	data = append(data, '|')
	data = append(data, []byte(metricType)...)
	data = c.appendTags(data, append(c.defaultTags, tags...)...)
	return data
}

func (c *Client) appendFloat(data []byte, metricType, name string, value float64, tags ...string) []byte {
	data = append(data, []byte(name)...)
	data = append(data, ':')
	data = strconv.AppendFloat(data, value, 'f', -1, 64)
	data = append(data, '|')
	data = append(data, []byte(metricType)...)
	data = c.appendTags(data, append(c.defaultTags, tags...)...)
	return data
}

func (c *Client) appendTags(data []byte, tags ...string) []byte {
	if len(tags) == 0 {
		return data
	}
	data = append(data, "|#"...)
	firstTag := true
	for _, tag := range tags {
		if !firstTag {
			data = append(data, ',')
		}
		data = append(data, strings.TrimSpace(tag)...)
		firstTag = false
	}
	return data
}

func (c *Client) appendMetricSeparator(data []byte) []byte {
	return append(data, '\n')
}

func (c *Client) shouldSend() bool {
	if c.SampleProvider == nil {
		return true
	}
	return c.SampleProvider()
}

func (c *Client) flushBuffer() error {
	if err := c.send(c.buffer); err != nil {
		return err
	}
	c.bufferCount = 0
	c.buffer = nil
	return nil
}

func (c *Client) send(data []byte) error {
	if c.MaxPacketSize > 0 && len(data) > c.MaxPacketSize {
		return ex.New(ErrMaxPacketSize)
	}

	c.connMu.Lock()
	defer c.connMu.Unlock()

	_, err := c.conn.Write(append(data, '\n'))
	if err != nil {
		return ex.New(err)
	}
	return nil
}
