package bufferutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestPool(t *testing.T) {
	assert := assert.New(t)

	pool := NewPool(1024)
	buf := pool.Get()
	assert.NotNil(buf)
	assert.Equal(1024, buf.Cap())
	assert.Zero(buf.Len())
	pool.Put(buf)
}
