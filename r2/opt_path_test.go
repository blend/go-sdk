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
