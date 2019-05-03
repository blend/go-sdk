package r2

import (
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestOptOnResponse(t *testing.T) {
	assert := assert.New(t)

	r := New("http://foo.com",
		OptOnResponse(func(_ *http.Request, _ *http.Response, _ time.Time, _ error) error { return nil }),
		OptOnResponse(func(_ *http.Request, _ *http.Response, _ time.Time, _ error) error { return nil }),
	)
	assert.Len(r.OnResponse, 2)
}
