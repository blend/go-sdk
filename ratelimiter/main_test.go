package ratelimiter

import (
	"time"
)

// Clock is a helper function.
func Clock(t time.Time, offset time.Duration) func() time.Time {
	return func() time.Time {
		return t.Add(offset)
	}
}
