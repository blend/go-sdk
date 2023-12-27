/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"crypto/tls"
)

// OptTLSClientCertObject adds a client cert to the request
func OptTLSClientCertObject(cert *tls.Certificate) Option {
	return func(r *Request) error {
		transport, err := EnsureHTTPTransport(r)
		if err != nil {
			return err
		}

		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}
		transport.TLSClientConfig.Certificates = append(transport.TLSClientConfig.Certificates, *cert)
		return nil
	}
}
