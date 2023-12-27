/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"net"
	"time"
)

// DialOption is a mutator for a net.Dialer
type DialOption func(*net.Dialer)

// OptDialTimeout sets the dial timeout.
func OptDialTimeout(d time.Duration) DialOption {
	return func(dialer *net.Dialer) {
		dialer.Timeout = d
	}
}

// OptDialKeepAlive sets the dial keep alive duration.
// Only use this if you know what you're doing, the defaults are typically sufficient.
func OptDialKeepAlive(d time.Duration) DialOption {
	return func(dialer *net.Dialer) {
		dialer.KeepAlive = d
	}
}
