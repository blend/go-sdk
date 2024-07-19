/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package proxyprotocol

import (
	"crypto/tls"
	"time"
)

// CreateListenerOptions are the options for creating listeners.
type CreateListenerOptions struct {
	TLSConfig        *tls.Config
	UseProxyProtocol bool
	KeepAlive        bool
	KeepAlivePeriod  time.Duration
}

// CreateListenerOption is a mutator for the options used when creating a listener.
type CreateListenerOption func(*CreateListenerOptions) error

// OptTLSConfig sets the listener tls config.
func OptTLSConfig(tlsConfig *tls.Config) CreateListenerOption {
	return func(clo *CreateListenerOptions) error {
		clo.TLSConfig = tlsConfig
		return nil
	}
}

// OptUseProxyProtocol sets if we should decode proxy protocol or not.
func OptUseProxyProtocol(useProxyProtocol bool) CreateListenerOption {
	return func(clo *CreateListenerOptions) error {
		clo.UseProxyProtocol = useProxyProtocol
		return nil
	}
}

// OptKeepAlive sets if we should keep TCP connections alive or not.
func OptKeepAlive(keepAlive bool) CreateListenerOption {
	return func(clo *CreateListenerOptions) error {
		clo.KeepAlive = keepAlive
		return nil
	}
}

// OptKeepAlivePeriod sets the duration we should keep connections alive for.
func OptKeepAlivePeriod(keepAlivePeriod time.Duration) CreateListenerOption {
	return func(clo *CreateListenerOptions) error {
		clo.KeepAlivePeriod = keepAlivePeriod
		return nil
	}
}
