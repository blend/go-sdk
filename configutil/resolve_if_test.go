/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestResolveIf(t *testing.T) {
	assert := assert.New(t)

	var called bool
	assert.Nil(ResolveIf(false, func(_ context.Context) error {
		called = true
		return nil
	})(context.Background()))

	assert.False(called)

	assert.Nil(ResolveIf(true, func(_ context.Context) error {
		called = true
		return nil
	})(context.Background()))
	assert.True(called)
}

func TestResolveIfFunc(t *testing.T) {
	assert := assert.New(t)

	returnBool := func(v bool) func(context.Context) bool {
		return func(_ context.Context) bool {
			return v
		}
	}

	var called bool
	assert.Nil(ResolveIfFunc(returnBool(false), func(_ context.Context) error {
		called = true
		return nil
	})(context.Background()))

	assert.False(called)

	assert.Nil(ResolveIfFunc(returnBool(true), func(_ context.Context) error {
		called = true
		return nil
	})(context.Background()))
	assert.True(called)
}
