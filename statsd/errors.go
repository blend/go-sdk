/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
