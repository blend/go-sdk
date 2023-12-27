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

func TestOptNoFollow(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptNoFollow())

	assert.NotNil(r.Client)
	assert.NotNil(r.Client.CheckRedirect)
	assert.Equal(http.ErrUseLastResponse, r.Client.CheckRedirect(nil, nil))
}
