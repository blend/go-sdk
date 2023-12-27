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

func TestOptHeader(t *testing.T) {
	assert := assert.New(t)

	opt := OptHeader(http.Header{"Foo": []string{"bar"}})

	req := New("https://foo.bar.local")
	assert.Nil(opt(req))

	assert.NotNil(req.Request.Header)
	assert.Equal("bar", req.Request.Header.Get("foo"))
}

func TestOptHeaderValue(t *testing.T) {
	assert := assert.New(t)

	opt := OptHeaderValue("Foo", "bar")

	req := New("https://foo.bar.local")
	assert.Nil(opt(req))

	assert.NotNil(req.Request.Header)
	assert.Equal("bar", req.Request.Header.Get("foo"))
}
