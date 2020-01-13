package r2

import (
	"net/http"
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestOptHeader(t *testing.T) {
	assert := assert.New(t)

	opt := OptHeader(http.Header{"Foo": []string{"bar"}})

	req := New("https://foo.bar.local")
	assert.Nil(opt(req))

	assert.NotNil(req.Header)
	assert.Equal("bar", req.Header.Get("foo"))
}

func TestOptHeaderValue(t *testing.T) {
	assert := assert.New(t)

	opt := OptHeaderValue("Foo", "bar")

	req := New("https://foo.bar.local")
	assert.Nil(opt(req))

	assert.NotNil(req.Header)
	assert.Equal("bar", req.Header.Get("foo"))
}
