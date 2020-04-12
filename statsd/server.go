package statsd

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Server is a listener for statsd metrics.
// It is meant to be used for diagnostic purposes, and is not suitable for
// production anything.
type Server struct {
	Addr     string
	Listener net.PacketConn
	Handler  func(...Metric)
}

// Start starts the server. This call blocks.
func (s *Server) Start() error {
	var err error
	if s.Listener == nil && s.Addr != "" {
		s.Listener, err = NewUDPListener(s.Addr)
		if err != nil {
			return err
		}
	}
	if s.Listener == nil {
		return fmt.Errorf("server cannot start; no listener or addr provided")
	}

	data := make([]byte, DefaultMaxPacketSize)
	var metrics []Metric
	var n int
	for {
		n, _, err = s.Listener.ReadFrom(data)
		if err != nil {
			return err
		}
		metrics, err = s.parseMetrics(data[:n])
		if err != nil {
			return err
		}
		go s.Handler(metrics...)
	}
}

// Stop closes the server with a given context.
func (s *Server) Stop() error {
	if s.Listener == nil {
		return nil
	}
	return s.Listener.Close()
}

func (s *Server) parseMetrics(data []byte) ([]Metric, error) {
	var metrics []Metric
	for index := 0; index < len(data); index++ {
		m, err := s.parseMetric(&index, data)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

// parseMetric parses a metric from a given data packet.
func (s *Server) parseMetric(index *int, data []byte) (m Metric, err error) {
	var name []byte
	var metricType []byte
	var value []byte
	var tags []byte

	var state int
	var b byte
	for ; *index < len(data); (*index)++ {
		b = data[*index]
		switch state {
		case 0: //name
			if b == ':' {
				state = 1
				continue
			}
			name = append(name, b)
			continue
		case 1: //value
			if b == '|' {
				state = 2
				continue
			}
			value = append(value, b)
			continue
		case 2:
			if b == '|' {
				state = 3
				continue
			}
			metricType = append(metricType, b)
			continue
		case 3: // tags
			if b == '#' {
				continue
			}
			if b == '\n' {
				break // drop out at newline
			}
			tags = append(tags, b)
		}
	}

	m.Name = string(name)
	m.Type = string(metricType)
	m.Value = string(value)
	m.Tags = strings.Split(string(tags), ",")
	return
}

// NewUDPListener returns a new UDP listener for a given address.
func NewUDPListener(addr string) (net.PacketConn, error) {
	listener, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

// Metric is a statsd metric.
type Metric struct {
	Name  string
	Type  string
	Value string
	Tags  []string
}

// Float64 returns the value parsed as a float64.s
func (m Metric) Float64() (float64, error) {
	return strconv.ParseFloat(m.Value, 64)
}

// Int64 returns the value parsed as an int64.
func (m Metric) Int64() (int64, error) {
	return strconv.ParseInt(m.Value, 10, 64)
}

// Duration is the value parsed as a duration assuming
// it was a float64 of milliseconds.
func (m Metric) Duration() (time.Duration, error) {
	f64, err := m.Float64()
	if err != nil {
		return 0, err
	}
	return time.Duration(f64 * float64(time.Millisecond)), nil
}
