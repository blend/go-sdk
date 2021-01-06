package r2

import (
	"net"
	"time"

	"github.com/blend/go-sdk/webutil"
)

// OptDialTimeout sets the dial timeout.
func OptDialTimeout(d time.Duration) DialOption {
	return func(dialer *net.Dialer) {
		webutil.OptDialTimeout(d)(dialer)
	}
}
