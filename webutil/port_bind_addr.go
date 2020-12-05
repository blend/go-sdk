package webutil

import (
	"strconv"
	"strings"
)

// PortFromBindAddr returns a port number as an integer from a bind addr.
func PortFromBindAddr(bindAddr string) (port int32) {
	if len(bindAddr) == 0 {
		return 0
	}
	parts := strings.SplitN(bindAddr, ":", 2)
	// NOTE: We don't check that `len(parts) > 0`. An empty (or `nil`) slice
	//       can **only** be returned if the `n` argument to `SplitN` is equal
	//       to `0`.
	if len(parts) < 2 {
		output, _ := strconv.ParseInt(parts[0], 10, 64)
		port = int32(output)
		return
	}
	output, _ := strconv.ParseInt(parts[1], 10, 64)
	port = int32(output)
	return
}
