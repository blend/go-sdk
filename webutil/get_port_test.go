/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestGetPort(t *testing.T) {
	assert := assert.New(t)

	assert.Empty(GetPort(nil))
	assert.Empty(GetPort(&http.Request{}))
	assert.Equal("8443", GetPort(&http.Request{
		Header: http.Header{
			HeaderXForwardedPort: {"8443"},
		},
	}), "should use existing header if found")
	assert.Equal("8443", GetPort(&http.Request{
		Header: http.Header{
			HeaderXForwardedPort: {"9090,8443"},
		},
	}), "should use existing header last chunk if found")
}
