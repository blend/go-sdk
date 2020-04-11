package statsd

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Server is a listener for statsd metrics.
type Server struct {
	Addr     string
	Listener *net.UDPConn
	Handler  func(Metric)
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
	var metric Metric
	for {
		_, err = s.Listener.Read(data)
		if err != nil {
			return err
		}
		println("server read", string(data))
		err = s.parseMetric(data, &metric)
		if err != nil {
			return err
		}
		s.Handler(metric)
	}
}

// Stop closes the server with a given context.
func (s *Server) Stop() error {
	if s.Listener == nil {
		return nil
	}
	return s.Listener.Close()
}

// parseMetric parses a metric from a given data packet.
func (s *Server) parseMetric(data []byte, m *Metric) error {

	var name []byte
	var metricType []byte
	var value []byte
	var tags []byte

	var state int
	var b byte
	for x := 0; x < len(data); x++ {
		b = data[x]
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
			tags = append(tags, b)
		}
	}

	m.Name = string(name)
	m.Type = string(metricType)
	m.Value = string(value)
	m.Tags = strings.Split(string(tags), ",")
	return nil
}

// NewUDPListener returns a new UDP listener for a given address.
func NewUDPListener(addr string) (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenUDP("udp", udpAddr)
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
