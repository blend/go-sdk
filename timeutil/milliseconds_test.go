/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package timeutil

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestMilliseconds(t *testing.T) {
	assert := assert.New(t)
	d := time.Millisecond + time.Microsecond

	assert.Equal(1.001, Milliseconds(d))
}

func TestFromMilliseconds(t *testing.T) {
	assert := assert.New(t)
	expected := time.Millisecond + time.Microsecond
	assert.Equal(expected, FromMilliseconds(1.001))
}
