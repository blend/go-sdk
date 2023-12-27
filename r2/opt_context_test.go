/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptContext(t *testing.T) {
	assert := assert.New(t)

	opt := OptContext(context.TODO())

	req := New("https://foo.bar.local")
	assert.Nil(opt(req))
	assert.NotNil(req.Request.Context())
}
