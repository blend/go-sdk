/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

type mockResponseWriter struct {
	Headers		http.Header
	StatusCode	int
	Output		io.Writer
}

// Header returns the response headers.
func (mrw mockResponseWriter) Header() http.Header {
	return mrw.Headers
}

// WriteHeader writes the status code.
func (mrw mockResponseWriter) WriteHeader(code int) {
	mrw.StatusCode = code
}

// Write writes data.
func (mrw mockResponseWriter) Write(contents []byte) (int, error) {
	return mrw.Output.Write(contents)
}

func Test_StatusResponseWriter(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	rw := NewStatusResponseWriter(mockResponseWriter{Output: output, Headers: http.Header{}})

	rw.Header().Set("foo", "bar")
	rw.WriteHeader(http.StatusOK)
	_, err := rw.Write([]byte("this is a test"))
	assert.Nil(err)

	assert.Equal(http.StatusOK, rw.StatusCode())
	assert.Equal("this is a test", output.String())
}

func Test_StatusResponseWriter_self(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	rw := NewStatusResponseWriter(mockResponseWriter{Output: output, Headers: http.Header{}})

	srw := NewStatusResponseWriter(rw)

	srw.Header().Set("foo", "bar")
	srw.WriteHeader(http.StatusOK)
	_, err := srw.Write([]byte("this is a test"))
	assert.Nil(err)

	assert.Equal(http.StatusOK, rw.StatusCode())
	assert.Equal("this is a test", output.String())
}
