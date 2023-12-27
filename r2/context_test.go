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

func Test_WithParameterizedPath(t *testing.T) {
	its := assert.New(t)

	ctx := WithParameterizedPath(context.Background(), "/foo/:id")
	its.Equal("/foo/:id", GetParameterizedPath(ctx))
}
