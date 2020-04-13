package statsd

import (
	"strconv"
	"time"
)

// Metric is a statsd metric.
type Metric struct {
	Name  string
	Type  string
	Value string
	Tags  []string
}

// Float64 returns the value parsed as a float64.
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
