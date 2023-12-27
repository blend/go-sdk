/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package statsd

import (
	"net"
	"strings"
)

// NewUDPListener returns a new UDP listener for a given address.
func NewUDPListener(addr string) (net.PacketConn, error) {
	listener, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

// NewUnixgramListener returns a new unixgram listener for a given path.
func NewUnixgramListener(path string) (net.PacketConn, error) {
	path = strings.TrimPrefix(path, "unix://")
	listener, err := net.ListenPacket("unixgram", path)
	if err != nil {
		return nil, err
	}
	return listener, nil
}
