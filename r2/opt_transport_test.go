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

func TestOptTransport(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptTransport(&http.Transport{}))
	assert.NotNil(r.Client.Timeout)
}
