/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package statsd

import (
	"fmt"
	"log"
	"net"
	"strings"
)

// Server is a listener for statsd metrics.
// It is meant to be used for diagnostic purposes, and is not suitable for
// production anything.
type Server struct {
	Addr		string
	Log		*log.Logger
	MaxPacketSize	int
	Listener	net.PacketConn
	Handler		func(...Metric)
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
	if s.Handler == nil {
		return fmt.Errorf("server cannot start; no handler provided")
	}
	if s.Listener == nil && s.Addr != "" {
		s.Listener, err = NewUDPListener(s.Addr)
		if err != nil {
			return err
		}
	}
	if s.Listener == nil {
		return fmt.Errorf("server cannot start; no listener or addr provided")
	}

	s.logf("statsd server listening: %s", s.Listener.LocalAddr().String())
	data := make([]byte, s.MaxPacketSizeOrDefault())
	var metrics []Metric
	var n int
	for {
		n, _, err = s.Listener.ReadFrom(data)
		if IsErrUseOfClosedNetworkConnection(err) {
			return nil
		} else if err != nil {
			return err
		}
		metrics, err = s.parseMetrics(data[:n])
		if err != nil {
			return err
		}
		s.Handler(metrics...)
	}
}

// Stop closes the server.
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

	const (
		stateName	= iota
		stateValue	= iota
		stateMetricType	= iota
		stateTags	= iota
		stateTagValues	= iota
	)

	var b byte
	state := stateName
	for ; *index < len(data); (*index)++ {
		b = data[*index]

		if b == '\n' {
			break	// drop out at newline
		}

		switch state {
		case stateName:	//name
			if b == ':' {
				state = stateValue
				continue
			}
			name = append(name, b)
			continue
		case stateValue:	//value
			if b == '|' {
				state = stateMetricType
				continue
			}
			value = append(value, b)
			continue
		case stateMetricType:	// metric type
			if b == '|' {
				state = stateTags
				continue
			}
			metricType = append(metricType, b)
			continue
		case stateTags:	// tags
			if b == '#' {
				state = stateTagValues
				continue
			}
			err = fmt.Errorf("invalid metric; tags should be marked with '#'")
			return
		case stateTagValues:
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

//
// logging
//

func (s *Server) logf(format string, args ...interface{}) {
	if s.Log != nil {
		format = strings.TrimSpace(format)
		s.Log.Printf(format+"\n", args...)
	}
}

func (s *Server) logln(args ...interface{}) {
	if s.Log != nil {
		s.Log.Println(args...)
	}
}

// FormatContent
