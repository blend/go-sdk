package stringutil

import (
	"bytes"
	"fmt"

	"github.com/blend/go-sdk/ex"
)

const newline = "\n"

// NewBuffer creates a new buffer.
func NewBuffer(input []byte) *Buffer {
	return &Buffer{
		Buffer: bytes.NewBuffer(input),
	}
}

// Buffer is an extension of bytes.Buffer with some helpers.
type Buffer struct {
	*bytes.Buffer
}

// Writeline is a macro for writing a string to the buffer that ends in \n.
func (b *Buffer) Writeline(contents ...interface{}) (int, error) {
	n, err := b.WriteString(fmt.Sprint(contents...) + newline)
	return n, ex.New(err)
}

// Writelinef is a macro for writing a string to the buffer with a given format and args that ends in \n.
func (b *Buffer) Writelinef(format string, args ...interface{}) (int, error) {
	n, err := b.WriteString(fmt.Sprintf(format, args...) + newline)
	return n, ex.New(err)
}
