/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import "crypto/tls"

// OptTLSClientCertFilePair adds a client cert and key to the request.
func OptTLSClientCertFilePair(certFile, keyFile string) Option {
	return func(r *Request) error {
		transport, err := EnsureHTTPTransport(r)
		if err != nil {
			return err
		}

		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return err
		}
		transport.TLSClientConfig.Certificates = append(transport.TLSClientConfig.Certificates, cert)
		return nil
	}
}
