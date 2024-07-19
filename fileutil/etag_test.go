/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package fileutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestETag(t *testing.T) {
	assert := assert.New(t)

	corpus := []byte(`the quick brown fox jumpst over the lazy dog`)
	etag, err := ETag(corpus)
	assert.Nil(err)
	assert.Equal("10cb95681f0bf2e2c3263a6ea222c463", etag)
}
