package r2

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptTracer(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptTracer(MockTracer{}))
	assert.NotNil(r.Tracer)
}
