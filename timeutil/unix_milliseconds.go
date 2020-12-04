package timeutil

import "time"

// UnixMilliseconds returns the time in unix (seconds) format
// with a floating point remainder for subsecond fraction.
func UnixMilliseconds(t time.Time) float64 {
	nanosPerSecond := float64(time.Second) / float64(time.Nanosecond)
	return float64(t.UnixNano()) / nanosPerSecond
}
