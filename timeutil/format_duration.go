package timeutil

import (
	"fmt"
	"time"
)

// FormatDuration formats a duration to it's nearest major increment.
func FormatDuration(d time.Duration) string {
	if d >= time.Hour {
		return fmt.Sprintf("%dh", d.Round(time.Hour)/time.Hour)
	}
	if d >= time.Minute {
		return fmt.Sprintf("%dm", d.Round(time.Minute)/time.Minute)
	}
	if d >= time.Second {
		return fmt.Sprintf("%ds", d.Round(time.Second)/time.Second)
	}
	if d >= time.Millisecond {
		return fmt.Sprintf("%dms", d.Round(time.Millisecond)/time.Millisecond)
	}
	if d >= time.Microsecond {
		return fmt.Sprint(d.Round(time.Microsecond))
	}
	return fmt.Sprint(d)
}
