package r2

import (
	"net/http"
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestOptDisableKeepAlives(t *testing.T) {
	assert := assert.New(t)

	r := New("https://foo.bar.local", OptDisableKeepAlives(true))
	assert.True(r.Client.Transport.(*http.Transport).DisableKeepAlives)
}
