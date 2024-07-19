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

func TestOptNoFollow(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptNoFollow())

	assert.NotNil(r.Client)
	assert.NotNil(r.Client.CheckRedirect)
	assert.Equal(http.ErrUseLastResponse, r.Client.CheckRedirect(nil, nil))
}
