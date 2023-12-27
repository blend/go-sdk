/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptCookie(t *testing.T) {
	assert := assert.New(t)

	opt := OptCookie(&http.Cookie{Name: "foo", Value: "bar"})

	req := New("https://foo.bar.local")
	assert.Nil(opt(req))

	assert.NotNil(req.Request.Header)
	assert.NotEmpty(req.Request.Header.Get("Cookie"))
}

func TestOptCookieValue(t *testing.T) {
	assert := assert.New(t)

	opt := OptCookieValue("foo", "bar")

	req := New("https://foo.bar.local")
	assert.Nil(opt(req))

	assert.NotNil(req.Request.Header)
	assert.NotEmpty(req.Request.Header.Get("Cookie"))
}
