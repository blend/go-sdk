/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptTLSCLientConfig(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptTLSClientConfig(&tls.Config{}))
	assert.NotNil(r.Client.Transport.(*http.Transport).TLSClientConfig)
}
