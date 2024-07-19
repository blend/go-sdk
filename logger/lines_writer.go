/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"bytes"
	"io"

	"github.com/blend/go-sdk/ex"
)

// NOTE: Ensure that
//   - `LinesWriter` satisfies `io.Writer`.
var (
	_ io.Writer = (*LinesWriter)(nil)
)

// LinesWriter is a writer that writes one line at a time, i.e. if `Write()` is called
// with multiple lines, `w.Write()` will be called for each line.
type LinesWriter struct {
	w io.Writer
}

// NewLinesWriter returns a new line writer.
func NewLinesWriter(w io.Writer) *LinesWriter {
	return &LinesWriter{
		w: w,
	}
}

// Write implements io.Writer.
func (lw *LinesWriter) Write(p []byte) (int, error) {
	n := 0
	for {
		eol := bytes.Index(p, []byte("\n"))
		if eol == -1 {
			break
		}
		lineN, err := lw.w.Write(p[:eol])
		n += lineN + 1
		if err != nil {
			return n, ex.New(err)
		}
		p = p[eol+1:]
	}
	lineN, err := lw.w.Write(p)
	return n + lineN, ex.New(err)
}
