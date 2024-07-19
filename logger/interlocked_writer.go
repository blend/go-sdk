/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"io"
	"sync"
)

// NewInterlockedWriter returns a new interlocked writer.
func NewInterlockedWriter(output io.Writer) *InterlockedWriter {
	if typed, ok := output.(*InterlockedWriter); ok {
		return typed
	}
	return &InterlockedWriter{
		Output: output,
	}
}

// InterlockedWriter is a writer that serializes access to the Write() method.
type InterlockedWriter struct {
	sync.Mutex

	Output io.Writer
}

// Write writes the given bytes to the inner writer.
func (iw *InterlockedWriter) Write(buffer []byte) (count int, err error) {
	iw.Lock()

	count, err = iw.Output.Write(buffer)
	if err != nil {
		iw.Unlock()
		return
	}

	iw.Unlock()
	return
}

// Close closes any outputs that are io.WriteCloser's.
func (iw *InterlockedWriter) Close() (err error) {
	iw.Lock()
	defer iw.Unlock()

	if typed, isTyped := iw.Output.(io.WriteCloser); isTyped {
		err = typed.Close()
		if err != nil {
			return
		}
	}
	return
}
