package bindata

import (
	"fmt"
	"io"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// NewByteWriter returns a new byte writer.
func NewByteWriter(wr io.Writer) *ByteWriter {
	return &ByteWriter{Writer: wr}
}

// ByteWriter writes escaped bytes to a writer.
type ByteWriter struct {
	io.Writer
	Indent []byte
	c      int
}

func (w *ByteWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}

	for n = range p {
		if w.c%96 == 0 {
			if n > 0 {
				_, _ = w.Writer.Write(newline)
			}
			_, _ = w.Writer.Write(w.Indent)
			w.c = 0
		} else {
			_, _ = w.Writer.Write(space)
		}

		fmt.Fprintf(w.Writer, "0x%02x,", p[n])
		w.c++
	}

	n++
	return
}
