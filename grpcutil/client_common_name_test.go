/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package grpcutil

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestClientCommonName(t *testing.T) {
	assert := assert.New(t)

	ctx := context.TODO()
	name := GetClientCommonName(ctx)
	assert.Equal("", name)

	ctx = WithClientCommonName(ctx, "name")
	name = GetClientCommonName(ctx)
	assert.Equal("name", name)
}
