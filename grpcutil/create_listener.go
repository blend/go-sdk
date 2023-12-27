/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package grpcutil

import (
	"net"
	"strings"

	"github.com/blend/go-sdk/ex"
)

// CreateListener creates a net listener for a given bind address.
// It handles detecting if we should create a unix socket address.
func CreateListener(bindAddr string) (net.Listener, error) {
	var socketListener net.Listener
	var err error
	if strings.HasPrefix(bindAddr, "unix://") {
		socketListener, err = net.Listen("unix", strings.TrimPrefix(bindAddr, "unix://"))
		if typed, ok := socketListener.(*net.UnixListener); ok {
			typed.SetUnlinkOnClose(true)
		}
	} else {
		socketListener, err = net.Listen("tcp", bindAddr)
	}
	return socketListener, ex.New(err)
}
