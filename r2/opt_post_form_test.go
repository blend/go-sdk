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
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptPostForm(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptPostForm(url.Values{"bar": []string{"baz, buzz"}}))
	assert.NotNil(r.Request.PostForm)
	assert.NotEmpty(r.Request.PostForm.Get("bar"))
}

func TestOptPostFormValue(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptPostFormValue("bar", "baz"))
	assert.NotNil(r.Request.PostForm)
	assert.Equal("baz", r.Request.PostForm.Get("bar"))
}
