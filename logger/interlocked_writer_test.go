/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"bytes"
	"io"
	"testing"

	"github.com/blend/go-sdk/assert"
)

var (
	_ io.WriteCloser = (*mockWriter)(nil)
)

// mockWriter is a stub for a io.WriteCloser.
type mockWriter struct {
	WriteHandler func([]byte) (int, error)
	CloseHandler func() error
}

// Write implements io.Writer
func (mw mockWriter) Write(data []byte) (int, error) {
	if mw.WriteHandler != nil {
		return mw.WriteHandler(data)
	}
	return 0, nil
}

// Close implements io.Closer.
func (mw mockWriter) Close() error {
	if mw.CloseHandler != nil {
		return mw.CloseHandler()
	}
	return nil
}

func TestInterlockedWriter(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	var didWrite, didClose bool
	mw := mockWriter{
		WriteHandler: func(data []byte) (int, error) {
			defer func() {
				didWrite = true
			}()
			return buf.Write(data)
		},
		CloseHandler: func() error {
			defer func() {
				didClose = true
			}()
			return nil
		},
	}

	iw := NewInterlockedWriter(mw)
	data := []byte("this is a test")
	written, err := iw.Write(data)
	assert.Nil(err)
	assert.Equal(len(data), written)
	assert.True(didWrite)
	assert.Nil(iw.Close())
	assert.True(didClose)

	assert.Equal("this is a test", buf.String())
}
