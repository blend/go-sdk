/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package proxyprotocol

import (
	"context"
	"net"
)

// NewDialer returns a new proxy protocol dialer.
func NewDialer(opts ...DialerOption) *Dialer {
	d := &Dialer{
		Dialer:		new(net.Dialer),
		HeaderProvider:	func(_ context.Context, _ net.Conn) *Header { return nil },
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// OptDialerHeaderProvider sets the header provider.
func OptDialerHeaderProvider(provider func(context.Context, net.Conn) *Header) DialerOption {
	return func(d *Dialer) {
		d.HeaderProvider = provider
	}
}

// OptDialerConstSourceAdddr sets the header provider to be a constant source.
func OptDialerConstSourceAdddr(addr net.Addr) DialerOption {
	return func(d *Dialer) {
		d.HeaderProvider = func(_ context.Context, conn net.Conn) *Header {
			return &Header{
				Version:		1,
				Command:		ProtocolVersionAndCommandProxy,
				TransportProtocol:	AddressFamilyAndProtocolTCPv4,
				SourceAddr:		addr,
				DestinationAddr:	conn.RemoteAddr(),
			}
		}
	}
}

// DialerOption mutates a dialer.
type DialerOption func(*Dialer)

// Dialer wraps a dialer with proxy protocol header injection.
type Dialer struct {
	*net.Dialer
	HeaderProvider	func(context.Context, net.Conn) *Header
}

// Dial implements the dialer, calling `HeaderProvider` for a the context passed to it.
func (d *Dialer) Dial(network, addr string) (net.Conn, error) {
	return d.DialContext(context.Background(), network, addr)
}

// DialContext implements the dialer, calling `HeaderProvider` for a the context passed to it.
func (d *Dialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	conn, err := d.Dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	header := d.HeaderProvider(ctx, conn)
	if header == nil {
		return conn, nil
	}
	_, err = header.WriteTo(conn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
