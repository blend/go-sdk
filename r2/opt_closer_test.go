package r2

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptCloser(t *testing.T) {
	it := assert.New(t)

	server := mockServerOK()
	defer server.Close()

	var didCallCloser bool
	req := New(server.URL, OptCloser(func() error {
		didCallCloser = true
		return fmt.Errorf("closer test")
	}))
	it.NotNil(req.Closer)

	meta, err := req.Discard()
	it.True(didCallCloser)
	it.NotNil(err)
	it.Equal("closer test", err.Error())
	it.Equal(http.StatusOK, meta.StatusCode)
}
