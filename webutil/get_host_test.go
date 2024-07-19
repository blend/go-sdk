/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestGetHost(t *testing.T) {
	assert := assert.New(t)

	r := http.Request{
		Host: "local.test.com",
	}
	assert.Equal("local.test.com", GetHost(&r))

	r = http.Request{
		Host: "local.test.com:8080",
	}
	assert.Equal("local.test.com", GetHost(&r))

	r = http.Request{
		URL:  &url.URL{},
		Host: "local.test.com:8080",
	}
	assert.Equal("local.test.com", GetHost(&r))

	r = http.Request{
		URL:  &url.URL{Host: "local.foo.com"},
		Host: "local.test.com:8080",
	}
	assert.Equal("local.foo.com", GetHost(&r))

	headers := http.Header{}
	headers.Set("X-Forwarded-Proto", "spdy,https")
	headers.Set("X-Forwarded-Host", "local.bar.com")
	r = http.Request{
		URL:    &url.URL{Host: "local.foo.com"},
		Host:   "local.test.com:8080",
		Header: headers,
	}
	assert.Equal("local.bar.com", GetHost(&r))
}
