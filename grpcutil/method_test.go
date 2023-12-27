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

func TestMethod(t *testing.T) {
	assert := assert.New(t)

	ctx := context.TODO()
	m := GetMethod(ctx)
	assert.Equal("", m)

	ctx = WithMethod(ctx, "DoSomething")
	m = GetMethod(ctx)
	assert.Equal("DoSomething", m)
}
