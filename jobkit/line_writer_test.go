package jobkit

import (
	"encoding/json"
	"io"
	"testing"
	"time"

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
	assert.Len(lw.Current.Line, 19)
	assert.Empty(lw.Lines)

	written, err = lw.Write([]byte("this is just a test\n"))
	assert.Nil(err)
	assert.Equal(20, written)
	assert.True(lw.Current.Timestamp.IsZero())
	assert.Empty(lw.Current.Line)
	assert.Len(lw.Lines, 1)
	assert.Equal("this is just a testthis is just a test", string(lw.Lines[0].Line))

	written, err = lw.Write([]byte("this is another test"))
	assert.Nil(err)
	assert.Equal(20, written)
	assert.Len(lw.Current.Line, 20)

	written, err = lw.Write([]byte("this is another test\n"))
	assert.Nil(err)
	assert.Equal(21, written)
	assert.Empty(lw.Current.Line)
	assert.Len(lw.Lines, 2)
	assert.Equal("this is just a testthis is just a test", string(lw.Lines[0].Line))
	assert.Equal("this is another testthis is another test", string(lw.Lines[1].Line))
}

func TestLineWriterLinesJSON(t *testing.T) {
	assert := assert.New(t)

	lines := []Line{
		{Timestamp: time.Date(2019, 9, 4, 12, 11, 10, 9, time.UTC), Line: []byte("this is just a test")},
		{Timestamp: time.Date(2019, 9, 5, 12, 11, 10, 9, time.UTC), Line: []byte("this is just another test")},
	}

	contents, err := json.Marshal(lines)
	assert.Nil(err)

	var verify []Line
	err = json.Unmarshal(contents, &verify)
	assert.Nil(err)
	assert.Len(verify, 2)
}
