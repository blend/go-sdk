/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"bytes"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestBufferedCompressedWriter(t *testing.T) {
	assert := assert.New(t)

	buf := bytes.NewBuffer(nil)
	mockedWriter := NewMockResponse(buf)
	bufferedWriter := NewGZipResponseWriter(mockedWriter)

	written, err := bufferedWriter.Write([]byte("ok"))
	assert.Nil(err)
	assert.NotZero(written)
}
