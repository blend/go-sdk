package grpcutil

import (
	"testing"

	"google.golang.org/grpc/metadata"

	"github.com/blend/go-sdk/assert"
)

func TestMetaValue(t *testing.T) {
	assert := assert.New(t)
	md := metadata.New(map[string]string{"testkey": "val"})
	assert.Equal("", MetaValue(md, "missingkey"))
	assert.Equal("val", MetaValue(md, "testkey"))
}
