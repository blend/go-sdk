/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestOptTimeout(t *testing.T) {
	assert := assert.New(t)

	r := New("https://foo.bar.local", OptTimeout(time.Second))
	assert.Equal(time.Second, r.Client.Timeout)
}
