/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"crypto/x509"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptTLSRootCAs(t *testing.T) {
	assert := assert.New(t)

	pool, err := x509.SystemCertPool()
	assert.Nil(err)
	r := New(TestURL, OptTLSRootCAs(pool))
	assert.NotNil(r.Client.Transport.(*http.Transport).TLSClientConfig.RootCAs)
}

func TestOptTLSRootCAsWithNilTransport(t *testing.T) {
	assert := assert.New(t)

	var transport *http.Transport
	certPool, err := x509.SystemCertPool()
	assert.Nil(err)

	req := New(
		TestURL,
		// NOTE: Transport **must** come before the root CAs since the CAs get set
		//       **on** the transport.
		OptTransport(transport),
		OptTLSRootCAs(certPool),
	)

	assert.NotNil(req.Client)
	assert.NotNil(req.Client.Transport)
	typed, ok := req.Client.Transport.(*http.Transport)
	assert.True(ok)
	assert.NotNil(typed)
	assert.NotNil(typed.TLSClientConfig)
	assert.NotNil(typed.TLSClientConfig.RootCAs)
}
