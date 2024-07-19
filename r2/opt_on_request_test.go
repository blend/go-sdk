/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptOnRequest(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL,
		OptOnRequest(func(_ *http.Request) error { return nil }),
		OptOnRequest(func(_ *http.Request) error { return nil }),
	)
	assert.Len(r.OnRequest, 2)
}
