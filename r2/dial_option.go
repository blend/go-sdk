package r2

import "net"

// DialOption is an option for the net dialer.
type DialOption func(*net.Dialer)
