/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
	if len(parts) == 0 {
		return 0
	}
	if len(parts) < 2 {
		output, _ := strconv.ParseInt(parts[0], 10, 64)
		port = int32(output)
		return
	}
	output, _ := strconv.ParseInt(parts[1], 10, 64)
	port = int32(output)
	return
}
