package grpcutil

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestClientCommonName(t *testing.T) {
	assert := assert.New(t)

	ctx := context.TODO()
	name := GetClientCommonName(ctx)
	assert.Equal("", name)

	ctx = WithClientCommonName(ctx, "name")
	name = GetClientCommonName(ctx)
	assert.Equal("name", name)
}
