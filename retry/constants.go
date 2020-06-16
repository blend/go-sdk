package retry

import "time"

// Defaults
const (
	DefaultMaxAttempts = 5
	DefaultRetryDelay  = time.Second
)
