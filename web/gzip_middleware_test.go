package web

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"testing"

	"github.com/blend/go-sdk/r2"

	"github.com/blend/go-sdk/assert"
)

func TestGZipMiddlewarePlaintext(t *testing.T) {
	assert := assert.New(t)

	app := New()
	app.Use(GZip)

	app.GET("/", ok)

	req := MockGet(app, "/")
	assert.Nil(req.Err)
	resBody, err := req.Bytes()
	assert.Nil(err)
	assert.Equal("\"OK!\"\n", string(resBody))
}

func TestGZipMiddlewareCompressed(t *testing.T) {
	assert := assert.New(t)

	app := New()
	app.Use(GZip)

	app.GET("/", ok)

	req := MockGet(app, "/", r2.OptHeaderValue(HeaderAcceptEncoding, "gzip"))
	assert.Nil(req.Err)
	body, meta, err := req.BytesWithResponse()

	assert.Equal("gzip", meta.Header.Get(HeaderContentEncoding))
	assert.Equal("Accept-Encoding", meta.Header.Get(HeaderVary))
	assert.Nil(err)

	decompressor, err := gzip.NewReader(bytes.NewBuffer(body))
	assert.Nil(err)
	decompressed, err := ioutil.ReadAll(decompressor)
	assert.Nil(err)

	assert.Equal("\"OK!\"\n", string(decompressed))
}
