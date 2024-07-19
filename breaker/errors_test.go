/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package breaker

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestErrIsOpen(t *testing.T) {
	assert := assert.New(t)

	assert.True(ErrIsOpen(ex.New(ErrOpenState)))
	assert.False(ErrIsOpen(nil))
	assert.False(ErrIsOpen(ex.New(ErrTooManyRequests)))
}

func TestErrIsTooManyRequests(t *testing.T) {
	assert := assert.New(t)

	assert.True(ErrIsTooManyRequests(ex.New(ErrTooManyRequests)))
	assert.False(ErrIsTooManyRequests(nil))
	assert.False(ErrIsTooManyRequests(ex.New(ErrOpenState)))
}
