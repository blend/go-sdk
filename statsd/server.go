package statsd

import (
	"fmt"
	"net"
	"time"
)

// Server is a listener for statsd metrics.
// It is meant to be used for diagnostic purposes, and is not suitable for
// production anything.
type Server struct {
	Addr          string
	ReadDeadline  time.Duration
	MaxPacketSize int
	Listener      net.PacketConn
	Handler       func(...Metric)
}

// MaxPacketSizeOrDefault returns the max packet size or a default.
func (s *Server) MaxPacketSizeOrDefault() int {
	if s.MaxPacketSize > 0 {
		return s.MaxPacketSize
	}
	return DefaultMaxPacketSize
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

	data := make([]byte, s.MaxPacketSizeOrDefault())
	var metrics []Metric
	var n int
	for {
		if s.ReadDeadline > 0 {
			if err := s.Listener.SetReadDeadline(time.Now().Add(s.ReadDeadline)); err != nil {
				return err
			}
		}
		n, _, err = s.Listener.ReadFrom(data)
		if err != nil {
			return err
		}
		metrics, err = s.parseMetrics(data[:n])
		if err != nil {
			return err
		}
		s.Handler(metrics...)
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
	var tag []byte

	var state int
	var b byte
	for ; *index < len(data); (*index)++ {
		b = data[*index]

		if b == '\n' {
			break // drop out at newline
		}

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
		case 2: // metric type
			if b == '|' {
				state = 3
				continue
			}
			metricType = append(metricType, b)
			continue
		case 3: // tags
			if b == '#' {
				state = 4
				continue
			}
			err = fmt.Errorf("invalid metric; tags should be marked with '#'")
			return
		case 4:
			if b == ',' {
				m.Tags = append(m.Tags, string(tag))
				tag = nil
				continue
			}
			tag = append(tag, b)
		}
	}
	if len(tag) > 0 {
		m.Tags = append(m.Tags, string(tag))
	}

	m.Name = string(name)
	m.Type = string(metricType)
	m.Value = string(value)
	return
}
