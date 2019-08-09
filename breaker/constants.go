package breaker

import "time"

// Constants
const (
	DefaultClosedExpiryInterval = 5 * time.Second
	DefaultOpenExpiryInterval   = 60 * time.Second
	DefaultHalfOpenMaxRequests  = 1
)
