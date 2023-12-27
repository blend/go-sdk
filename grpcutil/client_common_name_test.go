/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
