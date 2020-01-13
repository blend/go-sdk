package r2

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptContext(t *testing.T) {
	assert := assert.New(t)

	opt := OptContext(context.TODO())

	req := New("https://foo.bar.local")
	opt(req)
	assert.NotNil(req.Context())
}
