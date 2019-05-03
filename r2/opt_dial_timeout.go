package r2

import (
	"net"
	"time"
)

// OptDialTimeout sets the dial timeout.
func OptDialTimeout(d time.Duration) DialOption {
	return func(dialer *net.Dialer) {
		dialer.Timeout = d
	}
}
