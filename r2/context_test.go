/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
