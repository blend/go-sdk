package webutil

import (
	"bytes"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMockResponseWriterContentLength(t *testing.T) {
	assert := assert.New(t)

	mrw := MockResponseWriter{contentLength: 1337}
	assert.Equal(1337, mrw.ContentLength())
}

func TestMockResponseWriterBytes(t *testing.T) {
	assert := assert.New(t)

	b := []byte("not too much here")
	mrw := MockResponseWriter{contents: bytes.NewBuffer(b)}
	assert.Equal(b, mrw.Bytes())
}

func TestMockResponseWriterClose(t *testing.T) {
	assert := assert.New(t)

	mrw := MockResponseWriter{}
	assert.Nil(mrw.Close())
}
