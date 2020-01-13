package r2

import (
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestOptURL(t *testing.T) {
	assert := assert.New(t)

	r := New("http://foo.com", OptURL("https://foo.bar.com/buzz?a=b"))
	assert.NotNil(r.URL)
	assert.Equal("https://foo.bar.com/buzz?a=b", r.URL.String())

	var unset Request
	OptURL("https://foo.bar.com/buzz?a=b")(&unset)
	assert.Equal("https://foo.bar.com/buzz?a=b", unset.URL.String())
}
