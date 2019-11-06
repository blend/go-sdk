package webutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestETag(t *testing.T) {
	assert := assert.New(t)

	etag, err := ETag([]byte("a quick brown fox jumps over the something cool"))
	assert.Nil(err)
	assert.Equal("4743a94a6030d34968f838c94cf4a6fd", etag)
}
