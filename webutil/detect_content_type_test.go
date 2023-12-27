/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestDetectContentType(t *testing.T) {
	assert := assert.New(t)

	contentType, err := DetectContentType("foo.jpg")
	assert.Equal("image/jpeg", contentType)
	assert.Nil(err)

	contentType, err = DetectContentType("testdata/blank.pdf")
	assert.Equal("application/pdf", contentType)
	assert.Nil(err)

	contentType, err = DetectContentType("invalid_path.pdf")
	assert.Equal("", contentType)
	assert.NotNil(err)
}
