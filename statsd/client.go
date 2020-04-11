package statsd

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/blend/go-sdk/ex"
)

// Error classes.
const (
	ErrAddrUnset         ex.Class = "statsd client address unset"
	ErrMaxPacketSize     ex.Class = "statsd max packet size exceeded"
	ErrSampleRateInvalid ex.Class = "statsd invalid sample rate"
)

// New creates a new statsd client and opens
// the underlying UDP connection.
func New(opts ...ClientOpt) (*Client, error) {
	client := Client{
		DialTimeout:   DefaultDialTimeout,
		MaxPacketSize: DefaultMaxPacketSize,
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

// OptConfig sets fields on a client from a given config.
func OptConfig(cfg Config) ClientOpt {
	return func(c *Client) error {
		c.Addr = cfg.Addr
		c.DialTimeout = cfg.DialTimeout
		c.MaxPacketSize = cfg.MaxPacketSize
		for key, value := range cfg.DefaultTags {
			c.defaultTags = append(c.defaultTags, Tag(key, value))
		}
		return nil
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
		if rate == 1.0 { // unset on 100%
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
	Addr           string
	DialTimeout    time.Duration
	MaxPacketSize  int
	SampleProvider func() bool

	conn        net.Conn
	mu          sync.Mutex
	defaultTags []string
}

// AddDefaultTag adds a default tag.
func (c *Client) AddDefaultTag(key, value string) {
	c.defaultTags = append(c.defaultTags, fmt.Sprintf("%s:%s", key, value))
}

// DefaultTags returns the default tags.
func (c *Client) DefaultTags() []string {
	return c.defaultTags
}

// Count sends a count message.
func (c *Client) Count(name string, value int64, tags ...string) error {
	if !c.shouldSend() {
		return nil
	}
	return c.send(c.appendInt([]byte{}, "c", name, value, tags...))
}

// Increment sends a count message with a value of (1).
func (c *Client) Increment(name string, tags ...string) error {
	return c.Count(name, 1, tags...)
}

// Gauge sends a point in time value.
func (c *Client) Gauge(name string, value float64, tags ...string) error {
	if !c.shouldSend() {
		return nil
	}
	return c.send(c.appendFloat([]byte{}, "g", name, value, tags...))
}

// Histogram is an no-op for raw statsd.
func (c *Client) Histogram(name string, value float64, tags ...string) error {
	return nil
}

// TimeInMilliseconds sends a gauge method with a given value represented in milliseconds.
func (c *Client) TimeInMilliseconds(name string, value time.Duration, tags ...string) error {
	return c.Gauge(name, float64(value)/float64(time.Millisecond), tags...)
}

// Flush is a no-op.
func (c *Client) Flush() error {
	return nil
}

// Close closes the underlying connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) appendInt(data []byte, metricType, name string, value int64, tags ...string) []byte {
	data = append(data, []byte(name)...)
	data = append(data, ':')
	data = strconv.AppendInt(data, value, 10)
	data = append(data, '|')
	data = append(data, []byte(metricType)...)
	data = c.appendTags(data, tags...)
	return data
}

func (c *Client) appendFloat(data []byte, metricType, name string, value float64, tags ...string) []byte {
	data = append(data, []byte(name)...)
	data = append(data, ':')
	data = strconv.AppendFloat(data, value, 'f', -1, 64)
	data = append(data, '|')
	data = append(data, []byte(metricType)...)
	data = c.appendTags(data, tags...)
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

func (c *Client) shouldSend() bool {
	if c.SampleProvider == nil {
		return true
	}
	return c.SampleProvider()
}

func (c *Client) send(data []byte) error {
	if len(data) > c.MaxPacketSize {
		return ex.New(ErrMaxPacketSize)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.conn.Write(data)
	if err != nil {
		return ex.New(err)
	}
	return nil
}
