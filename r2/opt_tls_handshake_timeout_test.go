/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestOptTLSHandshakeTimeout(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptTLSHandshakeTimeout(time.Second))
	assert.Equal(time.Second, r.Client.Transport.(*http.Transport).TLSHandshakeTimeout)
}

func TestOptTLSHandshakeTimeoutWithNilTransport(t *testing.T) {
	assert := assert.New(t)

	var transport *http.Transport
	req := New(
		TestURL,
		// NOTE: Transport **must** come before the root CAs since the CAs get set
		//       **on** the transport.
		OptTransport(transport),
		OptTLSHandshakeTimeout(time.Second),
	)

	assert.NotNil(req.Client)
	assert.NotNil(req.Client.Transport)
	typed, ok := req.Client.Transport.(*http.Transport)
	assert.True(ok)
	assert.NotNil(typed)
	assert.Equal(time.Second, typed.TLSHandshakeTimeout)
}
