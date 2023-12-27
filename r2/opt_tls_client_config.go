/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

import (
	"crypto/tls"
)

// OptTLSClientConfig sets the tls config for the request.
// It will create a client, and a transport if unset.
func OptTLSClientConfig(cfg *tls.Config) Option {
	return func(r *Request) error {
		transport, err := EnsureHTTPTransport(r)
		if err != nil {
			return err
		}
		transport.TLSClientConfig = cfg
		return nil
	}
}
