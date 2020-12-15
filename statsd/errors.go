package statsd

import "net"

// IsErrUseOfClosedNetworkConnection is an error class checker.
func IsErrUseOfClosedNetworkConnection(err error) bool {
	if err == nil {
		return false
	}
	typed, ok := err.(*net.OpError)
	if !ok {
		return false
	}
	if typed.Temporary() || typed.Timeout() {
		return false
	}
	return typed.Err.Error() == "use of closed network connection"
}
