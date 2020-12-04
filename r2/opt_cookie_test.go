package r2

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptCookie(t *testing.T) {
	assert := assert.New(t)

	opt := OptCookie(&http.Cookie{Name: "foo", Value: "bar"})

	req := New("https://foo.bar.local")
	assert.Nil(opt(req))

	assert.NotNil(req.Request.Header)
	assert.NotEmpty(req.Request.Header.Get("Cookie"))
}

func TestOptCookieValue(t *testing.T) {
	assert := assert.New(t)

	opt := OptCookieValue("foo", "bar")

	req := New("https://foo.bar.local")
	assert.Nil(opt(req))

	assert.NotNil(req.Request.Header)
	assert.NotEmpty(req.Request.Header.Get("Cookie"))
}
