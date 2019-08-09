package breaker

import "time"

// Config is the breaker config.
type Config struct {
	HalfOpenMaxRequests  int64         `json:"halfOpenMaxRequests" yaml:"halfOpenMaxRequests"`
	ClosedExpiryInterval time.Duration `json:"closedExpiryInterval" yaml:"closedExpiryInterval"`
	OpenExpiryInterval   time.Duration `json:"openExpiryInterval" yaml:"openExpiryInterval"`
}
