package webutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestETag(t *testing.T) {
	assert := assert.New(t)

	etag := ETag([]byte("a quick brown fox jumps over the something cool"))
	assert.Equal("4743a94a6030d34968f838c94cf4a6fd", etag)

	etag = ETag([]byte("something else that is really cool"))
	assert.Equal("a8c90c3202be46c1d766b2c63d38332b", etag)
}
