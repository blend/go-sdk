package cron

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

func TestOutputBuffer(t *testing.T) {
	assert := assert.New(t)

	lw := new(OutputBuffer)

	assert.Nil(justError(io.WriteString(lw, "this is a test\n")))
	assert.Nil(justError(io.WriteString(lw, "this is another test\n")))
	assert.Nil(justError(io.WriteString(lw, "this is a test\n")))
	assert.Nil(justError(io.WriteString(lw, "this is another test\n")))
	assert.Nil(justError(io.WriteString(lw, "this is another test\r\nand another\n")))
	assert.Len(lw.Chunks, 5)
}

func TestOutputBufferWritten(t *testing.T) {
	assert := assert.New(t)

	ob := new(OutputBuffer)

	written, err := ob.Write([]byte("this is just a test"))
	assert.Nil(err)
	assert.Equal(19, written)
	assert.Len(ob.Chunks, 1)

	written, err = ob.Write([]byte("this is just a test\n"))
	assert.Nil(err)
	assert.Equal(20, written)
	assert.Len(ob.Chunks, 2)
	assert.Equal("this is just a testthis is just a test\n", ob.String())
}

func TestOutputChunkJSON(t *testing.T) {
	assert := assert.New(t)

	chunk := OutputChunk{
		Timestamp: time.Date(2019, 9, 21, 12, 11, 10, 9, time.UTC),
		Data:      []byte("this is just a test"),
	}

	jsonContents, err := json.Marshal(chunk)
	assert.Nil(err)
	assert.NotEmpty(jsonContents)

	var verify OutputChunk
	assert.Nil(json.Unmarshal(jsonContents, &verify))
	assert.Equal(chunk.Timestamp, verify.Timestamp)
	assert.Equal(chunk.Data, verify.Data)
}
