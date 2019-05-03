package r2

import (
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptOnRequest(t *testing.T) {
	assert := assert.New(t)

	r := New("http://foo.com",
		OptOnRequest(func(_ *http.Request) error { return nil }),
		OptOnRequest(func(_ *http.Request) error { return nil }),
	)
	assert.Len(r.OnRequest, 2)
}
