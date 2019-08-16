package breaker

import "time"

// Config is the breaker config.
type Config struct {
	HalfOpenMaxActions   int64         `json:"halfOpenMaxActions" yaml:"halfOpenMaxActions"`
	ClosedExpiryInterval time.Duration `json:"closedExpiryInterval" yaml:"closedExpiryInterval"`
	OpenExpiryInterval   time.Duration `json:"openExpiryInterval" yaml:"openExpiryInterval"`
}
