package r2

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptPath(t *testing.T) {
	assert := assert.New(t)

	r := New("http://foo.com", OptPath("/not-foo"))
	assert.Equal("/not-foo", r.Request.URL.Path)
}
