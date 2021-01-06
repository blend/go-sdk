package web

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/webutil"
)

func TestPathRedirectHandler(t *testing.T) {
	assert := assert.New(t)

	redirect := PathRedirectHandler("/foo")

	newURL := redirect(NewCtx(nil, webutil.NewMockRequest("GET", "/notfoo")))
	assert.Equal("/foo", newURL.Path)
}

func TestBase64URL(t *testing.T) {
	assert := assert.New(t)
	bs := []byte("hello")
	enc := Base64URLEncode(bs)
	assert.NotEmpty(enc)

	out, err := Base64URLDecode(enc)
	assert.Nil(err)
	assert.Equal(string(bs), string(out))
}

func TestNewCookie(t *testing.T) {
	assert := assert.New(t)
	c := NewCookie("hello", "world")
	assert.NotNil(c)
	assert.Equal("hello", c.Name)
	assert.Equal("world", c.Value)
}

func TestMergeHeaders(t *testing.T) {
	assert := assert.New(t)

	a := map[string][]string{
		"Foo": {"foo1a", "foo2a"},
		"Bar": {"bar1a", "bar2a"},
	}

	b := map[string][]string{
		"Foo":            {"foo1b", "foo2b", "foo3b"},
		"example-string": {"dog"},
	}

	c := map[string][]string{
		"Bar":  {"bar1c", "bar2c"},
		"Buzz": {"fuzz"},
	}

	merged := MergeHeaders(a, b, c)

	assert.Equal(
		[]string{"foo1a", "foo2a", "foo1b", "foo2b", "foo3b"},
		merged["Foo"],
	)

	assert.Equal(
		[]string{"bar1a", "bar2a", "bar1c", "bar2c"},
		merged["Bar"],
	)

	assert.Equal(
		[]string{"dog"},
		merged[http.CanonicalHeaderKey("example-string")],
	)

	assert.Equal(
		[]string{"fuzz"},
		merged["Buzz"],
	)
}
