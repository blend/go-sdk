/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
