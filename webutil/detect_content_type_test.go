/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
