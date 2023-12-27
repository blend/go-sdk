/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package reverseproxy

import (
	"net"
	"time"

	"github.com/blend/go-sdk/webutil"
)

// DialOption is a mutator for a net.Dialer.
type DialOption = webutil.DialOption

// OptDialTimeout sets the dial timeout.
func OptDialTimeout(d time.Duration) DialOption {
	return func(dialer *net.Dialer) {
		webutil.OptDialTimeout(d)(dialer)
	}
}

// OptDialKeepAlive sets the dial keep alive duration.
// Only use this if you know what you're doing, the defaults are typically sufficient.
func OptDialKeepAlive(d time.Duration) DialOption {
	return func(dialer *net.Dialer) {
		webutil.OptDialKeepAlive(d)(dialer)
	}
}
