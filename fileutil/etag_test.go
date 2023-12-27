/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
