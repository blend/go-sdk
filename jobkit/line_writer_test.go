package jobkit

import (
	"io"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func justError(_ interface{}, err error) error {
	return err
}

func TestLineWriter(t *testing.T) {
	assert := assert.New(t)

	lw := new(LineWriter)

	assert.Nil(justError(io.WriteString(lw, "this is a test\n")))
	assert.Nil(justError(io.WriteString(lw, "this is another test\n")))
	assert.Nil(justError(io.WriteString(lw, "this is a test\n")))
	assert.Nil(justError(io.WriteString(lw, "this is another test\n")))
	assert.Len(lw.Lines, 4)
}

func TestLineWriterWritten(t *testing.T) {
	assert := assert.New(t)

	lw := new(LineWriter)

	written, err := lw.Write([]byte("this is just a test"))
	assert.Nil(err)
	assert.Equal(19, written)
	assert.Len(lw.Buffer, 19)
	assert.Empty(lw.Lines)

	written, err = lw.Write([]byte("this is just a test\n"))
	assert.Nil(err)
	assert.Equal(20, written)
	assert.Empty(lw.Buffer)
	assert.Len(lw.Lines, 1)
	assert.Equal("this is just a testthis is just a test", string(lw.Lines[0].Line))

	written, err = lw.Write([]byte("this is another test"))
	assert.Nil(err)
	assert.Equal(20, written)
	assert.Len(lw.Buffer, 20)

	written, err = lw.Write([]byte("this is another test\n"))
	assert.Nil(err)
	assert.Equal(21, written)
	assert.Empty(lw.Buffer)
	assert.Len(lw.Lines, 2)
	assert.Equal("this is just a testthis is just a test", string(lw.Lines[0].Line))
	assert.Equal("this is another testthis is another test", string(lw.Lines[1].Line))
}
