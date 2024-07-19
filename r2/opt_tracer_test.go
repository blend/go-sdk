/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptTracer(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptTracer(MockTracer{}))
	assert.NotNil(r.Tracer)
}
