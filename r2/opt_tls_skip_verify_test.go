/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptTLSSkipVerify(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptTLSSkipVerify(true))
	assert.NotNil(r.Client.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify)
}
