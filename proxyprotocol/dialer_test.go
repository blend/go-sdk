/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package proxyprotocol

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_Dialer(t *testing.T) {
	its := assert.New(t)

	listener, err := CreateListener("tcp4", "127.0.0.1:0",
		OptUseProxyProtocol(true),
	)
	its.Nil(err)
	defer listener.Close()

	sourceAddr := &net.TCPAddr{
		IP:	net.ParseIP("192.168.0.7"),
		Port:	31234,
	}
	dialer := NewDialer(
		OptDialerConstSourceAdddr(sourceAddr),
	)

	go func() {
		conn, err := dialer.DialContext(context.Background(), "tcp4", listener.Addr().String())
		if err != nil {
			panic(err)
		}
		defer conn.Close()
	}()

	conn, err := listener.Accept()
	its.Nil(err)
	its.Equal("192.168.0.7:31234", conn.RemoteAddr().String(), fmt.Sprintf("listener addr: %v", listener.Addr()))
}
