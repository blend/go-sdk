/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptURL(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptURL("https://foo.bar.com/buzz?a=b"))
	assert.NotNil(r.Request.URL)
	assert.Equal("https://foo.bar.com/buzz?a=b", r.Request.URL.String())

	var unset Request
	assert.NotNil(OptURL("https://foo.bar.com/buzz?a=b")(&unset))
}
