package r2

import (
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestOptPath(t *testing.T) {
	assert := assert.New(t)

	r := New("http://foo.com", OptPath("/not-foo"))
	assert.Equal("/not-foo", r.Request.URL.Path)

	var unset Request
	OptPath("/not-foo")(&unset)
	assert.Equal("/not-foo", unset.URL.Path)
}

func TestOptPathf(t *testing.T) {
	assert := assert.New(t)

	r := New("http://foo.com", OptPathf("/not-foo/%s", "bar"))
	assert.Equal("/not-foo/bar", r.Request.URL.Path)

	var unset Request
	OptPathf("/not-foo/%s", "bar")(&unset)
	assert.Equal("/not-foo/bar", unset.URL.Path)
}
