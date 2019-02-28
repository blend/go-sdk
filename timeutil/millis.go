package timeutil

import "time"

// Millis returns the given duration as milliseconds
func Millis(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}
